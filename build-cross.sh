u#!/bin/bash
# Cross-platform build script for Servin - generates binaries for multiple platforms

set -e

# Configuration
BUILD_DIR="build"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date +%Y-%m-%dT%H:%M:%S)
LDFLAGS="-ldflags \"-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME\""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_header() {
    echo -e "${CYAN}================================================${NC}"
    echo -e "${CYAN}   Servin Cross-Platform Build Script${NC}"
    echo -e "${CYAN}================================================${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${BLUE}→ $1${NC}"
}

# Define target platforms
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

# Parse command line arguments
BUILD_ALL=false
BUILD_CURRENT=false
SPECIFIC_PLATFORMS=()

while [[ $# -gt 0 ]]; do
    case $1 in
        --all)
            BUILD_ALL=true
            shift
            ;;
        --current)
            BUILD_CURRENT=true
            shift
            ;;
        --platform)
            SPECIFIC_PLATFORMS+=("$2")
            shift 2
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --all                 Build for all supported platforms"
            echo "  --current            Build for current platform only"
            echo "  --platform PLATFORM  Build for specific platform (format: os/arch)"
            echo "  --help, -h           Show this help message"
            echo ""
            echo "Supported platforms:"
            for platform in "${PLATFORMS[@]}"; do
                echo "  - $platform"
            done
            echo ""
            echo "Examples:"
            echo "  $0 --current                    # Build for current platform"
            echo "  $0 --all                        # Build for all platforms"
            echo "  $0 --platform linux/amd64      # Build for Linux AMD64"
            echo "  $0 --platform darwin/arm64 --platform windows/amd64  # Multiple platforms"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Default to current platform if no arguments
if [[ $BUILD_ALL == false && $BUILD_CURRENT == false && ${#SPECIFIC_PLATFORMS[@]} -eq 0 ]]; then
    BUILD_CURRENT=true
fi

build_platform() {
    local goos=$1
    local goarch=$2
    local platform_name="$goos-$goarch"
    local platform_dir="$BUILD_DIR/$platform_name"
    
    # Set file extension for Windows
    local ext=""
    if [[ "$goos" == "windows" ]]; then
        ext=".exe"
    fi
    
    print_info "Building for $platform_name..."
    
    # Create platform directory
    mkdir -p "$platform_dir"
    
    # Build main servin binary
    if GOOS=$goos GOARCH=$goarch CGO_ENABLED=0 go build -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME" -o "$platform_dir/servin$ext" . 2>/dev/null; then
        print_success "Built servin binary for $platform_name"
    else
        print_error "Failed to build servin binary for $platform_name"
        return 1
    fi
    
    # Build TUI desktop binary
    if GOOS=$goos GOARCH=$goarch CGO_ENABLED=0 go build -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME" -o "$platform_dir/servin-desktop$ext" ./cmd/servin-desktop 2>/dev/null; then
        print_success "Built servin-desktop TUI binary for $platform_name"
    else
        print_error "Failed to build servin-desktop TUI binary for $platform_name"
    fi
    
    # Build GUI binary (requires CGO, might fail for cross-compilation)
    if [[ "$goos" == "$(go env GOOS)" && "$goarch" == "$(go env GOARCH)" ]]; then
        # Native build - try with CGO
        if CGO_ENABLED=1 go build -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME" -o "$platform_dir/servin-gui$ext" ./cmd/servin-gui 2>/dev/null; then
            print_success "Built servin-gui binary for $platform_name"
        else
            print_warning "GUI binary build skipped for $platform_name (CGO dependencies missing)"
        fi
    else
        # Cross-compilation - try without CGO first, then with if it fails
        if GOOS=$goos GOARCH=$goarch CGO_ENABLED=0 go build -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME" -o "$platform_dir/servin-gui$ext" ./cmd/servin-gui 2>/dev/null; then
            print_success "Built servin-gui binary for $platform_name (no CGO)"
        else
            print_warning "GUI binary build skipped for $platform_name (cross-compilation limitations)"
        fi
    fi
    
    # Create platform-specific README
    cat > "$platform_dir/README.md" << EOF
# Servin Binaries for $platform_name

This directory contains the compiled binaries for the Servin Container Runtime built for $platform_name.

## Generated Binaries

### \`servin$ext\`
- **Description**: Main Servin container runtime binary
- **Usage**: \`./servin$ext [command] [flags]\`
- **Help**: \`./servin$ext --help\`

### \`servin-desktop$ext\`
- **Description**: Terminal User Interface (TUI) for Servin  
- **Usage**: \`./servin-desktop$ext\`
- **Features**: Interactive terminal-based management interface

### \`servin-gui$ext\`
- **Description**: Graphical User Interface (GUI) for Servin
- **Usage**: \`./servin-gui$ext\`
- **Features**: Native desktop application using Fyne framework

## Build Information

- **Platform**: $platform_name
- **Version**: $VERSION
- **Build Time**: $BUILD_TIME
- **Built with**: Go $(go version | awk '{print $3}')

## Quick Start

1. **Container Management**:
   \`\`\`bash
   ./servin$ext run ubuntu:latest
   ./servin$ext ls
   ./servin$ext stop <container-id>
   \`\`\`

2. **GUI Interface**: \`./servin-gui$ext\`
3. **TUI Interface**: \`./servin-desktop$ext\`

## Cross-Platform Note

This binary was built for $platform_name. Make sure you're running it on a compatible system.
EOF
    
    echo ""
}

# Start build process
print_header

# Create/clean build directory
print_info "Creating build directory structure..."
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

# Determine which platforms to build
TARGETS=()

if [[ $BUILD_CURRENT == true ]]; then
    current_os=$(go env GOOS)
    current_arch=$(go env GOARCH)
    TARGETS+=("$current_os/$current_arch")
    print_info "Building for current platform: $current_os/$current_arch"
fi

if [[ $BUILD_ALL == true ]]; then
    TARGETS=("${PLATFORMS[@]}")
    print_info "Building for all supported platforms"
fi

if [[ ${#SPECIFIC_PLATFORMS[@]} -gt 0 ]]; then
    TARGETS+=("${SPECIFIC_PLATFORMS[@]}")
    print_info "Building for specified platforms: ${SPECIFIC_PLATFORMS[*]}"
fi

# Build for each target platform
echo ""
for target in "${TARGETS[@]}"; do
    IFS='/' read -r goos goarch <<< "$target"
    build_platform "$goos" "$goarch"
done

# Create main build directory README
cat > "$BUILD_DIR/README.md" << EOF
# Servin Build Directory

This directory contains platform-specific builds of the Servin Container Runtime.

## Directory Structure

\`\`\`
build/
$(for target in "${TARGETS[@]}"; do
    platform_name="${target//\//-}"
    echo "├── $platform_name/               # Binaries for $target"
done)
\`\`\`

## Available Platforms

$(for target in "${TARGETS[@]}"; do
    platform_name="${target//\//-}"
    echo "- **$platform_name**: Built on $(date)"
done)

## Usage

Navigate to your platform directory and run the binaries:

\`\`\`bash
# Example for your current platform
cd \$(go env GOOS)-\$(go env GOARCH)
./servin --help
./servin-desktop
./servin-gui
\`\`\`

## Rebuilding

To rebuild:
- Current platform: \`../build-cross.sh --current\`
- All platforms: \`../build-cross.sh --all\`
- Specific platform: \`../build-cross.sh --platform linux/amd64\`

Or use the local build script: \`../build-local.sh\` (current platform only)
EOF

echo ""
print_success "Build completed!"
print_info "Generated directories:"
ls -la "$BUILD_DIR"

echo ""
print_info "To use the binaries:"
for target in "${TARGETS[@]}"; do
    platform_name="${target//\//-}"
    ext=""
    if [[ "$target" == *"windows"* ]]; then
        ext=".exe"
    fi
    echo "  $platform_name: ./build/$platform_name/servin$ext --help"
done
