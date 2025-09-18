#!/bin/bash

# Direct QEMU VM Console SSH Setup Script
# This script provides multiple methods to configure SSH in the running Alpine VM

echo "🚀 Automated SSH Setup for Alpine Linux VM"
echo "=========================================="
echo ""

# Check VM status
echo "📊 Checking VM status..."
VM_STATUS=$(./servin vm status | grep "VM Status" | awk '{print $3}')
echo "VM Status: $VM_STATUS"

if [ "$VM_STATUS" != "running" ]; then
    echo "❌ VM is not running. Please start it first:"
    echo "   ./servin vm start"
    exit 1
fi

echo "✅ VM is running"
echo ""

# Method 1: Check if SSH is already working
echo "🔍 Testing SSH connectivity..."
if ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=2 -o BatchMode=yes root@localhost 'echo "SSH working"' 2>/dev/null; then
    echo "✅ SSH is already configured and working!"
    echo "   Connection: ssh root@localhost -p 2222"
    echo ""
    echo "🎯 Ready for container testing:"
    echo "   ./servin run nginx:alpine"
    echo "   ./servin exec <container> sh"
    exit 0
fi

echo "❌ SSH not configured yet"
echo ""

# Method 2: Create a setup script to copy to VM
echo "🛠️  Creating SSH setup script..."
cat > /tmp/vm_ssh_setup.sh << 'EOF'
#!/bin/sh
# SSH Setup Script for Alpine Linux

echo "Starting SSH setup in Alpine Linux..."

# Update package repository
apk update

# Install SSH and sudo
apk add openssh sudo curl

# Configure SSH for remote access
echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config
echo 'PasswordAuthentication yes' >> /etc/ssh/sshd_config
echo 'ClientAliveInterval 60' >> /etc/ssh/sshd_config

# Set root password
echo 'root:servin123' | chpasswd

# Enable SSH service
rc-update add sshd default
rc-service sshd start

# Create servin user for container management
adduser -D -s /bin/ash servin
echo 'servin:servin123' | chpasswd
echo 'servin ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

echo "✅ SSH setup complete!"
echo "SSH is now available on port 22"
echo "Root password: servin123"
echo "User 'servin' password: servin123"
EOF

chmod +x /tmp/vm_ssh_setup.sh

echo "✅ Setup script created: /tmp/vm_ssh_setup.sh"
echo ""

# Method 3: Provide detailed manual instructions
echo "📋 MANUAL SETUP INSTRUCTIONS"
echo "============================="
echo ""
echo "Since the VM console is attached to the QEMU process, you need to:"
echo ""
echo "1. Find the QEMU terminal/console:"
QEMU_PID=$(ps aux | grep qemu-system-aarch64 | grep -v grep | awk '{print $2}')
echo "   QEMU PID: $QEMU_PID"
echo ""

echo "2. The VM console should show: 'localhost login:'"
echo "   Login as: root (press Enter for password)"
echo ""

echo "3. Copy and paste these commands in the VM console:"
echo ""
cat << 'EOF'
   apk update
   apk add openssh sudo
   echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config
   echo 'PasswordAuthentication yes' >> /etc/ssh/sshd_config
   echo 'root:servin123' | chpasswd
   rc-update add sshd default
   rc-service sshd start
EOF

echo ""
echo "4. Test SSH access:"
echo "   ssh root@localhost -p 2222"
echo "   Password: servin123"
echo ""

# Method 4: Alternative restart approach with SSH
echo "🔄 ALTERNATIVE: Restart VM with SSH pre-configured"
echo "================================================"
echo ""
echo "If manual setup is difficult, you can:"
echo "1. Stop current VM:  ./servin vm stop"
echo "2. Restart VM:       ./servin vm start"
echo "3. The new VM will have the same setup requirements"
echo ""

echo "🎯 VERIFICATION"
echo "==============="
echo ""
echo "After SSH setup, verify with:"
echo "   ssh root@localhost -p 2222 'uname -a'"
echo ""
echo "Then test Servin containers:"
echo "   ./servin run hello-world"
echo "   ./servin list"
echo ""

echo "📞 TROUBLESHOOTING"
echo "=================="
echo ""
echo "If SSH still doesn't work:"
echo "• Check VM logs: journalctl -u sshd"
echo "• Verify SSH is running: rc-service sshd status"
echo "• Check network: ip addr show"
echo "• Test local connection: ssh localhost (inside VM)"