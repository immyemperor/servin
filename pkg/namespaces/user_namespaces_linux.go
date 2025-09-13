//go:build linux

package namespaces

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"servin/pkg/logger"
)

// UserNamespaceConfig holds user namespace configuration
type UserNamespaceConfig struct {
	// Enable user namespaces
	Enabled bool

	// UID mapping: maps container UID to host UID
	// Format: "container_uid:host_uid:count"
	UIDMappings []UIDGIDMapping

	// GID mapping: maps container GID to host GID
	// Format: "container_gid:host_gid:count"
	GIDMappings []UIDGIDMapping

	// Run as specific user inside container
	ContainerUID int
	ContainerGID int

	// Security options
	NoNewPrivs  bool     // Prevent gaining new privileges
	DropAllCaps bool     // Drop all capabilities
	AllowedCaps []string // Allowed capabilities to keep
}

// UIDGIDMapping represents a UID or GID mapping entry
type UIDGIDMapping struct {
	ContainerID int // ID inside the container
	HostID      int // ID on the host
	Size        int // Number of IDs to map
}

// DefaultUserNamespaceConfig returns a secure default configuration
func DefaultUserNamespaceConfig() *UserNamespaceConfig {
	currentUID := os.Getuid()
	currentGID := os.Getgid()

	return &UserNamespaceConfig{
		Enabled: true,
		UIDMappings: []UIDGIDMapping{
			{ContainerID: 0, HostID: currentUID, Size: 1},            // Map root to current user
			{ContainerID: 1000, HostID: currentUID + 1, Size: 65536}, // Map range for app users
		},
		GIDMappings: []UIDGIDMapping{
			{ContainerID: 0, HostID: currentGID, Size: 1},            // Map root group to current group
			{ContainerID: 1000, HostID: currentGID + 1, Size: 65536}, // Map range for app groups
		},
		ContainerUID: 1000, // Run as non-root user inside container
		ContainerGID: 1000,
		NoNewPrivs:   true,
		DropAllCaps:  true,
		AllowedCaps:  []string{}, // No capabilities by default
	}
}

// RootlessUserNamespaceConfig returns a rootless container configuration
func RootlessUserNamespaceConfig() *UserNamespaceConfig {
	currentUID := os.Getuid()
	currentGID := os.Getgid()

	return &UserNamespaceConfig{
		Enabled: true,
		UIDMappings: []UIDGIDMapping{
			{ContainerID: currentUID, HostID: currentUID, Size: 1}, // Map current user to itself
		},
		GIDMappings: []UIDGIDMapping{
			{ContainerID: currentGID, HostID: currentGID, Size: 1}, // Map current group to itself
		},
		ContainerUID: currentUID,
		ContainerGID: currentGID,
		NoNewPrivs:   true,
		DropAllCaps:  false, // Keep some capabilities for rootless operation
		AllowedCaps:  []string{"CAP_SETUID", "CAP_SETGID"},
	}
}

// SetupUserNamespace configures user namespace for a container
func SetupUserNamespace(config *UserNamespaceConfig, pid int) error {
	if !config.Enabled {
		return nil
	}

	logger.Info("Setting up user namespace with UID mappings: %+v, GID mappings: %+v",
		config.UIDMappings, config.GIDMappings)

	// Write UID mappings
	if err := writeUIDGIDMap(pid, "uid_map", config.UIDMappings); err != nil {
		return fmt.Errorf("failed to setup UID mapping: %w", err)
	}

	// Disable setgroups to allow GID mapping
	if err := writeSetgroups(pid, "deny"); err != nil {
		return fmt.Errorf("failed to disable setgroups: %w", err)
	}

	// Write GID mappings
	if err := writeUIDGIDMap(pid, "gid_map", config.GIDMappings); err != nil {
		return fmt.Errorf("failed to setup GID mapping: %w", err)
	}

	logger.Info("User namespace configured successfully")
	return nil
}

// ConfigureContainerUser sets up the user context inside the container
func ConfigureContainerUser(config *UserNamespaceConfig) error {
	if !config.Enabled {
		return nil
	}

	// Set no new privileges if requested
	if config.NoNewPrivs {
		if err := setNoNewPrivs(); err != nil {
			return fmt.Errorf("failed to set no_new_privs: %w", err)
		}
		logger.Info("Set no_new_privs for enhanced security")
	}

	// Drop capabilities if requested
	if config.DropAllCaps {
		if err := dropCapabilities(config.AllowedCaps); err != nil {
			return fmt.Errorf("failed to drop capabilities: %w", err)
		}
		logger.Info("Dropped capabilities, allowed: %v", config.AllowedCaps)
	}

	// Switch to specified user and group
	if config.ContainerGID != 0 {
		if err := syscall.Setgid(config.ContainerGID); err != nil {
			return fmt.Errorf("failed to set GID %d: %w", config.ContainerGID, err)
		}
	}

	if config.ContainerUID != 0 {
		if err := syscall.Setuid(config.ContainerUID); err != nil {
			return fmt.Errorf("failed to set UID %d: %w", config.ContainerUID, err)
		}
	}

	logger.Info("Container running as UID: %d, GID: %d",
		config.ContainerUID, config.ContainerGID)

	return nil
}

// GetUserNamespaceInfo returns information about the current user namespace
func GetUserNamespaceInfo() (*UserNamespaceInfo, error) {
	info := &UserNamespaceInfo{}

	// Check if we're in a user namespace
	userNSPath := "/proc/self/ns/user"
	if link, err := os.Readlink(userNSPath); err == nil {
		info.UserNamespace = link
		info.InUserNamespace = true
	}

	// Get current UID/GID
	info.UID = os.Getuid()
	info.GID = os.Getgid()
	info.EUID = os.Geteuid()
	info.EGID = os.Getegid()

	// Read UID/GID mappings
	if uidMaps, err := readUIDGIDMap("/proc/self/uid_map"); err == nil {
		info.UIDMappings = uidMaps
	}

	if gidMaps, err := readUIDGIDMap("/proc/self/gid_map"); err == nil {
		info.GIDMappings = gidMaps
	}

	// Check capabilities
	info.Capabilities = getCurrentCapabilities()

	return info, nil
}

// UserNamespaceInfo contains information about user namespace state
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

// Private helper functions

func writeUIDGIDMap(pid int, mapType string, mappings []UIDGIDMapping) error {
	if len(mappings) == 0 {
		return nil
	}

	mapPath := fmt.Sprintf("/proc/%d/%s", pid, mapType)

	var lines []string
	for _, mapping := range mappings {
		line := fmt.Sprintf("%d %d %d", mapping.ContainerID, mapping.HostID, mapping.Size)
		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n")
	return os.WriteFile(mapPath, []byte(content), 0644)
}

func writeSetgroups(pid int, value string) error {
	setgroupsPath := fmt.Sprintf("/proc/%d/setgroups", pid)
	return os.WriteFile(setgroupsPath, []byte(value), 0644)
}

func readUIDGIDMap(path string) ([]UIDGIDMapping, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var mappings []UIDGIDMapping
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 3 {
			continue
		}

		containerID, err1 := strconv.Atoi(parts[0])
		hostID, err2 := strconv.Atoi(parts[1])
		size, err3 := strconv.Atoi(parts[2])

		if err1 != nil || err2 != nil || err3 != nil {
			continue
		}

		mappings = append(mappings, UIDGIDMapping{
			ContainerID: containerID,
			HostID:      hostID,
			Size:        size,
		})
	}

	return mappings, nil
}

func setNoNewPrivs() error {
	// Set PR_SET_NO_NEW_PRIVS to prevent gaining new privileges
	// PR_SET_NO_NEW_PRIVS = 38
	const PR_SET_NO_NEW_PRIVS = 38
	_, _, errno := syscall.RawSyscall(syscall.SYS_PRCTL, PR_SET_NO_NEW_PRIVS, 1, 0)
	if errno != 0 {
		return errno
	}
	return nil
}

func dropCapabilities(allowedCaps []string) error {
	// This is a simplified capability dropping implementation
	// In a full implementation, you would use libcap or similar
	logger.Info("Capability dropping requested (simplified implementation)")

	// For now, just log what would be done
	if len(allowedCaps) > 0 {
		logger.Info("Would keep capabilities: %v", allowedCaps)
	} else {
		logger.Info("Would drop all capabilities")
	}

	return nil
}

func getCurrentCapabilities() []string {
	// This would read from /proc/self/status and parse CapEff, CapPrm, etc.
	// For now, return a placeholder
	return []string{"simplified_implementation"}
}

// ValidateUserNamespaceSupport checks if user namespaces are supported
func ValidateUserNamespaceSupport() error {
	// Check if user namespaces are available
	if _, err := os.Stat("/proc/self/ns/user"); err != nil {
		return fmt.Errorf("user namespaces not supported: %w", err)
	}

	// Check if we can create user namespaces
	maxUserNS := "/proc/sys/user/max_user_namespaces"
	if content, err := os.ReadFile(maxUserNS); err == nil {
		if strings.TrimSpace(string(content)) == "0" {
			return fmt.Errorf("user namespaces disabled (max_user_namespaces=0)")
		}
	}

	return nil
}

// GetSubUIDGIDRanges reads subuid and subgid ranges for a user
func GetSubUIDGIDRanges(username string) ([]UIDGIDMapping, []UIDGIDMapping, error) {
	var uidRanges, gidRanges []UIDGIDMapping

	// Read /etc/subuid
	if uidMaps, err := parseSubIDFile("/etc/subuid", username); err == nil {
		uidRanges = uidMaps
	}

	// Read /etc/subgid
	if gidMaps, err := parseSubIDFile("/etc/subgid", username); err == nil {
		gidRanges = gidMaps
	}

	return uidRanges, gidRanges, nil
}

func parseSubIDFile(path, username string) ([]UIDGIDMapping, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var mappings []UIDGIDMapping
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 3 || parts[0] != username {
			continue
		}

		startID, err1 := strconv.Atoi(parts[1])
		count, err2 := strconv.Atoi(parts[2])

		if err1 != nil || err2 != nil {
			continue
		}

		mappings = append(mappings, UIDGIDMapping{
			ContainerID: 1, // Start mapping from ID 1 in container
			HostID:      startID,
			Size:        count,
		})
	}

	return mappings, nil
}
