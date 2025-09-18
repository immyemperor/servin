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

echo [DEBUG] NSIS version information:
makensis /VERSION 2>nul || echo NSIS version check failed

REM Check if executables exist
if not exist "servin.exe" (
    echo [ERROR] servin.exe not found
    echo Please build the Servin executables first
    pause
    exit /b 1
)

if not exist "servin-tui.exe" (
    echo [ERROR] servin-tui.exe not found
    echo Please build the Servin TUI executable first
    pause
    exit /b 1
)

if not exist "servin-gui.exe" (
    echo [WARNING] servin-gui.exe not found - GUI components will not be included
    echo This is optional if building CLI-only version
) else (
    echo [INFO] servin-gui.exe found - GUI components will be included
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
echo [DEBUG] Current directory: %CD%
echo [DEBUG] Available files:
dir

echo [DEBUG] Running NSIS with command: makensis /NOCD servin-installer.nsi
echo [DEBUG] NSIS output:
makensis /NOCD servin-installer.nsi 2>&1 | tee build.log

echo [DEBUG] NSIS exit code: %errorlevel%
echo [DEBUG] Build log contents:
type build.log 2>nul || echo No build log found

if %errorlevel% eq 0 (
    echo.
    echo [SUCCESS] NSIS compilation completed successfully
    
    REM Check if the installer file was actually created
    if exist "servin-installer-1.0.0.exe" (
        echo [SUCCESS] Installer file created: servin-installer-1.0.0.exe
        dir /b *installer*.exe
        
        REM Copy to standardized name for workflow compatibility
        copy "servin-installer-1.0.0.exe" "Servin-Installer-1.0.0.exe"
        echo [INFO] Created standardized installer name: Servin-Installer-1.0.0.exe
        
        echo.
        echo The installer includes:
        echo - Servin Container Runtime executables
        echo - VM prerequisites installation (Chocolatey, Python, VM providers)
        echo - Windows Service integration
        echo - Desktop and Start Menu shortcuts
        echo - File associations and context menus
        echo.
        echo To test the installer, run as Administrator:
        echo Servin-Installer-1.0.0.exe
    ) else (
        echo [ERROR] NSIS compilation succeeded but installer file not found
        echo [DEBUG] Expected file: servin-installer-1.0.0.exe
        echo [DEBUG] Files matching pattern:
        dir /b *installer*.exe 2>nul || echo No installer files found
        echo [DEBUG] All .exe files:
        dir /b *.exe 2>nul || echo No exe files found
        exit /b 1
    )
) else (
    echo.
    echo [ERROR] NSIS compilation failed with exit code %errorlevel%
    echo [DEBUG] Build log contents:
    type build.log 2>nul || echo No build log available
    echo [DEBUG] Common NSIS issues:
    echo   - Missing required files (check servin.exe, servin-tui.exe)
    echo   - Syntax errors in .nsi file
    echo   - Missing NSIS plugins or includes
    echo   - Icon file format issues
    echo   - License file encoding problems
    pause
    exit /b 1
)

pause