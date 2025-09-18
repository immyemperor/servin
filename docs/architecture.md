---
layout: default
title: Architecture
permalink: /architecture/
---

# 🏗 Architecture

## Revolutionary Dual-Mode Architecture

Servin Container Runtime features a **revolutionary dual-mode architecture** that provides both native Linux containerization and universal VM-based containerization across all platforms.

### 🎯 **Containerization Modes**

1. **Native Mode** (Linux): Direct kernel integration for maximum performance
2. **VM Mode** (Universal): Linux VM providing true containerization on any platform

<div class="architecture-diagram">
┌─────────────────────────────────────────────────────────────┐
│                    Servin Container Runtime                  │
├─────────────────────────────────────────────────────────────┤
│  Dual-Mode Engine                                           │
│  ┌─────────────────────┐ ┌─────────────────────────────────┐│
│  │    Native Mode      │ │        VM Mode                  ││
│  │   (Linux Only)      │ │   (Windows/macOS/Linux)         ││
│  │                     │ │                                 ││
│  │ ┌─────────────────┐ │ │ ┌─────────────────────────────┐ ││
│  │ │ Direct Kernel   │ │ │ │    Linux VM Container       │ ││
│  │ │ Namespaces     │ │ │ │      Engine                │ ││
│  │ │ + cgroups      │ │ │ │                             │ ││
│  │ └─────────────────┘ │ │ │ ┌─────────────────────────┐ │ ││
│  └─────────────────────┘ │ │ │  KVM/Hyper-V/VMware    │ │ ││
│                          │ │ │  Virtualization.framework│ │ ││
│                          │ │ └─────────────────────────┘ │ ││
│                          │ └─────────────────────────────┘ ││
│                          └─────────────────────────────────┘│
├─────────────────────────────────────────────────────────────┤
│  User Interfaces (Common Across Both Modes)                 │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────────┐│
│  │    CLI      │ │     TUI     │ │      Desktop GUI        ││
│  │  Command    │ │  Terminal   │ │   Flask + pywebview     ││
│  │   Line      │ │ Interface   │ │   Binary Distribution   ││
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
│  Platform Integration Layer                                 │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────────┐│
│  │   Linux     │ │   Windows   │ │        macOS            ││
│  │   Native    │ │  VM Mode    │ │      VM Mode            ││
│  │   + VM      │ │  (Hyper-V)  │ │ (Virtualization.fwk)    ││
│  └─────────────┘ └─────────────┘ └─────────────────────────┘│
└─────────────────────────────────────────────────────────────┘
</div>

## Platform-Specific Implementation

### 🐧 **Linux**: Native + VM Modes
- **Native Mode** (Default): Direct kernel namespaces, cgroups, capabilities
- **VM Mode** (Optional): KVM/QEMU-based Linux VM for enhanced isolation
- **Automatic Selection**: Native preferred, VM available for security

### 🪟 **Windows**: VM Mode Only  
- **VM Engine**: Hyper-V or WSL2-based Linux VM
- **Automatic**: VM mode initializes seamlessly on first run
- **Integration**: Native Windows GUI with Linux container engine

### 🍎 **macOS**: VM Mode Only
- **VM Engine**: Virtualization.framework-based Linux VM
- **Universal**: Intel and Apple Silicon support
- **Seamless**: Native macOS experience with Linux containers

## Component Architecture

### 🎯 Dual-Mode Runtime Engine

#### Native Mode Engine (Linux)
- **Direct Kernel Access**: Linux namespaces (PID, Network, Mount, UTS, IPC, User)
- **Resource Control**: cgroups v1/v2 for CPU, memory, I/O limits
- **Security**: Capabilities, SELinux/AppArmor integration
- **Performance**: Zero virtualization overhead
- **Compatibility**: Full Docker and OCI compatibility

#### VM Mode Engine (Universal)
- **Linux VM**: Lightweight Linux VM for universal containerization
- **VM Backends**: 
  - **Windows**: Hyper-V, WSL2
  - **macOS**: Virtualization.framework, QEMU
  - **Linux**: KVM/QEMU (optional for enhanced isolation)
- **VM Management**: Automatic VM lifecycle, state persistence
- **Bridge Integration**: Seamless host-VM communication
- **Resource Efficiency**: Optimized VM with minimal overhead

### 🔧 Container Management Core

#### Container Lifecycle Management
- **Create, Start, Stop, Delete**: Full container lifecycle control across both modes
- **Process Management**: Handle container processes and signal forwarding  
- **Resource Management**: CPU, memory, and I/O resource allocation
- **State Persistence**: Container state maintained across VM restarts
- **Security**: Mode-appropriate isolation (namespaces or VM boundaries)

#### Image Manager
- **OCI Image Support**: Full OCI image specification compliance
- **Layer Management**: Efficient layer storage and deduplication
- **Cross-Mode Sharing**: Images work identically in native and VM modes
- **Multi-architecture**: Support for ARM64 and AMD64 architectures
- **Registry Integration**: Pull/push from any OCI-compatible registry

#### Volume Manager
- **Universal Volumes**: Consistent volume behavior across modes
- **Bind Mounts**: Host directory mounting with proper permissions
- **Named Volumes**: Persistent volume creation and management
- **VM Volume Bridge**: Seamless host-VM volume sharing in VM mode
- **Storage Drivers**: Pluggable storage backend support

#### Network Manager
- **Mode-Adaptive Networking**:
  - **Native Mode**: Direct Linux bridge networks, namespaces
  - **VM Mode**: VM-bridged networking with host integration
- **Port Management**: Port forwarding and publishing across VM boundaries
- **DNS Resolution**: Container name resolution in both modes
- **Network Isolation**: Security through network segmentation

#### Registry Client
- **Authentication**: Registry login and credential management
- **Push/Pull Operations**: Efficient image transfer
- **Manifest Handling**: Image manifest processing  
- **Mirror Support**: Registry mirror configuration

### 🔌 Universal Interface Layer

#### CLI Interface (Identical Across Modes)
```bash
# These commands work identically in native and VM modes:
servin run alpine:latest
servin ps
servin images  
servin networks ls
servin vm status     # VM mode specific
servin vm start      # VM mode specific  
```

#### Terminal UI (TUI)
- **Mode-Aware Interface**: Shows current mode (Native/VM)
- **Real-time Monitoring**: Container status regardless of mode
- **VM Management**: VM-specific controls when in VM mode
- **Cross-platform**: Identical experience on all platforms

#### Desktop GUI
- **Universal Web Interface**: Flask backend + pywebview frontend
- **Mode Indicator**: Clear indication of current containerization mode
- **VM Controls**: VM start/stop/status when in VM mode
- **Cross-platform Binary**: PyInstaller distribution for all platforms

### 🌐 API Layer

#### CRI Server (gRPC)
- Full Kubernetes CRI v1alpha2 implementation
- Pod sandbox management
- Container lifecycle operations
- Image service operations

#### HTTP API (REST)
- Docker-compatible REST API
- Authentication and authorization
- Rate limiting and throttling
- OpenAPI documentation

#### Internal gRPC Services
- Inter-component communication
- Service discovery and health checks
- Distributed operations coordination
- Event streaming and notifications

### 💾 Storage Layer

#### Container Storage
- Container filesystem layers
- Read-write layer management
- Snapshot support
- Copy-on-write optimization

#### Image Store
- OCI image storage
- Layer deduplication
- Garbage collection
- Metadata indexing

#### Configuration & Metadata
- YAML/JSON configuration files
- Container metadata database
- Network configuration
- Volume metadata

## Data Flow

### Container Creation Flow

1. **Request Received**: CLI/API receives container creation request
2. **Image Resolution**: Image manager resolves and pulls required image
3. **Network Setup**: Network manager creates container networking
4. **Volume Preparation**: Volume manager prepares storage mounts
5. **Container Creation**: Runtime engine creates container with specified configuration
6. **Process Start**: Container process is started with proper isolation
7. **Status Update**: Container status is updated and stored

### Image Pull Flow

1. **Registry Authentication**: Authenticate with target registry
2. **Manifest Download**: Download image manifest and layer information
3. **Layer Download**: Download missing layers with deduplication
4. **Layer Extraction**: Extract and store layers in storage backend
5. **Image Indexing**: Update image metadata and make available
6. **Cleanup**: Remove temporary files and update cache

## Security Architecture

### Container Isolation
- **Namespaces**: PID, network, mount, user, UTS, IPC isolation
- **Cgroups**: Resource limitation and accounting
- **Capabilities**: Fine-grained privilege control
- **Seccomp**: System call filtering

### Network Security
- **Network Isolation**: Separate network namespaces per container
- **Firewall Integration**: iptables/netfilter rule management
- **TLS Encryption**: Secure registry communication
- **Certificate Management**: PKI infrastructure support

### Storage Security
- **Filesystem Permissions**: Proper file ownership and permissions
- **Encryption Support**: At-rest and in-transit encryption
- **Integrity Checking**: Image and layer verification
- **Access Control**: Role-based access to storage resources

## Performance Considerations

### Optimization Strategies
- **Layer Caching**: Intelligent layer caching and reuse
- **Parallel Operations**: Concurrent image pulls and container operations
- **Memory Management**: Efficient memory usage and garbage collection
- **I/O Optimization**: Optimized filesystem operations

### Scalability Features
- **Horizontal Scaling**: Multiple daemon instances
- **Load Balancing**: Request distribution across instances
- **Resource Pooling**: Shared resource management
- **Async Operations**: Non-blocking operation handling

## Binary Distribution Architecture

### PyInstaller Integration
- **Single-File Executables**: Complete Python runtime embedded in 13MB binary
- **Cross-Platform Support**: Native binaries for Windows, Linux, and macOS
- **No Dependencies**: Self-contained executables require no Python installation
- **Optimized Performance**: Faster startup times compared to Python source execution

### Build System
```bash
# Cross-platform build orchestration
./build-all.sh

# Platform-specific outputs:
# ├── dist/windows/servin-gui.exe     # Windows executable
# ├── dist/linux/servin-gui          # Linux binary  
# └── dist/mac/servin-gui             # macOS universal binary
```

### Professional Distribution
- **macOS .dmg Creation**: Professional disk image with app bundle structure
- **Windows NSIS Installer**: Complete installation wizard with system integration
- **Linux Package Distribution**: Tar.gz archives with installation scripts
- **GitHub Releases Integration**: Automated release creation and artifact upload

### Installation Wizards
- **Cross-Platform Installers**: Python/Tkinter-based wizards for all platforms
- **Privilege Escalation**: Proper sudo/administrator privilege handling
- **Timeout Protection**: Robust subprocess management with comprehensive timeouts
- **User Consent Flows**: Interactive privilege escalation with clear explanations
- **Error Recovery**: Graceful handling of installation failures and cancellations

## Directory Structure

```
servin/
├── cmd/                          # Application entry points
│   ├── servin-desktop/          # Terminal UI application  
│   ├── servin-gui/              # GUI command integration
│   └── gui.go                   # GUI launcher implementation
├── webview_gui/                 # Desktop GUI application
│   ├── main.py                  # PyInstaller entry point
│   ├── app.py                   # Flask backend API
│   ├── servin_client.py         # Servin runtime interface
│   ├── servin-gui.spec          # PyInstaller build specification
│   ├── requirements.txt         # Python dependencies (Flask, pywebview, etc)
│   ├── templates/               # HTML templates for web interface
│   └── static/                  # CSS, JavaScript, and assets
├── pkg/                         # Core packages
│   ├── cri/                     # Container Runtime Interface
│   │   ├── server/             # CRI gRPC server
│   │   ├── sandbox/            # Pod sandbox management
│   │   └── image/              # CRI image service
│   ├── container/              # Container management
│   │   ├── runtime/            # Container runtime operations
│   │   ├── lifecycle/          # Lifecycle management
│   │   └── exec/               # Container exec operations
│   ├── image/                  # Image management
│   │   ├── store/              # Image storage backend
│   │   ├── registry/           # Registry operations
│   │   └── builder/            # Image building
│   ├── volume/                 # Volume management
│   │   ├── drivers/            # Storage drivers
│   │   └── manager/            # Volume lifecycle
│   ├── network/                # Network management
│   │   ├── bridge/             # Bridge driver
│   │   ├── dns/                # DNS resolution
│   │   └── firewall/           # Firewall integration
│   └── registry/               # Registry operations
│       ├── client/             # Registry client
│       ├── auth/               # Authentication
│       └── mirror/             # Registry mirrors
├── internal/                   # Internal packages
│   ├── config/                # Configuration management
│   ├── metrics/               # Metrics collection
│   ├── logging/               # Logging infrastructure
│   └── storage/               # Storage backends
├── api/                       # API definitions
│   ├── http/                  # REST API handlers
│   ├── grpc/                  # gRPC service definitions
│   └── swagger/               # API documentation
├── installers/                # Installation wizards
│   ├── windows/              # NSIS-based Windows installer + Python wizard
│   ├── linux/                # Python/Tkinter Linux installer wizard
│   └── macos/                # Python/Tkinter macOS installer wizard
├── build/                    # Build artifacts and platform binaries
├── dist/                     # PyInstaller distribution packages
│   ├── windows/              # Windows servin-gui.exe
│   ├── linux/                # Linux servin-gui binary
│   └── mac/                  # macOS servin-gui universal binary
└── docs/                     # Documentation
```

## Next Steps

- [Features Overview]({{ '/features' | relative_url }}) - Detailed feature descriptions
- [Installation Guide]({{ '/installation' | relative_url }}) - Get Servin running
- [Configuration]({{ '/configuration' | relative_url }}) - Configure Servin for your environment

[View Features →]({{ '/features' | relative_url }}){: .btn .btn-primary}
[Install Now →]({{ '/installation' | relative_url }}){: .btn .btn-outline}
