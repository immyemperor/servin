@echo off
title Docker Desktop GUI
echo Starting Docker Desktop GUI...
echo.

REM Change to the docker_gui directory
cd /d "%~dp0"

REM Check if virtual environment exists
if not exist "..\\.venv\\Scripts\\python.exe" (
    echo Error: Python virtual environment not found.
    echo Please run the setup first.
    pause
    exit /b 1
)

REM Activate virtual environment and run the application
echo Activating Python environment...
call "..\\.venv\\Scripts\\activate.bat"

echo Starting application...
python main.py

echo.
echo Application has exited.
pause
