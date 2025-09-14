//go:build linux

package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
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

	// Set up the container environment using namespaces
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
