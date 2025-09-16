# Servin Container Runtime - Windows PowerShell Installer
# Run this script as Administrator

param(
    [switch]$Uninstall,
    [string]$InstallDir = "C:\Program Files\Servin",
    [string]$DataDir = "C:\ProgramData\Servin"
)

$ErrorActionPreference = "Stop"

# Check if running as administrator
if (-NOT ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Error "This script must be run as Administrator. Right-click PowerShell and select 'Run as Administrator'"
    exit 1
}

function Write-ColoredText {
    param([string]$Text, [string]$Color = "Green")
    Write-Host $Text -ForegroundColor $Color
}

function Install-Servin {
    Write-ColoredText "================================================" "Cyan"
    Write-ColoredText "   Servin Container Runtime - PowerShell Installer" "Cyan"
    Write-ColoredText "================================================" "Cyan"
    Write-Host ""

    $LogDir = Join-Path $DataDir "logs"
    $ConfigDir = Join-Path $DataDir "config"
    $VolumesDir = Join-Path $DataDir "volumes"
    $ImagesDir = Join-Path $DataDir "images"

    Write-ColoredText "Installing Servin to: $InstallDir"
    Write-ColoredText "Data directory: $DataDir"

    # Create directories
    Write-ColoredText "Creating directories..."
    @($InstallDir, $DataDir, $LogDir, $ConfigDir, $VolumesDir, $ImagesDir) | ForEach-Object {
        if (-not (Test-Path $_)) {
            New-Item -ItemType Directory -Path $_ -Force | Out-Null
            Write-Host "  Created: $_"
        }
    }

    # Copy executables
    Write-ColoredText "Copying executables..."
    $executables = @("servin.exe", "servin-tui.exe", "servin-gui.exe")
    foreach ($exe in $executables) {
        if (Test-Path $exe) {
            Copy-Item $exe -Destination (Join-Path $InstallDir $exe) -Force
            Write-Host "  Copied: $exe"
        } else {
            Write-Warning "Executable not found: $exe"
        }
    }

    if (-not (Test-Path (Join-Path $InstallDir "servin.exe"))) {
        Write-Error "servin.exe not found. Please run this installer from the directory containing the executables."
    }

    # Set permissions
    Write-ColoredText "Setting permissions..."
    $acl = Get-Acl $InstallDir
    $accessRule = New-Object System.Security.AccessControl.FileSystemAccessRule("SYSTEM","FullControl","ContainerInherit,ObjectInherit","None","Allow")
    $acl.SetAccessRule($accessRule)
    $accessRule = New-Object System.Security.AccessControl.FileSystemAccessRule("Administrators","FullControl","ContainerInherit,ObjectInherit","None","Allow")
    $acl.SetAccessRule($accessRule)
    Set-Acl -Path $InstallDir -AclObject $acl

    # Add to PATH
    Write-ColoredText "Adding to PATH..."
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "Machine")
    if ($currentPath -notlike "*$InstallDir*") {
        [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$InstallDir", "Machine")
        Write-Host "  Added to system PATH"
    }

    # Create default configuration
    Write-ColoredText "Creating default configuration..."
    $configContent = @"
# Servin Configuration File
# Data directory
data_dir=$DataDir

# Log settings
log_level=info
log_file=$LogDir\servin.log

# Runtime settings
runtime=native

# Network settings
bridge_name=servin0

# CRI settings
cri_port=10250
cri_enabled=false
"@
    $configPath = Join-Path $ConfigDir "servin.conf"
    Set-Content -Path $configPath -Value $configContent
    Write-Host "  Created: $configPath"

    # Create Windows Service
    Write-ColoredText "Creating Windows Service..."
    
    # Create service wrapper
    $serviceScript = @"
@echo off
cd /d "$InstallDir"
"$InstallDir\servin.exe" daemon --config "$ConfigDir\servin.conf"
"@
    $servicePath = Join-Path $InstallDir "servin-service.bat"
    Set-Content -Path $servicePath -Value $serviceScript

    # Install service
    try {
        $service = Get-Service -Name "ServinRuntime" -ErrorAction SilentlyContinue
        if ($service) {
            Write-Host "  Removing existing service..."
            Stop-Service -Name "ServinRuntime" -Force -ErrorAction SilentlyContinue
            sc.exe delete "ServinRuntime" | Out-Null
        }

        $result = sc.exe create "ServinRuntime" binPath= "`"$servicePath`"" DisplayName= "Servin Container Runtime" Description= "Servin Container Runtime Service providing Docker-compatible container management" start= auto depend= "Tcpip"
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  Windows Service created successfully"
            
            # Set service to restart on failure
            sc.exe failure "ServinRuntime" reset= 86400 actions= restart/60000/restart/60000/restart/60000 | Out-Null
        } else {
            Write-Warning "Failed to create Windows Service. You can manually start Servin using: $InstallDir\servin.exe"
        }
    } catch {
        Write-Warning "Error creating service: $($_.Exception.Message)"
    }

    # Create uninstaller
    Write-ColoredText "Creating uninstaller..."
    $uninstallScript = @"
# Servin Uninstaller
Write-Host "Uninstalling Servin Container Runtime..." -ForegroundColor Yellow

# Stop and remove service
try {
    Stop-Service -Name "ServinRuntime" -Force -ErrorAction SilentlyContinue
    sc.exe delete "ServinRuntime" | Out-Null
    Write-Host "Service removed"
} catch { }

# Remove from PATH
`$currentPath = [Environment]::GetEnvironmentVariable("PATH", "Machine")
`$newPath = `$currentPath -replace [regex]::Escape(";$InstallDir"), ""
[Environment]::SetEnvironmentVariable("PATH", `$newPath, "Machine")

# Remove directories
Remove-Item -Path "$InstallDir" -Recurse -Force -ErrorAction SilentlyContinue
Remove-Item -Path "$DataDir" -Recurse -Force -ErrorAction SilentlyContinue

# Remove shortcuts
Remove-Item -Path "`$env:PUBLIC\Desktop\Servin GUI.lnk" -ErrorAction SilentlyContinue
Remove-Item -Path "`$env:ProgramData\Microsoft\Windows\Start Menu\Programs\Servin" -Recurse -Force -ErrorAction SilentlyContinue

Write-Host "Servin has been uninstalled." -ForegroundColor Green
Read-Host "Press Enter to continue"
"@
    $uninstallPath = Join-Path $InstallDir "uninstall.ps1"
    Set-Content -Path $uninstallPath -Value $uninstallScript

    # Create shortcuts
    Write-ColoredText "Creating shortcuts..."
    $WshShell = New-Object -comObject WScript.Shell

    # Desktop shortcut
    $Shortcut = $WshShell.CreateShortcut("$env:PUBLIC\Desktop\Servin GUI.lnk")
    $Shortcut.TargetPath = Join-Path $InstallDir "servin-gui.exe"
    $Shortcut.Description = "Servin Container Runtime GUI"
    $Shortcut.Save()

    # Start Menu entries
    $startMenuDir = "$env:ProgramData\Microsoft\Windows\Start Menu\Programs\Servin"
    if (-not (Test-Path $startMenuDir)) {
        New-Item -ItemType Directory -Path $startMenuDir -Force | Out-Null
    }

    $Shortcut = $WshShell.CreateShortcut("$startMenuDir\Servin GUI.lnk")
    $Shortcut.TargetPath = Join-Path $InstallDir "servin-gui.exe"
    $Shortcut.Description = "Servin Container Runtime GUI"
    $Shortcut.Save()

    $Shortcut = $WshShell.CreateShortcut("$startMenuDir\Uninstall Servin.lnk")
    $Shortcut.TargetPath = "powershell.exe"
    $Shortcut.Arguments = "-ExecutionPolicy Bypass -File `"$uninstallPath`""
    $Shortcut.Description = "Uninstall Servin Container Runtime"
    $Shortcut.Save()

    Write-Host ""
    Write-ColoredText "================================================" "Green"
    Write-ColoredText "   Installation completed successfully!" "Green"
    Write-ColoredText "================================================" "Green"
    Write-Host ""
    Write-Host "Installation directory: $InstallDir"
    Write-Host "Data directory: $DataDir"
    Write-Host "Configuration: $ConfigDir\servin.conf"
    Write-Host ""
    Write-ColoredText "Next steps:" "Yellow"
    Write-Host "1. Start the service: Start-Service 'ServinRuntime'"
    Write-Host "2. Or run manually: $InstallDir\servin.exe"
    Write-Host "3. Use GUI: $InstallDir\servin-gui.exe"
    Write-Host ""
    Write-Host "The service is configured to start automatically on boot."
    Write-Host "Check logs at: $LogDir\servin.log"
}

function Uninstall-Servin {
    Write-ColoredText "Uninstalling Servin Container Runtime..." "Yellow"
    
    # Execute the uninstaller if it exists
    $uninstallPath = Join-Path $InstallDir "uninstall.ps1"
    if (Test-Path $uninstallPath) {
        & $uninstallPath
    } else {
        Write-Warning "Uninstaller not found. Manual cleanup may be required."
    }
}

# Main execution
try {
    if ($Uninstall) {
        Uninstall-Servin
    } else {
        Install-Servin
    }
} catch {
    Write-Error "Installation failed: $($_.Exception.Message)"
    exit 1
}
