# Servin Examples

This directory contains examples and documentation for using the Servin container runtime.

## Basic Usage Examples

### 1. Running a Simple Command

```bash
# Run a simple echo command
sudo ./servin run alpine echo "Hello from container!"

# Run an interactive shell (when implemented)
sudo ./servin run alpine /bin/sh

# Run with resource limits
sudo ./servin run --memory 128m --name mycontainer alpine echo "Limited memory container"
```

### 2. Container Management

```bash
# List all containers
sudo ./servin ls

# Stop a running container
sudo ./servin stop mycontainer

# Execute command in running container
sudo ./servin exec mycontainer /bin/ls -la
```

## Example Container Setups

### 1. Basic Linux Container

Create a minimal Linux environment:

```bash
# Create a basic rootfs directory
mkdir -p /tmp/basic-rootfs/{bin,etc,proc,sys,tmp,var,usr/bin}

# Copy essential binaries
cp /bin/sh /tmp/basic-rootfs/bin/
cp /bin/ls /tmp/basic-rootfs/bin/
cp /bin/cat /tmp/basic-rootfs/bin/

# Run the container
sudo ./servin run basic /bin/sh
```

### 2. Development Container

Set up a container for development:

```bash
# Run with volume mount and environment variables
sudo ./servin run \
  --volume /home/user/code:/workspace \
  --env TERM=xterm \
  --workdir /workspace \
  --name devcontainer \
  alpine /bin/sh
```

## Architecture Overview

```
┌─────────────────────────────────────────┐
│                Servin                   │
├─────────────────────────────────────────┤
│  CLI Layer (Cobra)                     │
│  ├── run, stop, exec, ls commands      │
│  └── Argument parsing & validation     │
├─────────────────────────────────────────┤
│  Container Management                   │
│  ├── Container lifecycle               │
│  ├── Configuration management          │
│  └── State persistence                 │
├─────────────────────────────────────────┤
│  Core Isolation Components             │
│  ├── Namespaces (PID, UTS, IPC, etc.) │
│  ├── RootFS (chroot, filesystem)       │
│  ├── CGroups (resource limits)         │
│  └── Network (veth, bridges)           │
├─────────────────────────────────────────┤
│  Linux Kernel Features                 │
│  ├── clone() syscalls                  │
│  ├── mount() operations                │
│  ├── cgroup filesystem                 │
│  └── network namespaces                │
└─────────────────────────────────────────┘
```

## Security Considerations

### Root Privileges Required
- Servin requires root privileges to create namespaces and cgroups
- Always run with `sudo` on a trusted system
- Be cautious when running unknown container images

### Isolation Levels
- **Process Isolation**: PID namespaces prevent process visibility
- **Filesystem Isolation**: chroot prevents access to host filesystem
- **Network Isolation**: Network namespaces isolate network stack
- **Resource Isolation**: cgroups limit memory, CPU, and process count

### Best Practices
1. Always set resource limits to prevent resource exhaustion
2. Use minimal base images to reduce attack surface
3. Avoid running containers with --privileged flags
4. Regularly update the container runtime

## Troubleshooting

### Common Issues

1. **Permission Denied**
   ```
   Error: this command requires root privileges
   Solution: Run with sudo
   ```

2. **Namespace Creation Failed**
   ```
   Error: failed to create namespace
   Solution: Ensure kernel supports namespaces (Linux 3.8+)
   ```

3. **CGroup Mount Failed**
   ```
   Error: failed to create cgroup
   Solution: Ensure cgroup filesystem is mounted at /sys/fs/cgroup
   ```

### Debugging Tips

1. **Verbose Output**: Use `-v` flag for detailed logging
2. **Check Kernel Features**: Verify namespace and cgroup support
3. **Resource Monitoring**: Check system resources before running containers

## Performance Notes

### Memory Usage
- Each container has its own rootfs (isolated filesystem)
- Memory overhead is typically 2-10MB per container
- Resource limits should be set based on application needs

### CPU Performance
- Namespace creation adds minimal CPU overhead
- CGroup enforcement may impact performance under heavy load
- Network isolation adds latency for network operations

## Development

### Building from Source
```bash
cd servin
go build -o servin .
sudo mv servin /usr/local/bin/
```

### Running Tests
```bash
# Unit tests
go test ./pkg/...

# Integration tests (requires root)
sudo go test ./tests/...
```

### Contributing
1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

## References

- [Linux Namespaces](https://man7.org/linux/man-pages/man7/namespaces.7.html)
- [Control Groups (cgroups)](https://www.kernel.org/doc/Documentation/cgroup-v1/cgroups.txt)
- [Container Security](https://kubernetes.io/docs/concepts/security/)
- [Go System Programming](https://golang.org/pkg/syscall/)
