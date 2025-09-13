# Servin Container Runtime - Installation Guide

This document provides comprehensive installation instructions for Servin Container Runtime across all supported platforms.

## Overview

Servin is a lightweight, Docker-compatible container runtime with a modern GUI interface. The installation packages include:

- **Core Runtime**: Container management engine with CRI compatibility
- **GUI Application**: Desktop interface for container management
- **System Service**: Background service for automatic startup
- **Command Line Tools**: Full CLI interface for scripting and automation

## Platform Support

| Platform | Architecture | GUI Support | Service Support | Package Format |
|----------|-------------|-------------|-----------------|----------------|
| Windows  | x64         | ✅ Fyne     | ✅ Windows Service | ZIP + PowerShell |
| Linux    | x64         | ✅ Fyne*    | ✅ systemd/SysV | tar.gz + Bash |
| macOS    | x64         | ✅ Fyne     | ✅ launchd | tar.gz + Bash |

*Linux GUI requires X11 development libraries

## Windows Installation

### Prerequisites
- Windows 10/11 or Windows Server 2019/2022
- Administrator privileges
- PowerShell 5.1 or later

### Installation Steps

1. **Download Package**
   ```
   servin-windows-1.0.0.zip
   ```

2. **Extract Archive**
   - Right-click the ZIP file
   - Select "Extract All..."
   - Choose destination folder

3. **Run Installer**
   - Right-click PowerShell
   - Select "Run as Administrator"
   - Navigate to extracted folder
   - Run: `.\install.ps1`

### What Gets Installed

| Component | Location | Description |
|-----------|----------|-------------|
| Binaries | `C:\Program Files\Servin\` | servin.exe, servin-gui.exe |
| Data | `C:\ProgramData\Servin\` | Container data, volumes, images |
| Config | `C:\ProgramData\Servin\config\` | servin.conf |
| Logs | `C:\ProgramData\Servin\logs\` | Runtime logs |
| Service | Windows Services | "ServinRuntime" service |
| Shortcuts | Desktop & Start Menu | GUI shortcuts |

### Service Management
```powershell
# Start service
Start-Service ServinRuntime

# Stop service
Stop-Service ServinRuntime

# Check status
Get-Service ServinRuntime

# View logs
Get-Content "C:\ProgramData\Servin\logs\servin.log" -Tail 50
```

### Uninstallation
- From Start Menu: "Uninstall Servin"
- Or run: `C:\Program Files\Servin\uninstall.ps1`

## Linux Installation

### Prerequisites
- Ubuntu 18.04+, CentOS 7+, Debian 9+, or compatible
- Root privileges (sudo)
- systemd or SysV init system

### Installation Steps

1. **Download Package**
   ```bash
   wget https://releases.servin.io/servin-linux-1.0.0.tar.gz
   ```

2. **Extract Archive**
   ```bash
   tar -xzf servin-linux-1.0.0.tar.gz
   cd servin-linux-1.0.0
   ```

3. **Run Installer**
   ```bash
   sudo ./install.sh
   ```

### What Gets Installed

| Component | Location | Description |
|-----------|----------|-------------|
| Binaries | `/usr/local/bin/` | servin, servin-gui |
| Data | `/var/lib/servin/` | Container data, volumes, images |
| Config | `/etc/servin/` | servin.conf |
| Logs | `/var/log/servin/` | Runtime logs |
| Service | systemd/SysV | servin.service |
| Desktop | `/usr/share/applications/` | servin-gui.desktop |
| User | System user | `servin` user account |

### Service Management

**systemd (Ubuntu, CentOS 7+, Debian 9+):**
```bash
# Enable and start
sudo systemctl enable servin
sudo systemctl start servin

# Check status
sudo systemctl status servin

# View logs
sudo journalctl -u servin -f
```

**SysV (older systems):**
```bash
# Start service
sudo service servin start

# Check status
sudo service servin status

# View logs
sudo tail -f /var/log/servin/servin.log
```

### GUI Prerequisites (Optional)
For GUI support on Linux:
```bash
# Ubuntu/Debian
sudo apt-get install libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libgl1-mesa-dev

# CentOS/RHEL
sudo yum install libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel mesa-libGL-devel
```

### Uninstallation
```bash
sudo /usr/local/bin/servin-uninstall
```

## macOS Installation

### Prerequisites
- macOS 10.12 (Sierra) or later
- Administrator privileges
- Xcode Command Line Tools (for GUI)

### Installation Steps

1. **Download Package**
   ```bash
   curl -L -o servin-macos-1.0.0.tar.gz https://releases.servin.io/servin-macos-1.0.0.tar.gz
   ```

2. **Extract Archive**
   ```bash
   tar -xzf servin-macos-1.0.0.tar.gz
   cd servin-macos-1.0.0
   ```

3. **Run Installer**
   ```bash
   sudo ./install.sh
   ```

### What Gets Installed

| Component | Location | Description |
|-----------|----------|-------------|
| Binaries | `/usr/local/bin/` | servin, servin-gui |
| Data | `/usr/local/var/lib/servin/` | Container data, volumes, images |
| Config | `/usr/local/etc/servin/` | servin.conf |
| Logs | `/usr/local/var/log/servin/` | Runtime logs |
| Service | launchd | com.servin.runtime.plist |
| App Bundle | `/Applications/` | Servin GUI.app |
| User | System user | `_servin` user account |

### Service Management
```bash
# Check status
sudo launchctl list | grep servin

# Stop service
sudo launchctl unload /Library/LaunchDaemons/com.servin.runtime.plist

# Start service
sudo launchctl load /Library/LaunchDaemons/com.servin.runtime.plist

# View logs
tail -f /usr/local/var/log/servin/servin.log
```

### Uninstallation
```bash
sudo /usr/local/bin/servin-uninstall
```

## Docker Installation

### Quick Start
```bash
# Load Docker image
docker load < servin-docker-1.0.0.tar.gz

# Run container
docker run -d -p 10250:10250 --name servin servin:1.0.0

# Check status
docker logs servin
```

### Docker Compose
```yaml
version: '3.8'
services:
  servin:
    image: servin:1.0.0
    ports:
      - "10250:10250"
    volumes:
      - servin_data:/var/lib/servin
    restart: unless-stopped

volumes:
  servin_data:
```

## Configuration

### Default Configuration File
All platforms use a similar configuration format:

```ini
# Servin Configuration File

# Data directory
data_dir=/var/lib/servin

# Log settings
log_level=info
log_file=/var/log/servin/servin.log

# Runtime settings
runtime=native

# Network settings
bridge_name=servin0

# CRI settings
cri_port=10250
cri_enabled=false
```

### Configuration Locations

| Platform | Configuration File |
|----------|-------------------|
| Windows | `C:\ProgramData\Servin\config\servin.conf` |
| Linux | `/etc/servin/servin.conf` |
| macOS | `/usr/local/etc/servin/servin.conf` |

### Common Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `data_dir` | Directory for container data | Platform-specific |
| `log_level` | Logging verbosity (debug, info, warn, error) | `info` |
| `log_file` | Log file location | Platform-specific |
| `runtime` | Container runtime backend | `native` |
| `bridge_name` | Default network bridge name | `servin0` |
| `cri_port` | CRI server port | `10250` |
| `cri_enabled` | Enable CRI server | `false` |

## Usage

### GUI Interface
- **Windows**: Desktop shortcut or Start Menu
- **Linux**: Applications menu or `servin-gui` command
- **macOS**: Applications folder or `servin-gui` command

### Command Line Interface
```bash
# Show help
servin --help

# List containers
servin ps

# Run a container
servin run alpine:latest echo "Hello World"

# Start daemon mode
servin daemon

# Show version
servin version
```

### API Usage
When CRI is enabled, Servin provides a Kubernetes-compatible API:
```bash
# Health check
curl http://localhost:10250/healthz

# Runtime info
curl http://localhost:10250/api/v1/runtime/status
```

## Troubleshooting

### Common Issues

**Windows: "Execution Policy" Error**
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

**Linux: GUI Not Starting**
```bash
# Install X11 libraries
sudo apt-get install libx11-6 libxcursor1 libxrandr2 libxinerama1 libxi6 libgl1-mesa-glx

# Check DISPLAY variable
echo $DISPLAY
```

**macOS: "App Can't Be Opened" Error**
```bash
# Remove quarantine attribute
sudo xattr -rd com.apple.quarantine "/Applications/Servin GUI.app"
```

**Service Won't Start**
1. Check configuration file syntax
2. Verify permissions on data directories
3. Check port availability (10250)
4. Review log files for errors

### Log Locations

| Platform | Log Location |
|----------|--------------|
| Windows | `C:\ProgramData\Servin\logs\servin.log` |
| Linux | `/var/log/servin/servin.log` |
| macOS | `/usr/local/var/log/servin/servin.log` |

### Support

For issues and support:
- GitHub Issues: https://github.com/yourusername/servin/issues
- Documentation: https://docs.servin.io
- Community: https://discord.gg/servin

## Building from Source

See `build.sh` (Linux/macOS) or `build.ps1` (Windows) for building distributables:

```bash
# Build all platforms
./build.sh

# Build specific platform
./build.sh windows
./build.sh linux
./build.sh macos
```

## Security Considerations

1. **Run as Non-Root**: Services run as dedicated users (`servin`, `_servin`)
2. **Directory Permissions**: Restricted access to data directories
3. **Network Isolation**: Default bridge network configuration
4. **Log Security**: Logs stored in protected system directories
5. **Service Hardening**: Security features enabled in service configurations

## License

Servin Container Runtime is released under the Apache 2.0 License.
See LICENSE file for details.
