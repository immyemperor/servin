---
layout: default
title: Desktop GUI Application
permalink: /gui/
---

# 🖱️ Desktop GUI Application

The Servin Desktop GUI provides a modern, web-based interface for container management with real-time updates and responsive design. Built with Flask backend and pywebview frontend, it offers native desktop integration across Windows, Linux, and macOS while maintaining the flexibility of web technologies.

## 🚀 Getting Started

### **Binary Distribution**
The GUI is distributed as a compiled binary (`servin-gui` / `servin-gui.exe`) for better performance and easier deployment:

- **Windows**: `servin-gui.exe` - Single executable, no Python required
- **Linux**: `servin-gui` - Native binary with embedded Python runtime  
- **macOS**: `servin-gui` - Universal binary compatible with Intel and Apple Silicon

### **Launching the GUI**
- **Windows**: Start Menu → Servin Container Runtime → Servin GUI
- **Linux**: Applications → Development → Servin GUI  
- **macOS**: Applications → Servin GUI
- **Command Line**: `servin gui` or `servin-gui` directly

### **Interface Overview**
The GUI features a single-page web application with real-time sections:
- **📦 Containers** - Container lifecycle management with live status updates
- **🖼️ Images** - Image management and operations
- **� Volumes** - Persistent volume management
- **� System Info** - Runtime information and statistics

## 📦 Container Management

### **Container Dashboard**
Displays all containers in a responsive table layout with intelligent action buttons:
- **Status Indicators**: Real-time visual badges (🟢 running, 🔴 stopped, 🟡 paused)
- **Container Information**: Name, image, status, ports, and creation time
- **Smart Actions**: Context-aware buttons that adapt to container state
- **Real-time Updates**: Live status monitoring with WebSocket integration
- **Detailed View**: Click any container for comprehensive details

### **Container Details View**
Enhanced container inspection with tabbed interface:
- **📊 Overview** - Complete container metadata and configuration
- **📝 Logs** - Real-time log streaming with auto-scroll and download
- **📁 Files** - Container filesystem browser and file operations
- **💻 Terminal** - Interactive shell access with auto-connect
- **🔧 Environment** - Environment variables display and management
- **� Volumes** - Mount points and volume information
- **🌐 Network** - Networking configuration and port mappings
- **� Statistics** - Resource usage monitoring and metrics

### **Intelligent Container Actions**
Action buttons adapt dynamically based on container state:

**For Running Containers:**
- **⏸️ Stop** - Gracefully stop the container
- **🔄 Restart** - Restart the container with current configuration

**For Stopped Containers:**
- **▶️ Start** - Start the container
- **🗑️ Remove** - Delete the container with confirmation

### **Interactive Terminal**
Full-featured terminal interface with:
- **Auto-Connect**: Automatically connects when accessing terminal tab
- **Real Shell Prompt**: Displays `user@hostname:path$ ` format
- **Command History**: Navigate previous commands with arrow keys
- **Live Output**: Real-time command execution and output display
- **Multiple Shells**: Supports bash, sh, and zsh shells

### **Real-time Log Streaming**
Advanced logging capabilities:
- **Live Streaming**: WebSocket-based real-time log updates
- **Persistent Logs**: Logs maintain content when switching tabs
- **Auto-scroll**: Automatic scrolling to latest log entries
- **Download Logs**: Export logs to local file system
- **Clear Display**: Clean, formatted log presentation

## 🖼️ Image Management

### **Image Gallery**
Displays available images in a grid layout:
- **📦 Image Cards**: Repository name, tag, and size information
- **Creation Info**: Build date and image ID
- **Action Buttons**: Remove and inspect operations
- **Import Function**: Drag-and-drop or file browser support

### **Image Operations**
- **📁 Import** - Import image from tarball file
- **🗑️ Remove** - Delete unused images with confirmation
- **ℹ️ Inspect** - View detailed image metadata
- **🔄 Auto-refresh** - Live updates when images change

### **Image Cards**
```
┌─────────────────────────────────────┐
│ 🖼️ nginx:latest                     │
│ Size: 142.8 MB                       │
│ Created: 3 days ago                   │
│ [🗑️] [ℹ️]                          │
└─────────────────────────────────────┘
```

## 💾 Volume Management

### **Volume Dashboard**
Shows persistent volumes with:
- **Volume Name**: User-defined volume identifier
- **Mount Path**: Container mount point information
- **Usage Status**: Whether volume is currently in use
- **Creation Time**: When volume was created

### **Volume Operations**
- **➕ Create** - Create new named volume
- **🗑️ Remove** - Delete unused volumes
- **� Inspect** - View volume details and mount information
- **🔄 Refresh** - Update volume list

## 📊 System Information

### **Runtime Status**
Displays system information:
- **Servin Version**: Runtime version and build info
- **Platform**: Operating system and architecture
- **Container Count**: Total containers (running/stopped)
- **Image Count**: Total images stored locally
- **Volume Count**: Total persistent volumes

### **Resource Usage**
- **Storage Info**: Space used by containers and images
- **Runtime Path**: Servin installation and data directories
- **Configuration**: Current runtime settings and paths

## ⚙️ Technical Architecture

### **Flask Backend**
- **RESTful API**: Clean HTTP endpoints for all operations
- **Real-time Updates**: Automatic polling for live status updates
- **Error Handling**: Comprehensive error messages and validation
- **CORS Support**: Cross-origin requests for development

### **Frontend Architecture**
- **Modular Components**: 7 specialized JavaScript components for maintainable code
  - `ServinGUI.js`: Main application controller and coordinator
  - `ContainerDetails.js`: Container inspection and management
  - `Logs.js`: Real-time log streaming and display
  - `Terminal.js`: Interactive container terminal sessions
  - `FileExplorer.js`: Container filesystem navigation
  - `APIClient.js`: HTTP API communication layer
  - `SocketManager.js`: WebSocket connection management
- **CSS Framework**: 8 dedicated CSS modules for consistent styling
  - Component-specific stylesheets with CSS custom properties
  - Responsive design patterns and mobile-first approach
  - Dark theme with consistent color palette and typography

### **pywebview Integration**
- **Native Desktop**: Desktop window with web technologies
- **Cross-platform**: Consistent experience across operating systems
- **Web Standards**: Modern HTML5, CSS3, and ES6+ JavaScript
- **Responsive Design**: Adapts to different window sizes and screen densities

### **Binary Distribution**
- **PyInstaller Compilation**: Single-file executable with embedded Python
- **No Dependencies**: Runs without Python installation on target system
- **Platform Optimized**: Native binaries for Windows, Linux, and macOS
- **Reduced Size**: Optimized 13MB executable with all dependencies

## 🛠️ Development & Deployment

### **Development Mode**
For developers working on the GUI:

```bash
# Clone and setup development environment
cd webview_gui
python -m venv venv
source venv/bin/activate  # Linux/macOS
pip install -r requirements.txt

# Run in development mode
python main.py

# Run web-only demo
python demo.py
```

### **Binary Building**
Automated build process creates platform-specific executables:

```bash
# Build all platforms (uses PyInstaller)
./build-all.sh

# Creates:
# - dist/windows/servin-gui.exe
# - dist/linux/servin-gui  
# - dist/mac/servin-gui
```

### **API Endpoints**
The GUI communicates with Servin through these endpoints:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/containers` | GET | List all containers |
| `/api/containers/{id}/start` | POST | Start container |
| `/api/containers/{id}/stop` | POST | Stop container |  
| `/api/containers/{id}/remove` | DELETE | Remove container |
| `/api/images` | GET | List all images |
| `/api/images/{id}/remove` | DELETE | Remove image |
| `/api/volumes` | GET | List volumes |
| `/api/volumes` | POST | Create volume |
| `/api/system/info` | GET | System information |

## 🎨 User Experience

### **Design Principles**
- **Dark Theme**: Modern dark interface with blue accents
- **Responsive Layout**: Adapts to different window sizes and resolutions
- **Intuitive Icons**: Clear visual indicators for all operations
- **Immediate Feedback**: Instant visual confirmation of actions
- **Error Prevention**: Confirmation dialogs for destructive operations

### **Accessibility**
- **Keyboard Navigation**: Tab through all interactive elements
- **Screen Reader Support**: Semantic HTML with proper labels
- **High Contrast**: Clear color contrast for visibility
- **Responsive Text**: Scales with system font size preferences

### **Performance Features**
- **Auto-refresh**: Live updates every 5 seconds without user intervention
- **Lazy Loading**: Efficient data fetching and caching
- **Local State**: Maintains UI state between operations
- **Fast Startup**: Binary launches in under 2 seconds

## 🆕 Recent Enhancements

### **UI/UX Improvements**
- **Responsive Headers**: Container details headers now properly handle text overflow with ellipsis
- **Enhanced Buttons**: Container action buttons adapt intelligently to container state
- **Improved Styling**: Consistent spacing, typography, and color scheme across all components
- **Mobile Responsive**: Better layout adaptation for different screen sizes

### **Terminal Enhancements**
- **Auto-Connect**: Terminal automatically connects without manual shell selection
- **Realistic Prompt**: Enhanced shell prompt display with proper user@container format
- **Better Session Management**: Improved connection handling and error recovery
- **Command History**: Previous commands are preserved during session

### **Log Management**
- **Persistent Logs**: Log content persists when switching between tabs
- **Streaming Optimization**: Real-time log streaming with efficient data handling
- **Better Error Handling**: Improved handling of log retrieval failures
- **Search Integration**: Log content remains searchable and scrollable

### **Container Details**
- **Enhanced Tabs**: All tabs (logs, files, exec, env, volumes, network, stats) now work consistently
- **Environment Variables**: Fixed display of complex environment objects
- **Action Intelligence**: Container buttons show relevant actions based on current state
- **Data Persistence**: Tab content maintains state when switching views

### **Architecture Updates**
- **Modular Components**: 7 specialized JavaScript components for maintainability
- **CSS Framework**: 8 dedicated CSS modules for consistent styling
- **API Optimization**: Improved data fetching and response handling
- **Component Communication**: Enhanced inter-component messaging and state management

## 🚀 Installation Methods

### **From GitHub Releases (Recommended)**
Download pre-built binaries from [GitHub Releases](https://github.com/immyemperor/servin/releases/latest):

1. **Windows**: Download `servin-windows-amd64.zip`, extract, and run installer
2. **Linux**: Download `servin-linux-amd64.tar.gz`, extract, and run installer
3. **macOS**: Download `Servin-Container-Runtime.dmg` or `servin-macos-universal.tar.gz`

### **Command Line Access**
After installation, launch from anywhere:

```bash
# Launch GUI directly
servin-gui

# Launch through main CLI
servin gui

# Launch with custom port
servin gui --port 8080

# Launch in development mode
servin gui --dev
```

## 🔧 Configuration & Troubleshooting

### **Configuration**
The GUI automatically detects Servin binary location:
1. Platform-specific build directory
2. System PATH
3. Same directory as GUI executable

### **Common Issues**

#### **Binary Not Found**
- Ensure `servin` binary is in PATH or same directory
- Check file permissions (executable bit on Linux/macOS)
- Verify platform compatibility (ARM64 vs AMD64)

#### **Port Conflicts**
- Default port 5555 may be in use
- Use `--port` flag to specify different port
- Check firewall settings for local connections

#### **Web Engine Issues**
- GUI falls back to default browser if pywebview fails
- Install system web engine dependencies
- Use `demo.py` for browser-only testing

### **Logging & Debugging**
- GUI logs displayed in system console
- Enable debug mode with environment variable
- Check GUI process with system task manager
- View network requests in browser developer tools

## 📚 Additional Resources

- **Source Code**: `/webview_gui/` directory in Servin repository
- **API Documentation**: Built-in `/api/` endpoints documentation
- **Development Guide**: `/webview_gui/README.md`
- **Build Scripts**: Cross-platform build automation in `/build-all.sh`
- **Issue Reporting**: GitHub Issues for bug reports and feature requests

---

## 📚 Next Steps

- **[CLI Reference]({{ '/cli' | relative_url }})** - Learn command-line operations
- **[TUI Guide]({{ '/tui' | relative_url }})** - Explore terminal interface  
- **[Installation Guide]({{ '/installation' | relative_url }})** - Download and install Servin
- **[Configuration]({{ '/configuration' | relative_url }})** - Customize your setup

<div class="gui-tips">
  <h3>� GUI Pro Tips</h3>
  <ul>
    <li><strong>Real-time Updates:</strong> The interface automatically refreshes every 5 seconds for live status</li>
    <li><strong>Confirmation Dialogs:</strong> All destructive operations require confirmation to prevent accidents</li>
    <li><strong>Browser Fallback:</strong> If desktop app fails, GUI falls back to default browser automatically</li>
    <li><strong>Development Mode:</strong> Use <code>--dev</code> flag for development and testing</li>
    <li><strong>Custom Ports:</strong> Use <code>--port</code> to avoid conflicts with other services</li>
  </ul>
</div>
