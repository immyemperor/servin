//go:build !linux

package cgroups

import "fmt"

// CGroup manages container resource limits (placeholder for non-Linux)
type CGroup struct {
	ContainerID string
	Path        string
}

// New creates a new CGroup manager (non-Linux placeholder)
func New(containerID string) *CGroup {
	return &CGroup{
		ContainerID: containerID,
		Path:        "",
	}
}

// Create returns an error on non-Linux platforms
func (c *CGroup) Create() error {
	return fmt.Errorf("cgroups are only supported on Linux")
}

// SetMemoryLimit returns an error on non-Linux platforms
func (c *CGroup) SetMemoryLimit(limitBytes int64) error {
	return fmt.Errorf("cgroups are only supported on Linux")
}

// SetCPULimit returns an error on non-Linux platforms
func (c *CGroup) SetCPULimit(shares int) error {
	return fmt.Errorf("cgroups are only supported on Linux")
}

// SetPIDLimit returns an error on non-Linux platforms
func (c *CGroup) SetPIDLimit(max int) error {
	return fmt.Errorf("cgroups are only supported on Linux")
}

// AddProcess returns an error on non-Linux platforms
func (c *CGroup) AddProcess(pid int) error {
	return fmt.Errorf("cgroups are only supported on Linux")
}

// GetStats returns an error on non-Linux platforms
func (c *CGroup) GetStats() (map[string]string, error) {
	return nil, fmt.Errorf("cgroups are only supported on Linux")
}

// Cleanup returns an error on non-Linux platforms
func (c *CGroup) Cleanup() error {
	return fmt.Errorf("cgroups are only supported on Linux")
}

// ParseMemoryString converts memory strings (cross-platform)
func ParseMemoryString(memStr string) (int64, error) {
	return 0, fmt.Errorf("memory parsing not implemented for non-Linux")
}
