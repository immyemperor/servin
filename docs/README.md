# Servin Container Runtime Documentation

This directory contains comprehensive documentation for the Servin Container Runtime project.

## 📚 Documentation Files

### **HTML Wiki**
- **`wiki.html`** - Complete interactive HTML wiki with navigation
- **Features**: Responsive design, search functionality, syntax highlighting, copy buttons
- **Sections**: Installation, usage guides, API reference, troubleshooting, and more

### **Wiki Server Scripts**
- **`serve-wiki.py`** - Python HTTP server to serve the wiki locally
- **`serve-wiki.bat`** - Windows batch script to start the wiki server
- **`serve-wiki.sh`** - Linux/macOS shell script to start the wiki server

## 🚀 Quick Start

### **View the Wiki Locally**

#### **Windows**
```cmd
# Navigate to docs directory
cd docs

# Double-click serve-wiki.bat or run:
serve-wiki.bat
```

#### **Linux/macOS**
```bash
# Navigate to docs directory
cd docs

# Make executable and run
chmod +x serve-wiki.sh
./serve-wiki.sh

# Or run Python script directly
python3 serve-wiki.py
```

### **Direct File Access**
You can also open `wiki.html` directly in any web browser:
- Right-click `wiki.html` → "Open with" → Browser
- Or navigate to file:///path/to/servin/docs/wiki.html

## 🌟 Wiki Features

### **🎯 Interactive Navigation**
- **Sidebar navigation** with categorized sections
- **Search functionality** to quickly find topics
- **Smooth scrolling** and section highlighting
- **Responsive design** for desktop and mobile

### **💻 Enhanced Code Blocks**
- **Syntax highlighting** for shell commands and code
- **Copy buttons** on all code blocks
- **Multiple language support** (Bash, PowerShell, Go, YAML)
- **Proper formatting** with monospace fonts

### **📱 User Experience**
- **Professional styling** with modern design
- **Fast loading** with optimized CSS and JavaScript
- **Accessibility** features for screen readers
- **Cross-browser compatibility**

## 📖 Wiki Content

### **📋 Main Sections**

1. **🚀 Project Overview**
   - Features and capabilities
   - Target users and use cases
   - Key benefits

2. **🏗 Architecture**
   - System overview and components
   - Directory structure
   - Component relationships

3. **✨ Features**
   - Core runtime features
   - Kubernetes CRI integration
   - User interface options

4. **🛠 Installation**
   - Quick installation guides
   - Building from source
   - Platform-specific instructions

5. **💻 User Interfaces**
   - Command Line Interface (CLI)
   - Terminal User Interface (TUI)
   - Desktop GUI Application

6. **🎯 Core Features**
   - Container Management
   - Image Management
   - Volume Management
   - Registry Operations

7. **🔌 Integration**
   - Kubernetes CRI implementation
   - API reference (REST and gRPC)
   - Logging and monitoring

8. **👨‍💻 Development**
   - Building from source
   - Development environment setup
   - Contributing guidelines

9. **🔧 Support**
   - Troubleshooting guide
   - FAQ section
   - Community resources

## 🎨 Customization

### **Styling**
The wiki uses CSS custom properties (variables) for easy theming:
```css
:root {
    --primary-color: #2563eb;
    --secondary-color: #1e40af;
    --accent-color: #3b82f6;
    /* ... more variables */
}
```

### **Content Updates**
To update the wiki content:
1. Edit the HTML sections in `wiki.html`
2. Update navigation links if adding new sections
3. Test locally using the serve scripts
4. Commit changes to version control

### **Adding Sections**
To add new sections:
1. Add navigation item in the sidebar
2. Create corresponding content section with unique ID
3. Update the JavaScript navigation handler
4. Test navigation and search functionality

## 📊 Server Configuration

### **Default Settings**
- **Port**: 8080
- **Host**: localhost (127.0.0.1)
- **Auto-open**: Yes (opens browser automatically)
- **Directory**: Current docs/ directory

### **Custom Configuration**
Edit `serve-wiki.py` to customize:
```python
PORT = 8080          # Change port
# Add authentication, HTTPS, etc.
```

## 🔍 Search Functionality

The wiki includes real-time search that filters navigation items:
- **Type** in the search box to filter topics
- **Clear** the search to show all topics
- **Case-insensitive** matching
- **Instant** results as you type

## 📱 Mobile Support

The wiki is fully responsive and includes:
- **Collapsible sidebar** on mobile devices
- **Touch-friendly** navigation
- **Optimized fonts** and spacing
- **Fast loading** on slower connections

## 🚀 Production Deployment

For production deployment:

### **Static Hosting**
Upload `wiki.html` to any static hosting service:
- GitHub Pages
- Netlify
- Vercel
- AWS S3 + CloudFront

### **Web Server**
Serve with any web server:
```bash
# Nginx
server {
    listen 80;
    server_name wiki.example.com;
    root /path/to/docs;
    index wiki.html;
}

# Apache
DocumentRoot /path/to/docs
DirectoryIndex wiki.html
```

## 📝 License

This documentation is part of the Servin Container Runtime project and is licensed under the Apache License 2.0.

---

**💡 Tip**: Bookmark `http://localhost:8080/wiki.html` for quick access to your local wiki during development!
