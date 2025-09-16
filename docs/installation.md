---
layout: default
title: Installation
permalink: /installation/
---

# 🛠 Installation

## Quick Installation

Choose your platform and follow the installation instructions:

<div class="feature-grid">
  <div class="feature-box">
    <div class="feature-icon">🪟</div>
    <h4>Windows</h4>
    <p>Professional NSIS installer with service integration</p>
    <span class="badge badge-success">Wizard installer</span>
  </div>
  <div class="feature-box">
    <div class="feature-icon">🐧</div>
    <h4>Linux</h4>
    <p>GUI installer with systemd integration</p>
    <span class="badge badge-success">GUI installer</span>
  </div>
  <div class="feature-box">
    <div class="feature-icon">🍎</div>
    <h4>macOS</h4>
    <p>Native installer following Apple HIG</p>
    <span class="badge badge-success">Native installer</span>
  </div>
</div>

## Platform-Specific Instructions

### 🪟 Windows Installation

#### Option 1: Using the Installer (Recommended)

1. **Download** the latest installer:
   ```
   servin_1.0.0_windows_amd64_installer.exe
   ```

2. **Run as Administrator**:
   - Right-click the installer
   - Select "Run as administrator"

3. **Follow the Installation Wizard**:
   - Choose installation directory
   - Select components to install
   - Configure service options

4. **Launch from Start Menu**:
   - Find "Servin Container Runtime" in Start Menu
   - Or run `servin` from Command Prompt

#### What's Included:
- ✅ Servin CLI (`servin.exe`)
- ✅ Desktop GUI (`servin-gui.exe`)
- ✅ Terminal UI (`servin-desktop.exe`)
- ✅ Windows Service integration
- ✅ Start Menu shortcuts
- ✅ Add/Remove Programs entry
- ✅ Automatic PATH configuration

### 🐧 Linux Installation

#### Option 1: Using the GUI Installer (Recommended)

```bash
# Download installer
wget https://github.com/immyemperor/servin/releases/latest/download/servin_1.0.0_linux_amd64.tar.gz

# Extract and run installer
tar -xzf servin_1.0.0_linux_amd64.tar.gz
cd servin-linux-amd64
sudo ./ServinInstaller.sh
```

#### Option 2: Manual Installation

```bash
# Download binary package
wget https://github.com/immyemperor/servin/releases/latest/download/servin_1.0.0_linux_amd64.tar.gz

# Extract package
tar -xzf servin_1.0.0_linux_amd64.tar.gz

# Move binaries to system path
sudo cp servin* /usr/local/bin/

# Create systemd service
sudo systemctl enable servin
sudo systemctl start servin
```

#### Supported Distributions:
- ✅ Ubuntu 20.04+
- ✅ Debian 11+
- ✅ CentOS 8+
- ✅ Fedora 35+
- ✅ Arch Linux
- ✅ openSUSE Leap 15+

### 🍎 macOS Installation

#### Option 1: Using the GUI Installer (Recommended)

```bash
# Download installer
curl -O https://github.com/immyemperor/servin/releases/latest/download/servin_1.0.0_macos_universal.tar.gz

# Extract and run installer
tar -xzf servin_1.0.0_macos_universal.tar.gz
cd servin-macos-universal
sudo ./ServinInstaller.command
```

#### Option 2: Using Homebrew

```bash
# Add tap
brew tap immyemperor/servin

# Install servin
brew install servin

# Start service
brew services start servin
```

#### Requirements:
- ✅ macOS 10.15 (Catalina) or later
- ✅ Python 3.8+ (for installer)
- ✅ Command Line Tools for Xcode

## Building from Source

### Prerequisites

#### Required Software
- **Go 1.24+** - Latest Go version with module support
- **Git** - Version control for source code
- **CGO enabled** - Required for GUI compilation

#### Platform-Specific Requirements

**Windows:**
```powershell
# Install Go
winget install GoLang.Go

# Install Git
winget install Git.Git

# Install MinGW-w64 UCRT64 (for CGO)
# Download from: https://www.mingw-w64.org/downloads/
```

**Linux:**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go git build-essential

# CentOS/RHEL/Fedora
sudo dnf install golang git gcc

# Arch Linux
sudo pacman -S go git gcc
```

**macOS:**
```bash
# Using Homebrew
brew install go git

# Using MacPorts
sudo port install go git
```

### Build Commands

#### 🔨 Clone and Build

```bash
# Clone repository
git clone https://github.com/immyemperor/servin.git
cd servin

# Build all components
make build                     # Linux/macOS
.\build.ps1                    # Windows PowerShell

# Build specific components
go build -o servin main.go                    # CLI only
go build -o servin-desktop cmd/servin-desktop/ # TUI only
go build -o servin-gui cmd/servin-gui/         # GUI only
```

#### 🎯 Cross-Platform Building

```bash
# Build for all platforms using the new build system
./build-all.sh                    # All platforms with distribution packages

# Build for specific platforms
PLATFORM=windows ./build-all.sh   # Windows with ZIP and NSIS installer
PLATFORM=linux ./build-all.sh     # Linux with TAR.GZ and wizard installer
PLATFORM=mac ./build-all.sh       # macOS universal binary with wizard installer

# Clean previous builds
./build-all.sh --clean-all

# Using Go directly for development
GOOS=windows GOARCH=amd64 go build -o servin.exe main.go
GOOS=linux GOARCH=amd64 go build -o servin main.go
GOOS=darwin GOARCH=amd64 go build -o servin main.go
```

#### 📦 Distribution Packages

The new build system creates professional distribution packages:

```bash
# Build all distributions
./build-all.sh

# Outputs:
# - dist/servin_1.0.0_windows_amd64.zip           (Windows ZIP archive)
# - dist/servin_1.0.0_windows_amd64_installer.exe (Windows NSIS installer)
# - dist/servin_1.0.0_linux_amd64.tar.gz         (Linux distribution)
# - dist/servin_1.0.0_macos_universal.tar.gz     (macOS universal binary)
```

**What's included in each package:**
- ✅ **servin** - CLI container runtime
- ✅ **servin-desktop** - Terminal User Interface (TUI)
- ✅ **servin-webview** - Modern WebView GUI interface
- ✅ **Wizard installers** - Interactive GUI installation wizards
- ✅ **Professional icons** - Multi-format icon set
- ✅ **Documentation** - README, LICENSE, and usage guides

### Development Build

For development and testing:

```bash
# Install in development mode
go install ./...

# Run tests
go test ./...

# Run with live reload
go run main.go daemon --dev

# Build with debug information
go build -gcflags="all=-N -l" -o servin-debug main.go
```

## Post-Installation Setup

### 1. Verify Installation

```bash
# Check version
servin version

# Verify daemon can start
servin daemon --dry-run

# Test basic functionality
servin info
```

### 2. Configure Environment

```bash
# Set data directory (optional)
export SERVIN_DATA_ROOT="/var/lib/servin"

# Set log level
export SERVIN_LOG_LEVEL="info"

# Add to shell profile
echo 'export SERVIN_DATA_ROOT="/var/lib/servin"' >> ~/.bashrc
```

### 3. Start Services

#### Linux (systemd)
```bash
# Enable and start service
sudo systemctl enable servin
sudo systemctl start servin

# Check status
sudo systemctl status servin
```

#### Windows (Service)
```powershell
# Service is automatically installed and started
# Check status in Services.msc or:
sc query servin
```

#### macOS (launchd)
```bash
# Service is automatically installed
# Check status:
launchctl list | grep servin
```

### 4. Test Installation

```bash
# Pull test image
servin pull hello-world

# Run test container
servin run hello-world

# Launch GUI (if installed)
servin-gui

# Launch TUI (if installed)
servin-desktop
```

## Verification

### Basic Functionality Test

```bash
# Check daemon status
servin info

# List containers
servin ps

# List images
servin images

# Test container operations
servin run --rm alpine:latest echo "Hello, Servin!"

# Test network operations
servin network ls

# Test volume operations
servin volume ls
```

### GUI Applications Test

```bash
# Test Terminal UI
servin-desktop

# Test Desktop GUI
servin-gui
```

### CRI Integration Test (if enabled)

```bash
# Test CRI endpoint
crictl --runtime-endpoint unix:///var/run/servin.sock version

# List CRI pods
crictl --runtime-endpoint unix:///var/run/servin.sock pods
```

## Troubleshooting Installation

### Common Issues

#### Permission Errors
```bash
# Linux: Add user to servin group
sudo usermod -aG servin $USER
newgrp servin

# Windows: Run as Administrator
```

#### Path Issues
```bash
# Check if servin is in PATH
which servin          # Linux/macOS
where servin          # Windows

# Add to PATH if needed
export PATH=$PATH:/usr/local/bin  # Linux/macOS
```

#### Service Issues
```bash
# Linux: Check service logs
sudo journalctl -u servin -f

# Check configuration
servin daemon --config-check

# Reset configuration
servin config reset
```

### Getting Help

If you encounter issues during installation:

1. **Check the logs**:
   ```bash
   servin logs
   ```

2. **Run diagnostics**:
   ```bash
   servin doctor
   ```

3. **Visit our troubleshooting guide**: [Troubleshooting]({{ '/troubleshooting' | relative_url }})

4. **Open an issue on GitHub**: [GitHub Issues]({{ site.github.repository_url }}/issues)

## Next Steps

After successful installation:

1. **Configure Servin**: [Configuration Guide]({{ '/configuration' | relative_url }})
2. **Try the Quick Start**: [Quick Start Guide]({{ '/quick-start' | relative_url }})
3. **Explore Features**: [Features Overview]({{ '/features' | relative_url }})

[Configure Servin →]({{ '/configuration' | relative_url }}){: .btn .btn-primary}
[Quick Start →]({{ '/quick-start' | relative_url }}){: .btn .btn-outline}
