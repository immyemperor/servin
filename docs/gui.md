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
Displays all containers in a responsive card layout:
- **Status Indicators**: Visual badges (🟢 running, 🔴 stopped, 🟡 paused)
- **Container Information**: Name, image, status, and creation time
- **Quick Actions**: Start, stop, remove buttons with confirmation dialogs
- **Real-time Updates**: Automatic refresh every 5 seconds

### **Container Operations**
Interactive buttons for each container:
- **▶️ Start** - Start a stopped container with instant feedback
- **⏸️ Stop** - Gracefully stop a running container  
- **🗑️ Remove** - Delete container with confirmation dialog
- **📋 Logs** - View real-time container logs (if supported)
- **🔄 Auto-refresh** - Live status updates without manual refresh

### **Container Cards**
```
┌─────────────────────────────────────┐
│ 📦 container-name          🟢 Running │
│ Image: nginx:latest                   │
│ Created: 2 hours ago                  │
│ [▶️] [⏸️] [🗑️] [📋]                 │
└─────────────────────────────────────┘
```

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

### **pywebview Frontend**
- **Native Integration**: Desktop window with web technologies
- **Cross-platform**: Consistent experience across operating systems
- **Web Standards**: Modern HTML5, CSS3, and JavaScript
- **Responsive Design**: Adapts to different window sizes

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
