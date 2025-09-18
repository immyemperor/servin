# Servin Documentation

This directory contains the complete documentation for Servin Container Runtime, built with Jekyll and hosted on GitHub Pages.

[![Contributors](https://img.shields.io/github/contributors/immyemperor/servin?style=flat-square&logo=github)](https://github.com/immyemperor/servin/graphs/contributors)
[![Documentation](https://img.shields.io/badge/docs-live-brightgreen?style=flat-square&logo=github-pages)](https://immyemperor.github.io/servin)
[![Issues](https://img.shields.io/github/issues/immyemperor/servin?style=flat-square&logo=github)](https://github.com/immyemperor/servin/issues)
[![License](https://img.shields.io/github/license/immyemperor/servin?style=flat-square)](https://github.com/immyemperor/servin/blob/master/LICENSE)

## 🌐 Live Documentation

Visit the live documentation at: **https://immyemperor.github.io/servin**

## 📁 Structure

```
docs/
├── _config.yml              # Jekyll configuration
├── _layouts/                # Page layouts
│   └── default.html        # Main layout with sidebar
├── _includes/               # Reusable components
│   ├── head.html           # HTML head section
│   ├── header.html         # Site header
│   └── footer.html         # Site footer
├── assets/                  # CSS and other assets
│   └── sidebar.css         # Custom styles for sidebar
├── index.md                 # Home page
├── overview.md              # Project overview
├── architecture.md          # Architecture documentation
├── installation.md          # Installation guide
└── [other-pages].md        # Additional documentation pages
```

## 🚀 Local Development

### Prerequisites

- Ruby 3.0+
- Bundler
- Jekyll

### Setup

1. **Install dependencies:**
   ```bash
   cd docs
   bundle install
   ```

2. **Start development server:**
   ```bash
   bundle exec jekyll serve
   ```

3. **Open in browser:**
   ```
   http://localhost:4000/servin
   ```

### Live Reload

The development server automatically reloads when you make changes to:
- Markdown files
- HTML layouts
- CSS files
- Configuration

## ✨ Features

### 📱 Responsive Design
- Mobile-friendly sidebar navigation
- Responsive grid layouts
- Touch-friendly interface

### 🔍 Search Functionality
- Client-side search in sidebar
- Filter navigation items
- Fast and responsive

### 🎨 Professional Styling
- Modern, clean design
- Syntax highlighting for code blocks
- Professional typography
- GitHub-compatible markdown

### 📊 Navigation
- Hierarchical sidebar navigation
- Active page highlighting
- Smooth scrolling
- Keyboard navigation support

## 📝 Content Management

### Adding New Pages

1. **Create markdown file:**
   ```bash
   touch docs/new-page.md
   ```

2. **Add front matter:**
   ```yaml
   ---
   layout: default
   title: New Page
   permalink: /new-page/
   ---
   ```

3. **Add to navigation:**
   Edit `_layouts/default.html` to add navigation link

### Writing Content

- Use standard markdown syntax
- Add front matter to all pages
- Use relative URLs: `{{ '/page' | relative_url }}`
- Include code syntax highlighting: ```language

### Styling Components

Available CSS classes:
- `.feature-grid` - Responsive grid layout
- `.feature-box` - Feature highlight boxes
- `.badge` - Status badges
- `.btn` - Button styling
- `.architecture-diagram` - Monospace diagrams

## 🔧 Configuration

### Jekyll Configuration (`_config.yml`)

Key settings:
- Site title and description
- GitHub repository information
- Navigation order
- Plugin configuration

### Custom Styling (`assets/sidebar.css`)

- Sidebar navigation styling
- Responsive breakpoints
- Color scheme variables
- Component styles

## 📤 Deployment

### Automatic Deployment

Documentation is automatically deployed via GitHub Actions when:
- Changes are pushed to `main` or `master` branch
- Files in `docs/` directory are modified

### Manual Deployment

1. **Build site:**
   ```bash
   cd docs
   bundle exec jekyll build
   ```

2. **Deploy to GitHub Pages:**
   - Enable GitHub Pages in repository settings
   - Select "GitHub Actions" as source
   - Push changes to trigger deployment

## 🎯 GitHub Pages Setup

1. **Enable GitHub Pages:**
   - Go to repository Settings
   - Navigate to Pages section
   - Select "GitHub Actions" as source

2. **Configure custom domain (optional):**
   - Add `CNAME` file with your domain
   - Configure DNS settings

3. **Enable HTTPS:**
   - GitHub Pages automatically provides HTTPS
   - Check "Enforce HTTPS" option

## 🔍 SEO Optimization

The documentation includes:
- Meta tags for social sharing
- Structured data markup
- Sitemap generation
- Search engine optimization
- Fast loading times

## 📱 Mobile Experience

Optimized for mobile devices:
- Collapsible sidebar navigation
- Touch-friendly interface
- Responsive typography
- Fast loading on mobile networks

## 🤝 Contributing

To contribute to the documentation:

1. Fork the repository
2. Create a feature branch
3. Make your changes in the `docs/` directory
4. Test locally with Jekyll
5. Submit a pull request

### Documentation Guidelines

- Write clear, concise content
- Use proper markdown formatting
- Include code examples where helpful
- Test all links and references
- Follow the existing style and structure

## 📞 Support

- **Documentation Issues**: [GitHub Issues](https://github.com/immyemperor/servin/issues)
- **General Support**: See main README
- **Feature Requests**: Submit via GitHub Issues

---

**Built with ❤️ using Jekyll and GitHub Pages**

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

## � Contributors

Servin Container Runtime is built and maintained by an amazing community of developers:

### **Core Team**
- **[Brijesh Kumar](https://github.com/immyemperor)** - Project Lead & Architecture
- **[Abhishek Kumar](https://github.com/abhishek-kumar3)** - Core Development & Features

### **Community Contributions**
We welcome contributions from developers worldwide! Here's how you can contribute:

- **📖 Documentation**: Improve guides, fix typos, add examples
- **🐛 Bug Reports**: Help us identify and fix issues
- **💡 Feature Requests**: Suggest new capabilities and improvements
- **💻 Code Contributions**: Implement features, optimize performance
- **🧪 Testing**: Test on different platforms and report compatibility
- **🌍 Translations**: Help make Servin accessible globally

### **Special Thanks**
- All community members who report issues and provide feedback
- Beta testers helping validate VM mode across different platforms
- Documentation contributors improving our guides and examples

### **Contributing Guidelines**
Ready to contribute? Check out our:
- [Contributing Guide]({{ '/development' | relative_url }}) - Development setup and guidelines
- [GitHub Issues](https://github.com/immyemperor/servin/issues) - Bug reports and feature requests
- [GitHub Discussions](https://github.com/immyemperor/servin/discussions) - Community discussions

## �📝 License

This documentation is part of the Servin Container Runtime project and is licensed under the Apache License 2.0.

---

**💡 Tip**: Bookmark `http://localhost:8080/wiki.html` for quick access to your local wiki during development!
