#!/usr/bin/env python3
"""
Servin Container Runtime - Linux Installation Wizard
A graphical installer for Linux systems using tkinter
"""

import tkinter as tk
from tkinter import ttk, messagebox, filedialog
import os
import sys
import subprocess
import pwd
import grp
import shutil
import tempfile
import tarfile
from pathlib import Path

class ServinInstaller:
    def __init__(self):
        self.root = tk.Tk()
        self.root.title("Servin Container Runtime - Installation Wizard")
        self.root.geometry("600x500")
        self.root.resizable(False, False)
        
        # Installation variables
        self.install_dir = tk.StringVar(value="/usr/local/bin")
        self.data_dir = tk.StringVar(value="/var/lib/servin")
        self.config_dir = tk.StringVar(value="/etc/servin")
        self.install_gui = tk.BooleanVar(value=True)
        self.install_service = tk.BooleanVar(value=True)
        self.create_shortcuts = tk.BooleanVar(value=True)
        
        # Check if running as root
        self.is_root = os.geteuid() == 0
        
        # Current page
        self.current_page = 0
        self.pages = []
        
        self.setup_ui()
        self.show_page(0)
    
    def setup_ui(self):
        # Main frame
        main_frame = ttk.Frame(self.root, padding="20")
        main_frame.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        
        # Header
        header_frame = ttk.Frame(main_frame)
        header_frame.grid(row=0, column=0, columnspan=2, sticky=(tk.W, tk.E), pady=(0, 20))
        
        title_label = ttk.Label(header_frame, text="Servin Container Runtime", 
                               font=("Arial", 16, "bold"))
        title_label.grid(row=0, column=0, sticky=tk.W)
        
        subtitle_label = ttk.Label(header_frame, text="Installation Wizard", 
                                  font=("Arial", 10))
        subtitle_label.grid(row=1, column=0, sticky=tk.W)
        
        # Content frame
        self.content_frame = ttk.Frame(main_frame)
        self.content_frame.grid(row=1, column=0, columnspan=2, sticky=(tk.W, tk.E, tk.N, tk.S))
        
        # Button frame
        button_frame = ttk.Frame(main_frame)
        button_frame.grid(row=2, column=0, columnspan=2, sticky=(tk.W, tk.E), pady=(20, 0))
        
        self.back_button = ttk.Button(button_frame, text="< Back", command=self.go_back)
        self.back_button.grid(row=0, column=0, padx=(0, 10))
        
        self.next_button = ttk.Button(button_frame, text="Next >", command=self.go_next)
        self.next_button.grid(row=0, column=1, padx=(0, 10))
        
        self.cancel_button = ttk.Button(button_frame, text="Cancel", command=self.cancel_install)
        self.cancel_button.grid(row=0, column=2)
        
        # Setup pages
        self.setup_pages()
    
    def setup_pages(self):
        self.pages = [
            self.create_welcome_page,
            self.create_license_page,
            self.create_options_page,
            self.create_directories_page,
            self.create_confirmation_page,
            self.create_installation_page,
            self.create_finish_page
        ]
    
    def create_welcome_page(self):
        page = ttk.Frame(self.content_frame)
        
        welcome_text = """Welcome to the Servin Container Runtime Installation Wizard

Servin is a lightweight, Docker-compatible container runtime with a modern GUI interface.

This wizard will guide you through the installation process and configure:
• Core container runtime
• Desktop GUI application (optional)
• System service integration (optional)
• Command-line tools

Click Next to continue."""
        
        label = ttk.Label(page, text=welcome_text, font=("Arial", 10), justify=tk.LEFT)
        label.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N), pady=20)
        
        # System requirements
        req_frame = ttk.LabelFrame(page, text="System Requirements", padding="10")
        req_frame.grid(row=1, column=0, sticky=(tk.W, tk.E), pady=20)
        
        requirements = [
            "• Linux kernel 3.10+ (64-bit)",
            "• 512MB RAM (2GB+ recommended)",
            "• 1GB free disk space",
            "• Root privileges (sudo access)"
        ]
        
        for i, req in enumerate(requirements):
            ttk.Label(req_frame, text=req).grid(row=i, column=0, sticky=tk.W)
        
        # Root check warning
        if not self.is_root:
            warning_frame = ttk.Frame(page)
            warning_frame.grid(row=2, column=0, sticky=(tk.W, tk.E), pady=10)
            
            warning_label = ttk.Label(warning_frame, 
                                    text="⚠ Warning: Root privileges required for installation", 
                                    foreground="red", font=("Arial", 10, "bold"))
            warning_label.grid(row=0, column=0, sticky=tk.W)
            
            note_label = ttk.Label(warning_frame, 
                                 text="Please run: sudo python3 servin-installer.py", 
                                 font=("Arial", 9))
            note_label.grid(row=1, column=0, sticky=tk.W)
        
        return page
    
    def create_license_page(self):
        page = ttk.Frame(self.content_frame)
        
        ttk.Label(page, text="License Agreement", font=("Arial", 12, "bold")).grid(row=0, column=0, sticky=tk.W, pady=(0, 10))
        
        # License text area
        text_frame = ttk.Frame(page)
        text_frame.grid(row=1, column=0, sticky=(tk.W, tk.E, tk.N, tk.S), pady=(0, 10))
        
        text_widget = tk.Text(text_frame, height=15, width=70, wrap=tk.WORD)
        scrollbar = ttk.Scrollbar(text_frame, orient=tk.VERTICAL, command=text_widget.yview)
        text_widget.configure(yscrollcommand=scrollbar.set)
        
        license_text = """Apache License 2.0

Copyright 2025 Servin Project

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

This software is provided "AS IS" without warranty of any kind."""
        
        text_widget.insert(tk.END, license_text)
        text_widget.configure(state=tk.DISABLED)
        
        text_widget.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        scrollbar.grid(row=0, column=1, sticky=(tk.N, tk.S))
        
        # Accept checkbox
        self.accept_license = tk.BooleanVar()
        accept_check = ttk.Checkbutton(page, text="I accept the terms of the License Agreement", 
                                      variable=self.accept_license, command=self.update_next_button)
        accept_check.grid(row=2, column=0, sticky=tk.W, pady=10)
        
        return page
    
    def create_options_page(self):
        page = ttk.Frame(self.content_frame)
        
        ttk.Label(page, text="Installation Options", font=("Arial", 12, "bold")).grid(row=0, column=0, sticky=tk.W, pady=(0, 20))
        
        options_frame = ttk.LabelFrame(page, text="Components to Install", padding="10")
        options_frame.grid(row=1, column=0, sticky=(tk.W, tk.E), pady=(0, 20))
        
        ttk.Checkbutton(options_frame, text="Core Runtime (required)", 
                       state=tk.DISABLED, variable=tk.BooleanVar(value=True)).grid(row=0, column=0, sticky=tk.W, pady=2)
        
        ttk.Checkbutton(options_frame, text="Desktop GUI Application", 
                       variable=self.install_gui).grid(row=1, column=0, sticky=tk.W, pady=2)
        
        ttk.Checkbutton(options_frame, text="System Service (systemd/SysV)", 
                       variable=self.install_service).grid(row=2, column=0, sticky=tk.W, pady=2)
        
        ttk.Checkbutton(options_frame, text="Desktop Shortcuts", 
                       variable=self.create_shortcuts).grid(row=3, column=0, sticky=tk.W, pady=2)
        
        # Advanced options
        advanced_frame = ttk.LabelFrame(page, text="Advanced Options", padding="10")
        advanced_frame.grid(row=2, column=0, sticky=(tk.W, tk.E))
        
        self.create_user = tk.BooleanVar(value=True)
        ttk.Checkbutton(advanced_frame, text="Create 'servin' system user", 
                       variable=self.create_user).grid(row=0, column=0, sticky=tk.W, pady=2)
        
        self.add_to_path = tk.BooleanVar(value=True)
        ttk.Checkbutton(advanced_frame, text="Add to system PATH", 
                       variable=self.add_to_path).grid(row=1, column=0, sticky=tk.W, pady=2)
        
        return page
    
    def create_directories_page(self):
        page = ttk.Frame(self.content_frame)
        
        ttk.Label(page, text="Installation Directories", font=("Arial", 12, "bold")).grid(row=0, column=0, columnspan=3, sticky=tk.W, pady=(0, 20))
        
        # Installation directory
        ttk.Label(page, text="Installation Directory:").grid(row=1, column=0, sticky=tk.W, pady=5)
        ttk.Entry(page, textvariable=self.install_dir, width=40).grid(row=1, column=1, sticky=(tk.W, tk.E), padx=(10, 5))
        ttk.Button(page, text="Browse", command=lambda: self.browse_directory(self.install_dir)).grid(row=1, column=2)
        
        # Data directory
        ttk.Label(page, text="Data Directory:").grid(row=2, column=0, sticky=tk.W, pady=5)
        ttk.Entry(page, textvariable=self.data_dir, width=40).grid(row=2, column=1, sticky=(tk.W, tk.E), padx=(10, 5))
        ttk.Button(page, text="Browse", command=lambda: self.browse_directory(self.data_dir)).grid(row=2, column=2)
        
        # Config directory
        ttk.Label(page, text="Configuration Directory:").grid(row=3, column=0, sticky=tk.W, pady=5)
        ttk.Entry(page, textvariable=self.config_dir, width=40).grid(row=3, column=1, sticky=(tk.W, tk.E), padx=(10, 5))
        ttk.Button(page, text="Browse", command=lambda: self.browse_directory(self.config_dir)).grid(row=3, column=2)
        
        # Space requirements
        space_frame = ttk.LabelFrame(page, text="Disk Space Requirements", padding="10")
        space_frame.grid(row=4, column=0, columnspan=3, sticky=(tk.W, tk.E), pady=20)
        
        space_info = [
            "• Binaries: ~50MB",
            "• Configuration: <1MB",
            "• Logs: 10-100MB (varies by usage)",
            "• Container data: Varies by containers"
        ]
        
        for i, info in enumerate(space_info):
            ttk.Label(space_frame, text=info).grid(row=i, column=0, sticky=tk.W)
        
        return page
    
    def create_confirmation_page(self):
        page = ttk.Frame(self.content_frame)
        
        ttk.Label(page, text="Ready to Install", font=("Arial", 12, "bold")).grid(row=0, column=0, sticky=tk.W, pady=(0, 20))
        
        # Summary frame
        summary_frame = ttk.LabelFrame(page, text="Installation Summary", padding="10")
        summary_frame.grid(row=1, column=0, sticky=(tk.W, tk.E), pady=(0, 20))
        
        self.summary_text = tk.Text(summary_frame, height=15, width=70, wrap=tk.WORD)
        summary_scrollbar = ttk.Scrollbar(summary_frame, orient=tk.VERTICAL, command=self.summary_text.yview)
        self.summary_text.configure(yscrollcommand=summary_scrollbar.set)
        
        self.summary_text.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        summary_scrollbar.grid(row=0, column=1, sticky=(tk.N, tk.S))
        
        return page
    
    def create_installation_page(self):
        page = ttk.Frame(self.content_frame)
        
        ttk.Label(page, text="Installing Servin Container Runtime", font=("Arial", 12, "bold")).grid(row=0, column=0, sticky=tk.W, pady=(0, 20))
        
        # Progress bar
        self.progress = ttk.Progressbar(page, mode='determinate', length=400)
        self.progress.grid(row=1, column=0, sticky=(tk.W, tk.E), pady=(0, 10))
        
        # Status label
        self.status_label = ttk.Label(page, text="Preparing installation...")
        self.status_label.grid(row=2, column=0, sticky=tk.W, pady=(0, 20))
        
        # Log area
        log_frame = ttk.LabelFrame(page, text="Installation Log", padding="10")
        log_frame.grid(row=3, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        
        self.log_text = tk.Text(log_frame, height=10, width=70, wrap=tk.WORD)
        log_scrollbar = ttk.Scrollbar(log_frame, orient=tk.VERTICAL, command=self.log_text.yview)
        self.log_text.configure(yscrollcommand=log_scrollbar.set)
        
        self.log_text.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        log_scrollbar.grid(row=0, column=1, sticky=(tk.N, tk.S))
        
        return page
    
    def create_finish_page(self):
        page = ttk.Frame(self.content_frame)
        
        self.finish_title = ttk.Label(page, text="Installation Complete!", font=("Arial", 12, "bold"))
        self.finish_title.grid(row=0, column=0, sticky=tk.W, pady=(0, 20))
        
        self.finish_text = ttk.Label(page, text="", font=("Arial", 10), justify=tk.LEFT)
        self.finish_text.grid(row=1, column=0, sticky=(tk.W, tk.E, tk.N), pady=(0, 20))
        
        # Launch options
        launch_frame = ttk.LabelFrame(page, text="What would you like to do next?", padding="10")
        launch_frame.grid(row=2, column=0, sticky=(tk.W, tk.E))
        
        self.launch_gui = tk.BooleanVar(value=True)
        self.start_service = tk.BooleanVar(value=True)
        
        ttk.Checkbutton(launch_frame, text="Launch Servin GUI", variable=self.launch_gui).grid(row=0, column=0, sticky=tk.W, pady=2)
        ttk.Checkbutton(launch_frame, text="Start Servin service", variable=self.start_service).grid(row=1, column=0, sticky=tk.W, pady=2)
        
        return page
    
    def browse_directory(self, var):
        directory = filedialog.askdirectory(initialdir=var.get())
        if directory:
            var.set(directory)
    
    def update_next_button(self):
        if self.current_page == 1:  # License page
            self.next_button.configure(state=tk.NORMAL if self.accept_license.get() else tk.DISABLED)
    
    def show_page(self, page_num):
        # Clear content frame
        for widget in self.content_frame.winfo_children():
            widget.destroy()
        
        # Create and show new page
        page = self.pages[page_num]()
        page.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        
        # Update buttons
        self.back_button.configure(state=tk.NORMAL if page_num > 0 else tk.DISABLED)
        
        if page_num == len(self.pages) - 1:  # Last page
            self.next_button.configure(text="Finish", command=self.finish_install)
        elif page_num == len(self.pages) - 2:  # Installation page
            self.next_button.configure(text="Install", command=self.start_installation)
        else:
            self.next_button.configure(text="Next >", command=self.go_next)
        
        # Special handling for license page
        if page_num == 1:
            self.next_button.configure(state=tk.DISABLED)
        
        # Update confirmation page with summary
        if page_num == 4:  # Confirmation page
            self.update_summary()
        
        self.current_page = page_num
    
    def update_summary(self):
        summary = f"""Installation Directory: {self.install_dir.get()}
Data Directory: {self.data_dir.get()}
Configuration Directory: {self.config_dir.get()}

Components to install:
• Core Runtime: Yes
• Desktop GUI: {'Yes' if self.install_gui.get() else 'No'}
• System Service: {'Yes' if self.install_service.get() else 'No'}
• Desktop Shortcuts: {'Yes' if self.create_shortcuts.get() else 'No'}

Advanced Options:
• Create system user: {'Yes' if self.create_user.get() else 'No'}
• Add to PATH: {'Yes' if self.add_to_path.get() else 'No'}

Click 'Install' to begin the installation process."""
        
        self.summary_text.delete(1.0, tk.END)
        self.summary_text.insert(tk.END, summary)
        self.summary_text.configure(state=tk.DISABLED)
    
    def go_back(self):
        if self.current_page > 0:
            self.show_page(self.current_page - 1)
    
    def go_next(self):
        if self.current_page < len(self.pages) - 1:
            self.show_page(self.current_page + 1)
    
    def cancel_install(self):
        if messagebox.askyesno("Cancel Installation", "Are you sure you want to cancel the installation?"):
            self.root.quit()
    
    def start_installation(self):
        if not self.is_root:
            messagebox.showerror("Error", "Root privileges are required for installation.\nPlease run: sudo python3 servin-installer.py")
            return
        
        self.show_page(len(self.pages) - 2)  # Installation page
        self.next_button.configure(state=tk.DISABLED)
        self.back_button.configure(state=tk.DISABLED)
        
        # Start installation in background
        self.root.after(100, self.perform_installation)
    
    def log_message(self, message):
        self.log_text.insert(tk.END, f"{message}\n")
        self.log_text.see(tk.END)
        self.root.update()
    
    def perform_installation(self):
        try:
            self.progress['value'] = 0
            self.status_label.configure(text="Creating directories...")
            self.log_message("Starting installation...")
            
            # Create directories
            directories = [
                self.install_dir.get(),
                self.data_dir.get(),
                self.config_dir.get(),
                f"{self.data_dir.get()}/volumes",
                f"{self.data_dir.get()}/images",
                "/var/log/servin"
            ]
            
            for directory in directories:
                os.makedirs(directory, exist_ok=True)
                self.log_message(f"Created directory: {directory}")
            
            self.progress['value'] = 20
            
            # Install binaries (assuming they're in the same directory as this script)
            self.status_label.configure(text="Installing binaries...")
            script_dir = os.path.dirname(os.path.abspath(__file__))
            
            binaries = ["servin", "servin-desktop"]
            if self.install_gui.get():
                binaries.append("servin-gui")
            
            for binary in binaries:
                src = os.path.join(script_dir, binary)
                dst = os.path.join(self.install_dir.get(), binary)
                if os.path.exists(src):
                    shutil.copy2(src, dst)
                    os.chmod(dst, 0o755)
                    self.log_message(f"Installed: {binary}")
                else:
                    self.log_message(f"Warning: {binary} not found")
            
            self.progress['value'] = 40
            
            # Create configuration
            self.status_label.configure(text="Creating configuration...")
            config_content = f"""# Servin Configuration File
data_dir={self.data_dir.get()}
log_level=info
log_file=/var/log/servin/servin.log
runtime=native
bridge_name=servin0
cri_port=10250
cri_enabled=false"""
            
            with open(f"{self.config_dir.get()}/servin.conf", 'w') as f:
                f.write(config_content)
            self.log_message("Created configuration file")
            
            self.progress['value'] = 60
            
            # Create user
            if self.create_user.get():
                self.status_label.configure(text="Creating system user...")
                try:
                    subprocess.run(["useradd", "--system", "--no-create-home", "--shell", "/bin/false", "servin"], 
                                 check=False, capture_output=True)
                    self.log_message("Created system user: servin")
                except:
                    self.log_message("User 'servin' may already exist")
            
            self.progress['value'] = 80
            
            # Install service
            if self.install_service.get():
                self.status_label.configure(text="Installing service...")
                self.install_systemd_service()
            
            # Set permissions
            self.status_label.configure(text="Setting permissions...")
            if self.create_user.get():
                try:
                    uid = pwd.getpwnam("servin").pw_uid
                    gid = grp.getgrnam("servin").gr_gid
                    for path in [self.data_dir.get(), "/var/log/servin"]:
                        os.chown(path, uid, gid)
                        for root, dirs, files in os.walk(path):
                            for d in dirs:
                                os.chown(os.path.join(root, d), uid, gid)
                            for f in files:
                                os.chown(os.path.join(root, f), uid, gid)
                    self.log_message("Set directory permissions")
                except:
                    self.log_message("Warning: Could not set all permissions")
            
            self.progress['value'] = 100
            self.status_label.configure(text="Installation complete!")
            self.log_message("Installation completed successfully!")
            
            # Enable next button to go to finish page
            self.next_button.configure(state=tk.NORMAL, text="Next >", command=self.go_next)
            
        except Exception as e:
            self.log_message(f"Error: {str(e)}")
            messagebox.showerror("Installation Error", f"Installation failed: {str(e)}")
    
    def install_systemd_service(self):
        service_content = f"""[Unit]
Description=Servin Container Runtime
After=network.target

[Service]
Type=simple
User=servin
Group=servin
ExecStart={self.install_dir.get()}/servin daemon --config {self.config_dir.get()}/servin.conf
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target"""
        
        with open("/etc/systemd/system/servin.service", 'w') as f:
            f.write(service_content)
        
        subprocess.run(["systemctl", "daemon-reload"], check=False)
        subprocess.run(["systemctl", "enable", "servin"], check=False)
        self.log_message("Installed systemd service")
    
    def finish_install(self):
        if self.start_service.get() and self.install_service.get():
            try:
                subprocess.run(["systemctl", "start", "servin"], check=True)
                messagebox.showinfo("Service Started", "Servin service has been started successfully!")
            except:
                messagebox.showwarning("Service", "Could not start Servin service. You can start it manually with:\nsudo systemctl start servin")
        
        if self.launch_gui.get() and self.install_gui.get():
            try:
                gui_path = os.path.join(self.install_dir.get(), "servin-gui")
                subprocess.Popen([gui_path], start_new_session=True)
            except:
                messagebox.showwarning("GUI Launch", "Could not launch Servin GUI automatically.")
        
        messagebox.showinfo("Installation Complete", 
                           "Servin Container Runtime has been installed successfully!\n\n"
                           "You can now use 'servin' command from the terminal or launch the GUI from your applications menu.")
        self.root.quit()
    
    def run(self):
        self.root.mainloop()


if __name__ == "__main__":
    # Check if tkinter is available
    try:
        import tkinter
    except ImportError:
        print("Error: tkinter is not available. Please install it:")
        print("Ubuntu/Debian: sudo apt-get install python3-tk")
        print("CentOS/RHEL: sudo yum install tkinter")
        print("Or use the command-line installer: sudo ./install.sh")
        sys.exit(1)
    
    installer = ServinInstaller()
    installer.run()
