# Servin Container Runtime - Enhanced Windows Installer with VM Prerequisites
# Run this script as Administrator for full installation

param(
    [switch]$Uninstall,
    [string]$InstallDir = "C:\Program Files\Servin",
    [string]$DataDir = "C:\ProgramData\Servin",
    [switch]$SkipPrerequisites,
    [switch]$SkipVMSetup,
    [switch]$Silent = $false
)

$ErrorActionPreference = "Stop"

# Color output functions
function Write-Success { param([string]$Text) Write-Host "✓ $Text" -ForegroundColor Green }
function Write-Warning { param([string]$Text) Write-Host "⚠ $Text" -ForegroundColor Yellow }
function Write-Error { param([string]$Text) Write-Host "✗ $Text" -ForegroundColor Red }
function Write-Info { param([string]$Text) Write-Host "→ $Text" -ForegroundColor Cyan }
function Write-Header { param([string]$Text) Write-Host "`n$('=' * 60)`n$Text`n$('=' * 60)" -ForegroundColor Cyan }

# Check if running as administrator
if (-NOT ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Error "This script must be run as Administrator."
    Write-Info "Right-click PowerShell and select 'Run as Administrator', then run this script again."
    exit 1
}

function Test-Prerequisites {
    Write-Header "Checking System Prerequisites"
    
    $prerequisites = @{
        "CPU_Virtualization" = $false
        "Windows_Version" = $false
        "PowerShell_Version" = $false
        "Internet_Connection" = $false
        "Disk_Space" = $false
    }
    
    # Check Windows version
    Write-Info "Checking Windows version..."
    $winVersion = [System.Environment]::OSVersion.Version
    if ($winVersion.Major -ge 10) {
        Write-Success "Windows 10/11 detected (Version: $($winVersion.Major).$($winVersion.Minor))"
        $prerequisites["Windows_Version"] = $true
    } else {
        Write-Error "Windows 10 or later required (Current: $($winVersion.Major).$($winVersion.Minor))"
    }
    
    # Check PowerShell version
    Write-Info "Checking PowerShell version..."
    if ($PSVersionTable.PSVersion.Major -ge 5) {
        Write-Success "PowerShell 5.0+ detected (Version: $($PSVersionTable.PSVersion))"
        $prerequisites["PowerShell_Version"] = $true
    } else {
        Write-Error "PowerShell 5.0+ required (Current: $($PSVersionTable.PSVersion))"
    }
    
    # Check CPU virtualization
    Write-Info "Checking CPU virtualization support..."
    try {
        $cpu = Get-WmiObject -Class Win32_Processor | Select-Object -First 1
        if ($cpu.VMMonitorModeExtensions -eq $true) {
            Write-Success "CPU supports hardware virtualization"
            $prerequisites["CPU_Virtualization"] = $true
        } else {
            Write-Warning "CPU may not support hardware virtualization"
            Write-Info "VM features will use software emulation"
        }
    } catch {
        Write-Warning "Could not detect CPU virtualization support"
    }
    
    # Check internet connection
    Write-Info "Checking internet connectivity..."
    try {
        $response = Invoke-WebRequest -Uri "https://www.google.com" -TimeoutSec 10 -UseBasicParsing
        if ($response.StatusCode -eq 200) {
            Write-Success "Internet connection available"
            $prerequisites["Internet_Connection"] = $true
        }
    } catch {
        Write-Warning "Internet connection not available - some features may not work"
    }
    
    # Check disk space (5GB minimum)
    Write-Info "Checking available disk space..."
    $systemDrive = Get-PSDrive -Name ($env:SystemDrive.Replace(":", ""))
    $freeSpaceGB = [math]::Round($systemDrive.Free / 1GB, 2)
    if ($freeSpaceGB -ge 5) {
        Write-Success "Sufficient disk space available ($freeSpaceGB GB free)"
        $prerequisites["Disk_Space"] = $true
    } else {
        Write-Error "Insufficient disk space. Need 5GB, have $freeSpaceGB GB"
    }
    
    return $prerequisites
}

function Install-Chocolatey {
    Write-Header "Installing Chocolatey Package Manager"
    
    # Check if Chocolatey is already installed
    if (Get-Command choco -ErrorAction SilentlyContinue) {
        Write-Success "Chocolatey already installed"
        return $true
    }
    
    Write-Info "Installing Chocolatey..."
    try {
        Set-ExecutionPolicy Bypass -Scope Process -Force
        [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
        Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
        
        # Refresh environment variables
        $env:ChocolateyInstall = Convert-Path "$((Get-Command choco).Path)\..\.."
        Import-Module "$env:ChocolateyInstall\helpers\chocolateyProfile.psm1"
        refreshenv
        
        Write-Success "Chocolatey installed successfully"
        return $true
    } catch {
        Write-Error "Failed to install Chocolatey: $($_.Exception.Message)"
        return $false
    }
}

function Install-VMPrerequisites {
    Write-Header "Installing VM Prerequisites"
    
    $vmCapabilities = @{
        "Hyper_V" = $false
        "VirtualBox" = $false
        "WSL2" = $false
    }
    
    # Install Chocolatey first
    if (-not (Install-Chocolatey)) {
        Write-Error "Cannot proceed without Chocolatey package manager"
        return $vmCapabilities
    }
    
    # Check and enable Hyper-V
    Write-Info "Checking Hyper-V availability..."
    try {
        $hypervFeature = Get-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V-All -ErrorAction SilentlyContinue
        if ($hypervFeature) {
            if ($hypervFeature.State -eq "Enabled") {
                Write-Success "Hyper-V is already enabled"
                $vmCapabilities["Hyper_V"] = $true
            } elseif ($hypervFeature.State -eq "Disabled") {
                Write-Info "Enabling Hyper-V..."
                if (-not $Silent) {
                    $response = Read-Host "Enable Hyper-V? This requires a restart (y/N)"
                    if ($response -eq "y" -or $response -eq "Y") {
                        Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V-All -All
                        Write-Success "Hyper-V enabled (restart required)"
                        $vmCapabilities["Hyper_V"] = $true
                    }
                } else {
                    Write-Info "Skipping Hyper-V in silent mode"
                }
            }
        } else {
            Write-Warning "Hyper-V not available on this Windows edition"
            Write-Info "Requires Windows 10/11 Pro, Enterprise, or Education"
        }
    } catch {
        Write-Warning "Could not check Hyper-V status: $($_.Exception.Message)"
    }
    
    # Install VirtualBox
    Write-Info "Installing VirtualBox..."
    try {
        choco install virtualbox --params "/NoDesktopShortcut /NoQuickLaunchShortcut /NoStartupShortcut" -y --no-progress --limit-output
        
        # Verify VirtualBox installation
        if (Get-Command VBoxManage -ErrorAction SilentlyContinue) {
            $vboxVersion = & VBoxManage --version 2>$null
            Write-Success "VirtualBox installed successfully (Version: $vboxVersion)"
            $vmCapabilities["VirtualBox"] = $true
        } else {
            Write-Warning "VirtualBox installation may have failed"
        }
    } catch {
        Write-Warning "Failed to install VirtualBox: $($_.Exception.Message)"
    }
    
    # Check and install WSL2
    Write-Info "Checking WSL2 availability..."
    try {
        $wslStatus = wsl --status 2>$null
        if ($LASTEXITCODE -eq 0) {
            Write-Success "WSL2 is already available"
            $vmCapabilities["WSL2"] = $true
        } else {
            Write-Info "Installing WSL2..."
            if (-not $Silent) {
                $response = Read-Host "Install WSL2? This may require a restart (y/N)"
                if ($response -eq "y" -or $response -eq "Y") {
                    wsl --install --no-distribution
                    Write-Success "WSL2 installation initiated (restart may be required)"
                    $vmCapabilities["WSL2"] = $true
                }
            } else {
                Write-Info "Skipping WSL2 in silent mode"
            }
        }
    } catch {
        Write-Warning "Could not install WSL2: $($_.Exception.Message)"
    }
    
    # Install additional development tools
    Write-Info "Installing development tools..."
    try {
        # NSIS for installer creation
        choco install nsis -y --no-progress --limit-output
        Write-Success "NSIS installed"
        
        # Packer for VM image building
        choco install packer -y --no-progress --limit-output
        Write-Success "Packer installed"
        
        # Install Packer QEMU plugin
        packer plugins install github.com/hashicorp/qemu 2>$null
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Packer QEMU plugin installed"
        }
    } catch {
        Write-Warning "Some development tools may not have installed correctly"
    }
    
    return $vmCapabilities
}

function Install-PythonDependencies {
    Write-Header "Installing Python Dependencies"
    
    # Check if Python is available
    if (-not (Get-Command python -ErrorAction SilentlyContinue)) {
        Write-Info "Installing Python..."
        try {
            choco install python -y --no-progress --limit-output
            # Refresh PATH
            $env:PATH = [Environment]::GetEnvironmentVariable("PATH", "Machine") + ";" + [Environment]::GetEnvironmentVariable("PATH", "User")
        } catch {
            Write-Warning "Failed to install Python automatically"
            Write-Info "Please install Python manually from https://python.org"
            return $false
        }
    }
    
    # Install Python packages for GUI
    Write-Info "Installing Python WebView dependencies..."
    try {
        python -m pip install --upgrade pip --quiet
        python -m pip install pywebview flask flask-cors flask-socketio eventlet gevent pyinstaller --quiet
        
        # Test imports
        python -c "import webview; print('✓ pywebview available')"
        python -c "import flask_socketio; print('✓ flask-socketio available')"
        python -c "import eventlet; print('✓ eventlet available')"
        
        Write-Success "Python dependencies installed successfully"
        return $true
    } catch {
        Write-Warning "Failed to install Python dependencies: $($_.Exception.Message)"
        return $false
    }
}

function Install-Servin {
    Write-Header "Installing Servin Container Runtime"
    
    $LogDir = Join-Path $DataDir "logs"
    $ConfigDir = Join-Path $DataDir "config"
    $VolumesDir = Join-Path $DataDir "volumes"
    $ImagesDir = Join-Path $DataDir "images"
    $VMDir = Join-Path $DataDir "vm"
    $VMImagesDir = Join-Path $VMDir "images"
    $VMInstancesDir = Join-Path $VMDir "instances"

    Write-Info "Installation directory: $InstallDir"
    Write-Info "Data directory: $DataDir"

    # Create directories
    Write-Info "Creating directories..."
    @($InstallDir, $DataDir, $LogDir, $ConfigDir, $VolumesDir, $ImagesDir, $VMDir, $VMImagesDir, $VMInstancesDir) | ForEach-Object {
        if (-not (Test-Path $_)) {
            New-Item -ItemType Directory -Path $_ -Force | Out-Null
            Write-Success "Created: $_"
        }
    }

    # Copy executables
    Write-Info "Installing executables..."
    $executables = @("servin.exe", "servin-tui.exe", "servin-gui.exe")
    $installedCount = 0
    
    foreach ($exe in $executables) {
        if (Test-Path $exe) {
            Copy-Item $exe -Destination (Join-Path $InstallDir $exe) -Force
            Write-Success "Installed: $exe"
            $installedCount++
        } else {
            Write-Warning "Executable not found: $exe"
        }
    }

    if ($installedCount -eq 0) {
        Write-Error "No executables found. Please run this installer from the directory containing the Servin executables."
        return $false
    }

    # Set permissions
    Write-Info "Setting permissions..."
    try {
        $acl = Get-Acl $InstallDir
        $accessRule = New-Object System.Security.AccessControl.FileSystemAccessRule("SYSTEM","FullControl","ContainerInherit,ObjectInherit","None","Allow")
        $acl.SetAccessRule($accessRule)
        $accessRule = New-Object System.Security.AccessControl.FileSystemAccessRule("Administrators","FullControl","ContainerInherit,ObjectInherit","None","Allow")
        $acl.SetAccessRule($accessRule)
        Set-Acl -Path $InstallDir -AclObject $acl
        Write-Success "Permissions configured"
    } catch {
        Write-Warning "Could not set permissions: $($_.Exception.Message)"
    }

    # Add to PATH
    Write-Info "Adding to system PATH..."
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "Machine")
    if ($currentPath -notlike "*$InstallDir*") {
        [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$InstallDir", "Machine")
        Write-Success "Added to system PATH"
    } else {
        Write-Success "Already in system PATH"
    }

    # Create VM configuration
    Write-Info "Creating VM configuration..."
    $vmConfigContent = @"
vm:
  platform: windows
  providers:
    - name: hyperv
      priority: 1
      enabled: true
      acceleration: true
    - name: virtualbox
      priority: 2
      enabled: true
      acceleration: false
    - name: wsl2
      priority: 3
      enabled: true
      acceleration: true
  default_provider: hyperv
  image_cache: "$VMImagesDir"
  vm_storage: "$VMInstancesDir"
  max_memory: "2GB"
  default_memory: "1GB"
  max_cpu_cores: 2
"@
    $vmConfigPath = Join-Path $ConfigDir "vm-config.yaml"
    $vmConfigContent | Out-File -FilePath $vmConfigPath -Encoding UTF8
    Write-Success "VM configuration created: $vmConfigPath"

    # Create main configuration
    Write-Info "Creating main configuration..."
    $configContent = @"
# Servin Configuration File
data_dir=$DataDir
log_level=info
log_file=$LogDir\servin.log
runtime=vm
bridge_name=servin0
cri_port=10250
cri_enabled=false
vm_enabled=true
vm_config=$vmConfigPath
"@
    $configPath = Join-Path $ConfigDir "servin.conf"
    $configContent | Out-File -FilePath $configPath -Encoding UTF8
    Write-Success "Main configuration created: $configPath"

    return $true
}

function Initialize-VM {
    Write-Header "Initializing VM Support"
    
    $servinExe = Join-Path $InstallDir "servin.exe"
    if (-not (Test-Path $servinExe)) {
        Write-Error "Servin executable not found: $servinExe"
        return $false
    }
    
    Write-Info "Initializing VM directories..."
    try {
        & $servinExe vm init
        Write-Success "VM directories initialized"
    } catch {
        Write-Warning "VM initialization may have failed: $($_.Exception.Message)"
    }
    
    Write-Info "Testing VM providers..."
    try {
        & $servinExe vm list-providers
        Write-Success "VM providers detected"
    } catch {
        Write-Warning "VM provider detection failed: $($_.Exception.Message)"
    }
    
    return $true
}

function Show-InstallationSummary {
    param([hashtable]$Prerequisites, [hashtable]$VMCapabilities)
    
    Write-Header "Installation Summary"
    
    Write-Host "`nPrerequisites Status:" -ForegroundColor Yellow
    foreach ($prereq in $Prerequisites.GetEnumerator()) {
        $status = if ($prereq.Value) { "✓ PASS" } else { "✗ FAIL" }
        $color = if ($prereq.Value) { "Green" } else { "Red" }
        Write-Host "$status - $($prereq.Key -replace '_', ' ')" -ForegroundColor $color
    }
    
    Write-Host "`nVM Capabilities:" -ForegroundColor Yellow
    foreach ($capability in $VMCapabilities.GetEnumerator()) {
        $status = if ($capability.Value) { "✓ AVAILABLE" } else { "✗ NOT AVAILABLE" }
        $color = if ($capability.Value) { "Green" } else { "Yellow" }
        Write-Host "$status - $($capability.Key -replace '_', ' ')" -ForegroundColor $color
    }
    
    Write-Host "`nNext Steps:" -ForegroundColor Cyan
    Write-Host "1. Restart your computer if Hyper-V or WSL2 were installed"
    Write-Host "2. Open a new PowerShell/Command Prompt window"
    Write-Host "3. Run: servin vm init"
    Write-Host "4. Run: servin vm enable"
    Write-Host "5. Test: servin run --vm alpine echo 'Hello from VM!'"
    
    Write-Host "`nServin GUI: Run 'servin-gui' to open the graphical interface"
    Write-Host "Documentation: See VM_PREREQUISITES.md for detailed setup guide"
    
    Write-Success "`nServin Container Runtime installation completed!"
}

# Main installation flow
try {
    Write-Header "Servin Container Runtime - Windows Installer"
    Write-Host "This installer will set up Servin with VM containerization support.`n"
    
    if ($Uninstall) {
        Write-Info "Uninstallation requested - this feature is not yet implemented"
        exit 0
    }
    
    # Check prerequisites
    $prerequisites = Test-Prerequisites
    $failedCount = ($prerequisites.Values | Where-Object { $_ -eq $false }).Count
    
    if ($failedCount -gt 0 -and -not $SkipPrerequisites) {
        Write-Warning "$failedCount prerequisites failed"
        if (-not $Silent) {
            $response = Read-Host "Continue anyway? (y/N)"
            if ($response -ne "y" -and $response -ne "Y") {
                exit 1
            }
        }
    }
    
    # Install VM prerequisites
    $vmCapabilities = @{}
    if (-not $SkipVMSetup) {
        $vmCapabilities = Install-VMPrerequisites
    }
    
    # Install Python dependencies
    Install-PythonDependencies | Out-Null
    
    # Install Servin
    if (-not (Install-Servin)) {
        Write-Error "Servin installation failed"
        exit 1
    }
    
    # Initialize VM support
    if (-not $SkipVMSetup) {
        Initialize-VM | Out-Null
    }
    
    # Show summary
    Show-InstallationSummary -Prerequisites $prerequisites -VMCapabilities $vmCapabilities
    
} catch {
    Write-Error "Installation failed: $($_.Exception.Message)"
    Write-Info "Check the error message above and try running the installer again"
    exit 1
}