# Servin Container Runtime - Complete Package System

This directory contains comprehensive installer packages for Servin Container Runtime that include all VM dependencies and prerequisites embedded within the installers themselves.

## 🎯 Overview

Each platform-specific installer is a **complete, self-contained package** that:
- Embeds all VM dependencies (QEMU, Hyper-V, VirtualBox, etc.)
- Includes automatic system requirement checking
- Provides native platform integration
- Requires no internet connection during installation
- Works offline with all dependencies included

## 📦 Package Types

### Windows - NSIS Installer
- **File**: `Servin-Installer-1.0.0.exe`
- **Format**: NSIS executable installer
- **Size**: ~200-300MB (includes VM providers)
- **Features**:
  - VM Provider Selection (Hyper-V/VirtualBox/WSL2/All)
  - System Requirements Validation
  - Automatic Chocolatey/Python Installation
  - Windows Service Configuration
  - Desktop Integration (shortcuts, file associations)
  - Complete Uninstaller

### Linux - AppImage
- **File**: `Servin-1.0.0-x86_64.AppImage`
- **Format**: Portable AppImage executable
- **Size**: ~150-250MB (includes QEMU/KVM)
- **Features**:
  - Portable, no installation required
  - Embedded QEMU/KVM dependencies
  - Desktop integration when installed
  - Automatic VM prerequisite checking
  - Enhanced installer integration

### macOS - PKG Package
- **File**: `Servin-1.0.0-arm64.pkg` / `Servin-1.0.0-amd64.pkg`
- **Format**: Native macOS installer package
- **Size**: ~100-200MB (includes QEMU)
- **Features**:
  - Native macOS .app bundle
  - Embedded QEMU distribution
  - Automatic Homebrew integration
  - LaunchDaemons service setup
  - Virtualization.framework support

## 🚀 Quick Start

### Build All Packages
```bash
# Build packages for all platforms
./build-packages.sh

# Build specific platform only
./build-packages.sh --windows
./build-packages.sh --linux
./build-packages.sh --macos
```

### Install Platform-Specific Packages

#### Windows
```cmd
# Run as Administrator
Servin-Installer-1.0.0.exe
```

#### Linux
```bash
# Make executable and run
chmod +x Servin-1.0.0-x86_64.AppImage
./Servin-1.0.0-x86_64.AppImage

# Or install system-wide
./install-servin-appimage.sh
```

#### macOS
```bash
# Install package
sudo installer -pkg Servin-1.0.0-arm64.pkg -target /

# Or double-click the .pkg file
```

## 📁 Directory Structure

```
installers/
├── build-packages.sh           # Cross-platform build coordinator
├── PACKAGE_README.md          # This file
│
├── windows/                   # Windows NSIS Installer
│   ├── build-installer.bat    # Windows build script
│   ├── servin-installer.nsi   # NSIS installer script (500+ lines)
│   ├── EnvVarUpdate.nsh      # Environment variable management
│   └── Dockerfile.nsis       # Docker build environment
│
├── linux/                    # Linux AppImage
│   ├── build-appimage.sh     # AppImage build script (400+ lines)
│   ├── install-with-vm.sh    # Enhanced installer integration
│   └── Dockerfile.appimage   # Docker build environment
│
├── macos/                     # macOS Package
│   ├── build-package.sh      # PKG build script (500+ lines)
│   ├── install-with-vm.sh    # Enhanced installer integration
│   └── Servin.app structure  # Native app bundle
│
└── dist/                      # Distribution packages
    └── servin-1.0.0-complete.tar.gz
```

## 🔧 Build Requirements

### Prerequisites
- **Go 1.19+** for cross-compilation
- **Docker** (optional, for cross-platform builds)
- Platform-specific tools (see below)

### Platform-Specific Tools

#### Windows
- **NSIS 3.0+** (Nullsoft Scriptable Install System)
- **makensis** command available in PATH
- **Wine** (for cross-platform building on Linux/macOS)

#### Linux
- **AppImage tools** (downloaded automatically)
- **linuxdeploy** and **appimagetool**
- **QEMU development files** (for bundling)

#### macOS
- **Xcode Command Line Tools**
- **pkgbuild** and **productbuild**
- **hdiutil** (for DMG creation)
- **iconutil** (for icon generation)

## 🎨 Package Features

### Windows NSIS Installer Features
```nsis
# VM Provider Selection Dialog
!include "MUI2.nsh"
!insertmacro MUI_PAGE_COMPONENTS  # VM provider selection

# System Requirements Checking
!define MIN_WIN_VERSION "10.0"
!define MIN_RAM_GB 4
!define MIN_DISK_GB 10

# Automatic Dependencies
- Chocolatey installation
- Python 3.x with required packages
- VM provider installation (Hyper-V/VirtualBox/WSL2)
- Windows Service creation
```

### Linux AppImage Features
```bash
# Portable VM Environment
export APPDIR="$(dirname "$(readlink -f "${0}")")"
export PATH="$APPDIR/usr/bin:$APPDIR/opt/servin/bin:$PATH"

# Embedded Dependencies
- QEMU/KVM binaries and libraries
- Python environment for GUI
- Desktop integration files
- Automatic VM prerequisite checking
```

### macOS Package Features
```xml
<!-- Native App Bundle -->
<key>CFBundleIdentifier</key>
<string>com.servin.containerruntime</string>

<!-- VM Integration -->
- Embedded QEMU distribution
- Virtualization.framework support
- Automatic Homebrew installation
- LaunchDaemons service setup
```

## 🔍 Testing & Validation

### Automated Testing
```bash
# Test build process
./build-packages.sh --help

# Validate package integrity
cd installers/dist/servin-1.0.0-*/
sha256sum -c checksums.txt
```

### Manual Testing
1. **Windows**: Run installer in clean VM, verify all features
2. **Linux**: Test AppImage on different distributions
3. **macOS**: Verify on both Intel and Apple Silicon

## 📊 Package Sizes & Dependencies

| Platform | Base Size | With VM Deps | Total Features |
|----------|-----------|--------------|----------------|
| Windows  | ~50MB     | ~250MB       | Complete VM stack |
| Linux    | ~30MB     | ~200MB       | QEMU/KVM embedded |
| macOS    | ~40MB     | ~150MB       | QEMU + VZ framework |

## 🔐 Security & Signing

### Code Signing (Future Enhancement)
```bash
# Windows: signtool.exe
signtool sign /f certificate.p12 /p password Servin-Installer.exe

# macOS: codesign
codesign --sign "Developer ID" Servin.app
productbuild --sign "Developer ID Installer" ...

# Linux: GPG signing
gpg --detach-sign --armor Servin.AppImage
```

## 🚢 Distribution Strategy

### Release Process
1. **Build**: Run `./build-packages.sh` on each platform
2. **Test**: Validate installers on clean systems
3. **Package**: Create unified distribution archive
4. **Upload**: Release to GitHub/distribution channels

### Download Locations
- **Primary**: GitHub Releases
- **Mirror**: Direct download links
- **Packages**: Platform-specific repositories

## 🔧 Troubleshooting

### Common Build Issues

#### Windows
```bash
# NSIS not found
choco install nsis

# Build fails on non-Windows
docker run --rm -v $(pwd):/build servin-nsis-builder
```

#### Linux
```bash
# AppImage tools missing
wget https://github.com/linuxdeploy/linuxdeploy/releases/download/continuous/linuxdeploy-x86_64.AppImage

# Docker build alternative
docker run --rm -v $(pwd):/build servin-appimage-builder
```

#### macOS
```bash
# Xcode tools missing
xcode-select --install

# Package build fails
sudo xcode-select --reset
```

### Installation Issues

#### Windows
- **Antivirus blocking**: Add installer to exclusions
- **Permission denied**: Run as Administrator
- **VM provider conflicts**: Choose single provider in installer

#### Linux
- **AppImage won't run**: Check FUSE availability
- **Missing libraries**: Use Docker build for better compatibility
- **VM access denied**: Add user to kvm group

#### macOS
- **Gatekeeper blocking**: `sudo spctl --master-disable`
- **Notarization required**: Use signed packages for distribution
- **VM framework missing**: Requires macOS 11+ for full features

## 📈 Future Enhancements

### Planned Features
- [ ] **Auto-updater integration**
- [ ] **Digital signatures for all platforms**
- [ ] **Chocolatey/Homebrew/Snap packages**
- [ ] **Container image pre-loading**
- [ ] **Enterprise deployment scripts**
- [ ] **Silent installation modes**

### Build System Improvements
- [ ] **CI/CD integration for automated builds**
- [ ] **Cross-compilation optimization**
- [ ] **Incremental build support**
- [ ] **Package size optimization**
- [ ] **Dependency caching**

## 📝 Contributing

To contribute to the package system:

1. **Test on target platforms**
2. **Submit platform-specific improvements**
3. **Report packaging issues**
4. **Enhance build automation**

## 📄 License

All installer packages inherit the Servin Container Runtime license (MIT).

---

**Built with ❤️ for seamless cross-platform container runtime deployment**