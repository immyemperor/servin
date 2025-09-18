#!/bin/bash

# Quick Build Script - Build Servin and start VM
# Usage: ./quick-build.sh

set -e

echo "🚀 Quick Servin Build & VM Start"
echo "================================"

# Build
echo "📦 Building..."
go build -o servin main.go

# Start VM (with automated SSH setup)
echo "🚀 Starting VM..."
./servin vm start

echo ""
echo "✅ Quick build complete!"
echo "========================"
echo ""
echo "VM is starting with automated SSH setup..."
echo "Wait ~30-60 seconds for SSH to be available."
echo ""
echo "Test with: ssh root@localhost -p 2222"
echo "Password: servin123"
echo ""
echo "Or use the full build script: ./build-with-vm.sh"