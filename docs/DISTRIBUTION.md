# Servin Distribution and Packaging Documentation

## Overview

This document describes the comprehensive distribution strategy for Servin, a universal container runtime with VM-based containerization support. The distribution system handles cross-platform packaging, installation, and deployment across multiple channels.

## Architecture

### Distribution Components

```
Servin Distribution System
├── Core Binaries
│   ├── servin (main CLI)
│   ├── servin-desktop (GUI)
│   └── servin-daemon (background service)
├── VM Images
│   ├── Lightweight Linux VM (Ubuntu/Alpine)
│   ├── Container runtime (Docker/containerd)
│   └── Network/storage drivers
├── Installation System
│   ├── Universal installer script
│   ├── Platform-specific packages
│   └── Package manager integrations
└── Distribution Channels
    ├── Direct downloads
    ├── Package managers
    ├── Container registries
    └── Cloud marketplaces
```

### Target Platforms

#### Primary Platforms
- **Linux**: x86_64, ARM64 (Ubuntu, Debian, CentOS, Fedora, SUSE)
- **macOS**: x86_64, ARM64 (macOS 11+)
- **Windows**: x86_64, ARM64 (Windows 10+)

#### Package Formats
- **Linux**: DEB, RPM, Snap, AppImage, Flatpak
- **macOS**: PKG, Homebrew, MacPorts
- **Windows**: MSI, Winget, Chocolatey, Scoop

## Build System

### Cross-Platform Build Process

The build system is orchestrated by `build-vm-distribution.sh`:

```bash
# Complete build for all platforms
./build-vm-distribution.sh --all

# Platform-specific builds
./build-vm-distribution.sh --platform linux --arch amd64
./build-vm-distribution.sh --platform darwin --arch arm64
./build-vm-distribution.sh --platform windows --arch amd64
```

#### Build Stages

1. **Environment Setup**
   - Platform detection
   - Dependency verification
   - Build tool preparation

2. **Binary Compilation**
   - Cross-compilation for target platforms
   - Static linking for portability
   - Optimization for size and performance

3. **VM Image Creation**
   - Lightweight Linux VM building
   - Container runtime integration
   - Image optimization and compression

4. **Package Creation**
   - Platform-specific package formats
   - Metadata and dependency management
   - Digital signing and verification

5. **Installer Generation**
   - Universal installation scripts
   - Platform-specific installers
   - Automated testing and validation

### Build Artifacts

```
dist/
├── binaries/
│   ├── linux-amd64/
│   │   ├── servin
│   │   ├── servin-desktop
│   │   └── servin-daemon
│   ├── linux-arm64/
│   ├── darwin-amd64/
│   ├── darwin-arm64/
│   └── windows-amd64/
├── vm-images/
│   ├── servin-vm-linux-amd64.qcow2
│   ├── servin-vm-linux-arm64.qcow2
│   └── checksums.txt
├── packages/
│   ├── deb/
│   │   ├── servin_1.0.0_amd64.deb
│   │   └── servin_1.0.0_arm64.deb
│   ├── rpm/
│   │   ├── servin-1.0.0-1.x86_64.rpm
│   │   └── servin-1.0.0-1.aarch64.rpm
│   ├── snap/
│   │   └── servin_1.0.0_amd64.snap
│   ├── homebrew/
│   │   └── servin.rb
│   └── winget/
│       └── servin.yaml
├── installers/
│   ├── install.sh (universal)
│   ├── servin-linux-installer.sh
│   ├── servin-macos-installer.pkg
│   └── servin-windows-installer.exe
└── documentation/
    ├── README.md
    ├── INSTALL.md
    └── CHANGELOG.md
```

## Installation Methods

### Universal Installer

The universal installer script (`scripts/install.sh`) provides a one-command installation:

```bash
curl -fsSL https://get.servin.dev | sh
```

Features:
- Automatic platform detection
- Dependency checking and installation
- VM capability verification
- PATH configuration
- Post-install verification

### Package Manager Integrations

#### Homebrew (macOS/Linux)
```bash
brew tap immyemperor/servin
brew install servin
```

#### APT (Debian/Ubuntu)
```bash
curl -fsSL https://packages.servin.dev/gpg | sudo apt-key add -
echo "deb https://packages.servin.dev/apt stable main" | sudo tee /etc/apt/sources.list.d/servin.list
sudo apt update && sudo apt install servin
```

#### RPM (RedHat/CentOS/SUSE)
```bash
sudo rpm --import https://packages.servin.dev/gpg
sudo yum-config-manager --add-repo https://packages.servin.dev/yum/servin.repo
sudo yum install servin
```

#### Winget (Windows)
```bash
winget install servin
```

#### Snap (Linux)
```bash
sudo snap install servin
```

### Direct Downloads

Binary releases are available for direct download:
- GitHub Releases: `https://github.com/immyemperor/servin/releases`
- Official Site: `https://servin.dev/download`

## VM Image Distribution

### VM Image Strategy

Servin uses lightweight Linux VMs for true containerization across all platforms:

#### Image Characteristics
- **Base OS**: Ubuntu 22.04 LTS or Alpine Linux
- **Size**: ~100MB compressed
- **Boot Time**: <5 seconds
- **Memory**: 512MB default, configurable
- **Storage**: Dynamic allocation

#### Distribution Channels
- **Built-in**: Embedded in installation packages
- **Download**: On-demand from CDN
- **Container Registry**: OCI-compatible images
- **Local Cache**: Cached after first download

#### Image Management
```bash
# List available VM images
servin vm images list

# Download specific image
servin vm images pull ubuntu-22.04-amd64

# Update to latest
servin vm images update

# Clean old images
servin vm images prune
```

## Release Process

### Automated Release Pipeline

The release process is managed by `scripts/create-release.sh`:

```bash
# Create release with full automation
./scripts/create-release.sh v1.0.0

# Create pre-release
./scripts/create-release.sh v1.0.0-beta.1 --prerelease

# Create draft release
./scripts/create-release.sh v1.0.0 --draft
```

#### Release Steps
1. **Pre-release validation**
   - Test suite execution
   - Security scanning
   - Performance benchmarks

2. **Build generation**
   - Cross-platform compilation
   - Package creation
   - Installer generation

3. **Release creation**
   - Git tag creation
   - GitHub release
   - Artifact upload

4. **Distribution updates**
   - Package manager updates
   - Website deployment
   - Documentation refresh

### Release Channels

#### Stable Channel
- Thoroughly tested releases
- Production-ready
- Quarterly release cycle

#### Beta Channel
- Feature-complete pre-releases
- Community testing
- Monthly release cycle

#### Nightly Channel
- Daily builds from main branch
- Latest features
- Development testing

## Package Manager Integration

### Homebrew Integration

#### Formula Structure
```ruby
class Servin < Formula
  desc "Universal container runtime with VM-based containerization"
  homepage "https://servin.dev"
  url "https://github.com/immyemperor/servin/archive/v1.0.0.tar.gz"
  sha256 "abc123..."
  license "MIT"

  depends_on "go" => :build
  depends_on "qemu"

  def install
    system "make", "build"
    bin.install "servin"
    bin.install "servin-desktop"
  end

  test do
    assert_match "Servin", shell_output("#{bin}/servin --version")
  end
end
```

#### Update Process
```bash
# Automated formula update
./scripts/update-homebrew.sh v1.0.0
```

### APT Repository Management

#### Repository Structure
```
packages.servin.dev/apt/
├── dists/
│   └── stable/
│       ├── main/
│       │   └── binary-amd64/
│       └── Release
├── pool/
│   └── main/
│       └── s/
│           └── servin/
└── gpg/
    └── servin-archive-keyring.gpg
```

#### Package Metadata
```control
Package: servin
Version: 1.0.0
Architecture: amd64
Maintainer: Servin Team <team@servin.dev>
Depends: libc6, qemu-system-x86
Description: Universal container runtime
 Servin provides VM-based containerization for true isolation
 across all platforms including macOS and Windows.
```

### Windows Package Management

#### Winget Manifest
```yaml
PackageIdentifier: Servin.Servin
PackageVersion: 1.0.0
PackageName: Servin
Publisher: Servin
License: MIT
ShortDescription: Universal container runtime
Installers:
- Architecture: x64
  InstallerType: msi
  InstallerUrl: https://github.com/immyemperor/servin/releases/download/v1.0.0/servin-1.0.0-windows-amd64.msi
  InstallerSha256: abc123...
```

## Quality Assurance

### Testing Strategy

#### Automated Testing
- Unit tests for core functionality
- Integration tests for VM operations
- End-to-end installation testing
- Cross-platform compatibility testing

#### Manual Testing
- Installation process validation
- GUI functionality verification
- Container operation testing
- Performance benchmarking

### Security Measures

#### Code Signing
- macOS: Apple Developer certificate
- Windows: Code signing certificate
- Linux: GPG signing for packages

#### Supply Chain Security
- Dependency scanning
- Vulnerability assessments
- Reproducible builds
- SBOM generation

## Monitoring and Analytics

### Distribution Metrics

#### Download Statistics
- Total downloads by platform
- Geographic distribution
- Version adoption rates
- Installation success rates

#### User Analytics
- Active user counts
- Feature usage patterns
- Error reporting and crash analytics
- Performance metrics

### Monitoring Infrastructure

#### Services
- Download CDN monitoring
- Package repository health
- Installation success tracking
- User feedback collection

## Troubleshooting

### Common Installation Issues

#### Platform-Specific Issues
- **macOS**: Gatekeeper blocking unsigned binaries
- **Windows**: Windows Defender SmartScreen warnings
- **Linux**: Missing virtualization support

#### VM Issues
- Hardware virtualization not enabled
- Insufficient system resources
- Network connectivity problems
- Storage space limitations

### Support Channels

#### Documentation
- Installation guide: `https://servin.dev/docs/install`
- Troubleshooting: `https://servin.dev/docs/troubleshooting`
- FAQ: `https://servin.dev/docs/faq`

#### Community Support
- GitHub Issues: Bug reports and feature requests
- Discussions: Community Q&A
- Discord: Real-time support chat

## Future Enhancements

### Planned Improvements

#### Distribution
- ARM64 Windows support
- Additional Linux distributions
- Container registry integration
- Cloud marketplace listings

#### Technology
- Improved VM image compression
- Faster installation process
- Enhanced package validation
- Automated rollback capabilities

#### User Experience
- GUI installer improvements
- Better error messages
- Progress indicators
- Installation customization options

## References

### Documentation
- [Installation Guide](INSTALL.md)
- [Release Process](RELEASE.md)
- [VM Containerization](VM_CONTAINERIZATION.md)

### Scripts and Tools
- [Universal Installer](../scripts/install.sh)
- [Release Creator](../scripts/create-release.sh)
- [Distribution Builder](../build-vm-distribution.sh)