#!/bin/bash
# Shell script to serve the Servin wiki on Linux/macOS

echo "🚀 Servin Container Runtime Wiki Server"
echo "=============================================="

# Check if Python is installed
if ! command -v python3 &> /dev/null; then
    echo "❌ Error: Python 3 is not installed or not in PATH"
    echo "   Please install Python 3 and try again"
    exit 1
fi

# Check if wiki file exists
if [[ ! -f "wiki.html" ]]; then
    echo "❌ Error: wiki.html not found"
    echo "   Please run this script from the docs/ directory"
    exit 1
fi

echo "📚 Starting wiki server..."
echo "🌐 Wiki will open in your browser automatically"
echo "📝 Press Ctrl+C to stop the server"
echo "=============================================="

# Start the Python server
python3 serve-wiki.py
