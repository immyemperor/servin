//go:build linux

package vm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// KVMProvider implements VM operations using Linux KVM/QEMU
type KVMProvider struct {
	config  *VMConfig
	vmPath  string
	sshPort int
	running bool
}

// NewKVMProvider creates a new KVM provider
func NewKVMProvider(config *VMConfig) (VMProvider, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	vmPath := filepath.Join(homeDir, ".servin", "vms", config.Name)

	return &KVMProvider{
		config:  config,
		vmPath:  vmPath,
		sshPort: config.SSHPort,
		running: false,
	}, nil
}

// Create creates a new VM using KVM/QEMU
func (p *KVMProvider) Create(config *VMConfig) error {
	// Check if KVM is available
	if !p.isKVMAvailable() {
		return fmt.Errorf("KVM is not available on this system")
	}

	// Ensure VM directory exists
	if err := os.MkdirAll(p.vmPath, 0755); err != nil {
		return fmt.Errorf("failed to create VM directory: %v", err)
	}

	return p.createKVMVM(config)
}

// Start starts the VM
func (p *KVMProvider) Start() error {
	if p.running {
		return nil
	}

	return p.startKVMVM()
}

// Stop stops the VM
func (p *KVMProvider) Stop() error {
	if !p.running {
		return nil
	}

	// Send shutdown signal via SSH
	cmd := exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"root@localhost",
		"shutdown -h now")

	cmd.Run() // Ignore errors as VM might shutdown before SSH responds

	// Wait for VM to stop
	time.Sleep(5 * time.Second)
	p.running = false

	return nil
}

// Destroy removes the VM completely
func (p *KVMProvider) Destroy() error {
	if p.running {
		p.Stop()
	}

	return os.RemoveAll(p.vmPath)
}

// IsRunning checks if the VM is currently running
func (p *KVMProvider) IsRunning() bool {
	// Check if we can connect via SSH
	cmd := exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "ConnectTimeout=2",
		"root@localhost",
		"echo 'alive'")

	err := cmd.Run()
	p.running = (err == nil)
	return p.running
}

// GetInfo returns VM information
func (p *KVMProvider) GetInfo() (*VMInfo, error) {
	return &VMInfo{
		Name:       p.config.Name,
		Status:     p.getStatus(),
		Platform:   "Linux",
		Provider:   "KVM/QEMU",
		CPUs:       p.config.CPUs,
		Memory:     p.config.Memory,
		IPAddress:  "127.0.0.1",
		SSHPort:    p.sshPort,
		DockerPort: p.config.DockerPort,
		Capabilities: map[string]bool{
			"containers":   true,
			"networking":   true,
			"volumes":      true,
			"port_forward": true,
			"nested_virt":  true, // KVM supports nested virtualization
		},
	}, nil
}

// RunContainer runs a container inside the VM
func (p *KVMProvider) RunContainer(config *ContainerConfig) (*ContainerResult, error) {
	if !p.IsRunning() {
		return nil, fmt.Errorf("VM is not running")
	}

	// Build docker run command
	dockerCmd := p.buildDockerCommand(config)

	// Execute via SSH
	cmd := exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"root@localhost",
		dockerCmd)

	output, err := cmd.CombinedOutput()

	result := &ContainerResult{
		Name:   config.Name,
		Output: string(output),
	}

	if err != nil {
		result.Error = err.Error()
		result.ExitCode = 1
	} else {
		result.Status = "running"
		result.ExitCode = 0
		// Extract container ID from output
		if lines := strings.Split(strings.TrimSpace(string(output)), "\n"); len(lines) > 0 {
			result.ID = strings.TrimSpace(lines[len(lines)-1])
		}
	}

	return result, nil
}

// ListContainers lists containers in the VM
func (p *KVMProvider) ListContainers() ([]*ContainerInfo, error) {
	if !p.IsRunning() {
		return nil, fmt.Errorf("VM is not running")
	}

	cmd := exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"root@localhost",
		"docker ps -a --format 'table {{.ID}}\\t{{.Names}}\\t{{.Image}}\\t{{.Status}}\\t{{.CreatedAt}}\\t{{.Command}}'")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %v", err)
	}

	return p.parseContainerList(string(output)), nil
}

// StopContainer stops a container in the VM
func (p *KVMProvider) StopContainer(id string) error {
	return p.executeDockerCommand(fmt.Sprintf("docker stop %s", id))
}

// RemoveContainer removes a container in the VM
func (p *KVMProvider) RemoveContainer(id string) error {
	return p.executeDockerCommand(fmt.Sprintf("docker rm %s", id))
}

// CopyToVM copies files from host to VM
func (p *KVMProvider) CopyToVM(hostPath, vmPath string) error {
	cmd := exec.Command("scp",
		"-P", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		hostPath,
		fmt.Sprintf("root@localhost:%s", vmPath))

	return cmd.Run()
}

// CopyFromVM copies files from VM to host
func (p *KVMProvider) CopyFromVM(vmPath, hostPath string) error {
	cmd := exec.Command("scp",
		"-P", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		fmt.Sprintf("root@localhost:%s", vmPath),
		hostPath)

	return cmd.Run()
}

// ForwardPort forwards a port from host to VM
func (p *KVMProvider) ForwardPort(hostPort, vmPort int) error {
	// Port forwarding is typically configured during VM creation with QEMU
	// For dynamic forwarding, we'd need to use iptables or similar
	return fmt.Errorf("dynamic port forwarding not implemented")
}

// RemovePortForward removes a port forward
func (p *KVMProvider) RemovePortForward(hostPort int) error {
	return fmt.Errorf("dynamic port forwarding not implemented")
}

// Helper methods

func (p *KVMProvider) isKVMAvailable() bool {
	// Check if KVM device exists
	if _, err := os.Stat("/dev/kvm"); err != nil {
		return false
	}

	// Check if qemu-system-x86_64 is available
	cmd := exec.Command("which", "qemu-system-x86_64")
	return cmd.Run() == nil
}

func (p *KVMProvider) createKVMVM(config *VMConfig) error {
	diskPath := filepath.Join(p.vmPath, "disk.qcow2")

	// Create disk image
	cmd := exec.Command("qemu-img", "create", "-f", "qcow2", diskPath, fmt.Sprintf("%dG", config.DiskSize))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create disk image: %v", err)
	}

	// Note: In a real implementation, we would:
	// 1. Download appropriate Linux ISO (Alpine, Ubuntu, etc.)
	// 2. Create cloud-init configuration for automated setup
	// 3. Install Docker/containerd inside the VM
	// 4. Configure SSH keys and networking

	return nil
}

func (p *KVMProvider) startKVMVM() error {
	diskPath := filepath.Join(p.vmPath, "disk.qcow2")

	// QEMU command with KVM acceleration
	args := []string{
		"-enable-kvm",
		"-M", "pc",
		"-cpu", "host",
		"-smp", strconv.Itoa(p.config.CPUs),
		"-m", strconv.Itoa(p.config.Memory),
		"-drive", fmt.Sprintf("file=%s,if=virtio", diskPath),
		"-netdev", fmt.Sprintf("user,id=net0,hostfwd=tcp::%d-:22,hostfwd=tcp::%d-:2375", p.sshPort, p.config.DockerPort),
		"-device", "virtio-net-pci,netdev=net0",
		"-nographic",
		"-daemonize",
	}

	cmd := exec.Command("qemu-system-x86_64", args...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start QEMU VM: %v", err)
	}

	// Wait for VM to boot
	for i := 0; i < 30; i++ {
		if p.IsRunning() {
			p.running = true
			return nil
		}
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("VM failed to start within timeout")
}

func (p *KVMProvider) getStatus() string {
	if p.IsRunning() {
		return "running"
	}
	return "stopped"
}

func (p *KVMProvider) buildDockerCommand(config *ContainerConfig) string {
	cmd := []string{"docker", "run"}

	if config.Detached {
		cmd = append(cmd, "-d")
	}

	if config.Name != "" {
		cmd = append(cmd, "--name", config.Name)
	}

	// Add environment variables
	for key, value := range config.Environment {
		cmd = append(cmd, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	// Add port mappings
	for hostPort, containerPort := range config.Ports {
		cmd = append(cmd, "-p", fmt.Sprintf("%s:%s", hostPort, containerPort))
	}

	// Add volumes
	for hostPath, containerPath := range config.Volumes {
		cmd = append(cmd, "-v", fmt.Sprintf("%s:%s", hostPath, containerPath))
	}

	// Add working directory
	if config.WorkDir != "" {
		cmd = append(cmd, "-w", config.WorkDir)
	}

	// Add image
	cmd = append(cmd, config.Image)

	// Add command
	cmd = append(cmd, config.Command...)

	return strings.Join(cmd, " ")
}

func (p *KVMProvider) executeDockerCommand(dockerCmd string) error {
	cmd := exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"root@localhost",
		dockerCmd)

	return cmd.Run()
}

func (p *KVMProvider) parseContainerList(output string) []*ContainerInfo {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) <= 1 { // Skip header line
		return []*ContainerInfo{}
	}

	containers := make([]*ContainerInfo, 0, len(lines)-1)

	for _, line := range lines[1:] { // Skip header
		fields := strings.Split(line, "\t")
		if len(fields) >= 6 {
			containers = append(containers, &ContainerInfo{
				ID:      fields[0],
				Name:    fields[1],
				Image:   fields[2],
				Status:  fields[3],
				Created: fields[4],
				Command: fields[5],
			})
		}
	}

	return containers
}
