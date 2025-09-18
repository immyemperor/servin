#!/bin/bash

# Servin Container Runtime - Smart Installation Wizard
# Auto-detects platform and installs prerequisites if needed
# Run with: ./install-wizard.sh

set -e

# Color output functions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

print_success() { echo -e "${GREEN}✓ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠ $1${NC}"; }
print_error() { echo -e "${RED}✗ $1${NC}"; }
print_info() { echo -e "${BLUE}→ $1${NC}"; }
print_header() { echo -e "\n${CYAN}${BOLD}$1${NC}"; }
print_banner() {
    echo -e "${CYAN}${BOLD}"
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║                    Servin Container Runtime                    ║"
    echo "║                  Smart Installation Wizard                    ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}\n"
}

# Detect platform
detect_platform() {
    case "$(uname -s)" in
        Linux*)
            PLATFORM="linux"
            DISTRO_ID=""
            if [[ -f /etc/os-release ]]; then
                . /etc/os-release
                DISTRO_ID=$ID
            fi
            ;;
        Darwin*)
            PLATFORM="macos"
            MACOS_VERSION=$(sw_vers -productVersion)
            ARCH=$(uname -m)
            ;;
        CYGWIN*|MINGW*|MSYS*)
            print_error "Please run install-wizard.ps1 on Windows"
            exit 1
            ;;
        *)
            print_error "Unsupported platform: $(uname -s)"
            exit 1
            ;;
    esac
    
    print_info "Detected platform: $PLATFORM"
    if [[ "$PLATFORM" == "linux" && -n "$DISTRO_ID" ]]; then
        print_info "Distribution: $DISTRO_ID"
    elif [[ "$PLATFORM" == "macos" ]]; then
        print_info "macOS version: $MACOS_VERSION ($ARCH)"
    fi
}

# Check if we're in the right directory
check_directory() {
    local script_dir=$(dirname "$(readlink -f "$0" 2>/dev/null || realpath "$0")")
    
    if [[ ! -f "$script_dir/README.md" ]] || ! grep -q "Servin Container Runtime" "$script_dir/README.md" 2>/dev/null; then
        print_error "This script must be run from the Servin installers directory"
        print_info "Expected to find installers/README.md with Servin content"
        exit 1
    fi
    
    INSTALLER_DIR="$script_dir"
    print_success "Running from correct directory: $INSTALLER_DIR"
}

# Check for existing Servin installation
check_existing_installation() {
    print_header "Checking Existing Installation"
    
    local servin_found=false
    local servin_paths=("/usr/local/bin/servin" "/usr/bin/servin" "$HOME/.local/bin/servin")
    
    for path in "${servin_paths[@]}"; do
        if [[ -x "$path" ]]; then
            print_success "Found existing Servin installation: $path"
            local version=$("$path" version 2>/dev/null | head -1 || echo "Unknown version")
            print_info "Version: $version"
            servin_found=true
            EXISTING_SERVIN="$path"
            break
        fi
    done
    
    if [[ "$servin_found" == "false" ]]; then
        print_info "No existing Servin installation found"
        EXISTING_SERVIN=""
    fi
}

# Check VM prerequisites
check_vm_prerequisites() {
    print_header "Checking VM Prerequisites"
    
    local prereq_missing=false
    
    case "$PLATFORM" in
        "linux")
            # Check KVM/QEMU
            if ! command -v qemu-system-x86_64 >/dev/null 2>&1; then
                print_warning "QEMU not found"
                prereq_missing=true
            else
                print_success "QEMU available"
            fi
            
            if [[ ! -e /dev/kvm ]]; then
                print_warning "KVM device not available"
                prereq_missing=true
            else
                print_success "KVM device available"
            fi
            
            if ! command -v virsh >/dev/null 2>&1; then
                print_warning "libvirt not found"
                prereq_missing=true
            else
                print_success "libvirt available"
            fi
            ;;
            
        "macos")
            # Check QEMU
            if ! command -v qemu-system-x86_64 >/dev/null 2>&1; then
                print_warning "QEMU not found"
                prereq_missing=true
            else
                print_success "QEMU available"
            fi
            
            # Check Homebrew
            if ! command -v brew >/dev/null 2>&1; then
                print_warning "Homebrew not found"
                prereq_missing=true
            else
                print_success "Homebrew available"
            fi
            
            # Check Xcode Command Line Tools
            if ! xcode-select -p >/dev/null 2>&1; then
                print_warning "Xcode Command Line Tools not found"
                prereq_missing=true
            else
                print_success "Xcode Command Line Tools available"
            fi
            ;;
    esac
    
    # Check Python
    if ! command -v python3 >/dev/null 2>&1; then
        print_warning "Python 3 not found"
        prereq_missing=true
    else
        local python_version=$(python3 --version 2>&1 | cut -d' ' -f2)
        print_success "Python 3 available: $python_version"
    fi
    
    PREREQ_MISSING=$prereq_missing
}

# Check system resources
check_system_resources() {
    print_header "Checking System Resources"
    
    local resource_warnings=0
    
    # Check memory
    case "$PLATFORM" in
        "linux")
            local mem_gb=$(free -g | awk 'NR==2{printf "%.1f", $2}')
            ;;
        "macos")
            local mem_gb=$(( $(sysctl -n hw.memsize) / 1024 / 1024 / 1024 ))
            ;;
    esac
    
    if (( $(echo "$mem_gb >= 8" | bc -l 2>/dev/null || echo "0") )); then
        print_success "Memory: ${mem_gb}GB (excellent)"
    elif (( $(echo "$mem_gb >= 4" | bc -l 2>/dev/null || echo "0") )); then
        print_warning "Memory: ${mem_gb}GB (minimum met, 8GB+ recommended)"
        resource_warnings=$((resource_warnings + 1))
    else
        print_error "Memory: ${mem_gb}GB (insufficient, 4GB minimum required)"
        resource_warnings=$((resource_warnings + 1))
    fi
    
    # Check disk space
    case "$PLATFORM" in
        "linux")
            local disk_gb=$(df / | awk 'NR==2 {printf "%.1f", $4/1024/1024}')
            ;;
        "macos")
            local disk_gb=$(df -g / | awk 'NR==2 {print $4}')
            ;;
    esac
    
    if (( $(echo "$disk_gb >= 10" | bc -l 2>/dev/null || echo "0") )); then
        print_success "Disk space: ${disk_gb}GB free (excellent)"
    elif (( $(echo "$disk_gb >= 5" | bc -l 2>/dev/null || echo "0") )); then
        print_warning "Disk space: ${disk_gb}GB free (minimum met, 10GB+ recommended)"
        resource_warnings=$((resource_warnings + 1))
    else
        print_error "Disk space: ${disk_gb}GB free (insufficient, 5GB minimum required)"
        resource_warnings=$((resource_warnings + 1))
    fi
    
    RESOURCE_WARNINGS=$resource_warnings
}

# Get enhanced installer path
get_enhanced_installer() {
    case "$PLATFORM" in
        "linux")
            ENHANCED_INSTALLER="$INSTALLER_DIR/linux/install-with-vm.sh"
            ;;
        "macos")
            ENHANCED_INSTALLER="$INSTALLER_DIR/macos/install-with-vm.sh"
            ;;
    esac
    
    if [[ ! -f "$ENHANCED_INSTALLER" ]]; then
        print_error "Enhanced installer not found: $ENHANCED_INSTALLER"
        print_info "Please ensure you have the complete Servin installer package"
        exit 1
    fi
    
    print_success "Enhanced installer found: $ENHANCED_INSTALLER"
}

# Show installation options
show_installation_options() {
    print_header "Installation Options"
    
    echo -e "\n${BOLD}What would you like to do?${NC}\n"
    
    if [[ -n "$EXISTING_SERVIN" ]]; then
        echo "1) Update existing installation (preserve configuration)"
        echo "2) Fresh installation (reset configuration)"
        echo "3) Install VM prerequisites only"
        echo "4) Exit"
        echo
        read -p "Choose option (1-4): " -n 1 -r
        echo
        INSTALL_OPTION=$REPLY
    else
        echo "1) Full installation with VM prerequisites"
        echo "2) Basic installation (skip VM setup)"
        echo "3) Install VM prerequisites only"
        echo "4) Exit"
        echo
        read -p "Choose option (1-4): " -n 1 -r
        echo
        INSTALL_OPTION=$REPLY
    fi
}

# Confirm installation
confirm_installation() {
    print_header "Installation Summary"
    
    echo -e "\n${BOLD}Installation Details:${NC}"
    echo "• Platform: $PLATFORM"
    echo "• Enhanced installer: $(basename "$ENHANCED_INSTALLER")"
    
    case "$INSTALL_OPTION" in
        1)
            if [[ -n "$EXISTING_SERVIN" ]]; then
                echo "• Action: Update existing installation"
            else
                echo "• Action: Full installation with VM prerequisites"
            fi
            ;;
        2)
            if [[ -n "$EXISTING_SERVIN" ]]; then
                echo "• Action: Fresh installation"
            else
                echo "• Action: Basic installation (no VM setup)"
            fi
            ;;
        3)
            echo "• Action: Install VM prerequisites only"
            ;;
    esac
    
    if [[ "$PREREQ_MISSING" == "true" ]]; then
        echo "• VM Prerequisites: Will be installed"
    else
        echo "• VM Prerequisites: Already available"
    fi
    
    if [[ $RESOURCE_WARNINGS -gt 0 ]]; then
        echo -e "• ${YELLOW}Warnings: $RESOURCE_WARNINGS resource warnings detected${NC}"
    fi
    
    echo
    read -p "Continue with installation? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Installation cancelled"
        exit 0
    fi
}

# Check permissions
check_permissions() {
    case "$PLATFORM" in
        "linux")
            if [[ $EUID -ne 0 ]] && [[ "$INSTALL_OPTION" =~ ^[13]$ ]]; then
                print_error "Root privileges required for Linux installation"
                print_info "Please run: sudo ./install-wizard.sh"
                exit 1
            fi
            ;;
        "macos")
            if [[ $EUID -eq 0 ]]; then
                print_warning "Running as root on macOS is not recommended"
                print_info "The installer will prompt for admin password when needed"
            fi
            ;;
    esac
}

# Run enhanced installer
run_enhanced_installer() {
    print_header "Running Enhanced Installer"
    
    local installer_args=""
    
    case "$INSTALL_OPTION" in
        2)
            if [[ -z "$EXISTING_SERVIN" ]]; then
                installer_args="--no-vm"
            fi
            ;;
        3)
            # For VM prerequisites only, we'll run the full installer but explain it will install Servin too
            print_info "Note: This will install VM prerequisites along with Servin"
            ;;
    esac
    
    print_info "Executing: $ENHANCED_INSTALLER $installer_args"
    print_info "This may take several minutes..."
    echo
    
    # Make installer executable
    chmod +x "$ENHANCED_INSTALLER"
    
    # Run the enhanced installer
    if [[ "$PLATFORM" == "linux" ]]; then
        "$ENHANCED_INSTALLER" $installer_args
    else
        "$ENHANCED_INSTALLER" $installer_args
    fi
    
    local exit_code=$?
    
    if [[ $exit_code -eq 0 ]]; then
        print_success "Enhanced installer completed successfully!"
    else
        print_error "Enhanced installer failed with exit code: $exit_code"
        return $exit_code
    fi
}

# Post-installation verification
post_installation_verification() {
    print_header "Post-Installation Verification"
    
    # Check if Servin is now available
    local servin_cmd=""
    for path in "/usr/local/bin/servin" "/usr/bin/servin" "$HOME/.local/bin/servin"; do
        if [[ -x "$path" ]]; then
            servin_cmd="$path"
            break
        fi
    done
    
    if [[ -z "$servin_cmd" ]]; then
        print_warning "Servin command not found in standard locations"
        print_info "You may need to restart your terminal or add Servin to your PATH"
        return 1
    fi
    
    # Test basic functionality
    print_info "Testing Servin CLI..."
    if "$servin_cmd" version >/dev/null 2>&1; then
        local version=$("$servin_cmd" version 2>/dev/null | head -1)
        print_success "Servin CLI working: $version"
    else
        print_warning "Servin CLI test failed"
    fi
    
    # Test VM functionality if installed
    if [[ "$INSTALL_OPTION" != "2" ]]; then
        print_info "Testing VM functionality..."
        if "$servin_cmd" vm status >/dev/null 2>&1; then
            print_success "VM subsystem working"
        else
            print_warning "VM subsystem not fully functional"
            print_info "This may be normal on first run - try: $servin_cmd vm init"
        fi
    fi
}

# Show completion message
show_completion() {
    print_header "Installation Complete!"
    
    echo -e "\n${GREEN}${BOLD}🎉 Servin Container Runtime has been successfully installed!${NC}\n"
    
    echo -e "${BOLD}Next Steps:${NC}"
    
    case "$PLATFORM" in
        "linux")
            echo "1. Logout and login again (for group membership changes)"
            echo "2. Initialize VM support: servin vm init"
            echo "3. Test installation: servin run --vm alpine echo 'Hello!'"
            ;;
        "macos")
            echo "1. Restart your terminal (for PATH updates)"
            echo "2. Initialize VM support: servin vm init"  
            echo "3. Test installation: servin run --vm alpine echo 'Hello!'"
            ;;
    esac
    
    echo
    echo -e "${BOLD}Available Commands:${NC}"
    echo "• servin version       - Show version information"
    echo "• servin vm status     - Check VM subsystem status"
    echo "• servin-gui          - Launch graphical interface"
    echo "• servin-tui          - Launch terminal interface"
    
    echo
    echo -e "${BOLD}Documentation:${NC}"
    echo "• Complete guide: installers/VM_PREREQUISITES.md"
    echo "• CLI reference: docs/cli.md"
    echo "• Troubleshooting: docs/troubleshooting.md"
    
    echo
    print_success "Installation wizard completed successfully!"
}

# Main execution flow
main() {
    print_banner
    
    print_info "Starting Servin Container Runtime installation wizard..."
    echo
    
    detect_platform
    check_directory
    check_existing_installation
    check_vm_prerequisites
    check_system_resources
    get_enhanced_installer
    show_installation_options
    
    case "$INSTALL_OPTION" in
        1|2|3)
            confirm_installation
            check_permissions
            run_enhanced_installer
            post_installation_verification
            show_completion
            ;;
        4)
            print_info "Installation cancelled"
            exit 0
            ;;
        *)
            print_error "Invalid option selected"
            exit 1
            ;;
    esac
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "Servin Container Runtime - Smart Installation Wizard"
        echo
        echo "Usage: $0 [options]"
        echo
        echo "Options:"
        echo "  --help, -h     Show this help message"
        echo "  --auto         Run automated installation (no prompts)"
        echo "  --vm-only      Install VM prerequisites only"
        echo "  --no-vm        Skip VM setup"
        echo
        echo "This wizard will:"
        echo "1. Detect your platform (Linux/macOS)"
        echo "2. Check for existing installations"
        echo "3. Verify system prerequisites"
        echo "4. Run the appropriate enhanced installer"
        echo "5. Verify the installation"
        echo
        echo "Enhanced installers will be used:"
        echo "• Linux: linux/install-with-vm.sh"
        echo "• macOS: macos/install-with-vm.sh"
        exit 0
        ;;
    --auto)
        AUTO_MODE=1
        INSTALL_OPTION=1
        ;;
    --vm-only)
        AUTO_MODE=1
        INSTALL_OPTION=3
        ;;
    --no-vm)
        AUTO_MODE=1
        INSTALL_OPTION=2
        ;;
esac

# Run main installation
main "$@"