package volume

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"servin/pkg/errors"
	"servin/pkg/logger"
)

// Volume represents a managed volume
type Volume struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Mountpoint string            `json:"mountpoint"`
	CreatedAt  time.Time         `json:"created_at"`
	Labels     map[string]string `json:"labels"`
	Options    map[string]string `json:"options"`
	Scope      string            `json:"scope"`
	Status     map[string]string `json:"status"`
}

// Manager manages container volumes
type Manager struct {
	volumeDir string
	indexPath string
}

// NewManager creates a new volume manager
func NewManager() *Manager {
	var volumeDir string

	switch runtime.GOOS {
	case "windows":
		// Windows: Use user home directory
		homeDir, _ := os.UserHomeDir()
		volumeDir = filepath.Join(homeDir, ".servin", "volumes")
	case "darwin":
		// macOS: Use user home directory
		homeDir, _ := os.UserHomeDir()
		volumeDir = filepath.Join(homeDir, ".servin", "volumes")
	case "linux":
		// Linux: Use system directory
		volumeDir = "/var/lib/servin/volumes"
	default:
		// Other Unix-like systems: Use /var/lib
		volumeDir = "/var/lib/servin/volumes"
	}

	return &Manager{
		volumeDir: volumeDir,
		indexPath: filepath.Join(volumeDir, "index.json"),
	}
}

// ensureVolumeDir creates the volume directory if it doesn't exist
func (m *Manager) ensureVolumeDir() error {
	return os.MkdirAll(m.volumeDir, 0755)
}

// ListVolumes returns all available volumes
func (m *Manager) ListVolumes() ([]*Volume, error) {
	logger.Debug("Listing volumes from directory: %s", m.volumeDir)

	if err := m.ensureVolumeDir(); err != nil {
		logger.Error("Failed to ensure volume directory: %v", err)
		return nil, errors.WrapError(err, errors.ErrTypeIO, "ListVolumes", "failed to ensure volume directory")
	}

	// Check if index file exists
	if _, err := os.Stat(m.indexPath); os.IsNotExist(err) {
		logger.Debug("Volume index file does not exist, returning empty list")
		return []*Volume{}, nil
	}

	data, err := os.ReadFile(m.indexPath)
	if err != nil {
		logger.Error("Failed to read volume index file: %v", err)
		return nil, errors.WrapError(err, errors.ErrTypeIO, "ListVolumes", "failed to read volume index file").
			WithContext("index_path", m.indexPath)
	}

	var volumes []*Volume
	if err := json.Unmarshal(data, &volumes); err != nil {
		logger.Error("Failed to parse volume index JSON: %v", err)
		return nil, errors.WrapError(err, errors.ErrTypeIO, "ListVolumes", "failed to parse volume index JSON").
			WithContext("index_path", m.indexPath)
	}

	logger.Debug("Successfully loaded %d volumes from index", len(volumes))
	return volumes, nil
}

// GetVolume retrieves a volume by name
func (m *Manager) GetVolume(name string) (*Volume, error) {
	volumes, err := m.ListVolumes()
	if err != nil {
		return nil, err
	}

	for _, vol := range volumes {
		if vol.Name == name {
			return vol, nil
		}
	}

	return nil, fmt.Errorf("volume '%s' not found", name)
}

// CreateVolume creates a new volume
func (m *Manager) CreateVolume(name string, driver string, options map[string]string, labels map[string]string) (*Volume, error) {
	logger.Debug("Creating volume: %s (driver: %s)", name, driver)

	if err := m.ensureVolumeDir(); err != nil {
		logger.Error("Failed to ensure volume directory: %v", err)
		return nil, errors.WrapError(err, errors.ErrTypeIO, "CreateVolume", "failed to ensure volume directory")
	}

	// Check if volume already exists
	if _, err := m.GetVolume(name); err == nil {
		logger.Warn("Attempted to create volume that already exists: %s", name)
		return nil, errors.NewConflictError("CreateVolume", fmt.Sprintf("volume '%s' already exists", name)).
			WithContext("volume_name", name)
	}

	// Validate volume name
	if strings.Contains(name, "/") || strings.Contains(name, "\\") {
		logger.Error("Invalid volume name contains path separators: %s", name)
		return nil, errors.NewValidationError("CreateVolume", "volume name cannot contain path separators").
			WithContext("volume_name", name)
	}

	// Create volume directory
	volumePath := filepath.Join(m.volumeDir, name)
	logger.Debug("Creating volume directory: %s", volumePath)
	if err := os.MkdirAll(volumePath, 0755); err != nil {
		logger.Error("Failed to create volume directory: %v", err)
		return nil, errors.WrapError(err, errors.ErrTypeIO, "CreateVolume", "failed to create volume directory").
			WithContext("volume_path", volumePath).
			WithContext("volume_name", name)
	}

	// Set default driver if not specified
	if driver == "" {
		driver = "local"
		logger.Debug("Using default driver 'local' for volume: %s", name)
	}

	// Initialize labels and options if nil
	if labels == nil {
		labels = make(map[string]string)
	}
	if options == nil {
		options = make(map[string]string)
	}

	volume := &Volume{
		Name:       name,
		Driver:     driver,
		Mountpoint: volumePath,
		CreatedAt:  time.Now(),
		Labels:     labels,
		Options:    options,
		Scope:      "local",
		Status:     map[string]string{"state": "ready"},
	}

	// Save volume to index
	if err := m.SaveVolume(volume); err != nil {
		// Clean up created directory on failure
		os.RemoveAll(volumePath)
		return nil, fmt.Errorf("failed to save volume: %v", err)
	}

	return volume, nil
}

// SaveVolume saves a volume to the index
func (m *Manager) SaveVolume(vol *Volume) error {
	if err := m.ensureVolumeDir(); err != nil {
		return fmt.Errorf("failed to ensure volume directory: %v", err)
	}

	volumes, err := m.ListVolumes()
	if err != nil {
		return err
	}

	// Check if volume already exists and update it
	found := false
	for i, existingVol := range volumes {
		if existingVol.Name == vol.Name {
			volumes[i] = vol
			found = true
			break
		}
	}

	// If not found, add as new volume
	if !found {
		volumes = append(volumes, vol)
	}

	// Save updated index
	data, err := json.MarshalIndent(volumes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal volume index: %v", err)
	}

	if err := os.WriteFile(m.indexPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write volume index: %v", err)
	}

	return nil
}

// RemoveVolume removes a volume by name
func (m *Manager) RemoveVolume(name string, force bool) error {
	volume, err := m.GetVolume(name)
	if err != nil {
		return err
	}

	// TODO: Check if volume is in use by any containers
	// For now, we'll proceed with removal

	// Get all volumes
	volumes, err := m.ListVolumes()
	if err != nil {
		return err
	}

	// Remove the volume from the list
	var updatedVolumes []*Volume
	for _, vol := range volumes {
		if vol.Name != name {
			updatedVolumes = append(updatedVolumes, vol)
		}
	}

	// Remove volume directory
	if err := os.RemoveAll(volume.Mountpoint); err != nil && !force {
		return fmt.Errorf("failed to remove volume directory: %v", err)
	}

	// Save updated index
	data, err := json.MarshalIndent(updatedVolumes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal volume index: %v", err)
	}

	if err := os.WriteFile(m.indexPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write volume index: %v", err)
	}

	return nil
}

// RemoveAllVolumes removes all volumes
func (m *Manager) RemoveAllVolumes(force bool) error {
	volumes, err := m.ListVolumes()
	if err != nil {
		return err
	}

	var errors []string

	for _, volume := range volumes {
		if err := m.RemoveVolume(volume.Name, force); err != nil {
			errors = append(errors, fmt.Sprintf("failed to remove volume '%s': %v", volume.Name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors occurred while removing volumes:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

// GetVolumeDir returns the volume directory path
func (m *Manager) GetVolumeDir() string {
	return m.volumeDir
}

// PruneVolumes removes all unused volumes
func (m *Manager) PruneVolumes() ([]string, error) {
	// TODO: Implement volume usage checking
	// For now, this is a placeholder that returns empty list
	return []string{}, nil
}
