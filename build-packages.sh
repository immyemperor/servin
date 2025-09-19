#!/bin/bash

# Servin Container Runtime - Cross-Platform Package Builder
# Builds complete installer packages for Windows, Linux, and macOS with embedded VM dependencies

# Color output functions (defined first)
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
NC='\033[0m'

print_success() { echo -e "${GREEN}✓ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠ $1${NC}"; }
print_error() { echo -e "${RED}✗ $1${NC}"; }
print_info() { echo -e "${BLUE}→ $1${NC}"; }
print_header() { echo -e "\n${CYAN}${BOLD}$1${NC}"; }
print_platform() { echo -e "${MAGENTA}${BOLD}$1${NC}"; }

# Enhanced error handling
set -euo pipefail
IFS=$'\n\t'

# Error trap for debugging (now print_error is available)
error_exit() {
    local line_no=$1
    local error_code=$2
    print_error "Script failed at line $line_no with exit code $error_code"
    print_error "Command: ${BASH_COMMAND}"
    print_error "Working directory: $(pwd)"
    print_error "Available files:"
    ls -la 2>/dev/null || echo "Cannot list files"
    exit $error_code
}

trap 'error_exit ${LINENO} $?' ERR

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VERSION="1.0.0"
BUILD_DATE=$(date +"%Y%m%d")
BUILD_TIME=$(date +"%H%M%S")

print_info "Build started at $(date)"
print_info "Script directory: $SCRIPT_DIR"
print_info "Working directory: $(pwd)"

print_banner() {
    echo -e "${CYAN}${BOLD}"
    echo "╔══════════════════════════════════════════════════════════════════════╗"
    echo "║              Servin Container Runtime - Package Builder             ║"
    echo "║                     Cross-Platform Installer Creator                ║"
    echo "╚══════════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}\n"
}

# Detect current platform
detect_platform() {
    case "$(uname -s)" in
        Linux*)     PLATFORM="linux";;
        Darwin*)    PLATFORM="macos";;
        CYGWIN*|MINGW*|MSYS*) PLATFORM="windows";;
        *)          PLATFORM="unknown";;
    esac
    
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64|amd64) ARCH="amd64";;
        aarch64|arm64) ARCH="arm64";;
        armv7l) ARCH="arm";;
        i386|i686) ARCH="386";;
    esac
}

# Check prerequisites
check_prerequisites() {
    print_header "Checking Prerequisites"
    
    # Check if Servin is built (check multiple possible locations)
    local servin_binary=""
    if [[ -f "$SCRIPT_DIR/servin" ]]; then
        servin_binary="$SCRIPT_DIR/servin"
    elif [[ -f "$SCRIPT_DIR/servin" ]]; then
        servin_binary="$SCRIPT_DIR/servin"
    elif [[ -f "$SCRIPT_DIR/build/*/servin" ]]; then
        servin_binary=$(find "$SCRIPT_DIR/build" -name "servin" -type f | head -1)
    fi
    
    if [[ -z "$servin_binary" || ! -f "$servin_binary" ]]; then
        print_error "Servin executable not found. Building..."
        cd "$SCRIPT_DIR"
        
        # Ensure we're in the correct directory with go.mod
        if [[ ! -f "go.mod" ]]; then
            print_error "go.mod not found. This script must be run from the project root."
            exit 1
        fi
        
        # Build for current platform using proper Go module syntax
        print_info "Building Servin using Go modules..."
        if [[ -f "Makefile" ]]; then
            make build
        elif [[ -f "build.sh" ]]; then
            ./build.sh
        else
            # Use proper Go module build command
            go build -o servin .
        fi
        
        if [[ ! -f "servin" ]]; then
            print_error "Failed to build Servin"
            exit 1
        fi
        print_success "Servin built successfully"
    fi
    
    print_success "Prerequisites check passed"
}

# Build cross-platform executables
build_executables() {
    local target_platforms=("$@")
    
    if [[ ${#target_platforms[@]} -eq 0 ]]; then
        # Default: build all platforms
        target_platforms=("windows/amd64" "linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64")
        print_header "Building Cross-Platform Executables (All Platforms)"
    else
        print_header "Building Executables for Selected Platforms"
    fi
    
    cd "$SCRIPT_DIR"
    
    for target in "${target_platforms[@]}"; do
        local goos=$(echo "$target" | cut -d'/' -f1)
        local goarch=$(echo "$target" | cut -d'/' -f2)
        local output_dir="build/${goos}-${goarch}"
        local binary_name="servin"
        
        if [[ "$goos" == "windows" ]]; then
            binary_name="servin.exe"
        fi
        
        print_info "Building for $goos/$goarch..."
        
        mkdir -p "$output_dir"
        
        # Build main executable
        GOOS="$goos" GOARCH="$goarch" go build -ldflags="-s -w" -o "$output_dir/$binary_name" .
        
        # Build TUI if cmd/servin-tui exists
        if [[ -d "cmd/servin-tui" ]]; then
            local tui_name="servin-tui"
            if [[ "$goos" == "windows" ]]; then
                tui_name="servin-tui.exe"
            fi
            GOOS="$goos" GOARCH="$goarch" go build -ldflags="-s -w" -o "$output_dir/$tui_name" ./cmd/servin-tui/
        fi
        
        # Build GUI for Windows (requires Python environment)
        if [[ "$goos" == "windows" ]] && [[ -d "webview_gui" ]]; then
            local gui_name="servin-gui.exe"
            print_info "Building GUI executable for Windows..."
            
            # Check if Python and PyInstaller are available
            if command -v python3 >/dev/null 2>&1; then
                print_info "Python3 found: $(python3 --version)"
                
                if python3 -c "import PyInstaller" 2>/dev/null; then
                    print_info "PyInstaller available"
                    cd webview_gui
                    
                    # Install GUI dependencies if needed
                    print_info "Installing GUI requirements..."
                    python3 -m pip install -r requirements.txt --quiet 2>/dev/null || {
                        print_warning "Failed to install GUI requirements, trying without upgrade"
                        python3 -m pip install --no-deps -r requirements.txt --quiet 2>/dev/null || print_warning "GUI requirements installation failed"
                    }
                    
                    # Build GUI executable with better error handling
                    print_info "Building GUI with PyInstaller..."
                    if python3 -m PyInstaller --clean --onefile servin-gui.spec --distpath "../$output_dir" --workpath "../build/gui-work" --specpath . --log-level WARN 2>&1; then
                        if [[ -f "../$output_dir/$gui_name" ]]; then
                            print_success "GUI executable built: $gui_name"
                        else
                            print_warning "PyInstaller completed but GUI executable not found at expected location"
                            ls -la "../$output_dir/" || true
                        fi
                    else
                        print_warning "PyInstaller build failed, installer will be CLI-only"
                        print_info "This is acceptable - Windows installer will work without GUI"
                    fi
                    
                    cd ..
                else
                    print_warning "PyInstaller not available, skipping GUI build"
                fi
            else
                print_warning "Python3 not available, skipping GUI build"
            fi
        fi
        
        print_success "Built $goos/$goarch"
    done
    
    print_success "Cross-platform executables built"
}

# Build Windows installer
build_windows_installer() {
    print_platform "Building Windows Installer (NSIS)"
    
    local windows_dir="$SCRIPT_DIR/installers/windows"
    
    # Ensure Windows executables exist
    local main_exe_path=""
    if [[ -f "$SCRIPT_DIR/build/windows-amd64/servin.exe" ]]; then
        main_exe_path="$SCRIPT_DIR/build/windows-amd64/servin.exe"
    elif [[ -f "$SCRIPT_DIR/build/windows/servin.exe" ]]; then
        main_exe_path="$SCRIPT_DIR/build/windows/servin.exe"
    fi
    
    if [[ -z "$main_exe_path" ]]; then
        print_error "servin.exe not found in either location:"
        print_error "  - $SCRIPT_DIR/build/windows-amd64/servin.exe"
        print_error "  - $SCRIPT_DIR/build/windows/servin.exe"
        print_error "Build executables first with: ./build-all.sh"
        return 1
    fi
    
    # Copy Windows executables
    print_info "Copying Windows executables..."
    cp "$main_exe_path" "$windows_dir/"
    print_success "Main executable copied from: $main_exe_path"
    
    # Copy TUI executable if available
    local tui_copied=false
    if [[ -f "$SCRIPT_DIR/build/windows-amd64/servin-tui.exe" ]]; then
        cp "$SCRIPT_DIR/build/windows-amd64/servin-tui.exe" "$windows_dir/"
        tui_copied=true
        print_success "TUI executable copied (from windows-amd64)"
    elif [[ -f "$SCRIPT_DIR/build/windows/servin-tui.exe" ]]; then
        cp "$SCRIPT_DIR/build/windows/servin-tui.exe" "$windows_dir/"
        tui_copied=true
        print_success "TUI executable copied (from windows)"
    fi
    
    if [[ "$tui_copied" == "false" ]]; then
        print_warning "servin-tui.exe not found - installer will be main CLI only"
    fi
    
    # Copy GUI executable if it was built
    local gui_copied=false
    
    # Check both possible locations for GUI executable
    if [[ -f "$SCRIPT_DIR/build/windows-amd64/servin-gui.exe" ]]; then
        cp "$SCRIPT_DIR/build/windows-amd64/servin-gui.exe" "$windows_dir/"
        gui_copied=true
        print_success "GUI executable included in Windows installer (from windows-amd64)"
    elif [[ -f "$SCRIPT_DIR/build/windows/servin-gui.exe" ]]; then
        cp "$SCRIPT_DIR/build/windows/servin-gui.exe" "$windows_dir/"
        gui_copied=true
        print_success "GUI executable included in Windows installer (from windows)"
    fi
    
    if [[ "$gui_copied" == "false" ]]; then
        print_warning "GUI executable not found in either location:"
        print_warning "  - $SCRIPT_DIR/build/windows-amd64/servin-gui.exe"
        print_warning "  - $SCRIPT_DIR/build/windows/servin-gui.exe"
        print_info "Building CLI-only installer"
        print_info "Windows installer will work without GUI components"
    fi
    
    # Copy required files
    print_info "Copying configuration files..."
    [[ -f "$SCRIPT_DIR/LICENSE" ]] && cp "$SCRIPT_DIR/LICENSE" "$windows_dir/LICENSE.txt"
    
    cd "$windows_dir"
    
    # Check if we can build on this platform
    if [[ "$PLATFORM" == "windows" ]] && command -v makensis >/dev/null 2>&1; then
        print_info "Building NSIS installer natively..."
        
        # Enhanced Windows command execution with better debugging
        if [[ "$OS" == "Windows_NT" ]] || command -v cmd.exe >/dev/null 2>&1; then
            print_info "Running Windows batch file via cmd.exe..."
            
            # Direct NSIS execution with enhanced error detection
            print_info "Attempting direct NSIS compilation with error detection..."
            
            # Check NSIS availability and path
            print_info "Checking NSIS installation..."
            
            # Try common NSIS locations on Windows
            NSIS_PATHS=(
                "/c/Program Files (x86)/NSIS"
                "/c/Program Files/NSIS"
                "/cygdrive/c/Program Files (x86)/NSIS"
                "/cygdrive/c/Program Files/NSIS"
            )
            
            MAKENSIS_PATH=""
            
            # Check if makensis is in PATH
            if cmd.exe /c "where makensis" >/dev/null 2>&1; then
                MAKENSIS_PATH="makensis"
                print_info "NSIS makensis found in PATH"
            else
                # Try to find NSIS in common locations
                for nsis_path in "${NSIS_PATHS[@]}"; do
                    if [[ -f "$nsis_path/makensis.exe" ]]; then
                        MAKENSIS_PATH="\"$nsis_path/makensis.exe\""
                        print_info "NSIS found at: $nsis_path/makensis.exe"
                        break
                    fi
                done
            fi
            
            if [[ -n "$MAKENSIS_PATH" ]]; then
                # Get NSIS version for verification
                print_info "Checking NSIS version..."
                cmd.exe /c "$MAKENSIS_PATH /VERSION" 2>/dev/null || print_warning "Could not get NSIS version"
                
                print_info "NSIS found, attempting compilation..."
                
                # Create a temporary batch file with proper error handling
                cat > "build-debug.bat" << 'BATCH_EOF'
@echo off
setlocal enabledelayedexpansion

echo ========================================
echo NSIS Build Debug Script
echo ========================================
echo Current directory: %CD%
echo Date/Time: %DATE% %TIME%
echo.

echo === File Check ===
dir *.nsi *.exe 2>nul

echo.
echo === NSIS Version Check ===
makensis /VERSION 2>nul
if !errorlevel! neq 0 (
    echo ERROR: makensis not found in PATH
    echo Trying common NSIS locations...
    if exist "C:\Program Files (x86)\NSIS\makensis.exe" (
        set "MAKENSIS=C:\Program Files (x86)\NSIS\makensis.exe"
        echo Found NSIS at: !MAKENSIS!
    ) else if exist "C:\Program Files\NSIS\makensis.exe" (
        set "MAKENSIS=C:\Program Files\NSIS\makensis.exe"
        echo Found NSIS at: !MAKENSIS!
    ) else (
        echo ERROR: NSIS not found in common locations
        exit /b 1
    )
) else (
    set "MAKENSIS=makensis"
    echo NSIS found in PATH
)

echo.
echo === Building NSIS Installer ===
echo Running: "!MAKENSIS!" /V4 servin-installer.nsi
"!MAKENSIS!" /V4 servin-installer.nsi
set NSIS_EXIT=!errorlevel!
echo NSIS Exit Code: !NSIS_EXIT!

echo.
echo === Build Results ===
if !NSIS_EXIT! equ 0 (
    echo NSIS compilation completed successfully
    dir *installer*.exe 2>nul
    if exist "servin-installer-1.0.0.exe" (
        echo SUCCESS: servin-installer-1.0.0.exe created
        copy "servin-installer-1.0.0.exe" "Servin-Installer-1.0.0.exe" >nul 2>&1
        echo Copied to: Servin-Installer-1.0.0.exe
        exit /b 0
    ) else (
        echo WARNING: NSIS succeeded but expected file not found
        echo Available files:
        dir *.exe
        exit /b 2
    )
) else (
    echo ERROR: NSIS compilation failed with exit code !NSIS_EXIT!
    echo Available files:
    dir *.exe 2>nul
    echo.
    echo Checking required files for NSIS:
    if exist "servin.exe" (echo ✓ servin.exe found) else (echo ✗ servin.exe missing)
    if exist "servin-tui.exe" (echo ✓ servin-tui.exe found) else (echo ✗ servin-tui.exe missing)
    if exist "servin-gui.exe" (echo ✓ servin-gui.exe found) else (echo ✗ servin-gui.exe missing)
    if exist "servin.conf" (echo ✓ servin.conf found) else (echo ✗ servin.conf missing)
    if exist "LICENSE.txt" (echo ✓ LICENSE.txt found) else (echo ✗ LICENSE.txt missing)
    if exist "servin-installer.nsi" (echo ✓ servin-installer.nsi found) else (echo ✗ servin-installer.nsi missing)
    exit /b !NSIS_EXIT!
)
BATCH_EOF

                # Run the debug batch file with explicit output capture
                print_info "Running debug batch file with comprehensive error checking..."
                echo "=== BATCH FILE OUTPUT START ==="
                
                # Execute batch file and capture output explicitly
                if cmd.exe /c "build-debug.bat" > nsis-output.log 2>&1; then
                    echo "Batch file execution completed"
                    echo "=== CAPTURED NSIS OUTPUT ==="
                    cat nsis-output.log
                    echo "=== END NSIS OUTPUT ==="
                    
                    # Check if installer was actually created
                    if ls -la *installer*.exe 2>/dev/null | grep -q installer; then
                        print_success "NSIS installer build completed successfully"
                    else
                        print_error "NSIS batch completed but no installer file created"
                        echo "NSIS output log contents:"
                        cat nsis-output.log
                        
                        # Try direct NSIS compilation as fallback
                        print_info "Attempting direct NSIS compilation as fallback..."
                        if cmd.exe /c "makensis /V4 servin-installer.nsi" > direct-nsis.log 2>&1; then
                            print_info "Direct NSIS compilation attempted"
                            echo "Direct NSIS output:"
                            cat direct-nsis.log
                            
                            if ls -la *installer*.exe 2>/dev/null | grep -q installer; then
                                print_success "Direct NSIS compilation succeeded"
                            else
                                print_error "Direct NSIS compilation also failed"
                                print_info "Creating minimal NSIS script as last resort..."
                                
                                # Create a minimal NSIS script
                                cat > "servin-minimal.nsi" << 'MINIMAL_NSI'
!define PRODUCT_NAME "Servin Container Runtime"
!define PRODUCT_VERSION "1.0.0"

Name "${PRODUCT_NAME}"
OutFile "servin-installer-${PRODUCT_VERSION}.exe"
InstallDir "$PROGRAMFILES\Servin"
RequestExecutionLevel admin

Page directory
Page instfiles

Section "MainSection" SEC01
  SetOutPath "$INSTDIR"
  File "servin.exe"
  File "servin-tui.exe"
  File "servin-gui.exe"
  File "servin.conf"
  CreateDirectory "$SMPROGRAMS\Servin"
  CreateShortCut "$SMPROGRAMS\Servin\Servin.lnk" "$INSTDIR\servin.exe"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin" "DisplayName" "${PRODUCT_NAME}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin" "UninstallString" "$INSTDIR\uninstall.exe"
  WriteUninstaller "$INSTDIR\uninstall.exe"
SectionEnd

Section "Uninstall"
  Delete "$INSTDIR\servin.exe"
  Delete "$INSTDIR\servin-tui.exe"
  Delete "$INSTDIR\servin-gui.exe"
  Delete "$INSTDIR\servin.conf"
  Delete "$INSTDIR\uninstall.exe"
  Delete "$SMPROGRAMS\Servin\Servin.lnk"
  RMDir "$SMPROGRAMS\Servin"
  RMDir "$INSTDIR"
  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin"
SectionEnd
MINIMAL_NSI
                                
                                print_info "Attempting minimal NSIS compilation..."
                                if cmd.exe /c "makensis /V4 servin-minimal.nsi" > minimal-nsis.log 2>&1; then
                                    cat minimal-nsis.log
                                    if ls -la *installer*.exe 2>/dev/null | grep -q installer; then
                                        print_success "Minimal NSIS compilation succeeded"
                                    else
                                        print_error "Even minimal NSIS compilation failed"
                                        cat minimal-nsis.log
                                    fi
                                else
                                    print_error "Minimal NSIS compilation failed"
                                    cat minimal-nsis.log
                                fi
                            fi
                        else
                            print_error "Direct NSIS compilation failed"
                            cat direct-nsis.log
                        fi
                    fi
                else
                    print_error "NSIS installer build failed"
                    echo "NSIS error output:"
                    cat nsis-output.log
                fi
                echo "=== BATCH FILE OUTPUT END ==="
                
                # Additional verification after batch execution
                print_info "Post-execution verification:"
                echo "Files matching *installer*.exe:"
                ls -la *installer*.exe 2>/dev/null || echo "No installer files found"
                echo "Files matching servin-installer*:"
                ls -la servin-installer* 2>/dev/null || echo "No servin-installer files found"
                echo "Files matching Servin-Installer*:"
                ls -la Servin-Installer* 2>/dev/null || echo "No Servin-Installer files found"
                
                # Clean up temporary file
                rm -f "build-debug.bat"
                
            else
                print_error "NSIS (makensis) not found in PATH"
            fi
            
        else
            print_info "Running batch file directly..."
            chmod +x build-simple.bat
            if ./build-simple.bat; then
                print_success "Simple NSIS build completed successfully"
            else
                print_warning "Simple build failed, trying main build for enhanced features..."
                chmod +x build-installer.bat
                ./build-installer.bat || {
                    print_error "Both build methods failed"
                }
            fi
        fi
        
        # Verify installer was created and copy to dist
        local installer_created=false
        mkdir -p "$SCRIPT_DIR/dist"
        
        # Enhanced debugging for NSIS build failures
        print_info "Checking NSIS build results..."
        echo "Files in installers/windows directory:"
        ls -la
        
        if [[ -f "build.log" ]]; then
            print_info "NSIS build log found:"
            echo "=== NSIS BUILD LOG ==="
            cat build.log
            echo "=== END BUILD LOG ==="
        else
            print_warning "No NSIS build.log found"
        fi
        
        if [[ -f "Servin-Installer-1.0.0.exe" ]]; then
            print_success "Windows installer built successfully: Servin-Installer-1.0.0.exe"
            cp "Servin-Installer-1.0.0.exe" "$SCRIPT_DIR/dist/servin_${VERSION}_installer.exe"
            print_success "Windows installer copied to dist/servin_${VERSION}_installer.exe"
            installer_created=true
        elif [[ -f "servin-installer-1.0.0.exe" ]]; then
            print_success "Windows installer built successfully: servin-installer-1.0.0.exe"
            cp "servin-installer-1.0.0.exe" "$SCRIPT_DIR/dist/servin_${VERSION}_installer.exe"
            print_success "Windows installer copied to dist/servin_${VERSION}_installer.exe"
            installer_created=true
        fi
        
        if [[ "$installer_created" == "true" ]]; then
            ls -la *installer*.exe
        else
            print_error "Windows installer was not created"
            echo "Directory contents:"
            ls -la
            return 1
        fi
        
    elif command -v docker >/dev/null 2>&1; then
        print_info "Building NSIS installer using Docker..."
        
        # Create Dockerfile for NSIS build
        cat > "Dockerfile.nsis" << 'EOF'
FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    nsis \
    nsis-pluginapi \
    wine \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build
COPY . .

CMD ["makensis", "/DVERSION=1.0.0", "servin-installer.nsi"]
EOF
        
        cd "$windows_dir"
        docker build -f Dockerfile.nsis -t servin-nsis-builder .
        docker run --rm -v "$(pwd):/build" servin-nsis-builder
        
        if [[ -f "Servin-Installer-1.0.0.exe" ]]; then
            print_success "Windows installer built using Docker"
        else
            print_warning "Windows installer build failed"
        fi
    else
        print_warning "Cannot build Windows installer on this platform (NSIS not available)"
        print_info "Windows installer files prepared in: $windows_dir"
    fi
}

# Build Linux AppImage
build_linux_appimage() {
    print_platform "Building Linux AppImage"
    
    local linux_dir="$SCRIPT_DIR/installers/linux"
    
    # Copy Linux executables
    cp "$SCRIPT_DIR/build/linux-amd64/servin" "$linux_dir/"
    cp "$SCRIPT_DIR/build/linux-amd64/servin-tui" "$linux_dir/" 2>/dev/null || true
    
    if [[ "$PLATFORM" == "linux" ]]; then
        print_info "Building AppImage natively..."
        cd "$linux_dir"
        chmod +x build-appimage.sh
        ./build-appimage.sh
        
        # Copy AppImage to dist immediately after successful build
        mkdir -p "$SCRIPT_DIR/dist"
        if ls build/Servin-*.AppImage >/dev/null 2>&1; then
            local appimage_file=$(ls build/Servin-*.AppImage | head -1)
            cp "$appimage_file" "$SCRIPT_DIR/dist/servin_${VERSION}_installer.AppImage"
            print_success "Linux AppImage built and copied to dist/servin_${VERSION}_installer.AppImage"
            
            # Show size for verification
            local size=$(du -h "$SCRIPT_DIR/dist/servin_${VERSION}_installer.AppImage" | cut -f1)
            print_info "AppImage size: $size"
        else
            print_success "Linux AppImage built"
        fi
    elif command -v docker >/dev/null 2>&1; then
        print_info "Building AppImage using Docker..."
        
        # Create Dockerfile for AppImage build
        cat > "$linux_dir/Dockerfile.appimage" << 'EOF'
FROM ubuntu:20.04

RUN apt-get update && apt-get install -y \
    wget \
    curl \
    tar \
    file \
    desktop-file-utils \
    imagemagick \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build
COPY . .

CMD ["./build-appimage.sh"]
EOF
        
        cd "$linux_dir"
        docker build -f Dockerfile.appimage -t servin-appimage-builder .
        docker run --rm -v "$(pwd):/build" servin-appimage-builder
        
        if ls build/Servin-*.AppImage >/dev/null 2>&1; then
            print_success "Linux AppImage built using Docker"
        else
            print_warning "Linux AppImage build failed"
        fi
    else
        print_warning "Cannot build Linux AppImage on this platform"
        print_info "Linux AppImage files prepared in: $linux_dir"
    fi
}

# Build macOS package
build_macos_package() {
    print_platform "Building macOS Package"
    
    local macos_dir="$SCRIPT_DIR/installers/macos"
    
    # Copy macOS executables
    if [[ "$ARCH" == "arm64" ]]; then
        cp "$SCRIPT_DIR/build/darwin-arm64/servin" "$macos_dir/"
        cp "$SCRIPT_DIR/build/darwin-arm64/servin-tui" "$macos_dir/" 2>/dev/null || true
    else
        cp "$SCRIPT_DIR/build/darwin-amd64/servin" "$macos_dir/"
        cp "$SCRIPT_DIR/build/darwin-amd64/servin-tui" "$macos_dir/" 2>/dev/null || true
    fi
    
    if [[ "$PLATFORM" == "macos" ]]; then
        print_info "Building macOS package natively..."
        cd "$macos_dir"
        chmod +x build-package.sh
        ./build-package.sh
        print_success "macOS package built"
    else
        print_warning "Cannot build macOS package on this platform (requires macOS)"
        print_info "macOS package files prepared in: $macos_dir"
    fi
}

# Create unified distribution package
create_distribution() {
    print_header "Creating Distribution Package"
    
    local dist_dir="$SCRIPT_DIR/dist"
    local release_dir="$dist_dir/servin-$VERSION-$BUILD_DATE"
    
    # Preserve any immediate copies that were already placed in dist/
    local temp_dir="/tmp/servin-dist-backup-$$"
    if [[ -d "$dist_dir" ]]; then
        mkdir -p "$temp_dir"
        cp -r "$dist_dir"/* "$temp_dir/" 2>/dev/null || true
    fi
    
    rm -rf "$dist_dir"
    mkdir -p "$release_dir"/{windows,linux,macos,docs}
    mkdir -p "$dist_dir"  # Ensure root dist exists
    
    # Restore immediate copies for GitHub Actions verification
    if [[ -d "$temp_dir" ]]; then
        cp "$temp_dir"/* "$dist_dir/" 2>/dev/null || true
        rm -rf "$temp_dir"
    fi
    
    # Copy installers to structured release directory
    print_info "Collecting installer packages..."
    
    # Windows
    if [[ -f "$SCRIPT_DIR/installers/windows/Servin-Installer-$VERSION.exe" ]]; then
        cp "$SCRIPT_DIR/installers/windows/Servin-Installer-$VERSION.exe" "$release_dir/windows/"
        print_success "Windows installer included in release"
    elif [[ -f "$SCRIPT_DIR/installers/windows/servin-installer-$VERSION.exe" ]]; then
        cp "$SCRIPT_DIR/installers/windows/servin-installer-$VERSION.exe" "$release_dir/windows/"
        print_success "Windows installer included in release"
    fi
    
    # Linux
    if ls "$SCRIPT_DIR/installers/linux/build/Servin-"*.AppImage >/dev/null 2>&1; then
        cp "$SCRIPT_DIR/installers/linux/build/Servin-"*.AppImage "$release_dir/linux/"
        cp "$SCRIPT_DIR/installers/linux/build/install-servin-appimage.sh" "$release_dir/linux/" 2>/dev/null || true
        print_success "Linux AppImage included in release"
    fi
    
    # macOS
    if ls "$SCRIPT_DIR/installers/macos/build/Servin-"*.pkg >/dev/null 2>&1; then
        cp "$SCRIPT_DIR/installers/macos/build/Servin-"*.pkg "$release_dir/macos/"
        cp "$SCRIPT_DIR/installers/macos/build/Servin-"*.dmg "$release_dir/macos/" 2>/dev/null || true
        print_success "macOS package included in release"
    fi
    
    # Copy documentation
    print_info "Including documentation..."
    cp "$SCRIPT_DIR/README.md" "$release_dir/docs/"
    cp "$SCRIPT_DIR/LICENSE" "$release_dir/docs/" 2>/dev/null || true
    cp "$SCRIPT_DIR/INSTALL.md" "$release_dir/docs/" 2>/dev/null || true
    cp "$SCRIPT_DIR/VM_PREREQUISITES.md" "$release_dir/docs/" 2>/dev/null || true
    
    # Create installation guide
    cat > "$release_dir/INSTALLATION_GUIDE.md" << EOF
# Servin Container Runtime - Installation Guide

## Platform-Specific Installers

### Windows
- **Installer**: \`windows/Servin-Installer-$VERSION.exe\`
- **Requirements**: Windows 10/11 (64-bit)
- **Installation**: Run the installer as Administrator
- **Features**: 
  - Automatic VM provider installation (Hyper-V/VirtualBox/WSL2)
  - Desktop integration and Start Menu shortcuts
  - Automatic dependency management
  - Windows Service configuration

### Linux
- **AppImage**: \`linux/Servin-$VERSION-x86_64.AppImage\`
- **Requirements**: Linux with GLIBC 2.28+ (Ubuntu 20.04+, CentOS 8+)
- **Installation**: 
  \`\`\`bash
  chmod +x Servin-$VERSION-x86_64.AppImage
  ./linux/install-servin-appimage.sh
  \`\`\`
- **Features**:
  - Portable executable with all dependencies
  - Automatic QEMU/KVM installation
  - Desktop integration
  - No system-wide installation required

### macOS
- **Package**: \`macos/Servin-$VERSION-arm64.pkg\` (Apple Silicon) or \`macos/Servin-$VERSION-amd64.pkg\` (Intel)
- **Requirements**: macOS 10.15+ (Catalina or later)
- **Installation**: 
  \`\`\`bash
  sudo installer -pkg Servin-$VERSION-*.pkg -target /
  \`\`\`
- **Features**:
  - Native macOS app bundle
  - Automatic QEMU installation via Homebrew
  - Virtualization.framework support
  - LaunchDaemons integration

## Quick Start

After installation:

1. **Verify Installation**:
   \`\`\`bash
   servin version
   \`\`\`

2. **Initialize Servin**:
   \`\`\`bash
   servin init
   \`\`\`

3. **Start a Container**:
   \`\`\`bash
   servin run -it ubuntu:latest /bin/bash
   \`\`\`

4. **Launch GUI** (if available):
   - Windows: Start Menu → Servin Container Runtime
   - Linux: Applications → Servin or \`servin gui\`
   - macOS: Applications → Servin.app

## VM Prerequisites

All installers include automatic VM dependency installation:

- **Windows**: Hyper-V, VirtualBox, or WSL2
- **Linux**: QEMU/KVM with hardware acceleration
- **macOS**: QEMU with Virtualization.framework

## Troubleshooting

- **VM Issues**: Run \`servin vm status\` to check VM provider
- **Logs**: Check \`~/.servin/logs/\` for detailed logs
- **Support**: See docs/README.md for comprehensive documentation

## Build Information

- Version: $VERSION
- Build Date: $(date)
- Platforms: Windows (amd64), Linux (amd64/arm64), macOS (amd64/arm64)

EOF
    
    # Create checksums
    print_info "Generating checksums..."
    cd "$release_dir"
    find . -type f -name "*.exe" -o -name "*.AppImage" -o -name "*.pkg" -o -name "*.dmg" | xargs sha256sum > checksums.txt 2>/dev/null || true
    
    # Copy installers to root dist/ directory for GitHub Actions compatibility
    print_info "Creating GitHub Actions compatible installer copies..."
    
    # Windows installer
    print_info "Looking for Windows installer..."
    local windows_installer_found=false
    
    # Check multiple possible locations and names
    for installer_path in \
        "$release_dir/windows/Servin-Installer-$VERSION.exe" \
        "$release_dir/windows/servin-installer-$VERSION.exe" \
        "$SCRIPT_DIR/installers/windows/Servin-Installer-$VERSION.exe" \
        "$SCRIPT_DIR/installers/windows/servin-installer-$VERSION.exe"; do
        
        if [[ -f "$installer_path" ]]; then
            cp "$installer_path" "$dist_dir/servin_${VERSION}_installer.exe"
            print_success "Windows installer copied from $installer_path to dist/servin_${VERSION}_installer.exe"
            windows_installer_found=true
            break
        fi
    done
    
    if [[ "$windows_installer_found" == "false" ]]; then
        print_warning "Windows installer not found in expected locations"
        print_info "Searched locations:"
        print_info "  - $release_dir/windows/Servin-Installer-$VERSION.exe"
        print_info "  - $release_dir/windows/servin-installer-$VERSION.exe"
        print_info "  - $SCRIPT_DIR/installers/windows/Servin-Installer-$VERSION.exe"
        print_info "  - $SCRIPT_DIR/installers/windows/servin-installer-$VERSION.exe"
    fi
    
    # Linux AppImage
    print_info "Looking for Linux AppImage..."
    local appimage_found=false
    
    # Check multiple possible locations
    for appimage_path in \
        "$release_dir/linux/Servin-"*.AppImage \
        "$SCRIPT_DIR/installers/linux/build/Servin-"*.AppImage; do
        
        if ls $appimage_path >/dev/null 2>&1; then
            APPIMAGE_FILE=$(ls $appimage_path | head -1)
            APPIMAGE_NAME=$(basename "$APPIMAGE_FILE")
            cp "$APPIMAGE_FILE" "$dist_dir/servin_${VERSION}_installer.AppImage"
            print_success "Linux AppImage copied from $APPIMAGE_FILE to dist/servin_${VERSION}_installer.AppImage"
            appimage_found=true
            break
        fi
    done
    
    if [[ "$appimage_found" == "false" ]]; then
        print_warning "Linux AppImage not found in expected locations"
        print_info "Searched locations:"
        print_info "  - $release_dir/linux/Servin-*.AppImage"
        print_info "  - $SCRIPT_DIR/installers/linux/build/Servin-*.AppImage"
    fi
    
    # macOS PKG
    if ls "$release_dir/macos/Servin-"*.pkg >/dev/null 2>&1; then
        PKG_FILE=$(ls "$release_dir/macos/Servin-"*.pkg | head -1)
        PKG_NAME=$(basename "$PKG_FILE")
        cp "$PKG_FILE" "$dist_dir/servin_${VERSION}_installer.pkg"
        print_info "macOS PKG copied to dist/servin_${VERSION}_installer.pkg"
    fi
    
    # Create archive
    print_info "Creating distribution archive..."
    cd "$dist_dir"
    tar -czf "servin-$VERSION-$BUILD_DATE-complete.tar.gz" "servin-$VERSION-$BUILD_DATE"
    
    print_success "Distribution package created"
    
    # Show summary
    print_header "Distribution Summary"
    echo
    print_info "Distribution directory: $release_dir"
    print_info "Archive: $dist_dir/servin-$VERSION-$BUILD_DATE-complete.tar.gz"
    echo
    print_info "Contents:"
    find "$release_dir" -type f | sort | sed 's/^/  • /'
    echo
    
    if [[ -f "$release_dir/checksums.txt" ]]; then
        print_info "Checksums (SHA256):"
        cat "$release_dir/checksums.txt" | sed 's/^/  /'
    fi
}

# Show build summary
show_summary() {
    print_header "Build Summary"
    echo
    print_success "Cross-platform package build completed!"
    echo
    
    local built_packages=()
    
    # Check what was built
    if [[ -f "$SCRIPT_DIR/windows/Servin-Installer-$VERSION.exe" ]]; then
        built_packages+=("✓ Windows NSIS Installer")
    else
        built_packages+=("⚠ Windows NSIS Installer (files prepared)")
    fi
    
    if ls "$SCRIPT_DIR/linux/build/Servin-"*.AppImage >/dev/null 2>&1; then
        built_packages+=("✓ Linux AppImage")
    else
        built_packages+=("⚠ Linux AppImage (files prepared)")
    fi
    
    if ls "$SCRIPT_DIR/macos/build/Servin-"*.pkg >/dev/null 2>&1; then
        built_packages+=("✓ macOS Package")
    else
        built_packages+=("⚠ macOS Package (files prepared)")
    fi
    
    print_info "Built packages:"
    for package in "${built_packages[@]}"; do
        echo "  $package"
    done
    echo
    
    print_info "Current platform: $PLATFORM/$ARCH"
    print_info "To build missing packages, run this script on the target platform"
    echo
    
    if [[ -d "$SCRIPT_DIR/dist" ]]; then
        print_info "Distribution package available in: $SCRIPT_DIR/dist/"
    fi
}

# Main execution
main() {
    print_banner
    
    # Parse command line arguments
    local build_all=true
    local build_windows=false
    local build_linux=false
    local build_macos=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --windows) build_windows=true; build_all=false; shift ;;
            --linux) build_linux=true; build_all=false; shift ;;
            --macos) build_macos=true; build_all=false; shift ;;
            --help) 
                echo "Usage: $0 [--windows] [--linux] [--macos]"
                echo "  --windows   Build only Windows installer"
                echo "  --linux     Build only Linux AppImage"
                echo "  --macos     Build only macOS package"
                echo "  (no args)   Build all platforms"
                exit 0
                ;;
            *) 
                print_error "Unknown option: $1"
                exit 1
                ;;
        esac
    done
    
    detect_platform
    check_prerequisites
    
    # Build executables for selected platforms only
    if [[ "$build_all" == "true" ]]; then
        build_executables
    else
        local target_platforms=()
        if [[ "$build_windows" == "true" ]]; then
            target_platforms+=("windows/amd64")
        fi
        if [[ "$build_linux" == "true" ]]; then
            target_platforms+=("linux/amd64" "linux/arm64")
        fi
        if [[ "$build_macos" == "true" ]]; then
            target_platforms+=("darwin/amd64" "darwin/arm64")
        fi
        build_executables "${target_platforms[@]}"
    fi
    
    if [[ "$build_all" == "true" || "$build_windows" == "true" ]]; then
        build_windows_installer
    fi
    
    if [[ "$build_all" == "true" || "$build_linux" == "true" ]]; then
        build_linux_appimage
    fi
    
    if [[ "$build_all" == "true" || "$build_macos" == "true" ]]; then
        build_macos_package
    fi
    
    create_distribution
    show_summary
}

# Run main function
main "$@"