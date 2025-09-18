package container

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"time"

	"servin/pkg/cgroups"
	"servin/pkg/namespaces"
	"servin/pkg/network"
	"servin/pkg/rootfs"
	"servin/pkg/state"
)

// Config represents container configuration
type Config struct {
	Image        string
	Command      string
	Args         []string
	Name         string
	WorkDir      string
	Hostname     string
	Env          map[string]string
	Volumes      map[string]string
	NetworkMode  string
	Memory       string
	CPUs         string
	PortMappings []network.PortMapping
}

// Container represents a running container
type Container struct {
	ID             string
	Config         *Config
	PID            int
	Status         string
	RootPath       string
	RootFS         *rootfs.RootFS
	CGroup         *cgroups.CGroup
	StateManager   *state.StateManager
	NetworkManager *network.NetworkManager
	ContainerNet   *network.ContainerNetwork
}

// New creates a new container with the given configuration
func New(config *Config) (*Container, error) {
	// Generate container ID
	id, err := generateID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate container ID: %v", err)
	}

	// Set default name if not provided
	if config.Name == "" {
		config.Name = id[:12]
	}

	// Set default hostname if not provided
	if config.Hostname == "" {
		config.Hostname = config.Name
	}

	// Create RootFS manager
	rfs := rootfs.New(id, config.Image)

	// Create CGroup manager
	cg := cgroups.New(id)

	// Create state manager
	sm := state.NewStateManager()

	// Create network manager
	nm := network.NewNetworkManager()

	container := &Container{
		ID:             id,
		Config:         config,
		Status:         "created",
		RootPath:       fmt.Sprintf("/var/lib/servin/containers/%s", id),
		RootFS:         rfs,
		CGroup:         cg,
		StateManager:   sm,
		NetworkManager: nm,
	}

	// Save initial container state
	if err := container.SaveState(); err != nil {
		return nil, fmt.Errorf("failed to save container state: %v", err)
	}

	return container, nil
}

// Run starts the container with namespace isolation, filesystem isolation, and resource limits
func (c *Container) Run() error {
	fmt.Printf("Running container %s (%s)\n", c.Config.Name, c.ID[:12])

	// Create the container's root filesystem
	if err := c.RootFS.Create(); err != nil {
		return fmt.Errorf("failed to create rootfs: %v", err)
	}

	// Setup filesystem mounts for the container
	if err := c.RootFS.SetupMounts(); err != nil {
		fmt.Printf("Warning: failed to setup mounts: %v\n", err)
	}

	// Prepare rootfs environment (sets SERVIN_ROOTFS env var)
	if err := c.RootFS.Enter(); err != nil {
		return fmt.Errorf("failed to prepare rootfs environment: %v", err)
	}

	// Create cgroups for resource control
	if err := c.CGroup.Create(); err != nil {
		fmt.Printf("Warning: failed to create cgroups: %v\n", err)
	} else {
		// Set resource limits if specified
		if c.Config.Memory != "" {
			if memBytes, err := cgroups.ParseMemoryString(c.Config.Memory); err == nil && memBytes > 0 {
				if err := c.CGroup.SetMemoryLimit(memBytes); err != nil {
					fmt.Printf("Warning: failed to set memory limit: %v\n", err)
				} else {
					fmt.Printf("Set memory limit to %d bytes\n", memBytes)
				}
			}
		}

		// Set default PID limit to prevent fork bombs
		if err := c.CGroup.SetPIDLimit(1024); err != nil {
			fmt.Printf("Warning: failed to set PID limit: %v\n", err)
		}
	}

	// Clean up on exit
	defer func() {
		if err := c.RootFS.Cleanup(); err != nil {
			fmt.Printf("Warning: failed to cleanup rootfs: %v\n", err)
		}
		if err := c.CGroup.Cleanup(); err != nil {
			fmt.Printf("Warning: failed to cleanup cgroups: %v\n", err)
		}
		if c.ContainerNet != nil {
			if err := c.NetworkManager.DetachContainerFromNetwork(c.ContainerNet); err != nil {
				fmt.Printf("Warning: failed to cleanup network: %v\n", err)
			}
		}
	}()

	// Set up networking if not in host mode
	if c.Config.NetworkMode != "host" && c.Config.NetworkMode != "none" {
		containerNet, err := c.NetworkManager.CreateVethPair(c.ID)
		if err != nil {
			fmt.Printf("Warning: failed to create network interface: %v\n", err)
		} else {
			c.ContainerNet = containerNet
			fmt.Printf("Created network interface for container\n")
		}
	}

	// Create log directory for container output
	sm := state.NewStateManager()
	logDir := filepath.Join(filepath.Dir(sm.GetStateDir()), "logs", c.ID)

	// Create namespace configuration
	nsConfig := &namespaces.ContainerConfig{
		Command:     c.Config.Command,
		Args:        c.Config.Args,
		Hostname:    c.Config.Hostname,
		WorkDir:     c.Config.WorkDir,
		LogDir:      logDir,
		RootFS:      c.RootPath + "/rootfs", // Pass the rootfs path
		Environment: c.Config.Env,           // Pass environment variables
		OnExit: func(err error) {
			// Update container status when process exits
			c.UpdateStatus("exited")
			if err != nil {
				fmt.Printf("Container %s exited with error: %v\n", c.Config.Name, err)
			} else {
				fmt.Printf("Container %s exited successfully\n", c.Config.Name)
			}
		},
		Namespaces: []namespaces.NamespaceFlags{
			namespaces.CLONE_NEWPID, // New PID namespace
			namespaces.CLONE_NEWUTS, // New UTS namespace (hostname)
			namespaces.CLONE_NEWIPC, // New IPC namespace
			namespaces.CLONE_NEWNS,  // New mount namespace
			namespaces.CLONE_NEWNET, // New network namespace
		},
	}

	c.Status = "running"
	c.UpdateStatus("running")

	err := namespaces.CreateContainer(nsConfig)

	if err != nil {
		c.UpdateStatus("exited")
		return fmt.Errorf("container failed: %v", err)
	}

	// Don't set to "exited" immediately - the container is now running
	// The status will be updated to "exited" when the process actually terminates
	return nil
}

// GetStats returns container resource usage statistics
func (c *Container) GetStats() (map[string]string, error) {
	return c.CGroup.GetStats()
}

// SaveState saves the container state to persistent storage
func (c *Container) SaveState() error {
	if c.StateManager == nil {
		return fmt.Errorf("state manager not initialized")
	}

	containerState := &state.ContainerState{
		ID:           c.ID,
		Name:         c.Config.Name,
		Image:        c.Config.Image,
		Command:      c.Config.Command,
		Args:         c.Config.Args,
		Status:       c.Status,
		PID:          c.PID,
		Created:      time.Now(),
		RootPath:     c.RootPath,
		Hostname:     c.Config.Hostname,
		WorkDir:      c.Config.WorkDir,
		Env:          c.Config.Env,
		Volumes:      c.Config.Volumes,
		NetworkMode:  c.Config.NetworkMode,
		PortMappings: c.Config.PortMappings,
		Memory:       c.Config.Memory,
		CPUs:         c.Config.CPUs,
	}

	return c.StateManager.SaveContainer(containerState)
}

// UpdateStatus updates the container status in persistent storage
func (c *Container) UpdateStatus(status string) error {
	c.Status = status
	if c.StateManager != nil {
		return c.StateManager.UpdateContainerStatus(c.ID, status)
	}
	return nil
}

// UpdatePID updates the container PID in persistent storage
func (c *Container) UpdatePID(pid int) error {
	c.PID = pid
	if c.StateManager != nil {
		return c.StateManager.UpdateContainerPID(c.ID, pid)
	}
	return nil
}

// generateID creates a random container ID
func generateID() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
