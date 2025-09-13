@echo off
rem Servin GUI Launcher (Latest Build)
rem Sets up the CGO environment and runs the latest GUI build

echo Starting Servin Desktop GUI (Latest Build)...

rem Add MinGW to PATH for runtime (in case any dependencies need it)
set PATH=C:\msys64\ucrt64\bin;C:\msys64\usr\bin;%PATH%

rem Change to the directory containing this script
cd /d "%~dp0"

rem Run the latest GUI application
echo Working directory: %CD%
echo Starting latest GUI build with all fixes...
"%~dp0servin-gui-latest.exe"

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo GUI exited with error code: %ERRORLEVEL%
    echo.
)

pause
