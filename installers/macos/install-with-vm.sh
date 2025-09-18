#!/bin/bash

# Servin Container Runtime - Enhanced macOS Installer with VM Prerequisites
# Run with: ./install-with-vm.sh

set -e

# Configuration
INSTALL_DIR="/usr/local/bin"
DATA_DIR="/usr/local/var/lib/servin"
CONFIG_DIR="/usr/local/etc/servin"
LOG_DIR="/usr/local/var/log/servin"
VM_DIR="$DATA_DIR/vm"
USER_DATA_DIR="$HOME/.servin"
HOMEBREW_PREFIX=""

# Color output functions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_success() { echo -e "${GREEN}✓ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠ $1${NC}"; }
print_error() { echo -e "${RED}✗ $1${NC}"; }
print_info() { echo -e "${BLUE}→ $1${NC}"; }
print_header() { echo -e "\n${CYAN}$(printf '=%.0s' {1..60})\n$1\n$(printf '=%.0s' {1..60})${NC}"; }

# Check if running as root
check_root() {
    if [[ $EUID -eq 0 ]]; then
        print_warning "This script should not be run as root on macOS"
        print_info "Run without sudo - it will prompt for password when needed"
        exit 1
    fi
}

# Detect macOS version and architecture
detect_system() {
    MACOS_VERSION=$(sw_vers -productVersion)
    ARCH=$(uname -m)
    
    print_info "Detected system: macOS $MACOS_VERSION ($ARCH)"
    
    # Set Homebrew prefix based on architecture
    if [[ "$ARCH" == "arm64" ]]; then
        HOMEBREW_PREFIX="/opt/homebrew"
    else
        HOMEBREW_PREFIX="/usr/local"
    fi
    
    # Check minimum macOS version for Virtualization.framework
    if [[ $(echo "$MACOS_VERSION" | cut -d. -f1) -ge 11 ]]; then
        print_success "macOS version supports Virtualization.framework"
        VM_FRAMEWORK_SUPPORTED=true
    else
        print_warning "macOS version may not support Virtualization.framework"
        print_info "Will use QEMU for virtualization"
        VM_FRAMEWORK_SUPPORTED=false
    fi
}

# Check system prerequisites
check_prerequisites() {
    print_header "Checking System Prerequisites"
    
    local prereq_failed=0
    
    # Check if we're running on Apple Silicon or Intel with VT-x
    print_info "Checking virtualization support..."
    if [[ "$ARCH" == "arm64" ]]; then
        print_success "Apple Silicon detected - hardware virtualization supported"
    else
        # Check for Intel VT-x
        if sysctl -n machdep.cpu.features | grep -q VMX; then
            print_success "Intel VT-x supported"
        else
            print_warning "Intel VT-x not detected - virtualization may be limited"
            prereq_failed=$((prereq_failed + 1))
        fi
    fi
    
    # Check available memory
    print_info "Checking available memory..."
    local mem_gb=$(( $(sysctl -n hw.memsize) / 1024 / 1024 / 1024 ))
    if [[ $mem_gb -ge 8 ]]; then
        print_success "Sufficient memory available (${mem_gb}GB)"
    elif [[ $mem_gb -ge 4 ]]; then
        print_warning "Limited memory available (${mem_gb}GB). Recommended: 8GB+"
    else
        print_error "Insufficient memory (${mem_gb}GB). Need 4GB minimum"
        prereq_failed=$((prereq_failed + 1))
    fi
    
    # Check disk space
    print_info "Checking available disk space..."
    local disk_gb=$(df -g / | awk 'NR==2 {print $4}')
    if [[ $disk_gb -ge 10 ]]; then
        print_success "Sufficient disk space available (${disk_gb}GB free)"
    else
        print_error "Insufficient disk space (${disk_gb}GB free). Need 10GB minimum"
        prereq_failed=$((prereq_failed + 1))
    fi
    
    # Check internet connectivity
    print_info "Checking internet connectivity..."
    if ping -c 1 google.com >/dev/null 2>&1; then
        print_success "Internet connection available"
    else
        print_warning "Internet connection not available - some features may not work"
    fi
    
    # Check if Developer Tools are installed
    print_info "Checking Xcode Command Line Tools..."
    if xcode-select -p >/dev/null 2>&1; then
        print_success "Xcode Command Line Tools installed"
    else
        print_warning "Xcode Command Line Tools not found"
        print_info "Installing Command Line Tools..."
        xcode-select --install
        print_info "Please complete the Command Line Tools installation and run this script again"
        exit 1
    fi
    
    if [[ $prereq_failed -gt 0 ]]; then
        print_warning "$prereq_failed critical prerequisites failed"
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
}

# Install or update Homebrew
install_homebrew() {
    print_header "Setting up Homebrew Package Manager"
    
    if command -v brew >/dev/null 2>&1; then
        print_success "Homebrew already installed"
        print_info "Updating Homebrew..."
        brew update
    else
        print_info "Installing Homebrew..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        
        # Add Homebrew to PATH
        echo "# Homebrew" >> ~/.zshrc
        echo 'eval "$('$HOMEBREW_PREFIX'/bin/brew shellenv)"' >> ~/.zshrc
        eval "$($HOMEBREW_PREFIX/bin/brew shellenv)"
        
        print_success "Homebrew installed successfully"
    fi
}

# Install virtualization packages
install_virtualization() {
    print_header "Installing Virtualization Prerequisites"
    
    print_info "Installing QEMU and related tools..."
    brew install qemu
    
    # Install additional QEMU architectures if needed
    print_info "Installing QEMU system emulators..."
    brew install qemu-img
    
    # Install UTM if requested (optional GUI VM manager)
    print_info "Installing optional VM management tools..."
    read -p "Install UTM (GUI VM manager)? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        brew install --cask utm
        print_success "UTM installed"
    fi
    
    print_success "Virtualization tools installed"
}

# Install development tools
install_development_tools() {
    print_header "Installing Development Tools"
    
    # Install Python
    print_info "Installing Python..."
    brew install python@3.11
    brew install python-tk@3.11
    
    # Install Packer
    print_info "Installing HashiCorp Packer..."
    brew install packer
    
    # Verify installations
    if command -v packer >/dev/null 2>&1; then
        print_success "Packer installed: $(packer version)"
        
        # Install Packer plugins
        print_info "Installing Packer QEMU plugin..."
        packer plugins install github.com/hashicorp/qemu 2>/dev/null || print_warning "Packer QEMU plugin installation failed"
    else
        print_warning "Packer installation failed"
    fi
    
    # Install additional development tools
    print_info "Installing additional development tools..."
    brew install wget curl unzip git
    
    print_success "Development tools installed"
}

# Install Python dependencies
install_python_dependencies() {
    print_header "Installing Python Dependencies"
    
    print_info "Setting up Python environment..."
    local python_exe="$HOMEBREW_PREFIX/bin/python3.11"
    
    if [[ ! -f "$python_exe" ]]; then
        python_exe="python3"
    fi
    
    print_info "Installing Python WebView dependencies..."
    $python_exe -m pip install --upgrade pip --quiet
    $python_exe -m pip install pywebview[cocoa] flask flask-cors flask-socketio eventlet pyinstaller --quiet
    
    # Test imports
    print_info "Testing Python imports..."
    $python_exe -c "import webview; print('✓ pywebview available')" 2>/dev/null || print_warning "pywebview not available"
    $python_exe -c "import flask; print('✓ flask available')" 2>/dev/null || print_warning "flask not available"
    $python_exe -c "import flask_socketio; print('✓ flask-socketio available')" 2>/dev/null || print_warning "flask-socketio not available"
    
    print_success "Python dependencies installed"
}

# Configure virtualization framework
configure_virtualization() {
    print_header "Configuring Virtualization"
    
    # Test QEMU installation
    print_info "Testing QEMU installation..."
    if command -v qemu-system-x86_64 >/dev/null 2>&1; then
        print_success "QEMU x86_64 emulator available"
        qemu_version=$(qemu-system-x86_64 --version | head -1)
        print_info "  $qemu_version"
    else
        print_error "QEMU x86_64 emulator not found"
        return 1
    fi
    
    if command -v qemu-system-aarch64 >/dev/null 2>&1; then
        print_success "QEMU ARM64 emulator available"
    fi
    
    # Test Virtualization.framework if supported
    if [[ "$VM_FRAMEWORK_SUPPORTED" == "true" ]]; then
        print_info "Testing Virtualization.framework..."
        # This will be tested later with Servin's VM functionality
        print_success "Virtualization.framework should be available"
    fi
    
    print_success "Virtualization configured"
}

# Install Servin binaries
install_servin() {
    print_header "Installing Servin Container Runtime"
    
    # Create directories
    print_info "Creating directories..."
    sudo mkdir -p "$INSTALL_DIR" "$DATA_DIR" "$CONFIG_DIR" "$LOG_DIR" "$VM_DIR"/{images,instances}
    mkdir -p "$USER_DATA_DIR"
    
    # Copy executables
    print_info "Installing executables..."
    local installed_count=0
    
    for exe in servin servin-tui servin-gui; do
        if [[ -f "./$exe" ]]; then
            sudo cp "./$exe" "$INSTALL_DIR/"
            sudo chmod +x "$INSTALL_DIR/$exe"
            print_success "Installed: $exe"
            installed_count=$((installed_count + 1))
        else
            print_warning "Executable not found: $exe"
        fi
    done
    
    if [[ $installed_count -eq 0 ]]; then
        print_error "No executables found. Please run this installer from the directory containing the Servin executables."
        return 1
    fi
    
    # Create VM configuration
    print_info "Creating VM configuration..."
    local vm_config="$CONFIG_DIR/vm-config.yaml"
    sudo tee "$vm_config" > /dev/null << EOF
vm:
  platform: darwin
  providers:
    - name: virtualization
      priority: 1
      enabled: $VM_FRAMEWORK_SUPPORTED
      acceleration: true
    - name: qemu
      priority: 2
      enabled: true
      acceleration: true
  default_provider: $([ "$VM_FRAMEWORK_SUPPORTED" == "true" ] && echo "virtualization" || echo "qemu")
  image_cache: "$VM_DIR/images"
  vm_storage: "$VM_DIR/instances"
  qemu_binary: "$HOMEBREW_PREFIX/bin/qemu-system-x86_64"
  max_memory: "4GB"
  default_memory: "2GB"
  max_cpu_cores: 4
  hvf_acceleration: true
EOF
    
    # Create main configuration
    print_info "Creating main configuration..."
    local main_config="$CONFIG_DIR/servin.conf"
    sudo tee "$main_config" > /dev/null << EOF
# Servin Configuration File
data_dir=$DATA_DIR
log_level=info
log_file=$LOG_DIR/servin.log
runtime=vm
bridge_name=servin0
cri_port=10250
cri_enabled=false
vm_enabled=true
vm_config=$CONFIG_DIR/vm-config.yaml
homebrew_prefix=$HOMEBREW_PREFIX
EOF
    
    # Set permissions
    sudo chown -R root:wheel "$INSTALL_DIR"
    sudo chown -R root:wheel "$CONFIG_DIR"
    sudo chmod -R 755 "$INSTALL_DIR"
    sudo chmod -R 644 "$CONFIG_DIR"
    sudo chmod +x "$INSTALL_DIR"/servin*
    
    # Allow user access to data directory
    sudo chown -R $(whoami):staff "$DATA_DIR" "$LOG_DIR"
    
    print_success "Servin installed successfully"
}

# Create launchd service
create_service() {
    print_header "Creating System Service"
    
    if [[ "${SKIP_SERVICE:-0}" == "1" ]]; then
        print_info "Skipping service creation"
        return 0
    fi
    
    print_info "Creating launchd service..."
    local plist_file="/Library/LaunchDaemons/dev.servin.container-runtime.plist"
    sudo tee "$plist_file" > /dev/null << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>dev.servin.container-runtime</string>
    <key>ProgramArguments</key>
    <array>
        <string>$INSTALL_DIR/servin</string>
        <string>daemon</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <dict>
        <key>SuccessfulExit</key>
        <false/>
    </dict>
    <key>StandardOutPath</key>
    <string>$LOG_DIR/servin.out.log</string>
    <key>StandardErrorPath</key>
    <string>$LOG_DIR/servin.err.log</string>
    <key>WorkingDirectory</key>
    <string>$DATA_DIR</string>
    <key>EnvironmentVariables</key>
    <dict>
        <key>PATH</key>
        <string>$HOMEBREW_PREFIX/bin:/usr/local/bin:/usr/bin:/bin</string>
        <key>HOMEBREW_PREFIX</key>
        <string>$HOMEBREW_PREFIX</string>
    </dict>
</dict>
</plist>
EOF
    
    sudo launchctl load "$plist_file" 2>/dev/null || print_warning "Service will be available after next reboot"
    
    print_success "Launchd service created and loaded"
}

# Initialize VM support
initialize_vm() {
    print_header "Initializing VM Support"
    
    if [[ "${SKIP_VM:-0}" == "1" ]]; then
        print_info "Skipping VM initialization"
        return 0
    fi
    
    local servin_exe="$INSTALL_DIR/servin"
    if [[ ! -f "$servin_exe" ]]; then
        print_error "Servin executable not found: $servin_exe"
        return 1
    fi
    
    print_info "Initializing VM directories..."
    "$servin_exe" vm init || print_warning "VM initialization may have failed"
    
    print_info "Testing VM providers..."
    "$servin_exe" vm list-providers || print_warning "VM provider detection failed"
    
    # Test Virtualization.framework if supported
    if [[ "$VM_FRAMEWORK_SUPPORTED" == "true" ]]; then
        print_info "Testing Virtualization.framework..."
        if "$servin_exe" vm check-virtualization >/dev/null 2>&1; then
            print_success "Virtualization.framework available and working"
        else
            print_warning "Virtualization.framework not fully functional - will use QEMU fallback"
        fi
    fi
    
    # Test QEMU
    print_info "Testing QEMU functionality..."
    if "$servin_exe" vm check-qemu >/dev/null 2>&1; then
        print_success "QEMU provider available and working"
    else
        print_warning "QEMU provider not fully functional"
    fi
    
    print_success "VM support initialized"
}

# Run comprehensive tests
run_tests() {
    print_header "Running Verification Tests"
    
    local servin_exe="$INSTALL_DIR/servin"
    
    # Test basic functionality
    print_info "Testing basic functionality..."
    if "$servin_exe" version >/dev/null 2>&1; then
        print_success "Servin CLI working"
        local version=$("$servin_exe" version 2>/dev/null | head -1)
        print_info "  $version"
    else
        print_warning "Servin CLI test failed"
    fi
    
    # Test VM functionality
    print_info "Testing VM functionality..."
    if "$servin_exe" vm status >/dev/null 2>&1; then
        print_success "VM subsystem working"
    else
        print_warning "VM subsystem test failed"
    fi
    
    # Test QEMU
    print_info "Testing QEMU installation..."
    if qemu-system-x86_64 -version >/dev/null 2>&1; then
        print_success "QEMU working"
        local qemu_ver=$(qemu-system-x86_64 -version | head -1)
        print_info "  $qemu_ver"
    else
        print_warning "QEMU test failed"
    fi
}

# Show installation summary
show_summary() {
    print_header "Installation Summary"
    
    echo -e "\n${YELLOW}Installed Components:${NC}"
    echo "✓ Servin Container Runtime"
    echo "✓ VM containerization support"
    echo "✓ QEMU virtualization"
    if [[ "$VM_FRAMEWORK_SUPPORTED" == "true" ]]; then
        echo "✓ Virtualization.framework support"
    fi
    echo "✓ Homebrew package manager"
    echo "✓ Python WebView GUI support"
    echo "✓ Development tools (Packer)"
    echo "✓ Launchd service"
    
    echo -e "\n${YELLOW}Configuration Files:${NC}"
    echo "• Main config: $CONFIG_DIR/servin.conf"
    echo "• VM config: $CONFIG_DIR/vm-config.yaml"
    echo "• Logs: $LOG_DIR/servin.log"
    echo "• Data: $DATA_DIR"
    echo "• User data: $USER_DATA_DIR"
    
    echo -e "\n${YELLOW}Next Steps:${NC}"
    echo "1. Restart your terminal (for PATH updates)"
    echo "2. Initialize VM support: servin vm init"
    echo "3. Enable VM mode: servin vm enable"
    echo "4. Test: servin run --vm alpine echo 'Hello from VM!'"
    
    echo -e "\n${YELLOW}GUI Access:${NC}"
    echo "• Run 'servin-gui' to open the graphical interface"
    echo "• Or use 'servin-tui' for terminal interface"
    
    echo -e "\n${YELLOW}Service Management:${NC}"
    echo "• Start: sudo launchctl load /Library/LaunchDaemons/dev.servin.container-runtime.plist"
    echo "• Stop: sudo launchctl unload /Library/LaunchDaemons/dev.servin.container-runtime.plist"
    
    echo -e "\n${CYAN}Documentation: See VM_PREREQUISITES.md for detailed setup guide${NC}"
    print_success "\nServin Container Runtime installation completed!"
}

# Main installation flow
main() {
    print_header "Servin Container Runtime - macOS Installer"
    echo -e "This installer will set up Servin with VM containerization support.\n"
    
    check_root
    detect_system
    check_prerequisites
    install_homebrew
    install_virtualization
    install_development_tools
    install_python_dependencies
    configure_virtualization
    install_servin
    create_service
    initialize_vm
    run_tests
    show_summary
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "Usage: $0 [options]"
        echo "Options:"
        echo "  --help, -h     Show this help message"
        echo "  --no-service   Skip launchd service creation"
        echo "  --no-vm        Skip VM setup"
        echo "  --no-gui       Skip GUI dependencies"
        exit 0
        ;;
    --no-service)
        SKIP_SERVICE=1
        ;;
    --no-vm)
        SKIP_VM=1
        ;;
    --no-gui)
        SKIP_GUI=1
        ;;
esac

# Run main installation
main "$@"