# Servin GUI Launcher PowerShell Script
# Sets up the CGO environment and runs the GUI

Write-Host "Starting Servin Desktop GUI..." -ForegroundColor Green

# Get the script directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition

# Add MinGW to PATH for this session (in case any dependencies need it)
$env:PATH = "C:\msys64\ucrt64\bin;C:\msys64\usr\bin;" + $env:PATH

# Set working directory to script location so servin.exe can be found
Set-Location $ScriptDir

# Run the GUI application
$GuiPath = Join-Path $ScriptDir "servin-gui.exe"

if (Test-Path $GuiPath) {
    Write-Host "Launching GUI from: $GuiPath" -ForegroundColor Yellow
    Write-Host "Working directory: $ScriptDir" -ForegroundColor Cyan
    & $GuiPath
} else {
    Write-Host "GUI executable not found at: $GuiPath" -ForegroundColor Red
    Write-Host "Please build the GUI first with: go build -o servin-gui.exe ./cmd/servin-gui" -ForegroundColor Yellow
}

# Keep window open if run directly
if ($Host.Name -eq "ConsoleHost") {
    Write-Host "Press any key to continue..."
    $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
}
