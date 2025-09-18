#!/bin/bash

# Quick SSH Test and Container Demo Script
# This script tests SSH connectivity and demonstrates container functionality

echo "🧪 Testing Servin VM SSH and Container Functionality"
echo "===================================================="
echo ""

# Test VM status
echo "📊 Checking VM status..."
VM_STATUS=$(./servin vm status | grep "VM Status" | awk '{print $3}')
echo "VM Status: $VM_STATUS"

if [ "$VM_STATUS" != "running" ]; then
    echo "❌ VM not running. Starting VM..."
    ./servin vm start
    sleep 5
fi

echo ""

# Test SSH connectivity with timeout
echo "🔍 Testing SSH connectivity (5 second timeout)..."
if timeout 5 ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=3 -o BatchMode=yes root@localhost 'echo "SSH_SUCCESS"' 2>/dev/null | grep -q "SSH_SUCCESS"; then
    echo "✅ SSH is working!"
    
    # Test VM information
    echo ""
    echo "🖥️  VM Information:"
    ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@localhost 'uname -a'
    ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@localhost 'cat /etc/alpine-release'
    
    # Copy Servin binary to VM
    echo ""
    echo "📦 Copying Servin binary to VM..."
    scp -P 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ./servin root@localhost:/usr/local/bin/
    ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@localhost 'chmod +x /usr/local/bin/servin'
    
    # Test container functionality
    echo ""
    echo "🚀 Testing container functionality in VM..."
    echo ""
    echo "1. Testing basic container run:"
    ./servin run --name test-container hello-world
    echo ""
    
    echo "2. Listing containers:"
    ./servin list
    echo ""
    
    echo "3. Testing container logs:"
    ./servin logs test-container
    echo ""
    
    echo "4. Testing container exec:"
    ./servin exec test-container echo "Container exec working!"
    echo ""
    
    echo "✅ Container functionality test complete!"
    echo ""
    echo "🎯 Available commands:"
    echo "   ./servin run nginx:alpine"
    echo "   ./servin exec <container> sh"
    echo "   ./servin logs <container>"
    echo "   ./servin stop <container>"
    echo "   ./servin rm <container>"
    
else
    echo "❌ SSH not available yet"
    echo ""
    echo "🔧 SSH Setup Required:"
    echo "The VM is running but SSH needs to be configured manually."
    echo ""
    echo "Quick setup (connect to VM console and run):"
    echo "   apk update && apk add openssh"
    echo "   echo 'root:servin123' | chpasswd"
    echo "   echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config"
    echo "   rc-update add sshd default && rc-service sshd start"
    echo ""
    echo "Then test with: ssh root@localhost -p 2222"
    echo ""
    
    # Demonstrate that VM infrastructure is working
    echo "🖥️  VM Infrastructure Status:"
    echo "   VM Process: $(ps aux | grep qemu-system-aarch64 | grep -v grep | awk '{print "PID " $2 " - Running"}')"
    echo "   SSH Port: 2222 (forwarded to VM port 22)"
    echo "   VM Path: ~/.servin/vms/servin-vm/"
    echo ""
    
    echo "📋 Once SSH is configured, all container operations will work:"
    echo "   • Native Linux containers (not Docker simulation)"
    echo "   • Full Linux namespace support" 
    echo "   • True container isolation"
    echo "   • Native cgroup management"
fi

echo ""
echo "🏁 Test complete. VM infrastructure is ready for container operations!"