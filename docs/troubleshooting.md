---
layout: default
title: Troubleshooting
permalink: /troubleshooting/
---

# Troubleshooting

Comprehensive troubleshooting guide for Servin Container Runtime issues, debugging techniques, and performance optimization.

## Common Issues

### Installation Issues

#### Binary Not Found

**Problem**: `servin: command not found`

**Solutions:**
```bash
# Check if binary is installed
which servin
which servind

# Check PATH
echo $PATH

# Add to PATH if installed in custom location
export PATH=$PATH:/usr/local/bin

# Reinstall if missing
curl -sSL https://get.servin.dev | sh
```

#### Permission Denied

**Problem**: `permission denied while trying to connect to servin daemon`

**Solutions:**
```bash
# Check if user is in servin group
groups $USER

# Add user to servin group
sudo usermod -aG servin $USER

# Logout and login again
# Or run with sudo (not recommended for production)
sudo servin ps

# Check daemon socket permissions
ls -la /var/run/servin.sock
```

#### Daemon Start Failure

**Problem**: `failed to start servin daemon`

**Solutions:**
```bash
# Check daemon status
systemctl status servind

# Check daemon logs
journalctl -u servind -f

# Check configuration
servin system info

# Validate configuration file
sudo servind --validate-config

# Check available ports
netstat -tulpn | grep :2375
netstat -tulpn | grep :2376

# Start daemon manually for debugging
sudo servind --debug --log-level debug
```

### Container Issues

#### Container Won't Start

**Problem**: Container fails to start with various errors

**Diagnosis:**
```bash
# Check container status
servin ps -a

# Inspect container
servin inspect container-name

# Check container logs
servin logs container-name

# Check image exists
servin images | grep image-name

# Check resource availability
df -h
free -h
```

**Common Solutions:**
```bash
# Insufficient disk space
servin system prune -a
servin volume prune

# Memory issues
servin run --memory 512m image-name

# Port already in use
servin run -p 8081:80 image-name  # Use different port

# Missing image
servin pull image-name

# Invalid command
servin run image-name /bin/bash  # Use valid command
```

#### Container Exits Immediately

**Problem**: Container starts and exits immediately

**Diagnosis:**
```bash
# Check exit code
servin ps -a

# Check logs for errors
servin logs container-name

# Check image entrypoint
servin inspect image-name

# Test image manually
servin run -it --entrypoint /bin/sh image-name
```

**Solutions:**
```bash
# Keep container running
servin run -d image-name sleep infinity

# Use interactive mode for debugging
servin run -it image-name /bin/bash

# Override entrypoint
servin run --entrypoint /bin/bash image-name

# Check application configuration
servin run -e DEBUG=1 image-name
```

#### Container Performance Issues

**Problem**: Container running slowly or consuming too many resources

**Diagnosis:**
```bash
# Check resource usage
servin stats container-name

# Monitor system resources
top
htop
iotop

# Check container limits
servin inspect container-name | grep -i memory
servin inspect container-name | grep -i cpu

# Check disk I/O
iostat -x 1
```

**Solutions:**
```bash
# Set resource limits
servin run --memory 1g --cpu 0.5 image-name

# Use faster storage
servin run --tmpfs /tmp image-name

# Optimize application
servin run -e JAVA_OPTS="-Xmx512m" image-name

# Check network performance
servin run --network host image-name
```

### Image Issues

#### Image Pull Failures

**Problem**: Cannot pull images from registry

**Diagnosis:**
```bash
# Test connectivity
ping registry.company.com
curl -I https://registry.company.com/v2/

# Check authentication
servin login registry.company.com

# Check registry configuration
servin system info | grep -i registry

# Test with different registry
servin pull docker.io/library/alpine:latest
```

**Solutions:**
```bash
# Fix DNS resolution
echo "nameserver 8.8.8.8" >> /etc/resolv.conf

# Use insecure registry (development only)
sudo echo '{"insecure-registries": ["registry.company.com:5000"]}' > /etc/servin/daemon.json
sudo systemctl restart servind

# Fix authentication
servin logout registry.company.com
servin login registry.company.com

# Use registry mirror
sudo echo '{"registry-mirrors": ["https://mirror.company.com"]}' > /etc/servin/daemon.json
sudo systemctl restart servind
```

#### Image Build Failures

**Problem**: Image build fails with various errors

**Diagnosis:**
```bash
# Check Dockerfile syntax
servin build --no-cache .

# Build with debug output
servin build --progress=plain .

# Check build context size
du -sh .

# Check available disk space
df -h /var/lib/servin
```

**Solutions:**
```bash
# Clean build cache
servin builder prune

# Use .dockerignore
echo "*.log" >> .dockerignore
echo "node_modules" >> .dockerignore

# Reduce context size
servin build -f Dockerfile.prod .

# Fix network issues in build
servin build --network host .

# Use multi-stage builds
# See development.md for examples
```

### Network Issues

#### Container Connectivity Problems

**Problem**: Containers cannot communicate or reach external services

**Diagnosis:**
```bash
# Test container networking
servin exec container-name ping 8.8.8.8
servin exec container-name nslookup google.com

# Check network configuration
servin network ls
servin network inspect bridge

# Check iptables rules
sudo iptables -L -n
sudo iptables -t nat -L -n

# Check IP forwarding
cat /proc/sys/net/ipv4/ip_forward
```

**Solutions:**
```bash
# Enable IP forwarding
echo 'net.ipv4.ip_forward=1' >> /etc/sysctl.conf
sysctl -p

# Restart networking
sudo systemctl restart servind

# Use host networking for testing
servin run --network host image-name

# Create custom network
servin network create mynetwork
servin run --network mynetwork image-name

# Check DNS configuration
servin run --dns 8.8.8.8 image-name
```

#### Port Binding Issues

**Problem**: Cannot bind to host ports

**Diagnosis:**
```bash
# Check port availability
netstat -tulpn | grep :8080
ss -tulpn | grep :8080

# Check running containers
servin ps

# Check firewall rules
sudo ufw status
sudo iptables -L
```

**Solutions:**
```bash
# Use different port
servin run -p 8081:80 image-name

# Stop conflicting service
sudo systemctl stop apache2
sudo systemctl stop nginx

# Bind to specific interface
servin run -p 127.0.0.1:8080:80 image-name

# Use random port
servin run -P image-name
servin port container-name
```

### Storage Issues

#### Disk Space Problems

**Problem**: No space left on device

**Diagnosis:**
```bash
# Check disk usage
df -h
du -sh /var/lib/servin/*

# Check image sizes
servin images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"

# Check container sizes
servin ps -s

# Check volume usage
servin system df
```

**Solutions:**
```bash
# Clean up images
servin image prune -a

# Clean up containers
servin container prune

# Clean up volumes
servin volume prune

# Clean up everything
servin system prune -a --volumes

# Remove large images
servin rmi large-image:tag

# Configure log rotation
echo '{"log-driver":"json-file","log-opts":{"max-size":"10m","max-file":"3"}}' > /etc/servin/daemon.json
sudo systemctl restart servind
```

#### Volume Mount Issues

**Problem**: Volume mounts not working or permission denied

**Diagnosis:**
```bash
# Check mount points
servin exec container-name mount | grep /app/data

# Check permissions
servin exec container-name ls -la /app/data
ls -la /host/path

# Check SELinux context
ls -Z /host/path
```

**Solutions:**
```bash
# Fix permissions
sudo chown -R 1000:1000 /host/path
sudo chmod -R 755 /host/path

# Use bind mount with proper options
servin run -v /host/path:/app/data:Z image-name

# Create volume instead of bind mount
servin volume create myvolume
servin run -v myvolume:/app/data image-name

# Fix SELinux labels
sudo setsebool -P container_manage_cgroup 1
```

### Performance Issues

#### Slow Container Operations

**Problem**: Container operations are slow

**Diagnosis:**
```bash
# Check system load
uptime
top
iostat

# Check storage driver performance
servin system info | grep -i storage

# Check overlay2 usage
df -h /var/lib/servin/overlay2

# Monitor daemon logs
journalctl -u servind -f
```

**Solutions:**
```bash
# Optimize storage driver
echo '{"storage-driver":"overlay2","storage-opts":["overlay2.override_kernel_check=true"]}' > /etc/servin/daemon.json

# Use faster storage
# Move /var/lib/servin to SSD

# Increase file descriptors
echo 'DefaultLimitNOFILE=1048576' >> /etc/systemd/system.conf
systemctl daemon-reexec

# Tune kernel parameters
echo 'vm.max_map_count=262144' >> /etc/sysctl.conf
sysctl -p
```

#### High Resource Usage

**Problem**: Servin daemon using too many resources

**Diagnosis:**
```bash
# Check daemon resource usage
ps aux | grep servind
top -p $(pgrep servind)

# Check memory usage
cat /proc/$(pgrep servind)/status | grep Vm

# Check open files
lsof -p $(pgrep servind) | wc -l

# Monitor with htop
htop -p $(pgrep servind)
```

**Solutions:**
```bash
# Adjust daemon configuration
echo '{"max-concurrent-downloads":3,"max-concurrent-uploads":3}' > /etc/servin/daemon.json

# Reduce log level
echo '{"log-level":"warn"}' > /etc/servin/daemon.json

# Enable live restore
echo '{"live-restore":true}' > /etc/servin/daemon.json

# Restart daemon
sudo systemctl restart servind
```

## Debugging Techniques

### Daemon Debugging

#### Enable Debug Mode

```bash
# Start daemon in debug mode
sudo servind --debug --log-level debug

# Or configure in daemon.json
echo '{"debug":true,"log-level":"debug"}' > /etc/servin/daemon.json
sudo systemctl restart servind

# Watch debug logs
journalctl -u servind -f
```

#### Debug API Calls

```bash
# Enable API debugging
export SERVIN_DEBUG=1

# Use curl for API testing
curl -v -X GET http://localhost:2375/version

# Monitor API calls
sudo tcpdump -i lo -A 'port 2375'

# Use strace for system calls
sudo strace -p $(pgrep servind) -f -e trace=network
```

### Container Debugging

#### Interactive Debugging

```bash
# Start container in debug mode
servin run -it --entrypoint /bin/bash image-name

# Attach to running container
servin exec -it container-name /bin/bash

# Debug init process
servin run --init image-name

# Override entrypoint for debugging
servin run -it --entrypoint /bin/sh image-name
```

#### Process Debugging

```bash
# Check container processes
servin exec container-name ps aux

# Monitor process activity
servin exec container-name top

# Check process tree
servin exec container-name pstree

# Debug with strace
servin exec container-name strace -p 1

# Check process limits
servin exec container-name cat /proc/1/limits
```

#### Network Debugging

```bash
# Debug network connectivity
servin exec container-name ping 8.8.8.8
servin exec container-name nslookup google.com
servin exec container-name telnet host 80

# Check network interfaces
servin exec container-name ip addr show
servin exec container-name ip route show

# Monitor network traffic
servin exec container-name tcpdump -i eth0

# Check DNS resolution
servin exec container-name cat /etc/resolv.conf
servin exec container-name nslookup host.domain.com
```

### Log Analysis

#### Container Logs

```bash
# View logs with timestamps
servin logs --timestamps container-name

# Follow logs in real-time
servin logs -f container-name

# View logs from specific time
servin logs --since 2024-01-01T00:00:00Z container-name

# View last N lines
servin logs --tail 100 container-name

# Save logs to file
servin logs container-name > container.log
```

#### System Logs

```bash
# View daemon logs
journalctl -u servind

# Follow daemon logs
journalctl -u servind -f

# View logs from specific time
journalctl -u servind --since "1 hour ago"

# View logs with specific priority
journalctl -u servind -p err

# Export logs
journalctl -u servind --no-pager > servind.log
```

### Performance Profiling

#### CPU Profiling

```bash
# Enable profiling in daemon
echo '{"debug":true,"experimental":true}' > /etc/servin/daemon.json
sudo systemctl restart servind

# Profile CPU usage
go tool pprof http://localhost:2375/debug/pprof/profile?seconds=30

# Profile goroutines
go tool pprof http://localhost:2375/debug/pprof/goroutine

# Profile memory
go tool pprof http://localhost:2375/debug/pprof/heap
```

#### I/O Profiling

```bash
# Monitor disk I/O
iostat -x 1

# Monitor per-process I/O
iotop -p $(pgrep servind)

# Check file descriptor usage
lsof -p $(pgrep servind) | wc -l

# Monitor storage driver
watch -n 1 'df -h /var/lib/servin'
```

## Diagnostic Tools

### System Information Collection

```bash
#!/bin/bash
# collect-diagnostics.sh

echo "=== System Information ==="
uname -a
lsb_release -a
uptime

echo "=== Servin Version ==="
servin version
servind --version

echo "=== System Resources ==="
free -h
df -h
iostat

echo "=== Network Configuration ==="
ip addr show
ip route show
cat /etc/resolv.conf

echo "=== Servin Configuration ==="
cat /etc/servin/daemon.json
servin system info

echo "=== Running Containers ==="
servin ps -a

echo "=== Images ==="
servin images

echo "=== Networks ==="
servin network ls

echo "=== Volumes ==="
servin volume ls

echo "=== System Events ==="
servin system events --since 1h --until now

echo "=== Daemon Logs ==="
journalctl -u servind --since "1 hour ago" --no-pager

echo "=== Kernel Messages ==="
dmesg | tail -50
```

### Automated Health Checks

```bash
#!/bin/bash
# health-check.sh

# Check daemon status
if ! systemctl is-active --quiet servind; then
    echo "ERROR: Servin daemon is not running"
    exit 1
fi

# Check API connectivity
if ! curl -f http://localhost:2375/version >/dev/null 2>&1; then
    echo "ERROR: Cannot connect to Servin API"
    exit 1
fi

# Check disk space
USAGE=$(df /var/lib/servin | awk 'NR==2 {print $5}' | sed 's/%//')
if [ $USAGE -gt 90 ]; then
    echo "WARNING: Disk usage is at ${USAGE}%"
fi

# Check memory usage
MEM_USAGE=$(free | awk 'NR==2{printf "%.2f", $3*100/$2}')
if (( $(echo "$MEM_USAGE > 90" | bc -l) )); then
    echo "WARNING: Memory usage is at ${MEM_USAGE}%"
fi

# Test container creation
if ! servin run --rm alpine:latest echo "Health check OK" >/dev/null 2>&1; then
    echo "ERROR: Cannot create test container"
    exit 1
fi

echo "Health check passed"
```

## Support Resources

### Community Support

- **GitHub Issues**: [https://github.com/servin-dev/servin/issues](https://github.com/servin-dev/servin/issues)
- **Discussion Forum**: [https://github.com/servin-dev/servin/discussions](https://github.com/servin-dev/servin/discussions)
- **Discord Server**: [https://discord.gg/servin](https://discord.gg/servin)
- **Stack Overflow**: Tag questions with `servin`

### Documentation

- **Official Docs**: [https://servin.dev/docs](https://servin.dev/docs)
- **API Reference**: [https://servin.dev/docs/api](https://servin.dev/docs/api)
- **Examples**: [https://github.com/servin-dev/examples](https://github.com/servin-dev/examples)

### Bug Reports

When reporting bugs, include:

1. **Environment Information**
   - OS and version
   - Servin version
   - Kernel version
   - Hardware specifications

2. **Configuration**
   - Daemon configuration
   - Network setup
   - Storage configuration

3. **Reproduction Steps**
   - Exact commands used
   - Expected behavior
   - Actual behavior

4. **Logs and Output**
   - Daemon logs
   - Container logs
   - Error messages

5. **Diagnostic Information**
   - `servin system info`
   - `servin version`
   - System resource usage

### Emergency Procedures

#### Daemon Recovery

```bash
# Stop daemon gracefully
sudo systemctl stop servind

# Backup data
sudo tar -czf servin-backup-$(date +%Y%m%d).tar.gz /var/lib/servin

# Clean temporary files
sudo rm -rf /var/run/servin/*

# Reset daemon state
sudo servind --reset

# Start daemon
sudo systemctl start servind
```

#### Data Recovery

```bash
# Stop all containers
servin stop $(servin ps -q)

# Backup container data
sudo cp -r /var/lib/servin/containers /backup/

# Restore from backup
sudo tar -xzf servin-backup.tar.gz -C /

# Restart daemon
sudo systemctl restart servind
```

This comprehensive troubleshooting guide covers the most common issues, debugging techniques, and support resources for Servin Container Runtime.
