package cmd

import (
	"fmt"

	"servin/pkg/state"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "rm [OPTIONS] CONTAINER [CONTAINER...]",
	Aliases: []string{"remove"},
	Short:   "Remove one or more containers",
	Long: `Remove one or more containers. By default, running containers cannot be removed.
Use the --force flag to stop and remove running containers.`,
	Args: func(cmd *cobra.Command, args []string) error {
		// If --all flag is used, we don't need container arguments
		if removeAll {
			return nil
		}
		// Otherwise, require at least one container argument
		return cobra.MinimumNArgs(1)(cmd, args)
	},
	RunE: removeContainers,
}

var (
	forceRemove bool
	removeAll   bool
)

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().BoolVarP(&forceRemove, "force", "f", false, "Force removal of running containers")
	removeCmd.Flags().BoolVarP(&removeAll, "all", "a", false, "Remove all stopped containers")
}

func removeContainers(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	// Create state manager
	sm := state.NewStateManager()

	var containersToRemove []string

	// If --all flag is used, get all stopped containers
	if removeAll {
		containers, err := sm.ListContainers()
		if err != nil {
			return fmt.Errorf("failed to list containers: %v", err)
		}

		for _, container := range containers {
			if container.Status != "running" {
				containersToRemove = append(containersToRemove, container.ID)
			}
		}

		if len(containersToRemove) == 0 {
			fmt.Println("No stopped containers to remove")
			return nil
		}
	} else {
		// Use provided container references
		for _, containerRef := range args {
			containerID, err := resolveContainerRef(sm, containerRef)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			containersToRemove = append(containersToRemove, containerID)
		}
	}

	// Remove each container
	var removedCount int
	for _, containerID := range containersToRemove {
		if err := removeContainer(sm, containerID, forceRemove); err != nil {
			fmt.Printf("Error removing container %s: %v\n", containerID[:12], err)
		} else {
			removedCount++
		}
	}

	if removeAll {
		fmt.Printf("Removed %d stopped containers\n", removedCount)
	}

	return nil
}

func removeContainer(sm *state.StateManager, containerID string, force bool) error {
	// Load container state
	container, err := sm.LoadContainer(containerID)
	if err != nil {
		return fmt.Errorf("container not found: %v", err)
	}

	// Check if container is running
	if container.Status == "running" {
		if !force {
			return fmt.Errorf("cannot remove running container %s. Stop the container before removing or use --force", container.Name)
		}

		// Force stop the container first
		fmt.Printf("Stopping running container %s...\n", container.Name)
		if container.PID > 0 {
			if err := stopContainerProcess(container.PID); err != nil {
				fmt.Printf("Warning: failed to stop container process: %v\n", err)
			}
		}

		// Update status to stopped
		if err := sm.UpdateContainerStatus(containerID, "stopped"); err != nil {
			fmt.Printf("Warning: failed to update container status: %v\n", err)
		}
	}

	// Remove container resources
	if err := cleanupContainerResources(container); err != nil {
		fmt.Printf("Warning: failed to cleanup container resources: %v\n", err)
	}

	// Remove container state file
	if err := sm.DeleteContainer(containerID); err != nil {
		return fmt.Errorf("failed to remove container state: %v", err)
	}

	fmt.Printf("Removed container %s (%s)\n", container.Name, containerID[:12])
	return nil
}

// cleanupContainerResources removes container-specific resources
func cleanupContainerResources(container *state.ContainerState) error {
	// This function would clean up:
	// 1. Container rootfs directory
	// 2. Container network interfaces
	// 3. Container cgroups
	// 4. Container volumes

	// For now, we'll just log what would be cleaned up
	fmt.Printf("  Cleaning up resources for container %s\n", container.Name)

	if container.RootPath != "" {
		fmt.Printf("  - Would remove rootfs: %s\n", container.RootPath)
	}

	if len(container.Volumes) > 0 {
		fmt.Printf("  - Would unmount %d volumes\n", len(container.Volumes))
	}

	if container.NetworkMode == "bridge" {
		fmt.Printf("  - Would cleanup network interfaces\n")
	}

	// In a full implementation, you would:
	// - Remove the container's rootfs directory
	// - Remove network interfaces (veth pairs)
	// - Remove cgroups
	// - Unmount volumes
	// - Remove any temporary files

	return nil
}
