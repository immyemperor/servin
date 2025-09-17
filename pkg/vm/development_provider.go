//go:build ignore

package vm

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// DevelopmentVMProvider provides a simulated VM for development and testing
// This allows testing VM functionality without requiring actual VM setup
type DevelopmentVMProvider struct {
	config     *VMConfig
	vmPath     string
	sshPort    int
	running    bool
	containers map[string]*ContainerInfo
}

// NewDevelopmentVMProvider creates a new development VM provider
func NewDevelopmentVMProvider(config *VMConfig) (VMProvider, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	vmPath := filepath.Join(homeDir, ".servin", "dev-vm", config.Name)

	provider := &DevelopmentVMProvider{
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
}

// Create creates a simulated VM
func (p *DevelopmentVMProvider) Create(config *VMConfig) error {
	// Ensure VM directory exists
	if err := os.MkdirAll(p.vmPath, 0755); err != nil {
		return fmt.Errorf("failed to create VM directory: %v", err)
	}

	// Create a marker file to indicate VM is "created"
	markerFile := filepath.Join(p.vmPath, "vm-created")
	if err := os.WriteFile(markerFile, []byte("development-vm"), 0644); err != nil {
		return fmt.Errorf("failed to create VM marker: %v", err)
	}

	fmt.Printf("Development VM created at %s\n", p.vmPath)
	fmt.Println("Note: This is a simulated VM for development/testing purposes")

	return nil
}

// Start starts the simulated VM
func (p *DevelopmentVMProvider) Start() error {
	if p.running {
		return nil
	}

	// Ensure VM is created first
	markerFile := filepath.Join(p.vmPath, "vm-created")
	if _, err := os.Stat(markerFile); os.IsNotExist(err) {
		// Auto-create VM if not exists
		if err := p.Create(p.config); err != nil {
			return fmt.Errorf("failed to auto-create VM: %v", err)
		}
	}

	// Simulate VM startup time
	fmt.Println("Simulating VM startup (development mode)...")
	time.Sleep(2 * time.Second)

	p.running = true
	p.saveRunningState() // Save state to disk
	fmt.Printf("Development VM started successfully!\n")
	fmt.Printf("VM is simulated - containers will run in isolated namespaces\n")

	return nil
}

// Stop stops the simulated VM
func (p *DevelopmentVMProvider) Stop() error {
	if !p.running {
		return nil
	}

	fmt.Println("Stopping development VM...")
	time.Sleep(1 * time.Second)

	p.running = false
	p.saveRunningState() // Save state to disk
	fmt.Println("Development VM stopped")

	return nil
}

// Destroy removes the simulated VM
func (p *DevelopmentVMProvider) Destroy() error {
	if p.running {
		p.Stop()
	}

	fmt.Println("Destroying development VM...")
	err := os.RemoveAll(p.vmPath)
	if err != nil {
		return fmt.Errorf("failed to remove VM directory: %v", err)
	}

	p.containers = make(map[string]*ContainerInfo)
	fmt.Println("Development VM destroyed")

	return nil
}

// IsRunning checks if the simulated VM is running
func (p *DevelopmentVMProvider) IsRunning() bool {
	return p.running
}

// GetInfo returns VM information
func (p *DevelopmentVMProvider) GetInfo() (*VMInfo, error) {
	status := "stopped"
	if p.running {
		status = "running"
	}

	return &VMInfo{
		Name:       p.config.Name,
		Status:     status,
		Platform:   "macOS",
		Provider:   "Development (Simulated)",
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
			"development":  true,
		},
	}, nil
}

// RunContainer runs a simulated container
func (p *DevelopmentVMProvider) RunContainer(config *ContainerConfig) (*ContainerResult, error) {
	if !p.IsRunning() {
		return nil, fmt.Errorf("VM is not running")
	}

	// Generate a container ID
	containerID := fmt.Sprintf("dev_%d", time.Now().UnixNano())

	// Simulate container creation
	fmt.Printf("Simulating container creation: %s\n", config.Image)
	time.Sleep(1 * time.Second)

	// Store container info
	p.containers[containerID] = &ContainerInfo{
		ID:      containerID,
		Name:    config.Name,
		Image:   config.Image,
		Status:  "running",
		Created: time.Now().Format(time.RFC3339),
		Command: fmt.Sprintf("%s", config.Command),
	}

	result := &ContainerResult{
		ID:       containerID,
		Name:     config.Name,
		Status:   "running",
		Output:   fmt.Sprintf("Development container %s started successfully", config.Name),
		ExitCode: 0,
	}

	fmt.Printf("Development container %s (%s) created successfully\n", config.Name, containerID[:12])

	return result, nil
}

// ListContainers lists simulated containers
func (p *DevelopmentVMProvider) ListContainers() ([]*ContainerInfo, error) {
	if !p.IsRunning() {
		return nil, fmt.Errorf("VM is not running")
	}

	containers := make([]*ContainerInfo, 0, len(p.containers))
	for _, container := range p.containers {
		containers = append(containers, container)
	}

	return containers, nil
}

// StopContainer stops a simulated container
func (p *DevelopmentVMProvider) StopContainer(id string) error {
	if !p.IsRunning() {
		return fmt.Errorf("VM is not running")
	}

	container, exists := p.containers[id]
	if !exists {
		return fmt.Errorf("container %s not found", id)
	}

	container.Status = "stopped"
	fmt.Printf("Development container %s stopped\n", id[:12])

	return nil
}

// RemoveContainer removes a simulated container
func (p *DevelopmentVMProvider) RemoveContainer(id string) error {
	if !p.IsRunning() {
		return fmt.Errorf("VM is not running")
	}

	_, exists := p.containers[id]
	if !exists {
		return fmt.Errorf("container %s not found", id)
	}

	delete(p.containers, id)
	fmt.Printf("Development container %s removed\n", id[:12])

	return nil
}

// CopyToVM simulates copying files to VM
func (p *DevelopmentVMProvider) CopyToVM(hostPath, vmPath string) error {
	if !p.IsRunning() {
		return fmt.Errorf("VM is not running")
	}

	fmt.Printf("Simulating copy: %s -> VM:%s\n", hostPath, vmPath)
	return nil
}

// CopyFromVM simulates copying files from VM
func (p *DevelopmentVMProvider) CopyFromVM(vmPath, hostPath string) error {
	if !p.IsRunning() {
		return fmt.Errorf("VM is not running")
	}

	fmt.Printf("Simulating copy: VM:%s -> %s\n", vmPath, hostPath)
	return nil
}

// ForwardPort simulates port forwarding
func (p *DevelopmentVMProvider) ForwardPort(hostPort, vmPort int) error {
	fmt.Printf("Simulating port forward: %d -> %d\n", hostPort, vmPort)
	return nil
}

// loadRunningState loads the VM running state from disk
func (p *DevelopmentVMProvider) loadRunningState() {
	stateFile := filepath.Join(p.vmPath, "vm-running")
	if _, err := os.Stat(stateFile); err == nil {
		p.running = true
	} else {
		p.running = false
	}
}

// saveRunningState saves the VM running state to disk
func (p *DevelopmentVMProvider) saveRunningState() error {
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

// RemovePortForward simulates removing port forwarding
func (p *DevelopmentVMProvider) RemovePortForward(hostPort int) error {
	fmt.Printf("Simulating removal of port forward: %d\n", hostPort)
	return nil
}
