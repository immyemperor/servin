# Servin Desktop GUI

## Overview

Servin Desktop is a modern graphical user interface for the Servin container runtime, built using Go and the Fyne cross-platform GUI framework. It provides a Docker Desktop-like experience with comprehensive container, image, and volume management capabilities.

## Features

### 🐳 Container Management
- **Lifecycle Operations**: Start, stop, restart, and remove containers
- **Container Creation**: Run new containers with custom configurations
- **Log Viewing**: Real-time and historical container logs
- **Shell Access**: Execute commands inside running containers
- **Status Monitoring**: Real-time container status updates

### 📦 Image Management
- **Image Operations**: Import, build, tag, and remove images
- **Build Support**: Create images from Dockerfiles
- **Image Inspection**: Detailed image metadata and layer information
- **Registry Integration**: Pull and push images to registries

### 💾 Volume Management
- **Volume Operations**: Create, remove, and inspect volumes
- **Volume Pruning**: Clean up unused volumes
- **Bind Mount Support**: Host directory mounting
- **Volume Drivers**: Support for various volume drivers

### 🔧 CRI Server Integration
- **Kubernetes Compatibility**: Full CRI (Container Runtime Interface) support
- **Server Management**: Start/stop CRI server
- **Health Monitoring**: Connection testing and status checking
- **Configuration**: Customizable server settings

### 📋 Logging & Monitoring
- **Activity Logs**: Complete operation history
- **Status Updates**: Real-time status bar updates
- **Export Logs**: Save logs to files
- **Refresh Controls**: Automatic and manual data refreshing

## Architecture

### Core Components

1. **ServinDesktopGUI**: Main application structure
   - Window management and layout
   - Data containers for runtime state
   - UI component organization

2. **Tab-based Interface**: Organized feature areas
   - Containers tab for container management
   - Images tab for image operations
   - Volumes tab for storage management
   - CRI Server tab for Kubernetes integration
   - Logs tab for activity monitoring
   - About tab for application information

3. **Action Handlers**: Business logic implementation
   - Container operations (start, stop, remove, exec, logs)
   - Image operations (import, build, tag, remove, inspect)
   - Volume operations (create, remove, inspect, prune)
   - CRI server operations (start, stop, test, status)

### Technology Stack

- **Language**: Go 1.24+
- **GUI Framework**: Fyne v2.6.3
- **Container Runtime**: Servin CLI integration
- **Platform Support**: Windows, macOS, Linux

## Installation & Setup

### Prerequisites

1. **Go Development Environment**
   ```bash
   go version  # Should be 1.24 or later
   ```

2. **CGO Support** (Required for Fyne)
   ```bash
   # Windows: Install TDM-GCC or MinGW-w64
   # macOS: Install Xcode Command Line Tools
   # Linux: Install build-essential
   ```

3. **Servin CLI** (Must be in PATH)
   ```bash
   servin --version
   ```

### Building the GUI

1. **Install Dependencies**
   ```bash
   cd gocontainer
   go mod tidy
   ```

2. **Enable CGO** (Windows)
   ```bash
   set CGO_ENABLED=1
   go env -w CGO_ENABLED=1
   ```

3. **Build Application**
   ```bash
   go build -o servin-gui.exe cmd/servin-gui/main.go cmd/servin-gui/actions.go
   ```

### Running the GUI

```bash
# Direct execution
./servin-gui.exe

# Or through Servin CLI
servin gui
```

## User Interface Guide

### Main Window Layout

```
┌─────────────────────────────────────────────────────────────┐
│ [Refresh] [Settings] [Info]                     Toolbar     │
├─────────────────────────────────────────────────────────────┤
│ [🐳 Containers] [📦 Images] [💾 Volumes] [🔧 CRI] [📋 Logs] │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│                    Tab Content Area                         │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│ Status: Ready - Servin Desktop v1.0           Status Bar   │
└─────────────────────────────────────────────────────────────┘
```

### Container Tab Features

- **Container List**: Displays all containers with status indicators
- **Action Buttons**: 
  - Run: Create and start new containers
  - Start/Stop: Control container lifecycle
  - Remove: Delete containers
  - Logs: View container output
  - Exec: Execute commands in containers

### Image Tab Features

- **Image List**: Shows available images with size information
- **Action Buttons**:
  - Import: Load images from files
  - Build: Create images from Dockerfiles
  - Tag: Add additional tags to images
  - Remove: Delete images
  - Inspect: View detailed image information

### Volume Tab Features

- **Volume List**: Displays configured volumes
- **Action Buttons**:
  - Create: Make new volumes
  - Remove: Delete volumes
  - Inspect: View volume details
  - Prune: Clean up unused volumes

### CRI Server Tab Features

- **Status Display**: Shows server running state
- **Control Panel**: Start/stop server operations
- **Testing Tools**: Connection verification
- **Configuration**: Server settings management

## Integration with Servin CLI

The GUI application integrates seamlessly with the Servin CLI:

1. **Command Execution**: All operations use `servin` CLI commands
2. **Data Synchronization**: Real-time updates from CLI operations
3. **Configuration Sharing**: Uses same configuration as CLI
4. **Fallback Support**: Graceful degradation if GUI unavailable

### CLI Integration Examples

```go
// Container operations
exec.Command("servin", "run", "--name", "myapp", "alpine:latest")
exec.Command("servin", "start", containerID)
exec.Command("servin", "logs", containerID)

// Image operations  
exec.Command("servin", "image", "build", "-t", "myapp:latest", ".")
exec.Command("servin", "image", "ls", "--json")

// Volume operations
exec.Command("servin", "volume", "create", volumeName)
exec.Command("servin", "volume", "inspect", volumeName)

// CRI operations
exec.Command("servin", "cri", "start", "--port", "8080")
exec.Command("servin", "cri", "status")
```

## Development Notes

### Code Organization

```
cmd/servin-gui/
├── main.go          # Main application and UI layout
├── actions.go       # Action handlers and business logic
├── main_old.go      # Previous version backup
└── actions_old.go   # Previous version backup
```

### Key Design Patterns

1. **MVC Architecture**: Separation of UI, data, and business logic
2. **Event-Driven**: Reactive UI updates based on user actions
3. **Command Pattern**: CLI integration through command execution
4. **Observer Pattern**: Status updates and refresh mechanisms

### Future Enhancements

1. **Enhanced Error Handling**: Better error recovery and user feedback
2. **Configuration Management**: Advanced settings and preferences
3. **Theme Support**: Light/dark mode and customization
4. **Plugin System**: Extension support for additional features
5. **Multi-Server Support**: Manage multiple Servin instances
6. **Performance Monitoring**: Resource usage and metrics display

## Troubleshooting

### Common Issues

1. **CGO Build Errors**
   - Ensure CGO is enabled: `go env -w CGO_ENABLED=1`
   - Install appropriate C compiler for your platform
   - Verify Fyne dependencies are correctly installed

2. **Servin CLI Not Found**
   - Add Servin CLI to system PATH
   - Verify installation: `servin --version`
   - Check permissions for CLI execution

3. **GUI Performance Issues**
   - Adjust refresh interval in settings
   - Close unused tabs and windows
   - Check system resources and GPU drivers

### Platform-Specific Notes

#### Windows
- Requires TDM-GCC or MinGW-w64 for CGO support
- May need Visual Studio Build Tools
- OpenGL drivers must be up to date

#### macOS  
- Requires Xcode Command Line Tools
- May need to allow app in Security & Privacy settings
- Metal graphics acceleration supported

#### Linux
- Requires build-essential package
- OpenGL development libraries needed
- Wayland and X11 both supported

## Contributing

The Servin Desktop GUI is part of the larger Servin project. Contributions are welcome:

1. **Bug Reports**: Submit issues with detailed reproduction steps
2. **Feature Requests**: Propose new functionality with use cases
3. **Code Contributions**: Follow Go best practices and Fyne guidelines
4. **Documentation**: Improve user guides and developer documentation

## License

This project is part of Servin and follows the same licensing terms as the main project.
