---
layout: default
title: Home
permalink: /
---

# 🚀 Servin Container Runtime

**Revolutionary dual-mode container runtime with universal cross-platform containerization through VM-based Linux containers.**

<div class="feature-grid">
  <div class="feature-box">
    <div class="feature-icon">🔓</div>
    <h4>Open Source</h4>
    <p>Free and open-source with full transparency, community-driven development, and no vendor lock-in.</p>
  </div>
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

### � Open Source Features
- <span class="badge badge-success">✓</span> Free and open-source license
- <span class="badge badge-success">✓</span> Full source code transparency
- <span class="badge badge-success">✓</span> Community-driven development
- <span class="badge badge-success">✓</span> No vendor lock-in
- <span class="badge badge-success">✓</span> Educational and learning-friendly

### �🔌 Integration Features
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

## 🎉 What's New

<div class="feature-box" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 20px; margin: 20px 0; border-radius: 10px;">
  <div style="display: flex; align-items: center; gap: 15px;">
    <div class="feature-icon" style="background: rgba(255,255,255,0.2); color: white;">🚀</div>
    <div>
      <h4 style="margin: 0; color: white;">Enhanced VM Engine Management</h4>
      <p style="margin: 5px 0 0 0; opacity: 0.9;">Real-time status indicators, auto-connect terminal, and universal development provider with state persistence.</p>
      <a href="{{ '/recent-enhancements' | relative_url }}" style="color: #ffd700; text-decoration: none; font-weight: bold;">Learn More →</a>
    </div>
  </div>
</div>

## Quick Start

### 📦 Download & Install

Choose your platform and download the appropriate package:

<div class="feature-grid">
  <div class="feature-box">
    <div class="feature-icon">🪟</div>
    <h4>Windows</h4>
    <p><code>servin_1.0.0_windows_amd64_installer.exe</code></p>
    <span class="badge badge-primary">NSIS Installer</span>
  </div>
  <div class="feature-box">
    <div class="feature-icon">🐧</div>
    <h4>Linux</h4>
    <p><code>servin_1.0.0_linux_amd64.tar.gz</code></p>
    <span class="badge badge-primary">TAR.GZ + Installer</span>
  </div>
  <div class="feature-box">
    <div class="feature-icon">🍎</div>
    <h4>macOS</h4>
    <p><code>servin_1.0.0_macos_universal.tar.gz</code></p>
    <span class="badge badge-primary">Universal Binary</span>
  </div>
</div>

### 🚀 Get Started

1. **Install Servin** following our [Installation Guide]({{ '/installation' | relative_url }})
2. **Start the daemon:**
   ```bash
   servin daemon
   ```
3. **Run your first container:**
   ```bash
   servin run hello-world
   ```

### 🎯 Interface Options

- **CLI**: `servin` - Full command-line interface
- **TUI**: `servin-desktop` - Terminal user interface  
- **GUI**: `servin-gui` - Desktop application with Fyne
- **WebView**: `servin-webview` - Modern web interface

[Get Started →]({{ '/installation' | relative_url }}){: .btn .btn-primary}

### 🤝 Join the Community

Servin is built with ❤️ by the open-source community. We welcome contributions, feedback, and collaboration from developers worldwide.

- **🐛 Report Issues**: Found a bug? [Open an issue]({{ site.github.repository_url }}/issues)
- **💡 Feature Requests**: Have an idea? [Start a discussion]({{ site.github.repository_url }}/discussions)
- **🛠️ Contribute**: Check our [Contributing Guide]({{ '/development' | relative_url }})
- **📖 Documentation**: Help improve our docs
- **⭐ Star Us**: Show your support on GitHub

[View on GitHub →]({{ site.github.repository_url }}){: .btn .btn-outline}
