#!/bin/bash
# Servin universal installer script
# Downloads and installs Servin with VM containerization support

set -e

# Configuration
INSTALL_URL="https://get.servin.dev"
GITHUB_REPO="immyemperor/servin"
INSTALL_DIR=""
FORCE_INSTALL=false
ENABLE_VM=true

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

print_header() {
    echo -e "${CYAN}================================================${NC}"
    echo -e "${CYAN}   Servin Universal Installer${NC}"
    echo -e "${CYAN}   VM-based Container Runtime${NC}" 
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

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Detect platform and architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case $os in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        *)
            print_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac
    
    case $arch in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
    
    PLATFORM="${OS}-${ARCH}"
    print_info "Detected platform: $PLATFORM"
}

# Check if running as root
check_root() {
    if [ "$EUID" -eq 0 ]; then
        IS_ROOT=true
        INSTALL_DIR="/usr/local"
        print_info "Installing system-wide as root"
    else
        IS_ROOT=false
        INSTALL_DIR="$HOME/.local"
        print_info "Installing to user directory"
    fi
}

# Check dependencies
check_dependencies() {
    print_info "Checking dependencies..."
    
    # Check curl or wget
    if command -v curl >/dev/null 2>&1; then
        DOWNLOAD_CMD="curl -fsSL"
    elif command -v wget >/dev/null 2>&1; then
        DOWNLOAD_CMD="wget -qO-"
    else
        print_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
    
    # Check for virtualization support
    check_virtualization_support
    
    print_success "Dependencies checked"
}

check_virtualization_support() {
    case $OS in
        linux)
            if [ -e /dev/kvm ]; then
                print_success "KVM virtualization support detected"
                VM_PROVIDER="kvm"
            elif command -v qemu-system-x86_64 >/dev/null 2>&1; then
                print_warning "KVM not available, QEMU will be used (slower)"
                VM_PROVIDER="qemu"
            else
                print_warning "No virtualization support detected"
                print_info "Install qemu-system-x86 for VM containerization"
                VM_PROVIDER="none"
            fi
            ;;
        darwin)
            if sw_vers -productVersion | grep -E "^1[1-9]\.|^[2-9][0-9]\." >/dev/null 2>&1; then
                print_success "Virtualization.framework support detected"
                VM_PROVIDER="hvf"
            else
                print_warning "macOS 11+ required for optimal VM support"
                VM_PROVIDER="qemu"
            fi
            ;;
    esac
}

# Get latest version
get_latest_version() {
    print_info "Getting latest version..."
    
    VERSION=$(${DOWNLOAD_CMD} "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | \
        grep '"tag_name":' | \
        sed -E 's/.*"tag_name": "([^"]+)".*/\1/')
    
    if [ -z "$VERSION" ]; then
        print_error "Failed to get latest version"
        exit 1
    fi
    
    print_success "Latest version: $VERSION"
}

# Download and extract Servin
download_servin() {
    print_info "Downloading Servin $VERSION for $PLATFORM..."
    
    local download_url="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/servin-${VERSION}-${PLATFORM}.tar.gz"
    local temp_dir=$(mktemp -d)
    local archive_path="$temp_dir/servin.tar.gz"
    
    # Download
    if ! ${DOWNLOAD_CMD} "$download_url" > "$archive_path"; then
        print_error "Failed to download Servin"
        rm -rf "$temp_dir"
        exit 1
    fi
    
    print_success "Downloaded Servin"
    
    # Extract
    print_info "Extracting archive..."
    if ! tar -xzf "$archive_path" -C "$temp_dir"; then
        print_error "Failed to extract archive"
        rm -rf "$temp_dir"
        exit 1
    fi
    
    EXTRACT_DIR="$temp_dir/servin-${VERSION}-${PLATFORM}"
    
    if [ ! -d "$EXTRACT_DIR" ]; then
        print_error "Expected directory not found after extraction"
        rm -rf "$temp_dir"
        exit 1
    fi
    
    print_success "Extracted Servin"
}

# Install Servin
install_servin() {
    print_info "Installing Servin to $INSTALL_DIR..."
    
    # Create directories
    mkdir -p "$INSTALL_DIR/bin"
    mkdir -p "$INSTALL_DIR/share/servin"
    
    if [ "$IS_ROOT" = true ]; then
        mkdir -p "$INSTALL_DIR/etc/servin"
        CONFIG_DIR="$INSTALL_DIR/etc/servin"
    else
        mkdir -p "$HOME/.config/servin"
        CONFIG_DIR="$HOME/.config/servin"
    fi
    
    # Copy binaries
    cp "$EXTRACT_DIR/servin" "$INSTALL_DIR/bin/"
    chmod +x "$INSTALL_DIR/bin/servin"
    
    if [ -f "$EXTRACT_DIR/servin-desktop" ]; then
        cp "$EXTRACT_DIR/servin-desktop" "$INSTALL_DIR/bin/"
        chmod +x "$INSTALL_DIR/bin/servin-desktop"
    fi
    
    # Copy resources
    if [ -d "$EXTRACT_DIR/webview_gui" ]; then
        cp -r "$EXTRACT_DIR/webview_gui" "$INSTALL_DIR/share/servin/"
    fi
    
    # Copy configuration
    if [ -f "$EXTRACT_DIR/vm-config.json" ]; then
        cp "$EXTRACT_DIR/vm-config.json" "$CONFIG_DIR/"
    fi
    
    # Copy documentation
    if [ -f "$EXTRACT_DIR/README.md" ]; then
        cp "$EXTRACT_DIR/README.md" "$INSTALL_DIR/share/servin/"
    fi
    
    print_success "Installed Servin binaries and resources"
}

# Configure PATH
configure_path() {
    if [[ ":$PATH:" != *":$INSTALL_DIR/bin:"* ]]; then
        print_info "Adding Servin to PATH..."
        
        # Add to shell profiles
        for shell_profile in "$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.profile"; do
            if [ -f "$shell_profile" ] || [ "$shell_profile" = "$HOME/.profile" ]; then
                if ! grep -q "$INSTALL_DIR/bin" "$shell_profile" 2>/dev/null; then
                    echo "export PATH=\"$INSTALL_DIR/bin:\$PATH\"" >> "$shell_profile"
                fi
            fi
        done
        
        # Export for current session
        export PATH="$INSTALL_DIR/bin:$PATH"
        
        print_success "Added Servin to PATH"
        print_warning "You may need to restart your shell or run: source ~/.bashrc"
    else
        print_success "Servin already in PATH"
    fi
}

# Install VM dependencies
install_vm_dependencies() {
    if [ "$VM_PROVIDER" = "none" ]; then
        print_warning "Skipping VM dependency installation"
        return
    fi
    
    print_info "Installing VM dependencies for $OS..."
    
    case $OS in
        linux)
            if command -v apt >/dev/null 2>&1; then
                print_info "Installing QEMU via apt..."
                if [ "$IS_ROOT" = true ]; then
                    apt update && apt install -y qemu-system-x86 qemu-utils
                else
                    print_warning "Run 'sudo apt install qemu-system-x86 qemu-utils' to install VM support"
                fi
            elif command -v yum >/dev/null 2>&1; then
                print_info "Installing QEMU via yum..."
                if [ "$IS_ROOT" = true ]; then
                    yum install -y qemu-kvm qemu-img
                else
                    print_warning "Run 'sudo yum install qemu-kvm qemu-img' to install VM support"
                fi
            else
                print_warning "Please install QEMU manually for VM support"
            fi
            ;;
        darwin)
            if command -v brew >/dev/null 2>&1; then
                print_info "Installing QEMU via Homebrew..."
                brew install qemu
            else
                print_warning "Install Homebrew and run 'brew install qemu' for VM support"
            fi
            ;;
    esac
}

# Setup VM mode
setup_vm_mode() {
    if [ "$ENABLE_VM" != true ]; then
        return
    fi
    
    print_info "Setting up VM mode..."
    
    # Create VM directory
    mkdir -p "$HOME/.servin/vms"
    
    # Enable VM mode
    export SERVIN_VM_MODE=true
    echo "export SERVIN_VM_MODE=true" >> "$HOME/.bashrc"
    echo "export SERVIN_VM_MODE=true" >> "$HOME/.zshrc" 2>/dev/null || true
    
    print_success "VM mode enabled"
    print_info "Run 'servin vm enable' to complete VM setup"
}

# Verify installation
verify_installation() {
    print_info "Verifying installation..."
    
    if ! command -v servin >/dev/null 2>&1; then
        print_error "Servin not found in PATH"
        print_info "You may need to restart your shell or run: source ~/.bashrc"
        return 1
    fi
    
    # Test basic functionality
    if ! servin --version >/dev/null 2>&1; then
        print_error "Servin version check failed"
        return 1
    fi
    
    local installed_version=$(servin --version | head -n1 | cut -d' ' -f2)
    print_success "Servin $installed_version installed successfully"
    
    # Test VM capabilities
    if [ "$ENABLE_VM" = true ]; then
        print_info "Testing VM capabilities..."
        if servin vm status >/dev/null 2>&1; then
            print_success "VM functionality available"
        else
            print_warning "VM functionality limited (install dependencies for full support)"
        fi
    fi
    
    return 0
}

# Print next steps
print_next_steps() {
    echo ""
    echo -e "${CYAN}================================================${NC}"
    echo -e "${CYAN}   Installation Complete!${NC}"
    echo -e "${CYAN}================================================${NC}"
    echo ""
    echo -e "${GREEN}Next steps:${NC}"
    echo ""
    
    if [[ ":$PATH:" != *":$INSTALL_DIR/bin:"* ]]; then
        echo "1. Restart your shell or run:"
        echo -e "   ${YELLOW}source ~/.bashrc${NC}"
        echo ""
    fi
    
    if [ "$ENABLE_VM" = true ]; then
        echo "2. Enable VM-based containerization:"
        echo -e "   ${YELLOW}servin vm enable${NC}"
        echo -e "   ${YELLOW}servin vm start${NC}"
        echo ""
    fi
    
    echo "3. Run your first container:"
    echo -e "   ${YELLOW}servin run alpine echo \"Hello from Servin!\"${NC}"
    echo ""
    
    echo "4. Check status:"
    echo -e "   ${YELLOW}servin vm status${NC}"
    echo -e "   ${YELLOW}servin ls${NC}"
    echo ""
    
    echo -e "${GREEN}Documentation:${NC}"
    echo "• Getting Started: https://servin.dev/docs/getting-started"
    echo "• VM Guide: https://servin.dev/docs/vm-containerization"
    echo "• GitHub: https://github.com/$GITHUB_REPO"
    echo ""
    
    echo -e "${GREEN}Support:${NC}"
    echo "• Issues: https://github.com/$GITHUB_REPO/issues"
    echo "• Discussions: https://github.com/$GITHUB_REPO/discussions"
    echo ""
}

# Cleanup
cleanup() {
    if [ -n "$temp_dir" ] && [ -d "$temp_dir" ]; then
        rm -rf "$temp_dir"
    fi
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --install-dir)
                INSTALL_DIR="$2"
                shift 2
                ;;
            --force)
                FORCE_INSTALL=true
                shift
                ;;
            --no-vm)
                ENABLE_VM=false
                shift
                ;;
            --help)
                show_help
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

show_help() {
    cat << EOF
Servin Universal Installer

USAGE:
    curl -fsSL https://get.servin.dev | sh
    
    or with options:
    curl -fsSL https://get.servin.dev | sh -s -- [OPTIONS]

OPTIONS:
    --install-dir DIR    Install to specific directory (default: auto-detect)
    --force             Force installation even if already installed  
    --no-vm             Skip VM mode setup
    --help              Show this help message

EXAMPLES:
    # Basic installation
    curl -fsSL https://get.servin.dev | sh
    
    # Install to custom directory
    curl -fsSL https://get.servin.dev | sh -s -- --install-dir /opt/servin
    
    # Install without VM mode
    curl -fsSL https://get.servin.dev | sh -s -- --no-vm

EOF
}

# Main installation flow
main() {
    # Setup trap for cleanup
    trap cleanup EXIT
    
    print_header
    
    parse_args "$@"
    detect_platform
    check_root
    check_dependencies
    get_latest_version
    download_servin
    install_servin
    configure_path
    install_vm_dependencies
    setup_vm_mode
    
    if verify_installation; then
        print_next_steps
    else
        print_error "Installation verification failed"
        exit 1
    fi
}

# Run main function with all arguments
main "$@"