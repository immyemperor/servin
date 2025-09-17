"""
Servin Desktop GUI - Main Application
Tkinter wrapper that embeds the web interface using pywebview
"""

import os
import sys
import threading
import time
import tkinter as tk
from tkinter import ttk, messagebox
import webbrowser
import subprocess
import webview

# Add the current directory to Python path for imports
# Handle PyInstaller bundled executable paths
if hasattr(sys, '_MEIPASS'):
    # Running from PyInstaller bundle
    current_dir = sys._MEIPASS
    print(f"🚀 Running from PyInstaller bundle: {current_dir}")
else:
    # Running from source
    current_dir = os.path.dirname(os.path.abspath(__file__))
    print(f"🐍 Running from source: {current_dir}")

if current_dir not in sys.path:
    sys.path.insert(0, current_dir)

# Debug: Print Python path and current directory for troubleshooting
print(f"🐍 Python executable: {sys.executable}")
print(f"📂 Current directory: {current_dir}")
print(f"📦 Python path: {sys.path[:3]}...")  # Show first 3 entries

# Import Flask app with error handling
try:
    from app import app
    print("✅ Successfully imported Flask app")
except ImportError as e:
    print(f"❌ Failed to import app module: {e}")
    print(f"📂 Looking for app.py in: {current_dir}")
    print(f"📁 Directory contents: {os.listdir(current_dir)}")
    # Try to help with debugging
    if hasattr(sys, '_MEIPASS'):
        print(f"🔧 PyInstaller temp dir: {sys._MEIPASS}")
        try:
            print(f"🔧 PyInstaller temp contents: {os.listdir(sys._MEIPASS)}")
        except:
            print("🔧 Could not list PyInstaller temp contents")
    raise

class ServinDesktopGUI:
    def __init__(self):
        self.flask_thread = None
        self.flask_running = False
        self.webview_window = None
        self.root = None
        
    def check_servin_installed(self):
        """Check if Servin is installed and accessible"""
        try:
            import platform
            
            # Look for servin binary in the parent directory
            parent_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
            
            # Platform-specific binary detection
            system = platform.system().lower()
            machine = platform.machine().lower()
            
            search_paths = []
            
            if system == "windows":
                search_paths.extend([
                    os.path.join(parent_dir, "build", "windows-amd64", "servin.exe"),
                    os.path.join(parent_dir, "servin.exe"),
                    os.path.join(parent_dir, "servin")
                ])
            elif system == "darwin":
                if machine in ["arm64", "aarch64"]:
                    search_paths.extend([
                        os.path.join(parent_dir, "build", "darwin-arm64", "servin"),
                        os.path.join(parent_dir, "build", "darwin-amd64", "servin"),
                    ])
                else:
                    search_paths.extend([
                        os.path.join(parent_dir, "build", "darwin-amd64", "servin"),
                        os.path.join(parent_dir, "build", "darwin-arm64", "servin"),
                    ])
                search_paths.extend([
                    os.path.join(parent_dir, "servin"),
                    os.path.join(parent_dir, "servin.exe")
                ])
            elif system == "linux":
                if machine in ["aarch64", "arm64"]:
                    search_paths.extend([
                        os.path.join(parent_dir, "build", "linux-arm64", "servin"),
                        os.path.join(parent_dir, "build", "linux-amd64", "servin"),
                    ])
                else:
                    search_paths.extend([
                        os.path.join(parent_dir, "build", "linux-amd64", "servin"),
                        os.path.join(parent_dir, "build", "linux-arm64", "servin"),
                    ])
                search_paths.extend([
                    os.path.join(parent_dir, "servin"),
                    os.path.join(parent_dir, "servin.exe")
                ])
            
            # Find the first existing and executable binary
            for servin_path in search_paths:
                if os.path.exists(servin_path) and os.access(servin_path, os.X_OK):
                    result = subprocess.run([servin_path, '--help'], 
                                          capture_output=True, text=True, timeout=10)
                    if result.returncode == 0:
                        return True
            
            return False
        except (subprocess.TimeoutExpired, FileNotFoundError):
            return False
    
    def start_flask_server(self):
        """Start the Flask server in a separate thread"""
        def run_server():
            try:
                self.flask_running = True
                # Run Flask app on localhost:5555
                app.run(host='127.0.0.1', port=5555, debug=False, use_reloader=False)
            except Exception as e:
                print(f"Flask server error: {e}")
                self.flask_running = False
        
        self.flask_thread = threading.Thread(target=run_server, daemon=True)
        self.flask_thread.start()
        
        # Wait a moment for the server to start
        time.sleep(2)
    
    def create_webview_window(self):
        """Create the webview window with the Servin GUI"""
        try:
            # Try to create the webview window
            self.webview_window = webview.create_window(
                title='Servin Desktop GUI',
                url='http://127.0.0.1:5555',
                width=1200,
                height=800,
                min_size=(900, 600),
                resizable=True,
                fullscreen=False,
                on_top=False
            )
            
            # Start the webview (this will block until the window is closed)
            webview.start(debug=False)
            
        except ImportError as e:
            print(f"Webview import error: {e}")
            print("Falling back to browser...")
            self.open_in_browser()
            self.show_fallback_ui()
        except Exception as e:
            print(f"Webview error: {e}")
            print("This might be due to missing Edge WebView2 runtime on Windows.")
            print("Falling back to browser...")
            self.open_in_browser()
            self.show_fallback_ui()
    
    def show_fallback_ui(self):
        """Show a fallback Tkinter UI if webview fails"""
        self.root = tk.Tk()
        self.root.title("Servin Desktop GUI")
        self.root.geometry("600x400")
        self.root.configure(bg='#2d2d30')
        
        # Main frame
        main_frame = ttk.Frame(self.root)
        main_frame.pack(fill=tk.BOTH, expand=True, padx=20, pady=20)
        
        # Title
        title_label = ttk.Label(main_frame, text="Servin Desktop GUI", 
                               font=('Arial', 16, 'bold'))
        title_label.pack(pady=(0, 20))
        
        # Status
        status_frame = ttk.Frame(main_frame)
        status_frame.pack(fill=tk.X, pady=(0, 20))
        
        ttk.Label(status_frame, text="Status:").pack(side=tk.LEFT)
        
        if self.flask_running:
            status_text = "Flask Server: Running"
            status_color = "green"
        else:
            status_text = "Flask Server: Not Running"
            status_color = "red"
        
        status_label = ttk.Label(status_frame, text=status_text, foreground=status_color)
        status_label.pack(side=tk.LEFT, padx=(5, 0))
        
        # Buttons
        button_frame = ttk.Frame(main_frame)
        button_frame.pack(fill=tk.X, pady=(0, 20))
        
        if self.flask_running:
            open_browser_btn = ttk.Button(button_frame, text="Open in Browser", 
                                         command=self.open_in_browser)
            open_browser_btn.pack(side=tk.LEFT, padx=(0, 10))
        
        refresh_btn = ttk.Button(button_frame, text="Refresh", 
                                command=self.refresh_status)
        refresh_btn.pack(side=tk.LEFT, padx=(0, 10))
        
        exit_btn = ttk.Button(button_frame, text="Exit", 
                             command=self.root.quit)
        exit_btn.pack(side=tk.RIGHT)
        
        # Instructions
        instructions = (
            "Servin Desktop GUI is running.\n"
            "\n"
            "Features:\n"
            "• View and manage Servin containers\n"
            "• Browse and import container images\n"
            "• Manage Servin volumes\n"
            "• Real-time status updates\n"
            "\n"
            "If the embedded browser doesn't work,\n"
            "click 'Open in Browser' to use your default browser."
        )
        
        instructions_label = ttk.Label(main_frame, text=instructions, 
                                      justify=tk.LEFT, wraplength=500)
        instructions_label.pack(pady=(0, 20))
        
        # Servin status
        servin_frame = ttk.LabelFrame(main_frame, text="Servin Status")
        servin_frame.pack(fill=tk.X, pady=(0, 20))
        
        if self.check_servin_installed():
            servin_status = "✓ Servin is installed and accessible"
            servin_color = "green"
        else:
            servin_status = "✗ Servin not found or not accessible"
            servin_color = "red"
        
        servin_label = ttk.Label(servin_frame, text=servin_status, 
                                 foreground=servin_color)
        servin_label.pack(pady=10)
        
        self.root.mainloop()
    
    def open_in_browser(self):
        """Open the GUI in the default web browser"""
        webbrowser.open('http://127.0.0.1:5555')
    
    def refresh_status(self):
        """Refresh the application status"""
        if self.root:
            self.root.destroy()
        self.show_fallback_ui()
    
    def run(self):
        """Main entry point for the application"""
        print("Starting Servin Desktop GUI...")
        
        # Check if Servin is available
        if not self.check_servin_installed():
            print("Warning: Servin not found or not accessible")
            print("Make sure Servin binary is available in the parent directory")
        
        # Start Flask server
        print("Starting Flask server...")
        self.start_flask_server()
        
        if not self.flask_running:
            messagebox.showerror("Error", "Failed to start Flask server")
            return
        
        print("Flask server started successfully")
        print("Starting GUI...")
        
        # Try to create webview window first
        try:
            self.create_webview_window()
        except Exception as e:
            print(f"Webview failed: {e}")
            print("Falling back to Tkinter UI")
            self.show_fallback_ui()

def main():
    """Main function to run the Servin Desktop GUI"""
    # Add the parent directory to the Python path
    parent_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    if parent_dir not in sys.path:
        sys.path.insert(0, parent_dir)
    
    try:
        app = ServinDesktopGUI()
        app.run()
    except KeyboardInterrupt:
        print("\nApplication interrupted by user")
    except Exception as e:
        print(f"Application error: {e}")
        messagebox.showerror("Error", f"Application error: {e}")

if __name__ == "__main__":
    main()
