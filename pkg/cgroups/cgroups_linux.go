//go:build linux

package cgroups

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// CGroup manages container resource limits
type CGroup struct {
	ContainerID string
	Path        string
}

// New creates a new CGroup manager
func New(containerID string) *CGroup {
	// Use cgroup v1 path structure
	basePath := "/sys/fs/cgroup"
	return &CGroup{
		ContainerID: containerID,
		Path:        filepath.Join(basePath, "servin", containerID),
	}
}

// Create sets up cgroup directories and files
func (c *CGroup) Create() error {
	// Create cgroup directories for different subsystems
	subsystems := []string{"memory", "cpu", "pids"}

	for _, subsystem := range subsystems {
		subsystemPath := filepath.Join("/sys/fs/cgroup", subsystem, "servin", c.ContainerID)
		if err := os.MkdirAll(subsystemPath, 0755); err != nil {
			return fmt.Errorf("failed to create cgroup directory %s: %v", subsystemPath, err)
		}
		fmt.Printf("Created cgroup: %s\n", subsystemPath)
	}

	return nil
}

// SetMemoryLimit sets the memory limit for the container
func (c *CGroup) SetMemoryLimit(limitBytes int64) error {
	memoryPath := filepath.Join("/sys/fs/cgroup", "memory", "servin", c.ContainerID, "memory.limit_in_bytes")
	return writeToFile(memoryPath, strconv.FormatInt(limitBytes, 10))
}

// SetCPULimit sets the CPU limit for the container (in CPU shares)
func (c *CGroup) SetCPULimit(shares int) error {
	cpuPath := filepath.Join("/sys/fs/cgroup", "cpu", "servin", c.ContainerID, "cpu.shares")
	return writeToFile(cpuPath, strconv.Itoa(shares))
}

// SetPIDLimit sets the maximum number of processes
func (c *CGroup) SetPIDLimit(max int) error {
	pidsPath := filepath.Join("/sys/fs/cgroup", "pids", "servin", c.ContainerID, "pids.max")
	return writeToFile(pidsPath, strconv.Itoa(max))
}

// AddProcess adds a process to the cgroup
func (c *CGroup) AddProcess(pid int) error {
	subsystems := []string{"memory", "cpu", "pids"}
	pidStr := strconv.Itoa(pid)

	for _, subsystem := range subsystems {
		tasksPath := filepath.Join("/sys/fs/cgroup", subsystem, "servin", c.ContainerID, "tasks")
		if err := writeToFile(tasksPath, pidStr); err != nil {
			return fmt.Errorf("failed to add process %d to %s cgroup: %v", pid, subsystem, err)
		}
	}

	fmt.Printf("Added process %d to cgroups\n", pid)
	return nil
}

// GetStats returns resource usage statistics
func (c *CGroup) GetStats() (map[string]string, error) {
	stats := make(map[string]string)

	// Memory usage
	memUsagePath := filepath.Join("/sys/fs/cgroup", "memory", "servin", c.ContainerID, "memory.usage_in_bytes")
	if usage, err := readFromFile(memUsagePath); err == nil {
		stats["memory_usage"] = strings.TrimSpace(usage)
	}

	// CPU usage
	cpuUsagePath := filepath.Join("/sys/fs/cgroup", "cpu", "servin", c.ContainerID, "cpuacct.usage")
	if usage, err := readFromFile(cpuUsagePath); err == nil {
		stats["cpu_usage"] = strings.TrimSpace(usage)
	}

	// PID count
	pidsCurrentPath := filepath.Join("/sys/fs/cgroup", "pids", "servin", c.ContainerID, "pids.current")
	if current, err := readFromFile(pidsCurrentPath); err == nil {
		stats["pids_current"] = strings.TrimSpace(current)
	}

	return stats, nil
}

// Cleanup removes the cgroup directories
func (c *CGroup) Cleanup() error {
	subsystems := []string{"memory", "cpu", "pids"}

	for _, subsystem := range subsystems {
		subsystemPath := filepath.Join("/sys/fs/cgroup", subsystem, "servin", c.ContainerID)
		if err := os.RemoveAll(subsystemPath); err != nil {
			fmt.Printf("Warning: failed to remove cgroup %s: %v\n", subsystemPath, err)
		}
	}

	return nil
}

// ParseMemoryString converts memory strings like "128m", "1g" to bytes
func ParseMemoryString(memStr string) (int64, error) {
	if memStr == "" {
		return 0, nil
	}

	memStr = strings.ToLower(strings.TrimSpace(memStr))

	var multiplier int64 = 1
	var numStr string

	if strings.HasSuffix(memStr, "k") {
		multiplier = 1024
		numStr = strings.TrimSuffix(memStr, "k")
	} else if strings.HasSuffix(memStr, "m") {
		multiplier = 1024 * 1024
		numStr = strings.TrimSuffix(memStr, "m")
	} else if strings.HasSuffix(memStr, "g") {
		multiplier = 1024 * 1024 * 1024
		numStr = strings.TrimSuffix(memStr, "g")
	} else {
		numStr = memStr
	}

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid memory format: %s", memStr)
	}

	return int64(num * float64(multiplier)), nil
}

// Helper functions
func writeToFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func readFromFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
