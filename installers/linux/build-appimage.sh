#!/bin/bash

# Servin Container Runtime - Linux AppImage Builder
# Creates a portable Linux application with all VM dependencies included

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_DIR="$SCRIPT_DIR/build"
APPDIR="$BUILD_DIR/Servin.AppDir"
ARCH=$(uname -m)
VERSION="1.0.0"

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
    echo "║              Servin Container Runtime AppImage Builder        ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}\n"
}

# Check prerequisites
check_prerequisites() {
    print_header "Checking Prerequisites"
    
    # Check for required tools
    local missing_tools=()
    
    for tool in wget curl tar; do
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

# Download AppImage tools
download_tools() {
    print_header "Downloading AppImage Tools"
    
    local tools_dir="$BUILD_DIR/tools"
    mkdir -p "$tools_dir"
    
    # Download linuxdeploy
    local linuxdeploy_url="https://github.com/linuxdeploy/linuxdeploy/releases/download/continuous/linuxdeploy-${ARCH}.AppImage"
    if [[ ! -f "$tools_dir/linuxdeploy" ]]; then
        print_info "Downloading linuxdeploy..."
        wget -q "$linuxdeploy_url" -O "$tools_dir/linuxdeploy"
        chmod +x "$tools_dir/linuxdeploy"
        print_success "linuxdeploy downloaded"
    fi
    
    # Download appimagetool
    local appimagetool_url="https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-${ARCH}.AppImage"
    if [[ ! -f "$tools_dir/appimagetool" ]]; then
        print_info "Downloading appimagetool..."
        wget -q "$appimagetool_url" -O "$tools_dir/appimagetool"
        chmod +x "$tools_dir/appimagetool"
        print_success "appimagetool downloaded"
    fi
}

# Create AppDir structure
create_appdir() {
    print_header "Creating AppDir Structure"
    
    rm -rf "$APPDIR"
    mkdir -p "$APPDIR"/{usr/bin,usr/lib,usr/share/{applications,icons/hicolor/256x256/apps},etc,opt/servin}
    
    # Copy Servin executables
    print_info "Copying Servin executables..."
    cp "$SCRIPT_DIR/../../servin" "$APPDIR/usr/bin/"
    cp "$SCRIPT_DIR/../../servin-tui" "$APPDIR/usr/bin/" 2>/dev/null || print_warning "servin-tui not found"
    cp "$SCRIPT_DIR/../../servin-gui" "$APPDIR/usr/bin/" 2>/dev/null || print_warning "servin-gui not found"
    chmod +x "$APPDIR/usr/bin/servin"*
    
    # Create desktop file
    print_info "Creating desktop file..."
    cat > "$APPDIR/usr/share/applications/servin.desktop" << EOF
[Desktop Entry]
Type=Application
Name=Servin Container Runtime
Comment=Docker-compatible container runtime with VM support
Exec=servin-gui
Icon=servin
Categories=Development;System;
Terminal=false
StartupNotify=true
MimeType=application/x-servin-config;
EOF
    
    # Create AppRun script
    print_info "Creating AppRun script..."
    cat > "$APPDIR/AppRun" << 'EOF'
#!/bin/bash

# Servin Container Runtime AppImage Launcher

set -e

# Get the directory where this AppImage is mounted
APPDIR="$(dirname "$(readlink -f "${0}")")"

# Set up environment
export PATH="$APPDIR/usr/bin:$APPDIR/opt/servin/bin:$PATH"
export LD_LIBRARY_PATH="$APPDIR/usr/lib:$LD_LIBRARY_PATH"
export SERVIN_DATA_DIR="${SERVIN_DATA_DIR:-$HOME/.servin}"
export SERVIN_CONFIG_DIR="${SERVIN_CONFIG_DIR:-$HOME/.config/servin}"

# Create directories if they don't exist
mkdir -p "$SERVIN_DATA_DIR"/{vm/{images,instances},logs}
mkdir -p "$SERVIN_CONFIG_DIR"

# Check for VM prerequisites and offer to install if missing
check_vm_prerequisites() {
    local missing_prereq=false
    
    if ! command -v qemu-system-x86_64 >/dev/null 2>&1; then
        missing_prereq=true
    fi
    
    if [[ ! -e /dev/kvm ]]; then
        missing_prereq=true
    fi
    
    if [[ "$missing_prereq" == "true" ]]; then
        if command -v zenity >/dev/null 2>&1; then
            if zenity --question --text="VM prerequisites are missing. Would you like to install them now?\n\nThis will run the system installer and may require your password."; then
                # Try to find and run the enhanced installer
                if [[ -f "$APPDIR/opt/servin/install-with-vm.sh" ]]; then
                    exec "$APPDIR/opt/servin/install-with-vm.sh"
                else
                    zenity --error --text="Enhanced installer not found. Please install VM prerequisites manually."
                fi
            fi
        else
            echo "Warning: VM prerequisites are missing. Please install QEMU/KVM manually."
        fi
    fi
}

# Determine what to run based on arguments
if [[ $# -eq 0 ]]; then
    # No arguments - check for GUI, fall back to TUI, then CLI
    if [[ -x "$APPDIR/usr/bin/servin-gui" ]]; then
        check_vm_prerequisites
        exec "$APPDIR/usr/bin/servin-gui" "$@"
    elif [[ -x "$APPDIR/usr/bin/servin-tui" ]]; then
        check_vm_prerequisites
        exec "$APPDIR/usr/bin/servin-tui" "$@"
    else
        exec "$APPDIR/usr/bin/servin" "$@"
    fi
elif [[ "$1" == "gui" ]]; then
    check_vm_prerequisites
    exec "$APPDIR/usr/bin/servin-gui" "${@:2}"
elif [[ "$1" == "tui" ]]; then
    check_vm_prerequisites
    exec "$APPDIR/usr/bin/servin-tui" "${@:2}"
else
    # Pass through to main servin binary
    exec "$APPDIR/usr/bin/servin" "$@"
fi
EOF
    chmod +x "$APPDIR/AppRun"
    
    # Copy enhanced installer
    print_info "Including enhanced installer..."
    mkdir -p "$APPDIR/opt/servin"
    cp "$SCRIPT_DIR/../linux/install-with-vm.sh" "$APPDIR/opt/servin/" 2>/dev/null || print_warning "Enhanced installer not found"
    
    print_success "AppDir structure created"
}

# Create icon
create_icon() {
    print_header "Creating Application Icon"
    
    # Create a simple SVG icon
    cat > "$APPDIR/usr/share/icons/hicolor/256x256/apps/servin.svg" << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<svg width="256" height="256" viewBox="0 0 256 256" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <linearGradient id="grad1" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" style="stop-color:#4A90E2;stop-opacity:1" />
      <stop offset="100%" style="stop-color:#357ABD;stop-opacity:1" />
    </linearGradient>
  </defs>
  <rect width="256" height="256" rx="32" fill="url(#grad1)"/>
  <rect x="48" y="80" width="160" height="96" rx="8" fill="white" opacity="0.9"/>
  <rect x="64" y="96" width="32" height="16" rx="2" fill="#4A90E2"/>
  <rect x="104" y="96" width="32" height="16" rx="2" fill="#4A90E2"/>
  <rect x="144" y="96" width="32" height="16" rx="2" fill="#4A90E2"/>
  <rect x="64" y="120" width="32" height="16" rx="2" fill="#7BB3F0"/>
  <rect x="104" y="120" width="32" height="16" rx="2" fill="#7BB3F0"/>
  <rect x="144" y="120" width="32" height="16" rx="2" fill="#7BB3F0"/>
  <rect x="64" y="144" width="32" height="16" rx="2" fill="#A0C8F5"/>
  <rect x="104" y="144" width="32" height="16" rx="2" fill="#A0C8F5"/>
  <rect x="144" y="144" width="32" height="16" rx="2" fill="#A0C8F5"/>
  <text x="128" y="210" text-anchor="middle" font-family="Arial, sans-serif" font-size="20" font-weight="bold" fill="white">SERVIN</text>
</svg>
EOF
    
    # Convert SVG to PNG if imagemagick is available
    if command -v convert >/dev/null 2>&1; then
        convert "$APPDIR/usr/share/icons/hicolor/256x256/apps/servin.svg" "$APPDIR/usr/share/icons/hicolor/256x256/apps/servin.png"
        cp "$APPDIR/usr/share/icons/hicolor/256x256/apps/servin.png" "$APPDIR/servin.png"
    else
        cp "$APPDIR/usr/share/icons/hicolor/256x256/apps/servin.svg" "$APPDIR/servin.svg"
    fi
    
    print_success "Application icon created"
}

# Bundle dependencies
bundle_dependencies() {
    print_header "Bundling Dependencies"
    
    # Use linuxdeploy to bundle dependencies
    print_info "Running linuxdeploy..."
    
    cd "$BUILD_DIR"
    ./tools/linuxdeploy --appdir "$APPDIR" \
        --executable "$APPDIR/usr/bin/servin" \
        --desktop-file "$APPDIR/usr/share/applications/servin.desktop" \
        --icon-file "$APPDIR/usr/share/icons/hicolor/256x256/apps/servin"* \
        --output appimage \
        --verbosity 2
    
    print_success "Dependencies bundled"
}

# Create Python environment (optional)
create_python_env() {
    print_header "Creating Python Environment"
    
    # Check if Python is available
    if ! command -v python3 >/dev/null 2>&1; then
        print_warning "Python3 not found, skipping Python environment creation"
        return
    fi
    
    # Create minimal Python environment for GUI
    local python_dir="$APPDIR/opt/python"
    mkdir -p "$python_dir"
    
    print_info "Installing Python packages for GUI..."
    python3 -m pip install --target "$python_dir" \
        pywebview flask flask-cors flask-socketio eventlet 2>/dev/null || \
        print_warning "Failed to install Python packages"
    
    # Add Python to AppRun environment
    sed -i '/export PATH=/c\export PATH="$APPDIR/usr/bin:$APPDIR/opt/servin/bin:$PATH"\nexport PYTHONPATH="$APPDIR/opt/python:$PYTHONPATH"' "$APPDIR/AppRun"
    
    print_success "Python environment created"
}

# Create final AppImage
create_appimage() {
    print_header "Creating Final AppImage"
    
    cd "$BUILD_DIR"
    
    # Create the AppImage
    print_info "Building AppImage..."
    ARCH="$ARCH" ./tools/appimagetool "$APPDIR" "Servin-${VERSION}-${ARCH}.AppImage"
    
    if [[ -f "Servin-${VERSION}-${ARCH}.AppImage" ]]; then
        print_success "AppImage created successfully: Servin-${VERSION}-${ARCH}.AppImage"
        
        # Make executable
        chmod +x "Servin-${VERSION}-${ARCH}.AppImage"
        
        # Show file info
        local size=$(du -h "Servin-${VERSION}-${ARCH}.AppImage" | cut -f1)
        print_info "AppImage size: $size"
        
        # Test the AppImage
        print_info "Testing AppImage..."
        if "./Servin-${VERSION}-${ARCH}.AppImage" version >/dev/null 2>&1; then
            print_success "AppImage test passed"
        else
            print_warning "AppImage test failed"
        fi
    else
        print_error "Failed to create AppImage"
        exit 1
    fi
}

# Create installation script
create_installer_script() {
    print_header "Creating Installation Script"
    
    cat > "$BUILD_DIR/install-servin-appimage.sh" << 'EOF'
#!/bin/bash

# Servin Container Runtime - AppImage Installation Script

set -e

APPIMAGE_FILE=""
INSTALL_DIR="$HOME/.local/bin"
DESKTOP_DIR="$HOME/.local/share/applications"
ICON_DIR="$HOME/.local/share/icons/hicolor/256x256/apps"

# Find AppImage file
for file in Servin-*.AppImage; do
    if [[ -f "$file" ]]; then
        APPIMAGE_FILE="$file"
        break
    fi
done

if [[ -z "$APPIMAGE_FILE" ]]; then
    echo "Error: No Servin AppImage found"
    exit 1
fi

echo "Installing Servin Container Runtime AppImage..."

# Create directories
mkdir -p "$INSTALL_DIR" "$DESKTOP_DIR" "$ICON_DIR"

# Copy AppImage
cp "$APPIMAGE_FILE" "$INSTALL_DIR/servin"
chmod +x "$INSTALL_DIR/servin"

# Extract and install desktop file
"$INSTALL_DIR/servin" --appimage-extract usr/share/applications/servin.desktop >/dev/null 2>&1
if [[ -f "squashfs-root/usr/share/applications/servin.desktop" ]]; then
    sed "s|Exec=servin-gui|Exec=$INSTALL_DIR/servin gui|g" \
        "squashfs-root/usr/share/applications/servin.desktop" > "$DESKTOP_DIR/servin.desktop"
    chmod +x "$DESKTOP_DIR/servin.desktop"
    rm -rf squashfs-root
fi

# Extract and install icon
"$INSTALL_DIR/servin" --appimage-extract usr/share/icons/hicolor/256x256/apps/servin.* >/dev/null 2>&1
if [[ -d "squashfs-root/usr/share/icons/hicolor/256x256/apps" ]]; then
    cp squashfs-root/usr/share/icons/hicolor/256x256/apps/servin.* "$ICON_DIR/" 2>/dev/null || true
    rm -rf squashfs-root
fi

# Update desktop database
update-desktop-database "$DESKTOP_DIR" 2>/dev/null || true

echo "✓ Servin Container Runtime installed successfully!"
echo ""
echo "Usage:"
echo "  servin                 - Command line interface"
echo "  servin gui             - Graphical interface"
echo "  servin tui             - Terminal interface"
echo ""
echo "The application is also available in your application menu."
EOF
    
    chmod +x "$BUILD_DIR/install-servin-appimage.sh"
    print_success "Installation script created"
}

# Main execution
main() {
    print_banner
    
    check_prerequisites
    download_tools
    create_appdir
    create_icon
    create_python_env
    bundle_dependencies
    create_appimage
    create_installer_script
    
    print_header "Build Complete!"
    echo
    print_success "Servin AppImage built successfully!"
    echo
    print_info "Files created:"
    echo "  • $BUILD_DIR/Servin-${VERSION}-${ARCH}.AppImage"
    echo "  • $BUILD_DIR/install-servin-appimage.sh"
    echo
    print_info "To install:"
    echo "  cd $BUILD_DIR && ./install-servin-appimage.sh"
    echo
    print_info "To run directly:"
    echo "  ./Servin-${VERSION}-${ARCH}.AppImage"
}

# Run main function
main "$@"