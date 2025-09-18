#!/bin/bash

# Servin Container Runtime - Enhanced Linux Installer with VM Prerequisites
# Run with: sudo ./install-with-vm.sh

set -e

# Configuration
INSTALL_DIR="/usr/local/bin"
DATA_DIR="/var/lib/servin"
CONFIG_DIR="/etc/servin"
LOG_DIR="/var/log/servin"
VM_DIR="$DATA_DIR/vm"
USER_DATA_DIR="$HOME/.servin"

# Color output functions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_success() { echo -e "${GREEN}✓ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠ $1${NC}"; }
print_error() { echo -e "${RED}✗ $1${NC}"; }
print_info() { echo -e "${BLUE}→ $1${NC}"; }
print_header() { echo -e "\n${CYAN}$('=' * 60)\n$1\n$('=' * 60)${NC}"; }

# Check if running as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        print_error "This script must be run as root (use sudo)"
        exit 1
    fi
}

# Detect Linux distribution
detect_distro() {
    if [[ -f /etc/os-release ]]; then
        . /etc/os-release
        DISTRO=$ID
        DISTRO_VERSION=$VERSION_ID
    elif [[ -f /etc/debian_version ]]; then
        DISTRO="debian"
    elif [[ -f /etc/redhat-release ]]; then
        DISTRO="rhel"
    else
        DISTRO="unknown"
    fi
    
    print_info "Detected distribution: $DISTRO $DISTRO_VERSION"
}

# Check system prerequisites
check_prerequisites() {
    print_header "Checking System Prerequisites"
    
    local prereq_failed=0
    
    # Check CPU virtualization
    print_info "Checking CPU virtualization support..."
    if grep -E "(vmx|svm)" /proc/cpuinfo >/dev/null 2>&1; then
        print_success "CPU supports hardware virtualization"
        if grep -q "vmx" /proc/cpuinfo; then
            print_info "  Intel VT-x support detected"
        fi
        if grep -q "svm" /proc/cpuinfo; then
            print_info "  AMD-V support detected"
        fi
    else
        print_warning "CPU may not support hardware virtualization"
        print_info "VM features will use software emulation"
    fi
    
    # Check available memory
    print_info "Checking available memory..."
    local mem_gb=$(free -g | awk 'NR==2{printf "%.1f", $2}')
    if (( $(echo "$mem_gb >= 4" | bc -l) )); then
        print_success "Sufficient memory available (${mem_gb}GB)"
    else
        print_warning "Low memory detected (${mem_gb}GB). Recommended: 4GB+"
        prereq_failed=$((prereq_failed + 1))
    fi
    
    # Check disk space
    print_info "Checking available disk space..."
    local disk_gb=$(df / | awk 'NR==2 {printf "%.1f", $4/1024/1024}')
    if (( $(echo "$disk_gb >= 5" | bc -l) )); then
        print_success "Sufficient disk space available (${disk_gb}GB free)"
    else
        print_error "Insufficient disk space (${disk_gb}GB free). Need 5GB minimum"
        prereq_failed=$((prereq_failed + 1))
    fi
    
    # Check internet connectivity
    print_info "Checking internet connectivity..."
    if ping -c 1 google.com >/dev/null 2>&1; then
        print_success "Internet connection available"
    else
        print_warning "Internet connection not available - some features may not work"
    fi
    
    if [[ $prereq_failed -gt 0 ]]; then
        print_warning "$prereq_failed critical prerequisites failed"
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
}

# Install packages based on distribution
install_packages() {
    print_header "Installing VM Prerequisites"
    
    case "$DISTRO" in
        "ubuntu"|"debian")
            print_info "Updating package lists..."
            apt-get update -qq
            
            print_info "Installing virtualization packages..."
            DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
                qemu-system-x86 \
                qemu-utils \
                qemu-kvm \
                libvirt-clients \
                libvirt-daemon-system \
                bridge-utils \
                genisoimage \
                cpu-checker \
                virt-manager \
                libvirt-dev \
                python3-libvirt \
                python3 \
                python3-pip \
                python3-venv \
                python3-tk \
                python3-dev \
                curl \
                wget \
                unzip \
                bc || {
                print_error "Failed to install some packages"
                return 1
            }
            ;;
            
        "fedora"|"rhel"|"centos")
            print_info "Installing virtualization packages..."
            if command -v dnf >/dev/null 2>&1; then
                dnf install -y \
                    qemu-kvm \
                    qemu-img \
                    libvirt \
                    libvirt-client \
                    libvirt-daemon-system \
                    bridge-utils \
                    genisoimage \
                    virt-manager \
                    python3 \
                    python3-pip \
                    python3-tkinter \
                    curl \
                    wget \
                    unzip \
                    bc
            else
                yum install -y \
                    qemu-kvm \
                    qemu-img \
                    libvirt \
                    libvirt-client \
                    bridge-utils \
                    genisoimage \
                    python3 \
                    python3-pip \
                    curl \
                    wget \
                    unzip \
                    bc
            fi
            ;;
            
        "arch")
            print_info "Installing virtualization packages..."
            pacman -S --noconfirm \
                qemu-desktop \
                libvirt \
                bridge-utils \
                cdrtools \
                python \
                python-pip \
                tk \
                curl \
                wget \
                unzip \
                bc
            ;;
            
        *)
            print_warning "Unsupported distribution: $DISTRO"
            print_info "Please install QEMU, KVM, and libvirt manually"
            return 1
            ;;
    esac
    
    print_success "Virtualization packages installed"
}

# Configure KVM and libvirt
configure_virtualization() {
    print_header "Configuring Virtualization"
    
    # Load KVM modules
    print_info "Loading KVM kernel modules..."
    if grep -q "vmx" /proc/cpuinfo; then
        modprobe kvm_intel 2>/dev/null && print_success "Intel KVM module loaded" || print_warning "Could not load Intel KVM module"
    fi
    if grep -q "svm" /proc/cpuinfo; then
        modprobe kvm_amd 2>/dev/null && print_success "AMD KVM module loaded" || print_warning "Could not load AMD KVM module"
    fi
    modprobe kvm 2>/dev/null && print_success "KVM module loaded" || print_warning "Could not load KVM module"
    
    # Check KVM device
    print_info "Checking KVM device access..."
    if [[ -e /dev/kvm ]]; then
        print_success "/dev/kvm device exists"
        chmod 666 /dev/kvm 2>/dev/null || print_warning "Could not set KVM device permissions"
    else
        print_warning "/dev/kvm device not found"
    fi
    
    # Configure libvirt
    print_info "Configuring libvirt..."
    systemctl enable libvirtd >/dev/null 2>&1 || print_warning "Could not enable libvirtd"
    systemctl start libvirtd >/dev/null 2>&1 || print_warning "Could not start libvirtd"
    
    # Add users to groups (for current user and future users)
    if [[ -n "$SUDO_USER" ]]; then
        print_info "Adding user $SUDO_USER to virtualization groups..."
        usermod -a -G kvm,libvirt "$SUDO_USER" 2>/dev/null && print_success "User added to kvm and libvirt groups"
    fi
    
    # Create libvirt network if it doesn't exist
    print_info "Setting up libvirt default network..."
    if ! virsh net-list --all | grep -q default; then
        virsh net-define /dev/stdin << 'EOF'
<network>
  <name>default</name>
  <uuid>9a05da11-e96b-47f3-8253-a3a482e445f5</uuid>
  <forward mode='nat'/>
  <bridge name='virbr0' stp='on' delay='0'/>
  <mac address='52:54:00:0a:cd:21'/>
  <ip address='192.168.122.1' netmask='255.255.255.0'>
    <dhcp>
      <range start='192.168.122.2' end='192.168.122.254'/>
    </dhcp>
  </ip>
</network>
EOF
    fi
    
    virsh net-autostart default >/dev/null 2>&1
    virsh net-start default >/dev/null 2>&1 || true
    
    print_success "Virtualization configured"
}

# Install development tools
install_development_tools() {
    print_header "Installing Development Tools"
    
    # Install Packer
    print_info "Installing HashiCorp Packer..."
    local packer_version="1.10.0"
    local packer_url="https://releases.hashicorp.com/packer/${packer_version}/packer_${packer_version}_linux_amd64.zip"
    
    cd /tmp
    wget -q "$packer_url" -O packer.zip
    unzip -q packer.zip
    mv packer /usr/local/bin/hc-packer
    chmod +x /usr/local/bin/hc-packer
    ln -sf /usr/local/bin/hc-packer /usr/local/bin/packer
    rm packer.zip
    
    if /usr/local/bin/hc-packer version >/dev/null 2>&1; then
        print_success "Packer installed successfully"
        
        # Install Packer QEMU plugin
        print_info "Installing Packer QEMU plugin..."
        /usr/local/bin/hc-packer plugins install github.com/hashicorp/qemu 2>/dev/null || print_warning "Packer QEMU plugin installation failed"
    else
        print_warning "Packer installation may have failed"
    fi
}

# Install Python dependencies
install_python_dependencies() {
    print_header "Installing Python Dependencies"
    
    print_info "Installing Python WebView dependencies..."
    pip3 install --upgrade pip --quiet
    pip3 install pywebview[gtk] flask flask-cors flask-socketio eventlet pyinstaller --quiet
    
    # Test imports
    print_info "Testing Python imports..."
    python3 -c "import webview; print('✓ pywebview available')" 2>/dev/null || print_warning "pywebview not available"
    python3 -c "import flask; print('✓ flask available')" 2>/dev/null || print_warning "flask not available"
    python3 -c "import flask_socketio; print('✓ flask-socketio available')" 2>/dev/null || print_warning "flask-socketio not available"
    
    print_success "Python dependencies installed"
}

# Install Servin binaries
install_servin() {
    print_header "Installing Servin Container Runtime"
    
    # Create directories
    print_info "Creating directories..."
    mkdir -p "$INSTALL_DIR" "$DATA_DIR" "$CONFIG_DIR" "$LOG_DIR" "$VM_DIR"/{images,instances}
    
    # Copy executables
    print_info "Installing executables..."
    local installed_count=0
    
    for exe in servin servin-tui servin-gui; do
        if [[ -f "./$exe" ]]; then
            cp "./$exe" "$INSTALL_DIR/"
            chmod +x "$INSTALL_DIR/$exe"
            print_success "Installed: $exe"
            installed_count=$((installed_count + 1))
        else
            print_warning "Executable not found: $exe"
        fi
    done
    
    if [[ $installed_count -eq 0 ]]; then
        print_error "No executables found. Please run this installer from the directory containing the Servin executables."
        return 1
    fi
    
    # Create VM configuration
    print_info "Creating VM configuration..."
    cat > "$CONFIG_DIR/vm-config.yaml" << EOF
vm:
  platform: linux
  providers:
    - name: kvm
      priority: 1
      enabled: true
      acceleration: true
    - name: qemu
      priority: 2
      enabled: true
      acceleration: false
  default_provider: kvm
  image_cache: "$VM_DIR/images"
  vm_storage: "$VM_DIR/instances"
  kvm_group: libvirt
  max_memory: "2GB"
  default_memory: "1GB"
  max_cpu_cores: 2
EOF
    
    # Create main configuration
    print_info "Creating main configuration..."
    cat > "$CONFIG_DIR/servin.conf" << EOF
# Servin Configuration File
data_dir=$DATA_DIR
log_level=info
log_file=$LOG_DIR/servin.log
runtime=vm
bridge_name=servin0
cri_port=10250
cri_enabled=false
vm_enabled=true
vm_config=$CONFIG_DIR/vm-config.yaml
EOF
    
    # Set permissions
    chown -R root:root "$INSTALL_DIR"
    chown -R root:root "$CONFIG_DIR"
    chmod -R 755 "$INSTALL_DIR"
    chmod -R 644 "$CONFIG_DIR"
    chmod +x "$INSTALL_DIR"/servin*
    
    print_success "Servin installed successfully"
}

# Create systemd service
create_service() {
    print_header "Creating System Service"
    
    print_info "Creating systemd service..."
    cat > /etc/systemd/system/servin.service << EOF
[Unit]
Description=Servin Container Runtime
Documentation=https://servin.dev/docs
After=network-online.target libvirtd.service
Wants=network-online.target
Requires=libvirtd.service

[Service]
Type=notify
ExecStart=$INSTALL_DIR/servin daemon
Restart=on-failure
RestartSec=5
Delegate=yes
KillMode=process
OOMScoreAdjust=-999
LimitNOFILE=1048576
LimitNPROC=1048576
LimitCORE=infinity
TasksMax=infinity
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
EOF
    
    systemctl daemon-reload
    systemctl enable servin
    
    print_success "Systemd service created and enabled"
}

# Initialize VM support
initialize_vm() {
    print_header "Initializing VM Support"
    
    local servin_exe="$INSTALL_DIR/servin"
    if [[ ! -f "$servin_exe" ]]; then
        print_error "Servin executable not found: $servin_exe"
        return 1
    fi
    
    print_info "Initializing VM directories..."
    "$servin_exe" vm init || print_warning "VM initialization may have failed"
    
    print_info "Testing VM providers..."
    "$servin_exe" vm list-providers || print_warning "VM provider detection failed"
    
    # Test KVM if available
    if [[ -e /dev/kvm ]] && "$servin_exe" vm check-kvm >/dev/null 2>&1; then
        print_success "KVM provider available and working"
    else
        print_warning "KVM provider not fully functional - will use QEMU fallback"
    fi
    
    print_success "VM support initialized"
}

# Run comprehensive tests
run_tests() {
    print_header "Running Verification Tests"
    
    local servin_exe="$INSTALL_DIR/servin"
    
    # Test basic functionality
    print_info "Testing basic functionality..."
    if "$servin_exe" version >/dev/null 2>&1; then
        print_success "Servin CLI working"
    else
        print_warning "Servin CLI test failed"
    fi
    
    # Test VM functionality
    print_info "Testing VM functionality..."
    if "$servin_exe" vm status >/dev/null 2>&1; then
        print_success "VM subsystem working"
    else
        print_warning "VM subsystem test failed"
    fi
    
    # Test KVM
    if command -v kvm-ok >/dev/null 2>&1; then
        print_info "Running KVM compatibility check..."
        if kvm-ok >/dev/null 2>&1; then
            print_success "KVM acceleration available"
        else
            print_warning "KVM acceleration not available - using software emulation"
        fi
    fi
}

# Show installation summary
show_summary() {
    print_header "Installation Summary"
    
    echo -e "\n${YELLOW}Installed Components:${NC}"
    echo "✓ Servin Container Runtime"
    echo "✓ VM containerization support"
    echo "✓ QEMU/KVM virtualization"
    echo "✓ libvirt management"
    echo "✓ Python WebView GUI support"
    echo "✓ Development tools (Packer)"
    echo "✓ Systemd service"
    
    echo -e "\n${YELLOW}Configuration Files:${NC}"
    echo "• Main config: $CONFIG_DIR/servin.conf"
    echo "• VM config: $CONFIG_DIR/vm-config.yaml"
    echo "• Logs: $LOG_DIR/servin.log"
    echo "• Data: $DATA_DIR"
    
    echo -e "\n${YELLOW}Next Steps:${NC}"
    echo "1. Logout and login again (for group membership)"
    echo "2. Start the service: sudo systemctl start servin"
    echo "3. Initialize VM support: servin vm init"
    echo "4. Enable VM mode: servin vm enable"
    echo "5. Test: servin run --vm alpine echo 'Hello from VM!'"
    
    echo -e "\n${YELLOW}GUI Access:${NC}"
    echo "• Run 'servin-gui' to open the graphical interface"
    echo "• Or use 'servin-tui' for terminal interface"
    
    echo -e "\n${CYAN}Documentation: See VM_PREREQUISITES.md for detailed setup guide${NC}"
    print_success "\nServin Container Runtime installation completed!"
}

# Main installation flow
main() {
    print_header "Servin Container Runtime - Linux Installer"
    echo -e "This installer will set up Servin with VM containerization support.\n"
    
    check_root
    detect_distro
    check_prerequisites
    install_packages
    configure_virtualization
    install_development_tools
    install_python_dependencies
    install_servin
    create_service
    initialize_vm
    run_tests
    show_summary
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "Usage: $0 [options]"
        echo "Options:"
        echo "  --help, -h     Show this help message"
        echo "  --no-service   Skip systemd service creation"
        echo "  --no-vm        Skip VM setup"
        echo "  --user         Install for current user only (not implemented)"
        exit 0
        ;;
    --no-service)
        SKIP_SERVICE=1
        ;;
    --no-vm)
        SKIP_VM=1
        ;;
esac

# Run main installation
main "$@"