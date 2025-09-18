# Servin Container Runtime - Cross-Pl## Platform-Specific Enhanced Installers

If you prefer to run the enhanced installers directly (or the wizard detects missing prerequisites):

### 🐧 Linux - `linux/install-with-vm.sh`orm Installation Guide

This directory contains enhanced installers for all supported platforms with automatic VM prerequisite installation.

## 🚀 Quick Installation (Recommended)

### Smart Wizard Installers
The easiest way to install Servin is using our smart wizard installers that automatically detect prerequisites and install them if needed:

**🐧 Linux & 🍎 macOS:**
```bash
# Download and run universal installer
curl -sSL https://install.servin.dev | bash

# Or download wizard manually
wget https://github.com/immyemperor/servin/releases/latest/download/install-wizard.sh
chmod +x install-wizard.sh
./install-wizard.sh
```

**🪟 Windows:**
```powershell
# Download and run wizard (as Administrator)
Invoke-WebRequest -Uri "https://github.com/immyemperor/servin/releases/latest/download/install-wizard.ps1" -OutFile "install-wizard.ps1"
Set-ExecutionPolicy Bypass -Scope Process -Force
.\install-wizard.ps1
```

### Smart Wizard Features
- ✅ **Auto-detection**: Automatically detects your platform and existing installations
- ✅ **Prerequisite checking**: Verifies VM requirements and system resources
- ✅ **Smart installation**: Runs enhanced installers only if prerequisites are missing
- ✅ **Interactive or automated**: Supports both interactive and unattended installation
- ✅ **Comprehensive verification**: Tests installation and provides troubleshooting guidance

---

## Platform-Specific Installers

### � Linux - `linux/install-with-vm.sh`
**Features:**
- Automatic distribution detection (Ubuntu, Debian, Fedora, CentOS, Arch)
- KVM/QEMU installation and configuration
- libvirt setup with default network
- Python WebView dependencies
- Development tools (Packer, Python packages)
- Systemd service creation
- Comprehensive hardware checking

**Requirements:**
- Root access (sudo)
- 4GB+ RAM (8GB+ recommended)
- 5GB+ free disk space
- CPU with virtualization support

**Installation:**
```bash
sudo chmod +x linux/install-with-vm.sh
sudo ./linux/install-with-vm.sh
```

**Supported Distributions:**
- Ubuntu 18.04+ / Debian 9+
- Fedora 32+ / CentOS 8+ / RHEL 8+
- Arch Linux (current)

---

### 🍎 macOS - `macos/install-with-vm.sh`
**Features:**
- Homebrew automatic installation
- QEMU virtualization setup
- Virtualization.framework support (macOS 11+)
- Python WebView with Cocoa backend
- Development tools via Homebrew
- Launchd service configuration
- Apple Silicon and Intel support

**Requirements:**
- macOS 10.15+ (Catalina or later)
- Xcode Command Line Tools
- 8GB+ RAM (16GB+ recommended)
- 10GB+ free disk space
- Admin access

**Installation:**
```bash
chmod +x macos/install-with-vm.sh
./macos/install-with-vm.sh
```

**Architecture Support:**
- Apple Silicon (M1/M2/M3) - Native Virtualization.framework
- Intel x86_64 - QEMU with HVF acceleration

---

### 🪟 Windows - `windows/install-with-vm.ps1`
**Features:**
- Chocolatey package manager setup
- Multi-provider VM support (Hyper-V, VirtualBox, WSL2)
- Python WebView with tkinter backend
- Automatic Windows feature enabling
- Development tools installation
- Windows Service creation
- PowerShell 5.1+ and PowerShell Core support

**Requirements:**
- Windows 10/11 (Build 1903+)
- Administrator privileges
- 8GB+ RAM (16GB+ recommended)
- 10GB+ free disk space
- PowerShell 5.1+ or PowerShell Core

**Installation:**
```powershell
# Run as Administrator
Set-ExecutionPolicy Bypass -Scope Process -Force
.\windows\install-with-vm.ps1
```

**VM Provider Options:**
- **Hyper-V** (Windows Pro/Enterprise/Education)
- **VirtualBox** (All Windows editions)
- **WSL2** (Windows 10 2004+/Windows 11)

## Wizard Installer Options

All wizard installers support automated installation modes:

**Linux/macOS:**
```bash
# Automated full installation
./install-wizard.sh --auto

# Install VM prerequisites only
./install-wizard.sh --vm-only

# Skip VM setup
./install-wizard.sh --no-vm
```

**Windows:**
```powershell
# Automated full installation
.\install-wizard.ps1 -Auto

# Install VM prerequisites only
.\install-wizard.ps1 -VmOnly

# Skip VM setup  
.\install-wizard.ps1 -NoVm

# Skip confirmation prompts
.\install-wizard.ps1 -Auto -Force
```

---

## Common Installation Flow

All installers follow a similar comprehensive installation process:

### 1. **System Prerequisites Check**
- Hardware virtualization support
- Available memory and disk space
- Operating system compatibility
- Network connectivity
- Required permissions

### 2. **Package Manager Setup**
- **Linux**: Native package managers (apt, yum, dnf, pacman)
- **macOS**: Homebrew installation and configuration
- **Windows**: Chocolatey installation and setup

### 3. **Virtualization Prerequisites**
- **Linux**: KVM, QEMU, libvirt installation and configuration
- **macOS**: QEMU and Virtualization.framework setup
- **Windows**: Hyper-V, VirtualBox, or WSL2 installation

### 4. **Development Tools Installation**
- Python 3.8+ with GUI frameworks
- HashiCorp Packer for VM image building
- Essential development utilities

### 5. **Servin Installation**
- Binary installation to system directories
- Configuration file creation
- VM provider configuration
- Service registration

### 6. **System Integration**
- **Linux**: systemd service setup
- **macOS**: launchd service configuration  
- **Windows**: Windows Service installation

### 7. **Verification & Testing**
- VM provider functionality testing
- Basic CLI command verification
- Hardware acceleration validation

---

## Post-Installation Steps

After running any installer:

### 1. **Environment Setup**
```bash
# Linux/macOS: Logout and login for group membership
# Windows: Restart PowerShell session

# Test basic functionality
servin version
servin vm status
```

### 2. **VM Initialization**
```bash
# Initialize VM support
servin vm init

# Enable VM mode
servin vm enable

# List available providers
servin vm list-providers
```

### 3. **Test VM Functionality**
```bash
# Test with a simple container
servin run --vm alpine echo "Hello from VM!"

# Test with interactive container
servin run --vm -it alpine sh
```

### 4. **GUI Access** (if installed)
```bash
# Launch graphical interface
servin-gui

# Or use terminal interface
servin-tui
```

---

## Troubleshooting

### Common Issues

**VM Provider Not Available:**
- Ensure virtualization is enabled in BIOS/UEFI
- Check that hardware acceleration is working
- Verify user permissions (groups on Linux, admin on Windows)

**Installation Permissions:**
- **Linux**: Run with `sudo`
- **macOS**: Run as normal user (will prompt for admin when needed)
- **Windows**: Run PowerShell as Administrator

**Missing Dependencies:**
- Run the installer again - it will detect and install missing components
- Check internet connectivity for package downloads
- Verify system meets minimum requirements

### Platform-Specific Troubleshooting

#### Linux
```bash
# Check KVM device
ls -la /dev/kvm

# Verify libvirt service
systemctl status libvirtd

# Test QEMU
qemu-system-x86_64 --version
```

#### macOS
```bash
# Check Homebrew
brew doctor

# Verify QEMU
qemu-system-x86_64 --version

# Check Virtualization.framework (macOS 11+)
servin vm check-virtualization
```

#### Windows
```powershell
# Check Hyper-V status
Get-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V

# Verify VirtualBox
VBoxManage --version

# Check WSL2
wsl --status
```

---

## Advanced Configuration

### Custom Installation Directories

**Linux:**
```bash
export INSTALL_DIR="/opt/servin"
export DATA_DIR="/opt/servin/data"
sudo ./install-with-vm.sh
```

**macOS:**
```bash
# Edit installer script to modify directories
# Default: /usr/local/bin, /usr/local/var/lib/servin
```

**Windows:**
```powershell
# Use parameters
.\install-with-vm.ps1 -InstallPath "C:\Tools\Servin" -DataPath "C:\Data\Servin"
```

### Skip Components

All installers support skipping certain components:

```bash
# Linux
sudo ./install-with-vm.sh --no-service --no-vm

# macOS  
./install-with-vm.sh --no-service --no-gui

# Windows
.\install-with-vm.ps1 -SkipService -SkipVM
```

---

## Documentation References

- **Complete Setup Guide**: [VM_PREREQUISITES.md](../VM_PREREQUISITES.md)
- **CLI Documentation**: [docs/cli.md](../docs/cli.md)
- **VM Integration**: [docs/VM_INTEGRATION.md](../docs/VM_INTEGRATION.md)
- **Troubleshooting**: [docs/troubleshooting.md](../docs/troubleshooting.md)

---

## Support & Contributing

For issues with the installers:
1. Check the troubleshooting section above
2. Review the detailed prerequisites in `VM_PREREQUISITES.md`
3. Submit issues with full system information and installer logs

The installers are designed to be idempotent - you can run them multiple times safely to fix issues or update components.
# Build all platforms
.\build.ps1

# Build specific platform
.\build.ps1 -Target windows
.\build.ps1 -Target linux
.\build.ps1 -Target macos
```

### Linux/macOS (Bash)
```bash
# Build all platforms
./build.sh

# Build specific platform
./build.sh windows
./build.sh linux
./build.sh macos
```

## 🎯 Features

### Core Runtime
- **Docker-Compatible API** - Drop-in replacement for basic Docker workflows
- **Container Management** - Run, stop, list, remove containers
- **Image Handling** - Build, pull, push, tag images
- **Volume Management** - Persistent data storage
- **Network Isolation** - Container networking with bridges
- **Security** - User namespaces, capabilities, seccomp

### GUI Interface
- **Container Dashboard** - Visual container management
- **Image Browser** - Explore and manage images
- **Volume Manager** - Handle persistent storage
- **Log Viewer** - Real-time container logs
- **Resource Monitoring** - CPU, memory, network stats

### System Integration
- **Windows Service** - `ServinRuntime` service with auto-start
- **Linux systemd/SysV** - Native service integration
- **macOS launchd** - Background daemon support
- **PATH Integration** - Command-line access from anywhere
- **Desktop Shortcuts** - Quick GUI access

## 🔧 Configuration

Default configuration locations:
- **Windows**: `C:\ProgramData\Servin\config\servin.conf`
- **Linux**: `/etc/servin/servin.conf`
- **macOS**: `/usr/local/etc/servin/servin.conf`

Key settings:
```ini
# Data directory
data_dir=/var/lib/servin

# Logging
log_level=info
log_file=/var/log/servin/servin.log

# Runtime
runtime=native

# Network
bridge_name=servin0

# CRI Server (Kubernetes compatibility)
cri_port=10250
cri_enabled=false
```

## 🚨 Service Management

### Windows
```powershell
# Start/stop service
Start-Service ServinRuntime
Stop-Service ServinRuntime

# Check status
Get-Service ServinRuntime

# View logs
Get-Content "C:\ProgramData\Servin\logs\servin.log" -Tail 50
```

### Linux (systemd)
```bash
# Enable and start
sudo systemctl enable servin
sudo systemctl start servin

The installers are designed to be idempotent - you can run them multiple times safely to fix issues or update components.
```

### macOS (launchd)
```bash
# Check status
sudo launchctl list | grep servin

# View logs
tail -f /usr/local/var/log/servin/servin.log
```

## 📋 Command Line Usage

```bash
# Run a container
servin run alpine:latest echo "Hello from Servin!"

# List containers
servin ls

# Build an image
servin build -t myapp:latest .

# Manage volumes
servin volume create mydata
servin volume ls

# View logs
servin logs container_name

# Launch GUI
servin gui

# Start daemon mode
servin daemon
```

## 🔒 Security Features

- **Non-root execution** - Services run as dedicated users
- **Directory isolation** - Secure data directory permissions
- **Network isolation** - Default bridge network separation
- **Resource limits** - CPU, memory, and I/O constraints
- **Capability dropping** - Minimal privilege containers

## 📊 System Requirements

### Minimum
- **RAM**: 512MB available
- **Storage**: 1GB free space
- **CPU**: Single core (x64 architecture)

### Recommended
- **RAM**: 2GB+ for GUI and multiple containers
- **Storage**: 10GB+ for images and container data
- **CPU**: Multi-core for better performance

### Platform Specific
- **Windows**: Windows 10/11, Windows Server 2019/2022
- **Linux**: Ubuntu 18.04+, CentOS 7+, Debian 9+
- **macOS**: macOS 10.12 (Sierra) or later

## 🐛 Troubleshooting

### Common Issues

**Windows: PowerShell Execution Policy**
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

**Linux: Missing GUI Dependencies**
```bash
# Ubuntu/Debian
sudo apt-get install libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libgl1-mesa-dev

# CentOS/RHEL
sudo yum install libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel mesa-libGL-devel
```

**Service Won't Start**
1. Check configuration file syntax
2. Verify data directory permissions
3. Ensure port 10250 is available
4. Review log files for errors

### Getting Help
- **Documentation**: See `INSTALL.md` for detailed instructions
- **Logs**: Check platform-specific log locations
- **Issues**: Submit bug reports with log output
- **Community**: Join discussions and get support

## 📈 Performance Tuning

### For High Container Density
```ini
# Increase file descriptor limits
# In service configuration or system limits

# Optimize logging
log_level=warn

# Adjust resource limits
# Configure cgroups appropriately
```

### For GUI Performance
- Ensure graphics drivers are up to date
- Allocate sufficient RAM for desktop environment
- Consider running GUI on systems with dedicated graphics

## 🔄 Upgrade Process

1. Stop Servin service
2. Backup configuration and data directories
3. Run new installer (overwrites binaries)
4. Start service
5. Verify functionality

Configuration files are preserved during upgrades.

## 📜 License

Servin Container Runtime is released under the Apache 2.0 License.

---

**Built with ❤️ for containerization simplicity**
