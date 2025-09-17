package container

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"servin/pkg/network"
	"servin/pkg/vm"
)

// Helper function to convert PortMappings to map[string]string
func convertPortMappings(portMappings []network.PortMapping) map[string]string {
	ports := make(map[string]string)
	for _, pm := range portMappings {
		ports[fmt.Sprintf("%d", pm.HostPort)] = fmt.Sprintf("%d", pm.ContainerPort)
	}
	return ports
}

// VMContainerManager manages containers through VMs for true cross-platform containerization
type VMContainerManager struct {
	vmManager *vm.VMManager
	vmConfig  *vm.VMConfig
	enabled   bool
}

// NewVMContainerManager creates a new VM-based container manager
func NewVMContainerManager() (*VMContainerManager, error) {
	// Check if VM mode should be enabled
	enabled := shouldEnableVMMode()

	if !enabled {
		return &VMContainerManager{enabled: false}, nil
	}

	// Create VM configuration
	vmConfig := vm.DefaultVMConfig("servin-vm")

	// Customize VM config based on environment
	if customConfig := getCustomVMConfig(); customConfig != nil {
		vmConfig = customConfig
	}

	// Create VM manager
	vmManager, err := vm.NewVMManager(vmConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create VM manager: %v", err)
	}

	return &VMContainerManager{
		vmManager: vmManager,
		vmConfig:  vmConfig,
		enabled:   true,
	}, nil
}

// IsEnabled returns whether VM mode is enabled
func (vcm *VMContainerManager) IsEnabled() bool {
	return vcm.enabled
}

// EnsureVMRunning ensures the VM is running and ready for containers
func (vcm *VMContainerManager) EnsureVMRunning() error {
	if !vcm.enabled {
		return fmt.Errorf("VM mode is not enabled")
	}

	return vcm.vmManager.EnsureRunning()
}

// RunContainer runs a container in the VM
func (vcm *VMContainerManager) RunContainer(container *Container) (*VMContainerResult, error) {
	if !vcm.enabled {
		return nil, fmt.Errorf("VM mode is not enabled")
	}

	// Ensure VM is running
	if err := vcm.EnsureVMRunning(); err != nil {
		return nil, fmt.Errorf("failed to ensure VM is running: %v", err)
	}

	// Convert Servin container config to VM container config
	vmContainerConfig := &vm.ContainerConfig{
		Image:       container.Config.Image,
		Name:        container.Config.Name,
		Command:     append([]string{container.Config.Command}, container.Config.Args...),
		Environment: container.Config.Env,
		Ports:       convertPortMappings(container.Config.PortMappings),
		Volumes:     container.Config.Volumes,
		WorkDir:     container.Config.WorkDir,
		Detached:    true, // Always run detached in VM
	}

	// Run container in VM
	result, err := vcm.vmManager.RunContainer(vmContainerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to run container in VM: %v", err)
	}

	// Convert VM result to Servin result
	return &VMContainerResult{
		ContainerID: result.ID,
		Name:        result.Name,
		Status:      result.Status,
		Output:      result.Output,
		Error:       result.Error,
		ExitCode:    result.ExitCode,
		VMInfo:      vcm.getVMInfo(),
	}, nil
}

// GetVMInfo returns information about the VM
func (vcm *VMContainerManager) GetVMInfo() (*vm.VMInfo, error) {
	if !vcm.enabled {
		return nil, fmt.Errorf("VM mode is not enabled")
	}

	return vcm.vmManager.Provider.GetInfo()
}

// ListVMContainers lists containers running in the VM
func (vcm *VMContainerManager) ListVMContainers() ([]*vm.ContainerInfo, error) {
	if !vcm.enabled {
		return nil, fmt.Errorf("VM mode is not enabled")
	}

	return vcm.vmManager.Provider.ListContainers()
}

// StopVMContainer stops a container in the VM
func (vcm *VMContainerManager) StopVMContainer(containerID string) error {
	if !vcm.enabled {
		return fmt.Errorf("VM mode is not enabled")
	}

	return vcm.vmManager.Provider.StopContainer(containerID)
}

// RemoveVMContainer removes a container in the VM
func (vcm *VMContainerManager) RemoveVMContainer(containerID string) error {
	if !vcm.enabled {
		return fmt.Errorf("VM mode is not enabled")
	}

	return vcm.vmManager.Provider.RemoveContainer(containerID)
}

// Shutdown gracefully shuts down the VM
func (vcm *VMContainerManager) Shutdown() error {
	if !vcm.enabled {
		return nil
	}

	return vcm.vmManager.Shutdown()
}

// VMContainerResult represents the result of running a container in a VM
type VMContainerResult struct {
	ContainerID string  `json:"container_id"`
	Name        string  `json:"name"`
	Status      string  `json:"status"`
	Output      string  `json:"output"`
	Error       string  `json:"error"`
	ExitCode    int     `json:"exit_code"`
	VMInfo      *VMInfo `json:"vm_info"`
}

// VMInfo provides information about the VM state
type VMInfo struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	Platform   string `json:"platform"`
	Provider   string `json:"provider"`
	IPAddress  string `json:"ip_address"`
	SSHPort    int    `json:"ssh_port"`
	DockerPort int    `json:"docker_port"`
}

// Helper functions

func shouldEnableVMMode() bool {
	// Check environment variable
	if mode := os.Getenv("SERVIN_VM_MODE"); mode != "" {
		return strings.ToLower(mode) == "true" || mode == "1"
	}

	// Check configuration file
	if hasVMConfig() {
		return true
	}

	// Auto-detect based on platform capabilities
	return autoDetectVMMode()
}

func hasVMConfig() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	configPath := filepath.Join(homeDir, ".servin", "vm-config.json")
	_, err = os.Stat(configPath)
	return err == nil
}

func autoDetectVMMode() bool {
	// On non-Linux platforms, default to VM mode for true containerization
	// On Linux, users can choose between native containers and VM containers
	return true // Enable VM mode by default for consistency
}

func getCustomVMConfig() *vm.VMConfig {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	_ = filepath.Join(homeDir, ".servin", "vm-config.json")

	// In a real implementation, we would load and parse the JSON config
	// For now, return nil to use defaults

	return nil
}

func (vcm *VMContainerManager) getVMInfo() *VMInfo {
	info, err := vcm.vmManager.Provider.GetInfo()
	if err != nil {
		return &VMInfo{
			Name:     "unknown",
			Status:   "error",
			Platform: "unknown",
		}
	}

	return &VMInfo{
		Name:       info.Name,
		Status:     info.Status,
		Platform:   info.Platform,
		Provider:   info.Provider,
		IPAddress:  info.IPAddress,
		SSHPort:    info.SSHPort,
		DockerPort: info.DockerPort,
	}
}

// Integration with existing Container.Run method

// RunWithVM runs a container using VM if enabled, falls back to native if not
func (c *Container) RunWithVM() error {
	// Try VM mode first if available
	vmManager, err := NewVMContainerManager()
	if err == nil && vmManager.IsEnabled() {
		result, vmErr := vmManager.RunContainer(c)
		if vmErr == nil {
			// Update container with VM result
			c.ID = result.ContainerID
			c.Status = result.Status
			c.UpdateStatus(result.Status)

			fmt.Printf("Container %s running in VM (%s)\n", c.Config.Name, result.VMInfo.Provider)
			return nil
		}

		fmt.Printf("VM mode failed, falling back to native: %v\n", vmErr)
	}

	// Fall back to original implementation
	fmt.Printf("Running container natively (limited containerization)\n")
	return c.Run() // Original method
}
