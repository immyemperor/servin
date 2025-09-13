package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// getServinExecutable finds the servin executable path
func getServinExecutable() string {
	// Get the current executable directory
	exePath, err := os.Executable()
	if err != nil {
		return "servin" // fallback to PATH lookup
	}

	exeDir := filepath.Dir(exePath)

	// Check for servin.exe in the same directory as the GUI
	servinPath := filepath.Join(exeDir, "servin.exe")
	if _, err := os.Stat(servinPath); err == nil {
		return servinPath
	}

	// Check for servin in the same directory (Unix-style)
	servinPath = filepath.Join(exeDir, "servin")
	if _, err := os.Stat(servinPath); err == nil {
		return servinPath
	}

	// Fallback to PATH lookup
	return "servin"
}

// ContainerData represents the JSON response from servin ls --json
type ContainerData struct {
	ID       string   `json:"id"`
	Names    []string `json:"names"`
	Image    string   `json:"image"`
	Command  string   `json:"command"`
	Created  string   `json:"created"`
	Status   string   `json:"status"`
	State    string   `json:"state"`
	Ports    string   `json:"ports"`
	Networks string   `json:"networks"`
	Mounts   string   `json:"mounts"`
	PID      string   `json:"pid"`
	Uptime   string   `json:"uptime"`
}

// refreshContainers gets the current list of containers
func (gui *ServinDesktopGUI) refreshContainers() {
	// Try to get real data from servin CLI
	servinCmd := getServinExecutable()
	cmd := exec.Command(servinCmd, "ls", "--json")
	output, err := cmd.Output()

	if err != nil {
		// Use demo data if command fails
		gui.updateStatus("Using demo data (servin CLI not available)")
		gui.containers = []ContainerInfo{
			{
				ID: "abc123def456", Name: "test-container", Image: "alpine", Tag: "latest",
				Status: "running", State: "running", Ports: "8080:80", Created: "2 hours ago",
				Uptime: "2h 15m", PID: "1234", Command: "/bin/sh", Networks: "bridge", Mounts: "/data:/app/data",
			},
			{
				ID: "def456ghi789", Name: "web-server", Image: "nginx", Tag: "latest",
				Status: "stopped", State: "exited", Ports: "-", Created: "1 day ago",
				Uptime: "0", PID: "-", Command: "nginx -g daemon off;", Networks: "bridge", Mounts: "/etc/nginx:/etc/nginx:ro",
			},
			{
				ID: "ghi789abc123", Name: "database", Image: "postgres", Tag: "13",
				Status: "running", State: "running", Ports: "5432:5432", Created: "3 hours ago",
				Uptime: "3h 42m", PID: "5678", Command: "postgres", Networks: "bridge", Mounts: "/var/lib/postgresql/data:/data",
			},
			{
				ID: "jkl012mno345", Name: "redis-cache", Image: "redis", Tag: "alpine",
				Status: "running", State: "running", Ports: "6379:6379", Created: "1 hour ago",
				Uptime: "1h 23m", PID: "9012", Command: "redis-server", Networks: "bridge", Mounts: "-",
			},
		}
	} else {
		// Parse real data from JSON output
		var containerData []ContainerData
		if err := json.Unmarshal(output, &containerData); err != nil {
			// If JSON parsing fails, try to parse text output
			gui.parseTextContainers(string(output))
		} else {
			// Convert JSON data to ContainerInfo
			gui.containers = make([]ContainerInfo, len(containerData))
			for i, container := range containerData {
				name := container.ID[:12] // Use short ID as fallback
				if len(container.Names) > 0 {
					name = container.Names[0]
				}

				// Split image into name and tag
				imageParts := strings.Split(container.Image, ":")
				imageName := imageParts[0]
				imageTag := "latest"
				if len(imageParts) > 1 {
					imageTag = imageParts[1]
				}

				gui.containers[i] = ContainerInfo{
					ID:       container.ID,
					Name:     name,
					Image:    imageName,
					Tag:      imageTag,
					Status:   container.Status,
					State:    container.State,
					Ports:    container.Ports,
					Created:  container.Created,
					Uptime:   container.Uptime,
					PID:      container.PID,
					Command:  container.Command,
					Networks: container.Networks,
					Mounts:   container.Mounts,
				}
			}
		}
		gui.updateStatus(fmt.Sprintf("Loaded %d containers from Servin", len(gui.containers)))
	}

	if gui.containerList != nil {
		gui.containerList.Refresh()
	}
}

// parseTextContainers parses the text output from servin ls command
func (gui *ServinDesktopGUI) parseTextContainers(output string) {
	lines := strings.Split(output, "\n")
	gui.containers = []ContainerInfo{}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip header, empty lines, and info messages
		if line == "" || strings.HasPrefix(line, "CONTAINER ID") ||
			strings.HasPrefix(line, "State directory:") ||
			strings.Contains(line, "[INFO]") ||
			strings.Contains(line, "Note:") {
			continue
		}

		// Parse container line (space-separated)
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			container := ContainerInfo{
				ID:      fields[0],
				Image:   fields[1],
				Status:  fields[4],
				Created: strings.Join(fields[3:5], " "),
			}

			// Extract container name (last field)
			if len(fields) > 5 {
				container.Name = fields[len(fields)-1]
			} else {
				container.Name = container.ID[:12]
			}

			// Determine state from status
			if strings.Contains(container.Status, "running") {
				container.State = "running"
			} else if strings.Contains(container.Status, "exited") {
				container.State = "exited"
			} else {
				container.State = container.Status
			}

			gui.containers = append(gui.containers, container)
		}
	}
}

// parseTextImages parses the text output from servin images command
func (gui *ServinDesktopGUI) parseTextImages(output string) {
	lines := strings.Split(output, "\n")
	gui.images = []ImageInfo{}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip header, empty lines, and info messages
		if line == "" || strings.HasPrefix(line, "REPOSITORY") ||
			strings.Contains(line, "[INFO]") ||
			strings.Contains(line, "Note:") {
			continue
		}

		// Parse image line (space-separated)
		fields := strings.Fields(line)
		if len(fields) >= 5 {
			repository := fields[0]
			tag := fields[1]
			imageID := fields[2]
			created := strings.Join(fields[3:len(fields)-1], " ")
			size := fields[len(fields)-1]

			// Handle <none> tags
			if tag == "<none>" {
				tag = "none"
			}

			image := ImageInfo{
				ID:      imageID,
				Name:    repository,
				Tag:     tag,
				Size:    size,
				Created: created,
			}

			gui.images = append(gui.images, image)
		}
	}
}

func (gui *ServinDesktopGUI) runContainer() {
	// Create form for running new container
	imageEntry := widget.NewEntry()
	imageEntry.SetPlaceHolder("alpine:latest")

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("my-container")

	commandEntry := widget.NewEntry()
	commandEntry.SetPlaceHolder("echo 'Hello World'")

	form := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Image", imageEntry),
			widget.NewFormItem("Name (optional)", nameEntry),
			widget.NewFormItem("Command (optional)", commandEntry),
		},
		OnSubmit: func() {
			image := imageEntry.Text
			if image == "" {
				dialog.ShowError(fmt.Errorf("Image is required"), gui.window)
				return
			}

			args := []string{"run"}
			if nameEntry.Text != "" {
				args = append(args, "--name", nameEntry.Text)
			}
			args = append(args, image)
			if commandEntry.Text != "" {
				args = append(args, strings.Fields(commandEntry.Text)...)
			}

			cmd := exec.Command("servin", args...)
			if err := cmd.Run(); err != nil {
				dialog.ShowError(fmt.Errorf("Failed to run container: %v", err), gui.window)
				return
			}

			gui.updateStatus(fmt.Sprintf("Container started from image %s", image))
			gui.refreshContainers()
		},
	}

	dialog.ShowForm("Run New Container", "Run", "Cancel", form.Items, func(confirm bool) {
		if confirm {
			form.OnSubmit()
		}
	}, gui.window)
}

func (gui *ServinDesktopGUI) startContainer() {
	containerInfo := gui.getSelectedContainer()
	if containerInfo == nil {
		dialog.ShowInformation("No Selection", "Please select a container to start", gui.window)
		return
	}

	cmd := exec.Command("servin", "start", containerInfo.ID)
	if err := cmd.Run(); err != nil {
		dialog.ShowError(fmt.Errorf("Failed to start container: %v", err), gui.window)
		return
	}

	gui.updateStatus(fmt.Sprintf("Container %s started", containerInfo.Name))
	gui.refreshContainers()
}

func (gui *ServinDesktopGUI) stopContainer() {
	containerInfo := gui.getSelectedContainer()
	if containerInfo == nil {
		dialog.ShowInformation("No Selection", "Please select a container to stop", gui.window)
		return
	}

	cmd := exec.Command("servin", "stop", containerInfo.ID)
	if err := cmd.Run(); err != nil {
		dialog.ShowError(fmt.Errorf("Failed to stop container: %v", err), gui.window)
		return
	}

	gui.updateStatus(fmt.Sprintf("Container %s stopped", containerInfo.Name))
	gui.refreshContainers()
}

func (gui *ServinDesktopGUI) restartContainer() {
	containerInfo := gui.getSelectedContainer()
	if containerInfo == nil {
		dialog.ShowInformation("No Selection", "Please select a container to restart", gui.window)
		return
	}

	cmd := exec.Command("servin", "restart", containerInfo.ID)
	if err := cmd.Run(); err != nil {
		dialog.ShowError(fmt.Errorf("Failed to restart container: %v", err), gui.window)
		return
	}

	gui.updateStatus(fmt.Sprintf("Container %s restarted", containerInfo.Name))
	gui.refreshContainers()
}

func (gui *ServinDesktopGUI) removeContainer() {
	containerInfo := gui.getSelectedContainer()
	if containerInfo == nil {
		dialog.ShowInformation("No Selection", "Please select a container to remove", gui.window)
		return
	}

	confirm := dialog.NewConfirm(
		"Remove Container",
		fmt.Sprintf("Are you sure you want to remove container '%s'?", containerInfo.Name),
		func(confirmed bool) {
			if confirmed {
				cmd := exec.Command("servin", "rm", containerInfo.ID)
				if err := cmd.Run(); err != nil {
					dialog.ShowError(fmt.Errorf("Failed to remove container: %v", err), gui.window)
					return
				}

				gui.updateStatus(fmt.Sprintf("Container %s removed", containerInfo.Name))
				gui.refreshContainers()
			}
		},
		gui.window,
	)
	confirm.Show()
}

func (gui *ServinDesktopGUI) showContainerLogs() {
	containerInfo := gui.getSelectedContainer()
	if containerInfo == nil {
		dialog.ShowInformation("No Selection", "Please select a container to view logs", gui.window)
		return
	}

	cmd := exec.Command("servin", "logs", containerInfo.ID)
	output, err := cmd.Output()
	if err != nil {
		dialog.ShowError(fmt.Errorf("Failed to get logs: %v", err), gui.window)
		return
	}

	// Show logs in new window
	logWindow := gui.app.NewWindow(fmt.Sprintf("Logs - %s", containerInfo.Name))
	logWindow.Resize(fyne.NewSize(800, 600))

	logText := widget.NewMultiLineEntry()
	logText.SetText(string(output))
	logText.Wrapping = fyne.TextWrapWord

	logWindow.SetContent(container.NewScroll(logText))
	logWindow.Show()
}

func (gui *ServinDesktopGUI) inspectContainer() {
	containerInfo := gui.getSelectedContainer()
	if containerInfo == nil {
		dialog.ShowInformation("No Selection", "Please select a container to inspect", gui.window)
		return
	}

	cmd := exec.Command("servin", "inspect", containerInfo.ID)
	output, err := cmd.Output()
	if err != nil {
		dialog.ShowError(fmt.Errorf("Failed to inspect container: %v", err), gui.window)
		return
	}

	// Show inspection in new window
	inspectWindow := gui.app.NewWindow(fmt.Sprintf("Inspect - %s", containerInfo.Name))
	inspectWindow.Resize(fyne.NewSize(800, 600))

	inspectText := widget.NewMultiLineEntry()
	inspectText.SetText(string(output))
	inspectText.Wrapping = fyne.TextWrapWord

	inspectWindow.SetContent(container.NewScroll(inspectText))
	inspectWindow.Show()
}

// Image management methods
func (gui *ServinDesktopGUI) refreshImages() {
	// Try to get real data from servin CLI
	servinCmd := getServinExecutable()
	cmd := exec.Command(servinCmd, "images")
	output, err := cmd.Output()

	if err != nil {
		// Use demo data if command fails
		gui.updateStatus("Using demo data (servin CLI not available)")
		gui.images = []ImageInfo{
			{ID: "sha256:abc123", Name: "alpine", Tag: "latest", Size: "5.6MB", Created: "2 hours ago"},
			{ID: "sha256:def456", Name: "nginx", Tag: "latest", Size: "133MB", Created: "1 day ago"},
			{ID: "sha256:ghi789", Name: "postgres", Tag: "13", Size: "314MB", Created: "3 days ago"},
			{ID: "sha256:jkl012", Name: "redis", Tag: "alpine", Size: "28MB", Created: "5 hours ago"},
		}
	} else {
		// Parse real data from text output
		gui.parseTextImages(string(output))
		gui.updateStatus(fmt.Sprintf("Loaded %d images from Servin", len(gui.images)))
	}

	if gui.imageList != nil {
		gui.imageList.Refresh()
	}
}

func (gui *ServinDesktopGUI) pullImage() {
	// Create form for pulling new image
	imageEntry := widget.NewEntry()
	imageEntry.SetPlaceHolder("alpine:latest")

	form := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Image Name", imageEntry),
		},
		OnSubmit: func() {
			imageName := imageEntry.Text
			if imageName == "" {
				dialog.ShowError(fmt.Errorf("Image name is required"), gui.window)
				return
			}

			// Show progress dialog
			progress := dialog.NewProgressInfinite("Pulling Image", fmt.Sprintf("Pulling %s...", imageName), gui.window)
			progress.Show()

			go func() {
				cmd := exec.Command("servin", "pull", imageName)
				err := cmd.Run()
				progress.Hide()

				if err != nil {
					dialog.ShowError(fmt.Errorf("Failed to pull image: %v", err), gui.window)
					return
				}

				gui.updateStatus(fmt.Sprintf("Image %s pulled successfully", imageName))
				gui.refreshImages()
			}()
		},
	}

	dialog.ShowForm("Pull Image", "Pull", "Cancel", form.Items, func(confirm bool) {
		if confirm {
			form.OnSubmit()
		}
	}, gui.window)
}

func (gui *ServinDesktopGUI) removeImage() {
	imageInfo := gui.getSelectedImage()
	if imageInfo == nil {
		dialog.ShowInformation("No Selection", "Please select an image to remove", gui.window)
		return
	}

	confirm := dialog.NewConfirm(
		"Remove Image",
		fmt.Sprintf("Are you sure you want to remove image '%s:%s'?", imageInfo.Name, imageInfo.Tag),
		func(confirmed bool) {
			if confirmed {
				cmd := exec.Command("servin", "rmi", imageInfo.ID)
				if err := cmd.Run(); err != nil {
					dialog.ShowError(fmt.Errorf("Failed to remove image: %v", err), gui.window)
					return
				}

				gui.updateStatus(fmt.Sprintf("Image %s:%s removed", imageInfo.Name, imageInfo.Tag))
				gui.refreshImages()
			}
		},
		gui.window,
	)
	confirm.Show()
}

func (gui *ServinDesktopGUI) inspectImage() {
	imageInfo := gui.getSelectedImage()
	if imageInfo == nil {
		dialog.ShowInformation("No Selection", "Please select an image to inspect", gui.window)
		return
	}

	cmd := exec.Command("servin", "inspect", imageInfo.ID)
	output, err := cmd.Output()
	if err != nil {
		dialog.ShowError(fmt.Errorf("Failed to inspect image: %v", err), gui.window)
		return
	}

	// Show inspection in new window
	inspectWindow := gui.app.NewWindow(fmt.Sprintf("Inspect - %s:%s", imageInfo.Name, imageInfo.Tag))
	inspectWindow.Resize(fyne.NewSize(800, 600))

	inspectText := widget.NewMultiLineEntry()
	inspectText.SetText(string(output))
	inspectText.Wrapping = fyne.TextWrapWord

	inspectWindow.SetContent(container.NewScroll(inspectText))
	inspectWindow.Show()
}

// Volume management methods
func (gui *ServinDesktopGUI) refreshVolumes() {
	// Try to get real data, but fall back to mock data if command fails
	cmd := exec.Command("servin", "volumes", "--json")
	_, err := cmd.Output()
	if err != nil {
		// Use mock data for demo purposes
		gui.updateStatus("Using demo data (servin CLI not available)")
	} else {
		gui.updateStatus("Volumes refreshed")
	}

	// Always populate with demo data for now
	gui.volumes = []VolumeInfo{
		{Name: "app-data", Driver: "local", Mountpoint: "/var/lib/servin/volumes/app-data", Created: "2 hours ago"},
		{Name: "db-data", Driver: "local", Mountpoint: "/var/lib/servin/volumes/db-data", Created: "1 day ago"},
		{Name: "log-data", Driver: "local", Mountpoint: "/var/lib/servin/volumes/log-data", Created: "3 days ago"},
		{Name: "cache-data", Driver: "local", Mountpoint: "/var/lib/servin/volumes/cache-data", Created: "4 hours ago"},
	}

	if gui.volumeList != nil {
		gui.volumeList.Refresh()
	}
}

func (gui *ServinDesktopGUI) createVolume() {
	// Create form for new volume
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("my-volume")

	driverEntry := widget.NewEntry()
	driverEntry.SetText("local")

	form := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Volume Name", nameEntry),
			widget.NewFormItem("Driver", driverEntry),
		},
		OnSubmit: func() {
			volumeName := nameEntry.Text
			if volumeName == "" {
				dialog.ShowError(fmt.Errorf("Volume name is required"), gui.window)
				return
			}

			args := []string{"volume", "create"}
			if driverEntry.Text != "" && driverEntry.Text != "local" {
				args = append(args, "--driver", driverEntry.Text)
			}
			args = append(args, volumeName)

			cmd := exec.Command("servin", args...)
			if err := cmd.Run(); err != nil {
				dialog.ShowError(fmt.Errorf("Failed to create volume: %v", err), gui.window)
				return
			}

			gui.updateStatus(fmt.Sprintf("Volume %s created", volumeName))
			gui.refreshVolumes()
		},
	}

	dialog.ShowForm("Create Volume", "Create", "Cancel", form.Items, func(confirm bool) {
		if confirm {
			form.OnSubmit()
		}
	}, gui.window)
}

func (gui *ServinDesktopGUI) removeVolume() {
	volumeInfo := gui.getSelectedVolume()
	if volumeInfo == nil {
		dialog.ShowInformation("No Selection", "Please select a volume to remove", gui.window)
		return
	}

	confirm := dialog.NewConfirm(
		"Remove Volume",
		fmt.Sprintf("Are you sure you want to remove volume '%s'?", volumeInfo.Name),
		func(confirmed bool) {
			if confirmed {
				cmd := exec.Command("servin", "volume", "rm", volumeInfo.Name)
				if err := cmd.Run(); err != nil {
					dialog.ShowError(fmt.Errorf("Failed to remove volume: %v", err), gui.window)
					return
				}

				gui.updateStatus(fmt.Sprintf("Volume %s removed", volumeInfo.Name))
				gui.refreshVolumes()
			}
		},
		gui.window,
	)
	confirm.Show()
}

func (gui *ServinDesktopGUI) inspectVolume() {
	volumeInfo := gui.getSelectedVolume()
	if volumeInfo == nil {
		dialog.ShowInformation("No Selection", "Please select a volume to inspect", gui.window)
		return
	}

	cmd := exec.Command("servin", "volume", "inspect", volumeInfo.Name)
	output, err := cmd.Output()
	if err != nil {
		dialog.ShowError(fmt.Errorf("Failed to inspect volume: %v", err), gui.window)
		return
	}

	// Show inspection in new window
	inspectWindow := gui.app.NewWindow(fmt.Sprintf("Inspect Volume - %s", volumeInfo.Name))
	inspectWindow.Resize(fyne.NewSize(800, 600))

	inspectText := widget.NewMultiLineEntry()
	inspectText.SetText(string(output))
	inspectText.Wrapping = fyne.TextWrapWord

	inspectWindow.SetContent(container.NewScroll(inspectText))
	inspectWindow.Show()
}

// CRI server management
func (gui *ServinDesktopGUI) startCRIServer() {
	// Start CRI server in background
	go func() {
		cmd := exec.Command("servin", "cri", "start", "--port", "10250")
		if err := cmd.Run(); err != nil {
			gui.updateStatus("Failed to start CRI server")
			return
		}

		gui.criServerRunning = true
		gui.updateStatus("CRI server started on port 10250")
		gui.refreshCRIStatus()
	}()
}

func (gui *ServinDesktopGUI) stopCRIServer() {
	// Stop CRI server
	go func() {
		cmd := exec.Command("servin", "cri", "stop")
		if err := cmd.Run(); err != nil {
			gui.updateStatus("Failed to stop CRI server")
			return
		}

		gui.criServerRunning = false
		gui.updateStatus("CRI server stopped")
		gui.refreshCRIStatus()
	}()
}

func (gui *ServinDesktopGUI) toggleCRIServer() {
	if gui.criServerRunning {
		gui.stopCRIServer()
	} else {
		gui.startCRIServer()
	}
}

// File management and export/import
func (gui *ServinDesktopGUI) exportImage() {
	imageInfo := gui.getSelectedImage()
	if imageInfo == nil {
		dialog.ShowInformation("No Selection", "Please select an image to export", gui.window)
		return
	}

	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, gui.window)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()

		// Export image
		cmd := exec.Command("servin", "save", "-o", writer.URI().Path(), imageInfo.ID)
		if err := cmd.Run(); err != nil {
			dialog.ShowError(fmt.Errorf("Failed to export image: %v", err), gui.window)
			return
		}

		gui.updateStatus(fmt.Sprintf("Image %s:%s exported", imageInfo.Name, imageInfo.Tag))
	}, gui.window)
}

func (gui *ServinDesktopGUI) importImage() {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, gui.window)
			return
		}
		if reader == nil {
			return
		}
		defer reader.Close()

		// Import image
		cmd := exec.Command("servin", "load", "-i", reader.URI().Path())
		if err := cmd.Run(); err != nil {
			dialog.ShowError(fmt.Errorf("Failed to import image: %v", err), gui.window)
			return
		}

		gui.updateStatus("Image imported successfully")
		gui.refreshImages()
	}, gui.window)
}

// Registry management
func (gui *ServinDesktopGUI) loginToRegistry() {
	// Create form for registry login
	registryEntry := widget.NewEntry()
	registryEntry.SetPlaceHolder("docker.io")

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("username")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("password")

	form := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Registry", registryEntry),
			widget.NewFormItem("Username", usernameEntry),
			widget.NewFormItem("Password", passwordEntry),
		},
		OnSubmit: func() {
			registry := registryEntry.Text
			if registry == "" {
				registry = "docker.io"
			}

			username := usernameEntry.Text
			password := passwordEntry.Text

			if username == "" || password == "" {
				dialog.ShowError(fmt.Errorf("Username and password are required"), gui.window)
				return
			}

			// Note: In a real implementation, you'd handle password securely
			cmd := exec.Command("servin", "login", "-u", username, "-p", password, registry)
			if err := cmd.Run(); err != nil {
				dialog.ShowError(fmt.Errorf("Failed to login: %v", err), gui.window)
				return
			}

			gui.updateStatus(fmt.Sprintf("Logged into %s", registry))
		},
	}

	dialog.ShowForm("Registry Login", "Login", "Cancel", form.Items, func(confirm bool) {
		if confirm {
			form.OnSubmit()
		}
	}, gui.window)
}

func (gui *ServinDesktopGUI) pushImage() {
	imageInfo := gui.getSelectedImage()
	if imageInfo == nil {
		dialog.ShowInformation("No Selection", "Please select an image to push", gui.window)
		return
	}

	// Create form for push destination
	tagEntry := widget.NewEntry()
	tagEntry.SetText(fmt.Sprintf("%s:%s", imageInfo.Name, imageInfo.Tag))

	form := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Tag", tagEntry),
		},
		OnSubmit: func() {
			tag := tagEntry.Text
			if tag == "" {
				dialog.ShowError(fmt.Errorf("Tag is required"), gui.window)
				return
			}

			// Show progress dialog
			progress := dialog.NewProgressInfinite("Pushing Image", fmt.Sprintf("Pushing %s...", tag), gui.window)
			progress.Show()

			go func() {
				cmd := exec.Command("servin", "push", tag)
				err := cmd.Run()
				progress.Hide()

				if err != nil {
					dialog.ShowError(fmt.Errorf("Failed to push image: %v", err), gui.window)
					return
				}

				gui.updateStatus(fmt.Sprintf("Image %s pushed successfully", tag))
			}()
		},
	}

	dialog.ShowForm("Push Image", "Push", "Cancel", form.Items, func(confirm bool) {
		if confirm {
			form.OnSubmit()
		}
	}, gui.window)
}
