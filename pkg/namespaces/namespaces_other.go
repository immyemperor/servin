//go:build !linux

package namespaces

import "fmt"

// NamespaceFlags represents the Linux namespace types (placeholder for non-Linux)
type NamespaceFlags int

// Placeholder constants for non-Linux platforms
const (
	CLONE_NEWPID  NamespaceFlags = 0
	CLONE_NEWUTS  NamespaceFlags = 0
	CLONE_NEWIPC  NamespaceFlags = 0
	CLONE_NEWNET  NamespaceFlags = 0
	CLONE_NEWNS   NamespaceFlags = 0
	CLONE_NEWUSER NamespaceFlags = 0
)

// ContainerConfig holds namespace configuration (placeholder for non-Linux)
type ContainerConfig struct {
	Command    string
	Args       []string
	Namespaces []NamespaceFlags
	Hostname   string
	WorkDir    string
	LogDir     string // Directory to store container logs

	// User namespace configuration
	UserNamespace *UserNamespaceConfig
}

// CreateContainer returns an error on non-Linux platforms
func CreateContainer(config *ContainerConfig) error {
	return fmt.Errorf("namespace isolation is only supported on Linux")
}

// SetupNamespace returns an error on non-Linux platforms
func SetupNamespace(hostname string) error {
	return fmt.Errorf("namespace setup is only supported on Linux")
}

// GetNamespaceInfo returns an error on non-Linux platforms
func GetNamespaceInfo() (map[string]string, error) {
	return nil, fmt.Errorf("namespace information is only available on Linux")
}
