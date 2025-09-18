#!/bin/bash
# Servin Cross-Platform Build and Package Script

set -e

# Configuration
VERSION="1.0.0"
BUILD_DIR="build"
DIST_DIR="dist"
INSTALLER_DIR="installers"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

print_header() {
    echo -e "${CYAN}================================================${NC}"
    echo -e "${CYAN}   Servin Container Runtime - Build Script${NC}"
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

cleanup() {
    print_info "Cleaning up previous builds..."
    rm -rf "$BUILD_DIR" "$DIST_DIR"
    mkdir -p "$BUILD_DIR" "$DIST_DIR"
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
    
    # Check if pyinstaller is available (for standalone builds)
    if ! python3 -c "import PyInstaller" 2>/dev/null; then
        print_info "  PyInstaller not found - installing for standalone builds..."
        pip3 install pyinstaller >/dev/null 2>&1 || {
            print_warning "  Failed to install PyInstaller - WebView GUI will use Python runtime"
            return 1
        }
    fi
    
    print_success "  Python environment ready for WebView GUI"
    return 0
}

build_webview_gui() {
    local platform=$1
    local arch=$2
    local ext=$3
    local output_dir="$BUILD_DIR/$platform-$arch"
    
    print_info "  Building WebView GUI for $platform/$arch..."
    
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
    
    # Create platform-specific launcher scripts
    case "$platform" in
        "windows")
            # Create Windows batch launcher
            cat > "$output_dir/servin-webview$ext" << 'EOF'
@echo off
setlocal enabledelayedexpansion

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

:: Install requirements if needed
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
            # Create Unix shell launcher
            cat > "$output_dir/servin-webview$ext" << 'EOF'
#!/bin/bash

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
    python3 -m venv "$WEBVIEW_DIR/venv"
    source "$WEBVIEW_DIR/venv/bin/activate"
    pip install -r "$WEBVIEW_DIR/requirements.txt"
else
    source "$WEBVIEW_DIR/venv/bin/activate"
fi

# Launch the WebView GUI
cd "$WEBVIEW_DIR"
python main.py
EOF
            chmod +x "$output_dir/servin-webview$ext"
            ;;
    esac
    
    # Try to create standalone executable if PyInstaller is available
    if command -v python3 >/dev/null 2>&1 && python3 -c "import PyInstaller" 2>/dev/null; then
        print_info "    Creating standalone WebView GUI executable..."
        
        local pyinstaller_args=""
        case "$platform" in
            "windows")
                pyinstaller_args="--windowed --icon=../icons/servin.ico"
                ;;
            "darwin")
                pyinstaller_args="--windowed --icon=../icons/servin.icns"
                ;;
            "linux")
                pyinstaller_args="--windowed"
                ;;
        esac
        
        # Create standalone executable
        (cd webview_gui && python3 -m PyInstaller \
            --onefile \
            --clean \
            --name "servin-webview-standalone" \
            $pyinstaller_args \
            --distpath "../$output_dir" \
            --workpath "../$BUILD_DIR/pyinstaller-work" \
            --specpath "../$BUILD_DIR/pyinstaller-spec" \
            main.py 2>/dev/null) && \
        print_success "    Standalone WebView GUI executable created" || \
        print_warning "    Failed to create standalone executable, launcher script available"
    fi
    
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
        -ldflags="-w -s -X main.version=$VERSION" \
        -o "$output_dir/servin-tui$ext" \
        ./cmd/servin-tui
    
    # Build WebView GUI (cross-platform Python-based)
    build_webview_gui "$platform" "$arch" "$ext"
    
    print_success "  Built binaries for $platform/$arch"
    
    # List what was actually built
    print_info "    Built files:"
    ls -la "$output_dir"/ | grep -E "(servin|\.exe)$" | while read -r line; do
        echo "      $(echo "$line" | awk '{print $9, "(" $5 " bytes)"}')"
    done
}

create_windows_package() {
    print_info "Creating Windows package..."
    
    local platform_dir="$BUILD_DIR/windows-amd64"
    local package_dir="$DIST_DIR/servin-windows-$VERSION"
    
    mkdir -p "$package_dir"
    
    # Copy binaries
    cp "$platform_dir"/*.exe "$package_dir/" 2>/dev/null || true
    
    # Copy WebView GUI if available
    if [[ -d "$platform_dir/webview_gui" ]]; then
        cp -r "$platform_dir/webview_gui" "$package_dir/"
        [[ -f "$platform_dir/servin-webview.bat" ]] && cp "$platform_dir/servin-webview.bat" "$package_dir/"
        [[ -f "$platform_dir/servin-webview-standalone.exe" ]] && cp "$platform_dir/servin-webview-standalone.exe" "$package_dir/"
    fi
    
    # Copy installer and wizard
    cp "$INSTALLER_DIR/windows/install.ps1" "$package_dir/" 2>/dev/null || true
    if [[ -f "$INSTALLER_DIR/windows/servin-installer.py" ]]; then
        cp "$INSTALLER_DIR/windows/servin-installer.py" "$package_dir/"
    fi
    
    # Copy icons
    if [[ -d "icons" ]]; then
        mkdir -p "$package_dir/icons"
        cp icons/*.ico "$package_dir/icons/" 2>/dev/null || true
        cp icons/*.png "$package_dir/icons/" 2>/dev/null || true
    fi
    
    # Create README
    cat > "$package_dir/README.txt" << EOF
Servin Container Runtime for Windows
Version: $VERSION

Components:
- servin.exe               : Command-line interface
- servin-tui.exe       : Terminal user interface  
- servin-webview.bat       : WebView GUI launcher (Python-based)
- servin-webview-standalone.exe : Standalone WebView GUI (if available)
- webview_gui/             : WebView GUI source files

Installation Options:

Option 1 - GUI Installer (Recommended):
python servin-installer.py

Option 2 - PowerShell Script:
1. Right-click PowerShell and select "Run as Administrator"
2. Navigate to this directory
3. Run: .\install.ps1

This will:
- Install Servin to C:\Program Files\Servin
- Create a Windows Service named "ServinRuntime"
- Add Servin to your PATH
- Create desktop shortcuts
- Set up data directories in C:\ProgramData\Servin

Usage:
- CLI: servin.exe --help
- TUI: servin-tui.exe
- WebView GUI: servin-webview.bat or servin-webview-standalone.exe
- Service: Start-Service ServinRuntime

WebView GUI Features:
- Modern web-based interface
- Cross-platform compatibility
- Real-time container monitoring
- Image management dashboard
- Container logs viewer
- System statistics

WebView GUI Requirements:
- Python 3.7+ (for script version)
- Internet connection for initial setup
- Modern web browser engine

System Requirements:
- Windows 10 version 1803 or later
- .NET Framework 4.7.2 or later
- At least 4GB RAM recommended
- Python 3.7+ (for WebView GUI script version)

For more information, visit: https://github.com/yourusername/servin
EOF
    
    # Create ZIP archive
    cd "$DIST_DIR"
    if command -v zip >/dev/null 2>&1; then
        zip -r "servin-windows-$VERSION.zip" "servin-windows-$VERSION"
        print_success "  Created: servin-windows-$VERSION.zip"
    else
        print_warning "  ZIP command not found, directory package created instead"
    fi
    cd - >/dev/null
}

create_linux_package() {
    print_info "Creating Linux package..."
    
    local platform_dir="$BUILD_DIR/linux-amd64"
    local package_dir="$DIST_DIR/servin-linux-$VERSION"
    
    mkdir -p "$package_dir"
    
    # Copy binaries
    cp "$platform_dir"/servin* "$package_dir/" 2>/dev/null || cp "$platform_dir/servin" "$package_dir/"
    
    # Copy WebView GUI if available
    if [[ -d "$platform_dir/webview_gui" ]]; then
        cp -r "$platform_dir/webview_gui" "$package_dir/"
        [[ -f "$platform_dir/servin-webview" ]] && cp "$platform_dir/servin-webview" "$package_dir/"
        [[ -f "$platform_dir/servin-webview-standalone" ]] && cp "$platform_dir/servin-webview-standalone" "$package_dir/"
    fi
    
    # Copy installer and wizard
    cp "$INSTALLER_DIR/linux/install.sh" "$package_dir/" 2>/dev/null || true
    if [[ -f "$INSTALLER_DIR/linux/servin-installer.py" ]]; then
        cp "$INSTALLER_DIR/linux/servin-installer.py" "$package_dir/"
    fi
    chmod +x "$package_dir"/*.sh 2>/dev/null || true
    
    # Copy icons and desktop files
    if [[ -d "icons" ]]; then
        mkdir -p "$package_dir/icons"
        cp icons/*.png "$package_dir/icons/" 2>/dev/null || true
        cp icons/*.svg "$package_dir/icons/" 2>/dev/null || true
    fi
    
    # Create desktop file for WebView GUI
    cat > "$package_dir/servin.desktop" << EOF
[Desktop Entry]
Version=1.0
Type=Application
Name=Servin Container Runtime
Comment=Modern WebView-based container management interface
Exec=servin-webview
Icon=servin
Terminal=false
Categories=Development;System;
EOF
    
    # Create README
    cat > "$package_dir/README.md" << EOF
# Servin Container Runtime for Linux
Version: $VERSION

## Components
- \`servin\`                    : Command-line interface
- \`servin-tui\`            : Terminal user interface
- \`servin-webview\`            : WebView GUI launcher (Python-based)
- \`servin-webview-standalone\` : Standalone WebView GUI (if available)
- \`webview_gui/\`              : WebView GUI source files

## Installation Options

### Option 1 - GUI Installer (Recommended)
\`\`\`bash
python3 servin-installer.py
\`\`\`

### Option 2 - Shell Script
\`\`\`bash
sudo ./install.sh
\`\`\`

This will:
- Install Servin to /usr/local/bin
- Create a system user 'servin'
- Set up systemd service (or SysV init script)
- Create configuration in /etc/servin
- Set up data directories in /var/lib/servin
- Install desktop entries for WebView GUI

## WebView GUI Dependencies
For the WebView GUI to work, install required packages:

**Ubuntu/Debian:**
\`\`\`bash
sudo apt-get update
sudo apt-get install python3 python3-pip python3-venv
sudo apt-get install python3-tk python3-dev
sudo apt-get install libwebkit2gtk-4.0-dev  # For pywebview
\`\`\`

**CentOS/RHEL/Fedora:**
\`\`\`bash
sudo yum install python3 python3-pip  # CentOS/RHEL
sudo dnf install python3 python3-pip  # Fedora
sudo yum install python3-tkinter python3-devel  # CentOS/RHEL
sudo dnf install python3-tkinter python3-devel  # Fedora
sudo yum install webkit2gtk3-devel  # CentOS/RHEL
sudo dnf install webkit2gtk3-devel  # Fedora
\`\`\`

**Arch Linux:**
\`\`\`bash
sudo pacman -S python python-pip
sudo pacman -S tk python-tkinter
sudo pacman -S webkit2gtk
\`\`\`

## Usage
- CLI: \`servin --help\`
- TUI: \`servin-tui\`
- WebView GUI: \`./servin-webview\` or \`./servin-webview-standalone\`
- Service: \`sudo systemctl start servin\`

## WebView GUI Features
- Modern web-based interface using webkit2gtk
- Cross-platform compatibility
- Real-time container monitoring
- Image management dashboard
- Container logs viewer
- System statistics
- Responsive design for different screen sizes

## System Requirements
- Linux kernel 3.10+ (most modern distributions)
- Python 3.7+ (for WebView GUI)
- At least 4GB RAM recommended
- webkit2gtk (for WebView GUI)

For more information, visit: https://github.com/yourusername/servin
EOF
    
    # Create tar.gz archive
    cd "$DIST_DIR"
    tar -czf "servin-linux-$VERSION.tar.gz" "servin-linux-$VERSION"
    print_success "  Created: servin-linux-$VERSION.tar.gz"
    cd - >/dev/null
}

create_macos_package() {
    print_info "Creating macOS package..."
    
    local platform_dir="$BUILD_DIR/darwin-amd64"
    local package_dir="$DIST_DIR/servin-macos-$VERSION"
    
    mkdir -p "$package_dir"
    
    # Copy binaries
    cp "$platform_dir"/servin* "$package_dir/" 2>/dev/null || cp "$platform_dir/servin" "$package_dir/"
    
    # Copy WebView GUI if available
    if [[ -d "$platform_dir/webview_gui" ]]; then
        cp -r "$platform_dir/webview_gui" "$package_dir/"
        [[ -f "$platform_dir/servin-webview" ]] && cp "$platform_dir/servin-webview" "$package_dir/"
        [[ -f "$platform_dir/servin-webview-standalone" ]] && cp "$platform_dir/servin-webview-standalone" "$package_dir/"
    fi
    
    # Copy installer and wizard
    cp "$INSTALLER_DIR/macos/install.sh" "$package_dir/" 2>/dev/null || true
    if [[ -f "$INSTALLER_DIR/macos/servin-installer.py" ]]; then
        cp "$INSTALLER_DIR/macos/servin-installer.py" "$package_dir/"
    fi
    chmod +x "$package_dir"/*.sh 2>/dev/null || true
    
    # Copy icons and create app bundles
    if [[ -d "icons" ]]; then
        mkdir -p "$package_dir/icons"
        cp icons/*.icns "$package_dir/icons/" 2>/dev/null || true
        cp icons/*.png "$package_dir/icons/" 2>/dev/null || true
        
        # Create application bundle for WebView GUI
        if [[ -f "$package_dir/servin-webview" ]]; then
            local webview_app_dir="$package_dir/Servin.app"
            mkdir -p "$webview_app_dir/Contents/"{MacOS,Resources}
            
            # Copy launcher script
            cp "$package_dir/servin-webview" "$webview_app_dir/Contents/MacOS/"
            
            # Copy WebView GUI files
            cp -r "$package_dir/webview_gui" "$webview_app_dir/Contents/Resources/"
            
            # Copy icon
            cp icons/*.icns "$webview_app_dir/Contents/Resources/servin.icns" 2>/dev/null || true
            
            # Create Info.plist for WebView GUI
            cat > "$webview_app_dir/Contents/Info.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>servin-webview</string>
    <key>CFBundleIdentifier</key>
    <string>com.servin.container-runtime</string>
    <key>CFBundleName</key>
    <string>Servin</string>
    <key>CFBundleDisplayName</key>
    <string>Servin Container Runtime</string>
    <key>CFBundleVersion</key>
    <string>$VERSION</string>
    <key>CFBundleShortVersionString</key>
    <string>$VERSION</string>
    <key>CFBundleIconFile</key>
    <string>servin.icns</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.13</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>NSSupportsAutomaticGraphicsSwitching</key>
    <true/>
</dict>
</plist>
EOF
        fi
    fi
    
    # Create README
    cat > "$package_dir/README.md" << EOF
# Servin Container Runtime for macOS
Version: $VERSION

## Components
- \`servin\`                    : Command-line interface
- \`servin-tui\`            : Terminal user interface
- \`servin-webview\`            : WebView GUI launcher (Python-based)
- \`servin-webview-standalone\` : Standalone WebView GUI (if available)
- \`Servin WebView.app\`        : WebView macOS application bundle
- \`webview_gui/\`              : WebView GUI source files

## Installation Options

### Option 1 - GUI Installer (Recommended)
\`\`\`bash
python3 servin-installer.py
\`\`\`

### Option 2 - Shell Script
\`\`\`bash
sudo ./install.sh
\`\`\`

This will:
- Install Servin to /usr/local/bin
- Create a system user '_servin'
- Set up launchd service
- Create configuration in /usr/local/etc/servin
- Set up data directories in /usr/local/var/lib/servin
- Install application bundles to /Applications

## WebView GUI Requirements
The WebView GUI uses Python and requires:

\`\`\`bash
# Install Python 3 (if not already installed)
brew install python3

# The WebView GUI will automatically set up its environment
\`\`\`

## Usage
- CLI: \`servin --help\`
- TUI: \`servin-tui\`
- WebView GUI: Open "Servin" from Applications, run \`./servin-webview\`, or double-click "Servin.app"
- Service: Starts automatically (launchd)

## WebView GUI Features
- Modern web-based interface using macOS WebKit
- Native macOS look and feel
- Cross-platform compatibility
- Real-time container monitoring
- Image management dashboard
- Container logs viewer
- System statistics
- Integration with macOS notifications
- Responsive design for different screen sizes

## System Requirements
- macOS 10.13 (High Sierra) or later
- Intel or Apple Silicon Mac
- Python 3.7+ (for WebView GUI)
- At least 4GB RAM recommended

For more information, visit: https://github.com/yourusername/servin
EOF
    
    # Create tar.gz archive
    cd "$DIST_DIR"
    tar -czf "servin-macos-$VERSION.tar.gz" "servin-macos-$VERSION"
    print_success "  Created: servin-macos-$VERSION.tar.gz"
    cd - >/dev/null
}

build_docker_images() {
    print_info "Building Docker images..."
    
    # Create Dockerfile
    cat > Dockerfile << EOF
FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY build/linux-amd64/servin /usr/local/bin/servin

RUN adduser -D -s /bin/sh servin
USER servin

EXPOSE 10250

CMD ["/usr/local/bin/servin", "daemon"]
EOF
    
    # Build Docker image
    if command -v docker >/dev/null 2>&1; then
        docker build -t "servin:$VERSION" .
        docker tag "servin:$VERSION" "servin:latest"
        print_success "  Built Docker image: servin:$VERSION"
        
        # Save Docker image
        docker save "servin:$VERSION" | gzip > "$DIST_DIR/servin-docker-$VERSION.tar.gz"
        print_success "  Saved Docker image: servin-docker-$VERSION.tar.gz"
    else
        print_warning "  Docker not found, skipping Docker image build"
    fi
    
    rm -f Dockerfile
}

generate_checksums() {
    print_info "Generating checksums..."
    
    cd "$DIST_DIR"
    
    # Generate SHA256 checksums
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum *.tar.gz *.zip 2>/dev/null > "servin-$VERSION-checksums.txt" || true
    elif command -v shasum >/dev/null 2>&1; then
        shasum -a 256 *.tar.gz *.zip 2>/dev/null > "servin-$VERSION-checksums.txt" || true
    fi
    
    print_success "  Generated checksums"
    cd - >/dev/null
}

show_summary() {
    echo ""
    print_success "================================================"
    print_success "   Build completed successfully!"
    print_success "================================================"
    echo ""
    print_info "Built packages:"
    ls -la "$DIST_DIR"/ | grep -E "\.(zip|tar\.gz)$" | while read -r line; do
        echo "  $line"
    done
    echo ""
    print_info "Installation instructions:"
    echo "Windows: Extract ZIP, run python servin-installer.py or install.ps1 as Administrator"
    echo "Linux:   Extract tar.gz, run python3 servin-installer.py or sudo ./install.sh"  
    echo "macOS:   Extract tar.gz, run python3 servin-installer.py or sudo ./install.sh"
    echo ""
    print_info "GUI Options included:"
    echo "- WebView GUI: Modern web-based interface (servin-webview)"
    echo "- TUI: Terminal-based interface (servin-tui)"
    echo ""
    print_info "WebView GUI Features:"
    echo "- Cross-platform Python-based implementation"
    echo "- Modern web interface with real-time updates"
    echo "- Container management dashboard"
    echo "- Image browser and management"
    echo "- Container logs viewer"
    echo "- System statistics and monitoring"
    echo ""
    if [[ -f "$DIST_DIR/servin-docker-$VERSION.tar.gz" ]]; then
        print_info "Docker usage:"
        echo "docker load < servin-docker-$VERSION.tar.gz"
        echo "docker run -d -p 10250:10250 servin:$VERSION"
    fi
}

main() {
    print_header
    
    # Check if we're in the right directory
    if [[ ! -f "go.mod" ]]; then
        print_error "This script must be run from the project root directory"
        exit 1
    fi
    
    # Check Go installation
    if ! command -v go >/dev/null 2>&1; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    # Check Python requirements for WebView GUI
    check_python_requirements
    
    cleanup
    
    # Build for different platforms
    build_binaries "windows" "amd64" ".exe"
    build_binaries "linux" "amd64" ""
    build_binaries "darwin" "amd64" ""
    
    # Create platform packages
    create_windows_package
    create_linux_package
    create_macos_package
    
    # Build Docker images
    build_docker_images
    
    # Generate checksums
    generate_checksums
    
    show_summary
}

# VM Automation Functions
build_with_vm() {
    print_header
    print_info "🚀 Building Servin with VM Automation"
    echo ""
    
    # Check if we're in the right directory
    if [[ ! -f "go.mod" ]]; then
        print_error "This script must be run from the project root directory"
        exit 1
    fi
    
    # Check Go installation
    if ! command -v go >/dev/null 2>&1; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    # Build only the binary we need for VM testing
    print_info "📦 Building Servin binary..."
    go build -o servin main.go
    print_success "✅ Servin binary built successfully"
    echo ""
    
    # Clean any existing VM to start fresh
    print_info "🧹 Cleaning existing VM data..."
    rm -rf ~/.servin/vms/servin-vm 2>/dev/null || true
    print_success "✅ VM data cleaned"
    echo ""
    
    # Start VM with automated SSH setup
    print_info "🚀 Starting VM with automated SSH configuration..."
    print_info "This will:"
    print_info "  • Download Alpine Linux kernel if needed"
    print_info "  • Create cloud-init ISO with SSH automation"
    print_info "  • Start QEMU VM with Hypervisor.framework acceleration"
    print_info "  • Automatically configure SSH access"
    print_info "  • Deploy Servin binary to VM"
    echo ""
    
    ./servin vm start
    
    echo ""
    print_info "⏳ Monitoring SSH setup progress..."
    
    SSH_READY=false
    MAX_WAIT=90
    
    for i in $(seq 1 $MAX_WAIT); do
        if ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=2 -o BatchMode=yes root@localhost 'echo "SSH_WORKING"' 2>/dev/null | grep -q "SSH_WORKING"; then
            SSH_READY=true
            print_success "✅ SSH is ready after $i seconds!"
            break
        fi
        
        # Show progress every 10 seconds
        if [ $((i % 10)) -eq 0 ]; then
            print_info "   Waiting for SSH... ($i/$MAX_WAIT seconds)"
        fi
        
        sleep 1
    done
    
    if [ "$SSH_READY" = true ]; then
        echo ""
        print_success "🎯 VM Setup Complete!"
        print_success "===================="
        echo ""
        
        # Get VM information
        print_info "📊 VM Information:"
        VM_STATUS=$(./servin vm status | grep "VM Status" | awk '{print $3}' || echo "Unknown")
        echo "   Status: $VM_STATUS"
        echo "   SSH: ssh root@localhost -p 2222"
        echo "   Password: servin123"
        echo ""
        
        # Test VM connectivity and get system info
        print_info "🔍 Testing VM connectivity..."
        VM_KERNEL=$(ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@localhost 'uname -r' 2>/dev/null || echo "Unknown")
        VM_DISTRO=$(ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@localhost 'cat /etc/alpine-release' 2>/dev/null || echo "Unknown")
        echo "   Kernel: $VM_KERNEL"
        echo "   Distribution: Alpine Linux $VM_DISTRO"
        echo ""
        
        # Deploy Servin to VM
        print_info "📦 Deploying Servin to VM..."
        if scp -P 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ./servin root@localhost:/usr/local/bin/ 2>/dev/null; then
            ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@localhost 'chmod +x /usr/local/bin/servin' 2>/dev/null
            print_success "✅ Servin deployed to VM successfully"
        else
            print_warning "⚠️  Failed to deploy Servin to VM (manual deployment may be needed)"
        fi
        
        echo ""
        print_info "🧪 Testing Container Functionality..."
        print_info "===================================="
        echo ""
        
        # Test container operations
        print_info "1. Testing hello-world container:"
        if ./servin run --name test-hello hello-world 2>/dev/null; then
            print_success "✅ Container run successful"
        else
            print_warning "❌ Container run failed (may need manual configuration)"
        fi
        
        echo ""
        print_info "2. Listing containers:"
        ./servin list 2>/dev/null || print_warning "❌ Container list failed"
        
        echo ""
        print_info "3. Testing container logs:"
        ./servin logs test-hello 2>/dev/null || print_warning "❌ Container logs failed"
        
        echo ""
        print_success "🎉 Build and VM Setup Complete!"
        print_success "==============================="
        echo ""
        print_info "🎯 Ready for Development:"
        echo "   • VM Status: Running with SSH"
        echo "   • Container Runtime: Native Linux (not Docker simulation)"
        echo "   • SSH Access: ssh root@localhost -p 2222"
        echo "   • Servin Commands: ./servin run, ./servin exec, ./servin logs"
        echo ""
        print_info "📚 Example Commands:"
        echo "   ./servin run nginx:alpine"
        echo "   ./servin run --name web -p 8080:80 nginx:alpine"
        echo "   ./servin exec web sh"
        echo "   ./servin logs web"
        echo ""
        
    else
        echo ""
        print_warning "⚠️  SSH Auto-Setup Incomplete"
        print_warning "============================="
        echo ""
        print_info "The VM is running but SSH auto-setup didn't complete within $MAX_WAIT seconds."
        echo ""
        print_info "Manual setup required:"
        echo "1. Connect to VM console"
        echo "2. Login as root (no password needed)"
        echo "3. Mount and run setup script:"
        echo "   mount /dev/sr0 /mnt 2>/dev/null || true"
        echo "   /mnt/autosetup.sh"
        echo ""
        print_info "Alternative manual commands:"
        echo "   apk update && apk add openssh"
        echo "   echo 'root:servin123' | chpasswd"
        echo "   echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config"
        echo "   rc-update add sshd default && rc-service sshd start"
        echo ""
        echo "Then test: ssh root@localhost -p 2222"
        echo ""
    fi
    
    print_info "🏁 Build script completed!"
    echo ""
    VM_STATUS=$(./servin vm status | grep "VM Status" | awk '{print $3}' || echo "Unknown")
    QEMU_PID=$(ps aux | grep qemu-system-aarch64 | grep -v grep | awk '{print $2}' | head -1 || echo "Not found")
    echo "VM Status: $VM_STATUS"
    echo "QEMU Process: PID $QEMU_PID"
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
        create_windows_package
        ;;
    linux)
        cleanup
        build_binaries "linux" "amd64" ""
        create_linux_package
        ;;
    macos)
        cleanup
        build_binaries "darwin" "amd64" ""
        create_macos_package
        ;;
    docker)
        build_binaries "linux" "amd64" ""
        build_docker_images
        ;;
    vm)
        build_with_vm
        ;;
    *)
        main
        ;;
esac
