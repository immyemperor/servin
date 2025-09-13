---
layout: default
title: Configuration
---

# Configuration

Comprehensive guide to configuring Servin Container Runtime for optimal performance and security.

## Configuration Architecture

### Configuration Hierarchy

Servin uses a layered configuration approach:

```
Configuration Priority (highest to lowest):
┌─────────────────────────────────────────────────────────┐
│                 Command Line Flags                      │
├─────────────────────────────────────────────────────────┤
│                Environment Variables                    │
├─────────────────────────────────────────────────────────┤
│               User Configuration Files                  │
├─────────────────────────────────────────────────────────┤
│              System Configuration Files                 │
├─────────────────────────────────────────────────────────┤
│                   Default Values                       │
└─────────────────────────────────────────────────────────┘
```

### Configuration Locations

- **Global Config**: `/etc/servin/daemon.json`
- **User Config**: `~/.servin/config.json`
- **Runtime Config**: `/var/lib/servin/config.json`
- **Environment**: Environment variables with `SERVIN_` prefix

## Daemon Configuration

### Basic Daemon Configuration

Configure the Servin daemon:

```json
{
  "data-root": "/var/lib/servin",
  "exec-root": "/var/run/servin",
  "storage-driver": "overlay2",
  "storage-opts": [
    "overlay2.override_kernel_check=true"
  ],
  "runtime": "runc",
  "default-runtime": "runc",
  "runtimes": {
    "runc": {
      "path": "/usr/bin/runc",
      "runtime-args": ["--systemd-cgroup"]
    },
    "crun": {
      "path": "/usr/bin/crun"
    }
  },
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "5"
  },
  "debug": false,
  "experimental": false
}
```

### Advanced Daemon Settings

Configure advanced daemon options:

```json
{
  "api-cors-header": "*",
  "authorization-plugins": ["authz-broker"],
  "bip": "192.168.1.1/24",
  "bridge": "servin0",
  "cgroup-parent": "/servin",
  "cluster-store": "consul://localhost:8500",
  "cluster-store-opts": {
    "discovery.heartbeat": 20,
    "discovery.ttl": 60
  },
  "cluster-advertise": "192.168.1.100:2376",
  "default-gateway": "192.168.1.1",
  "default-gateway-v6": "2001:db8::1",
  "dns": ["8.8.8.8", "8.8.4.4"],
  "dns-opts": ["ndots:2", "edns0"],
  "dns-search": ["company.com"],
  "exec-opts": ["native.cgroupdriver=systemd"],
  "fixed-cidr": "172.17.0.0/16",
  "fixed-cidr-v6": "2001:db8::/64",
  "group": "servin",
  "hosts": [
    "unix:///var/run/servin.sock",
    "tcp://0.0.0.0:2376"
  ],
  "icc": true,
  "init": true,
  "init-path": "/usr/libexec/servin-init",
  "insecure-registries": ["myregistry:5000"],
  "ip": "0.0.0.0",
  "ip-forward": true,
  "ip-masq": true,
  "iptables": true,
  "ipv6": false,
  "labels": ["environment=production", "datacenter=us-west"],
  "live-restore": true,
  "log-level": "info",
  "max-concurrent-downloads": 3,
  "max-concurrent-uploads": 5,
  "max-download-attempts": 5,
  "mtu": 1500,
  "no-new-privileges": false,
  "oom-score-adjust": -500,
  "pidfile": "/var/run/servin.pid",
  "raw-logs": false,
  "registry-mirrors": ["https://registry.company.com"],
  "selinux-enabled": false,
  "shutdown-timeout": 15,
  "storage-driver": "overlay2",
  "swarm-default-advertise-addr": "192.168.1.100",
  "tls": true,
  "tlscert": "/etc/servin/server.pem",
  "tlskey": "/etc/servin/server-key.pem",
  "tlsverify": true,
  "tlscacert": "/etc/servin/ca.pem",
  "userland-proxy": true,
  "userns-remap": "default"
}
```

## Storage Configuration

### Storage Drivers

Configure storage backend:

```json
{
  "storage-driver": "overlay2",
  "storage-opts": [
    "overlay2.override_kernel_check=true",
    "overlay2.size=120G"
  ]
}
```

#### Overlay2 Storage Driver

```json
{
  "storage-driver": "overlay2",
  "storage-opts": [
    "overlay2.override_kernel_check=true",
    "overlay2.size=50G",
    "overlay2.use_composefs=true"
  ]
}
```

#### DeviceMapper Storage Driver

```json
{
  "storage-driver": "devicemapper",
  "storage-opts": [
    "dm.thinpooldev=/dev/mapper/servin-thinpool",
    "dm.use_deferred_removal=true",
    "dm.use_deferred_deletion=true"
  ]
}
```

#### Btrfs Storage Driver

```json
{
  "storage-driver": "btrfs",
  "storage-opts": [
    "btrfs.min_space=20G"
  ]
}
```

#### ZFS Storage Driver

```json
{
  "storage-driver": "zfs",
  "storage-opts": [
    "zfs.fsname=zroot/servin"
  ]
}
```

## Network Configuration

### Default Network Settings

Configure default networking:

```json
{
  "bridge": "servin0",
  "bip": "172.17.0.1/16",
  "fixed-cidr": "172.17.0.0/16",
  "fixed-cidr-v6": "2001:db8::/64",
  "default-gateway": "172.17.0.1",
  "default-gateway-v6": "2001:db8::1",
  "mtu": 1500,
  "icc": true,
  "ip-forward": true,
  "ip-masq": true,
  "iptables": true,
  "ipv6": false,
  "userland-proxy": true
}
```

### Custom Network Configuration

Advanced network settings:

```json
{
  "default-address-pools": [
    {
      "base": "172.30.0.0/16",
      "size": 24
    },
    {
      "base": "172.31.0.0/16", 
      "size": 24
    }
  ],
  "dns": ["8.8.8.8", "1.1.1.1"],
  "dns-opts": ["ndots:2", "edns0"],
  "dns-search": ["company.local"],
  "default-network-opts": {
    "bridge": {
      "com.docker.network.bridge.default_bridge": "true",
      "com.docker.network.bridge.enable_icc": "true",
      "com.docker.network.bridge.enable_ip_masquerade": "true",
      "com.docker.network.bridge.host_binding_ipv4": "0.0.0.0",
      "com.docker.network.bridge.name": "servin0",
      "com.docker.network.driver.mtu": "1500"
    }
  }
}
```

## Security Configuration

### TLS Configuration

Enable TLS for daemon communication:

```json
{
  "hosts": ["tcp://0.0.0.0:2376"],
  "tls": true,
  "tlscert": "/etc/servin/server.pem",
  "tlskey": "/etc/servin/server-key.pem",
  "tlsverify": true,
  "tlscacert": "/etc/servin/ca.pem"
}
```

#### Generate TLS Certificates

```bash
# Create CA key
openssl genrsa -aes256 -out ca-key.pem 4096

# Create CA certificate
openssl req -new -x509 -days 365 -key ca-key.pem -sha256 -out ca.pem

# Create server key
openssl genrsa -out server-key.pem 4096

# Create server certificate signing request
openssl req -subj "/CN=servind" -sha256 -new -key server-key.pem -out server.csr

# Sign server certificate
echo subjectAltName = DNS:servind,IP:10.10.10.20,IP:127.0.0.1 >> extfile.cnf
echo extendedKeyUsage = serverAuth >> extfile.cnf
openssl x509 -req -days 365 -sha256 -in server.csr -CA ca.pem -CAkey ca-key.pem -out server-cert.pem -extfile extfile.cnf -CAcreateserial

# Create client key
openssl genrsa -out key.pem 4096

# Create client certificate signing request
openssl req -subj '/CN=client' -new -key key.pem -out client.csr

# Sign client certificate
echo extendedKeyUsage = clientAuth >> extfile-client.cnf
openssl x509 -req -days 365 -sha256 -in client.csr -CA ca.pem -CAkey ca-key.pem -out cert.pem -extfile extfile-client.cnf -CAcreateserial

# Set permissions
chmod -v 0400 ca-key.pem key.pem server-key.pem
chmod -v 0444 ca.pem server-cert.pem cert.pem
```

### User Namespace Configuration

Enable user namespace remapping:

```json
{
  "userns-remap": "default"
}
```

#### Setup User Namespace

```bash
# Create subuid and subgid mappings
echo 'servind:165536:65536' >> /etc/subuid
echo 'servind:165536:65536' >> /etc/subgid

# Restart daemon
systemctl restart servind
```

### Authorization Plugins

Configure authorization plugins:

```json
{
  "authorization-plugins": ["authz-broker"],
  "authz-broker-config": {
    "broker-url": "https://authz.company.com",
    "broker-tls-verify": true,
    "broker-cert": "/etc/servin/authz-cert.pem",
    "broker-key": "/etc/servin/authz-key.pem",
    "broker-ca": "/etc/servin/authz-ca.pem"
  }
}
```

## Logging Configuration

### Log Drivers

Configure different log drivers:

```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3",
    "compress": "true"
  }
}
```

#### Syslog Driver

```json
{
  "log-driver": "syslog",
  "log-opts": {
    "syslog-address": "udp://192.168.1.42:514",
    "syslog-facility": "daemon",
    "syslog-tag": "servin"
  }
}
```

#### Fluentd Driver

```json
{
  "log-driver": "fluentd",
  "log-opts": {
    "fluentd-address": "fluentd.company.com:24224",
    "fluentd-tag": "servin.{{.Name}}",
    "fluentd-async-connect": "true",
    "fluentd-buffer-limit": "1MB"
  }
}
```

#### Splunk Driver

```json
{
  "log-driver": "splunk",
  "log-opts": {
    "splunk-token": "your-splunk-token",
    "splunk-url": "https://splunk.company.com:8088",
    "splunk-source": "servin",
    "splunk-sourcetype": "servin_logs",
    "splunk-index": "main",
    "splunk-capath": "/path/to/ca.pem",
    "splunk-caname": "SplunkServerDefaultCert",
    "splunk-insecureskipverify": "false",
    "splunk-format": "json",
    "splunk-verify-connection": "true",
    "splunk-gzip": "false",
    "splunk-gzip-level": "1"
  }
}
```

## Resource Management

### Cgroup Configuration

Configure cgroup settings:

```json
{
  "cgroup-parent": "/servin",
  "cgroups-namespace": "host",
  "exec-opts": ["native.cgroupdriver=systemd"],
  "default-ulimits": {
    "memlock": {
      "Hard": -1,
      "Name": "memlock",
      "Soft": -1
    },
    "nofile": {
      "Hard": 64000,
      "Name": "nofile", 
      "Soft": 64000
    }
  }
}
```

### Resource Limits

Set default resource limits:

```json
{
  "default-ulimits": {
    "nproc": {
      "Hard": 65535,
      "Name": "nproc",
      "Soft": 65535
    },
    "memlock": {
      "Hard": -1,
      "Name": "memlock",
      "Soft": -1
    },
    "stack": {
      "Hard": 67108864,
      "Name": "stack",
      "Soft": 8388608
    }
  },
  "oom-score-adjust": -500
}
```

## Registry Configuration

### Registry Settings

Configure registry access:

```json
{
  "registry-mirrors": [
    "https://registry.company.com",
    "https://mirror.company.com"
  ],
  "insecure-registries": [
    "internal-registry:5000",
    "192.168.1.100:5000"
  ],
  "disable-legacy-registry": true
}
```

### Registry Authentication

Configure registry credentials:

```json
{
  "auths": {
    "https://index.docker.io/v1/": {
      "auth": "base64-encoded-auth-string"
    },
    "registry.company.com": {
      "auth": "base64-encoded-auth-string",
      "email": "user@company.com"
    }
  },
  "credHelpers": {
    "gcr.io": "gcloud",
    "us.gcr.io": "gcloud",
    "asia.gcr.io": "gcloud",
    "staging-k8s.gcr.io": "gcloud"
  },
  "credsStore": "osxkeychain"
}
```

## Runtime Configuration

### Container Runtime Configuration

Configure OCI runtime:

```json
{
  "default-runtime": "runc",
  "runtimes": {
    "runc": {
      "path": "/usr/bin/runc",
      "runtime-args": [
        "--systemd-cgroup"
      ]
    },
    "crun": {
      "path": "/usr/bin/crun",
      "runtime-args": []
    },
    "kata": {
      "path": "/usr/bin/kata-runtime"
    },
    "gvisor": {
      "path": "/usr/local/bin/runsc",
      "runtime-args": [
        "--platform=ptrace"
      ]
    }
  }
}
```

### Runtime Security

Configure runtime security options:

```json
{
  "no-new-privileges": true,
  "seccomp-profile": "/etc/servin/seccomp.json",
  "selinux-enabled": true,
  "default-runtime": "runc",
  "runtimes": {
    "runc": {
      "path": "/usr/bin/runc",
      "runtime-args": [
        "--systemd-cgroup",
        "--seccomp",
        "/etc/servin/seccomp.json"
      ]
    }
  }
}
```

## Performance Tuning

### Performance Configuration

Optimize daemon performance:

```json
{
  "max-concurrent-downloads": 6,
  "max-concurrent-uploads": 5,
  "max-download-attempts": 5,
  "storage-opts": [
    "overlay2.override_kernel_check=true"
  ],
  "exec-opts": [
    "native.cgroupdriver=systemd"
  ],
  "live-restore": true,
  "init": true,
  "userland-proxy": false
}
```

### Memory Management

Configure memory settings:

```json
{
  "default-ulimits": {
    "memlock": {
      "Hard": -1,
      "Name": "memlock",
      "Soft": -1
    }
  },
  "oom-score-adjust": -500,
  "default-shm-size": "64M"
}
```

## Environment Variables

### Daemon Environment Variables

Configure daemon via environment variables:

```bash
# Data and execution directories
export SERVIN_DATA_ROOT=/var/lib/servin
export SERVIN_EXEC_ROOT=/var/run/servin

# Storage configuration
export SERVIN_STORAGE_DRIVER=overlay2
export SERVIN_STORAGE_OPTS="overlay2.override_kernel_check=true"

# Network configuration
export SERVIN_BRIDGE=servin0
export SERVIN_BIP=172.17.0.1/16

# Security configuration
export SERVIN_TLS=true
export SERVIN_TLSCERT=/etc/servin/server.pem
export SERVIN_TLSKEY=/etc/servin/server-key.pem
export SERVIN_TLSVERIFY=true
export SERVIN_TLSCACERT=/etc/servin/ca.pem

# Logging configuration
export SERVIN_LOG_DRIVER=json-file
export SERVIN_LOG_LEVEL=info

# Registry configuration
export SERVIN_REGISTRY_MIRROR=https://registry.company.com
export SERVIN_INSECURE_REGISTRY=internal-registry:5000

# Resource configuration
export SERVIN_DEFAULT_RUNTIME=runc
export SERVIN_LIVE_RESTORE=true
```

### Client Environment Variables

Configure client behavior:

```bash
# Connection settings
export SERVIN_HOST=tcp://servin.company.com:2376
export SERVIN_TLS_VERIFY=1
export SERVIN_CERT_PATH=/etc/servin/client

# API version
export SERVIN_API_VERSION=1.41

# Default configuration
export SERVIN_CONFIG=/home/user/.servin
export SERVIN_CONTEXT=production
```

## Configuration Validation

### Validate Configuration

Check configuration validity:

```bash
# Validate daemon configuration
servin daemon --validate-config

# Test configuration without starting daemon
servin daemon --config-file /etc/servin/daemon.json --validate

# Check configuration syntax
jq . /etc/servin/daemon.json

# Verify TLS certificates
openssl x509 -in /etc/servin/server.pem -text -noout
openssl verify -CAfile /etc/servin/ca.pem /etc/servin/server.pem
```

### Configuration Testing

Test configuration changes:

```bash
# Backup current configuration
cp /etc/servin/daemon.json /etc/servin/daemon.json.backup

# Test new configuration
servin daemon --config-file /etc/servin/daemon-test.json --validate

# Apply configuration
systemctl reload servind

# Verify daemon status
systemctl status servind
servin system info
```

## Configuration Management

### Configuration Templates

Use configuration templates:

```bash
#!/bin/bash
# generate-config.sh

ENVIRONMENT=${1:-development}
DATA_ROOT=${2:-/var/lib/servin}
LOG_LEVEL=${3:-info}

cat > /etc/servin/daemon.json << EOF
{
  "data-root": "${DATA_ROOT}",
  "storage-driver": "overlay2",
  "log-level": "${LOG_LEVEL}",
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "5"
  },
  "live-restore": true,
  "labels": ["environment=${ENVIRONMENT}"]
}
EOF

echo "Configuration generated for ${ENVIRONMENT} environment"
```

### Configuration Automation

Automate configuration management:

```yaml
# ansible-playbook.yml
---
- hosts: servin_nodes
  become: yes
  vars:
    servin_config:
      data-root: /var/lib/servin
      storage-driver: overlay2
      log-level: info
      log-driver: json-file
      log-opts:
        max-size: 10m
        max-file: "5"
      live-restore: true
      
  tasks:
    - name: Create servin config directory
      file:
        path: /etc/servin
        state: directory
        
    - name: Deploy servin configuration
      copy:
        content: "{{ servin_config | to_nice_json }}"
        dest: /etc/servin/daemon.json
        backup: yes
      notify: restart servin
      
  handlers:
    - name: restart servin
      systemd:
        name: servind
        state: restarted
```

## Best Practices

### Security Best Practices

1. **Enable TLS**: Always use TLS for daemon communication
2. **User Namespaces**: Enable user namespace remapping
3. **Resource Limits**: Set appropriate resource limits
4. **SELinux/AppArmor**: Enable mandatory access controls
5. **Regular Updates**: Keep daemon and runtime updated

### Performance Best Practices

1. **Storage Driver**: Choose appropriate storage driver
2. **Log Rotation**: Configure log rotation to prevent disk filling
3. **Resource Monitoring**: Monitor system resources
4. **Network Optimization**: Optimize network settings for workload
5. **Cgroup Management**: Use systemd cgroup driver

### Operational Best Practices

1. **Configuration Management**: Use version control for configs
2. **Backup**: Regular backup of configuration and data
3. **Monitoring**: Monitor daemon health and performance
4. **Documentation**: Document configuration changes
5. **Testing**: Test configuration changes in staging

This comprehensive configuration guide covers all aspects of configuring Servin Container Runtime for different environments and use cases.
