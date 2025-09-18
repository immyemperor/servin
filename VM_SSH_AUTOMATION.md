# Servin VM Automation Guide

## SSH Setup Automation in Build Scripts

This guide explains how SSH setup has been automated in the Servin build process to eliminate manual configuration steps.

## Overview

The SSH automation is now fully integrated into the VM creation process, making it seamless to set up a complete Linux container environment on macOS.

## Automated Build Scripts

### 1. Quick Build Script
```bash
./quick-build.sh
```
- Builds Servin binary
- Starts VM with automated SSH setup
- Takes ~30-60 seconds for SSH to be ready

### 2. Full Build Script with VM
```bash
./build.sh vm
```
- Complete build and VM setup with monitoring
- Tests SSH connectivity
- Deploys Servin to VM
- Runs container functionality tests
- Provides detailed status and examples

### 3. Comprehensive Build Script
```bash
./build-with-vm.sh
```
- Enhanced build with full VM automation
- Real-time SSH setup monitoring
- Automatic container testing
- Complete status reporting

## How SSH Automation Works

### 1. Cloud-Init Integration
The VM provider now creates a cloud-init ISO with an automated setup script that:
- Updates Alpine Linux packages
- Installs OpenSSH server
- Creates root password (servin123)
- Configures SSH for root login
- Starts SSH service
- Sets up container environment

### 2. Automatic SSH Detection
The build scripts monitor SSH connectivity:
- Test SSH connection every second
- Wait up to 90 seconds for SSH readiness
- Provide progress updates every 10 seconds
- Automatically proceed when SSH is ready

### 3. Seamless Deployment
Once SSH is ready:
- Automatically deploy Servin binary to VM
- Test basic container functionality
- Provide complete status information

## VM Provider Enhancements

### New Methods Added
- `testSSHConnectivity()`: Tests SSH connection to VM
- `deployServinToVM()`: Automatically deploys Servin binary
- Enhanced `createCloudInitISO()`: Creates comprehensive auto-setup script
- Enhanced `startQEMUVM()`: Monitors SSH and deploys automatically

### Auto-Setup Script Features
The cloud-init ISO contains a script that:
```bash
#!/bin/sh
apk update
apk add openssh sudo curl wget

# Configure SSH
echo 'root:servin123' | chpasswd
echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config
echo 'PasswordAuthentication yes' >> /etc/ssh/sshd_config

# Start SSH service
rc-update add sshd default
rc-service sshd start

# Setup container environment
modprobe overlay
modprobe bridge
echo overlay >> /etc/modules
echo bridge >> /etc/modules

# Success indicator
touch /tmp/ssh-setup-complete
```

## Usage Examples

### Quick Start
```bash
# Build and start VM with automation
./quick-build.sh

# Wait for completion message, then use:
ssh root@localhost -p 2222
# Password: servin123
```

### Full Development Setup
```bash
# Complete build with testing
./build.sh vm

# Or use the comprehensive script
./build-with-vm.sh
```

### Manual Build
```bash
# Traditional build
go build -o servin main.go

# Start VM (now with automation)
./servin vm start

# SSH will be ready automatically in ~30-60 seconds
ssh root@localhost -p 2222
```

## Container Testing

After SSH automation completes, you can immediately test containers:

```bash
# Basic container
./servin run hello-world

# Web server with port mapping
./servin run --name web -p 8080:80 nginx:alpine

# Interactive container
./servin run -it alpine:latest sh

# Container management
./servin list
./servin logs web
./servin exec web sh
./servin stop web
./servin remove web
```

## Troubleshooting

### SSH Setup Takes Too Long
If SSH setup exceeds 90 seconds:
1. Check VM console for boot messages
2. Verify QEMU process is running: `ps aux | grep qemu`
3. Try manual setup:
   ```bash
   # Connect to VM console
   # Login as root (no password)
   mount /dev/sr0 /mnt 2>/dev/null || true
   /mnt/autosetup.sh
   ```

### SSH Connection Fails
If SSH connection fails after setup:
1. Check VM is running: `./servin vm status`
2. Verify port forwarding: `lsof -i :2222`
3. Try manual SSH configuration:
   ```bash
   ssh root@localhost -p 2222
   # If connection refused, manually configure SSH in VM
   ```

### Container Operations Fail
If containers don't work after SSH setup:
1. Verify Servin is deployed: `ssh root@localhost -p 2222 'which servin'`
2. Check kernel modules: `ssh root@localhost -p 2222 'lsmod | grep overlay'`
3. Test basic functionality: `./servin version`

## Architecture Benefits

### Automation Advantages
- **Zero Manual Steps**: Complete VM setup without user intervention
- **Consistent Environment**: Same setup every time
- **Fast Development**: Ready for container testing in minutes
- **Native Linux**: Real Linux namespaces, not Docker simulation

### Technical Implementation
- **QEMU Integration**: Uses Hypervisor.framework for performance
- **Alpine Linux**: Lightweight, fast-booting distribution
- **Cloud-Init**: Industry-standard VM configuration
- **SSH Automation**: Reliable, scriptable access

## Next Steps

1. **Enhanced Testing**: Add more comprehensive container tests
2. **GUI Integration**: Connect WebView GUI to automated VM
3. **Development Workflow**: Streamline development with VM automation
4. **Documentation**: Expand examples and use cases

The SSH automation in build scripts eliminates the complexity of manual VM setup while providing a full Linux container environment for development and testing.