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

build_binaries() {
    local platform=$1
    local arch=$2
    local ext=$3
    
    print_info "Building for $platform/$arch..."
    
    local output_dir="$BUILD_DIR/$platform-$arch"
    mkdir -p "$output_dir"
    
    # Build main servin binary
    GOOS=$platform GOARCH=$arch CGO_ENABLED=0 go build -ldflags="-w -s -X main.version=$VERSION" -o "$output_dir/servin$ext" .
    
    # Build TUI binary (servin-desktop)
    GOOS=$platform GOARCH=$arch CGO_ENABLED=0 go build -ldflags="-w -s -X main.version=$VERSION" -o "$output_dir/servin-desktop$ext" ./cmd/servin-desktop
    
    # Build GUI binary (only for platforms that support it)
    if [[ "$platform" != "linux" ]] || command -v pkg-config >/dev/null 2>&1; then
        if [[ "$platform" == "windows" ]]; then
            GOOS=$platform GOARCH=$arch CGO_ENABLED=1 go build -ldflags="-w -s -X main.version=$VERSION" -o "$output_dir/servin-gui$ext" ./cmd/servin-gui
        elif [[ "$platform" == "darwin" ]]; then
            GOOS=$platform GOARCH=$arch CGO_ENABLED=1 go build -ldflags="-w -s -X main.version=$VERSION" -o "$output_dir/servin-gui$ext" ./cmd/servin-gui
        elif [[ "$platform" == "linux" ]]; then
            # Linux GUI build requires X11 dependencies
            if pkg-config --exists x11; then
                GOOS=$platform GOARCH=$arch CGO_ENABLED=1 go build -ldflags="-w -s -X main.version=$VERSION" -o "$output_dir/servin-gui$ext" ./cmd/servin-gui
            else
                print_warning "  Skipping GUI build - X11 development libraries not found"
            fi
        fi
    fi
    
    print_success "  Built binaries for $platform/$arch"
}

create_windows_package() {
    print_info "Creating Windows package..."
    
    local platform_dir="$BUILD_DIR/windows-amd64"
    local package_dir="$DIST_DIR/servin-windows-$VERSION"
    
    mkdir -p "$package_dir"
    
    # Copy binaries
    cp "$platform_dir"/*.exe "$package_dir/"
    
    # Copy installer
    cp "$INSTALLER_DIR/windows/install.ps1" "$package_dir/"
    
    # Create README
    cat > "$package_dir/README.txt" << EOF
Servin Container Runtime for Windows
Version: $VERSION

Installation:
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
- GUI: servin-gui.exe or use desktop shortcut
- CLI: servin.exe --help
- Service: Start-Service ServinRuntime

Uninstallation:
Run the uninstaller from the Start Menu or:
C:\Program Files\Servin\uninstall.ps1

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
    
    # Copy installer
    cp "$INSTALLER_DIR/linux/install.sh" "$package_dir/"
    chmod +x "$package_dir/install.sh"
    
    # Create README
    cat > "$package_dir/README.md" << EOF
# Servin Container Runtime for Linux
Version: $VERSION

## Installation
\`\`\`bash
sudo ./install.sh
\`\`\`

This will:
- Install Servin to /usr/local/bin
- Create a system user 'servin'
- Set up systemd service (or SysV init script)
- Create configuration in /etc/servin
- Set up data directories in /var/lib/servin

## Usage
- GUI: \`servin-gui\` (if installed)
- CLI: \`servin --help\`
- Service: \`sudo systemctl start servin\`

## Uninstallation
\`\`\`bash
sudo /usr/local/bin/servin-uninstall
\`\`\`

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
    
    # Copy installer
    cp "$INSTALLER_DIR/macos/install.sh" "$package_dir/"
    chmod +x "$package_dir/install.sh"
    
    # Create README
    cat > "$package_dir/README.md" << EOF
# Servin Container Runtime for macOS
Version: $VERSION

## Installation
\`\`\`bash
sudo ./install.sh
\`\`\`

This will:
- Install Servin to /usr/local/bin
- Create a system user '_servin'
- Set up launchd service
- Create configuration in /usr/local/etc/servin
- Set up data directories in /usr/local/var/lib/servin
- Create application bundle (if GUI available)

## Usage
- GUI: Open "Servin GUI" from Applications or run \`servin-gui\`
- CLI: \`servin --help\`
- Service: Starts automatically (launchd)

## Uninstallation
\`\`\`bash
sudo /usr/local/bin/servin-uninstall
\`\`\`

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
    echo "Windows: Extract ZIP, run install.ps1 as Administrator"
    echo "Linux:   Extract tar.gz, run sudo ./install.sh"
    echo "macOS:   Extract tar.gz, run sudo ./install.sh"
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
    *)
        main
        ;;
esac
