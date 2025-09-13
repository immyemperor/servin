package registry

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"servin/pkg/logger"
)

// LocalRegistry implements a simple HTTP-based local registry
type LocalRegistry struct {
	dataDir string
	port    int
	server  *http.Server
}

// NewLocalRegistry creates a new local registry instance
func NewLocalRegistry(dataDir string, port int) *LocalRegistry {
	return &LocalRegistry{
		dataDir: dataDir,
		port:    port,
	}
}

// Start starts the local registry HTTP server
func (lr *LocalRegistry) Start() error {
	// Ensure data directory exists
	if err := os.MkdirAll(lr.dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create registry data directory: %w", err)
	}

	// Create subdirectories
	for _, dir := range []string{"images", "manifests", "blobs"} {
		if err := os.MkdirAll(filepath.Join(lr.dataDir, dir), 0755); err != nil {
			return fmt.Errorf("failed to create %s directory: %w", dir, err)
		}
	}

	mux := http.NewServeMux()

	// Registry API endpoints
	mux.HandleFunc("/v2/", lr.handleRoot)
	mux.HandleFunc("/v2/_catalog", lr.handleCatalog)
	mux.HandleFunc("/v2/", lr.handleImageOperations)

	// Management endpoints
	mux.HandleFunc("/health", lr.handleHealth)
	mux.HandleFunc("/info", lr.handleInfo)

	lr.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", lr.port),
		Handler: mux,
	}

	logger.Info("Starting local registry server on port %d", lr.port)
	logger.Info("Registry data directory: %s", lr.dataDir)

	return lr.server.ListenAndServe()
}

// Stop stops the local registry server
func (lr *LocalRegistry) Stop() error {
	if lr.server != nil {
		logger.Info("Stopping local registry server")
		return lr.server.Close()
	}
	return nil
}

// Push uploads an image to the local registry
func (lr *LocalRegistry) Push(imageName, tag string, imageData []byte) error {
	// Calculate digest
	hash := sha256.Sum256(imageData)
	digest := fmt.Sprintf("sha256:%x", hash)

	// Store blob
	blobPath := filepath.Join(lr.dataDir, "blobs", digest)
	if err := ioutil.WriteFile(blobPath, imageData, 0644); err != nil {
		return fmt.Errorf("failed to write blob: %w", err)
	}

	// Create manifest
	manifest := &Manifest{
		SchemaVersion: 2,
		MediaType:     "application/vnd.docker.distribution.manifest.v2+json",
		Name:          imageName,
		Tag:           tag,
		Size:          int64(len(imageData)),
		Digest:        digest,
		CreatedAt:     time.Now(),
		Layers: []Layer{
			{
				MediaType: "application/vnd.docker.image.rootfs.diff.tar",
				Size:      int64(len(imageData)),
				Digest:    digest,
			},
		},
	}

	// Store manifest
	manifestData, err := json.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	manifestPath := filepath.Join(lr.dataDir, "manifests", fmt.Sprintf("%s_%s.json", imageName, tag))
	if err := ioutil.WriteFile(manifestPath, manifestData, 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	logger.Info("Pushed image %s:%s to local registry", imageName, tag)
	return nil
}

// Pull downloads an image from the local registry
func (lr *LocalRegistry) Pull(imageName, tag string) ([]byte, error) {
	// Load manifest
	manifestPath := filepath.Join(lr.dataDir, "manifests", fmt.Sprintf("%s_%s.json", imageName, tag))
	manifestData, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("image %s:%s not found in registry: %w", imageName, tag, err)
	}

	var manifest Manifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	// Load blob
	blobPath := filepath.Join(lr.dataDir, "blobs", manifest.Digest)
	imageData, err := ioutil.ReadFile(blobPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	logger.Info("Pulled image %s:%s from local registry", imageName, tag)
	return imageData, nil
}

// List returns available images and tags
func (lr *LocalRegistry) List() (map[string][]string, error) {
	images := make(map[string][]string)

	manifestsDir := filepath.Join(lr.dataDir, "manifests")
	files, err := ioutil.ReadDir(manifestsDir)
	if err != nil {
		// Return empty list if directory doesn't exist
		return images, nil
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		// Parse filename: imagename_tag.json
		name := strings.TrimSuffix(file.Name(), ".json")
		parts := strings.Split(name, "_")
		if len(parts) >= 2 {
			imageName := strings.Join(parts[:len(parts)-1], "_")
			tag := parts[len(parts)-1]

			if _, exists := images[imageName]; !exists {
				images[imageName] = []string{}
			}
			images[imageName] = append(images[imageName], tag)
		}
	}

	return images, nil
}

// Delete removes an image from the local registry
func (lr *LocalRegistry) Delete(imageName, tag string) error {
	manifestPath := filepath.Join(lr.dataDir, "manifests", fmt.Sprintf("%s_%s.json", imageName, tag))

	// Load manifest to get blob digest
	manifestData, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("image %s:%s not found: %w", imageName, tag, err)
	}

	var manifest Manifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}

	// Remove manifest
	if err := os.Remove(manifestPath); err != nil {
		return fmt.Errorf("failed to remove manifest: %w", err)
	}

	// Remove blob (TODO: implement reference counting to avoid removing shared blobs)
	blobPath := filepath.Join(lr.dataDir, "blobs", manifest.Digest)
	if err := os.Remove(blobPath); err != nil {
		logger.Warn("Failed to remove blob %s: %v", manifest.Digest, err)
	}

	logger.Info("Deleted image %s:%s from local registry", imageName, tag)
	return nil
}

// GetManifest returns the manifest for an image
func (lr *LocalRegistry) GetManifest(imageName, tag string) (*Manifest, error) {
	manifestPath := filepath.Join(lr.dataDir, "manifests", fmt.Sprintf("%s_%s.json", imageName, tag))
	manifestData, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("manifest for %s:%s not found: %w", imageName, tag, err)
	}

	var manifest Manifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	return &manifest, nil
}

// HTTP handlers

func (lr *LocalRegistry) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Servin Local Registry",
		"version": "v2",
	})
}

func (lr *LocalRegistry) handleCatalog(w http.ResponseWriter, r *http.Request) {
	images, err := lr.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var repositories []string
	for imageName := range images {
		repositories = append(repositories, imageName)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"repositories": repositories,
	})
}

func (lr *LocalRegistry) handleImageOperations(w http.ResponseWriter, r *http.Request) {
	// Parse URL path for image operations
	path := strings.TrimPrefix(r.URL.Path, "/v2/")

	if path == "" {
		lr.handleRoot(w, r)
		return
	}

	// Handle image-specific operations based on path and method
	parts := strings.Split(path, "/")
	if len(parts) >= 1 {
		imageName := parts[0]

		if len(parts) >= 3 && parts[1] == "manifests" {
			tag := parts[2]

			switch r.Method {
			case http.MethodGet:
				lr.handleGetManifest(w, r, imageName, tag)
			case http.MethodPut:
				lr.handlePutManifest(w, r, imageName, tag)
			case http.MethodDelete:
				lr.handleDeleteManifest(w, r, imageName, tag)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	}
}

func (lr *LocalRegistry) handleGetManifest(w http.ResponseWriter, r *http.Request, imageName, tag string) {
	manifest, err := lr.GetManifest(imageName, tag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")
	json.NewEncoder(w).Encode(manifest)
}

func (lr *LocalRegistry) handlePutManifest(w http.ResponseWriter, r *http.Request, imageName, tag string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// For simplicity, treat the entire body as image data
	if err := lr.Push(imageName, tag, body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (lr *LocalRegistry) handleDeleteManifest(w http.ResponseWriter, r *http.Request, imageName, tag string) {
	if err := lr.Delete(imageName, tag); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (lr *LocalRegistry) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func (lr *LocalRegistry) handleInfo(w http.ResponseWriter, r *http.Request) {
	images, _ := lr.List()

	info := map[string]interface{}{
		"registry_type": "local",
		"data_dir":      lr.dataDir,
		"port":          lr.port,
		"images":        images,
		"uptime":        time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}
