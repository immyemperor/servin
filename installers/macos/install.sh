#!/bin/bash
# Servin Container Runtime - macOS Installer Script
# Updated for platform-organized build system

set -e

# Configuration
INSTALL_DIR="/usr/local/bin"
DATA_DIR="/usr/local/var/lib/servin"
CONFIG_DIR="/usr/local/etc/servin"
LOG_DIR="/usr/local/var/log/servin"
LAUNCHD_DIR="/Library/LaunchDaemons"
SERVICE_NAME="com.servin.runtime"
USER="_servin"

# Build system configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"
PLATFORM="darwin-$(uname -m)"
BUILD_DIR="$PROJECT_ROOT/build/$PLATFORM"

# Check for Apple Silicon vs Intel
if [[ "$(uname -m)" == "arm64" ]]; then
    PLATFORM="darwin-arm64"
else
    PLATFORM="darwin-amd64"
fi

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${CYAN}================================================${NC}"
    echo -e "${CYAN}   Servin Container Runtime - macOS Installer${NC}"
    echo -e "${CYAN}================================================${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}$1${NC}"
}

print_warning() {
    echo -e "${YELLOW}$1${NC}"
}

print_error() {
    echo -e "${RED}$1${NC}"
}

print_info() {
    echo -e "${BLUE}$1${NC}"
}

check_root() {
    if [[ $EUID -ne 0 ]]; then
        print_error "This script must be run as root (use sudo)"
        exit 1
    fi
}

check_macos() {
    if [[ "$(uname)" != "Darwin" ]]; then
        print_error "This installer is only for macOS"
        exit 1
    fi
    
    local version=$(sw_vers -productVersion)
    print_info "macOS version: $version"
}

create_user() {
    print_info "Creating system user: $USER"
    
    # Check if user already exists
    if dscl . -read /Users/$USER >/dev/null 2>&1; then
        print_warning "  User $USER already exists"
        return
    fi
    
    # Find next available UID in system range (< 500)
    local uid
    for ((uid=200; uid<500; uid++)); do
        if ! dscl . -read /Users -uid $uid >/dev/null 2>&1; then
            break
        fi
    done
    
    # Create user
    dscl . -create /Users/$USER
    dscl . -create /Users/$USER UserShell /usr/bin/false
    dscl . -create /Users/$USER RealName "Servin Runtime User"
    dscl . -create /Users/$USER UniqueID $uid
    dscl . -create /Users/$USER PrimaryGroupID 20
    dscl . -create /Users/$USER NFSHomeDirectory /var/empty
    
    print_success "  User $USER created with UID $uid"
}

create_directories() {
    print_info "Creating directories..."
    
    directories=("$DATA_DIR" "$CONFIG_DIR" "$LOG_DIR" "$DATA_DIR/volumes" "$DATA_DIR/images")
    
    for dir in "${directories[@]}"; do
        if [[ ! -d "$dir" ]]; then
            mkdir -p "$dir"
            print_success "  Created: $dir"
        fi
    done
    
    # Set ownership and permissions
    chown -R "$USER:staff" "$DATA_DIR" "$LOG_DIR"
    chmod 755 "$DATA_DIR" "$LOG_DIR"
    chmod 755 "$CONFIG_DIR"
}

install_binaries() {
    print_info "Installing binaries from platform build directory..."
    
    # Check if platform build directory exists
    if [[ ! -d "$BUILD_DIR" ]]; then
        print_error "Platform build directory not found: $BUILD_DIR"
        print_info "Please run './build-local.sh' from the project root first"
        exit 1
    fi
    
    # Check for required binaries
    if [[ ! -f "$BUILD_DIR/servin" ]]; then
        print_error "servin binary not found in $BUILD_DIR"
        print_info "Please run './build-local.sh' from the project root first"
        exit 1
    fi
    
    print_info "Installing from: $BUILD_DIR"
    
    # Install main runtime binary
    cp "$BUILD_DIR/servin" "$INSTALL_DIR/"
    chmod 755 "$INSTALL_DIR/servin"
    print_success "  Installed: $INSTALL_DIR/servin"
    
    # Install TUI binary if available
    if [[ -f "$BUILD_DIR/servin-desktop" ]]; then
        cp "$BUILD_DIR/servin-desktop" "$INSTALL_DIR/"
        chmod 755 "$INSTALL_DIR/servin-desktop"
        print_success "  Installed: $INSTALL_DIR/servin-desktop"
    fi
    
    # Install GUI binary if available
    if [[ -f "$BUILD_DIR/servin-gui" ]]; then
        cp "$BUILD_DIR/servin-gui" "$INSTALL_DIR/"
        chmod 755 "$INSTALL_DIR/servin-gui"
        print_success "  Installed: $INSTALL_DIR/servin-gui"
        
        # Create macOS App Bundle for GUI
        create_app_bundle
    fi
}

create_app_bundle() {
    print_info "Creating macOS App Bundle for Servin GUI..."
    
    local app_name="Servin Container Runtime"
    local app_dir="/Applications/${app_name}.app"
    local contents_dir="$app_dir/Contents"
    local macos_dir="$contents_dir/MacOS"
    local resources_dir="$contents_dir/Resources"
    
    # Create app bundle structure
    mkdir -p "$macos_dir"
    mkdir -p "$resources_dir"
    
    # Copy GUI binary to app bundle
    cp "$BUILD_DIR/servin-gui" "$macos_dir/ServinGUI"
    chmod 755 "$macos_dir/ServinGUI"
    
    # Create Info.plist
    cat > "$contents_dir/Info.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>ServinGUI</string>
    <key>CFBundleIdentifier</key>
    <string>com.servin.gui</string>
    <key>CFBundleName</key>
    <string>${app_name}</string>
    <key>CFBundleDisplayName</key>
    <string>${app_name}</string>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.0</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleSignature</key>
    <string>SRVN</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.12</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>NSRequiresAquaSystemAppearance</key>
    <false/>
</dict>
</plist>
EOF
    
    # Create app icon if available
    if [[ -f "$PROJECT_ROOT/icons/tool_icon.png" ]]; then
        cp "$PROJECT_ROOT/icons/tool_icon.png" "$resources_dir/AppIcon.png"
    fi
    
    # Set permissions
    chown -R root:wheel "$app_dir"
    chmod -R 755 "$app_dir"
    
    print_success "  Created: $app_dir"
}

create_config() {
    print_info "Creating configuration..."
    
    cat > "$CONFIG_DIR/servin.conf" << EOF
# Servin Configuration File
# Data directory
data_dir=$DATA_DIR

# Log settings
log_level=info
log_file=$LOG_DIR/servin.log

# Runtime settings
runtime=native

# Network settings
bridge_name=servin0

# CRI settings
cri_port=10250
cri_enabled=false
EOF
    
    chown root:wheel "$CONFIG_DIR/servin.conf"
    chmod 644 "$CONFIG_DIR/servin.conf"
    print_success "  Created: $CONFIG_DIR/servin.conf"
}

create_launchd_service() {
    print_info "Creating launchd service..."
    
    cat > "$LAUNCHD_DIR/$SERVICE_NAME.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>$SERVICE_NAME</string>
    <key>ProgramArguments</key>
    <array>
        <string>$INSTALL_DIR/servin</string>
        <string>daemon</string>
        <string>--config</string>
        <string>$CONFIG_DIR/servin.conf</string>
    </array>
    <key>UserName</key>
    <string>$USER</string>
    <key>GroupName</key>
    <string>staff</string>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <dict>
        <key>SuccessfulExit</key>
        <false/>
        <key>Crashed</key>
        <true/>
    </dict>
    <key>StandardOutPath</key>
    <string>$LOG_DIR/servin.stdout.log</string>
    <key>StandardErrorPath</key>
    <string>$LOG_DIR/servin.stderr.log</string>
    <key>WorkingDirectory</key>
    <string>$DATA_DIR</string>
    <key>EnvironmentVariables</key>
    <dict>
        <key>PATH</key>
        <string>/usr/local/bin:/usr/bin:/bin</string>
    </dict>
    <key>SoftResourceLimits</key>
    <dict>
        <key>NumberOfFiles</key>
        <integer>65536</integer>
        <key>NumberOfProcesses</key>
        <integer>2048</integer>
    </dict>
    <key>ThrottleInterval</key>
    <integer>10</integer>
</dict>
</plist>
EOF
    
    chown root:wheel "$LAUNCHD_DIR/$SERVICE_NAME.plist"
    chmod 644 "$LAUNCHD_DIR/$SERVICE_NAME.plist"
    
    # Load the service
    launchctl load "$LAUNCHD_DIR/$SERVICE_NAME.plist"
    
    print_success "  Created and loaded launchd service"
}

# App bundle creation is handled in install_binaries function above

create_uninstaller() {
    print_info "Creating uninstaller..."
    
    cat > "$INSTALL_DIR/servin-uninstall" << EOF
#!/bin/bash
# Servin Uninstaller for macOS

echo "Uninstalling Servin Container Runtime..."

# Stop and unload launchd service
launchctl unload "$LAUNCHD_DIR/$SERVICE_NAME.plist" 2>/dev/null || true
rm -f "$LAUNCHD_DIR/$SERVICE_NAME.plist"

# Remove user
dscl . -delete /Users/$USER 2>/dev/null || true

# Remove files and directories
rm -f "$INSTALL_DIR/servin"
rm -f "$INSTALL_DIR/servin-gui"
rm -f "$INSTALL_DIR/servin-uninstall"
rm -rf "$CONFIG_DIR"
rm -rf "$DATA_DIR"
rm -rf "$LOG_DIR"
rm -rf "/Applications/Servin GUI.app"

echo "Servin has been uninstalled."
EOF
    
    chmod +x "$INSTALL_DIR/servin-uninstall"
    print_success "  Created uninstaller: $INSTALL_DIR/servin-uninstall"
}

main() {
    print_header
    
    print_info "Detected platform: $PLATFORM"
    print_info "Build directory: $BUILD_DIR"
    
    check_root
    check_macos
    
    create_user
    create_directories
    install_binaries
    create_config
    create_launchd_service
    create_uninstaller
    
    echo ""
    print_success "================================================"
    print_success "   Installation completed successfully!"
    print_success "================================================"
    echo ""
    echo "Installation directory: $INSTALL_DIR"
    echo "Data directory: $DATA_DIR"
    echo "Configuration: $CONFIG_DIR/servin.conf"
    echo ""
    print_info "Next steps:"
    echo "1. Service is already running (started automatically)"
    echo "2. Check status: sudo launchctl list | grep servin"
    echo "3. View logs: tail -f $LOG_DIR/servin.log"
    echo "4. Use CLI: servin --help"
    
    if [[ -f "$INSTALL_DIR/servin-desktop" ]]; then
        echo "5. Use TUI: servin-desktop"
    fi
    
    if [[ -f "$INSTALL_DIR/servin-gui" ]]; then
        echo "6. Use GUI: Open 'Servin Container Runtime' from Applications folder"
        echo "   Or run: servin-gui"
    fi
    
    echo ""
    print_info "Installed binaries:"
    echo "  • servin (CLI runtime)"
    if [[ -f "$INSTALL_DIR/servin-desktop" ]]; then
        echo "  • servin-desktop (Terminal UI)"
    fi
    if [[ -f "$INSTALL_DIR/servin-gui" ]]; then
        echo "  • servin-gui (Graphical UI)"
    fi
    echo ""
    echo "Service will start automatically on boot."
    echo "To uninstall: sudo $INSTALL_DIR/servin-uninstall"
    
    # Show service status
    echo ""
    print_info "Service status:"
    if launchctl list | grep -q "$SERVICE_NAME"; then
        print_success "✓ Service is running"
    else
        print_warning "⚠ Service may not be running properly"
        echo "  Check logs at: $LOG_DIR/"
    fi
}

# Handle command line arguments
case "${1:-}" in
    uninstall)
        if [[ -f "$INSTALL_DIR/servin-uninstall" ]]; then
            exec "$INSTALL_DIR/servin-uninstall"
        else
            print_error "Uninstaller not found"
            exit 1
        fi
        ;;
    *)
        main
        ;;
esac
