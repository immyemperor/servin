---
layout: default
title: Desktop GUI Application
permalink: /gui/
---

# 🖱️ Desktop GUI Application

The Servin Desktop GUI provides a professional, visual interface for container management with modern design principles and intuitive workflows. Built with Fyne v2.6, it offers native performance across Windows, Linux, and macOS.

## 🚀 Getting Started

### **Launching the GUI**
- **Windows**: Start Menu → Servin Container Runtime
- **Linux**: Applications → Development → Servin Runtime  
- **macOS**: Applications → Servin Runtime
- **Command Line**: `servin gui` or `servin-gui`

### **First Launch Setup**
1. **Connection Configuration** - Configure runtime connection
2. **Theme Selection** - Choose light/dark theme
3. **Layout Preferences** - Customize interface layout
4. **Notification Settings** - Configure alerts and notifications

## 🎨 Interface Overview

### **Main Window Layout**
```
┌─────────────────────────────────────────────────────────────────┐
│ File  Edit  View  Container  Image  Volume  Network  Help     │
├─────────────────────────────────────────────────────────────────┤
│ 🏠 📦 🖼️ 💾 🌐 ⚙️    Search: [nginx________] 🔍 [Filter▼]      │
├─────────────────────────────────────────────────────────────────┤
│ ┌─ Sidebar ─┐ ┌────────── Main Content Area ──────────────────┐ │
│ │           │ │                                              │ │
│ │ 📦 Containers │ │  Container List / Details View             │ │
│ │   Running     │ │                                              │ │
│ │   Stopped     │ │  ┌─────────────────────────────────────┐    │ │
│ │   Paused      │ │  │ web-server     nginx:latest    🟢   │    │ │
│ │               │ │  │ Created: 2h ago  CPU: 5%  Mem: 128MB│    │ │
│ │ 🖼️ Images      │ │  └─────────────────────────────────────┘    │ │
│ │   Local       │ │                                              │ │
│ │   Registry    │ │  ┌─────────────────────────────────────┐    │ │
│ │               │ │  │ api-service    node:16         🟡   │    │ │
│ │ 💾 Volumes     │ │  │ Created: 1d ago  CPU: 2%  Mem: 256MB│    │ │
│ │   Mounted     │ │  └─────────────────────────────────────┘    │ │
│ │   Unused      │ │                                              │ │
│ │               │ │  Actions: [Start] [Stop] [Restart] [⋮]      │ │
│ │ 🌐 Networks    │ │                                              │ │
│ │   Bridge      │ │                                              │ │
│ │   Custom      │ │                                              │ │
│ │               │ │                                              │ │
│ │ 📊 System      │ │                                              │ │
│ │   Overview    │ │                                              │ │
│ │   Logs        │ │                                              │ │
│ │   Events      │ │                                              │ │
│ └───────────────┘ └──────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────────┤
│ Status: Connected | Containers: 12 | Images: 25 | CPU: 15%   │
└─────────────────────────────────────────────────────────────────┘
```

## 📦 Container Management

### **Container List View**
The container view displays all containers with visual status indicators:

- **🟢 Running** - Container is actively running
- **🟡 Paused** - Container is paused
- **🔴 Stopped** - Container is stopped
- **🟠 Restarting** - Container is restarting
- **⚫ Unknown** - Status unknown

### **Container Cards**
Each container is displayed as a card showing:
- **Name and Image** - Container identification
- **Status Badge** - Visual status indicator
- **Resource Usage** - CPU, memory, network stats
- **Uptime/Age** - How long container has been running
- **Port Mappings** - Exposed ports and bindings
- **Quick Actions** - Start, stop, restart, logs

### **Container Details Panel**
Click any container to open the detailed view:

#### **Overview Tab**
- **Basic Information**
  - Container ID and name
  - Image name and tag
  - Creation time and uptime
  - Current status and exit code
  - Restart policy and count

- **Resource Usage**
  - Real-time CPU usage graph
  - Memory usage with limits
  - Network I/O statistics
  - Disk I/O metrics
  - Process count

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
