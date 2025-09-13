#!/bin/bash
# Servin Container Runtime - macOS Installer Script

set -e

# Configuration
INSTALL_DIR="/usr/local/bin"
DATA_DIR="/usr/local/var/lib/servin"
CONFIG_DIR="/usr/local/etc/servin"
LOG_DIR="/usr/local/var/log/servin"
LAUNCHD_DIR="/Library/LaunchDaemons"
SERVICE_NAME="com.servin.runtime"
USER="_servin"

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
    print_info "Installing binaries..."
    
    # Check if binaries exist in current directory
    if [[ ! -f "servin" ]]; then
        print_error "servin binary not found in current directory"
        exit 1
    fi
    
    # Copy binaries
    cp servin "$INSTALL_DIR/"
    chmod 755 "$INSTALL_DIR/servin"
    print_success "  Installed: $INSTALL_DIR/servin"
    
    if [[ -f "servin-gui" ]]; then
        cp servin-gui "$INSTALL_DIR/"
        chmod 755 "$INSTALL_DIR/servin-gui"
        print_success "  Installed: $INSTALL_DIR/servin-gui"
    fi
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

create_app_bundle() {
    print_info "Creating application bundle..."
    
    if [[ ! -f "$INSTALL_DIR/servin-gui" ]]; then
        print_warning "  servin-gui not found, skipping app bundle creation"
        return
    fi
    
    local app_dir="/Applications/Servin GUI.app"
    local contents_dir="$app_dir/Contents"
    local macos_dir="$contents_dir/MacOS"
    local resources_dir="$contents_dir/Resources"
    
    # Create directory structure
    mkdir -p "$macos_dir" "$resources_dir"
    
    # Copy executable
    cp "$INSTALL_DIR/servin-gui" "$macos_dir/Servin GUI"
    chmod +x "$macos_dir/Servin GUI"
    
    # Create Info.plist
    cat > "$contents_dir/Info.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>Servin GUI</string>
    <key>CFBundleIdentifier</key>
    <string>com.servin.gui</string>
    <key>CFBundleName</key>
    <string>Servin GUI</string>
    <key>CFBundleDisplayName</key>
    <string>Servin GUI</string>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.0</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleSignature</key>
    <string>SERV</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.12</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>LSApplicationCategoryType</key>
    <string>public.app-category.developer-tools</string>
</dict>
</plist>
EOF
    
    # Create a simple icon (ASCII art in icns format would be complex, so we'll skip for now)
    # In production, you'd want to include a proper .icns file
    
    print_success "  Created application bundle: $app_dir"
}

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
    
    check_root
    check_macos
    
    create_user
    create_directories
    install_binaries
    create_config
    create_launchd_service
    create_app_bundle
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
    if [[ -f "$INSTALL_DIR/servin-gui" ]]; then
        echo "4. Use GUI: Open 'Servin GUI' from Applications folder"
        echo "   Or run: servin-gui"
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
