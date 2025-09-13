---
layout: default
title: Development Guide
---

# Development Guide

Complete guide for developing, contributing, and extending Servin Container Runtime.

## Development Environment

### Prerequisites

Set up development environment:

```bash
# Required software
- Go 1.21+ (latest stable recommended)
- Git 2.30+
- Make 4.0+
- Docker/Podman (for testing)
- Protocol Buffers compiler (protoc)
- gRPC tools

# Development tools
- golangci-lint (code linting)
- goreleaser (release automation)
- testify (testing framework)
- mockgen (mock generation)
- protoc-gen-go (protobuf generation)
```

### Environment Setup

```bash
# Install Go (Linux/macOS)
curl -L https://go.dev/dl/go1.21.5.linux-amd64.tar.gz | sudo tar -xzC /usr/local
export PATH=$PATH:/usr/local/go/bin

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/golang/mock/mockgen@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Clone repository
git clone https://github.com/servin-dev/servin.git
cd servin

# Install dependencies
go mod download

# Build from source
make build

# Run tests
make test

# Install locally
make install
```

## Project Structure

### Source Code Organization

```
servin/
├── cmd/                    # Command line applications
│   ├── servind/           # Daemon main
│   ├── servin/            # CLI client
│   ├── servin-tui/        # Terminal UI
│   └── servin-gui/        # Desktop GUI
├── pkg/                   # Public libraries
│   ├── api/               # API definitions
│   ├── client/            # Client libraries
│   ├── container/         # Container management
│   ├── image/             # Image management
│   ├── network/           # Network management
│   ├── storage/           # Storage drivers
│   ├── runtime/           # OCI runtime integration
│   └── util/              # Utility functions
├── internal/              # Private libraries
│   ├── daemon/            # Daemon implementation
│   ├── server/            # API server
│   ├── cri/               # CRI implementation
│   └── config/            # Configuration
├── api/                   # API specifications
│   ├── proto/             # Protocol Buffers
│   └── openapi/           # OpenAPI specifications
├── docs/                  # Documentation
├── test/                  # Test files
│   ├── unit/              # Unit tests
│   ├── integration/       # Integration tests
│   └── e2e/               # End-to-end tests
├── scripts/               # Build and utility scripts
├── deployments/           # Deployment configurations
├── examples/              # Example configurations
├── Makefile              # Build automation
├── go.mod                # Go module
├── go.sum                # Go dependencies
├── Dockerfile            # Container build
└── README.md             # Project documentation
```

### Core Components

#### Daemon Architecture

```go
// internal/daemon/daemon.go
package daemon

import (
    "context"
    "github.com/servin-dev/servin/pkg/container"
    "github.com/servin-dev/servin/pkg/image"
    "github.com/servin-dev/servin/pkg/network"
    "github.com/servin-dev/servin/pkg/storage"
)

type Daemon struct {
    containerManager *container.Manager
    imageManager     *image.Manager
    networkManager   *network.Manager
    storageManager   *storage.Manager
    config          *Config
}

func New(cfg *Config) (*Daemon, error) {
    d := &Daemon{
        config: cfg,
    }
    
    // Initialize managers
    d.containerManager = container.NewManager(cfg.ContainerConfig)
    d.imageManager = image.NewManager(cfg.ImageConfig)
    d.networkManager = network.NewManager(cfg.NetworkConfig)
    d.storageManager = storage.NewManager(cfg.StorageConfig)
    
    return d, nil
}

func (d *Daemon) Start(ctx context.Context) error {
    // Start all managers
    if err := d.storageManager.Start(ctx); err != nil {
        return err
    }
    
    if err := d.networkManager.Start(ctx); err != nil {
        return err
    }
    
    if err := d.imageManager.Start(ctx); err != nil {
        return err
    }
    
    if err := d.containerManager.Start(ctx); err != nil {
        return err
    }
    
    return nil
}
```

#### Container Manager

```go
// pkg/container/manager.go
package container

import (
    "context"
    "github.com/opencontainers/runtime-spec/specs-go"
)

type Manager struct {
    runtime Runtime
    store   Store
    config  *Config
}

type Container struct {
    ID       string
    Name     string
    Image    string
    Config   *specs.Spec
    State    State
    Metadata map[string]string
}

type State int

const (
    StateCreated State = iota
    StateRunning
    StatePaused
    StateStopped
    StateUnknown
)

func (m *Manager) Create(ctx context.Context, req *CreateRequest) (*Container, error) {
    // Validate request
    if err := req.Validate(); err != nil {
        return nil, err
    }
    
    // Generate container ID
    id := generateID()
    
    // Create OCI spec
    spec, err := m.buildSpec(req)
    if err != nil {
        return nil, err
    }
    
    // Create container
    container := &Container{
        ID:     id,
        Name:   req.Name,
        Image:  req.Image,
        Config: spec,
        State:  StateCreated,
    }
    
    // Store container
    if err := m.store.Save(container); err != nil {
        return nil, err
    }
    
    return container, nil
}

func (m *Manager) Start(ctx context.Context, id string) error {
    container, err := m.store.Get(id)
    if err != nil {
        return err
    }
    
    if container.State != StateCreated {
        return ErrInvalidState
    }
    
    // Start container via runtime
    if err := m.runtime.Start(ctx, container); err != nil {
        return err
    }
    
    container.State = StateRunning
    return m.store.Save(container)
}
```

## Building from Source

### Build System

Makefile targets for development:

```makefile
# Makefile
.PHONY: build test clean install

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build info
VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Build flags
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

# Targets
all: build

build: build-daemon build-client build-tui build-gui

build-daemon:
    $(GOBUILD) $(LDFLAGS) -o bin/servind cmd/servind/main.go

build-client:
    $(GOBUILD) $(LDFLAGS) -o bin/servin cmd/servin/main.go

build-tui:
    $(GOBUILD) $(LDFLAGS) -o bin/servin-tui cmd/servin-tui/main.go

build-gui:
    $(GOBUILD) $(LDFLAGS) -o bin/servin-gui cmd/servin-gui/main.go

test:
    $(GOTEST) -v -race -cover ./...

test-integration:
    $(GOTEST) -v -tags=integration ./test/integration/...

test-e2e:
    $(GOTEST) -v -tags=e2e ./test/e2e/...

lint:
    golangci-lint run

format:
    gofmt -s -w .
    goimports -w .

clean:
    $(GOCLEAN)
    rm -rf bin/

install: build
    sudo cp bin/servind /usr/local/bin/
    sudo cp bin/servin /usr/local/bin/
    sudo cp bin/servin-tui /usr/local/bin/
    sudo cp bin/servin-gui /usr/local/bin/

docker-build:
    docker build -t servin:dev .

proto:
    protoc --go_out=. --go-grpc_out=. api/proto/*.proto

deps:
    $(GOMOD) download
    $(GOMOD) verify

update-deps:
    $(GOMOD) get -u ./...
    $(GOMOD) tidy

generate:
    $(GOCMD) generate ./...

release:
    goreleaser release --clean

snapshot:
    goreleaser release --snapshot --clean
```

### Cross-platform Build

Build for multiple platforms:

```bash
# Build for all platforms
make cross-compile

# Or manually
GOOS=linux GOARCH=amd64 go build -o bin/servin-linux-amd64 cmd/servin/main.go
GOOS=windows GOARCH=amd64 go build -o bin/servin-windows-amd64.exe cmd/servin/main.go
GOOS=darwin GOARCH=amd64 go build -o bin/servin-darwin-amd64 cmd/servin/main.go
GOOS=darwin GOARCH=arm64 go build -o bin/servin-darwin-arm64 cmd/servin/main.go
```

### Container Build

Build Servin in containers:

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o servind cmd/servind/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o servin cmd/servin/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /src/servind .
COPY --from=builder /src/servin .

EXPOSE 2375 2376

CMD ["./servind"]
```

## Testing

### Unit Tests

Write comprehensive unit tests:

```go
// pkg/container/manager_test.go
package container_test

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/servin-dev/servin/pkg/container"
)

func TestManager_Create(t *testing.T) {
    tests := []struct {
        name    string
        request *container.CreateRequest
        want    *container.Container
        wantErr bool
    }{
        {
            name: "valid container creation",
            request: &container.CreateRequest{
                Name:  "test-container",
                Image: "nginx:latest",
            },
            want: &container.Container{
                Name:  "test-container",
                Image: "nginx:latest",
                State: container.StateCreated,
            },
            wantErr: false,
        },
        {
            name: "invalid container name",
            request: &container.CreateRequest{
                Name:  "",
                Image: "nginx:latest",
            },
            want:    nil,
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            manager := setupTestManager(t)
            
            got, err := manager.Create(context.Background(), tt.request)
            
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, got)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, got)
                assert.Equal(t, tt.want.Name, got.Name)
                assert.Equal(t, tt.want.Image, got.Image)
                assert.Equal(t, tt.want.State, got.State)
            }
        })
    }
}

func setupTestManager(t *testing.T) *container.Manager {
    mockRuntime := &MockRuntime{}
    mockStore := &MockStore{}
    
    config := &container.Config{
        Runtime: "test",
    }
    
    return container.NewManager(mockRuntime, mockStore, config)
}
```

### Integration Tests

Test component integration:

```go
// test/integration/container_test.go
// +build integration

package integration

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/suite"
    "github.com/servin-dev/servin/pkg/container"
)

type ContainerIntegrationSuite struct {
    suite.Suite
    manager *container.Manager
    ctx     context.Context
}

func (s *ContainerIntegrationSuite) SetupSuite() {
    s.ctx = context.Background()
    
    // Setup real container manager with test config
    config := &container.Config{
        Runtime:   "runc",
        StateDir:  "/tmp/servin-test",
        ImageDir:  "/tmp/servin-images",
    }
    
    manager, err := container.NewManager(config)
    s.Require().NoError(err)
    
    s.manager = manager
}

func (s *ContainerIntegrationSuite) TestContainerLifecycle() {
    // Create container
    req := &container.CreateRequest{
        Name:  "test-integration",
        Image: "alpine:latest",
        Cmd:   []string{"sleep", "30"},
    }
    
    container, err := s.manager.Create(s.ctx, req)
    s.Require().NoError(err)
    s.Assert().Equal("test-integration", container.Name)
    
    // Start container
    err = s.manager.Start(s.ctx, container.ID)
    s.Require().NoError(err)
    
    // Wait for container to be running
    time.Sleep(1 * time.Second)
    
    // Check status
    status, err := s.manager.Status(s.ctx, container.ID)
    s.Require().NoError(err)
    s.Assert().Equal(container.StateRunning, status.State)
    
    // Stop container
    err = s.manager.Stop(s.ctx, container.ID, 10)
    s.Require().NoError(err)
    
    // Remove container
    err = s.manager.Remove(s.ctx, container.ID)
    s.Require().NoError(err)
}

func TestContainerIntegration(t *testing.T) {
    suite.Run(t, new(ContainerIntegrationSuite))
}
```

### End-to-End Tests

Test complete workflows:

```go
// test/e2e/cli_test.go
// +build e2e

package e2e

import (
    "os/exec"
    "strings"
    "testing"
    
    "github.com/stretchr/testify/assert"
)

func TestCLIWorkflow(t *testing.T) {
    // Start daemon
    daemon := exec.Command("servind", "--config", "test-config.json")
    err := daemon.Start()
    assert.NoError(t, err)
    defer daemon.Process.Kill()
    
    // Wait for daemon to start
    time.Sleep(3 * time.Second)
    
    // Test container creation
    cmd := exec.Command("servin", "run", "-d", "--name", "test-e2e", "alpine:latest", "sleep", "30")
    output, err := cmd.CombinedOutput()
    assert.NoError(t, err)
    
    containerID := strings.TrimSpace(string(output))
    assert.NotEmpty(t, containerID)
    
    // Test container listing
    cmd = exec.Command("servin", "ps")
    output, err = cmd.CombinedOutput()
    assert.NoError(t, err)
    assert.Contains(t, string(output), "test-e2e")
    
    // Test container stop
    cmd = exec.Command("servin", "stop", "test-e2e")
    _, err = cmd.CombinedOutput()
    assert.NoError(t, err)
    
    // Test container removal
    cmd = exec.Command("servin", "rm", "test-e2e")
    _, err = cmd.CombinedOutput()
    assert.NoError(t, err)
}
```

## Contributing Guidelines

### Code Style

Follow Go best practices:

```go
// Good: Clear function names and documentation
// CreateContainer creates a new container with the specified configuration.
// It validates the input parameters and returns an error if validation fails.
func (m *Manager) CreateContainer(ctx context.Context, req *CreateContainerRequest) (*Container, error) {
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }
    
    // Implementation...
}

// Good: Use meaningful variable names
containerManager := container.NewManager(config)
networkInterface := "eth0"
maxRetries := 3

// Bad: Unclear names
mgr := container.NewManager(config)
iface := "eth0"
max := 3

// Good: Handle errors appropriately
result, err := operation()
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// Good: Use context for cancellation
func (s *Server) handleRequest(ctx context.Context, req *Request) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // Handle request
    }
}
```

### Commit Messages

Use conventional commit format:

```bash
# Format: type(scope): description
feat(container): add support for health checks
fix(api): resolve memory leak in container stats
docs(readme): update installation instructions
test(integration): add network isolation tests
refactor(storage): simplify driver interface
perf(runtime): optimize container startup time

# Breaking changes
feat(api)!: change container create endpoint response format

BREAKING CHANGE: Container create endpoint now returns detailed container info instead of just ID
```

### Pull Request Process

1. **Fork and Branch**
```bash
git fork https://github.com/servin-dev/servin
git checkout -b feature/my-feature
```

2. **Make Changes**
```bash
# Write code
# Add tests
# Update documentation
```

3. **Test Changes**
```bash
make test
make test-integration
make lint
```

4. **Commit and Push**
```bash
git add .
git commit -m "feat(container): add health check support"
git push origin feature/my-feature
```

5. **Create Pull Request**
- Use descriptive title
- Include detailed description
- Reference related issues
- Add screenshots for UI changes

### Code Review Checklist

- [ ] Code follows style guidelines
- [ ] Tests are included and passing
- [ ] Documentation is updated
- [ ] No breaking changes (or properly documented)
- [ ] Performance impact considered
- [ ] Security implications reviewed
- [ ] Error handling is appropriate
- [ ] Logging is adequate

## API Development

### Adding New Endpoints

Add new API endpoints:

```go
// internal/server/routes.go
func (s *Server) setupRoutes() {
    // Container routes
    s.router.Handle("/containers/json", s.handleContainerList).Methods("GET")
    s.router.Handle("/containers/create", s.handleContainerCreate).Methods("POST")
    s.router.Handle("/containers/{id}/start", s.handleContainerStart).Methods("POST")
    
    // New endpoint
    s.router.Handle("/containers/{id}/health", s.handleContainerHealth).Methods("GET")
}

func (s *Server) handleContainerHealth(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    
    health, err := s.daemon.GetContainerHealth(r.Context(), id)
    if err != nil {
        s.writeError(w, err)
        return
    }
    
    s.writeJSON(w, health)
}
```

### Protocol Buffers

Define gRPC services:

```protobuf
// api/proto/container.proto
syntax = "proto3";

package servin.container.v1;

option go_package = "github.com/servin-dev/servin/api/container/v1";

service ContainerService {
  rpc CreateContainer(CreateContainerRequest) returns (CreateContainerResponse);
  rpc StartContainer(StartContainerRequest) returns (StartContainerResponse);
  rpc StopContainer(StopContainerRequest) returns (StopContainerResponse);
  rpc ListContainers(ListContainersRequest) returns (ListContainersResponse);
  rpc GetContainer(GetContainerRequest) returns (GetContainerResponse);
  rpc DeleteContainer(DeleteContainerRequest) returns (DeleteContainerResponse);
}

message Container {
  string id = 1;
  string name = 2;
  string image = 3;
  ContainerState state = 4;
  map<string, string> labels = 5;
  int64 created_at = 6;
  int64 started_at = 7;
  int64 finished_at = 8;
}

enum ContainerState {
  CONTAINER_STATE_UNSPECIFIED = 0;
  CONTAINER_STATE_CREATED = 1;
  CONTAINER_STATE_RUNNING = 2;
  CONTAINER_STATE_PAUSED = 3;
  CONTAINER_STATE_STOPPED = 4;
  CONTAINER_STATE_REMOVING = 5;
}

message CreateContainerRequest {
  string name = 1;
  string image = 2;
  repeated string command = 3;
  repeated string args = 4;
  map<string, string> env = 5;
  map<string, string> labels = 6;
  ContainerConfig config = 7;
}

message CreateContainerResponse {
  Container container = 1;
}
```

## Performance Optimization

### Profiling

Profile application performance:

```go
// Enable profiling
import _ "net/http/pprof"

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // Main application
}

// Profile CPU usage
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

// Profile memory usage
go tool pprof http://localhost:6060/debug/pprof/heap

// Profile goroutines
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

### Benchmarking

Write performance benchmarks:

```go
// pkg/container/manager_bench_test.go
func BenchmarkContainerCreate(b *testing.B) {
    manager := setupBenchmarkManager(b)
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        req := &container.CreateRequest{
            Name:  fmt.Sprintf("bench-container-%d", i),
            Image: "alpine:latest",
        }
        
        _, err := manager.Create(context.Background(), req)
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkContainerStart(b *testing.B) {
    manager := setupBenchmarkManager(b)
    containers := make([]*container.Container, b.N)
    
    // Setup containers
    for i := 0; i < b.N; i++ {
        req := &container.CreateRequest{
            Name:  fmt.Sprintf("bench-container-%d", i),
            Image: "alpine:latest",
        }
        
        c, err := manager.Create(context.Background(), req)
        if err != nil {
            b.Fatal(err)
        }
        containers[i] = c
    }
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        err := manager.Start(context.Background(), containers[i].ID)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## Release Process

### Versioning

Follow semantic versioning:

- **Major**: Breaking changes (v1.0.0 → v2.0.0)
- **Minor**: New features (v1.0.0 → v1.1.0)
- **Patch**: Bug fixes (v1.0.0 → v1.0.1)

### Release Automation

Use GoReleaser for releases:

```yaml
# .goreleaser.yml
project_name: servin

before:
  hooks:
    - go mod tidy
    - go generate ./...

builds:
  - id: servind
    main: ./cmd/servind
    binary: servind
    goos:
      - linux
      - windows
      - darwin
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w
      - -X main.version={{.Version}}
      - -X main.commit={{.Commit}}
      - -X main.date={{.Date}}

  - id: servin
    main: ./cmd/servin
    binary: servin
    goos:
      - linux
      - windows
      - darwin
    goarch:
      - amd64
      - arm64

archives:
  - id: default
    builds:
      - servind
      - servin
    name_template: "servin_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: 'checksums.txt'

changelog:
  sort: asc
  filters:
    exclude:
      - '^docs:'
      - '^test:'

dockers:
  - image_templates:
      - "servin/servin:{{ .Version }}"
      - "servin/servin:latest"
    dockerfile: Dockerfile
    build_flag_templates:
      - "--platform=linux/amd64"

release:
  github:
    owner: servin-dev
    name: servin
  draft: false
  prerelease: auto
```

This comprehensive development guide covers all aspects of contributing to and extending Servin Container Runtime, from initial setup to production releases.
