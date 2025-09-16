//go:build linux

package cmd

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// setupBasicContainerEnvironment provides basic container setup functionality
// This is a simpler alternative to the full setupContainerEnvironment in init_linux.go
func setupBasicContainerEnvironment() error {
	// Set hostname
	hostname := "container"
	if err := unix.Sethostname([]byte(hostname)); err != nil {
		return fmt.Errorf("failed to set hostname: %v", err)
	}

	// Mount proc filesystem
	if err := unix.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		// Non-fatal if proc is already mounted
		fmt.Printf("Warning: failed to mount /proc: %v\n", err)
	}

	// Change to root directory
	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("failed to change to root directory: %v", err)
	}

	return nil
}
