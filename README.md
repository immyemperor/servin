# Servin - Open-Source Cross-Platform Container Runtime

## Overview

Servin is a **free and open-source** lightweight container runtime built from scratch in Go that provides comprehensive containerization capabilities with cross-platform support for development and production environments. As an open-source project, Servin includes advanced security features such as user namespaces, rootless containers, and local registry support, with full transparency and community-driven development.

**🔓 Open Source License**: Servin is released under an open-source license, allowing free use, modification, and distribution for both personal and commercial projects.

## � Quick Installation

### Download from GitHub Releases

Get the latest release from: **[GitHub Releases](https://github.com/immyemperor/servin/releases/latest)**

#### 🍎 macOS
```bash
# Download and install via DMG (Recommended)
# 1. Download Servin-Container-Runtime.dmg from releases
# 2. Double-click to mount
# 3. Drag Servin.app to Applications

# Or use the installer wizard
curl -O https://github.com/immyemperor/servin/releases/latest/download/servin-macos-universal.tar.gz
tar -xzf servin-macos-universal.tar.gz
cd servin-macos-universal
sudo ./ServinInstaller.command
```

#### 🐧 Linux
```bash
# Download and extract
wget https://github.com/immyemperor/servin/releases/latest/download/servin-linux-amd64.tar.gz
tar -xzf servin-linux-amd64.tar.gz
cd servin-linux-amd64

# Run installer wizard
sudo ./ServinInstaller.sh
```

#### 🪟 Windows
```powershell
# Download servin-windows-amd64.zip from releases
# Extract and run installer
.\ServinSetup.exe
```

## �📚 Documentation

**Complete documentation is available at: [https://immyemperor.github.io/servin](https://immyemperor.github.io/servin)**

The documentation includes:
- 🛠️ **Installation guides** for Windows, Linux, and macOS
- 🖥️ **User interface documentation** (CLI, TUI, Desktop GUI)
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

## 🎯 Core Features

### 🛡️ Security First
- **User Namespaces**: Complete user and group ID mapping
- **Rootless Containers**: Run containers without root privileges
- **PID Isolation**: Process ID namespace isolation
- **Network Isolation**: Dedicated network namespaces per container
- **cgroup v2 Integration**: Advanced resource management and limits

### 🏗️ Container Management
- **Lifecycle Control**: Create, start, stop, pause, resume containers
- **Multi-format Support**: OCI, Docker images, and custom rootfs
- **Volume Management**: Bind mounts and named volumes
- **Network Management**: Custom networking with IPAM
- **Compose Support**: Multi-container application orchestration

### 🖥️ Multiple User Interfaces
- **CLI**: Command-line interface for automation and scripting
- **TUI**: Text-based user interface for interactive management
- **Desktop GUI**: Native desktop application (Flask + pywebview, distributed as binary)
- **CRI Support**: Kubernetes Container Runtime Interface compatible

### 🚀 Recent Improvements
- **Binary Distribution**: Desktop GUI now ships as a compiled binary for better performance
- **Professional Installers**: Cross-platform installer wizards with proper privilege handling
- **macOS .dmg Support**: Professional disk image distribution for macOS
- **Enhanced Security**: Improved privilege escalation and user consent flows
- **Better Error Handling**: Comprehensive timeout mechanisms and UI responsiveness

### Desktop Interface
- **Terminal User Interface (TUI)**: Full-featured text-based interface for all Servin operations
- **Interactive menus**: Easy navigation through container, image, CRI, volume, and registry management
- **Real-time feedback**: Command output and status updates displayed directly in terminal
- **Cross-platform**: Works on Windows, Linux, and macOS with no additional dependencies
- **Docker Desktop-like experience**: Familiar interface for container management workflows

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
| **Core Runtime** | | | |
| Namespaces | ✅ | ❌ (simulated) | ❌ (simulated) |
| Cgroups | ✅ | ❌ | ❌ |
| User Namespaces | ✅ | ❌ | ❌ |
| Rootless Containers | ✅ | ❌ | ❌ |
| Security Isolation | ✅ | ⚠️ (basic) | ⚠️ (basic) |
| Container Management | ✅ | ✅ | ✅ |
| Container Simulation | ✅ | ✅ | ✅ |
| **Images & Registry** | | | |
| Image Management | ✅ | ✅ | ✅ |
| Image Building | ✅ | ✅ | ✅ |
| Multi-Architecture | ✅ | ✅ | ✅ |
| Local Registry | ✅ | ✅ | ✅ |
| Registry Push/Pull | ✅ | ✅ | ✅ |
| Image Security Scan | ✅ | ✅ | ✅ |
| **Storage & Network** | | | |
| Volume Management | ✅ | ✅ | ✅ |
| Bridge Networking | ✅ | ❌ | ❌ |
| Port Management | ✅ | ✅ | ✅ |
| Network Isolation | ✅ | ❌ | ❌ |
| **Orchestration** | | | |
| Compose Orchestration | ✅ | ✅ | ✅ |
| Multi-Container Apps | ✅ | ✅ | ✅ |
| Service Discovery | ✅ | ✅ | ✅ |
| **Kubernetes Integration** | | | |
| CRI v1alpha2 | ✅ | ✅ | ✅ |
| Pod Sandbox Management | ✅ | ✅ | ✅ |
| gRPC API Server | ✅ | ✅ | ✅ |
| Kubelet Integration | ✅ | ✅ | ✅ |
| **VM Engine** | | | |
| VM Management | ✅ | ✅ | ✅ |
| VM Status Monitoring | ✅ | ✅ | ✅ |
| Cross-Platform VMs | ✅ | ✅ | ✅ |
| VM Configuration | ✅ | ✅ | ✅ |
| **User Interfaces** | | | |
| CLI Interface | ✅ | ✅ | ✅ |
| Terminal UI (TUI) | ✅ | ✅ | ✅ |
| Desktop GUI | ✅ | ✅ | ✅ |
| WebView Interface | ✅ | ✅ | ✅ |
| **Security Features** | | | |
| Capability Management | ✅ | ❌ | ❌ |
| Security Policies | ✅ | ⚠️ (basic) | ⚠️ (basic) |
| Security Testing | ✅ | ✅ | ✅ |
| Privilege Dropping | ✅ | ❌ | ❌ |
| **Monitoring & Logging** | | | |
| Container Logs | ✅ | ✅ | ✅ |
| Log Streaming | ✅ | ⚠️ (limited) | ⚠️ (limited) |
| Health Checks | ✅ | ✅ | ✅ |
| Metrics Export | ✅ | ✅ | ✅ |
| Performance Monitoring | ✅ | ✅ | ✅ |
| Prometheus Integration | ✅ | ✅ | ✅ |
| **Development & DevOps** | | | |
| State Persistence | ✅ | ✅ | ✅ |
| Development Mode | ✅ | ✅ | ✅ |
| Cross-Platform Testing | ✅ | ✅ | ✅ |
| Professional Installers | ✅ | ✅ | ✅ |

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

### Recommended: Download from Releases
**Get the latest release from: [GitHub Releases](https://github.com/immyemperor/servin/releases/latest)**

Pre-built binaries are available for:
- **macOS**: Universal binary + professional .dmg installer
- **Linux**: AMD64 binary + installer wizard
- **Windows**: AMD64 binary + setup wizard

### Building from Source (Development)
For development or custom builds:

#### Prerequisites
- **Go 1.21+** for building from source
- **Python 3.8+** for desktop GUI development
- **Linux kernel 3.8+** for full containerization features
- **Root privileges** for production Linux deployment

#### Quick Start
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

## � Open Source & Community

### **Why Open Source?**
Servin is committed to open-source principles, providing:
- **🔍 Full Transparency**: Complete source code visibility and audit capability
- **🤝 Community-Driven**: Development guided by community needs and contributions
- **📚 Educational Value**: Learn containerization by studying real implementation
- **🔒 No Vendor Lock-in**: Freedom to modify, extend, and distribute
- **🆓 Always Free**: No licensing fees, premium tiers, or usage restrictions

### **Contributing to Servin**
We welcome contributions from developers of all skill levels:
- **🐛 Bug Reports**: Help improve stability and reliability
- **💡 Feature Requests**: Suggest new capabilities and enhancements
- **📝 Documentation**: Improve guides, examples, and explanations
- **💻 Code Contributions**: Implement features, fix bugs, optimize performance
- **🧪 Testing**: Cross-platform testing and validation
- **🌐 Translations**: Help make Servin accessible globally

### **Repository & Development**
- **📦 Source Code**: [https://github.com/immyemperor/servin](https://github.com/immyemperor/servin)
- **🐛 Issue Tracker**: Report bugs and request features on GitHub
- **📋 Project Board**: Track development progress and roadmap
- **🔄 Pull Requests**: Contribute code improvements and new features
- **📞 Discussions**: Join community discussions and ask questions

### **License & Usage**
Servin is released under an open-source license that permits:
- ✅ **Personal Use**: Free for individual developers and personal projects
- ✅ **Commercial Use**: No restrictions for business and enterprise usage
- ✅ **Modification**: Adapt and customize for specific needs
- ✅ **Distribution**: Share modified versions with the community

## �📚 Learn More

- **📖 Full Documentation**: [https://immyemperor.github.io/servin](https://immyemperor.github.io/servin)
- **🛠️ Installation Guide**: [Installation Instructions](https://immyemperor.github.io/servin/installation/)
- **🖥️ User Interfaces**: [CLI](https://immyemperor.github.io/servin/cli/), [TUI](https://immyemperor.github.io/servin/tui/), [GUI](https://immyemperor.github.io/servin/gui/)
- **🏗️ Architecture**: [Technical Overview](https://immyemperor.github.io/servin/architecture/)
- **🔧 Configuration**: [Setup and Configuration](https://immyemperor.github.io/servin/configuration/)
- **❓ Troubleshooting**: [Common Issues](https://immyemperor.github.io/servin/troubleshooting/)
