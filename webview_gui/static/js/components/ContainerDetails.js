/**
 * Container Details Component
 * Handles the detailed view of individual containers
 */

class ContainerDetails {
    constructor(apiClient, socketManager) {
        this.apiClient = apiClient;
        this.socketManager = socketManager;
        this.currentContainerId = null;
        this.activeTab = 'logs';
        
        this.init();
    }

    init() {
        this.setupEventListeners();
        this.setupSocketHandlers();
    }

    setupEventListeners() {
        // Back to containers button
        const backBtn = document.getElementById('backToContainers');
        if (backBtn) {
            backBtn.addEventListener('click', () => this.hide());
        }

        // Tab navigation
        document.querySelectorAll('.tab-btn').forEach(tabBtn => {
            tabBtn.addEventListener('click', (e) => {
                const tabName = e.target.dataset.tab;
                this.switchTab(tabName);
            });
        });

        // Container action buttons
        this.setupActionButtons();
    }

    setupActionButtons() {
        const startBtn = document.getElementById('startContainerBtn');
        const stopBtn = document.getElementById('stopContainerBtn');
        const restartBtn = document.getElementById('restartContainerBtn');
        const removeBtn = document.getElementById('removeContainerBtn');

        if (startBtn) {
            startBtn.addEventListener('click', () => this.startContainer());
        }

        if (stopBtn) {
            stopBtn.addEventListener('click', () => this.stopContainer());
        }

        if (restartBtn) {
            restartBtn.addEventListener('click', () => this.restartContainer());
        }

        if (removeBtn) {
            removeBtn.addEventListener('click', () => this.removeContainer());
        }
    }

    setupSocketHandlers() {
        // Socket handlers for real-time updates will be set up here
        // This is where log streaming, exec sessions, etc. will be handled
    }

    async show(containerId) {
        console.log('Showing details for container:', containerId);
        this.currentContainerId = containerId;
        
        try {
            // Show the details view
            document.getElementById('containerDetails').style.display = 'block';
            document.getElementById('containersList').style.display = 'none';
            
            // Load container details
            const container = await this.apiClient.getContainerDetails(containerId);
            this.renderContainerInfo(container);
            
            // Setup tabs and load default content
            this.setupTabEventListeners();
            this.switchTab('logs');
            
        } catch (error) {
            console.error('Failed to load container details:', error);
            UIHelpers.showToast('Failed to load container details', 'error');
        }
    }

    hide() {
        console.log('Hiding container details');
        document.getElementById('containerDetails').style.display = 'none';
        document.getElementById('containersList').style.display = 'block';
        this.currentContainerId = null;
    }

    renderContainerInfo(container) {
        // Update header information
        const nameElement = document.getElementById('detailsContainerName');
        const idElement = document.getElementById('detailsContainerId');

        if (nameElement) {
            nameElement.textContent = container.name || 'Unnamed';
        }
        
        if (idElement) {
            idElement.textContent = container.id.substring(0, 12); // Show short ID
        }

        // Update overview tab with detailed information
        this.renderOverview(container);
        
        // Update action buttons based on container status
        this.updateActionButtons(container.status);
    }

    renderOverview(container) {
        const overviewContainer = document.getElementById('overviewTab');
        if (!overviewContainer) return;

        overviewContainer.innerHTML = `
            <div class="overview-grid">
                <div class="overview-card">
                    <h4>Basic Information</h4>
                    <div class="info-grid">
                        <div class="info-item">
                            <label>Container ID:</label>
                            <span>${container.id}</span>
                        </div>
                        <div class="info-item">
                            <label>Name:</label>
                            <span>${container.name || '-'}</span>
                        </div>
                        <div class="info-item">
                            <label>Image:</label>
                            <span>${container.image}</span>
                        </div>
                        <div class="info-item">
                            <label>Status:</label>
                            <span class="status-badge status-${container.status.toLowerCase()}">${container.status}</span>
                        </div>
                        <div class="info-item">
                            <label>Created:</label>
                            <span>${UIHelpers.formatDate(container.created)}</span>
                        </div>
                        <div class="info-item">
                            <label>Started:</label>
                            <span>${UIHelpers.formatDate(container.started_at)}</span>
                        </div>
                    </div>
                </div>
                
                <div class="overview-card">
                    <h4>Network</h4>
                    <div class="info-grid">
                        <div class="info-item">
                            <label>Ports:</label>
                            <span>${container.ports || '-'}</span>
                        </div>
                        <div class="info-item">
                            <label>Network Mode:</label>
                            <span>${container.network_mode || '-'}</span>
                        </div>
                        <div class="info-item">
                            <label>IP Address:</label>
                            <span>${container.ip_address || '-'}</span>
                        </div>
                    </div>
                </div>
                
                <div class="overview-card">
                    <h4>Resources</h4>
                    <div class="info-grid">
                        <div class="info-item">
                            <label>Memory Limit:</label>
                            <span>${container.memory_limit ? UIHelpers.formatBytes(container.memory_limit) : 'No limit'}</span>
                        </div>
                        <div class="info-item">
                            <label>CPU Limit:</label>
                            <span>${container.cpu_limit || 'No limit'}</span>
                        </div>
                        <div class="info-item">
                            <label>Restart Policy:</label>
                            <span>${container.restart_policy || '-'}</span>
                        </div>
                    </div>
                </div>
            </div>
        `;

        // Re-setup action button listeners after rendering
        this.setupActionButtons();
    }

    updateActionButtons(containerStatus) {
        const startBtn = document.getElementById('startContainerBtn');
        const stopBtn = document.getElementById('stopContainerBtn');
        const restartBtn = document.getElementById('restartContainerBtn');
        const removeBtn = document.getElementById('removeContainerBtn');

        // Hide all buttons first
        [startBtn, stopBtn, restartBtn, removeBtn].forEach(btn => {
            if (btn) btn.style.display = 'none';
        });

        // Show appropriate buttons based on container status
        if (containerStatus === 'running') {
            // Running container: show restart and stop buttons only
            if (restartBtn) restartBtn.style.display = 'inline-flex';
            if (stopBtn) stopBtn.style.display = 'inline-flex';
        } else {
            // Stopped container: show start (play) and delete buttons only
            if (startBtn) startBtn.style.display = 'inline-flex';
            if (removeBtn) removeBtn.style.display = 'inline-flex';
        }
    }

    setupTabEventListeners() {
        console.log('Setting up tab event listeners');
        
        // Remove existing listeners by replacing elements
        document.querySelectorAll('.tab-btn').forEach(btn => {
            const newBtn = btn.cloneNode(true);
            btn.parentNode.replaceChild(newBtn, btn);
        });
        
        // Add new event listeners with proper context binding
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const tabName = e.target.dataset.tab;
                this.switchTab(tabName);
            });
        });
    }

    switchTab(tabName) {
        console.log('Switching to tab:', tabName);
        this.activeTab = tabName;
        
        // Update tab buttons
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.classList.remove('active');
        });
        
        const targetTab = document.querySelector(`[data-tab="${tabName}"]`);
        if (targetTab) {
            targetTab.classList.add('active');
        }
        
        // Update tab content
        document.querySelectorAll('.tab-pane').forEach(pane => {
            pane.classList.remove('active');
            pane.style.display = 'none';
        });
        
        const targetPane = document.getElementById(`${tabName}Tab`);
        if (targetPane) {
            targetPane.classList.add('active');
            targetPane.style.display = 'flex';
        }
        
        // Load tab-specific content
        this.loadTabContent(tabName);
    }

    async loadTabContent(tabName) {
        if (!this.currentContainerId) return;
        
        try {
            switch (tabName) {
                case 'logs':
                    await this.loadLogs();
                    break;
                case 'files':
                    await this.loadFiles();
                    break;
                case 'exec':
                    this.setupExec();
                    break;
                case 'env':
                    await this.loadEnvironment();
                    break;
                case 'volumes':
                    await this.loadVolumes();
                    break;
                case 'network':
                    await this.loadNetwork();
                    break;
                case 'stats':
                    await this.loadStats();
                    break;
            }
        } catch (error) {
            console.error(`Failed to load ${tabName} content:`, error);
            UIHelpers.showToast(`Failed to load ${tabName} content`, 'error');
        }
    }

    async loadLogs() {
        // Logs functionality will be handled by the Logs component
        // For now, show a placeholder
        const logsContent = document.getElementById('logsContent');
        if (logsContent) {
            logsContent.innerHTML = '<div class="loading">Loading logs...</div>';
        }
    }

    async loadFiles(path = '/') {
        // File explorer functionality will be handled by the FileExplorer component
        const filesContent = document.getElementById('filesContent');
        if (filesContent) {
            filesContent.innerHTML = '<div class="loading">Loading files...</div>';
        }
    }

    setupExec() {
        // Initialize Terminal component for this container
        if (window.terminal && this.currentContainerId) {
            window.terminal.setupTerminal(this.currentContainerId);
        }
    }

    async loadEnvironment() {
        try {
            const response = await this.apiClient.getContainerEnvironment(this.currentContainerId);
            // Extract the environment data from the response
            const envVars = response.environment || response;
            this.renderEnvironmentVariables(envVars);
        } catch (error) {
            console.error('Failed to load environment variables:', error);
            const envContent = document.getElementById('envContent');
            if (envContent) {
                envContent.innerHTML = '<div class="error">Failed to load environment variables</div>';
            }
        }
    }

    renderEnvironmentVariables(envVars) {
        const envContent = document.getElementById('envContent');
        if (!envContent) return;

        if (!envVars || (Array.isArray(envVars) && envVars.length === 0) || 
            (typeof envVars === 'object' && Object.keys(envVars).length === 0)) {
            envContent.innerHTML = '<div class="empty">No environment variables found</div>';
            return;
        }

        let html = '<div class="env-list">';
        
        if (Array.isArray(envVars)) {
            envVars.forEach(envVar => {
                if (typeof envVar === 'string') {
                    const [key, ...valueParts] = envVar.split('=');
                    const value = valueParts.join('=');
                    html += `
                        <div class="env-item">
                            <span class="env-key">${this.escapeHtml(key || '')}</span>
                            <span class="env-value">${this.escapeHtml(value || '')}</span>
                        </div>
                    `;
                } else if (typeof envVar === 'object' && envVar !== null) {
                    // Handle object format like {key: "KEY", value: "VALUE"}
                    const key = envVar.key || envVar.name || '';
                    const value = envVar.value || '';
                    html += `
                        <div class="env-item">
                            <span class="env-key">${this.escapeHtml(String(key))}</span>
                            <span class="env-value">${this.escapeHtml(String(value))}</span>
                        </div>
                    `;
                }
            });
        } else if (typeof envVars === 'object' && envVars !== null) {
            Object.entries(envVars).forEach(([key, value]) => {
                html += `
                    <div class="env-item">
                        <span class="env-key">${this.escapeHtml(String(key))}</span>
                        <span class="env-value">${this.escapeHtml(String(value))}</span>
                    </div>
                `;
            });
        }
        
        html += '</div>';
        envContent.innerHTML = html;
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    async loadVolumes() {
        // Placeholder for volumes content
        const volumesContent = document.getElementById('volumesContent');
        if (volumesContent) {
            volumesContent.innerHTML = '<div class="placeholder">Volume information will be displayed here</div>';
        }
    }

    async loadNetwork() {
        // Placeholder for network content
        const networkContent = document.getElementById('networkContent');
        if (networkContent) {
            networkContent.innerHTML = '<div class="placeholder">Network information will be displayed here</div>';
        }
    }

    async loadStats() {
        // Placeholder for stats content
        const statsContent = document.getElementById('statsContent');
        if (statsContent) {
            statsContent.innerHTML = '<div class="placeholder">Container statistics will be displayed here</div>';
        }
    }

    // Container action methods
    async startContainer() {
        if (!this.currentContainerId) return;
        
        try {
            const response = await this.apiClient.startContainer(this.currentContainerId);
            
            // Check if we got the new container information
            if (response.success && response.new_container) {
                UIHelpers.showToast('Container started successfully', 'success');
                
                // Update to show the new container
                this.currentContainerId = response.new_container.id;
                
                // Refresh container details with the new container ID
                setTimeout(() => this.show(this.currentContainerId), 1000);
                
                // Also notify parent component to refresh the container list
                if (this.onContainerUpdated) {
                    this.onContainerUpdated();
                }
            } else {
                // Fallback for old response format or no new container info
                UIHelpers.showToast('Container started successfully', 'success');
                setTimeout(() => this.show(this.currentContainerId), 1000);
            }
        } catch (error) {
            console.error('Failed to start container:', error);
            UIHelpers.showToast('Failed to start container', 'error');
        }
    }

    async stopContainer() {
        if (!this.currentContainerId) return;
        
        try {
            await this.apiClient.stopContainer(this.currentContainerId);
            UIHelpers.showToast('Container stopped successfully', 'success');
            // Refresh container details
            setTimeout(() => this.show(this.currentContainerId), 1000);
        } catch (error) {
            console.error('Failed to stop container:', error);
            UIHelpers.showToast('Failed to stop container', 'error');
        }
    }

    async restartContainer() {
        if (!this.currentContainerId) return;
        
        try {
            await this.apiClient.restartContainer(this.currentContainerId);
            UIHelpers.showToast('Container restarted successfully', 'success');
            // Refresh container details
            setTimeout(() => this.show(this.currentContainerId), 1000);
        } catch (error) {
            console.error('Failed to restart container:', error);
            UIHelpers.showToast('Failed to restart container', 'error');
        }
    }

    async removeContainer() {
        if (!this.currentContainerId) return;
        
        if (!confirm('Are you sure you want to remove this container?')) {
            return;
        }
        
        try {
            await this.apiClient.removeContainer(this.currentContainerId);
            UIHelpers.showToast('Container removed successfully', 'success');
            // Go back to containers list
            this.hide();
        } catch (error) {
            console.error('Failed to remove container:', error);
            UIHelpers.showToast('Failed to remove container', 'error');
        }
    }
}

// Export the component
window.ContainerDetails = ContainerDetails;