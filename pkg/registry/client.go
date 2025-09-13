package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"servin/pkg/logger"
)

// Client handles communication with registries
type Client struct {
	config     *RegistryConfig
	httpClient *http.Client
	dataDir    string
}

// NewClient creates a new registry client
func NewClient(dataDir string) (*Client, error) {
	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	config, err := loadConfig(dataDir)
	if err != nil {
		// Create default config if none exists
		config = getDefaultConfig(dataDir)
		if err := saveConfig(dataDir, config); err != nil {
			logger.Warn("Failed to save default registry config: %v", err)
		}
	}

	return &Client{
		config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		dataDir:    dataDir,
	}, nil
}

// PushImage pushes an image to a registry
func (c *Client) PushImage(imageName, tag string, registryURL string, options *PushOptions) error {
	if options == nil {
		options = &PushOptions{}
	}

	// Determine target registry
	targetRegistry := registryURL
	if targetRegistry == "" {
		targetRegistry = c.config.DefaultRegistry
	}
	if targetRegistry == "" {
		targetRegistry = "localhost:" + fmt.Sprintf("%d", c.config.LocalPort)
	}

	// Load image from local image directory (simplified approach)
	imagePath := filepath.Join(c.dataDir, "images", fmt.Sprintf("%s_%s.tar", imageName, tag))
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return fmt.Errorf("failed to load image %s:%s: %w", imageName, tag, err)
	}

	if !options.Quiet {
		logger.Info("Pushing %s:%s to %s", imageName, tag, targetRegistry)
	}

	// Check if it's a local registry
	if strings.Contains(targetRegistry, "localhost") || strings.Contains(targetRegistry, "127.0.0.1") {
		return c.pushToLocal(imageName, tag, targetRegistry, imageData, options)
	}

	// Push to remote registry
	return c.pushToRemote(imageName, tag, targetRegistry, imageData, options)
}

// PullImage pulls an image from a registry
func (c *Client) PullImage(imageName, tag string, registryURL string, options *PullOptions) error {
	if options == nil {
		options = &PullOptions{}
	}

	// Determine source registry
	sourceRegistry := registryURL
	if sourceRegistry == "" {
		sourceRegistry = c.config.DefaultRegistry
	}
	if sourceRegistry == "" {
		sourceRegistry = "localhost:" + fmt.Sprintf("%d", c.config.LocalPort)
	}

	if !options.Quiet {
		logger.Info("Pulling %s:%s from %s", imageName, tag, sourceRegistry)
	}

	var imageData []byte
	var err error

	// Check if it's a local registry
	if strings.Contains(sourceRegistry, "localhost") || strings.Contains(sourceRegistry, "127.0.0.1") {
		imageData, err = c.pullFromLocal(imageName, tag, sourceRegistry, options)
	} else {
		imageData, err = c.pullFromRemote(imageName, tag, sourceRegistry, options)
	}

	if err != nil {
		return err
	}

	// Save image to local image directory
	imagePath := filepath.Join(c.dataDir, "images", fmt.Sprintf("%s_%s.tar", imageName, tag))
	return os.WriteFile(imagePath, imageData, 0644)
}

// ListRegistryImages lists images available in a registry
func (c *Client) ListRegistryImages(registryURL string) (map[string][]string, error) {
	targetRegistry := registryURL
	if targetRegistry == "" {
		targetRegistry = "localhost:" + fmt.Sprintf("%d", c.config.LocalPort)
	}

	// Check if it's a local registry
	if strings.Contains(targetRegistry, "localhost") || strings.Contains(targetRegistry, "127.0.0.1") {
		return c.listFromLocal(targetRegistry)
	}

	return c.listFromRemote(targetRegistry)
}

// LoginToRegistry authenticates with a registry
func (c *Client) LoginToRegistry(registryURL, username, password, email string) error {
	if c.config.Credentials == nil {
		c.config.Credentials = make(map[string]Auth)
	}

	c.config.Credentials[registryURL] = Auth{
		Username: username,
		Password: password,
		Email:    email,
	}

	// Save updated config
	return saveConfig(c.dataDir, c.config)
}

// LogoutFromRegistry removes authentication for a registry
func (c *Client) LogoutFromRegistry(registryURL string) error {
	if c.config.Credentials != nil {
		delete(c.config.Credentials, registryURL)
	}

	// Save updated config
	return saveConfig(c.dataDir, c.config)
}

// GetRegistryInfo returns information about configured registries
func (c *Client) GetRegistryInfo() ([]*RegistryInfo, error) {
	var registries []*RegistryInfo

	// Add local registry
	localInfo := &RegistryInfo{
		Name:   "local",
		URL:    "localhost:" + fmt.Sprintf("%d", c.config.LocalPort),
		Type:   "local",
		Status: "unknown",
	}

	// Check if local registry is running
	if c.isRegistryHealthy(localInfo.URL) {
		localInfo.Status = "healthy"
	} else {
		localInfo.Status = "stopped"
	}
	localInfo.LastCheck = time.Now()

	registries = append(registries, localInfo)

	// Add configured remote registries
	for name, url := range c.config.Registries {
		info := &RegistryInfo{
			Name:   name,
			URL:    url,
			Type:   "remote",
			Status: "unknown",
		}

		if c.isRegistryHealthy(url) {
			info.Status = "healthy"
		} else {
			info.Status = "unreachable"
		}
		info.LastCheck = time.Now()

		registries = append(registries, info)
	}

	return registries, nil
}

// Private methods

func (c *Client) pushToLocal(imageName, tag, registryURL string, imageData []byte, options *PushOptions) error {
	url := fmt.Sprintf("http://%s/v2/%s/manifests/%s", registryURL, imageName, tag)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(imageData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to push to local registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("push failed with status %d: %s", resp.StatusCode, string(body))
	}

	if !options.Quiet {
		logger.Info("Successfully pushed %s:%s", imageName, tag)
	}
	return nil
}

func (c *Client) pullFromLocal(imageName, tag, registryURL string, options *PullOptions) ([]byte, error) {
	url := fmt.Sprintf("http://%s/v2/%s/manifests/%s", registryURL, imageName, tag)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to pull from local registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pull failed with status %d", resp.StatusCode)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	if !options.Quiet {
		logger.Info("Successfully pulled %s:%s", imageName, tag)
	}
	return imageData, nil
}

func (c *Client) pushToRemote(imageName, tag, registryURL string, imageData []byte, options *PushOptions) error {
	// TODO: Implement Docker Registry HTTP API v2 for remote push
	return fmt.Errorf("remote registry push not yet implemented")
}

func (c *Client) pullFromRemote(imageName, tag, registryURL string, options *PullOptions) ([]byte, error) {
	// TODO: Implement Docker Registry HTTP API v2 for remote pull
	return nil, fmt.Errorf("remote registry pull not yet implemented")
}

func (c *Client) listFromLocal(registryURL string) (map[string][]string, error) {
	url := fmt.Sprintf("http://%s/v2/_catalog", registryURL)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to list from local registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list failed with status %d", resp.StatusCode)
	}

	var catalog struct {
		Repositories []string `json:"repositories"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&catalog); err != nil {
		return nil, fmt.Errorf("failed to decode catalog response: %w", err)
	}

	// Convert to format expected by caller
	result := make(map[string][]string)
	for _, repo := range catalog.Repositories {
		result[repo] = []string{"latest"} // Simplified - would need tag enumeration API
	}

	return result, nil
}

func (c *Client) listFromRemote(registryURL string) (map[string][]string, error) {
	// TODO: Implement Docker Registry HTTP API v2 for remote listing
	return nil, fmt.Errorf("remote registry listing not yet implemented")
}

func (c *Client) isRegistryHealthy(registryURL string) bool {
	url := fmt.Sprintf("http://%s/health", registryURL)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// Configuration management

func loadConfig(storageRoot string) (*RegistryConfig, error) {
	configPath := filepath.Join(storageRoot, "registry-config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config RegistryConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse registry config: %w", err)
	}

	return &config, nil
}

func saveConfig(storageRoot string, config *RegistryConfig) error {
	configPath := filepath.Join(storageRoot, "registry-config.json")

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry config: %w", err)
	}

	return os.WriteFile(configPath, data, 0644)
}

func getDefaultConfig(storageRoot string) *RegistryConfig {
	return &RegistryConfig{
		LocalPort:          5000,
		LocalDataDir:       filepath.Join(storageRoot, "registry"),
		DefaultRegistry:    "",
		Registries:         make(map[string]string),
		Credentials:        make(map[string]Auth),
		InsecureRegistries: []string{"localhost:5000"},
		CertificateDir:     filepath.Join(storageRoot, "certs"),
	}
}
