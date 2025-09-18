#!/bin/bash
# Enhanced cross-platform build script for Servin with VM distribution

set -e

# Configuration
BUILD_DIR="dist"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date +%Y-%m-%dT%H%M%S)

# VM Configuration
VM_IMAGES_DIR="vm-images"
VM_VERSION="1.0.0"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_header() {
    echo -e "${CYAN}================================================${NC}"
    echo -e "${CYAN}   Servin Distribution Build System${NC}"
    echo -e "${CYAN}   Universal VM Containerization${NC}"
    echo -e "${CYAN}================================================${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_info() {
    echo -e "${BLUE}→ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# Enhanced platform definitions with VM support
PLATFORMS=(
    "linux/amd64:kvm,qemu"
    "linux/arm64:kvm,qemu" 
    "darwin/amd64:hvf,qemu"
    "darwin/arm64:hvf,qemu"
    "windows/amd64:hyperv,vbox"
)

# VM image configurations
VM_CONFIGS=(
    "alpine:3.18:docker:512M:2G"
    "ubuntu:22.04:containerd:1G:4G"
    "debian:12:podman:512M:2G"
)

create_build_structure() {
    print_info "Creating distribution structure..."
    
    rm -rf "$BUILD_DIR"
    mkdir -p "$BUILD_DIR"/{binaries,packages,installers,vm-images,docs}
    mkdir -p "$BUILD_DIR"/vm-images/{alpine,ubuntu,debian}
    mkdir -p "$BUILD_DIR"/installers/{linux,macos,windows}
    mkdir -p "$BUILD_DIR"/packages/{deb,rpm,homebrew,winget}
    
    print_success "Distribution structure created"
}

build_binaries() {
    print_info "Building cross-platform binaries..."
    
    for platform_info in "${PLATFORMS[@]}"; do
        IFS=':' read -r platform vm_support <<< "$platform_info"
        IFS='/' read -r os arch <<< "$platform"
        
        binary_name="servin"
        if [ "$os" = "windows" ]; then
            binary_name="servin.exe"
        fi
        
        output_dir="$BUILD_DIR/binaries/$os-$arch"
        mkdir -p "$output_dir"
        
        print_info "Building for $os/$arch with VM support: $vm_support"
        
        # Build main binary with VM tags
        GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build \
            -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -s -w" \
            -tags "vm_enabled,$vm_support" \
            -o "$output_dir/$binary_name" .
        
        # Build desktop GUI
        if [ "$os" != "linux" ] || [ "$arch" = "amd64" ]; then
            GOOS=$os GOARCH=$arch CGO_ENABLED=1 go build \
                -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -s -w" \
                -tags "desktop,vm_enabled" \
                -o "$output_dir/servin-desktop$([ "$os" = "windows" ] && echo ".exe")" \
                ./cmd/servin-desktop
        fi
        
        # Copy WebView GUI
        cp -r webview_gui "$output_dir/"
        
        # Create platform-specific VM configs
        create_vm_config "$os" "$arch" "$vm_support" "$output_dir"
        
        print_success "Built $os/$arch"
    done
}

create_vm_config() {
    local os=$1
    local arch=$2
    local vm_support=$3
    local output_dir=$4
    
    cat > "$output_dir/vm-config.json" << EOF
{
    "platform": "$os",
    "architecture": "$arch", 
    "vm_providers": [$(echo "$vm_support" | sed 's/,/", "/g' | sed 's/^/"/;s/$/"/')],
    "default_config": {
        "name": "servin-vm",
        "cpus": 2,
        "memory_mb": 2048,
        "disk_size_gb": 20,
        "linux_distro": "alpine",
        "container_runtime": "docker",
        "ssh_port": 2222,
        "docker_port": 2375
    },
    "vm_images": {
        "alpine": {
            "version": "3.18",
            "url": "https://releases.servin.dev/vm-images/alpine-$VM_VERSION.qcow2",
            "checksum": "sha256:..."
        },
        "ubuntu": {
            "version": "22.04", 
            "url": "https://releases.servin.dev/vm-images/ubuntu-$VM_VERSION.qcow2",
            "checksum": "sha256:..."
        }
    }
}
EOF
}

build_vm_images() {
    print_info "Building VM images..."
    
    for vm_config in "${VM_CONFIGS[@]}"; do
        IFS=':' read -r distro version runtime memory disk <<< "$vm_config"
        
        print_info "Building $distro $version VM image..."
        
        # Create cloud-init configuration
        create_cloud_init "$distro" "$version" "$runtime" "$memory" "$disk"
        
        # Build VM image using Packer (if available) or custom script
        if command -v packer >/dev/null 2>&1; then
            build_with_packer "$distro" "$version" "$runtime"
        else
            build_custom_vm "$distro" "$version" "$runtime"
        fi
        
        print_success "Built $distro VM image"
    done
}

create_cloud_init() {
    local distro=$1
    local version=$2
    local runtime=$3
    local memory=$4
    local disk=$5
    
    local vm_dir="$BUILD_DIR/vm-images/$distro"
    
    # Create cloud-init user-data
    cat > "$vm_dir/user-data" << EOF
#cloud-config
hostname: servin-vm
users:
  - name: servin
    sudo: ALL=(ALL) NOPASSWD:ALL
    shell: /bin/bash
    ssh_authorized_keys:
      - ssh-rsa AAAAB3NzaC1yc2E... # Servin VM key

package_update: true
package_upgrade: true

packages:
  - curl
  - wget
  - openssh-server
  - ca-certificates
  - gnupg
  - lsb-release

runcmd:
  # Install container runtime
  - curl -fsSL https://get.docker.com | sh
  - usermod -aG docker servin
  - systemctl enable docker
  - systemctl start docker
  
  # Configure SSH
  - systemctl enable ssh
  - systemctl start ssh
  
  # Install container tools
  - docker pull hello-world
  - docker pull alpine:latest
  - docker pull nginx:alpine
  
  # Optimize for container workloads
  - echo 'vm.swappiness=1' >> /etc/sysctl.conf
  - echo 'net.ipv4.ip_forward=1' >> /etc/sysctl.conf
  
  # Signal VM is ready
  - touch /var/lib/cloud/servin-vm-ready

final_message: "Servin VM is ready for containers!"
EOF

    # Create meta-data
    cat > "$vm_dir/meta-data" << EOF
instance-id: servin-vm-$version
local-hostname: servin-vm
EOF

    # Create network config
    cat > "$vm_dir/network-config" << EOF
version: 2
ethernets:
  eth0:
    dhcp4: true
EOF
}

build_with_packer() {
    local distro=$1
    local version=$2
    local runtime=$3
    
    # Create Packer template for automated VM building
    cat > "$BUILD_DIR/vm-images/$distro/packer.json" << EOF
{
    "builders": [
        {
            "type": "qemu",
            "iso_url": "$(get_iso_url $distro $version)",
            "iso_checksum": "$(get_iso_checksum $distro $version)",
            "output_directory": "$BUILD_DIR/vm-images/$distro",
            "vm_name": "servin-$distro-$VM_VERSION",
            "disk_size": "4096M",
            "memory": 2048,
            "cpus": 2,
            "accelerator": "kvm",
            "ssh_username": "servin",
            "ssh_timeout": "20m",
            "shutdown_command": "sudo shutdown -P now",
            "cd_files": [
                "$BUILD_DIR/vm-images/$distro/user-data",
                "$BUILD_DIR/vm-images/$distro/meta-data",
                "$BUILD_DIR/vm-images/$distro/network-config"
            ]
        }
    ],
    "provisioners": [
        {
            "type": "shell",
            "inline": [
                "while [ ! -f /var/lib/cloud/servin-vm-ready ]; do sleep 5; done",
                "sudo cloud-init clean",
                "sudo rm -rf /var/lib/cloud/instances/*"
            ]
        }
    ]
}
EOF

    # Build with Packer
    cd "$BUILD_DIR/vm-images/$distro" && packer build packer.json
}

get_iso_url() {
    local distro=$1
    local version=$2
    
    case $distro in
        alpine)
            echo "https://dl-cdn.alpinelinux.org/alpine/v${version}/releases/x86_64/alpine-virt-${version}.0-x86_64.iso"
            ;;
        ubuntu)
            echo "https://releases.ubuntu.com/${version}/ubuntu-${version}-live-server-amd64.iso"
            ;;
        debian)
            echo "https://cdimage.debian.org/debian-cd/current/amd64/iso-cd/debian-${version}.0-amd64-netinst.iso"
            ;;
    esac
}

create_installers() {
    print_info "Creating platform installers..."
    
    # macOS installer
    create_macos_installer
    
    # Windows installer  
    create_windows_installer
    
    # Linux packages
    create_linux_packages
    
    print_success "Platform installers created"
}

create_macos_installer() {
    local installer_dir="$BUILD_DIR/installers/macos"
    
    # Create Homebrew formula
    cat > "$installer_dir/servin.rb" << 'EOF'
class Servin < Formula
  desc "Universal container runtime with VM-based true containerization"
  homepage "https://github.com/immyemperor/servin"
  url "https://github.com/immyemperor/servin/releases/download/v${VERSION}/servin-darwin-universal.tar.gz"
  sha256 "..."
  license "MIT"

  depends_on "qemu"
  
  def install
    bin.install "servin"
    bin.install "servin-desktop"
    share.install "webview_gui"
    etc.install "vm-config.json" => "servin/vm-config.json"
  end

  service do
    run [opt_bin/"servin", "vm", "start"]
    keep_alive false
    log_path var/"log/servin.log"
    error_log_path var/"log/servin.log"
  end

  test do
    system "#{bin}/servin", "--version"
  end
end
EOF

    # Create .pkg installer script
    cat > "$installer_dir/build-pkg.sh" << 'EOF'
#!/bin/bash
# Build macOS .pkg installer

PACKAGE_NAME="Servin"
PACKAGE_VERSION="${VERSION}"
PACKAGE_ID="dev.servin.servin"

# Create package structure
mkdir -p pkg-root/usr/local/bin
mkdir -p pkg-root/usr/local/share/servin
mkdir -p pkg-root/etc/servin

# Copy binaries and resources
cp ../binaries/darwin-*/servin pkg-root/usr/local/bin/
cp ../binaries/darwin-*/servin-desktop pkg-root/usr/local/bin/
cp -r ../binaries/darwin-*/webview_gui pkg-root/usr/local/share/servin/
cp ../binaries/darwin-*/vm-config.json pkg-root/etc/servin/

# Build package
pkgbuild --root pkg-root \
         --identifier "$PACKAGE_ID" \
         --version "$PACKAGE_VERSION" \
         --install-location / \
         --scripts scripts \
         "${PACKAGE_NAME}-${PACKAGE_VERSION}.pkg"
EOF

    chmod +x "$installer_dir/build-pkg.sh"
}

create_windows_installer() {
    local installer_dir="$BUILD_DIR/installers/windows"
    
    # Create NSIS installer script
    cat > "$installer_dir/servin-installer.nsi" << 'EOF'
!define PRODUCT_NAME "Servin"
!define PRODUCT_VERSION "${VERSION}"
!define PRODUCT_PUBLISHER "Servin Project"
!define PRODUCT_WEB_SITE "https://github.com/immyemperor/servin"

SetCompressor lzma

!include "MUI2.nsh"

!define MUI_ABORTWARNING
!define MUI_ICON "..\..\icons\servin.ico"

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "..\..\LICENSE"
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_COMPONENTS
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_WELCOME
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_UNPAGE_FINISH

!insertmacro MUI_LANGUAGE "English"

Name "${PRODUCT_NAME} ${PRODUCT_VERSION}"
OutFile "Servin-${PRODUCT_VERSION}-Setup.exe"
InstallDir "$PROGRAMFILES64\Servin"
ShowInstDetails show
ShowUnInstDetails show

Section "Core Runtime" SEC01
  SectionIn RO
  SetOutPath "$INSTDIR"
  File "..\binaries\windows-amd64\servin.exe"
  File "..\binaries\windows-amd64\vm-config.json"
  
  SetOutPath "$INSTDIR\webview_gui"
  File /r "..\binaries\windows-amd64\webview_gui\*"
  
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin" "DisplayName" "Servin Container Runtime"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin" "UninstallString" "$INSTDIR\uninstall.exe"
  WriteUninstaller "$INSTDIR\uninstall.exe"
SectionEnd

Section "Desktop GUI" SEC02
  SetOutPath "$INSTDIR"
  File "..\binaries\windows-amd64\servin-desktop.exe"
  CreateShortCut "$DESKTOP\Servin.lnk" "$INSTDIR\servin-desktop.exe"
  CreateDirectory "$SMPROGRAMS\Servin"
  CreateShortCut "$SMPROGRAMS\Servin\Servin.lnk" "$INSTDIR\servin-desktop.exe"
SectionEnd

Section "VM Images" SEC03
  SetOutPath "$INSTDIR\vm-images"
  ; Download VM images during installation
  inetc::get "https://releases.servin.dev/vm-images/alpine-${VM_VERSION}.qcow2" "$INSTDIR\vm-images\alpine.qcow2"
SectionEnd
EOF

    # Create WiX installer for Windows Store
    cat > "$installer_dir/servin.wxs" << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<Wix xmlns="http://schemas.microsoft.com/wix/2006/wi">
  <Product Id="*" Name="Servin" Language="1033" Version="${VERSION}" Manufacturer="Servin Project" UpgradeCode="12345678-1234-1234-1234-123456789012">
    <Package InstallerVersion="200" Compressed="yes" InstallScope="perMachine" />
    
    <MajorUpgrade DowngradeErrorMessage="A newer version of [ProductName] is already installed." />
    <MediaTemplate EmbedCab="yes" />
    
    <Feature Id="ProductFeature" Title="Servin" Level="1">
      <ComponentGroupRef Id="ProductComponents" />
    </Feature>
  </Product>
  
  <Fragment>
    <Directory Id="TARGETDIR" Name="SourceDir">
      <Directory Id="ProgramFiles64Folder">
        <Directory Id="INSTALLFOLDER" Name="Servin" />
      </Directory>
    </Directory>
  </Fragment>
  
  <Fragment>
    <ComponentGroup Id="ProductComponents" Directory="INSTALLFOLDER">
      <Component Id="ServinExe" Guid="*">
        <File Id="ServinExe" Source="..\binaries\windows-amd64\servin.exe" />
      </Component>
    </ComponentGroup>
  </Fragment>
</Wix>
EOF
}

create_linux_packages() {
    local packages_dir="$BUILD_DIR/packages"
    
    # Create Debian package
    create_deb_package "$packages_dir/deb"
    
    # Create RPM package  
    create_rpm_package "$packages_dir/rpm"
    
    # Create Snap package
    create_snap_package "$packages_dir/snap"
    
    # Create Flatpak
    create_flatpak_package "$packages_dir/flatpak"
}

create_deb_package() {
    local deb_dir=$1
    mkdir -p "$deb_dir/servin-$VERSION"/{DEBIAN,usr/bin,usr/share/servin,etc/servin}
    
    # Control file
    cat > "$deb_dir/servin-$VERSION/DEBIAN/control" << EOF
Package: servin
Version: $VERSION
Section: utils
Priority: optional
Architecture: amd64
Depends: qemu-system-x86, qemu-utils
Maintainer: Servin Project <contact@servin.dev>
Description: Universal container runtime with VM-based containerization
 Servin provides true containerization across all platforms by running
 containers inside lightweight Linux VMs. Supports Docker API compatibility
 with full process, network, and filesystem isolation.
EOF

    # Post-installation script
    cat > "$deb_dir/servin-$VERSION/DEBIAN/postinst" << 'EOF'
#!/bin/bash
set -e

# Enable VM mode by default on installation
export SERVIN_VM_MODE=true

# Create servin user and group
if ! getent group servin >/dev/null; then
    addgroup --system servin
fi

if ! getent passwd servin >/dev/null; then
    adduser --system --group --no-create-home --shell /bin/false servin
fi

# Set permissions
chown -R servin:servin /etc/servin
chmod 755 /usr/bin/servin

# Initialize VM if requested
if [ "$1" = "configure" ]; then
    echo "Servin installed successfully!"
    echo "Run 'servin vm enable' to set up containerization VM"
fi
EOF

    chmod 755 "$deb_dir/servin-$VERSION/DEBIAN/postinst"
    
    # Copy files
    cp "$BUILD_DIR/binaries/linux-amd64/servin" "$deb_dir/servin-$VERSION/usr/bin/"
    cp -r "$BUILD_DIR/binaries/linux-amd64/webview_gui" "$deb_dir/servin-$VERSION/usr/share/servin/"
    cp "$BUILD_DIR/binaries/linux-amd64/vm-config.json" "$deb_dir/servin-$VERSION/etc/servin/"
    
    # Build package
    dpkg-deb --build "$deb_dir/servin-$VERSION"
}

create_rpm_package() {
    local rpm_dir=$1
    mkdir -p "$rpm_dir"/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
    
    # RPM spec file
    cat > "$rpm_dir/SPECS/servin.spec" << EOF
Name:           servin
Version:        $VERSION
Release:        1%{?dist}
Summary:        Universal container runtime with VM-based containerization

License:        MIT
URL:            https://github.com/immyemperor/servin
Source0:        servin-$VERSION.tar.gz

BuildRequires:  golang >= 1.21
Requires:       qemu-system-x86, qemu-img

%description
Servin provides true containerization across all platforms by running
containers inside lightweight Linux VMs. Supports Docker API compatibility
with full process, network, and filesystem isolation.

%prep
%setup -q

%build
go build -ldflags "-X main.Version=%{version}" -o servin .

%install
mkdir -p %{buildroot}%{_bindir}
mkdir -p %{buildroot}%{_datadir}/servin
mkdir -p %{buildroot}%{_sysconfdir}/servin

install -m 755 servin %{buildroot}%{_bindir}/
cp -r webview_gui %{buildroot}%{_datadir}/servin/
install -m 644 vm-config.json %{buildroot}%{_sysconfdir}/servin/

%files
%{_bindir}/servin
%{_datadir}/servin/
%config(noreplace) %{_sysconfdir}/servin/vm-config.json

%changelog
* $(date +'%a %b %d %Y') Servin Project <contact@servin.dev> - $VERSION-1
- Universal VM containerization release
EOF
}

create_snap_package() {
    local snap_dir=$1
    mkdir -p "$snap_dir"
    
    cat > "$snap_dir/snapcraft.yaml" << EOF
name: servin
version: '$VERSION'
summary: Universal container runtime with VM-based containerization
description: |
  Servin provides true containerization across all platforms by running
  containers inside lightweight Linux VMs. Supports Docker API compatibility
  with full process, network, and filesystem isolation.

grade: stable
confinement: classic

parts:
  servin:
    plugin: go
    source: ../../../
    build-snaps: [go/1.21/stable]
    build-packages: [gcc]
    
apps:
  servin:
    command: bin/servin
    plugs: [home, network, network-bind, removable-media]
EOF
}

package_distributions() {
    print_info "Creating distribution packages..."
    
    for platform_info in "${PLATFORMS[@]}"; do
        IFS=':' read -r platform vm_support <<< "$platform_info"
        IFS='/' read -r os arch <<< "$platform"
        
        package_dir="$BUILD_DIR/packages/servin-$VERSION-$os-$arch"
        mkdir -p "$package_dir"
        
        # Copy binaries
        cp -r "$BUILD_DIR/binaries/$os-$arch"/* "$package_dir/"
        
        # Create README
        create_package_readme "$package_dir" "$os" "$arch" "$vm_support"
        
        # Create install script
        create_install_script "$package_dir" "$os" "$arch"
        
        # Create archive
        cd "$BUILD_DIR/packages"
        if [ "$os" = "windows" ]; then
            zip -r "servin-$VERSION-$os-$arch.zip" "servin-$VERSION-$os-$arch/"
        else
            tar -czf "servin-$VERSION-$os-$arch.tar.gz" "servin-$VERSION-$os-$arch/"
        fi
        cd - > /dev/null
        
        print_success "Packaged $os/$arch"
    done
}

create_package_readme() {
    local package_dir=$1
    local os=$2
    local arch=$3
    local vm_support=$4
    
    cat > "$package_dir/README.md" << EOF
# Servin $VERSION - $os/$arch

Universal container runtime with VM-based true containerization.

## Installation

### Quick Install
\`\`\`bash
$(get_install_command "$os")
\`\`\`

### Manual Install
1. Extract this package to your desired location
2. Run the install script: \`./install.sh\` (Linux/macOS) or \`install.bat\` (Windows)
3. Enable VM mode: \`servin vm enable\`
4. Start containerization: \`servin vm start\`

## VM Support
This build supports: $vm_support

## Usage
\`\`\`bash
# Enable VM-based containerization
servin vm enable
servin vm start

# Run containers with true isolation
servin run alpine echo "Hello World"
servin run nginx -p 80:80

# Check VM status
servin vm status
\`\`\`

## Requirements
$(get_requirements "$os")

## Documentation
- [VM Containerization Guide](https://github.com/immyemperor/servin/blob/master/docs/VM_CONTAINERIZATION.md)
- [Installation Guide](https://github.com/immyemperor/servin/blob/master/INSTALL.md)
- [User Manual](https://github.com/immyemperor/servin/blob/master/README.md)
EOF
}

get_install_command() {
    local os=$1
    case $os in
        linux)
            echo "curl -fsSL https://get.servin.dev | sh"
            ;;
        darwin)
            echo "brew install servin/tap/servin"
            ;;
        windows)
            echo "winget install Servin.Servin"
            ;;
    esac
}

get_requirements() {
    local os=$1
    case $os in
        linux)
            echo "- QEMU/KVM support"
            echo "- 2GB RAM minimum"
            echo "- 10GB disk space"
            ;;
        darwin)
            echo "- macOS 11+ (Big Sur)"
            echo "- Virtualization.framework support"
            echo "- 4GB RAM minimum"
            echo "- 15GB disk space"
            ;;
        windows)
            echo "- Windows 10/11 Pro or Enterprise"
            echo "- Hyper-V support (or VirtualBox fallback)"
            echo "- 4GB RAM minimum"
            echo "- 15GB disk space"
            ;;
    esac
}

create_install_script() {
    local package_dir=$1
    local os=$2
    local arch=$3
    
    if [ "$os" = "windows" ]; then
        cat > "$package_dir/install.bat" << 'EOF'
@echo off
echo Installing Servin...

REM Create installation directory
mkdir "%PROGRAMFILES%\Servin" 2>nul

REM Copy files
copy servin.exe "%PROGRAMFILES%\Servin\"
xcopy /E /I webview_gui "%PROGRAMFILES%\Servin\webview_gui"
copy vm-config.json "%PROGRAMFILES%\Servin\"

REM Add to PATH
setx PATH "%PATH%;%PROGRAMFILES%\Servin" /M

echo Servin installed successfully!
echo Run 'servin vm enable' to set up containerization
pause
EOF
    else
        cat > "$package_dir/install.sh" << 'EOF'
#!/bin/bash
set -e

echo "Installing Servin..."

# Determine install location
if [ "$EUID" -eq 0 ]; then
    INSTALL_DIR="/usr/local"
    SUDO=""
else
    INSTALL_DIR="$HOME/.local"
    SUDO=""
    mkdir -p "$INSTALL_DIR/bin"
fi

# Copy binaries
$SUDO cp servin "$INSTALL_DIR/bin/"
if [ -f servin-desktop ]; then
    $SUDO cp servin-desktop "$INSTALL_DIR/bin/"
fi

# Copy resources
$SUDO mkdir -p "$INSTALL_DIR/share/servin"
$SUDO cp -r webview_gui "$INSTALL_DIR/share/servin/"

# Copy config
$SUDO mkdir -p "$INSTALL_DIR/etc/servin" 2>/dev/null || mkdir -p "$HOME/.config/servin"
if [ "$EUID" -eq 0 ]; then
    $SUDO cp vm-config.json "$INSTALL_DIR/etc/servin/"
else
    cp vm-config.json "$HOME/.config/servin/"
fi

# Make executable
$SUDO chmod +x "$INSTALL_DIR/bin/servin"
[ -f "$INSTALL_DIR/bin/servin-desktop" ] && $SUDO chmod +x "$INSTALL_DIR/bin/servin-desktop"

# Add to PATH if needed
if [[ ":$PATH:" != *":$INSTALL_DIR/bin:"* ]]; then
    echo "export PATH=\"$INSTALL_DIR/bin:\$PATH\"" >> "$HOME/.bashrc"
    echo "export PATH=\"$INSTALL_DIR/bin:\$PATH\"" >> "$HOME/.zshrc" 2>/dev/null || true
fi

echo "Servin installed successfully!"
echo "Run 'servin vm enable' to set up containerization VM"
echo "You may need to restart your shell or run: source ~/.bashrc"
EOF
        chmod +x "$package_dir/install.sh"
    fi
}

create_documentation() {
    print_info "Creating distribution documentation..."
    
    local docs_dir="$BUILD_DIR/docs"
    
    # Copy existing docs
    cp -r docs/* "$docs_dir/"
    
    # Create release notes
    cat > "$docs_dir/RELEASE_NOTES.md" << EOF
# Servin $VERSION Release Notes

## Universal VM Containerization

This release introduces universal VM-based containerization, providing true container isolation across all platforms (macOS, Windows, and Linux).

### Key Features

- **True Containerization**: Full process, network, and filesystem isolation on all platforms
- **Universal Compatibility**: Identical container behavior across macOS, Windows, and Linux
- **Lightweight VMs**: Optimized Linux VMs with minimal overhead
- **Docker API Compatibility**: Full compatibility with Docker commands and images
- **Automatic VM Management**: Seamless VM lifecycle management
- **Multiple VM Providers**: Platform-native virtualization (Virtualization.framework, Hyper-V, KVM)

### What's New

- Universal Linux VM implementation for true containerization
- Platform-specific VM providers for optimal performance
- Enhanced container management with VM integration
- Improved GUI with VM status monitoring
- Comprehensive distribution system with automated installers

### Installation

Choose your platform:

#### macOS
\`\`\`bash
brew install servin/tap/servin
# or download from releases
\`\`\`

#### Windows
\`\`\`powershell
winget install Servin.Servin
# or run the .exe installer
\`\`\`

#### Linux
\`\`\`bash
# Ubuntu/Debian
wget https://github.com/immyemperor/servin/releases/download/v$VERSION/servin_$VERSION_amd64.deb
sudo dpkg -i servin_$VERSION_amd64.deb

# CentOS/RHEL/Fedora
sudo yum install https://github.com/immyemperor/servin/releases/download/v$VERSION/servin-$VERSION-1.x86_64.rpm

# Manual install
curl -fsSL https://get.servin.dev | sh
\`\`\`

### Getting Started

1. **Enable VM Mode**:
   \`\`\`bash
   servin vm enable
   servin vm start
   \`\`\`

2. **Run Containers**:
   \`\`\`bash
   servin run alpine echo "Hello from true container!"
   servin run nginx -p 8080:80
   \`\`\`

3. **Monitor VMs**:
   \`\`\`bash
   servin vm status
   servin ls
   \`\`\`

### Migration Guide

Existing Servin users can seamlessly upgrade:

1. **Automatic Fallback**: Existing containers continue to work with VFS mode
2. **Opt-in VM Mode**: Enable VM mode when ready for true containerization
3. **No Breaking Changes**: All existing commands and configurations work

### Platform Requirements

#### macOS
- macOS 11+ (Big Sur or later)
- Apple Silicon or Intel processors
- 4GB RAM minimum, 8GB recommended
- 15GB disk space for VM images

#### Windows
- Windows 10 Pro/Enterprise or Windows 11
- Hyper-V support (or VirtualBox as fallback)
- 4GB RAM minimum, 8GB recommended
- 15GB disk space for VM images

#### Linux
- x86_64 or ARM64 architecture
- KVM virtualization support
- 2GB RAM minimum, 4GB recommended
- 10GB disk space for VM images

### Known Issues

- Initial VM setup requires internet connection for image download
- VM startup time: 10-30 seconds on first run
- Some container networking features may require additional configuration

### Breaking Changes

None. This release maintains full backward compatibility.

### Documentation

- [VM Containerization Guide](VM_CONTAINERIZATION.md)
- [Installation Guide](../INSTALL.md)
- [User Manual](../README.md)
- [API Reference](API.md)

### Support

- GitHub Issues: https://github.com/immyemperor/servin/issues
- Documentation: https://servin.dev/docs
- Community: https://discord.gg/servin
EOF

    print_success "Distribution documentation created"
}

main() {
    # Parse command line arguments
    ENHANCE_EXISTING=false
    PLATFORM_FILTER=""
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --enhance-existing)
                ENHANCE_EXISTING=true
                shift
                ;;
            --platform)
                PLATFORM_FILTER="$2"
                shift 2
                ;;
            --help)
                show_help
                exit 0
                ;;
            *)
                shift
                ;;
        esac
    done
    
    print_header
    
    if [ "$ENHANCE_EXISTING" = true ]; then
        print_info "Enhancing existing build with VM capabilities..."
        enhance_existing_build
    else
        print_info "Building Servin distribution with universal VM containerization..."
        echo "Version: $VERSION"
        echo "Build Time: $BUILD_TIME" 
        echo "VM Version: $VM_VERSION"
        echo ""
        
        create_build_structure
        build_binaries
        build_vm_images
        create_installers
        package_distributions
        create_documentation
        
        print_success "Distribution build complete!"
        echo ""
        echo "Generated distributions:"
        find "$BUILD_DIR" -name "*.tar.gz" -o -name "*.zip" -o -name "*.deb" -o -name "*.rpm" | sort
        echo ""
        echo "Ready for release to GitHub, package managers, and distribution channels."
    fi
}

enhance_existing_build() {
    print_info "Enhancing existing build with VM capabilities..."
    
    # Map platform names from GitHub Actions to our format
    case "$PLATFORM_FILTER" in
        mac|macos) PLATFORM_FILTER="darwin" ;;
        windows) PLATFORM_FILTER="windows" ;;
        linux) PLATFORM_FILTER="linux" ;;
    esac
    
    # Check if traditional build directories exist
    if [ -d "build/$PLATFORM_FILTER" ]; then
        print_info "Found existing build for $PLATFORM_FILTER, adding VM support..."
        
        # Enhance binaries with VM tags
        enhance_binaries_with_vm "$PLATFORM_FILTER"
        
        # Add VM components to distribution
        enhance_distribution_with_vm "$PLATFORM_FILTER"
        
        print_success "VM enhancement complete for $PLATFORM_FILTER"
    else
        print_warning "No existing build found for $PLATFORM_FILTER, running full build..."
        # Fall back to full build
        main
    fi
}

enhance_binaries_with_vm() {
    local platform="$1"
    local build_dir="build/$platform"
    local dist_dir="dist/$platform"
    
    print_info "Adding VM capabilities to existing binaries..."
    
    # Determine OS and architecture
    case "$platform" in
        darwin|mac|macos)
            local os="darwin"
            local arch="universal"  # GitHub Actions builds universal
            ;;
        linux)
            local os="linux" 
            local arch="amd64"
            ;;
        windows)
            local os="windows"
            local arch="amd64"
            ;;
    esac
    
    # Rebuild main binary with VM support
    local binary_name="servin"
    if [ "$os" = "windows" ]; then
        binary_name="servin.exe"
    fi
    
    if [ -f "$build_dir/$binary_name" ]; then
        print_info "Rebuilding $binary_name with VM support..."
        
        # Determine VM providers for platform
        local vm_tags=""
        case "$os" in
            linux) vm_tags="vm_enabled,kvm,qemu" ;;
            darwin) vm_tags="vm_enabled,hvf,qemu" ;;
            windows) vm_tags="vm_enabled,hyperv,vbox" ;;
        esac
        
        # Rebuild with VM support
        GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build \
            -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -s -w" \
            -tags "$vm_tags" \
            -o "$build_dir/$binary_name" .
            
        print_success "Enhanced $binary_name with VM support"
    fi
    
    # Add VM helper binaries if not present
    add_vm_helpers "$build_dir" "$os" "$arch"
}

add_vm_helpers() {
    local build_dir="$1"
    local os="$2" 
    local arch="$3"
    
    print_info "Adding VM helper components..."
    
    # Add VM configuration
    cat > "$build_dir/vm-config.json" << EOF
{
    "version": "$VM_VERSION",
    "providers": {
        "primary": "$(get_primary_vm_provider $os)",
        "fallback": "qemu"
    },
    "images": {
        "default": "alpine:3.18",
        "alternatives": ["ubuntu:22.04", "debian:12"]
    },
    "resources": {
        "memory": "512M",
        "disk": "2G",
        "cpus": 2
    }
}
EOF

    # Add VM startup script
    case "$os" in
        windows)
            cat > "$build_dir/start-vm.bat" << 'EOF'
@echo off
echo Starting Servin VM...
servin.exe vm start
EOF
            ;;
        *)
            cat > "$build_dir/start-vm.sh" << 'EOF'
#!/bin/bash
echo "Starting Servin VM..."
./servin vm start
EOF
            chmod +x "$build_dir/start-vm.sh"
            ;;
    esac
    
    print_success "VM helper components added"
}

get_primary_vm_provider() {
    case "$1" in
        linux) echo "kvm" ;;
        darwin) echo "hvf" ;; 
        windows) echo "hyperv" ;;
        *) echo "qemu" ;;
    esac
}

enhance_distribution_with_vm() {
    local platform="$1"
    local dist_dir="dist/$platform"
    
    if [ ! -d "$dist_dir" ]; then
        print_warning "Distribution directory not found for $platform"
        return
    fi
    
    print_info "Adding VM documentation and guides..."
    
    # Add VM setup guide
    cat > "$dist_dir/VM_SETUP.md" << 'EOF'
# VM Containerization Setup

Servin now includes VM-based containerization for true container isolation across all platforms.

## Quick Start

1. Enable VM mode:
   ```
   servin vm enable
   ```

2. Start the VM:
   ```
   servin vm start
   ```

3. Run containers:
   ```
   servin run alpine echo "Hello from VM!"
   ```

## Platform-Specific Notes

### macOS
- Uses Apple's Virtualization.framework (macOS 11+)
- Fallback to QEMU for older systems

### Linux  
- Prefers KVM for hardware acceleration
- Falls back to QEMU if KVM unavailable

### Windows
- Uses Hyper-V on Windows 10/11 Pro
- VirtualBox fallback for Home editions

## Troubleshooting

If you encounter issues:
- Check virtualization is enabled in BIOS
- Ensure sufficient RAM (minimum 2GB free)
- Verify network connectivity for image downloads

For more help: https://servin.dev/docs/vm-troubleshooting
EOF

    # Update main README if it exists
    if [ -f "$dist_dir/README.md" ]; then
        # Add VM section to existing README
        cat >> "$dist_dir/README.md" << 'EOF'

## 🚀 VM-Based Containerization

This version includes revolutionary VM-based containerization:

- **True Isolation**: Full process, network, and filesystem isolation
- **Cross-Platform**: Same container behavior on macOS, Linux, and Windows
- **Native Performance**: Hardware-accelerated virtualization
- **Easy Setup**: Automatic VM management and configuration

### Quick Start with VMs

```bash
# Enable VM mode (one-time setup)
servin vm enable

# Start containers with VM isolation
servin run --vm alpine echo "Hello World!"

# Check VM status
servin vm status
```

See `VM_SETUP.md` for detailed configuration options.
EOF
    fi
    
    print_success "VM documentation added to distribution"
}

show_help() {
    cat << EOF
Servin VM Distribution Builder

USAGE:
    $0 [OPTIONS]

OPTIONS:
    --enhance-existing     Enhance existing build with VM capabilities
    --platform PLATFORM   Target specific platform (linux/darwin/windows)
    --help                Show this help message

EXAMPLES:
    # Full VM distribution build
    $0
    
    # Enhance existing GitHub Actions build
    $0 --enhance-existing --platform linux
    
    # Platform-specific enhancement  
    $0 --enhance-existing --platform darwin

EOF
}

# Run main function
main "$@"