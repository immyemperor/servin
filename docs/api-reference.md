---
layout: default
title: API Reference
permalink: /api-reference/
---

# API Reference

Complete REST API documentation for Servin Container Runtime.

## API Overview

### API Architecture

Servin provides a comprehensive REST API compatible with Docker API v1.41+:

```
API Architecture:
┌─────────────────────────────────────────────────────────┐
│                    Client Applications                  │
├─────────────────────────────────────────────────────────┤
│                     HTTP/HTTPS                         │
├─────────────────────────────────────────────────────────┤
│                   Servin REST API                      │
├─────────────────────────────────────────────────────────┤
│                 Authentication Layer                   │
├─────────────────────────────────────────────────────────┤
│                Authorization Layer                     │
├─────────────────────────────────────────────────────────┤
│                  Servin Core Engine                    │
└─────────────────────────────────────────────────────────┘
```

### API Endpoints

Base URL: `http://localhost:2375` (insecure) or `https://localhost:2376` (TLS)

- **System**: System information and events
- **Containers**: Container lifecycle management
- **Images**: Image operations and management
- **Networks**: Network configuration and management
- **Volumes**: Volume operations and management
- **Exec**: Container command execution
- **Swarm**: (Future) Swarm mode management

## Authentication

### API Authentication

Servin supports multiple authentication methods:

#### No Authentication (Development)

```bash
# Configure daemon for no auth (development only)
{
  "hosts": ["tcp://0.0.0.0:2375"]
}
```

#### TLS Authentication

```bash
# Client certificate authentication
curl --cert client-cert.pem \
     --key client-key.pem \
     --cacert ca.pem \
     https://servin.example.com:2376/version
```

#### Bearer Token Authentication

```bash
# JWT token authentication
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
     https://servin.example.com:2376/version
```

#### API Key Authentication

```bash
# API key in header
curl -H "X-API-Key: your-api-key" \
     https://servin.example.com:2376/version
```

## System API

### System Information

Get system-wide information:

```http
GET /version
```

**Response:**
```json
{
  "Version": "1.0.0",
  "ApiVersion": "1.41",
  "MinAPIVersion": "1.12",
  "GitCommit": "a1b2c3d",
  "GoVersion": "go1.21.0",
  "Os": "linux",
  "Arch": "amd64",
  "KernelVersion": "5.15.0",
  "BuildTime": "2024-01-01T00:00:00Z"
}
```

```http
GET /info
```

**Response:**
```json
{
  "ID": "7TRN:IPZB:QYBB:VPBQ:UMPP:KARE:6ZNR:XE6T:7EWV:PKF4:ZOJD:TPYS",
  "Containers": 5,
  "ContainersRunning": 3,
  "ContainersPaused": 0,
  "ContainersStopped": 2,
  "Images": 12,
  "Driver": "overlay2",
  "DriverStatus": [
    ["Backing Filesystem", "ext4"],
    ["Supports d_type", "true"]
  ],
  "SystemStatus": null,
  "Plugins": {
    "Volume": ["local"],
    "Network": ["bridge", "host", "none", "overlay"],
    "Authorization": null,
    "Log": ["awslogs", "fluentd", "json-file", "journald", "splunk", "syslog"]
  },
  "MemoryLimit": true,
  "SwapLimit": true,
  "KernelMemory": true,
  "CpuCfsPeriod": true,
  "CpuCfsQuota": true,
  "CPUShares": true,
  "CPUSet": true,
  "PidsLimit": true,
  "IPv4Forwarding": true,
  "BridgeNfIptables": true,
  "BridgeNfIp6tables": true,
  "Debug": false,
  "NFd": 23,
  "OomKillDisable": true,
  "NGoroutines": 135,
  "SystemTime": "2024-01-01T12:00:00.000000000Z",
  "LoggingDriver": "json-file",
  "CgroupDriver": "systemd",
  "NEventsListener": 0,
  "KernelVersion": "5.15.0-generic",
  "OperatingSystem": "Ubuntu 22.04 LTS",
  "OSType": "linux",
  "Architecture": "x86_64",
  "IndexServerAddress": "https://index.docker.io/v1/",
  "RegistryConfig": {
    "AllowNondistributableArtifactsCIDRs": [],
    "AllowNondistributableArtifactsHostnames": [],
    "InsecureRegistryCIDRs": ["127.0.0.0/8"],
    "IndexConfigs": {
      "docker.io": {
        "Name": "docker.io",
        "Mirrors": [],
        "Secure": true,
        "Official": true
      }
    },
    "Mirrors": []
  },
  "NCPU": 4,
  "MemTotal": 8374108160,
  "GenericResources": null,
  "DockerRootDir": "/var/lib/servin",
  "HttpProxy": "",
  "HttpsProxy": "",
  "NoProxy": "",
  "Name": "servin-host",
  "Labels": ["environment=production"],
  "ExperimentalBuild": false,
  "ServerVersion": "1.0.0",
  "ClusterStore": "",
  "ClusterAdvertise": "",
  "Runtimes": {
    "runc": {
      "path": "/usr/bin/runc"
    }
  },
  "DefaultRuntime": "runc",
  "Swarm": {
    "NodeID": "",
    "NodeAddr": "",
    "LocalNodeState": "inactive",
    "ControlAvailable": false,
    "Error": "",
    "RemoteManagers": null
  },
  "LiveRestoreEnabled": true,
  "Isolation": "",
  "InitBinary": "servin-init",
  "ContainerdCommit": {
    "ID": "a1b2c3d4",
    "Expected": "a1b2c3d4"
  },
  "RuncCommit": {
    "ID": "a1b2c3d4",
    "Expected": "a1b2c3d4"
  },
  "InitCommit": {
    "ID": "a1b2c3d4",
    "Expected": "a1b2c3d4"
  },
  "SecurityOptions": [
    "name=apparmor",
    "name=seccomp,profile=default"
  ]
}
```

### System Events

Stream system events:

```http
GET /events
```

**Query Parameters:**
- `since`: Show events created since this timestamp
- `until`: Show events created until this timestamp
- `filters`: Filter events by type, container, image, etc.

**Response (Server-Sent Events):**
```json
{
  "Type": "container",
  "Action": "start",
  "Actor": {
    "ID": "container_id",
    "Attributes": {
      "image": "nginx:latest",
      "name": "web-server"
    }
  },
  "time": 1640995200,
  "timeNano": 1640995200000000000
}
```

### System Data Usage

Get filesystem usage information:

```http
GET /system/df
```

**Response:**
```json
{
  "LayersSize": 1234567890,
  "Images": [
    {
      "Id": "sha256:abc123...",
      "ParentId": "",
      "RepoTags": ["nginx:latest"],
      "RepoDigests": ["nginx@sha256:def456..."],
      "Created": 1640995200,
      "Size": 142123456,
      "SharedSize": 0,
      "VirtualSize": 142123456,
      "Labels": {},
      "Containers": 2
    }
  ],
  "Containers": [
    {
      "Id": "container_id",
      "Names": ["/web-server"],
      "Image": "nginx:latest",
      "ImageID": "sha256:abc123...",
      "Command": "nginx -g 'daemon off;'",
      "Created": 1640995200,
      "State": "running",
      "Status": "Up 5 minutes",
      "Ports": [{"PrivatePort": 80, "PublicPort": 8080, "Type": "tcp"}],
      "SizeRw": 1024,
      "SizeRootFs": 142124480,
      "Labels": {},
      "NetworkMode": "default",
      "Mounts": []
    }
  ],
  "Volumes": [
    {
      "Name": "my-volume",
      "Driver": "local",
      "Mountpoint": "/var/lib/servin/volumes/my-volume/_data",
      "CreatedAt": "2024-01-01T00:00:00Z",
      "Status": {},
      "Labels": {},
      "Scope": "local",
      "Options": {},
      "UsageData": {
        "Size": 1048576,
        "RefCount": 1
      }
    }
  ],
  "BuildCache": []
}
```

## Container API

### List Containers

List containers:

```http
GET /containers/json
```

**Query Parameters:**
- `all`: Show all containers (default shows just running)
- `limit`: Last n containers
- `size`: Show container sizes
- `filters`: Filter containers

**Response:**
```json
[
  {
    "Id": "container_id",
    "Names": ["/web-server"],
    "Image": "nginx:latest",
    "ImageID": "sha256:abc123...",
    "Command": "nginx -g 'daemon off;'",
    "Created": 1640995200,
    "Ports": [
      {
        "PrivatePort": 80,
        "PublicPort": 8080,
        "Type": "tcp",
        "IP": "0.0.0.0"
      }
    ],
    "SizeRw": 1024,
    "SizeRootFs": 142124480,
    "Labels": {
      "environment": "production"
    },
    "State": "running",
    "Status": "Up 5 minutes",
    "HostConfig": {
      "NetworkMode": "default"
    },
    "NetworkSettings": {
      "Networks": {
        "bridge": {
          "IPAMConfig": null,
          "Links": null,
          "Aliases": null,
          "NetworkID": "network_id",
          "EndpointID": "endpoint_id",
          "Gateway": "172.17.0.1",
          "IPAddress": "172.17.0.2",
          "IPPrefixLen": 16,
          "IPv6Gateway": "",
          "GlobalIPv6Address": "",
          "GlobalIPv6PrefixLen": 0,
          "MacAddress": "02:42:ac:11:00:02"
        }
      }
    },
    "Mounts": []
  }
]
```

### Create Container

Create a new container:

```http
POST /containers/create
```

**Request Body:**
```json
{
  "Image": "nginx:latest",
  "Cmd": ["nginx", "-g", "daemon off;"],
  "ExposedPorts": {
    "80/tcp": {}
  },
  "Env": [
    "NGINX_HOST=localhost",
    "NGINX_PORT=80"
  ],
  "HostConfig": {
    "PortBindings": {
      "80/tcp": [{"HostPort": "8080"}]
    },
    "RestartPolicy": {
      "Name": "unless-stopped",
      "MaximumRetryCount": 0
    },
    "LogConfig": {
      "Type": "json-file",
      "Config": {
        "max-size": "10m",
        "max-file": "5"
      }
    },
    "Binds": [
      "/host/path:/container/path:rw"
    ],
    "Memory": 536870912,
    "CpuShares": 512
  },
  "NetworkingConfig": {
    "EndpointsConfig": {
      "bridge": {
        "IPAMConfig": {
          "IPv4Address": "172.17.0.100"
        }
      }
    }
  },
  "Labels": {
    "environment": "production",
    "service": "web"
  }
}
```

**Response:**
```json
{
  "Id": "container_id",
  "Warnings": []
}
```

### Inspect Container

Get detailed container information:

```http
GET /containers/{id}/json
```

**Response:**
```json
{
  "Id": "container_id",
  "Created": "2024-01-01T00:00:00.000000000Z",
  "Path": "nginx",
  "Args": ["-g", "daemon off;"],
  "State": {
    "Status": "running",
    "Running": true,
    "Paused": false,
    "Restarting": false,
    "OOMKilled": false,
    "Dead": false,
    "Pid": 1234,
    "ExitCode": 0,
    "Error": "",
    "StartedAt": "2024-01-01T00:00:01.000000000Z",
    "FinishedAt": "0001-01-01T00:00:00Z",
    "Health": {
      "Status": "healthy",
      "FailingStreak": 0,
      "Log": [
        {
          "Start": "2024-01-01T00:01:00.000000000Z",
          "End": "2024-01-01T00:01:00.100000000Z",
          "ExitCode": 0,
          "Output": "OK"
        }
      ]
    }
  },
  "Image": "sha256:abc123...",
  "ResolvConfPath": "/var/lib/servin/containers/container_id/resolv.conf",
  "HostnamePath": "/var/lib/servin/containers/container_id/hostname",
  "HostsPath": "/var/lib/servin/containers/container_id/hosts",
  "LogPath": "/var/lib/servin/containers/container_id/container_id-json.log",
  "Name": "/web-server",
  "RestartCount": 0,
  "Driver": "overlay2",
  "Platform": "linux",
  "MountLabel": "",
  "ProcessLabel": "",
  "AppArmorProfile": "docker-default",
  "ExecIDs": null,
  "HostConfig": {
    "Binds": ["/host/path:/container/path:rw"],
    "ContainerIDFile": "",
    "LogConfig": {
      "Type": "json-file",
      "Config": {
        "max-file": "5",
        "max-size": "10m"
      }
    },
    "NetworkMode": "default",
    "PortBindings": {
      "80/tcp": [{"HostIp": "", "HostPort": "8080"}]
    },
    "RestartPolicy": {
      "Name": "unless-stopped",
      "MaximumRetryCount": 0
    },
    "AutoRemove": false,
    "VolumeDriver": "",
    "VolumesFrom": null,
    "CapAdd": null,
    "CapDrop": null,
    "CgroupnsMode": "host",
    "Dns": [],
    "DnsOptions": [],
    "DnsSearch": [],
    "ExtraHosts": null,
    "GroupAdd": null,
    "IpcMode": "private",
    "Cgroup": "",
    "Links": null,
    "OomScoreAdj": 0,
    "PidMode": "",
    "Privileged": false,
    "PublishAllPorts": false,
    "ReadonlyRootfs": false,
    "SecurityOpt": null,
    "UTSMode": "",
    "UsernsMode": "",
    "ShmSize": 67108864,
    "Runtime": "runc",
    "ConsoleSize": [0, 0],
    "Isolation": "",
    "CpuShares": 512,
    "Memory": 536870912,
    "NanoCpus": 0,
    "CgroupParent": "",
    "BlkioWeight": 0,
    "BlkioWeightDevice": [],
    "BlkioDeviceReadBps": null,
    "BlkioDeviceWriteBps": null,
    "BlkioDeviceReadIOps": null,
    "BlkioDeviceWriteIOps": null,
    "CpuPeriod": 0,
    "CpuQuota": 0,
    "CpuRealtimePeriod": 0,
    "CpuRealtimeRuntime": 0,
    "CpusetCpus": "",
    "CpusetMems": "",
    "Devices": [],
    "DeviceCgroupRules": null,
    "DeviceRequests": null,
    "KernelMemory": 0,
    "KernelMemoryTCP": 0,
    "MemoryReservation": 0,
    "MemorySwap": 0,
    "MemorySwappiness": null,
    "OomKillDisable": false,
    "PidsLimit": null,
    "Ulimits": null,
    "CpuCount": 0,
    "CpuPercent": 0,
    "IOMaximumIOps": 0,
    "IOMaximumBandwidth": 0,
    "MaskedPaths": [
      "/proc/asound",
      "/proc/acpi",
      "/proc/kcore",
      "/proc/keys",
      "/proc/latency_stats",
      "/proc/timer_list",
      "/proc/timer_stats",
      "/proc/sched_debug",
      "/proc/scsi",
      "/sys/firmware"
    ],
    "ReadonlyPaths": [
      "/proc/bus",
      "/proc/fs",
      "/proc/irq",
      "/proc/sys",
      "/proc/sysrq-trigger"
    ]
  },
  "GraphDriver": {
    "Data": {
      "LowerDir": "/var/lib/servin/overlay2/l/ABC123:...",
      "MergedDir": "/var/lib/servin/overlay2/merged",
      "UpperDir": "/var/lib/servin/overlay2/diff",
      "WorkDir": "/var/lib/servin/overlay2/work"
    },
    "Name": "overlay2"
  },
  "Mounts": [
    {
      "Type": "bind",
      "Source": "/host/path",
      "Destination": "/container/path",
      "Mode": "rw",
      "RW": true,
      "Propagation": "rprivate"
    }
  ],
  "Config": {
    "Hostname": "container_id",
    "Domainname": "",
    "User": "",
    "AttachStdin": false,
    "AttachStdout": true,
    "AttachStderr": true,
    "ExposedPorts": {
      "80/tcp": {}
    },
    "Tty": false,
    "OpenStdin": false,
    "StdinOnce": false,
    "Env": [
      "NGINX_HOST=localhost",
      "NGINX_PORT=80",
      "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
    ],
    "Cmd": ["nginx", "-g", "daemon off;"],
    "Healthcheck": {
      "Test": ["CMD", "curl", "-f", "http://localhost/health"],
      "Interval": 30000000000,
      "Timeout": 10000000000,
      "Retries": 3,
      "StartPeriod": 60000000000
    },
    "ArgsEscaped": true,
    "Image": "nginx:latest",
    "Volumes": null,
    "WorkingDir": "",
    "Entrypoint": null,
    "NetworkDisabled": false,
    "MacAddress": "",
    "OnBuild": null,
    "Labels": {
      "environment": "production",
      "service": "web"
    },
    "StopSignal": "SIGQUIT",
    "StopTimeout": 10,
    "Shell": ["/bin/sh", "-c"]
  },
  "NetworkSettings": {
    "Bridge": "",
    "SandboxID": "sandbox_id",
    "HairpinMode": false,
    "LinkLocalIPv6Address": "",
    "LinkLocalIPv6PrefixLen": 0,
    "Ports": {
      "80/tcp": [
        {
          "HostIp": "0.0.0.0",
          "HostPort": "8080"
        }
      ]
    },
    "SandboxKey": "/var/run/servin/netns/sandbox_id",
    "SecondaryIPAddresses": null,
    "SecondaryIPv6Addresses": null,
    "EndpointID": "endpoint_id",
    "Gateway": "172.17.0.1",
    "GlobalIPv6Address": "",
    "GlobalIPv6PrefixLen": 0,
    "IPAddress": "172.17.0.2",
    "IPPrefixLen": 16,
    "IPv6Gateway": "",
    "MacAddress": "02:42:ac:11:00:02",
    "Networks": {
      "bridge": {
        "IPAMConfig": null,
        "Links": null,
        "Aliases": null,
        "NetworkID": "network_id",
        "EndpointID": "endpoint_id",
        "Gateway": "172.17.0.1",
        "IPAddress": "172.17.0.2",
        "IPPrefixLen": 16,
        "IPv6Gateway": "",
        "GlobalIPv6Address": "",
        "GlobalIPv6PrefixLen": 0,
        "MacAddress": "02:42:ac:11:00:02",
        "DriverOpts": null
      }
    }
  }
}
```

### Container Operations

Start, stop, restart containers:

```http
POST /containers/{id}/start
POST /containers/{id}/stop
POST /containers/{id}/restart
POST /containers/{id}/kill
POST /containers/{id}/pause
POST /containers/{id}/unpause
```

**Query Parameters (stop/restart):**
- `t`: Number of seconds to wait before killing

**Query Parameters (kill):**
- `signal`: Signal to send (default: SIGKILL)

### Container Logs

Get container logs:

```http
GET /containers/{id}/logs
```

**Query Parameters:**
- `follow`: Follow log output
- `stdout`: Show stdout
- `stderr`: Show stderr
- `since`: Show logs since timestamp
- `until`: Show logs until timestamp
- `timestamps`: Show timestamps
- `tail`: Number of lines to show from end

### Container Stats

Get container resource usage statistics:

```http
GET /containers/{id}/stats
```

**Query Parameters:**
- `stream`: Stream stats continuously (default: true)

**Response:**
```json
{
  "read": "2024-01-01T12:00:00.000000000Z",
  "preread": "2024-01-01T11:59:59.000000000Z",
  "pids_stats": {
    "current": 10,
    "limit": 1024
  },
  "blkio_stats": {
    "io_service_bytes_recursive": [
      {
        "major": 8,
        "minor": 0,
        "op": "Read",
        "value": 1048576
      },
      {
        "major": 8,
        "minor": 0,
        "op": "Write",
        "value": 2097152
      }
    ],
    "io_serviced_recursive": [
      {
        "major": 8,
        "minor": 0,
        "op": "Read",
        "value": 100
      },
      {
        "major": 8,
        "minor": 0,
        "op": "Write",
        "value": 200
      }
    ]
  },
  "num_procs": 0,
  "storage_stats": {},
  "cpu_stats": {
    "cpu_usage": {
      "total_usage": 25902091,
      "percpu_usage": [15478871, 10423220],
      "usage_in_kernelmode": 5000000,
      "usage_in_usermode": 20000000
    },
    "system_cpu_usage": 1000000000,
    "online_cpus": 2,
    "throttling_data": {
      "periods": 0,
      "throttled_periods": 0,
      "throttled_time": 0
    }
  },
  "precpu_stats": {
    "cpu_usage": {
      "total_usage": 25000000,
      "percpu_usage": [15000000, 10000000],
      "usage_in_kernelmode": 4500000,
      "usage_in_usermode": 19500000
    },
    "system_cpu_usage": 999000000,
    "online_cpus": 2,
    "throttling_data": {
      "periods": 0,
      "throttled_periods": 0,
      "throttled_time": 0
    }
  },
  "memory_stats": {
    "usage": 104857600,
    "max_usage": 209715200,
    "stats": {
      "active_anon": 52428800,
      "active_file": 10485760,
      "cache": 20971520,
      "dirty": 0,
      "hierarchical_memory_limit": 536870912,
      "hierarchical_memsw_limit": 1073741824,
      "inactive_anon": 0,
      "inactive_file": 10485760,
      "mapped_file": 5242880,
      "pgfault": 12345,
      "pgmajfault": 123,
      "pgpgin": 6789,
      "pgpgout": 5432,
      "rss": 52428800,
      "rss_huge": 0,
      "total_active_anon": 52428800,
      "total_active_file": 10485760,
      "total_cache": 20971520,
      "total_dirty": 0,
      "total_inactive_anon": 0,
      "total_inactive_file": 10485760,
      "total_mapped_file": 5242880,
      "total_pgfault": 12345,
      "total_pgmajfault": 123,
      "total_pgpgin": 6789,
      "total_pgpgout": 5432,
      "total_rss": 52428800,
      "total_rss_huge": 0,
      "total_unevictable": 0,
      "total_writeback": 0,
      "unevictable": 0,
      "writeback": 0
    },
    "limit": 536870912
  },
  "name": "/web-server",
  "id": "container_id",
  "networks": {
    "eth0": {
      "rx_bytes": 1048576,
      "rx_packets": 1000,
      "rx_errors": 0,
      "rx_dropped": 0,
      "tx_bytes": 2097152,
      "tx_packets": 2000,
      "tx_errors": 0,
      "tx_dropped": 0
    }
  }
}
```

## Image API

### List Images

List images:

```http
GET /images/json
```

**Query Parameters:**
- `all`: Show all images (default hides intermediate images)
- `filters`: Filter images
- `digests`: Show image digests

**Response:**
```json
[
  {
    "Id": "sha256:abc123...",
    "ParentId": "",
    "RepoTags": [
      "nginx:latest",
      "nginx:1.21"
    ],
    "RepoDigests": [
      "nginx@sha256:def456..."
    ],
    "Created": 1640995200,
    "Size": 142123456,
    "VirtualSize": 142123456,
    "SharedSize": 0,
    "Labels": {
      "maintainer": "NGINX Docker Maintainers <docker-maint@nginx.com>"
    },
    "Containers": 2
  }
]
```

### Build Image

Build image from Dockerfile:

```http
POST /build
```

**Query Parameters:**
- `dockerfile`: Dockerfile name (default: Dockerfile)
- `t`: Repository name and tag
- `extrahosts`: Extra hosts to add
- `remote`: Remote context URL
- `q`: Suppress output
- `nocache`: Do not use cache
- `cachefrom`: Images to consider as cache sources
- `pull`: Always pull newer version of base image
- `rm`: Remove intermediate containers
- `forcerm`: Always remove intermediate containers
- `memory`: Memory limit for build
- `memswap`: Total memory limit (memory + swap)
- `cpushares`: CPU shares
- `cpusetcpus`: CPUs to use
- `cpuperiod`: CPU period
- `cpuquota`: CPU quota
- `buildargs`: JSON map of build arguments
- `shmsize`: Size of shared memory
- `squash`: Squash layers
- `labels`: JSON map of labels
- `networkmode`: Network mode
- `platform`: Platform for build
- `target`: Build target stage
- `outputs`: Build outputs

**Request Body:** tar archive containing Dockerfile and context

### Pull Image

Pull image from registry:

```http
POST /images/create
```

**Query Parameters:**
- `fromImage`: Image name
- `fromSrc`: Source to import
- `repo`: Repository name
- `tag`: Tag name
- `message`: Commit message
- `platform`: Platform for image

### Push Image

Push image to registry:

```http
POST /images/{name}/push
```

**Query Parameters:**
- `tag`: Tag to push

**Headers:**
- `X-Registry-Auth`: Base64 encoded auth config

### Image Operations

Tag, remove images:

```http
POST /images/{name}/tag
DELETE /images/{name}
```

**Query Parameters (tag):**
- `repo`: Repository name
- `tag`: Tag name

**Query Parameters (remove):**
- `force`: Force removal
- `noprune`: Do not delete untagged parents

## Volume API

### List Volumes

List volumes:

```http
GET /volumes
```

**Query Parameters:**
- `filters`: Filter volumes

**Response:**
```json
{
  "Volumes": [
    {
      "CreatedAt": "2024-01-01T00:00:00Z",
      "Driver": "local",
      "Labels": {
        "environment": "production"
      },
      "Mountpoint": "/var/lib/servin/volumes/my-volume/_data",
      "Name": "my-volume",
      "Options": {},
      "Scope": "local",
      "Status": {},
      "UsageData": {
        "RefCount": 1,
        "Size": 1048576
      }
    }
  ],
  "Warnings": []
}
```

### Create Volume

Create a new volume:

```http
POST /volumes/create
```

**Request Body:**
```json
{
  "Name": "my-volume",
  "Driver": "local",
  "DriverOpts": {
    "type": "ext4",
    "device": "/dev/sdb1"
  },
  "Labels": {
    "environment": "production"
  }
}
```

### Volume Operations

Inspect, remove volumes:

```http
GET /volumes/{name}
DELETE /volumes/{name}
```

**Query Parameters (remove):**
- `force`: Force removal

## Network API

### List Networks

List networks:

```http
GET /networks
```

**Query Parameters:**
- `filters`: Filter networks

### Create Network

Create a new network:

```http
POST /networks/create
```

**Request Body:**
```json
{
  "Name": "my-network",
  "Driver": "bridge",
  "IPAM": {
    "Driver": "default",
    "Config": [
      {
        "Subnet": "172.20.0.0/16",
        "Gateway": "172.20.0.1"
      }
    ]
  },
  "Options": {
    "com.docker.network.bridge.enable_icc": "true",
    "com.docker.network.bridge.enable_ip_masquerade": "true"
  },
  "Labels": {
    "environment": "production"
  }
}
```

### Network Operations

Connect/disconnect containers, inspect, remove networks:

```http
POST /networks/{id}/connect
POST /networks/{id}/disconnect
GET /networks/{id}
DELETE /networks/{id}
```

## Exec API

### Create Exec

Create exec instance:

```http
POST /containers/{id}/exec
```

**Request Body:**
```json
{
  "AttachStdin": true,
  "AttachStdout": true,
  "AttachStderr": true,
  "DetachKeys": "ctrl-p,ctrl-q",
  "Tty": true,
  "Cmd": ["/bin/bash"],
  "Env": ["VAR=value"],
  "User": "root",
  "Privileged": false,
  "WorkingDir": "/app"
}
```

### Start Exec

Start exec instance:

```http
POST /exec/{id}/start
```

**Request Body:**
```json
{
  "Detach": false,
  "Tty": true
}
```

### Inspect Exec

Get exec instance information:

```http
GET /exec/{id}/json
```

## Error Handling

### HTTP Status Codes

- `200 OK`: Successful operation
- `201 Created`: Resource created successfully
- `204 No Content`: Successful operation with no content
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Access denied
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource conflict
- `500 Internal Server Error`: Server error

### Error Response Format

```json
{
  "message": "Error description",
  "code": "ERROR_CODE",
  "details": {
    "field": "Additional error details"
  }
}
```

## API Clients

### Official Clients

- **Go**: `github.com/servin/servin-go`
- **Python**: `servin-py`
- **Node.js**: `servin-node`
- **Java**: `servin-java`

### Example Usage

#### Python Client

```python
import servin

client = servin.APIClient(base_url='http://localhost:2375')

# List containers
containers = client.containers.list()

# Create container
container = client.containers.create(
    image='nginx:latest',
    name='web-server',
    ports={'80/tcp': 8080}
)

# Start container
container.start()

# Get logs
logs = container.logs(follow=True)
```

#### Go Client

```go
package main

import (
    "context"
    "github.com/servin/servin-go"
)

func main() {
    client, err := servin.NewClientWithOpts(servin.FromEnv)
    if err != nil {
        panic(err)
    }

    // List containers
    containers, err := client.ContainerList(context.Background(), types.ContainerListOptions{})
    if err != nil {
        panic(err)
    }

    // Create container
    resp, err := client.ContainerCreate(context.Background(), &container.Config{
        Image: "nginx:latest",
        ExposedPorts: nat.PortSet{
            "80/tcp": struct{}{},
        },
    }, &container.HostConfig{
        PortBindings: nat.PortMap{
            "80/tcp": []nat.PortBinding{{HostPort: "8080"}},
        },
    }, nil, nil, "web-server")
}
```

This comprehensive API reference covers all major endpoints and operations available in the Servin Container Runtime REST API.
