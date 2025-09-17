# Servin Container Runtime Makefile with VM Distribution

# Build variables
BINARY_NAME=servin
TUI_BINARY=servin-tui
DESKTOP_BINARY=servin-desktop
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date +%Y-%m-%dT%H%M%S)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# VM Configuration
VM_VERSION=1.0.0
DIST_DIR=dist

# Platform detection
ifeq ($(OS),Windows_NT)
    BINARY_EXT=.exe
    TUI_BINARY_EXT=.exe
    DESKTOP_BINARY_EXT=.exe
else
    BINARY_EXT=
    TUI_BINARY_EXT=
    DESKTOP_BINARY_EXT=
endif

.PHONY: build build-tui build-desktop build-all build-vm clean test help install run-tui setup-webview-gui dist-all dist-clean vm-images packages installers vm-dist-enhance vm-dist-full

# Default target
all: build build-tui build-desktop

# Build main servin binary with enhanced VM support
build:
	@echo "Building Servin runtime with VM containerization..."
	go build $(LDFLAGS) -tags "vm_enabled,kvm,qemu,hvf,hyperv,vbox" -o $(BINARY_NAME)$(BINARY_EXT) .

# Build TUI component
build-tui:
	@echo "Building Servin TUI..."
	go build $(LDFLAGS) -o $(TUI_BINARY)$(TUI_BINARY_EXT) ./cmd/servin-tui

# Build Desktop GUI component with VM support
build-desktop:
	@echo "Building Servin Desktop GUI with VM monitoring..."
	go build $(LDFLAGS) -tags "desktop,vm_enabled,kvm,qemu,hvf,hyperv" -o $(DESKTOP_BINARY)$(DESKTOP_BINARY_EXT) ./cmd/servin-desktop

# Build with enhanced VM support
build-vm: build build-desktop
	@echo "✓ Built Servin with universal VM containerization support"
	@echo "  Supported VM providers: KVM, QEMU, Virtualization.framework, Hyper-V, VirtualBox"

# Setup WebView GUI dependencies
setup-webview-gui:
	@echo "Setting up WebView GUI dependencies..."
	@if [ -d "webview_gui" ]; then \
		cd webview_gui && \
		python3 -m venv venv 2>/dev/null || python -m venv venv && \
		. venv/bin/activate && \
		pip install -r requirements.txt; \
	else \
		echo "WebView GUI directory not found"; \
	fi

# Build all components
build-all: build build-tui build-desktop setup-webview-gui

# Cross-platform distribution build with VM support
dist-all:
	@echo "Building universal VM distribution for all platforms..."
	./build-vm-distribution.sh

# Enhance existing build with VM capabilities (for GitHub Actions)
vm-dist-enhance:
	@echo "Enhancing existing build with VM capabilities..."
	./build-vm-distribution.sh --enhance-existing

# Full VM distribution build
vm-dist-full: clean
	@echo "Building complete VM distribution from scratch..."
	./build-vm-distribution.sh --all

# Build VM images only
vm-images:
	@echo "Building VM images for containerization..."
	@mkdir -p $(DIST_DIR)/vm-images
	@if [ -f "scripts/build-vm-image.sh" ]; then \
		./scripts/build-vm-image.sh alpine 3.18 $(DIST_DIR)/vm-images/servin-alpine-$(VM_VERSION).qcow2; \
		./scripts/build-vm-image.sh ubuntu 22.04 $(DIST_DIR)/vm-images/servin-ubuntu-$(VM_VERSION).qcow2; \
		cd $(DIST_DIR)/vm-images && sha256sum *.qcow2 > checksums.txt; \
	else \
		echo "VM image build script not found"; \
	fi

# Create platform packages with VM support
packages: vm-dist-full
	@echo "Creating platform packages with VM containerization..."
	@echo "✓ Packages created in $(DIST_DIR)/packages/"

# Create installers with VM components
installers: packages
	@echo "Creating platform installers with VM support..."
	@echo "✓ Installers created in $(DIST_DIR)/installers/"

# Clean all build artifacts including distribution
dist-clean: clean
	@echo "Cleaning distribution artifacts..."
	rm -rf $(DIST_DIR)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)$(BINARY_EXT)
	rm -f $(TUI_BINARY)$(TUI_BINARY_EXT) 
	rm -f $(DESKTOP_BINARY)$(DESKTOP_BINARY_EXT)
	@if [ -d "webview_gui/venv" ]; then \
		echo "Cleaning WebView GUI virtual environment..."; \
		rm -rf webview_gui/venv; \
	fi

# Run TUI
run-tui: build build-tui
	./$(BINARY_NAME)$(BINARY_EXT) gui --tui

# Quick test
test-quick: build build-tui
	@echo "Quick functionality test..."
	./$(BINARY_NAME)$(BINARY_EXT) --version
	./$(BINARY_NAME)$(BINARY_EXT) gui --help
	rm -f servin servin-linux servin.exe

# Run tests
test:
	@echo "Running tests..."
	go test ./pkg/...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Install to system (Linux only)
install: build-linux
	@echo "Installing Servin to /usr/local/bin..."
	sudo cp servin-linux /usr/local/bin/servin
	sudo chmod +x /usr/local/bin/servin

# Run example tests (requires root on Linux)
test-containers:
	@echo "Running container tests..."
	@if [ "$(shell uname)" = "Linux" ]; then \
		sudo ./examples/test_containers.sh; \
	else \
		echo "Container tests can only run on Linux"; \
	fi

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	go mod tidy
	go get -u github.com/spf13/cobra@latest
	go get -u golang.org/x/sys/unix

# Check code quality
lint:
	@echo "Running linting..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Show help
help:
	@echo "Servin Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  build         Build for current platform"
	@echo "  build-linux   Build for Linux (for deployment)"
	@echo "  clean         Remove build artifacts"
	@echo "  test          Run unit tests"
	@echo "  deps          Install/update dependencies"
	@echo "  install       Install to system (Linux, requires sudo)"
	@echo "  test-containers Run integration tests (Linux, requires sudo)"
	@echo "  dev-setup     Set up development environment"
	@echo "  lint          Run code linting"
	@echo "  help          Show this help message"
	@echo ""
	@echo "Example usage:"
	@echo "  make build-linux    # Build for Linux deployment"
	@echo "  make install        # Install on Linux system"
	@echo "  sudo servin run alpine echo 'Hello!'"
