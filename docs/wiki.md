# 🚀 Servin Container Runtime - Complete Wiki

**Version 1.0.0**

A Docker-compatible container runtime with Kubernetes CRI support and professional desktop interface.

---

## Table of Contents

### Overview
- [Project Overview](#project-overview)
- [Architecture](#architecture)
- [Features](#features)

### Getting Started
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)

### User Interfaces
- [Command Line (CLI)](#command-line-cli)
- [Terminal UI (TUI)](#terminal-ui-tui)
- [Desktop GUI](#desktop-gui)

### Core Features
- [Container Management](#container-management)
- [Image Management](#image-management)
- [Volume Management](#volume-management)
- [Registry Operations](#registry-operations)

### Integration
- [Kubernetes CRI](#kubernetes-cri)
- [API Reference](#api-reference)
- [Logging & Monitoring](#logging--monitoring)

### Development
- [Building from Source](#building-from-source)
- [Development Guide](#development-guide)
- [Contributing](#contributing)

### Support
- [Troubleshooting](#troubleshooting)
- [FAQ](#faq)
- [Resources](#resources)

---

## Project Overview

Servin Container Runtime is a Docker-compatible container runtime with Kubernetes CRI support and professional desktop interface.

### Key Capabilities

#### 🎯 Core Runtime Features
- ✅ Container lifecycle management
- ✅ Image management (pull, push, build)
- ✅ Volume management
- ✅ Network management
- ✅ Registry operations

#### 🔌 Integration Features
- ✅ Kubernetes CRI v1alpha2
- ✅ Cross-platform support
- ✅ Service integration
- ✅ Professional installers
- ✅ REST and gRPC APIs

### Target Users

| User Type | Use Case |
|-----------|----------|
| 👨‍💻 **Developers** | Container-based application development |
| ⚙️ **DevOps Engineers** | Container orchestration and deployment |
| 🔧 **System Admins** | Container infrastructure management |
| ☸️ **Kubernetes Users** | CRI-compatible runtime for clusters |

---

## Architecture

### System Overview

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

### Component Architecture

#### 🎯 Core Components
- **Runtime Engine** - Container lifecycle management
- **Image Manager** - OCI image handling and storage
- **Volume Manager** - Persistent storage management
- **Network Manager** - Container networking
- **Registry Client** - Image registry operations

#### 🔌 Interface Layer
- **CLI Interface** - Command-line tool
- **TUI Interface** - Interactive terminal UI
- **GUI Interface** - Desktop application
- **HTTP API** - REST endpoints
- **gRPC CRI** - Kubernetes integration

### Directory Structure

```
servin/
├── cmd/                          # Application entry points
│   ├── servin-desktop/          # Terminal UI application
│   ├── servin-gui/              # Desktop GUI application
│   └── gui.go                   # GUI command integration
├── pkg/                         # Core packages
│   ├── cri/                     # Container Runtime Interface
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
└── docs/                      # Documentation
```

---

## Features

### 🎯 Core Runtime Features

#### Container Operations
- ✅ Create, start, stop, delete containers
- ✅ Container logs and monitoring
- ✅ Resource management (CPU, memory)
- ✅ Health checks and restart policies
- ✅ Container exec and attach

#### Image Operations
- ✅ Pull, push, build, tag images
- ✅ Multi-architecture support
- ✅ Layer caching and optimization
- ✅ Registry authentication
- ✅ Image inspection and history

### 🔌 Kubernetes Integration

#### CRI v1alpha2 Implementation

| Component | Features |
|-----------|----------|
| **Pod Management** | Pod sandbox lifecycle<br>Container creation/deletion<br>Pod networking setup |
| **Image Service** | Image pulling and removal<br>Image status queries<br>Filesystem usage reporting |
| **Runtime Service** | Container execution<br>Log streaming<br>Resource monitoring |

### 🖥 User Interface Options

| Interface | Description | Status |
|-----------|-------------|---------|
| 💻 **Command Line Interface** | Comprehensive CLI with Docker-compatible commands for automation and scripting | Full featured |
| 📱 **Terminal User Interface** | Interactive menu-driven interface for visual container management in the terminal | Cross-platform |
| 🖼️ **Desktop GUI Application** | Professional desktop application built with Fyne framework for visual management | Native look |

### 📦 Installation & Distribution

#### 🪟 Windows
- **NSIS Installer** - Professional wizard
- **Service Integration** - Windows Service
- **Start Menu** - Shortcuts and uninstaller
- **Registry Integration** - Add/Remove Programs

#### 🐧 Linux
- **GUI Installer** - Python/Tkinter wizard
- **systemd Service** - Background service
- **Desktop Integration** - Application menus
- **Cross-distribution** - Ubuntu, CentOS, Debian

#### 🍎 macOS
- **Native Installer** - Apple HIG compliant
- **launchd Service** - Background daemon
- **App Bundle** - Professional .app package
- **Retina Ready** - High-DPI displays

---

## Installation

### Quick Installation

#### 🪟 Windows Installation
1. Download `ServinSetup-1.0.0.exe`
2. Run as Administrator
3. Follow installation wizard
4. Launch from Start Menu

*Status: Wizard installer*

#### 🐧 Linux Installation
```bash
# Download and run installer
sudo python3 servin-installer.py
```

*Status: GUI installer*

#### 🍎 macOS Installation
```bash
# Download and run installer
sudo python3 servin-installer.py
```

*Status: Native installer*

### Building from Source

#### Prerequisites

**Required Software:**
- **Go 1.24+** - Latest Go version
- **Git** - Source code management
- **CGO enabled** - For GUI compilation

**Windows Additional Requirements:**
- **MinGW-w64 UCRT64** - GCC toolchain
- **NSIS** - For installer creation
- **PATH configuration** - Add mingw64/bin to PATH

### Build Commands

#### 🔨 Building Applications

```bash
# Clone repository
git clone https://github.com/yourusername/servin
cd servin

# Build all components
.\build.ps1                    # Windows
./build.sh                     # Linux/macOS

# Build specific components
go build -o servin main.go                    # CLI
go build -o servin-desktop cmd/servin-desktop/ # TUI
go build -o servin-gui cmd/servin-gui/         # GUI
```

#### 📦 Building Installers

```bash
# Windows NSIS installer
.\build-installer.ps1

# Creates: dist\ServinSetup-1.0.0.exe

# Cross-platform packages
.\build.ps1 -Target all        # All platforms
.\build.ps1 -Target windows    # Windows only
.\build.ps1 -Target linux      # Linux only
.\build.ps1 -Target macos      # macOS only
```

### Verification

#### Testing Installation

```bash
# Check version
servin version

# Test CLI
servin containers ls
servin images ls

# Test TUI
servin-desktop

# Test GUI
servin-gui

# Test CRI (if enabled)
crictl --runtime-endpoint unix:///var/run/servin.sock version
```

---

## Configuration

### Configuration Overview

Servin uses YAML configuration files to manage runtime settings, daemon options, registry configurations, and network settings. Configuration can be specified via files, environment variables, or command-line flags.

### Configuration File Locations

#### 🪟 Windows
- **System:** `C:\ProgramData\Servin\config.yaml`
- **User:** `%USERPROFILE%\.servin\config.yaml`
- **Working Dir:** `.\servin.yaml`

*Priority: Highest to lowest*

#### 🐧 Linux
- **System:** `/etc/servin/config.yaml`
- **User:** `~/.config/servin/config.yaml`
- **Working Dir:** `./servin.yaml`

*Priority: Highest to lowest*

#### 🍎 macOS
- **System:** `/etc/servin/config.yaml`
- **User:** `~/Library/Application Support/servin/config.yaml`
- **Working Dir:** `./servin.yaml`

*Priority: Highest to lowest*

### Main Configuration File

#### 📄 config.yaml

```yaml
# Servin Configuration File
# /etc/servin/config.yaml

# Daemon settings
daemon:
  # Data directory for containers, images, volumes
  data_root: "/var/lib/servin"
  
  # Runtime configuration
  runtime: "runc"
  runtime_args: []
  
  # Logging configuration
  log_level: "info"
  log_file: "/var/log/servin/servin.log"
  log_format: "text"  # text or json
  
  # Process limits
  max_concurrent_downloads: 6
  max_concurrent_uploads: 6

# Network configuration
network:
  # Default network for containers
  default_network: "bridge"
  
  # Bridge network settings
  bridge:
    name: "servin0"
    subnet: "172.17.0.0/16"
    gateway: "172.17.0.1"
    ip_masq: true
    icc: true
    mtu: 1500

# Storage configuration
storage:
  # Storage driver
  driver: "overlay2"
  
  # Storage options
  opts:
    overlay2.override_kernel_check: "true"
    overlay2.size: "120G"
  
  # Image storage settings
  image_store: "/var/lib/servin/images"
  
  # Volume storage settings
  volume_store: "/var/lib/servin/volumes"

# Registry configuration
registries:
  # Default registry
  default: "docker.io"
  
  # Registry mirrors
  mirrors:
    "docker.io":
      - "https://mirror.example.com"
  
  # Insecure registries (HTTP)
  insecure:
    - "registry.internal:5000"
  
  # Registry authentication
  auths:
    "registry.example.com":
      username: "user"
      password: "pass"

# CRI configuration
cri:
  # Enable CRI server
  enabled: true
  
  # CRI socket path
  socket_path: "/var/run/servin/servin.sock"
  
  # CRI server settings
  server:
    address: "127.0.0.1"
    port: 10010
    
  # Image service configuration
  image_service:
    max_parallel_downloads: 6
    
  # Runtime service configuration
  runtime_service:
    cgroup_manager: "systemd"
    default_runtime: "runc"

# API server configuration
api:
  # Enable REST API server
  enabled: true
  
  # API server settings
  server:
    address: "0.0.0.0"
    port: 8080
    
  # TLS configuration
  tls:
    enabled: false
    cert_file: "/etc/servin/tls/server.crt"
    key_file: "/etc/servin/tls/server.key"
    
  # Authentication
  auth:
    method: "none"  # none, basic, token, jwt
    token_file: "/etc/servin/tokens.json"

# Metrics and monitoring
metrics:
  # Enable metrics endpoint
  enabled: true
  
  # Metrics server settings
  server:
    address: "127.0.0.1"
    port: 9090
    path: "/metrics"
```

### Environment Variables

#### 🌍 Runtime Environment

```bash
# Data directory
export SERVIN_DATA_ROOT="/var/lib/servin"

# Log level
export SERVIN_LOG_LEVEL="debug"

# Log file
export SERVIN_LOG_FILE="/var/log/servin/servin.log"

# Configuration file
export SERVIN_CONFIG="/etc/servin/config.yaml"

# Runtime
export SERVIN_RUNTIME="runc"

# Network settings
export SERVIN_BRIDGE_NAME="servin0"
export SERVIN_BRIDGE_SUBNET="172.17.0.0/16"
```

#### 🔐 Security Environment

```bash
# API authentication
export SERVIN_API_TOKEN="your-api-token"

# Registry credentials
export SERVIN_REGISTRY_USER="username"
export SERVIN_REGISTRY_PASS="password"

# TLS certificates
export SERVIN_TLS_CERT="/path/to/cert.pem"
export SERVIN_TLS_KEY="/path/to/key.pem"

# Container runtime
export SERVIN_RUNTIME_ROOT="/run/servin"
export SERVIN_RUNTIME_STATE="/var/run/servin"
```

### Command Line Flags

#### 🚩 Daemon Flags

```bash
# Start daemon with custom configuration
servin daemon \
  --config /etc/servin/config.yaml \
  --data-root /var/lib/servin \
  --log-level debug \
  --log-file /var/log/servin/servin.log \
  --cri \
  --cri-port 10010 \
  --api \
  --api-port 8080 \
  --metrics \
  --metrics-port 9090

# Override specific settings
servin daemon \
  --runtime crun \
  --storage-driver overlay2 \
  --network-bridge servin0 \
  --registry-mirror https://mirror.example.com

# Security settings
servin daemon \
  --tls \
  --tls-cert /etc/servin/tls/server.crt \
  --tls-key /etc/servin/tls/server.key \
  --auth-method token \
  --auth-file /etc/servin/tokens.json
```

### Registry Configuration

#### 🏪 Registry Settings

```yaml
# Registry configuration file
# ~/.config/servin/registries.yaml

registries:
  # Docker Hub configuration
  "docker.io":
    mirrors:
      - "https://mirror.gcr.io"
      - "https://mirror.example.com"
    
  # Private registry
  "registry.company.com":
    tls:
      insecure: false
      ca_file: "/etc/ssl/certs/company-ca.crt"
    auth:
      username: "user"
      password: "secret"
      
  # Insecure registry
  "registry.internal:5000":
    tls:
      insecure: true
```

#### 🔑 Authentication

```bash
# Login to registries
servin login docker.io
servin login registry.company.com
```

```json
# Configure authentication file
# ~/.config/servin/auth.json
{
  "auths": {
    "docker.io": {
      "username": "user",
      "password": "pass",
      "email": "user@example.com",
      "auth": "dXNlcjpwYXNz"
    },
    "registry.company.com": {
      "username": "employee",
      "password": "secret",
      "auth": "ZW1wbG95ZWU6c2VjcmV0"
    }
  }
}
```

### Network Configuration

#### 🌐 Network Settings

**Bridge Network:**
```yaml
# Default bridge configuration
network:
  bridge:
    name: "servin0"
    subnet: "172.17.0.0/16"
    gateway: "172.17.0.1"
    ip_masq: true
    icc: true
    mtu: 1500
```

```bash
# Custom bridge
servin network create \
  --driver bridge \
  --subnet 192.168.100.0/24 \
  --gateway 192.168.100.1 \
  custom-network
```

**Advanced Network Options:**
```yaml
# Network configuration
network:
  # Enable IPv6
  ipv6: true
  
  # Fixed CIDR
  fixed_cidr: "172.17.0.0/16"
  
  # Default address pools
  default_address_pools:
    - base: "172.80.0.0/16"
      size: 24
    - base: "172.90.0.0/16"
      size: 24
      
  # DNS settings
  dns:
    - "8.8.8.8"
    - "8.8.4.4"
```

### Performance Tuning

#### ⚡ Performance Settings

```yaml
# Performance configuration
daemon:
  # Concurrent operations
  max_concurrent_downloads: 10
  max_concurrent_uploads: 5
  
  # Container limits
  default_ulimits:
    - name: "nofile"
      soft: 65536
      hard: 65536
```

---

## Quick Start

### First Steps

1. **Install Servin** following the [Installation](#installation) guide
2. **Start the daemon:**
   ```bash
   servin daemon
   ```
3. **Verify installation:**
   ```bash
   servin version
   ```

### Basic Container Operations

```bash
# Pull an image
servin pull hello-world

# Run a container
servin run hello-world

# List containers
servin ps

# List images
servin images

# Stop a container
servin stop <container-id>

# Remove a container
servin rm <container-id>
```

### Using the GUI

```bash
# Launch desktop GUI
servin-gui

# Launch terminal UI
servin-desktop
```

---

## Command Line (CLI)

*[This section would contain comprehensive CLI documentation]*

## Terminal UI (TUI)

*[This section would contain TUI interface documentation]*

## Desktop GUI

*[This section would contain GUI application documentation]*

## Container Management

*[This section would contain detailed container management documentation]*

## Image Management

*[This section would contain detailed image management documentation]*

## Volume Management

*[This section would contain detailed volume management documentation]*

## Registry Operations

*[This section would contain detailed registry operations documentation]*

## Kubernetes CRI

*[This section would contain detailed CRI integration documentation]*

## API Reference

*[This section would contain detailed API reference documentation]*

## Logging & Monitoring

*[This section would contain detailed logging and monitoring documentation]*

## Building from Source

*[This section would contain detailed build instructions]*

## Development Guide

*[This section would contain detailed development guide]*

## Contributing

*[This section would contain contribution guidelines]*

## Troubleshooting

*[This section would contain troubleshooting information]*

## FAQ

*[This section would contain frequently asked questions]*

## Resources

*[This section would contain additional resources and links]*

---

**© 2025 Servin Container Runtime. All rights reserved.**
