---
layout: default
title: Container Management
permalink: /container-management/
---

# Container Management

Comprehensive container lifecycle management with Servin Container Runtime.

## Container Operations

### Creating Containers

Create containers from images with flexible configuration options:

```bash
# Basic container creation
servin create nginx:latest

# Create with custom name
servin create --name web-server nginx:latest

# Create with environment variables
servin create --env DATABASE_URL=postgres://localhost/app nginx:latest

# Create with port mappings
servin create --port 8080:80 nginx:latest

# Create with volume mounts
servin create --volume ./data:/app/data nginx:latest

# Create with resource limits
servin create --memory 512m --cpu 0.5 nginx:latest
```

### Starting Containers

Start containers with various options:

```bash
# Start a container
servin start web-server

# Start multiple containers
servin start web-server db-server cache-server

# Start with attached output
servin start --attach web-server

# Start in detached mode (default)
servin start --detach web-server
```

### Running Containers

Create and start containers in one command:

```bash
# Run basic container
servin run nginx:latest

# Run with interactive terminal
servin run -it ubuntu:latest /bin/bash

# Run in background
servin run -d --name web nginx:latest

# Run with complete configuration
servin run -d \
  --name production-web \
  --port 80:80 \
  --port 443:443 \
  --env NODE_ENV=production \
  --volume ./config:/app/config \
  --memory 1g \
  nginx:latest
```

## Container Status and Information

### Listing Containers

View running and stopped containers:

```bash
# List running containers
servin ps

# List all containers (running and stopped)
servin ps --all

# List with detailed information
servin ps --format table

# List specific containers
servin ps web-server db-server

# Filter containers
servin ps --filter status=running
servin ps --filter name=web*
```

### Container Inspection

Get detailed container information:

```bash
# Inspect container configuration
servin inspect web-server

# Get specific information
servin inspect --format "{{.State.Status}}" web-server

# Inspect multiple containers
servin inspect web-server db-server cache-server
```

### Container Statistics

Monitor resource usage:

```bash
# Real-time stats for all containers
servin stats

# Stats for specific containers
servin stats web-server db-server

# One-time stats (no streaming)
servin stats --no-stream

# Stats with custom format
servin stats --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"
```

## Container Control

### Stopping Containers

Gracefully stop containers:

```bash
# Stop a container
servin stop web-server

# Stop with timeout
servin stop --time 30 web-server

# Force stop (SIGKILL)
servin stop --force web-server

# Stop multiple containers
servin stop web-server db-server cache-server
```

### Restarting Containers

Restart containers with various policies:

```bash
# Restart a container
servin restart web-server

# Restart with timeout
servin restart --time 10 web-server

# Restart multiple containers
servin restart web-server db-server
```

### Pausing and Unpausing

Suspend and resume container execution:

```bash
# Pause container execution
servin pause web-server

# Resume container execution
servin unpause web-server

# Pause multiple containers
servin pause web-server worker-1 worker-2
```

### Killing Containers

Send signals to containers:

```bash
# Kill container (SIGKILL)
servin kill web-server

# Send specific signal
servin kill --signal SIGTERM web-server
servin kill --signal SIGUSR1 web-server

# Kill multiple containers
servin kill web-server worker-1 worker-2
```

## Container Interaction

### Executing Commands

Run commands inside containers:

```bash
# Execute command
servin exec web-server ls -la /app

# Interactive shell
servin exec -it web-server /bin/bash

# Execute as specific user
servin exec --user nginx web-server whoami

# Execute with environment variables
servin exec --env DEBUG=true web-server npm test

# Execute in working directory
servin exec --workdir /app web-server npm start
```

### Copying Files

Copy files between host and containers:

```bash
# Copy file to container
servin cp ./config.yml web-server:/app/config.yml

# Copy file from container
servin cp web-server:/app/logs/app.log ./logs/

# Copy directory to container
servin cp ./static/ web-server:/app/public/

# Copy directory from container
servin cp web-server:/app/generated/ ./output/
```

### Viewing Logs

Access container logs:

```bash
# View logs
servin logs web-server

# Follow logs (tail -f)
servin logs --follow web-server

# Show timestamps
servin logs --timestamps web-server

# Show last N lines
servin logs --tail 100 web-server

# Show logs since specific time
servin logs --since 2024-01-01T00:00:00Z web-server

# Show logs for multiple containers
servin logs web-server db-server
```

## Container Cleanup

### Removing Containers

Clean up containers:

```bash
# Remove stopped container
servin rm web-server

# Force remove running container
servin rm --force web-server

# Remove with volumes
servin rm --volumes web-server

# Remove multiple containers
servin rm web-server db-server cache-server

# Remove all stopped containers
servin container prune

# Remove containers with confirmation
servin container prune --filter until=24h
```

### Auto-removal

Configure containers for automatic cleanup:

```bash
# Remove container when it exits
servin run --rm ubuntu:latest echo "Hello World"

# Remove container after specific time
servin run --rm --stop-timeout 60 nginx:latest
```

## Container Networking

### Network Configuration

Configure container networking:

```bash
# Create container with custom network
servin run --network mynetwork nginx:latest

# Connect container to network
servin network connect mynetwork web-server

# Disconnect container from network
servin network disconnect mynetwork web-server

# Run with no network
servin run --network none alpine:latest
```

### Port Management

Manage container ports:

```bash
# Map single port
servin run -p 8080:80 nginx:latest

# Map multiple ports
servin run -p 80:80 -p 443:443 nginx:latest

# Map port to specific interface
servin run -p 127.0.0.1:8080:80 nginx:latest

# Map random port
servin run -P nginx:latest

# View port mappings
servin port web-server
```

## Resource Management

### CPU and Memory

Control resource allocation:

```bash
# Set memory limit
servin run --memory 512m nginx:latest

# Set CPU limit
servin run --cpu 0.5 nginx:latest

# Set both CPU and memory
servin run --memory 1g --cpu 1.0 nginx:latest

# Set CPU quota and period
servin run --cpu-quota 50000 --cpu-period 100000 nginx:latest
```

### Storage Configuration

Manage container storage:

```bash
# Set storage driver options
servin run --storage-opt size=10G nginx:latest

# Set temporary filesystem
servin run --tmpfs /tmp:rw,size=1g nginx:latest

# Set device mappings
servin run --device /dev/sda:/dev/xvda nginx:latest
```

## Health Monitoring

### Health Checks

Configure container health monitoring:

```bash
# Run with health check
servin run --health-cmd "curl -f http://localhost/health" \
  --health-interval 30s \
  --health-timeout 10s \
  --health-retries 3 \
  nginx:latest

# View health status
servin inspect --format "{{.State.Health.Status}}" web-server

# View health check logs
servin inspect --format "{{.State.Health.Log}}" web-server
```

## Container Templates

### Dockerfile Integration

Work with container definitions:

```dockerfile
# Example Dockerfile for Servin
FROM ubuntu:22.04

LABEL maintainer="admin@example.com"
LABEL version="1.0"

# Install dependencies
RUN apt-get update && apt-get install -y \
    nginx \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Configure application
COPY nginx.conf /etc/nginx/nginx.conf
COPY app/ /var/www/html/

# Set working directory
WORKDIR /var/www/html

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost/health || exit 1

# Start command
CMD ["nginx", "-g", "daemon off;"]
```

### Compose Integration

Use with container orchestration:

```yaml
# docker-compose.yml compatible with Servin
version: '3.8'

services:
  web:
    image: nginx:latest
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./html:/usr/share/nginx/html
    environment:
      - NGINX_HOST=localhost
      - NGINX_PORT=80
    depends_on:
      - database
    networks:
      - frontend
      - backend

  database:
    image: postgres:13
    environment:
      - POSTGRES_DB=myapp
      - POSTGRES_USER=appuser
      - POSTGRES_PASSWORD=secret
    volumes:
      - db_data:/var/lib/postgresql/data
    networks:
      - backend

volumes:
  db_data:

networks:
  frontend:
  backend:
```

## Advanced Features

### Container Commit

Create images from containers:

```bash
# Commit container to image
servin commit web-server myapp:v1.0

# Commit with metadata
servin commit --author "John Doe <john@example.com>" \
  --message "Add configuration files" \
  web-server myapp:v1.1
```

### Container Export/Import

Export and import containers:

```bash
# Export container to tar
servin export web-server > web-server.tar

# Import container from tar
cat web-server.tar | servin import - myapp:imported
```

### Container Diff

View container changes:

```bash
# Show changes in container
servin diff web-server

# Show specific types of changes
servin diff --type A web-server  # Added files
servin diff --type D web-server  # Deleted files
servin diff --type C web-server  # Changed files
```

## Best Practices

### Security

- Run containers as non-root users
- Use read-only root filesystems when possible
- Limit resource usage with CPU and memory constraints
- Use health checks for monitoring
- Regularly update base images
- Scan images for vulnerabilities

### Performance

- Use multi-stage builds to reduce image size
- Optimize Dockerfile layer caching
- Use appropriate restart policies
- Monitor resource usage with stats
- Clean up unused containers regularly

### Monitoring

- Implement comprehensive health checks
- Use structured logging
- Monitor resource consumption
- Set up alerts for container failures
- Use centralized log aggregation

## Integration Examples

### Web Application Stack

```bash
# Database container
servin run -d \
  --name postgres-db \
  --env POSTGRES_DB=webapp \
  --env POSTGRES_USER=appuser \
  --env POSTGRES_PASSWORD=secret \
  --volume db_data:/var/lib/postgresql/data \
  postgres:13

# Redis cache
servin run -d \
  --name redis-cache \
  --volume redis_data:/data \
  redis:6-alpine

# Web application
servin run -d \
  --name web-app \
  --port 80:8000 \
  --env DATABASE_URL=postgres://appuser:secret@postgres-db/webapp \
  --env REDIS_URL=redis://redis-cache:6379 \
  --link postgres-db \
  --link redis-cache \
  myapp:latest
```

## 🖥️ GUI Management Interface

### Desktop GUI Features

The Servin Desktop GUI provides a comprehensive visual interface for container management:

#### **Container Dashboard**
- **Live Status Display**: Real-time container status with automatic refresh
- **Responsive Table**: Container list with name, image, status, ports, and created time
- **Smart Actions**: Context-aware buttons that adapt to container state
  - **Running containers**: Stop, Restart, View Details
  - **Stopped containers**: Start, Delete, View Details
- **Quick Operations**: One-click container management without command line

#### **Container Details View**
- **Tabbed Interface**: Organized information across multiple tabs
  - **Logs**: Real-time log streaming with search and filtering
  - **Files**: Container filesystem browser and file operations
  - **Exec**: Interactive terminal sessions within containers
  - **Environment**: Environment variables and configuration display
  - **Volumes**: Volume mount information and management
  - **Network**: Network configuration and port mappings
  - **Stats**: Resource usage monitoring and performance metrics

#### **Interactive Terminal**
- **Auto-Connect**: Automatic connection to container shell
- **Realistic Prompt**: Enhanced shell prompt with proper user@container format
- **Command History**: Previous commands preserved during session
- **Session Management**: Robust connection handling and error recovery

#### **Real-time Log Streaming**
- **Live Updates**: Continuous log streaming from running containers
- **Persistent Display**: Log content persists when switching between tabs
- **Search Integration**: Searchable and scrollable log content
- **Error Handling**: Graceful handling of log retrieval failures

#### **Enhanced User Experience**
- **Responsive Design**: Adapts to different window sizes and screen densities
- **Dark Theme**: Modern dark interface with consistent styling
- **Mobile Support**: Touch-friendly interface for tablet use
- **Keyboard Navigation**: Full keyboard accessibility support

### Launching the GUI

Access the desktop interface through multiple methods:

```bash
# Launch GUI directly
servin-gui

# Launch through main CLI
servin gui

# Launch with custom port
servin gui --port 8080

# Launch in development mode
servin gui --dev
```

### Integration with CLI

The GUI seamlessly integrates with CLI operations:
- **Real-time Sync**: Changes made via CLI are immediately reflected in GUI
- **Bi-directional Control**: Perform operations through either interface
- **Consistent State**: Shared state management between CLI and GUI
- **API Compatibility**: GUI uses the same APIs as CLI commands

This comprehensive container management guide covers all aspects of working with containers in Servin, from command-line operations to modern GUI interfaces, providing flexibility for users of all preferences and workflows.
