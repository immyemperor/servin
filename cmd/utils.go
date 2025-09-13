package cmd

import (
	"fmt"
	"os"
	"syscall"

	"servin/pkg/state"
)

// resolveContainerRef resolves a container reference (ID, short ID, or name) to a full ID
func resolveContainerRef(sm *state.StateManager, ref string) (string, error) {
	// Try as full ID first
	if _, err := sm.LoadContainer(ref); err == nil {
		return ref, nil
	}

	// Try as short ID
	if fullID, err := sm.FindContainerByShortID(ref); err == nil {
		return fullID, nil
	}

	// Try as name
	if fullID, err := sm.FindContainerByName(ref); err == nil {
		return fullID, nil
	}

	return "", fmt.Errorf("container '%s' not found", ref)
}

// stopContainerProcess stops a container process by PID
func stopContainerProcess(pid int) error {
	// First try SIGTERM for graceful shutdown
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("process %d not found: %v", pid, err)
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		// If SIGTERM fails, try SIGKILL
		if err := process.Signal(syscall.SIGKILL); err != nil {
			return fmt.Errorf("failed to kill process %d: %v", pid, err)
		}
	}

	// Wait for process to exit (with timeout would be better)
	state, err := process.Wait()
	if err != nil {
		return fmt.Errorf("failed to wait for process %d: %v", pid, err)
	}

	fmt.Printf("Process %d exited with status: %s\n", pid, state)
	return nil
}
