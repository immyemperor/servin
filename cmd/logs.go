package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"servin/pkg/errors"
	"servin/pkg/logger"
	"servin/pkg/state"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs [OPTIONS] CONTAINER",
	Short: "Fetch the logs of a container",
	Long: `Fetch and display the logs of a running or stopped container.
The logs command retrieves stdout and stderr output from the container.`,
	Args: cobra.ExactArgs(1),
	RunE: showContainerLogs,
}

var (
	follow     bool
	timestamps bool
	tail       string
	since      string
	until      string
)

func init() {
	rootCmd.AddCommand(logsCmd)

	// Add flags for log options
	logsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output")
	logsCmd.Flags().BoolVarP(&timestamps, "timestamps", "t", false, "Show timestamps")
	logsCmd.Flags().StringVar(&tail, "tail", "all", "Number of lines to show from the end of the logs")
	logsCmd.Flags().StringVar(&since, "since", "", "Show logs since timestamp (e.g. 2013-01-02T13:23:37Z) or relative (e.g. 42m for 42 minutes)")
	logsCmd.Flags().StringVar(&until, "until", "", "Show logs before a timestamp (e.g. 2013-01-02T13:23:37Z) or relative (e.g. 42m for 42 minutes)")
}

func showContainerLogs(cmd *cobra.Command, args []string) error {
	containerIDOrName := args[0]

	logger.Debug("Showing logs for container: %s", containerIDOrName)

	// Create state manager
	sm := state.NewStateManager()

	// Find container by ID or name
	container, err := sm.LoadContainer(containerIDOrName)
	if err != nil {
		// Try to find by name if ID lookup failed
		containers, listErr := sm.ListContainers()
		if listErr != nil {
			logger.Error("Failed to list containers: %v", listErr)
			return errors.NewNotFoundError("container", "container not found and unable to search by name")
		}

		var found *state.ContainerState
		for _, c := range containers {
			if c.Name == containerIDOrName {
				found = c
				break
			}
		}

		if found == nil {
			logger.Error("Container not found: %s", containerIDOrName)
			return errors.NewNotFoundError("container", "no container found with this ID or name")
		}
		container = found
	}

	logger.Debug("Found container: %s (status: %s)", container.ID, container.Status)

	// Get log file paths
	logDir := getContainerLogDir(container.ID)
	stdoutPath := filepath.Join(logDir, "stdout.log")
	stderrPath := filepath.Join(logDir, "stderr.log")

	logger.Debug("Looking for log files in: %s", logDir)

	// Check if log files exist
	if _, err := os.Stat(stdoutPath); os.IsNotExist(err) {
		if _, err := os.Stat(stderrPath); os.IsNotExist(err) {
			logger.Warn("No log files found for container: %s", container.ID)
			fmt.Printf("No logs available for container %s\n", containerIDOrName)
			return nil
		}
	}

	// Parse tail option
	tailLines := -1 // -1 means show all lines
	if tail != "all" {
		if n, err := strconv.Atoi(tail); err == nil && n >= 0 {
			tailLines = n
		} else {
			logger.Warn("Invalid tail value: %s, showing all lines", tail)
		}
	}

	// Parse time filters
	var sinceTime, untilTime time.Time
	if since != "" {
		if t, err := parseTimeOption(since); err == nil {
			sinceTime = t
		} else {
			logger.Warn("Invalid since time format: %s", since)
		}
	}
	if until != "" {
		if t, err := parseTimeOption(until); err == nil {
			untilTime = t
		} else {
			logger.Warn("Invalid until time format: %s", until)
		}
	}

	// Display logs
	if follow && container.Status == state.StatusRunning {
		logger.Debug("Following logs for running container")
		return followLogs(stdoutPath, stderrPath, timestamps, sinceTime, untilTime)
	} else {
		logger.Debug("Displaying static logs (tail: %d)", tailLines)
		return displayLogs(stdoutPath, stderrPath, timestamps, tailLines, sinceTime, untilTime)
	}
}

// getContainerLogDir returns the directory where container logs are stored
func getContainerLogDir(containerID string) string {
	// Use the same base directory as state manager but for logs
	sm := state.NewStateManager()
	stateDir := sm.GetStateDir()

	// Replace "containers" with "logs" in the path
	logDir := filepath.Join(filepath.Dir(stateDir), "logs", containerID)
	return logDir
}

// displayLogs shows static logs from log files
func displayLogs(stdoutPath, stderrPath string, showTimestamps bool, tailLines int, since, until time.Time) error {
	var lines []LogLine

	// Read stdout logs
	if _, err := os.Stat(stdoutPath); err == nil {
		stdoutLines, err := readLogFile(stdoutPath, "stdout", since, until)
		if err != nil {
			logger.Error("Failed to read stdout log: %v", err)
		} else {
			lines = append(lines, stdoutLines...)
		}
	}

	// Read stderr logs
	if _, err := os.Stat(stderrPath); err == nil {
		stderrLines, err := readLogFile(stderrPath, "stderr", since, until)
		if err != nil {
			logger.Error("Failed to read stderr log: %v", err)
		} else {
			lines = append(lines, stderrLines...)
		}
	}

	// Sort lines by timestamp
	sortLogLines(lines)

	// Apply tail limit
	if tailLines >= 0 && tailLines < len(lines) {
		lines = lines[len(lines)-tailLines:]
	}

	// Display lines
	for _, line := range lines {
		if showTimestamps {
			fmt.Printf("%s [%s] %s\n", line.Timestamp.Format(time.RFC3339), line.Stream, line.Content)
		} else {
			fmt.Printf("%s\n", line.Content)
		}
	}

	return nil
}

// followLogs tails log files in real-time
func followLogs(stdoutPath, stderrPath string, showTimestamps bool, since, until time.Time) error {
	// This is a simplified implementation
	// In a production system, you might use inotify or similar for efficient file watching

	// For now, we'll just display existing logs and then poll for new content
	err := displayLogs(stdoutPath, stderrPath, showTimestamps, -1, since, until)
	if err != nil {
		return err
	}

	// TODO: Implement real-time following
	// This would require tracking file positions and polling for new content
	fmt.Println("\n[Following logs - press Ctrl+C to exit...]")

	// Simple polling implementation (not efficient, but functional)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var lastStdoutSize, lastStderrSize int64

	// Get initial file sizes
	if stat, err := os.Stat(stdoutPath); err == nil {
		lastStdoutSize = stat.Size()
	}
	if stat, err := os.Stat(stderrPath); err == nil {
		lastStderrSize = stat.Size()
	}

	for {
		select {
		case <-ticker.C:
			// Check for new content in stdout
			if stat, err := os.Stat(stdoutPath); err == nil && stat.Size() > lastStdoutSize {
				if err := displayNewLogContent(stdoutPath, lastStdoutSize, "stdout", showTimestamps); err != nil {
					logger.Error("Failed to read new stdout content: %v", err)
				}
				lastStdoutSize = stat.Size()
			}

			// Check for new content in stderr
			if stat, err := os.Stat(stderrPath); err == nil && stat.Size() > lastStderrSize {
				if err := displayNewLogContent(stderrPath, lastStderrSize, "stderr", showTimestamps); err != nil {
					logger.Error("Failed to read new stderr content: %v", err)
				}
				lastStderrSize = stat.Size()
			}
		}
	}
}

// LogLine represents a single log line with metadata
type LogLine struct {
	Timestamp time.Time
	Stream    string // "stdout" or "stderr"
	Content   string
}

// readLogFile reads and parses a log file
func readLogFile(path, stream string, since, until time.Time) ([]LogLine, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []LogLine
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Parse timestamp if present (format: YYYY-MM-DDTHH:MM:SS.sssZ message)
		timestamp := time.Now() // Default to current time
		content := line

		// Try to parse RFC3339 timestamp at the beginning
		if len(line) > 20 && line[19] == '.' || len(line) > 19 && (line[19] == 'Z' || line[19] == '+' || line[19] == '-') {
			if t, err := time.Parse(time.RFC3339Nano, line[:20]+"Z"); err == nil {
				timestamp = t
				if len(line) > 21 {
					content = line[21:] // Skip timestamp and space
				}
			} else if t, err := time.Parse(time.RFC3339, line[:20]); err == nil {
				timestamp = t
				if len(line) > 21 {
					content = line[21:] // Skip timestamp and space
				}
			}
		}

		// Apply time filters
		if !since.IsZero() && timestamp.Before(since) {
			continue
		}
		if !until.IsZero() && timestamp.After(until) {
			continue
		}

		lines = append(lines, LogLine{
			Timestamp: timestamp,
			Stream:    stream,
			Content:   content,
		})
	}

	return lines, scanner.Err()
}

// sortLogLines sorts log lines by timestamp
func sortLogLines(lines []LogLine) {
	// Simple bubble sort (fine for typical log volumes)
	n := len(lines)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if lines[j].Timestamp.After(lines[j+1].Timestamp) {
				lines[j], lines[j+1] = lines[j+1], lines[j]
			}
		}
	}
}

// displayNewLogContent shows new content from a specific position in the file
func displayNewLogContent(path string, startPos int64, stream string, showTimestamps bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Seek to the start position
	if _, err := file.Seek(startPos, 0); err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if showTimestamps {
			fmt.Printf("%s [%s] %s\n", time.Now().Format(time.RFC3339), stream, line)
		} else {
			fmt.Printf("%s\n", line)
		}
	}

	return scanner.Err()
}

// parseTimeOption parses various time formats for since/until options
func parseTimeOption(timeStr string) (time.Time, error) {
	// Try RFC3339 format first
	if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
		return t, nil
	}

	// Try RFC3339 without timezone
	if t, err := time.Parse("2006-01-02T15:04:05", timeStr); err == nil {
		return t, nil
	}

	// Try relative time format (e.g., "1h", "30m", "45s")
	if duration, err := time.ParseDuration(timeStr); err == nil {
		return time.Now().Add(-duration), nil
	}

	return time.Time{}, fmt.Errorf("invalid time format: %s", timeStr)
}
