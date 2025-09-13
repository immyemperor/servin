#!/bin/bash

# Servin Test Script
# This script demonstrates basic Servin functionality
# Run with: sudo ./test_containers.sh

set -e

SERVIN="./servin"
TEST_IMAGE="testimage"

echo "=== Servin Test Suite ==="

# Check if servin binary exists
if [ ! -f "$SERVIN" ]; then
    echo "Error: servin binary not found. Please build first with 'go build .'"
    exit 1
fi

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "Error: This script must be run as root (use sudo)"
    exit 1
fi

echo "1. Testing basic container execution..."
$SERVIN run $TEST_IMAGE echo "Hello from container!" || echo "Test 1 failed (expected on first run)"

echo ""
echo "2. Testing container with hostname..."
$SERVIN run --hostname testhost $TEST_IMAGE hostname || echo "Test 2 failed (expected on first run)"

echo ""
echo "3. Testing container with memory limit..."
$SERVIN run --memory 64m --name memtest $TEST_IMAGE echo "Memory limited container" || echo "Test 3 failed (expected on first run)"

echo ""
echo "4. Testing container listing..."
$SERVIN ls

echo ""
echo "5. Testing help output..."
$SERVIN --help

echo ""
echo "=== Test Suite Complete ==="
echo "Note: Some tests may fail until image management is implemented"
echo "The core container runtime infrastructure is in place!"
