# Servin Container Runtime - Smart Installation Wizard for Windows
# Auto-detects platform and installs prerequisites if needed
# Run as Administrator: powershell -ExecutionPolicy Bypass -File install-wizard.ps1

param(
    [switch]$Help,
    [switch]$Auto,
    [switch]$VmOnly,
    [switch]$NoVm,
    [switch]$Force
)

# Configuration
$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

# Colors for console output
function Write-Success { param($Message) Write-Host "✓ $Message" -ForegroundColor Green }
function Write-Warning { param($Message) Write-Host "⚠ $Message" -ForegroundColor Yellow }
function Write-Error { param($Message) Write-Host "✗ $Message" -ForegroundColor Red }
function Write-Info { param($Message) Write-Host "→ $Message" -ForegroundColor Blue }
function Write-Header { param($Message) Write-Host "`n$Message" -ForegroundColor Cyan -BackgroundColor DarkBlue }

function Show-Banner {
    Write-Host "`n" -ForegroundColor Cyan
    Write-Host "╔════════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
    Write-Host "║                    Servin Container Runtime                    ║" -ForegroundColor Cyan
    Write-Host "║                  Smart Installation Wizard                    ║" -ForegroundColor Cyan
    Write-Host "╚════════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
    Write-Host ""
}

# Check if running as Administrator
function Test-Administrator {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

# Detect Windows version and architecture
function Get-SystemInfo {
    Write-Header "Detecting System Information"
    
    $osInfo = Get-CimInstance Win32_OperatingSystem
    $procInfo = Get-CimInstance Win32_Processor
    
    $script:WindowsVersion = $osInfo.Version
    $script:WindowsBuild = $osInfo.BuildNumber
    $script:Architecture = $procInfo.Architecture
    $script:ProcessorName = $procInfo.Name
    
    Write-Info "Windows Version: $($osInfo.Caption) (Build $WindowsBuild)"
    Write-Info "Architecture: $(if ($Architecture -eq 9) { 'x64' } elseif ($Architecture -eq 12) { 'ARM64' } else { 'x86' })"
    Write-Info "Processor: $ProcessorName"
    
    # Check minimum Windows version (Windows 10 1903 / Build 18362)
    if ($WindowsBuild -lt 18362) {
        Write-Error "Windows 10 version 1903 (Build 18362) or later is required"
        Write-Info "Current build: $WindowsBuild"
        exit 1
    } else {
        Write-Success "Windows version supported"
    }
}

# Check current directory and find enhanced installer
function Test-InstallerDirectory {
    Write-Header "Validating Installation Directory"
    
    $currentDir = Get-Location
    $readmePath = Join-Path $currentDir "README.md"
    
    if (-not (Test-Path $readmePath)) {
        Write-Error "This script must be run from the Servin installers directory"
        Write-Info "Expected to find installers/README.md"
        exit 1
    }
    
    $readmeContent = Get-Content $readmePath -Raw
    if ($readmeContent -notmatch "Servin Container Runtime") {
        Write-Error "Invalid installer directory - README.md doesn't contain Servin content"
        exit 1
    }
    
    $script:InstallerDir = $currentDir
    $script:EnhancedInstaller = Join-Path $InstallerDir "windows\install-with-vm.ps1"
    
    if (-not (Test-Path $EnhancedInstaller)) {
        Write-Error "Enhanced installer not found: $EnhancedInstaller"
        Write-Info "Please ensure you have the complete Servin installer package"
        exit 1
    }
    
    Write-Success "Running from correct directory: $InstallerDir"
    Write-Success "Enhanced installer found: $EnhancedInstaller"
}

# Check for existing Servin installation
function Test-ExistingInstallation {
    Write-Header "Checking Existing Installation"
    
    $servinPaths = @(
        "${env:ProgramFiles}\Servin\servin.exe",
        "${env:ProgramFiles(x86)}\Servin\servin.exe",
        "${env:LOCALAPPDATA}\Servin\servin.exe"
    )
    
    $script:ExistingServin = $null
    
    foreach ($path in $servinPaths) {
        if (Test-Path $path) {
            Write-Success "Found existing Servin installation: $path"
            try {
                $version = & $path version 2>$null | Select-Object -First 1
                Write-Info "Version: $version"
            } catch {
                Write-Info "Version: Unable to determine"
            }
            $script:ExistingServin = $path
            break
        }
    }
    
    if (-not $ExistingServin) {
        Write-Info "No existing Servin installation found"
    }
    
    # Check if Servin is in PATH
    try {
        $pathServin = Get-Command servin -ErrorAction SilentlyContinue
        if ($pathServin) {
            Write-Info "Servin found in PATH: $($pathServin.Source)"
        }
    } catch { }
}

# Check VM prerequisites
function Test-VmPrerequisites {
    Write-Header "Checking VM Prerequisites"
    
    $script:PrereqMissing = $false
    
    # Check Hyper-V
    Write-Info "Checking Hyper-V availability..."
    try {
        $hyperv = Get-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V-All -ErrorAction SilentlyContinue
        if ($hyperv -and $hyperv.State -eq "Enabled") {
            Write-Success "Hyper-V is enabled"
            $script:HyperVAvailable = $true
        } else {
            Write-Warning "Hyper-V not enabled"
            $script:HyperVAvailable = $false
            $script:PrereqMissing = $true
        }
    } catch {
        Write-Warning "Could not check Hyper-V status"
        $script:HyperVAvailable = $false
        $script:PrereqMissing = $true
    }
    
    # Check VirtualBox
    Write-Info "Checking VirtualBox..."
    try {
        $vboxPath = "${env:ProgramFiles}\Oracle\VirtualBox\VBoxManage.exe"
        if (Test-Path $vboxPath) {
            $vboxVersion = & $vboxPath --version 2>$null
            Write-Success "VirtualBox available: $vboxVersion"
            $script:VirtualBoxAvailable = $true
        } else {
            Write-Warning "VirtualBox not found"
            $script:VirtualBoxAvailable = $false
            $script:PrereqMissing = $true
        }
    } catch {
        Write-Warning "VirtualBox not accessible"
        $script:VirtualBoxAvailable = $false
        $script:PrereqMissing = $true
    }
    
    # Check WSL2
    Write-Info "Checking WSL2..."
    try {
        $wslStatus = wsl --status 2>$null
        if ($LASTEXITCODE -eq 0) {
            Write-Success "WSL2 available"
            $script:WSL2Available = $true
        } else {
            Write-Warning "WSL2 not available"
            $script:WSL2Available = $false
            $script:PrereqMissing = $true
        }
    } catch {
        Write-Warning "WSL2 not accessible"
        $script:WSL2Available = $false
        $script:PrereqMissing = $true
    }
    
    # Check Chocolatey
    Write-Info "Checking Chocolatey..."
    try {
        $chocoVersion = choco --version 2>$null
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Chocolatey available: $chocoVersion"
            $script:ChocolateyAvailable = $true
        } else {
            Write-Warning "Chocolatey not found"
            $script:ChocolateyAvailable = $false
            $script:PrereqMissing = $true
        }
    } catch {
        Write-Warning "Chocolatey not accessible"
        $script:ChocolateyAvailable = $false
        $script:PrereqMissing = $true
    }
    
    # Check Python
    Write-Info "Checking Python..."
    try {
        $pythonVersion = python --version 2>$null
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Python available: $pythonVersion"
            $script:PythonAvailable = $true
        } else {
            Write-Warning "Python not found"
            $script:PythonAvailable = $false
            $script:PrereqMissing = $true
        }
    } catch {
        Write-Warning "Python not accessible"
        $script:PythonAvailable = $false
        $script:PrereqMissing = $true
    }
    
    # Summary
    if ($PrereqMissing) {
        Write-Warning "Some VM prerequisites are missing and will be installed"
    } else {
        Write-Success "All VM prerequisites are available"
    }
}

# Check system resources
function Test-SystemResources {
    Write-Header "Checking System Resources"
    
    $script:ResourceWarnings = 0
    
    # Check memory
    Write-Info "Checking available memory..."
    try {
        $memory = Get-CimInstance Win32_ComputerSystem
        $memoryGB = [math]::Round($memory.TotalPhysicalMemory / 1GB, 1)
        
        if ($memoryGB -ge 16) {
            Write-Success "Memory: ${memoryGB}GB (excellent)"
        } elseif ($memoryGB -ge 8) {
            Write-Success "Memory: ${memoryGB}GB (good)"
        } elseif ($memoryGB -ge 4) {
            Write-Warning "Memory: ${memoryGB}GB (minimum met, 8GB+ recommended)"
            $script:ResourceWarnings++
        } else {
            Write-Error "Memory: ${memoryGB}GB (insufficient, 4GB minimum required)"
            $script:ResourceWarnings++
        }
    } catch {
        Write-Warning "Could not check memory"
        $script:ResourceWarnings++
    }
    
    # Check disk space
    Write-Info "Checking available disk space..."
    try {
        $systemDrive = Get-CimInstance Win32_LogicalDisk | Where-Object { $_.DeviceID -eq $env:SystemDrive }
        $freeSpaceGB = [math]::Round($systemDrive.FreeSpace / 1GB, 1)
        
        if ($freeSpaceGB -ge 20) {
            Write-Success "Disk space: ${freeSpaceGB}GB free (excellent)"
        } elseif ($freeSpaceGB -ge 10) {
            Write-Success "Disk space: ${freeSpaceGB}GB free (good)"
        } elseif ($freeSpaceGB -ge 5) {
            Write-Warning "Disk space: ${freeSpaceGB}GB free (minimum met, 10GB+ recommended)"
            $script:ResourceWarnings++
        } else {
            Write-Error "Disk space: ${freeSpaceGB}GB free (insufficient, 5GB minimum required)"
            $script:ResourceWarnings++
        }
    } catch {
        Write-Warning "Could not check disk space"
        $script:ResourceWarnings++
    }
    
    # Check CPU virtualization
    Write-Info "Checking CPU virtualization support..."
    try {
        $cpu = Get-CimInstance Win32_Processor
        $vmFeatures = $cpu.VirtualizationFirmwareEnabled
        
        if ($vmFeatures) {
            Write-Success "CPU virtualization supported and enabled"
        } else {
            Write-Warning "CPU virtualization not enabled in BIOS/UEFI"
            $script:ResourceWarnings++
        }
    } catch {
        Write-Warning "Could not check CPU virtualization features"
        $script:ResourceWarnings++
    }
}

# Show installation options
function Show-InstallationOptions {
    Write-Header "Installation Options"
    
    Write-Host "`nWhat would you like to do?`n" -ForegroundColor White
    
    if ($ExistingServin) {
        Write-Host "1) Update existing installation (preserve configuration)"
        Write-Host "2) Fresh installation (reset configuration)"
        Write-Host "3) Install VM prerequisites only"
        Write-Host "4) Exit"
        Write-Host ""
        
        if (-not $Auto) {
            $script:InstallOption = Read-Host "Choose option (1-4)"
        } else {
            $script:InstallOption = "1"
            Write-Info "Auto mode: Selected option 1 (Update existing installation)"
        }
    } else {
        Write-Host "1) Full installation with VM prerequisites"
        Write-Host "2) Basic installation (skip VM setup)"
        Write-Host "3) Install VM prerequisites only"
        Write-Host "4) Exit"
        Write-Host ""
        
        if (-not $Auto) {
            $script:InstallOption = Read-Host "Choose option (1-4)"
        } else {
            if ($VmOnly) {
                $script:InstallOption = "3"
                Write-Info "Auto mode: Selected option 3 (VM prerequisites only)"
            } elseif ($NoVm) {
                $script:InstallOption = "2"
                Write-Info "Auto mode: Selected option 2 (Basic installation)"
            } else {
                $script:InstallOption = "1"
                Write-Info "Auto mode: Selected option 1 (Full installation)"
            }
        }
    }
}

# Confirm installation
function Confirm-Installation {
    Write-Header "Installation Summary"
    
    Write-Host "`nInstallation Details:" -ForegroundColor White
    Write-Host "• Platform: Windows $WindowsBuild"
    Write-Host "• Enhanced installer: $(Split-Path $EnhancedInstaller -Leaf)"
    
    switch ($InstallOption) {
        "1" {
            if ($ExistingServin) {
                Write-Host "• Action: Update existing installation"
            } else {
                Write-Host "• Action: Full installation with VM prerequisites"
            }
        }
        "2" {
            if ($ExistingServin) {
                Write-Host "• Action: Fresh installation"
            } else {
                Write-Host "• Action: Basic installation (no VM setup)"
            }
        }
        "3" {
            Write-Host "• Action: Install VM prerequisites only"
        }
    }
    
    if ($PrereqMissing) {
        Write-Host "• VM Prerequisites: Will be installed"
    } else {
        Write-Host "• VM Prerequisites: Already available"
    }
    
    if ($ResourceWarnings -gt 0) {
        Write-Host "• Warnings: $ResourceWarnings resource warnings detected" -ForegroundColor Yellow
    }
    
    Write-Host ""
    
    if (-not $Auto -and -not $Force) {
        $continue = Read-Host "Continue with installation? (y/N)"
        if ($continue -notmatch "^[Yy]") {
            Write-Info "Installation cancelled"
            exit 0
        }
    } else {
        Write-Info "Auto mode: Proceeding with installation"
    }
}

# Run enhanced installer
function Invoke-EnhancedInstaller {
    Write-Header "Running Enhanced Installer"
    
    $installerArgs = @()
    
    switch ($InstallOption) {
        "2" {
            if (-not $ExistingServin) {
                $installerArgs += "-SkipVM"
            }
        }
        "3" {
            Write-Info "Note: This will install VM prerequisites along with Servin"
        }
    }
    
    Write-Info "Executing: $EnhancedInstaller $($installerArgs -join ' ')"
    Write-Info "This may take several minutes..."
    Write-Host ""
    
    try {
        # Run the enhanced installer
        $startInfo = New-Object System.Diagnostics.ProcessStartInfo
        $startInfo.FileName = "powershell.exe"
        $startInfo.Arguments = "-ExecutionPolicy Bypass -File `"$EnhancedInstaller`" $($installerArgs -join ' ')"
        $startInfo.RedirectStandardOutput = $false
        $startInfo.RedirectStandardError = $false
        $startInfo.UseShellExecute = $false
        $startInfo.CreateNoWindow = $false
        
        $process = New-Object System.Diagnostics.Process
        $process.StartInfo = $startInfo
        $process.Start()
        $process.WaitForExit()
        
        $exitCode = $process.ExitCode
        
        if ($exitCode -eq 0) {
            Write-Success "Enhanced installer completed successfully!"
        } else {
            Write-Error "Enhanced installer failed with exit code: $exitCode"
            throw "Installation failed"
        }
    } catch {
        Write-Error "Failed to run enhanced installer: $_"
        throw
    }
}

# Post-installation verification
function Test-PostInstallation {
    Write-Header "Post-Installation Verification"
    
    # Check if Servin is now available
    $servinPaths = @(
        "${env:ProgramFiles}\Servin\servin.exe",
        "${env:ProgramFiles(x86)}\Servin\servin.exe",
        "${env:LOCALAPPDATA}\Servin\servin.exe"
    )
    
    $servinCmd = $null
    foreach ($path in $servinPaths) {
        if (Test-Path $path) {
            $servinCmd = $path
            break
        }
    }
    
    if (-not $servinCmd) {
        Write-Warning "Servin executable not found in standard locations"
        Write-Info "You may need to restart your PowerShell session"
        return $false
    }
    
    # Test basic functionality
    Write-Info "Testing Servin CLI..."
    try {
        $version = & $servinCmd version 2>$null | Select-Object -First 1
        Write-Success "Servin CLI working: $version"
    } catch {
        Write-Warning "Servin CLI test failed"
        return $false
    }
    
    # Test VM functionality if installed
    if ($InstallOption -ne "2") {
        Write-Info "Testing VM functionality..."
        try {
            & $servinCmd vm status 2>$null | Out-Null
            if ($LASTEXITCODE -eq 0) {
                Write-Success "VM subsystem working"
            } else {
                Write-Warning "VM subsystem not fully functional"
                Write-Info "This may be normal on first run - try: $servinCmd vm init"
            }
        } catch {
            Write-Warning "VM subsystem test failed"
        }
    }
    
    return $true
}

# Show completion message
function Show-Completion {
    Write-Header "Installation Complete!"
    
    Write-Host "`n🎉 Servin Container Runtime has been successfully installed!`n" -ForegroundColor Green
    
    Write-Host "Next Steps:" -ForegroundColor White
    
    Write-Host "1. Restart your PowerShell session (for PATH updates)"
    Write-Host "2. Initialize VM support: servin vm init"
    Write-Host "3. Test installation: servin run --vm alpine echo 'Hello!'"
    
    Write-Host "`nAvailable Commands:" -ForegroundColor White
    Write-Host "• servin version       - Show version information"
    Write-Host "• servin vm status     - Check VM subsystem status"
    Write-Host "• servin-gui          - Launch graphical interface"
    Write-Host "• servin-tui          - Launch terminal interface"
    
    Write-Host "`nService Management:" -ForegroundColor White
    Write-Host "• Start: Start-Service ServinRuntime"
    Write-Host "• Stop: Stop-Service ServinRuntime"
    Write-Host "• Status: Get-Service ServinRuntime"
    
    Write-Host "`nDocumentation:" -ForegroundColor White
    Write-Host "• Complete guide: installers\VM_PREREQUISITES.md"
    Write-Host "• CLI reference: docs\cli.md"
    Write-Host "• Troubleshooting: docs\troubleshooting.md"
    
    Write-Host ""
    Write-Success "Installation wizard completed successfully!"
}

# Show help
function Show-Help {
    Write-Host "Servin Container Runtime - Smart Installation Wizard for Windows"
    Write-Host ""
    Write-Host "Usage: .\install-wizard.ps1 [options]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -Help          Show this help message"
    Write-Host "  -Auto          Run automated installation (no prompts)"
    Write-Host "  -VmOnly        Install VM prerequisites only"
    Write-Host "  -NoVm          Skip VM setup"
    Write-Host "  -Force         Skip confirmation prompts"
    Write-Host ""
    Write-Host "This wizard will:"
    Write-Host "1. Detect your Windows version and architecture"
    Write-Host "2. Check for existing installations"
    Write-Host "3. Verify system prerequisites"
    Write-Host "4. Run the enhanced Windows installer"
    Write-Host "5. Verify the installation"
    Write-Host ""
    Write-Host "Enhanced installer used:"
    Write-Host "• Windows: windows\install-with-vm.ps1"
    Write-Host ""
    Write-Host "Requirements:"
    Write-Host "• Windows 10 version 1903+ or Windows 11"
    Write-Host "• Administrator privileges"
    Write-Host "• PowerShell 5.1+ or PowerShell Core"
    Write-Host "• Internet connection for downloads"
}

# Main execution flow
function Main {
    try {
        Show-Banner
        
        Write-Info "Starting Servin Container Runtime installation wizard..."
        Write-Host ""
        
        Get-SystemInfo
        Test-InstallerDirectory
        Test-ExistingInstallation
        Test-VmPrerequisites
        Test-SystemResources
        Show-InstallationOptions
        
        switch ($InstallOption) {
            { $_ -in @("1", "2", "3") } {
                Confirm-Installation
                Invoke-EnhancedInstaller
                Test-PostInstallation
                Show-Completion
            }
            "4" {
                Write-Info "Installation cancelled"
                exit 0
            }
            default {
                Write-Error "Invalid option selected: $InstallOption"
                exit 1
            }
        }
    } catch {
        Write-Error "Installation wizard failed: $_"
        Write-Info "Check the error message above and try again"
        exit 1
    }
}

# Handle command line arguments
if ($Help) {
    Show-Help
    exit 0
}

# Check administrator privileges
if (-not (Test-Administrator)) {
    Write-Error "This script must be run as Administrator"
    Write-Info "Right-click PowerShell and select 'Run as Administrator'"
    exit 1
}

# Run main installation
Main