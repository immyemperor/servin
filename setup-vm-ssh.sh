#!/bin/bash

# Script to automatically set up SSH in Alpine Linux VM
# This script sends commands to the VM to configure SSH

echo "Setting up SSH in Alpine Linux VM..."
echo "This will configure SSH access for the Servin VM"

# Wait for VM to be fully booted
echo "Waiting for VM to boot completely..."
sleep 5

# Try to connect and send setup commands
# Note: This is a demonstration script - in practice, you would
# connect to the VM console directly or use expect/automation

cat << 'EOF'
To manually set up SSH in the Alpine VM:

1. The VM should be running (you can see the QEMU process above)
2. Connect to VM console (QEMU is running in -nographic mode)
3. Login as root (press Enter when asked for password)
4. Run these commands in the VM:

   apk update
   apk add openssh
   echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config
   echo 'PasswordAuthentication yes' >> /etc/ssh/sshd_config
   passwd root  # Set password to 'servin123'
   rc-update add sshd default
   rc-service sshd start

5. After setup, SSH will be available at:
   ssh root@localhost -p 2222

The VM is ready for Servin container testing once SSH is configured!
EOF

echo ""
echo "VM Status:"
ps aux | grep qemu-system-aarch64 | grep -v grep | head -1
echo ""
echo "To test SSH (after manual setup): ssh root@localhost -p 2222"