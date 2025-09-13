# Servin Container Runtime - Complete Wiki

## Table of Contents
1. [Platform Requirements](#platform-requirements)
2. [Getting Started](#getting-started)
3. [Installation Guide](#installation-guide)
4. [Building from Source](#building-from-source)
5. [Usage Guide](#usage-guide)
6. [Troubleshooting](#troubleshooting)
7. [API Reference](#api-reference)
8. [Contributing](#contributing)

## Platform Requirements

### Windows
**Minimum Requirements:**
- Windows 10/11 (64-bit)
- Go 1.24 or later
- 4GB RAM minimum, 8GB recommended
- 2GB free disk space

**Development Requirements:**
- **CGO Support**: Required for GUI compilation
- **C Compiler**: MinGW-w64 or MSYS2 with UCRT64
- **OpenGL**: For Fyne GUI framework

**Recommended Setup:**
```bash
# Install MSYS2
# Download from: https://www.msys2.org/
# Add to PATH: C:\msys64\ucrt64\bin

# Enable CGO
go env -w CGO_ENABLED=1

# Install GCC toolchain
pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-pkg-config
```

### Linux
**Minimum Requirements:**
- Ubuntu 20.04+ / CentOS 8+ / Debian 11+
- Go 1.24 or later
- 2GB RAM minimum, 4GB recommended
- 1GB free disk space

**Development Requirements:**
```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install build-essential libgl1-mesa-dev xorg-dev

# CentOS/RHEL
sudo yum groupinstall "Development Tools"
sudo yum install mesa-libGL-devel libXrandr-devel libXinerama-devel libXcursor-devel libXi-devel

# Arch Linux
sudo pacman -S base-devel mesa libxrandr libxinerama libxcursor libxi
```

### macOS
**Minimum Requirements:**
- macOS 10.15 (Catalina) or later
- Go 1.24 or later
- Xcode Command Line Tools
- 4GB RAM minimum, 8GB recommended

**Development Requirements:**
```bash
# Install Xcode Command Line Tools
xcode-select --install

# Install Homebrew (optional, for dependencies)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

## Getting Started

### Quick Start (CLI Only)
```bash
# Clone the repository
git clone https://github.com/yourorg/servin
cd servin/gocontainer

# Build CLI
go build -o servin main.go

# Test installation
./servin --version
./servin --help
```

### Quick Start (Full GUI)
```bash
# 1. Set up development environment (see platform requirements above)

# 2. Enable CGO (required for GUI)
go env -w CGO_ENABLED=1

# 3. Install dependencies
go mod tidy

# 4. Build all components
make build-all

# 5. Launch GUI
./servin gui
```

## Installation Guide

### Option 1: Pre-built Binaries (Recommended)
```bash
# Download latest release
curl -L https://github.com/yourorg/servin/releases/latest/download/servin-linux-amd64.tar.gz | tar xz

# Or for Windows
curl -L https://github.com/yourorg/servin/releases/latest/download/servin-windows-amd64.zip -o servin.zip
unzip servin.zip

# Add to PATH
echo 'export PATH=$PATH:/path/to/servin' >> ~/.bashrc
source ~/.bashrc
```

### Option 2: Package Managers
```bash
# Homebrew (macOS/Linux)
brew install servin

# Chocolatey (Windows)
choco install servin

# Snap (Linux)
sudo snap install servin

# APT (Ubuntu/Debian)
curl -fsSL https://packages.servin.io/gpg | sudo apt-key add -
echo "deb https://packages.servin.io/apt stable main" | sudo tee /etc/apt/sources.list.d/servin.list
sudo apt update && sudo apt install servin
```

### Option 3: Docker
```bash
# Run Servin in container
docker run -it --rm servin/servin:latest

# With volume mounting
docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock servin/servin:latest
```

## Building from Source

### Prerequisites Check
```bash
# Verify Go installation
go version  # Should be 1.24+

# Check CGO support (for GUI)
go env CGO_ENABLED  # Should be 1

# Verify C compiler (for GUI)
gcc --version  # Should exist
```

### Build Steps

#### 1. Clone and Setup
```bash
git clone https://github.com/yourorg/servin
cd servin/gocontainer
go mod download
```

#### 2. Build CLI Only
```bash
go build -o servin main.go
./servin --version
```

#### 3. Build TUI (Terminal Interface)
```bash
go build -o servin-desktop cmd/servin-desktop/main.go
./servin-desktop
```

#### 4. Build GUI (Visual Interface)
```bash
# Ensure CGO is enabled
export CGO_ENABLED=1

# Add compiler to PATH (Windows example)
export PATH="/c/msys64/ucrt64/bin:$PATH"

# Build GUI
go build -o servin-gui cmd/servin-gui/main.go cmd/servin-gui/actions.go

# Test GUI
./servin-gui
```

#### 5. Build All Components
```bash
make build-all
# Or manually:
make build build-tui build-gui
```

### Makefile Targets
```bash
make help          # Show all available targets
make build         # Build CLI only
make build-tui     # Build Terminal UI
make build-gui     # Build Visual GUI
make build-all     # Build everything
make test          # Run tests
make clean         # Clean build artifacts
make install       # Install to system
```

## Usage Guide

### Command Line Interface

#### Container Operations
```bash
# Run containers
servin run alpine:latest
servin run --name myapp -p 8080:80 nginx:latest
servin run -it --rm ubuntu:20.04 /bin/bash

# List containers
servin ls
servin ls -a  # Include stopped containers

# Container lifecycle
servin start myapp
servin stop myapp
servin restart myapp
servin rm myapp

# Container information
servin logs myapp
servin exec -it myapp /bin/sh
servin inspect myapp
```

#### Image Operations
```bash
# List images
servin image ls
servin image ls --json

# Import/Export
servin image import myimage.tar
servin image export myimage:latest myimage.tar

# Build images
servin image build -t myapp:v1.0 .
servin image build -f Dockerfile.prod -t myapp:prod .

# Image management
servin image tag myapp:latest myapp:v1.0
servin image rm myapp:old
servin image inspect myapp:latest
```

#### Volume Operations
```bash
# Volume management
servin volume create myvolume
servin volume ls
servin volume inspect myvolume
servin volume rm myvolume
servin volume prune  # Remove unused volumes
```

#### CRI Server (Kubernetes Integration)
```bash
# Start CRI server
servin cri start --port 8080
servin cri start --host 0.0.0.0 --port 8080

# Monitor CRI server
servin cri status
servin cri health

# Stop CRI server
servin cri stop
```

### Graphical User Interface

#### Launching GUI
```bash
# Launch visual GUI
servin gui

# Launch terminal UI
servin gui --tui
servin desktop  # Alternative command

# Development mode
servin gui --dev --port 8081
```

#### GUI Features
- **Container Tab**: Complete lifecycle management
- **Images Tab**: Import, build, tag, remove operations
- **Volumes Tab**: Create, inspect, remove, prune
- **CRI Server Tab**: Start/stop, monitor, test connections
- **Logs Tab**: Activity monitoring and export
- **About Tab**: Application information

### Terminal User Interface (TUI)

#### Navigation
```
Arrow Keys    - Navigate menus
Enter         - Select option
Esc           - Go back/exit
Space         - Toggle selections
Tab           - Switch between panels
```

#### Menu Structure
```
Main Menu
├── Container Management
│   ├── List Containers
│   ├── Run Container
│   ├── Start/Stop/Remove
│   └── Container Logs
├── Image Management
│   ├── List Images
│   ├── Import Image
│   ├── Build Image
│   └── Remove Image
├── Volume Management
│   ├── List Volumes
│   ├── Create Volume
│   └── Remove Volume
├── CRI Server
│   ├── Start Server
│   ├── Stop Server
│   └── Server Status
└── System Information
```

## Troubleshooting

### Common Issues

#### CGO Compilation Errors (Windows)
```
Error: cgo: C compiler "gcc" not found
```
**Solution:**
```bash
# Install MinGW-w64
# Add to PATH: C:\msys64\ucrt64\bin
export PATH="/c/msys64/ucrt64/bin:$PATH"
go env -w CGO_ENABLED=1
```

#### Container.Box Type Errors
```
Error: container.Box is not a type
```
**Solution:**
```bash
# Update Fyne import in GUI files
# Change: container.Box to container.NewHBox() or container.NewVBox()
```

#### Permission Denied
```
Error: permission denied
```
**Solution:**
```bash
# Linux/macOS
sudo chmod +x servin
sudo chown $USER:$USER servin

# Windows (Run as Administrator)
# Right-click PowerShell -> "Run as Administrator"
```

#### Port Already in Use
```
Error: address already in use
```
**Solution:**
```bash
# Find process using port
netstat -tulpn | grep :8080
lsof -i :8080

# Kill process or use different port
servin cri start --port 8081
```

### Debugging

#### Enable Debug Logging
```bash
# CLI debug mode
servin --verbose --log-level debug command

# GUI debug mode
FYNE_DEBUG=1 ./servin-gui
servin gui --dev

# TUI debug mode
servin desktop --verbose
```

#### Log Files
```bash
# Default log locations
# Linux: ~/.local/share/servin/logs/
# macOS: ~/Library/Logs/servin/
# Windows: %APPDATA%\servin\logs\

# Custom log file
servin --log-file /path/to/custom.log command
```

### Performance Tuning

#### Resource Limits
```bash
# Set memory limits
export SERVIN_MEMORY_LIMIT=2G

# Set CPU limits
export SERVIN_CPU_LIMIT=2

# Configure cache
export SERVIN_CACHE_SIZE=1G
```

#### Network Configuration
```bash
# Custom network settings
servin config set network.bridge servin0
servin config set network.subnet 172.17.0.0/16
```

## API Reference

### REST API Endpoints

#### CRI Server API
```
GET    /health                 - Health check
GET    /cri/runtime/version    - Runtime version
GET    /cri/runtime/containers - List containers
POST   /cri/runtime/containers - Create container
GET    /cri/image/list        - List images
POST   /cri/image/pull        - Pull image
```

#### Internal API
```
GET    /api/v1/containers     - List containers
POST   /api/v1/containers     - Create container
GET    /api/v1/images         - List images
GET    /api/v1/volumes        - List volumes
GET    /api/v1/system/info    - System information
```

### Configuration Files

#### Main Configuration
```yaml
# ~/.servin/config.yaml
runtime:
  driver: native
  data_root: /var/lib/servin
  exec_root: /var/run/servin

network:
  bridge: servin0
  subnet: 172.17.0.0/16
  
logging:
  level: info
  file: /var/log/servin.log

gui:
  theme: auto
  port: 8081
  refresh_interval: 5s
```

#### CRI Configuration
```yaml
# ~/.servin/cri.yaml
server:
  listen: "0.0.0.0:8080"
  socket: "/var/run/servin/servin.sock"
  
kubernetes:
  version: "v1alpha2"
  namespace: default
```

## Contributing

### Development Setup
```bash
# Fork and clone
git clone https://github.com/yourusername/servin
cd servin/gocontainer

# Install development dependencies
go mod download
make dev-setup

# Run tests
make test
go test ./...

# Format code
make fmt
go fmt ./...
```

### Code Guidelines
- Follow Go best practices and conventions
- Add tests for new features
- Update documentation
- Use semantic versioning
- Write clear commit messages

### Submitting Changes
1. Create feature branch: `git checkout -b feature/amazing-feature`
2. Make changes and add tests
3. Run tests: `make test`
4. Commit changes: `git commit -m "Add amazing feature"`
5. Push branch: `git push origin feature/amazing-feature`
6. Open Pull Request

### Architecture Overview
```
├── cmd/                    # Command line interfaces
│   ├── main.go            # Main CLI
│   ├── gui.go             # GUI command
│   ├── servin-gui/        # Visual GUI application
│   └── servin-desktop/    # Terminal UI application
├── pkg/                   # Core packages
│   ├── container/         # Container management
│   ├── image/             # Image operations
│   ├── volume/            # Volume management
│   └── cri/               # CRI server implementation
├── docs/                  # Documentation
├── examples/              # Usage examples
└── test/                  # Test files
```

## License
This project is licensed under the Apache License 2.0. See [LICENSE](LICENSE) file for details.

## Support
- GitHub Issues: https://github.com/yourorg/servin/issues
- Documentation: https://servin.io/docs
- Community: https://discord.gg/servin
- Email: support@servin.io
