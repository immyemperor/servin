---
layout: default
title: Command Line Interface
permalink: /cli/
---

# 💻 Command Line Interface (CLI)

The Servin CLI provides a comprehensive command-line interface for container management, image operations, and system administration. Designed for developers, DevOps engineers, and system administrators who prefer command-line workflows.

## 🚀 Getting Started

### **Installation Verification**
```bash
# Check Servin version
servin --version

# Display help information
servin --help

# Check service status
servin service status
```

### **Basic Configuration**
```bash
# Initialize Servin configuration
servin config init

# Set default registry
servin config set registry.default docker.io

# View current configuration
servin config list
```

### **Global Options**
```bash
# Verbose output for debugging
servin --verbose command
servin -v command

# Development mode (skip root checks)
servin --dev command

# Custom log level
servin --log-level debug command
servin --log-level warn command

# Custom log file location
servin --log-file /path/to/logfile command

# Combine global options
servin --verbose --dev --log-level debug containers ls
```

#### **Available Log Levels**
- **debug** - Detailed debugging information
- **info** - General information (default)
- **warn** - Warning messages only
- **error** - Error messages only

## 📦 Container Management

### **Container Lifecycle**

#### **Creating Containers**
```bash
# Create a simple container
servin containers create ubuntu:latest

# Create with custom name
servin containers create --name web-server nginx:latest

# Create with port mapping
servin containers create --name web-app -p 8080:80 nginx:latest

# Create with environment variables
servin containers create --name app -e NODE_ENV=production -e PORT=3000 node:16

# Create with volume mounting
servin containers create --name db -v /data:/var/lib/mysql mysql:8.0

# Create with resource limits
servin containers create --name limited --memory 512m --cpus 1.0 alpine:latest
```

#### **Running Containers**
```bash
# Run container in background
servin run -d --name web-server nginx:latest

# Run container interactively
servin run -it ubuntu:latest bash

# Run with automatic removal
servin run --rm alpine:latest echo "Hello World"

# Run with custom command
servin run ubuntu:latest ls -la /

# Run with working directory
servin run -w /app node:16 npm start

# Run with user specification
servin run --user 1000:1000 ubuntu:latest whoami
```

#### **Container Control**
```bash
# Start containers
servin containers start web-server
servin containers start web-server db-server  # Multiple containers

# Stop containers
servin containers stop web-server
servin containers stop --time 30 db-server    # Graceful shutdown with timeout

# Restart containers
servin containers restart web-server
servin containers restart --time 10 app       # Custom restart timeout

# Pause/Unpause containers
servin containers pause web-server
servin containers unpause web-server

# Kill containers
servin containers kill web-server
servin containers kill --signal SIGTERM app   # Custom signal
```

#### **Container Information**
```bash
# List all containers
servin containers ls

# List all containers (including stopped)
servin containers ls -a

# List with custom format
servin containers ls --format "table {{.ID}}\t{{.Names}}\t{{.Status}}\t{{.Ports}}"

# List with filters
servin containers ls --filter "status=running"
servin containers ls --filter "name=web"
servin containers ls --filter "ancestor=nginx"

# Container inspection
servin containers inspect web-server
servin containers inspect --format '{{.State.Status}}' web-server

# Container statistics
servin containers stats                        # Real-time stats for all
servin containers stats web-server db-server  # Specific containers
servin containers stats --no-stream          # One-time stats
```

#### **Container Logs**
```bash
# View container logs
servin containers logs web-server

# Follow log output
servin containers logs -f web-server

# Show timestamps
servin containers logs --timestamps web-server

# Show last N lines
servin containers logs --tail 100 web-server

# Show logs since timestamp
servin containers logs --since 2024-01-01T00:00:00Z web-server

# Show logs with details
servin containers logs --details web-server
```

#### **Container Cleanup**
```bash
# Remove stopped containers
servin containers rm web-server

# Force remove running container
servin containers rm -f web-server

# Remove multiple containers
servin containers rm web-server db-server app

# Remove all stopped containers
servin containers prune

# Remove containers with confirmation
servin containers prune --force
```

## 🖼️ Image Management

### **Image Operations**

#### **Pulling Images**
```bash
# Pull latest image
servin images pull ubuntu

# Pull specific tag
servin images pull ubuntu:20.04

# Pull from specific registry
servin images pull docker.io/library/nginx:latest

# Pull with authentication
servin images pull --username myuser private-registry.com/myapp:latest

# Pull all tags
servin images pull --all-tags ubuntu

# Pull with platform specification
servin images pull --platform linux/amd64 ubuntu:latest
```

#### **Building Images**
```bash
# Build from Buildfile (separate command)
servin build .

# Build with tag
servin build -t myapp:latest .

# Build with custom Buildfile name
servin build -f MyBuildfile -t myapp .

# Build with build arguments
servin build --build-arg NODE_ENV=production -t myapp .

# Build with no cache
servin build --no-cache -t myapp .

# Build with labels
servin build --label version=1.0.0 --label maintainer=team@company.com -t myapp .

# Build quietly (only show image ID)
servin build -q -t myapp .

# Alternative: Build using image subcommand
servin images build -t myapp:latest .
servin images build -f Dockerfile.prod -t myapp:prod .
servin images build --target production -t myapp:prod .
```

#### **Buildfile Format**
```dockerfile
# Example Buildfile (similar to Dockerfile)
FROM alpine:latest
RUN apk add --no-cache curl
COPY . /app
WORKDIR /app
CMD ["./app"]
```

#### **Image Information**
```bash
# List images
servin images ls

# List with filters
servin images ls --filter "dangling=true"
servin images ls --filter "reference=ubuntu:*"
servin images ls --filter "before=nginx:latest"

# Image inspection
servin images inspect ubuntu:latest
servin images inspect --format '{{.Config.Env}}' ubuntu:latest

# Image history
servin images history ubuntu:latest

# Image size and usage
servin system df
servin system df -v  # Verbose output
```

#### **Image Tagging**
```bash
# Tag image
servin images tag ubuntu:latest myregistry.com/ubuntu:latest

# Tag with multiple tags
servin images tag myapp:latest myapp:v1.0.0
servin images tag myapp:latest myregistry.com/myapp:stable
```

#### **Pushing Images**
```bash
# Push to registry
servin images push myapp:latest

# Push to specific registry
servin images push myregistry.com/myapp:latest

# Push with authentication
servin images push --username myuser myregistry.com/myapp:latest

# Push all tags
servin images push --all-tags myapp
```

#### **Image Cleanup**
```bash
# Remove image
servin images rm ubuntu:20.04

# Force remove image
servin images rm -f myapp:latest

# Remove multiple images
servin images rm ubuntu:20.04 nginx:latest alpine:3.14

# Remove dangling images
servin images prune

# Remove all unused images
servin images prune -a

# Remove images with force
servin images prune -f
```

## 💾 Volume Management

### **Volume Operations**

#### **Creating Volumes**
```bash
# Create simple volume
servin volumes create data-volume

# Create with driver
servin volumes create --driver local data-volume

# Create with options
servin volumes create --opt type=nfs --opt device=server:/path nfs-volume

# Create with labels
servin volumes create --label environment=production --label team=backend data-volume
```

#### **Volume Information**
```bash
# List volumes
servin volumes ls

# List with filters
servin volumes ls --filter "dangling=true"
servin volumes ls --filter "driver=local"
servin volumes ls --filter "label=environment=production"

# Volume inspection
servin volumes inspect data-volume
servin volumes inspect --format '{{.Mountpoint}}' data-volume
```

#### **Volume Usage**
```bash
# Mount volume in container
servin run -v data-volume:/data ubuntu:latest

# Mount with read-only access
servin run -v data-volume:/data:ro ubuntu:latest

# Mount host directory
servin run -v /host/path:/container/path ubuntu:latest

# Mount with specific options
servin run -v data-volume:/data:Z ubuntu:latest  # SELinux label
```

#### **Volume Cleanup**
```bash
# Remove volume
servin volumes rm data-volume

# Remove multiple volumes
servin volumes rm vol1 vol2 vol3

# Remove all unused volumes
servin volumes prune

# Remove with force
servin volumes prune -f
```

## 🌐 Network Management

### **Network Operations**

#### **Creating Networks**
```bash
# Create bridge network
servin networks create mynetwork

# Create with custom driver
servin networks create --driver bridge mynetwork

# Create with subnet
servin networks create --subnet 172.20.0.0/16 mynetwork

# Create with gateway
servin networks create --subnet 172.20.0.0/16 --gateway 172.20.0.1 mynetwork

# Create with options
servin networks create --opt com.docker.network.bridge.name=mybr0 mynetwork
```

#### **Network Information**
```bash
# List networks
servin networks ls

# List with filters
servin networks ls --filter "driver=bridge"
servin networks ls --filter "type=custom"

# Network inspection
servin networks inspect mynetwork
servin networks inspect --format '{{.IPAM.Config}}' mynetwork
```

#### **Network Usage**
```bash
# Connect container to network
servin networks connect mynetwork web-server

# Connect with IP address
servin networks connect --ip 172.20.0.10 mynetwork web-server

# Connect with alias
servin networks connect --alias webserver mynetwork web-server

# Disconnect from network
servin networks disconnect mynetwork web-server

# Run container with custom network
servin run --network mynetwork nginx:latest
```

#### **Network Cleanup**
```bash
# Remove network
servin networks rm mynetwork

# Remove multiple networks
servin networks rm net1 net2 net3

# Remove all unused networks
servin networks prune

# Remove with force
servin networks prune -f
```

## 🏪 Registry Operations

### **Registry Authentication**
```bash
# Login to registry
servin login docker.io
servin login --username myuser --password-stdin docker.io < password.txt

# Login to private registry
servin login myregistry.com
servin login --username myuser myregistry.com

# Logout from registry
servin logout docker.io
servin logout myregistry.com

# View stored credentials
servin system info | grep -A 10 "Registry Mirrors"
```

### **Image Distribution**
```bash
# Search for images
servin search ubuntu
servin search --limit 25 --filter stars=3 nginx

# Save image to tar file
servin images save -o ubuntu.tar ubuntu:latest
servin images save ubuntu:latest | gzip > ubuntu.tar.gz

# Load image from tar file
servin images load -i ubuntu.tar
cat ubuntu.tar.gz | gunzip | servin images load

# Export container as tar
servin containers export web-server -o webserver.tar

# Import tar as image
servin images import webserver.tar myapp:from-container
```

## ⚙️ System Management

### **Service Control**
```bash
# Service management
servin service start
servin service stop
servin service restart
servin service status

# Service installation (requires admin)
servin service install
servin service uninstall

# Enable auto-start
servin service enable
servin service disable
```

### **System Information**
```bash
# System information
servin system info
servin system info --format '{{.ServerVersion}}'

# System events
servin system events
servin system events --filter container=web-server
servin system events --filter type=image --since 24h

# Disk usage
servin system df
servin system df -v

# System prune
servin system prune              # Remove unused data
servin system prune -a           # Remove all unused data
servin system prune --volumes    # Include volumes
servin system prune -f           # Force without confirmation
```

### **Configuration Management**
```bash
# Configuration commands
servin config init               # Initialize config
servin config list               # Show all settings
servin config get registry.default
servin config set registry.default docker.io
servin config unset registry.default

# Configuration file location
servin config path              # Show config file path
servin config edit              # Edit config with default editor
```

## 🔧 Advanced Features

### **Container Execution**
```bash
# Execute commands in running container
servin exec web-server ps aux
servin exec -it web-server bash
servin exec --user root web-server apt update

# Execute with environment variables
servin exec -e VAR=value web-server env

# Execute with working directory
servin exec -w /app web-server npm test
```

### **File Operations**
```bash
# Copy files to/from containers
servin cp file.txt web-server:/tmp/
servin cp web-server:/var/log/app.log ./app.log

# Copy with archive mode
servin cp -a folder/ web-server:/opt/

# Copy following symlinks
servin cp -L web-server:/etc/resolv.conf ./
```

### **Container Commit**
```bash
# Create image from container
servin containers commit web-server myapp:v1.0.0

# Commit with message and author
servin containers commit --message "Added configurations" --author "Developer <dev@company.com>" web-server myapp:v1.0.1

# Commit with changes
servin containers commit --change "ENV DEBUG=true" web-server myapp:debug
```

### **Resource Monitoring**
```bash
# Real-time container stats
servin stats
servin stats --no-stream
servin stats --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}"

# Container processes
servin top web-server
servin top web-server -eo pid,cmd

# Container port information
servin port web-server
servin port web-server 80
```

## 📋 Output Formatting

### **Format Options**
```bash
# Table format
servin containers ls --format "table {{.ID}}\t{{.Names}}\t{{.Status}}"

# JSON output
servin containers ls --format json
servin containers inspect --format json web-server

# Custom templates
servin containers ls --format "{{.Names}}: {{.Status}}"
servin images ls --format "{{.Repository}}:{{.Tag}} ({{.Size}})"

# CSV format for scripting
servin containers ls --format "{{.Names}},{{.Status}},{{.Ports}}"
```

### **Filtering Examples**
```bash
# Container filters
servin containers ls --filter "status=running"
servin containers ls --filter "name=web"
servin containers ls --filter "ancestor=nginx"
servin containers ls --filter "exited=0"
servin containers ls --filter "health=healthy"

# Image filters
servin images ls --filter "dangling=true"
servin images ls --filter "reference=ubuntu:*"
servin images ls --filter "before=nginx:latest"
servin images ls --filter "since=alpine:latest"
servin images ls --filter "label=version=1.0"

# Volume filters
servin volumes ls --filter "dangling=true"
servin volumes ls --filter "driver=local"
servin volumes ls --filter "label=environment=prod"

# Network filters
servin networks ls --filter "driver=bridge"
servin networks ls --filter "type=custom"
servin networks ls --filter "label=team=frontend"
```

## 🚀 Scripting and Automation

### **Shell Scripting Examples**
```bash
#!/bin/bash

# Stop and remove all containers
servin containers stop $(servin containers ls -aq)
servin containers rm $(servin containers ls -aq)

# Remove all images
servin images rm $(servin images ls -aq)

# Clean up everything
servin system prune -a -f --volumes

# Deploy application stack
servin run -d --name db -e MYSQL_ROOT_PASSWORD=secret mysql:8.0
servin run -d --name web --link db:database -p 80:80 nginx:latest

# Health check script
if servin containers inspect web-server --format '{{.State.Health.Status}}' | grep -q healthy; then
    echo "Container is healthy"
    exit 0
else
    echo "Container is unhealthy"
    exit 1
fi
```

### **PowerShell Examples**
```powershell
# Windows PowerShell scripting
$containers = servin containers ls --format json | ConvertFrom-Json
foreach ($container in $containers) {
    Write-Host "Container: $($container.Names) - Status: $($container.Status)"
}

# Bulk operations
$stopped = servin containers ls --filter "status=exited" --format "{{.Names}}"
if ($stopped) {
    servin containers rm $stopped.Split("`n")
}
```

## 🖥️ User Interface Commands

### **Desktop GUI**
```bash
# Launch desktop GUI application
servin gui

# Launch Terminal UI instead
servin gui --tui

# Launch GUI in development mode
servin gui --dev

# Custom port for web interface
servin gui --port 8081 --host localhost
```

### **Available GUI Features**
- **Container Management** - Visual container lifecycle control
- **Image Operations** - Import, remove, tag, and inspect images
- **CRI Server Control** - Start/stop Kubernetes CRI server
- **System Monitoring** - Real-time logs and status updates

## 🔗 Container Runtime Interface (CRI)

### **CRI Server Management**
```bash
# Start CRI HTTP server
servin cri start

# Start on custom port
servin cri start --port 9090

# Start with verbose logging
servin cri start --port 8080 -v

# Check CRI server status
servin cri status

# Test CRI endpoints
servin cri test
```

### **CRI Features**
- **Kubernetes Compatibility** - Full CRI v1alpha2 specification
- **HTTP Endpoints** - RESTful API at `/v1/runtime/` and `/v1/image/`
- **Pod Sandbox Operations** - Complete pod lifecycle management
- **Health Monitoring** - Built-in health checks at `/health`

## 🐋 Compose Orchestration

### **Multi-Service Applications**
```bash
# Start services from servin-compose.yml
servin compose up

# Start in detached mode
servin compose up -d

# Start with custom file
servin compose -f custom-compose.yml up

# Start with project name
servin compose --project-name myapp up

# Stop and remove services
servin compose down

# Stop and remove with volumes
servin compose down --volumes

# View service status
servin compose ps

# View service logs
servin compose logs

# Execute command in service
servin compose exec web bash
```

### **Compose File Format**
```yaml
# servin-compose.yml
version: '3'
services:
  web:
    image: nginx:latest
    ports:
      - "80:80"
    volumes:
      - ./html:/usr/share/nginx/html
  db:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: secret
    volumes:
      - db_data:/var/lib/mysql
volumes:
  db_data:
```

## 🔒 Security Management

### **Security Features Check**
```bash
# Check all security features
servin security check

# Check user namespace support only
servin security check --user-ns

# Show security configuration
servin security info

# Show security info for specific container
servin security info --container web-server
```

### **Security Features**
- **User Namespaces** - Enhanced privilege isolation
- **UID/GID Mapping** - Secure user mapping between host and container
- **Capability Management** - Fine-grained privilege control
- **Rootless Containers** - Run containers without root privileges

---

## 📚 Next Steps

- **[Terminal UI Guide]({{ '/tui' | relative_url }})** - Learn the interactive terminal interface
- **[Desktop GUI]({{ '/gui' | relative_url }})** - Explore the visual desktop application
- **[API Reference]({{ '/api' | relative_url }})** - Integrate with Servin programmatically
- **[Configuration]({{ '/configuration' | relative_url }})** - Customize Servin for your needs

<div class="cli-help">
  <h3>💡 Pro Tips</h3>
  <ul>
    <li>Use <code>servin COMMAND --help</code> for detailed command help</li>
    <li>Set up shell completion: <code>servin completion bash &gt; /etc/bash_completion.d/servin</code></li>
    <li>Use aliases for common commands: <code>alias sl='servin containers ls'</code></li>
    <li>Combine with tools like <code>jq</code> for JSON processing: <code>servin containers ls --format json | jq '.[0].Names'</code></li>
  </ul>
</div>
