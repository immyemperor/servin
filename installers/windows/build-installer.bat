@echo off
REM Build script for Servin NSIS Installer
REM Requires: NSIS installed, servin executables built

echo.
echo Building Servin Container Runtime Installer...
echo.

REM Check if NSIS is installed
where makensis >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] NSIS (makensis) not found in PATH
    echo Please install NSIS from https://nsis.sourceforge.io/
    echo Or install via Chocolatey: choco install nsis
    pause
    exit /b 1
)

REM Check if executables exist
if not exist "servin.exe" (
    echo [ERROR] servin.exe not found
    echo Please build the Servin executables first
    pause
    exit /b 1
)

if not exist "servin-gui.exe" (
    echo [ERROR] servin-gui.exe not found
    echo Please build the Servin GUI executable first
    pause
    exit /b 1
)

if not exist "servin-tui.exe" (
    echo [ERROR] servin-tui.exe not found
    echo Please build the Servin TUI executable first
    pause
    exit /b 1
)

REM Check for required files
if not exist "LICENSE.txt" (
    echo [WARNING] LICENSE.txt not found, creating placeholder
    echo Servin Container Runtime License > LICENSE.txt
    echo Please replace with actual license file >> LICENSE.txt
)

if not exist "servin.conf" (
    echo [WARNING] servin.conf not found, creating default
    echo # Servin Configuration > servin.conf
    echo data_dir=%APPDATA%\Servin\data >> servin.conf
    echo log_level=info >> servin.conf
    echo vm_enabled=true >> servin.conf
)

if not exist "servin.ico" (
    echo [WARNING] servin.ico not found, using default Windows icon
)

echo [INFO] Building installer with NSIS...
makensis /NOCD servin-installer.nsi

if %errorlevel% eq 0 (
    echo.
    echo [SUCCESS] Installer built successfully: servin-installer-1.0.0.exe
    echo.
    echo The installer includes:
    echo - Servin Container Runtime executables
    echo - VM prerequisites installation (Chocolatey, Python, VM providers)
    echo - Windows Service integration
    echo - Desktop and Start Menu shortcuts
    echo - File associations and context menus
    echo.
    echo To test the installer, run as Administrator:
    echo servin-installer-1.0.0.exe
) else (
    echo.
    echo [ERROR] Installer build failed
    echo Check the output above for errors
    pause
    exit /b 1
)

pause