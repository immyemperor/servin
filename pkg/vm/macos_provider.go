//go:build darwin

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

// VirtualizationFrameworkProvider implements VM operations using macOS Virtualization.framework
type VirtualizationFrameworkProvider struct {
	config  *VMConfig
	vmPath  string
	sshPort int
	running bool
}

// NewVirtualizationFrameworkProvider creates a new Virtualization.framework provider
func NewVirtualizationFrameworkProvider(config *VMConfig) (VMProvider, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	vmPath := filepath.Join(homeDir, ".servin", "vms", config.Name)

	return &VirtualizationFrameworkProvider{
		config:  config,
		vmPath:  vmPath,
		sshPort: config.SSHPort,
		running: false,
	}, nil
}

// Create creates a new VM using UTM (which uses Virtualization.framework)
func (p *VirtualizationFrameworkProvider) Create(config *VMConfig) error {
	// Ensure VM directory exists
	if err := os.MkdirAll(p.vmPath, 0755); err != nil {
		return fmt.Errorf("failed to create VM directory: %v", err)
	}

	// Check if UTM is available, fall back to QEMU if not
	if p.isUTMAvailable() {
		return p.createUTMVM(config)
	}

	return p.createQEMUVM(config)
}

// Start starts the VM
func (p *VirtualizationFrameworkProvider) Start() error {
	if p.running {
		return nil
	}

	// Try UTM first, fall back to QEMU
	if p.isUTMAvailable() {
		return p.startUTMVM()
	}

	return p.startQEMUVM()
}

// Stop stops the VM
func (p *VirtualizationFrameworkProvider) Stop() error {
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
func (p *VirtualizationFrameworkProvider) Destroy() error {
	if p.running {
		p.Stop()
	}

	return os.RemoveAll(p.vmPath)
}

// IsRunning checks if the VM is currently running
func (p *VirtualizationFrameworkProvider) IsRunning() bool {
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
func (p *VirtualizationFrameworkProvider) GetInfo() (*VMInfo, error) {
	return &VMInfo{
		Name:       p.config.Name,
		Status:     p.getStatus(),
		Platform:   "macOS",
		Provider:   "Virtualization.framework",
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
		},
	}, nil
}

// RunContainer runs a container inside the VM
func (p *VirtualizationFrameworkProvider) RunContainer(config *ContainerConfig) (*ContainerResult, error) {
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
func (p *VirtualizationFrameworkProvider) ListContainers() ([]*ContainerInfo, error) {
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
func (p *VirtualizationFrameworkProvider) StopContainer(id string) error {
	return p.executeDockerCommand(fmt.Sprintf("docker stop %s", id))
}

// RemoveContainer removes a container in the VM
func (p *VirtualizationFrameworkProvider) RemoveContainer(id string) error {
	return p.executeDockerCommand(fmt.Sprintf("docker rm %s", id))
}

// CopyToVM copies files from host to VM
func (p *VirtualizationFrameworkProvider) CopyToVM(hostPath, vmPath string) error {
	cmd := exec.Command("scp",
		"-P", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		hostPath,
		fmt.Sprintf("root@localhost:%s", vmPath))

	return cmd.Run()
}

// CopyFromVM copies files from VM to host
func (p *VirtualizationFrameworkProvider) CopyFromVM(vmPath, hostPath string) error {
	cmd := exec.Command("scp",
		"-P", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		fmt.Sprintf("root@localhost:%s", vmPath),
		hostPath)

	return cmd.Run()
}

// ForwardPort forwards a port from host to VM
func (p *VirtualizationFrameworkProvider) ForwardPort(hostPort, vmPort int) error {
	// Port forwarding is typically configured during VM creation
	// For dynamic forwarding, we'd need to use SSH tunneling
	return fmt.Errorf("dynamic port forwarding not implemented")
}

// RemovePortForward removes a port forward
func (p *VirtualizationFrameworkProvider) RemovePortForward(hostPort int) error {
	return fmt.Errorf("dynamic port forwarding not implemented")
}

// Helper methods

func (p *VirtualizationFrameworkProvider) isUTMAvailable() bool {
	cmd := exec.Command("which", "utm")
	return cmd.Run() == nil
}

func (p *VirtualizationFrameworkProvider) createUTMVM(config *VMConfig) error {
	// UTM VM creation would require UTM CLI or API
	// For now, fall back to QEMU
	return p.createQEMUVM(config)
}

func (p *VirtualizationFrameworkProvider) createQEMUVM(config *VMConfig) error {
	// Download Alpine Linux ISO if not exists
	diskPath := filepath.Join(p.vmPath, "disk.qcow2")

	// Create disk image
	cmd := exec.Command("qemu-img", "create", "-f", "qcow2", diskPath, fmt.Sprintf("%dG", config.DiskSize))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create disk image: %v", err)
	}

	// Note: In a real implementation, we'd download and configure the Linux ISO
	// with cloud-init or similar for automated setup

	return nil
}

func (p *VirtualizationFrameworkProvider) startUTMVM() error {
	// Start UTM VM
	return fmt.Errorf("UTM integration not implemented")
}

func (p *VirtualizationFrameworkProvider) startQEMUVM() error {
	diskPath := filepath.Join(p.vmPath, "disk.qcow2")

	// QEMU command with Virtualization.framework acceleration
	args := []string{
		"-M", "virt,accel=hvf", // Use Hypervisor.framework
		"-cpu", "host",
		"-smp", strconv.Itoa(p.config.CPUs),
		"-m", strconv.Itoa(p.config.Memory),
		"-drive", fmt.Sprintf("file=%s,if=virtio", diskPath),
		"-netdev", fmt.Sprintf("user,id=net0,hostfwd=tcp::%d-:22,hostfwd=tcp::%d-:2375", p.sshPort, p.config.DockerPort),
		"-device", "virtio-net-pci,netdev=net0",
		"-nographic",
		"-daemonize",
	}

	cmd := exec.Command("qemu-system-aarch64", args...)
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

func (p *VirtualizationFrameworkProvider) getStatus() string {
	if p.IsRunning() {
		return "running"
	}
	return "stopped"
}

func (p *VirtualizationFrameworkProvider) buildDockerCommand(config *ContainerConfig) string {
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

func (p *VirtualizationFrameworkProvider) executeDockerCommand(dockerCmd string) error {
	cmd := exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"root@localhost",
		dockerCmd)

	return cmd.Run()
}

func (p *VirtualizationFrameworkProvider) parseContainerList(output string) []*ContainerInfo {
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
