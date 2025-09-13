//go:build linux

package namespaces

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

// NamespaceFlags represents the Linux namespace types
type NamespaceFlags int

const (
	// CLONE_NEWPID creates a new PID namespace
	CLONE_NEWPID NamespaceFlags = unix.CLONE_NEWPID
	// CLONE_NEWUTS creates a new UTS namespace (hostname, domainname)
	CLONE_NEWUTS NamespaceFlags = unix.CLONE_NEWUTS
	// CLONE_NEWIPC creates a new IPC namespace
	CLONE_NEWIPC NamespaceFlags = unix.CLONE_NEWIPC
	// CLONE_NEWNET creates a new network namespace
	CLONE_NEWNET NamespaceFlags = unix.CLONE_NEWNET
	// CLONE_NEWNS creates a new mount namespace
	CLONE_NEWNS NamespaceFlags = unix.CLONE_NEWNS
	// CLONE_NEWUSER creates a new user namespace
	CLONE_NEWUSER NamespaceFlags = unix.CLONE_NEWUSER
)

// ContainerConfig holds namespace configuration
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

// CreateContainer creates a new container with the specified namespaces
func CreateContainer(config *ContainerConfig) error {
	// Combine all namespace flags
	var cloneFlags uintptr
	for _, ns := range config.Namespaces {
		cloneFlags |= uintptr(ns)
	}

	// Add user namespace if configured
	if config.UserNamespace != nil && config.UserNamespace.Enabled {
		cloneFlags |= uintptr(CLONE_NEWUSER)
	}

	// Create the container process
	cmd := exec.Command("/proc/self/exe", append([]string{"init"}, config.Command)...)
	cmd.Args = append(cmd.Args, config.Args...)
	cmd.Stdin = os.Stdin

	// Set up log redirection if LogDir is specified
	if config.LogDir != "" {
		if err := setupLogRedirection(cmd, config.LogDir); err != nil {
			return fmt.Errorf("failed to setup log redirection: %v", err)
		}
	} else {
		// Default to direct output
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// Set up the clone flags for namespace creation
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	// Start the process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start container process: %v", err)
	}

	// Setup user namespace if configured
	if config.UserNamespace != nil && config.UserNamespace.Enabled {
		if err := SetupUserNamespace(config.UserNamespace, cmd.Process.Pid); err != nil {
			cmd.Process.Kill()
			return fmt.Errorf("failed to setup user namespace: %v", err)
		}
	}

	fmt.Printf("Creating container with namespaces: %v\n", config.Namespaces)
	if config.UserNamespace != nil && config.UserNamespace.Enabled {
		fmt.Printf("User namespace enabled with UID mappings: %+v\n", config.UserNamespace.UIDMappings)
	}

	// Wait for the process to complete
	return cmd.Wait()
}

// SetupNamespace configures the namespace environment
func SetupNamespace(hostname string) error {
	// Set the container hostname
	if hostname != "" {
		if err := unix.Sethostname([]byte(hostname)); err != nil {
			return fmt.Errorf("failed to set hostname: %v", err)
		}
		fmt.Printf("Set container hostname to: %s\n", hostname)
	}

	// Mount proc filesystem in the new PID namespace
	if err := unix.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		// This might fail if /proc is already mounted, which is okay
		fmt.Printf("Warning: failed to mount /proc: %v\n", err)
	}

	return nil
}

// GetNamespaceInfo returns information about the current namespaces
func GetNamespaceInfo() (map[string]string, error) {
	info := make(map[string]string)

	// Read namespace information from /proc/self/ns/
	namespaces := []string{"pid", "uts", "ipc", "net", "mnt", "user"}

	for _, ns := range namespaces {
		path := fmt.Sprintf("/proc/self/ns/%s", ns)
		if link, err := os.Readlink(path); err == nil {
			info[ns] = link
		}
	}

	return info, nil
}

// setupLogRedirection sets up stdout and stderr redirection to log files
func setupLogRedirection(cmd *exec.Cmd, logDir string) error {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// Create stdout log file
	stdoutPath := filepath.Join(logDir, "stdout.log")
	stdoutFile, err := os.OpenFile(stdoutPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create stdout log file: %v", err)
	}

	// Create stderr log file
	stderrPath := filepath.Join(logDir, "stderr.log")
	stderrFile, err := os.OpenFile(stderrPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		stdoutFile.Close()
		return fmt.Errorf("failed to create stderr log file: %v", err)
	}

	// Create timestamped writers
	cmd.Stdout = &TimestampedWriter{file: stdoutFile}
	cmd.Stderr = &TimestampedWriter{file: stderrFile}

	return nil
}

// TimestampedWriter wraps a file writer to add timestamps to each line
type TimestampedWriter struct {
	file *os.File
}

// Write implements io.Writer interface with timestamping
func (tw *TimestampedWriter) Write(p []byte) (n int, err error) {
	timestamp := time.Now().Format(time.RFC3339Nano)

	// Split input into lines and add timestamp to each
	lines := strings.Split(string(p), "\n")

	for i, line := range lines {
		if i == len(lines)-1 && line == "" {
			// Don't write empty line at the end
			break
		}

		timestampedLine := fmt.Sprintf("%s %s\n", timestamp, line)
		if _, writeErr := tw.file.WriteString(timestampedLine); writeErr != nil {
			return 0, writeErr
		}
	}

	// Flush to ensure data is written
	tw.file.Sync()

	return len(p), nil
}
