package registry

import (
	"time"
)

// Registry represents a container image registry
type Registry interface {
	// Push uploads an image to the registry
	Push(imageName, tag string, imageData []byte) error

	// Pull downloads an image from the registry
	Pull(imageName, tag string) ([]byte, error)

	// List returns available images and tags
	List() (map[string][]string, error)

	// Delete removes an image from the registry
	Delete(imageName, tag string) error

	// GetManifest returns the manifest for an image
	GetManifest(imageName, tag string) (*Manifest, error)
}

// RegistryConfig holds configuration for registry operations
type RegistryConfig struct {
	// Local registry settings
	LocalPort    int    `json:"local_port"`
	LocalDataDir string `json:"local_data_dir"`

	// Remote registry settings
	DefaultRegistry string            `json:"default_registry"`
	Registries      map[string]string `json:"registries"`
	Credentials     map[string]Auth   `json:"credentials"`

	// TLS settings
	InsecureRegistries []string `json:"insecure_registries"`
	CertificateDir     string   `json:"certificate_dir"`
}

// Auth holds authentication credentials for a registry
type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email,omitempty"`
	Token    string `json:"token,omitempty"`
}

// Manifest represents an image manifest
type Manifest struct {
	SchemaVersion int       `json:"schemaVersion"`
	MediaType     string    `json:"mediaType"`
	Name          string    `json:"name"`
	Tag           string    `json:"tag"`
	Architecture  string    `json:"architecture"`
	Size          int64     `json:"size"`
	Digest        string    `json:"digest"`
	CreatedAt     time.Time `json:"createdAt"`

	// Layers information
	Layers []Layer `json:"layers"`
}

// Layer represents a layer in an image
type Layer struct {
	MediaType string `json:"mediaType"`
	Size      int64  `json:"size"`
	Digest    string `json:"digest"`
}

// PushOptions contains options for pushing images
type PushOptions struct {
	Registry string
	Force    bool
	Quiet    bool
}

// PullOptions contains options for pulling images
type PullOptions struct {
	Registry string
	Quiet    bool
	Platform string
}

// RegistryInfo contains information about a registry
type RegistryInfo struct {
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Type      string    `json:"type"` // "local" or "remote"
	Status    string    `json:"status"`
	LastCheck time.Time `json:"last_check"`
}
