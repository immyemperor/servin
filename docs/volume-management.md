---
layout: default
title: Volume Management
---

# Volume Management

Comprehensive guide to managing persistent storage with Servin Container Runtime.

## Volume Fundamentals

### Understanding Volumes

Volumes provide persistent storage that outlives container lifecycles:

- **Named Volumes**: Managed by Servin with automatic lifecycle
- **Bind Mounts**: Direct host filesystem mounts
- **Anonymous Volumes**: Temporary volumes tied to container lifecycle
- **tmpfs Mounts**: In-memory temporary filesystems

### Volume Types

```bash
# Named volume
servin run -v mydata:/app/data nginx:latest

# Bind mount
servin run -v /host/path:/container/path nginx:latest

# Anonymous volume
servin run -v /app/data nginx:latest

# tmpfs mount
servin run --tmpfs /tmp:rw,size=1g nginx:latest
```

## Creating and Managing Volumes

### Creating Volumes

Create named volumes for persistent storage:

```bash
# Create basic volume
servin volume create mydata

# Create volume with driver options
servin volume create mydata --driver local

# Create volume with labels
servin volume create mydata \
  --label project=webapp \
  --label environment=production

# Create volume with options
servin volume create mydata \
  --opt type=ext4 \
  --opt device=/dev/sdb1

# Create volume with size limit
servin volume create mydata \
  --opt size=10G
```

### Listing Volumes

View available volumes:

```bash
# List all volumes
servin volume ls

# List with detailed information
servin volume ls --format table

# Filter volumes by label
servin volume ls --filter label=project=webapp

# Filter volumes by driver
servin volume ls --filter driver=local

# Show dangling volumes
servin volume ls --filter dangling=true

# Custom output format
servin volume ls --format "{{.Name}}\t{{.Driver}}\t{{.Mountpoint}}"
```

### Inspecting Volumes

Get detailed volume information:

```bash
# Inspect volume details
servin volume inspect mydata

# Get specific information
servin volume inspect --format "{{.Mountpoint}}" mydata

# Inspect multiple volumes
servin volume inspect mydata logs cache

# Show volume configuration
servin volume inspect --format "{{.Options}}" mydata
```

## Volume Operations

### Mounting Volumes

Attach volumes to containers:

```bash
# Mount named volume
servin run -v mydata:/app/data nginx:latest

# Mount multiple volumes
servin run \
  -v mydata:/app/data \
  -v logs:/app/logs \
  -v config:/app/config \
  nginx:latest

# Mount with read-only access
servin run -v mydata:/app/data:ro nginx:latest

# Mount with specific mount options
servin run -v mydata:/app/data:rw,z nginx:latest

# Mount bind mount
servin run -v /host/data:/app/data nginx:latest
```

### Volume Drivers

Use different storage drivers:

```bash
# Local driver (default)
servin volume create mydata --driver local

# Network storage driver
servin volume create shared-data --driver nfs \
  --opt addr=nfs-server.company.com \
  --opt export=/shared/data

# Cloud storage driver
servin volume create cloud-data --driver cloudstor \
  --opt size=100GB \
  --opt type=ssd

# Encrypted volume
servin volume create secure-data --driver encrypted \
  --opt key=myencryptionkey
```

## Bind Mounts

### Host Directory Mounting

Mount host directories directly:

```bash
# Basic bind mount
servin run -v /host/path:/container/path nginx:latest

# Bind mount with read-only access
servin run -v /host/config:/app/config:ro nginx:latest

# Bind mount with specific options
servin run -v /host/data:/app/data:rw,z nginx:latest

# Multiple bind mounts
servin run \
  -v /host/app:/app \
  -v /host/config:/app/config:ro \
  -v /host/logs:/app/logs \
  nginx:latest
```

### File-level Mounting

Mount individual files:

```bash
# Mount configuration file
servin run -v /host/nginx.conf:/etc/nginx/nginx.conf:ro nginx:latest

# Mount secret files
servin run \
  -v /host/ssl/cert.pem:/etc/ssl/cert.pem:ro \
  -v /host/ssl/key.pem:/etc/ssl/key.pem:ro \
  nginx:latest

# Mount script files
servin run -v /host/scripts/startup.sh:/app/startup.sh:ro \
  nginx:latest
```

## tmpfs Mounts

### Memory-based Storage

Use in-memory filesystems:

```bash
# Basic tmpfs mount
servin run --tmpfs /tmp nginx:latest

# tmpfs with size limit
servin run --tmpfs /tmp:rw,size=512m nginx:latest

# tmpfs with specific options
servin run --tmpfs /tmp:rw,size=1g,uid=1000,gid=1000 nginx:latest

# Multiple tmpfs mounts
servin run \
  --tmpfs /tmp:rw,size=512m \
  --tmpfs /var/cache:rw,size=256m \
  nginx:latest
```

### Performance Optimization

Use tmpfs for temporary data:

```bash
# Cache directory in memory
servin run --tmpfs /app/cache:rw,size=1g nginx:latest

# Session storage in memory
servin run --tmpfs /app/sessions:rw,size=512m nginx:latest

# Temporary uploads
servin run --tmpfs /app/uploads/temp:rw,size=2g nginx:latest
```

## Volume Backup and Restore

### Backing Up Volumes

Create backups of volume data:

```bash
# Backup volume to tar file
servin run --rm \
  -v mydata:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/mydata-backup.tar.gz -C /data .

# Backup with timestamp
servin run --rm \
  -v mydata:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/mydata-$(date +%Y%m%d-%H%M%S).tar.gz -C /data .

# Backup multiple volumes
servin run --rm \
  -v data:/volumes/data \
  -v logs:/volumes/logs \
  -v config:/volumes/config \
  -v $(pwd):/backup \
  alpine tar czf /backup/full-backup.tar.gz -C /volumes .
```

### Restoring Volumes

Restore data from backups:

```bash
# Restore volume from backup
servin run --rm \
  -v mydata:/data \
  -v $(pwd):/backup \
  alpine tar xzf /backup/mydata-backup.tar.gz -C /data

# Restore with ownership
servin run --rm \
  -v mydata:/data \
  -v $(pwd):/backup \
  alpine sh -c "tar xzf /backup/mydata-backup.tar.gz -C /data && chown -R 1000:1000 /data"

# Restore specific files
servin run --rm \
  -v mydata:/data \
  -v $(pwd):/backup \
  alpine tar xzf /backup/mydata-backup.tar.gz -C /data ./specific-file.txt
```

### Automated Backup Scripts

Create automated backup solutions:

```bash
#!/bin/bash
# volume-backup.sh

BACKUP_DIR="/backups"
DATE=$(date +%Y%m%d-%H%M%S)

# Backup application data
servin run --rm \
  -v app-data:/data \
  -v $BACKUP_DIR:/backup \
  alpine tar czf /backup/app-data-$DATE.tar.gz -C /data .

# Backup database
servin run --rm \
  -v db-data:/data \
  -v $BACKUP_DIR:/backup \
  alpine tar czf /backup/db-data-$DATE.tar.gz -C /data .

# Clean old backups (keep 7 days)
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete

echo "Backup completed: $DATE"
```

## Volume Cleanup

### Removing Volumes

Clean up unused volumes:

```bash
# Remove specific volume
servin volume rm mydata

# Force remove volume
servin volume rm --force mydata

# Remove multiple volumes
servin volume rm data logs cache

# Remove dangling volumes
servin volume prune

# Remove all unused volumes
servin volume prune --all

# Remove with confirmation
servin volume prune --filter label=project=old
```

### Volume Pruning

Automated volume cleanup:

```bash
# Remove volumes not used by containers
servin volume prune --force

# Remove volumes older than specific time
servin volume prune --filter until=72h

# Remove volumes with specific labels
servin volume prune --filter label=temporary=true

# Scheduled cleanup
crontab -e
# Add: 0 2 * * * /usr/local/bin/servin volume prune --force
```

## Volume Sharing

### Container-to-Container Sharing

Share volumes between containers:

```bash
# Create shared volume
servin volume create shared-data

# Container 1: Producer
servin run -d --name producer \
  -v shared-data:/app/output \
  myapp:producer

# Container 2: Consumer
servin run -d --name consumer \
  -v shared-data:/app/input \
  myapp:consumer

# Container 3: Monitor
servin run -d --name monitor \
  -v shared-data:/app/data:ro \
  myapp:monitor
```

### Volume Containers

Use containers as volume providers:

```bash
# Create data container
servin create --name data-container \
  -v /app/data \
  -v /app/config \
  alpine:latest

# Use volumes from data container
servin run --volumes-from data-container nginx:latest

# Share volumes between multiple containers
servin run --volumes-from data-container app1:latest
servin run --volumes-from data-container app2:latest
```

## Performance Optimization

### Storage Performance

Optimize volume performance:

```bash
# Use local SSD storage
servin volume create fast-data --driver local \
  --opt type=ext4 \
  --opt device=/dev/nvme0n1

# Configure I/O options
servin run -v data:/app/data \
  --device-read-bps /dev/sda:50mb \
  --device-write-bps /dev/sda:50mb \
  nginx:latest

# Use memory-backed storage for temporary data
servin run --tmpfs /app/temp:rw,size=2g nginx:latest
```

### Caching Strategies

Implement effective caching:

```bash
# Application cache volume
servin volume create app-cache
servin run -v app-cache:/app/cache nginx:latest

# Database cache
servin volume create db-cache
servin run -v db-cache:/var/lib/mysql/cache mysql:8

# Build cache for development
servin volume create build-cache
servin run -v build-cache:/app/node_modules node:16
```

## Advanced Volume Features

### Volume Plugins

Use third-party volume drivers:

```bash
# Install volume plugin
servin plugin install store/sumologic/docker-sumologic-syslog

# Create volume with plugin
servin volume create myvolume --driver sumologic

# List available plugins
servin plugin ls

# Configure plugin options
servin volume create logs --driver sumologic \
  --opt sumo-url=https://collectors.sumologic.com/...
```

### Volume Encryption

Implement encrypted volumes:

```bash
# Create encrypted volume
servin volume create secure-data --driver encrypted \
  --opt encryption-key=/path/to/key \
  --opt cipher=aes-256-xts

# Mount encrypted volume
servin run -v secure-data:/app/secure nginx:latest

# Rotate encryption keys
servin volume update secure-data --opt new-key=/path/to/new-key
```

### Network Volumes

Use network-attached storage:

```bash
# NFS volume
servin volume create nfs-data --driver nfs \
  --opt addr=nfs.company.com \
  --opt export=/shared/data \
  --opt version=4

# CIFS/SMB volume
servin volume create smb-data --driver cifs \
  --opt addr=smb.company.com \
  --opt share=data \
  --opt username=user \
  --opt password=pass

# S3-compatible storage
servin volume create s3-data --driver s3 \
  --opt endpoint=s3.company.com \
  --opt bucket=app-data \
  --opt region=us-west-2
```

## Monitoring and Troubleshooting

### Volume Monitoring

Monitor volume usage and performance:

```bash
# Check volume disk usage
servin system df

# Monitor volume I/O
servin stats --format "table {{.Name}}\t{{.BlockIO}}"

# Check volume health
servin volume inspect --format "{{.Status}}" mydata

# View volume events
servin system events --filter type=volume
```

### Troubleshooting

Common volume issues and solutions:

```bash
# Check volume permissions
servin run --rm -v mydata:/data alpine ls -la /data

# Fix ownership issues
servin run --rm -v mydata:/data alpine chown -R 1000:1000 /data

# Check volume mount points
servin inspect container-name --format "{{.Mounts}}"

# Verify volume driver
servin volume inspect mydata --format "{{.Driver}}"

# Test volume accessibility
servin run --rm -v mydata:/test alpine touch /test/test-file
```

## Integration Examples

### Database Persistence

Persistent database storage:

```bash
# PostgreSQL with persistent data
servin volume create postgres-data
servin run -d \
  --name postgres \
  -v postgres-data:/var/lib/postgresql/data \
  -e POSTGRES_DB=myapp \
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=pass \
  postgres:13

# MySQL with persistent data and configuration
servin volume create mysql-data
servin volume create mysql-config
servin run -d \
  --name mysql \
  -v mysql-data:/var/lib/mysql \
  -v mysql-config:/etc/mysql/conf.d \
  -e MYSQL_ROOT_PASSWORD=rootpass \
  -e MYSQL_DATABASE=myapp \
  mysql:8
```

### Web Application Stack

Complete web stack with volumes:

```bash
# Create volumes
servin volume create app-data
servin volume create nginx-config
servin volume create app-logs

# Nginx reverse proxy
servin run -d \
  --name nginx \
  -p 80:80 \
  -p 443:443 \
  -v nginx-config:/etc/nginx/conf.d \
  -v app-logs:/var/log/nginx \
  nginx:alpine

# Application server
servin run -d \
  --name app \
  -v app-data:/app/data \
  -v app-logs:/app/logs \
  myapp:latest

# Log aggregator
servin run -d \
  --name fluentd \
  -v app-logs:/fluentd/log \
  fluent/fluentd:latest
```

This comprehensive volume management guide covers all aspects of persistent storage in Servin, from basic volume operations to advanced network storage and monitoring techniques.
