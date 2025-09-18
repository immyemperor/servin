# Servin Dockerfile - Hybrid Deployment Mode

This Dockerfile allows running Servin as a containerized daemon while it manages VM-based containers.

## Use Cases

### ✅ **When to use this Dockerfile:**
- **Kubernetes Deployment:** Running Servin daemon as a Kubernetes pod
- **Hybrid Infrastructure:** Docker for orchestration, VMs for workloads  
- **Development/Testing:** Quick Servin daemon setup for testing
- **Service Mesh Integration:** Servin as a containerized service

### ❌ **When NOT to use this Dockerfile:**
- **Pure VM Mode:** Use native installation for best performance
- **Docker Replacement:** Use direct Servin installation instead
- **Single-host Development:** Use `./servin` binary directly

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Docker Host   │    │  Servin Container│    │   VM Workloads  │
│                 │───▶│    (daemon)      │───▶│   (containers)  │
│   Kubernetes    │    │                  │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## Usage

### Build and Run
```bash
# Build the image (done automatically by build.sh)
docker build -t servin:latest .

# Run Servin daemon in Docker
docker run -d \
  --name servin-daemon \
  --privileged \
  -v /var/run:/var/run \
  -p 10250:10250 \
  servin:latest

# Use Servin to manage VM-based containers
servin run nginx:alpine
```

### Kubernetes Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: servin-daemon
spec:
  replicas: 1
  selector:
    matchLabels:
      app: servin
  template:
    metadata:
      labels:
        app: servin
    spec:
      containers:
      - name: servin
        image: servin:latest
        ports:
        - containerPort: 10250
        securityContext:
          privileged: true
```

## Alternative: Pure VM Mode

For pure VM-based containerization without Docker:

```bash
# Install Servin natively
curl -L https://github.com/immyemperor/servin/releases/latest/download/servin-linux.tar.gz | tar xz
sudo mv servin /usr/local/bin/

# Use VM mode directly  
servin vm start
servin run nginx:alpine
```

## Conclusion

This Dockerfile enables **hybrid deployment** where Servin daemon runs in Docker but manages VM-based workloads. Choose based on your infrastructure needs:

- **Container Orchestration:** Use this Dockerfile
- **Pure VM Experience:** Use native installation