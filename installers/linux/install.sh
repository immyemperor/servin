#!/bin/bash
# Servin Container Runtime - Linux Installer Script

set -e

# Configuration
INSTALL_DIR="/usr/local/bin"
DATA_DIR="/var/lib/servin"
CONFIG_DIR="/etc/servin"
LOG_DIR="/var/log/servin"
SERVICE_NAME="servin"
USER="servin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${CYAN}================================================${NC}"
    echo -e "${CYAN}   Servin Container Runtime - Linux Installer${NC}"
    echo -e "${CYAN}================================================${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}$1${NC}"
}

print_warning() {
    echo -e "${YELLOW}$1${NC}"
}

print_error() {
    echo -e "${RED}$1${NC}"
}

print_info() {
    echo -e "${BLUE}$1${NC}"
}

check_root() {
    if [[ $EUID -ne 0 ]]; then
        print_error "This script must be run as root (use sudo)"
        exit 1
    fi
}

detect_init_system() {
    if command -v systemctl >/dev/null 2>&1; then
        echo "systemd"
    elif command -v service >/dev/null 2>&1; then
        echo "sysv"
    else
        echo "unknown"
    fi
}

create_user() {
    print_info "Creating system user: $USER"
    if ! id "$USER" &>/dev/null; then
        useradd --system --no-create-home --shell /bin/false "$USER"
        print_success "  User $USER created"
    else
        print_warning "  User $USER already exists"
    fi
}

create_directories() {
    print_info "Creating directories..."
    
    directories=("$DATA_DIR" "$CONFIG_DIR" "$LOG_DIR" "$DATA_DIR/volumes" "$DATA_DIR/images")
    
    for dir in "${directories[@]}"; do
        if [[ ! -d "$dir" ]]; then
            mkdir -p "$dir"
            print_success "  Created: $dir"
        fi
    done
    
    # Set ownership and permissions
    chown -R "$USER:$USER" "$DATA_DIR" "$LOG_DIR"
    chmod 755 "$DATA_DIR" "$LOG_DIR"
    chmod 755 "$CONFIG_DIR"
}

install_binaries() {
    print_info "Installing binaries..."
    
    # Check if binaries exist in current directory
    if [[ ! -f "servin" ]]; then
        print_error "servin binary not found in current directory"
        exit 1
    fi
    
    # Copy binaries
    cp servin "$INSTALL_DIR/"
    chmod 755 "$INSTALL_DIR/servin"
    print_success "  Installed: $INSTALL_DIR/servin"
    
    if [[ -f "servin-gui" ]]; then
        cp servin-gui "$INSTALL_DIR/"
        chmod 755 "$INSTALL_DIR/servin-gui"
        print_success "  Installed: $INSTALL_DIR/servin-gui"
    fi
}

create_config() {
    print_info "Creating configuration..."
    
    cat > "$CONFIG_DIR/servin.conf" << EOF
# Servin Configuration File
# Data directory
data_dir=$DATA_DIR

# Log settings
log_level=info
log_file=$LOG_DIR/servin.log

# Runtime settings
runtime=native

# Network settings
bridge_name=servin0

# CRI settings
cri_port=10250
cri_enabled=false
EOF
    
    chown root:root "$CONFIG_DIR/servin.conf"
    chmod 644 "$CONFIG_DIR/servin.conf"
    print_success "  Created: $CONFIG_DIR/servin.conf"
}

create_systemd_service() {
    print_info "Creating systemd service..."
    
    cat > "/etc/systemd/system/$SERVICE_NAME.service" << EOF
[Unit]
Description=Servin Container Runtime
Documentation=https://github.com/yourusername/servin
After=network.target
Wants=network.target

[Service]
Type=simple
User=$USER
Group=$USER
ExecStart=$INSTALL_DIR/servin daemon --config $CONFIG_DIR/servin.conf
ExecReload=/bin/kill -HUP \$MAINPID
KillMode=process
Restart=on-failure
RestartSec=5
TimeoutStopSec=30
LimitNOFILE=1048576
LimitNPROC=1048576
LimitCORE=infinity
TasksMax=infinity
Delegate=yes
OOMScoreAdjust=-500

# Security settings
NoNewPrivileges=yes
ProtectHome=yes
ProtectSystem=strict
ReadWritePaths=$DATA_DIR $LOG_DIR

[Install]
WantedBy=multi-user.target
EOF
    
    systemctl daemon-reload
    print_success "  Created systemd service"
}

create_sysv_service() {
    print_info "Creating SysV init script..."
    
    cat > "/etc/init.d/$SERVICE_NAME" << 'EOF'
#!/bin/bash
# Servin Container Runtime
# chkconfig: 35 99 99
# description: Servin Container Runtime

. /etc/rc.d/init.d/functions

USER="servin"
DAEMON="servin"
ROOT_DIR="/var/lib/servin"

SERVER="$ROOT_DIR/$DAEMON"
LOCK_FILE="/var/lock/subsys/$DAEMON"

start() {
    echo -n $"Starting $DAEMON: "
    daemon --user "$USER" --pidfile="$LOCK_FILE" "$SERVER" daemon --config /etc/servin/servin.conf
    RETVAL=$?
    echo
    [ $RETVAL -eq 0 ] && touch $LOCK_FILE
    return $RETVAL
}

stop() {
    echo -n $"Shutting down $DAEMON: "
    pid=`ps -aefw | grep "$DAEMON" | grep -v " grep " | awk '{print $2}'`
    kill -9 $pid > /dev/null 2>&1
    [ $? -eq 0 ] && echo_success || echo_failure
    echo
    [ $RETVAL -eq 0 ] && rm -f $LOCK_FILE
    return $RETVAL
}

restart() {
    stop
    start
}

status() {
    if [ -f $LOCK_FILE ]; then
        echo "$DAEMON is running."
    else
        echo "$DAEMON is stopped."
    fi
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    status)
        status
        ;;
    restart)
        restart
        ;;
    *)
        echo "Usage: {start|stop|status|restart}"
        exit 1
        ;;
esac

exit $?
EOF
    
    chmod +x "/etc/init.d/$SERVICE_NAME"
    
    # Enable service for different distributions
    if command -v chkconfig >/dev/null 2>&1; then
        chkconfig --add "$SERVICE_NAME"
        chkconfig "$SERVICE_NAME" on
    elif command -v update-rc.d >/dev/null 2>&1; then
        update-rc.d "$SERVICE_NAME" defaults
    fi
    
    print_success "  Created SysV init script"
}

create_desktop_entry() {
    print_info "Creating desktop entry..."
    
    if [[ -f "$INSTALL_DIR/servin-gui" ]]; then
        cat > "/usr/share/applications/servin-gui.desktop" << EOF
[Desktop Entry]
Name=Servin GUI
Comment=Servin Container Runtime GUI
Exec=$INSTALL_DIR/servin-gui
Icon=container
Terminal=false
Type=Application
Categories=System;
StartupNotify=true
EOF
        
        chmod 644 "/usr/share/applications/servin-gui.desktop"
        print_success "  Created desktop entry"
    fi
}

create_uninstaller() {
    print_info "Creating uninstaller..."
    
    cat > "$INSTALL_DIR/servin-uninstall" << EOF
#!/bin/bash
# Servin Uninstaller

echo "Uninstalling Servin Container Runtime..."

# Stop and disable service
if command -v systemctl >/dev/null 2>&1; then
    systemctl stop $SERVICE_NAME 2>/dev/null || true
    systemctl disable $SERVICE_NAME 2>/dev/null || true
    rm -f /etc/systemd/system/$SERVICE_NAME.service
    systemctl daemon-reload
elif command -v service >/dev/null 2>&1; then
    service $SERVICE_NAME stop 2>/dev/null || true
    if command -v chkconfig >/dev/null 2>&1; then
        chkconfig --del $SERVICE_NAME 2>/dev/null || true
    elif command -v update-rc.d >/dev/null 2>&1; then
        update-rc.d -f $SERVICE_NAME remove 2>/dev/null || true
    fi
    rm -f /etc/init.d/$SERVICE_NAME
fi

# Remove user
userdel $USER 2>/dev/null || true

# Remove files and directories
rm -f $INSTALL_DIR/servin
rm -f $INSTALL_DIR/servin-gui
rm -f $INSTALL_DIR/servin-uninstall
rm -rf $CONFIG_DIR
rm -rf $DATA_DIR
rm -rf $LOG_DIR
rm -f /usr/share/applications/servin-gui.desktop

echo "Servin has been uninstalled."
EOF
    
    chmod +x "$INSTALL_DIR/servin-uninstall"
    print_success "  Created uninstaller: $INSTALL_DIR/servin-uninstall"
}

main() {
    print_header
    
    check_root
    
    INIT_SYSTEM=$(detect_init_system)
    print_info "Detected init system: $INIT_SYSTEM"
    
    create_user
    create_directories
    install_binaries
    create_config
    
    case "$INIT_SYSTEM" in
        systemd)
            create_systemd_service
            ;;
        sysv)
            create_sysv_service
            ;;
        *)
            print_warning "Unknown init system. Service not installed."
            ;;
    esac
    
    create_desktop_entry
    create_uninstaller
    
    echo ""
    print_success "================================================"
    print_success "   Installation completed successfully!"
    print_success "================================================"
    echo ""
    echo "Installation directory: $INSTALL_DIR"
    echo "Data directory: $DATA_DIR"
    echo "Configuration: $CONFIG_DIR/servin.conf"
    echo ""
    print_info "Next steps:"
    
    case "$INIT_SYSTEM" in
        systemd)
            echo "1. Enable service: sudo systemctl enable $SERVICE_NAME"
            echo "2. Start service: sudo systemctl start $SERVICE_NAME"
            echo "3. Check status: sudo systemctl status $SERVICE_NAME"
            ;;
        sysv)
            echo "1. Start service: sudo service $SERVICE_NAME start"
            echo "2. Check status: sudo service $SERVICE_NAME status"
            ;;
        *)
            echo "1. Run manually: sudo -u $USER $INSTALL_DIR/servin daemon --config $CONFIG_DIR/servin.conf"
            ;;
    esac
    
    if [[ -f "$INSTALL_DIR/servin-gui" ]]; then
        echo "4. Use GUI: servin-gui (or find 'Servin GUI' in applications menu)"
    fi
    echo ""
    echo "Check logs at: $LOG_DIR/servin.log"
    echo "To uninstall: sudo $INSTALL_DIR/servin-uninstall"
}

# Handle command line arguments
case "${1:-}" in
    uninstall)
        if [[ -f "$INSTALL_DIR/servin-uninstall" ]]; then
            exec "$INSTALL_DIR/servin-uninstall"
        else
            print_error "Uninstaller not found"
            exit 1
        fi
        ;;
    *)
        main
        ;;
esac
