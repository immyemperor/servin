# Windows PyInstaller Build Troubleshooting

## Issue: ModuleNotFoundError: No module named 'app'

This error occurs when PyInstaller doesn't properly include the `app.py` module in the Windows executable.

## Solutions

### 1. Use the Updated Spec File

The issue has been fixed by creating a proper `servin-gui.spec` file that explicitly includes all required modules. Make sure you're using the latest build script.

### 2. Manual Build for Testing

If the automated build fails, you can manually build on Windows:

```bash
# Navigate to webview_gui directory
cd webview_gui

# Create virtual environment
python -m venv venv
venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt

# Test imports first
python test_imports.py

# Build with explicit spec file
pyinstaller --clean --distpath=dist --workpath=build servin-gui.spec

# Test the built executable
python test_executable.py dist\servin-gui.exe
```

### 3. Debug Build

For debugging, you can create a console version to see error messages:

```bash
# Modify the spec file temporarily for debugging
# Change: console=False to console=True
# Change: windowed=True to windowed=False

pyinstaller --clean servin-gui.spec
```

### 4. Verify Module Inclusion

Check if modules are properly included:

```bash
# List contents of the built executable
pyi-archive_viewer dist\servin-gui.exe
```

## Updated Files

The following files have been updated to fix this issue:

1. **webview_gui/servin-gui.spec** - Explicit module inclusion
2. **build-all.sh** - Updated PyInstaller command to use spec file
3. **webview_gui/main.py** - Better error handling and path detection
4. **webview_gui/test_imports.py** - Test script to verify modules
5. **webview_gui/test_executable.py** - Test script to verify built executable

## Verification

After building, run:

```bash
# Test the executable
dist\servin-gui.exe

# Should see debug output like:
# 🚀 Running from PyInstaller bundle: C:\Users\...\servin-gui\_internal
# 🐍 Python executable: C:\Users\...\servin-gui\servin-gui.exe
# ✅ Successfully imported Flask app
```

## Cross-Platform Build Limitation

Note: PyInstaller cannot cross-compile from macOS/Linux to Windows. Windows executables must be built on Windows.

If you need Windows builds and don't have access to Windows:
- Use GitHub Actions Windows runner
- Use Windows VM or container
- Use cross-compilation tools like Wine (limited support)