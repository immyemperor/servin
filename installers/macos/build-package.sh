#!/bin/bash

# Servin Container Runtime - macOS Package Builder
# Creates a native macOS installer package (.pkg) with all VM dependencies included

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_DIR="$SCRIPT_DIR/build"
PKG_NAME="Servin"
PKG_IDENTIFIER="com.servin.containerruntime"
VERSION="1.0.0"
ARCH=$(uname -m)

# Color output functions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

print_success() { echo -e "${GREEN}✓ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠ $1${NC}"; }
print_error() { echo -e "${RED}✗ $1${NC}"; }
print_info() { echo -e "${BLUE}→ $1${NC}"; }
print_header() { echo -e "\n${CYAN}${BOLD}$1${NC}"; }

print_banner() {
    echo -e "${CYAN}${BOLD}"
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║              Servin Container Runtime macOS Builder           ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}\n"
}

# Check prerequisites
check_prerequisites() {
    print_header "Checking Prerequisites"
    
    # Check if we're on macOS
    if [[ "$(uname)" != "Darwin" ]]; then
        print_error "This script must be run on macOS"
        exit 1
    fi
    
    # Check for required tools
    local missing_tools=()
    
    for tool in pkgbuild productbuild curl; do
        if ! command -v "$tool" >/dev/null 2>&1; then
            missing_tools+=("$tool")
        fi
    done
    
    if [[ ${#missing_tools[@]} -gt 0 ]]; then
        print_error "Missing required tools: ${missing_tools[*]}"
        exit 1
    fi
    
    # Check for servin executables
    if [[ ! -f "$SCRIPT_DIR/../../servin" ]]; then
        print_error "servin executable not found. Please build it first."
        exit 1
    fi
    
    print_success "Prerequisites check passed"
}

# Create package structure
create_package_structure() {
    print_header "Creating Package Structure"
    
    rm -rf "$BUILD_DIR"
    mkdir -p "$BUILD_DIR"/{root,scripts,resources,payload}
    
    # Create directory structure
    mkdir -p "$BUILD_DIR/root"/{usr/local/bin,Applications,Library/{LaunchDaemons,Application\ Support/Servin}}
    
    # Copy Servin executables
    print_info "Copying Servin executables..."
    cp "$SCRIPT_DIR/../../servin" "$BUILD_DIR/root/usr/local/bin/"
    cp "$SCRIPT_DIR/../../servin-tui" "$BUILD_DIR/root/usr/local/bin/" 2>/dev/null || print_warning "servin-tui not found"
    
    # Create Servin.app bundle
    print_info "Creating Servin.app bundle..."
    local app_dir="$BUILD_DIR/root/Applications/Servin.app"
    mkdir -p "$app_dir/Contents"/{MacOS,Resources}
    
    # Create Info.plist
    cat > "$app_dir/Contents/Info.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleDisplayName</key>
    <string>Servin Container Runtime</string>
    <key>CFBundleExecutable</key>
    <string>Servin</string>
    <key>CFBundleIconFile</key>
    <string>servin.icns</string>
    <key>CFBundleIdentifier</key>
    <string>$PKG_IDENTIFIER</string>
    <key>CFBundleName</key>
    <string>Servin</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>$VERSION</string>
    <key>CFBundleVersion</key>
    <string>$VERSION</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.15</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>NSRequiresAquaSystemAppearance</key>
    <false/>
    <key>CFBundleDocumentTypes</key>
    <array>
        <dict>
            <key>CFBundleTypeExtensions</key>
            <array>
                <string>servin</string>
            </array>
            <key>CFBundleTypeName</key>
            <string>Servin Configuration</string>
            <key>CFBundleTypeRole</key>
            <string>Editor</string>
        </dict>
    </array>
</dict>
</plist>
EOF
    
    # Create app launcher script
    cat > "$app_dir/Contents/MacOS/Servin" << 'EOF'
#!/bin/bash

# Servin macOS App Launcher

# Set up environment
export PATH="/usr/local/bin:$PATH"
export SERVIN_DATA_DIR="${SERVIN_DATA_DIR:-$HOME/.servin}"
export SERVIN_CONFIG_DIR="${SERVIN_CONFIG_DIR:-$HOME/.config/servin}"

# Create directories if they don't exist
mkdir -p "$SERVIN_DATA_DIR"/{vm/{images,instances},logs}
mkdir -p "$SERVIN_CONFIG_DIR"

# Check for VM prerequisites
check_vm_prerequisites() {
    local missing_prereq=false
    local install_needed=()
    
    # Check for QEMU
    if ! command -v qemu-system-x86_64 >/dev/null 2>&1; then
        missing_prereq=true
        install_needed+=("QEMU")
    fi
    
    # Check for Virtualization.framework (macOS 11+)
    if [[ $(sw_vers -productVersion | cut -d. -f1) -lt 11 ]]; then
        missing_prereq=true
        install_needed+=("macOS 11+ for Virtualization.framework")
    fi
    
    if [[ "$missing_prereq" == "true" ]]; then
        local message="VM prerequisites are missing:\n${install_needed[*]}\n\nWould you like to install them now?"
        
        # Use AppleScript for native dialog
        if osascript -e "display dialog \"$message\" buttons {\"Cancel\", \"Install\"} default button \"Install\"" >/dev/null 2>&1; then
            # Run the enhanced installer
            if [[ -f "/Library/Application Support/Servin/install-with-vm.sh" ]]; then
                osascript -e 'do shell script "/Library/Application Support/Servin/install-with-vm.sh" with administrator privileges'
            else
                osascript -e 'display dialog "Enhanced installer not found. Please install VM prerequisites manually." buttons {"OK"} default button "OK"'
            fi
        fi
    fi
}

# Check for GUI availability and launch appropriate interface
if [[ -f "/usr/local/bin/servin-gui" ]]; then
    check_vm_prerequisites
    exec "/usr/local/bin/servin-gui" "$@"
elif [[ -f "/usr/local/bin/servin-tui" ]]; then
    # Launch TUI in Terminal
    check_vm_prerequisites
    osascript -e 'tell application "Terminal" to do script "/usr/local/bin/servin-tui"'
else
    # Launch CLI in Terminal
    osascript -e 'tell application "Terminal" to do script "/usr/local/bin/servin"'
fi
EOF
    chmod +x "$app_dir/Contents/MacOS/Servin"
    
    print_success "Package structure created"
}

# Create app icon
create_app_icon() {
    print_header "Creating Application Icon"
    
    local app_dir="$BUILD_DIR/root/Applications/Servin.app"
    local iconset_dir="$BUILD_DIR/servin.iconset"
    
    mkdir -p "$iconset_dir"
    
    # Create SVG icon first
    cat > "$BUILD_DIR/servin.svg" << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<svg width="1024" height="1024" viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <linearGradient id="grad1" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" style="stop-color:#4A90E2;stop-opacity:1" />
      <stop offset="100%" style="stop-color:#357ABD;stop-opacity:1" />
    </linearGradient>
  </defs>
  <rect width="1024" height="1024" rx="180" fill="url(#grad1)"/>
  <rect x="192" y="320" width="640" height="384" rx="32" fill="white" opacity="0.9"/>
  <rect x="256" y="384" width="128" height="64" rx="8" fill="#4A90E2"/>
  <rect x="416" y="384" width="128" height="64" rx="8" fill="#4A90E2"/>
  <rect x="576" y="384" width="128" height="64" rx="8" fill="#4A90E2"/>
  <rect x="256" y="480" width="128" height="64" rx="8" fill="#7BB3F0"/>
  <rect x="416" y="480" width="128" height="64" rx="8" fill="#7BB3F0"/>
  <rect x="576" y="480" width="128" height="64" rx="8" fill="#7BB3F0"/>
  <rect x="256" y="576" width="128" height="64" rx="8" fill="#A0C8F5"/>
  <rect x="416" y="576" width="128" height="64" rx="8" fill="#A0C8F5"/>
  <rect x="576" y="576" width="128" height="64" rx="8" fill="#A0C8F5"/>
  <text x="512" y="840" text-anchor="middle" font-family="SF Pro Display, Arial, sans-serif" font-size="80" font-weight="bold" fill="white">SERVIN</text>
</svg>
EOF
    
    # Generate icon sizes using sips if available
    if command -v sips >/dev/null 2>&1; then
        # Convert SVG to high-res PNG first
        if command -v qlmanage >/dev/null 2>&1; then
            qlmanage -t -s 1024 -o "$BUILD_DIR" "$BUILD_DIR/servin.svg" >/dev/null 2>&1
            mv "$BUILD_DIR/servin.svg.png" "$BUILD_DIR/servin-1024.png" 2>/dev/null || true
        fi
        
        # Generate all required icon sizes
        local sizes=(16 32 64 128 256 512 1024)
        for size in "${sizes[@]}"; do
            if [[ -f "$BUILD_DIR/servin-1024.png" ]]; then
                sips -z $size $size "$BUILD_DIR/servin-1024.png" --out "$iconset_dir/icon_${size}x${size}.png" >/dev/null 2>&1
                # Create @2x versions
                if [[ $size -lt 512 ]]; then
                    local double_size=$((size * 2))
                    sips -z $double_size $double_size "$BUILD_DIR/servin-1024.png" --out "$iconset_dir/icon_${size}x${size}@2x.png" >/dev/null 2>&1
                fi
            fi
        done
        
        # Create icns file
        if command -v iconutil >/dev/null 2>&1; then
            iconutil -c icns "$iconset_dir" -o "$app_dir/Contents/Resources/servin.icns"
            print_success "Application icon created"
        else
            print_warning "iconutil not available, using default icon"
        fi
    else
        print_warning "sips not available, using default icon"
    fi
    
    # Cleanup
    rm -rf "$iconset_dir" "$BUILD_DIR/servin.svg" "$BUILD_DIR/servin-1024.png" 2>/dev/null || true
}

# Download and bundle VM dependencies
bundle_vm_dependencies() {
    print_header "Bundling VM Dependencies"
    
    local vm_dir="$BUILD_DIR/root/Library/Application Support/Servin/vm"
    mkdir -p "$vm_dir"
    
    # Download QEMU for macOS
    print_info "Downloading QEMU for macOS..."
    
    local qemu_version="8.1.0"
    local qemu_url=""
    
    if [[ "$ARCH" == "arm64" ]]; then
        qemu_url="https://github.com/kholia/qemu-builds/releases/download/v${qemu_version}/qemu-${qemu_version}-arm64.tar.xz"
    else
        qemu_url="https://github.com/kholia/qemu-builds/releases/download/v${qemu_version}/qemu-${qemu_version}-x86_64.tar.xz"
    fi
    
    # Download and extract QEMU
    print_info "Attempting to download QEMU from: $qemu_url"
    
    if curl -L --max-time 300 --retry 3 "$qemu_url" -o "$vm_dir/qemu.tar.xz"; then
        print_info "QEMU download successful, extracting..."
        cd "$vm_dir"
        
        if tar -xf qemu.tar.xz --strip-components=1; then
            rm -f qemu.tar.xz
            
            # Verify QEMU extraction was successful
            if [[ -f "bin/qemu-system-x86_64" ]] || [[ -f "bin/qemu-system-aarch64" ]]; then
                print_success "QEMU bundled successfully"
                
                # Show bundled QEMU size for verification
                QEMU_SIZE=$(du -sh . | cut -f1)
                print_info "Bundled QEMU size: $QEMU_SIZE"
            else
                print_warning "QEMU extraction verification failed - binary not found"
                # Create a fallback marker
                echo "QEMU_FALLBACK=homebrew" > "$vm_dir/qemu-fallback.txt"
            fi
        else
            print_warning "Failed to extract QEMU archive"
            rm -f qemu.tar.xz
            echo "QEMU_FALLBACK=homebrew" > "$vm_dir/qemu-fallback.txt"
        fi
    else
        print_warning "Failed to download QEMU from $qemu_url"
        print_info "Package will use system QEMU installation via Homebrew"
        echo "QEMU_FALLBACK=homebrew" > "$vm_dir/qemu-fallback.txt"
    fi
    
    # Create VM setup script
    cat > "$vm_dir/setup-vm.sh" << 'EOF'
#!/bin/bash

# Servin VM Setup Script for macOS

set -e

# Add bundled QEMU to PATH if it exists
if [[ -d "/Library/Application Support/Servin/vm/bin" ]]; then
    export PATH="/Library/Application Support/Servin/vm/bin:$PATH"
fi

# Check for system QEMU installation
if ! command -v qemu-system-x86_64 >/dev/null 2>&1; then
    echo "Installing QEMU via Homebrew..."
    
    # Check if Homebrew is installed
    if ! command -v brew >/dev/null 2>&1; then
        echo "Installing Homebrew..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        
        # Add Homebrew to PATH
        if [[ -f "/opt/homebrew/bin/brew" ]]; then
            eval "$(/opt/homebrew/bin/brew shellenv)"
        elif [[ -f "/usr/local/bin/brew" ]]; then
            eval "$(/usr/local/bin/brew shellenv)"
        fi
    fi
    
    # Install QEMU
    brew install qemu
fi

# Verify installation
if command -v qemu-system-x86_64 >/dev/null 2>&1; then
    echo "✓ QEMU installed successfully"
    qemu-system-x86_64 --version | head -1
else
    echo "✗ QEMU installation failed"
    exit 1
fi

# Set up VM directories
mkdir -p "$HOME/.servin/vm"/{images,instances}

echo "✓ VM environment set up successfully"
EOF
    chmod +x "$vm_dir/setup-vm.sh"
    
    print_success "VM dependencies bundled"
}

# Create installation scripts
create_install_scripts() {
    print_header "Creating Installation Scripts"
    
    # Pre-install script
    cat > "$BUILD_DIR/scripts/preinstall" << 'EOF'
#!/bin/bash

# Servin Container Runtime - Pre-installation Script

# Stop any running Servin services
pkill -f servin 2>/dev/null || true

# Create application support directory
mkdir -p "/Library/Application Support/Servin"

exit 0
EOF
    
    # Post-install script
    cat > "$BUILD_DIR/scripts/postinstall" << 'EOF'
#!/bin/bash

# Servin Container Runtime - Post-installation Script

set -e

# Set up permissions
chown -R root:wheel "/usr/local/bin/servin"*
chmod +x "/usr/local/bin/servin"*

# Set up Application Support directory
chown -R root:admin "/Library/Application Support/Servin"
chmod -R 755 "/Library/Application Support/Servin"

# Copy enhanced installer
if [[ -f "/tmp/servin-installer/install-with-vm.sh" ]]; then
    cp "/tmp/servin-installer/install-with-vm.sh" "/Library/Application Support/Servin/"
    chmod +x "/Library/Application Support/Servin/install-with-vm.sh"
fi

# Add /usr/local/bin to PATH for all users if not already present
if ! grep -q "/usr/local/bin" /etc/paths; then
    echo "/usr/local/bin" >> /etc/paths
fi

# Create launchd service for system-wide access
cat > "/Library/LaunchDaemons/com.servin.daemon.plist" << 'PLIST'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.servin.daemon</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/servin</string>
        <string>daemon</string>
    </array>
    <key>RunAtLoad</key>
    <false/>
    <key>KeepAlive</key>
    <false/>
    <key>StandardOutPath</key>
    <string>/var/log/servin.log</string>
    <key>StandardErrorPath</key>
    <string>/var/log/servin.log</string>
</dict>
</plist>
PLIST

# Set permissions on launchd service
chown root:wheel "/Library/LaunchDaemons/com.servin.daemon.plist"
chmod 644 "/Library/LaunchDaemons/com.servin.daemon.plist"

# Register with launchctl
launchctl load "/Library/LaunchDaemons/com.servin.daemon.plist" 2>/dev/null || true

# Run VM setup if requested
if [[ -f "/Library/Application Support/Servin/vm/setup-vm.sh" ]]; then
    echo "Setting up VM environment..."
    "/Library/Application Support/Servin/vm/setup-vm.sh" || echo "VM setup completed with warnings"
fi

echo "✓ Servin Container Runtime installed successfully!"

exit 0
EOF
    
    # Make scripts executable
    chmod +x "$BUILD_DIR/scripts"/*
    
    # Copy enhanced installer to temp location for postinstall
    mkdir -p "/tmp/servin-installer"
    cp "$SCRIPT_DIR/../macos/install-with-vm.sh" "/tmp/servin-installer/" 2>/dev/null || print_warning "Enhanced installer not found"
    
    print_success "Installation scripts created"
}

# Create distribution file
create_distribution() {
    print_header "Creating Distribution File"
    
    cat > "$BUILD_DIR/distribution.xml" << EOF
<?xml version="1.0" encoding="utf-8"?>
<installer-gui-script minSpecVersion="1">
    <title>Servin Container Runtime</title>
    <organization>$PKG_IDENTIFIER</organization>
    <domains enable_anywhere="true"/>
    <options customize="never" require-scripts="false" rootVolumeOnly="true" hostArchitectures="$ARCH"/>
    
    <!-- Welcome -->
    <welcome file="welcome.html" mime-type="text/html"/>
    
    <!-- License -->
    <license file="license.txt" mime-type="text/plain"/>
    
    <!-- Background -->
    <background file="background.png" mime-type="image/png" alignment="topleft" scaling="tofit"/>
    
    <!-- Choices -->
    <choices-outline>
        <line choice="default">
            <line choice="com.servin.core"/>
            <line choice="com.servin.vm"/>
        </line>
    </choices-outline>
    
    <choice id="default"/>
    <choice id="com.servin.core" visible="false">
        <pkg-ref id="com.servin.core"/>
    </choice>
    <choice id="com.servin.vm" title="VM Support" description="Install VM dependencies (QEMU) for container virtualization">
        <pkg-ref id="com.servin.vm"/>
    </choice>
    
    <pkg-ref id="com.servin.core" version="$VERSION" onConclusion="none">Servin-Core.pkg</pkg-ref>
    <pkg-ref id="com.servin.vm" version="$VERSION" onConclusion="none">Servin-VM.pkg</pkg-ref>
    
    <!-- Product Definition -->
    <product id="$PKG_IDENTIFIER" version="$VERSION"/>
</installer-gui-script>
EOF
    
    print_success "Distribution file created"
}

# Create resources
create_resources() {
    print_header "Creating Installer Resources"
    
    # Welcome HTML
    cat > "$BUILD_DIR/resources/welcome.html" << 'EOF'
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Welcome to Servin</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, sans-serif; margin: 20px; }
        h1 { color: #4A90E2; }
        .feature { margin: 10px 0; }
        .feature::before { content: "✓ "; color: #4A90E2; font-weight: bold; }
    </style>
</head>
<body>
    <h1>Welcome to Servin Container Runtime</h1>
    <p>Servin is a Docker-compatible container runtime that provides seamless container virtualization on macOS.</p>
    
    <h2>Features:</h2>
    <div class="feature">Docker-compatible API and CLI</div>
    <div class="feature">Native macOS integration</div>
    <div class="feature">VM-based container isolation</div>
    <div class="feature">Graphical and terminal interfaces</div>
    <div class="feature">Automatic dependency management</div>
    
    <p>This installer will set up Servin Container Runtime and optionally install VM dependencies for container virtualization.</p>
</body>
</html>
EOF
    
    # License text
    cat > "$BUILD_DIR/resources/license.txt" << 'EOF'
MIT License

Copyright (c) 2024 Servin Container Runtime

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
EOF
    
    # Create a simple background image (base64 encoded 1x1 pixel)
    echo "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==" | base64 -d > "$BUILD_DIR/resources/background.png"
    
    print_success "Installer resources created"
}

# Build packages
build_packages() {
    print_header "Building Packages"
    
    # Build core package
    print_info "Building core package..."
    pkgbuild --root "$BUILD_DIR/root" \
             --scripts "$BUILD_DIR/scripts" \
             --identifier "$PKG_IDENTIFIER.core" \
             --version "$VERSION" \
             --install-location "/" \
             "$BUILD_DIR/Servin-Core.pkg"
    
    # Build VM package (just the VM dependencies)
    print_info "Building VM package..."
    local vm_root="$BUILD_DIR/vm-root"
    mkdir -p "$vm_root"
    cp -R "$BUILD_DIR/root/Library" "$vm_root/" 2>/dev/null || true
    
    pkgbuild --root "$vm_root" \
             --identifier "$PKG_IDENTIFIER.vm" \
             --version "$VERSION" \
             --install-location "/" \
             "$BUILD_DIR/Servin-VM.pkg"
    
    # Build final product
    print_info "Building final installer..."
    productbuild --distribution "$BUILD_DIR/distribution.xml" \
                 --resources "$BUILD_DIR/resources" \
                 --package-path "$BUILD_DIR" \
                 "$BUILD_DIR/Servin-$VERSION-$ARCH.pkg"
    
    if [[ -f "$BUILD_DIR/Servin-$VERSION-$ARCH.pkg" ]]; then
        print_success "macOS package created successfully: Servin-$VERSION-$ARCH.pkg"
        
        # Show file info
        local size=$(du -h "$BUILD_DIR/Servin-$VERSION-$ARCH.pkg" | cut -f1)
        print_info "Package size: $size"
    else
        print_error "Failed to create macOS package"
        exit 1
    fi
}

# Create DMG (optional)
create_dmg() {
    print_header "Creating DMG (Optional)"
    
    if ! command -v hdiutil >/dev/null 2>&1; then
        print_warning "hdiutil not available, skipping DMG creation"
        return
    fi
    
    print_info "Creating disk image..."
    
    local dmg_dir="$BUILD_DIR/dmg"
    mkdir -p "$dmg_dir"
    
    # Copy installer to DMG
    cp "$BUILD_DIR/Servin-$VERSION-$ARCH.pkg" "$dmg_dir/"
    
    # Create alias to Applications folder
    ln -s /Applications "$dmg_dir/Applications"
    
    # Create DMG
    hdiutil create -volname "Servin Container Runtime" \
                   -srcfolder "$dmg_dir" \
                   -ov -format UDZO \
                   "$BUILD_DIR/Servin-$VERSION-$ARCH.dmg"
    
    if [[ -f "$BUILD_DIR/Servin-$VERSION-$ARCH.dmg" ]]; then
        print_success "DMG created: Servin-$VERSION-$ARCH.dmg"
    fi
}

# Main execution
main() {
    print_banner
    
    check_prerequisites
    create_package_structure
    create_app_icon
    bundle_vm_dependencies
    create_install_scripts
    create_distribution
    create_resources
    build_packages
    create_dmg
    
    print_header "Build Complete!"
    echo
    print_success "Servin macOS package built successfully!"
    echo
    print_info "Files created:"
    echo "  • $BUILD_DIR/Servin-$VERSION-$ARCH.pkg (Main installer)"
    if [[ -f "$BUILD_DIR/Servin-$VERSION-$ARCH.dmg" ]]; then
        echo "  • $BUILD_DIR/Servin-$VERSION-$ARCH.dmg (Disk image)"
    fi
    echo
    print_info "To install:"
    echo "  Double-click the .pkg file or run:"
    echo "  sudo installer -pkg $BUILD_DIR/Servin-$VERSION-$ARCH.pkg -target /"
    echo
    print_info "The application will be available as:"
    echo "  • /Applications/Servin.app (GUI)"
    echo "  • servin command in Terminal"
}

# Run main function
main "$@"