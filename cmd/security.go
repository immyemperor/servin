package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"servin/pkg/logger"
	"servin/pkg/namespaces"

	"github.com/spf13/cobra"
)

var securityCmd = &cobra.Command{
	Use:   "security",
	Short: "Manage container security settings",
	Long: `Security commands allow you to:
- Check user namespace support
- Configure UID/GID mappings
- Inspect security configurations
- Validate security settings`,
}

var securityCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check security features availability",
	Long: `Check if security features like user namespaces are available on this system.

Examples:
  servin security check                    # Check all security features
  servin security check --user-ns         # Check only user namespace support`,
	RunE: runSecurityCheck,
}

var securityInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show current security configuration",
	Long: `Display information about the current security configuration including:
- User namespace status
- UID/GID mappings
- Capabilities
- Security constraints

Examples:
  servin security info                     # Show all security info
  servin security info --container ID     # Show security info for specific container`,
	RunE: runSecurityInfo,
}

var securityConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure security settings for containers",
	Long: `Configure security settings such as:
- User namespace mappings
- Capability restrictions
- Security policies

Examples:
  servin security config --uid-map "0:1000:1"     # Map root to user 1000
  servin security config --rootless               # Enable rootless mode
  servin security config --no-new-privs           # Prevent privilege escalation`,
	RunE: runSecurityConfig,
}

var securityTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test security configuration with a simple container",
	Long: `Run a test container with specified security settings to verify configuration.

Examples:
  servin security test                             # Test with default security
  servin security test --user-ns                  # Test with user namespaces
  servin security test --rootless                 # Test rootless configuration`,
	RunE: runSecurityTest,
}

func init() {
	rootCmd.AddCommand(securityCmd)

	securityCmd.AddCommand(securityCheckCmd)
	securityCmd.AddCommand(securityInfoCmd)
	securityCmd.AddCommand(securityConfigCmd)
	securityCmd.AddCommand(securityTestCmd)

	// Security check flags
	securityCheckCmd.Flags().Bool("user-ns", false, "Check only user namespace support")
	securityCheckCmd.Flags().Bool("capabilities", false, "Check capability support")

	// Security info flags
	securityInfoCmd.Flags().String("container", "", "Show security info for specific container")
	securityInfoCmd.Flags().Bool("verbose", false, "Show detailed information")

	// Security config flags
	securityConfigCmd.Flags().StringSlice("uid-map", []string{}, "UID mapping (format: container_uid:host_uid:count)")
	securityConfigCmd.Flags().StringSlice("gid-map", []string{}, "GID mapping (format: container_gid:host_gid:count)")
	securityConfigCmd.Flags().Bool("rootless", false, "Configure for rootless operation")
	securityConfigCmd.Flags().Bool("no-new-privs", false, "Prevent privilege escalation")
	securityConfigCmd.Flags().Bool("drop-caps", false, "Drop all capabilities")
	securityConfigCmd.Flags().StringSlice("cap-add", []string{}, "Add specific capabilities")
	securityConfigCmd.Flags().String("user", "", "Run as specific user (uid:gid)")

	// Security test flags
	securityTestCmd.Flags().Bool("user-ns", false, "Test with user namespaces")
	securityTestCmd.Flags().Bool("rootless", false, "Test rootless configuration")
	securityTestCmd.Flags().String("image", "alpine", "Image to use for testing")
	securityTestCmd.Flags().String("command", "id", "Command to run in test container")
}

func runSecurityCheck(cmd *cobra.Command, args []string) error {
	userNSOnly, _ := cmd.Flags().GetBool("user-ns")
	capsOnly, _ := cmd.Flags().GetBool("capabilities")

	fmt.Println("Security Feature Availability Check")
	fmt.Println("====================================")

	if !userNSOnly && !capsOnly {
		// Check all features
		fmt.Println("Checking all security features...")
	}

	if !capsOnly {
		// Check user namespace support
		fmt.Print("User Namespaces: ")
		if err := namespaces.ValidateUserNamespaceSupport(); err != nil {
			fmt.Printf("❌ Not Available (%v)\n", err)
		} else {
			fmt.Println("✅ Available")
		}
	}

	if !userNSOnly {
		// Check other security features
		fmt.Print("Capabilities: ")
		if checkCapabilitySupport() {
			fmt.Println("✅ Available")
		} else {
			fmt.Println("❌ Not Available")
		}

		fmt.Print("No New Privileges: ")
		if checkNoNewPrivsSupport() {
			fmt.Println("✅ Available")
		} else {
			fmt.Println("❌ Not Available")
		}
	}

	// Show recommendations
	fmt.Println("\nRecommendations:")
	if err := namespaces.ValidateUserNamespaceSupport(); err != nil {
		fmt.Println("- Enable user namespaces in kernel configuration")
		fmt.Println("- Check if max_user_namespaces > 0")
	} else {
		fmt.Println("- User namespaces are available for enhanced security")
		fmt.Println("- Consider using rootless containers for better isolation")
	}

	return nil
}

func runSecurityInfo(cmd *cobra.Command, args []string) error {
	containerID, _ := cmd.Flags().GetString("container")
	verbose, _ := cmd.Flags().GetBool("verbose")

	fmt.Println("Current Security Configuration")
	fmt.Println("==============================")

	// Get user namespace info
	userNSInfo, err := namespaces.GetUserNamespaceInfo()
	if err != nil {
		logger.Warn("Failed to get user namespace info: %v", err)
	} else {
		fmt.Printf("In User Namespace: %t\n", userNSInfo.InUserNamespace)
		fmt.Printf("Current UID: %d\n", userNSInfo.UID)
		fmt.Printf("Current GID: %d\n", userNSInfo.GID)
		fmt.Printf("Effective UID: %d\n", userNSInfo.EUID)
		fmt.Printf("Effective GID: %d\n", userNSInfo.EGID)

		if len(userNSInfo.UIDMappings) > 0 {
			fmt.Println("\nUID Mappings:")
			for _, mapping := range userNSInfo.UIDMappings {
				fmt.Printf("  Container UID %d -> Host UID %d (count: %d)\n",
					mapping.ContainerID, mapping.HostID, mapping.Size)
			}
		}

		if len(userNSInfo.GIDMappings) > 0 {
			fmt.Println("\nGID Mappings:")
			for _, mapping := range userNSInfo.GIDMappings {
				fmt.Printf("  Container GID %d -> Host GID %d (count: %d)\n",
					mapping.ContainerID, mapping.HostID, mapping.Size)
			}
		}

		if verbose && len(userNSInfo.Capabilities) > 0 {
			fmt.Println("\nCapabilities:")
			for _, cap := range userNSInfo.Capabilities {
				fmt.Printf("  %s\n", cap)
			}
		}
	}

	if containerID != "" {
		fmt.Printf("\nContainer %s Security Info:\n", containerID)
		fmt.Println("(Container-specific security info not yet implemented)")
	}

	return nil
}

func runSecurityConfig(cmd *cobra.Command, args []string) error {
	uidMaps, _ := cmd.Flags().GetStringSlice("uid-map")
	gidMaps, _ := cmd.Flags().GetStringSlice("gid-map")
	rootless, _ := cmd.Flags().GetBool("rootless")
	noNewPrivs, _ := cmd.Flags().GetBool("no-new-privs")
	dropCaps, _ := cmd.Flags().GetBool("drop-caps")
	capAdd, _ := cmd.Flags().GetStringSlice("cap-add")
	user, _ := cmd.Flags().GetString("user")

	fmt.Println("Configuring Security Settings")
	fmt.Println("=============================")

	var config *namespaces.UserNamespaceConfig

	if rootless {
		config = namespaces.RootlessUserNamespaceConfig()
		fmt.Println("✅ Configured for rootless operation")
	} else {
		config = namespaces.DefaultUserNamespaceConfig()
	}

	// Parse UID mappings
	if len(uidMaps) > 0 {
		config.UIDMappings = []namespaces.UIDGIDMapping{}
		for _, mapping := range uidMaps {
			parsed, err := parseUIDGIDMapping(mapping)
			if err != nil {
				return fmt.Errorf("invalid UID mapping '%s': %w", mapping, err)
			}
			config.UIDMappings = append(config.UIDMappings, parsed)
		}
		fmt.Printf("✅ Configured UID mappings: %+v\n", config.UIDMappings)
	}

	// Parse GID mappings
	if len(gidMaps) > 0 {
		config.GIDMappings = []namespaces.UIDGIDMapping{}
		for _, mapping := range gidMaps {
			parsed, err := parseUIDGIDMapping(mapping)
			if err != nil {
				return fmt.Errorf("invalid GID mapping '%s': %w", mapping, err)
			}
			config.GIDMappings = append(config.GIDMappings, parsed)
		}
		fmt.Printf("✅ Configured GID mappings: %+v\n", config.GIDMappings)
	}

	// Configure security options
	if noNewPrivs {
		config.NoNewPrivs = true
		fmt.Println("✅ Enabled no_new_privs")
	}

	if dropCaps {
		config.DropAllCaps = true
		config.AllowedCaps = capAdd
		fmt.Printf("✅ Configured to drop all capabilities, allowed: %v\n", capAdd)
	}

	// Parse user specification
	if user != "" {
		if err := parseUserSpec(user, config); err != nil {
			return fmt.Errorf("invalid user specification '%s': %w", user, err)
		}
		fmt.Printf("✅ Configured to run as UID: %d, GID: %d\n",
			config.ContainerUID, config.ContainerGID)
	}

	// Save configuration (simplified - would normally save to config file)
	fmt.Println("\nConfiguration Summary:")
	fmt.Printf("  Enabled: %t\n", config.Enabled)
	fmt.Printf("  Container UID: %d\n", config.ContainerUID)
	fmt.Printf("  Container GID: %d\n", config.ContainerGID)
	fmt.Printf("  No New Privileges: %t\n", config.NoNewPrivs)
	fmt.Printf("  Drop Capabilities: %t\n", config.DropAllCaps)

	return nil
}

func runSecurityTest(cmd *cobra.Command, args []string) error {
	userNS, _ := cmd.Flags().GetBool("user-ns")
	rootless, _ := cmd.Flags().GetBool("rootless")
	image, _ := cmd.Flags().GetString("image")
	command, _ := cmd.Flags().GetString("command")

	fmt.Println("Testing Security Configuration")
	fmt.Println("==============================")

	var configDescription string

	if rootless {
		configDescription = "rootless configuration"
	} else if userNS {
		configDescription = "default user namespace configuration"
	} else {
		configDescription = "standard configuration (no user namespaces)"
	}

	fmt.Printf("Using %s\n", configDescription)
	fmt.Printf("Test image: %s\n", image)
	fmt.Printf("Test command: %s\n", command)

	// This would create a test container with the security configuration
	fmt.Println("\n🚧 Security test implementation in progress")
	fmt.Println("This would:")
	fmt.Println("1. Create a test container with specified security settings")
	fmt.Println("2. Run the test command inside the container")
	fmt.Println("3. Verify the security constraints are working")
	fmt.Println("4. Report the results")

	return nil
}

// Helper functions

func parseUIDGIDMapping(mapping string) (namespaces.UIDGIDMapping, error) {
	parts := strings.Split(mapping, ":")
	if len(parts) != 3 {
		return namespaces.UIDGIDMapping{}, fmt.Errorf("mapping must be in format container_id:host_id:count")
	}

	containerID, err := strconv.Atoi(parts[0])
	if err != nil {
		return namespaces.UIDGIDMapping{}, fmt.Errorf("invalid container ID: %w", err)
	}

	hostID, err := strconv.Atoi(parts[1])
	if err != nil {
		return namespaces.UIDGIDMapping{}, fmt.Errorf("invalid host ID: %w", err)
	}

	size, err := strconv.Atoi(parts[2])
	if err != nil {
		return namespaces.UIDGIDMapping{}, fmt.Errorf("invalid count: %w", err)
	}

	return namespaces.UIDGIDMapping{
		ContainerID: containerID,
		HostID:      hostID,
		Size:        size,
	}, nil
}

func parseUserSpec(userSpec string, config *namespaces.UserNamespaceConfig) error {
	parts := strings.Split(userSpec, ":")

	uid, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid UID: %w", err)
	}
	config.ContainerUID = uid

	if len(parts) > 1 {
		gid, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("invalid GID: %w", err)
		}
		config.ContainerGID = gid
	}

	return nil
}

func checkCapabilitySupport() bool {
	// Check if capabilities are supported (simplified check)
	_, err := os.Stat("/proc/self/status")
	return err == nil
}

func checkNoNewPrivsSupport() bool {
	// Check if no_new_privs is supported (simplified check)
	_, err := os.Stat("/proc/sys/kernel/seccomp")
	return err == nil
}
