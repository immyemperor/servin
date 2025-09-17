package container

import (
	"fmt"
	"runtime"
)

// ContainerBackend represents different containerization approaches
type ContainerBackend string

const (
	// NativeBackend uses platform-native containerization (Linux namespaces)
	NativeBackend ContainerBackend = "native"
	// SimulatedBackend uses filesystem isolation only (current macOS implementation)
	SimulatedBackend ContainerBackend = "simulated"
	// DockerBackend delegates to Docker Desktop for true containerization
	DockerBackend ContainerBackend = "docker"
	// PodmanBackend delegates to Podman for containerization
	PodmanBackend ContainerBackend = "podman"
	// LimaBackend uses Lima VM for Linux containers on macOS
	LimaBackend ContainerBackend = "lima"
)

// BackendCapabilities describes what each backend can provide
type BackendCapabilities struct {
	ProcessIsolation     bool
	NetworkIsolation     bool
	FilesystemIsolation  bool
	ResourceLimits       bool
	TrueContainerization bool
	Platform             string
}

// GetAvailableBackends returns container backends available on current platform
func GetAvailableBackends() map[ContainerBackend]BackendCapabilities {
	backends := make(map[ContainerBackend]BackendCapabilities)

	switch runtime.GOOS {
	case "linux":
		backends[NativeBackend] = BackendCapabilities{
			ProcessIsolation:     true,
			NetworkIsolation:     true,
			FilesystemIsolation:  true,
			ResourceLimits:       true,
			TrueContainerization: true,
			Platform:             "linux",
		}

	case "darwin": // macOS
		backends[SimulatedBackend] = BackendCapabilities{
			ProcessIsolation:     false,
			NetworkIsolation:     false,
			FilesystemIsolation:  true, // Via VFS
			ResourceLimits:       false,
			TrueContainerization: false,
			Platform:             "macOS",
		}

		// Check if Docker Desktop is available
		if isDockerAvailable() {
			backends[DockerBackend] = BackendCapabilities{
				ProcessIsolation:     true,
				NetworkIsolation:     true,
				FilesystemIsolation:  true,
				ResourceLimits:       true,
				TrueContainerization: true,
				Platform:             "macOS (via Docker Desktop)",
			}
		}

		// Check if Podman is available
		if isPodmanAvailable() {
			backends[PodmanBackend] = BackendCapabilities{
				ProcessIsolation:     true,
				NetworkIsolation:     true,
				FilesystemIsolation:  true,
				ResourceLimits:       true,
				TrueContainerization: true,
				Platform:             "macOS (via Podman)",
			}
		}

		// Check if Lima is available
		if isLimaAvailable() {
			backends[LimaBackend] = BackendCapabilities{
				ProcessIsolation:     true,
				NetworkIsolation:     true,
				FilesystemIsolation:  true,
				ResourceLimits:       true,
				TrueContainerization: true,
				Platform:             "macOS (via Lima VM)",
			}
		}

	case "windows":
		backends[SimulatedBackend] = BackendCapabilities{
			ProcessIsolation:     false,
			NetworkIsolation:     false,
			FilesystemIsolation:  true,
			ResourceLimits:       false,
			TrueContainerization: false,
			Platform:             "Windows",
		}

		if isDockerAvailable() {
			backends[DockerBackend] = BackendCapabilities{
				ProcessIsolation:     true,
				NetworkIsolation:     true,
				FilesystemIsolation:  true,
				ResourceLimits:       true,
				TrueContainerization: true,
				Platform:             "Windows (via Docker Desktop)",
			}
		}
	}

	return backends
}

// RecommendBackend suggests the best backend for the current platform
func RecommendBackend() (ContainerBackend, string) {
	backends := GetAvailableBackends()

	// Prefer true containerization when available
	for backend, caps := range backends {
		if caps.TrueContainerization {
			return backend, fmt.Sprintf("Recommended: %s provides true containerization with %s", backend, caps.Platform)
		}
	}

	// Fall back to simulation
	for backend, caps := range backends {
		if !caps.TrueContainerization {
			return backend, fmt.Sprintf("Using: %s (filesystem isolation only) on %s", backend, caps.Platform)
		}
	}

	return SimulatedBackend, "Using simulated containerization (limited capabilities)"
}

// Helper functions to check runtime availability
func isDockerAvailable() bool {
	// Implementation would check if docker command is available
	return false // Placeholder
}

func isPodmanAvailable() bool {
	// Implementation would check if podman command is available
	return false // Placeholder
}

func isLimaAvailable() bool {
	// Implementation would check if limactl command is available
	return false // Placeholder
}
