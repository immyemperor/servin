@echo off
REM Windows batch script to serve the Servin wiki

echo 🚀 Servin Container Runtime Wiki Server
echo ==============================================

REM Check if Python is installed
python --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Error: Python is not installed or not in PATH
    echo    Please install Python 3 and try again
    pause
    exit /b 1
)

REM Check if wiki file exists
if not exist "wiki.html" (
    echo ❌ Error: wiki.html not found
    echo    Please run this script from the docs\ directory
    pause
    exit /b 1
)

echo 📚 Starting wiki server...
echo 🌐 Wiki will open in your browser automatically
echo 📝 Press Ctrl+C to stop the server
echo ==============================================

REM Start the Python server
python serve-wiki.py

pause
