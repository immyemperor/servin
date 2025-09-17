//go:build windows

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

// HyperVProvider implements VM operations using Windows Hyper-V
type HyperVProvider struct {
	config  *VMConfig
	vmPath  string
	sshPort int
	running bool
}

// NewHyperVProvider creates a new Hyper-V provider
func NewHyperVProvider(config *VMConfig) (VMProvider, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	vmPath := filepath.Join(homeDir, ".servin", "vms", config.Name)

	return &HyperVProvider{
		config:  config,
		vmPath:  vmPath,
		sshPort: config.SSHPort,
		running: false,
	}, nil
}

// Create creates a new VM using Hyper-V
func (p *HyperVProvider) Create(config *VMConfig) error {
	// Check if Hyper-V is available
	if !p.isHyperVAvailable() {
		// Fall back to VirtualBox if available
		if p.isVirtualBoxAvailable() {
			return p.createVirtualBoxVM(config)
		}
		return fmt.Errorf("neither Hyper-V nor VirtualBox is available")
	}

	// Ensure VM directory exists
	if err := os.MkdirAll(p.vmPath, 0755); err != nil {
		return fmt.Errorf("failed to create VM directory: %v", err)
	}

	return p.createHyperVVM(config)
}

// Start starts the VM
func (p *HyperVProvider) Start() error {
	if p.running {
		return nil
	}

	if p.isHyperVAvailable() {
		return p.startHyperVVM()
	}

	return p.startVirtualBoxVM()
}

// Stop stops the VM
func (p *HyperVProvider) Stop() error {
	if !p.running {
		return nil
	}

	if p.isHyperVAvailable() {
		return p.stopHyperVVM()
	}

	return p.stopVirtualBoxVM()
}

// Destroy removes the VM completely
func (p *HyperVProvider) Destroy() error {
	if p.running {
		p.Stop()
	}

	if p.isHyperVAvailable() {
		// Remove Hyper-V VM
		cmd := exec.Command("powershell", "-Command",
			fmt.Sprintf("Remove-VM -Name '%s' -Force", p.config.Name))
		cmd.Run() // Ignore errors
	}

	return os.RemoveAll(p.vmPath)
}

// IsRunning checks if the VM is currently running
func (p *HyperVProvider) IsRunning() bool {
	// Check if we can connect via SSH
	cmd := exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=NUL",
		"-o", "ConnectTimeout=2",
		"root@localhost",
		"echo alive")

	err := cmd.Run()
	p.running = (err == nil)
	return p.running
}

// GetInfo returns VM information
func (p *HyperVProvider) GetInfo() (*VMInfo, error) {
	provider := "Hyper-V"
	if !p.isHyperVAvailable() {
		provider = "VirtualBox"
	}

	return &VMInfo{
		Name:       p.config.Name,
		Status:     p.getStatus(),
		Platform:   "Windows",
		Provider:   provider,
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
func (p *HyperVProvider) RunContainer(config *ContainerConfig) (*ContainerResult, error) {
	if !p.IsRunning() {
		return nil, fmt.Errorf("VM is not running")
	}

	// Build docker run command
	dockerCmd := p.buildDockerCommand(config)

	// Execute via SSH
	cmd := exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=NUL",
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
func (p *HyperVProvider) ListContainers() ([]*ContainerInfo, error) {
	if !p.IsRunning() {
		return nil, fmt.Errorf("VM is not running")
	}

	cmd := exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=NUL",
		"root@localhost",
		"docker ps -a --format \"table {{.ID}}\\t{{.Names}}\\t{{.Image}}\\t{{.Status}}\\t{{.CreatedAt}}\\t{{.Command}}\"")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %v", err)
	}

	return p.parseContainerList(string(output)), nil
}

// StopContainer stops a container in the VM
func (p *HyperVProvider) StopContainer(id string) error {
	return p.executeDockerCommand(fmt.Sprintf("docker stop %s", id))
}

// RemoveContainer removes a container in the VM
func (p *HyperVProvider) RemoveContainer(id string) error {
	return p.executeDockerCommand(fmt.Sprintf("docker rm %s", id))
}

// CopyToVM copies files from host to VM
func (p *HyperVProvider) CopyToVM(hostPath, vmPath string) error {
	cmd := exec.Command("scp",
		"-P", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=NUL",
		hostPath,
		fmt.Sprintf("root@localhost:%s", vmPath))

	return cmd.Run()
}

// CopyFromVM copies files from VM to host
func (p *HyperVProvider) CopyFromVM(vmPath, hostPath string) error {
	cmd := exec.Command("scp",
		"-P", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=NUL",
		fmt.Sprintf("root@localhost:%s", vmPath),
		hostPath)

	return cmd.Run()
}

// ForwardPort forwards a port from host to VM
func (p *HyperVProvider) ForwardPort(hostPort, vmPort int) error {
	return fmt.Errorf("dynamic port forwarding not implemented")
}

// RemovePortForward removes a port forward
func (p *HyperVProvider) RemovePortForward(hostPort int) error {
	return fmt.Errorf("dynamic port forwarding not implemented")
}

// Helper methods

func (p *HyperVProvider) isHyperVAvailable() bool {
	cmd := exec.Command("powershell", "-Command", "Get-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "Enabled")
}

func (p *HyperVProvider) isVirtualBoxAvailable() bool {
	cmd := exec.Command("VBoxManage", "--version")
	return cmd.Run() == nil
}

func (p *HyperVProvider) createHyperVVM(config *VMConfig) error {
	vhdPath := filepath.Join(p.vmPath, "disk.vhdx")

	// Create VHD
	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf("New-VHD -Path '%s' -SizeBytes %dGB -Dynamic", vhdPath, config.DiskSize))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create VHD: %v", err)
	}

	// Create VM
	cmd = exec.Command("powershell", "-Command",
		fmt.Sprintf("New-VM -Name '%s' -MemoryStartupBytes %dMB -VHDPath '%s' -Generation 2",
			config.Name, config.Memory, vhdPath))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create VM: %v", err)
	}

	// Configure VM
	cmd = exec.Command("powershell", "-Command",
		fmt.Sprintf("Set-VM -Name '%s' -ProcessorCount %d", config.Name, config.CPUs))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to configure VM: %v", err)
	}

	return nil
}

func (p *HyperVProvider) createVirtualBoxVM(config *VMConfig) error {
	// Create VirtualBox VM
	cmd := exec.Command("VBoxManage", "createvm",
		"--name", config.Name,
		"--ostype", "Linux_64",
		"--register")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create VirtualBox VM: %v", err)
	}

	// Configure VM
	commands := [][]string{
		{"modifyvm", config.Name, "--memory", strconv.Itoa(config.Memory)},
		{"modifyvm", config.Name, "--cpus", strconv.Itoa(config.CPUs)},
		{"modifyvm", config.Name, "--nic1", "nat"},
		{"modifyvm", config.Name, "--natpf1", fmt.Sprintf("ssh,tcp,,%d,,22", p.sshPort)},
		{"modifyvm", config.Name, "--natpf1", fmt.Sprintf("docker,tcp,,%d,,2375", config.DockerPort)},
	}

	for _, cmdArgs := range commands {
		cmd := exec.Command("VBoxManage", cmdArgs...)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to configure VM: %v", err)
		}
	}

	return nil
}

func (p *HyperVProvider) startHyperVVM() error {
	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf("Start-VM -Name '%s'", p.config.Name))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start Hyper-V VM: %v", err)
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

func (p *HyperVProvider) startVirtualBoxVM() error {
	cmd := exec.Command("VBoxManage", "startvm", p.config.Name, "--type", "headless")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start VirtualBox VM: %v", err)
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

func (p *HyperVProvider) stopHyperVVM() error {
	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf("Stop-VM -Name '%s' -Force", p.config.Name))
	err := cmd.Run()
	p.running = false
	return err
}

func (p *HyperVProvider) stopVirtualBoxVM() error {
	cmd := exec.Command("VBoxManage", "controlvm", p.config.Name, "poweroff")
	err := cmd.Run()
	p.running = false
	return err
}

func (p *HyperVProvider) getStatus() string {
	if p.IsRunning() {
		return "running"
	}
	return "stopped"
}

func (p *HyperVProvider) buildDockerCommand(config *ContainerConfig) string {
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

func (p *HyperVProvider) executeDockerCommand(dockerCmd string) error {
	cmd := exec.Command("ssh",
		"-p", strconv.Itoa(p.sshPort),
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=NUL",
		"root@localhost",
		dockerCmd)

	return cmd.Run()
}

func (p *HyperVProvider) parseContainerList(output string) []*ContainerInfo {
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
