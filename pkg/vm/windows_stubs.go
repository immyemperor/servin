//go:build !windows

package vm

import "fmt"

// NewHyperVProvider stub for non-Windows platforms
func NewHyperVProvider(config *VMConfig) (VMProvider, error) {
	return nil, fmt.Errorf("Hyper-V provider only available on Windows")
}
