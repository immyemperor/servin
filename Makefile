# Servin Container Runtime Makefile

# Build variables
BINARY_NAME=servin
TUI_BINARY=servin-tui
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date +%Y-%m-%dT%H:%M:%S)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Platform detection
ifeq ($(OS),Windows_NT)
    BINARY_EXT=.exe
    TUI_BINARY_EXT=.exe
else
    BINARY_EXT=
    TUI_BINARY_EXT=
endif

.PHONY: build build-tui build-all clean test help install run-tui setup-webview-gui

# Default target
all: build build-tui

# Build main servin binary
build:
	@echo "Building Servin runtime..."
	go build $(LDFLAGS) -o $(BINARY_NAME)$(BINARY_EXT) .

# Build TUI component
build-tui:
	@echo "Building Servin Desktop TUI..."
	go build $(LDFLAGS) -o $(TUI_BINARY)$(TUI_BINARY_EXT) ./cmd/servin-tui

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
build-all: build build-tui setup-webview-gui

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)$(BINARY_EXT)
	rm -f $(TUI_BINARY)$(TUI_BINARY_EXT)
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
