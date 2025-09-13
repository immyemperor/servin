package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec [FLAGS] CONTAINER COMMAND [ARG...]",
	Short: "Execute a command in a running container",
	Args:  cobra.MinimumNArgs(2),
	RunE:  execInContainer,
}

func init() {
	rootCmd.AddCommand(execCmd)

	// Add flags for exec command
	execCmd.Flags().BoolP("interactive", "i", false, "Keep STDIN open even if not attached")
	execCmd.Flags().BoolP("tty", "t", false, "Allocate a pseudo-TTY")
}

func execInContainer(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	containerID := args[0]
	command := args[1]
	commandArgs := args[2:]

	fmt.Printf("Executing '%s %v' in container %s...\n", command, commandArgs, containerID)

	// TODO: Implement exec functionality
	fmt.Println("Exec functionality not yet implemented")

	return nil
}
