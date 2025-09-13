#!/usr/bin/env python3
"""
Servin Container Runtime - macOS Installation Wizard
A native macOS installer using tkinter with Cocoa styling
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
import plistlib
from pathlib import Path

class ServinMacInstaller:
    def __init__(self):
        self.root = tk.Tk()
        self.root.title("Servin Container Runtime Installer")
        self.root.geometry("700x550")
        self.root.resizable(False, False)
        
        # macOS styling
        try:
            self.root.tk.call('tk', 'scaling', 1.4)  # Better for Retina displays
        except:
            pass
        
        # Installation variables
        self.install_dir = tk.StringVar(value="/usr/local/bin")
        self.data_dir = tk.StringVar(value="/usr/local/var/lib/servin")
        self.config_dir = tk.StringVar(value="/usr/local/etc/servin")
        self.install_gui = tk.BooleanVar(value=True)
        self.install_service = tk.BooleanVar(value=True)
        self.create_app_bundle = tk.BooleanVar(value=True)
        
        # Check if running as root
        self.is_root = os.geteuid() == 0
        
        # Current page
        self.current_page = 0
        self.pages = []
        
        self.setup_ui()
        self.show_page(0)
    
    def setup_ui(self):
        # Configure style for macOS
        style = ttk.Style()
        try:
            style.theme_use('aqua')  # macOS native theme
        except:
            pass
        
        # Main frame with padding
        main_frame = ttk.Frame(self.root, padding="30")
        main_frame.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        
        self.root.columnconfigure(0, weight=1)
        self.root.rowconfigure(0, weight=1)
        main_frame.columnconfigure(0, weight=1)
        main_frame.rowconfigure(1, weight=1)
        
        # Header with macOS-style appearance
        header_frame = ttk.Frame(main_frame)
        header_frame.grid(row=0, column=0, sticky=(tk.W, tk.E), pady=(0, 30))
        
        # App icon placeholder (you'd include an actual icon here)
        icon_label = ttk.Label(header_frame, text="📦", font=("SF Pro Display", 24))
        icon_label.grid(row=0, column=0, rowspan=2, padx=(0, 15))
        
        title_label = ttk.Label(header_frame, text="Servin Container Runtime", 
                               font=("SF Pro Display", 18, "bold"))
        title_label.grid(row=0, column=1, sticky=tk.W)
        
        subtitle_label = ttk.Label(header_frame, text="Installation Assistant", 
                                  font=("SF Pro Display", 12))
        subtitle_label.grid(row=1, column=1, sticky=tk.W)
        
        # Content frame
        self.content_frame = ttk.Frame(main_frame)
        self.content_frame.grid(row=1, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        self.content_frame.columnconfigure(0, weight=1)
        self.content_frame.rowconfigure(0, weight=1)
        
        # Button frame with macOS-style spacing
        button_frame = ttk.Frame(main_frame)
        button_frame.grid(row=2, column=0, sticky=(tk.W, tk.E), pady=(30, 0))
        
        # macOS-style button layout (right-aligned)
        button_container = ttk.Frame(button_frame)
        button_container.grid(row=0, column=1, sticky=tk.E)
        
        self.cancel_button = ttk.Button(button_container, text="Cancel", command=self.cancel_install)
        self.cancel_button.grid(row=0, column=0, padx=(0, 10))
        
        self.back_button = ttk.Button(button_container, text="Go Back", command=self.go_back)
        self.back_button.grid(row=0, column=1, padx=(0, 10))
        
        self.next_button = ttk.Button(button_container, text="Continue", command=self.go_next)
        self.next_button.grid(row=0, column=2)
        
        button_frame.columnconfigure(1, weight=1)
        
        # Setup pages
        self.setup_pages()
    
    def setup_pages(self):
        self.pages = [
            self.create_welcome_page,
            self.create_license_page,
            self.create_destination_page,
            self.create_installation_type_page,
            self.create_summary_page,
            self.create_installation_page,
            self.create_success_page
        ]
    
    def create_welcome_page(self):
        page = ttk.Frame(self.content_frame, padding="20")
        
        # Welcome message with macOS styling
        welcome_frame = ttk.Frame(page)
        welcome_frame.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N), pady=(0, 30))
        
        title = ttk.Label(welcome_frame, text="Welcome to the Servin Container Runtime Installer", 
                         font=("SF Pro Display", 16, "bold"))
        title.grid(row=0, column=0, sticky=tk.W, pady=(0, 15))
        
        description = """This installer will guide you through the installation of Servin Container Runtime, a lightweight Docker-compatible container runtime with an intuitive GUI interface.

Servin provides:
• Docker-compatible container management
• Native macOS GUI application
• Background service integration with launchd
• Command-line tools for automation
• Kubernetes CRI compatibility

The installation requires administrator privileges and approximately 100MB of disk space."""
        
        desc_label = ttk.Label(welcome_frame, text=description, font=("SF Pro Text", 11), 
                              justify=tk.LEFT, wraplength=600)
        desc_label.grid(row=1, column=0, sticky=(tk.W, tk.E))
        
        # System info
        system_frame = ttk.LabelFrame(page, text="System Requirements", padding="15")
        system_frame.grid(row=1, column=0, sticky=(tk.W, tk.E), pady=(0, 20))
        
        # Get macOS version
        try:
            version_output = subprocess.check_output(['sw_vers', '-productVersion'], text=True).strip()
            build_output = subprocess.check_output(['sw_vers', '-buildVersion'], text=True).strip()
            system_info = f"macOS {version_output} (Build {build_output})"
        except:
            system_info = "macOS (version detection failed)"
        
        ttk.Label(system_frame, text=f"Current System: {system_info}").grid(row=0, column=0, sticky=tk.W)
        ttk.Label(system_frame, text="Required: macOS 10.12 (Sierra) or later").grid(row=1, column=0, sticky=tk.W)
        ttk.Label(system_frame, text="Architecture: 64-bit Intel or Apple Silicon").grid(row=2, column=0, sticky=tk.W)
        
        if not self.is_root:
            warning_frame = ttk.Frame(page)
            warning_frame.grid(row=2, column=0, sticky=(tk.W, tk.E))
            
            warning_label = ttk.Label(warning_frame, 
                                    text="⚠️ Administrator privileges required", 
                                    foreground="#ff6b35", font=("SF Pro Text", 12, "bold"))
            warning_label.grid(row=0, column=0, sticky=tk.W)
            
            note_label = ttk.Label(warning_frame, 
                                 text="Please run: sudo python3 servin-installer.py", 
                                 font=("SF Mono", 10))
            note_label.grid(row=1, column=0, sticky=tk.W, pady=(5, 0))
        
        return page
    
    def create_license_page(self):
        page = ttk.Frame(self.content_frame, padding="20")
        
        ttk.Label(page, text="Software License Agreement", 
                 font=("SF Pro Display", 14, "bold")).grid(row=0, column=0, sticky=tk.W, pady=(0, 15))
        
        license_frame = ttk.Frame(page)
        license_frame.grid(row=1, column=0, sticky=(tk.W, tk.E, tk.N, tk.S), pady=(0, 15))
        license_frame.columnconfigure(0, weight=1)
        license_frame.rowconfigure(0, weight=1)
        
        # License text with native scrolling
        text_widget = tk.Text(license_frame, height=18, wrap=tk.WORD, 
                             font=("SF Mono", 10), borderwidth=1, relief=tk.SOLID)
        scrollbar = ttk.Scrollbar(license_frame, orient=tk.VERTICAL, command=text_widget.yview)
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

DISCLAIMER: This software is provided "AS IS" without warranty of any kind, 
express or implied, including but not limited to the warranties of 
merchantability, fitness for a particular purpose and noninfringement."""
        
        text_widget.insert(tk.END, license_text)
        text_widget.configure(state=tk.DISABLED)
        
        text_widget.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        scrollbar.grid(row=0, column=1, sticky=(tk.N, tk.S))
        
        # Agreement checkbox
        agreement_frame = ttk.Frame(page)
        agreement_frame.grid(row=2, column=0, sticky=tk.W, pady=(15, 0))
        
        self.accept_license = tk.BooleanVar()
        accept_check = ttk.Checkbutton(agreement_frame, 
                                      text="I have read and agree to the terms of the software license agreement", 
                                      variable=self.accept_license, command=self.update_continue_button)
        accept_check.grid(row=0, column=0, sticky=tk.W)
        
        return page
    
    def create_destination_page(self):
        page = ttk.Frame(self.content_frame, padding="20")
        
        ttk.Label(page, text="Select Installation Destination", 
                 font=("SF Pro Display", 14, "bold")).grid(row=0, column=0, sticky=tk.W, pady=(0, 20))
        
        # Destination selection
        dest_frame = ttk.LabelFrame(page, text="Installation Location", padding="15")
        dest_frame.grid(row=1, column=0, sticky=(tk.W, tk.E), pady=(0, 20))
        dest_frame.columnconfigure(1, weight=1)
        
        ttk.Label(dest_frame, text="Install to:").grid(row=0, column=0, sticky=tk.W, pady=5)
        install_entry = ttk.Entry(dest_frame, textvariable=self.install_dir, width=50)
        install_entry.grid(row=0, column=1, sticky=(tk.W, tk.E), padx=(10, 5))
        ttk.Button(dest_frame, text="Choose...", command=lambda: self.browse_directory(self.install_dir)).grid(row=0, column=2)
        
        ttk.Label(dest_frame, text="Data storage:").grid(row=1, column=0, sticky=tk.W, pady=5)
        data_entry = ttk.Entry(dest_frame, textvariable=self.data_dir, width=50)
        data_entry.grid(row=1, column=1, sticky=(tk.W, tk.E), padx=(10, 5))
        ttk.Button(dest_frame, text="Choose...", command=lambda: self.browse_directory(self.data_dir)).grid(row=1, column=2)
        
        ttk.Label(dest_frame, text="Configuration:").grid(row=2, column=0, sticky=tk.W, pady=5)
        config_entry = ttk.Entry(dest_frame, textvariable=self.config_dir, width=50)
        config_entry.grid(row=2, column=1, sticky=(tk.W, tk.E), padx=(10, 5))
        ttk.Button(dest_frame, text="Choose...", command=lambda: self.browse_directory(self.config_dir)).grid(row=2, column=2)
        
        # Space calculation
        space_frame = ttk.LabelFrame(page, text="Disk Space", padding="15")
        space_frame.grid(row=2, column=0, sticky=(tk.W, tk.E))
        
        try:
            statvfs = os.statvfs(self.install_dir.get().split('/')[1] if self.install_dir.get().startswith('/') else '/')
            free_space = statvfs.f_bavail * statvfs.f_frsize / (1024 * 1024 * 1024)  # GB
            space_text = f"Available space: {free_space:.1f} GB"
        except:
            space_text = "Available space: Unable to calculate"
        
        ttk.Label(space_frame, text="Space required: ~100 MB").grid(row=0, column=0, sticky=tk.W)
        ttk.Label(space_frame, text=space_text).grid(row=1, column=0, sticky=tk.W)
        
        return page
    
    def create_installation_type_page(self):
        page = ttk.Frame(self.content_frame, padding="20")
        
        ttk.Label(page, text="Installation Type", 
                 font=("SF Pro Display", 14, "bold")).grid(row=0, column=0, sticky=tk.W, pady=(0, 20))
        
        # Standard installation
        standard_frame = ttk.LabelFrame(page, text="Components", padding="15")
        standard_frame.grid(row=1, column=0, sticky=(tk.W, tk.E), pady=(0, 20))
        
        ttk.Checkbutton(standard_frame, text="Core Runtime (required)", 
                       state=tk.DISABLED, variable=tk.BooleanVar(value=True)).grid(row=0, column=0, sticky=tk.W, pady=3)
        
        ttk.Checkbutton(standard_frame, text="Desktop GUI Application", 
                       variable=self.install_gui).grid(row=1, column=0, sticky=tk.W, pady=3)
        
        ttk.Checkbutton(standard_frame, text="Background Service (launchd)", 
                       variable=self.install_service).grid(row=2, column=0, sticky=tk.W, pady=3)
        
        ttk.Checkbutton(standard_frame, text="Create Application Bundle", 
                       variable=self.create_app_bundle).grid(row=3, column=0, sticky=tk.W, pady=3)
        
        # Advanced options
        advanced_frame = ttk.LabelFrame(page, text="Advanced Options", padding="15")
        advanced_frame.grid(row=2, column=0, sticky=(tk.W, tk.E))
        
        self.create_user = tk.BooleanVar(value=True)
        ttk.Checkbutton(advanced_frame, text="Create '_servin' system user", 
                       variable=self.create_user).grid(row=0, column=0, sticky=tk.W, pady=3)
        
        self.add_to_path = tk.BooleanVar(value=True)
        ttk.Checkbutton(advanced_frame, text="Add to PATH environment", 
                       variable=self.add_to_path).grid(row=1, column=0, sticky=tk.W, pady=3)
        
        self.install_cli_tools = tk.BooleanVar(value=True)
        ttk.Checkbutton(advanced_frame, text="Install command-line tools", 
                       variable=self.install_cli_tools).grid(row=2, column=0, sticky=tk.W, pady=3)
        
        return page
    
    def create_summary_page(self):
        page = ttk.Frame(self.content_frame, padding="20")
        
        ttk.Label(page, text="Installation Summary", 
                 font=("SF Pro Display", 14, "bold")).grid(row=0, column=0, sticky=tk.W, pady=(0, 15))
        
        # Summary text area
        summary_frame = ttk.Frame(page)
        summary_frame.grid(row=1, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        summary_frame.columnconfigure(0, weight=1)
        summary_frame.rowconfigure(0, weight=1)
        
        self.summary_text = tk.Text(summary_frame, height=18, wrap=tk.WORD, 
                                   font=("SF Mono", 10), borderwidth=1, relief=tk.SOLID,
                                   state=tk.DISABLED)
        summary_scrollbar = ttk.Scrollbar(summary_frame, orient=tk.VERTICAL, command=self.summary_text.yview)
        self.summary_text.configure(yscrollcommand=summary_scrollbar.set)
        
        self.summary_text.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        summary_scrollbar.grid(row=0, column=1, sticky=(tk.N, tk.S))
        
        return page
    
    def create_installation_page(self):
        page = ttk.Frame(self.content_frame, padding="20")
        
        ttk.Label(page, text="Installing Servin Container Runtime", 
                 font=("SF Pro Display", 14, "bold")).grid(row=0, column=0, sticky=tk.W, pady=(0, 20))
        
        # Progress indicator
        progress_frame = ttk.Frame(page)
        progress_frame.grid(row=1, column=0, sticky=(tk.W, tk.E), pady=(0, 15))
        progress_frame.columnconfigure(0, weight=1)
        
        self.progress = ttk.Progressbar(progress_frame, mode='indeterminate', length=500)
        self.progress.grid(row=0, column=0, sticky=(tk.W, tk.E))
        
        self.status_label = ttk.Label(progress_frame, text="Preparing installation...", 
                                     font=("SF Pro Text", 11))
        self.status_label.grid(row=1, column=0, sticky=tk.W, pady=(10, 0))
        
        # Installation log
        log_frame = ttk.LabelFrame(page, text="Installation Log", padding="10")
        log_frame.grid(row=2, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        log_frame.columnconfigure(0, weight=1)
        log_frame.rowconfigure(0, weight=1)
        
        self.log_text = tk.Text(log_frame, height=12, wrap=tk.WORD, 
                               font=("SF Mono", 9), borderwidth=1, relief=tk.SOLID)
        log_scrollbar = ttk.Scrollbar(log_frame, orient=tk.VERTICAL, command=self.log_text.yview)
        self.log_text.configure(yscrollcommand=log_scrollbar.set)
        
        self.log_text.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        log_scrollbar.grid(row=0, column=1, sticky=(tk.N, tk.S))
        
        return page
    
    def create_success_page(self):
        page = ttk.Frame(self.content_frame, padding="20")
        
        # Success icon and message
        success_frame = ttk.Frame(page)
        success_frame.grid(row=0, column=0, sticky=(tk.W, tk.E), pady=(20, 30))
        
        success_icon = ttk.Label(success_frame, text="✅", font=("SF Pro Display", 36))
        success_icon.grid(row=0, column=0, padx=(0, 15))
        
        success_text = ttk.Label(success_frame, text="Installation Completed Successfully!", 
                                font=("SF Pro Display", 16, "bold"))
        success_text.grid(row=0, column=1, sticky=tk.W)
        
        # Summary of what was installed
        installed_frame = ttk.LabelFrame(page, text="Installed Components", padding="15")
        installed_frame.grid(row=1, column=0, sticky=(tk.W, tk.E), pady=(0, 20))
        
        self.installed_summary = ttk.Label(installed_frame, text="", font=("SF Pro Text", 11), 
                                          justify=tk.LEFT)
        self.installed_summary.grid(row=0, column=0, sticky=tk.W)
        
        # Next steps
        next_frame = ttk.LabelFrame(page, text="What's Next?", padding="15")
        next_frame.grid(row=2, column=0, sticky=(tk.W, tk.E))
        
        self.launch_gui = tk.BooleanVar(value=True)
        self.start_service = tk.BooleanVar(value=True)
        self.open_applications = tk.BooleanVar(value=False)
        
        ttk.Checkbutton(next_frame, text="Launch Servin GUI", 
                       variable=self.launch_gui).grid(row=0, column=0, sticky=tk.W, pady=2)
        ttk.Checkbutton(next_frame, text="Start background service", 
                       variable=self.start_service).grid(row=1, column=0, sticky=tk.W, pady=2)
        ttk.Checkbutton(next_frame, text="Open Applications folder", 
                       variable=self.open_applications).grid(row=2, column=0, sticky=tk.W, pady=2)
        
        return page
    
    def browse_directory(self, var):
        directory = filedialog.askdirectory(initialdir=var.get())
        if directory:
            var.set(directory)
    
    def update_continue_button(self):
        if self.current_page == 1:  # License page
            self.next_button.configure(state=tk.NORMAL if self.accept_license.get() else tk.DISABLED)
    
    def show_page(self, page_num):
        # Clear content frame
        for widget in self.content_frame.winfo_children():
            widget.destroy()
        
        # Create and show new page
        page = self.pages[page_num]()
        page.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        
        # Update buttons based on page
        self.back_button.configure(state=tk.NORMAL if page_num > 0 else tk.DISABLED)
        
        if page_num == 0:  # Welcome
            self.next_button.configure(text="Continue")
        elif page_num == len(self.pages) - 1:  # Success page
            self.next_button.configure(text="Close", command=self.finish_install)
        elif page_num == len(self.pages) - 2:  # Installation page
            self.next_button.configure(text="Install", command=self.start_installation)
        else:
            self.next_button.configure(text="Continue", command=self.go_next)
        
        # Special handling for license page
        if page_num == 1:
            self.next_button.configure(state=tk.DISABLED)
        
        # Update summary page
        if page_num == 4:  # Summary page
            self.update_summary()
        
        self.current_page = page_num
    
    def update_summary(self):
        components = []
        if self.install_gui.get():
            components.append("Desktop GUI Application")
        if self.install_service.get():
            components.append("Background Service (launchd)")
        if self.create_app_bundle.get():
            components.append("Application Bundle")
        
        summary = f"""Installation Configuration:

Installation Directory: {self.install_dir.get()}
Data Directory: {self.data_dir.get()}
Configuration Directory: {self.config_dir.get()}

Components to Install:
• Core Runtime (servin)
"""
        
        for component in components:
            summary += f"• {component}\n"
        
        if self.create_user.get():
            summary += "\nAdvanced Options:\n• Create system user '_servin'\n"
        if self.add_to_path.get():
            summary += "• Add to PATH environment\n"
        if self.install_cli_tools.get():
            summary += "• Install command-line tools\n"
        
        summary += "\nClick 'Install' to begin the installation process."
        
        self.summary_text.configure(state=tk.NORMAL)
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
        if messagebox.askyesno("Cancel Installation", 
                              "Are you sure you want to cancel the installation?"):
            self.root.quit()
    
    def start_installation(self):
        if not self.is_root:
            messagebox.showerror("Administrator Required", 
                               "Administrator privileges are required for installation.\n\n"
                               "Please run: sudo python3 servin-installer.py")
            return
        
        self.show_page(len(self.pages) - 2)  # Installation page
        self.next_button.configure(state=tk.DISABLED)
        self.back_button.configure(state=tk.DISABLED)
        self.cancel_button.configure(state=tk.DISABLED)
        
        # Start progress animation
        self.progress.start(10)
        
        # Start installation
        self.root.after(500, self.perform_installation)
    
    def log_message(self, message):
        self.log_text.insert(tk.END, f"{message}\n")
        self.log_text.see(tk.END)
        self.root.update()
    
    def perform_installation(self):
        try:
            self.status_label.configure(text="Creating directories...")
            self.log_message("🚀 Starting Servin installation for macOS")
            
            # Create directories
            directories = [
                self.install_dir.get(),
                self.data_dir.get(),
                self.config_dir.get(),
                f"{self.data_dir.get()}/volumes",
                f"{self.data_dir.get()}/images",
                "/usr/local/var/log/servin"
            ]
            
            for directory in directories:
                os.makedirs(directory, exist_ok=True)
                self.log_message(f"📁 Created directory: {directory}")
            
            # Install binaries
            self.status_label.configure(text="Installing binaries...")
            script_dir = os.path.dirname(os.path.abspath(__file__))
            
            binaries = ["servin"]
            if self.install_gui.get():
                binaries.append("servin-gui")
            
            for binary in binaries:
                src = os.path.join(script_dir, binary)
                dst = os.path.join(self.install_dir.get(), binary)
                if os.path.exists(src):
                    shutil.copy2(src, dst)
                    os.chmod(dst, 0o755)
                    self.log_message(f"📦 Installed: {binary}")
                else:
                    self.log_message(f"⚠️  Warning: {binary} not found in installer package")
            
            # Create configuration
            self.status_label.configure(text="Creating configuration...")
            config_content = f"""# Servin Configuration File for macOS
data_dir={self.data_dir.get()}
log_level=info
log_file=/usr/local/var/log/servin/servin.log
runtime=native
bridge_name=servin0
cri_port=10250
cri_enabled=false
gui_theme=auto
enable_notifications=true"""
            
            with open(f"{self.config_dir.get()}/servin.conf", 'w') as f:
                f.write(config_content)
            self.log_message("⚙️  Created configuration file")
            
            # Create system user
            if self.create_user.get():
                self.status_label.configure(text="Creating system user...")
                self.create_system_user()
            
            # Install launchd service
            if self.install_service.get():
                self.status_label.configure(text="Installing service...")
                self.install_launchd_service()
            
            # Create application bundle
            if self.create_app_bundle.get() and self.install_gui.get():
                self.status_label.configure(text="Creating application bundle...")
                self.create_application_bundle()
            
            self.progress.stop()
            self.status_label.configure(text="Installation completed successfully!")
            self.log_message("✅ Installation completed successfully!")
            
            # Update success page
            self.update_success_page()
            
            # Enable next button
            self.next_button.configure(state=tk.NORMAL, text="Continue", command=self.go_next)
            
        except Exception as e:
            self.progress.stop()
            self.log_message(f"❌ Error: {str(e)}")
            messagebox.showerror("Installation Error", f"Installation failed:\n\n{str(e)}")
            self.cancel_button.configure(state=tk.NORMAL)
    
    def create_system_user(self):
        try:
            # Find next available UID in system range
            uid = 200
            while uid < 500:
                try:
                    pwd.getpwuid(uid)
                    uid += 1
                except KeyError:
                    break
            
            # Create user using dscl
            commands = [
                ['dscl', '.', '-create', '/Users/_servin'],
                ['dscl', '.', '-create', '/Users/_servin', 'UserShell', '/usr/bin/false'],
                ['dscl', '.', '-create', '/Users/_servin', 'RealName', 'Servin Runtime User'],
                ['dscl', '.', '-create', '/Users/_servin', 'UniqueID', str(uid)],
                ['dscl', '.', '-create', '/Users/_servin', 'PrimaryGroupID', '20'],
                ['dscl', '.', '-create', '/Users/_servin', 'NFSHomeDirectory', '/var/empty']
            ]
            
            for cmd in commands:
                subprocess.run(cmd, check=True, capture_output=True)
            
            self.log_message(f"👤 Created system user '_servin' with UID {uid}")
            
        except Exception as e:
            self.log_message(f"⚠️  Warning: Could not create system user: {str(e)}")
    
    def install_launchd_service(self):
        plist_path = "/Library/LaunchDaemons/com.servin.runtime.plist"
        
        plist_data = {
            'Label': 'com.servin.runtime',
            'ProgramArguments': [
                f"{self.install_dir.get()}/servin",
                'daemon',
                '--config',
                f"{self.config_dir.get()}/servin.conf"
            ],
            'UserName': '_servin' if self.create_user.get() else 'root',
            'GroupName': 'staff',
            'RunAtLoad': True,
            'KeepAlive': {
                'SuccessfulExit': False,
                'Crashed': True
            },
            'StandardOutPath': '/usr/local/var/log/servin/servin.stdout.log',
            'StandardErrorPath': '/usr/local/var/log/servin/servin.stderr.log',
            'WorkingDirectory': self.data_dir.get(),
            'EnvironmentVariables': {
                'PATH': '/usr/local/bin:/usr/bin:/bin'
            },
            'ThrottleInterval': 10
        }
        
        with open(plist_path, 'wb') as f:
            plistlib.dump(plist_data, f)
        
        os.chmod(plist_path, 0o644)
        subprocess.run(['launchctl', 'load', plist_path], check=False)
        
        self.log_message("🔧 Installed launchd service")
    
    def create_application_bundle(self):
        app_path = "/Applications/Servin GUI.app"
        contents_path = f"{app_path}/Contents"
        macos_path = f"{contents_path}/MacOS"
        resources_path = f"{contents_path}/Resources"
        
        # Create directory structure
        for path in [macos_path, resources_path]:
            os.makedirs(path, exist_ok=True)
        
        # Copy executable
        gui_exe = f"{self.install_dir.get()}/servin-gui"
        if os.path.exists(gui_exe):
            shutil.copy2(gui_exe, f"{macos_path}/Servin GUI")
            os.chmod(f"{macos_path}/Servin GUI", 0o755)
        
        # Create Info.plist
        info_plist = {
            'CFBundleExecutable': 'Servin GUI',
            'CFBundleIdentifier': 'com.servin.gui',
            'CFBundleName': 'Servin GUI',
            'CFBundleDisplayName': 'Servin GUI',
            'CFBundleVersion': '1.0.0',
            'CFBundleShortVersionString': '1.0.0',
            'CFBundlePackageType': 'APPL',
            'CFBundleSignature': 'SERV',
            'LSMinimumSystemVersion': '10.12',
            'NSHighResolutionCapable': True,
            'LSApplicationCategoryType': 'public.app-category.developer-tools'
        }
        
        with open(f"{contents_path}/Info.plist", 'wb') as f:
            plistlib.dump(info_plist, f)
        
        self.log_message("🍎 Created application bundle")
    
    def update_success_page(self):
        components = ["Core Runtime (servin)"]
        
        if self.install_gui.get():
            components.append("Desktop GUI Application")
        if self.install_service.get():
            components.append("Background Service (launchd)")
        if self.create_app_bundle.get():
            components.append("Application Bundle")
        
        summary_text = "Successfully installed:\n" + "\n".join(f"• {comp}" for comp in components)
        
        if hasattr(self, 'installed_summary'):
            self.installed_summary.configure(text=summary_text)
    
    def finish_install(self):
        # Perform final actions
        if self.start_service.get() and self.install_service.get():
            try:
                subprocess.run(['launchctl', 'start', 'com.servin.runtime'], check=True)
                self.log_message("✅ Started Servin service")
            except:
                messagebox.showwarning("Service", "Could not start Servin service automatically.\n"
                                     "You can start it manually from System Preferences.")
        
        if self.launch_gui.get() and self.install_gui.get():
            try:
                if self.create_app_bundle.get():
                    subprocess.Popen(['open', '/Applications/Servin GUI.app'], start_new_session=True)
                else:
                    gui_path = f"{self.install_dir.get()}/servin-gui"
                    subprocess.Popen([gui_path], start_new_session=True)
            except:
                pass
        
        if self.open_applications.get():
            subprocess.Popen(['open', '/Applications'], start_new_session=True)
        
        messagebox.showinfo("Installation Complete", 
                           "Servin Container Runtime has been successfully installed!\n\n"
                           "You can now use the 'servin' command in Terminal or launch "
                           "the GUI from your Applications folder.")
        self.root.quit()
    
    def run(self):
        # Center window on screen
        self.root.eval('tk::PlaceWindow . center')
        self.root.mainloop()


if __name__ == "__main__":
    # Check if running on macOS
    if os.uname().sysname != 'Darwin':
        print("Error: This installer is designed for macOS only.")
        print("Please use the appropriate installer for your platform.")
        sys.exit(1)
    
    # Check if tkinter is available
    try:
        import tkinter
    except ImportError:
        print("Error: tkinter is not available.")
        print("Please install Python with tkinter support or use the command-line installer:")
        print("sudo ./install.sh")
        sys.exit(1)
    
    installer = ServinMacInstaller()
    installer.run()
