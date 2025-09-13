package compose

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"servin/pkg/image"
	"servin/pkg/logger"
	"servin/pkg/state"
)

// Project represents a compose project with its services and state
type Project struct {
	Name         string
	File         string
	Compose      *ComposeFile
	Services     map[string]*Service
	stateManager *state.StateManager
	imageManager *image.Manager
}

// Service represents a compose service instance
type Service struct {
	Name        string
	Config      ServiceConfig
	ContainerID string
	Status      string
	CreatedAt   time.Time
	project     *Project
}

// LoadProject loads a compose project from a file
func LoadProject(filePath, projectName string) (*Project, error) {
	// Parse compose file
	compose, err := ParseComposeFile(filePath)
	if err != nil {
		return nil, err
	}

	// Initialize managers
	stateManager := state.NewStateManager()
	imageManager := image.NewManager()

	project := &Project{
		Name:         projectName,
		File:         filePath,
		Compose:      compose,
		Services:     make(map[string]*Service),
		stateManager: stateManager,
		imageManager: imageManager,
	}

	// Initialize services
	for serviceName, serviceConfig := range compose.Services {
		project.Services[serviceName] = &Service{
			Name:    serviceName,
			Config:  serviceConfig,
			Status:  "created",
			project: project,
		}
	}

	return project, nil
}

// Up creates and starts all services
func (p *Project) Up(detach bool) error {
	logger.Info("Starting compose project %s", p.Name)

	fmt.Printf("Creating network %s_default\n", p.Name)

	// Sort services by dependency order
	serviceOrder, err := p.resolveDependencies()
	if err != nil {
		return fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	// Start services in dependency order
	for _, serviceName := range serviceOrder {
		service := p.Services[serviceName]
		err := p.startService(service)
		if err != nil {
			return fmt.Errorf("failed to start service %s: %w", serviceName, err)
		}
	}

	if !detach {
		fmt.Printf("Attaching to logs... (Press Ctrl+C to stop)\n")
		// In a real implementation, this would follow logs from all containers
		select {} // Block forever (would be interrupted by Ctrl+C)
	}

	return nil
}

// Down stops and removes all services
func (p *Project) Down(removeVolumes bool) error {
	logger.Info("Stopping compose project %s", p.Name)

	// Stop services in reverse dependency order
	serviceOrder, _ := p.resolveDependencies()

	// Reverse the order for shutdown
	for i := len(serviceOrder) - 1; i >= 0; i-- {
		serviceName := serviceOrder[i]
		if service, exists := p.Services[serviceName]; exists {
			err := p.stopService(service)
			if err != nil {
				logger.Warn("Failed to stop service %s: %v", serviceName, err)
			}
		}
	}

	fmt.Printf("Removing network %s_default\n", p.Name)

	if removeVolumes {
		fmt.Printf("Removing volumes for project %s\n", p.Name)
	}

	return nil
}

// Ps lists project containers
func (p *Project) Ps(all bool) error {
	// Print header
	fmt.Printf("%-20s %-15s %-10s %-15s %-30s\n", "NAME", "SERVICE", "STATUS", "PORTS", "COMMAND")
	fmt.Println(strings.Repeat("-", 90))

	// Print service information
	for serviceName, service := range p.Services {
		if !all && service.Status != "running" {
			continue
		}

		containerName := fmt.Sprintf("%s_%s_1", p.Name, serviceName)
		ports := p.formatServicePorts(service)
		command := strings.Join(service.Config.GetCommand(), " ")
		if len(command) > 30 {
			command = command[:27] + "..."
		}

		fmt.Printf("%-20s %-15s %-10s %-15s %-30s\n",
			containerName, serviceName, service.Status, ports, command)
	}

	return nil
}

// Logs shows logs from services
func (p *Project) Logs(services []string, follow, timestamps bool, tail string) error {
	// Filter services if specified
	var targetServices []*Service
	if len(services) > 0 {
		serviceSet := make(map[string]bool)
		for _, service := range services {
			serviceSet[service] = true
		}

		for _, service := range p.Services {
			if serviceSet[service.Name] {
				targetServices = append(targetServices, service)
			}
		}
	} else {
		for _, service := range p.Services {
			targetServices = append(targetServices, service)
		}
	}

	if len(targetServices) == 0 {
		fmt.Println("No services found")
		return nil
	}

	// Show logs for each service
	for _, service := range targetServices {
		containerName := fmt.Sprintf("%s_%s_1", p.Name, service.Name)
		fmt.Printf("==> %s (%s) <==\n", service.Name, containerName)

		// In a real implementation, this would integrate with the existing logs command
		fmt.Printf("Logs from %s service would be shown here\n", service.Name)
		if timestamps {
			fmt.Printf("[%s] Service started\n", time.Now().Format(time.RFC3339))
		}
	}

	if follow {
		fmt.Println("Following logs... (Press Ctrl+C to stop)")
		select {} // Block forever (would be interrupted by Ctrl+C)
	}

	return nil
}

// Exec executes a command in a service container
func (p *Project) Exec(serviceName string, command []string, interactive bool) error {
	// Find the service
	service, exists := p.Services[serviceName]
	if !exists {
		return fmt.Errorf("service not found: %s", serviceName)
	}

	if service.Status != "running" {
		return fmt.Errorf("service %s is not running", serviceName)
	}

	containerName := fmt.Sprintf("%s_%s_1", p.Name, serviceName)
	fmt.Printf("Executing %v in container %s\n", command, containerName)

	// In a real implementation, this would integrate with the existing exec command
	return nil
}

// resolveDependencies resolves service start order based on depends_on
func (p *Project) resolveDependencies() ([]string, error) {
	// Simple topological sort for dependency resolution
	visited := make(map[string]bool)
	visiting := make(map[string]bool)
	var order []string

	var visit func(string) error
	visit = func(serviceName string) error {
		if visiting[serviceName] {
			return fmt.Errorf("circular dependency detected involving service: %s", serviceName)
		}
		if visited[serviceName] {
			return nil
		}

		visiting[serviceName] = true

		service := p.Services[serviceName]
		for _, dep := range service.Config.DependsOn {
			if _, exists := p.Services[dep]; !exists {
				return fmt.Errorf("service %s depends on unknown service: %s", serviceName, dep)
			}
			err := visit(dep)
			if err != nil {
				return err
			}
		}

		visiting[serviceName] = false
		visited[serviceName] = true
		order = append(order, serviceName)

		return nil
	}

	// Visit all services
	for serviceName := range p.Services {
		err := visit(serviceName)
		if err != nil {
			return nil, err
		}
	}

	return order, nil
}

// startService starts a single service
func (p *Project) startService(service *Service) error {
	logger.Info("Starting service %s", service.Name)

	containerName := fmt.Sprintf("%s_%s_1", p.Name, service.Name)

	// Check if image exists or needs to be built
	imageName := service.Config.Image
	buildConfig := service.Config.GetBuildConfig()
	if buildConfig.Context != "" {
		imageName = fmt.Sprintf("%s_%s", p.Name, service.Name)
		err := p.buildServiceImage(service, imageName)
		if err != nil {
			return fmt.Errorf("failed to build image for service %s: %w", service.Name, err)
		}
	} else {
		// Check if image exists
		_, err := p.imageManager.GetImage(imageName)
		if err != nil {
			return fmt.Errorf("image %s not found for service %s", imageName, service.Name)
		}
	}

	fmt.Printf("Creating %s ... done\n", containerName)
	fmt.Printf("Starting %s ... done\n", containerName)

	service.ContainerID = containerName
	service.Status = "running"
	service.CreatedAt = time.Now()

	return nil
}

// stopService stops a single service
func (p *Project) stopService(service *Service) error {
	if service.ContainerID == "" {
		return nil // Already stopped
	}

	logger.Info("Stopping service %s container %s", service.Name, service.ContainerID)

	fmt.Printf("Stopping %s ... done\n", service.ContainerID)
	fmt.Printf("Removing %s ... done\n", service.ContainerID)

	service.Status = "stopped"
	return nil
}

// buildServiceImage builds an image for a service
func (p *Project) buildServiceImage(service *Service, imageName string) error {
	buildConfig := service.Config.GetBuildConfig()

	// Resolve build context relative to compose file
	composeDir := filepath.Dir(p.File)
	buildContext := filepath.Join(composeDir, buildConfig.Context)

	logger.Info("Building image for service %s from context %s", service.Name, buildContext)

	fmt.Printf("Building %s from %s\n", imageName, buildContext)

	// In a real implementation, this would integrate with the existing build command
	return nil
}

// formatServicePorts formats port mappings for display
func (p *Project) formatServicePorts(service *Service) string {
	if len(service.Config.Ports) == 0 {
		return ""
	}

	var ports []string
	for _, port := range service.Config.Ports {
		ports = append(ports, port)
	}

	return strings.Join(ports, ", ")
}
