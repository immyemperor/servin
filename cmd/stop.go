package cmd

import (
	"fmt"

	"servin/pkg/state"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop CONTAINER [CONTAINER...]",
	Short: "Stop one or more running containers",
	Args:  cobra.MinimumNArgs(1),
	RunE:  stopContainers,
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func stopContainers(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	// Create state manager
	sm := state.NewStateManager()

	for _, containerRef := range args {
		fmt.Printf("Stopping container %s...\n", containerRef)

		// Find the container (could be ID, short ID, or name)
		containerID, err := resolveContainerRef(sm, containerRef)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// Load container state
		container, err := sm.LoadContainer(containerID)
		if err != nil {
			fmt.Printf("Error: failed to load container %s: %v\n", containerRef, err)
			continue
		}

		// Check if container is running
		if container.Status != "running" {
			fmt.Printf("Container %s is not running (status: %s)\n", containerRef, container.Status)
			continue
		}

		// Stop the container process
		if container.PID > 0 {
			if err := stopContainerProcess(container.PID); err != nil {
				fmt.Printf("Error stopping container %s: %v\n", containerRef, err)
				continue
			}
		}

		// Update container status
		if err := sm.UpdateContainerStatus(containerID, "stopped"); err != nil {
			fmt.Printf("Warning: failed to update container status: %v\n", err)
		}

		fmt.Printf("Container %s stopped\n", containerRef)
	}

	return nil
}
