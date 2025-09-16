#!/bin/bash
# Build all Servin components for distribution

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check for help or clean options
if [[ "$1" == "--help" || "$1" == "-h" ]]; then
    echo "🐳 Servin Build System"
    echo ""
    echo "Usage:"
    echo "  $0                    # Auto-detect platform and build"
    echo "  PLATFORM=mac $0      # Build for specific platform"
    echo "  PLATFORM=linux $0    # Build for Linux"
    echo "  PLATFORM=windows $0  # Build for Windows"
    echo "  $0 --clean-all       # Clean all build artifacts"
    echo ""
    echo "Supported platforms: mac (Universal), linux, windows"
    echo "Build artifacts: build/[platform]/"
    echo "Distribution packages: dist/[platform]/"
    exit 0
fi

if [[ "$1" == "--clean-all" ]]; then
    echo -e "${YELLOW}🧹 Cleaning all build artifacts${NC}"
    rm -rf build/ dist/
    echo -e "${GREEN}✅ All build artifacts cleaned${NC}"
    exit 0
fi

# Detect platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Allow platform override via environment variable
if [[ -n "${PLATFORM:-}" ]]; then
    # Use provided platform override
    echo -e "${BLUE}🎯 Using platform override: ${PLATFORM}${NC}"
else
    case $OS in
        darwin) PLATFORM="mac" ;;
        linux)  PLATFORM="linux" ;;
        mingw*|cygwin*|msys*) PLATFORM="windows" ;;
        *)      PLATFORM="other" ;;
    esac
fi

case $ARCH in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) ARCH="amd64" ;;
esac

# For macOS, build both architectures
if [[ "$PLATFORM" == "mac" ]]; then
    ARCHITECTURES=("amd64" "arm64")
    echo -e "${BLUE}🚀 Building Servin for ${PLATFORM} (Universal: amd64 + arm64)${NC}"
elif [[ "$PLATFORM" == "windows" ]]; then
    # Windows primarily uses amd64
    ARCHITECTURES=("amd64")
    echo -e "${BLUE}🚀 Building Servin for ${PLATFORM}-amd64${NC}"
else
    ARCHITECTURES=("$ARCH")
    echo -e "${BLUE}🚀 Building Servin for ${PLATFORM}-${ARCH}${NC}"
fi

BUILD_DIR="build/${PLATFORM}"
DIST_DIR="dist/${PLATFORM}"

# Clean previous build artifacts
echo -e "${YELLOW}🧹 Cleaning previous build artifacts${NC}"
rm -rf "$BUILD_DIR" "$DIST_DIR"

# Create build directories
mkdir -p "$BUILD_DIR"
mkdir -p "$DIST_DIR"

# Function to build for a specific architecture
build_for_arch() {
    local target_arch=$1
    local arch_suffix=""
    local goos=""
    local exe_ext=""
    
    # Set platform-specific variables
    case $PLATFORM in
        mac)
            arch_suffix="-${target_arch}"
            goos="darwin"
            echo -e "${YELLOW}📦 Building for ${target_arch}${NC}"
            ;;
        linux)
            goos="linux"
            echo -e "${YELLOW}📦 Building for ${target_arch}${NC}"
            ;;
        windows)
            goos="windows"
            exe_ext=".exe"
            echo -e "${YELLOW}📦 Building for ${target_arch}${NC}"
            ;;
    esac
    
    local arch_build_dir="${BUILD_DIR}${arch_suffix}"
    mkdir -p "$arch_build_dir"
    
    # Build main CLI
    echo -e "${YELLOW}  📦 Building CLI (servin${arch_suffix}${exe_ext})${NC}"
    GOOS=${goos} GOARCH=${target_arch} go build -ldflags="-s -w" -o "${arch_build_dir}/servin${exe_ext}" ./main.go
    
    # Build Desktop
    echo -e "${YELLOW}  📦 Building Desktop (servin-desktop${arch_suffix}${exe_ext})${NC}"
    GOOS=${goos} GOARCH=${target_arch} go build -ldflags="-s -w" -o "${arch_build_dir}/servin-desktop${exe_ext}" ./cmd/servin-desktop/main.go
    
    return 0
}

# Function to build WebView GUI (architecture-independent)
build_webview_gui() {
    echo -e "${YELLOW}📦 Building WebView GUI${NC}"
    
    # Check if webview_gui directory exists
    if [[ ! -d "webview_gui" ]]; then
        echo -e "${YELLOW}  ⚠️ WebView GUI source not found, skipping...${NC}"
        return
    fi
    
    # Copy WebView GUI source files to build directory
    local webview_build_dir="$BUILD_DIR/webview_gui"
    mkdir -p "$webview_build_dir"
    cp -r webview_gui/* "$webview_build_dir/"
    
    # Create platform-specific launcher script in build directory
    case $PLATFORM in
        mac)
            cat > "$BUILD_DIR/servin-webview" << 'EOF'
#!/bin/bash

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WEBVIEW_DIR="$SCRIPT_DIR/webview_gui"

# Check if Python 3 is available
if ! command -v python3 >/dev/null 2>&1; then
    echo "Error: Python 3 is not installed or not in PATH"
    echo "Please install Python 3.7+ using Homebrew: brew install python3"
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
            chmod +x "$BUILD_DIR/servin-webview"
            ;;
        linux)
            cat > "$BUILD_DIR/servin-webview" << 'EOF'
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
            chmod +x "$BUILD_DIR/servin-webview"
            ;;
        windows)
            cat > "$BUILD_DIR/servin-webview.bat" << 'EOF'
@echo off
REM Servin WebView GUI Launcher for Windows

REM Get the directory where this script is located
set "SCRIPT_DIR=%~dp0"
set "WEBVIEW_DIR=%SCRIPT_DIR%webview_gui"

REM Check if Python 3 is available
python --version >nul 2>&1
if errorlevel 1 (
    echo Error: Python 3 is not installed or not in PATH
    echo Please install Python 3.7+ from https://python.org
    pause
    exit /b 1
)

REM Setup virtual environment if it doesn't exist
if not exist "%WEBVIEW_DIR%\venv" (
    echo Setting up WebView GUI environment...
    python -m venv "%WEBVIEW_DIR%\venv"
    call "%WEBVIEW_DIR%\venv\Scripts\activate.bat"
    pip install -r "%WEBVIEW_DIR%\requirements.txt"
) else (
    call "%WEBVIEW_DIR%\venv\Scripts\activate.bat"
)

REM Launch the WebView GUI
cd /d "%WEBVIEW_DIR%"
python main.py
EOF
            ;;
    esac
    
    echo -e "${GREEN}  ✅ WebView GUI built${NC}"
}

# Function to create distribution packages
create_distribution_package() {
    echo -e "${YELLOW}📦 Creating distribution package${NC}"
    
    case $PLATFORM in
        mac)
            # Create macOS installer package
            cat > "$DIST_DIR/install-servin.sh" << 'EOF'
#!/bin/bash
# Servin Container Runtime - macOS Installer
# Universal binary for Intel and Apple Silicon Macs

set -e

INSTALL_DIR="/usr/local/bin"
DATA_DIR="/usr/local/var/lib/servin"
CONFIG_DIR="/usr/local/etc/servin"
LAUNCHD_PLIST="/Library/LaunchDaemons/com.servin.runtime.plist"

echo "🐳 Installing Servin Container Runtime for macOS..."
echo ""

# Check for admin privileges
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root (use sudo)" 
   exit 1
fi

# Get the script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Install binaries
echo "📦 Installing binaries to $INSTALL_DIR..."
cp "$SCRIPT_DIR/servin" "$INSTALL_DIR/"
cp "$SCRIPT_DIR/servin-desktop" "$INSTALL_DIR/"
cp "$SCRIPT_DIR/servin-webview" "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR"/servin*

# Install WebView GUI
echo "🌐 Installing WebView GUI..."
mkdir -p "$DATA_DIR"
cp -r "$SCRIPT_DIR/webview_gui" "$DATA_DIR/"

# Create directories
echo "📁 Creating directories..."
mkdir -p "$DATA_DIR"/{containers,images,volumes,logs}
mkdir -p "$CONFIG_DIR"

# Create configuration
echo "⚙️ Creating configuration..."
cat > "$CONFIG_DIR/servin.conf" << 'CONF'
# Servin Configuration
data_dir = /usr/local/var/lib/servin
log_level = info
runtime = native
CONF

echo ""
echo "✅ Servin Container Runtime installed successfully!"
echo ""
echo "Usage:"
echo "  servin --help                 # CLI interface"
echo "  servin-desktop               # Terminal UI"  
echo "  servin-webview               # Modern web GUI"
echo ""
echo "The WebView GUI provides the best user experience with:"
echo "- Real-time container monitoring"
echo "- Image management dashboard"
echo "- Container logs viewer"
echo "- System statistics"
EOF
            chmod +x "$DIST_DIR/install-servin.sh"
            
            # Copy Python installer
            if [[ -f "installers/macos/servin-installer.py" ]]; then
                cp "installers/macos/servin-installer.py" "$DIST_DIR/"
            fi
            ;;
            
        linux)
            # Create Linux installer package
            cat > "$DIST_DIR/install-servin.sh" << 'EOF'
#!/bin/bash
# Servin Container Runtime - Linux Installer

set -e

INSTALL_DIR="/usr/local/bin"
DATA_DIR="/var/lib/servin"
CONFIG_DIR="/etc/servin"
SYSTEMD_SERVICE="/etc/systemd/system/servin.service"

echo "🐳 Installing Servin Container Runtime for Linux..."
echo ""

# Check for admin privileges
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root (use sudo)" 
   exit 1
fi

# Get the script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Install binaries
echo "📦 Installing binaries to $INSTALL_DIR..."
cp "$SCRIPT_DIR/servin" "$INSTALL_DIR/"
cp "$SCRIPT_DIR/servin-desktop" "$INSTALL_DIR/"
cp "$SCRIPT_DIR/servin-webview" "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR"/servin*

# Install WebView GUI
echo "🌐 Installing WebView GUI..."
mkdir -p "$DATA_DIR"
cp -r "$SCRIPT_DIR/webview_gui" "$DATA_DIR/"

# Create servin user
echo "👤 Creating servin user..."
if ! id "servin" &>/dev/null; then
    useradd -r -s /bin/false -d "$DATA_DIR" servin
fi

# Create directories
echo "📁 Creating directories..."
mkdir -p "$DATA_DIR"/{containers,images,volumes,logs}
mkdir -p "$CONFIG_DIR"
chown -R servin:servin "$DATA_DIR"

# Create configuration
echo "⚙️ Creating configuration..."
cat > "$CONFIG_DIR/servin.conf" << 'CONF'
# Servin Configuration
data_dir = /var/lib/servin
log_level = info
runtime = native
CONF

# Create systemd service
echo "🔧 Creating systemd service..."
cat > "$SYSTEMD_SERVICE" << 'SERVICE'
[Unit]
Description=Servin Container Runtime
After=network.target

[Service]
Type=simple
User=servin
ExecStart=/usr/local/bin/servin daemon --config /etc/servin/servin.conf
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
SERVICE

systemctl daemon-reload
systemctl enable servin

echo ""
echo "✅ Servin Container Runtime installed successfully!"
echo ""
echo "Usage:"
echo "  servin --help                 # CLI interface"
echo "  servin-desktop               # Terminal UI"  
echo "  servin-webview               # Modern web GUI"
echo "  sudo systemctl start servin  # Start service"
echo ""
echo "The WebView GUI provides the best user experience with:"
echo "- Real-time container monitoring"
echo "- Image management dashboard"
echo "- Container logs viewer"
echo "- System statistics"
EOF
            chmod +x "$DIST_DIR/install-servin.sh"
            
            # Copy Python installer
            if [[ -f "installers/linux/servin-installer.py" ]]; then
                cp "installers/linux/servin-installer.py" "$DIST_DIR/"
            fi
            ;;
        windows)
            # Create Windows installer package
            cat > "$DIST_DIR/install-servin.bat" << 'EOF'
@echo off
REM Servin Container Runtime - Windows Installer

setlocal EnableDelayedExpansion

set "INSTALL_DIR=%ProgramFiles%\Servin"
set "DATA_DIR=%ProgramData%\Servin"
set "CONFIG_DIR=%INSTALL_DIR%\config"

echo 🐳 Installing Servin Container Runtime for Windows...
echo.

REM Check for admin privileges
net session >nul 2>&1
if !errorlevel! neq 0 (
    echo This script must be run as Administrator
    echo Right-click and select "Run as administrator"
    pause
    exit /b 1
)

REM Get the script directory
set "SCRIPT_DIR=%~dp0"

REM Install binaries
echo 📦 Installing binaries to "%INSTALL_DIR%"...
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"
copy "%SCRIPT_DIR%servin.exe" "%INSTALL_DIR%\" > nul
copy "%SCRIPT_DIR%servin-desktop.exe" "%INSTALL_DIR%\" > nul
copy "%SCRIPT_DIR%servin-webview.bat" "%INSTALL_DIR%\" > nul

REM Install WebView GUI
echo 🌐 Installing WebView GUI...
if not exist "%DATA_DIR%" mkdir "%DATA_DIR%"
xcopy "%SCRIPT_DIR%webview_gui" "%DATA_DIR%\webview_gui" /E /I /Q > nul

REM Create directories
echo 📁 Creating directories...
if not exist "%DATA_DIR%\containers" mkdir "%DATA_DIR%\containers"
if not exist "%DATA_DIR%\images" mkdir "%DATA_DIR%\images"
if not exist "%DATA_DIR%\volumes" mkdir "%DATA_DIR%\volumes"
if not exist "%DATA_DIR%\logs" mkdir "%DATA_DIR%\logs"
if not exist "%CONFIG_DIR%" mkdir "%CONFIG_DIR%"

REM Create configuration
echo ⚙️ Creating configuration...
echo # Servin Configuration > "%CONFIG_DIR%\servin.conf"
echo data_dir = %DATA_DIR:\=/% >> "%CONFIG_DIR%\servin.conf"
echo log_level = info >> "%CONFIG_DIR%\servin.conf"
echo runtime = native >> "%CONFIG_DIR%\servin.conf"

REM Add to PATH
echo 🔧 Adding to system PATH...
setx PATH "%PATH%;%INSTALL_DIR%" /M > nul

REM Create Start Menu shortcuts
echo 📱 Creating Start Menu shortcuts...
set "START_MENU=%ProgramData%\Microsoft\Windows\Start Menu\Programs\Servin"
if not exist "%START_MENU%" mkdir "%START_MENU%"

REM Create CLI shortcut
echo @echo off > "%START_MENU%\Servin CLI.bat"
echo "%INSTALL_DIR%\servin.exe" %%* >> "%START_MENU%\Servin CLI.bat"

REM Create Desktop GUI shortcut
echo @echo off > "%START_MENU%\Servin Desktop.bat"
echo "%INSTALL_DIR%\servin-desktop.exe" >> "%START_MENU%\Servin Desktop.bat"

REM Create WebView GUI shortcut
echo @echo off > "%START_MENU%\Servin WebView.bat"
echo cd /d "%DATA_DIR%" >> "%START_MENU%\Servin WebView.bat"
echo "%INSTALL_DIR%\servin-webview.bat" >> "%START_MENU%\Servin WebView.bat"

echo.
echo ✅ Servin Container Runtime installed successfully!
echo.
echo Usage:
echo   servin --help                 # CLI interface
echo   servin-desktop               # Terminal UI  
echo   servin-webview               # Modern web GUI
echo.
echo The WebView GUI provides the best user experience with:
echo - Real-time container monitoring
echo - Image management dashboard
echo - Container logs viewer
echo - System statistics
echo.
echo Note: You may need to restart your command prompt to use servin commands.
pause
EOF
            
            # Copy Windows installer files
            if [[ -f "installers/windows/servin-installer.nsi" ]]; then
                cp "installers/windows/servin-installer.nsi" "$DIST_DIR/"
            fi
            if [[ -f "installers/windows/install.bat" ]]; then
                cp "installers/windows/install.bat" "$DIST_DIR/"
            fi
            ;;
    esac
    
    # Copy README and LICENSE to distribution
    if [[ -f "README.md" ]]; then
        cp "README.md" "$DIST_DIR/"
    fi
    if [[ -f "LICENSE" ]]; then
        cp "LICENSE" "$DIST_DIR/"
    fi
    
    echo -e "${GREEN}  ✅ Distribution package created${NC}"
}

# Function to create wizard installer packages
create_wizard_installer() {
    echo -e "${YELLOW}🧙 Creating wizard installer package${NC}"
    
    case $PLATFORM in
        mac)
            # Create macOS wizard installer package
            local installer_dir="$DIST_DIR/installer"
            mkdir -p "$installer_dir"
            
            # Copy installer wizard
            if [[ -f "installers/macos/servin-installer.py" ]]; then
                cp "installers/macos/servin-installer.py" "$installer_dir/"
                chmod +x "$installer_dir/servin-installer.py"
                
                # Create installer package directory with all components
                mkdir -p "$installer_dir/package"
                cp "$DIST_DIR/servin" "$installer_dir/package/"
                cp "$DIST_DIR/servin-desktop" "$installer_dir/package/"
                cp "$DIST_DIR/servin-webview" "$installer_dir/package/"
                cp -r "$DIST_DIR/webview_gui" "$installer_dir/package/"
                cp -r "$DIST_DIR/icons" "$installer_dir/package/"
                cp "$DIST_DIR/README.md" "$installer_dir/package/"
                cp "$DIST_DIR/LICENSE" "$installer_dir/package/"
                
                # Create launcher script for the wizard
                cat > "$DIST_DIR/ServinInstaller.command" << 'EOF'
#!/bin/bash
# Servin Installation Wizard Launcher
cd "$(dirname "$0")/installer"
python3 servin-installer.py
EOF
                chmod +x "$DIST_DIR/ServinInstaller.command"
                
                echo -e "${GREEN}  ✅ macOS wizard installer created: ServinInstaller.command${NC}"
            fi
            ;;
            
        linux)
            # Create Linux wizard installer package
            local installer_dir="$DIST_DIR/installer"
            mkdir -p "$installer_dir"
            
            # Copy installer wizard
            if [[ -f "installers/linux/servin-installer.py" ]]; then
                cp "installers/linux/servin-installer.py" "$installer_dir/"
                chmod +x "$installer_dir/servin-installer.py"
                
                # Create installer package directory with all components
                mkdir -p "$installer_dir/package"
                cp "$DIST_DIR/servin" "$installer_dir/package/"
                cp "$DIST_DIR/servin-desktop" "$installer_dir/package/"
                cp "$DIST_DIR/servin-webview" "$installer_dir/package/"
                cp -r "$DIST_DIR/webview_gui" "$installer_dir/package/"
                cp -r "$DIST_DIR/icons" "$installer_dir/package/"
                cp "$DIST_DIR/README.md" "$installer_dir/package/"
                cp "$DIST_DIR/LICENSE" "$installer_dir/package/"
                
                # Create launcher script for the wizard
                cat > "$DIST_DIR/ServinInstaller.sh" << 'EOF'
#!/bin/bash
# Servin Installation Wizard Launcher
cd "$(dirname "$0")/installer"
python3 servin-installer.py
EOF
                chmod +x "$DIST_DIR/ServinInstaller.sh"
                
                echo -e "${GREEN}  ✅ Linux wizard installer created: ServinInstaller.sh${NC}"
            fi
            ;;
            
        windows)
            # Create Windows wizard installer package
            local installer_dir="$DIST_DIR/installer"
            mkdir -p "$installer_dir"
            
            # Copy all installer components
            if [[ -f "installers/windows/servin-installer.nsi" ]]; then
                cp "installers/windows/"* "$installer_dir/"
                
                # Create installer package directory with all components
                mkdir -p "$installer_dir/package"
                cp "$DIST_DIR/servin.exe" "$installer_dir/package/"
                cp "$DIST_DIR/servin-desktop.exe" "$installer_dir/package/"
                cp "$DIST_DIR/servin-webview.bat" "$installer_dir/package/"
                cp -r "$DIST_DIR/webview_gui" "$installer_dir/package/"
                cp -r "$DIST_DIR/icons" "$installer_dir/package/"
                cp "$DIST_DIR/README.md" "$installer_dir/package/"
                cp "$DIST_DIR/LICENSE" "$installer_dir/package/"
                
                # Create build script for NSIS installer
                cat > "$installer_dir/build-installer.bat" << 'EOF'
@echo off
REM Build Servin Windows Installer
echo Building Servin Windows Installer...

REM Check if NSIS is available
where makensis >nul 2>&1
if errorlevel 1 (
    echo Error: NSIS not found in PATH
    echo Please install NSIS from https://nsis.sourceforge.io/
    echo Or install via Chocolatey: choco install nsis
    pause
    exit /b 1
)

REM Build the installer
makensis servin-installer.nsi

if exist "ServinSetup-1.0.0.exe" (
    echo.
    echo ✅ Installer built successfully: ServinSetup-1.0.0.exe
) else (
    echo.
    echo ❌ Installer build failed
)
pause
EOF
                
                # Create instructions
                cat > "$installer_dir/README-INSTALLER.txt" << 'EOF'
Servin Windows Installer Package
================================

This directory contains everything needed to create a Windows installer.

Quick Start:
1. Install NSIS from https://nsis.sourceforge.io/
2. Run: build-installer.bat
3. The installer ServinSetup-1.0.0.exe will be created

Advanced:
- Edit servin-installer.nsi to customize the installer
- All binaries are in the package/ directory
- Icons and resources are included

Requirements:
- NSIS (Nullsoft Scriptable Install System)
- Windows with PowerShell/Command Prompt
EOF
                
                echo -e "${GREEN}  ✅ Windows wizard installer package created in installer/${NC}"
                echo -e "${YELLOW}  📝 Run installer/build-installer.bat to create ServinSetup.exe${NC}"
            fi
            ;;
    esac
    
    echo -e "${GREEN}  ✅ Wizard installer package ready${NC}"
}

# Build for each architecture
for arch in "${ARCHITECTURES[@]}"; do
    build_for_arch "$arch"
done

# For macOS, create universal binaries
if [[ "$PLATFORM" == "mac" ]]; then
    echo -e "${YELLOW}🔗 Creating universal binaries${NC}"
    
    # Create universal servin binary
    echo -e "${YELLOW}  🔗 Creating universal servin${NC}"
    lipo -create \
        "${BUILD_DIR}-amd64/servin" \
        "${BUILD_DIR}-arm64/servin" \
        -output "${BUILD_DIR}/servin"
    
    # Create universal servin-desktop binary
    echo -e "${YELLOW}  🔗 Creating universal servin-desktop${NC}"
    lipo -create \
        "${BUILD_DIR}-amd64/servin-desktop" \
        "${BUILD_DIR}-arm64/servin-desktop" \
        -output "${BUILD_DIR}/servin-desktop"
    
    echo -e "${GREEN}✅ Universal binaries created${NC}"
else
    # For other platforms, organize build directory
    echo -e "${YELLOW}📁 Organizing build directory${NC}"
    
    # For single-architecture platforms, binaries are already in the right place
    # Just clean up any duplicate files if they exist
    if [[ "$PLATFORM" == "windows" ]]; then
        # Remove Unix-style binaries if they exist (keep only .exe)
        rm -f "${BUILD_DIR}/servin" "${BUILD_DIR}/servin-desktop" 2>/dev/null || true
    else
        # Remove Windows-style binaries if they exist (keep only Unix)
        rm -f "${BUILD_DIR}/servin.exe" "${BUILD_DIR}/servin-desktop.exe" 2>/dev/null || true
    fi
fi

# Build WebView GUI for all platforms
build_webview_gui

# Copy binaries and WebView GUI to distribution
echo -e "${YELLOW}📦 Copying binaries to distribution${NC}"

# Set file extensions based on platform
if [[ "$PLATFORM" == "windows" ]]; then
    cp "$BUILD_DIR/servin.exe" "$DIST_DIR/"
    if [[ -f "$BUILD_DIR/servin-desktop.exe" ]]; then
        cp "$BUILD_DIR/servin-desktop.exe" "$DIST_DIR/"
    fi
    cp "$BUILD_DIR/servin-webview.bat" "$DIST_DIR/"
else
    cp "$BUILD_DIR/servin" "$DIST_DIR/"
    if [[ -f "$BUILD_DIR/servin-desktop" ]]; then
        cp "$BUILD_DIR/servin-desktop" "$DIST_DIR/"
    fi
    cp "$BUILD_DIR/servin-webview" "$DIST_DIR/"
fi

cp -r "$BUILD_DIR/webview_gui" "$DIST_DIR/"

# Copy icons
echo -e "${YELLOW}🎨 Copying icons${NC}"
if [[ -d "icons" ]]; then
    mkdir -p "$DIST_DIR/icons"
    cp icons/* "$DIST_DIR/icons/" 2>/dev/null || true
fi

# Create distribution package with installers
create_distribution_package

# Create wizard installer packages
create_wizard_installer

# Copy documentation
cp README.md "${DIST_DIR}/"
cp LICENSE "${DIST_DIR}/"

# Make all executables
chmod +x "${DIST_DIR}"/servin*

echo -e "${GREEN}✅ Build completed successfully!${NC}"
echo ""
echo -e "${BLUE}📁 Build artifacts are in: ${BUILD_DIR}${NC}"
echo -e "${BLUE}📦 Distribution package is in: ${DIST_DIR}${NC}"
echo ""
echo -e "${YELLOW}Built binaries:${NC}"
ls -la "${BUILD_DIR}"/servin*

echo ""
echo -e "${YELLOW}Distribution package contents:${NC}"
ls -la "${DIST_DIR}"

# Show sizes and architecture info
echo ""
echo -e "${YELLOW}Binary sizes:${NC}"
du -h "${BUILD_DIR}"/servin*

# For macOS, show architecture information
if [[ "$PLATFORM" == "mac" ]]; then
    echo ""
    echo -e "${YELLOW}Architecture information:${NC}"
    for binary in "${BUILD_DIR}"/servin*; do
        if [[ -f "$binary" && -x "$binary" && ! "$binary" =~ webview ]]; then
            echo -e "${BLUE}$(basename "$binary"):${NC}"
            lipo -info "$binary" 2>/dev/null || file "$binary"
        fi
    done
fi

echo ""
echo -e "${GREEN}🚀 Ready for distribution!${NC}"
echo ""
echo -e "${BLUE}📦 Distribution Options:${NC}"

if [[ "$PLATFORM" == "windows" ]]; then
    echo -e "${YELLOW}  🧙 Wizard Installer: ${DIST_DIR}/installer/ (Run build-installer.bat to create .exe)${NC}"
    echo -e "${YELLOW}  📜 Quick Install: ${DIST_DIR}/install-servin.bat${NC}"
elif [[ "$PLATFORM" == "mac" ]]; then
    echo -e "${YELLOW}  🧙 Wizard Installer: ${DIST_DIR}/ServinInstaller.command (Double-click to run)${NC}"
    echo -e "${YELLOW}  📜 Quick Install: ${DIST_DIR}/install-servin.sh${NC}"
else
    echo -e "${YELLOW}  🧙 Wizard Installer: ${DIST_DIR}/ServinInstaller.sh${NC}"
    echo -e "${YELLOW}  📜 Quick Install: ${DIST_DIR}/install-servin.sh${NC}"
fi

echo ""
echo -e "${GREEN}✨ The wizard installer provides the best user experience with:${NC}"
echo -e "${GREEN}  • Interactive GUI installation${NC}"
echo -e "${GREEN}  • Automatic dependency checking${NC}"
echo -e "${GREEN}  • Custom installation paths${NC}"
echo -e "${GREEN}  • Professional system integration${NC}"