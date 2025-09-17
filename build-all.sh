#!/bin/bash
# Build all Servin components for distribution

set -e

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

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
            arch_suffix="-${target_arch}"
            goos="linux"
            echo -e "${YELLOW}📦 Building for ${target_arch}${NC}"
            ;;
        windows)
            arch_suffix="-${target_arch}"
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
    echo -e "${YELLOW}  📦 Building Desktop (servin-tui${arch_suffix}${exe_ext})${NC}"
    GOOS=${goos} GOARCH=${target_arch} go build -ldflags="-s -w" -o "${arch_build_dir}/servin-tui${exe_ext}" ./cmd/servin-tui/main.go
    
    return 0
}

# Function to build WebView GUI (architecture-independent)
build_webview_gui() {
    echo -e "${YELLOW}📦 Building WebView GUI with PyInstaller${NC}"
    
    # Check if webview_gui directory exists
    if [[ ! -d "webview_gui" ]]; then
        echo -e "${YELLOW}  ⚠️ WebView GUI source not found, skipping...${NC}"
        return
    fi
    
    # Check for cross-compilation limitations
    local host_os="$(uname -s | tr '[:upper:]' '[:lower:]')"
    if [[ "$PLATFORM" == "windows" && "$host_os" != "mingw"* && "$host_os" != "msys"* ]]; then
        echo -e "${YELLOW}  ⚠️ PyInstaller cannot cross-compile from $host_os to Windows${NC}"
        echo -e "${YELLOW}  ⚠️ Skipping WebView GUI build for Windows target...${NC}"
        echo -e "${YELLOW}  💡 To build Windows GUI: run this script on Windows${NC}"
        return
    fi
    
    # Set platform-specific variables
    local exe_ext=""
    local gui_name="servin-gui"
    
    case $PLATFORM in
        windows)
            exe_ext=".exe"
            gui_name="servin-gui.exe"
            ;;
    esac
    
    # Create temporary build directory for PyInstaller
    local temp_build_dir="$BUILD_DIR/webview_build_temp"
    mkdir -p "$temp_build_dir"
    
    # Copy webview_gui source to temp directory
    cp -r webview_gui/* "$temp_build_dir/"
    
    echo -e "${YELLOW}  🐍 Setting up Python environment...${NC}"
    cd "$temp_build_dir"
    
    # Check if Python 3 is available
    local python_cmd=""
    if command -v python3 >/dev/null 2>&1; then
        python_cmd="python3"
    elif command -v python >/dev/null 2>&1; then
        python_cmd="python"
    else
        echo -e "${RED}  ❌ Python not found, skipping WebView GUI build...${NC}"
        cd "$SCRIPT_DIR"
        return
    fi
    
    # Create virtual environment and install dependencies
    if ! $python_cmd -m venv venv 2>/dev/null; then
        echo -e "${RED}  ❌ Failed to create virtual environment, skipping WebView GUI build...${NC}"
        cd "$SCRIPT_DIR"
        return
    fi
    
    # Activate virtual environment (based on host OS, not target platform)
    if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
        source venv/Scripts/activate
    else
        source venv/bin/activate
    fi
    
    # Install dependencies including PyInstaller
    echo -e "${YELLOW}  📦 Installing dependencies...${NC}"
    if ! pip install -r requirements.txt >/dev/null 2>&1; then
        echo -e "${RED}  ❌ Failed to install dependencies, skipping WebView GUI build...${NC}"
        cd "$SCRIPT_DIR"
        return
    fi
    
    # Build with PyInstaller using spec file
    echo -e "${YELLOW}  🔨 Building executable with PyInstaller...${NC}"
    echo -e "${YELLOW}  📂 Building from directory: $(pwd)${NC}"
    echo -e "${YELLOW}  📋 Available files: $(ls -la)${NC}"
    
    # Enable verbose output for debugging
    if pyinstaller --clean --distpath=dist --workpath=build --log-level=INFO servin-gui.spec; then
        
        # Copy the built executable to the build directory
        if [[ -f "dist/$gui_name" ]]; then
            echo -e "${YELLOW}  📋 Copying $gui_name from $(pwd)/dist/ to $SCRIPT_DIR/$BUILD_DIR/${NC}"
            cp "dist/$gui_name" "$SCRIPT_DIR/$BUILD_DIR/"
            chmod +x "$SCRIPT_DIR/$BUILD_DIR/$gui_name"
            echo -e "${GREEN}  ✅ WebView GUI executable built: $gui_name${NC}"
        else
            echo -e "${RED}  ❌ PyInstaller succeeded but executable not found, skipping WebView GUI...${NC}"
            cd "$SCRIPT_DIR"
            return
        fi
    else
        echo -e "${RED}  ❌ PyInstaller failed, skipping WebView GUI build...${NC}"
        cd "$SCRIPT_DIR"
        return
    fi
    
    # Clean up temp directory
    cd "$SCRIPT_DIR"
    rm -rf "$temp_build_dir"
    
    echo -e "${GREEN}  ✅ WebView GUI build completed${NC}"
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

# Install binaries
echo "📦 Installing binaries to $INSTALL_DIR..."
cp "$SCRIPT_DIR/servin" "$INSTALL_DIR/"
cp "$SCRIPT_DIR/servin-tui" "$INSTALL_DIR/"
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
echo "  servin-tui               # Terminal UI"  
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

# Install binaries
echo "📦 Installing binaries to $INSTALL_DIR..."
cp "$SCRIPT_DIR/servin" "$INSTALL_DIR/"
cp "$SCRIPT_DIR/servin-tui" "$INSTALL_DIR/"
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
echo "  servin-tui               # Terminal UI"  
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
copy "%SCRIPT_DIR%servin-tui.exe" "%INSTALL_DIR%\" > nul
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
echo "%INSTALL_DIR%\servin-tui.exe" >> "%START_MENU%\Servin Desktop.bat"

REM Create WebView GUI shortcut
echo @echo off > "%START_MENU%\Servin WebView.bat"
echo cd /d "%DATA_DIR%" >> "%START_MENU%\Servin WebView.bat"
echo "%INSTALL_DIR%\servin-webview.bat" >> "%START_MENU%\Servin WebView.bat"

echo.
echo ✅ Servin Container Runtime installed successfully!
echo.
echo Usage:
echo   servin --help                 # CLI interface
echo   servin-tui               # Terminal UI  
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

# Function to create macOS .dmg file
create_macos_dmg() {
    if [[ "$PLATFORM" != "mac" ]]; then
        return
    fi
    
    echo -e "${YELLOW}📦 Creating macOS .dmg file${NC}"
    
    local dmg_name="Servin-Container-Runtime"
    local dmg_temp_dir="$BUILD_DIR/dmg_temp"
    local dmg_file="$DIST_DIR/${dmg_name}.dmg"
    
    # Clean up any existing DMG temp directory
    rm -rf "$dmg_temp_dir"
    mkdir -p "$dmg_temp_dir"
    
    # Create Servin.app bundle structure
    local app_dir="$dmg_temp_dir/Servin.app"
    mkdir -p "$app_dir/Contents/MacOS"
    mkdir -p "$app_dir/Contents/Resources"
    
    # Copy binaries to app bundle
    cp "$DIST_DIR/servin" "$app_dir/Contents/MacOS/"
    cp "$DIST_DIR/servin-tui" "$app_dir/Contents/MacOS/"
    if [[ -f "$DIST_DIR/servin-gui" ]]; then
        cp "$DIST_DIR/servin-gui" "$app_dir/Contents/MacOS/"
    fi
    
    # Create Info.plist for the app bundle
    cat > "$app_dir/Contents/Info.plist" << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>servin-gui</string>
    <key>CFBundleIdentifier</key>
    <string>com.servin.runtime</string>
    <key>CFBundleName</key>
    <string>Servin</string>
    <key>CFBundleDisplayName</key>
    <string>Servin Container Runtime</string>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.0</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleSignature</key>
    <string>????</string>
    <key>NSHighResolutionCapable</key>
    <true/>
</dict>
</plist>
EOF
    
    # Copy icon if available
    if [[ -f "icons/tool_icon.png" ]]; then
        cp "icons/tool_icon.png" "$app_dir/Contents/Resources/AppIcon.png"
    fi
    
    # Copy documentation
    cp "$DIST_DIR/README.md" "$dmg_temp_dir/"
    cp "$DIST_DIR/LICENSE" "$dmg_temp_dir/"
    
    # Create Applications symlink
    ln -sf /Applications "$dmg_temp_dir/Applications"
    
    # Create the DMG
    echo -e "${YELLOW}  🔨 Building DMG file...${NC}"
    
    # Remove any existing DMG
    rm -f "$dmg_file"
    
    # Create the DMG directly from the temp directory
    if hdiutil create -srcfolder "$dmg_temp_dir" -volname "Servin Container Runtime" -format UDZO -imagekey zlib-level=9 "$dmg_file" >/dev/null 2>&1; then
        local dmg_size=$(ls -lh "$dmg_file" | awk '{print $5}')
        echo -e "${GREEN}  ✅ macOS DMG created: ${dmg_name}.dmg (${dmg_size})${NC}"
    else
        echo -e "${RED}  ❌ Failed to create DMG file${NC}"
    fi
    
    # Clean up
    rm -rf "$dmg_temp_dir"
}

# Function to create wizard installer packages
create_wizard_installer() {
    echo -e "${YELLOW}🧙 Creating wizard installer package${NC}"
    
    case $PLATFORM in
        mac)
            # Create .dmg file first
            create_macos_dmg
            
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
                cp "$DIST_DIR/servin-tui" "$installer_dir/package/"
                if [[ -f "$DIST_DIR/servin-gui" ]]; then
                    cp "$DIST_DIR/servin-gui" "$installer_dir/package/"
                fi
                cp -r "$DIST_DIR/icons" "$installer_dir/package/"
                cp "$DIST_DIR/README.md" "$installer_dir/package/"
                cp "$DIST_DIR/LICENSE" "$installer_dir/package/"
                
                # Create launcher script for the wizard
                cat > "$DIST_DIR/ServinInstaller.command" << 'EOF'
#!/bin/bash
# Servin Installation Wizard Launcher

# Get the directory containing this script
SCRIPT_DIR="$(dirname "$0")"
INSTALLER_DIR="$SCRIPT_DIR/installer"

echo "🐳 Servin Container Runtime Installer"
echo "====================================="
echo ""

# Check if running with sudo
if [[ $EUID -eq 0 ]]; then
    echo "✅ Running with administrator privileges"
    cd "$INSTALLER_DIR"
    python3 servin-installer.py
else
    echo "⚠️  Administrator privileges required for system installation"
    echo ""
    echo "This installer will:"
    echo "• Install Servin binaries to /usr/local/bin"
    echo "• Create system directories in /usr/local"
    echo "• Install a system service (launchd)"
    echo "• Create an Application bundle"
    echo ""
    read -p "Continue with installation? (y/N): " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "🔐 Requesting administrator privileges..."
        sudo bash "$0"
    else
        echo "Installation cancelled."
        exit 0
    fi
fi
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
                cp "$DIST_DIR/servin-tui" "$installer_dir/package/"
                if [[ -f "$DIST_DIR/servin-gui" ]]; then
                    cp "$DIST_DIR/servin-gui" "$installer_dir/package/"
                fi
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
                cp "$DIST_DIR/servin-tui.exe" "$installer_dir/package/"
                if [[ -f "$DIST_DIR/servin-gui.exe" ]]; then
                    cp "$DIST_DIR/servin-gui.exe" "$installer_dir/package/"
                fi
                cp -r "$DIST_DIR/icons" "$installer_dir/package/"
                
                # Copy and rename files to match NSIS expectations
                cp "$DIST_DIR/README.md" "$installer_dir/package/README.txt"
                cp "$DIST_DIR/LICENSE" "$installer_dir/package/LICENSE.txt"
                
                # Copy configuration file if it exists
                if [[ -f "installers/windows/servin.conf" ]]; then
                    cp "installers/windows/servin.conf" "$installer_dir/package/"
                fi
                
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
echo "DEBUG: About to check platform for binary organization"
echo "DEBUG: PLATFORM = '$PLATFORM'"

if [[ "$PLATFORM" == "mac" ]]; then
    echo -e "${YELLOW}🔗 Creating universal binaries${NC}"
    
    # Create universal servin binary
    echo -e "${YELLOW}  🔗 Creating universal servin${NC}"
    lipo -create \
        "${BUILD_DIR}-amd64/servin" \
        "${BUILD_DIR}-arm64/servin" \
        -output "${BUILD_DIR}/servin"
    
    # Create universal servin-tui binary
    echo -e "${YELLOW}  🔗 Creating universal servin-tui${NC}"
    lipo -create \
        "${BUILD_DIR}-amd64/servin-tui" \
        "${BUILD_DIR}-arm64/servin-tui" \
        -output "${BUILD_DIR}/servin-tui"
    
    echo -e "${GREEN}✅ Universal binaries created${NC}"
else
    # For other platforms, organize build directory
    echo -e "${YELLOW}📁 Organizing build directory${NC}"
    echo "DEBUG: In else block, PLATFORM = '$PLATFORM'"
    
    # For single-architecture platforms, copy binaries from arch-specific directories
    if [[ "$PLATFORM" == "windows" ]]; then
        # Copy Windows binaries from arch-specific directory to main build directory
        echo -e "${YELLOW}  📁 Copying Windows binaries${NC}"
        echo "    Source: ${BUILD_DIR}-amd64/servin.exe"
        echo "    Target: ${BUILD_DIR}/servin.exe"
        echo "DEBUG: About to check if source file exists"
        
        if [[ ! -f "${BUILD_DIR}-amd64/servin.exe" ]]; then
            echo "    ❌ ERROR: Source file does not exist: ${BUILD_DIR}-amd64/servin.exe"
            echo "    Directory contents of ${BUILD_DIR}-amd64/:"
            ls -la "${BUILD_DIR}-amd64/" 2>/dev/null || echo "    Directory does not exist"
            return 1
        fi
        
        echo "DEBUG: Source file exists, about to copy"
        cp "${BUILD_DIR}-amd64/servin.exe" "${BUILD_DIR}/" || {
            echo "    ❌ ERROR: Failed to copy servin.exe"
            return 1
        }
        echo "    ✅ servin.exe copied successfully"
        
        if [[ -f "${BUILD_DIR}-amd64/servin-tui.exe" ]]; then
            cp "${BUILD_DIR}-amd64/servin-tui.exe" "${BUILD_DIR}/" || {
                echo "    ❌ ERROR: Failed to copy servin-tui.exe"
                return 1
            }
            echo "    ✅ servin-tui.exe copied successfully"
        fi
        # Remove Unix-style binaries if they exist (keep only .exe)
        rm -f "${BUILD_DIR}/servin" "${BUILD_DIR}/servin-tui" 2>/dev/null || true
    else
        # Copy Linux binaries from arch-specific directory to main build directory  
        echo -e "${YELLOW}  📁 Copying Linux binaries${NC}"
        arch_dir="${BUILD_DIR}-${ARCHITECTURES[0]}"
        if [[ -d "$arch_dir" ]]; then
            cp "${arch_dir}/servin" "${BUILD_DIR}/"
            if [[ -f "${arch_dir}/servin-tui" ]]; then
                cp "${arch_dir}/servin-tui" "${BUILD_DIR}/"
            fi
        fi
        # Remove Windows-style binaries if they exist (keep only Unix)
        rm -f "${BUILD_DIR}/servin.exe" "${BUILD_DIR}/servin-tui.exe" 2>/dev/null || true
    fi
fi

# Build WebView GUI for all platforms
build_webview_gui

# Copy binaries and WebView GUI to distribution
echo -e "${YELLOW}📦 Copying binaries to distribution${NC}"

# Set file extensions based on platform
if [[ "$PLATFORM" == "windows" ]]; then
    echo "  Checking for Windows binaries in: $BUILD_DIR"
    
    # Check main build directory first
    if [[ -f "$BUILD_DIR/servin.exe" ]]; then
        echo "  ✅ Found servin.exe in main build directory"
        cp "$BUILD_DIR/servin.exe" "$DIST_DIR/"
        echo "  ✅ servin.exe copied to distribution"
    elif [[ -f "${BUILD_DIR}-amd64/servin.exe" ]]; then
        echo "  ⚠️ servin.exe not in main directory, copying from architecture-specific directory"
        cp "${BUILD_DIR}-amd64/servin.exe" "$DIST_DIR/"
        echo "  ✅ servin.exe copied to distribution from ${BUILD_DIR}-amd64/"
    else
        echo "  ❌ ERROR: servin.exe not found in either location"
        echo "  Main directory contents:"
        ls -la "$BUILD_DIR/" 2>/dev/null || echo "  Main build directory does not exist"
        echo "  Architecture directory contents:"
        ls -la "${BUILD_DIR}-amd64/" 2>/dev/null || echo "  Architecture directory does not exist"
        exit 1
    fi
    
    # Check for desktop binary
    if [[ -f "$BUILD_DIR/servin-tui.exe" ]]; then
        cp "$BUILD_DIR/servin-tui.exe" "$DIST_DIR/"
        echo "  ✅ servin-tui.exe copied to distribution"
    elif [[ -f "${BUILD_DIR}-amd64/servin-tui.exe" ]]; then
        echo "  ⚠️ servin-tui.exe not in main directory, copying from architecture-specific directory"
        cp "${BUILD_DIR}-amd64/servin-tui.exe" "$DIST_DIR/"
        echo "  ✅ servin-tui.exe copied to distribution from ${BUILD_DIR}-amd64/"
    fi
    
    # Copy GUI binary if it exists
    if [[ -f "$BUILD_DIR/servin-gui.exe" ]]; then
        cp "$BUILD_DIR/servin-gui.exe" "$DIST_DIR/"
        echo "  ✅ servin-gui.exe copied to distribution"
    fi
else
    cp "$BUILD_DIR/servin" "$DIST_DIR/"
    if [[ -f "$BUILD_DIR/servin-tui" ]]; then
        cp "$BUILD_DIR/servin-tui" "$DIST_DIR/"
    fi
    
    # Copy GUI binary if it exists
    if [[ -f "$BUILD_DIR/servin-gui" ]]; then
        cp "$BUILD_DIR/servin-gui" "$DIST_DIR/"
        echo "  ✅ servin-gui copied to distribution"
    fi
fi

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
    echo -e "${YELLOW}  📦 macOS Disk Image: ${DIST_DIR}/Servin-Container-Runtime.dmg (Drag & Drop installer)${NC}"
    echo -e "${YELLOW}  🧙 Wizard Installer: ${DIST_DIR}/ServinInstaller.command (Double-click to run)${NC}"
    echo -e "${YELLOW}  📜 Quick Install: ${DIST_DIR}/install-servin.sh${NC}"
else
    echo -e "${YELLOW}  🧙 Wizard Installer: ${DIST_DIR}/ServinInstaller.sh${NC}"
    echo -e "${YELLOW}  📜 Quick Install: ${DIST_DIR}/install-servin.sh${NC}"
fi

echo "All wizards created successfully!"
echo -e "${GREEN}✨ The wizard installer provides the best user experience with:${NC}"
echo -e "${GREEN}  • Interactive GUI installation${NC}"
echo -e "${GREEN}  • Automatic dependency checking${NC}"
echo -e "${GREEN}  • Custom installation paths${NC}"
echo -e "${GREEN}  • Professional system integration${NC}"