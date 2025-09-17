package vm

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// VMProvider represents different virtualization backends per platform
type VMProvider interface {
	// VM lifecycle management
	Create(config *VMConfig) error
	Start() error
	Stop() error
	Destroy() error

	// VM state management
	IsRunning() bool
	GetInfo() (*VMInfo, error)

	// Container operations inside VM
	RunContainer(config *ContainerConfig) (*ContainerResult, error)
	ListContainers() ([]*ContainerInfo, error)
	StopContainer(id string) error
	RemoveContainer(id string) error

	// File operations
	CopyToVM(hostPath, vmPath string) error
	CopyFromVM(vmPath, hostPath string) error

	// Network operations
	ForwardPort(hostPort, vmPort int) error
	RemovePortForward(hostPort int) error
}

// VMConfig represents VM configuration
type VMConfig struct {
	Name             string            `json:"name"`
	CPUs             int               `json:"cpus"`
	Memory           int               `json:"memory_mb"`
	DiskSize         int               `json:"disk_size_gb"`
	LinuxDistro      string            `json:"linux_distro"`      // "alpine", "ubuntu", "debian"
	ContainerRuntime string            `json:"container_runtime"` // "docker", "containerd", "podman"
	SSHPort          int               `json:"ssh_port"`
	DockerPort       int               `json:"docker_port"`
	WorkDir          string            `json:"work_dir"`
	Environment      map[string]string `json:"environment"`
}

// VMInfo represents VM status and information
type VMInfo struct {
	Name         string          `json:"name"`
	Status       string          `json:"status"`
	Platform     string          `json:"platform"`
	Provider     string          `json:"provider"`
	CPUs         int             `json:"cpus"`
	Memory       int             `json:"memory_mb"`
	DiskUsage    int             `json:"disk_usage_mb"`
	Uptime       string          `json:"uptime"`
	IPAddress    string          `json:"ip_address"`
	SSHPort      int             `json:"ssh_port"`
	DockerPort   int             `json:"docker_port"`
	Capabilities map[string]bool `json:"capabilities"`
}

// ContainerConfig represents container configuration for VM
type ContainerConfig struct {
	Image       string            `json:"image"`
	Name        string            `json:"name"`
	Command     []string          `json:"command"`
	Environment map[string]string `json:"environment"`
	Ports       map[string]string `json:"ports"`
	Volumes     map[string]string `json:"volumes"`
	WorkDir     string            `json:"workdir"`
	Detached    bool              `json:"detached"`
}

// ContainerResult represents container execution result
type ContainerResult struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	ExitCode int    `json:"exit_code"`
	Output   string `json:"output"`
	Error    string `json:"error"`
}

// ContainerInfo represents container information
type ContainerInfo struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Image   string            `json:"image"`
	Status  string            `json:"status"`
	Ports   map[string]string `json:"ports"`
	Created string            `json:"created"`
	Command string            `json:"command"`
}

// GetVMProvider returns the appropriate VM provider for the current platform
func GetVMProvider(config *VMConfig) (VMProvider, error) {
	// Check if we're in development mode
	if isDevelopmentMode() {
		// Use universal development provider for all platforms in dev mode
		return NewDevelopmentVMProvider(config)
	}

	// Production VM providers
	switch runtime.GOOS {
	case "darwin": // macOS
		return NewVirtualizationFrameworkProvider(config)
	case "windows":
		return NewHyperVProvider(config)
	case "linux":
		return NewKVMProvider(config)
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// isDevelopmentMode checks if we're running in development mode
func isDevelopmentMode() bool {
	// Check environment variable
	if mode := os.Getenv("SERVIN_DEV_MODE"); mode != "" {
		return strings.ToLower(mode) == "true" || mode == "1"
	}

	// Check if --dev flag was used (this is a simple heuristic)
	for _, arg := range os.Args {
		if arg == "--dev" {
			return true
		}
	}

	return false
}

// DefaultVMConfig returns a sensible default VM configuration
func DefaultVMConfig(name string) *VMConfig {
	return &VMConfig{
		Name:             name,
		CPUs:             2,
		Memory:           2048, // 2GB
		DiskSize:         20,   // 20GB
		LinuxDistro:      "alpine",
		ContainerRuntime: "docker",
		SSHPort:          2222,
		DockerPort:       2375,
		WorkDir:          "/servin",
		Environment: map[string]string{
			"SERVIN_VM": "true",
		},
	}
}

// VMManager manages VM lifecycle and operations
type VMManager struct {
	Provider VMProvider
	Config   *VMConfig
}

// NewVMManager creates a new VM manager
func NewVMManager(config *VMConfig) (*VMManager, error) {
	provider, err := GetVMProvider(config)
	if err != nil {
		return nil, fmt.Errorf("failed to get VM provider: %v", err)
	}

	return &VMManager{
		Provider: provider,
		Config:   config,
	}, nil
}

// EnsureRunning ensures the VM is created and running
func (vm *VMManager) EnsureRunning() error {
	if vm.Provider.IsRunning() {
		return nil
	}

	// Check if VM exists, create if not
	info, err := vm.Provider.GetInfo()
	if err != nil || info == nil {
		if err := vm.Provider.Create(vm.Config); err != nil {
			return fmt.Errorf("failed to create VM: %v", err)
		}
	}

	// Start the VM
	if err := vm.Provider.Start(); err != nil {
		return fmt.Errorf("failed to start VM: %v", err)
	}

	return nil
}

// RunContainer runs a container inside the VM
func (vm *VMManager) RunContainer(config *ContainerConfig) (*ContainerResult, error) {
	if err := vm.EnsureRunning(); err != nil {
		return nil, err
	}

	return vm.Provider.RunContainer(config)
}

// Shutdown gracefully shuts down the VM
func (vm *VMManager) Shutdown() error {
	return vm.Provider.Stop()
}
