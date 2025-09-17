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
        if (refreshBtn) {
            refreshBtn.addEventListener('click', () => this.refresh());
        }
    }

    async loadFiles(containerId, path = '/') {
        console.log('FileExplorer: Loading files for container:', containerId, 'path:', path);
        
        this.currentContainerId = containerId;
        this.currentPath = path;
        
        const filesContent = document.getElementById('filesContent');
        const currentPathElement = document.getElementById('currentPath');
        
        if (!filesContent) return;
        
        filesContent.innerHTML = '<div class="loading">Loading files...</div>';
        
        if (currentPathElement) {
            currentPathElement.textContent = path;
        }

        // Setup file controls
        this.setupFileControls();

        try {
            const files = await this.apiClient.getContainerFiles(containerId, path);
            this.renderFiles(files, path);
        } catch (error) {
            console.error('Failed to load files:', error);
            filesContent.innerHTML = '<div class="error">Failed to load files</div>';
        }
    }

    setupFileControls() {
        const filesHeader = document.querySelector('.files-header');
        if (!filesHeader) return;

        filesHeader.innerHTML = `
            <div class="files-toolbar">
                <div class="files-navigation">
                    <div class="files-path">
                        <i class="fas fa-folder"></i>
                        <span id="currentPath">${this.currentPath}</span>
                    </div>
                    <button class="btn btn-secondary btn-sm" id="refreshFilesBtn">
                        <i class="fas fa-sync-alt"></i> Refresh
                    </button>
                </div>
            </div>
        `;

        // Re-setup refresh button
        const refreshBtn = document.getElementById('refreshFilesBtn');
        if (refreshBtn) {
            refreshBtn.addEventListener('click', () => this.refresh());
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
                    <i class="fas fa-level-up-alt"></i>
                    <span>.. (Parent Directory)</span>
                    <span class="file-size"></span>
                </div>
            `;
        }

        if (!files || files.length === 0) {
            html += '<div class="empty">No files found in this directory</div>';
        } else {
            // Sort files: directories first, then files
            const sortedFiles = [...files].sort((a, b) => {
                if (a.type === 'directory' && b.type !== 'directory') return -1;
                if (a.type !== 'directory' && b.type === 'directory') return 1;
                return a.name.localeCompare(b.name);
            });

            sortedFiles.forEach(file => {
                const isDirectory = file.type === 'directory';
                const icon = isDirectory ? 'fa-folder' : 'fa-file';
                const clickHandler = isDirectory ? 
                    `onclick="window.fileExplorer.navigateToPath('${this.joinPath(currentPath, file.name)}')"` : 
                    '';
                
                html += `
                    <div class="file-item ${isDirectory ? 'directory' : 'file'}" ${clickHandler}>
                        <i class="fas ${icon}"></i>
                        <span class="file-name">${file.name}</span>
                        <span class="file-size">${isDirectory ? '-' : UIHelpers.formatFileSize(file.size || 0)}</span>
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