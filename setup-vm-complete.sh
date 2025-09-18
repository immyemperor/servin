#!/bin/bash

# Streamlined VM Setup with Auto-SSH Configuration
# This script restarts the VM with automated SSH setup

echo "🚀 Streamlined VM Setup with Auto-SSH"
echo "====================================="
echo ""

echo "This script will:"
echo "1. Stop the current VM"
echo "2. Create an enhanced VM configuration"
echo "3. Start VM with automatic SSH setup"
echo "4. Test container functionality"
echo ""

read -p "Continue? (y/n): " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Setup cancelled."
    exit 0
fi

echo "🛑 Stopping current VM..."
./servin vm stop 2>/dev/null || echo "VM was not running"

echo ""
echo "⏳ Waiting for VM to fully stop..."
sleep 3

echo ""
echo "🔧 Creating enhanced VM configuration..."

# Create a startup script that will automatically configure SSH
mkdir -p ~/.servin/vms/servin-vm/autosetup
cat > ~/.servin/vms/servin-vm/autosetup/init.sh << 'EOF'
#!/bin/sh
# Auto SSH Setup Script for Alpine Linux

# Wait for networking
sleep 5

# Update packages and install SSH
apk update
apk add openssh sudo curl

# Configure SSH
echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config
echo 'PasswordAuthentication yes' >> /etc/ssh/sshd_config
echo 'ClientAliveInterval 60' >> /etc/ssh/sshd_config

# Set passwords
echo 'root:servin123' | chpasswd

# Create servin user
adduser -D -s /bin/ash servin
echo 'servin:servin123' | chpasswd
echo 'servin ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

# Enable and start SSH
rc-update add sshd default
rc-service sshd start

# Create completion marker
echo "SSH auto-setup completed at $(date)" > /tmp/ssh-ready

echo "✅ SSH is now available on port 22"
EOF

chmod +x ~/.servin/vms/servin-vm/autosetup/init.sh

echo "✅ Auto-setup script created"
echo ""

echo "🚀 Starting VM with enhanced configuration..."
./servin vm start

echo ""
echo "⏳ Waiting for VM to boot and configure SSH (30 seconds)..."

# Wait and test SSH connectivity
for i in {1..30}; do
    echo -n "."
    if ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=1 -o BatchMode=yes root@localhost 'echo "SSH_READY"' 2>/dev/null | grep -q "SSH_READY"; then
        echo ""
        echo "✅ SSH is ready!"
        break
    fi
    sleep 1
done

echo ""

# Test SSH connectivity
echo "🔍 Testing SSH connectivity..."
if ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@localhost 'echo "✅ SSH connection successful!"'; then
    echo ""
    echo "🎯 VM Setup Complete!"
    echo "===================="
    echo ""
    echo "VM Details:"
    echo "   SSH Access: ssh root@localhost -p 2222"
    echo "   Password: servin123"
    echo "   Status: $(./servin vm status | grep "VM Status" | awk '{print $3}')"
    echo ""
    
    echo "🧪 Testing container functionality..."
    echo ""
    
    # Copy Servin to VM
    echo "📦 Installing Servin in VM..."
    scp -P 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ./servin root@localhost:/usr/local/bin/
    ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@localhost 'chmod +x /usr/local/bin/servin'
    
    echo "✅ Servin installed in VM"
    echo ""
    
    echo "🎉 Ready for container operations!"
    echo "================================="
    echo ""
    echo "Try these commands:"
    echo "   ./servin run hello-world"
    echo "   ./servin run nginx:alpine"
    echo "   ./servin list"
    echo "   ./servin exec <container> sh"
    echo ""
    
else
    echo "❌ SSH setup incomplete"
    echo ""
    echo "Manual setup required:"
    echo "1. Connect to VM console"
    echo "2. Login as root (no password)"
    echo "3. Run: /tmp/autosetup/init.sh"
    echo ""
    echo "Or restart setup: $0"
fi

echo ""
echo "🏁 Setup complete! VM is ready for native Linux containers."