//go:build !linux

package namespaces

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

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
	Command     string
	Args        []string
	Namespaces  []NamespaceFlags
	Hostname    string
	WorkDir     string
	LogDir      string            // Directory to store container logs
	RootFS      string            // RootFS path for the container
	Environment map[string]string // Environment variables
	OnExit      func(error)       // Callback when process exits

	// User namespace configuration
	UserNamespace *UserNamespaceConfig
}

// CreateContainer runs the command directly on non-Linux platforms (simulated containerization)
func CreateContainer(config *ContainerConfig) error {
	fmt.Printf("Running command directly (no containerization): %s %v\n", config.Command, config.Args)

	// Create a simple command execution without namespaces
	cmd := exec.Command(config.Command, config.Args...)

	// Set working directory if specified
	if config.WorkDir != "" {
		cmd.Dir = config.WorkDir
	}

	// Set environment variables
	cmd.Env = os.Environ()
	if config.Hostname != "" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("HOSTNAME=%s", config.Hostname))
	}
	for key, value := range config.Environment {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	// Set up logging if log directory is specified
	if config.LogDir != "" {
		// Create log directory if it doesn't exist
		if err := os.MkdirAll(config.LogDir, 0755); err != nil {
			fmt.Printf("Warning: failed to create log directory: %v\n", err)
		} else {
			// Create log files for stdout and stderr
			stdoutFile, err := os.Create(filepath.Join(config.LogDir, "stdout.log"))
			if err != nil {
				fmt.Printf("Warning: failed to create stdout log: %v\n", err)
			} else {
				cmd.Stdout = stdoutFile
			}

			stderrFile, err := os.Create(filepath.Join(config.LogDir, "stderr.log"))
			if err != nil {
				fmt.Printf("Warning: failed to create stderr log: %v\n", err)
			} else {
				cmd.Stderr = stderrFile
			}
		}
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}

	// Return immediately - don't wait for the command to finish
	// This allows the container to be in "running" state while the command executes
	fmt.Printf("Command started with PID %d, returning to allow container to run\n", cmd.Process.Pid)

	// Start a goroutine to wait for the process and handle cleanup
	go func() {
		err := cmd.Wait()
		if err != nil {
			fmt.Printf("Container command exited with error: %v\n", err)
		} else {
			fmt.Printf("Container command completed successfully\n")
		}

		// Call the exit callback if provided
		if config.OnExit != nil {
			config.OnExit(err)
		}
	}()

	return nil
}

// SetupNamespace returns an error on non-Linux platforms
func SetupNamespace(hostname string) error {
	return fmt.Errorf("namespace setup is only supported on Linux")
}

// GetNamespaceInfo returns an error on non-Linux platforms
func GetNamespaceInfo() (map[string]string, error) {
	return nil, fmt.Errorf("namespace information is only available on Linux")
}
