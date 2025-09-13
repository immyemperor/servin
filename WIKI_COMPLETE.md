# Servin Container Runtime - Complete Wiki

Welcome to the comprehensive wiki for Servin Container Runtime - a Docker-compatible container runtime with Kubernetes CRI support and professional desktop interface.

## 📖 Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture](#architecture)
3. [Features](#features)
4. [Installation](#installation)
5. [User Interfaces](#user-interfaces)
6. [Container Runtime Interface (CRI)](#container-runtime-interface-cri)
7. [Container Management](#container-management)
8. [Image Management](#image-management)
9. [Volume Management](#volume-management)
10. [Registry Operations](#registry-operations)
11. [Logging & Monitoring](#logging--monitoring)
12. [Configuration](#configuration)
13. [API Reference](#api-reference)
14. [Development](#development)
15. [Troubleshooting](#troubleshooting)
16. [Contributing](#contributing)

---

## 🚀 Project Overview

Servin Container Runtime is a comprehensive container management solution that provides Docker-compatible functionality with Kubernetes CRI support. Built in Go, it offers multiple interfaces including CLI, Terminal UI (TUI), and Desktop GUI applications.

### **Key Capabilities**
- **Container Runtime Interface (CRI)** - Full Kubernetes v1alpha2 compatibility
- **Docker-Compatible API** - Seamless migration from Docker workflows  
- **Multiple Interfaces** - CLI, TUI, and Desktop GUI applications
- **Cross-Platform** - Windows, Linux, and macOS support
- **Professional Installation** - Wizard-based installers for all platforms
- **Service Integration** - Background service with auto-start capabilities

### **Target Users**
- **Developers** - Container-based application development
- **DevOps Engineers** - Container orchestration and deployment
- **System Administrators** - Container infrastructure management
- **Kubernetes Users** - CRI-compatible runtime for clusters

---

## 🏗 Architecture

### **Core Components**

```
┌─────────────────────────────────────────────────────────────┐
│                    Servin Container Runtime                  │
├─────────────────────────────────────────────────────────────┤
│  Interfaces                                                 │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────────┐│
│  │    CLI      │ │     TUI     │ │      Desktop GUI        ││
│  │  Command    │ │  Terminal   │ │   Fyne-based Visual     ││
│  │   Line      │ │ Interface   │ │      Application        ││
│  └─────────────┘ └─────────────┘ └─────────────────────────┘│
├─────────────────────────────────────────────────────────────┤
│  Core Runtime Services                                      │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────────┐│
│  │ Container   │ │   Image     │ │       Volume            ││
│  │ Management  │ │ Management  │ │     Management          ││
│  └─────────────┘ └─────────────┘ └─────────────────────────┘│
├─────────────────────────────────────────────────────────────┤
│  API Layer                                                  │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────────┐│
│  │ CRI Server  │ │ HTTP API    │ │    gRPC Services        ││
│  │ (gRPC)      │ │ (REST)      │ │   (Internal Comms)      ││
│  └─────────────┘ └─────────────┘ └─────────────────────────┘│
├─────────────────────────────────────────────────────────────┤
│  Storage & Runtime                                          │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────────┐│
│  │ Container   │ │   Image     │ │      Configuration      ││
│  │   Storage   │ │   Store     │ │       & Metadata        ││
│  └─────────────┘ └─────────────┘ └─────────────────────────┘│
└─────────────────────────────────────────────────────────────┘
```

### **Directory Structure**

```
servin/
├── cmd/                          # Application entry points
│   ├── servin-desktop/          # Terminal UI application
│   ├── servin-gui/              # Desktop GUI application
│   └── gui.go                   # GUI command integration
├── pkg/                         # Core packages
│   ├── cri/                     # Container Runtime Interface
│   │   ├── server.go           # CRI gRPC server
│   │   ├── runtime.go          # Runtime service implementation
│   │   ├── image.go            # Image service implementation
│   │   └── types.go            # CRI type definitions
│   ├── container/              # Container management
│   ├── image/                  # Image management
│   ├── volume/                 # Volume management
│   └── registry/               # Registry operations
├── installers/                 # Installation wizards
│   ├── windows/               # NSIS-based Windows installer
│   ├── linux/                 # Python/Tkinter Linux installer
│   └── macos/                 # Python/Tkinter macOS installer
├── build/                     # Build artifacts
├── dist/                      # Distribution packages
├── docs/                      # Documentation
└── examples/                  # Usage examples
```

---

## ✨ Features

### **🎯 Core Runtime Features**
- ✅ Container lifecycle management (create, start, stop, delete)
- ✅ Image management (pull, push, build, tag, remove)
- ✅ Volume management (create, mount, unmount, remove)
- ✅ Network management (bridge, host, none)
- ✅ Registry operations (login, logout, push, pull)
- ✅ Container logs and monitoring
- ✅ Resource management (CPU, memory limits)
- ✅ Health checks and restart policies

### **🔌 Kubernetes Integration**
- ✅ Full CRI v1alpha2 implementation
- ✅ Pod sandbox management
- ✅ Container runtime service
- ✅ Image service with pull/list/remove
- ✅ gRPC API server on port 10010
- ✅ Compatible with kubelet

### **🖥 User Interfaces**

#### **Command Line Interface (CLI)**
```bash
servin --help                    # Show help
servin containers ls             # List containers
servin images pull ubuntu:latest # Pull image
servin run ubuntu:latest bash    # Run container
```

#### **Terminal User Interface (TUI)**
- Interactive menu-driven interface
- Real-time container status
- Visual container management
- Cross-platform terminal compatibility

#### **Desktop GUI Application**
- **Framework**: Fyne v2.6.3 (Go-native GUI)
- **Features**: Visual container management, image browser, log viewer
- **Platforms**: Windows, Linux, macOS
- **Architecture**: Multi-threaded with proper UI thread safety

### **📦 Installation Options**

#### **Windows**
- **NSIS Wizard Installer** - Professional Windows installer
- **Features**: Service installation, Start Menu integration, Uninstaller
- **Requirements**: Windows 10+, Administrator privileges

#### **Linux**
- **Python/Tkinter GUI Installer** - Cross-distribution compatibility
- **Features**: systemd service setup, desktop integration
- **Supported**: Ubuntu, CentOS, Debian, Fedora, Arch Linux

#### **macOS**
- **Python/Tkinter Native Installer** - Apple HIG compliant
- **Features**: Application bundle creation, launchd service
- **Requirements**: macOS 10.12+, Python 3 with tkinter

---

## 🛠 Installation

### **Quick Installation**

#### **Windows (Installer Wizard)**
1. Download `ServinSetup-1.0.0.exe`
2. Run as Administrator
3. Follow installation wizard
4. Launch from Start Menu

#### **Linux (GUI Installer)**
```bash
# Download and run installer
sudo python3 servin-installer.py
```

#### **macOS (GUI Installer)**
```bash
# Download and run installer
sudo python3 servin-installer.py
```

### **Building from Source**

#### **Prerequisites**
- **Go 1.24+** - Latest Go version
- **CGO enabled** - For GUI compilation
- **Git** - Source code management

#### **Windows Build Requirements**
```powershell
# Install MinGW-w64 UCRT64 toolchain
# Download from: https://www.mingw-w64.org/downloads/
# Add to PATH: C:\mingw64\bin

# Verify installation
gcc --version
go version
```

#### **Build Commands**
```bash
# Clone repository
git clone https://github.com/yourusername/servin
cd servin

# Build all components
.\build.ps1                    # Windows
./build.sh                     # Linux/macOS

# Build specific targets
go build -o servin main.go                    # CLI
go build -o servin-desktop cmd/servin-desktop/ # TUI
go build -o servin-gui cmd/servin-gui/         # GUI

# Build installer wizards
.\build-installer.ps1          # Windows NSIS installer
```

---

## 🖱 User Interfaces

### **1. Command Line Interface (CLI)**

The CLI provides the most comprehensive access to Servin functionality:

#### **Container Operations**
```bash
# Container management
servin run ubuntu:latest bash           # Run interactive container
servin run -d nginx:latest             # Run detached container
servin containers ls                   # List containers
servin containers stop <container-id>  # Stop container
servin containers rm <container-id>    # Remove container

# Container inspection
servin containers inspect <container-id>
servin containers logs <container-id>
servin containers stats <container-id>
```

#### **Image Operations**
```bash
# Image management
servin images pull ubuntu:latest       # Pull image
servin images ls                       # List images
servin images rm ubuntu:latest         # Remove image
servin images build -t myapp .         # Build image

# Image inspection
servin images inspect ubuntu:latest
servin images history ubuntu:latest
```

#### **Volume Operations**
```bash
# Volume management
servin volumes create myvolume         # Create volume
servin volumes ls                      # List volumes
servin volumes rm myvolume             # Remove volume
servin volumes inspect myvolume        # Inspect volume
```

### **2. Terminal User Interface (TUI)**

Interactive menu-driven interface for visual container management:

```
┌─────────────────────────────────────────────────────────────┐
│              Servin Container Runtime                       │
│                   Desktop Interface                         │
├─────────────────────────────────────────────────────────────┤
│  [1] Container Management                                   │
│      ├── List Containers                                   │
│      ├── Start Container                                   │
│      ├── Stop Container                                    │
│      └── Remove Container                                  │
│                                                             │
│  [2] Image Management                                       │
│      ├── List Images                                       │
│      ├── Pull Image                                        │
│      ├── Remove Image                                      │
│      └── Build Image                                       │
│                                                             │
│  [3] Volume Management                                      │
│      ├── List Volumes                                      │
│      ├── Create Volume                                     │
│      └── Remove Volume                                     │
│                                                             │
│  [4] Registry Operations                                    │
│      ├── Login to Registry                                 │
│      ├── Push Image                                        │
│      └── Logout from Registry                              │
│                                                             │
│  [5] System Information                                     │
│  [0] Exit                                                   │
└─────────────────────────────────────────────────────────────┘
```

**Launch TUI:**
```bash
servin-desktop    # Launch desktop TUI
servin desktop    # Alternative command
```

### **3. Desktop GUI Application**

Professional desktop application built with Fyne framework:

#### **Features**
- **Container Overview** - Visual container grid with status indicators
- **Image Browser** - Searchable image catalog with metadata
- **Volume Manager** - Drag-and-drop volume operations
- **Log Viewer** - Real-time container log streaming
- **Registry Explorer** - Browse and manage registry connections

#### **GUI Components**
```go
// Main application structure
type App struct {
    Window   fyne.Window
    Tabs     *container.AppTabs
    
    // Container management
    ContainerList *widget.List
    ContainerGrid *container.GridWithColumns
    
    // Image management  
    ImageList     *widget.List
    ImageSearch   *widget.Entry
    
    // Volume management
    VolumeList    *widget.List
    VolumeTree    *widget.Tree
    
    // Log viewer
    LogText       *widget.RichText
    LogScroll     *container.Scroll
}
```

**Launch GUI:**
```bash
servin-gui        # Launch desktop GUI
servin gui        # Alternative command
```

---

## 🔧 Container Runtime Interface (CRI)

Servin implements the Kubernetes Container Runtime Interface (CRI) v1alpha2 specification, providing full compatibility with kubelet.

### **CRI Server Configuration**

#### **Server Details**
- **Protocol**: gRPC
- **Port**: 10010 (default)
- **Socket**: `/var/run/servin.sock` (Linux/macOS)
- **Named Pipe**: `\\.\pipe\servin` (Windows)

#### **Start CRI Server**
```bash
servin daemon --cri-port 10010        # Start with custom port
servin daemon --cri-socket /custom/path/servin.sock
```

### **CRI Services Implementation**

#### **Runtime Service**
```go
type RuntimeService interface {
    // Pod sandbox management
    RunPodSandbox(config *runtimeapi.PodSandboxConfig) (string, error)
    StopPodSandbox(podSandboxID string) error
    RemovePodSandbox(podSandboxID string) error
    PodSandboxStatus(podSandboxID string) (*runtimeapi.PodSandboxStatus, error)
    
    // Container management
    CreateContainer(config *runtimeapi.ContainerConfig) (string, error)
    StartContainer(containerID string) error
    StopContainer(containerID string, timeout int64) error
    RemoveContainer(containerID string) error
    
    // Container inspection
    ContainerStatus(containerID string) (*runtimeapi.ContainerStatus, error)
    ListContainers(filter *runtimeapi.ContainerFilter) ([]*runtimeapi.Container, error)
    
    // Exec and logs
    ExecSync(containerID string, cmd []string, timeout time.Duration) (*runtimeapi.ExecSyncResponse, error)
    Exec(req *runtimeapi.ExecRequest) (*runtimeapi.ExecResponse, error)
}
```

#### **Image Service**
```go
type ImageService interface {
    // Image management
    PullImage(spec *runtimeapi.ImageSpec, auth *runtimeapi.AuthConfig) (string, error)
    RemoveImage(spec *runtimeapi.ImageSpec) error
    ImageStatus(spec *runtimeapi.ImageSpec) (*runtimeapi.Image, error)
    ListImages(filter *runtimeapi.ImageFilter) ([]*runtimeapi.Image, error)
    
    // Image filesystem
    ImageFsInfo() ([]*runtimeapi.FilesystemUsage, error)
}
```

### **Kubernetes Integration**

#### **kubelet Configuration**
```yaml
# /var/lib/kubelet/config.yaml
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
containerRuntime: remote
containerRuntimeEndpoint: unix:///var/run/servin.sock
imageServiceEndpoint: unix:///var/run/servin.sock
```

#### **Test CRI Connectivity**
```bash
# Install crictl (CRI CLI tool)
# Test runtime
crictl --runtime-endpoint unix:///var/run/servin.sock version
crictl --runtime-endpoint unix:///var/run/servin.sock info

# Test containers
crictl --runtime-endpoint unix:///var/run/servin.sock ps
crictl --runtime-endpoint unix:///var/run/servin.sock images
```

---

## 📦 Container Management

### **Container Lifecycle**

#### **1. Container Creation**
```bash
# Basic container creation
servin run ubuntu:latest
servin run -d --name webserver nginx:latest
servin run -it ubuntu:latest bash

# Advanced options
servin run \
  --name myapp \
  --port 8080:80 \
  --volume /data:/app/data \
  --env NODE_ENV=production \
  --memory 512m \
  --cpu 0.5 \
  node:16 npm start
```

#### **2. Container Operations**
```bash
# Container control
servin containers start <container-id>
servin containers stop <container-id>
servin containers restart <container-id>
servin containers pause <container-id>
servin containers unpause <container-id>

# Container cleanup
servin containers rm <container-id>
servin containers prune              # Remove stopped containers
```

#### **3. Container Inspection**
```bash
# Get container information
servin containers inspect <container-id>
servin containers logs <container-id>
servin containers stats <container-id>

# Execute commands in container
servin containers exec <container-id> bash
servin containers exec -it <container-id> sh
```

### **Container States**

| State | Description | Transitions |
|-------|-------------|-------------|
| `created` | Container created but not started | → `running`, `deleted` |
| `running` | Container is actively running | → `paused`, `stopped`, `deleted` |
| `paused` | Container is paused | → `running`, `deleted` |
| `stopped` | Container has stopped | → `running`, `deleted` |
| `deleted` | Container has been removed | (terminal state) |

### **Resource Management**

#### **CPU Limits**
```bash
# CPU allocation
servin run --cpu 1.5 ubuntu:latest      # 1.5 CPU cores
servin run --cpu-shares 512 ubuntu:latest # Relative CPU weight
```

#### **Memory Limits**
```bash
# Memory allocation
servin run --memory 512m ubuntu:latest   # 512 MB memory
servin run --memory 2g ubuntu:latest     # 2 GB memory
servin run --oom-kill-disable ubuntu:latest # Disable OOM killer
```

#### **Storage Limits**
```bash
# Disk I/O limits
servin run --device-read-bps /dev/sda:1mb ubuntu:latest
servin run --device-write-bps /dev/sda:1mb ubuntu:latest
```

---

## 🖼 Image Management

### **Image Operations**

#### **1. Image Retrieval**
```bash
# Pull images from registry
servin images pull ubuntu:latest
servin images pull docker.io/library/nginx:1.21
servin images pull gcr.io/my-project/my-app:v1.0

# Pull all tags
servin images pull --all-tags ubuntu
```

#### **2. Image Building**
```bash
# Build from Dockerfile
servin images build -t myapp:latest .
servin images build -t myapp:v1.0 -f Dockerfile.prod .

# Build with build arguments
servin images build \
  --build-arg NODE_VERSION=16 \
  --build-arg ENV=production \
  -t myapp:latest .
```

#### **3. Image Management**
```bash
# List images
servin images ls
servin images ls --filter "dangling=true"
servin images ls --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"

# Remove images
servin images rm ubuntu:latest
servin images rm $(servin images ls -q)  # Remove all images
servin images prune                       # Remove unused images
```

### **Image Formats & Standards**

#### **Supported Formats**
- **OCI Image Format** - Standard container image format
- **Docker Image Format** - Docker-compatible image format
- **Multi-architecture Images** - ARM64, AMD64, etc.

#### **Image Layers**
```bash
# Inspect image layers
servin images inspect ubuntu:latest
servin images history ubuntu:latest

# Show layer information
{
  "RootFS": {
    "Type": "layers",
    "Layers": [
      "sha256:...",
      "sha256:...",
      "sha256:..."
    ]
  }
}
```

### **Image Tagging & Distribution**

#### **Tagging**
```bash
# Tag images
servin images tag ubuntu:latest myregistry.com/ubuntu:latest
servin images tag myapp:latest myapp:v1.0
servin images tag myapp:latest myregistry.com/myapp:latest
```

#### **Registry Operations**
```bash
# Push to registry
servin images push myregistry.com/myapp:latest

# Save/Load images
servin images save -o myapp.tar myapp:latest
servin images load -i myapp.tar
```

---

## 💾 Volume Management

### **Volume Types**

#### **1. Named Volumes**
```bash
# Create named volumes
servin volumes create mydata
servin volumes create --driver local myvolume

# Use named volumes
servin run -v mydata:/data ubuntu:latest
```

#### **2. Bind Mounts**
```bash
# Mount host directories
servin run -v /host/path:/container/path ubuntu:latest
servin run -v "C:\Host\Path:/container/path" ubuntu:latest  # Windows
```

#### **3. Anonymous Volumes**
```bash
# Create anonymous volumes
servin run -v /container/path ubuntu:latest
```

### **Volume Operations**

#### **Volume Management**
```bash
# List volumes
servin volumes ls
servin volumes ls --filter "dangling=true"

# Inspect volumes
servin volumes inspect myvolume

# Remove volumes
servin volumes rm myvolume
servin volumes prune                    # Remove unused volumes
```

#### **Volume Configuration**
```json
{
  "Name": "myvolume",
  "Driver": "local",
  "Mountpoint": "/var/lib/servin/volumes/myvolume/_data",
  "Labels": {},
  "Options": {},
  "Scope": "local"
}
```

### **Volume Drivers**

#### **Local Driver**
- **Storage**: Host filesystem
- **Features**: Fast access, persistence across container restarts
- **Use Cases**: Development, single-node deployments

#### **Network Drivers** (Future)
- **NFS**: Network File System
- **CIFS/SMB**: Windows file sharing
- **GlusterFS**: Distributed storage

### **Volume Best Practices**

#### **Data Persistence**
```bash
# Database volumes
servin run -v postgres_data:/var/lib/postgresql/data postgres:13

# Application data
servin run -v app_data:/app/data -v app_logs:/app/logs myapp:latest
```

#### **Configuration Management**
```bash
# Configuration files
servin run -v ./config:/app/config:ro myapp:latest  # Read-only mount

# Secrets management
servin run -v ./secrets:/run/secrets:ro myapp:latest
```

---

## 🏪 Registry Operations

### **Registry Configuration**

#### **Default Registries**
- **Docker Hub** - `docker.io` (default)
- **GitHub Container Registry** - `ghcr.io`
- **Google Container Registry** - `gcr.io`
- **Amazon ECR** - `<account>.dkr.ecr.<region>.amazonaws.com`

#### **Registry Authentication**
```bash
# Login to registries
servin registry login docker.io
servin registry login -u username -p password registry.example.com

# Use environment variables
export SERVIN_REGISTRY_USER=username
export SERVIN_REGISTRY_PASSWORD=password
servin registry login registry.example.com

# Use token authentication
servin registry login -u oauth2accesstoken --password-stdin gcr.io
```

### **Image Distribution**

#### **Push Images**
```bash
# Tag and push
servin images tag myapp:latest registry.example.com/myapp:latest
servin images push registry.example.com/myapp:latest

# Push all tags
servin images push --all-tags registry.example.com/myapp
```

#### **Pull Images**
```bash
# Pull from specific registry
servin images pull registry.example.com/myapp:latest
servin images pull gcr.io/my-project/my-app:v1.0

# Pull with authentication
servin images pull --username myuser --password mypass private-registry.com/app:latest
```

### **Registry Management**

#### **Registry Configuration File**
```json
{
  "registries": {
    "docker.io": {
      "mirrors": ["https://mirror1.docker.io", "https://mirror2.docker.io"],
      "secure": true,
      "ca_file": "/path/to/ca.crt"
    },
    "registry.example.com": {
      "secure": true,
      "username": "myuser",
      "password": "mypass"
    },
    "insecure-registry.local": {
      "secure": false
    }
  }
}
```

#### **Registry Operations**
```bash
# List registry credentials
servin registry ls

# Logout from registry
servin registry logout docker.io
servin registry logout --all

# Test registry connectivity
servin registry ping registry.example.com
```

---

## 📊 Logging & Monitoring

### **Container Logs**

#### **Log Retrieval**
```bash
# View container logs
servin containers logs <container-id>
servin containers logs -f <container-id>          # Follow logs
servin containers logs --tail 100 <container-id>  # Last 100 lines
servin containers logs --since 1h <container-id>  # Last hour

# Filter logs
servin containers logs --timestamps <container-id>
servin containers logs --details <container-id>
```

#### **Log Drivers**
```bash
# Configure log driver
servin run --log-driver json-file --log-opt max-size=10m ubuntu:latest
servin run --log-driver syslog --log-opt syslog-address=udp://192.168.1.42:514 ubuntu:latest
```

### **System Monitoring**

#### **Container Statistics**
```bash
# Real-time stats
servin containers stats
servin containers stats <container-id>
servin containers stats --no-stream

# Stats output format
CONTAINER ID   NAME        CPU %    MEM USAGE / LIMIT     MEM %    NET I/O       BLOCK I/O
abc123def456   webserver   0.50%    256MiB / 2GiB        12.80%   1.2kB / 648B  0B / 0B
```

#### **System Information**
```bash
# System overview
servin system info
servin system df                    # Disk usage
servin system events               # System events
servin system prune                # Clean up system

# Version information
servin version
servin version --format json
```

### **Health Monitoring**

#### **Health Checks**
```bash
# Configure health checks
servin run \
  --health-cmd "curl -f http://localhost:8080/health" \
  --health-interval 30s \
  --health-timeout 3s \
  --health-retries 3 \
  myapp:latest
```

#### **Container Health Status**
```bash
# Check container health
servin containers inspect <container-id> | grep -A 10 '"Health"'

# Health states: starting, healthy, unhealthy
```

---

## ⚙️ Configuration

### **Global Configuration**

#### **Configuration File Locations**
- **Linux**: `/etc/servin/config.yaml` or `~/.config/servin/config.yaml`
- **macOS**: `/etc/servin/config.yaml` or `~/Library/Application Support/servin/config.yaml`
- **Windows**: `C:\ProgramData\Servin\config\servin.conf` or `%APPDATA%\servin\config.yaml`

#### **Configuration Format**
```yaml
# Servin Container Runtime Configuration
version: "1.0"

# Runtime configuration
runtime:
  data_dir: "/var/lib/servin"
  state_dir: "/run/servin"
  log_level: "info"
  max_containers: 1000

# CRI configuration
cri:
  enabled: true
  listen_address: "0.0.0.0:10010"
  socket_path: "/var/run/servin.sock"
  
# Registry configuration
registries:
  default: "docker.io"
  insecure:
    - "registry.local:5000"
  mirrors:
    "docker.io":
      - "https://mirror1.docker.io"
      - "https://mirror2.docker.io"

# Network configuration
networking:
  default_bridge: "servin0"
  bridge_ip: "172.17.0.1/16"
  enable_ipv6: false

# Storage configuration
storage:
  driver: "overlay2"
  root: "/var/lib/servin"
  options:
    overlay2.size: "10G"

# Logging configuration
logging:
  driver: "json-file"
  level: "info"
  max_size: "10m"
  max_files: 5
  
# Security configuration
security:
  enable_user_namespaces: true
  default_ulimits:
    nofile: 1024:4096
    nproc: 2048:4096
```

### **Environment Variables**

#### **Runtime Variables**
```bash
# Data directories
export SERVIN_ROOT="/var/lib/servin"
export SERVIN_STATE_DIR="/run/servin"
export SERVIN_CONFIG_DIR="/etc/servin"

# Logging
export SERVIN_LOG_LEVEL="debug"
export SERVIN_LOG_FORMAT="json"

# CRI settings
export SERVIN_CRI_LISTEN="0.0.0.0:10010"
export SERVIN_CRI_SOCKET="/var/run/servin.sock"

# Registry authentication
export SERVIN_REGISTRY_USER="username"
export SERVIN_REGISTRY_PASSWORD="password"
export SERVIN_REGISTRY_AUTH_FILE="/etc/servin/auth.json"
```

### **Service Configuration**

#### **systemd Service (Linux)**
```ini
# /etc/systemd/system/servin.service
[Unit]
Description=Servin Container Runtime
Documentation=https://github.com/yourusername/servin
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/servin daemon
ExecReload=/bin/kill -HUP $MAINPID
KillMode=process
Restart=on-failure
RestartSec=5
LimitNOFILE=1048576
LimitNPROC=1048576
LimitCORE=infinity
TasksMax=infinity
TimeoutStartSec=0
Delegate=yes
OOMScoreAdjust=-999

[Install]
WantedBy=multi-user.target
```

#### **Windows Service**
```xml
<!-- Servin Windows Service Configuration -->
<service>
  <id>ServinRuntime</id>
  <name>Servin Container Runtime</name>
  <description>Servin Container Runtime Service</description>
  <executable>C:\Program Files\Servin\servin.exe</executable>
  <arguments>daemon --service</arguments>
  <startmode>Automatic</startmode>
  <logpath>C:\ProgramData\Servin\logs</logpath>
  <logmode>rotate</logmode>
</service>
```

---

## 🔌 API Reference

### **HTTP REST API**

#### **Base URL**
```
http://localhost:8080/api/v1
```

#### **Authentication**
```bash
# API token authentication
curl -H "Authorization: Bearer <token>" \
     http://localhost:8080/api/v1/containers
```

### **Container Endpoints**

#### **List Containers**
```http
GET /api/v1/containers
```

**Query Parameters:**
- `all` (boolean) - Show all containers (default: only running)
- `limit` (integer) - Maximum number of containers to return
- `filters` (string) - JSON encoded filters

**Response:**
```json
{
  "containers": [
    {
      "id": "abc123def456",
      "name": "webserver",
      "image": "nginx:latest",
      "state": "running",
      "status": "Up 2 hours",
      "ports": [
        {
          "private_port": 80,
          "public_port": 8080,
          "type": "tcp"
        }
      ],
      "created": "2025-09-13T10:00:00Z"
    }
  ]
}
```

#### **Create Container**
```http
POST /api/v1/containers/create
```

**Request Body:**
```json
{
  "name": "mycontainer",
  "image": "ubuntu:latest",
  "cmd": ["/bin/bash"],
  "env": ["NODE_ENV=production"],
  "ports": {
    "80/tcp": [{"HostPort": "8080"}]
  },
  "volumes": {
    "/host/path": {
      "bind": "/container/path",
      "mode": "rw"
    }
  },
  "restart_policy": {
    "name": "unless-stopped"
  }
}
```

#### **Start Container**
```http
POST /api/v1/containers/{id}/start
```

#### **Stop Container**
```http
POST /api/v1/containers/{id}/stop?t=10
```

#### **Remove Container**
```http
DELETE /api/v1/containers/{id}?force=true
```

### **Image Endpoints**

#### **List Images**
```http
GET /api/v1/images
```

#### **Pull Image**
```http
POST /api/v1/images/create?fromImage=ubuntu&tag=latest
```

#### **Remove Image**
```http
DELETE /api/v1/images/{name}
```

### **Volume Endpoints**

#### **List Volumes**
```http
GET /api/v1/volumes
```

#### **Create Volume**
```http
POST /api/v1/volumes/create
```

**Request Body:**
```json
{
  "name": "myvolume",
  "driver": "local",
  "labels": {
    "environment": "production"
  }
}
```

### **gRPC CRI API**

#### **Runtime Service**
```protobuf
service RuntimeService {
    rpc Version(VersionRequest) returns (VersionResponse) {}
    rpc RunPodSandbox(RunPodSandboxRequest) returns (RunPodSandboxResponse) {}
    rpc StopPodSandbox(StopPodSandboxRequest) returns (StopPodSandboxResponse) {}
    rpc RemovePodSandbox(RemovePodSandboxRequest) returns (RemovePodSandboxResponse) {}
    rpc PodSandboxStatus(PodSandboxStatusRequest) returns (PodSandboxStatusResponse) {}
    rpc ListPodSandbox(ListPodSandboxRequest) returns (ListPodSandboxResponse) {}
    
    rpc CreateContainer(CreateContainerRequest) returns (CreateContainerResponse) {}
    rpc StartContainer(StartContainerRequest) returns (StartContainerResponse) {}
    rpc StopContainer(StopContainerRequest) returns (StopContainerResponse) {}
    rpc RemoveContainer(RemoveContainerRequest) returns (RemoveContainerResponse) {}
    rpc ListContainers(ListContainersRequest) returns (ListContainersResponse) {}
    rpc ContainerStatus(ContainerStatusRequest) returns (ContainerStatusResponse) {}
    
    rpc ExecSync(ExecSyncRequest) returns (ExecSyncResponse) {}
    rpc Exec(ExecRequest) returns (ExecResponse) {}
}
```

#### **Image Service**
```protobuf
service ImageService {
    rpc ListImages(ListImagesRequest) returns (ListImagesResponse) {}
    rpc ImageStatus(ImageStatusRequest) returns (ImageStatusResponse) {}
    rpc PullImage(PullImageRequest) returns (PullImageResponse) {}
    rpc RemoveImage(RemoveImageRequest) returns (RemoveImageResponse) {}
    rpc ImageFsInfo(ImageFsInfoRequest) returns (ImageFsInfoResponse) {}
}
```

---

## 🛠 Development

### **Development Setup**

#### **Prerequisites**
- **Go 1.24+** - Latest Go version with generics support
- **Git** - Version control
- **Make** - Build automation (optional)
- **CGO** - Required for GUI components

#### **Development Environment**
```bash
# Clone repository
git clone https://github.com/yourusername/servin
cd servin

# Install dependencies
go mod download
go mod verify

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/goreleaser/goreleaser@latest
```

### **Project Structure**

#### **Core Packages**
```
pkg/
├── cri/                    # Container Runtime Interface
│   ├── server.go          # gRPC server implementation
│   ├── runtime.go         # Runtime service
│   ├── image.go           # Image service
│   └── types.go           # CRI type definitions
├── container/             # Container management
│   ├── manager.go         # Container lifecycle
│   ├── runtime.go         # Runtime operations
│   └── config.go          # Container configuration
├── image/                 # Image management
│   ├── manager.go         # Image operations
│   ├── builder.go         # Image building
│   └── registry.go        # Registry operations
├── volume/                # Volume management
│   ├── manager.go         # Volume operations
│   └── drivers/           # Volume drivers
└── network/               # Network management
    ├── bridge.go          # Bridge networking
    └── manager.go         # Network operations
```

### **Building Components**

#### **CLI Application**
```bash
# Build CLI
go build -o servin main.go

# Build with version info
go build -ldflags "-X main.version=1.0.0 -X main.commit=$(git rev-parse HEAD)" -o servin main.go
```

#### **Desktop Applications**
```bash
# Build TUI
go build -o servin-desktop cmd/servin-desktop/

# Build GUI (requires CGO)
CGO_ENABLED=1 go build -o servin-gui cmd/servin-gui/
```

#### **Cross-Platform Building**
```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o servin.exe main.go

# Linux
GOOS=linux GOARCH=amd64 go build -o servin main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o servin main.go
```

### **Testing**

#### **Unit Tests**
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./pkg/container/
go test ./pkg/cri/
```

#### **Integration Tests**
```bash
# Run integration tests
go test -tags=integration ./test/integration/

# Run CRI tests
go test ./test/cri/
```

#### **End-to-End Tests**
```bash
# Run E2E tests
./test/e2e/run-tests.sh
```

### **Code Quality**

#### **Linting**
```bash
# Run linter
golangci-lint run

# Run with auto-fix
golangci-lint run --fix
```

#### **Formatting**
```bash
# Format code
go fmt ./...
goimports -w .
```

#### **Documentation**
```bash
# Generate documentation
godoc -http=:6060

# View package documentation
go doc pkg/container
go doc pkg/cri.RuntimeService
```

---

## 🔧 Troubleshooting

### **Common Issues**

#### **Installation Problems**

**Issue: NSIS not found**
```
Error: NSIS (Nullsoft Scriptable Install System) is not installed!
```
**Solution:**
1. Download NSIS from https://nsis.sourceforge.io/Download
2. Install latest version
3. Verify installation: `makensis.exe` in PATH

**Issue: CGO compilation failed**
```
# github.com/fyne-io/fyne/v2
cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in %PATH%
```
**Solution (Windows):**
1. Install MinGW-w64 UCRT64 toolchain
2. Add `C:\mingw64\bin` to PATH
3. Verify: `gcc --version`

#### **Runtime Issues**

**Issue: Permission denied**
```
Error: permission denied: cannot connect to servin daemon
```
**Solution:**
- **Linux/macOS**: Run with `sudo` or add user to `servin` group
- **Windows**: Run as Administrator

**Issue: Port already in use**
```
Error: bind: address already in use (port 10010)
```
**Solution:**
```bash
# Find process using port
netstat -tlnp | grep 10010    # Linux
netstat -ano | findstr 10010  # Windows

# Kill process or use different port
servin daemon --cri-port 10011
```

**Issue: Container fails to start**
```
Error: container failed to start: executable file not found in $PATH
```
**Solution:**
- Check image exists: `servin images ls`
- Verify command path: `servin containers inspect <id>`
- Use absolute paths in container commands

#### **CRI Issues**

**Issue: kubelet cannot connect**
```
kubelet: Failed to connect to servin CRI: connection refused
```
**Solution:**
1. Verify CRI server is running: `servin daemon --cri-port 10010`
2. Check socket permissions: `ls -l /var/run/servin.sock`
3. Test connectivity: `crictl --runtime-endpoint unix:///var/run/servin.sock version`

**Issue: Image pull fails**
```
Error: failed to pull image: unauthorized
```
**Solution:**
```bash
# Login to registry
servin registry login docker.io
servin registry login -u username -p password registry.example.com

# Check credentials
servin registry ls
```

### **Debug Mode**

#### **Enable Debug Logging**
```bash
# Environment variable
export SERVIN_LOG_LEVEL=debug

# Command line flag
servin --log-level debug daemon

# Configuration file
# config.yaml
logging:
  level: debug
```

#### **Debug Information**
```bash
# System information
servin system info
servin system df
servin version

# Container debugging
servin containers inspect <container-id>
servin containers logs --details <container-id>

# Process information
ps aux | grep servin
systemctl status servin    # Linux
sc query ServinRuntime     # Windows
```

### **Performance Issues**

#### **High Memory Usage**
```bash
# Check container memory usage
servin containers stats
servin system df

# Limit container memory
servin run --memory 512m ubuntu:latest
```

#### **Slow Image Operations**
```bash
# Clean up unused resources
servin system prune
servin images prune
servin volumes prune

# Use registry mirrors
# config.yaml
registries:
  mirrors:
    "docker.io":
      - "https://mirror.gcr.io"
```

### **Log Analysis**

#### **Log Locations**
- **Linux**: `/var/log/servin/`, `journalctl -u servin`
- **macOS**: `/usr/local/var/log/servin/`, `launchctl logs servin`  
- **Windows**: `C:\ProgramData\Servin\logs\`, Event Viewer

#### **Log Commands**
```bash
# Follow logs
tail -f /var/log/servin/servin.log         # Linux
Get-Content -Wait C:\ProgramData\Servin\logs\servin.log  # Windows

# System logs
journalctl -u servin -f                    # Linux systemd
launchctl logs -f servin                   # macOS
```

---

## 🤝 Contributing

### **Contribution Guidelines**

#### **Getting Started**
1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/yourusername/servin`
3. **Create** feature branch: `git checkout -b feature/amazing-feature`
4. **Make** your changes
5. **Test** thoroughly
6. **Submit** pull request

#### **Development Workflow**
```bash
# Setup development environment
git clone https://github.com/yourusername/servin
cd servin
go mod download

# Create feature branch
git checkout -b feature/new-feature

# Make changes and test
go test ./...
golangci-lint run

# Commit changes
git add .
git commit -m "feat: add amazing new feature"

# Push and create PR
git push origin feature/new-feature
```

### **Code Standards**

#### **Go Style Guide**
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` and `goimports` for formatting
- Write comprehensive tests
- Document public APIs

#### **Commit Messages**
```
type(scope): description

feat(container): add container restart functionality
fix(cri): resolve pod sandbox creation issue  
docs(api): update REST API documentation
test(image): add image building test cases
```

#### **Pull Request Template**
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature  
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added/updated
```

### **Areas for Contribution**

#### **High Priority**
- **Network Plugins** - Additional networking drivers
- **Storage Drivers** - Cloud storage integration
- **Security Features** - Enhanced security and isolation
- **Performance** - Optimization and profiling

#### **Medium Priority**
- **Registry Support** - Additional registry backends
- **GUI Enhancements** - Advanced desktop features
- **CLI Improvements** - Better user experience
- **Documentation** - Tutorials and guides

#### **Future Features**
- **Swarm Mode** - Multi-node orchestration
- **Compose Support** - Docker Compose compatibility
- **Plugin System** - Extensible architecture
- **Monitoring** - Metrics and alerting

### **Community**

#### **Communication**
- **GitHub Issues** - Bug reports and feature requests
- **GitHub Discussions** - General questions and discussions
- **Discord** - Real-time community chat
- **Stack Overflow** - Technical questions (tag: servin-runtime)

#### **Resources**
- **Documentation** - Complete API and usage docs
- **Examples** - Sample configurations and use cases
- **Blog** - Development updates and tutorials
- **Roadmap** - Future development plans

---

## 📚 Additional Resources

### **Documentation**
- [Installation Guide](INSTALL.md)
- [CRI Implementation](CRI.md)
- [Desktop GUI Guide](DESKTOP_GUI.md)
- [API Reference](docs/api/)
- [Configuration Reference](docs/config/)

### **Examples**
- [Basic Container Usage](examples/basic/)
- [Kubernetes Integration](examples/kubernetes/)
- [Multi-Container Applications](examples/multi-container/)
- [Registry Integration](examples/registry/)

### **Tools & Integrations**
- **crictl** - CRI debugging tool
- **kubectl** - Kubernetes command line
- **Docker Compose** - Multi-container applications (planned)
- **Portainer** - Web-based management (planned)

---

## 📄 License

Servin Container Runtime is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- **Kubernetes SIG Node** - CRI specification
- **containerd** - Container runtime architecture inspiration
- **Docker** - Container ecosystem and APIs
- **Fyne** - Cross-platform GUI framework
- **Go Community** - Programming language and ecosystem

---

**Built with ❤️ using Go and modern container technologies**
wiki pahe here

*For support, please open an issue on GitHub or join our community discussions.*
