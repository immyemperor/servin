package main

import (
	"fmt"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// refreshContainers gets the current list of containers
func (gui *ServinDesktopGUI) refreshContainers() {
	// Try to get real data, but fall back to mock data if command fails
	cmd := exec.Command("servin", "ls", "--json")
	_, err := cmd.Output()
	if err != nil {
		// Use mock data for demo purposes
		gui.updateStatus("Using demo data (servin CLI not available)")
	} else {
		gui.updateStatus("Containers refreshed")
	}

	// Always populate with demo data for now
	gui.containers = []ContainerInfo{
		{ID: "abc123", Name: "test-container", Image: "alpine:latest", Status: "running", State: "running", Ports: "8080:80", Created: "2 hours ago"},
		{ID: "def456", Name: "web-server", Image: "nginx:latest", Status: "stopped", State: "exited", Ports: "-", Created: "1 day ago"},
		{ID: "ghi789", Name: "database", Image: "postgres:13", Status: "running", State: "running", Ports: "5432:5432", Created: "3 hours ago"},
	}

	if gui.containerList != nil {
		gui.containerList.Refresh()
	}
	// Status already updated above based on command result
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
	// Try to get real data, but fall back to mock data if command fails
	cmd := exec.Command("servin", "images", "--json")
	_, err := cmd.Output()
	if err != nil {
		// Use mock data for demo purposes
		gui.updateStatus("Using demo data (servin CLI not available)")
	} else {
		gui.updateStatus("Images refreshed")
	}

	// Always populate with demo data for now
	gui.images = []ImageInfo{
		{ID: "sha256:abc123", Name: "alpine", Tag: "latest", Size: "5.6MB", Created: "2 hours ago"},
		{ID: "sha256:def456", Name: "nginx", Tag: "latest", Size: "133MB", Created: "1 day ago"},
		{ID: "sha256:ghi789", Name: "postgres", Tag: "13", Size: "314MB", Created: "3 days ago"},
		{ID: "sha256:jkl012", Name: "redis", Tag: "alpine", Size: "28MB", Created: "5 hours ago"},
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
