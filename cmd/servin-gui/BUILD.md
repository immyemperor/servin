# Servin Desktop - Build Instructions

## Platform Setup Requirements

### Windows (Current Environment)

**Issue**: CGO is disabled and OpenGL/C toolchain is missing

**Required Setup**:
1. Install TDM-GCC or MinGW-w64:
   ```
   Download from: https://jmeubank.github.io/tdm-gcc/
   Or: https://www.mingw-w64.org/downloads/
   ```

2. Enable CGO:
   ```bash
   go env -w CGO_ENABLED=1
   set CGO_ENABLED=1
   ```

3. Install OpenGL headers (if needed):
   ```bash
   # Usually included with graphics drivers
   # May need Windows SDK
   ```

4. Build with CGO:
   ```bash
   go build -o servin-gui.exe cmd/servin-gui/main.go cmd/servin-gui/actions.go
   ```

### Linux (Recommended for Development)

**Setup**:
```bash
# Ubuntu/Debian
sudo apt-get install build-essential libgl1-mesa-dev xorg-dev

# CentOS/RHEL
sudo yum groupinstall "Development Tools"
sudo yum install mesa-libGL-devel libXrandr-devel libXinerama-devel libXcursor-devel libXi-devel

# Build
go build -o servin-gui cmd/servin-gui/main.go cmd/servin-gui/actions.go
```

### macOS

**Setup**:
```bash
# Install Xcode Command Line Tools
xcode-select --install

# Build
go build -o servin-gui cmd/servin-gui/main.go cmd/servin-gui/actions.go
```

## Alternative: Use Go Build Tags

For environments without GUI support, create console-only version:

```go
//go:build !gui

package main

import "fmt"

func main() {
    fmt.Println("GUI not available - use 'servin desktop' for terminal interface")
}
```

## Docker Build Environment

Create containerized build environment:

```dockerfile
FROM golang:1.24-bullseye

RUN apt-get update && apt-get install -y \
    build-essential \
    libgl1-mesa-dev \
    xorg-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .
RUN go build -o servin-gui cmd/servin-gui/main.go cmd/servin-gui/actions.go
```

## Testing the GUI

Once built successfully:

```bash
# Direct execution
./servin-gui

# Through CLI integration  
./servin gui

# With debug output
FYNE_DEBUG=1 ./servin-gui
```
