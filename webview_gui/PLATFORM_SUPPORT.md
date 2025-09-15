# Servin Desktop GUI - Cross-Platform Support

The Servin Desktop GUI automatically detects and uses the appropriate servin binary for your platform.

## Supported Platforms

- **Windows** (amd64)
- **Linux** (amd64, arm64)
- **macOS** (amd64, arm64)

## Binary Detection Logic

The GUI follows this priority order when looking for the servin binary:

### Windows
1. `build/windows-amd64/servin.exe` (platform-specific build)
2. `servin.exe` (root directory)
3. `servin` (fallback)

### macOS
1. `build/darwin-arm64/servin` (for Apple Silicon)
2. `build/darwin-amd64/servin` (for Intel Macs)
3. `servin` (root directory fallback)
4. `servin.exe` (fallback)

### Linux
1. `build/linux-arm64/servin` (for ARM64 systems)
2. `build/linux-amd64/servin` (for x86-64 systems)
3. `servin` (root directory fallback)
4. `servin.exe` (fallback)

## Building for All Platforms

To build servin for all platforms, use the build script:

### Windows (PowerShell)
```powershell
.\build.ps1 -Target all
```

### Linux/macOS (Bash)
```bash
./build-cross.sh --all
```

This will create platform-specific binaries in the `build/` directory:
- `build/windows-amd64/servin.exe`
- `build/linux-amd64/servin`
- `build/linux-arm64/servin`
- `build/darwin-amd64/servin`
- `build/darwin-arm64/servin`

## Running the GUI

The GUI will automatically detect your platform and use the appropriate binary:

```bash
# All platforms - GUI with embedded browser
python main.py

# All platforms - Browser-based demo
python demo.py
```

## Fallback Behavior

If the platform-specific binary is not found, the GUI will:
1. Fall back to the root directory binary
2. If still not found, use the mock client for demonstration
3. Display appropriate error messages

## Requirements

- Python 3.8+
- Flask 3.0.3+
- pywebview 5.1 (for embedded browser, optional)
- Appropriate servin binary for your platform

## Architecture Independence

The GUI automatically detects:
- Operating system (Windows, Linux, macOS)
- Architecture (amd64, arm64)
- Available binaries
- Executable permissions

This ensures the GUI works seamlessly across different development and deployment environments.
