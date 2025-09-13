# Servin Cross-Platform Build and Package Script for Windows PowerShell

param(
    [string]$Target = "all",
    [string]$Version = "1.0.0"
)

$ErrorActionPreference = "Stop"

# Configuration
$BUILD_DIR = "build"
$DIST_DIR = "dist"
$INSTALLER_DIR = "installers"

function Write-ColoredText {
    param([string]$Text, [string]$Color = "Green")
    Write-Host $Text -ForegroundColor $Color
}

function Write-Header {
    Write-ColoredText "================================================" "Cyan"
    Write-ColoredText "   Servin Container Runtime - Build Script" "Cyan"
    Write-ColoredText "================================================" "Cyan"
    Write-Host ""
}

function Test-Prerequisites {
    Write-ColoredText "Checking prerequisites..."
    
    # Check Go installation
    try {
        $goVersion = go version
        Write-Host "  Go: $goVersion"
    } catch {
        Write-Error "Go is not installed or not in PATH"
    }
    
    # Check if we're in the right directory
    if (-not (Test-Path "go.mod")) {
        Write-Error "This script must be run from the project root directory"
    }
    
    Write-ColoredText "  Prerequisites OK" "Green"
}

function Clear-BuildDirectories {
    Write-ColoredText "Cleaning up previous builds..."
    
    if (Test-Path $BUILD_DIR) {
        Remove-Item -Path $BUILD_DIR -Recurse -Force
    }
    if (Test-Path $DIST_DIR) {
        Remove-Item -Path $DIST_DIR -Recurse -Force
    }
    
    New-Item -ItemType Directory -Path $BUILD_DIR -Force | Out-Null
    New-Item -ItemType Directory -Path $DIST_DIR -Force | Out-Null
}

function Build-Binaries {
    param(
        [string]$Platform,
        [string]$Arch,
        [string]$Extension = ""
    )
    
    Write-ColoredText "Building for $Platform/$Arch..."
    
    $outputDir = Join-Path $BUILD_DIR "$Platform-$Arch"
    New-Item -ItemType Directory -Path $outputDir -Force | Out-Null
    
    # Set environment variables for cross-compilation
    $env:GOOS = $Platform
    $env:GOARCH = $Arch
    
    # Build main servin binary
    $env:CGO_ENABLED = "0"
    $ldflags = "-w -s -X main.version=$Version"
    $servinOutput = Join-Path $outputDir "servin$Extension"
    
    go build -ldflags="$ldflags" -o $servinOutput .
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to build servin binary"
    }
    
    # Build GUI binary (requires CGO)
    if ($Platform -eq "windows") {
        try {
            $env:CGO_ENABLED = "1"
            $guiOutput = Join-Path $outputDir "servin-gui$Extension"
            go build -ldflags="$ldflags" -o $guiOutput ./cmd/servin-gui
            if ($LASTEXITCODE -eq 0) {
                Write-Host "  Built GUI binary"
            } else {
                Write-Warning "  Failed to build GUI binary (CGO compilation issue)"
            }
        } catch {
            Write-Warning "  Skipping GUI build - CGO compilation failed"
        }
    }
    
    Write-ColoredText "  Built binaries for $Platform/$Arch" "Green"
}

function New-WindowsPackage {
    Write-ColoredText "Creating Windows package..."
    
    $platformDir = Join-Path $BUILD_DIR "windows-amd64"
    $packageDir = Join-Path $DIST_DIR "servin-windows-$Version"
    
    New-Item -ItemType Directory -Path $packageDir -Force | Out-Null
    
    # Copy binaries
    Get-ChildItem -Path $platformDir -Filter "*.exe" | Copy-Item -Destination $packageDir
    
    # Copy installer scripts
    $installerSource = Join-Path $INSTALLER_DIR "windows\install.ps1"
    Copy-Item -Path $installerSource -Destination $packageDir
    
    # Copy NSIS installer files for wizard creation
    $nsisFiles = @("servin-installer.nsi", "LICENSE.txt", "servin.conf")
    foreach ($file in $nsisFiles) {
        $nsisFile = Join-Path $INSTALLER_DIR "windows\$file"
        if (Test-Path $nsisFile) {
            Copy-Item -Path $nsisFile -Destination $packageDir
        }
    }
    
    # Create README
    $readmeContent = @"
Servin Container Runtime for Windows
Version: $Version

Installation Options:

OPTION 1: Installation Wizard (Recommended)
1. Install NSIS from: https://nsis.sourceforge.io/Download
2. Run: .\build-installer.ps1 (creates ServinSetup-$Version.exe)
3. Run ServinSetup-$Version.exe as Administrator

OPTION 2: PowerShell Script
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
- If installed via wizard: Use Windows "Add or Remove Programs"
- If installed via script: Run C:\Program Files\Servin\uninstall.ps1

For more information, visit: https://github.com/yourusername/servin
"@
    
    Set-Content -Path (Join-Path $packageDir "README.txt") -Value $readmeContent
    
    # Create ZIP archive
    $zipPath = Join-Path $DIST_DIR "servin-windows-$Version.zip"
    try {
        Compress-Archive -Path $packageDir -DestinationPath $zipPath -Force
        Write-ColoredText "  Created: servin-windows-$Version.zip" "Green"
    } catch {
        Write-Warning "  Failed to create ZIP archive, directory package created instead"
    }
}

function Build-LinuxBinaries {
    Write-ColoredText "Building Linux binaries (cross-compile)..."
    
    # Linux build without GUI (no CGO)
    Build-Binaries "linux" "amd64" ""
    
    $platformDir = Join-Path $BUILD_DIR "linux-amd64"
    $packageDir = Join-Path $DIST_DIR "servin-linux-$Version"
    
    New-Item -ItemType Directory -Path $packageDir -Force | Out-Null
    
    # Copy binaries
    Copy-Item -Path (Join-Path $platformDir "servin") -Destination $packageDir
    
    # Copy installer
    $installerSource = Join-Path $INSTALLER_DIR "linux\install.sh"
    Copy-Item -Path $installerSource -Destination $packageDir
    
    # Create README
    $readmeContent = @"
# Servin Container Runtime for Linux
Version: $Version

## Installation
``````bash
sudo ./install.sh
``````

This will:
- Install Servin to /usr/local/bin
- Create a system user 'servin'
- Set up systemd service (or SysV init script)
- Create configuration in /etc/servin
- Set up data directories in /var/lib/servin

## Usage
- CLI: ``servin --help``
- Service: ``sudo systemctl start servin``

## Uninstallation
``````bash
sudo /usr/local/bin/servin-uninstall
``````

For more information, visit: https://github.com/yourusername/servin
"@
    
    Set-Content -Path (Join-Path $packageDir "README.md") -Value $readmeContent
    
    Write-ColoredText "  Created Linux package directory" "Green"
}

function Build-MacOSBinaries {
    Write-ColoredText "Building macOS binaries (cross-compile)..."
    
    # macOS build without GUI (no CGO)
    Build-Binaries "darwin" "amd64" ""
    
    $platformDir = Join-Path $BUILD_DIR "darwin-amd64"
    $packageDir = Join-Path $DIST_DIR "servin-macos-$Version"
    
    New-Item -ItemType Directory -Path $packageDir -Force | Out-Null
    
    # Copy binaries
    Copy-Item -Path (Join-Path $platformDir "servin") -Destination $packageDir
    
    # Copy installer
    $installerSource = Join-Path $INSTALLER_DIR "macos\install.sh"
    Copy-Item -Path $installerSource -Destination $packageDir
    
    # Create README
    $readmeContent = @"
# Servin Container Runtime for macOS
Version: $Version

## Installation
``````bash
sudo ./install.sh
``````

This will:
- Install Servin to /usr/local/bin
- Create a system user '_servin'
- Set up launchd service
- Create configuration in /usr/local/etc/servin
- Set up data directories in /usr/local/var/lib/servin

## Usage
- CLI: ``servin --help``
- Service: Starts automatically (launchd)

## Uninstallation
``````bash
sudo /usr/local/bin/servin-uninstall
``````

For more information, visit: https://github.com/yourusername/servin
"@
    
    Set-Content -Path (Join-Path $packageDir "README.md") -Value $readmeContent
    
    Write-ColoredText "  Created macOS package directory" "Green"
}

function Generate-Checksums {
    Write-ColoredText "Generating checksums..."
    
    $files = Get-ChildItem -Path $DIST_DIR -Filter "*.zip"
    if ($files.Count -gt 0) {
        $checksumFile = Join-Path $DIST_DIR "servin-$Version-checksums.txt"
        $checksums = @()
        
        foreach ($file in $files) {
            $hash = Get-FileHash -Path $file.FullName -Algorithm SHA256
            $checksums += "$($hash.Hash.ToLower())  $($file.Name)"
        }
        
        Set-Content -Path $checksumFile -Value $checksums
        Write-ColoredText "  Generated checksums" "Green"
    }
}

function Show-Summary {
    Write-Host ""
    Write-ColoredText "================================================" "Green"
    Write-ColoredText "   Build completed successfully!" "Green"
    Write-ColoredText "================================================" "Green"
    Write-Host ""
    Write-ColoredText "Built packages:" "Blue"
    
    Get-ChildItem -Path $DIST_DIR | ForEach-Object {
        Write-Host "  $($_.Name)"
    }
    
    Write-Host ""
    Write-ColoredText "Installation instructions:" "Blue"
    Write-Host "Windows: Extract ZIP, run install.ps1 as Administrator"
    Write-Host "Linux:   Extract and run sudo ./install.sh (on Linux system)"
    Write-Host "macOS:   Extract and run sudo ./install.sh (on macOS system)"
}

function Main {
    Write-Header
    Test-Prerequisites
    Clear-BuildDirectories
    
    switch ($Target.ToLower()) {
        "windows" {
            Build-Binaries "windows" "amd64" ".exe"
            New-WindowsPackage
        }
        "linux" {
            Build-LinuxBinaries
        }
        "macos" {
            Build-MacOSBinaries
        }
        "clean" {
            Write-ColoredText "Cleaned build and dist directories" "Green"
            return
        }
        default {
            # Build all platforms
            Build-Binaries "windows" "amd64" ".exe"
            New-WindowsPackage
            Build-LinuxBinaries
            Build-MacOSBinaries
        }
    }
    
    Generate-Checksums
    Show-Summary
}

# Execute main function
try {
    Main
} catch {
    Write-Error "Build failed: $($_.Exception.Message)"
    exit 1
} finally {
    # Reset environment variables
    Remove-Item Env:GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
    Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue
}
