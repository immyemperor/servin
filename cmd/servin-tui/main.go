package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ServinTUI represents a Terminal User Interface for Servin
type ServinTUI struct {
	running bool
}

// NewServinTUI creates a new TUI instance
func NewServinTUI() *ServinTUI {
	return &ServinTUI{running: true}
}

// Run starts the TUI application
func (tui *ServinTUI) Run() {
	for tui.running {
		tui.clearScreen()
		tui.showBanner()
		tui.showMainMenu()
		choice := tui.getInput("Select an option: ")
		tui.handleMainMenu(choice)
	}
}

// showBanner displays the Servin banner
func (tui *ServinTUI) showBanner() {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                         Servin Desktop                         ║")
	fmt.Println("║                Container Runtime Management                    ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

// clearScreen clears the terminal screen
func (tui *ServinTUI) clearScreen() {
	// Clear screen command for Unix/Windows
	cmd := exec.Command("clear")
	if os.Getenv("OS") == "Windows_NT" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// showMainMenu displays the main menu options
func (tui *ServinTUI) showMainMenu() {
	fmt.Println("┌─────────────────── Main Menu ────────────────────┐")
	fmt.Println("│  1. Container Management                          │")
	fmt.Println("│  2. Image Management                              │")
	fmt.Println("│  3. CRI Server Control                            │")
	fmt.Println("│  4. Volume Management                             │")
	fmt.Println("│  5. Registry Operations                           │")
	fmt.Println("│  6. System Information                            │")
	fmt.Println("│  7. Exit                                          │")
	fmt.Println("└───────────────────────────────────────────────────┘")
	fmt.Println()
}

// handleMainMenu processes main menu selections
func (tui *ServinTUI) handleMainMenu(choice string) {
	switch choice {
	case "1":
		tui.containerMenu()
	case "2":
		tui.imageMenu()
	case "3":
		tui.criMenu()
	case "4":
		tui.volumeMenu()
	case "5":
		tui.registryMenu()
	case "6":
		tui.systemInfo()
		if tui.running {
			tui.getInput("Press Enter to continue...")
		}
	case "7":
		tui.running = false
		fmt.Println("Thank you for using Servin Desktop!")
	default:
		fmt.Println("Invalid option. Please try again.")
		if tui.running {
			tui.getInput("Press Enter to continue...")
		}
	}
	fmt.Println()
}

// containerMenu handles container management
func (tui *ServinTUI) containerMenu() {
	for {
		tui.clearScreen()
		fmt.Println("┌─────────────── Container Management ──────────────┐")
		fmt.Println("│  1. List Containers                               │")
		fmt.Println("│  2. Run New Container                             │")
		fmt.Println("│  3. Start Container                               │")
		fmt.Println("│  4. Stop Container                                │")
		fmt.Println("│  5. Remove Container                              │")
		fmt.Println("│  6. View Container Logs                           │")
		fmt.Println("│  7. Execute Command in Container                  │")
		fmt.Println("│  8. Back to Main Menu                             │")
		fmt.Println("└───────────────────────────────────────────────────┘")

		choice := tui.getInput("Select an option: ")

		switch choice {
		case "1":
			tui.runCommand("servin", "ls")
			tui.getInput("Press Enter to continue...")
		case "2":
			tui.runNewContainer()
			tui.getInput("Press Enter to continue...")
		case "3":
			tui.startContainer()
			tui.getInput("Press Enter to continue...")
		case "4":
			tui.stopContainer()
			tui.getInput("Press Enter to continue...")
		case "5":
			tui.removeContainer()
			tui.getInput("Press Enter to continue...")
		case "6":
			tui.viewContainerLogs()
			tui.getInput("Press Enter to continue...")
		case "7":
			tui.execInContainer()
			tui.getInput("Press Enter to continue...")
		case "8":
			return
		default:
			fmt.Println("Invalid option. Please try again.")
			tui.getInput("Press Enter to continue...")
		}
		fmt.Println()
	}
}

// imageMenu handles image management
func (tui *ServinTUI) imageMenu() {
	for {
		fmt.Println("┌──────────────── Image Management ─────────────────┐")
		fmt.Println("│  1. List Images                                   │")
		fmt.Println("│  2. Import Image                                  │")
		fmt.Println("│  3. Remove Image                                  │")
		fmt.Println("│  4. Tag Image                                     │")
		fmt.Println("│  5. Inspect Image                                 │")
		fmt.Println("│  6. Build Image                                   │")
		fmt.Println("│  7. Back to Main Menu                             │")
		fmt.Println("└───────────────────────────────────────────────────┘")

		choice := tui.getInput("Select an option: ")

		switch choice {
		case "1":
			tui.runCommand("servin", "image", "ls")
		case "2":
			tui.importImage()
		case "3":
			tui.removeImage()
		case "4":
			tui.tagImage()
		case "5":
			tui.inspectImage()
		case "6":
			tui.buildImage()
		case "7":
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
		fmt.Println()
	}
}

// criMenu handles CRI server operations
func (tui *ServinTUI) criMenu() {
	for {
		fmt.Println("┌────────────────── CRI Server ─────────────────────┐")
		fmt.Println("│  1. Start CRI Server                             │")
		fmt.Println("│  2. Check CRI Server Status                      │")
		fmt.Println("│  3. Test CRI Connection                           │")
		fmt.Println("│  4. View CRI Endpoints                            │")
		fmt.Println("│  5. Back to Main Menu                             │")
		fmt.Println("└───────────────────────────────────────────────────┘")

		choice := tui.getInput("Select an option: ")

		switch choice {
		case "1":
			tui.startCRIServer()
		case "2":
			tui.runCommand("servin", "cri", "status")
		case "3":
			tui.testCRIConnection()
		case "4":
			tui.showCRIEndpoints()
		case "5":
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
		fmt.Println()
	}
}

// volumeMenu handles volume management
func (tui *ServinTUI) volumeMenu() {
	for {
		fmt.Println("┌─────────────── Volume Management ─────────────────┐")
		fmt.Println("│  1. List Volumes                                  │")
		fmt.Println("│  2. Create Volume                                 │")
		fmt.Println("│  3. Remove Volume                                 │")
		fmt.Println("│  4. Inspect Volume                                │")
		fmt.Println("│  5. Remove All Volumes                            │")
		fmt.Println("│  6. Back to Main Menu                             │")
		fmt.Println("└───────────────────────────────────────────────────┘")

		choice := tui.getInput("Select an option: ")

		switch choice {
		case "1":
			tui.runCommand("servin", "volume", "ls")
		case "2":
			tui.createVolume()
		case "3":
			tui.removeVolume()
		case "4":
			tui.inspectVolume()
		case "5":
			tui.removeAllVolumes()
		case "6":
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
		fmt.Println()
	}
}

// registryMenu handles registry operations
func (tui *ServinTUI) registryMenu() {
	for {
		fmt.Println("┌────────────── Registry Operations ────────────────┐")
		fmt.Println("│  1. Start Local Registry                          │")
		fmt.Println("│  2. Stop Local Registry                           │")
		fmt.Println("│  3. Push Image to Registry                        │")
		fmt.Println("│  4. Pull Image from Registry                      │")
		fmt.Println("│  5. List Registries                               │")
		fmt.Println("│  6. Back to Main Menu                             │")
		fmt.Println("└───────────────────────────────────────────────────┘")

		choice := tui.getInput("Select an option: ")

		switch choice {
		case "1":
			tui.runCommand("servin", "registry", "start")
		case "2":
			tui.runCommand("servin", "registry", "stop")
		case "3":
			tui.pushToRegistry()
		case "4":
			tui.pullFromRegistry()
		case "5":
			tui.runCommand("servin", "registry", "list")
		case "6":
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
		fmt.Println()
	}
}

// systemInfo displays system information
func (tui *ServinTUI) systemInfo() {
	fmt.Println("╔════════════════ System Information ═══════════════╗")

	// Get Servin version
	fmt.Println("║ Servin Runtime Information:                        ║")
	tui.runCommand("servin", "--version")

	// Show platform info
	fmt.Printf("║ Platform: %s\n", "Windows/Linux/macOS")
	fmt.Printf("║ Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	// Show container count
	fmt.Println("║                                                    ║")
	fmt.Println("║ Container Statistics:                              ║")
	tui.runCommand("servin", "ls")

	// Show image count
	fmt.Println("║                                                    ║")
	fmt.Println("║ Image Statistics:                                  ║")
	tui.runCommand("servin", "image", "ls")

	fmt.Println("╚════════════════════════════════════════════════════╝")
}

// Helper methods for specific operations

func (tui *ServinTUI) runNewContainer() {
	image := tui.getInput("Enter image name: ")
	if image == "" {
		fmt.Println("Image name is required.")
		return
	}

	name := tui.getInput("Enter container name (optional): ")
	command := tui.getInput("Enter command to run (optional): ")

	args := []string{"run"}
	if name != "" {
		args = append(args, "--name", name)
	}
	args = append(args, image)
	if command != "" {
		args = append(args, strings.Fields(command)...)
	}

	tui.runCommand("servin", args...)
}

func (tui *ServinTUI) startContainer() {
	containerID := tui.getInput("Enter container ID: ")
	if containerID == "" {
		fmt.Println("Container ID is required.")
		return
	}
	tui.runCommand("servin", "start", containerID)
}

func (tui *ServinTUI) stopContainer() {
	containerID := tui.getInput("Enter container ID: ")
	if containerID == "" {
		fmt.Println("Container ID is required.")
		return
	}
	tui.runCommand("servin", "stop", containerID)
}

func (tui *ServinTUI) removeContainer() {
	containerID := tui.getInput("Enter container ID: ")
	if containerID == "" {
		fmt.Println("Container ID is required.")
		return
	}

	confirm := tui.getInput("Are you sure? (y/N): ")
	if strings.ToLower(confirm) == "y" || strings.ToLower(confirm) == "yes" {
		tui.runCommand("servin", "rm", containerID)
	}
}

func (tui *ServinTUI) viewContainerLogs() {
	containerID := tui.getInput("Enter container ID: ")
	if containerID == "" {
		fmt.Println("Container ID is required.")
		return
	}

	follow := tui.getInput("Follow logs? (y/N): ")
	if strings.ToLower(follow) == "y" || strings.ToLower(follow) == "yes" {
		tui.runCommand("servin", "logs", "-f", containerID)
	} else {
		tui.runCommand("servin", "logs", containerID)
	}
}

func (tui *ServinTUI) execInContainer() {
	containerID := tui.getInput("Enter container ID: ")
	if containerID == "" {
		fmt.Println("Container ID is required.")
		return
	}

	command := tui.getInput("Enter command to execute: ")
	if command == "" {
		command = "sh"
	}

	args := append([]string{"exec", containerID}, strings.Fields(command)...)
	tui.runCommand("servin", args...)
}

func (tui *ServinTUI) importImage() {
	filepath := tui.getInput("Enter image file path: ")
	if filepath == "" {
		fmt.Println("File path is required.")
		return
	}
	tui.runCommand("servin", "image", "import", filepath)
}

func (tui *ServinTUI) removeImage() {
	imageName := tui.getInput("Enter image name: ")
	if imageName == "" {
		fmt.Println("Image name is required.")
		return
	}

	confirm := tui.getInput("Are you sure? (y/N): ")
	if strings.ToLower(confirm) == "y" || strings.ToLower(confirm) == "yes" {
		tui.runCommand("servin", "image", "rm", imageName)
	}
}

func (tui *ServinTUI) tagImage() {
	source := tui.getInput("Enter source image: ")
	if source == "" {
		fmt.Println("Source image is required.")
		return
	}

	target := tui.getInput("Enter target tag: ")
	if target == "" {
		fmt.Println("Target tag is required.")
		return
	}

	tui.runCommand("servin", "image", "tag", source, target)
}

func (tui *ServinTUI) inspectImage() {
	imageName := tui.getInput("Enter image name: ")
	if imageName == "" {
		fmt.Println("Image name is required.")
		return
	}
	tui.runCommand("servin", "image", "inspect", imageName)
}

func (tui *ServinTUI) buildImage() {
	path := tui.getInput("Enter build context path (default: .): ")
	if path == "" {
		path = "."
	}

	tag := tui.getInput("Enter image tag (optional): ")

	args := []string{"build"}
	if tag != "" {
		args = append(args, "-t", tag)
	}
	args = append(args, path)

	tui.runCommand("servin", args...)
}

func (tui *ServinTUI) startCRIServer() {
	port := tui.getInput("Enter port (default: 8080): ")
	if port == "" {
		port = "8080"
	}

	// Validate port
	if _, err := strconv.Atoi(port); err != nil {
		fmt.Println("Invalid port number.")
		return
	}

	fmt.Printf("Starting CRI server on port %s...\n", port)
	fmt.Println("Note: This will run in the background. Use Ctrl+C to stop.")
	tui.runCommand("servin", "cri", "start", "--port", port, "--verbose")
}

func (tui *ServinTUI) testCRIConnection() {
	fmt.Println("Testing CRI connection to http://localhost:8080...")

	// Test using curl if available, otherwise use servin cri test
	cmd := exec.Command("curl", "-s", "http://localhost:8080/health")
	if err := cmd.Run(); err != nil {
		fmt.Println("Curl not available, trying servin cri test...")
		tui.runCommand("servin", "cri", "test")
	} else {
		fmt.Println("✓ CRI server is responding")
	}
}

func (tui *ServinTUI) showCRIEndpoints() {
	fmt.Println("╔═══════════════ CRI Endpoints ════════════════════╗")
	fmt.Println("║ Health Check:                                     ║")
	fmt.Println("║   GET  /health                                    ║")
	fmt.Println("║                                                   ║")
	fmt.Println("║ Runtime Operations:                               ║")
	fmt.Println("║   POST /v1/runtime/version                        ║")
	fmt.Println("║   POST /v1/runtime/status                         ║")
	fmt.Println("║                                                   ║")
	fmt.Println("║ Pod Sandbox Operations:                           ║")
	fmt.Println("║   POST /v1/runtime/sandbox/list                   ║")
	fmt.Println("║   POST /v1/runtime/sandbox/create                 ║")
	fmt.Println("║   POST /v1/runtime/sandbox/start                  ║")
	fmt.Println("║   POST /v1/runtime/sandbox/stop                   ║")
	fmt.Println("║   POST /v1/runtime/sandbox/remove                 ║")
	fmt.Println("║                                                   ║")
	fmt.Println("║ Container Operations:                             ║")
	fmt.Println("║   POST /v1/runtime/container/list                 ║")
	fmt.Println("║   POST /v1/runtime/container/create               ║")
	fmt.Println("║   POST /v1/runtime/container/start                ║")
	fmt.Println("║   POST /v1/runtime/container/stop                 ║")
	fmt.Println("║   POST /v1/runtime/container/remove               ║")
	fmt.Println("║   POST /v1/runtime/container/status               ║")
	fmt.Println("║                                                   ║")
	fmt.Println("║ Image Operations:                                 ║")
	fmt.Println("║   POST /v1/image/list                             ║")
	fmt.Println("║   POST /v1/image/status                           ║")
	fmt.Println("║   POST /v1/image/pull                             ║")
	fmt.Println("║   POST /v1/image/remove                           ║")
	fmt.Println("║   POST /v1/image/fs                               ║")
	fmt.Println("╚═══════════════════════════════════════════════════╝")
}

func (tui *ServinTUI) createVolume() {
	name := tui.getInput("Enter volume name: ")
	if name == "" {
		fmt.Println("Volume name is required.")
		return
	}
	tui.runCommand("servin", "volume", "create", name)
}

func (tui *ServinTUI) removeVolume() {
	name := tui.getInput("Enter volume name: ")
	if name == "" {
		fmt.Println("Volume name is required.")
		return
	}

	confirm := tui.getInput("Are you sure? (y/N): ")
	if strings.ToLower(confirm) == "y" || strings.ToLower(confirm) == "yes" {
		tui.runCommand("servin", "volume", "rm", name)
	}
}

func (tui *ServinTUI) inspectVolume() {
	name := tui.getInput("Enter volume name: ")
	if name == "" {
		fmt.Println("Volume name is required.")
		return
	}
	tui.runCommand("servin", "volume", "inspect", name)
}

func (tui *ServinTUI) removeAllVolumes() {
	confirm := tui.getInput("Remove ALL volumes? This cannot be undone! (y/N): ")
	if strings.ToLower(confirm) == "y" || strings.ToLower(confirm) == "yes" {
		tui.runCommand("servin", "volume", "rm-all")
	}
}

func (tui *ServinTUI) pushToRegistry() {
	image := tui.getInput("Enter image name to push: ")
	if image == "" {
		fmt.Println("Image name is required.")
		return
	}

	registry := tui.getInput("Enter registry address (optional): ")

	if registry != "" {
		tui.runCommand("servin", "registry", "push", image, registry)
	} else {
		tui.runCommand("servin", "registry", "push", image)
	}
}

func (tui *ServinTUI) pullFromRegistry() {
	image := tui.getInput("Enter image name to pull: ")
	if image == "" {
		fmt.Println("Image name is required.")
		return
	}

	registry := tui.getInput("Enter registry address (optional): ")

	if registry != "" {
		tui.runCommand("servin", "registry", "pull", image, registry)
	} else {
		tui.runCommand("servin", "registry", "pull", image)
	}
}

// Utility methods

func (tui *ServinTUI) getInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func (tui *ServinTUI) runCommand(command string, args ...string) {
	fmt.Printf("Running: %s %s\n", command, strings.Join(args, " "))
	fmt.Println("─────────────────────────────────────────────────────")

	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("─────────────────────────────────────────────────────")
}

func main() {
	tui := NewServinTUI()
	tui.Run()
}
