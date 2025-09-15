package image

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// RegistryClient handles pulling images from Docker registries
type RegistryClient struct {
	registryURL string
	client      *http.Client
}

// NewRegistryClient creates a new registry client
func NewRegistryClient(registryURL string) *RegistryClient {
	if registryURL == "" {
		registryURL = "https://registry-1.docker.io"
	}

	return &RegistryClient{
		registryURL: registryURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ManifestV2 represents Docker Registry API v2 manifest
type ManifestV2 struct {
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"`
	Config        struct {
		MediaType string `json:"mediaType"`
		Size      int64  `json:"size"`
		Digest    string `json:"digest"`
	} `json:"config"`
	Layers []struct {
		MediaType string `json:"mediaType"`
		Size      int64  `json:"size"`
		Digest    string `json:"digest"`
	} `json:"layers"`
}

// ManifestList represents a manifest list (for multi-arch images)
type ManifestList struct {
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"`
	Manifests     []struct {
		MediaType string `json:"mediaType"`
		Size      int64  `json:"size"`
		Digest    string `json:"digest"`
		Platform  struct {
			Architecture string `json:"architecture"`
			OS           string `json:"os"`
		} `json:"platform"`
	} `json:"manifests"`
}

// ImageConfig represents the configuration from the config blob
type ImageConfigBlob struct {
	Config struct {
		Env        []string          `json:"Env"`
		Cmd        []string          `json:"Cmd"`
		Entrypoint []string          `json:"Entrypoint"`
		WorkingDir string            `json:"WorkingDir"`
		User       string            `json:"User"`
		Labels     map[string]string `json:"Labels"`
	} `json:"config"`
	RootFS struct {
		Type    string   `json:"type"`
		DiffIDs []string `json:"diff_ids"`
	} `json:"rootfs"`
}

// PullImage pulls an image from Docker Hub or another registry
func (m *Manager) PullImage(imageRef string) error {
	fmt.Printf("Pulling image %s from Docker Hub...\n", imageRef)

	// Parse image reference
	repo, tag := parseImageRef(imageRef)
	if tag == "" {
		tag = "latest"
	}

	fmt.Printf("Parsed image: repo=%s, tag=%s\n", repo, tag)

	// Create registry client
	client := NewRegistryClient("")

	// Get auth token for Docker Hub
	fmt.Printf("Getting auth token...\n")
	token, err := client.getAuthToken(repo)
	if err != nil {
		return fmt.Errorf("failed to get auth token: %v", err)
	}

	// Get image manifest
	fmt.Printf("Getting manifest...\n")
	manifest, err := client.getManifest(repo, tag, token)
	if err != nil {
		return fmt.Errorf("failed to get manifest: %v", err)
	}

	fmt.Printf("Manifest received: %d layers, config digest: %s\n", len(manifest.Layers), manifest.Config.Digest)

	// Get config blob
	fmt.Printf("Getting config blob...\n")
	configBlob, err := client.getConfigBlob(repo, manifest.Config.Digest, token)
	if err != nil {
		return fmt.Errorf("failed to get config blob: %v", err)
	}

	// Create image directory
	imageID := generateImageID(fmt.Sprintf("%s:%s", repo, tag), "")
	imageDir := filepath.Join(m.imageDir, imageID)
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		return fmt.Errorf("failed to create image directory: %v", err)
	}

	// Download and extract layers
	rootfsDir := filepath.Join(imageDir, "rootfs")
	if err := os.MkdirAll(rootfsDir, 0755); err != nil {
		return fmt.Errorf("failed to create rootfs directory: %v", err)
	}

	fmt.Printf("Downloading %d layers...\n", len(manifest.Layers))
	for i, layer := range manifest.Layers {
		fmt.Printf("Downloading layer %d/%d...\n", i+1, len(manifest.Layers))
		if err := client.downloadAndExtractLayer(repo, layer.Digest, rootfsDir, token); err != nil {
			return fmt.Errorf("failed to download layer %s: %v", layer.Digest, err)
		}
	}

	// Create image metadata
	img := &Image{
		ID:         imageID,
		RepoTags:   []string{imageRef},
		Created:    time.Now(),
		Size:       calculateLayersSizes(manifest.Layers),
		Layers:     extractLayerDigests(manifest.Layers),
		RootFSType: "layers",
		RootFSPath: rootfsDir,
		Config: ImageConfig{
			Env:        configBlob.Config.Env,
			Cmd:        configBlob.Config.Cmd,
			Entrypoint: configBlob.Config.Entrypoint,
			WorkingDir: configBlob.Config.WorkingDir,
			User:       configBlob.Config.User,
			Labels:     configBlob.Config.Labels,
		},
	}

	// Save image to index
	if err := m.saveImage(img); err != nil {
		return fmt.Errorf("failed to save image: %v", err)
	}

	fmt.Printf("Successfully pulled %s\n", imageRef)
	return nil
}

// getAuthToken gets an authentication token for Docker Hub
func (rc *RegistryClient) getAuthToken(repo string) (string, error) {
	// Docker Hub auth endpoint
	authURL := fmt.Sprintf("https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s:pull", repo)

	resp, err := rc.client.Get(authURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("auth request failed with status %d", resp.StatusCode)
	}

	var authResp struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", err
	}

	return authResp.Token, nil
}

// getManifest gets the image manifest, handling manifest lists
func (rc *RegistryClient) getManifest(repo, tag, token string) (*ManifestV2, error) {
	url := fmt.Sprintf("%s/v2/%s/manifests/%s", rc.registryURL, repo, tag)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	// Try multiple manifest formats including OCI
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json, application/vnd.docker.distribution.manifest.list.v2+json, application/vnd.oci.image.manifest.v1+json, application/vnd.oci.image.index.v1+json")

	resp, err := rc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("manifest request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body to determine manifest type
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest response: %v", err)
	}

	// Parse as generic manifest first to check type
	var genericManifest map[string]interface{}
	if err := json.Unmarshal(body, &genericManifest); err != nil {
		return nil, fmt.Errorf("failed to decode manifest: %v", err)
	}

	mediaType, _ := genericManifest["mediaType"].(string)
	fmt.Printf("Manifest media type: %s\n", mediaType)

	// Handle manifest list (multi-arch) - both Docker and OCI formats
	if mediaType == "application/vnd.docker.distribution.manifest.list.v2+json" ||
		mediaType == "application/vnd.oci.image.index.v1+json" {
		var manifestList ManifestList
		if err := json.Unmarshal(body, &manifestList); err != nil {
			return nil, fmt.Errorf("failed to decode manifest list: %v", err)
		}

		// Find the right manifest for our platform (prefer amd64/linux)
		var targetDigest string
		for _, manifest := range manifestList.Manifests {
			if manifest.Platform.OS == "linux" && manifest.Platform.Architecture == "amd64" {
				targetDigest = manifest.Digest
				break
			}
		}

		// If no amd64 found, use the first one
		if targetDigest == "" && len(manifestList.Manifests) > 0 {
			targetDigest = manifestList.Manifests[0].Digest
		}

		if targetDigest == "" {
			return nil, fmt.Errorf("no suitable manifest found in manifest list")
		}

		fmt.Printf("Found manifest list, using digest: %s\n", targetDigest)

		// Get the specific manifest
		return rc.getManifestByDigest(repo, targetDigest, token)
	}

	// Handle regular manifest
	var manifest ManifestV2
	if err := json.Unmarshal(body, &manifest); err != nil {
		return nil, fmt.Errorf("failed to decode manifest: %v", err)
	}

	return &manifest, nil
}

// getManifestByDigest gets a specific manifest by digest
func (rc *RegistryClient) getManifestByDigest(repo, digest, token string) (*ManifestV2, error) {
	url := fmt.Sprintf("%s/v2/%s/manifests/%s", rc.registryURL, repo, digest)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json, application/vnd.oci.image.manifest.v1+json")

	resp, err := rc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("manifest by digest request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var manifest ManifestV2
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, fmt.Errorf("failed to decode manifest: %v", err)
	}

	return &manifest, nil
}

// getConfigBlob gets the image configuration blob
func (rc *RegistryClient) getConfigBlob(repo, digest, token string) (*ImageConfigBlob, error) {
	url := fmt.Sprintf("%s/v2/%s/blobs/%s", rc.registryURL, repo, digest)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := rc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("config blob request failed with status %d: %s (URL: %s)", resp.StatusCode, string(body), url)
	}

	var configBlob ImageConfigBlob
	if err := json.NewDecoder(resp.Body).Decode(&configBlob); err != nil {
		return nil, fmt.Errorf("failed to decode config blob: %v", err)
	}

	return &configBlob, nil
}

// downloadAndExtractLayer downloads and extracts a layer to the rootfs
func (rc *RegistryClient) downloadAndExtractLayer(repo, digest, rootfsDir, token string) error {
	url := fmt.Sprintf("%s/v2/%s/blobs/%s", rc.registryURL, repo, digest)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := rc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("layer download failed with status %d", resp.StatusCode)
	}

	// Create gzip reader
	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzipReader.Close()

	// Create tar reader
	tarReader := tar.NewReader(gzipReader)

	// Extract tar contents
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("tar read error: %v", err)
		}

		targetPath := filepath.Join(rootfsDir, header.Name)

		// Security check: prevent path traversal
		if !strings.HasPrefix(targetPath, rootfsDir) {
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", targetPath, err)
			}
		case tar.TypeReg:
			// Create parent directories
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %v", targetPath, err)
			}

			// Create file
			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %v", targetPath, err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to write file %s: %v", targetPath, err)
			}
			outFile.Close()
		case tar.TypeSymlink:
			// Create symlink
			if err := os.Symlink(header.Linkname, targetPath); err != nil {
				// Ignore symlink errors for now
				continue
			}
		}
	}

	return nil
}

// parseImageRef parses an image reference into repository and tag
func parseImageRef(imageRef string) (repo, tag string) {
	parts := strings.Split(imageRef, ":")
	if len(parts) == 1 {
		// No tag specified, use latest
		repo = parts[0]
		tag = "latest"
	} else {
		repo = strings.Join(parts[:len(parts)-1], ":")
		tag = parts[len(parts)-1]
	}

	// Handle Docker Hub library images (e.g., "alpine" -> "library/alpine")
	if !strings.Contains(repo, "/") {
		repo = "library/" + repo
	}

	return repo, tag
}

// calculateLayersSizes calculates total size of all layers
func calculateLayersSizes(layers []struct {
	MediaType string `json:"mediaType"`
	Size      int64  `json:"size"`
	Digest    string `json:"digest"`
}) int64 {
	var total int64
	for _, layer := range layers {
		total += layer.Size
	}
	return total
}

// extractLayerDigests extracts layer digests
func extractLayerDigests(layers []struct {
	MediaType string `json:"mediaType"`
	Size      int64  `json:"size"`
	Digest    string `json:"digest"`
}) []string {
	var digests []string
	for _, layer := range layers {
		digests = append(digests, layer.Digest)
	}
	return digests
}

// saveImage saves an image to the index
func (m *Manager) saveImage(img *Image) error {
	images, err := m.ListImages()
	if err != nil {
		return err
	}

	// Add new image
	images = append(images, img)

	// Save updated index
	data, err := json.MarshalIndent(images, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.indexPath, data, 0644)
}
