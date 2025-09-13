@echo off
rem Servin GUI Launcher
rem Sets up the CGO environment and runs the GUI

echo Starting Servin Desktop GUI...

rem Add MinGW to PATH for runtime (in case any dependencies need it)
set PATH=C:\msys64\ucrt64\bin;C:\msys64\usr\bin;%PATH%

rem Change to the directory containing this script
cd /d "%~dp0"

rem Run the GUI application
echo Working directory: %CD%
"%~dp0servin-gui.exe"

pause
