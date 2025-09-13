---
layout: default
title: Registry Operations
permalink: /registry-operations/
---

# Registry Operations

Complete guide to container registry operations with Servin Container Runtime.

## Registry Fundamentals

### Understanding Registries

Container registries store and distribute container images:

- **Public Registries**: Docker Hub, GitHub Container Registry, Quay.io
- **Private Registries**: Self-hosted or cloud-managed private registries
- **Hybrid Registries**: Mix of public and private registry access
- **Mirror Registries**: Local caches of remote registries

### Registry Components

```
Registry Architecture:
┌─────────────────────────────────────────────────────────┐
│                    Registry Frontend                    │
├─────────────────────────────────────────────────────────┤
│                  Authentication Layer                   │
├─────────────────────────────────────────────────────────┤
│                   Authorization Layer                   │
├─────────────────────────────────────────────────────────┤
│                      Storage Backend                    │
├─────────────────────────────────────────────────────────┤
│              Metadata & Index Storage                   │
└─────────────────────────────────────────────────────────┘
```

## Registry Configuration

### Basic Registry Setup

Configure registry access and authentication:

```bash
# Configure default registry
servin config set registry.default docker.io

# Add private registry
servin config set registry.private registry.company.com

# Configure registry mirrors
servin config set registry.mirrors \
  "https://mirror1.company.com,https://mirror2.company.com"

# Set insecure registries
servin config set registry.insecure \
  "internal-registry.company.com:5000"
```

### Registry Authentication

Manage registry credentials:

```bash
# Login to public registry
servin login docker.io
servin login --username myuser docker.io

# Login to private registry
servin login registry.company.com
servin login --username serviceaccount registry.company.com

# Login with token
echo $REGISTRY_TOKEN | servin login --username oauth2accesstoken \
  --password-stdin gcr.io

# Login with credential helper
servin login --credential-helper docker-credential-gcloud gcr.io

# Store credentials securely
servin login --password-store pass registry.company.com
```

### Registry Configuration File

Advanced registry configuration:

```json
{
  "registries": {
    "docker.io": {
      "mirrors": [
        "https://registry-mirror.company.com",
        "https://backup-mirror.company.com"
      ],
      "resolve": true,
      "tls": {
        "insecure": false,
        "ca_file": "/etc/ssl/certs/ca-certificates.crt"
      }
    },
    "registry.company.com": {
      "auth": {
        "username": "serviceaccount",
        "password": "encrypted-password",
        "auth_file": "/etc/servin/registry-auth.json"
      },
      "tls": {
        "insecure": false,
        "cert_file": "/etc/ssl/registry/client.crt",
        "key_file": "/etc/ssl/registry/client.key",
        "ca_file": "/etc/ssl/registry/ca.crt"
      },
      "headers": {
        "User-Agent": "Servin/1.0",
        "Custom-Header": "value"
      }
    },
    "localhost:5000": {
      "tls": {
        "insecure": true
      },
      "resolve": false
    }
  },
  "unqualified_search_registries": [
    "docker.io",
    "registry.company.com"
  ],
  "short_name_mode": "enforcing"
}
```

## Private Registry Deployment

### Setting Up a Private Registry

Deploy a private Docker registry:

```bash
# Basic registry deployment
servin run -d \
  --name registry \
  --restart always \
  -p 5000:5000 \
  -v registry-data:/var/lib/registry \
  registry:2

# Registry with authentication
servin run -d \
  --name secure-registry \
  --restart always \
  -p 5000:5000 \
  -v registry-data:/var/lib/registry \
  -v registry-auth:/auth \
  -v registry-certs:/certs \
  -e REGISTRY_AUTH=htpasswd \
  -e REGISTRY_AUTH_HTPASSWD_REALM="Registry Realm" \
  -e REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd \
  -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/domain.crt \
  -e REGISTRY_HTTP_TLS_KEY=/certs/domain.key \
  registry:2
```

### Registry Configuration

Advanced registry setup with configuration file:

```yaml
# registry-config.yml
version: 0.1
log:
  level: info
  formatter: text
  fields:
    service: registry

storage:
  filesystem:
    rootdirectory: /var/lib/registry
  delete:
    enabled: true
  cache:
    blobdescriptor: inmemory

http:
  addr: :5000
  secret: registry-secret-key
  relativeurls: false
  draintimeout: 60s

auth:
  htpasswd:
    realm: basic-realm
    path: /auth/htpasswd

middleware:
  registry:
    - name: cloudfront
      options:
        baseurl: https://my.cloudfrontdomain.com/
        privatekey: /path/to/pem
        keypairid: cloudfrontkeypairid
        duration: 3000s

health:
  storagedriver:
    enabled: true
    interval: 10s
    threshold: 3

proxy:
  remoteurl: https://registry-1.docker.io
  username: username
  password: password
```

### Authentication Setup

Configure registry authentication:

```bash
# Create htpasswd file
servin run --rm \
  --entrypoint htpasswd \
  registry:2 -Bbn admin password > auth/htpasswd

# Add more users
servin run --rm \
  --entrypoint htpasswd \
  registry:2 -Bbn developer devpass >> auth/htpasswd

# Create SSL certificates
openssl req -newkey rsa:4096 -nodes -sha256 \
  -keyout certs/domain.key \
  -x509 -days 365 \
  -out certs/domain.crt \
  -subj "/CN=registry.company.com"
```

## Registry Operations

### Image Distribution

Push and pull images from registries:

```bash
# Tag image for private registry
servin tag myapp:latest registry.company.com/myapp:latest

# Push to private registry
servin push registry.company.com/myapp:latest

# Pull from private registry
servin pull registry.company.com/myapp:latest

# Push multiple tags
servin push --all-tags registry.company.com/myapp

# Push with signature
servin push --sign registry.company.com/myapp:latest
```

### Registry Search

Search for images in registries:

```bash
# Search Docker Hub
servin search nginx

# Search with filters
servin search --filter stars=3 nginx
servin search --filter is-official=true nginx

# Search private registry
servin search registry.company.com/nginx

# Limit search results
servin search --limit 10 nginx

# Search with no truncation
servin search --no-trunc nginx
```

### Registry Catalog

List repositories in a registry:

```bash
# List repositories
curl -X GET https://registry.company.com/v2/_catalog

# List tags for repository
curl -X GET https://registry.company.com/v2/myapp/tags/list

# Get image manifest
curl -X GET https://registry.company.com/v2/myapp/manifests/latest

# Get image config
curl -X GET https://registry.company.com/v2/myapp/blobs/sha256:...
```

## Registry Security

### Access Control

Implement registry access control:

```bash
# Role-based access control
servin run -d \
  --name registry-authz \
  -p 5001:5001 \
  -v registry-authz-config:/config \
  cesanta/docker_auth:1 \
  /config/auth_config.yml

# OAuth2 authentication
servin run -d \
  --name registry-oauth \
  -e OAUTH2_PROVIDER=github \
  -e OAUTH2_CLIENT_ID=your-client-id \
  -e OAUTH2_CLIENT_SECRET=your-client-secret \
  -p 5002:5002 \
  oauth2-proxy/oauth2-proxy
```

### Content Trust

Enable image signing and verification:

```bash
# Enable content trust
export SERVIN_CONTENT_TRUST=1

# Sign and push image
servin push registry.company.com/myapp:latest

# Rotate signing keys
servin trust key rotate registry.company.com/myapp

# View trust data
servin trust inspect registry.company.com/myapp:latest

# Add trusted signers
servin trust signer add --key cert.pem alice registry.company.com/myapp

# Remove trusted signers
servin trust signer remove alice registry.company.com/myapp
```

### Vulnerability Scanning

Implement registry-level scanning:

```bash
# Configure registry with scanning
servin run -d \
  --name registry-scanner \
  -e SCANNER_ENABLED=true \
  -e SCANNER_TRIVY_URL=http://trivy:8080 \
  -v scanner-data:/data \
  goharbor/harbor-core

# Scan images on push
servin push --scan registry.company.com/myapp:latest

# View scan results
servin scan results registry.company.com/myapp:latest

# Set security policies
servin policy create --severity HIGH --action BLOCK vulnerability-policy
```

## Registry Mirroring

### Pull-through Cache

Set up registry mirrors for improved performance:

```bash
# Deploy pull-through cache
servin run -d \
  --name registry-mirror \
  -p 5000:5000 \
  -v mirror-data:/var/lib/registry \
  -e REGISTRY_PROXY_REMOTEURL=https://registry-1.docker.io \
  -e REGISTRY_PROXY_USERNAME=dockerhub-user \
  -e REGISTRY_PROXY_PASSWORD=dockerhub-pass \
  registry:2

# Configure clients to use mirror
echo '{"registry-mirrors": ["http://mirror.company.com:5000"]}' > \
  /etc/servin/daemon.json
```

### Registry Replication

Replicate registries across regions:

```bash
# Primary registry
servin run -d \
  --name registry-primary \
  -p 5000:5000 \
  -v primary-data:/var/lib/registry \
  registry:2

# Replica registry
servin run -d \
  --name registry-replica \
  -p 5001:5000 \
  -v replica-data:/var/lib/registry \
  -e REGISTRY_PROXY_REMOTEURL=http://registry-primary:5000 \
  registry:2

# Sync script for replication
#!/bin/bash
# registry-sync.sh
servin run --rm \
  -v primary-data:/source:ro \
  -v replica-data:/dest \
  alpine:latest \
  sh -c "cp -r /source/* /dest/"
```

## Registry Maintenance

### Registry Cleanup

Manage registry storage and cleanup:

```bash
# Run garbage collection
servin exec registry \
  /bin/registry garbage-collect /etc/registry/config.yml

# Delete unused layers
servin exec registry \
  /bin/registry garbage-collect --delete-untagged /etc/registry/config.yml

# Dry run garbage collection
servin exec registry \
  /bin/registry garbage-collect --dry-run /etc/registry/config.yml

# Clean registry data
servin volume create temp-registry
servin run --rm \
  -v registry-data:/registry \
  -v temp-registry:/temp \
  alpine:latest \
  sh -c "find /registry -type f -mtime +30 -delete"
```

### Registry Backup

Backup registry data and metadata:

```bash
# Backup registry data
servin run --rm \
  -v registry-data:/data:ro \
  -v $(pwd)/backups:/backup \
  alpine:latest \
  tar czf /backup/registry-$(date +%Y%m%d).tar.gz -C /data .

# Backup registry configuration
cp /etc/registry/config.yml backups/config-$(date +%Y%m%d).yml

# Backup authentication data
servin run --rm \
  -v registry-auth:/auth:ro \
  -v $(pwd)/backups:/backup \
  alpine:latest \
  tar czf /backup/auth-$(date +%Y%m%d).tar.gz -C /auth .

# Automated backup script
#!/bin/bash
# registry-backup.sh
DATE=$(date +%Y%m%d-%H%M%S)
BACKUP_DIR="/backups/registry"

mkdir -p $BACKUP_DIR

# Backup data
servin run --rm \
  -v registry-data:/data:ro \
  -v $BACKUP_DIR:/backup \
  alpine:latest \
  tar czf /backup/data-$DATE.tar.gz -C /data .

# Backup config
cp /etc/registry/config.yml $BACKUP_DIR/config-$DATE.yml

# Clean old backups
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete
```

## Registry Monitoring

### Health Checks

Monitor registry health and performance:

```bash
# Registry health endpoint
curl -f http://registry.company.com:5000/debug/health

# Check registry version
curl http://registry.company.com:5000/v2/

# Monitor registry metrics
curl http://registry.company.com:5000/debug/metrics

# Custom health check
servin run --rm \
  appropriate/curl \
  curl -f http://registry:5000/debug/health || echo "Registry unhealthy"
```

### Logging and Metrics

Configure comprehensive monitoring:

```yaml
# registry-config.yml with logging
log:
  level: info
  formatter: json
  fields:
    service: registry
    environment: production

reporting:
  bugsnag:
    apikey: bugsnagapikey
    releasestage: production

health:
  storagedriver:
    enabled: true
    interval: 10s
    threshold: 3
  http:
    - file: /path/to/check
      interval: 5s
      threshold: 3
      timeout: 3s
      statuscode: 200
      body: |
        server is alive
```

### Performance Monitoring

Monitor registry performance:

```bash
# Registry statistics
servin stats registry

# Storage usage
servin exec registry du -sh /var/lib/registry

# Network metrics
servin exec registry netstat -an | grep :5000

# Process monitoring
servin exec registry ps aux

# Resource usage
servin exec registry top
```

## Registry High Availability

### Load Balancing

Set up load-balanced registry:

```bash
# HAProxy configuration for registry
# haproxy.cfg
global
    daemon

defaults
    mode http
    timeout connect 5000ms
    timeout client 50000ms
    timeout server 50000ms

frontend registry_frontend
    bind *:5000
    default_backend registry_backend

backend registry_backend
    balance roundrobin
    option httpchk GET /v2/
    server registry1 registry1:5000 check
    server registry2 registry2:5000 check
    server registry3 registry3:5000 check

# Deploy HAProxy
servin run -d \
  --name registry-lb \
  -p 5000:5000 \
  -v $(pwd)/haproxy.cfg:/usr/local/etc/haproxy/haproxy.cfg \
  haproxy:alpine
```

### Clustering

Deploy registry in cluster mode:

```bash
# Shared storage for cluster
servin volume create registry-shared

# Registry cluster node 1
servin run -d \
  --name registry-node1 \
  -p 5001:5000 \
  -v registry-shared:/var/lib/registry \
  -e REGISTRY_STORAGE_DELETE_ENABLED=true \
  registry:2

# Registry cluster node 2
servin run -d \
  --name registry-node2 \
  -p 5002:5000 \
  -v registry-shared:/var/lib/registry \
  -e REGISTRY_STORAGE_DELETE_ENABLED=true \
  registry:2

# Registry cluster node 3
servin run -d \
  --name registry-node3 \
  -p 5003:5000 \
  -v registry-shared:/var/lib/registry \
  -e REGISTRY_STORAGE_DELETE_ENABLED=true \
  registry:2
```

## Cloud Registry Integration

### AWS ECR

Integrate with Amazon Elastic Container Registry:

```bash
# Authenticate with ECR
aws ecr get-login-password --region us-west-2 | \
  servin login --username AWS --password-stdin \
  123456789012.dkr.ecr.us-west-2.amazonaws.com

# Create ECR repository
aws ecr create-repository --repository-name myapp

# Tag and push to ECR
servin tag myapp:latest \
  123456789012.dkr.ecr.us-west-2.amazonaws.com/myapp:latest
servin push 123456789012.dkr.ecr.us-west-2.amazonaws.com/myapp:latest

# Pull from ECR
servin pull 123456789012.dkr.ecr.us-west-2.amazonaws.com/myapp:latest
```

### Google Container Registry

Integrate with Google Container Registry:

```bash
# Authenticate with GCR
gcloud auth configure-docker

# Tag and push to GCR
servin tag myapp:latest gcr.io/project-id/myapp:latest
servin push gcr.io/project-id/myapp:latest

# Pull from GCR
servin pull gcr.io/project-id/myapp:latest

# Use service account
servin login --username _json_key --password-stdin gcr.io < key.json
```

### Azure Container Registry

Integrate with Azure Container Registry:

```bash
# Authenticate with ACR
az acr login --name myregistry

# Tag and push to ACR
servin tag myapp:latest myregistry.azurecr.io/myapp:latest
servin push myregistry.azurecr.io/myapp:latest

# Pull from ACR
servin pull myregistry.azurecr.io/myapp:latest

# Use service principal
servin login --username $SP_APP_ID --password $SP_PASSWD \
  myregistry.azurecr.io
```

This comprehensive registry operations guide covers all aspects of container registry management with Servin, from basic operations to advanced enterprise deployment and cloud integration.
