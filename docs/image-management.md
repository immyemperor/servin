---
layout: default
title: Image Management
---

# Image Management

Complete guide to managing container images with Servin Container Runtime.

## Image Operations

### Pulling Images

Download images from registries:

```bash
# Pull latest version
servin pull nginx

# Pull specific tag
servin pull nginx:1.21-alpine

# Pull from specific registry
servin pull docker.io/library/nginx:latest

# Pull from private registry
servin pull registry.company.com/app:v1.0

# Pull with authentication
servin pull --username myuser --password mypass private-registry.com/app:latest

# Pull multiple images
servin pull nginx redis postgres:13
```

### Building Images

Create images from Dockerfiles:

```bash
# Build from current directory
servin build .

# Build with tag
servin build -t myapp:v1.0 .

# Build with multiple tags
servin build -t myapp:v1.0 -t myapp:latest .

# Build from specific Dockerfile
servin build -f Dockerfile.prod -t myapp:prod .

# Build with build arguments
servin build --build-arg NODE_ENV=production -t myapp:prod .

# Build with context from URL
servin build -t myapp:latest https://github.com/user/repo.git#main

# Build with no cache
servin build --no-cache -t myapp:v1.0 .
```

### Listing Images

View available images:

```bash
# List all images
servin images

# List images with detailed information
servin images --format table

# List specific images
servin images nginx

# Filter images
servin images --filter dangling=true
servin images --filter before=nginx:latest
servin images --filter since=nginx:1.20

# Show image digests
servin images --digests

# Show all tags for images
servin images --all
```

## Image Information

### Inspecting Images

Get detailed image information:

```bash
# Inspect image configuration
servin inspect nginx:latest

# Get specific information
servin inspect --format "{{.Config.Cmd}}" nginx:latest

# View image layers
servin inspect --format "{{.RootFS.Layers}}" nginx:latest

# Show image metadata
servin inspect --format "{{.Config.Labels}}" nginx:latest
```

### Image History

View image layer history:

```bash
# Show image history
servin history nginx:latest

# Show history without truncation
servin history --no-trunc nginx:latest

# Show history in human-readable format
servin history --human nginx:latest

# Show quiet output (only image IDs)
servin history --quiet nginx:latest
```

### Image Size Analysis

Analyze image sizes:

```bash
# Show disk usage by images
servin system df

# Show detailed disk usage
servin system df --verbose

# Analyze specific image layers
servin inspect --format "{{.Size}}" nginx:latest

# Compare image sizes
servin images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"
```

## Image Registry Operations

### Pushing Images

Upload images to registries:

```bash
# Push to default registry
servin push myapp:latest

# Push to specific registry
servin push registry.company.com/myapp:v1.0

# Push with authentication
servin push --username myuser --password mypass private-registry.com/myapp:latest

# Push all tags
servin push --all-tags myapp
```

### Registry Authentication

Manage registry credentials:

```bash
# Login to registry
servin login registry.company.com
servin login --username myuser registry.company.com

# Login with password from stdin
echo "mypassword" | servin login --username myuser --password-stdin registry.company.com

# Logout from registry
servin logout registry.company.com

# View stored credentials
servin system info | grep "Registry Mirrors"
```

### Registry Configuration

Configure registry settings:

```json
# /etc/servin/registry.json
{
  "registries": {
    "docker.io": {
      "mirrors": [
        "https://mirror1.company.com",
        "https://mirror2.company.com"
      ]
    },
    "registry.company.com": {
      "auth": {
        "username": "serviceaccount",
        "password": "token123"
      },
      "tls": {
        "insecure": false,
        "ca_file": "/etc/ssl/company-ca.pem"
      }
    }
  }
}
```

## Image Manipulation

### Tagging Images

Create tags for images:

```bash
# Tag image with new name
servin tag nginx:latest myapp:base

# Tag for different registry
servin tag myapp:latest registry.company.com/myapp:v1.0

# Create multiple tags
servin tag myapp:latest myapp:stable myapp:v1.0.0
```

### Saving and Loading Images

Export and import images:

```bash
# Save image to tar file
servin save nginx:latest > nginx.tar

# Save multiple images
servin save nginx:latest redis:6-alpine > images.tar

# Save with compression
servin save nginx:latest | gzip > nginx.tar.gz

# Load image from tar file
servin load < nginx.tar

# Load with verbose output
servin load --input nginx.tar

# Load from compressed file
gunzip -c nginx.tar.gz | servin load
```

### Image Import/Export

Import and export image filesystems:

```bash
# Export container filesystem
servin export container-name > filesystem.tar

# Import filesystem as image
cat filesystem.tar | servin import - myapp:imported

# Import with metadata
cat filesystem.tar | servin import - myapp:v1.0 \
  --change 'CMD ["nginx", "-g", "daemon off;"]' \
  --change 'EXPOSE 80'
```

## Image Cleanup

### Removing Images

Clean up images:

```bash
# Remove image
servin rmi nginx:old-version

# Force remove image
servin rmi --force nginx:latest

# Remove multiple images
servin rmi nginx:old redis:old postgres:old

# Remove by image ID
servin rmi a1b2c3d4e5f6

# Remove dangling images
servin image prune

# Remove all unused images
servin image prune --all

# Remove images with filter
servin image prune --filter until=24h
servin image prune --filter label=version=old
```

### Automated Cleanup

Schedule regular image cleanup:

```bash
# Remove unused images older than 24 hours
servin image prune --filter until=24h --force

# Remove all unused images
servin image prune --all --force

# Clean system (images, containers, networks, volumes)
servin system prune --all --force

# Clean with size limit
servin system prune --filter until=72h
```

## Advanced Image Operations

### Multi-architecture Images

Work with multi-platform images:

```bash
# Build for multiple architectures
servin buildx build --platform linux/amd64,linux/arm64 -t myapp:latest .

# Pull specific architecture
servin pull --platform linux/arm64 nginx:latest

# Inspect architecture information
servin inspect --format "{{.Architecture}}" nginx:latest

# List available platforms
servin buildx ls
```

### Image Signing and Verification

Secure image operations:

```bash
# Sign image with cosign
servin images --format "{{.ID}}" myapp:latest | cosign sign

# Verify image signature
cosign verify registry.company.com/myapp:latest

# Scan image for vulnerabilities
servin scan myapp:latest

# Generate SBOM (Software Bill of Materials)
servin sbom myapp:latest
```

## Image Optimization

### Layer Optimization

Optimize image layers:

```dockerfile
# Bad: Multiple RUN commands create multiple layers
RUN apt-get update
RUN apt-get install -y curl
RUN apt-get install -y vim
RUN rm -rf /var/lib/apt/lists/*

# Good: Single RUN command with cleanup
RUN apt-get update && \
    apt-get install -y curl vim && \
    rm -rf /var/lib/apt/lists/*
```

### Multi-stage Builds

Reduce image size with multi-stage builds:

```dockerfile
# Build stage
FROM node:16-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

# Production stage
FROM node:16-alpine AS production
WORKDIR /app
COPY --from=builder /app/node_modules ./node_modules
COPY . .
EXPOSE 3000
CMD ["npm", "start"]
```

### Base Image Selection

Choose optimal base images:

```dockerfile
# Use minimal base images
FROM alpine:3.16          # ~5MB
FROM node:16-alpine       # ~110MB vs node:16 (~900MB)
FROM nginx:alpine         # ~23MB vs nginx:latest (~140MB)

# Use distroless for security
FROM gcr.io/distroless/node:16

# Use scratch for static binaries
FROM scratch
COPY myapp /
ENTRYPOINT ["/myapp"]
```

## Image Scanning and Security

### Vulnerability Scanning

Scan images for security vulnerabilities:

```bash
# Scan image for vulnerabilities
servin scan nginx:latest

# Scan with specific scanner
servin scan --scanner trivy nginx:latest

# Generate scan report
servin scan --output json nginx:latest > scan-report.json

# Scan during build
servin build --scan -t myapp:latest .

# Set security policies
servin scan --policy security-policy.yaml nginx:latest
```

### Security Best Practices

Implement secure image practices:

```dockerfile
# Use specific versions, not latest
FROM node:16.17.0-alpine

# Run as non-root user
RUN addgroup -g 1001 -S nodejs
RUN adduser -S nextjs -u 1001
USER nextjs

# Use read-only filesystem
USER 1001:1001
RUN chmod -R 755 /app && chown -R 1001:1001 /app

# Remove unnecessary packages
RUN apk add --no-cache curl && \
    apk del curl

# Set security labels
LABEL security.scan=true
LABEL security.policy=strict
```

## Registry Management

### Private Registry Setup

Set up and manage private registries:

```yaml
# registry-config.yml
version: 0.1
log:
  level: info
storage:
  filesystem:
    rootdirectory: /var/lib/registry
http:
  addr: :5000
  secret: registry-secret-key
auth:
  htpasswd:
    realm: basic-realm
    path: /etc/registry/htpasswd
```

### Registry Mirroring

Configure registry mirrors:

```json
{
  "registry-mirrors": [
    "https://mirror.company.com",
    "https://backup-mirror.company.com"
  ],
  "insecure-registries": [
    "internal-registry.company.com:5000"
  ]
}
```

### Content Trust

Enable image signing and verification:

```bash
# Enable content trust
export SERVIN_CONTENT_TRUST=1

# Sign and push image
servin push myapp:latest

# Verify signed image
servin pull myapp:latest

# Rotate signing keys
servin trust key rotate myapp

# View trust metadata
servin trust inspect myapp:latest
```

## Image Monitoring

### Registry Metrics

Monitor registry usage:

```bash
# Show registry statistics
servin system df

# Monitor pull/push rates
servin system events --filter type=image

# Generate usage reports
servin images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}"

# Track image lifecycle
servin events --filter image=nginx --since 24h
```

### Automation Scripts

Automate image management:

```bash
#!/bin/bash
# automated-cleanup.sh

# Remove old images (older than 7 days)
servin image prune --filter until=168h --force

# Remove unused containers
servin container prune --force

# Remove unused networks
servin network prune --force

# Remove unused volumes
servin volume prune --force

# Generate cleanup report
echo "Cleanup completed at $(date)" >> /var/log/servin-cleanup.log
```

## Integration Examples

### CI/CD Pipeline

Integrate with CI/CD workflows:

```yaml
# .github/workflows/build.yml
name: Build and Push Image

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Servin
        run: |
          curl -sSL https://get.servin.dev | sh
          
      - name: Build image
        run: |
          servin build -t myapp:${{ github.sha }} .
          servin tag myapp:${{ github.sha }} myapp:latest
          
      - name: Scan image
        run: |
          servin scan myapp:latest
          
      - name: Push image
        run: |
          echo ${{ secrets.REGISTRY_PASSWORD }} | servin login --username ${{ secrets.REGISTRY_USERNAME }} --password-stdin
          servin push myapp:${{ github.sha }}
          servin push myapp:latest
```

This comprehensive image management guide covers all aspects of working with container images in Servin, from basic operations to advanced security and optimization techniques.
