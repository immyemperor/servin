package vm

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// UniversalDevelopmentVMProvider provides a simulated VM for development and testing
// This allows testing VM functionality without requiring actual VM setup on non-macOS platforms
type UniversalDevelopmentVMProvider struct {
	config     *VMConfig
	vmPath     string
	sshPort    int
	running    bool
	containers map[string]*ContainerInfo
}

// NewDevelopmentVMProvider creates a new development VM provider for all platforms
func NewDevelopmentVMProvider(config *VMConfig) (VMProvider, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	vmPath := filepath.Join(homeDir, ".servin", "dev-vm", config.Name)

	provider := &UniversalDevelopmentVMProvider{
		config:     config,
		vmPath:     vmPath,
		sshPort:    config.SSHPort,
		running:    false,
		containers: make(map[string]*ContainerInfo),
	}

	// Check if VM is already running by reading state file
	provider.loadRunningState()

	return provider, nil
}

// Create creates a simulated VM setup
func (p *UniversalDevelopmentVMProvider) Create(config *VMConfig) error {
	fmt.Println("Creating development VM (simulated)...")

	// Ensure VM directory exists
	if err := os.MkdirAll(p.vmPath, 0755); err != nil {
		return fmt.Errorf("failed to create VM directory: %v", err)
	}

	fmt.Printf("Development VM created at: %s\n", p.vmPath)
	return nil
}

// Start starts the simulated VM
func (p *UniversalDevelopmentVMProvider) Start() error {
	if p.running {
		return fmt.Errorf("VM is already running")
	}

	fmt.Println("Starting development VM (simulated)...")
	p.running = true
	p.saveRunningState() // Save state to disk
	fmt.Println("Development VM started successfully")
	return nil
}

// Stop stops the simulated VM
func (p *UniversalDevelopmentVMProvider) Stop() error {
	if !p.running {
		return fmt.Errorf("VM is not running")
	}

	fmt.Println("Stopping development VM (simulated)...")
	p.running = false
	p.saveRunningState() // Save state to disk
	fmt.Println("Development VM stopped successfully")
	return nil
}

// Destroy destroys the simulated VM
func (p *UniversalDevelopmentVMProvider) Destroy() error {
	if p.running {
		if err := p.Stop(); err != nil {
			return err
		}
	}

	fmt.Println("Destroying development VM (simulated)...")
	if err := os.RemoveAll(p.vmPath); err != nil {
		return fmt.Errorf("failed to remove VM directory: %v", err)
	}

	fmt.Println("Development VM destroyed successfully")
	return nil
}

// IsRunning returns whether the simulated VM is running
func (p *UniversalDevelopmentVMProvider) IsRunning() bool {
	return p.running
}

// GetInfo returns information about the simulated VM
func (p *UniversalDevelopmentVMProvider) GetInfo() (*VMInfo, error) {
	status := "stopped"
	if p.running {
		status = "running"
	}

	return &VMInfo{
		Name:       p.config.Name,
		Status:     status,
		Provider:   "Development (Simulated)",
		Platform:   runtime.GOOS,
		CPUs:       p.config.CPUs,
		Memory:     p.config.Memory,
		IPAddress:  "127.0.0.1",
		SSHPort:    p.sshPort,
		DockerPort: p.config.DockerPort,
	}, nil
}

// RunContainer simulates running a container in the VM
func (p *UniversalDevelopmentVMProvider) RunContainer(config *ContainerConfig) (*ContainerResult, error) {
	if !p.running {
		return nil, fmt.Errorf("VM is not running")
	}

	containerID := fmt.Sprintf("dev-%d", time.Now().Unix())

	// Simulate container creation
	container := &ContainerInfo{
		ID:      containerID,
		Name:    config.Name,
		Image:   config.Image,
		Status:  "running",
		Created: time.Now().Format(time.RFC3339),
	}

	p.containers[containerID] = container

	fmt.Printf("Simulated container started: %s (ID: %s)\n", config.Name, containerID)

	return &ContainerResult{
		ID:       containerID,
		Name:     config.Name,
		Status:   "running",
		ExitCode: 0,
		Output:   "Container started successfully (simulated)",
	}, nil
}

// ListContainers returns list of simulated containers
func (p *UniversalDevelopmentVMProvider) ListContainers() ([]*ContainerInfo, error) {
	if !p.running {
		return nil, fmt.Errorf("VM is not running")
	}

	containers := make([]*ContainerInfo, 0, len(p.containers))
	for _, container := range p.containers {
		containers = append(containers, container)
	}

	return containers, nil
}

// StopContainer simulates stopping a container
func (p *UniversalDevelopmentVMProvider) StopContainer(id string) error {
	if !p.running {
		return fmt.Errorf("VM is not running")
	}

	container, exists := p.containers[id]
	if !exists {
		return fmt.Errorf("container not found: %s", id)
	}

	container.Status = "stopped"
	fmt.Printf("Simulated container stopped: %s\n", id)
	return nil
}

// RemoveContainer simulates removing a container
func (p *UniversalDevelopmentVMProvider) RemoveContainer(id string) error {
	if !p.running {
		return fmt.Errorf("VM is not running")
	}

	if _, exists := p.containers[id]; !exists {
		return fmt.Errorf("container not found: %s", id)
	}

	delete(p.containers, id)
	fmt.Printf("Simulated container removed: %s\n", id)
	return nil
}

// CopyToVM simulates copying files to VM
func (p *UniversalDevelopmentVMProvider) CopyToVM(hostPath, vmPath string) error {
	if !p.running {
		return fmt.Errorf("VM is not running")
	}

	fmt.Printf("Simulated copy to VM: %s -> %s\n", hostPath, vmPath)
	return nil
}

// CopyFromVM simulates copying files from VM
func (p *UniversalDevelopmentVMProvider) CopyFromVM(vmPath, hostPath string) error {
	if !p.running {
		return fmt.Errorf("VM is not running")
	}

	fmt.Printf("Simulated copy from VM: %s -> %s\n", vmPath, hostPath)
	return nil
}

// ForwardPort simulates port forwarding
func (p *UniversalDevelopmentVMProvider) ForwardPort(hostPort, vmPort int) error {
	if !p.running {
		return fmt.Errorf("VM is not running")
	}

	fmt.Printf("Simulated port forward: %d -> %d\n", hostPort, vmPort)
	return nil
}

// RemovePortForward simulates removing port forwarding
func (p *UniversalDevelopmentVMProvider) RemovePortForward(hostPort int) error {
	if !p.running {
		return fmt.Errorf("VM is not running")
	}

	fmt.Printf("Simulated port forward removed: %d\n", hostPort)
	return nil
}

// loadRunningState loads the VM running state from disk
func (p *UniversalDevelopmentVMProvider) loadRunningState() {
	stateFile := filepath.Join(p.vmPath, "vm-running")
	if _, err := os.Stat(stateFile); err == nil {
		p.running = true
	} else {
		p.running = false
	}
}

// saveRunningState saves the VM running state to disk
func (p *UniversalDevelopmentVMProvider) saveRunningState() error {
	if err := os.MkdirAll(p.vmPath, 0755); err != nil {
		return err
	}

	stateFile := filepath.Join(p.vmPath, "vm-running")
	if p.running {
		return os.WriteFile(stateFile, []byte("running"), 0644)
	} else {
		// Remove the state file if VM is stopped
		os.Remove(stateFile)
		return nil
	}
}
