//go:build linux

package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

var initCmd = &cobra.Command{
	Use:    "init",
	Short:  "Initialize container environment (internal command)",
	Hidden: true, // Hide from help as this is an internal command
	RunE:   initContainer,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initContainer(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("init command requires at least one argument")
	}

	fmt.Printf("Initializing container as PID %d\n", os.Getpid())

	// Set up the container environment
	if err := setupContainerEnvironment(); err != nil {
		return fmt.Errorf("failed to setup container environment: %v", err)
	}

	// Execute the target command
	command := args[0]
	commandArgs := args[1:]

	execCmd := exec.Command(command, commandArgs...)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	return execCmd.Run()
}

func setupContainerEnvironment() error {
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
