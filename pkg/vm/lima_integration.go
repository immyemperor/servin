//go:build darwin

package vm

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// LimaVM manages a Lima-based Linux VM for true containerization on macOS
type LimaVM struct {
	Name       string
	ConfigPath string
	Running    bool
}

// NewLimaVM creates a new Lima VM instance for Servin
func NewLimaVM(name string) *LimaVM {
	return &LimaVM{
		Name:       name,
		ConfigPath: filepath.Join("~/.lima", name, "lima.yaml"),
		Running:    false,
	}
}

// Start initializes and starts the Lima VM
func (vm *LimaVM) Start() error {
	// Check if Lima is installed
	if !vm.isLimaInstalled() {
		return fmt.Errorf("Lima is not installed. Please install with: brew install lima")
	}

	// Start the VM
	cmd := exec.Command("limactl", "start", vm.Name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start Lima VM: %v", err)
	}

	vm.Running = true
	return nil
}

// Stop shuts down the Lima VM
func (vm *LimaVM) Stop() error {
	cmd := exec.Command("limactl", "stop", vm.Name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop Lima VM: %v", err)
	}

	vm.Running = false
	return nil
}

// RunContainer executes a container inside the Lima VM
func (vm *LimaVM) RunContainer(containerCmd string) error {
	if !vm.Running {
		return fmt.Errorf("Lima VM is not running")
	}

	// Execute container command inside Lima VM
	cmd := exec.Command("limactl", "shell", vm.Name, containerCmd)
	return cmd.Run()
}

// IsRunning checks if the Lima VM is currently running
func (vm *LimaVM) IsRunning() bool {
	cmd := exec.Command("limactl", "list", "--format", "{{.Status}}", vm.Name)
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return string(output) == "Running\n"
}

// isLimaInstalled checks if Lima is available on the system
func (vm *LimaVM) isLimaInstalled() bool {
	cmd := exec.Command("which", "limactl")
	return cmd.Run() == nil
}

// GetVMInfo returns information about the Lima VM
func (vm *LimaVM) GetVMInfo() (map[string]string, error) {
	info := make(map[string]string)

	// Get VM status
	cmd := exec.Command("limactl", "list", "--format", "{{.Name}}\t{{.Status}}\t{{.Arch}}\t{{.CPUs}}\t{{.Memory}}", vm.Name)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get VM info: %v", err)
	}

	info["status"] = string(output)
	info["type"] = "Lima VM"
	return info, nil
}
