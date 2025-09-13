---
layout: default
title: Project Overview
permalink: /overview/
---

# 🚀 Project Overview

Servin Container Runtime is a comprehensive container management solution that provides Docker-compatible functionality with enhanced Kubernetes integration and professional user interfaces.

## What is Servin?

Servin is a modern container runtime designed to bridge the gap between Docker's ease of use and Kubernetes' enterprise requirements. It offers:

- **Full Docker API compatibility** for seamless migration
- **Native Kubernetes CRI support** for cluster integration
- **Multiple user interfaces** (CLI, TUI, GUI) for different workflows
- **Cross-platform support** with professional installers
- **Enterprise-grade features** for production environments

## Why Choose Servin?

### 🔄 **Seamless Migration**
Migrate from Docker without changing your existing workflows, scripts, or automation.

### 🏢 **Enterprise Ready**
Built with enterprise requirements in mind, including security, scalability, and monitoring.

### 🖥️ **User-Friendly Interfaces**
Choose from command-line, terminal UI, or desktop GUI based on your preference.

### ☸️ **Kubernetes Native**
First-class Kubernetes support with full CRI v1alpha2 implementation.

### 🔧 **Professional Installation**
Native installers for Windows, Linux, and macOS with service integration.

## Core Components

| Component | Description | Purpose |
|-----------|-------------|---------|
| **Runtime Engine** | Container lifecycle management | Create, start, stop, delete containers |
| **Image Manager** | OCI image handling and storage | Pull, push, build, tag images |
| **Volume Manager** | Persistent storage management | Manage container volumes and bind mounts |
| **Network Manager** | Container networking | Bridge networks, port forwarding |
| **Registry Client** | Image registry operations | Authenticate and communicate with registries |
| **CRI Server** | Kubernetes integration | gRPC server for kubelet communication |

## Use Cases

### Development Workflows
- Local application development and testing
- Container-based development environments
- CI/CD pipeline integration
- Multi-service application orchestration

### Production Deployment
- Kubernetes cluster container runtime
- Microservices architecture deployment
- Container-based application hosting
- Enterprise container management

### DevOps Operations
- Container infrastructure management
- Automated deployment pipelines
- Monitoring and logging integration
- Security and compliance management

## Getting Started

Ready to get started? Check out our [Installation Guide]({{ '/installation' | relative_url }}) to begin using Servin in your environment.

[Install Servin →]({{ '/installation' | relative_url }}){: .btn .btn-primary}
[View Architecture →]({{ '/architecture' | relative_url }}){: .btn .btn-outline}
