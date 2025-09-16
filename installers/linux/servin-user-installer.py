#!/usr/bin/env python3
"""
Servin Container Runtime - Linux User Installation Wizard
A graphical installer for Linux systems using tkinter
Installs to user directories without requiring sudo
"""

import tkinter as tk
from tkinter import ttk, messagebox, filedialog
import os
import sys
import subprocess
import shutil
import tempfile
from pathlib import Path

class ServinLinuxInstaller:
    def __init__(self):
        self.root = tk.Tk()
        self.root.title("Servin Container Runtime - User Installer")
        self.root.geometry("700x550")
        self.root.resizable(False, False)
        
        # Set window icon if available
        try:
            # Try to load icon from the same directory as the installer
            script_dir = os.path.dirname(os.path.abspath(__file__))
            icon_path = os.path.join(script_dir, "servin-icon-64.png")
            if os.path.exists(icon_path):
                from PIL import Image, ImageTk
                icon_image = Image.open(icon_path)
                self.icon_photo = ImageTk.PhotoImage(icon_image)
                self.root.iconphoto(True, self.icon_photo)
            else:
                # Fallback to trying different icon formats
                for icon_file in ["servin.ico", "servin-icon-48.png", "servin-icon-32.png"]:
                    icon_path = os.path.join(script_dir, icon_file)
                    if os.path.exists(icon_path):
                        if icon_file.endswith('.png'):
                            from PIL import Image, ImageTk
                            icon_image = Image.open(icon_path)
                            self.icon_photo = ImageTk.PhotoImage(icon_image)
                            self.root.iconphoto(True, self.icon_photo)
                        break
        except:
            pass  # No icon if PIL not available or icon not found
        
        # Get user home directory
        self.home_dir = Path.home()
        
        # Installation variables - all in user directories
        self.install_dir = tk.StringVar(value=str(self.home_dir / ".local" / "bin"))
        self.data_dir = tk.StringVar(value=str(self.home_dir / ".local" / "share" / "servin"))
        self.config_dir = tk.StringVar(value=str(self.home_dir / ".config" / "servin"))
        self.install_desktop = tk.BooleanVar(value=True)
        self.add_to_path = tk.BooleanVar(value=True)
        self.create_desktop_files = tk.BooleanVar(value=True)
        
        # Current page
        self.current_page = 0
        self.pages = []
        
        self.setup_ui()
        self.show_page(0)
    
    def setup_ui(self):
        # Main frame
        main_frame = ttk.Frame(self.root, padding="10")
        main_frame.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        
        # Header
        header_frame = ttk.Frame(main_frame)
        header_frame.grid(row=0, column=0, sticky=(tk.W, tk.E), pady=(0, 20))
        
        title_label = ttk.Label(header_frame, text="Servin Container Runtime", 
                               font=("Arial", 18, "bold"))
        title_label.grid(row=0, column=0, sticky=tk.W)
        
        subtitle_label = ttk.Label(header_frame, text="User Installation Wizard", 
                                  font=("Arial", 12))
        subtitle_label.grid(row=1, column=0, sticky=tk.W)
        
        # Content frame
        self.content_frame = ttk.Frame(main_frame)
        self.content_frame.grid(row=1, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        
        # Button frame
        button_frame = ttk.Frame(main_frame)
        button_frame.grid(row=2, column=0, sticky=(tk.W, tk.E), pady=(20, 0))
        
        button_container = ttk.Frame(button_frame)
        button_container.grid(row=0, column=1, sticky=tk.E)
        
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
            self.create_destination_page,
            self.create_components_page,
            self.create_summary_page,
            self.create_installation_page,
            self.create_success_page
        ]
    
    def create_welcome_page(self):
        page = ttk.Frame(self.content_frame, padding="20")
        
        # Welcome message
        welcome_frame = ttk.Frame(page)
        welcome_frame.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N), pady=(0, 30))
        
        title = ttk.Label(welcome_frame, text="Welcome to Servin Container Runtime", 
                         font=("Arial", 16, "bold"))
        title.grid(row=0, column=0, sticky=tk.W, pady=(0, 15))
        
        description = """This installer will set up Servin Container Runtime in your user directories. No root privileges required!

Servin provides:
• Docker-compatible container management
• Native Linux desktop application  
• Command-line tools for automation
• Lightweight containerization with Linux namespaces

Features of this user installation:
• Installs to ~/.local/bin (no sudo required)
• Automatic PATH configuration
• User-specific configuration
• Desktop integration via .desktop files"""
        
        desc_label = ttk.Label(welcome_frame, text=description, font=("Arial", 11), 
                              justify=tk.LEFT, wraplength=600)
        desc_label.grid(row=1, column=0, sticky=(tk.W, tk.E))
        
        # Benefits of user installation
        benefits_frame = ttk.LabelFrame(page, text="Benefits of User Installation", padding="15")
        benefits_frame.grid(row=1, column=0, sticky=(tk.W, tk.E), pady=(0, 20))
        
        benefits = [
            "✅ No root password required",
            "✅ Safe installation in user directories", 
            "✅ Won't affect system-wide software",
            "✅ Easy to uninstall or update",
            "✅ Automatic shell integration"
        ]
        
        for i, benefit in enumerate(benefits):
            ttk.Label(benefits_frame, text=benefit, font=("Arial", 10)).grid(row=i, column=0, sticky=tk.W, pady=2)
        
        return page
    
    def create_destination_page(self):
        page = ttk.Frame(self.content_frame, padding="20")
        
        ttk.Label(page, text="Installation Directories", 
                 font=("Arial", 14, "bold")).grid(row=0, column=0, sticky=tk.W, pady=(0, 20))
        
        # Destination selection
        dest_frame = ttk.LabelFrame(page, text="User Directory Locations", padding="15")
        dest_frame.grid(row=1, column=0, sticky=(tk.W, tk.E), pady=(0, 20))
        dest_frame.columnconfigure(1, weight=1)
        
        ttk.Label(dest_frame, text="Binaries:").grid(row=0, column=0, sticky=tk.W, pady=5)
        install_entry = ttk.Entry(dest_frame, textvariable=self.install_dir, width=50)
        install_entry.grid(row=0, column=1, sticky=(tk.W, tk.E), padx=(10, 5))
        ttk.Button(dest_frame, text="Choose...", command=lambda: self.browse_directory(self.install_dir)).grid(row=0, column=2)
        
        ttk.Label(dest_frame, text="Data:").grid(row=1, column=0, sticky=tk.W, pady=5)
        data_entry = ttk.Entry(dest_frame, textvariable=self.data_dir, width=50)
        data_entry.grid(row=1, column=1, sticky=(tk.W, tk.E), padx=(10, 5))
        ttk.Button(dest_frame, text="Choose...", command=lambda: self.browse_directory(self.data_dir)).grid(row=1, column=2)
        
        ttk.Label(dest_frame, text="Config:").grid(row=2, column=0, sticky=tk.W, pady=5)
        config_entry = ttk.Entry(dest_frame, textvariable=self.config_dir, width=50)
        config_entry.grid(row=2, column=1, sticky=(tk.W, tk.E), padx=(10, 5))
        ttk.Button(dest_frame, text="Choose...", command=lambda: self.browse_directory(self.config_dir)).grid(row=2, column=2)
        
        # Information note
        info_frame = ttk.LabelFrame(page, text="Path Information", padding="15")
        info_frame.grid(row=2, column=0, sticky=(tk.W, tk.E))
        
        info_text = f"""Installation follows XDG Base Directory specification:

• Binaries: {self.install_dir.get()}
• Data: {self.data_dir.get()}  
• Configuration: {self.config_dir.get()}

The installer will automatically add {self.install_dir.get()} to your PATH."""

        ttk.Label(info_frame, text=info_text, font=("Arial", 10), justify=tk.LEFT).grid(row=0, column=0, sticky=tk.W)
        
        return page
    
    def create_components_page(self):
        page = ttk.Frame(self.content_frame, padding="20")
        
        ttk.Label(page, text="Select Components", 
                 font=("Arial", 14, "bold")).grid(row=0, column=0, sticky=tk.W, pady=(0, 20))
        
        # Components
        components_frame = ttk.LabelFrame(page, text="Components to Install", padding="15")
        components_frame.grid(row=1, column=0, sticky=(tk.W, tk.E), pady=(0, 20))
        
        ttk.Checkbutton(components_frame, text="Core Runtime (servin) - Required", 
                       state=tk.DISABLED, variable=tk.BooleanVar(value=True)).grid(row=0, column=0, sticky=tk.W, pady=3)
        
        ttk.Checkbutton(components_frame, text="Desktop Application (servin-tui)", 
                       variable=self.install_desktop).grid(row=1, column=0, sticky=tk.W, pady=3)
        
        ttk.Checkbutton(components_frame, text="Create Desktop Files (.desktop)", 
                       variable=self.create_desktop_files).grid(row=2, column=0, sticky=tk.W, pady=3)
        
        # Configuration
        config_frame = ttk.LabelFrame(page, text="Configuration Options", padding="15")
        config_frame.grid(row=2, column=0, sticky=(tk.W, tk.E))
        
        ttk.Checkbutton(config_frame, text="Add to PATH environment (recommended)", 
                       variable=self.add_to_path).grid(row=0, column=0, sticky=tk.W, pady=3)
        
        # PATH explanation
        path_info = """Adding to PATH allows you to run 'servin' from any terminal window.
This modifies your shell configuration files (.bashrc, .zshrc, etc.)"""
        
        ttk.Label(config_frame, text=path_info, font=("Arial", 9), 
                 foreground="gray", justify=tk.LEFT).grid(row=1, column=0, sticky=tk.W, padx=(20, 0))
        
        return page
    
    def create_summary_page(self):
        page = ttk.Frame(self.content_frame, padding="20")
        
        ttk.Label(page, text="Installation Summary", 
                 font=("Arial", 14, "bold")).grid(row=0, column=0, sticky=tk.W, pady=(0, 15))
        
        # Summary content
        summary_frame = ttk.Frame(page)
        summary_frame.grid(row=1, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        summary_frame.columnconfigure(0, weight=1)
        summary_frame.rowconfigure(0, weight=1)
        
        self.summary_text = tk.Text(summary_frame, height=15, wrap=tk.WORD, 
                                   font=("monospace", 10), borderwidth=1, relief=tk.SOLID,
                                   bg="white", fg="black")
        scrollbar = ttk.Scrollbar(summary_frame, orient=tk.VERTICAL, command=self.summary_text.yview)
        self.summary_text.configure(yscrollcommand=scrollbar.set)
        
        self.summary_text.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        scrollbar.grid(row=0, column=1, sticky=(tk.N, tk.S))
        
        # Update summary content
        self.update_summary()
        
        return page
    
    def update_summary(self):
        """Update the summary text with current settings"""
        summary = f"""Servin Container Runtime Installation Summary

INSTALLATION DIRECTORIES:
• Binaries: {self.install_dir.get()}
• Data: {self.data_dir.get()}
• Configuration: {self.config_dir.get()}

COMPONENTS:
• Core Runtime (servin): ✓ Included
• Desktop Application: {'✓ Included' if self.install_desktop.get() else '✗ Skipped'}
• Desktop Files: {'✓ Included' if self.create_desktop_files.get() else '✗ Skipped'}

CONFIGURATION:
• Add to PATH: {'✓ Yes' if self.add_to_path.get() else '✗ No'}
• Installation Type: User installation (no sudo required)
• Shell Integration: Automatic

WHAT WILL BE INSTALLED:
• servin - Main container runtime CLI
{'• servin-tui - Desktop GUI application' if self.install_desktop.get() else ''}
• Configuration files and documentation
• Shell integration for PATH
{'• Desktop entries for application menu' if self.create_desktop_files.get() else ''}

POST-INSTALLATION:
• Commands will be available: servin
{'• Desktop app available in application menu' if self.create_desktop_files.get() else ''}
• Configuration files will be created on first run
• No system-wide changes will be made

The installation is completely contained within your user directories and can be
easily uninstalled by removing the installation directories."""

        self.summary_text.delete(1.0, tk.END)
        self.summary_text.insert(tk.END, summary)
        self.summary_text.configure(state=tk.DISABLED)
    
    def create_installation_page(self):
        page = ttk.Frame(self.content_frame, padding="20")
        
        ttk.Label(page, text="Installing Servin Container Runtime", 
                 font=("Arial", 14, "bold")).grid(row=0, column=0, sticky=tk.W, pady=(0, 20))
        
        # Progress bar
        self.progress = ttk.Progressbar(page, mode='indeterminate')
        self.progress.grid(row=1, column=0, sticky=(tk.W, tk.E), pady=(0, 10))
        
        # Status label
        self.status_label = ttk.Label(page, text="Preparing installation...")
        self.status_label.grid(row=2, column=0, sticky=tk.W, pady=(0, 10))
        
        # Log area
        log_frame = ttk.LabelFrame(page, text="Installation Log", padding="10")
        log_frame.grid(row=3, column=0, sticky=(tk.W, tk.E, tk.N, tk.S), pady=(10, 0))
        
        self.log_text = tk.Text(log_frame, height=12, font=("monospace", 9))
        log_scrollbar = ttk.Scrollbar(log_frame, orient=tk.VERTICAL, command=self.log_text.yview)
        self.log_text.configure(yscrollcommand=log_scrollbar.set)
        
        self.log_text.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        log_scrollbar.grid(row=0, column=1, sticky=(tk.N, tk.S))
        
        return page
    
    def create_success_page(self):
        page = ttk.Frame(self.content_frame, padding="20")
        
        # Success message
        success_frame = ttk.Frame(page)
        success_frame.grid(row=0, column=0, sticky=(tk.W, tk.E), pady=(0, 30))
        
        title = ttk.Label(success_frame, text="Installation Completed Successfully!", 
                         font=("Arial", 16, "bold"), foreground="green")
        title.grid(row=0, column=0, sticky=tk.W, pady=(0, 15))
        
        description = """Servin Container Runtime has been successfully installed to your user directories.

Next steps:
1. Open a new terminal window (to load PATH changes)
2. Run 'servin --help' to see available commands
3. Try pulling an image: 'servin image pull alpine'
4. Create and run a container: 'servin run alpine echo "Hello World"'

The desktop application is available in your application menu."""
        
        desc_label = ttk.Label(success_frame, text=description, font=("Arial", 11), 
                              justify=tk.LEFT, wraplength=600)
        desc_label.grid(row=1, column=0, sticky=(tk.W, tk.E))
        
        # Quick start commands
        commands_frame = ttk.LabelFrame(page, text="Quick Start Commands", padding="15")
        commands_frame.grid(row=1, column=0, sticky=(tk.W, tk.E))
        
        commands = [
            "servin --help",
            "servin image pull alpine",
            "servin run alpine echo 'Hello World'"
        ]
        
        for i, cmd in enumerate(commands):
            cmd_label = ttk.Label(commands_frame, text=f"$ {cmd}", font=("monospace", 10))
            cmd_label.grid(row=i, column=0, sticky=tk.W, pady=2)
        
        return page
    
    def show_page(self, page_num):
        # Clear content frame
        for widget in self.content_frame.winfo_children():
            widget.destroy()
        
        # Show current page
        if 0 <= page_num < len(self.pages):
            self.current_page = page_num
            page = self.pages[page_num]()
            page.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
            
            # Update buttons
            self.back_button.configure(state=tk.NORMAL if page_num > 0 else tk.DISABLED)
            
            if page_num == len(self.pages) - 2:  # Installation page
                self.next_button.configure(state=tk.DISABLED)
                self.start_installation()
            elif page_num == len(self.pages) - 1:  # Success page
                self.next_button.configure(text="Finish", command=self.finish_installation)
            else:
                self.next_button.configure(text="Continue", command=self.go_next, state=tk.NORMAL)
    
    def go_back(self):
        if self.current_page > 0:
            self.show_page(self.current_page - 1)
    
    def go_next(self):
        if self.current_page < len(self.pages) - 1:
            self.show_page(self.current_page + 1)
    
    def browse_directory(self, var):
        directory = filedialog.askdirectory(initialdir=var.get())
        if directory:
            var.set(directory)
    
    def start_installation(self):
        self.progress.start()
        self.root.after(100, self.perform_installation)
    
    def log_message(self, message):
        self.log_text.insert(tk.END, f"{message}\n")
        self.log_text.see(tk.END)
        self.root.update()
    
    def perform_installation(self):
        try:
            # Create directories
            self.status_label.configure(text="Creating directories...")
            self.log_message("🚀 Starting Servin installation")
            
            directories = [
                self.install_dir.get(),
                self.data_dir.get(),
                self.config_dir.get(),
                f"{self.data_dir.get()}/volumes",
                f"{self.data_dir.get()}/images",
                f"{self.data_dir.get()}/containers",
                f"{self.data_dir.get()}/logs"
            ]
            
            for directory in directories:
                os.makedirs(directory, exist_ok=True)
                self.log_message(f"📁 Created directory: {directory}")
            
            # Install binaries
            self.status_label.configure(text="Installing binaries...")
            script_dir = os.path.dirname(os.path.abspath(__file__))
            
            binaries = ["servin"]
            if self.install_desktop.get():
                binaries.append("servin-tui")
            
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
            config_content = f"""# Servin Configuration File
data_dir={self.data_dir.get()}
log_level=info
log_file={self.data_dir.get()}/logs/servin.log
runtime=native
bridge_name=servin0
gui_theme=auto
enable_notifications=true"""
            
            config_file = os.path.join(self.config_dir.get(), "servin.conf")
            with open(config_file, 'w') as f:
                f.write(config_content)
            self.log_message("⚙️  Created configuration file")
            
            # Setup PATH
            if self.add_to_path.get():
                self.status_label.configure(text="Setting up PATH...")
                self.setup_path()
            
            # Create desktop files
            if self.create_desktop_files.get() and self.install_desktop.get():
                self.status_label.configure(text="Creating desktop entries...")
                self.create_desktop_files_func()
            
            self.progress.stop()
            self.status_label.configure(text="Installation completed successfully!")
            self.log_message("✅ Installation completed successfully!")
            
            # Enable next button
            self.next_button.configure(state=tk.NORMAL)
            
        except Exception as e:
            self.progress.stop()
            self.log_message(f"❌ Error: {str(e)}")
            messagebox.showerror("Installation Error", f"Installation failed:\n\n{str(e)}")
    
    def setup_path(self):
        """Add installation directory to PATH by modifying shell configuration files"""
        install_bin = self.install_dir.get()
        
        # Shell configuration files to update
        shell_configs = [
            self.home_dir / ".bashrc",
            self.home_dir / ".zshrc",
            self.home_dir / ".profile"
        ]
        
        # PATH line to add
        path_line = f'export PATH="{install_bin}:$PATH"  # Added by Servin installer'
        
        for config_file in shell_configs:
            if config_file.exists() or config_file.name in [".bashrc", ".profile"]:
                try:
                    # Read existing content
                    if config_file.exists():
                        with open(config_file, 'r') as f:
                            content = f.read()
                    else:
                        content = ""
                    
                    # Check if already added
                    if "Added by Servin installer" not in content:
                        # Add PATH export
                        with open(config_file, 'a') as f:
                            f.write(f"\n# Servin Container Runtime\n{path_line}\n")
                        self.log_message(f"📝 Updated {config_file.name}")
                    else:
                        self.log_message(f"📝 {config_file.name} already configured")
                        
                except Exception as e:
                    self.log_message(f"⚠️  Warning: Could not update {config_file.name}: {e}")
    
    def create_desktop_files_func(self):
        """Create .desktop files for Linux desktop integration"""
        desktop_dir = self.home_dir / ".local" / "share" / "applications"
        desktop_dir.mkdir(parents=True, exist_ok=True)
        
        # Copy icon to user's icon directory
        icon_dir = self.home_dir / ".local" / "share" / "icons" / "hicolor"
        script_dir = os.path.dirname(os.path.abspath(__file__))
        
        icon_installed = False
        icon_name = "servin-tui"
        
        # Try to install icon in different sizes
        icon_sizes = [16, 32, 48, 64, 128, 256]
        for size in icon_sizes:
            icon_source = os.path.join(script_dir, f"servin-icon-{size}.png")
            if os.path.exists(icon_source):
                size_dir = icon_dir / f"{size}x{size}" / "apps"
                size_dir.mkdir(parents=True, exist_ok=True)
                icon_target = size_dir / f"{icon_name}.png"
                shutil.copy2(icon_source, icon_target)
                icon_installed = True
                self.log_message(f"🎨 Installed {size}x{size} icon")
        
        # Fallback to copying any available icon
        if not icon_installed:
            for icon_file in ["servin-icon-64.png", "servin-icon-48.png", "servin.ico"]:
                icon_source = os.path.join(script_dir, icon_file)
                if os.path.exists(icon_source) and icon_file.endswith('.png'):
                    apps_dir = icon_dir / "48x48" / "apps"
                    apps_dir.mkdir(parents=True, exist_ok=True)
                    icon_target = apps_dir / f"{icon_name}.png"
                    shutil.copy2(icon_source, icon_target)
                    icon_installed = True
                    self.log_message("🎨 Installed fallback icon")
                    break
        
        # Desktop file content
        icon_reference = icon_name if icon_installed else "application-x-executable"
        
        desktop_content = f"""[Desktop Entry]
Version=1.0
Type=Application
Name=Servin Desktop
Comment=Container Management with Servin
Exec={self.install_dir.get()}/servin-tui
Icon={icon_reference}
Terminal=false
Categories=Development;System;
Keywords=container;docker;runtime;
StartupNotify=true
"""
        
        desktop_file = desktop_dir / "servin-tui.desktop"
        with open(desktop_file, 'w') as f:
            f.write(desktop_content)
        
        # Make executable
        os.chmod(desktop_file, 0o755)
        
        self.log_message("🖥️  Created desktop entry")
    
    def finish_installation(self):
        self.root.quit()

if __name__ == "__main__":
    installer = ServinLinuxInstaller()
    installer.root.mainloop()