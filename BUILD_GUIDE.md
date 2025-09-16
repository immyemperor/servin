# Servin Build Scripts

This document describes the build scripts and processes for the Servin Container Runtime project.

## Build Scripts

### `build-local.sh` ⭐ **Recommended for Current Platform**
- **Purpose**: Builds all Servin binaries for the current platform into `build/<platform>/`
- **Output**: Creates `build/<platform>/servin`, `build/<platform>/servin-tui`, `build/<platform>/servin-gui`
- **Usage**: `./build-local.sh`
- **Features**:
  - Clean build environment (removes old build directory)
  - Platform-specific directory organization
  - Colorized output with progress indicators
  - Automatic version detection from git
  - Generates README.md in each platform directory
  - Error handling for GUI dependencies

### `build-cross.sh` ⭐ **Recommended for Multi-Platform**
- **Purpose**: Cross-platform build script for multiple platforms
- **Output**: Creates `build/<platform>/` directories for each target
- **Usage**: 
  - Current platform: `./build-cross.sh --current`
  - All platforms: `./build-cross.sh --all` 
  - Specific platforms: `./build-cross.sh --platform linux/amd64 --platform windows/amd64`
- **Features**:
  - Cross-compilation support
  - Multiple platform targeting
  - Automatic file extension handling (.exe for Windows)
  - CGO handling for GUI components
  - Platform-specific documentation

### `rebuild.sh`
- **Purpose**: Quick alias for `build-local.sh`
- **Usage**: `./rebuild.sh`

### `build.sh`
- **Purpose**: Original cross-platform build script for distribution packages
- **Output**: Creates platform-specific packages in `dist/` directory
- **Usage**: See script for detailed options

## Supported Platforms

- **linux/amd64** - Linux 64-bit Intel/AMD
- **linux/arm64** - Linux 64-bit ARM (Apple Silicon, Raspberry Pi, etc.)
- **darwin/amd64** - macOS Intel 64-bit
- **darwin/arm64** - macOS Apple Silicon (M1/M2/M3)
- **windows/amd64** - Windows 64-bit

## Makefile Targets

### `make build-local` ⭐ **Recommended**
- Runs `build-local.sh` script
- Builds for current platform only

### `make build`
- Builds only the main `servin` binary in project root

### `make build-tui`
- Builds only the `servin-tui` TUI binary in project root

### `make build-gui`
- Builds only the `servin-gui` GUI binary in project root

### `make build-all`
- Builds all three binaries in project root

### `make clean`
- Removes build artifacts from project root

## Build Output

### Directory Structure
```
build/
├── README.md                    # Main build documentation
├── darwin-arm64/               # macOS Apple Silicon binaries
│   ├── README.md               # Platform-specific docs
│   ├── servin                  # Main runtime
│   ├── servin-tui          # TUI application
│   └── servin-gui              # GUI application
├── linux-amd64/               # Linux 64-bit binaries
│   ├── README.md
│   ├── servin
│   ├── servin-tui
│   └── servin-gui
└── windows-amd64/              # Windows 64-bit binaries
    ├── README.md
    ├── servin.exe
    ├── servin-tui.exe
    └── servin-gui.exe
```

### Generated Binaries

| Binary | Size | Description |
|--------|------|-------------|
| `servin` | ~11MB | Main container runtime |
| `servin-tui` | ~3MB | Terminal UI application |
| `servin-gui` | ~31MB | Desktop GUI application |

## Quick Commands

```bash
# Build for current platform only
./build-local.sh
# OR
make build-local

# Build for multiple platforms
./build-cross.sh --all
./build-cross.sh --platform linux/amd64 --platform windows/amd64

# Quick rebuild current platform
./rebuild.sh

# Test the binaries (example for macOS)
./build/darwin-arm64/servin --help
./build/darwin-arm64/servin-tui
./build/darwin-arm64/servin-gui

# Cross-platform usage
./build/linux-amd64/servin --help          # Linux
./build/windows-amd64/servin.exe --help    # Windows
```

## Development Workflow

1. **Make changes** to source code
2. **Rebuild** with `./build-local.sh` for quick testing
3. **Test** binaries in `build/<current-platform>/` directory
4. **Cross-compile** with `./build-cross.sh --all` for final testing
5. **Iterate** as needed

## Cross-Compilation Notes

- **GUI components**: May not build for all platforms due to CGO dependencies
- **Windows builds**: Automatically get `.exe` extensions
- **Linux builds**: Require Linux-specific code paths
- **macOS builds**: Support both Intel and Apple Silicon

## Notes

- The `build/` directory is in `.gitignore` and won't be committed
- Build scripts automatically detect version from git tags
- Platform detection is automatic based on `go env GOOS` and `go env GOARCH`
- All scripts are designed to work on macOS, Linux, and Windows (with appropriate shell)
