"""
Demo/Test launcher for Servin Desktop GUI
This script will start the Flask server and open the GUI in browser for testing
"""

import os
import sys
import time
import threading
import webbrowser

# Add current directory to Python path
current_dir = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, current_dir)

from app import app

def main():
    print("Starting Servin Desktop GUI Demo...")
    print("This will start the Flask server and open the GUI in your browser.")
    print("Press Ctrl+C to stop the server.\n")
    
    # Start Flask server in a separate thread
    def run_server():
        app.run(host='127.0.0.1', port=5555, debug=False, use_reloader=False)
    
    server_thread = threading.Thread(target=run_server, daemon=True)
    server_thread.start()
    
    # Wait for server to start
    print("Starting Flask server...")
    time.sleep(3)
    
    # Open browser
    url = 'http://127.0.0.1:5555'
    print(f"Opening {url} in your default browser...")
    webbrowser.open(url)
    
    print("\nServin Desktop GUI is now running!")
    print("You can:")
    print("- View containers, images, and volumes")
    print("- Manage Servin resources through the web interface")
    print("- Use the browser interface if pywebview doesn't work")
    print("\nPress Ctrl+C to stop the server...")
    
    try:
        # Keep the main thread alive
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        print("\nShutting down server...")

if __name__ == "__main__":
    main()
