# Servin Container Runtime - Complete NSIS Installer Script
# Creates a comprehensive Windows installer with VM prerequisites and dependencies
# Build with: makensis servin-installer.nsi

!define PRODUCT_NAME "Servin Container Runtime"
!define PRODUCT_VERSION "1.0.0"
!define PRODUCT_PUBLISHER "Servin Team"
!define PRODUCT_WEB_SITE "https://servin.dev"
!define PRODUCT_DIR_REGKEY "Software\Microsoft\Windows\CurrentVersion\App Paths\servin.exe"
!define PRODUCT_UNINST_KEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}"
!define PRODUCT_UNINST_ROOT_KEY "HKLM"

# Include modern UI and required headers
!include "MUI2.nsh"
!include "WinVer.nsh"
!include "x64.nsh"
!include "LogicLib.nsh"
!include "FileFunc.nsh"

# Installer properties
Name "${PRODUCT_NAME}"
OutFile "servin-installer-${PRODUCT_VERSION}.exe"
InstallDir "$PROGRAMFILES\Servin"
InstallDirRegKey HKLM "${PRODUCT_DIR_REGKEY}" ""
ShowInstDetails show
ShowUnInstDetails show
RequestExecutionLevel admin

# Version information
VIProductVersion "1.0.0.0"
VIAddVersionKey "ProductName" "${PRODUCT_NAME}"
VIAddVersionKey "ProductVersion" "${PRODUCT_VERSION}"
VIAddVersionKey "CompanyName" "${PRODUCT_PUBLISHER}"
VIAddVersionKey "FileDescription" "Servin Container Runtime Complete Installer"
VIAddVersionKey "FileVersion" "1.0.0.0"

# Modern UI configuration
!define MUI_ABORTWARNING
!if /FILEEXISTS "servin.ico"
  !define MUI_ICON "servin.ico"
  !define MUI_UNICON "servin.ico"
!endif
!define MUI_HEADERIMAGE
!define MUI_WELCOMEFINISHPAGE_BITMAP_NOSTRETCH
!define MUI_COMPONENTSPAGE_SMALLDESC

# Installer pages
!insertmacro MUI_PAGE_WELCOME
!if /FILEEXISTS "LICENSE.txt"
  !insertmacro MUI_PAGE_LICENSE "LICENSE.txt"
!endif
!insertmacro MUI_PAGE_COMPONENTS
!insertmacro MUI_PAGE_DIRECTORY

# Custom page for VM provider selection
Page custom VMProviderPage VMProviderPageLeave

!insertmacro MUI_PAGE_INSTFILES

# Custom completion page
Page custom CompletionPage

!insertmacro MUI_PAGE_FINISH

# Uninstaller pages
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

# Languages
!insertmacro MUI_LANGUAGE "English"

; Request admin privileges
RequestExecutionLevel admin

; Compression
SetCompressor /SOLID lzma

;--------------------------------
; Variables
Var StartMenuFolder

;--------------------------------
; Version Information
# Variables for installation options
Var VMProvider
Var InstallVirtualBox
Var InstallHyperV
Var InstallWSL2
Var InstallChocolatey
Var SystemCheck
Var NeedsRestart

# VM Provider selection dialog
Function VMProviderPage
  !insertmacro MUI_HEADER_TEXT "VM Provider Selection" "Choose your preferred virtualization provider"
  
  nsDialogs::Create 1018
  Pop $0
  ${If} $0 == error
    Abort
  ${EndIf}
  
  ${NSD_CreateLabel} 0 0 100% 20u "Servin can use multiple VM providers. Select your preferred option:"
  
  ${NSD_CreateRadioButton} 10 30u 280u 15u "Hyper-V (Windows Pro/Enterprise/Education)"
  Pop $1
  ${NSD_SetState} $1 ${BST_CHECKED}
  
  ${NSD_CreateRadioButton} 10 50u 280u 15u "VirtualBox (All Windows editions)"
  Pop $2
  
  ${NSD_CreateRadioButton} 10 70u 280u 15u "WSL2 (Windows 10 2004+/Windows 11)"
  Pop $3
  
  ${NSD_CreateRadioButton} 10 90u 280u 15u "Install all available providers"
  Pop $4
  
  ${NSD_CreateCheckBox} 10 120u 280u 15u "Download and cache base VM images"
  Pop $5
  ${NSD_SetState} $5 ${BST_CHECKED}
  
  ${NSD_CreateLabel} 0 145u 100% 40u "Note: The installer will automatically detect your system capabilities and install compatible providers. You can change this selection later."
  
  nsDialogs::Show
FunctionEnd

Function VMProviderPageLeave
  # Determine selected VM provider
  ${NSD_GetState} $1 $0
  ${If} $0 == ${BST_CHECKED}
    StrCpy $VMProvider "hyperv"
  ${EndIf}
  
  ${NSD_GetState} $2 $0
  ${If} $0 == ${BST_CHECKED}
    StrCpy $VMProvider "virtualbox"
  ${EndIf}
  
  ${NSD_GetState} $3 $0
  ${If} $0 == ${BST_CHECKED}
    StrCpy $VMProvider "wsl2"
  ${EndIf}
  
  ${NSD_GetState} $4 $0
  ${If} $0 == ${BST_CHECKED}
    StrCpy $VMProvider "all"
  ${EndIf}
FunctionEnd

# System requirements check
Function CheckSystemRequirements
  DetailPrint "Checking system requirements..."
  
  # Check Windows version (minimum Windows 10 1903)
  ${IfNot} ${AtLeastWin10}
    MessageBox MB_ICONSTOP "Windows 10 version 1903 or later is required."
    Abort
  ${EndIf}
  
  # Check for 64-bit
  ${IfNot} ${RunningX64}
    MessageBox MB_ICONSTOP "64-bit Windows is required."
    Abort
  ${EndIf}
  
  # Check available memory (minimum 4GB)
  ${GetTotalPhysicalMemory} $0 "MB"
  IntOp $1 $0 / 1024
  ${If} $1 < 4
    MessageBox MB_ICONQUESTION|MB_YESNO "System has only ${1}GB RAM. 4GB minimum recommended. Continue anyway?" IDYES +2
    Abort
  ${EndIf}
  
  # Check disk space (minimum 5GB)
  ${GetDrives} "HDD" "CheckDriveSpace"
  
  DetailPrint "System requirements check passed (${1}GB RAM available)"
FunctionEnd

Function CheckDriveSpace
  ${GetSize} "$9" "/S=0K" $0 $1 $2
  IntOp $3 $0 / 1048576  # Convert to GB
  ${If} $3 < 5
    StrCpy $SystemCheck "insufficient_disk"
    MessageBox MB_ICONQUESTION|MB_YESNO "Drive $9 has only ${3}GB free. 5GB minimum required. Continue anyway?" IDYES +2
    Abort
  ${EndIf}
  Push $0
FunctionEnd

# Install Chocolatey package manager
Function InstallChocolatey
  DetailPrint "Checking for Chocolatey package manager..."
  
  # Check if Chocolatey is already installed
  nsExec::ExecToStack 'powershell -Command "Get-Command choco -ErrorAction SilentlyContinue"'
  Pop $0
  ${If} $0 == 0
    DetailPrint "Chocolatey already installed"
    Return
  ${EndIf}
  
  DetailPrint "Installing Chocolatey package manager..."
  nsExec::ExecToLog 'powershell -NoProfile -ExecutionPolicy Bypass -Command "Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))"'
# Install Python and dependencies
Function InstallPython
  DetailPrint "Installing Python and dependencies..."
  
  # Check if Python is already installed
  nsExec::ExecToStack 'python --version'
  Pop $0
  ${If} $0 == 0
    DetailPrint "Python already installed"
  ${Else}
    DetailPrint "Installing Python via Chocolatey..."
    nsExec::ExecToLog 'choco install python3 -y --no-progress'
    Pop $0
  ${EndIf}
  
  # Install Python packages for GUI
  DetailPrint "Installing Python WebView dependencies..."
  nsExec::ExecToLog 'python -m pip install --upgrade pip --quiet'
  nsExec::ExecToLog 'python -m pip install pywebview[cef] flask flask-cors flask-socketio eventlet pyinstaller --quiet'
  Pop $0
  
  DetailPrint "Python dependencies installation completed"
FunctionEnd

# Install Hyper-V
Function InstallHyperV
  DetailPrint "Checking and enabling Hyper-V..."
  
  # Check if system supports Hyper-V
  nsExec::ExecToStack 'powershell -Command "(Get-WmiObject -Class Win32_ComputerSystem).HypervisorPresent"'
  Pop $0
  Pop $1
  ${If} $1 == "False"
    DetailPrint "Hyper-V not supported on this system"
    Return
  ${EndIf}
  
  # Check if Hyper-V is already enabled
  nsExec::ExecToStack 'powershell -Command "(Get-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V-All).State"'
  Pop $0
  Pop $1
  ${If} $1 == "Enabled"
    DetailPrint "Hyper-V already enabled"
    Return
  ${EndIf}
  
  # Enable Hyper-V feature
  DetailPrint "Enabling Hyper-V Windows feature..."
  nsExec::ExecToLog 'powershell -Command "Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V -All -NoRestart"'
  Pop $0
  ${If} $0 == 0
    DetailPrint "Hyper-V enabled successfully (restart may be required)"
    StrCpy $NeedsRestart "true"
  ${Else}
    DetailPrint "Failed to enable Hyper-V (error code: $0)"
  ${EndIf}
FunctionEnd

# Install VirtualBox
Function InstallVirtualBox
  DetailPrint "Installing VirtualBox..."
  
  # Check if VirtualBox is already installed
  nsExec::ExecToStack 'powershell -Command "Get-ItemProperty HKLM:\Software\Microsoft\Windows\CurrentVersion\Uninstall\* | Where-Object {$_.DisplayName -like '*VirtualBox*'}"'
  Pop $0
  ${If} $0 == 0
    DetailPrint "VirtualBox already installed"
    Return
  ${EndIf}
  
  # Install VirtualBox via Chocolatey
  DetailPrint "Installing VirtualBox via Chocolatey..."
  nsExec::ExecToLog 'choco install virtualbox -y --no-progress'
  Pop $0
  ${If} $0 == 0
    DetailPrint "VirtualBox installed successfully"
  ${Else}
    DetailPrint "VirtualBox installation failed (error code: $0)"
  ${EndIf}
FunctionEnd

# Install WSL2
Function InstallWSL2
  DetailPrint "Installing WSL2..."
  
  # Check Windows version for WSL2 support
  ${If} ${AtLeastWin10}
    ${AndIf} ${AtLeastBuild} 19041
    DetailPrint "Windows version supports WSL2"
  ${Else}
    DetailPrint "Windows version does not support WSL2"
    Return
  ${EndIf}
  
  # Enable WSL feature
  DetailPrint "Enabling Windows Subsystem for Linux..."
  nsExec::ExecToLog 'powershell -Command "Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Windows-Subsystem-Linux -NoRestart"'
  Pop $0
  
  # Enable Virtual Machine Platform
  DetailPrint "Enabling Virtual Machine Platform..."
  nsExec::ExecToLog 'powershell -Command "Enable-WindowsOptionalFeature -Online -FeatureName VirtualMachinePlatform -NoRestart"'
  Pop $0
  
  # Set WSL2 as default version
  DetailPrint "Setting WSL2 as default version..."
  nsExec::ExecToLog 'wsl --set-default-version 2'
  Pop $0
  
  StrCpy $NeedsRestart "true"
  DetailPrint "WSL2 installation completed (restart required)"
FunctionEnd

# Install development tools
Function InstallDevTools
  DetailPrint "Installing development tools..."
  
  # Install Packer
  DetailPrint "Installing HashiCorp Packer..."
  nsExec::ExecToLog 'choco install packer -y --no-progress'
  Pop $0
  
  # Install NSIS for future builds
  DetailPrint "Installing NSIS..."
  nsExec::ExecToLog 'choco install nsis -y --no-progress'
  Pop $0
  
  # Install additional development tools
  DetailPrint "Installing additional development tools..."
  nsExec::ExecToLog 'choco install git -y --no-progress'
  nsExec::ExecToLog 'choco install curl -y --no-progress'
  nsExec::ExecToLog 'choco install 7zip -y --no-progress'
  
  DetailPrint "Development tools installation completed"
FunctionEnd

# Download VM images
Function DownloadVMImages
  DetailPrint "Downloading base VM images..."
  
  # Create VM images directory
  CreateDirectory "$APPDATA\Servin\vm\images"
  SetOutPath "$APPDATA\Servin\vm\images"
  
  # Download Alpine Linux base image
  DetailPrint "Downloading Alpine Linux VM image..."
  NSISdl::download "https://dl-cdn.alpinelinux.org/alpine/v3.19/releases/x86_64/alpine-virt-3.19.1-x86_64.iso" "alpine-virt-3.19.1-x86_64.iso"
  Pop $0
  ${If} $0 == "success"
    DetailPrint "Alpine Linux image downloaded successfully"
  ${Else}
    DetailPrint "Warning: Failed to download Alpine Linux image"
  ${EndIf}
  
  # Download Ubuntu Server image (smaller cloud image)
  DetailPrint "Downloading Ubuntu Server cloud image..."
  NSISdl::download "https://cloud-images.ubuntu.com/releases/22.04/release/ubuntu-22.04-server-cloudimg-amd64.img" "ubuntu-22.04-server-cloudimg-amd64.img"
  Pop $0
  ${If} $0 == "success"
    DetailPrint "Ubuntu Server image downloaded successfully"
  ${Else}
    DetailPrint "Warning: Failed to download Ubuntu Server image"
  ${EndIf}
FunctionEnd

# Create Windows Service
Function CreateService
  DetailPrint "Creating Servin Windows Service..."
  
  # Stop existing service if running
  nsExec::ExecToLog 'sc stop "ServinRuntime"'
  nsExec::ExecToLog 'sc delete "ServinRuntime"'
  
  # Create new service
  nsExec::ExecToLog 'sc create "ServinRuntime" binPath= "$INSTDIR\servin.exe daemon" DisplayName= "Servin Container Runtime" start= auto description= "Servin Container Runtime Service"'
  Pop $0
  ${If} $0 == 0
    DetailPrint "Windows Service created successfully"
    
    # Start the service
    nsExec::ExecToLog 'sc start "ServinRuntime"'
    Pop $0
    ${If} $0 == 0
      DetailPrint "Service started successfully"
    ${Else}
      DetailPrint "Service created but failed to start (will start on reboot)"
    ${EndIf}
  ${Else}
    DetailPrint "Warning: Service creation failed (error code: $0)"
  ${EndIf}
FunctionEnd

# Initialize Servin VM support
Function InitializeVM
  DetailPrint "Initializing Servin VM support..."
  
  # Create VM configuration
  SetOutPath "$INSTDIR\config"
  FileOpen $0 "$INSTDIR\config\vm-config.yaml" w
  FileWrite $0 "vm:$\r$\n"
  FileWrite $0 "  platform: windows$\r$\n"
  FileWrite $0 "  providers:$\r$\n"
  
  ${If} $VMProvider == "hyperv"
  ${OrIf} $VMProvider == "all"
    FileWrite $0 "    - name: hyperv$\r$\n"
    FileWrite $0 "      priority: 1$\r$\n"
    FileWrite $0 "      enabled: true$\r$\n"
  ${EndIf}
  
  ${If} $VMProvider == "virtualbox"
  ${OrIf} $VMProvider == "all"
    FileWrite $0 "    - name: virtualbox$\r$\n"
    FileWrite $0 "      priority: 2$\r$\n"
    FileWrite $0 "      enabled: true$\r$\n"
  ${EndIf}
  
  ${If} $VMProvider == "wsl2"
  ${OrIf} $VMProvider == "all"
    FileWrite $0 "    - name: wsl2$\r$\n"
    FileWrite $0 "      priority: 3$\r$\n"
    FileWrite $0 "      enabled: true$\r$\n"
  ${EndIf}
  
  FileWrite $0 "  image_cache: $APPDATA\Servin\vm\images$\r$\n"
  FileWrite $0 "  vm_storage: $APPDATA\Servin\vm\instances$\r$\n"
  FileWrite $0 "  max_memory: 4GB$\r$\n"
  FileWrite $0 "  default_memory: 2GB$\r$\n"
  FileWrite $0 "  max_cpu_cores: 4$\r$\n"
  FileClose $0
  
  # Initialize VM subsystem
  nsExec::ExecToLog '"$INSTDIR\servin.exe" vm init'
  Pop $0
  ${If} $0 == 0
    DetailPrint "VM subsystem initialized successfully"
  ${Else}
    DetailPrint "VM initialization completed with warnings"
  ${EndIf}
FunctionEnd

# Completion page
Function CompletionPage
  !insertmacro MUI_HEADER_TEXT "Installation Complete" "Servin Container Runtime has been successfully installed"
  
  nsDialogs::Create 1018
  Pop $0
  
  ${NSD_CreateLabel} 0 0 100% 30u "Servin Container Runtime has been successfully installed with the following components:"
  
  ${NSD_CreateLabel} 10 40u 280u 15u "✓ Servin Container Runtime (CLI, GUI, TUI)"
  
  ${If} $VMProvider == "hyperv"
  ${OrIf} $VMProvider == "all"
    ${NSD_CreateLabel} 10 60u 280u 15u "✓ Hyper-V virtualization support"
  ${EndIf}
  
  ${If} $VMProvider == "virtualbox"
  ${OrIf} $VMProvider == "all"
    ${NSD_CreateLabel} 10 80u 280u 15u "✓ VirtualBox virtualization support"
  ${EndIf}
  
  ${If} $VMProvider == "wsl2"
  ${OrIf} $VMProvider == "all"
    ${NSD_CreateLabel} 10 100u 280u 15u "✓ WSL2 virtualization support"
  ${EndIf}
  
  ${NSD_CreateLabel} 10 120u 280u 15u "✓ Python WebView GUI framework"
  ${NSD_CreateLabel} 10 140u 280u 15u "✓ Development tools and dependencies"
  ${NSD_CreateLabel} 10 160u 280u 15u "✓ Windows Service integration"
  
  ${If} $NeedsRestart == "true"
    ${NSD_CreateLabel} 0 190u 100% 20u "⚠ A system restart is required to complete the installation." 
  ${Else}
    ${NSD_CreateLabel} 0 190u 100% 20u "The installation is complete and ready to use."
  ${EndIf}
  
  nsDialogs::Show
FunctionEnd
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_UNPAGE_FINISH

;--------------------------------
; Languages
!insertmacro MUI_LANGUAGE "English"

;--------------------------------
; Functions
Function .onInit
  ; Check Windows version
  ${IfNot} ${AtLeastWin10}
    MessageBox MB_OK|MB_ICONSTOP "Servin Container Runtime requires Windows 10 or later."
    Abort
  ${EndIf}
  
  ; Check if 64-bit
  ${IfNot} ${RunningX64}
    MessageBox MB_OK|MB_ICONSTOP "Servin Container Runtime requires a 64-bit version of Windows."
    Abort
  ${EndIf}
  
  ; Check admin rights
  UserInfo::GetAccountType
  pop $0
  ${If} $0 != "admin"
    MessageBox MB_OK|MB_ICONSTOP "Administrator privileges are required to install Servin Container Runtime."
    Abort
  ${EndIf}
FunctionEnd

Function .onInstSuccess
  ; Start the service
  ExecWait 'sc start "ServinRuntime"'
FunctionEnd

Function un.onInit
  ; Check admin rights for uninstall
  UserInfo::GetAccountType
  pop $0
  ${If} $0 != "admin"
    MessageBox MB_OK|MB_ICONSTOP "Administrator privileges are required to uninstall Servin Container Runtime."
    Abort
  ${EndIf}
FunctionEnd

;--------------------------------
; Installation Sections

; Core Runtime (Required)
# Main installer sections
Section "Servin Core Runtime" SEC01
  SectionIn RO  # Required section
  
  Call CheckSystemRequirements
  
  SetOutPath "$INSTDIR"
  SetOverwrite ifnewer
  
  # Install main executables (these files need to be in the build directory)
  File "servin.exe"
  File "servin-tui.exe" 
  
  # Install GUI executable if available
  IfFileExists "servin-gui.exe" 0 +2
  File "servin-gui.exe"
  
  # Install configuration files
  SetOutPath "$INSTDIR\config"
  File "servin.conf"
  
  # Create data directories
  CreateDirectory "$APPDATA\Servin"
  CreateDirectory "$APPDATA\Servin\vm"
  CreateDirectory "$APPDATA\Servin\vm\images"
  CreateDirectory "$APPDATA\Servin\vm\instances"
  CreateDirectory "$APPDATA\Servin\config"
  CreateDirectory "$APPDATA\Servin\logs"
  
  # Create registry entries
  WriteRegStr HKLM "${PRODUCT_DIR_REGKEY}" "" "$INSTDIR\servin.exe"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "DisplayName" "$(^Name)"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "UninstallString" "$INSTDIR\uninst.exe"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "DisplayIcon" "$INSTDIR\servin.exe"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "DisplayVersion" "${PRODUCT_VERSION}"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "URLInfoAbout" "${PRODUCT_WEB_SITE}"
  WriteRegStr ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "Publisher" "${PRODUCT_PUBLISHER}"
  
  # Add to PATH environment variable
  ${EnvVarUpdate} $0 "PATH" "A" "HKLM" "$INSTDIR"
  
  Call CreateService
SectionEnd

Section "VM Prerequisites" SEC02
  Call InstallChocolatey
  Call InstallPython
  
  # Install selected VM provider(s)
  ${If} $VMProvider == "hyperv"
  ${OrIf} $VMProvider == "all"
    Call InstallHyperV
  ${EndIf}
  
  ${If} $VMProvider == "virtualbox" 
  ${OrIf} $VMProvider == "all"
    Call InstallVirtualBox
  ${EndIf}
  
  ${If} $VMProvider == "wsl2"
  ${OrIf} $VMProvider == "all"
    Call InstallWSL2
  ${EndIf}
  
  Call InstallDevTools
  Call InitializeVM
SectionEnd

Section "Desktop Integration" SEC03
  # Create desktop shortcuts - GUI only if available
  IfFileExists "$INSTDIR\servin-gui.exe" 0 +2
    CreateShortCut "$DESKTOP\Servin GUI.lnk" "$INSTDIR\servin-gui.exe" "" "$INSTDIR\servin.exe" 0
  
  # Create Start Menu folder
  CreateDirectory "$SMPROGRAMS\Servin"
  
  # Create GUI shortcut only if GUI is available
  IfFileExists "$INSTDIR\servin-gui.exe" 0 +2
    CreateShortCut "$SMPROGRAMS\Servin\Servin GUI.lnk" "$INSTDIR\servin-gui.exe" "" "$INSTDIR\servin.exe" 0
  
  CreateShortCut "$SMPROGRAMS\Servin\Servin TUI.lnk" "$INSTDIR\servin-tui.exe" "" "$INSTDIR\servin.exe" 1
  CreateShortCut "$SMPROGRAMS\Servin\Command Prompt.lnk" "$SYSDIR\cmd.exe" "/k cd /d $\"$INSTDIR$\"" "$SYSDIR\cmd.exe" 0
  CreateShortCut "$SMPROGRAMS\Servin\Uninstall.lnk" "$INSTDIR\uninst.exe" "" "$INSTDIR\uninst.exe" 0
  
  # Register file associations for .servin config files
  WriteRegStr HKCR ".servin" "" "ServinConfig"
  WriteRegStr HKCR "ServinConfig" "" "Servin Configuration File"
  WriteRegStr HKCR "ServinConfig\DefaultIcon" "" "$INSTDIR\servin.exe,0"
  WriteRegStr HKCR "ServinConfig\shell\open\command" "" '"$INSTDIR\servin.exe" config edit "%1"'
  
  # Add context menu for containers
  WriteRegStr HKCR "Directory\shell\ServinHere" "" "Open Servin Here"
  WriteRegStr HKCR "Directory\shell\ServinHere\command" "" '"$INSTDIR\servin-tui.exe" --workdir "%1"'
SectionEnd

Section /o "VM Images" SEC04
  Call DownloadVMImages
SectionEnd

Section "Documentation" SEC05
  SetOutPath "$INSTDIR\docs"
  
  # Copy documentation files (these need to be in the build directory)
  File /nonfatal "README.txt"
  File /nonfatal "VM_PREREQUISITES.md"
  File /nonfatal "INSTALL.md"
  
  # Create quick start guide
  FileOpen $0 "$INSTDIR\docs\QuickStart.txt" w
  FileWrite $0 "Servin Container Runtime - Quick Start Guide$\r$\n"
  FileWrite $0 "===========================================$\r$\n$\r$\n"
  FileWrite $0 "Getting Started:$\r$\n"
  FileWrite $0 "1. Launch Servin GUI from Start Menu or Desktop$\r$\n"
  FileWrite $0 "2. Initialize VM: servin vm init$\r$\n"
  FileWrite $0 "3. Run your first container: servin run --vm alpine echo 'Hello World!'$\r$\n$\r$\n"
  FileWrite $0 "Command Line Usage:$\r$\n"
  FileWrite $0 "- servin version          Show version$\r$\n"
  FileWrite $0 "- servin vm status        Check VM status$\r$\n"
  FileWrite $0 "- servin run alpine sh    Run interactive Alpine container$\r$\n"
  FileWrite $0 "- servin-gui              Launch GUI$\r$\n"
  FileWrite $0 "- servin-tui              Launch Terminal UI$\r$\n$\r$\n"
  FileWrite $0 "For more information, visit: ${PRODUCT_WEB_SITE}$\r$\n"
  FileClose $0
SectionEnd
  CreateDirectory "$APPDATA\Servin\logs"
  CreateDirectory "$APPDATA\Servin\volumes"
  CreateDirectory "$APPDATA\Servin\images"
  
  ; Install configuration file (if available)
  SetOutPath "$APPDATA\Servin\config"
  File /nonfatal "package\servin.conf"
  
  ; Install documentation
  SetOutPath "$INSTDIR"
  File "package\README.txt"
  File "package\LICENSE.txt"
  
  ; Create uninstaller
  WriteUninstaller "$INSTDIR\Uninstall.exe"
  
  ; Registry entries
  WriteRegStr HKLM "Software\Servin" "InstallPath" "$INSTDIR"
  WriteRegStr HKLM "Software\Servin" "Version" "1.0.0"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin" "DisplayName" "Servin Container Runtime"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin" "UninstallString" "$INSTDIR\Uninstall.exe"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin" "InstallLocation" "$INSTDIR"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin" "DisplayIcon" "$INSTDIR\servin.exe"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin" "Publisher" "Servin Project"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin" "DisplayVersion" "1.0.0"
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin" "NoModify" 1
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Servin" "NoRepair" 1
  
  ; Add to PATH
  Push "$INSTDIR"
  Call AddToPath
SectionEnd

; GUI Application
Section "Desktop GUI" SecGUI
  SetOutPath "$INSTDIR"
  File "package\servin-gui.exe"
  
  ; Create desktop shortcut
  CreateShortcut "$DESKTOP\Servin GUI.lnk" "$INSTDIR\servin-gui.exe" "" "$INSTDIR\servin-gui.exe" 0
SectionEnd

; Windows Service
Section "Windows Service" SecService
  ; Install service using Windows sc command
  ExecWait 'sc create "ServinRuntime" binpath= "$INSTDIR\servin.exe daemon" displayname= "Servin Container Runtime" start= auto'
  
  ; Set service description
  ExecWait 'sc description "ServinRuntime" "Servin Container Runtime Service providing Docker-compatible container management"'
  
  ; Set service to restart on failure
  ExecWait 'sc failure "ServinRuntime" reset= 86400 actions= restart/60000/restart/60000/restart/60000'
SectionEnd

; Start Menu Shortcuts
Section "Start Menu Shortcuts" SecStartMenu
  !insertmacro MUI_STARTMENU_WRITE_BEGIN Application
  
  CreateDirectory "$SMPROGRAMS\$StartMenuFolder"
  CreateShortcut "$SMPROGRAMS\$StartMenuFolder\Servin GUI.lnk" "$INSTDIR\servin-gui.exe"
  CreateShortcut "$SMPROGRAMS\$StartMenuFolder\Servin Desktop (TUI).lnk" "$INSTDIR\servin-tui.exe"
  CreateShortcut "$SMPROGRAMS\$StartMenuFolder\Servin Command Prompt.lnk" "cmd.exe" "/k cd /d $\"$INSTDIR$\""
  CreateShortcut "$SMPROGRAMS\$StartMenuFolder\Uninstall Servin.lnk" "$INSTDIR\Uninstall.exe"
  CreateShortcut "$SMPROGRAMS\$StartMenuFolder\README.lnk" "$INSTDIR\README.txt"
  
  !insertmacro MUI_STARTMENU_WRITE_END
SectionEnd

;--------------------------------
; Section Descriptions
!insertmacro MUI_FUNCTION_DESCRIPTION_BEGIN
  !insertmacro MUI_DESCRIPTION_TEXT ${SecCore} "Core Servin container runtime (required)"
  !insertmacro MUI_DESCRIPTION_TEXT ${SecGUI} "Desktop GUI application for managing containers"
  !insertmacro MUI_DESCRIPTION_TEXT ${SecService} "Install as Windows Service for automatic startup"
  !insertmacro MUI_DESCRIPTION_TEXT ${SecStartMenu} "Create Start Menu shortcuts"
!insertmacro MUI_FUNCTION_DESCRIPTION_END

;--------------------------------
; Uninstaller Section
Section "Uninstall"
  ; Stop and remove service
  ExecWait 'sc stop "ServinRuntime"'
  ExecWait 'sc delete "ServinRuntime"'
  
  ; Remove files
  Delete "$INSTDIR\servin.exe"
  Delete "$INSTDIR\servin-tui.exe"
  Delete "$INSTDIR\servin-gui.exe"
  Delete "$INSTDIR\README.txt"
  Delete "$INSTDIR\LICENSE.txt"
  Delete "$INSTDIR\Uninstall.exe"
  
  ; Remove directories
  RMDir "$INSTDIR"
  
  ; Remove data directories (ask user)
  MessageBox MB_YESNO|MB_ICONQUESTION "Do you want to remove all container data, images, and volumes? This cannot be undone." IDNO skip_data_removal
  RMDir /r "$APPDATA\Servin"
  skip_data_removal:
  
  ; Remove shortcuts
  Delete "$DESKTOP\Servin GUI.lnk"
  
  ; Remove Start Menu shortcuts
  !insertmacro MUI_STARTMENU_GETFOLDER Application $StartMenuFolder
  Delete "$SMPROGRAMS\$StartMenuFolder\Servin GUI.lnk"
  Delete "$SMPROGRAMS\$StartMenuFolder\Servin Command Prompt.lnk"
  Delete "$SMPROGRAMS\$StartMenuFolder\Uninstall Servin.lnk"
  Delete "$SMPROGRAMS\$StartMenuFolder\README.lnk"
  RMDir "$SMPROGRAMS\$StartMenuFolder"
  
  ; Remove from PATH
  Push "$INSTDIR"
  Call un.RemoveFromPath
# Component descriptions
!insertmacro MUI_FUNCTION_DESCRIPTION_BEGIN
  !insertmacro MUI_DESCRIPTION_TEXT ${SEC01} "Core Servin Container Runtime with CLI, GUI, and TUI interfaces. This component is required."
  !insertmacro MUI_DESCRIPTION_TEXT ${SEC02} "VM prerequisites including Python, virtualization providers (Hyper-V/VirtualBox/WSL2), and development tools."
  !insertmacro MUI_DESCRIPTION_TEXT ${SEC03} "Desktop shortcuts, Start Menu entries, file associations, and context menu integration."
  !insertmacro MUI_DESCRIPTION_TEXT ${SEC04} "Download base VM images for immediate use (Alpine Linux, Ubuntu Server). Requires internet connection."
  !insertmacro MUI_DESCRIPTION_TEXT ${SEC05} "Documentation files including Quick Start guide, README, and setup instructions."
!insertmacro MUI_FUNCTION_DESCRIPTION_END

# Installer initialization
Function .onInit
  # Initialize variables
  StrCpy $VMProvider "hyperv"
  StrCpy $NeedsRestart "false"
  StrCpy $SystemCheck "ok"
  
  # Check for existing installation
  ReadRegStr $R0 ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}" "UninstallString"
  StrCmp $R0 "" done
  
  MessageBox MB_OKCANCEL|MB_ICONEXCLAMATION \
  "${PRODUCT_NAME} is already installed.$\n$\nClick OK to remove the previous version or Cancel to cancel this upgrade." \
  IDOK uninst
  Abort
  
uninst:
  ClearErrors
  ExecWait '$R0 /S _?=$INSTDIR'
  
  IfErrors no_remove_uninstaller done
    Delete "$R0"
  no_remove_uninstaller:
  
done:
FunctionEnd

Function .onInstSuccess
  # Test basic functionality
  nsExec::ExecToLog '"$INSTDIR\servin.exe" version'
  
  ${If} $NeedsRestart == "true"
    MessageBox MB_ICONINFORMATION "Installation completed successfully!$\n$\nA system restart is required to enable some VM features. Please restart your computer to use all functionality."
  ${Else}
    MessageBox MB_ICONINFORMATION "Installation completed successfully!$\n$\nServin Container Runtime is now ready to use. You can launch the GUI from the Start Menu or desktop shortcut."
  ${EndIf}
FunctionEnd

# Uninstaller section
Section Uninstall
  # Stop and remove service
  nsExec::ExecToLog 'sc stop "ServinRuntime"'
  nsExec::ExecToLog 'sc delete "ServinRuntime"'
  
  # Remove from PATH
  ${un.EnvVarUpdate} $0 "PATH" "R" "HKLM" "$INSTDIR"
  
  # Remove files
  Delete "$INSTDIR\servin.exe"
  Delete "$INSTDIR\servin-tui.exe"
  Delete "$INSTDIR\servin-gui.exe"
  Delete "$INSTDIR\uninst.exe"
  
  # Remove directories
  RMDir /r "$INSTDIR\config"
  RMDir /r "$INSTDIR\docs"
  RMDir "$INSTDIR"
  
  # Remove shortcuts
  Delete "$DESKTOP\Servin GUI.lnk"
  RMDir /r "$SMPROGRAMS\Servin"
  
  # Remove registry keys
  DeleteRegKey ${PRODUCT_UNINST_ROOT_KEY} "${PRODUCT_UNINST_KEY}"
  DeleteRegKey HKLM "${PRODUCT_DIR_REGKEY}"
  DeleteRegKey HKCR ".servin"
  DeleteRegKey HKCR "ServinConfig"
  DeleteRegKey HKCR "Directory\shell\ServinHere"
  
  # Ask about removing user data
  MessageBox MB_ICONQUESTION|MB_YESNO "Do you want to remove user data and VM images?$\n$\nThis will delete all containers, images, and configuration files." IDNO skip_userdata
  
  RMDir /r "$APPDATA\Servin"
  
skip_userdata:
  SetAutoClose true
SectionEnd

# Include environment variable functions
!include "EnvVarUpdate.nsh"
    SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000
  ${EndIf}
  
  Pop $2
  Pop $1
  Pop $0
FunctionEnd

Function un.RemoveFromPath
  Exch $0
  Push $1
  Push $2
  
  Call un.IsNT
  Pop $1
  ${If} $1 == 1
    ReadRegStr $1 ${Environ} "PATH"
    Push $1
    Push $0
    Call un.StrStr
    Pop $2
    ${If} $2 == ""
      Goto done
    ${EndIf}
    
    StrLen $2 $0
    StrCpy $1 $1 $2
    StrCpy $0 $0 "" $2
    StrCpy $0 "$1$0"
    WriteRegExpandStr ${Environ} "PATH" $0
    SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000
  ${EndIf}
  
  done:
  Pop $2
  Pop $1
  Pop $0
FunctionEnd

Function IsNT
  Push $0
  ReadRegStr $0 HKLM "SOFTWARE\Microsoft\Windows NT\CurrentVersion" CurrentVersion
  ${If} $0 == ""
    Push 0
  ${Else}
    Push 1
  ${EndIf}
  Exch
  Pop $0
FunctionEnd

Function un.IsNT
  Push $0
  ReadRegStr $0 HKLM "SOFTWARE\Microsoft\Windows NT\CurrentVersion" CurrentVersion
  ${If} $0 == ""
    Push 0
  ${Else}
    Push 1
  ${EndIf}
  Exch
  Pop $0
FunctionEnd

Function un.StrStr
  Exch $R1
  Exch
  Exch $R2
  Push $R3
  Push $R4
  Push $R5
  StrLen $R3 $R1
  StrCpy $R4 0
  
  loop:
    StrCpy $R5 $R2 $R3 $R4
    ${If} $R5 == $R1
      StrCpy $R0 $R2 $R4
      Goto done
    ${EndIf}
    ${If} $R3 >= $R4
      StrCpy $R0 ""
      Goto done
    ${EndIf}
    IntOp $R4 $R4 + 1
    Goto loop
  
  done:
  Pop $R5
  Pop $R4
  Pop $R3
  Pop $R2
  Exch $R0
  Pop $R1
FunctionEnd
