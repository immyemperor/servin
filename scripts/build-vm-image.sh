#!/bin/bash
# Build lightweight VM images for Servin containerization

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
IMAGE_SIZE="2G"
MEMORY="512M"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_info() {
    echo -e "${BLUE}→ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

show_usage() {
    cat << EOF
Usage: $0 <distro> <version> <output_file>

Arguments:
    distro      Distribution (alpine, ubuntu, debian)
    version     Distribution version (e.g., 3.18, 22.04)
    output_file Output QCOW2 file path

Examples:
    $0 alpine 3.18 vm-images/servin-alpine.qcow2
    $0 ubuntu 22.04 vm-images/servin-ubuntu.qcow2

EOF
}

# Parse arguments
if [ $# -ne 3 ]; then
    show_usage
    exit 1
fi

DISTRO="$1"
VERSION="$2"
OUTPUT_FILE="$3"

# Validate distro
case "$DISTRO" in
    alpine|ubuntu|debian) ;;
    *)
        print_error "Unsupported distribution: $DISTRO"
        show_usage
        exit 1
        ;;
esac

print_info "Building $DISTRO $VERSION VM image..."

# Create output directory
mkdir -p "$(dirname "$OUTPUT_FILE")"

# Create temporary workspace
TEMP_DIR=$(mktemp -d)
trap 'rm -rf "$TEMP_DIR"' EXIT

build_alpine_image() {
    local version="$1"
    local output="$2"
    
    print_info "Building Alpine Linux $version VM image..."
    
    # Create base image
    qemu-img create -f qcow2 "$output" "$IMAGE_SIZE"
    
    # Download Alpine ISO
    local iso_url="https://dl-cdn.alpinelinux.org/alpine/v${version%.*}/releases/x86_64/alpine-virt-${version}-x86_64.iso"
    local iso_file="$TEMP_DIR/alpine.iso"
    
    print_info "Downloading Alpine ISO..."
    curl -L "$iso_url" -o "$iso_file"
    
    # Create cloud-init configuration for automated install
    cat > "$TEMP_DIR/user-data" << 'EOF'
#cloud-config
users:
  - name: servin
    sudo: ALL=(ALL) NOPASSWD:ALL
    shell: /bin/ash
    ssh_authorized_keys: []

packages:
  - docker
  - docker-compose
  - curl
  - wget

runcmd:
  - rc-update add docker default
  - service docker start
  - addgroup servin docker
EOF

    cat > "$TEMP_DIR/meta-data" << EOF
instance-id: servin-alpine-vm
local-hostname: servin-alpine
EOF

    # Create ISO with cloud-init data
    genisoimage -output "$TEMP_DIR/cloud-init.iso" -volid cidata -joliet -rock "$TEMP_DIR/user-data" "$TEMP_DIR/meta-data" 2>/dev/null || \
    mkisofs -o "$TEMP_DIR/cloud-init.iso" -V cidata -J -r "$TEMP_DIR/user-data" "$TEMP_DIR/meta-data"
    
    print_success "Alpine VM image created: $output"
}

build_ubuntu_image() {
    local version="$1"
    local output="$2"
    
    print_info "Building Ubuntu $version VM image..."
    
    # Create base image
    qemu-img create -f qcow2 "$output" "$IMAGE_SIZE"
    
    # Download Ubuntu cloud image
    local cloud_img_url="https://cloud-images.ubuntu.com/releases/${version}/release/ubuntu-${version}-server-cloudimg-amd64.img"
    local cloud_img="$TEMP_DIR/ubuntu-cloud.img"
    
    print_info "Downloading Ubuntu cloud image..."
    curl -L "$cloud_img_url" -o "$cloud_img"
    
    # Resize and convert to qcow2
    qemu-img convert -f qcow2 -O qcow2 "$cloud_img" "$output"
    qemu-img resize "$output" "$IMAGE_SIZE"
    
    # Create cloud-init configuration
    cat > "$TEMP_DIR/user-data" << 'EOF'
#cloud-config
users:
  - name: servin
    sudo: ALL=(ALL) NOPASSWD:ALL
    shell: /bin/bash
    ssh_authorized_keys: []

packages:
  - docker.io
  - docker-compose
  - curl
  - wget

runcmd:
  - systemctl enable docker
  - systemctl start docker
  - usermod -aG docker servin
  - apt-get clean
EOF

    cat > "$TEMP_DIR/meta-data" << EOF
instance-id: servin-ubuntu-vm
local-hostname: servin-ubuntu
EOF

    print_success "Ubuntu VM image created: $output"
}

build_debian_image() {
    local version="$1"
    local output="$2"
    
    print_info "Building Debian $version VM image..."
    
    # Create base image
    qemu-img create -f qcow2 "$output" "$IMAGE_SIZE"
    
    # Download Debian cloud image
    local cloud_img_url="https://cloud.debian.org/images/cloud/bookworm/latest/debian-12-generic-amd64.qcow2"
    local cloud_img="$TEMP_DIR/debian-cloud.qcow2"
    
    print_info "Downloading Debian cloud image..."
    curl -L "$cloud_img_url" -o "$cloud_img"
    
    # Copy and resize
    cp "$cloud_img" "$output"
    qemu-img resize "$output" "$IMAGE_SIZE"
    
    print_success "Debian VM image created: $output"
}

# Build the VM image based on distro
case "$DISTRO" in
    alpine)
        build_alpine_image "$VERSION" "$OUTPUT_FILE"
        ;;
    ubuntu)
        build_ubuntu_image "$VERSION" "$OUTPUT_FILE"
        ;;
    debian)
        build_debian_image "$VERSION" "$OUTPUT_FILE"
        ;;
esac

# Verify the image
if [ -f "$OUTPUT_FILE" ]; then
    print_info "Verifying VM image..."
    qemu-img info "$OUTPUT_FILE"
    print_success "VM image build complete: $OUTPUT_FILE"
else
    print_error "Failed to create VM image"
    exit 1
fi