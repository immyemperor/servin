# Servin Release Documentation

## Release Process Overview

This document outlines the complete release process for Servin, including building, packaging, testing, and distribution across all supported platforms.

## Release Workflow

### 1. Version Management

#### Semantic Versioning
- **Major** (X.0.0): Breaking changes, major architecture updates
- **Minor** (0.X.0): New features, backward compatible
- **Patch** (0.0.X): Bug fixes, security patches

#### Version Bumping Process
```bash
# Update version in main files
./scripts/version-bump.sh <new-version>

# Files updated:
# - go.mod
# - pkg/version/version.go
# - installers/*/VERSION
# - package.json (GUI components)
```

### 2. Pre-Release Testing

#### Automated Testing
```bash
# Run full test suite
make test-all

# Platform-specific testing
make test-linux
make test-macos
make test-windows

# VM integration testing
make test-vm-integration

# GUI testing
make test-gui
```

#### Manual Testing Checklist
- [ ] Basic container operations (run, stop, remove)
- [ ] VM initialization and management
- [ ] Cross-platform GUI functionality
- [ ] Installation scripts
- [ ] Package manager integrations
- [ ] Documentation completeness

### 3. Build Process

#### Complete Distribution Build
```bash
# Build everything for all platforms
./build-vm-distribution.sh --all

# Platform-specific builds
./build-vm-distribution.sh --platform linux
./build-vm-distribution.sh --platform darwin
./build-vm-distribution.sh --platform windows
```

#### Build Artifacts
```
dist/
├── binaries/
│   ├── linux-amd64/
│   ├── linux-arm64/
│   ├── darwin-amd64/
│   ├── darwin-arm64/
│   └── windows-amd64/
├── packages/
│   ├── deb/
│   ├── rpm/
│   ├── snap/
│   ├── homebrew/
│   └── winget/
├── installers/
│   ├── linux-installer.sh
│   ├── macos-installer.pkg
│   └── windows-installer.exe
└── vm-images/
    ├── servin-linux-amd64.qcow2
    ├── servin-linux-arm64.qcow2
    └── checksums.txt
```

### 4. Release Creation

#### GitHub Release
```bash
# Create release with automated script
./scripts/create-release.sh v1.0.0

# Manual process:
# 1. Create git tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 2. Upload artifacts to GitHub
gh release create v1.0.0 \
    --title "Servin v1.0.0" \
    --notes-file CHANGELOG.md \
    dist/packages/* \
    dist/installers/*
```

#### Release Notes Template
```markdown
# Servin v1.0.0

## What's New
- Universal VM-based containerization
- Cross-platform GUI improvements
- Enhanced container management

## Breaking Changes
- VM mode now default for all platforms
- Configuration file format updated

## Bug Fixes
- Fixed container networking on macOS
- Resolved GUI memory leaks

## Installation
- Universal installer: `curl -fsSL https://get.servin.dev | sh`
- Homebrew: `brew install servin`
- Winget: `winget install servin`

## Checksums
[Include artifact checksums]
```

### 5. Distribution Channels

#### Package Managers

##### Homebrew (macOS/Linux)
```bash
# Update homebrew-servin repository
cd homebrew-servin
./update-formula.sh v1.0.0
git commit -am "Update to v1.0.0"
git push origin main
```

##### Winget (Windows)
```bash
# Update winget-pkgs repository
cd winget-pkgs
./manifests/s/Servin/Servin/update-manifest.sh v1.0.0
# Submit PR to microsoft/winget-pkgs
```

##### APT Repository (Debian/Ubuntu)
```bash
# Upload to APT repository
./scripts/upload-apt.sh v1.0.0
# Updates packages.servin.dev
```

##### RPM Repository (RedHat/SUSE)
```bash
# Upload to RPM repository
./scripts/upload-rpm.sh v1.0.0
# Updates rpm.servin.dev
```

##### Snap Store (Linux)
```bash
# Build and upload snap
snapcraft upload servin_1.0.0_amd64.snap
snapcraft release servin 1.0.0 stable
```

#### Container Registries
```bash
# Push base VM images
./scripts/upload-vm-images.sh v1.0.0
# Uploads to registry.servin.dev
```

### 6. Website and Documentation Updates

#### Website Deployment
```bash
# Update servin.dev
cd servin-website
./deploy.sh v1.0.0
```

#### Documentation Updates
- [ ] Update installation instructions
- [ ] Refresh getting started guide
- [ ] Update API documentation
- [ ] Refresh changelog
- [ ] Update download links

### 7. Post-Release Activities

#### Monitoring
- [ ] Monitor download statistics
- [ ] Check package manager adoption
- [ ] Monitor error reporting
- [ ] Review user feedback

#### Community Updates
- [ ] Announce on GitHub Discussions
- [ ] Update README.md
- [ ] Social media announcements
- [ ] Blog post (if major release)

#### Issue Tracking
- [ ] Create milestone for next version
- [ ] Triage new issues
- [ ] Plan patch releases if needed

## Release Schedule

### Regular Releases
- **Patch releases**: Monthly (or as needed for critical bugs)
- **Minor releases**: Quarterly
- **Major releases**: Annually

### Emergency Releases
- Critical security vulnerabilities
- Major platform compatibility issues
- Data loss bugs

## Release Checklist

### Pre-Release (T-1 week)
- [ ] Feature freeze
- [ ] Documentation review
- [ ] Testing complete
- [ ] Security audit
- [ ] Performance benchmarks

### Release Day
- [ ] Final builds created
- [ ] All tests passing
- [ ] Artifacts uploaded
- [ ] Release notes published
- [ ] Package managers updated
- [ ] Website updated

### Post-Release (T+1 day)
- [ ] Download statistics reviewed
- [ ] Initial user feedback collected
- [ ] Critical issues triaged
- [ ] Next version planning started

## Rollback Procedures

### Package Manager Rollback
```bash
# Homebrew
brew uninstall servin
brew install servin@1.0.0

# APT
apt remove servin
apt install servin=1.0.0

# Winget
winget uninstall servin
winget install --version 1.0.0 servin
```

### Binary Rollback
```bash
# Download previous version
curl -fsSL https://github.com/immyemperor/servin/releases/download/v0.9.0/servin-0.9.0-linux-amd64.tar.gz | tar -xz
```

## Automation

### CI/CD Integration
```yaml
# .github/workflows/release.yml
name: Release
on:
  push:
    tags: ['v*']
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Build Distribution
        run: ./build-vm-distribution.sh --all
      - name: Create Release
        run: ./scripts/create-release.sh ${{ github.ref_name }}
```

### Automated Testing
```yaml
# .github/workflows/test-release.yml
name: Test Release
on:
  release:
    types: [published]
jobs:
  test-installation:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
    runs-on: ${{ matrix.os }}
    steps:
      - name: Test Installation
        run: |
          curl -fsSL https://get.servin.dev | sh
          servin --version
          servin vm status
```

## Metrics and Analytics

### Release Metrics
- Download counts by platform
- Installation success rates
- User adoption rates
- Performance benchmarks
- Error rates and patterns

### Tracking Tools
- GitHub release statistics
- Package manager analytics
- Website analytics
- Error reporting systems
- User feedback systems

## Documentation

### Release Documentation
- [ ] CHANGELOG.md updated
- [ ] VERSION file updated
- [ ] Migration guides (for breaking changes)
- [ ] Installation guides updated
- [ ] API documentation refreshed

### Communication
- [ ] Release announcement prepared
- [ ] User migration guides
- [ ] Developer documentation
- [ ] Support documentation updated