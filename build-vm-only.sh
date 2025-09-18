#!/bin/bash

# Servin Build Script - VM Mode Only (No Docker)
# This script builds Servin for VM-based containerization without Docker dependencies

set -e

# Configuration
VERSION="1.0.0"
BUILD_DIR="build"
DIST_DIR="dist"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

print_success() { echo -e "${GREEN}✓ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠ $1${NC}"; }
print_info() { echo -e "${BLUE}→ $1${NC}"; }
print_header() { echo -e "\n${CYAN}${1}${NC}"; }

print_banner() {
    echo -e "${CYAN}"
    echo "================================================"
    echo "   Servin Container Runtime - VM Mode Build"
    echo "         (No Docker Dependencies)"
    echo "================================================"
    echo -e "${NC}\n"
}

check_python_requirements() {
    print_info "Checking Python requirements for WebView GUI..."
    
    # Check if Python 3 is available
    if ! command -v python3 >/dev/null 2>&1; then
        print_warning "  Python 3 not found - WebView GUI packaging will be skipped"
        return 1
    fi
    
    # Check if pip is available
    if ! command -v pip3 >/dev/null 2>&1; then
        print_warning "  pip3 not found - WebView GUI packaging will be skipped"
        return 1
    fi
    
    print_success "  Python environment ready for WebView GUI"
    return 0
}

build_webview_gui_enhanced() {
    local platform=$1
    local arch=$2
    local ext=$3
    local output_dir="$BUILD_DIR/$platform-$arch"
    
    print_info "  Building Enhanced WebView GUI for $platform/$arch..."
    
    # Skip if webview_gui directory doesn't exist
    if [[ ! -d "webview_gui" ]]; then
        print_warning "    WebView GUI source not found, skipping..."
        return
    fi
    
    # Create webview GUI output directory
    local webview_dir="$output_dir/webview_gui"
    mkdir -p "$webview_dir"
    
    # Copy WebView GUI source files
    cp -r webview_gui/* "$webview_dir/"
    
    # Try to create standalone executable with detailed error reporting
    if command -v python3 >/dev/null 2>&1 && python3 -c "import PyInstaller" 2>/dev/null; then
        print_info "    Creating standalone WebView GUI executable..."
        
        local pyinstaller_args=""
        case "$platform" in
            "windows")
                pyinstaller_args="--windowed --icon=../icons/servin.ico"
                ;;
            "darwin")
                # Check if icon exists before using it
                if [[ -f "../icons/servin.icns" ]]; then
                    pyinstaller_args="--windowed --icon=../icons/servin.icns"
                else
                    pyinstaller_args="--windowed"
                    print_warning "    macOS icon not found, building without icon"
                fi
                ;;
            "linux")
                pyinstaller_args="--windowed"
                ;;
        esac
        
        # Create standalone executable with error capture
        local pyinstaller_log="$BUILD_DIR/pyinstaller_${platform}_${arch}.log"
        print_info "    PyInstaller log: $pyinstaller_log"
        print_info "    Working directory: $(pwd)"
        print_info "    Checking icon path from webview_gui/:"
        if [[ -f "webview_gui/../icons/servin.icns" ]]; then
            print_info "    ✅ Icon found at webview_gui/../icons/servin.icns"
        else
            print_warning "    ❌ Icon NOT found at webview_gui/../icons/servin.icns"
        fi
        
        if (cd webview_gui && python3 -m PyInstaller \
            --onefile \
            --clean \
            --name "servin-webview-standalone" \
            $pyinstaller_args \
            --distpath "../$output_dir" \
            --workpath "../$BUILD_DIR/pyinstaller-work" \
            --specpath "../$BUILD_DIR/pyinstaller-spec" \
            main.py > "../$pyinstaller_log" 2>&1); then
            print_success "    ✅ Standalone WebView GUI executable created successfully"
        else
            print_warning "    ❌ PyInstaller failed - checking log for details..."
            echo ""
            print_warning "    PyInstaller Error Details:"
            tail -20 "$pyinstaller_log" | while IFS= read -r line; do
                echo "      $line"
            done
            echo ""
            print_info "    📝 Full log available at: $pyinstaller_log"
            print_info "    🔄 Fallback: Using Python launcher script instead"
        fi
    else
        print_warning "    PyInstaller not available - using Python launcher script"
    fi
    
    # Create platform-specific launcher scripts (always created as fallback)
    case "$platform" in
        "windows")
            cat > "$output_dir/servin-webview$ext" << 'EOF'
@echo off
setlocal enabledelayedexpansion

echo Starting Servin WebView GUI...
echo.

:: Get the directory where this script is located
set "SCRIPT_DIR=%~dp0"
set "WEBVIEW_DIR=%SCRIPT_DIR%webview_gui"

:: Check if Python is available
python --version >nul 2>&1
if !errorlevel! neq 0 (
    echo Error: Python is not installed or not in PATH
    echo Please install Python 3.7+ from https://python.org
    pause
    exit /b 1
)

:: Setup virtual environment if needed
if not exist "%WEBVIEW_DIR%\venv" (
    echo Setting up WebView GUI environment...
    python -m venv "%WEBVIEW_DIR%\venv"
    call "%WEBVIEW_DIR%\venv\Scripts\activate.bat"
    pip install -r "%WEBVIEW_DIR%\requirements.txt"
) else (
    call "%WEBVIEW_DIR%\venv\Scripts\activate.bat"
)

:: Launch the WebView GUI
cd /d "%WEBVIEW_DIR%"
python main.py
EOF
            ;;
        "darwin"|"linux")
            cat > "$output_dir/servin-webview$ext" << 'EOF'
#!/bin/bash

echo "Starting Servin WebView GUI..."
echo ""

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WEBVIEW_DIR="$SCRIPT_DIR/webview_gui"

# Check if Python 3 is available
if ! command -v python3 >/dev/null 2>&1; then
    echo "Error: Python 3 is not installed or not in PATH"
    echo "Please install Python 3.7+ from your package manager"
    exit 1
fi

# Setup virtual environment if it doesn't exist
if [[ ! -d "$WEBVIEW_DIR/venv" ]]; then
    echo "Setting up WebView GUI environment..."
    cd "$WEBVIEW_DIR"
    python3 -m venv venv
    source venv/bin/activate
    pip install -r requirements.txt
else
    cd "$WEBVIEW_DIR"
    source venv/bin/activate
fi

# Launch the WebView GUI
echo "Launching WebView GUI..."
python main.py
EOF
            chmod +x "$output_dir/servin-webview$ext"
            ;;
    esac
    
    print_success "    WebView GUI built for $platform/$arch"
}

build_binaries() {
    local platform=$1
    local arch=$2
    local ext=$3
    
    print_info "Building for $platform/$arch..."
    
    local output_dir="$BUILD_DIR/$platform-$arch"
    mkdir -p "$output_dir"
    
    # Build main servin binary
    print_info "  Building CLI binary..."
    GOOS=$platform GOARCH=$arch CGO_ENABLED=0 go build \
        -ldflags="-w -s -X main.version=$VERSION" \
        -o "$output_dir/servin$ext" .
    
    # Build TUI binary (servin-tui)
    print_info "  Building TUI binary..."
    GOOS=$platform GOARCH=$arch CGO_ENABLED=0 go build \
        -ldflags="-w -s" \
        -o "$output_dir/servin-tui$ext" ./cmd/servin-tui/
    
    # Build WebView GUI with enhanced error reporting
    build_webview_gui_enhanced "$platform" "$arch" "$ext"
    
    print_success "  Built binaries for $platform/$arch"
    
    # List built files with sizes
    print_info "    Built files:"
    find "$output_dir" -maxdepth 1 -type f -executable -o -name "*.exe" | while read -r file; do
        local size=$(du -h "$file" | cut -f1)
        local filename=$(basename "$file")
        echo "      $filename ($size)"
    done
}

cleanup() {
    print_info "Cleaning up previous builds..."
    rm -rf "$BUILD_DIR" "$DIST_DIR"
    mkdir -p "$BUILD_DIR" "$DIST_DIR"
}

create_packages() {
    print_header "Creating Platform Packages (VM Mode)"
    
    # Create packages for each platform
    for platform_dir in "$BUILD_DIR"/*; do
        if [[ -d "$platform_dir" ]]; then
            local platform=$(basename "$platform_dir")
            print_info "Creating package for $platform..."
            
            case "$platform" in
                *windows*)
                    (cd "$BUILD_DIR" && zip -r "../$DIST_DIR/servin-vm-${platform}-${VERSION}.zip" "$platform/")
                    print_success "  Created: servin-vm-${platform}-${VERSION}.zip"
                    ;;
                *)
                    (cd "$BUILD_DIR" && tar -czf "../$DIST_DIR/servin-vm-${platform}-${VERSION}.tar.gz" "$platform/")
                    print_success "  Created: servin-vm-${platform}-${VERSION}.tar.gz"
                    ;;
            esac
        fi
    done
}

generate_checksums() {
    print_info "Generating checksums..."
    
    cd "$DIST_DIR"
    
    # Generate SHA256 checksums
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum *.tar.gz *.zip 2>/dev/null > "servin-vm-$VERSION-checksums.txt" || true
    elif command -v shasum >/dev/null 2>&1; then
        shasum -a 256 *.tar.gz *.zip 2>/dev/null > "servin-vm-$VERSION-checksums.txt" || true
    fi
    
    print_success "  Generated checksums"
    cd - >/dev/null
}

show_summary() {
    print_header "Build Summary"
    
    print_success "VM Mode Build completed successfully!"
    echo ""
    print_info "Built packages (VM Mode - No Docker):"
    ls -la "$DIST_DIR"/ | grep -E "\.(zip|tar\.gz)$" | while read -r line; do
        echo "  $line"
    done
    echo ""
    print_info "Installation instructions:"
    echo "Windows: Extract ZIP and run servin.exe"
    echo "Linux:   Extract tar.gz and run ./servin"  
    echo "macOS:   Extract tar.gz and run ./servin"
    echo ""
    print_info "WebView GUI Features:"
    echo "- Native VM-based containerization (no Docker required)"
    echo "- Cross-platform Python-based interface"
    echo "- Real-time container monitoring and management"
    echo "- Integrated terminal and file explorer"
    echo ""
    print_info "Usage Examples:"
    echo "./servin vm start                    # Start VM"
    echo "./servin run nginx:alpine            # Run container in VM"
    echo "./servin-webview                     # Launch GUI"
}

# Main execution function
main() {
    print_banner
    
    # Check if we're in the right directory
    if [[ ! -f "go.mod" ]]; then
        echo "Error: This script must be run from the project root directory"
        exit 1
    fi
    
    # Check Go installation
    if ! command -v go >/dev/null 2>&1; then
        echo "Error: Go is not installed or not in PATH"
        exit 1
    fi
    
    # Check Python requirements for WebView GUI
    check_python_requirements
    
    cleanup
    
    # Build for different platforms (VM mode only)
    build_binaries "windows" "amd64" ".exe"
    build_binaries "linux" "amd64" ""
    build_binaries "darwin" "amd64" ""
    
    # Create platform packages
    create_packages
    
    # Generate checksums
    generate_checksums
    
    show_summary
    
    print_header "Next Steps"
    print_info "1. Test VM functionality: ./servin vm start"
    print_info "2. Test container operations: ./servin run hello-world"
    print_info "3. Launch WebView GUI: ./servin-webview"
    print_info "4. Deploy packages from dist/ directory"
}

# Handle command line arguments
case "${1:-}" in
    clean)
        cleanup
        print_success "Cleaned build and dist directories"
        ;;
    windows)
        cleanup
        build_binaries "windows" "amd64" ".exe"
        create_packages
        ;;
    linux)
        cleanup
        build_binaries "linux" "amd64" ""
        create_packages
        ;;
    macos)
        cleanup
        build_binaries "darwin" "amd64" ""
        create_packages
        ;;
    test-pyinstaller)
        print_banner
        print_header "Testing PyInstaller Installation and Functionality"
        
        # Test PyInstaller
        if command -v python3 >/dev/null 2>&1; then
            if python3 -c "import PyInstaller" 2>/dev/null; then
                print_success "PyInstaller is installed and importable"
                python3 -c "import PyInstaller; print(f'PyInstaller version: {PyInstaller.__version__}')"
            else
                print_warning "PyInstaller not found - attempting installation..."
                pip3 install pyinstaller
                if python3 -c "import PyInstaller" 2>/dev/null; then
                    print_success "PyInstaller installed successfully"
                else
                    print_warning "PyInstaller installation failed"
                fi
            fi
        else
            print_warning "Python 3 not found"
        fi
        
        # Test webview_gui directory
        if [[ -d "webview_gui" ]]; then
            print_success "WebView GUI source directory found"
            if [[ -f "webview_gui/main.py" ]]; then
                print_success "main.py found"
            else
                print_warning "main.py not found in webview_gui/"
            fi
            if [[ -f "webview_gui/requirements.txt" ]]; then
                print_success "requirements.txt found"
                print_info "Dependencies:"
                cat webview_gui/requirements.txt | while IFS= read -r line; do
                    echo "  $line"
                done
            fi
        else
            print_warning "webview_gui directory not found"
        fi
        ;;
    *)
        main
        ;;
esac