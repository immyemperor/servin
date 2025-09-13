# Servin - A Lightweight Container Runtime

A Go-based containerization tool that creates, runs, and manages isolated Linux containers using core Linux kernel features.

## 🚀 Features

- **Self-Contained**: Implements its own container runtime logic without external dependencies
- **Linux Native**: Uses low-level kernel features (namespaces, cgroups, chroot)
- **CLI-Driven**: Familiar command-line interface similar to Docker
- **Lightweight**: Minimal overhead and resource usage
- **Educational**: Perfect for understanding how containers work under the hood

## 📋 Prerequisites

- **Linux operating system** (Ubuntu, CentOS, etc.)
- **Go 1.19 or later**
- **Root privileges** (required for namespace and cgroup operations)
- **Kernel support** for namespaces and cgroups (Linux 3.8+)

## 🔧 Installation

### Quick Install

```bash
git clone <repository-url>
cd servin
make build-linux
make install
```

### Manual Build

```bash
# Clone and build
git clone <repository-url>
cd servin
go mod tidy
GOOS=linux GOARCH=amd64 go build -o servin .

# Install
sudo cp servin /usr/local/bin/
sudo chmod +x /usr/local/bin/servin
```

### Windows Development

If developing on Windows, use the PowerShell build script:
```powershell
.\examples\build.ps1
```

## 🎯 Usage

### Basic Commands

```bash
# Run a container
sudo servin run <image> <command> [args...]

# List running containers
sudo servin ls

# Stop a container
sudo servin stop <container-id>

# Remove a container
sudo servin rm <container-id>

# Execute command in running container
sudo servin exec <container-id> <command>
```

### Advanced Usage

```bash
# Run with resource limits
sudo servin run --memory 128m --cpus 0.5 alpine echo "Hello!"

# Run with custom hostname and name
sudo servin run --hostname myhost --name mycontainer alpine /bin/sh

# Run with volume mounts
sudo servin run --volume /host/path:/container/path alpine ls /container/path

# Run with environment variables
sudo servin run --env VAR1=value1 --env VAR2=value2 alpine env
```

## 🏗️ Architecture

Servin implements containerization through multiple isolation layers:

### 1. **Namespaces** (`pkg/namespaces/`)
- **PID**: Process isolation - container sees only its own processes
- **UTS**: Hostname isolation - container can have its own hostname
- **IPC**: Inter-process communication isolation
- **Network**: Network stack isolation
- **Mount**: Filesystem isolation

### 2. **Root Filesystem** (`pkg/rootfs/`)
- **chroot**: Changes the apparent root directory
- **Essential directories**: Creates `/bin`, `/etc`, `/proc`, etc.
- **File copying**: Copies essential binaries from host

### 3. **Control Groups** (`pkg/cgroups/`)
- **Memory limits**: Prevents memory exhaustion
- **CPU limits**: Controls CPU usage
- **Process limits**: Prevents fork bombs

### 4. **Container Management** (`pkg/container/`)
- **Lifecycle management**: Create, start, stop, cleanup
- **Configuration**: Handles all container settings
- **State tracking**: Monitors container status

## 📁 Project Structure

```
servin/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command and global flags
│   ├── run.go             # Container run command
│   ├── init.go            # Container initialization (internal)
│   ├── list.go            # List containers
│   ├── stop.go            # Stop containers
│   └── exec.go            # Execute in containers
├── pkg/                   # Core packages
│   ├── container/         # Container management
│   ├── namespaces/        # Linux namespace handling
│   ├── rootfs/            # Root filesystem management
│   ├── cgroups/           # Control groups (resource limits)
│   ├── images/            # Image management (future)
│   └── network/           # Network configuration (future)
├── examples/              # Examples and documentation
│   ├── README.md          # Detailed examples
│   ├── test_containers.sh # Test script
│   └── build.ps1          # Windows build script
├── main.go                # Entry point
├── Makefile              # Build automation
└── README.md             # This file
```

## 🔒 Security Considerations

### Root Privileges
- Servin requires root privileges for kernel features
- Always run on trusted systems with `sudo`
- Be cautious with container images and commands

### Isolation Guarantees
- **Process isolation**: Containers cannot see host processes
- **Filesystem isolation**: Containers cannot access host files (outside chroot)
- **Resource isolation**: Memory and CPU limits prevent resource exhaustion
- **Network isolation**: Containers have isolated network stack

### Security Best Practices
1. **Set resource limits** to prevent DoS attacks
2. **Use minimal images** to reduce attack surface
3. **Avoid privileged containers** unless absolutely necessary
4. **Keep runtime updated** with latest security patches

## 🐛 Troubleshooting

### Common Issues

**"Permission denied" errors:**
```bash
# Solution: Run with sudo
sudo servin run alpine echo "Hello!"
```

**"Namespace creation failed":**
```bash
# Check kernel namespace support
ls /proc/self/ns/
# Should show: pid, uts, ipc, net, mnt, user
```

**"CGroup mount failed":**
```bash
# Check cgroup filesystem
ls /sys/fs/cgroup/
# Should show: memory, cpu, pids, etc.
```

### Debugging

```bash
# Enable verbose output
sudo servin -v run alpine echo "Debug info"

# Check system requirements
uname -r  # Kernel version (need 3.8+)
cat /proc/version  # Detailed kernel info
```

## 🚧 Current Status

### ✅ Implemented Features
- [x] CLI framework with Cobra
- [x] Linux namespace isolation (PID, UTS, IPC, NET, Mount)
- [x] Root filesystem isolation with chroot
- [x] Control groups for resource limits
- [x] Container lifecycle management
- [x] Basic container operations (run, list, stop, exec)
- [x] Memory and CPU limit configuration
- [x] Cross-platform build support

### 🔄 In Progress
- [ ] Image management (tar-based format)
- [ ] Container state persistence
- [ ] Network configuration with veth pairs
- [ ] Volume mounting

### 📋 Future Features
- [ ] Container registry integration
- [ ] Docker-compatible image format
- [ ] Web-based management UI
- [ ] Container orchestration
- [ ] Advanced networking (bridges, port forwarding)

## 🧪 Testing

### Unit Tests
```bash
make test
# or
go test ./pkg/...
```

### Integration Tests (Linux only)
```bash
make test-containers
# or
sudo ./examples/test_containers.sh
```

## 🤝 Contributing

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Development Setup
```bash
make dev-setup  # Install dependencies
make lint       # Run code linting
make test       # Run tests
```

## 📚 Learning Resources

- [Linux Namespaces](https://man7.org/linux/man-pages/man7/namespaces.7.html)
- [Control Groups](https://www.kernel.org/doc/Documentation/cgroup-v1/)
- [Container Security](https://kubernetes.io/docs/concepts/security/)
- [How containers work](https://jvns.ca/blog/2016/10/10/what-even-is-a-container/)

## 📄 License

MIT License - see LICENSE file for details.

## 🙏 Acknowledgments

- Inspired by Docker, containerd, and runc
- Built for educational purposes and deep container understanding
- Thanks to the Go community for excellent system programming support

---

**⚠️ Disclaimer**: This is an educational project. While functional, it's not intended for production use. Use Docker or other mature container runtimes for production workloads.
