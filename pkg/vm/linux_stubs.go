//go:build !linux

package vm

import "fmt"

// NewKVMProvider stub for non-Linux platforms
func NewKVMProvider(config *VMConfig) (VMProvider, error) {
	return nil, fmt.Errorf("KVM provider only available on Linux")
}
