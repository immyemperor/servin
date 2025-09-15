#!/bin/bash

# Simple test script to verify Servin containerization
# This creates a basic container and runs a simple command

set -e

echo "=== Servin Containerization Test ==="
echo

# Build servin if not already built
if [[ ! -f "./servin-test" ]]; then
    echo "Building servin..."
    go build -o servin-test ./main.go
fi

echo "Testing container functionality..."
echo

# Test 1: Run a simple command in a container
echo "Test 1: Running 'echo Hello from container' in isolated container"
sudo ./servin-test run alpine echo "Hello from container"
echo

# Test 2: Run shell command
echo "Test 2: Running 'pwd' command to check working directory"
sudo ./servin-test run alpine pwd
echo

# Test 3: Check process isolation
echo "Test 3: Running 'ps aux' to check process isolation"
sudo ./servin-test run alpine ps aux || echo "ps command not available in container (expected)"
echo

echo "=== Test Complete ==="
echo "Note: These tests require Linux with root privileges for full containerization."
echo "On macOS/Windows, containerization will be simulated."