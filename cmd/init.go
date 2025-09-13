//go:build linux

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"servin/pkg/namespaces"

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

func setupContainerEnvironment() error {
	// Basic setup that works on all platforms
	if runtime.GOOS != "linux" {
		return fmt.Errorf("this containerization tool only works on Linux")
	}

	// Set up namespaces (this will work on Linux)
	if err := namespaces.SetupNamespace("container"); err != nil {
		fmt.Printf("Warning: namespace setup failed: %v\n", err)
	}

	// Change to root directory
	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("failed to change to root directory: %v", err)
	}

	// Print namespace information for debugging
	if nsInfo, err := namespaces.GetNamespaceInfo(); err == nil {
		fmt.Println("Container namespace info:")
		for ns, info := range nsInfo {
			fmt.Printf("  %s: %s\n", ns, info)
		}
	}

	fmt.Println("Container environment initialized")
	return nil
}
