package cmd

import (
	"fmt"
	"os"
	"strings"

	"servin/pkg/logger"
	"servin/pkg/registry"

	"github.com/spf13/cobra"
)

var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Manage local and remote registries",
	Long: `Registry commands allow you to:
- Start a local registry server
- Push images to registries  
- Pull images from registries
- Login/logout from registries
- List available registries and images`,
}

var startRegistryCmd = &cobra.Command{
	Use:   "start",
	Short: "Start local registry server",
	Long: `Start a local registry server that can store and serve images.
This creates a Docker Registry-compatible HTTP API server running locally.

Examples:
  servin registry start                    # Start on default port 5000
  servin registry start --port 5001       # Start on custom port
  servin registry start --data-dir ./reg  # Use custom data directory`,
	RunE: runStartRegistry,
}

var stopRegistryCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop local registry server",
	Long:  `Stop the running local registry server.`,
	RunE:  runStopRegistry,
}

var pushCmd = &cobra.Command{
	Use:   "push IMAGE[:TAG] [REGISTRY_URL]",
	Short: "Push an image to a registry",
	Long: `Push a local image to a registry (local or remote).

Examples:
  servin registry push myapp:latest                    # Push to default registry
  servin registry push myapp:v1.0 localhost:5001      # Push to specific registry
  servin registry push myapp:latest docker.io/user/   # Push to Docker Hub`,
	Args: cobra.MinimumNArgs(1),
	RunE: runPush,
}

var pullCmd = &cobra.Command{
	Use:   "pull IMAGE[:TAG] [REGISTRY_URL]",
	Short: "Pull an image from a registry",
	Long: `Pull an image from a registry (local or remote) and store it locally.

Examples:
  servin registry pull myapp:latest                    # Pull from default registry
  servin registry pull myapp:v1.0 localhost:5001      # Pull from specific registry
  servin registry pull nginx:alpine docker.io         # Pull from Docker Hub`,
	Args: cobra.MinimumNArgs(1),
	RunE: runPull,
}

var loginCmd = &cobra.Command{
	Use:   "login REGISTRY_URL",
	Short: "Login to a registry",
	Long: `Authenticate with a registry to enable push/pull operations.

Examples:
  servin registry login docker.io
  servin registry login localhost:5001`,
	Args: cobra.ExactArgs(1),
	RunE: runLogin,
}

var logoutCmd = &cobra.Command{
	Use:   "logout REGISTRY_URL",
	Short: "Logout from a registry",
	Long: `Remove authentication credentials for a registry.

Examples:
  servin registry logout docker.io
  servin registry logout localhost:5001`,
	Args: cobra.ExactArgs(1),
	RunE: runLogout,
}

var registryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured registries and their status",
	Long:  `Display information about all configured registries including status and available images.`,
	RunE:  runRegistryList,
}

func init() {
	rootCmd.AddCommand(registryCmd)

	registryCmd.AddCommand(startRegistryCmd)
	registryCmd.AddCommand(stopRegistryCmd)
	registryCmd.AddCommand(pushCmd)
	registryCmd.AddCommand(pullCmd)
	registryCmd.AddCommand(loginCmd)
	registryCmd.AddCommand(logoutCmd)
	registryCmd.AddCommand(registryListCmd)

	// Start registry flags
	startRegistryCmd.Flags().Int("port", 5000, "Port for the registry server")
	startRegistryCmd.Flags().String("data-dir", "", "Data directory for registry storage")
	startRegistryCmd.Flags().Bool("detach", false, "Run registry in background")

	// Push flags
	pushCmd.Flags().Bool("force", false, "Force push even if image exists")
	pushCmd.Flags().BoolP("quiet", "q", false, "Suppress output")

	// Pull flags
	pullCmd.Flags().BoolP("quiet", "q", false, "Suppress output")
	pullCmd.Flags().String("platform", "", "Set platform if server is multi-platform capable")

	// Login flags
	loginCmd.Flags().StringP("username", "u", "", "Username for authentication")
	loginCmd.Flags().StringP("password", "p", "", "Password for authentication")
	loginCmd.Flags().String("email", "", "Email for authentication")
}

func runStartRegistry(cmd *cobra.Command, args []string) error {
	port, _ := cmd.Flags().GetInt("port")
	dataDir, _ := cmd.Flags().GetString("data-dir")
	detach, _ := cmd.Flags().GetBool("detach")

	if dataDir == "" {
		dataDir = getRegistryDataDir()
	}

	localRegistry := registry.NewLocalRegistry(dataDir, port)

	if detach {
		// TODO: Implement background/daemon mode
		fmt.Printf("Starting registry server in background on port %d\n", port)
		fmt.Printf("Data directory: %s\n", dataDir)
		return fmt.Errorf("background mode not yet implemented")
	}

	fmt.Printf("Starting registry server on port %d\n", port)
	fmt.Printf("Data directory: %s\n", dataDir)
	fmt.Printf("Registry URL: http://localhost:%d\n", port)
	fmt.Println("Press Ctrl+C to stop")

	return localRegistry.Start()
}

func runStopRegistry(cmd *cobra.Command, args []string) error {
	// TODO: Implement registry stop functionality
	return fmt.Errorf("registry stop not yet implemented")
}

func runPush(cmd *cobra.Command, args []string) error {
	imageArg := args[0]
	var registryURL string
	if len(args) > 1 {
		registryURL = args[1]
	}

	// Parse image name and tag
	imageName, tag := parseImageTag(imageArg)

	// Get flags
	force, _ := cmd.Flags().GetBool("force")
	quiet, _ := cmd.Flags().GetBool("quiet")

	// Create registry client
	client, err := registry.NewClient(getRegistryDataDir())
	if err != nil {
		return fmt.Errorf("failed to create registry client: %w", err)
	}

	options := &registry.PushOptions{
		Registry: registryURL,
		Force:    force,
		Quiet:    quiet,
	}

	return client.PushImage(imageName, tag, registryURL, options)
}

func runPull(cmd *cobra.Command, args []string) error {
	imageArg := args[0]
	var registryURL string
	if len(args) > 1 {
		registryURL = args[1]
	}

	// Parse image name and tag
	imageName, tag := parseImageTag(imageArg)

	// Get flags
	quiet, _ := cmd.Flags().GetBool("quiet")
	platform, _ := cmd.Flags().GetString("platform")

	// Create registry client
	client, err := registry.NewClient(getRegistryDataDir())
	if err != nil {
		return fmt.Errorf("failed to create registry client: %w", err)
	}

	options := &registry.PullOptions{
		Registry: registryURL,
		Quiet:    quiet,
		Platform: platform,
	}

	return client.PullImage(imageName, tag, registryURL, options)
}

func runLogin(cmd *cobra.Command, args []string) error {
	registryURL := args[0]

	// Get credentials from flags or prompt
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")
	email, _ := cmd.Flags().GetString("email")

	if username == "" {
		fmt.Print("Username: ")
		fmt.Scanln(&username)
	}

	if password == "" {
		fmt.Print("Password: ")
		// TODO: Use a proper password input library to hide input
		fmt.Scanln(&password)
	}

	// Create registry client
	client, err := registry.NewClient(getRegistryDataDir())
	if err != nil {
		return fmt.Errorf("failed to create registry client: %w", err)
	}

	if err := client.LoginToRegistry(registryURL, username, password, email); err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	fmt.Printf("Login succeeded for %s\n", registryURL)
	return nil
}

func runLogout(cmd *cobra.Command, args []string) error {
	registryURL := args[0]

	// Create registry client
	client, err := registry.NewClient(getRegistryDataDir())
	if err != nil {
		return fmt.Errorf("failed to create registry client: %w", err)
	}

	if err := client.LogoutFromRegistry(registryURL); err != nil {
		return fmt.Errorf("logout failed: %w", err)
	}

	fmt.Printf("Logout succeeded for %s\n", registryURL)
	return nil
}

func runRegistryList(cmd *cobra.Command, args []string) error {
	// Create registry client
	client, err := registry.NewClient(getRegistryDataDir())
	if err != nil {
		return fmt.Errorf("failed to create registry client: %w", err)
	}

	registries, err := client.GetRegistryInfo()
	if err != nil {
		return fmt.Errorf("failed to get registry information: %w", err)
	}

	// Display registry information
	fmt.Println("REGISTRY NAME        TYPE     STATUS      URL")
	fmt.Println("----------------------------------------------------")

	for _, reg := range registries {
		fmt.Printf("%-20s %-8s %-11s %s\n",
			reg.Name,
			reg.Type,
			reg.Status,
			reg.URL)
	}

	return nil
}

// Helper functions

func parseImageTag(imageArg string) (string, string) {
	parts := strings.Split(imageArg, ":")
	if len(parts) == 1 {
		return parts[0], "latest"
	}
	return strings.Join(parts[:len(parts)-1], ":"), parts[len(parts)-1]
}

func getRegistryDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Warn("Failed to get home directory: %v", err)
		return "./servin-registry"
	}
	return homeDir + "/.servin/registry"
}
