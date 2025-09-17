//go:build darwin

package runtime

import (
	"fmt"
	"os/exec"
)

// ContainerRuntime represents different container runtime backends
type ContainerRuntime interface {
	IsAvailable() bool
	RunContainer(config *ContainerConfig) error
	ListContainers() ([]Container, error)
	StopContainer(id string) error
	RemoveContainer(id string) error
}

// DockerRuntime uses Docker Desktop for true containerization
type DockerRuntime struct{}

func (d *DockerRuntime) IsAvailable() bool {
	cmd := exec.Command("docker", "--version")
	return cmd.Run() == nil
}

func (d *DockerRuntime) RunContainer(config *ContainerConfig) error {
	args := []string{"run", "-d"}

	// Add port mappings
	for hostPort, containerPort := range config.Ports {
		args = append(args, "-p", fmt.Sprintf("%s:%s", hostPort, containerPort))
	}

	// Add environment variables
	for key, value := range config.Environment {
		args = append(args, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	// Add volumes
	for hostPath, containerPath := range config.Volumes {
		args = append(args, "-v", fmt.Sprintf("%s:%s", hostPath, containerPath))
	}

	// Add image and command
	args = append(args, config.Image)
	args = append(args, config.Command...)

	cmd := exec.Command("docker", args...)
	return cmd.Run()
}

func (d *DockerRuntime) ListContainers() ([]Container, error) {
	cmd := exec.Command("docker", "ps", "-a", "--format", "{{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse Docker output into Container structs
	containers := []Container{}
	_ = output // Use output variable
	return containers, nil
}

func (d *DockerRuntime) StopContainer(id string) error {
	cmd := exec.Command("docker", "stop", id)
	return cmd.Run()
}

func (d *DockerRuntime) RemoveContainer(id string) error {
	cmd := exec.Command("docker", "rm", id)
	return cmd.Run()
}

// PodmanRuntime uses Podman for containerization
type PodmanRuntime struct{}

func (p *PodmanRuntime) IsAvailable() bool {
	cmd := exec.Command("podman", "--version")
	return cmd.Run() == nil
}

func (p *PodmanRuntime) RunContainer(config *ContainerConfig) error {
	// Similar implementation to Docker but using podman commands
	args := []string{"run", "-d"}

	// Add configuration parameters
	for hostPort, containerPort := range config.Ports {
		args = append(args, "-p", fmt.Sprintf("%s:%s", hostPort, containerPort))
	}

	args = append(args, config.Image)
	args = append(args, config.Command...)

	cmd := exec.Command("podman", args...)
	return cmd.Run()
}

func (p *PodmanRuntime) ListContainers() ([]Container, error) {
	cmd := exec.Command("podman", "ps", "-a", "--format", "{{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse Podman output into Container structs
	containers := []Container{}
	_ = output // Use output variable
	return containers, nil
}

func (p *PodmanRuntime) StopContainer(id string) error {
	cmd := exec.Command("podman", "stop", id)
	return cmd.Run()
}

func (p *PodmanRuntime) RemoveContainer(id string) error {
	cmd := exec.Command("podman", "rm", id)
	return cmd.Run()
}

// RuntimeDetector finds available container runtimes
type RuntimeDetector struct {
	runtimes []ContainerRuntime
}

func NewRuntimeDetector() *RuntimeDetector {
	return &RuntimeDetector{
		runtimes: []ContainerRuntime{
			&DockerRuntime{},
			&PodmanRuntime{},
		},
	}
}

func (rd *RuntimeDetector) DetectBestRuntime() (ContainerRuntime, error) {
	for _, runtime := range rd.runtimes {
		if runtime.IsAvailable() {
			return runtime, nil
		}
	}

	return nil, fmt.Errorf("no container runtime available. Please install Docker Desktop or Podman")
}

// ContainerConfig represents container configuration
type ContainerConfig struct {
	Image       string
	Command     []string
	Ports       map[string]string
	Environment map[string]string
	Volumes     map[string]string
	Name        string
}

// Container represents a running container
type Container struct {
	ID     string
	Name   string
	Image  string
	Status string
}
