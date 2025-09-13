# Build NSIS Installer for Servin Container Runtime

param(
    [string]$Version = "1.0.0",
    [string]$NSISPath = "C:\Program Files (x86)\NSIS\makensis.exe"
)

$ErrorActionPreference = "Stop"

function Write-ColoredText {
    param([string]$Text, [string]$Color = "Green")
    Write-Host $Text -ForegroundColor $Color
}

function Test-NSISInstalled {
    if (Test-Path $NSISPath) {
        return $true
    }
    
    # Try common NSIS locations
    $commonPaths = @(
        "C:\Program Files (x86)\NSIS\makensis.exe",
        "C:\Program Files\NSIS\makensis.exe",
        "${env:ProgramFiles(x86)}\NSIS\makensis.exe",
        "$env:ProgramFiles\NSIS\makensis.exe"
    )
    
    foreach ($path in $commonPaths) {
        if (Test-Path $path) {
            $script:NSISPath = $path
            return $true
        }
    }
    
    return $false
}

Write-ColoredText "================================================" "Cyan"
Write-ColoredText "   Building Servin NSIS Installer Wizard" "Cyan"
Write-ColoredText "================================================" "Cyan"
Write-Host ""

# Check if NSIS is installed
if (-not (Test-NSISInstalled)) {
    Write-ColoredText "NSIS (Nullsoft Scriptable Install System) is not installed!" "Red"
    Write-Host ""
    Write-Host "To build the installer wizard, please install NSIS:"
    Write-Host "1. Download from: https://nsis.sourceforge.io/Download"
    Write-Host "2. Install the latest version"
    Write-Host "3. Re-run this script"
    Write-Host ""
    Write-Host "Alternative: Use the PowerShell installer script instead:"
    Write-Host "   .\install.ps1"
    exit 1
}

Write-ColoredText "Found NSIS at: $NSISPath" "Green"

# Ensure we're in the right directory
$installerDir = "installers\windows"
if (-not (Test-Path $installerDir)) {
    Write-Error "Please run this script from the project root directory"
}

# Build the main executables first
Write-ColoredText "Building executables..." "Yellow"
try {
    .\build.ps1 -Target windows | Out-Null
    Write-ColoredText "  Executables built successfully" "Green"
} catch {
    Write-Error "Failed to build executables: $($_.Exception.Message)"
}

# Copy files to installer directory
Write-ColoredText "Preparing installer files..." "Yellow"

$sourceDir = "build\windows-amd64"
$files = @(
    @{Source = "$sourceDir\servin.exe"; Dest = "$installerDir\servin.exe"},
    @{Source = "$sourceDir\servin-gui.exe"; Dest = "$installerDir\servin-gui.exe"},
    @{Source = "icons\favicon.ico"; Dest = "$installerDir\servin.ico"}
)

foreach ($file in $files) {
    if (Test-Path $file.Source) {
        Copy-Item $file.Source $file.Dest -Force
        Write-Host "  Copied: $($file.Source)"
    } else {
        Write-Warning "  Missing: $($file.Source)"
    }
}

# Create README for installer
$readmeContent = @"
Servin Container Runtime - Installation Complete!

Thank you for installing Servin Container Runtime v$Version.

What's installed:
- Servin Container Runtime (servin.exe)
- Servin Desktop GUI (servin-gui.exe)
- Windows Service (ServinRuntime)
- Configuration files and data directories

Quick Start:
1. The Servin service is now running in the background
2. Launch the GUI from the Start Menu or desktop shortcut
3. Use 'servin --help' in Command Prompt for CLI usage

Configuration:
- Config file: C:\ProgramData\Servin\config\servin.conf
- Data directory: C:\ProgramData\Servin\data
- Log files: C:\ProgramData\Servin\logs

Documentation:
- Visit: https://github.com/yourusername/servin
- Support: https://github.com/yourusername/servin/issues

Enjoy containerizing with Servin!
"@

Set-Content -Path "$installerDir\README.txt" -Value $readmeContent

# Build the installer
Write-ColoredText "Building NSIS installer..." "Yellow"

Push-Location $installerDir
try {
    & $NSISPath "servin-installer.nsi"
    if ($LASTEXITCODE -eq 0) {
        Write-ColoredText "  Installer built successfully!" "Green"
        
        # Move installer to dist directory
        $installerFile = "ServinSetup-$Version.exe"
        if (Test-Path $installerFile) {
            Move-Item $installerFile "..\..\dist\$installerFile" -Force
            Write-ColoredText "  Installer saved to: dist\$installerFile" "Green"
        }
    } else {
        Write-Error "NSIS compilation failed with exit code $LASTEXITCODE"
    }
} finally {
    Pop-Location
}

Write-Host ""
Write-ColoredText "================================================" "Green"
Write-ColoredText "   Installer Wizard Build Complete!" "Green"
Write-ColoredText "================================================" "Green"
Write-Host ""
Write-Host "Installer wizard: dist\ServinSetup-$Version.exe"
Write-Host ""
Write-Host "To install Servin:"
Write-Host "1. Run ServinSetup-$Version.exe as Administrator"
Write-Host "2. Follow the installation wizard"
Write-Host "3. Launch Servin GUI from Start Menu or desktop"
