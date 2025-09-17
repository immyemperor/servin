//go:build !darwin

package vm

import "fmt"

// NewVirtualizationFrameworkProvider stub for non-macOS platforms
func NewVirtualizationFrameworkProvider(config *VMConfig) (VMProvider, error) {
	return nil, fmt.Errorf("Virtualization.framework provider only available on macOS")
}

// NewSimplifiedLinuxVMProvider stub for non-macOS platforms
func NewSimplifiedLinuxVMProvider(config *VMConfig) (VMProvider, error) {
	return nil, fmt.Errorf("SimplifiedLinuxVM provider only available on macOS")
}
