#!/bin/bash
# Servin Container Runtime - User Installation Script
# Installs Servin to user directories without requiring sudo

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
    *)      echo -e "${RED}❌ Unsupported platform: $OS${NC}"; exit 1 ;;
esac

case $ARCH in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo -e "${YELLOW}⚠️  Warning: Unsupported architecture $ARCH, defaulting to amd64${NC}"; ARCH="amd64" ;;
esac

# Installation directories
HOME_DIR="$HOME"
INSTALL_DIR="$HOME/.local/bin"
DATA_DIR="$HOME/.local/share/servin"
CONFIG_DIR="$HOME/.config/servin"

echo -e "${BLUE}🚀 Servin Container Runtime - User Installer${NC}"
echo -e "${BLUE}Platform: $PLATFORM-$ARCH${NC}"
echo ""

# Create directories
echo -e "${YELLOW}📁 Creating directories...${NC}"
mkdir -p "$INSTALL_DIR"
mkdir -p "$DATA_DIR"/{volumes,images,containers,logs}
mkdir -p "$CONFIG_DIR"

# Get script directory (where this installer is located)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Install binaries
echo -e "${YELLOW}📦 Installing binaries...${NC}"

# Core binary
if [ -f "$SCRIPT_DIR/servin" ]; then
    cp "$SCRIPT_DIR/servin" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/servin"
    echo -e "${GREEN}✓ Installed: servin${NC}"
else
    echo -e "${RED}❌ Error: servin binary not found in installer package${NC}"
    exit 1
fi

# Desktop binary (optional)
if [ -f "$SCRIPT_DIR/servin-tui" ]; then
    cp "$SCRIPT_DIR/servin-tui" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/servin-tui"
    echo -e "${GREEN}✓ Installed: servin-tui${NC}"
else
    echo -e "${YELLOW}⚠️  Warning: servin-tui not found, skipping${NC}"
fi

# Create configuration file
echo -e "${YELLOW}⚙️  Creating configuration...${NC}"
cat > "$CONFIG_DIR/servin.conf" << EOF
# Servin Configuration File
data_dir=$DATA_DIR
log_level=info
log_file=$DATA_DIR/logs/servin.log
runtime=native
bridge_name=servin0
gui_theme=auto
enable_notifications=true
EOF

echo -e "${GREEN}✓ Created configuration file${NC}"

# Setup PATH
echo -e "${YELLOW}📝 Setting up PATH...${NC}"

PATH_LINE="export PATH=\"$INSTALL_DIR:\$PATH\"  # Added by Servin installer"

# Shell configuration files to update
SHELL_CONFIGS=()
[ -f "$HOME/.bashrc" ] && SHELL_CONFIGS+=("$HOME/.bashrc")
[ -f "$HOME/.zshrc" ] && SHELL_CONFIGS+=("$HOME/.zshrc")
[ -f "$HOME/.profile" ] && SHELL_CONFIGS+=("$HOME/.profile")

# If no shell configs exist, create .profile
if [ ${#SHELL_CONFIGS[@]} -eq 0 ]; then
    SHELL_CONFIGS+=("$HOME/.profile")
fi

# On macOS, prefer .zshrc
if [ "$PLATFORM" = "mac" ]; then
    SHELL_CONFIGS=("$HOME/.zshrc")
fi

for config_file in "${SHELL_CONFIGS[@]}"; do
    if [ ! -f "$config_file" ] || ! grep -q "Added by Servin installer" "$config_file"; then
        echo "" >> "$config_file"
        echo "# Servin Container Runtime" >> "$config_file"
        echo "$PATH_LINE" >> "$config_file"
        echo -e "${GREEN}✓ Updated $(basename "$config_file")${NC}"
    else
        echo -e "${BLUE}ℹ️  $(basename "$config_file") already configured${NC}"
    fi
done

# Create desktop integration (Linux only)
if [ "$PLATFORM" = "linux" ] && [ -f "$INSTALL_DIR/servin-tui" ]; then
    echo -e "${YELLOW}🖥️  Creating desktop integration...${NC}"
    
    DESKTOP_DIR="$HOME/.local/share/applications"
    ICON_DIR="$HOME/.local/share/icons/hicolor"
    mkdir -p "$DESKTOP_DIR"
    
    # Install icon if available
    ICON_NAME="servin-tui"
    ICON_REFERENCE="application-x-executable"  # fallback
    
    # Try to install different icon sizes
    for size in 16 32 48 64 128 256; do
        if [ -f "$SCRIPT_DIR/servin-icon-${size}.png" ]; then
            SIZE_DIR="$ICON_DIR/${size}x${size}/apps"
            mkdir -p "$SIZE_DIR"
            cp "$SCRIPT_DIR/servin-icon-${size}.png" "$SIZE_DIR/${ICON_NAME}.png"
            ICON_REFERENCE="$ICON_NAME"
            echo -e "${GREEN}✓ Installed ${size}x${size} icon${NC}"
        fi
    done
    
    # Create desktop file
    cat > "$DESKTOP_DIR/servin-tui.desktop" << EOF
[Desktop Entry]
Version=1.0
Type=Application
Name=Servin Desktop
Comment=Container Management with Servin
Exec=$INSTALL_DIR/servin-tui
Icon=$ICON_REFERENCE
Terminal=false
Categories=Development;System;
Keywords=container;docker;runtime;
StartupNotify=true
EOF
    
    chmod +x "$DESKTOP_DIR/servin-tui.desktop"
    echo -e "${GREEN}✓ Created desktop entry${NC}"
fi

# Create macOS app bundle (macOS only)
if [ "$PLATFORM" = "mac" ] && [ -f "$INSTALL_DIR/servin-tui" ]; then
    echo -e "${YELLOW}🍎 Creating macOS application bundle...${NC}"
    
    APP_PATH="$HOME/Applications/Servin Desktop.app"
    CONTENTS_PATH="$APP_PATH/Contents"
    MACOS_PATH="$CONTENTS_PATH/MacOS"
    RESOURCES_PATH="$CONTENTS_PATH/Resources"
    
    mkdir -p "$MACOS_PATH"
    mkdir -p "$RESOURCES_PATH"
    
    # Copy executable
    cp "$INSTALL_DIR/servin-tui" "$MACOS_PATH/Servin Desktop"
    chmod +x "$MACOS_PATH/Servin Desktop"
    
    # Copy icon if available
    ICON_COPIED=false
    for icon_file in "servin.icns" "servin-icon-512.png" "servin-icon-256.png"; do
        if [ -f "$SCRIPT_DIR/$icon_file" ]; then
            if [[ "$icon_file" == *.icns ]]; then
                cp "$SCRIPT_DIR/$icon_file" "$RESOURCES_PATH/servin.icns"
                ICON_NAME="servin"
            else
                cp "$SCRIPT_DIR/$icon_file" "$RESOURCES_PATH/servin.png"
                ICON_NAME="servin"
            fi
            ICON_COPIED=true
            echo -e "${GREEN}✓ Installed app icon: $icon_file${NC}"
            break
        fi
    done
    
    # Create Info.plist
    cat > "$CONTENTS_PATH/Info.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>Servin Desktop</string>
    <key>CFBundleIdentifier</key>
    <string>com.servin.desktop</string>
    <key>CFBundleName</key>
    <string>Servin Desktop</string>
    <key>CFBundleDisplayName</key>
    <string>Servin Desktop</string>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.0</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.12</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>LSApplicationCategoryType</key>
    <string>public.app-category.developer-tools</string>$(if [ "$ICON_COPIED" = true ]; then echo "
    <key>CFBundleIconFile</key>
    <string>$ICON_NAME</string>"; fi)
</dict>
</plist>
EOF
    
    echo -e "${GREEN}✓ Created application bundle${NC}"
fi

echo ""
echo -e "${GREEN}🎉 Installation completed successfully!${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo -e "${BLUE}1. Open a new terminal window (to load PATH changes)${NC}"
echo -e "${BLUE}2. Run: servin --help${NC}"
echo -e "${BLUE}3. Try: servin image pull alpine${NC}"
echo -e "${BLUE}4. Try: servin run alpine echo 'Hello World'${NC}"
echo ""

if [ "$PLATFORM" = "mac" ]; then
    echo -e "${BLUE}The desktop app is available in your Applications folder${NC}"
elif [ "$PLATFORM" = "linux" ]; then
    echo -e "${BLUE}The desktop app is available in your application menu${NC}"
fi

echo ""
echo -e "${YELLOW}Installation locations:${NC}"
echo -e "${BLUE}• Binaries: $INSTALL_DIR${NC}"
echo -e "${BLUE}• Data: $DATA_DIR${NC}"
echo -e "${BLUE}• Config: $CONFIG_DIR${NC}"
echo ""
echo -e "${YELLOW}To uninstall, simply remove these directories.${NC}"