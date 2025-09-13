package gui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	fynecontainer "fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"servin/pkg/container"
	"servin/pkg/image"
	"servin/pkg/state"
)

// ServinGUI represents the main GUI application
type ServinGUI struct {
	app           fyne.App
	window        fyne.Window
	containerList *widget.List
	imageList     *widget.List
	logOutput     *widget.RichText
	statusBar     *widget.Label
	criServerBtn  *widget.Button
	criStatus     *widget.Label
	refreshTimer  *time.Ticker

	// Data
	containers []state.Container
	images     []image.Image
	criRunning bool
}

// NewServinGUI creates a new GUI application
func NewServinGUI() *ServinGUI {
	myApp := app.NewWithID("com.servin.desktop")
	myApp.SetIcon(theme.ComputerIcon())

	return &ServinGUI{
		app:          myApp,
		window:       myApp.NewWindow("Servin Desktop"),
		refreshTimer: time.NewTicker(2 * time.Second),
	}
}

// Run starts the GUI application
func (gui *ServinGUI) Run() {
	gui.setupWindow()
	gui.createLayout()
	gui.startRefreshTimer()

	gui.window.ShowAndRun()
}

// setupWindow configures the main window
func (gui *ServinGUI) setupWindow() {
	gui.window.Resize(fyne.NewSize(1200, 800))
	gui.window.CenterOnScreen()
	gui.window.SetMaster()
}

// createLayout builds the main UI layout
func (gui *ServinGUI) createLayout() {
	// Create tab container for main sections
	tabs := fynecontainer.NewAppTabs(
		fynecontainer.NewTabItem("Containers", gui.createContainerTab()),
		fynecontainer.NewTabItem("Images", gui.createImageTab()),
		fynecontainer.NewTabItem("CRI Server", gui.createCRITab()),
		fynecontainer.NewTabItem("Logs", gui.createLogTab()),
	)

	// Create status bar
	gui.statusBar = widget.NewLabel("Ready")
	gui.statusBar.TextStyle = fyne.TextStyle{Italic: true}

	// Main layout with status bar at bottom
	content := fynecontainer.NewBorder(
		nil,           // top
		gui.statusBar, // bottom
		nil,           // left
		nil,           // right
		tabs,          // center
	)

	gui.window.SetContent(content)
}

// createContainerTab creates the container management tab
func (gui *ServinGUI) createContainerTab() fyne.CanvasObject {
	// Container list
	gui.containerList = widget.NewList(
		func() int {
			return len(gui.containers)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.MediaPlayIcon()),
				widget.NewLabel("Container Name"),
				widget.NewLabel("Status"),
				widget.NewLabel("Image"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(gui.containers) {
				return
			}

			cont := gui.containers[id]
			hbox := obj.(*container.Box)

			// Update icon based on status
			icon := hbox.Objects[0].(*widget.Icon)
			if cont.Status == "running" {
				icon.SetResource(theme.MediaPlayIcon())
			} else {
				icon.SetResource(theme.MediaPauseIcon())
			}

			// Update labels
			nameLabel := hbox.Objects[1].(*widget.Label)
			statusLabel := hbox.Objects[2].(*widget.Label)
			imageLabel := hbox.Objects[3].(*widget.Label)

			nameLabel.SetText(cont.Name)
			statusLabel.SetText(cont.Status)
			imageLabel.SetText(cont.Image)
		},
	)

	// Container action buttons
	startBtn := widget.NewButtonWithIcon("Start", theme.MediaPlayIcon(), gui.startContainer)
	stopBtn := widget.NewButtonWithIcon("Stop", theme.MediaPauseIcon(), gui.stopContainer)
	removeBtn := widget.NewButtonWithIcon("Remove", theme.DeleteIcon(), gui.removeContainer)
	logsBtn := widget.NewButtonWithIcon("Logs", theme.DocumentIcon(), gui.showContainerLogs)
	refreshBtn := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), gui.refreshContainers)

	actionBar := container.NewHBox(
		startBtn, stopBtn, removeBtn, logsBtn,
		layout.NewSpacer(),
		refreshBtn,
	)

	return container.NewBorder(
		nil,               // top
		actionBar,         // bottom
		nil,               // left
		nil,               // right
		gui.containerList, // center
	)
}

// createImageTab creates the image management tab
func (gui *ServinGUI) createImageTab() fyne.CanvasObject {
	// Image list
	gui.imageList = widget.NewList(
		func() int {
			return len(gui.images)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.StorageIcon()),
				widget.NewLabel("Image Name"),
				widget.NewLabel("Tag"),
				widget.NewLabel("Size"),
				widget.NewLabel("Created"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(gui.images) {
				return
			}

			img := gui.images[id]
			hbox := obj.(*container.Box)

			// Update labels
			nameLabel := hbox.Objects[1].(*widget.Label)
			tagLabel := hbox.Objects[2].(*widget.Label)
			sizeLabel := hbox.Objects[3].(*widget.Label)
			createdLabel := hbox.Objects[4].(*widget.Label)

			nameLabel.SetText(img.Name)
			tagLabel.SetText(img.Tag)
			sizeLabel.SetText(formatSize(img.Size))
			createdLabel.SetText(img.Created.Format("2006-01-02"))
		},
	)

	// Image action buttons
	importBtn := widget.NewButtonWithIcon("Import", theme.FolderOpenIcon(), gui.importImage)
	removeBtn := widget.NewButtonWithIcon("Remove", theme.DeleteIcon(), gui.removeImage)
	tagBtn := widget.NewButtonWithIcon("Tag", theme.ContentAddIcon(), gui.tagImage)
	inspectBtn := widget.NewButtonWithIcon("Inspect", theme.InfoIcon(), gui.inspectImage)
	refreshBtn := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), gui.refreshImages)

	actionBar := container.NewHBox(
		importBtn, removeBtn, tagBtn, inspectBtn,
		layout.NewSpacer(),
		refreshBtn,
	)

	return container.NewBorder(
		nil,           // top
		actionBar,     // bottom
		nil,           // left
		nil,           // right
		gui.imageList, // center
	)
}

// createCRITab creates the CRI server management tab
func (gui *ServinGUI) createCRITab() fyne.CanvasObject {
	// CRI server status
	gui.criStatus = widget.NewLabel("Stopped")
	gui.criStatus.TextStyle = fyne.TextStyle{Bold: true}

	statusCard := widget.NewCard("CRI Server Status", "",
		container.NewVBox(
			gui.criStatus,
			widget.NewSeparator(),
			widget.NewLabel("Port: 8080"),
			widget.NewLabel("Endpoint: http://localhost:8080"),
		),
	)

	// CRI server controls
	gui.criServerBtn = widget.NewButtonWithIcon("Start CRI Server", theme.MediaPlayIcon(), gui.toggleCRIServer)
	testBtn := widget.NewButtonWithIcon("Test Connection", theme.ConfirmIcon(), gui.testCRIConnection)

	controlCard := widget.NewCard("Controls", "",
		container.NewVBox(
			gui.criServerBtn,
			testBtn,
		),
	)

	// CRI endpoints info
	endpointsInfo := widget.NewRichTextFromMarkdown(`
## Available Endpoints

### Runtime Operations
- POST /v1/runtime/version
- POST /v1/runtime/status

### Pod Sandbox Operations
- POST /v1/runtime/sandbox/list
- POST /v1/runtime/sandbox/create
- POST /v1/runtime/sandbox/start
- POST /v1/runtime/sandbox/stop
- POST /v1/runtime/sandbox/remove

### Container Operations
- POST /v1/runtime/container/list
- POST /v1/runtime/container/create
- POST /v1/runtime/container/start
- POST /v1/runtime/container/stop
- POST /v1/runtime/container/remove
- POST /v1/runtime/container/status

### Image Operations
- POST /v1/image/list
- POST /v1/image/status
- POST /v1/image/pull
- POST /v1/image/remove
- POST /v1/image/fs

### Health Check
- GET /health
`)

	endpointsCard := widget.NewCard("CRI Endpoints", "", endpointsInfo)

	return container.NewHSplit(
		container.NewVBox(statusCard, controlCard),
		endpointsCard,
	)
}

// createLogTab creates the log viewing tab
func (gui *ServinGUI) createLogTab() fyne.CanvasObject {
	gui.logOutput = widget.NewRichText(&widget.RichTextSegment{
		Text:  "Servin Desktop started...\n",
		Style: &widget.RichTextStyle{},
	})
	gui.logOutput.Wrapping = fyne.TextWrapWord

	clearBtn := widget.NewButtonWithIcon("Clear Logs", theme.ContentClearIcon(), func() {
		gui.logOutput.Segments = []*widget.RichTextSegment{}
		gui.logOutput.Refresh()
	})

	return container.NewBorder(
		nil,                                // top
		clearBtn,                           // bottom
		nil,                                // left
		nil,                                // right
		container.NewScroll(gui.logOutput), // center
	)
}

// Utility function to format file sizes
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Log adds a message to the log output
func (gui *ServinGUI) Log(message string) {
	timestamp := time.Now().Format("15:04:05")
	logMessage := fmt.Sprintf("[%s] %s\n", timestamp, message)

	gui.logOutput.Segments = append(gui.logOutput.Segments, &widget.RichTextSegment{
		Text:  logMessage,
		Style: &widget.RichTextStyle{},
	})
	gui.logOutput.Refresh()

	// Update status bar
	gui.statusBar.SetText(message)
}
