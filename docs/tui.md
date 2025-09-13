---
layout: default
title: Terminal User Interface
permalink: /tui/
---

# 📟 Terminal User Interface (TUI)

The Servin Terminal User Interface provides an interactive, menu-driven experience for container management directly in your terminal. Perfect for server environments, SSH sessions, and users who prefer keyboard-driven workflows.

## 🚀 Getting Started

### **Launching the TUI**
```bash
# Start the terminal interface
servin tui

# Start with specific view
servin tui --view containers
servin tui --view images
servin tui --view volumes

# Start with connection options
servin tui --host unix:///var/run/servin.sock
servin tui --host tcp://remote-host:2375
```

### **Navigation Basics**
- **Arrow Keys** - Navigate menus and lists
- **Tab/Shift+Tab** - Switch between panels
- **Enter** - Select item or confirm action
- **Escape** - Go back or cancel operation
- **q** - Quit application
- **/** - Search/filter
- **?** - Show help overlay

## 🖥️ Interface Overview

### **Main Dashboard**
```
┌─────────────────────────────────────────────────────────────────┐
│                     Servin Container Runtime                    │
├─────────────────────────────────────────────────────────────────┤
│ [Containers] [Images] [Volumes] [Networks] [System] [Help]     │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  System Status:                    Resource Usage:             │
│  ✅ Runtime: Active                 📊 CPU: 15%                  │
│  🔌 CRI Server: Running            💾 Memory: 2.1GB/8GB        │
│  📡 API Port: 10010                💽 Disk: 45GB/100GB         │
│                                                                 │
│  Quick Stats:                      Recent Activity:            │
│  📦 Containers: 12 (8 running)     🔄 nginx-web started        │
│  🖼️  Images: 25 (15 in use)        ⏹️  old-app stopped         │
│  💾 Volumes: 8 (6 mounted)         📥 ubuntu:latest pulled     │
│  🌐 Networks: 4 (3 active)         🗑️  temp-container removed   │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│ Press Tab to navigate | ? for help | q to quit                 │
└─────────────────────────────────────────────────────────────────┘
```

## 📦 Container Management

### **Container List View**
```
┌─────────────────── Containers ──────────────────┐
│ Search: [nginx___________] 🔍 Filter: [All____] │
├──────────────────────────────────────────────────┤
│ ID       │ Name        │ Image      │ Status    │
├──────────┼─────────────┼────────────┼───────────┤
│ 🟢 a1b2c3 │ web-server  │ nginx:1.21 │ Running   │
│ 🟡 d4e5f6 │ api-service │ node:16    │ Paused    │
│ 🔴 g7h8i9 │ old-worker  │ worker:2.0 │ Exited(0) │
│ 🟢 j1k2l3 │ database    │ mysql:8.0  │ Running   │
├──────────────────────────────────────────────────┤
│ ↑/↓: Navigate | Enter: Details | Space: Select  │
│ s: Start | t: Stop | r: Restart | d: Delete     │
│ l: Logs | e: Exec | i: Inspect | c: Create      │
└──────────────────────────────────────────────────┘
```

### **Container Details Panel**
```
┌─────────────── Container Details: web-server ───────────────┐
│                                                              │
│ General Information:                                         │
│ 🆔 ID: a1b2c3d4e5f6                                          │
│ 🏷️  Name: web-server                                        │
│ 🖼️  Image: nginx:1.21                                       │
│ ⚡ Status: Running (2h 15m)                                 │
│ 🔄 Restart Policy: unless-stopped                          │
│                                                              │
│ Network Configuration:                                       │
│ 🌐 IP Address: 172.17.0.2                                   │
│ 🔗 Ports: 80:8080/tcp, 443:8443/tcp                        │
│ 🌍 Networks: bridge, web-frontend                          │
│                                                              │
│ Resource Usage:                                              │
│ 📊 CPU: 5.2% (Limit: 1.0 CPU)                              │
│ 💾 Memory: 128MB/512MB (25%)                               │
│ 💽 Disk I/O: 1.2MB read, 850KB write                      │
│ 🌐 Network I/O: 15MB in, 25MB out                          │
│                                                              │
│ Environment Variables:                                       │
│ 🔧 NODE_ENV=production                                      │
│ 🔧 PORT=3000                                                │
│ 🔧 DATABASE_URL=mysql://db:3306/app                        │
│                                                              │
│ Volumes:                                                     │
│ 💾 /var/www/html → /app/public (ro)                        │
│ 💾 app-logs → /var/log/nginx                               │
│                                                              │
├──────────────────────────────────────────────────────────────┤
│ [Logs] [Exec] [Files] [Stats] [Edit] [Actions]             │
└──────────────────────────────────────────────────────────────┘
```

### **Container Actions Menu**
- **🔄 Lifecycle**
  - ▶️ Start - Start stopped container
  - ⏹️ Stop - Gracefully stop running container  
  - 🔄 Restart - Restart container
  - ⏸️ Pause - Pause container execution
  - ▶️ Unpause - Resume paused container
  - 💀 Kill - Force stop container

- **📊 Monitoring**
  - 📄 View Logs - Real-time log streaming
  - 📈 Live Stats - CPU, memory, network usage
  - 🔍 Inspect - Detailed configuration view
  - 📋 Processes - Running processes inside container

- **🔧 Management**
  - 💻 Execute Shell - Interactive shell access
  - 📁 Browse Files - Container filesystem explorer
  - 🏷️ Rename - Change container name
  - 📝 Edit Config - Modify container settings
  - 🗑️ Remove - Delete container

## 🖼️ Image Management

### **Image List View**
```
┌──────────────────── Images ─────────────────────┐
│ Search: [ubuntu______] 🔍 Sort: [Name_________] │
├─────────────────────────────────────────────────┤
│ Repository        │ Tag     │ Size   │ Created  │
├───────────────────┼─────────┼────────┼──────────┤
│ nginx             │ latest  │ 142MB  │ 2 days   │
│ nginx             │ 1.21    │ 140MB  │ 1 week   │
│ ubuntu            │ latest  │ 72MB   │ 3 days   │
│ ubuntu            │ 20.04   │ 72MB   │ 1 week   │
│ node              │ 16      │ 908MB  │ 5 days   │
│ mysql             │ 8.0     │ 516MB  │ 1 week   │
│ alpine            │ latest  │ 5MB    │ 6 days   │
│ 📦 <none>         │ <none>  │ 1.2GB  │ 2 weeks  │
├─────────────────────────────────────────────────┤
│ ↑/↓: Navigate | Enter: Details | Space: Select │
│ p: Pull | b: Build | t: Tag | h: Push | d: Del │
│ r: Run | s: Save | l: Load | i: Inspect        │
└─────────────────────────────────────────────────┘
```

### **Image Details Panel**
```
┌─────────────── Image Details: nginx:latest ──────────────┐
│                                                           │
│ Basic Information:                                        │
│ 🆔 ID: sha256:a1b2c3...                                   │
│ 🏷️  Repository: nginx                                    │
│ 🔖 Tag: latest                                           │
│ 📦 Size: 142MB (Virtual: 142MB)                          │
│ 📅 Created: 2 days ago                                   │
│ 👤 Author: NGINX Docker Maintainers                      │
│                                                           │
│ Configuration:                                            │
│ 🚪 Exposed Ports: 80/tcp                                │
│ 💻 Default CMD: ["nginx", "-g", "daemon off;"]          │
│ 📁 Working Dir: /                                        │
│ 👤 User: root                                            │
│ 🔧 Env: PATH=/usr/local/sbin:/usr/local/bin...          │
│                                                           │
│ Layer Information:                                        │
│ 📄 Layers: 6                                             │
│ 🔗 Parent: sha256:b1c2d3...                              │
│ 📊 Architecture: amd64                                    │
│ 🖥️  OS: linux                                            │
│                                                           │
│ Labels:                                                   │
│ 🏷️  maintainer=NGINX Docker Maintainers                 │
│ 🏷️  org.opencontainers.image.version=1.21.6            │
│                                                           │
│ Usage:                                                    │
│ 📦 Used by: 3 containers                                 │
│ 🔗 Children: 2 images                                    │
│                                                           │
├───────────────────────────────────────────────────────────┤
│ [History] [Layers] [Run] [Tag] [Push] [Export] [Delete] │
└───────────────────────────────────────────────────────────┘
```

### **Image Actions**
- **📥 Registry Operations**
  - 📤 Pull - Download image from registry
  - 📦 Push - Upload image to registry
  - 🔍 Search - Search registry for images
  - 🔑 Login - Authenticate with registry

- **🔨 Build Operations**
  - 🏗️ Build - Build image from Dockerfile
  - 🏷️ Tag - Add tags to image
  - 💾 Save - Export image to tar file
  - 📁 Load - Import image from tar file

- **🚀 Container Operations**
  - ▶️ Run - Create and start container
  - 🔧 Create - Create container without starting
  - 📋 Inspect - View detailed image information
  - 📜 History - View image layer history

## 💾 Volume Management

### **Volume List View**
```
┌─────────────────── Volumes ──────────────────────┐
│ Search: [data_______] 🔍 Filter: [All_________] │
├──────────────────────────────────────────────────┤
│ Name           │ Driver │ Size    │ Mount Point  │
├────────────────┼────────┼─────────┼──────────────┤
│ app-data       │ local  │ 2.1GB   │ /var/lib/... │
│ db-storage     │ local  │ 890MB   │ /var/lib/... │
│ logs-volume    │ local  │ 156MB   │ /var/lib/... │
│ config-files   │ local  │ 12MB    │ /var/lib/... │
│ 🔗 shared-nfs  │ nfs    │ 15GB    │ server:/data │
│ 📦 temp-cache  │ local  │ 500MB   │ /var/lib/... │
├──────────────────────────────────────────────────┤
│ ↑/↓: Navigate | Enter: Details | Space: Select  │
│ c: Create | d: Delete | p: Prune | b: Backup    │
│ m: Mount | u: Unmount | i: Inspect              │
└──────────────────────────────────────────────────┘
```

### **Volume Actions**
- **📦 Lifecycle**
  - 🆕 Create - Create new volume
  - 🗑️ Delete - Remove volume
  - 🧹 Prune - Remove unused volumes
  - 📋 Inspect - View volume details

- **🔗 Operations**
  - 📁 Browse - Explore volume contents
  - 💾 Backup - Create volume backup
  - 📥 Restore - Restore from backup
  - 📊 Usage - Show space usage

## 🌐 Network Management

### **Network List View**
```
┌─────────────────── Networks ─────────────────────┐
│ Search: [bridge____] 🔍 Filter: [All__________] │
├───────────────────────────────────────────────────┤
│ Name        │ Driver │ Scope │ Connected │ Subnet │
├─────────────┼────────┼───────┼───────────┼────────┤
│ bridge      │ bridge │ local │ 3         │ 172... │
│ host        │ host   │ local │ 0         │ -      │
│ none        │ null   │ local │ 0         │ -      │
│ web-net     │ bridge │ local │ 2         │ 192... │
│ api-network │ bridge │ local │ 4         │ 10.... │
├───────────────────────────────────────────────────┤
│ ↑/↓: Navigate | Enter: Details | Space: Select   │
│ c: Create | d: Delete | p: Prune | o: Connect    │
│ x: Disconnect | i: Inspect                       │
└───────────────────────────────────────────────────┘
```

## 📊 System Information

### **System Overview**
```
┌─────────────────── System Status ────────────────────┐
│                                                       │
│ Runtime Information:                                  │
│ 🆔 Version: Servin 1.0.0                             │
│ 🏗️  Build: go1.24.0 linux/amd64                      │
│ 📅 Started: 2024-01-15 10:30:15 (uptime: 2h 45m)    │
│ 🔌 CRI Server: Active on port 10010                  │
│ 📡 API Server: Active on unix socket                 │
│                                                       │
│ Resource Summary:                                     │
│ 📦 Containers: 12 total (8 running, 3 stopped, 1 paused) │
│ 🖼️  Images: 25 total (15 in use, 10 unused)          │
│ 💾 Volumes: 8 total (6 mounted, 2 unmounted)         │
│ 🌐 Networks: 5 total (3 active, 2 inactive)          │
│                                                       │
│ Storage Usage:                                        │
│ 📁 Images: 4.2GB                                     │
│ 📦 Containers: 1.8GB                                 │
│ 💾 Volumes: 3.5GB                                    │
│ 🏗️  Build Cache: 890MB                               │
│ 📊 Total: 10.39GB                                    │
│                                                       │
│ Performance Metrics:                                  │
│ 📊 CPU Usage: 15% (4 cores available)               │
│ 💾 Memory Usage: 2.1GB/8GB (26%)                    │
│ 💽 Disk I/O: 125MB/s read, 89MB/s write             │
│ 🌐 Network I/O: 45MB/s in, 32MB/s out               │
│                                                       │
├───────────────────────────────────────────────────────┤
│ [Events] [Logs] [Config] [Cleanup] [Export] [Quit]  │
└───────────────────────────────────────────────────────┘
```

### **Event Monitor**
```
┌─────────────────── Live Events ───────────────────────┐
│ 🟢 2024-01-15 13:45:23  container  web-server started │
│ 🟡 2024-01-15 13:44:58  image      nginx:latest pulled │
│ 🔴 2024-01-15 13:44:12  container  old-app stopped    │
│ 🟢 2024-01-15 13:43:45  volume     data-vol created   │
│ 🟡 2024-01-15 13:43:20  network    api-net connected  │
│ 🟢 2024-01-15 13:42:55  container  db-server started  │
│ 🔴 2024-01-15 13:42:30  container  temp-job exited(0) │
│ 🟡 2024-01-15 13:42:05  image      ubuntu:latest built │
├────────────────────────────────────────────────────────┤
│ ↑/↓: Scroll | f: Filter | c: Clear | s: Save | q: Quit │
└────────────────────────────────────────────────────────┘
```

## ⌨️ Keyboard Shortcuts

### **Global Navigation**
- **Tab** - Next panel/field
- **Shift+Tab** - Previous panel/field
- **Ctrl+C** - Cancel current operation
- **Escape** - Go back/cancel
- **q** - Quit application
- **?** - Show context help
- **/** - Search/filter current view

### **List Navigation**
- **↑/↓** - Move selection up/down
- **Page Up/Down** - Scroll page up/down
- **Home/End** - Go to first/last item
- **Enter** - Select item or open details
- **Space** - Toggle selection (multi-select)

### **Container Management**
- **s** - Start selected container
- **t** - Stop selected container
- **r** - Restart selected container
- **p** - Pause selected container
- **k** - Kill selected container
- **d** - Delete selected container
- **l** - View logs
- **e** - Execute shell
- **i** - Inspect details

### **Image Management**
- **p** - Pull image from registry
- **b** - Build image from Dockerfile
- **t** - Tag image
- **h** - Push image to registry
- **r** - Run container from image
- **d** - Delete image
- **s** - Save image to file
- **l** - Load image from file

### **System Operations**
- **F5** - Refresh current view
- **Ctrl+R** - Reload all data
- **Ctrl+L** - Clear screen
- **Ctrl+S** - Save current view to file
- **Ctrl+E** - Export system information

## 🎨 Customization

### **Theme Options**
The TUI supports multiple color schemes:
- **Default** - Standard terminal colors
- **Dark** - Dark theme with high contrast
- **Light** - Light theme for bright terminals
- **Monochrome** - Black and white for compatibility
- **Custom** - User-defined color scheme

### **Configuration**
```bash
# Set default theme
servin config set tui.theme dark

# Enable mouse support
servin config set tui.mouse true

# Set refresh interval
servin config set tui.refresh 2s

# Configure log tail lines
servin config set tui.logs.tail 100

# Set default container shell
servin config set tui.shell /bin/bash
```

### **Layout Customization**
- **Panel Arrangement** - Customize panel layout
- **Column Visibility** - Show/hide specific columns
- **Sort Options** - Default sorting preferences
- **Filter Presets** - Save commonly used filters
- **Hotkey Remapping** - Customize keyboard shortcuts

## 🔧 Advanced Features

### **Bulk Operations**
- **Multi-Selection** - Use Space to select multiple items
- **Bulk Actions** - Apply operations to selected items
- **Confirmation Dialogs** - Safety prompts for destructive actions
- **Progress Indicators** - Visual feedback for long operations

### **Search and Filtering**
- **Real-time Search** - Filter as you type
- **Advanced Filters** - Status, labels, dates
- **Regular Expressions** - Pattern-based filtering
- **Saved Searches** - Store frequently used filters

### **Integration Features**
- **Shell Integration** - Launch external commands
- **File Manager** - Browse container filesystems
- **Log Streaming** - Real-time log following
- **Statistics Charts** - ASCII-based performance graphs

---

## 📚 Next Steps

- **[Desktop GUI]({{ '/gui' | relative_url }})** - Explore the visual desktop interface
- **[CLI Reference]({{ '/cli' | relative_url }})** - Learn command-line operations
- **[Configuration]({{ '/configuration' | relative_url }})** - Customize your setup
- **[API Integration]({{ '/api' | relative_url }})** - Programmatic access

<div class="tui-tips">
  <h3>💡 TUI Pro Tips</h3>
  <ul>
    <li><strong>Mouse Support:</strong> Enable mouse support with <code>servin config set tui.mouse true</code></li>
    <li><strong>SSH Sessions:</strong> TUI works perfectly over SSH for remote management</li>
    <li><strong>Screen/Tmux:</strong> Run TUI in screen or tmux for persistent sessions</li>
    <li><strong>Context Help:</strong> Press <code>?</code> in any view for context-specific help</li>
    <li><strong>Log Monitoring:</strong> Use the log viewer for real-time container debugging</li>
  </ul>
</div>
