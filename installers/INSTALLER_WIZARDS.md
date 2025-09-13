# Servin Container Runtime - Installation Wizards

Professional installer wizards for easy Servin Container Runtime deployment across all platforms.

## 🎯 Overview

We've created user-friendly installation wizards that provide a guided installation experience instead of command-line scripts:

### **Windows Installer Wizard**
- **NSIS-based installer** - Professional Windows installer
- **GUI wizard interface** - Step-by-step guided installation
- **Component selection** - Choose what to install
- **Service integration** - Automatic Windows Service setup
- **Uninstaller included** - Clean removal via Control Panel

### **Linux Installation Wizard**
- **Python/Tkinter GUI** - Native Linux desktop installer
- **Cross-distribution support** - Works on Ubuntu, CentOS, Debian, etc.
- **Service detection** - Automatically detects systemd vs SysV
- **Permission handling** - Proper user and directory setup

### **macOS Installation Wizard**
- **Native macOS design** - Follows Apple Human Interface Guidelines
- **Application bundle creation** - Professional .app bundle
- **Launchd integration** - Background service with macOS best practices
- **Retina-ready interface** - Optimized for modern Mac displays

## 🚀 Quick Start

### Windows Installation Wizard

**Prerequisites:**
- Download and install [NSIS](https://nsis.sourceforge.io/Download)

**Build the installer:**
```powershell
# From project root (requires favicon.ico in icons/ directory)
.\build-installer.ps1
```

**This creates:** `dist\ServinSetup-1.0.0.exe` (with Servin icon)

**To install:**
1. Run `ServinSetup-1.0.0.exe` as Administrator
2. Follow the installation wizard
3. Launch Servin GUI from Start Menu

### Linux Installation Wizard

**Prerequisites:**
```bash
# Ubuntu/Debian
sudo apt-get install python3-tk

# CentOS/RHEL
sudo yum install tkinter

# Arch Linux
sudo pacman -S tk
```

**To install:**
```bash
# Make executable and run
chmod +x installers/linux/servin-installer.py
sudo python3 installers/linux/servin-installer.py
```

### macOS Installation Wizard

**Prerequisites:**
- Python 3 with tkinter (usually included)
- macOS 10.12 (Sierra) or later

**To install:**
```bash
# Make executable and run
chmod +x installers/macos/servin-installer.py
sudo python3 installers/macos/servin-installer.py
```

## 📦 Installer Features

### Windows NSIS Wizard

**Welcome Screen:**
- System requirements check
- Windows version validation
- Administrator privilege verification

**License Agreement:**
- Apache 2.0 license display
- Acceptance required to proceed

**Component Selection:**
- Core Runtime (required)
- Desktop GUI Application
- Windows Service
- Start Menu Shortcuts

**Installation Directory:**
- Default: `C:\Program Files\Servin`
- Customizable installation path
- Data directory: `C:\ProgramData\Servin`

**Installation Progress:**
- Real-time progress indicator
- Detailed installation log
- Error handling and reporting

**Completion:**
- Service start options
- GUI launch options
- Installation summary

### Linux GUI Installer

**Welcome Page:**
- System information display
- Requirements verification
- Root privilege check

**License Agreement:**
- Scrollable license text
- Acceptance checkbox

**Installation Options:**
- Component selection (Core, GUI, Service, Shortcuts)
- Advanced options (user creation, PATH integration)

**Directory Configuration:**
- Installation directory (`/usr/local/bin`)
- Data directory (`/var/lib/servin`)
- Configuration directory (`/etc/servin`)
- Browse button for custom paths

**Confirmation Summary:**
- Review all settings
- Disk space calculation
- Final confirmation before installation

**Installation Progress:**
- Progress bar with status
- Real-time log output
- Error handling

**Success Page:**
- Installation summary
- Launch options
- Next steps guidance

### macOS Native Installer

**Welcome Screen:**
- macOS version detection
- System information display
- Retina-optimized interface

**License Agreement:**
- Native macOS text widget
- Proper scroll behavior

**Destination Selection:**
- Installation path configuration
- Disk space calculation
- Standard macOS directory structure

**Installation Type:**
- Standard vs custom installation
- Component selection
- Advanced options

**Installation Summary:**
- Professional summary display
- Configuration review

**Installation Progress:**
- macOS-style progress indicator
- Native logging interface

**Success Screen:**
- Completion confirmation
- Launch options
- Applications folder integration

## 🛠 Building Installer Wizards

### Windows (NSIS)

```powershell
# Install NSIS first
# Then build the installer wizard
.\build-installer.ps1 -Version "1.0.0"
```

**Output:** `dist\ServinSetup-1.0.0.exe`

**Features:**
- Professional Windows installer
- Add/Remove Programs integration
- Automatic uninstaller creation
- PATH environment updates
- Service installation and configuration

### Linux (Python/Tkinter)

The Linux installer is self-contained Python script:

```bash
# Copy binaries to installer directory
cp dist/servin-linux-1.0.0/servin installers/linux/
cp dist/servin-linux-1.0.0/servin-gui installers/linux/

# The installer script is ready to use
sudo python3 installers/linux/servin-installer.py
```

### macOS (Python/Tkinter with Cocoa)

Similar to Linux but with macOS-specific features:

```bash
# Copy binaries to installer directory
cp dist/servin-macos-1.0.0/servin installers/macos/
cp dist/servin-macos-1.0.0/servin-gui installers/macos/

# Run the installer
sudo python3 installers/macos/servin-installer.py
```

## 🔧 Customization

### NSIS Installer Customization

Edit `installers/windows/servin-installer.nsi`:

```nsis
; Change installer appearance
!define MUI_ICON "custom-icon.ico"
!define MUI_HEADERIMAGE_BITMAP "custom-header.bmp"

; Modify installation directory
InstallDir "$PROGRAMFILES\YourCompany\Servin"

; Add custom components
Section "Custom Component" SecCustom
  ; Your custom installation code
SectionEnd
```

### Python Installer Customization

Modify the Python installer scripts:

```python
# Change default directories
self.install_dir = tk.StringVar(value="/opt/servin")
self.data_dir = tk.StringVar(value="/opt/servin/data")

# Add custom installation steps
def custom_installation_step(self):
    self.log_message("Performing custom setup...")
    # Your custom code here
```

## 📊 Comparison: Wizards vs Scripts

| Feature | Wizard Installer | Script Installer |
|---------|------------------|-------------------|
| **User Experience** | Guided GUI | Command-line |
| **Error Handling** | Visual feedback | Text output |
| **Customization** | Interactive selection | Edit script |
| **Progress Tracking** | Progress bars | Text messages |
| **Accessibility** | Mouse/keyboard | Keyboard only |
| **Professional Look** | Native OS styling | Terminal interface |
| **Size** | Larger (GUI framework) | Smaller (text only) |
| **Dependencies** | GUI libraries | Shell/PowerShell |

## 🔍 Troubleshooting

### Windows NSIS Issues

**"NSIS not found":**
```powershell
# Install NSIS from official site
# Or specify custom path
.\build-installer.ps1 -NSISPath "C:\Custom\Path\makensis.exe"
```

**Installer won't run:**
- Check Windows SmartScreen settings
- Ensure running as Administrator
- Verify digital signature (if signed)

### Linux GUI Issues

**"tkinter not available":**
```bash
# Ubuntu/Debian
sudo apt-get install python3-tk

# CentOS/RHEL
sudo yum install tkinter python3-tkinter

# Fedora
sudo dnf install python3-tkinter
```

**Display issues:**
```bash
# Set DISPLAY variable if using SSH
export DISPLAY=:0

# Or use X11 forwarding
ssh -X user@host
```

### macOS GUI Issues

**"Permission denied":**
```bash
# Ensure running with sudo
sudo python3 servin-installer.py

# Check script permissions
chmod +x servin-installer.py
```

**App won't launch:**
```bash
# Remove quarantine attribute
sudo xattr -rd com.apple.quarantine "/Applications/Servin GUI.app"
```

## 🎨 UI Screenshots

### Windows Installer
- Professional NSIS wizard interface
- Component selection dialog
- Progress indication with logs
- Completion screen with launch options

### Linux Installer
- Native GTK/Tkinter styling
- Cross-platform compatibility
- System integration options
- Service configuration

### macOS Installer
- Aqua theme integration
- Retina display optimization
- Application bundle creation
- Launchd service setup

## 📈 Benefits of Installer Wizards

1. **Professional Appearance** - Native OS look and feel
2. **User-Friendly** - No command-line knowledge required
3. **Error Prevention** - Validation and guided setup
4. **Comprehensive** - Handles all installation aspects
5. **Uninstaller** - Clean removal capabilities
6. **Branding** - Custom icons, images, and text
7. **Compliance** - Follows OS installation standards

## 🚀 Future Enhancements

- **Code signing** for Windows and macOS installers
- **Multi-language support** for international users
- **Update mechanism** for automatic updates
- **Custom branding** options for enterprise deployment
- **Silent installation** modes for automation
- **Rollback capabilities** for failed installations

---

**The installer wizards provide a professional, user-friendly way to deploy Servin Container Runtime while maintaining all the functionality of the script-based installers.**
