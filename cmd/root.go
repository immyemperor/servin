package cmd

import (
	"fmt"
	"os"
	"runtime"

	"servin/pkg/errors"
	"servin/pkg/logger"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "servin",
	Short: "A lightweight Go-based container runtime",
	Long: `Servin is a lightweight container runtime built from scratch in Go.
It implements core containerization features using Linux namespaces, cgroups,
and chroot without relying on external container runtimes.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().Bool("dev", false, "development mode (skip root check)")
	rootCmd.PersistentFlags().String("log-level", "info", "log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().String("log-file", "", "log file path (default: platform-specific)")

	// Initialize logging
	cobra.OnInitialize(initLogging)
}

// initLogging initializes the logging system
func initLogging() {
	verbose, _ := rootCmd.PersistentFlags().GetBool("verbose")
	logLevelStr, _ := rootCmd.PersistentFlags().GetString("log-level")
	logFile, _ := rootCmd.PersistentFlags().GetString("log-file")

	// Parse log level
	var logLevel logger.LogLevel
	switch logLevelStr {
	case "debug":
		logLevel = logger.DEBUG
	case "info":
		logLevel = logger.INFO
	case "warn":
		logLevel = logger.WARN
	case "error":
		logLevel = logger.ERROR
	default:
		logLevel = logger.INFO
	}

	// Use default log path if not specified
	if logFile == "" {
		logFile = logger.GetLogPath()
	}

	// Initialize logger
	if err := logger.InitLogger(logLevel, verbose, logFile); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize logging: %v\n", err)
	}

	logger.Debug("Logging initialized - level: %s, verbose: %v, file: %s", logLevelStr, verbose, logFile)
}

// checkRoot ensures the command is run with root privileges
func checkRoot() error {
	logger.Debug("Checking root privileges for platform: %s", runtime.GOOS)

	// Handle different platforms
	switch runtime.GOOS {
	case "windows":
		logger.Info("Running on Windows - containerization features limited")
		fmt.Println("Note: Running on Windows - containerization features limited")
		return nil
	case "darwin":
		logger.Info("Running on macOS - containerization features limited")
		fmt.Println("Note: Running on macOS - containerization features limited")
		// Check for development mode flag
		if devMode, _ := rootCmd.PersistentFlags().GetBool("dev"); devMode {
			logger.Info("Development mode enabled - skipping root check")
			fmt.Println("Note: Development mode - skipping root check")
			return nil
		}
		if os.Geteuid() != 0 {
			err := errors.NewPermissionError("checkRoot", "root privileges required on macOS")
			logger.Error("Root privilege check failed: %v", err)
			return fmt.Errorf("this command requires root privileges on macOS. Please run with sudo")
		}
		logger.Debug("Root privileges confirmed on macOS")
		return nil
	}

	// Check for development mode flag
	if devMode, _ := rootCmd.PersistentFlags().GetBool("dev"); devMode {
		logger.Info("Development mode enabled - skipping root check")
		fmt.Println("Note: Development mode - skipping root check")
		return nil
	}

	if os.Geteuid() != 0 {
		err := errors.NewPermissionError("checkRoot", "root privileges required on Linux")
		logger.Error("Root privilege check failed: %v", err)
		return fmt.Errorf("this command requires root privileges. Please run with sudo")
	}

	logger.Debug("Root privileges confirmed on Linux")
	return nil
}

// checkRootForContainerOps ensures root privileges only for container operations that need them
func checkRootForContainerOps() error {
	logger.Debug("Checking root privileges for container operations on platform: %s", runtime.GOOS)

	// Handle different platforms
	switch runtime.GOOS {
	case "windows":
		logger.Info("Running on Windows - containerization features limited")
		fmt.Println("Note: Running on Windows - true containerization not available, running process directly")
		return nil
	case "darwin":
		logger.Info("Running on macOS - containerization features limited")
		fmt.Println("Note: Running on macOS - true containerization requires Linux, running process with limited isolation")
		// Check for development mode flag
		if devMode, _ := rootCmd.PersistentFlags().GetBool("dev"); devMode {
			logger.Info("Development mode enabled - skipping root check")
			fmt.Println("Note: Development mode - running without full containerization")
			return nil
		}
		// On macOS, we can run without root but with warnings
		if os.Geteuid() != 0 {
			fmt.Println("Warning: Running without root privileges - containerization features will be simulated")
			return nil
		}
		logger.Debug("Root privileges confirmed on macOS")
		return nil
	}

	// On Linux, we need root for real containerization
	// Check for development mode flag
	if devMode, _ := rootCmd.PersistentFlags().GetBool("dev"); devMode {
		logger.Info("Development mode enabled - skipping root check")
		fmt.Println("Note: Development mode - containerization may not work properly without root")
		return nil
	}

	if os.Geteuid() != 0 {
		err := errors.NewPermissionError("checkRootForContainerOps", "root privileges required for containerization on Linux")
		logger.Error("Root privilege check failed: %v", err)
		return fmt.Errorf("containerization requires root privileges on Linux. Please run with sudo")
	}

	logger.Debug("Root privileges confirmed on Linux")
	return nil
}
