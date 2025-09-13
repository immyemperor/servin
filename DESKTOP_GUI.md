# Servin Desktop GUI

## Overview

Servin Desktop provides both graphical (GUI) and terminal (TUI) user interfaces for managing containers, images, CRI server, volumes, and registry operations. The interface is designed to be similar to Docker Desktop, providing an intuitive way to work with Servin container runtime.

## Features

### 🖥️ **Terminal User Interface (TUI)**
- **Cross-platform**: Works on Windows, Linux, and macOS
- **No dependencies**: Built with Go's standard library
- **Full functionality**: Access to all Servin features
- **Interactive menus**: Easy navigation with numbered options
- **Real-time feedback**: Command output displayed directly

### 🎨 **Graphical User Interface (GUI)**
- **Modern interface**: Using Fyne framework for native look and feel
- **Cross-platform**: Single binary for all platforms
- **Visual management**: Point-and-click container and image management
- **Real-time updates**: Automatic refresh of container/image status
- **Integrated logs**: Built-in log viewer for containers

## Usage

### Launching the Interface

```bash
# Launch GUI (if available)
servin gui

# Force Terminal UI
servin gui --tui

# Launch in development mode
servin gui --dev

# Launch GUI on custom port
servin gui --port 9090

# Get help
servin gui --help
```

### TUI Navigation

The Terminal UI uses numbered menus for navigation:

1. **Main Menu**
   - Container Management
   - Image Management  
   - CRI Server Control
   - Volume Management
   - Registry Operations
   - System Information
   - Exit

2. **Navigation**
   - Enter the number of your choice
   - Use option 8 (or similar) to go back to previous menu
   - Use option 7 in main menu to exit

### TUI Features

#### Container Management
```
┌─────────────── Container Management ──────────────┐
│  1. List Containers                               │
│  2. Run New Container                             │
│  3. Start Container                               │
│  4. Stop Container                                │
│  5. Remove Container                              │
│  6. View Container Logs                           │
│  7. Execute Command in Container                  │
│  8. Back to Main Menu                             │
└───────────────────────────────────────────────────┘
```

**Features:**
- **List containers**: Shows all containers with status
- **Run new container**: Interactive container creation with image, name, and command
- **Start/Stop**: Control container lifecycle with ID input
- **Remove**: Safe removal with confirmation prompt
- **View logs**: Real-time and historical log viewing with follow option
- **Execute commands**: Run commands inside containers interactively

#### Image Management
```
┌──────────────── Image Management ─────────────────┐
│  1. List Images                                   │
│  2. Import Image                                  │
│  3. Remove Image                                  │
│  4. Tag Image                                     │
│  5. Inspect Image                                 │
│  6. Build Image                                   │
│  7. Back to Main Menu                             │
└───────────────────────────────────────────────────┘
```

**Features:**
- **List images**: Display all available images with metadata
- **Import**: Import images from tar files with file path input
- **Remove**: Safe image removal with confirmation
- **Tag**: Create new tags for existing images
- **Inspect**: View detailed image information
- **Build**: Build images from Buildfiles with custom tags

#### CRI Server Control
```
┌────────────────── CRI Server ─────────────────────┐
│  1. Start CRI Server                             │
│  2. Check CRI Server Status                      │
│  3. Test CRI Connection                           │
│  4. View CRI Endpoints                            │
│  5. Back to Main Menu                             │
└───────────────────────────────────────────────────┘
```

**Features:**
- **Start server**: Launch CRI HTTP server with custom port
- **Status check**: Verify server health and configuration
- **Connection test**: Test server connectivity and responsiveness
- **Endpoint documentation**: Display all available CRI API endpoints

#### Volume Management
```
┌─────────────── Volume Management ─────────────────┐
│  1. List Volumes                                  │
│  2. Create Volume                                 │
│  3. Remove Volume                                 │
│  4. Inspect Volume                                │
│  5. Remove All Volumes                            │
│  6. Back to Main Menu                             │
└───────────────────────────────────────────────────┘
```

**Features:**
- **List volumes**: Show all named volumes
- **Create**: Create new volumes with custom names
- **Remove**: Individual volume removal with confirmation
- **Inspect**: View volume details and usage
- **Bulk removal**: Remove all volumes with safety confirmation

#### Registry Operations
```
┌────────────── Registry Operations ────────────────┐
│  1. Start Local Registry                          │
│  2. Stop Local Registry                           │
│  3. Push Image to Registry                        │
│  4. Pull Image from Registry                      │
│  5. List Registries                               │
│  6. Back to Main Menu                             │
└───────────────────────────────────────────────────┘
```

**Features:**
- **Local registry**: Start/stop built-in registry server
- **Push/Pull**: Registry image operations with custom registry support
- **Registry management**: View configured registries and status

#### System Information
```
╔════════════════ System Information ═══════════════╗
║ Servin Runtime Information:                        ║
║ Platform: Windows/Linux/macOS                      ║
║ Time: 2024-01-01 12:00:00                         ║
║ Container Statistics:                              ║
║ Image Statistics:                                  ║
╚════════════════════════════════════════════════════╝
```

**Features:**
- **Runtime info**: Servin version and platform details
- **Statistics**: Container and image counts
- **System status**: Overall health and configuration

## GUI Features (Future)

### Planned GUI Components

#### Main Window Layout
- **Sidebar Navigation**: Quick access to all sections
- **Container List**: Sortable table with status indicators
- **Image Gallery**: Visual image browser with thumbnails
- **Log Panel**: Integrated log viewer with search and filtering
- **Status Bar**: Real-time system status and notifications

#### Container Management
- **Visual Status**: Color-coded container states
- **Quick Actions**: Start/stop/remove buttons
- **Log Streaming**: Real-time log updates
- **Resource Monitoring**: CPU/memory usage graphs
- **Port Mapping**: Visual port configuration

#### Image Management
- **Image Browser**: Thumbnail view of container images
- **Import Wizard**: Drag-and-drop image import
- **Tag Management**: Visual tag editing
- **Size Visualization**: Storage usage charts
- **Build Progress**: Real-time build status

#### CRI Server Dashboard
- **Server Status**: Visual health indicators
- **Endpoint Browser**: Interactive API explorer
- **Performance Metrics**: Request/response statistics
- **Configuration Panel**: Server settings management

## Installation and Setup

### Prerequisites
- **Go 1.21+** for building from source
- **Terminal**: Any terminal emulator for TUI
- **Display**: GUI requires display server (X11/Wayland on Linux)

### Building
```bash
# Build main servin executable
go build -o servin.exe .

# Build TUI component
go build -o servin-desktop.exe ./cmd/servin-desktop

# Build GUI component (when available)
go build -o servin-gui.exe ./cmd/servin-gui
```

### Running
```bash
# Launch TUI
servin gui --tui

# Launch GUI (fallback to TUI if unavailable)
servin gui

# Run TUI directly
./servin-desktop.exe
```

## Platform Support

| Feature | Windows | Linux | macOS |
|---------|---------|-------|-------|
| TUI | ✅ | ✅ | ✅ |
| GUI | 🔄 | 🔄 | 🔄 |
| Container Ops | ✅ | ✅ | ✅ |
| Image Ops | ✅ | ✅ | ✅ |
| CRI Server | ✅ | ✅ | ✅ |
| Volume Ops | ✅ | ✅ | ✅ |
| Registry Ops | ✅ | ✅ | ✅ |

**Legend:**
- ✅ Fully supported
- 🔄 In development
- ❌ Not supported

## Configuration

### Environment Variables
```bash
# Set default GUI port
export SERVIN_GUI_PORT=8081

# Set default host
export SERVIN_GUI_HOST=localhost

# Enable development mode
export SERVIN_DEV_MODE=true
```

### Command Line Options
```bash
# GUI options
servin gui --port 8081          # Custom port
servin gui --host 0.0.0.0       # Bind to all interfaces
servin gui --dev                # Development mode
servin gui --tui                # Force TUI mode
```

## Keyboard Shortcuts (TUI)

| Key | Action |
|-----|--------|
| `1-9` | Select menu option |
| `Enter` | Confirm selection |
| `Ctrl+C` | Exit/Cancel |
| `Esc` | Back to previous menu |

## Troubleshooting

### Common Issues

#### "GUI not available"
- **Cause**: No display server detected
- **Solution**: Use `--tui` flag or set up display

#### "TUI executable not found"
- **Cause**: servin-desktop.exe not built
- **Solution**: Run `go build -o servin-desktop.exe ./cmd/servin-desktop`

#### "Command not found"
- **Cause**: servin executable not in PATH
- **Solution**: Use full path or add to PATH

### Debug Mode
```bash
# Enable verbose logging
servin gui --tui --verbose

# Check system info
servin gui --tui
# Select option 6 (System Information)
```

## Development

### Adding New Features

#### TUI Extensions
1. Add new menu option in `showMainMenu()`
2. Implement handler in `handleMainMenu()`
3. Create submenu function for detailed operations
4. Add helper functions for specific operations

#### GUI Extensions
1. Design UI components in Fyne
2. Add event handlers for user interactions
3. Integrate with Servin CLI commands
4. Test across platforms

### Contributing
1. Fork the repository
2. Create feature branch for GUI/TUI improvements
3. Test on multiple platforms
4. Submit pull request with screenshots/demos

## Future Enhancements

### TUI Improvements
- **Color support**: Syntax highlighting and status colors
- **Search functionality**: Filter containers and images
- **Bulk operations**: Multi-select for batch actions
- **Configuration**: Save user preferences
- **Themes**: Light/dark mode support

### GUI Implementation
- **Fyne framework**: Cross-platform native GUI
- **Web interface**: Browser-based management
- **System tray**: Background service management
- **Notifications**: Desktop alerts for events
- **Plugins**: Extension system for custom features

### Integration
- **Docker compatibility**: Import Docker containers/images
- **Kubernetes**: Native K8s cluster management
- **Cloud providers**: Integration with cloud container services
- **Monitoring**: Prometheus metrics and dashboards

The Servin Desktop interface provides a comprehensive, user-friendly way to manage containerized applications, making container operations accessible to both beginners and advanced users.
