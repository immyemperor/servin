# Servin - Cross-Platform Container Runtime

## Overview

Servin is a lightweight container runtime built from scratch in Go that provides comprehensive containerization capabilities with cross-platform support for development and production environments. It includes advanced security features such as user namespaces, rootless containers, and local registry support.

## 📚 Documentation

**Complete documentation is available at: [https://immyemperor.github.io/servin](https://immyemperor.github.io/servin)**

The documentation includes:
- 🛠️ **Installation guides** for Windows, Linux, and macOS
- 🖥️ **User interface documentation** (CLI, TUI, Desktop GUI, WebView)
- 📖 **API reference** and developer guides
- 🏗️ **Architecture overview** and technical details
- 🔧 **Configuration** and troubleshooting guides

## Platform Support

### 🐧 Linux (Production Ready)
- **Full containerization support** with namespaces (PID, UTS, IPC, NET, Mount, User)
- **Enhanced security isolation** with user namespaces and UID/GID mapping
- **Complete cgroups v1 integration** for resource management
- **Advanced networking** with bridge networks, veth pairs, and IPAM
- **Root filesystem isolation** using chroot
- **System-wide installation** in `/var/lib/servin`

### 🪟 Windows (Development Environment)
- **Cross-platform container simulation** for development workflows
- **Image management system** with full import/export capabilities
- **Container state persistence** in user home directory (`cd `)
- **Development-friendly** with automatic dev mode enabling
- **CLI compatibility** with all commands available

### 🍎 macOS (Development Environment)
- **Unix-compatible development environment** with proper privilege handling
- **Complete image management** matching Linux functionality
- **Homebrew-style user directory storage** (`~/.servin`)
- **Sudo requirement handling** with `--dev` flag bypass option
- **Native macOS path conventions**

## Core Features

### Container Management
- **Lifecycle management**: create, start, stop, remove containers
- **Process isolation**: namespace-based process separation (Linux)
- **Resource limits**: CPU and memory constraints via cgroups (Linux)
- **Enhanced security**: User namespaces with UID/GID mapping for privilege isolation
- **State persistence**: JSON-based container state tracking across platforms
- **Log capture**: Automatic stdout/stderr capture with timestamp support

### Image System
- **Tarball import/export**: Full image packaging system
- **Metadata management**: JSON-based image indexing
- **Multi-tag support**: Repository-style image tagging and re-tagging
- **Cross-platform storage**: Platform-appropriate paths for all systems
- **Buildfile support**: Docker-compatible image building from Buildfiles

### Networking (Linux Only)
- **Bridge networking**: Default `servin0` bridge with IPAM
- **Container connectivity**: veth pair creation and management
- **Port mapping**: Host-to-container port forwarding
- **NAT integration**: iptables-based network address translation

### Volume Management
- **Named volumes**: Create and manage persistent data volumes
- **Cross-platform storage**: Platform-appropriate volume paths
- **Label support**: Metadata and organization capabilities
- **Lifecycle management**: Create, list, inspect, and remove volumes
- **Bulk operations**: Remove all volumes with confirmation

### Security and Isolation
- **User namespaces**: Enhanced security with UID/GID mapping and privilege isolation
- **Rootless containers**: Run containers without root privileges for improved security
- **Capability management**: Fine-grained capability control and privilege dropping
- **Security policies**: Configurable security constraints and no-new-privileges support
- **Security validation**: Built-in security feature detection and testing

### Compose Orchestration
- **Multi-service applications**: Define and run applications with multiple containers
- **YAML configuration**: Docker Compose-compatible `servin-compose.yml` format
- **Service dependencies**: Automatic dependency resolution and startup ordering
- **Shared networking**: Services can communicate using service names
- **Volume sharing**: Named volumes and bind mounts across services
- **Environment management**: Service-specific environment variables and configuration
- **Lifecycle management**: Up/down operations for entire application stacks

### Container Runtime Interface (CRI)
- **Kubernetes compatibility**: HTTP-based CRI server implementing Kubernetes CRI specification v1alpha2
- **Pod sandbox operations**: Create, list, start, stop, and remove pod sandboxes for Kubernetes pod management
- **Container lifecycle**: Full CRI-compatible container operations including create, start, stop, remove, and status
- **Image service**: Comprehensive image operations with list, pull, remove, and status functionality
- **Health monitoring**: Built-in health endpoints for monitoring CRI server status
- **RESTful API**: HTTP endpoints at `/v1/runtime/` and `/v1/image/` for easy integration
- **Development friendly**: Simple HTTP-based approach instead of complex gRPC for easier testing and debugging

#### CRI Usage
```bash
# Start CRI server
servin cri start --port 8080

# Available endpoints:
# GET  /health                     - Health check
# POST /v1/runtime/version         - Get runtime version
# GET  /v1/runtime/status          - Get runtime status  
# POST /v1/runtime/sandbox/list    - List pod sandboxes
# POST /v1/runtime/sandbox/create  - Create pod sandbox
# POST /v1/runtime/sandbox/start   - Start pod sandbox
# POST /v1/runtime/sandbox/stop    - Stop pod sandbox
# POST /v1/runtime/sandbox/remove  - Remove pod sandbox
# POST /v1/runtime/container/list  - List containers
# POST /v1/runtime/container/create - Create container
# POST /v1/runtime/container/start - Start container
# POST /v1/runtime/container/stop  - Stop container
# POST /v1/runtime/container/remove - Remove container
# POST /v1/runtime/container/status - Get container status
# POST /v1/image/list              - List images
# POST /v1/image/status            - Get image status
# POST /v1/image/pull              - Pull image
# POST /v1/image/remove            - Remove image
# POST /v1/image/fs                - Get filesystem usage
```

### Desktop Interface
- **Terminal User Interface (TUI)**: Full-featured text-based interface for all Servin operations
- **Interactive menus**: Easy navigation through container, image, CRI, volume, and registry management
- **Real-time feedback**: Command output and status updates displayed directly in terminal
- **Cross-platform**: Works on Windows, Linux, and macOS with no additional dependencies
- **Docker Desktop-like experience**: Familiar interface for container management workflows
- **Future GUI support**: Planned graphical interface using Go-native frameworks

### Kubernetes Integration
- **Container Runtime Interface (CRI)**: HTTP-based CRI server for Kubernetes compatibility
- **Pod sandbox management**: Create, list, and remove pod sandboxes
- **Container lifecycle**: Full Kubernetes-compatible container operations
- **Image service**: List, pull, remove, and status operations for container images
- **Health monitoring**: Built-in health checks and status endpoints
- **RESTful API**: HTTP endpoints matching CRI specification

### Logging and Error Handling
- **Structured logging**: Multi-level logging with file and console output
- **Rich error context**: Categorized errors with contextual information
- **Debug support**: Verbose mode with caller information and stack traces
- **Cross-platform logs**: Platform-specific log file locations
- **Operational monitoring**: Comprehensive audit trail and troubleshooting

### Command Line Interface
```bash
# Container operations
servin run [--name NAME] IMAGE COMMAND [ARGS...]
servin ls                    # List containers
servin stop CONTAINER_ID     # Stop running container
servin rm CONTAINER_ID       # Remove container
servin exec CONTAINER_ID CMD # Execute command in container
servin logs CONTAINER_ID     # Fetch logs from container

# Image operations
servin image ls              # List images
servin image import FILE     # Import tarball as image
servin image rm IMAGE        # Remove image
servin image inspect IMAGE   # Inspect image details
servin image tag SOURCE TARGET # Tag an image with a new name
servin build PATH            # Build image from Buildfile

# CRI operations - Kubernetes Container Runtime Interface
servin cri start             # Start CRI HTTP server on port 8080
servin cri start --port 9090 # Start CRI server on custom port
servin cri status            # Check CRI server status
servin cri test              # Test CRI server connectivity

# GUI operations - Desktop interface
servin gui                   # Launch Servin Desktop GUI
servin gui --tui             # Launch Terminal User Interface
servin gui --dev             # Launch in development mode
servin gui --port 8081       # Launch GUI on custom port

# Network operations (Linux only)
servin network ls            # List networks
servin network create NAME   # Create network
servin network rm NAME       # Remove network

# Volume operations
servin volume ls             # List volumes
servin volume create NAME    # Create volume
servin volume rm VOLUME      # Remove volume
servin volume rm-all         # Remove all volumes
servin volume inspect VOLUME # Inspect volume details
servin volume prune          # Remove unused volumes

# Logs operations
servin logs CONTAINER        # Show container logs
servin logs -f CONTAINER     # Follow logs in real-time
servin logs -t CONTAINER     # Show logs with timestamps
servin logs --tail 10 CONTAINER  # Show last 10 lines
servin logs --since 1h CONTAINER # Show logs from last hour

# Build operations
servin build .               # Build image from Buildfile in current directory
servin build -t myapp:v1.0 . # Build and tag image
servin build -f MyBuildfile . # Build with custom Buildfile name
servin build --build-arg VERSION=1.0 . # Build with arguments
servin build -q .            # Quiet build (only show image ID)

# Compose operations
servin compose up            # Create and start services from servin-compose.yml
servin compose up -d         # Start services in detached mode
servin compose down          # Stop and remove services
servin compose down --volumes # Stop services and remove volumes
servin compose ps            # List running services
servin compose ps -a         # List all services (including stopped)
servin compose logs          # Show logs from all services
servin compose logs web      # Show logs from specific service
servin compose logs -f web   # Follow logs from specific service
servin compose exec web sh   # Execute command in running service
servin compose -f custom-compose.yml up # Use custom compose file
servin compose -p myproject up # Specify project name

# Registry operations
servin registry start        # Start local registry server on port 5000
servin registry start --port 5001 # Start on custom port
servin registry start --detach # Start in background (planned)
servin registry stop         # Stop local registry server
servin registry push myapp:latest # Push image to default registry
servin registry push myapp:v1.0 localhost:5001 # Push to specific registry
servin registry pull nginx:alpine # Pull image from default registry
servin registry pull myapp:latest localhost:5001 # Pull from specific registry
servin registry login docker.io # Authenticate with registry
servin registry logout docker.io # Remove authentication
servin registry list         # List configured registries and status

# Security operations
servin security check        # Check security feature availability
servin security info         # Display current security configuration
servin security config --user-ns --uid-map "0:1000:1" # Configure user namespace mapping
servin security config --rootless # Enable rootless container mode
servin security config --no-new-privs # Enable no-new-privileges policy
servin security test         # Test security isolation and namespace functionality
servin security test --user-ns # Test specific security feature

# Global logging and debugging flags
servin --verbose COMMAND     # Enable verbose output
servin --log-level debug COMMAND  # Set log level (debug, info, warn, error)
servin --log-file PATH COMMAND     # Specify custom log file
```

## Platform-Specific Behavior

### Storage Locations
| Platform | Container State | Images | Volumes | Networks | Logs | Registry |
|----------|----------------|--------|---------|----------|------|----------|
| Linux | `/var/lib/servin/containers` | `/var/lib/servin/images` | `/var/lib/servin/volumes` | `/var/lib/servin/networks` | `/var/lib/servin/logs` | `/var/lib/servin/registry` |
| Windows | `%USERPROFILE%\.servin\containers` | `%USERPROFILE%\.servin\images` | `%USERPROFILE%\.servin\volumes` | N/A | `%USERPROFILE%\.servin\logs` | `%USERPROFILE%\.servin\registry` |
| macOS | `~/.servin/containers` | `~/.servin/images` | `~/.servin/volumes` | N/A | `~/.servin/logs` | `~/.servin/registry` |

### Root Privileges
- **Linux**: Required for namespace and cgroup operations
- **Windows**: Development mode enabled by default
- **macOS**: Required but can be bypassed with `--dev` flag

### Feature Matrix
| Feature | Linux | Windows | macOS |
|---------|-------|---------|-------|
| Namespaces | ✅ | ❌ (simulated) | ❌ (simulated) |
| Cgroups | ✅ | ❌ | ❌ |
| Networking | ✅ | ❌ | ❌ |
| User Namespaces | ✅ | ❌ | ❌ |
| Security Isolation | ✅ | ⚠️ (basic) | ⚠️ (basic) |
| Rootless Containers | ✅ | ❌ | ❌ |
| Container Management | ✅ | ✅ | ✅ |
| Image Management | ✅ | ✅ | ✅ |
| Image Building | ✅ | ✅ | ✅ |
| Volume Management | ✅ | ✅ | ✅ |
| Compose Orchestration | ✅ | ✅ | ✅ |
| Local Registry | ✅ | ✅ | ✅ |
| Registry Push/Pull | ✅ | ✅ | ✅ |
| CRI Compatibility | ✅ | ✅ | ✅ |
| Desktop Interface | ✅ | ✅ | ✅ |
| State Persistence | ✅ | ✅ | ✅ |
| Log Capture | ✅ | ⚠️ (limited) | ⚠️ (limited) |
| Container Simulation | ✅ | ✅ | ✅ |

## Development Workflow

### Cross-Platform Development
1. **Develop on any platform** using full image and state management
2. **Test container logic** with simulated environments
3. **Deploy to Linux** for production containerization

### Build Instructions
```bash
# Build for current platform
go build -o servin .

# Cross-compile for different platforms
GOOS=linux go build -o servin-linux .
GOOS=windows go build -o servin-windows.exe .
GOOS=darwin go build -o servin-macos .
```

### Testing
```bash
# Run comprehensive cross-platform test
go run test-platform.go

# Test basic functionality
./servin run alpine echo "Hello World"
./servin image ls
./servin ls
```

## Architecture

### Core Packages
- **`cmd/`**: CLI command implementations using Cobra framework
- **`pkg/container/`**: Container lifecycle and process management
- **`pkg/image/`**: Image storage, import/export, and metadata management
- **`pkg/state/`**: Container state persistence and retrieval
- **`pkg/rootfs/`**: Root filesystem creation and management
- **`pkg/network/`**: Networking stack with bridge and veth support
- **`pkg/namespaces/`**: Linux namespace creation and management
- **`pkg/cgroups/`**: Resource limitation and monitoring
- **`pkg/cri/`**: Container Runtime Interface (CRI) server and Kubernetes integration

### Build Constraints
```go
//go:build linux
// Full implementation for Linux

//go:build !linux  
// Cross-platform stubs for Windows/macOS
```

## Installation & Usage

> **📖 For detailed installation instructions with professional installers, see the [Installation Guide](https://immyemperor.github.io/servin/installation/)**

### Prerequisites
- **Go 1.21+** for building from source
- **Linux kernel 3.8+** for full containerization features
- **Root privileges** for production Linux deployment

### Quick Start
```bash
# Clone and build
git clone <repository>
cd servin
go build -o servin .

# Import an image
servin image import alpine.tar

# Run a container
servin run alpine echo "Hello from Servin!"

# List containers
servin ls

# Clean up
servin rm <container_id>
```

## Limitations & Future Enhancements

### Current Limitations
- **Windows/macOS**: No true containerization (development simulation only)
- **Networking**: Linux-only bridge networking
- **Remote registries**: Docker Hub and other remote registries (implementation in progress)

### Planned Enhancements
- **Complete remote registry** support for Docker Hub and other registries
- **Windows Containers** integration
- **macOS containers** via hypervisor framework

## Conclusion

Servin provides a complete foundation for understanding and working with container technologies while offering practical cross-platform development capabilities. It bridges the gap between learning containerization concepts and building production-ready solutions.

## 📚 Learn More

- **📖 Full Documentation**: [https://immyemperor.github.io/servin](https://immyemperor.github.io/servin)
- **🛠️ Installation Guide**: [Installation Instructions](https://immyemperor.github.io/servin/installation/)
- **🖥️ User Interfaces**: [CLI](https://immyemperor.github.io/servin/cli/), [TUI](https://immyemperor.github.io/servin/tui/), [GUI](https://immyemperor.github.io/servin/gui/)
- **🏗️ Architecture**: [Technical Overview](https://immyemperor.github.io/servin/architecture/)
- **🔧 Configuration**: [Setup and Configuration](https://immyemperor.github.io/servin/configuration/)
- **❓ Troubleshooting**: [Common Issues](https://immyemperor.github.io/servin/troubleshooting/)
