# Windows PyInstaller Build Troubleshooting

## Issue: ModuleNotFoundError: No module named 'app'

This error occurs when PyInstaller doesn't properly include the `app.py` module in the Windows executable, specifically when building from the `webview_build_temp` directory.

## Root Cause

The build script creates a temporary directory `webview_build_temp` and copies source files there, but PyInstaller may not correctly resolve module paths from this temporary location.

## Solutions Applied

### 1. Enhanced PyInstaller Spec File

**File: `webview_gui/servin-gui.spec`**
- Added explicit `pathex=[current_dir]` to ensure current directory is in Python path
- Added Python files as data files: `('app.py', '.')`, `('servin_client.py', '.')`, etc.
- Added debug output to show build directory
- Added hooks directory for custom module detection

### 2. Robust Import Handling

**File: `webview_gui/main.py`**
- Enhanced PyInstaller bundle detection with `sys._MEIPASS`
- Multiple import strategies: direct import and `importlib.util`
- Comprehensive debugging output showing available files and paths
- Better error messages for troubleshooting

### 3. Enhanced Build Process

**File: `build-all.sh`**
- Added verbose PyInstaller output with `--log-level=INFO`
- Added debug output showing build directory and available files
- Better error handling and reporting

## Manual Build for Testing

If the automated build fails, you can manually build and debug on Windows:

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

# Build with verbose output for debugging
pyinstaller --clean --distpath=dist --workpath=build --log-level=DEBUG servin-gui.spec

# Test the built executable
python test_executable.py dist\servin-gui.exe
```

## Debug Build for Troubleshooting

To see detailed error messages, temporarily modify the spec file:

```python
# In servin-gui.spec, change:
debug=True,     # Enable debug mode
console=True,   # Show console for error messages
windowed=False, # Disable windowed mode for debugging
```

Then rebuild:
```bash
pyinstaller --clean servin-gui.spec
dist\servin-gui.exe  # Run and check console output
```

## Verification Steps

1. **Check module inclusion**: Verify `app.py` is bundled
   ```bash
   pyi-archive_viewer dist\servin-gui.exe
   # Look for 'app' in the archive contents
   ```

2. **Test import resolution**: The executable should show:
   ```
   [INFO] Running from PyInstaller bundle: C:\...\servin-gui\_internal
   [DEBUG] Available files: ['app.py', 'servin_client.py', ...]
   [SUCCESS] Successfully imported Flask app via direct import
   ```

3. **Check data files**: Verify templates and static files are included

## Enhanced Features

- **Multiple import methods**: Falls back to `importlib.util` if direct import fails
- **Explicit file inclusion**: Python modules added as data files
- **Better path handling**: Explicit `pathex` configuration
- **Comprehensive debugging**: Shows all paths and available files
- **Hook system**: Custom PyInstaller hooks for module detection

## Updated Files

1. **webview_gui/servin-gui.spec** - Enhanced spec with explicit paths and data files
2. **webview_gui/main.py** - Robust import handling with fallbacks
3. **webview_gui/hook-app.py** - Custom PyInstaller hook
4. **build-all.sh** - Verbose build output and debugging
5. **webview_gui/test_imports.py** - Windows-compatible testing
6. **webview_gui/test_executable.py** - Windows-compatible executable testing

The enhanced build process should now properly handle the `webview_build_temp` directory and ensure all modules are correctly bundled in the Windows executable.