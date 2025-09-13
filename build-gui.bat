@echo off
echo Building Servin GUI...
cd /d "C:\Users\immye\Desktop\servin\gocontainer"

rem Set up CGO environment
set PATH=C:\msys64\ucrt64\bin;C:\msys64\usr\bin;%PATH%

rem Clean and build
echo Cleaning build cache...
go clean -cache

echo Building GUI application...
go build -v -o servin-gui-fresh.exe ./cmd/servin-gui

if exist "servin-gui-fresh.exe" (
    echo ✅ Build successful!
    dir servin-gui-fresh.exe
) else (
    echo ❌ Build failed!
)

pause
