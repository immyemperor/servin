---
layout: default
title: Container Runtime Interface
permalink: /cri/
---

# 🔌 Container Runtime Interface (CRI)

Servin provides full Kubernetes Container Runtime Interface (CRI) support, enabling seamless integration with Kubernetes clusters as a drop-in replacement for other container runtimes.

## 🎯 Overview

The Container Runtime Interface (CRI) is a plugin interface that enables kubelet to use different container runtimes without needing to recompile. Servin implements the complete CRI v1alpha2 specification, providing:

- **Pod Sandbox Management** - Kubernetes pod lifecycle support
- **Container Runtime Service** - Container operations via gRPC
- **Image Service** - Image management for Kubernetes
- **Streaming Service** - Exec, attach, and port-forward support

## 🏗️ Architecture

### **CRI Server Architecture**
```
┌─────────────────────────────────────────────────────────────┐
│                        Kubernetes                           │
├─────────────────────────────────────────────────────────────┤
│                         kubelet                             │
├─────────────────────────────────────────────────────────────┤
│                     CRI gRPC API                            │
├─────────────────────────────────────────────────────────────┤
│                    Servin CRI Server                        │
│  ┌─────────────────┐ ┌─────────────────┐ ┌───────────────┐ │
│  │ Runtime Service │ │  Image Service  │ │ Stream Server │ │
│  │                 │ │                 │ │               │ │
│  │ • Pod Sandbox   │ │ • Pull Images   │ │ • Exec        │ │
│  │ • Containers    │ │ • List Images   │ │ • Attach      │ │
│  │ • Status        │ │ • Remove Images │ │ • Port Forward│ │
│  │ • Events        │ │ • Image Status  │ │ • Streaming   │ │
│  └─────────────────┘ └─────────────────┘ └───────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                 Servin Container Engine                     │
└─────────────────────────────────────────────────────────────┘
```

### **Component Interaction**
1. **kubelet** communicates with Servin via gRPC on port 10010
2. **Runtime Service** handles pod and container lifecycle operations
3. **Image Service** manages container images for Kubernetes
4. **Stream Server** provides exec, attach, and port-forward functionality
5. **Servin Engine** executes the actual container operations

## 🚀 Getting Started

### **CRI Server Configuration**
```yaml
# /etc/servin/cri-config.yaml
apiVersion: v1
kind: ServinCRIConfig
metadata:
  name: default
spec:
  # CRI Server Settings
  criSocket: unix:///var/run/servin/servin.sock
  criPort: 10010
  
  # Runtime Configuration
  runtime:
    name: servin
    version: v1.0.0
    
  # Image Service Settings
  imageService:
    defaultRegistry: docker.io
    pullTimeout: 300s
    parallelPulls: 3
    
  # Networking
  networkPlugin: cni
  cniConfigDir: /etc/cni/net.d
  cniBinDir: /opt/cni/bin
  
  # Pod Sandbox Settings
  podSandbox:
    imageRepository: k8s.gcr.io/pause
    imageTag: "3.6"
    
  # Security Settings
  securityContext:
    runAsUser: true
    readOnlyRootFilesystem: false
    privileged: false
```

### **Starting CRI Server**
```bash
# Start Servin with CRI enabled
servin daemon --enable-cri

# Start with custom CRI configuration
servin daemon --enable-cri --cri-config /etc/servin/cri-config.yaml

# Start with specific CRI socket
servin daemon --enable-cri --cri-socket unix:///var/run/servin/servin.sock

# Enable debug logging for CRI
servin daemon --enable-cri --log-level debug --cri-log-level trace
```

### **kubelet Configuration**
Configure kubelet to use Servin as the container runtime:

```yaml
# /etc/kubernetes/kubelet/kubelet-config.yaml
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
containerRuntimeEndpoint: unix:///var/run/servin/servin.sock
imageServiceEndpoint: unix:///var/run/servin/servin.sock
runtimeRequestTimeout: "15m"
imageGCHighThresholdPercent: 85
imageGCLowThresholdPercent: 80
```

## 🏃‍♂️ Runtime Service

### **Pod Sandbox Management**

#### **RunPodSandbox**
Creates and starts a pod sandbox environment:

```go
// CRI API Request
message RunPodSandboxRequest {
    PodSandboxConfig config = 1;
    string runtime_handler = 2;
}

// Example Pod Configuration
{
  "metadata": {
    "name": "nginx-pod",
    "namespace": "default",
    "uid": "12345678-1234-1234-1234-123456789012"
  },
  "hostname": "nginx-pod",
  "log_directory": "/var/log/pods/default_nginx-pod_12345678",
  "dns_config": {
    "servers": ["10.96.0.10"],
    "searches": ["default.svc.cluster.local", "svc.cluster.local", "cluster.local"]
  },
  "port_mappings": [
    {
      "protocol": 1,
      "container_port": 80,
      "host_port": 8080
    }
  ],
  "linux": {
    "cgroup_parent": "/kubepods/burstable/pod12345678-1234-1234-1234-123456789012"
  }
}
```

#### **StopPodSandbox**
Stops a running pod sandbox:

```bash
# Servin CRI stops all containers in the sandbox
# Performs graceful shutdown with SIGTERM
# Waits for grace period (default 30s)
# Sends SIGKILL if containers don't stop gracefully
```

#### **RemovePodSandbox**
Removes a stopped pod sandbox:

```bash
# Cleanup sequence:
# 1. Verify sandbox is stopped
# 2. Remove all containers in sandbox
# 3. Cleanup network configuration
# 4. Remove sandbox metadata
# 5. Cleanup storage mounts
```

#### **PodSandboxStatus**
Returns detailed sandbox status information:

```json
{
  "id": "sandbox123",
  "metadata": {
    "name": "nginx-pod",
    "namespace": "default",
    "uid": "12345678-1234-1234-1234-123456789012"
  },
  "state": "SANDBOX_READY",
  "created_at": 1642694400000000000,
  "network": {
    "ip": "10.244.0.5"
  },
  "linux": {
    "namespaces": {
      "options": {
        "ipc": "POD",
        "network": "POD",
        "pid": "CONTAINER"
      }
    }
  },
  "labels": {
    "io.kubernetes.pod.name": "nginx-pod",
    "io.kubernetes.pod.namespace": "default"
  }
}
```

### **Container Management**

#### **CreateContainer**
Creates a container within a pod sandbox:

```go
message CreateContainerRequest {
    string pod_sandbox_id = 1;
    ContainerConfig config = 2;
    PodSandboxConfig sandbox_config = 3;
}

// Container Configuration Example
{
  "metadata": {
    "name": "nginx-container"
  },
  "image": {
    "image": "nginx:1.21"
  },
  "command": ["/usr/sbin/nginx"],
  "args": ["-g", "daemon off;"],
  "working_dir": "/",
  "envs": [
    {
      "key": "PATH",
      "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
    }
  ],
  "mounts": [
    {
      "container_path": "/var/log/nginx",
      "host_path": "/var/log/pods/default_nginx-pod_12345/nginx",
      "readonly": false
    }
  ],
  "labels": {
    "io.kubernetes.container.name": "nginx-container",
    "io.kubernetes.pod.name": "nginx-pod"
  },
  "linux": {
    "resources": {
      "cpu_period": 100000,
      "cpu_quota": 50000,
      "memory_limit_in_bytes": 134217728
    },
    "security_context": {
      "run_as_user": {
        "value": 101
      },
      "readonly_root_filesystem": false
    }
  }
}
```

#### **StartContainer**
Starts a created container:

```bash
# Start sequence:
# 1. Validate container exists and is created
# 2. Setup networking (if not shared with sandbox)
# 3. Apply resource limits and security context
# 4. Start container process
# 5. Update container status to running
```

#### **StopContainer**
Stops a running container:

```bash
# Stop sequence:
# 1. Send SIGTERM to container process
# 2. Wait for timeout period (configurable, default 30s)
# 3. Send SIGKILL if still running
# 4. Update container status to exited
# 5. Preserve exit code and finish time
```

#### **RemoveContainer**
Removes a stopped container:

```bash
# Remove sequence:
# 1. Verify container is stopped
# 2. Cleanup container filesystem
# 3. Remove network configuration
# 4. Cleanup volume mounts
# 5. Remove container metadata
```

#### **ListContainers**
Lists containers with optional filtering:

```json
{
  "containers": [
    {
      "id": "container123",
      "pod_sandbox_id": "sandbox123",
      "metadata": {
        "name": "nginx-container"
      },
      "image": {
        "image": "nginx:1.21"
      },
      "image_ref": "sha256:abc123...",
      "state": "CONTAINER_RUNNING",
      "created_at": 1642694460000000000,
      "labels": {
        "io.kubernetes.container.name": "nginx-container"
      }
    }
  ]
}
```

#### **ContainerStatus**
Returns detailed container status:

```json
{
  "id": "container123",
  "metadata": {
    "name": "nginx-container"
  },
  "state": "CONTAINER_RUNNING",
  "created_at": 1642694460000000000,
  "started_at": 1642694461000000000,
  "finished_at": 0,
  "exit_code": 0,
  "image": {
    "image": "nginx:1.21"
  },
  "image_ref": "sha256:abc123...",
  "reason": "",
  "message": "",
  "labels": {
    "io.kubernetes.container.name": "nginx-container"
  },
  "mounts": [
    {
      "container_path": "/var/log/nginx",
      "host_path": "/var/log/pods/default_nginx-pod_12345/nginx",
      "readonly": false
    }
  ],
  "log_path": "/var/log/pods/default_nginx-pod_12345/nginx-container/0.log"
}
```

## 🖼️ Image Service

### **Image Operations**

#### **ListImages**
Lists available images with filtering support:

```json
{
  "images": [
    {
      "id": "sha256:abc123...",
      "repo_tags": ["nginx:1.21", "nginx:latest"],
      "repo_digests": ["nginx@sha256:def456..."],
      "size": 142000000,
      "uid": {
        "value": 0
      },
      "username": "",
      "spec": {
        "image": "nginx:1.21"
      }
    }
  ]
}
```

#### **ImageStatus**
Returns detailed image information:

```json
{
  "image": {
    "id": "sha256:abc123...",
    "repo_tags": ["nginx:1.21"],
    "repo_digests": ["nginx@sha256:def456..."],
    "size": 142000000
  },
  "info": {
    "created_at": 1642694400000000000,
    "config": {
      "env": ["PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"],
      "entrypoint": ["/docker-entrypoint.sh"],
      "cmd": ["nginx", "-g", "daemon off;"],
      "working_dir": "/",
      "exposed_ports": ["80/tcp"]
    }
  }
}
```

#### **PullImage**
Pulls an image from a registry:

```go
message PullImageRequest {
    ImageSpec image = 1;
    AuthConfig auth = 2;
    PodSandboxConfig sandbox_config = 3;
}

// Pull with authentication
{
  "image": {
    "image": "private-registry.com/myapp:v1.0.0"
  },
  "auth": {
    "username": "myuser",
    "password": "mypassword",
    "server_address": "private-registry.com"
  }
}
```

#### **RemoveImage**
Removes an image from local storage:

```bash
# Remove sequence:
# 1. Check if image is in use by containers
# 2. Remove image layers not shared with other images
# 3. Update image metadata
# 4. Cleanup orphaned layers
```

## 🌊 Streaming Service

### **Exec**
Executes commands in running containers:

```bash
# WebSocket connection for exec
ws://localhost:10250/exec/{token}

# Request Parameters
{
  "container_id": "container123",
  "cmd": ["/bin/bash"],
  "tty": true,
  "stdin": true,
  "stdout": true,
  "stderr": true
}
```

### **Attach**
Attaches to a running container:

```bash
# WebSocket connection for attach
ws://localhost:10250/attach/{token}

# Attach to container's main process
{
  "container_id": "container123",
  "tty": true,
  "stdin": true,
  "stdout": true,
  "stderr": true
}
```

### **PortForward**
Forwards ports from pods to the host:

```bash
# WebSocket connection for port forwarding
ws://localhost:10250/portforward/{token}

# Port forwarding configuration
{
  "pod_sandbox_id": "sandbox123",
  "port": [80, 443]
}
```

## ⚙️ CRI Configuration

### **Runtime Configuration**
```yaml
# Advanced CRI configuration
apiVersion: v1
kind: ServinCRIConfig
spec:
  # Runtime settings
  runtime:
    runtimeType: "servin"
    runtimeRoot: "/var/lib/servin"
    runtimeEngine: "native"
    
  # Container settings
  container:
    defaultRuntime: "runc"
    logDriver: "json-file"
    logMaxSize: "10MB"
    logMaxFiles: 3
    
  # Image settings
  image:
    defaultRegistry: "docker.io"
    registryMirrors:
      - "https://mirror.gcr.io"
      - "https://registry.docker-cn.com"
    insecureRegistries:
      - "localhost:5000"
    
  # Networking
  network:
    plugin: "cni"
    pluginBinDir: "/opt/cni/bin"
    pluginConfDir: "/etc/cni/net.d"
    
  # Security
  security:
    apparmor: true
    selinux: true
    seccomp: true
    noNewPrivileges: true
    
  # Resource management
  resources:
    cpuCfsPeriod: 100000
    cpuCfsQuota: -1
    memorySwap: -1
    oomScoreAdj: -999
```

### **CNI Integration**
```json
// CNI Configuration Example
{
  "cniVersion": "0.4.0",
  "name": "servin-bridge",
  "type": "bridge",
  "bridge": "servin0",
  "isDefaultGateway": true,
  "ipMasq": true,
  "hairpinMode": true,
  "ipam": {
    "type": "host-local",
    "subnet": "10.244.0.0/16",
    "routes": [
      { "dst": "0.0.0.0/0" }
    ]
  }
}
```

## 🔍 Monitoring and Debugging

### **CRI Logs**
```bash
# Enable CRI debug logging
servin daemon --cri-log-level debug

# View CRI server logs
journalctl -u servin -f | grep CRI

# Container logs for Kubernetes
crictl logs container123
crictl logs -f container123  # Follow logs
crictl logs --tail 100 container123  # Last 100 lines
```

### **CRI Tools Integration**
```bash
# Install crictl
curl -L https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.25.0/crictl-v1.25.0-linux-amd64.tar.gz | tar -xz
sudo mv crictl /usr/local/bin/

# Configure crictl for Servin
cat <<EOF > /etc/crictl.yaml
runtime-endpoint: unix:///var/run/servin/servin.sock
image-endpoint: unix:///var/run/servin/servin.sock
timeout: 10
debug: false
EOF

# List pods and containers
crictl pods
crictl ps
crictl images

# Pod operations
crictl runp pod-config.yaml
crictl stopp pod123
crictl rmp pod123

# Container operations
crictl create pod123 container-config.yaml pod-config.yaml
crictl start container123
crictl stop container123
crictl rm container123

# Image operations
crictl pull nginx:latest
crictl rmi nginx:latest

# Debugging
crictl exec -it container123 /bin/bash
crictl logs container123
crictl stats
```

### **Performance Monitoring**
```bash
# Monitor CRI server performance
servin cri stats

# Container resource usage
crictl stats
crictl stats container123

# System resource usage
servin system df
servin system events --filter type=container
```

## 🚀 Kubernetes Integration

### **Node Setup**
```bash
# Install Servin as CRI runtime
sudo servin service install --enable-cri

# Configure kubelet for Servin
sudo systemctl edit kubelet
# Add:
# [Service]
# Environment="KUBELET_EXTRA_ARGS=--container-runtime=remote --container-runtime-endpoint=unix:///var/run/servin/servin.sock"

# Restart services
sudo systemctl restart servin
sudo systemctl restart kubelet
```

### **Cluster Validation**
```bash
# Check node status
kubectl get nodes -o wide

# Verify container runtime
kubectl describe node | grep "Container Runtime"

# Run test pod
kubectl run test-pod --image=nginx:latest
kubectl get pod test-pod -o wide

# Check pod using crictl
crictl pods | grep test-pod
crictl ps | grep test-pod
```

### **Troubleshooting**
```bash
# Check CRI connectivity
crictl version

# Verify socket permissions
ls -la /var/run/servin/servin.sock

# Check firewall/SELinux
sudo firewall-cmd --list-all
sudo getenforce

# Restart services if needed
sudo systemctl restart servin
sudo systemctl restart kubelet
```

---

## 📚 Next Steps

- **[Container Management]({{ '/containers' | relative_url }})** - Learn container lifecycle operations
- **[Image Management]({{ '/images' | relative_url }})** - Understand image handling
- **[Configuration]({{ '/configuration' | relative_url }})** - Customize CRI settings
- **[Troubleshooting]({{ '/troubleshooting' | relative_url }})** - Resolve common issues

<div class="cri-reference">
  <h3>📖 CRI Specification</h3>
  <p>Servin implements the complete Kubernetes CRI v1alpha2 specification. For detailed API reference, see the <a href="https://github.com/kubernetes/cri-api">official CRI API documentation</a>.</p>
  
  <h3>🧪 Testing CRI Integration</h3>
  <ul>
    <li><strong>critest:</strong> Official CRI validation test suite</li>
    <li><strong>sonobuoy:</strong> Kubernetes conformance testing</li>
    <li><strong>e2e tests:</strong> End-to-end Kubernetes tests</li>
  </ul>
</div>
