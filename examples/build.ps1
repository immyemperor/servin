# Servin Build Script for Windows Development

# This PowerShell script helps build Servin on Windows for Linux deployment

Write-Host "=== Servin Build Script ===" -ForegroundColor Green

# Check Go installation
$goVersion = go version
if ($LASTEXITCODE -ne 0) {
    Write-Host "Error: Go is not installed or not in PATH" -ForegroundColor Red
    exit 1
}
Write-Host "Go version: $goVersion" -ForegroundColor Yellow

# Set environment for Linux builds
$env:GOOS = "linux"
$env:GOARCH = "amd64"

Write-Host "Building for Linux (GOOS=linux, GOARCH=amd64)..." -ForegroundColor Yellow

# Build the project
go build -o servin-linux .
if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful! Binary: servin-linux" -ForegroundColor Green
    Write-Host ""
    Write-Host "To deploy on Linux:" -ForegroundColor Cyan
    Write-Host "1. Copy servin-linux to your Linux system"
    Write-Host "2. chmod +x servin-linux"
    Write-Host "3. sudo mv servin-linux /usr/local/bin/servin"
    Write-Host "4. sudo servin run alpine echo 'Hello from container!'"
} else {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}

# Also build for current platform (Windows) for development
$env:GOOS = ""
$env:GOARCH = ""

Write-Host ""
Write-Host "Building for current platform (development)..." -ForegroundColor Yellow
go build -o servin.exe .
if ($LASTEXITCODE -eq 0) {
    Write-Host "Development build successful! Binary: servin.exe" -ForegroundColor Green
    Write-Host "Note: Windows binary is for development only. Container features require Linux."
} else {
    Write-Host "Development build failed!" -ForegroundColor Red
}

Write-Host ""
Write-Host "=== Build Complete ===" -ForegroundColor Green
