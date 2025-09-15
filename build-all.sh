#!/bin/bash
# Build all Servin components for distribution

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Detect platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $OS in
    darwin) PLATFORM="mac" ;;
    linux)  PLATFORM="linux" ;;
    *)      PLATFORM="other" ;;
esac

case $ARCH in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) ARCH="amd64" ;;
esac

BUILD_DIR="build/${PLATFORM}"
DIST_DIR="dist/${PLATFORM}"

echo -e "${BLUE}🚀 Building Servin for ${PLATFORM}-${ARCH}${NC}"

# Create build directories
mkdir -p "$BUILD_DIR"
mkdir -p "$DIST_DIR"

# Build main CLI
echo -e "${YELLOW}📦 Building CLI (servin)${NC}"
go build -ldflags="-s -w" -o "${BUILD_DIR}/servin" ./main.go

# Build GUI (temporarily disabled due to compatibility issues)
echo -e "${YELLOW}📦 Skipping GUI (compatibility issues)${NC}"
# go build -ldflags="-s -w" -o "${BUILD_DIR}/servin-gui" ./cmd/servin-gui/main.go

# Build Desktop
echo -e "${YELLOW}📦 Building Desktop (servin-desktop)${NC}"
go build -ldflags="-s -w" -o "${BUILD_DIR}/servin-desktop" ./cmd/servin-desktop/main.go

# Create distribution package
echo -e "${YELLOW}📁 Creating distribution package${NC}"
cp "${BUILD_DIR}"/* "${DIST_DIR}/"

# Copy icons
echo -e "${YELLOW}🎨 Copying icons${NC}"
cp icons/servin-icon-*.png "${DIST_DIR}/" 2>/dev/null || true
cp icons/servin.ico "${DIST_DIR}/" 2>/dev/null || true
cp icons/servin.icns "${DIST_DIR}/" 2>/dev/null || true

# Copy installer
case $PLATFORM in
    mac)
        cp installers/macos/servin-installer.py "${DIST_DIR}/"
        cp installers/macos/install.sh "${DIST_DIR}/"
        ;;
    linux)
        cp installers/linux/servin-installer.py "${DIST_DIR}/"
        cp installers/linux/install.sh "${DIST_DIR}/"
        ;;
esac

# Copy documentation
cp README.md "${DIST_DIR}/"
cp LICENSE "${DIST_DIR}/"

# Make all executables
chmod +x "${DIST_DIR}"/servin*

echo -e "${GREEN}✅ Build completed successfully!${NC}"
echo -e "${BLUE}📁 Distribution files are in: ${DIST_DIR}${NC}"
echo ""
echo -e "${YELLOW}Built components:${NC}"
ls -la "${DIST_DIR}"/servin*

# Show sizes
echo ""
echo -e "${YELLOW}Component sizes:${NC}"
du -h "${DIST_DIR}"/servin*