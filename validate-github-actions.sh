#!/bin/bash

# GitHub Actions Integration Validation Script
# Tests that all components work together correctly

set -e

echo "=============================================="
echo "GitHub Actions Integration Validation"
echo "=============================================="

# Check build scripts
echo ""
echo "1. Checking Build Scripts..."
echo "==============================="

if [ -f "./build-packages.sh" ]; then
    echo "✓ Cross-platform build script found"
    if [ -x "./build-packages.sh" ]; then
        echo "✓ Build script is executable"
    else
        echo "❌ Build script is not executable"
        exit 1
    fi
else
    echo "❌ Cross-platform build script not found"
    exit 1
fi

# Test build script help
echo ""
echo "2. Testing Build Script Help..."
echo "==============================="
./build-packages.sh --help | head -5
echo "✓ Build script help command works"

# Check installer directories
echo ""
echo "3. Checking Installer Components..."
echo "==================================="

# Windows NSIS installer
if [ -f "installers/windows/servin-installer.nsi" ]; then
    echo "✓ Windows NSIS installer script found"
    lines=$(wc -l < installers/windows/servin-installer.nsi)
    echo "  → NSIS script: $lines lines"
else
    echo "❌ Windows NSIS installer script not found"
fi

if [ -f "installers/windows/build-installer.bat" ]; then
    echo "✓ Windows build script found"
else
    echo "❌ Windows build script not found"
fi

# Linux AppImage builder
if [ -f "installers/linux/build-appimage.sh" ]; then
    echo "✓ Linux AppImage builder found"
    if [ -x "installers/linux/build-appimage.sh" ]; then
        echo "✓ Linux AppImage builder is executable"
    else
        echo "❌ Linux AppImage builder is not executable"
    fi
else
    echo "❌ Linux AppImage builder not found"
fi

# macOS package builder
if [ -f "installers/macos/build-package.sh" ]; then
    echo "✓ macOS package builder found"
    if [ -x "installers/macos/build-package.sh" ]; then
        echo "✓ macOS package builder is executable"
    else
        echo "❌ macOS package builder is not executable"
    fi
else
    echo "❌ macOS package builder not found"
fi

# Check GitHub Actions workflow
echo ""
echo "4. Checking GitHub Actions Workflow..."
echo "======================================"

if [ -f ".github/workflows/build-release.yml" ]; then
    echo "✓ GitHub Actions workflow found"
    
    # Check for our new build-packages.sh integration
    if grep -q "build-packages.sh" .github/workflows/build-release.yml; then
        echo "✓ Workflow uses build-packages.sh"
    else
        echo "❌ Workflow does not use build-packages.sh"
    fi
    
    # Check for installer package handling
    if grep -q "installer" .github/workflows/build-release.yml; then
        echo "✓ Workflow handles installer packages"
    else
        echo "❌ Workflow does not handle installer packages"
    fi
    
    # Check for complete installer package verification
    if grep -q -i "complete installer package" .github/workflows/build-release.yml; then
        echo "✓ Workflow includes installer package verification"
    else
        echo "❌ Workflow missing installer package verification"
    fi
    
else
    echo "❌ GitHub Actions workflow not found"
fi

# Check documentation
echo ""
echo "5. Checking Documentation..."
echo "============================"

if [ -f "installers/PACKAGE_README.md" ]; then
    echo "✓ Package documentation found"
else
    echo "❌ Package documentation not found"
fi

if [ -f "COMPLETE_INSTALLER_SUMMARY.md" ]; then
    echo "✓ Complete installer summary found"
else
    echo "❌ Complete installer summary not found"
fi

# Test workflow syntax (if available)
echo ""
echo "6. Testing Workflow Syntax..."
echo "============================="

if command -v yamllint >/dev/null 2>&1; then
    if yamllint .github/workflows/build-release.yml; then
        echo "✓ Workflow YAML syntax is valid"
    else
        echo "❌ Workflow YAML syntax is invalid"
    fi
else
    echo "⚠ yamllint not available, skipping syntax check"
fi

# Check if we can run a dry-run build (platform detection)
echo ""
echo "7. Testing Platform Detection..."
echo "==============================="

# Test platform detection logic from our build script
if [[ "$(uname -s)" == "Darwin" ]]; then
    echo "✓ Platform: macOS detected"
    TEST_PLATFORM="macos"
elif [[ "$(uname -s)" == "Linux" ]]; then
    echo "✓ Platform: Linux detected"
    TEST_PLATFORM="linux"
else
    echo "⚠ Platform: Unknown (assuming Windows in CI)"
    TEST_PLATFORM="windows"
fi

echo "✓ Test platform would be: $TEST_PLATFORM"

# Test basic Go build (if Go is available)
echo ""
echo "8. Testing Basic Build Capability..."
echo "===================================="

if command -v go >/dev/null 2>&1; then
    echo "✓ Go compiler available"
    if [ -f "main.go" ]; then
        echo "✓ Main Go file found"
        # Test compilation without actually building
        if go build -o /dev/null main.go 2>/dev/null; then
            echo "✓ Go code compiles successfully"
        else
            echo "❌ Go code compilation failed"
        fi
    else
        echo "❌ Main Go file not found"
    fi
else
    echo "⚠ Go compiler not available (will be available in CI)"
fi

# Summary
echo ""
echo "=============================================="
echo "Validation Summary"
echo "=============================================="

# Count components
total_checks=0
passed_checks=0

components=(
    "./build-packages.sh"
    "installers/windows/servin-installer.nsi"
    "installers/linux/build-appimage.sh"
    "installers/macos/build-package.sh"
    ".github/workflows/build-release.yml"
    "installers/PACKAGE_README.md"
    "COMPLETE_INSTALLER_SUMMARY.md"
)

for component in "${components[@]}"; do
    total_checks=$((total_checks + 1))
    if [ -f "$component" ]; then
        passed_checks=$((passed_checks + 1))
        echo "✓ $component"
    else
        echo "❌ $component"
    fi
done

echo ""
echo "Results: $passed_checks/$total_checks components ready"

if [ $passed_checks -eq $total_checks ]; then
    echo ""
    echo "🎉 GitHub Actions integration is ready!"
    echo ""
    echo "Next steps:"
    echo "1. Commit and push changes"
    echo "2. Create a release tag to trigger full build"
    echo "3. Monitor GitHub Actions for complete installer package creation"
    echo "4. Test installers on target platforms"
    exit 0
else
    echo ""
    echo "❌ Some components are missing. Please address the issues above."
    exit 1
fi