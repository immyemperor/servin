# macOS Support for Servin Container Runtime

## Overview

Servin now supports macOS as a development and testing platform alongside Windows and Linux. While full containerization features require Linux kernel capabilities, macOS support enables:

- ✅ **Image Management**: Import, list, inspect, and remove container images
- ✅ **RootFS Creation**: Create container filesystems from images or basic templates
- ✅ **State Management**: Track container state and metadata
- ✅ **CLI Interface**: Full command-line interface compatibility
- ⚠️ **Limited Isolation**: No namespace/cgroup isolation (development only)

## Platform-Specific Behavior

### Storage Locations

**macOS** uses the following directories:
```bash
~/Library/Application Support/.servin/
├── containers/          # Container state and rootfs
└── images/             # Image storage and metadata
```

**Alternative path** (fallback):
```bash
~/.servin/
├── containers/          # Container state and rootfs  
└── images/             # Image storage and metadata
```

### Root Privileges

macOS requires `sudo` for container operations (similar to Linux):
```bash
# Requires sudo on macOS
sudo servin run alpine:latest /bin/sh

# Development mode bypasses root check
servin run --dev alpine:latest /bin/sh
```

### System Integration Points (SIP)

macOS System Integrity Protection (SIP) affects:
- **chroot**: Not available without disabling SIP
- **Network interfaces**: Limited without root
- **System calls**: Restricted namespace operations

## Usage Examples

### Basic Container Operations

```bash
# Import an image (works on macOS)
servin image import alpine.tar.gz alpine:latest

# List images (works on macOS)  
servin image ls

# Create container with limited isolation (macOS)
sudo servin run alpine:latest /bin/sh
# or in dev mode:
servin run --dev alpine:latest /bin/sh
```

### Development Workflow

```bash
# 1. Import test images
servin image import test-app.tar.gz myapp:dev

# 2. Test container creation (dev mode)
servin run --dev --name test-container myapp:dev /bin/bash

# 3. Inspect container state
servin ls

# 4. Clean up
servin rm --all
servin image rm myapp:dev
```

## Platform Comparison

| Feature | Linux | macOS | Windows |
|---------|-------|-------|---------|
| **Image Management** | ✅ Full | ✅ Full | ✅ Full |
| **RootFS Creation** | ✅ Full | ✅ Full | ✅ Full |
| **Process Isolation** | ✅ Namespaces | ❌ Limited | ❌ Limited |
| **Filesystem Isolation** | ✅ chroot | ⚠️ Simulated | ⚠️ Simulated |
| **Network Isolation** | ✅ netns | ❌ Host only | ❌ Host only |
| **Resource Control** | ✅ cgroups | ❌ None | ❌ None |
| **Root Required** | ✅ Yes | ✅ Yes* | ❌ No |
| **Development Mode** | ✅ Bypass | ✅ Bypass | ✅ Default |

*macOS requires `sudo` but `--dev` flag bypasses this

## Technical Implementation

### Path Resolution
```go
switch runtime.GOOS {
case "darwin":
    // macOS: Use user home directory
    homeDir, _ := os.UserHomeDir()
    path = filepath.Join(homeDir, ".servin", "containers", containerID, "rootfs")
}
```

### Platform Detection
```go
func checkRoot() error {
    switch runtime.GOOS {
    case "darwin":
        fmt.Println("Note: Running on macOS - containerization features limited")
        // Handle macOS-specific requirements
    }
}
```

### Cross-Platform RootFS
```go
func (r *RootFS) Enter() error {
    switch runtime.GOOS {
    case "darwin":
        fmt.Println("macOS - chroot requires root privileges and SIP considerations")
        return nil // Simulate chroot
    }
}
```

## Limitations & Considerations

### Security Limitations
- **No Process Isolation**: Containers run as host processes
- **No Filesystem Isolation**: Simulated chroot only
- **No Network Isolation**: Containers use host network stack
- **No Resource Limits**: No memory/CPU enforcement

### Development Benefits
- **Faster Testing**: Quick image and rootfs testing without VMs
- **Cross-Platform Development**: Develop on macOS, deploy on Linux
- **Image Development**: Build and test image import/export workflows
- **CLI Testing**: Test command-line interface and workflows

### Production Deployment
⚠️ **macOS is for development only** - Deploy to Linux for production:

```bash
# Development on macOS
servin image import myapp.tar.gz myapp:latest --dev
servin run --dev myapp:latest /bin/sh

# Production on Linux
servin run myapp:latest /bin/sh  # Full isolation
```

## Installation & Setup

### Prerequisites
- macOS 10.15+ (Catalina or later)
- Go 1.19+ for building from source
- `sudo` access for container operations

### Build from Source
```bash
# Clone and build
git clone <servin-repo>
cd servin
go build -o servin .

# Install system-wide (optional)
sudo mv servin /usr/local/bin/
```

### Quick Start
```bash
# Test basic functionality
servin image ls --dev
servin run --dev ubuntu:latest echo "Hello macOS"
```

## Troubleshooting

### Common Issues

1. **Permission Denied**
   ```
   Error: this command requires root privileges on macOS
   Solution: Use sudo or --dev flag
   ```

2. **chroot Not Available**
   ```
   Note: macOS - chroot requires root privileges and SIP considerations
   Solution: Expected behavior - containers use host filesystem
   ```

3. **SIP Restrictions**
   ```
   Error: operation not permitted
   Solution: Use --dev mode for development testing
   ```

### macOS-Specific Tips

1. **Use Development Mode**: `--dev` flag bypasses most restrictions
2. **Homebrew Compatibility**: Works with Homebrew-installed dependencies
3. **File Permissions**: Ensure proper permissions on ~/.servin directory
4. **Network Access**: Containers inherit host network configuration

## Future Enhancements

### Planned macOS Improvements
- **launchd Integration**: macOS-native service management
- **Keychain Integration**: Secure credential storage
- **Notification Support**: macOS notification for container events
- **Finder Integration**: GUI tools for container management

### Advanced Features (Future)
- **Docker Desktop Compatibility**: Import Docker images
- **Kubernetes Support**: Local k8s development
- **Native App Bundling**: Package containers as .app bundles

This macOS support enables effective development and testing workflows while maintaining full compatibility with Linux production environments.
