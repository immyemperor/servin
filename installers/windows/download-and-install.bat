@echo off
REM Servin Container Runtime - Windows Wizard Installer Downloader
REM This batch file downloads and runs the smart installer wizard

echo.
echo ╔════════════════════════════════════════════════════════════════╗
echo ║                    Servin Container Runtime                    ║
echo ║                   Windows Installer Download                  ║
echo ╚════════════════════════════════════════════════════════════════╝
echo.

REM Check if running as Administrator
net session >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] This script must be run as Administrator
    echo Right-click on Command Prompt and select "Run as Administrator"
    pause
    exit /b 1
)

echo [INFO] Downloading Servin installer wizard...

REM Download the wizard installer
powershell -Command "& {Invoke-WebRequest -Uri 'https://github.com/immyemperor/servin/releases/latest/download/install-wizard.ps1' -OutFile 'install-wizard.ps1'}"

if %errorlevel% neq 0 (
    echo [ERROR] Failed to download installer wizard
    echo Please check your internet connection and try again
    pause
    exit /b 1
)

echo [SUCCESS] Download completed: install-wizard.ps1

echo.
echo [INFO] Starting installation wizard...
echo This will install Servin Container Runtime with VM prerequisites

REM Run the wizard installer
powershell -ExecutionPolicy Bypass -File install-wizard.ps1

echo.
echo [INFO] Installation wizard completed
pause