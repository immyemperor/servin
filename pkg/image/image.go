package image

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Image represents a container image
type Image struct {
	ID         string            `json:"id"`
	RepoTags   []string          `json:"repo_tags"`
	Created    time.Time         `json:"created"`
	Size       int64             `json:"size"`
	Layers     []string          `json:"layers"`
	Config     ImageConfig       `json:"config"`
	Metadata   map[string]string `json:"metadata"`
	RootFSType string            `json:"rootfs_type"`
	RootFSPath string            `json:"rootfs_path"`
}

// ImageConfig holds the configuration for the image
type ImageConfig struct {
	Env          []string            `json:"env"`
	Cmd          []string            `json:"cmd"`
	Entrypoint   []string            `json:"entrypoint"`
	WorkingDir   string              `json:"working_dir"`
	User         string              `json:"user"`
	ExposedPorts map[string]struct{} `json:"exposed_ports"`
	Labels       map[string]string   `json:"labels"`
}

// Manager manages container images
type Manager struct {
	imageDir  string
	indexPath string
}

// NewManager creates a new image manager
func NewManager() *Manager {
	var imageDir string

	switch runtime.GOOS {
	case "windows":
		// Windows: Use user home directory
		homeDir, _ := os.UserHomeDir()
		imageDir = filepath.Join(homeDir, ".servin", "images")
	case "darwin":
		// macOS: Use user home directory (similar to Windows but Unix paths)
		homeDir, _ := os.UserHomeDir()
		imageDir = filepath.Join(homeDir, ".servin", "images")
	case "linux":
		// Linux: Use system directory (though this shouldn't be called on Linux)
		imageDir = "/var/lib/servin/images"
	default:
		// Other Unix-like systems: Use /var/lib
		imageDir = "/var/lib/servin/images"
	}

	return &Manager{
		imageDir:  imageDir,
		indexPath: filepath.Join(imageDir, "index.json"),
	}
}

// ensureImageDir creates the image directory if it doesn't exist
func (m *Manager) ensureImageDir() error {
	return os.MkdirAll(m.imageDir, 0755)
}

// ListImages returns all available images
func (m *Manager) ListImages() ([]*Image, error) {
	if err := m.ensureImageDir(); err != nil {
		return nil, fmt.Errorf("failed to ensure image directory: %v", err)
	}

	// Check if index file exists
	if _, err := os.Stat(m.indexPath); os.IsNotExist(err) {
		return []*Image{}, nil
	}

	data, err := os.ReadFile(m.indexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image index: %v", err)
	}

	var images []*Image
	if err := json.Unmarshal(data, &images); err != nil {
		return nil, fmt.Errorf("failed to parse image index: %v", err)
	}

	return images, nil
}

// GetImage retrieves an image by reference (name:tag or ID)
func (m *Manager) GetImage(ref string) (*Image, error) {
	images, err := m.ListImages()
	if err != nil {
		return nil, err
	}

	// Try to find by exact match first
	for _, img := range images {
		// Check by ID
		if img.ID == ref || strings.HasPrefix(img.ID, ref) {
			return img, nil
		}
		// Check by repo tags
		for _, tag := range img.RepoTags {
			if tag == ref {
				return img, nil
			}
		}
	}

	return nil, fmt.Errorf("image '%s' not found", ref)
}

// SaveImage saves an image to the index
func (m *Manager) SaveImage(img *Image) error {
	if err := m.ensureImageDir(); err != nil {
		return fmt.Errorf("failed to ensure image directory: %v", err)
	}

	images, err := m.ListImages()
	if err != nil {
		return err
	}

	// Check if image already exists and update it
	found := false
	for i, existingImg := range images {
		if existingImg.ID == img.ID {
			images[i] = img
			found = true
			break
		}
	}

	// If not found, add as new image
	if !found {
		images = append(images, img)
	}

	// Save updated index
	data, err := json.MarshalIndent(images, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal image index: %v", err)
	}

	if err := os.WriteFile(m.indexPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write image index: %v", err)
	}

	return nil
}

// RemoveImage removes an image by reference
func (m *Manager) RemoveImage(ref string) error {
	images, err := m.ListImages()
	if err != nil {
		return err
	}

	// Find and remove the image
	var updatedImages []*Image
	found := false
	var removedImage *Image

	for _, img := range images {
		remove := false

		// Check by ID
		if img.ID == ref || strings.HasPrefix(img.ID, ref) {
			remove = true
		}

		// Check by repo tags
		for _, tag := range img.RepoTags {
			if tag == ref {
				remove = true
				break
			}
		}

		if remove {
			found = true
			removedImage = img
		} else {
			updatedImages = append(updatedImages, img)
		}
	}

	if !found {
		return fmt.Errorf("image '%s' not found", ref)
	}

	// Clean up image files
	if removedImage.RootFSPath != "" {
		if err := os.RemoveAll(removedImage.RootFSPath); err != nil {
			fmt.Printf("Warning: failed to remove image rootfs: %v\n", err)
		}
	}

	// Save updated index
	data, err := json.MarshalIndent(updatedImages, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal image index: %v", err)
	}

	if err := os.WriteFile(m.indexPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write image index: %v", err)
	}

	return nil
}

// TagImage adds a new tag to an existing image
func (m *Manager) TagImage(sourceRef, targetTag string) error {
	// Find the source image
	sourceImage, err := m.GetImage(sourceRef)
	if err != nil {
		return fmt.Errorf("source image not found: %v", err)
	}

	// Validate target tag format
	if !strings.Contains(targetTag, ":") {
		targetTag = targetTag + ":latest"
	}

	// Check if target tag already exists
	if existingImage, err := m.GetImage(targetTag); err == nil {
		return fmt.Errorf("tag '%s' already exists on image %s", targetTag, existingImage.ID[:12])
	}

	// Create a copy of the source image with the new tag
	newImage := &Image{
		ID:         sourceImage.ID,                          // Same image ID
		RepoTags:   append(sourceImage.RepoTags, targetTag), // Add new tag
		Created:    sourceImage.Created,
		Size:       sourceImage.Size,
		Layers:     sourceImage.Layers,
		Config:     sourceImage.Config,
		Metadata:   sourceImage.Metadata,
		RootFSType: sourceImage.RootFSType,
		RootFSPath: sourceImage.RootFSPath,
	}

	// Save the updated image
	if err := m.SaveImage(newImage); err != nil {
		return fmt.Errorf("failed to save tagged image: %v", err)
	}

	return nil
}

// GetImageDir returns the image directory path
func (m *Manager) GetImageDir() string {
	return m.imageDir
}

// CreateImageFromTarball creates an image from a tarball
func (m *Manager) CreateImageFromTarball(tarballPath, name, tag string) (*Image, error) {
	if err := m.ensureImageDir(); err != nil {
		return nil, fmt.Errorf("failed to ensure image directory: %v", err)
	}

	// Generate image ID
	imageID := generateImageID(name, tag)

	// Create image directory
	imagePath := filepath.Join(m.imageDir, imageID)
	if err := os.MkdirAll(imagePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create image directory: %v", err)
	}

	// Extract tarball to image directory
	if err := extractTarball(tarballPath, imagePath); err != nil {
		return nil, fmt.Errorf("failed to extract tarball: %v", err)
	}

	// Get tarball size
	stat, err := os.Stat(tarballPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat tarball: %v", err)
	}

	// Create image metadata
	repoTag := fmt.Sprintf("%s:%s", name, tag)
	image := &Image{
		ID:         imageID,
		RepoTags:   []string{repoTag},
		Created:    time.Now(),
		Size:       stat.Size(),
		Layers:     []string{imageID}, // Simple single-layer for now
		RootFSType: "tarball",
		RootFSPath: imagePath,
		Config: ImageConfig{
			Cmd:          []string{"/bin/sh"},
			WorkingDir:   "/",
			User:         "root",
			Env:          []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
			Labels:       make(map[string]string),
			ExposedPorts: make(map[string]struct{}),
		},
		Metadata: map[string]string{
			"source":        "tarball",
			"original_path": tarballPath,
		},
	}

	// Save image to index
	if err := m.SaveImage(image); err != nil {
		return nil, fmt.Errorf("failed to save image: %v", err)
	}

	return image, nil
}

// GetImageRootFS returns the rootfs path for an image
func (m *Manager) GetImageRootFS(ref string) (string, error) {
	image, err := m.GetImage(ref)
	if err != nil {
		return "", err
	}

	if image.RootFSPath == "" {
		return "", fmt.Errorf("image '%s' has no rootfs", ref)
	}

	return image.RootFSPath, nil
}
