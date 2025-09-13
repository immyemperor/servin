# Servin Container Runtime - Cross-Platform Installers

Comprehensive installation packages for Servin Container Runtime with service integration and GUI support.

## 🚀 Quick Start

### Windows
1. Download `servin-windows-1.0.0.zip`
2. Extract the archive
3. Right-click PowerShell → "Run as Administrator"
4. Navigate to extracted folder: `cd path\to\servin-windows-1.0.0`
5. Run installer: `.\install.ps1`

### Linux
1. Download `servin-linux-1.0.0.tar.gz`
2. Extract: `tar -xzf servin-linux-1.0.0.tar.gz`
3. Install: `cd servin-linux-1.0.0 && sudo ./install.sh`

### macOS
1. Download `servin-macos-1.0.0.tar.gz`
2. Extract: `tar -xzf servin-macos-1.0.0.tar.gz`
3. Install: `cd servin-macos-1.0.0 && sudo ./install.sh`

## 📦 What's Included

Each platform package contains:
- **Servin Runtime** (`servin`/`servin.exe`) - Core container management
- **GUI Application** (`servin-gui`/`servin-gui.exe`) - Desktop interface  
- **System Service** - Automatic startup integration
- **Installation Script** - Automated setup with permissions
- **Uninstaller** - Clean removal capability

## 🛠 Build from Source

### Prerequisites
- Go 1.24+ with CGO support
- Platform-specific GUI libraries (for GUI builds)

### Windows (PowerShell)
```powershell
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

# Check status
sudo systemctl status servin

# View logs
sudo journalctl -u servin -f
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
