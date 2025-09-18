#!/bin/bash

# Automated Build Script with SSH Setup for Servin VM
# This script builds Servin and automatically configures VM with SSH

set -e  # Exit on any error

echo "🚀 Servin Build Script with Automated VM SSH Setup"
echo "=================================================="
echo ""

# Build Servin binary
echo "📦 Building Servin binary..."
go build -o servin main.go
echo "✅ Servin binary built successfully"
echo ""

# Clean any existing VM to start fresh
echo "🧹 Cleaning existing VM data..."
rm -rf ~/.servin/vms/servin-vm 2>/dev/null || true
echo "✅ VM data cleaned"
echo ""

# Start VM with automated SSH setup
echo "🚀 Starting VM with automated SSH configuration..."
./servin vm start

echo ""

# Monitor SSH setup progress
echo "⏳ Monitoring SSH setup progress..."
SSH_READY=false
MAX_WAIT=90

for i in $(seq 1 $MAX_WAIT); do
    if ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=2 -o BatchMode=yes root@localhost 'echo "SSH_WORKING"' 2>/dev/null | grep -q "SSH_WORKING"; then
        SSH_READY=true
        echo "✅ SSH is ready after $i seconds!"
        break
    fi
    
    # Show progress every 10 seconds
    if [ $((i % 10)) -eq 0 ]; then
        echo "   Waiting for SSH... ($i/$MAX_WAIT seconds)"
    fi
    
    sleep 1
done

if [ "$SSH_READY" = true ]; then
    echo ""
    echo "🎯 VM Setup Complete!"
    echo "===================="
    echo ""
    
    # Get VM information
    echo "📊 VM Information:"
    echo "   Status: $(./servin vm status | grep "VM Status" | awk '{print $3}')"
    echo "   SSH: ssh root@localhost -p 2222"
    echo "   Password: servin123"
    echo ""
    
    # Test VM connectivity
    echo "🔍 Testing VM connectivity..."
    VM_KERNEL=$(ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@localhost 'uname -r' 2>/dev/null)
    VM_DISTRO=$(ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@localhost 'cat /etc/alpine-release' 2>/dev/null)
    echo "   Kernel: $VM_KERNEL"
    echo "   Distribution: Alpine Linux $VM_DISTRO"
    echo ""
    
    # Deploy Servin to VM
    echo "📦 Deploying Servin to VM..."
    if scp -P 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ./servin root@localhost:/usr/local/bin/ 2>/dev/null; then
        ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@localhost 'chmod +x /usr/local/bin/servin' 2>/dev/null
        echo "✅ Servin deployed to VM successfully"
    else
        echo "⚠️  Failed to deploy Servin to VM (manual deployment may be needed)"
    fi
    
    echo ""
    echo "🧪 Testing Container Functionality..."
    echo "===================================="
    echo ""
    
    # Test container operations
    echo "1. Testing hello-world container:"
    if ./servin run --name test-hello hello-world 2>/dev/null; then
        echo "✅ Container run successful"
    else
        echo "❌ Container run failed (may need manual SSH setup)"
    fi
    
    echo ""
    echo "2. Listing containers:"
    ./servin list 2>/dev/null || echo "❌ Container list failed"
    
    echo ""
    echo "3. Testing container logs:"
    ./servin logs test-hello 2>/dev/null || echo "❌ Container logs failed"
    
    echo ""
    echo "🎉 Build and VM Setup Complete!"
    echo "==============================="
    echo ""
    echo "🎯 Ready for Development:"
    echo "   • VM Status: Running with SSH"
    echo "   • Container Runtime: Native Linux (not Docker simulation)"
    echo "   • SSH Access: ssh root@localhost -p 2222"
    echo "   • Servin Commands: ./servin run, ./servin exec, ./servin logs"
    echo ""
    echo "📚 Example Commands:"
    echo "   ./servin run nginx:alpine"
    echo "   ./servin run --name web -p 8080:80 nginx:alpine"
    echo "   ./servin exec web sh"
    echo "   ./servin logs web"
    echo ""
    
else
    echo ""
    echo "⚠️  SSH Auto-Setup Incomplete"
    echo "============================="
    echo ""
    echo "The VM is running but SSH auto-setup didn't complete within $MAX_WAIT seconds."
    echo ""
    echo "Manual setup required:"
    echo "1. Connect to VM console"
    echo "2. Login as root (no password needed)"
    echo "3. Mount and run setup script:"
    echo "   mount /dev/sr0 /mnt 2>/dev/null || true"
    echo "   /mnt/autosetup.sh"
    echo ""
    echo "Alternative manual commands:"
    echo "   apk update && apk add openssh"
    echo "   echo 'root:servin123' | chpasswd"
    echo "   echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config"
    echo "   rc-update add sshd default && rc-service sshd start"
    echo ""
    echo "Then test: ssh root@localhost -p 2222"
    echo ""
fi

echo "🏁 Build script completed!"
echo ""
echo "VM Status: $(./servin vm status | grep "VM Status" | awk '{print $3}')"
echo "QEMU Process: $(ps aux | grep qemu-system-aarch64 | grep -v grep | awk '{print "PID " $2}' || echo "Not found")"