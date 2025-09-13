---
layout: default
title: Architecture
permalink: /architecture/
---

# 🏗 Architecture

## System Overview

Servin Container Runtime follows a modular architecture design that separates concerns while maintaining high performance and reliability.

<div class="architecture-diagram">
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
</div>

## Component Architecture

### 🎯 Core Components

#### Runtime Engine
- **Container Lifecycle Management**: Create, start, stop, pause, delete containers
- **Process Management**: Handle container processes and signal forwarding
- **Resource Management**: CPU, memory, and I/O resource allocation
- **Security**: Namespace isolation, capability management, SELinux/AppArmor

#### Image Manager
- **OCI Image Support**: Full OCI image specification compliance
- **Layer Management**: Efficient layer storage and deduplication
- **Image Operations**: Pull, push, build, tag, inspect operations
- **Multi-architecture**: Support for different CPU architectures

#### Volume Manager
- **Persistent Storage**: Named volumes with lifecycle management
- **Bind Mounts**: Host directory mounting with proper permissions
- **Tmpfs Mounts**: In-memory temporary storage
- **Storage Drivers**: Pluggable storage backend support

#### Network Manager
- **Bridge Networks**: Default container networking
- **Custom Networks**: User-defined networks with isolation
- **Port Management**: Port forwarding and publishing
- **DNS Resolution**: Container name resolution

#### Registry Client
- **Authentication**: Registry login and credential management
- **Push/Pull Operations**: Efficient image transfer
- **Manifest Handling**: Image manifest processing
- **Mirror Support**: Registry mirror configuration

### 🔌 Interface Layer

#### CLI Interface
```bash
servin run alpine:latest
servin ps
servin images
servin networks ls
```

#### Terminal UI (TUI)
- Interactive menu-driven interface
- Real-time container monitoring
- Visual network and volume management
- Cross-platform terminal support

#### Desktop GUI
- Native desktop application using Fyne framework
- Visual container management
- Real-time logs and monitoring
- Professional user experience

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

## Directory Structure

```
servin/
├── cmd/                          # Application entry points
│   ├── servin-desktop/          # Terminal UI application
│   ├── servin-gui/              # Desktop GUI application
│   └── gui.go                   # GUI command integration
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
│   ├── windows/              # NSIS-based Windows installer
│   ├── linux/                # Python/Tkinter Linux installer
│   └── macos/                # Python/Tkinter macOS installer
├── build/                    # Build artifacts
├── dist/                     # Distribution packages
└── docs/                     # Documentation
```

## Next Steps

- [Features Overview]({{ '/features' | relative_url }}) - Detailed feature descriptions
- [Installation Guide]({{ '/installation' | relative_url }}) - Get Servin running
- [Configuration]({{ '/configuration' | relative_url }}) - Configure Servin for your environment

[View Features →]({{ '/features' | relative_url }}){: .btn .btn-primary}
[Install Now →]({{ '/installation' | relative_url }}){: .btn .btn-outline}
