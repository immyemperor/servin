package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"servin/pkg/state"

	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect CONTAINER",
	Short: "Display detailed container information",
	Long:  "Display comprehensive information about a container including state, config, and runtime details",
	Args:  cobra.ExactArgs(1),
	RunE:  inspectContainer,
}

var psCmd = &cobra.Command{
	Use:   "ps CONTAINER",
	Short: "List processes running in container",
	Long:  "Display all processes currently running inside the specified container",
	Args:  cobra.ExactArgs(1),
	RunE:  listContainerProcesses,
}

var topCmd = &cobra.Command{
	Use:   "top CONTAINER",
	Short: "Display running processes in container",
	Long:  "Display a live view of processes running in the container (similar to htop)",
	Args:  cobra.ExactArgs(1),
	RunE:  showContainerTop,
}

var statsCmd = &cobra.Command{
	Use:   "stats [CONTAINER...]",
	Short: "Display container resource usage statistics",
	Long:  "Display live resource usage statistics for one or more containers",
	RunE:  showContainerStats,
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	rootCmd.AddCommand(psCmd)
	rootCmd.AddCommand(topCmd)
	rootCmd.AddCommand(statsCmd)

	// Add flags
	inspectCmd.Flags().BoolP("format", "f", false, "Format output as JSON")
	statsCmd.Flags().BoolP("no-stream", "n", false, "Disable streaming stats and only pull the first result")
	statsCmd.Flags().IntP("interval", "i", 1, "Refresh interval in seconds")
}

func inspectContainer(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	containerID := args[0]
	sm := state.NewStateManager()

	// Load container state
	container, err := sm.LoadContainer(containerID)
	if err != nil {
		return fmt.Errorf("container not found: %s", containerID)
	}

	format, _ := cmd.Flags().GetBool("format")

	if format {
		// JSON format output
		fmt.Printf(`{
  "Id": "%s",
  "Name": "%s",
  "Image": "%s",
  "Command": "%s",
  "Args": %q,
  "Status": "%s",
  "Created": "%s",
  "Started": "%s",
  "PID": %d,
  "RootFS": "%s",
  "NetworkMode": "%s"
}`, container.ID, container.Name, container.Image, container.Command,
			container.Args, container.Status, container.Created.Format(time.RFC3339),
			container.Started.Format(time.RFC3339), container.PID,
			getContainerRootFSPath(container.ID), container.NetworkMode)
	} else {
		// Human readable format
		fmt.Printf("Container ID: %s\n", container.ID)
		fmt.Printf("Name: %s\n", container.Name)
		fmt.Printf("Image: %s\n", container.Image)
		fmt.Printf("Command: %s %s\n", container.Command, strings.Join(container.Args, " "))
		fmt.Printf("Status: %s\n", container.Status)
		fmt.Printf("Created: %s\n", container.Created.Format(time.RFC3339))
		fmt.Printf("Started: %s\n", container.Started.Format(time.RFC3339))
		fmt.Printf("PID: %d\n", container.PID)
		fmt.Printf("Network Mode: %s\n", container.NetworkMode)

		// Show rootfs information
		rootfsPath := getContainerRootFSPath(container.ID)
		if stat, err := os.Stat(rootfsPath); err == nil {
			fmt.Printf("RootFS: %s (size: %d bytes)\n", rootfsPath, stat.Size())
		} else {
			fmt.Printf("RootFS: %s (not accessible)\n", rootfsPath)
		}

		// Show resource usage if available
		if container.PID > 0 {
			showProcessInfo(container.PID)
		}
	}

	return nil
}

func listContainerProcesses(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	containerID := args[0]
	sm := state.NewStateManager()

	// Load container state
	container, err := sm.LoadContainer(containerID)
	if err != nil {
		return fmt.Errorf("container not found: %s", containerID)
	}

	if container.Status != state.StatusRunning {
		return fmt.Errorf("container is not running")
	}

	// Use ps to find processes in the container's namespace
	// This is a simplified approach - in a full implementation you'd check the actual namespace
	fmt.Printf("%-8s %-8s %-8s %-20s %s\n", "PID", "PPID", "USER", "TIME", "COMMAND")

	if container.PID > 0 {
		// Show the main container process
		if processExists(container.PID) {
			cmd := exec.Command("ps", "-p", fmt.Sprintf("%d", container.PID), "-o", "pid,ppid,user,time,cmd", "--no-headers")
			output, err := cmd.Output()
			if err == nil {
				fmt.Print(string(output))
			}
		}

		// Try to find child processes
		cmd := exec.Command("pgrep", "-P", fmt.Sprintf("%d", container.PID))
		output, err := cmd.Output()
		if err == nil {
			pids := strings.Fields(string(output))
			for _, pid := range pids {
				cmd := exec.Command("ps", "-p", pid, "-o", "pid,ppid,user,time,cmd", "--no-headers")
				if childOutput, err := cmd.Output(); err == nil {
					fmt.Print(string(childOutput))
				}
			}
		}
	} else {
		fmt.Println("No processes found (container may not be running)")
	}

	return nil
}

func showContainerTop(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	containerID := args[0]
	sm := state.NewStateManager()

	// Load container state
	container, err := sm.LoadContainer(containerID)
	if err != nil {
		return fmt.Errorf("container not found: %s", containerID)
	}

	if container.Status != state.StatusRunning {
		return fmt.Errorf("container is not running")
	}

	// Simple top-like display
	fmt.Printf("Processes in container %s (PID: %d)\n", containerID[:12], container.PID)
	fmt.Printf("%-8s %-8s %-8s %-8s %-20s %s\n", "PID", "USER", "CPU%", "MEM%", "TIME", "COMMAND")

	if container.PID > 0 && processExists(container.PID) {
		// Show detailed process information
		cmd := exec.Command("ps", "-p", fmt.Sprintf("%d", container.PID), "-o", "pid,user,pcpu,pmem,time,cmd", "--no-headers")
		output, err := cmd.Output()
		if err == nil {
			fmt.Print(string(output))
		}
	} else {
		fmt.Println("No processes found")
	}

	return nil
}

func showContainerStats(cmd *cobra.Command, args []string) error {
	if err := checkRoot(); err != nil {
		return err
	}

	sm := state.NewStateManager()
	var containers []*state.ContainerState

	if len(args) == 0 {
		// Show stats for all running containers
		allContainers, err := sm.ListContainers()
		if err != nil {
			return err
		}
		for _, c := range allContainers {
			if c.Status == state.StatusRunning {
				containers = append(containers, c)
			}
		}
	} else {
		// Show stats for specified containers
		for _, containerID := range args {
			container, err := sm.LoadContainer(containerID)
			if err != nil {
				fmt.Printf("Warning: container %s not found\n", containerID)
				continue
			}
			containers = append(containers, container)
		}
	}

	if len(containers) == 0 {
		fmt.Println("No running containers found")
		return nil
	}

	noStream, _ := cmd.Flags().GetBool("no-stream")
	interval, _ := cmd.Flags().GetInt("interval")

	// Print header
	fmt.Printf("%-12s %-20s %-8s %-8s %-12s %-12s\n",
		"CONTAINER", "NAME", "CPU %", "MEM %", "MEM USAGE", "NET I/O")

	if noStream {
		// Single snapshot
		for _, container := range containers {
			showContainerResourceUsage(container)
		}
	} else {
		// Continuous monitoring
		for {
			fmt.Print("\033[2J\033[H") // Clear screen and move cursor to top
			fmt.Printf("%-12s %-20s %-8s %-8s %-12s %-12s\n",
				"CONTAINER", "NAME", "CPU %", "MEM %", "MEM USAGE", "NET I/O")

			for _, container := range containers {
				showContainerResourceUsage(container)
			}

			time.Sleep(time.Duration(interval) * time.Second)
		}
	}

	return nil
}

// Helper functions

func getContainerRootFSPath(containerID string) string {
	possiblePaths := []string{
		fmt.Sprintf("/var/lib/servin/containers/%s/rootfs", containerID),
		fmt.Sprintf("/tmp/servin/containers/%s/rootfs", containerID),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return fmt.Sprintf("/var/lib/servin/containers/%s/rootfs", containerID)
}

func processExists(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func showProcessInfo(pid int) {
	if !processExists(pid) {
		fmt.Printf("Process Status: Not running\n")
		return
	}

	// Read process stat file
	statFile := fmt.Sprintf("/proc/%d/stat", pid)
	if _, err := os.Stat(statFile); err == nil {
		fmt.Printf("Process Status: Running (PID: %d)\n", pid)

		// Try to get memory usage
		statusFile := fmt.Sprintf("/proc/%d/status", pid)
		if content, err := os.ReadFile(statusFile); err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "VmRSS:") {
					fmt.Printf("Memory Usage: %s\n", strings.TrimSpace(strings.TrimPrefix(line, "VmRSS:")))
					break
				}
			}
		}
	} else {
		fmt.Printf("Process Status: Running (limited info available)\n")
	}
}

func showContainerResourceUsage(container *state.ContainerState) {
	containerShort := container.ID
	if len(containerShort) > 12 {
		containerShort = containerShort[:12]
	}

	name := container.Name
	if len(name) > 20 {
		name = name[:17] + "..."
	}

	if container.PID <= 0 || !processExists(container.PID) {
		fmt.Printf("%-12s %-20s %-8s %-8s %-12s %-12s\n",
			containerShort, name, "0.00%", "0.00%", "0B / 0B", "0B / 0B")
		return
	}

	// Get CPU and memory usage
	cpuPercent := "0.00%"
	memPercent := "0.00%"
	memUsage := "0B / 0B"
	netIO := "0B / 0B"

	// Try to get actual stats from /proc
	statusFile := fmt.Sprintf("/proc/%d/status", container.PID)
	if content, err := os.ReadFile(statusFile); err == nil {
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "VmRSS:") {
				// Extract memory usage
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					memUsage = fmt.Sprintf("%s kB / -", fields[1])
				}
				break
			}
		}
	}

	fmt.Printf("%-12s %-20s %-8s %-8s %-12s %-12s\n",
		containerShort, name, cpuPercent, memPercent, memUsage, netIO)
}
