//go:build !linux

package namespaces

import "fmt"

// UserNamespaceConfig holds user namespace configuration (stub for non-Linux)
type UserNamespaceConfig struct {
	Enabled      bool
	UIDMappings  []UIDGIDMapping
	GIDMappings  []UIDGIDMapping
	ContainerUID int
	ContainerGID int
	NoNewPrivs   bool
	DropAllCaps  bool
	AllowedCaps  []string
}

// UIDGIDMapping represents a UID or GID mapping entry (stub for non-Linux)
type UIDGIDMapping struct {
	ContainerID int
	HostID      int
	Size        int
}

// UserNamespaceInfo contains information about user namespace state (stub for non-Linux)
type UserNamespaceInfo struct {
	InUserNamespace bool
	UserNamespace   string
	UID             int
	GID             int
	EUID            int
	EGID            int
	UIDMappings     []UIDGIDMapping
	GIDMappings     []UIDGIDMapping
	Capabilities    []string
}

// DefaultUserNamespaceConfig returns a disabled configuration for non-Linux
func DefaultUserNamespaceConfig() *UserNamespaceConfig {
	return &UserNamespaceConfig{
		Enabled: false,
	}
}

// RootlessUserNamespaceConfig returns a disabled configuration for non-Linux
func RootlessUserNamespaceConfig() *UserNamespaceConfig {
	return &UserNamespaceConfig{
		Enabled: false,
	}
}

// SetupUserNamespace is a no-op on non-Linux platforms
func SetupUserNamespace(config *UserNamespaceConfig, pid int) error {
	if config.Enabled {
		return fmt.Errorf("user namespaces not supported on this platform")
	}
	return nil
}

// ConfigureContainerUser is a no-op on non-Linux platforms
func ConfigureContainerUser(config *UserNamespaceConfig) error {
	if config.Enabled {
		return fmt.Errorf("user namespaces not supported on this platform")
	}
	return nil
}

// GetUserNamespaceInfo returns basic info on non-Linux platforms
func GetUserNamespaceInfo() (*UserNamespaceInfo, error) {
	return &UserNamespaceInfo{
		InUserNamespace: false,
		UserNamespace:   "not_supported",
		UID:             0,
		GID:             0,
		EUID:            0,
		EGID:            0,
		UIDMappings:     []UIDGIDMapping{},
		GIDMappings:     []UIDGIDMapping{},
		Capabilities:    []string{"not_supported"},
	}, nil
}

// ValidateUserNamespaceSupport always returns an error on non-Linux
func ValidateUserNamespaceSupport() error {
	return fmt.Errorf("user namespaces not supported on this platform")
}

// GetSubUIDGIDRanges returns empty ranges on non-Linux platforms
func GetSubUIDGIDRanges(username string) ([]UIDGIDMapping, []UIDGIDMapping, error) {
	return []UIDGIDMapping{}, []UIDGIDMapping{}, fmt.Errorf("subuid/subgid not supported on this platform")
}
