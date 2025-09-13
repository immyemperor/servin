---
layout: default
title: Home
permalink: /
---

# 🚀 Servin Container Runtime

**A Docker-compatible container runtime with Kubernetes CRI support and professional desktop interface.**

<div class="feature-grid">
  <div class="feature-box">
    <div class="feature-icon">🐳</div>
    <h4>Docker Compatible</h4>
    <p>Seamless migration from Docker workflows with full API compatibility.</p>
  </div>
  <div class="feature-box">
    <div class="feature-icon">☸️</div>
    <h4>Kubernetes CRI</h4>
    <p>Full Container Runtime Interface v1alpha2 implementation for Kubernetes clusters.</p>
  </div>
  <div class="feature-box">
    <div class="feature-icon">🖥️</div>
    <h4>Multiple Interfaces</h4>
    <p>CLI, Terminal UI, and Desktop GUI applications for every workflow.</p>
  </div>
</div>

## Key Capabilities

### 🎯 Core Runtime Features
- <span class="badge badge-success">✓</span> Container lifecycle management
- <span class="badge badge-success">✓</span> Image management (pull, push, build)
- <span class="badge badge-success">✓</span> Volume management
- <span class="badge badge-success">✓</span> Network management
- <span class="badge badge-success">✓</span> Registry operations

### 🔌 Integration Features
- <span class="badge badge-success">✓</span> Kubernetes CRI v1alpha2
- <span class="badge badge-success">✓</span> Cross-platform support
- <span class="badge badge-success">✓</span> Service integration
- <span class="badge badge-success">✓</span> Professional installers
- <span class="badge badge-success">✓</span> REST and gRPC APIs

## Target Users

<div class="feature-grid">
  <div class="feature-box">
    <div class="feature-icon">👨‍💻</div>
    <h4>Developers</h4>
    <p>Container-based application development</p>
  </div>
  <div class="feature-box">
    <div class="feature-icon">⚙️</div>
    <h4>DevOps Engineers</h4>
    <p>Container orchestration and deployment</p>
  </div>
  <div class="feature-box">
    <div class="feature-icon">🔧</div>
    <h4>System Admins</h4>
    <p>Container infrastructure management</p>
  </div>
  <div class="feature-box">
    <div class="feature-icon">☸️</div>
    <h4>Kubernetes Users</h4>
    <p>CRI-compatible runtime for clusters</p>
  </div>
</div>

## Quick Start

1. **Install Servin** following our [Installation Guide]({{ '/installation' | relative_url }})
2. **Start the daemon:**
   ```bash
   servin daemon
   ```
3. **Run your first container:**
   ```bash
   servin run hello-world
   ```

[Get Started →]({{ '/installation' | relative_url }}){: .btn .btn-primary}
[View on GitHub →]({{ site.github.repository_url }}){: .btn .btn-outline}
