#!/bin/bash

# Servin Container Runtime - Universal Installation Script
# Detects platform and runs appropriate wizard installer
# Usage: curl -sSL https://install.servin.dev | bash

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
    echo "║                   Universal Installer Script                  ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}\n"
}

# Detect platform
detect_platform() {
    case "$(uname -s)" in
        Linux*)
            PLATFORM="linux"
            if [[ -f /etc/os-release ]]; then
                . /etc/os-release
                DISTRO_ID=$ID
                print_info "Detected: Linux ($DISTRO_ID)"
            else
                print_info "Detected: Linux (unknown distribution)"
            fi
            ;;
        Darwin*)
            PLATFORM="macos"
            MACOS_VERSION=$(sw_vers -productVersion)
            ARCH=$(uname -m)
            print_info "Detected: macOS $MACOS_VERSION ($ARCH)"
            ;;
        CYGWIN*|MINGW*|MSYS*)
            print_error "Windows detected - please use install-wizard.ps1 instead"
            print_info "Download from: https://github.com/immyemperor/servin/releases"
            exit 1
            ;;
        *)
            print_error "Unsupported platform: $(uname -s)"
            exit 1
            ;;
    esac
}

# Check prerequisites
check_prerequisites() {
    print_header "Checking Prerequisites"
    
    # Check internet connectivity
    print_info "Checking internet connectivity..."
    if ! curl -s --connect-timeout 5 https://github.com >/dev/null; then
        print_error "Internet connection required for installation"
        exit 1
    fi
    print_success "Internet connection available"
    
    # Check required tools
    local missing_tools=()
    
    for tool in curl tar; do
        if ! command -v "$tool" >/dev/null 2>&1; then
            missing_tools+=("$tool")
        fi
    done
    
    if [[ ${#missing_tools[@]} -gt 0 ]]; then
        print_error "Missing required tools: ${missing_tools[*]}"
        print_info "Please install these tools and try again"
        exit 1
    fi
    
    print_success "Required tools available"
    
    # Check permissions for Linux
    if [[ "$PLATFORM" == "linux" ]]; then
        if [[ $EUID -ne 0 ]]; then
            print_warning "Root privileges will be required for installation"
            print_info "The installer will prompt for sudo password when needed"
        fi
    fi
}

# Download and extract installer package
download_installer() {
    print_header "Downloading Servin Installer Package"
    
    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    
    # Get latest release URL
    print_info "Fetching latest release information..."
    local release_url="https://api.github.com/repos/immyemperor/servin/releases/latest"
    local download_url
    
    case "$PLATFORM" in
        "linux")
            download_url=$(curl -s "$release_url" | grep "browser_download_url.*linux.*tar.gz" | cut -d '"' -f 4 | head -1)
            ;;
        "macos")
            download_url=$(curl -s "$release_url" | grep "browser_download_url.*macos.*tar.gz" | cut -d '"' -f 4 | head -1)
            ;;
    esac
    
    if [[ -z "$download_url" ]]; then
        print_error "Could not find installer package for $PLATFORM"
        print_info "Please download manually from: https://github.com/immyemperor/servin/releases"
        exit 1
    fi
    
    print_info "Downloading: $download_url"
    if ! curl -L -o installer.tar.gz "$download_url"; then
        print_error "Failed to download installer package"
        exit 1
    fi
    
    print_success "Download completed"
    
    # Extract installer
    print_info "Extracting installer package..."
    if ! tar -xzf installer.tar.gz; then
        print_error "Failed to extract installer package"
        exit 1
    fi
    
    # Find extracted directory
    local extracted_dir=$(find . -maxdepth 1 -type d -name "servin-*" | head -1)
    if [[ -z "$extracted_dir" ]]; then
        print_error "Could not find extracted installer directory"
        exit 1
    fi
    
    INSTALLER_DIR="$extracted_dir"
    print_success "Installer extracted to: $INSTALLER_DIR"
}

# Run platform-specific wizard installer
run_wizard_installer() {
    print_header "Running Platform-Specific Installer"
    
    cd "$INSTALLER_DIR"
    
    case "$PLATFORM" in
        "linux")
            local wizard_script="install-wizard.sh"
            ;;
        "macos")
            local wizard_script="install-wizard.sh"
            ;;
    esac
    
    if [[ ! -f "$wizard_script" ]]; then
        print_error "Wizard installer not found: $wizard_script"
        exit 1
    fi
    
    print_info "Executing wizard installer..."
    chmod +x "$wizard_script"
    
    # Run with auto mode if possible
    if [[ "$PLATFORM" == "linux" ]]; then
        if [[ $EUID -eq 0 ]]; then
            ./"$wizard_script" --auto
        else
            print_info "Switching to interactive mode (sudo required)"
            sudo ./"$wizard_script"
        fi
    else
        ./"$wizard_script" --auto
    fi
}

# Cleanup temporary files
cleanup() {
    if [[ -n "$TEMP_DIR" && -d "$TEMP_DIR" ]]; then
        print_info "Cleaning up temporary files..."
        rm -rf "$TEMP_DIR"
    fi
}

# Main installation flow
main() {
    print_banner
    
    print_info "This script will automatically install Servin Container Runtime"
    print_info "with VM containerization support for your platform."
    echo
    
    detect_platform
    check_prerequisites
    download_installer
    run_wizard_installer
    
    print_success "Universal installer completed successfully!"
    
    cleanup
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "Servin Container Runtime - Universal Installer"
        echo
        echo "Usage: $0 [options]"
        echo "   or: curl -sSL https://install.servin.dev | bash"
        echo
        echo "Options:"
        echo "  --help, -h     Show this help message"
        echo
        echo "This script will:"
        echo "1. Detect your platform (Linux/macOS)"
        echo "2. Download the appropriate installer package"
        echo "3. Run the smart wizard installer"
        echo "4. Install VM prerequisites automatically"
        echo
        echo "Supported platforms:"
        echo "• Linux: Ubuntu, Debian, Fedora, CentOS, Arch"
        echo "• macOS: 10.15+ (Intel and Apple Silicon)"
        echo
        echo "For Windows: Download and run install-wizard.ps1"
        echo "From: https://github.com/immyemperor/servin/releases"
        exit 0
        ;;
esac

# Set trap for cleanup
trap cleanup EXIT

# Run main installation
main "$@"