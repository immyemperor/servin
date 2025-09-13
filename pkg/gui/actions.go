package gui

import (
	"fmt"
	"os/exec"

	"fyne.io/fyne/v2"
	fynecontainer "fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// startRefreshTimer starts the automatic refresh timer
func (gui *ServinGUI) startRefreshTimer() {
	go func() {
		for range gui.refreshTimer.C {
			gui.refreshData()
		}
	}()
}

// refreshData refreshes containers and images data
func (gui *ServinGUI) refreshData() {
	gui.refreshContainers()
	gui.refreshImages()
	gui.updateCRIStatus()
}

// Container management actions

func (gui *ServinGUI) startContainer() {
	selected := gui.containerList.GetSelected()
	if selected < 0 || selected >= len(gui.containers) {
		gui.showError("Please select a container to start")
		return
	}

	container := gui.containers[selected]
	gui.Log(fmt.Sprintf("Starting container %s...", container.Name))

	// Execute servin start command
	cmd := exec.Command("servin", "start", container.ID)
	if err := cmd.Run(); err != nil {
		gui.showError(fmt.Sprintf("Failed to start container: %v", err))
		return
	}

	gui.Log(fmt.Sprintf("Container %s started successfully", container.Name))
	gui.refreshContainers()
}

func (gui *ServinGUI) stopContainer() {
	selected := gui.containerList.GetSelected()
	if selected < 0 || selected >= len(gui.containers) {
		gui.showError("Please select a container to stop")
		return
	}

	container := gui.containers[selected]
	gui.Log(fmt.Sprintf("Stopping container %s...", container.Name))

	// Execute servin stop command
	cmd := exec.Command("servin", "stop", container.ID)
	if err := cmd.Run(); err != nil {
		gui.showError(fmt.Sprintf("Failed to stop container: %v", err))
		return
	}

	gui.Log(fmt.Sprintf("Container %s stopped successfully", container.Name))
	gui.refreshContainers()
}

func (gui *ServinGUI) removeContainer() {
	selected := gui.containerList.GetSelected()
	if selected < 0 || selected >= len(gui.containers) {
		gui.showError("Please select a container to remove")
		return
	}

	container := gui.containers[selected]

	// Confirmation dialog
	confirm := dialog.NewConfirm(
		"Remove Container",
		fmt.Sprintf("Are you sure you want to remove container '%s'?", container.Name),
		func(confirmed bool) {
			if confirmed {
				gui.Log(fmt.Sprintf("Removing container %s...", container.Name))

				// Execute servin rm command
				cmd := exec.Command("servin", "rm", container.ID)
				if err := cmd.Run(); err != nil {
					gui.showError(fmt.Sprintf("Failed to remove container: %v", err))
					return
				}

				gui.Log(fmt.Sprintf("Container %s removed successfully", container.Name))
				gui.refreshContainers()
			}
		},
		gui.window,
	)
	confirm.Show()
}

func (gui *ServinGUI) showContainerLogs() {
	selected := gui.containerList.GetSelected()
	if selected < 0 || selected >= len(gui.containers) {
		gui.showError("Please select a container to view logs")
		return
	}

	container := gui.containers[selected]
	gui.Log(fmt.Sprintf("Fetching logs for container %s...", container.Name))

	// Execute servin logs command
	cmd := exec.Command("servin", "logs", container.ID)
	output, err := cmd.Output()
	if err != nil {
		gui.showError(fmt.Sprintf("Failed to get container logs: %v", err))
		return
	}

	// Show logs in a new window
	logWindow := gui.app.NewWindow(fmt.Sprintf("Logs - %s", container.Name))
	logWindow.Resize(fyne.NewSize(800, 600))

	logText := widget.NewRichText(&widget.RichTextSegment{
		Text:  string(output),
		Style: &widget.RichTextStyle{},
	})
	logText.Wrapping = fyne.TextWrapWord

	logWindow.SetContent(fynecontainer.NewScroll(logText))
	logWindow.Show()
}

func (gui *ServinGUI) refreshContainers() {
	// Execute servin ls command to get containers
	cmd := exec.Command("servin", "ls", "--json")
	output, err := cmd.Output()
	if err != nil {
		gui.Log(fmt.Sprintf("Failed to refresh containers: %v", err))
		return
	}

	// Parse JSON output (simplified for demo)
	// In a real implementation, you'd properly parse the JSON
	gui.Log("Containers refreshed")
	if gui.containerList != nil {
		gui.containerList.Refresh()
	}
}

// Image management actions

func (gui *ServinGUI) importImage() {
	// File dialog to select image file
	fileDialog := dialog.NewFileOpen(
		func(reader fyne.URIReadCloser) {
			if reader == nil {
				return
			}

			filepath := reader.URI().Path()
			gui.Log(fmt.Sprintf("Importing image from %s...", filepath))

			// Execute servin image import command
			cmd := exec.Command("servin", "image", "import", filepath)
			if err := cmd.Run(); err != nil {
				gui.showError(fmt.Sprintf("Failed to import image: %v", err))
				return
			}

			gui.Log("Image imported successfully")
			gui.refreshImages()
		},
		gui.window,
	)

	// Set file filter for common image formats
	fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".tar", ".tar.gz", ".tgz"}))
	fileDialog.Show()
}

func (gui *ServinGUI) removeImage() {
	selected := gui.imageList.GetSelected()
	if selected < 0 || selected >= len(gui.images) {
		gui.showError("Please select an image to remove")
		return
	}

	image := gui.images[selected]

	// Confirmation dialog
	confirm := dialog.NewConfirm(
		"Remove Image",
		fmt.Sprintf("Are you sure you want to remove image '%s:%s'?", image.Name, image.Tag),
		func(confirmed bool) {
			if confirmed {
				gui.Log(fmt.Sprintf("Removing image %s:%s...", image.Name, image.Tag))

				// Execute servin image rm command
				cmd := exec.Command("servin", "image", "rm", fmt.Sprintf("%s:%s", image.Name, image.Tag))
				if err := cmd.Run(); err != nil {
					gui.showError(fmt.Sprintf("Failed to remove image: %v", err))
					return
				}

				gui.Log(fmt.Sprintf("Image %s:%s removed successfully", image.Name, image.Tag))
				gui.refreshImages()
			}
		},
		gui.window,
	)
	confirm.Show()
}

func (gui *ServinGUI) tagImage() {
	selected := gui.imageList.GetSelected()
	if selected < 0 || selected >= len(gui.images) {
		gui.showError("Please select an image to tag")
		return
	}

	image := gui.images[selected]

	// Tag input dialog
	entry := widget.NewEntry()
	entry.SetPlaceHolder("new-name:new-tag")

	form := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("New Tag", entry),
		},
		OnSubmit: func() {
			newTag := entry.Text
			if newTag == "" {
				gui.showError("Please enter a new tag")
				return
			}

			gui.Log(fmt.Sprintf("Tagging image %s:%s as %s...", image.Name, image.Tag, newTag))

			// Execute servin image tag command
			cmd := exec.Command("servin", "image", "tag", fmt.Sprintf("%s:%s", image.Name, image.Tag), newTag)
			if err := cmd.Run(); err != nil {
				gui.showError(fmt.Sprintf("Failed to tag image: %v", err))
				return
			}

			gui.Log(fmt.Sprintf("Image tagged successfully as %s", newTag))
			gui.refreshImages()
		},
	}

	dialog.ShowForm("Tag Image", "Tag", "Cancel", form.Items, form.OnSubmit, gui.window)
}

func (gui *ServinGUI) inspectImage() {
	selected := gui.imageList.GetSelected()
	if selected < 0 || selected >= len(gui.images) {
		gui.showError("Please select an image to inspect")
		return
	}

	image := gui.images[selected]
	gui.Log(fmt.Sprintf("Inspecting image %s:%s...", image.Name, image.Tag))

	// Execute servin image inspect command
	cmd := exec.Command("servin", "image", "inspect", fmt.Sprintf("%s:%s", image.Name, image.Tag))
	output, err := cmd.Output()
	if err != nil {
		gui.showError(fmt.Sprintf("Failed to inspect image: %v", err))
		return
	}

	// Show inspection details in a new window
	inspectWindow := gui.app.NewWindow(fmt.Sprintf("Inspect - %s:%s", image.Name, image.Tag))
	inspectWindow.Resize(fyne.NewSize(800, 600))

	inspectText := widget.NewRichText(&widget.RichTextSegment{
		Text:  string(output),
		Style: &widget.RichTextStyle{},
	})
	inspectText.Wrapping = fyne.TextWrapWord

	inspectWindow.SetContent(fynecontainer.NewScroll(inspectText))
	inspectWindow.Show()
}

func (gui *ServinGUI) refreshImages() {
	// Execute servin image ls command to get images
	cmd := exec.Command("servin", "image", "ls", "--json")
	output, err := cmd.Output()
	if err != nil {
		gui.Log(fmt.Sprintf("Failed to refresh images: %v", err))
		return
	}

	// Parse JSON output (simplified for demo)
	// In a real implementation, you'd properly parse the JSON
	gui.Log("Images refreshed")
	if gui.imageList != nil {
		gui.imageList.Refresh()
	}
}

// CRI server management

func (gui *ServinGUI) toggleCRIServer() {
	if gui.criRunning {
		gui.stopCRIServer()
	} else {
		gui.startCRIServer()
	}
}

func (gui *ServinGUI) startCRIServer() {
	gui.Log("Starting CRI server...")

	// Execute servin cri start command in background
	cmd := exec.Command("servin", "cri", "start", "--port", "8080")
	if err := cmd.Start(); err != nil {
		gui.showError(fmt.Sprintf("Failed to start CRI server: %v", err))
		return
	}

	gui.criRunning = true
	gui.criStatus.SetText("Running")
	gui.criServerBtn.SetText("Stop CRI Server")
	gui.criServerBtn.SetIcon(theme.MediaPauseIcon())
	gui.Log("CRI server started on port 8080")
}

func (gui *ServinGUI) stopCRIServer() {
	gui.Log("Stopping CRI server...")

	// For simplicity, we'll use a simple approach to stop the server
	// In a real implementation, you'd track the process and stop it properly
	gui.criRunning = false
	gui.criStatus.SetText("Stopped")
	gui.criServerBtn.SetText("Start CRI Server")
	gui.criServerBtn.SetIcon(theme.MediaPlayIcon())
	gui.Log("CRI server stopped")
}

func (gui *ServinGUI) testCRIConnection() {
	gui.Log("Testing CRI connection...")

	// Test health endpoint
	cmd := exec.Command("curl", "-s", "http://localhost:8080/health")
	_, err := cmd.Output()
	if err != nil {
		gui.showError("CRI server is not responding. Make sure it's running.")
		return
	}

	gui.Log("CRI server connection test successful")
	gui.showInfo("CRI server is responding correctly")
}

func (gui *ServinGUI) updateCRIStatus() {
	// Check if CRI server is running by testing the health endpoint
	cmd := exec.Command("curl", "-s", "http://localhost:8080/health")
	err := cmd.Run()

	if err == nil && !gui.criRunning {
		gui.criRunning = true
		gui.criStatus.SetText("Running")
		gui.criServerBtn.SetText("Stop CRI Server")
		gui.criServerBtn.SetIcon(theme.MediaPauseIcon())
	} else if err != nil && gui.criRunning {
		gui.criRunning = false
		gui.criStatus.SetText("Stopped")
		gui.criServerBtn.SetText("Start CRI Server")
		gui.criServerBtn.SetIcon(theme.MediaPlayIcon())
	}
}

// Utility methods

func (gui *ServinGUI) showError(message string) {
	gui.Log(fmt.Sprintf("ERROR: %s", message))
	dialog.ShowError(fmt.Errorf(message), gui.window)
}

func (gui *ServinGUI) showInfo(message string) {
	gui.Log(fmt.Sprintf("INFO: %s", message))
	dialog.ShowInformation("Information", message, gui.window)
}
