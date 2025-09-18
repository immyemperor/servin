@echo off
REM Simple Windows installer build script for debugging
echo Building Servin Installer...

REM Basic checks
if not exist "servin.exe" (
    echo ERROR: servin.exe not found
    exit /b 1
)

if not exist "servin-tui.exe" (
    echo ERROR: servin-tui.exe not found
    exit /b 1
)

echo Files present:
echo - servin.exe: %~z1 bytes
echo - servin-tui.exe exists
if exist "servin-gui.exe" echo - servin-gui.exe exists

REM Try NSIS compilation
echo Running makensis...
makensis servin-installer.nsi

if exist "servin-installer-1.0.0.exe" (
    echo SUCCESS: Installer created
    copy "servin-installer-1.0.0.exe" "Servin-Installer-1.0.0.exe"
    echo Standardized installer name created
) else (
    echo FAILED: Installer not created
    exit /b 1
)