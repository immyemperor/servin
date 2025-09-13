# Volume Management Guide

This document provides comprehensive guidance on managing volumes in Servin.

## Overview

Volumes in Servin provide persistent storage that can be shared between containers and persist beyond container lifecycles. They are managed independently of containers and can be mounted into multiple containers simultaneously.

## Volume Commands

### List Volumes
```bash
# List all volumes
servin volume ls

# Example output:
VOLUME NAME    DRIVER  MOUNTPOINT                                    CREATED
data-volume    local   C:\Users\user\.servin\volumes\data-volume     2025-09-13 04:27:21
config-volume  local   C:\Users\user\.servin\volumes\config-volume   2025-09-13 04:27:30
```

### Create Volumes
```bash
# Create a simple volume
servin volume create myvolume

# Create volume with labels
servin volume create --label env=production --label app=web data-volume

# Create volume with driver options
servin volume create --driver local --opt size=10G large-volume
```

### Inspect Volumes
```bash
# Inspect a single volume
servin volume inspect data-volume

# Inspect multiple volumes
servin volume inspect volume1 volume2 volume3

# Example output:
Volume: data-volume
Driver: local
Mountpoint: C:\Users\user\.servin\volumes\data-volume
Created: 2025-09-13 04:27:21
Scope: local
Labels:
  app=web
  env=production
Status:
  state=ready
```

### Remove Volumes
```bash
# Remove a single volume
servin volume rm myvolume

# Remove multiple volumes
servin volume rm volume1 volume2 volume3

# Force remove volumes (ignore errors)
servin volume rm --force problematic-volume

# Remove ALL volumes (with confirmation)
servin volume rm-all

# Remove unused volumes
servin volume prune
```

## Using Volumes with Containers

### Mount Volumes in Containers
```bash
# Mount a named volume
servin run --volume myvolume:/data alpine ls /data

# Mount multiple volumes
servin run \
  --volume data-volume:/app/data \
  --volume config-volume:/app/config \
  myapp:latest

# Mix named volumes and bind mounts
servin run \
  --volume data-volume:/app/data \
  --volume /host/config:/app/config \
  myapp:latest
```

## Volume Drivers

Currently, Servin supports the following volume drivers:

### Local Driver (default)
- **Type**: `local`
- **Storage**: Local filesystem
- **Platforms**: Linux, Windows, macOS
- **Features**: Basic volume creation and management

## Platform-Specific Behavior

### Linux
```bash
# Volumes stored in: /var/lib/servin/volumes/
# Full file permissions and ownership support
# Can be mounted with specific user/group permissions
```

### Windows
```bash
# Volumes stored in: %USERPROFILE%\.servin\volumes\
# Windows file system permissions
# Suitable for development and testing
```

### macOS
```bash
# Volumes stored in: ~/.servin/volumes/
# Unix file permissions
# Suitable for development and testing
```

## Best Practices

### Volume Naming
```bash
# Use descriptive names
servin volume create postgres-data
servin volume create app-logs
servin volume create shared-config

# Use labels for organization
servin volume create --label project=webapp --label env=prod webapp-data
```

### Volume Lifecycle
```bash
# Create volumes before running containers
servin volume create app-data
servin run --volume app-data:/app/data myapp:latest

# Backup volume data before removal
# (copy files from mountpoint to backup location)
cp -r /var/lib/servin/volumes/app-data /backup/

# Clean up unused volumes periodically
servin volume prune
```

### Development Workflow
```bash
# Development environment setup
servin volume create dev-database
servin volume create dev-logs
servin volume create dev-cache

# Run development stack
servin run -d --name db --volume dev-database:/var/lib/postgresql postgres
servin run -d --name app --volume dev-logs:/app/logs --volume dev-cache:/app/cache myapp:dev

# Clean up development environment
servin rm app db
servin volume rm-all  # Remove all development volumes
```

## Volume Security

### Access Control
- Volumes inherit file system permissions from the host
- Container processes access volumes with their configured user/group
- Use appropriate labels to track sensitive volumes

### Data Protection
```bash
# Label sensitive volumes
servin volume create --label security=sensitive --label backup=required user-data

# Regular backups (manual for now)
tar -czf backup-$(date +%Y%m%d).tar.gz -C /var/lib/servin/volumes/user-data .
```

## Troubleshooting

### Common Issues

#### Volume Not Found
```bash
$ servin volume rm nonexistent
Error: volume 'nonexistent' not found

# Solution: Check volume name with 'servin volume ls'
```

#### Permission Denied
```bash
# Ensure proper permissions on volume directories
# Windows: Check user has access to .servin directory
# Linux/macOS: Ensure volume directory is accessible
```

#### Volume in Use
```bash
# Currently volumes can be removed even if in use
# Future versions will prevent removal of mounted volumes
```

### Volume Inspection
```bash
# Check volume details
servin volume inspect problematic-volume

# Check filesystem directly
ls -la /var/lib/servin/volumes/  # Linux/macOS
dir %USERPROFILE%\.servin\volumes\  # Windows
```

## Future Enhancements

### Planned Features
- **Usage tracking**: Detect which containers are using volumes
- **Remote drivers**: Support for network storage backends
- **Volume encryption**: Built-in encryption for sensitive data
- **Backup integration**: Automated backup and restore capabilities
- **Volume replication**: Multi-node volume synchronization

### Driver Extensibility
- Plugin system for custom volume drivers
- Cloud storage integration (AWS EBS, Azure Disk, GCP Persistent Disk)
- Network file system support (NFS, CIFS/SMB)

## Migration and Compatibility

### From Other Container Runtimes
```bash
# Import data from existing volumes
# 1. Copy data to Servin volume directory
# 2. Create volume with matching name
# 3. Update container configurations

# Example: Migrating from Docker
servin volume create legacy-data
cp -r /var/lib/docker/volumes/legacy-data/_data/* /var/lib/servin/volumes/legacy-data/
```

## API Reference

The volume management system is built on these core operations:

- `ListVolumes()` - Get all volumes
- `GetVolume(name)` - Get specific volume
- `CreateVolume(name, driver, options, labels)` - Create new volume
- `RemoveVolume(name, force)` - Remove volume
- `RemoveAllVolumes(force)` - Remove all volumes
- `SaveVolume(volume)` - Persist volume metadata
