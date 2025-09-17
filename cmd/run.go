package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"servin/pkg/container"
	"servin/pkg/network"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [flags] IMAGE COMMAND [ARG...]",
	Short: "Run a command in a new container",
	Long: `Create and run a new container from the specified image.
The container will be isolated using Linux namespaces and optionally
resource-limited using cgroups.`,
	Args: cobra.MinimumNArgs(2),
	RunE: runContainer,
}

var (
	// Container configuration flags
	containerName string
	memory        string
	cpus          string
	networkMode   string
	volumes       []string
	workdir       string
	env           []string
	hostname      string
	ports         []string
	detach        bool
)

func init() {
	rootCmd.AddCommand(runCmd)

	// Container configuration flags
	runCmd.Flags().StringVar(&containerName, "name", "", "Assign a name to the container")
	runCmd.Flags().StringVar(&memory, "memory", "", "Memory limit (e.g., 128m, 1g)")
	runCmd.Flags().StringVar(&cpus, "cpus", "", "CPU limit (e.g., 0.5, 2)")
	runCmd.Flags().StringVar(&networkMode, "network", "bridge", "Network mode (bridge, host, none)")
	runCmd.Flags().StringSliceVar(&volumes, "volume", []string{}, "Bind mount volumes (host:container)")
	runCmd.Flags().StringVar(&workdir, "workdir", "/", "Working directory inside container")
	runCmd.Flags().StringSliceVar(&env, "env", []string{}, "Set environment variables")
	runCmd.Flags().StringVar(&hostname, "hostname", "", "Container hostname")
	runCmd.Flags().StringSliceVarP(&ports, "publish", "p", []string{}, "Publish container ports (host:container or hostPort:containerPort/protocol)")
	runCmd.Flags().BoolVarP(&detach, "detach", "d", false, "Run container in background and print container ID")
}

func runContainer(cmd *cobra.Command, args []string) error {
	if err := checkRootForContainerOps(); err != nil {
		return err
	}

	image := args[0]
	command := args[1]
	commandArgs := args[2:]

	// Create container configuration
	config := &container.Config{
		Image:        image,
		Command:      command,
		Args:         commandArgs,
		Name:         containerName,
		WorkDir:      workdir,
		Hostname:     hostname,
		Env:          parseEnvVars(env),
		Volumes:      parseVolumes(volumes),
		NetworkMode:  networkMode,
		PortMappings: parsePortMappings(ports),
	}

	// Apply resource limits if specified
	if memory != "" {
		config.Memory = memory
	}
	if cpus != "" {
		config.CPUs = cpus
	}

	// Create and run the container
	c, err := container.New(config)
	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}

	verbose, _ := cmd.Flags().GetBool("verbose")
	if verbose {
		fmt.Printf("Starting container %s with image %s\n", c.ID, image)
		fmt.Printf("Command: %s %v\n", command, commandArgs)
	}

	if detach {
		// Run in background
		fmt.Printf("%s\n", c.ID)
		go func() {
			if err := c.RunWithVM(); err != nil {
				fmt.Printf("Container %s exited with error: %v\n", c.ID[:12], err)
			}
		}()
		return nil
	} else {
		// Show exit instructions for foreground runs
		fmt.Printf("Starting container... (Press Ctrl+C to exit)\n")
		return c.RunWithVM()
	}
}

// parseEnvVars parses environment variables from KEY=VALUE format
func parseEnvVars(envs []string) map[string]string {
	result := make(map[string]string)
	for _, env := range envs {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

// parseVolumes parses volume mounts from host:container format
func parseVolumes(vols []string) map[string]string {
	result := make(map[string]string)
	for _, vol := range vols {
		parts := strings.SplitN(vol, ":", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

// parsePortMappings parses port mappings from various formats
func parsePortMappings(portSpecs []string) []network.PortMapping {
	var mappings []network.PortMapping

	for _, spec := range portSpecs {
		mapping, err := parsePortMapping(spec)
		if err != nil {
			fmt.Printf("Warning: invalid port mapping '%s': %v\n", spec, err)
			continue
		}
		mappings = append(mappings, mapping)
	}

	return mappings
}

// parsePortMapping parses a single port mapping specification
// Supports formats: port, hostPort:containerPort, hostIP:hostPort:containerPort, hostPort:containerPort/protocol
func parsePortMapping(spec string) (network.PortMapping, error) {
	var mapping network.PortMapping

	// Split on '/' to separate protocol
	parts := strings.SplitN(spec, "/", 2)
	portPart := parts[0]
	protocol := "tcp"
	if len(parts) == 2 {
		protocol = strings.ToLower(parts[1])
	}

	// Split port part on ':'
	portParts := strings.Split(portPart, ":")

	switch len(portParts) {
	case 1:
		// Format: port (same for host and container)
		port, err := strconv.Atoi(portParts[0])
		if err != nil {
			return mapping, fmt.Errorf("invalid port number: %s", portParts[0])
		}
		mapping = network.PortMapping{
			HostPort:      port,
			ContainerPort: port,
			Protocol:      protocol,
		}
	case 2:
		// Format: hostPort:containerPort
		hostPort, err := strconv.Atoi(portParts[0])
		if err != nil {
			return mapping, fmt.Errorf("invalid host port: %s", portParts[0])
		}
		containerPort, err := strconv.Atoi(portParts[1])
		if err != nil {
			return mapping, fmt.Errorf("invalid container port: %s", portParts[1])
		}
		mapping = network.PortMapping{
			HostPort:      hostPort,
			ContainerPort: containerPort,
			Protocol:      protocol,
		}
	case 3:
		// Format: hostIP:hostPort:containerPort
		hostIP := portParts[0]
		hostPort, err := strconv.Atoi(portParts[1])
		if err != nil {
			return mapping, fmt.Errorf("invalid host port: %s", portParts[1])
		}
		containerPort, err := strconv.Atoi(portParts[2])
		if err != nil {
			return mapping, fmt.Errorf("invalid container port: %s", portParts[2])
		}
		mapping = network.PortMapping{
			HostIP:        hostIP,
			HostPort:      hostPort,
			ContainerPort: containerPort,
			Protocol:      protocol,
		}
	default:
		return mapping, fmt.Errorf("invalid port mapping format")
	}

	return mapping, nil
}
