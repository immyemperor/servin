package cmd

import (
	"fmt"
	"time"

	"servin/pkg/state"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "ps"},
	Short:   "List containers",
	Long:    "List all containers (running and stopped)",
	RunE:    listContainers,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("detailed", "d", false, "Show detailed container information including port mappings")
}

func listContainers(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	// Create state manager
	sm := state.NewStateManager()

	// Load all containers
	containers, err := sm.ListContainers()
	if err != nil {
		return fmt.Errorf("failed to list containers: %v", err)
	}

	if len(containers) == 0 {
		fmt.Println("CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS    NAMES")
		fmt.Println("(No containers found)")
		return nil
	}

	// Print header
	fmt.Printf("%-12s %-15s %-20s %-15s %-10s %s\n",
		"CONTAINER ID", "IMAGE", "COMMAND", "CREATED", "STATUS", "NAMES")

	// Print each container
	detailed, _ := cmd.Flags().GetBool("detailed")

	for _, container := range containers {
		shortID := container.ID[:12]
		image := truncateString(container.Image, 15)
		command := truncateString(container.Command, 20)
		created := formatTime(container.Created)
		status := container.Status
		name := container.Name

		fmt.Printf("%-12s %-15s %-20s %-15s %-10s %s\n",
			shortID, image, command, created, status, name)

		// Show detailed information if requested
		if detailed {
			if len(container.PortMappings) > 0 {
				fmt.Printf("  Ports: ")
				for i, port := range container.PortMappings {
					if i > 0 {
						fmt.Printf(", ")
					}
					if port.HostIP != "" {
						fmt.Printf("%s:", port.HostIP)
					}
					fmt.Printf("%d->%d/%s", port.HostPort, port.ContainerPort, port.Protocol)
				}
				fmt.Printf("\n")
			}
			if container.NetworkMode != "" && container.NetworkMode != "bridge" {
				fmt.Printf("  Network: %s\n", container.NetworkMode)
			}
			fmt.Printf("\n")
		}
	}

	// Show state directory info
	fmt.Printf("\nState directory: %s\n", sm.GetStateDir())

	return nil
}

// truncateString truncates a string to the specified length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// formatTime formats a time for display
func formatTime(t time.Time) string {
	if t.IsZero() {
		return "Unknown"
	}

	now := time.Now()
	duration := now.Sub(t)

	switch {
	case duration < time.Minute:
		return fmt.Sprintf("%d seconds ago", int(duration.Seconds()))
	case duration < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	case duration < 7*24*time.Hour:
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	default:
		return t.Format("2006-01-02")
	}
}
