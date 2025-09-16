#!/bin/bash
# Simple build script for Servin - generates binaries in build folder organized by platform

set -e

# Configuration
BUILD_DIR="build"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date +%Y-%m-%dT%H:%M:%S)
LDFLAGS="-ldflags \"-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME\""

# Platform detection
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Normalize architecture names
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    i386|i686)
        ARCH="386"
        ;;
esac

# Set platform-specific variables
case $OS in
    darwin)
        PLATFORM="darwin"
        EXT=""
        ;;
    linux)
        PLATFORM="linux"
        EXT=""
        ;;
    mingw*|cygwin*|msys*)
        PLATFORM="windows"
        EXT=".exe"
        ;;
    *)
        PLATFORM="unknown"
        EXT=""
        ;;
esac

PLATFORM_DIR="$BUILD_DIR/$PLATFORM-$ARCH"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_header() {
    echo -e "${CYAN}================================================${NC}"
    echo -e "${CYAN}   Servin Local Build Script${NC}"
    echo -e "${CYAN}================================================${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_info() {
    echo -e "${BLUE}→ $1${NC}"
}

# Start build process
print_header

print_info "Building for platform: $PLATFORM-$ARCH"

# Create/clean build directory
print_info "Creating build directory structure..."
rm -rf "$BUILD_DIR"
mkdir -p "$PLATFORM_DIR"

# Build main servin binary
print_info "Building servin runtime binary..."
eval "go build $LDFLAGS -o $PLATFORM_DIR/servin$EXT ."
print_success "Built servin binary"

# Build TUI desktop binary
print_info "Building servin-tui TUI binary..."
eval "go build $LDFLAGS -o $PLATFORM_DIR/servin-tui$EXT ./cmd/servin-tui"
print_success "Built servin-tui TUI binary"

# Build GUI binary (optional)
print_info "Building servin-gui binary..."
if eval "go build $LDFLAGS -o $PLATFORM_DIR/servin-gui$EXT ./cmd/servin-gui" 2>/dev/null; then
    print_success "Built servin-gui binary"
else
    echo -e "${YELLOW}⚠ GUI binary build skipped (dependencies may be missing)${NC}"
fi

# Display results
echo ""
print_info "Build completed! Generated files:"
ls -la "$PLATFORM_DIR"

# Create platform-specific README
cat > "$PLATFORM_DIR/README.md" << EOF
# Servin Binaries for $PLATFORM-$ARCH

This directory contains the compiled binaries for the Servin Container Runtime built for $PLATFORM-$ARCH.

## Generated Binaries

### \`servin$EXT\`
- **Description**: Main Servin container runtime binary
- **Usage**: \`./servin$EXT [command] [flags]\`
- **Help**: \`./servin$EXT --help\`

### \`servin-tui$EXT\`
- **Description**: Terminal User Interface (TUI) for Servin  
- **Usage**: \`./servin-tui$EXT\`
- **Features**: Interactive terminal-based management interface

### \`servin-gui$EXT\`
- **Description**: Graphical User Interface (GUI) for Servin
- **Usage**: \`./servin-gui$EXT\`
- **Features**: Native desktop application using Fyne framework

## Build Information

- **Platform**: $PLATFORM-$ARCH
- **Version**: $VERSION
- **Build Time**: $BUILD_TIME
- **Built with**: Go \$(go version | awk '{print \$3}')

## Quick Start

1. **Container Management**:
   \`\`\`bash
   ./servin$EXT run ubuntu:latest
   ./servin$EXT ls
   ./servin$EXT stop <container-id>
   \`\`\`

2. **GUI Interface**: \`./servin-gui$EXT\`
3. **TUI Interface**: \`./servin-tui$EXT\`

## Rebuilding

To rebuild: \`../../build-local.sh\` or \`make build-local\` from project root
EOF

# Create main build directory README
cat > "$BUILD_DIR/README.md" << EOF
# Servin Build Directory

This directory contains platform-specific builds of the Servin Container Runtime.

## Directory Structure

\`\`\`
build/
└── $PLATFORM-$ARCH/          # Current platform binaries
    ├── README.md             # Platform-specific documentation
    ├── servin$EXT            # Main runtime binary
    ├── servin-tui$EXT    # TUI binary
    └── servin-gui$EXT        # GUI binary (if available)
\`\`\`

## Available Platforms

- **$PLATFORM-$ARCH**: Current build ($(date))

## Usage

Navigate to your platform directory and run the binaries:

\`\`\`bash
cd $PLATFORM-$ARCH
./servin$EXT --help
./servin-tui$EXT
./servin-gui$EXT
\`\`\`

## Cross-Platform Building

To build for other platforms, you can use:

\`\`\`bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o build/linux-amd64/servin .

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o build/windows-amd64/servin.exe .

# macOS ARM64
GOOS=darwin GOARCH=arm64 go build -o build/darwin-arm64/servin .
\`\`\`

## Rebuilding

To rebuild: \`../build-local.sh\` or \`make build-local\` from project root
EOF

echo ""
print_success "All binaries are available in the '$PLATFORM_DIR' directory"
echo ""
print_info "Usage:"
echo "  Runtime: ./$PLATFORM_DIR/servin$EXT --help"
echo "  TUI:     ./$PLATFORM_DIR/servin-tui$EXT"
echo "  GUI:     ./$PLATFORM_DIR/servin-gui$EXT"
