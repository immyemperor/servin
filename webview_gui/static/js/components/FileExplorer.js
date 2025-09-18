/**
 * File Explorer Component
 * Handles container filesystem browsing
 */

class FileExplorer {
    constructor(apiClient) {
        this.apiClient = apiClient;
        this.currentContainerId = null;
        this.currentPath = '/';
        
        this.init();
    }

    init() {
        this.setupControls();
    }

    setupControls() {
        const refreshBtn = document.getElementById('refreshFilesBtn');
        const goBackBtn = document.getElementById('goBackBtn');
        const goToRootBtn = document.getElementById('goToRootBtn');
        
        if (refreshBtn) {
            refreshBtn.addEventListener('click', () => this.refresh());
        }
        
        if (goBackBtn) {
            goBackBtn.addEventListener('click', () => this.goBack());
        }
        
        if (goToRootBtn) {
            goToRootBtn.addEventListener('click', () => this.goToRoot());
        }
    }

    async loadFiles(containerId, path = '/') {
        console.log('FileExplorer: Loading files for container:', containerId, 'path:', path);
        
        this.currentContainerId = containerId;
        this.currentPath = path;
        
        const filesContent = document.getElementById('filesContent');
        
        if (!filesContent) return;
        
        filesContent.innerHTML = '<div class="loading">Loading files...</div>';
        
        // Update breadcrumb navigation
        this.updateBreadcrumb(path);
        
        // Update navigation buttons state
        this.updateNavigationButtons();

        try {
            const files = await this.apiClient.getContainerFiles(containerId, path);
            this.renderFiles(files, path);
        } catch (error) {
            console.error('Failed to load files:', error);
            filesContent.innerHTML = `
                <div class="error">
                    <i class="fas fa-exclamation-triangle"></i>
                    <p>Failed to load files</p>
                    <small>${error.message || 'Unknown error'}</small>
                </div>
            `;
        }
    }

    updateBreadcrumb(path) {
        const breadcrumbPath = document.getElementById('breadcrumbPath');
        if (!breadcrumbPath) return;
        
        // Clear existing breadcrumb
        breadcrumbPath.innerHTML = '';
        
        // Split path into segments
        const segments = path.split('/').filter(segment => segment !== '');
        
        // Add root segment
        const rootSegment = document.createElement('span');
        rootSegment.className = 'path-segment root';
        rootSegment.dataset.path = '/';
        rootSegment.innerHTML = '<i class="fas fa-home"></i>';
        rootSegment.addEventListener('click', () => this.navigateToPath('/'));
        breadcrumbPath.appendChild(rootSegment);
        
        // Add path segments
        let currentPath = '';
        segments.forEach((segment, index) => {
            currentPath += '/' + segment;
            
            // Add separator
            const separator = document.createElement('span');
            separator.className = 'path-separator';
            separator.textContent = '/';
            breadcrumbPath.appendChild(separator);
            
            // Add segment
            const segmentElement = document.createElement('span');
            segmentElement.className = 'path-segment';
            segmentElement.dataset.path = currentPath;
            segmentElement.textContent = segment;
            segmentElement.addEventListener('click', () => this.navigateToPath(currentPath));
            breadcrumbPath.appendChild(segmentElement);
        });
    }
    
    updateNavigationButtons() {
        const goBackBtn = document.getElementById('goBackBtn');
        const goToRootBtn = document.getElementById('goToRootBtn');
        
        if (goBackBtn) {
            goBackBtn.disabled = this.currentPath === '/';
        }
        
        if (goToRootBtn) {
            goToRootBtn.disabled = this.currentPath === '/';
        }
    }

    renderFiles(files, currentPath) {
        const filesContent = document.getElementById('filesContent');
        if (!filesContent) return;

        let html = '<div class="files-list">';

        // Add parent directory link if not root
        if (currentPath !== '/') {
            const parentPath = currentPath.split('/').slice(0, -1).join('/') || '/';
            html += `
                <div class="file-item parent-dir" onclick="window.fileExplorer.navigateToPath('${parentPath}')">
                    <div class="file-icon">
                        <i class="fas fa-level-up-alt"></i>
                    </div>
                    <div class="file-info">
                        <div class="file-name">.. (Parent Directory)</div>
                        <div class="file-meta">Go up one level</div>
                    </div>
                    <div class="file-actions"></div>
                </div>
            `;
        }

        if (!files || files.length === 0) {
            html += `
                <div class="empty-state">
                    <i class="fas fa-folder-open"></i>
                    <p>This directory is empty</p>
                </div>
            `;
        } else {
            // Sort files: directories first, then files
            const sortedFiles = [...files].sort((a, b) => {
                if (a.type === 'directory' && b.type !== 'directory') return -1;
                if (a.type !== 'directory' && b.type === 'directory') return 1;
                return a.name.localeCompare(b.name);
            });

            sortedFiles.forEach(file => {
                const isDirectory = file.type === 'directory';
                const isSymlink = file.type === 'symlink';
                
                let icon = 'fa-file';
                if (isDirectory) {
                    icon = 'fa-folder';
                } else if (isSymlink) {
                    icon = 'fa-link';
                } else {
                    // Set specific icons based on file extension
                    const ext = file.name.toLowerCase().split('.').pop();
                    switch (ext) {
                        case 'js': case 'json': icon = 'fa-file-code'; break;
                        case 'py': icon = 'fa-file-code'; break;
                        case 'html': case 'htm': icon = 'fa-file-code'; break;
                        case 'css': icon = 'fa-file-code'; break;
                        case 'txt': case 'md': case 'log': icon = 'fa-file-alt'; break;
                        case 'jpg': case 'jpeg': case 'png': case 'gif': icon = 'fa-file-image'; break;
                        case 'pdf': icon = 'fa-file-pdf'; break;
                        case 'zip': case 'tar': case 'gz': icon = 'fa-file-archive'; break;
                        default: icon = 'fa-file';
                    }
                }
                
                const clickHandler = isDirectory ? 
                    `onclick="window.fileExplorer.navigateToPath('${this.joinPath(currentPath, file.name)}')"` : 
                    '';
                
                const fileSize = isDirectory ? '-' : UIHelpers.formatFileSize(file.size || 0);
                const permissions = file.permissions || '-';
                
                html += `
                    <div class="file-item ${file.type}" ${clickHandler}>
                        <div class="file-icon">
                            <i class="fas ${icon}"></i>
                        </div>
                        <div class="file-info">
                            <div class="file-name">${this.escapeHtml(file.name)}</div>
                            <div class="file-meta">
                                <span class="file-permissions">${permissions}</span>
                                <span class="file-size">${fileSize}</span>
                            </div>
                        </div>
                        <div class="file-actions">
                            ${isDirectory ? '<i class="fas fa-chevron-right"></i>' : ''}
                        </div>
                    </div>
                `;
            });
        }

        html += '</div>';
        filesContent.innerHTML = html;
    }

    navigateToPath(path) {
        console.log('Navigating to path:', path);
        this.loadFiles(this.currentContainerId, path);
    }

    joinPath(basePath, fileName) {
        if (basePath === '/') {
            return '/' + fileName;
        }
        return basePath + '/' + fileName;
    }
    
    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    refresh() {
        if (this.currentContainerId) {
            this.loadFiles(this.currentContainerId, this.currentPath);
        }
    }

    goBack() {
        if (this.currentPath !== '/') {
            const parentPath = this.currentPath.split('/').slice(0, -1).join('/') || '/';
            this.navigateToPath(parentPath);
        }
    }

    goToRoot() {
        this.navigateToPath('/');
    }

    cleanup() {
        this.currentContainerId = null;
        this.currentPath = '/';
    }
}

// Export the component
window.FileExplorer = FileExplorer;