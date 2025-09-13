package state

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"servin/pkg/network"
)

// Container status constants
const (
	StatusCreated = "created"
	StatusRunning = "running"
	StatusStopped = "stopped"
	StatusExited  = "exited"
)

// ContainerState represents the persistent state of a container
type ContainerState struct {
	ID           string                `json:"id"`
	Name         string                `json:"name"`
	Image        string                `json:"image"`
	Command      string                `json:"command"`
	Args         []string              `json:"args"`
	Status       string                `json:"status"` // created, running, stopped, exited
	PID          int                   `json:"pid"`
	ExitCode     int                   `json:"exit_code"`
	Created      time.Time             `json:"created"`
	Started      time.Time             `json:"started,omitempty"`
	Finished     time.Time             `json:"finished,omitempty"`
	RootPath     string                `json:"root_path"`
	Hostname     string                `json:"hostname"`
	WorkDir      string                `json:"work_dir"`
	Env          map[string]string     `json:"env"`
	Volumes      map[string]string     `json:"volumes"`
	NetworkMode  string                `json:"network_mode"`
	PortMappings []network.PortMapping `json:"port_mappings"`
	Memory       string                `json:"memory"`
	CPUs         string                `json:"cpus"`
}

// StateManager manages container state persistence
type StateManager struct {
	stateDir string
}

// NewStateManager creates a new state manager
func NewStateManager() *StateManager {
	// Use different paths for different platforms
	var stateDir string

	switch runtime.GOOS {
	case "windows":
		// Windows: Use user home directory
		homeDir, _ := os.UserHomeDir()
		stateDir = filepath.Join(homeDir, ".servin", "containers")
	case "darwin":
		// macOS: Use user home directory (similar to Windows but Unix paths)
		homeDir, _ := os.UserHomeDir()
		stateDir = filepath.Join(homeDir, ".servin", "containers")
	case "linux":
		// Linux: Use system directory
		stateDir = "/var/lib/servin/containers"
	default:
		// Other Unix-like systems: Use /var/lib
		stateDir = "/var/lib/servin/containers"
	}

	return &StateManager{
		stateDir: stateDir,
	}
}

// ensureStateDir creates the state directory if it doesn't exist
func (sm *StateManager) ensureStateDir() error {
	return os.MkdirAll(sm.stateDir, 0755)
}

// SaveContainer saves container state to disk
func (sm *StateManager) SaveContainer(state *ContainerState) error {
	if err := sm.ensureStateDir(); err != nil {
		return fmt.Errorf("failed to create state directory: %v", err)
	}

	statePath := filepath.Join(sm.stateDir, state.ID+".json")
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal container state: %v", err)
	}

	if err := os.WriteFile(statePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write container state: %v", err)
	}

	fmt.Printf("Container state saved: %s\n", statePath)
	return nil
}

// LoadContainer loads container state from disk
func (sm *StateManager) LoadContainer(id string) (*ContainerState, error) {
	statePath := filepath.Join(sm.stateDir, id+".json")
	data, err := os.ReadFile(statePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read container state: %v", err)
	}

	var state ContainerState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal container state: %v", err)
	}

	return &state, nil
}

// ListContainers returns all containers (running and stopped)
func (sm *StateManager) ListContainers() ([]*ContainerState, error) {
	if err := sm.ensureStateDir(); err != nil {
		return nil, fmt.Errorf("failed to access state directory: %v", err)
	}

	var containers []*ContainerState

	err := filepath.WalkDir(sm.stateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".json") {
			data, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("Warning: failed to read %s: %v\n", path, err)
				return nil
			}

			var state ContainerState
			if err := json.Unmarshal(data, &state); err != nil {
				fmt.Printf("Warning: failed to parse %s: %v\n", path, err)
				return nil
			}

			containers = append(containers, &state)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan state directory: %v", err)
	}

	return containers, nil
}

// DeleteContainer removes container state from disk
func (sm *StateManager) DeleteContainer(id string) error {
	statePath := filepath.Join(sm.stateDir, id+".json")
	if err := os.Remove(statePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete container state: %v", err)
	}
	return nil
}

// UpdateContainerStatus updates just the status of a container
func (sm *StateManager) UpdateContainerStatus(id, status string) error {
	state, err := sm.LoadContainer(id)
	if err != nil {
		return err
	}

	state.Status = status

	// Update timestamps based on status
	switch status {
	case "running":
		if state.Started.IsZero() {
			state.Started = time.Now()
		}
	case "stopped", "exited":
		state.Finished = time.Now()
	}

	return sm.SaveContainer(state)
}

// UpdateContainerPID updates the PID of a container
func (sm *StateManager) UpdateContainerPID(id string, pid int) error {
	state, err := sm.LoadContainer(id)
	if err != nil {
		return err
	}

	state.PID = pid
	return sm.SaveContainer(state)
}

// FindContainerByName finds a container by name (returns ID)
func (sm *StateManager) FindContainerByName(name string) (string, error) {
	containers, err := sm.ListContainers()
	if err != nil {
		return "", err
	}

	for _, container := range containers {
		if container.Name == name {
			return container.ID, nil
		}
	}

	return "", fmt.Errorf("container with name '%s' not found", name)
}

// FindContainerByShortID finds a container by short ID (first 12 chars)
func (sm *StateManager) FindContainerByShortID(shortID string) (string, error) {
	containers, err := sm.ListContainers()
	if err != nil {
		return "", err
	}

	for _, container := range containers {
		if strings.HasPrefix(container.ID, shortID) {
			return container.ID, nil
		}
	}

	return "", fmt.Errorf("container with ID '%s*' not found", shortID)
}

// GetStateDir returns the state directory path
func (sm *StateManager) GetStateDir() string {
	return sm.stateDir
}
