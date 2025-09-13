package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Launch Servin Desktop GUI",
	Long: `Launch the Servin Desktop graphical user interface.

The GUI provides an easy-to-use interface for managing containers, images,
CRI server, volumes, and registry operations. It includes:

- Container management (start, stop, remove, logs)
- Image operations (import, remove, tag, inspect)
- CRI server control and monitoring
- Volume management
- Registry operations
- System information and monitoring

Examples:
  servin gui                    # Launch GUI
  servin gui --tui              # Launch Terminal UI instead
  servin gui --dev              # Launch in development mode`,
	RunE: runGUI,
}

var (
	useTUI  bool
	devMode bool
	guiPort int
	guiHost string
)

func init() {
	rootCmd.AddCommand(guiCmd)

	guiCmd.Flags().BoolVar(&useTUI, "tui", false, "Use Terminal User Interface instead of GUI")
	guiCmd.Flags().BoolVar(&devMode, "dev", false, "Launch in development mode")
	guiCmd.Flags().IntVar(&guiPort, "port", 8081, "Port for GUI web interface")
	guiCmd.Flags().StringVar(&guiHost, "host", "localhost", "Host for GUI web interface")
}

func runGUI(cmd *cobra.Command, args []string) error {
	if useTUI {
		return runTUI()
	}

	// Check if we can run GUI
	if !canRunGUI() {
		fmt.Println("GUI not available on this platform or display not detected.")
		fmt.Println("Falling back to Terminal UI...")
		return runTUI()
	}

	// Try to launch Fyne GUI first, fall back to TUI if it fails
	if err := runFyneGUI(); err != nil {
		fmt.Printf("GUI launch failed: %v\n", err)
		fmt.Println("Falling back to Terminal UI...")
		return runTUI()
	}

	return nil
}

func canRunGUI() bool {
	switch runtime.GOOS {
	case "windows":
		return true // Windows always has GUI capability
	case "darwin":
		return os.Getenv("DISPLAY") != "" || true // macOS usually has GUI
	case "linux":
		// Check for X11 or Wayland
		return os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != ""
	default:
		return false
	}
}

func runFyneGUI() error {
	fmt.Println("Starting Servin Desktop GUI...")

	// Get the path to the current executable
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// Look for servin-gui executable in the same directory
	guiPath := filepath.Join(filepath.Dir(execPath), "servin-gui")
	if runtime.GOOS == "windows" {
		guiPath += ".exe"
	}

	// Check if GUI executable exists
	if _, err := os.Stat(guiPath); os.IsNotExist(err) {
		return fmt.Errorf("GUI executable not found at %s", guiPath)
	}

	// Launch the GUI
	cmd := exec.Command(guiPath)
	if devMode {
		cmd.Args = append(cmd.Args, "--dev")
	}

	return cmd.Start()
}

func runTUI() error {
	fmt.Println("Starting Servin Desktop Terminal UI...")

	// Get the path to the current executable
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// Look for servin-desktop executable in the same directory
	tuiPath := filepath.Join(filepath.Dir(execPath), "servin-desktop")
	if runtime.GOOS == "windows" {
		tuiPath += ".exe"
	}

	// Check if TUI executable exists
	if _, err := os.Stat(tuiPath); os.IsNotExist(err) {
		return fmt.Errorf("TUI executable not found at %s", tuiPath)
	}

	// Launch the TUI
	cmd := exec.Command(tuiPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
