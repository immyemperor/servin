package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"servin/pkg/compose"
	"servin/pkg/logger"

	"github.com/spf13/cobra"
)

var composeCmd = &cobra.Command{
	Use:   "compose",
	Short: "Manage multi-container applications with Servin Compose",
	Long: `Servin Compose allows you to define and run multi-container applications.
Use a servin-compose.yml file to configure your application's services, networks, and volumes.

Available subcommands:
  up     - Create and start services
  down   - Stop and remove services
  ps     - List running services
  logs   - View output from services
  exec   - Execute a command in a running service`,
}

var composeUpCmd = &cobra.Command{
	Use:   "up [OPTIONS]",
	Short: "Create and start services",
	Long: `Create and start services defined in servin-compose.yml.
This command will:
1. Create any missing networks and volumes
2. Build or pull required images
3. Create and start containers for all services
4. Attach to service logs (unless -d is specified)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		detach, _ := cmd.Flags().GetBool("detach")
		file, _ := cmd.Flags().GetString("file")
		projectName, _ := cmd.Flags().GetString("project-name")

		return runComposeUp(file, projectName, detach)
	},
}

var composeDownCmd = &cobra.Command{
	Use:   "down [OPTIONS]",
	Short: "Stop and remove services",
	Long: `Stop and remove containers, networks, and volumes created by 'up'.
By default, only containers and networks are removed. Use --volumes to also remove volumes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		projectName, _ := cmd.Flags().GetString("project-name")
		removeVolumes, _ := cmd.Flags().GetBool("volumes")

		return runComposeDown(file, projectName, removeVolumes)
	},
}

var composePsCmd = &cobra.Command{
	Use:   "ps [OPTIONS]",
	Short: "List running services",
	Long:  `List containers for services defined in the compose file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		projectName, _ := cmd.Flags().GetString("project-name")
		all, _ := cmd.Flags().GetBool("all")

		return runComposePs(file, projectName, all)
	},
}

var composeLogsCmd = &cobra.Command{
	Use:   "logs [OPTIONS] [SERVICE...]",
	Short: "View output from services",
	Long:  `Display log output from services. If no services are specified, logs from all services are shown.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		projectName, _ := cmd.Flags().GetString("project-name")
		follow, _ := cmd.Flags().GetBool("follow")
		timestamps, _ := cmd.Flags().GetBool("timestamps")
		tail, _ := cmd.Flags().GetString("tail")

		return runComposeLogs(file, projectName, args, follow, timestamps, tail)
	},
}

var composeExecCmd = &cobra.Command{
	Use:   "exec [OPTIONS] SERVICE COMMAND [ARG...]",
	Short: "Execute a command in a running service",
	Long:  `Execute a command in a running service container.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("exec requires at least a service name and command")
		}

		file, _ := cmd.Flags().GetString("file")
		projectName, _ := cmd.Flags().GetString("project-name")
		interactive, _ := cmd.Flags().GetBool("interactive")

		serviceName := args[0]
		command := args[1:]

		return runComposeExec(file, projectName, serviceName, command, interactive)
	},
}

func init() {
	// Add compose command to root
	rootCmd.AddCommand(composeCmd)

	// Add subcommands
	composeCmd.AddCommand(composeUpCmd)
	composeCmd.AddCommand(composeDownCmd)
	composeCmd.AddCommand(composePsCmd)
	composeCmd.AddCommand(composeLogsCmd)
	composeCmd.AddCommand(composeExecCmd)

	// Global compose flags
	composeCmd.PersistentFlags().StringP("file", "f", "servin-compose.yml", "Specify an alternate compose file")
	composeCmd.PersistentFlags().StringP("project-name", "p", "", "Specify an alternate project name (default: directory name)")

	// Up command flags
	composeUpCmd.Flags().BoolP("detach", "d", false, "Detached mode: Run containers in the background")

	// Down command flags
	composeDownCmd.Flags().Bool("volumes", false, "Remove named volumes declared in the volumes section")

	// Ps command flags
	composePsCmd.Flags().BoolP("all", "a", false, "Show all containers (default shows just running)")

	// Logs command flags
	composeLogsCmd.Flags().Bool("follow", false, "Follow log output")
	composeLogsCmd.Flags().BoolP("timestamps", "t", false, "Show timestamps")
	composeLogsCmd.Flags().String("tail", "all", "Number of lines to show from the end of the logs")

	// Exec command flags
	composeExecCmd.Flags().BoolP("interactive", "i", false, "Keep STDIN open even if not attached")
}

func runComposeUp(file, projectName string, detach bool) error {
	// Resolve compose file path
	composeFile, err := resolveComposeFile(file)
	if err != nil {
		return err
	}

	// Determine project name
	if projectName == "" {
		projectName = getDefaultProjectName(composeFile)
	}

	logger.Info("Starting compose project %s from file %s", projectName, composeFile)

	// Parse compose file
	project, err := compose.LoadProject(composeFile, projectName)
	if err != nil {
		return fmt.Errorf("failed to load compose file: %w", err)
	}

	// Create and start services
	err = project.Up(detach)
	if err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}

	if detach {
		fmt.Printf("Started %d services for project %s\n", len(project.Services), projectName)
	}

	return nil
}

func runComposeDown(file, projectName string, removeVolumes bool) error {
	// Resolve compose file path
	composeFile, err := resolveComposeFile(file)
	if err != nil {
		return err
	}

	// Determine project name
	if projectName == "" {
		projectName = getDefaultProjectName(composeFile)
	}

	logger.Info("Stopping compose project %s from file %s", projectName, composeFile)

	// Parse compose file
	project, err := compose.LoadProject(composeFile, projectName)
	if err != nil {
		return fmt.Errorf("failed to load compose file: %w", err)
	}

	// Stop and remove services
	err = project.Down(removeVolumes)
	if err != nil {
		return fmt.Errorf("failed to stop services: %w", err)
	}

	fmt.Printf("Stopped and removed services for project %s\n", projectName)
	return nil
}

func runComposePs(file, projectName string, all bool) error {
	// Resolve compose file path
	composeFile, err := resolveComposeFile(file)
	if err != nil {
		return err
	}

	// Determine project name
	if projectName == "" {
		projectName = getDefaultProjectName(composeFile)
	}

	// Parse compose file
	project, err := compose.LoadProject(composeFile, projectName)
	if err != nil {
		return fmt.Errorf("failed to load compose file: %w", err)
	}

	// List services
	return project.Ps(all)
}

func runComposeLogs(file, projectName string, services []string, follow, timestamps bool, tail string) error {
	// Resolve compose file path
	composeFile, err := resolveComposeFile(file)
	if err != nil {
		return err
	}

	// Determine project name
	if projectName == "" {
		projectName = getDefaultProjectName(composeFile)
	}

	// Parse compose file
	project, err := compose.LoadProject(composeFile, projectName)
	if err != nil {
		return fmt.Errorf("failed to load compose file: %w", err)
	}

	// Show logs
	return project.Logs(services, follow, timestamps, tail)
}

func runComposeExec(file, projectName, serviceName string, command []string, interactive bool) error {
	// Resolve compose file path
	composeFile, err := resolveComposeFile(file)
	if err != nil {
		return err
	}

	// Determine project name
	if projectName == "" {
		projectName = getDefaultProjectName(composeFile)
	}

	// Parse compose file
	project, err := compose.LoadProject(composeFile, projectName)
	if err != nil {
		return fmt.Errorf("failed to load compose file: %w", err)
	}

	// Execute command
	return project.Exec(serviceName, command, interactive)
}

func resolveComposeFile(file string) (string, error) {
	// If absolute path, use as-is
	if filepath.IsAbs(file) {
		if _, err := os.Stat(file); err != nil {
			return "", fmt.Errorf("compose file not found: %s", file)
		}
		return file, nil
	}

	// Try current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	composeFile := filepath.Join(currentDir, file)
	if _, err := os.Stat(composeFile); err != nil {
		return "", fmt.Errorf("compose file not found: %s", composeFile)
	}

	return composeFile, nil
}

func getDefaultProjectName(composeFile string) string {
	// Use directory name as default project name
	dir := filepath.Dir(composeFile)
	return filepath.Base(dir)
}
