//go:build !linux

package namespaces

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
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

	// User namespace configuration
	UserNamespace *UserNamespaceConfig
}

// CreateContainer runs the command directly on non-Linux platforms (simulated containerization)
func CreateContainer(config *ContainerConfig) error {
	fmt.Printf("Running command directly (no containerization): %s %v\n", config.Command, config.Args)

	// Create a simple command execution without namespaces
	cmd := exec.Command(config.Command, config.Args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	// Wait for the process to complete or signal
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		// Process completed normally
		return err
	case sig := <-sigChan:
		// Received signal, terminate the process
		fmt.Printf("\nReceived signal %v, terminating process...\n", sig)

		// First try graceful termination
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			fmt.Printf("Failed to send SIGTERM: %v\n", err)
		}

		// Wait a bit for graceful shutdown
		timeout := time.After(2 * time.Second)
		select {
		case err := <-done:
			fmt.Println("Process terminated gracefully")
			return err
		case <-timeout:
			// Force kill if it doesn't respond
			fmt.Println("Process didn't respond to SIGTERM, force killing...")
			if err := cmd.Process.Kill(); err != nil {
				fmt.Printf("Failed to kill process: %v\n", err)
			}
			<-done // Wait for the process to actually exit
			return fmt.Errorf("process forcefully terminated")
		}
	}
}

// SetupNamespace returns an error on non-Linux platforms
func SetupNamespace(hostname string) error {
	return fmt.Errorf("namespace setup is only supported on Linux")
}

// GetNamespaceInfo returns an error on non-Linux platforms
func GetNamespaceInfo() (map[string]string, error) {
	return nil, fmt.Errorf("namespace information is only available on Linux")
}
