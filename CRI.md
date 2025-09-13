# Container Runtime Interface (CRI) Support

## Overview

Servin implements the Kubernetes Container Runtime Interface (CRI) specification v1alpha2, enabling it to be used as a container runtime for Kubernetes clusters. The implementation provides an HTTP-based CRI server that exposes all required CRI operations through RESTful endpoints.

## CRI Server Architecture

### Design Principles
- **HTTP-based API**: Simplified approach using HTTP/JSON instead of gRPC/protobuf for easier development and testing
- **Kubernetes Compatible**: Implements the complete CRI specification for pod and container lifecycle management
- **Cross-platform**: Works on Linux, Windows, and macOS development environments
- **Development Friendly**: Easy to test with curl, browser, or any HTTP client

### Core Components

#### 1. CRI Types (`pkg/cri/types.go`)
- **RuntimeService Interface**: Defines all runtime operations (version, status, pod/container lifecycle)
- **ImageService Interface**: Defines all image operations (list, pull, remove, status)
- **Data Structures**: Pod sandbox, container, and image metadata types

#### 2. Message Types (`pkg/cri/messages.go`)
- **Request/Response Messages**: JSON-serializable types for all CRI operations
- **Kubernetes Compatible**: Matches CRI specification message formats
- **Version Information**: Runtime name, version, and API version details

#### 3. Runtime Service (`pkg/cri/minimal_runtime.go`)
- **MinimalRuntimeService**: Implementation of the RuntimeService interface
- **Pod Sandbox Management**: Create, start, stop, remove, and list pod sandboxes
- **Container Lifecycle**: Full container management within pod sandboxes
- **State Tracking**: In-memory state management for pods and containers

#### 4. Image Service (`pkg/cri/image_service.go`)
- **ServinImageService**: Implementation of the ImageService interface
- **Image Operations**: List, pull, remove, and status operations
- **Servin Integration**: Bridges CRI image operations to Servin's image system
- **Filesystem Info**: Reports image storage usage and availability

#### 5. HTTP Server (`pkg/cri/server.go`)
- **CRIHTTPServer**: Main HTTP server implementing CRI endpoints
- **RESTful Endpoints**: Maps CRI operations to HTTP routes
- **JSON Marshaling**: Handles request/response serialization
- **Error Handling**: Proper HTTP status codes and error responses

#### 6. CLI Commands (`cmd/cri.go`)
- **`servin cri start`**: Start the CRI HTTP server
- **`servin cri status`**: Check server status and connectivity
- **`servin cri test`**: Test CRI server functionality

## CRI Endpoints

### Runtime Operations

#### Version Information
```bash
POST /v1/runtime/version
Content-Type: application/json

# Response
{
  "version": {
    "version": "v1alpha2",
    "runtime_name": "servin",
    "runtime_version": "0.1.0",
    "runtime_api_version": "v1alpha2"
  }
}
```

#### Runtime Status
```bash
POST /v1/runtime/status
Content-Type: application/json

# Response
{
  "status": {
    "conditions": [
      {
        "type": "RuntimeReady",
        "status": true,
        "reason": "RuntimeReady",
        "message": "Runtime is ready"
      }
    ]
  }
}
```

### Pod Sandbox Operations

#### List Pod Sandboxes
```bash
POST /v1/runtime/sandbox/list
Content-Type: application/json
{
  "filter": {}
}

# Response
{
  "items": [
    {
      "id": "sandbox-123",
      "metadata": {
        "name": "test-pod",
        "namespace": "default"
      },
      "state": "SANDBOX_READY",
      "created_at": "2024-01-01T12:00:00Z"
    }
  ]
}
```

#### Create Pod Sandbox
```bash
POST /v1/runtime/sandbox/create
Content-Type: application/json
{
  "config": {
    "metadata": {
      "name": "test-pod",
      "namespace": "default"
    },
    "hostname": "test-pod",
    "log_directory": "/var/log/pods"
  }
}

# Response
{
  "pod_sandbox_id": "sandbox-456"
}
```

### Container Operations

#### List Containers
```bash
POST /v1/runtime/container/list
Content-Type: application/json
{
  "filter": {}
}

# Response
{
  "containers": [
    {
      "id": "container-789",
      "pod_sandbox_id": "sandbox-123",
      "metadata": {
        "name": "test-container"
      },
      "image": {
        "image": "alpine:latest"
      },
      "state": "CONTAINER_RUNNING",
      "created_at": "2024-01-01T12:05:00Z"
    }
  ]
}
```

#### Create Container
```bash
POST /v1/runtime/container/create
Content-Type: application/json
{
  "pod_sandbox_id": "sandbox-123",
  "config": {
    "metadata": {
      "name": "test-container"
    },
    "image": {
      "image": "alpine:latest"
    },
    "command": ["/bin/sh"],
    "args": ["-c", "sleep 3600"]
  },
  "sandbox_config": {
    "metadata": {
      "name": "test-pod",
      "namespace": "default"
    }
  }
}

# Response
{
  "container_id": "container-789"
}
```

### Image Operations

#### List Images
```bash
POST /v1/image/list
Content-Type: application/json
{
  "filter": {}
}

# Response
{
  "images": [
    {
      "id": "sha256:abc123...",
      "repo_tags": ["alpine:latest"],
      "repo_digests": [],
      "size": 5242880,
      "uid": {
        "value": "123"
      }
    }
  ]
}
```

#### Pull Image
```bash
POST /v1/image/pull
Content-Type: application/json
{
  "image": {
    "image": "alpine:latest"
  },
  "auth": {},
  "sandbox_config": {}
}

# Response
{
  "image_ref": "alpine:latest"
}
```

## Usage Examples

### Starting the CRI Server
```bash
# Start on default port 8080
servin cri start

# Start on custom port with verbose logging
servin cri start --port 9090 --verbose

# Start in background (Linux/macOS)
servin cri start --port 8080 &
```

### Testing with curl
```bash
# Health check
curl http://localhost:8080/health

# Get runtime version
curl -X POST http://localhost:8080/v1/runtime/version \
  -H "Content-Type: application/json" \
  -d '{}'

# List images
curl -X POST http://localhost:8080/v1/image/list \
  -H "Content-Type: application/json" \
  -d '{"filter": {}}'
```

### Integration with Kubernetes

To use Servin as a CRI runtime with Kubernetes:

1. **Start the CRI server**:
   ```bash
   servin cri start --port 8080
   ```

2. **Configure kubelet** to use Servin CRI:
   ```yaml
   # kubelet configuration
   apiVersion: kubelet.config.k8s.io/v1beta1
   kind: KubeletConfiguration
   containerRuntimeEndpoint: "http://localhost:8080"
   imageServiceEndpoint: "http://localhost:8080"
   ```

3. **Start kubelet** with the configuration:
   ```bash
   kubelet --config=kubelet-config.yaml
   ```

## Development and Testing

### Local Development
```bash
# Build Servin with CRI support
go build -o servin .

# Start CRI server for testing
./servin cri start --port 8080 --verbose

# In another terminal, test endpoints
./servin cri test
```

### Custom Testing
```bash
# Test specific endpoints
curl -X POST http://localhost:8080/v1/runtime/version -H "Content-Type: application/json" -d '{}'
curl -X POST http://localhost:8080/v1/runtime/status -H "Content-Type: application/json" -d '{}'
curl -X POST http://localhost:8080/v1/image/list -H "Content-Type: application/json" -d '{"filter": {}}'
```

## Implementation Notes

### Simplified Architecture
- **HTTP instead of gRPC**: Easier to test and debug during development
- **JSON instead of protobuf**: Human-readable and simpler to work with
- **In-memory state**: Simplified state management for initial implementation

### Future Enhancements
- **gRPC Support**: Optional gRPC endpoints for full CRI compatibility
- **Persistent State**: File-based state persistence for server restarts  
- **Advanced Networking**: Integration with Servin's bridge networking
- **Security Integration**: User namespace and rootless container support
- **Metrics and Monitoring**: Prometheus metrics and health monitoring
- **Log Streaming**: Real-time log streaming for containers and pods

### Kubernetes Compatibility
The CRI implementation follows the Kubernetes CRI specification v1alpha2 and provides:
- Complete pod sandbox lifecycle management
- Full container lifecycle operations
- Image service operations
- Status and health reporting
- Error handling and status codes

This enables Servin to be used as a drop-in replacement for other container runtimes like containerd or CRI-O in Kubernetes environments.
