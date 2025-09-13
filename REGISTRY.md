# Servin Registry Usage Examples

This document provides examples of how to use the servin registry functionality.

## Starting the Local Registry

```bash
# Start registry on default port 5000
servin registry start

# Start registry on custom port
servin registry start --port 5001

# Start registry with custom data directory
servin registry start --data-dir ./my-registry

# Start registry in background (planned feature)
servin registry start --detach
```

## Managing Images with Registry

```bash
# List available registries
servin registry list

# Push a local image to registry
servin registry push myapp:latest

# Push to specific registry
servin registry push myapp:v1.0 localhost:5001

# Pull image from registry
servin registry pull nginx:alpine

# Pull from specific registry
servin registry pull myapp:latest localhost:5001
```

## Authentication

```bash
# Login to a registry
servin registry login docker.io
servin registry login localhost:5001

# Login with credentials via flags
servin registry login docker.io -u username -p password

# Logout from registry
servin registry logout docker.io
```

## Registry Configuration

The registry system automatically creates a configuration file at:
- Linux: `~/.servin/registry-config.json`
- Windows: `%USERPROFILE%\.servin\registry-config.json`
- macOS: `~/.servin/registry-config.json`

### Example Configuration
```json
{
  "local_port": 5000,
  "local_data_dir": "/home/user/.servin/registry",
  "default_registry": "",
  "registries": {
    "docker": "docker.io",
    "local": "localhost:5000"
  },
  "credentials": {
    "docker.io": {
      "username": "myuser",
      "password": "encrypted_password",
      "email": "user@example.com"
    }
  },
  "insecure_registries": ["localhost:5000"],
  "certificate_dir": "/home/user/.servin/certs"
}
```

## Registry API

The local registry implements a Docker Registry HTTP API v2 compatible interface:

- `GET /v2/` - Registry information
- `GET /v2/_catalog` - List repositories
- `GET /v2/{name}/manifests/{tag}` - Get image manifest
- `PUT /v2/{name}/manifests/{tag}` - Push image manifest
- `DELETE /v2/{name}/manifests/{tag}` - Delete image
- `GET /health` - Health check
- `GET /info` - Registry information

## Integration with Compose

You can use the registry with compose files:

```yaml
version: '3.8'
services:
  web:
    image: localhost:5000/myapp:latest
    ports:
      - "8080:80"
```

Then push your image to the local registry:
```bash
servin registry push myapp:latest localhost:5000
servin compose up
```

## Remote Registry Support

The current implementation provides:
- ✅ Local registry server with Docker-compatible API
- ✅ Push/pull to/from local registry
- ✅ Authentication management
- ✅ Multi-registry configuration
- 🚧 Remote registry integration (Docker Hub, etc.) - In Progress

Future enhancements will include full Docker Hub and private registry support.
