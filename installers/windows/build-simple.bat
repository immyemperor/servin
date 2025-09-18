@echo off
REM Simple Windows installer build script for debugging
echo =====================================
echo Simple NSIS Build Debug Script
echo =====================================

REM Basic checks
echo Checking required files...
if not exist "servin.exe" (
    echo ERROR: servin.exe not found
    exit /b 1
) else (
    echo OK: servin.exe found
)

if not exist "servin-tui.exe" (
    echo ERROR: servin-tui.exe not found
    exit /b 1
) else (
    echo OK: servin-tui.exe found
)

if exist "servin-gui.exe" (
    echo OK: servin-gui.exe found
) else (
    echo INFO: servin-gui.exe not found (optional)
)

if exist "LICENSE.txt" (
    echo OK: LICENSE.txt found
) else (
    echo INFO: LICENSE.txt not found (will be handled in script)
)

if exist "servin.ico" (
    echo OK: servin.ico found  
) else (
    echo INFO: servin.ico not found (will be handled in script)
)

REM Check NSIS
echo.
echo Checking NSIS installation...
where makensis >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: makensis not found in PATH
    exit /b 1
) else (
    echo OK: makensis found
    makensis /VERSION
)

REM Try NSIS compilation with detailed output
echo.
echo Running makensis with detailed output...
echo Command: makensis /V4 servin-installer.nsi
makensis /V4 servin-installer.nsi

echo.
echo makensis exit code: %errorlevel%
echo.

if exist "servin-installer-1.0.0.exe" (
    echo SUCCESS: Installer created - servin-installer-1.0.0.exe
    dir servin-installer*.exe
    
    REM Copy to standardized name
    copy "servin-installer-1.0.0.exe" "Servin-Installer-1.0.0.exe" >nul
    echo Copied to Servin-Installer-1.0.0.exe
) else (
    echo FAILED: Installer not created
    echo.
    echo Checking for any .exe files:
    dir *.exe 2>nul || echo No .exe files found
    exit /b 1
)