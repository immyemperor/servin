# Servin VM Containerization Prerequisites

This document outlines the prerequisites and setup requirements for Servin's VM-based containerization system across different platforms.

## Overview

Servin provides universal VM-based containerization that ensures consistent container behavior across macOS, Linux, and Windows. The system uses platform-specific virtualization technologies for optimal performance:

- **Linux**: KVM with QEMU acceleration
- **macOS**: Virtualization.framework with QEMU fallback
- **Windows**: Hyper-V, VirtualBox, or WSL2

## System Requirements

### Minimum Hardware Requirements

- **CPU**: x64 processor with virtualization support (Intel VT-x or AMD-V)
- **RAM**: 4GB minimum (2GB for host + 2GB for VMs)
- **Storage**: 5GB available disk space
- **Network**: Internet connection for VM image downloads

### Recommended Hardware Requirements

- **CPU**: Multi-core x64 processor with virtualization enabled in BIOS/UEFI
- **RAM**: 8GB or more for multiple concurrent containers
- **Storage**: SSD storage for optimal VM performance
- **Network**: Broadband connection for faster image downloads

## Platform-Specific Prerequisites

### 🐧 Linux Prerequisites

#### Required Packages
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install qemu-system-x86 qemu-utils qemu-kvm libvirt-clients libvirt-daemon-system bridge-utils genisoimage cpu-checker

# RHEL/CentOS/Fedora
sudo dnf install qemu-kvm qemu-img libvirt libvirt-client libvirt-daemon-system bridge-utils genisoimage

# Arch Linux
sudo pacman -S qemu-desktop libvirt bridge-utils cdrtools
```

#### Hardware Acceleration Setup
1. **Enable CPU Virtualization**:
   - Enter BIOS/UEFI settings during boot
   - Enable Intel VT-x or AMD-V
   - Enable VT-d/IOMMU if available

2. **Verify KVM Support**:
   ```bash
   # Check if CPU supports virtualization
   grep -E "(vmx|svm)" /proc/cpuinfo
   
   # Check KVM availability
   sudo kvm-ok
   ```

3. **Configure KVM Access**:
   ```bash
   # Add user to kvm and libvirt groups
   sudo usermod -a -G kvm,libvirt $USER
   
   # Start and enable libvirt service
   sudo systemctl enable --now libvirtd
   
   # Logout and login for group changes to take effect
   ```

4. **Verify Setup**:
   ```bash
   # Test KVM access
   ls -la /dev/kvm
   
   # Test libvirt connection
   virsh version
   ```

#### Troubleshooting Linux
- **KVM not available**: Check if virtualization is enabled in BIOS
- **Permission denied**: Ensure user is in kvm and libvirt groups
- **libvirtd not running**: `sudo systemctl start libvirtd`

### 🍎 macOS Prerequisites

#### System Requirements
- **macOS**: 11.0 (Big Sur) or later for Virtualization.framework
- **Hardware**: Intel Mac or Apple Silicon (M1/M2/M3)

#### Required Tools
1. **Homebrew Package Manager**:
   ```bash
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
   ```

2. **QEMU Installation**:
   ```bash
   brew install qemu
   ```

3. **Xcode Command Line Tools**:
   ```bash
   xcode-select --install
   ```

#### Hardware Acceleration Setup
1. **Check Virtualization.framework**:
   ```bash
   # Verify hypervisor support
   sysctl kern.hv_support
   ```

2. **System Preferences** (if needed):
   - Ensure SIP (System Integrity Protection) allows virtualization
   - Check Security & Privacy settings for hypervisor access

#### Troubleshooting macOS
- **Virtualization.framework not available**: Update to macOS 11.0+
- **Permission errors**: Check System Preferences > Security & Privacy
- **Performance issues**: Ensure not running under Rosetta translation

### 🪟 Windows Prerequisites

#### System Requirements
- **Windows**: 10 Pro/Enterprise/Education or Windows 11
- **Architecture**: x64 processor
- **RAM**: 4GB minimum (8GB recommended)

#### Option 1: Hyper-V (Recommended)
1. **Enable Hyper-V**:
   ```powershell
   # Run as Administrator
   Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V-All
   
   # Or via GUI: Control Panel > Programs > Windows Features
   ```

2. **Enable Hardware Virtualization**:
   - Enter BIOS/UEFI settings
   - Enable Intel VT-x or AMD-V
   - Enable SLAT (Second Level Address Translation)

3. **Verify Hyper-V**:
   ```powershell
   Get-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V-All
   ```

#### Option 2: VirtualBox
1. **Download VirtualBox**:
   - Visit https://www.virtualbox.org/
   - Download and install latest version

2. **Configure VirtualBox**:
   ```powershell
   # Verify installation
   VBoxManage --version
   ```

#### Option 3: WSL2
1. **Enable WSL2**:
   ```powershell
   # Run as Administrator
   wsl --install
   
   # Or manual setup
   dism.exe /online /enable-feature /featurename:Microsoft-Windows-Subsystem-Linux /all /norestart
   dism.exe /online /enable-feature /featurename:VirtualMachinePlatform /all /norestart
   ```

2. **Set WSL2 as Default**:
   ```powershell
   wsl --set-default-version 2
   ```

#### Troubleshooting Windows
- **Hyper-V not available**: Requires Windows 10 Pro/Enterprise/Education
- **Virtualization disabled**: Enable in BIOS/UEFI settings
- **WSL2 installation fails**: Ensure Windows is up to date

## Automated Prerequisites Checking

Servin includes built-in prerequisites checking:

```bash
# Check system capabilities
servin vm init

# List available providers
servin vm list-providers

# Test specific providers
servin vm check-kvm      # Linux
servin vm check-virtualization  # macOS
servin vm check-hyperv   # Windows
servin vm check-virtualbox      # All platforms
```

## Network Configuration

### Default Network Setup
Servin automatically configures networking for VM containers:

- **Bridge Networking**: For container-to-host communication
- **NAT Networking**: For internet access from containers
- **Port Forwarding**: Automatic port mapping

### Custom Network Configuration
```yaml
# ~/.servin/network-config.yaml
network:
  bridge: servin-br0
  subnet: 192.168.100.0/24
  dhcp_range: 192.168.100.100-192.168.100.200
  dns_servers:
    - 8.8.8.8
    - 8.8.4.4
```

## Storage Configuration

### Default Storage Locations
- **Linux**: `~/.servin/vm/`
- **macOS**: `~/Library/Application Support/Servin/VM/`
- **Windows**: `%APPDATA%\Servin\VM\`

### Custom Storage Configuration
```yaml
# ~/.servin/storage-config.yaml
storage:
  vm_images: "/custom/path/images"
  vm_instances: "/custom/path/instances"
  max_image_cache: "10GB"
  compression: true
```

## Security Considerations

### Firewall Configuration
Servin may require firewall rules for VM networking:

```bash
# Linux (iptables)
sudo iptables -A INPUT -i servin-br0 -j ACCEPT
sudo iptables -A FORWARD -i servin-br0 -j ACCEPT

# macOS (pfctl) - Usually automatic

# Windows (Windows Firewall) - Usually automatic via Hyper-V
```

### User Permissions
- **Linux**: User must be in `kvm` and `libvirt` groups
- **macOS**: User needs developer tools access
- **Windows**: Administrator privileges may be required for initial setup

## Performance Optimization

### CPU Configuration
```yaml
# ~/.servin/vm-config.yaml
vm:
  default_cpu_cores: 2
  max_cpu_cores: 4
  cpu_limit_percent: 80
```

### Memory Configuration
```yaml
vm:
  default_memory: "1GB"
  max_memory: "4GB"
  memory_ballooning: true
```

### Storage Optimization
- Use SSD storage for VM images
- Enable compression for image storage
- Regular cleanup of unused images

## Getting Help

### Diagnostic Commands
```bash
# System information
servin vm status

# Detailed diagnostics
servin vm diagnose

# Log files
tail -f ~/.servin/logs/vm.log
```

### Common Issues
1. **Slow VM Performance**: Check hardware acceleration is enabled
2. **Network Issues**: Verify bridge networking configuration
3. **Storage Full**: Clean up unused VM images
4. **Permission Errors**: Check user group memberships

### Support Resources
- **Documentation**: https://servin.dev/docs
- **GitHub Issues**: https://github.com/immyemperor/servin/issues
- **Community Forum**: https://community.servin.dev

## Quick Start Verification

After setting up prerequisites, verify your installation:

```bash
# Initialize VM support
servin vm init

# Enable VM mode
servin vm enable

# Test with a simple container
servin run --vm alpine echo "Hello from VM container!"

# Check VM status
servin vm status
```

If all commands complete successfully, your system is ready for VM-based containerization!