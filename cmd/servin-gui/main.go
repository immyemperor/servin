package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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

	// Selection tracking
	selectedContainer int
	selectedImage     int
	selectedVolume    int

	// Refresh timer
	refreshTimer *time.Ticker
}

// ContainerInfo represents container information for the GUI
type ContainerInfo struct {
	ID      string
	Name    string
	Image   string
	Status  string
	State   string
	Ports   string
	Created string
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

	// Do initial refresh on main thread
	gui.refreshAllData()

	// TODO: Re-enable timer with proper threading later
	// gui.startRefreshTimer()

	gui.window.ShowAndRun()
}

// setupWindow configures the main window
func (gui *ServinDesktopGUI) setupWindow() {
	gui.window.Resize(fyne.NewSize(1400, 900))
	gui.window.SetMaster()
	gui.window.CenterOnScreen()
}

// createMainLayout sets up the main application layout
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

	// Create main tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Containers", gui.createContainerTab()),
		container.NewTabItem("Images", gui.createImageTab()),
		container.NewTabItem("Volumes", gui.createVolumeTab()),
		container.NewTabItem("CRI Server", gui.createCRITab()),
		container.NewTabItem("Logs", gui.createLogTab()),
	)

	// Create main layout
	content := container.NewBorder(
		nil,             // top
		statusContainer, // bottom
		nil,             // left
		nil,             // right
		tabs,            // center
	)

	gui.window.SetContent(content)
}

// createContainerTab creates the container management tab
func (gui *ServinDesktopGUI) createContainerTab() fyne.CanvasObject {
	// Create list widget
	gui.containerList = widget.NewList(
		func() int {
			return len(gui.containers)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.DocumentIcon()),
				widget.NewLabel("Container Name"),
				widget.NewLabel("Status"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(gui.containers) {
				return
			}
			container_info := gui.containers[id]
			hbox := obj.(*fyne.Container)
			nameLabel := hbox.Objects[1].(*widget.Label)
			statusLabel := hbox.Objects[2].(*widget.Label)

			nameLabel.SetText(fmt.Sprintf("%s (%s)", container_info.Name, container_info.Image))
			statusLabel.SetText(container_info.Status)
		},
	)

	gui.containerList.OnSelected = func(id widget.ListItemID) {
		gui.selectedContainer = id
	}

	// Action buttons
	startBtn := widget.NewButton("Start", func() {
		gui.startContainer()
	})
	stopBtn := widget.NewButton("Stop", func() {
		gui.stopContainer()
	})
	restartBtn := widget.NewButton("Restart", func() {
		gui.restartContainer()
	})
	removeBtn := widget.NewButton("Remove", func() {
		gui.removeContainer()
	})
	inspectBtn := widget.NewButton("Inspect", func() {
		gui.inspectContainer()
	})
	refreshBtn := widget.NewButton("Refresh", func() {
		gui.refreshContainers()
	})

	buttonContainer := container.NewHBox(
		startBtn, stopBtn, restartBtn, removeBtn, inspectBtn, refreshBtn,
	)

	return container.NewBorder(
		nil,               // top
		buttonContainer,   // bottom
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

func main() {
	gui := NewServinDesktopGUI()
	gui.Run()
}
