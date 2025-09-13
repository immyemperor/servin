package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ServinDesktopGUI represents the main visual GUI application
type ServinDesktopGUI struct {
	app    fyne.App
	window fyne.Window

	// Data containers
	containers       []ContainerInfo
	images           []ImageInfo
	volumes          []VolumeInfo
	criServerRunning bool

	// UI Components
	containerList   *widget.List
	imageList       *widget.List
	volumeList      *widget.List
	logOutput       *widget.Entry
	statusBar       *widget.Label
	criStatusLabel  *widget.Label
	criToggleButton *widget.Button

	// Sidebar navigation
	sidebarTabs *container.AppTabs
	mainContent *fyne.Container

	// Container details
	containerDetailTabs *container.AppTabs
	containerLogs       *widget.Entry
	containerEnv        *widget.Entry
	containerFiles      *widget.Entry
	containerExec       *widget.Entry
	containerVolumes    *widget.Entry

	// Selection tracking
	selectedContainer int
	selectedImage     int
	selectedVolume    int
	currentSidebarTab string
	currentContainer  *ContainerInfo

	// Refresh timer
	refreshTimer *time.Ticker
}

// ContainerInfo represents container information for the GUI
type ContainerInfo struct {
	ID       string
	Name     string
	Image    string
	Tag      string
	Status   string
	State    string
	Ports    string
	Created  string
	Uptime   string
	PID      string
	Command  string
	Networks string
	Mounts   string
}

// ImageInfo represents image information for the GUI
type ImageInfo struct {
	ID      string
	Name    string
	Tag     string
	Size    string
	Created string
}

// VolumeInfo represents volume information for the GUI
type VolumeInfo struct {
	Name       string
	Driver     string
	Mountpoint string
	Created    string
}

// NewServinDesktopGUI creates a new GUI application
func NewServinDesktopGUI() *ServinDesktopGUI {
	myApp := app.NewWithID("com.servin.desktop")
	myApp.SetIcon(theme.ComputerIcon())

	gui := &ServinDesktopGUI{
		app:               myApp,
		window:            myApp.NewWindow("Servin Desktop"),
		selectedContainer: -1,
		selectedImage:     -1,
		selectedVolume:    -1,
		refreshTimer:      time.NewTicker(5 * time.Second),
	}

	return gui
}

// Run starts the GUI application
func (gui *ServinDesktopGUI) Run() {
	gui.setupWindow()
	gui.createMainLayout()
	// Use simple close instead of confirmation dialog for now
	gui.setupWindowEventsSimple()

	// Do initial refresh on main thread
	gui.refreshAllData()

	// TODO: Re-enable timer with proper threading later
	// gui.startRefreshTimer()

	gui.window.ShowAndRun()
}

// setupWindow configures the main window
func (gui *ServinDesktopGUI) setupWindow() {
	gui.window.SetTitle("Servin Desktop - Container Management")
	gui.window.Resize(fyne.NewSize(800, 600))

	// Set window properties to ensure proper window controls
	gui.window.SetFixedSize(false) // Allow resizing

	// Center the window on screen after all other properties are set
	gui.window.CenterOnScreen()
}

// setupWindowEvents configures window event handlers
func (gui *ServinDesktopGUI) setupWindowEvents() {
	// Set up close confirmation with proper callback handling
	gui.window.SetCloseIntercept(func() {
		// Create confirmation dialog
		confirm := dialog.NewConfirm("Exit Application",
			"Are you sure you want to exit Servin Desktop?",
			func(response bool) {
				if response {
					// Stop any timers first
					if gui.refreshTimer != nil {
						gui.refreshTimer.Stop()
					}
					// Then quit the application
					gui.app.Quit()
				}
				// If response is false, do nothing (stay open)
			}, gui.window)

		// Show the dialog
		confirm.Show()
	})

	// Add keyboard shortcuts
	gui.window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		// F5 to refresh
		if key.Name == fyne.KeyF5 {
			gui.refreshAllData()
		}
		// Escape to close confirmation dialogs
		if key.Name == fyne.KeyEscape {
			// This helps with dialog navigation
		}
	})
}

// Alternative: setupWindowEventsSimple - direct close without confirmation
func (gui *ServinDesktopGUI) setupWindowEventsSimple() {
	// Direct close without confirmation (for testing)
	gui.window.SetCloseIntercept(func() {
		// Stop any timers first
		if gui.refreshTimer != nil {
			gui.refreshTimer.Stop()
		}
		// Quit directly
		gui.app.Quit()
	})
}

// createMainLayout sets up the main application layout with vertical sidebar
func (gui *ServinDesktopGUI) createMainLayout() {
	// Create status bar
	gui.statusBar = widget.NewLabel("Ready")
	gui.criStatusLabel = widget.NewLabel("CRI: Stopped")
	gui.criToggleButton = widget.NewButton("Start CRI", func() {
		gui.toggleCRIServer()
	})

	statusContainer := container.NewHBox(
		gui.statusBar,
		widget.NewSeparator(),
		gui.criStatusLabel,
		gui.criToggleButton,
	)

	// Create vertical sidebar with tabs
	gui.sidebarTabs = container.NewAppTabs(
		container.NewTabItem("📦 Containers", gui.createContainerSidebar()),
		container.NewTabItem("💿 Images", gui.createImageTab()),
		container.NewTabItem("💾 Volumes", gui.createVolumeTab()),
		container.NewTabItem("⚙️ CRI Server", gui.createCRITab()),
		container.NewTabItem("📋 Logs", gui.createLogTab()),
	)

	// Set tabs to appear on the left (vertical)
	gui.sidebarTabs.SetTabLocation(container.TabLocationLeading)

	// Create main content area (initially empty)
	gui.mainContent = container.NewBorder(nil, nil, nil, nil,
		widget.NewLabel("Select a container to view details"))

	// Create main layout with sidebar and content
	mainLayout := container.NewHSplit(
		gui.sidebarTabs,
		gui.mainContent,
	)
	mainLayout.SetOffset(0.3) // 30% for sidebar, 70% for content

	// Create full layout with status bar
	content := container.NewBorder(
		nil,             // top
		statusContainer, // bottom
		nil,             // left
		nil,             // right
		mainLayout,      // center
	)

	gui.window.SetContent(content)
	gui.currentSidebarTab = "containers"
}

// createContainerSidebar creates the container list with detailed info for sidebar
func (gui *ServinDesktopGUI) createContainerSidebar() fyne.CanvasObject {
	// Create list widget with enhanced display
	gui.containerList = widget.NewList(
		func() int {
			return len(gui.containers)
		},
		func() fyne.CanvasObject {
			// Create a more detailed container item
			nameLabel := widget.NewLabel("Container Name")
			nameLabel.TextStyle.Bold = true

			tagLabel := widget.NewLabel("image:tag")
			tagLabel.TextStyle.Italic = true

			statusLabel := widget.NewLabel("Status")
			portsLabel := widget.NewLabel("Ports: 8080:80")
			uptimeLabel := widget.NewLabel("Uptime: 2h 30m")

			// Action menu button (3 dots)
			actionBtn := widget.NewButton("⋮", nil)
			actionBtn.Resize(fyne.NewSize(30, 30))

			// Container info in vertical layout
			infoBox := container.NewVBox(
				nameLabel,
				tagLabel,
				container.NewHBox(statusLabel, actionBtn),
				portsLabel,
				uptimeLabel,
			)

			return container.NewBorder(nil, nil, nil, nil, infoBox)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(gui.containers) {
				return
			}
			containerInfo := gui.containers[id]

			// Get the container from border layout
			border := obj.(*fyne.Container)
			infoBox := border.Objects[0].(*fyne.Container)

			nameLabel := infoBox.Objects[0].(*widget.Label)
			tagLabel := infoBox.Objects[1].(*widget.Label)
			statusRow := infoBox.Objects[2].(*fyne.Container)
			statusLabel := statusRow.Objects[0].(*widget.Label)
			actionBtn := statusRow.Objects[1].(*widget.Button)
			portsLabel := infoBox.Objects[3].(*widget.Label)
			uptimeLabel := infoBox.Objects[4].(*widget.Label)

			// Update labels with container info
			nameLabel.SetText(containerInfo.Name)
			tagLabel.SetText(fmt.Sprintf("%s:%s", containerInfo.Image, containerInfo.Tag))
			statusLabel.SetText(containerInfo.Status)
			portsLabel.SetText(fmt.Sprintf("Ports: %s", containerInfo.Ports))
			uptimeLabel.SetText(fmt.Sprintf("Uptime: %s", containerInfo.Uptime))

			// Set status color
			if containerInfo.Status == "running" {
				statusLabel.Importance = widget.SuccessImportance
			} else {
				statusLabel.Importance = widget.DangerImportance
			}

			// Setup action button with container-specific menu
			actionBtn.OnTapped = func() {
				gui.showContainerActionMenu(id, containerInfo)
			}
		},
	)

	gui.containerList.OnSelected = func(id widget.ListItemID) {
		gui.selectedContainer = id
		if id < len(gui.containers) {
			gui.currentContainer = &gui.containers[id]
			gui.showContainerDetails(gui.containers[id])
		}
	}

	// Refresh button at the bottom
	refreshBtn := widget.NewButton("🔄 Refresh Containers", func() {
		gui.refreshContainers()
	})

	return container.NewBorder(
		nil,               // top
		refreshBtn,        // bottom
		nil,               // left
		nil,               // right
		gui.containerList, // center
	)
}

// createImageTab creates the image management tab
func (gui *ServinDesktopGUI) createImageTab() fyne.CanvasObject {
	// Create list widget
	gui.imageList = widget.NewList(
		func() int {
			return len(gui.images)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.DocumentIcon()),
				widget.NewLabel("Image Name"),
				widget.NewLabel("Size"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(gui.images) {
				return
			}
			image := gui.images[id]
			hbox := obj.(*fyne.Container)
			nameLabel := hbox.Objects[1].(*widget.Label)
			sizeLabel := hbox.Objects[2].(*widget.Label)

			nameLabel.SetText(fmt.Sprintf("%s:%s", image.Name, image.Tag))
			sizeLabel.SetText(image.Size)
		},
	)

	gui.imageList.OnSelected = func(id widget.ListItemID) {
		gui.selectedImage = id
	}

	// Action buttons
	pullBtn := widget.NewButton("Pull", func() {
		gui.pullImage()
	})
	removeBtn := widget.NewButton("Remove", func() {
		gui.removeImage()
	})
	inspectBtn := widget.NewButton("Inspect", func() {
		gui.inspectImage()
	})
	refreshBtn := widget.NewButton("Refresh", func() {
		gui.refreshImages()
	})

	buttonContainer := container.NewHBox(
		pullBtn, removeBtn, inspectBtn, refreshBtn,
	)

	return container.NewBorder(
		nil,             // top
		buttonContainer, // bottom
		nil,             // left
		nil,             // right
		gui.imageList,   // center
	)
}

// createVolumeTab creates the volume management tab
func (gui *ServinDesktopGUI) createVolumeTab() fyne.CanvasObject {
	// Create list widget
	gui.volumeList = widget.NewList(
		func() int {
			return len(gui.volumes)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.StorageIcon()),
				widget.NewLabel("Volume Name"),
				widget.NewLabel("Driver"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(gui.volumes) {
				return
			}
			volume := gui.volumes[id]
			hbox := obj.(*fyne.Container)
			nameLabel := hbox.Objects[1].(*widget.Label)
			pathLabel := hbox.Objects[2].(*widget.Label)

			nameLabel.SetText(volume.Name)
			pathLabel.SetText(volume.Mountpoint)
		},
	)

	gui.volumeList.OnSelected = func(id widget.ListItemID) {
		gui.selectedVolume = id
	}

	// Action buttons
	createBtn := widget.NewButton("Create", func() {
		gui.createVolume()
	})
	removeBtn := widget.NewButton("Remove", func() {
		gui.removeVolume()
	})
	inspectBtn := widget.NewButton("Inspect", func() {
		gui.inspectVolume()
	})
	refreshBtn := widget.NewButton("Refresh", func() {
		gui.refreshVolumes()
	})

	buttonContainer := container.NewHBox(
		createBtn, removeBtn, inspectBtn, refreshBtn,
	)

	return container.NewBorder(
		nil,             // top
		buttonContainer, // bottom
		nil,             // left
		nil,             // right
		gui.volumeList,  // center
	)
}

// createCRITab creates the CRI server management tab
func (gui *ServinDesktopGUI) createCRITab() fyne.CanvasObject {
	statusLabel := widget.NewLabel("CRI Server Status:")
	statusValue := widget.NewLabel("Stopped")

	startBtn := widget.NewButton("Start CRI Server", func() {
		gui.startCRIServer()
	})
	stopBtn := widget.NewButton("Stop CRI Server", func() {
		gui.stopCRIServer()
	})

	portLabel := widget.NewLabel("Port:")
	portEntry := widget.NewEntry()
	portEntry.SetText("10250")

	infoText := widget.NewRichTextFromMarkdown(`
# CRI Server Configuration

The Container Runtime Interface (CRI) server provides Kubernetes compatibility.

**Default Configuration:**
- Port: 10250
- Protocol: HTTP/gRPC
- Endpoint: /runtime.v1alpha2.RuntimeService

**Usage:**
1. Start the CRI server
2. Configure kubelet to use this runtime
3. Deploy pods using Kubernetes
	`)

	form := container.NewVBox(
		container.NewHBox(statusLabel, statusValue),
		widget.NewSeparator(),
		container.NewHBox(portLabel, portEntry),
		container.NewHBox(startBtn, stopBtn),
		widget.NewSeparator(),
		container.NewScroll(infoText),
	)

	return form
}

// createLogTab creates the log viewing tab
func (gui *ServinDesktopGUI) createLogTab() fyne.CanvasObject {
	gui.logOutput = widget.NewMultiLineEntry()
	gui.logOutput.SetText("Application logs will appear here...\n")
	gui.logOutput.Wrapping = fyne.TextWrapWord

	clearBtn := widget.NewButton("Clear", func() {
		gui.logOutput.SetText("")
	})

	return container.NewBorder(
		nil,                                // top
		clearBtn,                           // bottom
		nil,                                // left
		nil,                                // right
		container.NewScroll(gui.logOutput), // center
	)
}

// startRefreshTimer starts the automatic refresh timer
func (gui *ServinDesktopGUI) startRefreshTimer() {
	go func() {
		for range gui.refreshTimer.C {
			// Only refresh if window is still open
			if gui.window != nil {
				go func() {
					gui.refreshContainers()
					gui.refreshImages()
					gui.refreshVolumes()
					gui.refreshCRIStatus()
				}()
			}
		}
	}()
}

// refreshAllData refreshes all data in the GUI
func (gui *ServinDesktopGUI) refreshAllData() {
	gui.refreshContainers()
	gui.refreshImages()
	gui.refreshVolumes()
	gui.refreshCRIStatus()
}

// refreshCRIStatus refreshes the CRI server status
func (gui *ServinDesktopGUI) refreshCRIStatus() {
	// Update status based on actual server state
	if gui.criServerRunning {
		if gui.criStatusLabel != nil {
			gui.criStatusLabel.SetText("CRI: Running")
		}
		if gui.criToggleButton != nil {
			gui.criToggleButton.SetText("Stop CRI")
		}
	} else {
		if gui.criStatusLabel != nil {
			gui.criStatusLabel.SetText("CRI: Stopped")
		}
		if gui.criToggleButton != nil {
			gui.criToggleButton.SetText("Start CRI")
		}
	}
}

// logMessage adds a message to the log output
func (gui *ServinDesktopGUI) logMessage(message string) {
	if gui.logOutput != nil {
		timestamp := time.Now().Format("15:04:05")
		gui.logOutput.SetText(gui.logOutput.Text + fmt.Sprintf("[%s] %s\n", timestamp, message))
	}
	if gui.statusBar != nil {
		gui.statusBar.SetText(message)
	}
}

// updateStatus updates the status bar (alias for logMessage for compatibility)
func (gui *ServinDesktopGUI) updateStatus(message string) {
	gui.logMessage(message)
}

// getSelectedContainer returns the currently selected container
func (gui *ServinDesktopGUI) getSelectedContainer() *ContainerInfo {
	if gui.selectedContainer >= 0 && gui.selectedContainer < len(gui.containers) {
		return &gui.containers[gui.selectedContainer]
	}
	return nil
}

// getSelectedImage returns the currently selected image
func (gui *ServinDesktopGUI) getSelectedImage() *ImageInfo {
	if gui.selectedImage >= 0 && gui.selectedImage < len(gui.images) {
		return &gui.images[gui.selectedImage]
	}
	return nil
}

// getSelectedVolume returns the currently selected volume
func (gui *ServinDesktopGUI) getSelectedVolume() *VolumeInfo {
	if gui.selectedVolume >= 0 && gui.selectedVolume < len(gui.volumes) {
		return &gui.volumes[gui.selectedVolume]
	}
	return nil
}

// showContainerActionMenu displays the action menu for a container
func (gui *ServinDesktopGUI) showContainerActionMenu(_ widget.ListItemID, containerInfo ContainerInfo) {
	// Create popup menu with container actions
	var menuItems []*fyne.MenuItem

	if containerInfo.Status == "running" {
		menuItems = append(menuItems,
			fyne.NewMenuItem("⏹️ Stop", func() {
				gui.stopContainerByID(containerInfo.ID)
			}),
			fyne.NewMenuItem("🔄 Restart", func() {
				gui.restartContainerByID(containerInfo.ID)
			}),
		)
	} else {
		menuItems = append(menuItems,
			fyne.NewMenuItem("▶️ Start", func() {
				gui.startContainerByID(containerInfo.ID)
			}),
		)
	}

	menuItems = append(menuItems,
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("🔍 Inspect", func() {
			gui.inspectContainerByID(containerInfo.ID)
		}),
		fyne.NewMenuItem("📋 Logs", func() {
			gui.showContainerLogsTab(containerInfo.ID)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("🗑️ Delete", func() {
			gui.deleteContainerByID(containerInfo.ID)
		}),
	)

	menu := fyne.NewMenu("Container Actions", menuItems...)
	popupMenu := widget.NewPopUpMenu(menu, gui.window.Canvas())
	popupMenu.ShowAtPosition(fyne.CurrentApp().Driver().AbsolutePositionForObject(gui.containerList))
}

// showContainerDetails displays detailed container information in the main content area
func (gui *ServinDesktopGUI) showContainerDetails(containerInfo ContainerInfo) {
	// Create container detail tabs
	gui.containerDetailTabs = container.NewAppTabs()

	// Overview tab
	overviewContent := gui.createContainerOverview(containerInfo)
	gui.containerDetailTabs.Append(container.NewTabItem("📊 Overview", overviewContent))

	// Logs tab
	gui.containerLogs = widget.NewMultiLineEntry()
	gui.containerLogs.SetText("Loading container logs...")
	gui.containerLogs.Wrapping = fyne.TextWrapWord
	logsScroll := container.NewScroll(gui.containerLogs)
	logsScroll.SetMinSize(fyne.NewSize(400, 300))

	logsRefreshBtn := widget.NewButton("🔄 Refresh Logs", func() {
		gui.refreshContainerLogs(containerInfo.ID)
	})
	logsContent := container.NewBorder(nil, logsRefreshBtn, nil, nil, logsScroll)
	gui.containerDetailTabs.Append(container.NewTabItem("📋 Logs", logsContent))

	// Exec tab
	gui.containerExec = widget.NewMultiLineEntry()
	gui.containerExec.SetPlaceHolder("Enter commands to execute in container...")
	execBtn := widget.NewButton("▶️ Execute", func() {
		gui.executeInContainer(containerInfo.ID, gui.containerExec.Text)
	})
	execContent := container.NewBorder(nil, execBtn, nil, nil, gui.containerExec)
	gui.containerDetailTabs.Append(container.NewTabItem("💻 Exec", execContent))

	// Files tab
	gui.containerFiles = widget.NewMultiLineEntry()
	gui.containerFiles.SetText("Container filesystem browser - Feature coming soon")
	gui.containerDetailTabs.Append(container.NewTabItem("📁 Files", gui.containerFiles))

	// Environment tab
	gui.containerEnv = widget.NewMultiLineEntry()
	gui.containerEnv.SetText("Loading environment variables...")
	gui.containerDetailTabs.Append(container.NewTabItem("🌍 Environment", gui.containerEnv))

	// Volumes tab
	gui.containerVolumes = widget.NewMultiLineEntry()
	gui.containerVolumes.SetText("Loading volume information...")
	gui.containerDetailTabs.Append(container.NewTabItem("💾 Volumes", gui.containerVolumes))

	// Update main content
	gui.mainContent.RemoveAll()
	gui.mainContent.Add(gui.containerDetailTabs)

	// Load initial data
	gui.refreshContainerLogs(containerInfo.ID)
	gui.refreshContainerEnv(containerInfo.ID)
	gui.refreshContainerVolumes(containerInfo.ID)
}

// createContainerOverview creates the overview tab content for a container
func (gui *ServinDesktopGUI) createContainerOverview(containerInfo ContainerInfo) fyne.CanvasObject {
	// Container info grid
	infoGrid := container.NewGridWithColumns(2,
		widget.NewLabel("Name:"), widget.NewLabel(containerInfo.Name),
		widget.NewLabel("ID:"), widget.NewLabel(containerInfo.ID[:12]+"..."),
		widget.NewLabel("Image:"), widget.NewLabel(fmt.Sprintf("%s:%s", containerInfo.Image, containerInfo.Tag)),
		widget.NewLabel("Status:"), widget.NewLabel(containerInfo.Status),
		widget.NewLabel("Ports:"), widget.NewLabel(containerInfo.Ports),
		widget.NewLabel("Uptime:"), widget.NewLabel(containerInfo.Uptime),
		widget.NewLabel("PID:"), widget.NewLabel(containerInfo.PID),
		widget.NewLabel("Command:"), widget.NewLabel(containerInfo.Command),
		widget.NewLabel("Networks:"), widget.NewLabel(containerInfo.Networks),
	)

	// Action buttons
	var actionButtons *fyne.Container
	if containerInfo.Status == "running" {
		actionButtons = container.NewHBox(
			widget.NewButton("⏹️ Stop", func() { gui.stopContainerByID(containerInfo.ID) }),
			widget.NewButton("🔄 Restart", func() { gui.restartContainerByID(containerInfo.ID) }),
			widget.NewButton("⏸️ Pause", func() { gui.pauseContainerByID(containerInfo.ID) }),
		)
	} else {
		actionButtons = container.NewHBox(
			widget.NewButton("▶️ Start", func() { gui.startContainerByID(containerInfo.ID) }),
			widget.NewButton("🗑️ Delete", func() { gui.deleteContainerByID(containerInfo.ID) }),
		)
	}

	return container.NewVBox(
		widget.NewCard("Container Information", "", infoGrid),
		widget.NewCard("Actions", "", actionButtons),
	)
}

// Container action functions (placeholder implementations)
func (gui *ServinDesktopGUI) startContainerByID(id string) {
	gui.logMessage(fmt.Sprintf("Starting container %s...", id[:12]))
	// TODO: Implement actual container start
}

func (gui *ServinDesktopGUI) stopContainerByID(id string) {
	gui.logMessage(fmt.Sprintf("Stopping container %s...", id[:12]))
	// TODO: Implement actual container stop
}

func (gui *ServinDesktopGUI) restartContainerByID(id string) {
	gui.logMessage(fmt.Sprintf("Restarting container %s...", id[:12]))
	// TODO: Implement actual container restart
}

func (gui *ServinDesktopGUI) pauseContainerByID(id string) {
	gui.logMessage(fmt.Sprintf("Pausing container %s...", id[:12]))
	// TODO: Implement actual container pause
}

func (gui *ServinDesktopGUI) deleteContainerByID(id string) {
	// Show confirmation dialog
	confirm := dialog.NewConfirm("Delete Container",
		fmt.Sprintf("Are you sure you want to delete container %s?", id[:12]),
		func(confirmed bool) {
			if confirmed {
				gui.logMessage(fmt.Sprintf("Deleting container %s...", id[:12]))
				// TODO: Implement actual container deletion
			}
		}, gui.window)
	confirm.Show()
}

func (gui *ServinDesktopGUI) inspectContainerByID(id string) {
	gui.logMessage(fmt.Sprintf("Inspecting container %s...", id[:12]))
	// TODO: Implement container inspection
}

func (gui *ServinDesktopGUI) showContainerLogsTab(id string) {
	// Switch to logs tab if container details are shown
	if gui.containerDetailTabs != nil {
		gui.containerDetailTabs.SelectTab(gui.containerDetailTabs.Items[1]) // Logs tab
		gui.refreshContainerLogs(id)
	}
}

func (gui *ServinDesktopGUI) refreshContainerLogs(id string) {
	if gui.containerLogs != nil {
		gui.containerLogs.SetText(fmt.Sprintf("Logs for container %s:\n\n[INFO] Container started\n[INFO] Application running on port 8080\n[DEBUG] Processing request...\n", id[:12]))
	}
}

func (gui *ServinDesktopGUI) refreshContainerEnv(id string) {
	if gui.containerEnv != nil {
		gui.containerEnv.SetText(fmt.Sprintf("Environment variables for container %s:\n\nPATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\nHOSTNAME=%s\nTERM=xterm\n", id[:12], id[:8]))
	}
}

func (gui *ServinDesktopGUI) refreshContainerVolumes(id string) {
	if gui.containerVolumes != nil {
		gui.containerVolumes.SetText(fmt.Sprintf("Volume mounts for container %s:\n\n/var/lib/data -> /host/data (rw)\n/etc/config -> /host/config (ro)\n", id[:12]))
	}
}

func (gui *ServinDesktopGUI) executeInContainer(id string, command string) {
	if command == "" {
		gui.logMessage("Please enter a command to execute")
		return
	}
	gui.logMessage(fmt.Sprintf("Executing '%s' in container %s...", command, id[:12]))
	// TODO: Implement actual command execution
}

func main() {
	gui := NewServinDesktopGUI()
	gui.Run()
}
