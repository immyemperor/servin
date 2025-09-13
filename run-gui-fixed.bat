@echo off
rem Servin GUI Launcher (Fixed Version)
rem Sets up the CGO environment and runs the GUI with proper window controls

echo Starting Servin Desktop GUI (Fixed)...

rem Add MinGW to PATH for runtime (in case any dependencies need it)
set PATH=C:\msys64\ucrt64\bin;C:\msys64\usr\bin;%PATH%

rem Change to the directory containing this script
cd /d "%~dp0"

rem Run the fixed GUI application
echo Working directory: %CD%
echo Starting GUI with proper window controls...
"%~dp0servin-gui-fixed.exe"

pause
