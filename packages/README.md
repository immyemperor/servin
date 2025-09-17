# Package Manager Integration for Servin VM Distribution

This directory contains configuration files and scripts for integrating Servin with various package managers across different platforms.

## Supported Package Managers

### macOS
- **Homebrew**: Formula for easy installation via `brew install servin`
- **MacPorts**: Portfile for MacPorts users

### Windows  
- **Winget**: Microsoft's package manager
- **Chocolatey**: Community package manager
- **Scoop**: User-space package manager

### Linux
- **APT**: Debian/Ubuntu packages (.deb)
- **YUM/DNF**: Red Hat/Fedora packages (.rpm) 
- **Snap**: Universal Linux packages
- **Flatpak**: Sandboxed application distribution
- **AUR**: Arch User Repository

### Universal
- **Docker Hub**: Container images for Servin itself
- **GitHub Releases**: Direct binary downloads

## Distribution Strategy

### Phase 1: Core Packages
1. **Homebrew** (macOS primary)
2. **APT Repository** (Ubuntu/Debian)
3. **GitHub Releases** (manual downloads)

### Phase 2: Extended Support
1. **Winget** (Windows primary)
2. **RPM Repository** (RHEL/Fedora)
3. **Snap Store** (universal Linux)

### Phase 3: Community Packages
1. **AUR** (Arch Linux)
2. **Chocolatey** (Windows community)
3. **Flatpak** (sandboxed Linux)

## Package Contents

Each package includes:
- **Servin binary** with VM support
- **Desktop GUI** (where applicable)
- **WebView GUI** resources
- **VM configuration** templates
- **Documentation** and examples
- **Install/uninstall scripts**

## VM Image Distribution

VM images are distributed separately via:
- **CDN**: Fast global delivery
- **Checksums**: Integrity verification
- **Incremental updates**: Delta downloads
- **Multiple distros**: Alpine, Ubuntu, Debian options

## Usage Examples

### macOS
```bash
# Install via Homebrew
brew tap servin/tap
brew install servin

# Enable VM mode
servin vm enable
servin vm start
```

### Windows
```powershell
# Install via Winget
winget install Servin.Servin

# Enable VM mode
servin vm enable
servin vm start
```

### Linux
```bash
# Ubuntu/Debian
sudo apt install servin

# Fedora/RHEL
sudo dnf install servin

# Arch Linux (AUR)
yay -S servin

# Universal (Snap)
sudo snap install servin --classic

# Manual
curl -fsSL https://get.servin.dev | sh
```

## Repository Structure

```
packages/
├── homebrew/
│   ├── servin.rb              # Homebrew formula
│   └── update-formula.sh      # Auto-update script
├── winget/
│   ├── servin.yaml           # Winget manifest
│   └── update-manifest.py    # Auto-update script
├── debian/
│   ├── control              # Package metadata
│   ├── postinst             # Post-install script
│   ├── prerm                # Pre-removal script
│   └── build-deb.sh         # Build script
├── rpm/
│   ├── servin.spec          # RPM specification
│   └── build-rpm.sh         # Build script
├── snap/
│   ├── snapcraft.yaml       # Snap configuration
│   └── build-snap.sh        # Build script
└── docker/
    ├── Dockerfile           # Container image
    └── build-image.sh       # Build script
```

## Automation

### CI/CD Integration
- **GitHub Actions**: Automated package building
- **Release Triggers**: Automatic updates on new releases
- **Testing**: Package installation verification
- **Signing**: Code signing for security

### Update Mechanism
- **Version Detection**: Automatic version checking
- **Delta Updates**: Efficient incremental updates
- **Rollback**: Safe rollback on update failures
- **Notifications**: User notification system

## Security

### Package Signing
- **GPG Signatures**: All packages signed with GPG
- **Checksums**: SHA-256 verification
- **Repository Keys**: Secure repository authentication

### VM Images
- **Verified Sources**: Only official distribution sources
- **Security Updates**: Regular security patching
- **Minimal Attack Surface**: Stripped-down VM images

## Distribution Channels

### Official Channels
1. **GitHub Releases**: Primary source
2. **Official Website**: https://servin.dev
3. **Package Repositories**: Platform-specific repos

### Community Channels  
1. **AUR**: Community-maintained Arch packages
2. **Chocolatey**: Community Windows packages
3. **Third-party RPMs**: Community RHEL/CentOS packages

## Quality Assurance

### Testing Matrix
- **Platform Testing**: All supported OS/arch combinations
- **Installation Testing**: Package manager installation
- **Functionality Testing**: VM containerization verification
- **Performance Testing**: VM overhead benchmarks

### Compatibility
- **Backward Compatibility**: Existing container support
- **Forward Compatibility**: Upgrade path planning
- **Cross-Platform**: Identical behavior verification

## Metrics and Analytics

### Download Tracking
- **Package Downloads**: Per-platform statistics
- **Version Adoption**: Update rate monitoring
- **Geographic Distribution**: Global usage patterns

### Performance Monitoring
- **Installation Success Rate**: Package installation metrics
- **VM Performance**: Containerization performance data
- **Error Reporting**: Anonymous error collection

## Support and Documentation

### User Support
- **Installation Guides**: Platform-specific instructions
- **Troubleshooting**: Common issue resolution
- **Community Forums**: User community support

### Developer Support
- **API Documentation**: VM integration APIs
- **Contributing Guidelines**: Package maintainer guides
- **Release Process**: How releases are managed