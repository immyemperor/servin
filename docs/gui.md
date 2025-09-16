---
layout: default
title: Desktop GUI Application
permalink: /gui/
---

# 🖱️ Desktop GUI Application

The Servin Desktop GUI provides a clean, intuitive interface for container management with modern design principles. Built with Fyne v2, it offers native performance across Windows, Linux, and macOS with automatic refresh and real-time status updates.

## 🚀 Getting Started

### **Launching the GUI**
- **Windows**: Start Menu → Servin Container Runtime → Servin GUI
- **Linux**: Applications → Development → Servin GUI  
- **macOS**: Applications → Servin GUI
- **Command Line**: `servin gui` or `servin-gui`

### **Interface Overview**
The GUI features a tabbed interface with four main sections:
- **📦 Containers** - Container lifecycle management
- **🖼️ Images** - Image management and operations
- **🔗 CRI Server** - Kubernetes CRI server control
- **📋 Logs** - Real-time application logs

## 📦 Container Management

### **Container List View**
Displays all containers with:
- **Status Icons**: Visual indicators (▶️ running, ⏸️ stopped)
- **Container Name**: User-friendly container identifier
- **Current Status**: Running, stopped, paused, etc.
- **Base Image**: Container image name and tag

### **Container Operations**
- **▶️ Start** - Start a stopped container
- **⏸️ Stop** - Stop a running container  
- **🗑️ Remove** - Delete container (with confirmation)
- **📋 Logs** - View container logs in popup
- **🔄 Refresh** - Update container list

### **Container Actions**
```
Action Bar: [Start] [Stop] [Remove] [Logs] [Refresh]
```

## 🖼️ Image Management

### **Image List View**
Shows available images with:
- **📦 Storage Icon** - Image type indicator
- **Image Name** - Repository name
- **Tag** - Version or tag
- **Size** - Formatted file size (MB/GB)
- **Created Date** - When image was built

### **Image Operations**
- **📁 Import** - Import image from file
- **🗑️ Remove** - Delete unused images
- **🏷️ Tag** - Add tags to images
- **ℹ️ Inspect** - View image details
- **🔄 Refresh** - Update image list

## 🔗 CRI Server Management

### **Server Control Panel**
The CRI (Container Runtime Interface) tab provides:

#### **Status Display**
- **Server Status**: Running/Stopped indicator
- **Port Information**: CRI server listening port
- **Health Status**: Connection health check

#### **Control Actions**
- **▶️ Start CRI Server** - Launch Kubernetes CRI server
- **⏸️ Stop CRI Server** - Shutdown CRI server
- **🔄 Restart** - Restart CRI server
- **📊 Status Check** - Verify server health

#### **API Endpoints Reference**
Complete list of supported CRI endpoints:
- **Runtime Service**: Sandbox and container operations
- **Image Service**: Image management operations
- **Health Check**: Service status verification

## 📋 Application Logs

### **Real-time Log Viewer**
- **Live Updates**: Automatic log streaming
- **Timestamps**: Each log entry with time information
- **Action Tracking**: GUI operations and status updates
- **Clear Function**: Clear log history

### **Log Features**
- **Scrollable View**: Navigate through log history
- **Rich Text**: Formatted log output
- **Status Integration**: Log messages update status bar
- **Auto-refresh**: Continuous log updates

## ⚙️ Technical Features

### **Auto-refresh System**
- **2-second intervals**: Automatic data updates
- **Background threading**: Non-blocking UI operations
- **Smart updates**: Only refresh when data changes

### **Error Handling**
- **User-friendly dialogs**: Clear error messages
- **Confirmation prompts**: Safety for destructive operations
- **Status feedback**: Real-time operation status

### **Cross-platform Support**
- **Native Look**: Platform-appropriate styling
- **Keyboard Shortcuts**: Standard platform shortcuts
- **File Dialogs**: Native file selection

## 🎨 User Interface Design

### **Layout Structure**
```
┌─────────────────────────────────────────────────────────────────┐
│ ┌─ Containers ─┐ ┌─ Images ─┐ ┌─ CRI Server ─┐ ┌─ Logs ─┐     │
│ │ Container     │ │ Image    │ │ Server       │ │ Log     │     │
│ │ List          │ │ List     │ │ Controls     │ │ Viewer  │     │
│ │               │ │          │ │              │ │         │     │
│ │ [Actions...]  │ │ [Actions]│ │ [Start/Stop] │ │ [Clear] │     │
│ └───────────────┘ └──────────┘ └──────────────┘ └─────────┘     │
├─────────────────────────────────────────────────────────────────┤
│ Status: Ready | Last Update: 15:04:05                          │
└─────────────────────────────────────────────────────────────────┘
```

### **Window Properties**
- **Size**: 1200x800 pixels (resizable)
- **Position**: Centered on screen
- **Theme**: System-appropriate (light/dark)
- **Icons**: Fyne theme icons for consistency

- **Configuration**
  - Environment variables
  - Command and arguments
  - Working directory
  - User and group settings
  - Security options

#### **Logs Tab**
- **Real-time Log Streaming** - Live log updates
- **Search and Filter** - Find specific log entries
- **Log Level Filtering** - Error, warning, info, debug
- **Timestamp Options** - Show/hide timestamps
- **Export Logs** - Save logs to file
- **Auto-scroll** - Follow new log entries

#### **Exec Tab**
- **Interactive Shell** - Built-in terminal emulator
- **Command Execution** - Run commands in container
- **File Browser** - Navigate container filesystem
- **Process Manager** - View running processes
- **Environment Editor** - View/edit environment variables

#### **Files Tab**
- **Filesystem Browser** - Navigate container files
- **File Operations** - Copy, edit, delete files
- **Upload/Download** - Transfer files to/from container
- **Permission Management** - View/change file permissions
- **Text Editor** - Edit files directly in GUI

#### **Environment Tab**
- **Environment Variables** - View all environment variables
- **Variable Editor** - Add/edit/remove variables
- **Import/Export** - Load variables from file
- **Templates** - Common environment configurations
- **Validation** - Check variable syntax

#### **Volumes Tab**
- **Mount Points** - View all mounted volumes
- **Volume Browser** - Explore mounted volume contents
- **Mount Editor** - Add/remove volume mounts
- **Volume Creation** - Create new volumes on-the-fly
- **Backup/Restore** - Volume backup operations

### **Container Actions Menu**
Right-click any container or use the action menu:

- **Lifecycle Operations**
  - ▶️ Start - Start stopped container
  - ⏹️ Stop - Gracefully stop container
  - 🔄 Restart - Restart container
  - ⏸️ Pause - Pause container execution
  - ▶️ Resume - Resume paused container
  - 💀 Kill - Force stop container

- **Management Operations**
  - 🏷️ Rename - Change container name
  - 📋 Duplicate - Create copy of container
  - 📤 Export - Export container to tar
  - 📥 Import - Import container from tar
  - 🔧 Edit Config - Modify container settings

- **Monitoring**
  - 📊 Statistics - Detailed resource usage
  - 📄 View Logs - Open log viewer
  - 🔍 Inspect - Raw container configuration
  - 📈 Performance - Performance metrics

## 🖼️ Image Management

### **Image Gallery View**
Images are displayed in a visual gallery format:

- **Image Thumbnails** - Visual representation of images
- **Repository and Tag** - Clear labeling
- **Size Information** - Image size and virtual size
- **Creation Date** - When image was created/pulled
- **Usage Count** - How many containers use this image
- **Security Status** - Vulnerability scan results

### **Image Details Panel**
Select any image to view detailed information:

#### **General Information**
- **Image ID** - Full SHA256 hash
- **Repository and Tags** - All associated tags
- **Architecture** - Target architecture (amd64, arm64)
- **Operating System** - Base OS (linux, windows)
- **Size Breakdown** - Compressed and uncompressed sizes
- **Creation Details** - Author, creation date, build info

#### **Configuration**
- **Entry Point** - Default command entry point
- **Command** - Default command arguments
- **Environment** - Default environment variables
- **Exposed Ports** - Ports exposed by image
- **Working Directory** - Default working directory
- **User** - Default user for containers

#### **Layer History**
- **Layer Visualization** - Visual layer stack
- **Layer Sizes** - Individual layer sizes
- **Commands** - Dockerfile commands for each layer
- **Created By** - Build step information
- **Layer Sharing** - Shared layers with other images

#### **Security Scan**
- **Vulnerability Report** - Known security issues
- **Risk Assessment** - Overall security score
- **Package List** - Installed packages and versions
- **Recommendations** - Security improvement suggestions

### **Image Actions**
- **Registry Operations**
  - 📥 Pull - Download newer version
  - 📤 Push - Upload to registry
  - 🔍 Search - Find related images
  - 🏷️ Tag - Add new tags

- **Container Operations**
  - ▶️ Run - Create and start container
  - 🔧 Create - Create container without starting
  - 📝 Run with Options - Advanced run dialog
  - 🚀 Quick Run - Run with common settings

- **Management**
  - 🗑️ Delete - Remove image
  - 💾 Export - Save image to tar file
  - 📊 Analyze - Detailed image analysis
  - 🔍 Inspect - Raw image configuration

## 💾 Volume Management

### **Volume List View**
Volumes are displayed with the following information:
- **Volume Name** - User-defined or generated name
- **Driver Type** - local, nfs, cloud storage drivers
- **Mount Point** - Host filesystem location
- **Size Usage** - Used space and total capacity
- **Container Count** - How many containers use it
- **Created Date** - When volume was created

### **Volume Details Panel**
- **Configuration**
  - Driver options and settings
  - Mount point and permissions
  - Labels and metadata
  - Backup configuration

- **Usage Information**
  - Connected containers
  - File system type
  - Access patterns
  - Performance metrics

- **File Browser**
  - Browse volume contents
  - File operations (copy, move, delete)
  - Permission management
  - Search functionality

### **Volume Actions**
- **Lifecycle**
  - 🆕 Create - Create new volume
  - 🗑️ Delete - Remove volume
  - 📋 Clone - Duplicate volume
  - 🧹 Cleanup - Remove unused volumes

- **Data Operations**
  - 💾 Backup - Create volume backup
  - 📥 Restore - Restore from backup
  - 📤 Export - Export volume data
  - 📥 Import - Import volume data

## 🌐 Network Management

### **Network Topology View**
Visual representation of container networks:
- **Network Diagrams** - Visual network topology
- **Container Connections** - Show container relationships
- **IP Address Mapping** - Network address assignments
- **Port Mappings** - Port forwarding visualization
- **Traffic Flow** - Network traffic indicators

### **Network Details**
- **Configuration**
  - Network driver and options
  - Subnet and gateway settings
  - DNS configuration
  - Security policies

- **Connected Containers**
  - Container list with IP addresses
  - Port mappings and aliases
  - Connection status
  - Network statistics

### **Network Actions**
- **Management**
  - 🆕 Create - Create custom network
  - 🗑️ Delete - Remove network
  - 🔧 Configure - Edit network settings
  - 📊 Monitor - Network traffic analysis

- **Container Operations**
  - 🔗 Connect - Connect container to network
  - ❌ Disconnect - Remove container from network
  - 🏷️ Alias - Manage container aliases
  - 🔍 Inspect - Network configuration details

## 📊 System Monitoring

### **Dashboard Overview**
- **Resource Meters**
  - CPU usage gauge with history
  - Memory usage with available/total
  - Disk I/O rates and usage
  - Network traffic rates

- **Quick Statistics**
  - Total containers (running/stopped)
  - Image count and total size
  - Volume count and usage
  - Network count and connections

- **Recent Activity**
  - Container lifecycle events
  - Image pull/push operations
  - System alerts and warnings
  - Performance notifications

### **Performance Monitoring**
- **Real-time Graphs**
  - CPU usage over time
  - Memory consumption trends
  - Network I/O patterns
  - Disk usage and I/O

- **Container Performance**
  - Per-container resource usage
  - Performance comparison charts
  - Resource limit visualization
  - Bottleneck identification

### **Event Logging**
- **System Events**
  - Container lifecycle events
  - Image operations
  - Volume operations
  - Network changes

- **Log Analysis**
  - Event filtering and search
  - Export to file formats
  - Real-time event streaming
  - Alert configuration

## ⚙️ Settings and Configuration

### **Application Settings**
- **Appearance**
  - Theme selection (light/dark/auto)
  - Font size and family
  - Color scheme customization
  - Layout preferences

- **Behavior**
  - Auto-refresh intervals
  - Confirmation dialogs
  - Default actions
  - Keyboard shortcuts

- **Notifications**
  - Desktop notifications
  - Sound alerts
  - Email notifications
  - Webhook integrations

### **Runtime Configuration**
- **Connection Settings**
  - Socket path or TCP endpoint
  - Authentication credentials
  - Timeout settings
  - TLS configuration

- **Default Options**
  - Default container settings
  - Image pull preferences
  - Volume creation options
  - Network configurations

### **Advanced Settings**
- **Debug Options**
  - Logging level
  - Debug mode
  - Performance profiling
  - Error reporting

- **Security Settings**
  - Access controls
  - Audit logging
  - Encryption options
  - Certificate management

## 🔧 Advanced Features

### **Batch Operations**
- **Multi-Selection** - Select multiple containers/images
- **Bulk Actions** - Apply operations to selected items
- **Templates** - Save common configurations
- **Automation** - Schedule recurring tasks

### **Import/Export**
- **Container Export** - Save containers as images
- **Configuration Backup** - Export settings
- **Bulk Import** - Import multiple containers
- **Migration Tools** - Move between hosts

### **Integration Features**
- **Docker Compose** - Import/export compose files
- **Kubernetes** - Generate Kubernetes manifests
- **CI/CD Integration** - Build pipeline integration
- **Registry Integration** - Multi-registry support

### **Development Tools**
- **Container Debugging** - Debug running containers
- **Log Analysis** - Advanced log parsing
- **Performance Profiling** - Container performance analysis
- **Resource Planning** - Capacity planning tools

## 🎨 Customization

### **Themes and Appearance**
- **Built-in Themes**
  - Light theme for bright environments
  - Dark theme for low-light conditions
  - High contrast for accessibility
  - Custom themes with full color control

### **Layout Customization**
- **Panel Arrangement** - Customize panel layout
- **Column Configuration** - Show/hide specific columns
- **Toolbar Customization** - Add/remove toolbar buttons
- **Keyboard Shortcuts** - Customize hotkeys

### **Workflow Optimization**
- **Quick Actions** - Configure favorite actions
- **Context Menus** - Customize right-click menus
- **Default Behaviors** - Set preferred defaults
- **Automation Rules** - Create custom automation

## 🔌 Plugin System

### **Available Plugins**
- **Registry Plugins** - Additional registry support
- **Monitoring Extensions** - Enhanced monitoring
- **Export Tools** - Additional export formats
- **Theme Packages** - Community themes

### **Plugin Development**
- **Plugin API** - Develop custom plugins
- **Extension Points** - Available customization points
- **Development Kit** - Tools for plugin creation
- **Community Gallery** - Share and discover plugins

---

## 📚 Next Steps

- **[CLI Reference]({{ '/cli' | relative_url }})** - Learn command-line operations
- **[TUI Guide]({{ '/tui' | relative_url }})** - Explore terminal interface
- **[Configuration]({{ '/configuration' | relative_url }})** - Customize your setup
- **[Development]({{ '/development' | relative_url }})** - Contribute to Servin

<div class="gui-tips">
  <h3>💡 GUI Pro Tips</h3>
  <ul>
    <li><strong>Keyboard Navigation:</strong> Most operations support keyboard shortcuts for power users</li>
    <li><strong>Multi-Monitor:</strong> GUI supports multiple monitors and custom window arrangements</li>
    <li><strong>Performance:</strong> Use the resource monitoring to optimize container performance</li>
    <li><strong>Automation:</strong> Set up automated tasks for routine container management</li>
    <li><strong>Backup:</strong> Regularly export your configurations for disaster recovery</li>
  </ul>
</div>
