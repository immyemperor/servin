# Servin Container Runtime - Hybrid Deployment Mode
# 
# This Dockerfile creates a containerized version of Servin daemon for hybrid deployments
# where Servin runs in Docker/Kubernetes but manages VM-based container workloads.
#
# WHEN TO USE THIS DOCKERFILE:
# ✅ Kubernetes/orchestration deployments
# ✅ Hybrid infrastructure (Docker for services, VMs for workloads)
# ✅ Development/testing environments
# ✅ Service mesh integration
#
# WHEN TO USE PURE VM MODE INSTEAD:
# ❌ Single-host development (use ./servin binary directly)
# ❌ Pure VM-based infrastructure (use native installation)
# ❌ Docker replacement scenarios (defeats the purpose)
#
# For pure VM mode: Download native installer from releases
# For hybrid mode: Use this Dockerfile

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY build/linux-amd64/servin /usr/local/bin/servin

RUN adduser -D -s /bin/sh servin
USER servin

# Servin daemon API port
EXPOSE 10250

CMD ["/usr/local/bin/servin", "daemon"]
