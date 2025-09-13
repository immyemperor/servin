#!/usr/bin/env python3
"""
Simple HTTP server to serve the Servin Container Runtime wiki
"""

import http.server
import socketserver
import webbrowser
import os
import sys
from pathlib import Path

def main():
    # Configuration
    PORT = 8081
    WIKI_DIR = Path(__file__).parent
    WIKI_FILE = "wiki.html"
    
    # Change to wiki directory
    os.chdir(WIKI_DIR)
    
    # Check if wiki file exists
    if not (WIKI_DIR / WIKI_FILE).exists():
        print(f"❌ Error: {WIKI_FILE} not found in {WIKI_DIR}")
        print(f"   Please run this script from the docs/ directory")
        sys.exit(1)
    
    # Create HTTP server
    handler = http.server.SimpleHTTPRequestHandler
    
    try:
        with socketserver.TCPServer(("", PORT), handler) as httpd:
            print("🚀 Servin Container Runtime Wiki Server")
            print("=" * 50)
            print(f"📚 Serving wiki at: http://localhost:{PORT}/{WIKI_FILE}")
            print(f"📂 Directory: {WIKI_DIR}")
            print(f"🌐 Opening in browser...")
            print("=" * 50)
            print("Press Ctrl+C to stop the server")
            
            # Open in browser
            webbrowser.open(f"http://localhost:{PORT}/{WIKI_FILE}")
            
            # Start serving
            httpd.serve_forever()
            
    except KeyboardInterrupt:
        print("\n\n👋 Wiki server stopped. Thanks for using Servin!")
    except OSError as e:
        if e.errno == 48:  # Address already in use
            print(f"❌ Error: Port {PORT} is already in use")
            print(f"   Try a different port or stop the process using port {PORT}")
        else:
            print(f"❌ Error starting server: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
