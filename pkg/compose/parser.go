package compose

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// ComposeFile represents the structure of a servin-compose.yml file
type ComposeFile struct {
	Version  string                   `yaml:"version"`
	Services map[string]ServiceConfig `yaml:"services"`
	Networks map[string]NetworkConfig `yaml:"networks,omitempty"`
	Volumes  map[string]VolumeConfig  `yaml:"volumes,omitempty"`
}

// ServiceConfig represents a service configuration
type ServiceConfig struct {
	Image         string      `yaml:"image,omitempty"`
	Build         interface{} `yaml:"build,omitempty"`
	Command       interface{} `yaml:"command,omitempty"`
	Entrypoint    interface{} `yaml:"entrypoint,omitempty"`
	Environment   interface{} `yaml:"environment,omitempty"`
	Ports         []string    `yaml:"ports,omitempty"`
	Volumes       []string    `yaml:"volumes,omitempty"`
	Networks      interface{} `yaml:"networks,omitempty"`
	DependsOn     []string    `yaml:"depends_on,omitempty"`
	Restart       string      `yaml:"restart,omitempty"`
	WorkingDir    string      `yaml:"working_dir,omitempty"`
	User          string      `yaml:"user,omitempty"`
	Hostname      string      `yaml:"hostname,omitempty"`
	Labels        interface{} `yaml:"labels,omitempty"` // Can be slice or map
	Expose        []string    `yaml:"expose,omitempty"`
	Links         []string    `yaml:"links,omitempty"`
	ExternalLinks []string    `yaml:"external_links,omitempty"`
}

// GetBuildConfig returns the normalized build configuration
func (s *ServiceConfig) GetBuildConfig() BuildConfig {
	return getBuildConfig(s.Build)
}

// BuildConfig represents build configuration
type BuildConfig struct {
	Context    string            `yaml:"context,omitempty"`
	Dockerfile string            `yaml:"dockerfile,omitempty"`
	Args       map[string]string `yaml:"args,omitempty"`
	Target     string            `yaml:"target,omitempty"`
}

// NetworkConfig represents network configuration
type NetworkConfig struct {
	Driver     string            `yaml:"driver,omitempty"`
	DriverOpts map[string]string `yaml:"driver_opts,omitempty"`
	IPAM       IPAMConfig        `yaml:"ipam,omitempty"`
	External   bool              `yaml:"external,omitempty"`
	Labels     interface{}       `yaml:"labels,omitempty"` // Can be slice or map
}

// IPAMConfig represents IPAM configuration
type IPAMConfig struct {
	Driver string       `yaml:"driver,omitempty"`
	Config []IPAMSubnet `yaml:"config,omitempty"`
}

// IPAMSubnet represents IPAM subnet configuration
type IPAMSubnet struct {
	Subnet  string `yaml:"subnet,omitempty"`
	Gateway string `yaml:"gateway,omitempty"`
}

// VolumeConfig represents volume configuration
type VolumeConfig struct {
	Driver     string            `yaml:"driver,omitempty"`
	DriverOpts map[string]string `yaml:"driver_opts,omitempty"`
	External   bool              `yaml:"external,omitempty"`
	Labels     interface{}       `yaml:"labels,omitempty"` // Can be slice or map
}

// ParseComposeFile parses a servin-compose.yml file
func ParseComposeFile(filePath string) (*ComposeFile, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read compose file: %w", err)
	}

	var compose ComposeFile
	err = yaml.Unmarshal(data, &compose)
	if err != nil {
		return nil, fmt.Errorf("failed to parse compose file: %w", err)
	}

	// Set default version if not specified
	if compose.Version == "" {
		compose.Version = "3.8"
	}

	// Validate and normalize the compose file
	err = validateAndNormalize(&compose)
	if err != nil {
		return nil, fmt.Errorf("invalid compose file: %w", err)
	}

	return &compose, nil
}

// validateAndNormalize validates and normalizes the compose file structure
func validateAndNormalize(compose *ComposeFile) error {
	// Validate services
	if len(compose.Services) == 0 {
		return fmt.Errorf("no services defined")
	}

	for serviceName, service := range compose.Services {
		// Validate that service has either image or build
		buildConfig := getBuildConfig(service.Build)
		if service.Image == "" && buildConfig.Context == "" {
			return fmt.Errorf("service '%s' must specify either 'image' or 'build'", serviceName)
		}

		// Normalize command and entrypoint to string slices
		normalizedService := service
		normalizedService.Command = normalizeStringSlice(service.Command)
		normalizedService.Entrypoint = normalizeStringSlice(service.Entrypoint)
		normalizedService.Environment = normalizeEnvironment(service.Environment)
		normalizedService.Networks = normalizeNetworks(service.Networks)
		normalizedService.Build = normalizeBuild(service.Build)
		normalizedService.Labels = normalizeLabels(service.Labels)

		compose.Services[serviceName] = normalizedService
	}

	// Normalize volumes
	for volumeName, volume := range compose.Volumes {
		normalizedVolume := volume
		normalizedVolume.Labels = normalizeLabels(volume.Labels)
		compose.Volumes[volumeName] = normalizedVolume
	}

	// Normalize networks
	for networkName, network := range compose.Networks {
		normalizedNetwork := network
		normalizedNetwork.Labels = normalizeLabels(network.Labels)
		compose.Networks[networkName] = normalizedNetwork
	}

	return nil
}

// normalizeStringSlice converts various formats to []string
func normalizeStringSlice(value interface{}) []string {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		// Split shell-style command
		return strings.Fields(v)
	case []interface{}:
		result := make([]string, len(v))
		for i, item := range v {
			result[i] = fmt.Sprintf("%v", item)
		}
		return result
	case []string:
		return v
	default:
		return []string{fmt.Sprintf("%v", value)}
	}
}

// normalizeEnvironment converts various environment formats to map[string]string
func normalizeEnvironment(value interface{}) map[string]string {
	if value == nil {
		return nil
	}

	result := make(map[string]string)

	switch v := value.(type) {
	case map[string]interface{}:
		for key, val := range v {
			result[key] = fmt.Sprintf("%v", val)
		}
	case map[interface{}]interface{}:
		for key, val := range v {
			result[fmt.Sprintf("%v", key)] = fmt.Sprintf("%v", val)
		}
	case []interface{}:
		for _, item := range v {
			if str, ok := item.(string); ok {
				if parts := strings.SplitN(str, "=", 2); len(parts) == 2 {
					result[parts[0]] = parts[1]
				} else {
					result[parts[0]] = ""
				}
			}
		}
	case []string:
		for _, item := range v {
			if parts := strings.SplitN(item, "=", 2); len(parts) == 2 {
				result[parts[0]] = parts[1]
			} else {
				result[parts[0]] = ""
			}
		}
	case map[string]string:
		return v
	}

	return result
}

// normalizeBuild normalizes build configuration from either string or object format
func normalizeBuild(build interface{}) BuildConfig {
	if build == nil {
		return BuildConfig{}
	}

	switch v := build.(type) {
	case string:
		// Build specified as string (e.g., "build: .")
		return BuildConfig{
			Context: v,
		}
	case map[interface{}]interface{}:
		// Build specified as object
		result := BuildConfig{}
		if context, ok := v["context"].(string); ok {
			result.Context = context
		}
		if dockerfile, ok := v["dockerfile"].(string); ok {
			result.Dockerfile = dockerfile
		}
		if args, ok := v["args"]; ok {
			result.Args = normalizeEnvironment(args)
		}
		if target, ok := v["target"].(string); ok {
			result.Target = target
		}
		return result
	default:
		return BuildConfig{}
	}
}

// getBuildConfig safely gets build config from interface
func getBuildConfig(build interface{}) BuildConfig {
	return normalizeBuild(build)
}

// normalizeLabels normalizes labels from either slice or map format
func normalizeLabels(labels interface{}) map[string]string {
	if labels == nil {
		return make(map[string]string)
	}

	switch v := labels.(type) {
	case []interface{}:
		// Labels specified as slice (e.g., ["key=value", "key2=value2"])
		result := make(map[string]string)
		for _, label := range v {
			if labelStr, ok := label.(string); ok {
				parts := strings.SplitN(labelStr, "=", 2)
				if len(parts) == 2 {
					result[parts[0]] = parts[1]
				} else {
					result[parts[0]] = ""
				}
			}
		}
		return result
	case []string:
		// Labels specified as string slice
		result := make(map[string]string)
		for _, label := range v {
			parts := strings.SplitN(label, "=", 2)
			if len(parts) == 2 {
				result[parts[0]] = parts[1]
			} else {
				result[parts[0]] = ""
			}
		}
		return result
	case map[interface{}]interface{}:
		// Labels specified as object
		result := make(map[string]string)
		for k, val := range v {
			if key, ok := k.(string); ok {
				if value, ok := val.(string); ok {
					result[key] = value
				}
			}
		}
		return result
	case map[string]interface{}:
		// Labels specified as string-keyed object
		result := make(map[string]string)
		for k, val := range v {
			if value, ok := val.(string); ok {
				result[k] = value
			}
		}
		return result
	case map[string]string:
		// Labels already in correct format
		return v
	default:
		return make(map[string]string)
	}
}

// normalizeNetworks converts various network formats to []string
func normalizeNetworks(value interface{}) []string {
	if value == nil {
		return []string{"default"}
	}

	switch v := value.(type) {
	case string:
		return []string{v}
	case []interface{}:
		result := make([]string, len(v))
		for i, item := range v {
			result[i] = fmt.Sprintf("%v", item)
		}
		return result
	case []string:
		return v
	case map[interface{}]interface{}:
		// Extract network names from map format
		var networks []string
		for key := range v {
			networks = append(networks, fmt.Sprintf("%v", key))
		}
		return networks
	default:
		return []string{"default"}
	}
}

// GetCommand returns the command as a string slice
func (s *ServiceConfig) GetCommand() []string {
	if cmd, ok := s.Command.([]string); ok {
		return cmd
	}
	return nil
}

// GetEntrypoint returns the entrypoint as a string slice
func (s *ServiceConfig) GetEntrypoint() []string {
	if ep, ok := s.Entrypoint.([]string); ok {
		return ep
	}
	return nil
}

// GetEnvironment returns the environment as a map
func (s *ServiceConfig) GetEnvironment() map[string]string {
	if env, ok := s.Environment.(map[string]string); ok {
		return env
	}
	return nil
}

// GetNetworks returns the networks as a string slice
func (s *ServiceConfig) GetNetworks() []string {
	if networks, ok := s.Networks.([]string); ok {
		return networks
	}
	return []string{"default"}
}

// ParsePortMapping parses port mapping strings like "8080:80" or "80"
func ParsePortMapping(portStr string) (hostPort, containerPort int, err error) {
	parts := strings.Split(portStr, ":")

	switch len(parts) {
	case 1:
		// Just container port, use same for host
		port, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid port: %s", parts[0])
		}
		return port, port, nil
	case 2:
		// Host:container format
		hostPort, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid host port: %s", parts[0])
		}
		containerPort, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid container port: %s", parts[1])
		}
		return hostPort, containerPort, nil
	default:
		return 0, 0, fmt.Errorf("invalid port mapping format: %s", portStr)
	}
}

// ParseVolumeMapping parses volume mapping strings like "/host/path:/container/path" or "volume_name:/container/path"
func ParseVolumeMapping(volumeStr string) (source, target, mode string, err error) {
	parts := strings.Split(volumeStr, ":")

	if len(parts) < 2 {
		return "", "", "", fmt.Errorf("invalid volume mapping format: %s", volumeStr)
	}

	source = parts[0]
	target = parts[1]
	mode = "rw" // default mode

	if len(parts) > 2 {
		mode = parts[2]
	}

	// Validate mode
	if mode != "rw" && mode != "ro" {
		return "", "", "", fmt.Errorf("invalid volume mode: %s (must be 'rw' or 'ro')", mode)
	}

	return source, target, mode, nil
}

// ResolveBuildContext resolves the build context relative to the compose file
func (b *BuildConfig) ResolveBuildContext(composeFileDir string) string {
	if b.Context == "" {
		return composeFileDir
	}

	if filepath.IsAbs(b.Context) {
		return b.Context
	}

	return filepath.Join(composeFileDir, b.Context)
}

// GetDockerfile returns the Dockerfile path, defaulting to "Buildfile"
func (b *BuildConfig) GetDockerfile() string {
	if b.Dockerfile == "" {
		return "Buildfile"
	}
	return b.Dockerfile
}
