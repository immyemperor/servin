/**
 * Docker Desktop GUI - Frontend JavaScript
 * Handles all UI interactions and API communications
 */

class DockerGUI {
    constructor() {
        this.apiBase = '';
        this.refreshInterval = null;
        this.currentSection = 'containers';
        this.currentContainerId = null;
        this.data = {
            containers: [],
            images: [],
            volumes: []
        };
        
        this.init();
    }
    
    init() {
        this.setupEventListeners();
        this.checkDockerConnection();
        this.startAutoRefresh();
        this.loadData();
    }
    
    setupEventListeners() {
        // Navigation
        document.querySelectorAll('.nav-item').forEach(item => {
            item.addEventListener('click', (e) => {
                const section = e.currentTarget.dataset.section;
                this.switchSection(section);
            });
        });
        
        // Refresh button
        const refreshBtn = document.getElementById('refreshBtn');
        if (refreshBtn) {
            refreshBtn.addEventListener('click', () => {
                this.loadData();
            });
        }
        
        // Modal controls
        this.setupModalControls();
        
        // Search functionality
        this.setupSearchFilters();
        
        // Form submissions
        this.setupFormHandlers();
    }
    
    setupModalControls() {
        // Pull Image Modal - only setup if elements exist
        const pullImageBtn = document.getElementById('pullImageBtn');
        const pullImageModal = document.getElementById('pullImageModal');
        const closePullModal = document.getElementById('closePullModal');
        const cancelPullBtn = document.getElementById('cancelPullBtn');
        
        if (pullImageBtn && pullImageModal) {
            pullImageBtn.addEventListener('click', () => {
                pullImageModal.style.display = 'block';
            });
            
            [closePullModal, cancelPullBtn].forEach(btn => {
                if (btn) {
                    btn.addEventListener('click', () => {
                        pullImageModal.style.display = 'none';
                        const form = document.getElementById('pullImageForm');
                        if (form) form.reset();
                    });
                }
            });
        }
        
        // Create Volume Modal - only setup if elements exist
        const createVolumeBtn = document.getElementById('createVolumeBtn');
        const createVolumeModal = document.getElementById('createVolumeModal');
        const closeVolumeModal = document.getElementById('closeVolumeModal');
        const cancelVolumeBtn = document.getElementById('cancelVolumeBtn');
        
        if (createVolumeBtn && createVolumeModal) {
            createVolumeBtn.addEventListener('click', () => {
                createVolumeModal.style.display = 'block';
            });
            
            [closeVolumeModal, cancelVolumeBtn].forEach(btn => {
                if (btn) {
                    btn.addEventListener('click', () => {
                        createVolumeModal.style.display = 'none';
                        const form = document.getElementById('createVolumeForm');
                        if (form) form.reset();
                    });
                }
            });
        }

        // Container Details functionality
        this.setupContainerDetailsControls();
        
        // Close modals when clicking outside
        window.addEventListener('click', (e) => {
            if (e.target.classList.contains('modal')) {
                e.target.style.display = 'none';
            }
        });
    }

    setupContainerDetailsControls() {
        // Back to containers button
        const backBtn = document.getElementById('backToContainers');
        if (backBtn) {
            backBtn.addEventListener('click', () => {
                this.hideContainerDetails();
            });
        }

        // Tab navigation
        document.querySelectorAll('.tab-btn').forEach(tabBtn => {
            tabBtn.addEventListener('click', (e) => {
                const tabName = e.currentTarget.dataset.tab;
                this.switchTab(tabName);
            });
        });

        // Container action buttons
        const startBtn = document.getElementById('startContainerBtn');
        const stopBtn = document.getElementById('stopContainerBtn');
        const restartBtn = document.getElementById('restartContainerBtn');
        const removeBtn = document.getElementById('removeContainerBtn');

        if (startBtn) {
            startBtn.addEventListener('click', () => {
                const containerId = this.currentContainerId;
                if (containerId) this.startContainer(containerId);
            });
        }

        if (stopBtn) {
            stopBtn.addEventListener('click', () => {
                const containerId = this.currentContainerId;
                if (containerId) this.stopContainer(containerId);
            });
        }

        if (restartBtn) {
            restartBtn.addEventListener('click', () => {
                const containerId = this.currentContainerId;
                if (containerId) this.restartContainer(containerId);
            });
        }

        if (removeBtn) {
            removeBtn.addEventListener('click', () => {
                const containerId = this.currentContainerId;
                if (containerId && confirm('Are you sure you want to remove this container?')) {
                    this.removeContainer(containerId);
                }
            });
        }
    }
    
    setupSearchFilters() {
        ['container', 'image', 'volume'].forEach(type => {
            const searchInput = document.getElementById(`${type}Search`);
            if (searchInput) {
                searchInput.addEventListener('input', (e) => {
                    this.filterTable(type, e.target.value);
                });
            }
        });
    }
    
    setupFormHandlers() {
        // Pull Image Form
        const pullImageForm = document.getElementById('pullImageForm');
        if (pullImageForm) {
            pullImageForm.addEventListener('submit', async (e) => {
                e.preventDefault();
                const imageName = document.getElementById('imageName').value.trim();
                if (imageName) {
                    await this.pullImage(imageName);
                    document.getElementById('pullImageModal').style.display = 'none';
                    document.getElementById('pullImageForm').reset();
                }
            });
        }
        
        // Create Volume Form
        const createVolumeForm = document.getElementById('createVolumeForm');
        if (createVolumeForm) {
            createVolumeForm.addEventListener('submit', async (e) => {
                e.preventDefault();
                const volumeName = document.getElementById('volumeName').value.trim();
                if (volumeName) {
                    await this.createVolume(volumeName);
                    document.getElementById('createVolumeModal').style.display = 'none';
                    document.getElementById('createVolumeForm').reset();
                }
            });
        }
    }
    
    switchSection(section) {
        // Update navigation
        document.querySelectorAll('.nav-item').forEach(item => {
            item.classList.remove('active');
        });
        document.querySelector(`[data-section="${section}"]`).classList.add('active');
        
        // Update content sections
        document.querySelectorAll('.content-section').forEach(sec => {
            sec.classList.remove('active');
        });
        document.getElementById(`${section}Section`).classList.add('active');
        
        this.currentSection = section;
    }
    
    async checkDockerConnection() {
        try {
            const response = await fetch(`${this.apiBase}/api/system/info`);
            const statusIndicator = document.getElementById('statusIndicator');
            const statusText = document.getElementById('statusText');
            
            if (response.ok) {
                statusIndicator.classList.add('connected');
                statusText.textContent = 'Connected';
            } else {
                statusIndicator.classList.remove('connected');
                statusText.textContent = 'Disconnected';
            }
        } catch (error) {
            const statusIndicator = document.getElementById('statusIndicator');
            const statusText = document.getElementById('statusText');
            statusIndicator.classList.remove('connected');
            statusText.textContent = 'Error';
        }
    }
    
    async loadData() {
        this.showLoading();
        
        try {
            await Promise.all([
                this.loadContainers(),
                this.loadImages(),
                this.loadVolumes()
            ]);
        } catch (error) {
            this.showToast('Error loading data', 'error');
        } finally {
            this.hideLoading();
        }
    }
    
    async loadContainers() {
        try {
            const response = await fetch(`${this.apiBase}/api/containers`);
            if (response.ok) {
                this.data.containers = await response.json();
                this.renderContainers();
                this.updateCounts();
            } else {
                throw new Error('Failed to load containers');
            }
        } catch (error) {
            console.error('Error loading containers:', error);
            this.showToast('Failed to load containers', 'error');
        }
    }
    
    async loadImages() {
        try {
            const response = await fetch(`${this.apiBase}/api/images`);
            if (response.ok) {
                this.data.images = await response.json();
                this.renderImages();
                this.updateCounts();
            } else {
                throw new Error('Failed to load images');
            }
        } catch (error) {
            console.error('Error loading images:', error);
            this.showToast('Failed to load images', 'error');
        }
    }
    
    async loadVolumes() {
        try {
            const response = await fetch(`${this.apiBase}/api/volumes`);
            if (response.ok) {
                this.data.volumes = await response.json();
                this.renderVolumes();
                this.updateCounts();
            } else {
                throw new Error('Failed to load volumes');
            }
        } catch (error) {
            console.error('Error loading volumes:', error);
            this.showToast('Failed to load volumes', 'error');
        }
    }
    
    renderContainers() {
        const tbody = document.getElementById('containersTableBody');
        const emptyState = document.getElementById('containersEmpty');
        const table = document.getElementById('containersTable');
        
        if (this.data.containers.length === 0) {
            table.style.display = 'none';
            emptyState.style.display = 'flex';
            return;
        }
        
        table.style.display = 'table';
        emptyState.style.display = 'none';
        
        console.log('=== Rendering containers ===');
        console.log('Container data:', this.data.containers);
        
        tbody.innerHTML = this.data.containers.map(container => {
            const isRunning = (container.state === 'running' || container.status === 'running');
            const status = container.status || container.state || 'unknown';
            
            console.log(`Container ${container.name}: isRunning=${isRunning}, status=${status}, state=${container.state}`);
            
            return `
            <tr data-id="${container.id}" class="container-row clickable" 
                onclick="dockerGUI.showContainerDetails('${container.id}')" title="Click to view details">
                <td>
                    <strong>${container.name}</strong>
                    <br>
                    <small class="text-muted">${container.id}</small>
                </td>
                <td>${container.image}</td>
                <td>
                    <span class="status-badge status-${status.toLowerCase()}">
                        ${status}
                    </span>
                </td>
                <td>${container.ports && container.ports.length > 0 ? container.ports.join(', ') : '-'}</td>
                <td>
                    <small class="text-muted">
                        ${this.formatDate(container.created)}
                    </small>
                </td>
                <td>
                    <div class="action-buttons" onclick="event.stopPropagation()">
                        ${isRunning
                            ? `<button class="action-btn stop" onclick="dockerGUI.stopContainer('${container.id}')" title="Stop">
                                <i class="fas fa-stop"></i>
                               </button>`
                            : `<button class="action-btn start" onclick="dockerGUI.startContainer('${container.id}')" title="Start">
                                <i class="fas fa-play"></i>
                               </button>`
                        }
                        <button class="action-btn" onclick="dockerGUI.restartContainer('${container.id}')" title="Restart">
                            <i class="fas fa-redo"></i>
                        </button>
                        <button class="action-btn remove" onclick="dockerGUI.removeContainer('${container.id}')" title="Remove">
                            <i class="fas fa-trash"></i>
                        </button>
                        <button class="action-btn details" onclick="dockerGUI.showContainerDetails('${container.id}')" title="Details">
                            <i class="fas fa-info-circle"></i>
                        </button>
                    </div>
                </td>
            </tr>`;
        }).join('');
        
        console.log('Containers rendered. HTML:', tbody.innerHTML.substring(0, 500) + '...');
    }
    
    renderImages() {
        const tbody = document.getElementById('imagesTableBody');
        const emptyState = document.getElementById('imagesEmpty');
        const table = document.getElementById('imagesTable');
        
        if (this.data.images.length === 0) {
            table.style.display = 'none';
            emptyState.style.display = 'flex';
            return;
        }
        
        table.style.display = 'table';
        emptyState.style.display = 'none';
        
        tbody.innerHTML = this.data.images.map(image => `
            <tr data-id="${image.id}">
                <td><strong>${image.repository}</strong></td>
                <td>${image.tag}</td>
                <td>
                    <small class="text-muted">${image.id}</small>
                </td>
                <td>
                    <small class="text-muted">
                        ${this.formatDate(image.created)}
                    </small>
                </td>
                <td>${this.formatBytes(image.size)}</td>
                <td>
                    <div class="action-buttons">
                        <button class="action-btn remove" onclick="dockerGUI.removeImage('${image.id}')" title="Remove">
                            <i class="fas fa-trash"></i>
                        </button>
                    </div>
                </td>
            </tr>
        `).join('');
    }
    
    renderVolumes() {
        const tbody = document.getElementById('volumesTableBody');
        const emptyState = document.getElementById('volumesEmpty');
        const table = document.getElementById('volumesTable');
        
        if (this.data.volumes.length === 0) {
            table.style.display = 'none';
            emptyState.style.display = 'flex';
            return;
        }
        
        table.style.display = 'table';
        emptyState.style.display = 'none';
        
        tbody.innerHTML = this.data.volumes.map(volume => `
            <tr data-name="${volume.name}">
                <td><strong>${volume.name}</strong></td>
                <td>${volume.driver}</td>
                <td>
                    <small class="text-muted">${volume.mountpoint}</small>
                </td>
                <td>
                    <small class="text-muted">
                        ${this.formatDate(volume.created)}
                    </small>
                </td>
                <td>
                    <div class="action-buttons">
                        <button class="action-btn remove" onclick="dockerGUI.removeVolume('${volume.name}')" title="Remove">
                            <i class="fas fa-trash"></i>
                        </button>
                    </div>
                </td>
            </tr>
        `).join('');
    }
    
    updateCounts() {
        document.getElementById('containerCount').textContent = this.data.containers.length;
        document.getElementById('imageCount').textContent = this.data.images.length;
        document.getElementById('volumeCount').textContent = this.data.volumes.length;
    }
    
    filterTable(type, searchTerm) {
        const tableBody = document.getElementById(`${type}sTableBody`);
        const rows = tableBody.querySelectorAll('tr');
        
        rows.forEach(row => {
            const text = row.textContent.toLowerCase();
            const match = text.includes(searchTerm.toLowerCase());
            row.style.display = match ? '' : 'none';
        });
    }
    
    // Container Actions
    async startContainer(containerId) {
        try {
            const response = await fetch(`${this.apiBase}/api/containers/${containerId}/start`, {
                method: 'POST'
            });
            
            if (response.ok) {
                this.showToast('Container started successfully', 'success');
                await this.loadContainers();
            } else {
                const error = await response.json();
                this.showToast(error.error || 'Failed to start container', 'error');
            }
        } catch (error) {
            this.showToast('Error starting container', 'error');
        }
    }
    
    async stopContainer(containerId) {
        try {
            const response = await fetch(`${this.apiBase}/api/containers/${containerId}/stop`, {
                method: 'POST'
            });
            
            if (response.ok) {
                this.showToast('Container stopped successfully', 'success');
                await this.loadContainers();
            } else {
                const error = await response.json();
                this.showToast(error.error || 'Failed to stop container', 'error');
            }
        } catch (error) {
            this.showToast('Error stopping container', 'error');
        }
    }
    
    async restartContainer(containerId) {
        try {
            const response = await fetch(`${this.apiBase}/api/containers/${containerId}/restart`, {
                method: 'POST'
            });
            
            if (response.ok) {
                this.showToast('Container restarted successfully', 'success');
                await this.loadContainers();
            } else {
                const error = await response.json();
                this.showToast(error.error || 'Failed to restart container', 'error');
            }
        } catch (error) {
            this.showToast('Error restarting container', 'error');
        }
    }
    
    async removeContainer(containerId) {
        if (!confirm('Are you sure you want to remove this container?')) {
            return;
        }
        
        try {
            const response = await fetch(`${this.apiBase}/api/containers/${containerId}/remove`, {
                method: 'DELETE'
            });
            
            if (response.ok) {
                this.showToast('Container removed successfully', 'success');
                await this.loadContainers();
            } else {
                const error = await response.json();
                this.showToast(error.error || 'Failed to remove container', 'error');
            }
        } catch (error) {
            this.showToast('Error removing container', 'error');
        }
    }
    
    // Image Actions
    async pullImage(imageName) {
        try {
            this.showLoading();
            
            const response = await fetch(`${this.apiBase}/api/images/pull`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ image: imageName })
            });
            
            if (response.ok) {
                this.showToast(`Image ${imageName} pulled successfully`, 'success');
                await this.loadImages();
            } else {
                const error = await response.json();
                this.showToast(error.error || 'Failed to pull image', 'error');
            }
        } catch (error) {
            this.showToast('Error pulling image', 'error');
        } finally {
            this.hideLoading();
        }
    }
    
    async removeImage(imageId) {
        if (!confirm('Are you sure you want to remove this image?')) {
            return;
        }
        
        try {
            const response = await fetch(`${this.apiBase}/api/images/${imageId}/remove`, {
                method: 'DELETE'
            });
            
            if (response.ok) {
                this.showToast('Image removed successfully', 'success');
                await this.loadImages();
            } else {
                const error = await response.json();
                this.showToast(error.error || 'Failed to remove image', 'error');
            }
        } catch (error) {
            this.showToast('Error removing image', 'error');
        }
    }
    
    // Volume Actions
    async createVolume(volumeName) {
        try {
            const response = await fetch(`${this.apiBase}/api/volumes/create`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ name: volumeName })
            });
            
            if (response.ok) {
                this.showToast(`Volume ${volumeName} created successfully`, 'success');
                await this.loadVolumes();
            } else {
                const error = await response.json();
                this.showToast(error.error || 'Failed to create volume', 'error');
            }
        } catch (error) {
            this.showToast('Error creating volume', 'error');
        }
    }
    
    async removeVolume(volumeName) {
        if (!confirm('Are you sure you want to remove this volume?')) {
            return;
        }
        
        try {
            const response = await fetch(`${this.apiBase}/api/volumes/${volumeName}/remove`, {
                method: 'DELETE'
            });
            
            if (response.ok) {
                this.showToast('Volume removed successfully', 'success');
                await this.loadVolumes();
            } else {
                const error = await response.json();
                this.showToast(error.error || 'Failed to remove volume', 'error');
            }
        } catch (error) {
            this.showToast('Error removing volume', 'error');
        }
    }
    
    // Utility Methods
    showToast(message, type = 'info') {
        const container = document.getElementById('toastContainer');
        const toast = document.createElement('div');
        toast.className = `toast ${type}`;
        
        const icon = this.getToastIcon(type);
        toast.innerHTML = `
            <i class="fas ${icon}"></i>
            <span>${message}</span>
        `;
        
        container.appendChild(toast);
        
        setTimeout(() => {
            toast.remove();
        }, 5000);
    }
    
    getToastIcon(type) {
        const icons = {
            success: 'fa-check-circle',
            error: 'fa-exclamation-circle',
            warning: 'fa-exclamation-triangle',
            info: 'fa-info-circle'
        };
        return icons[type] || icons.info;
    }
    
    showLoading() {
        document.getElementById('loadingOverlay').style.display = 'flex';
    }
    
    hideLoading() {
        document.getElementById('loadingOverlay').style.display = 'none';
    }
    
    formatDate(dateString) {
        const date = new Date(dateString);
        return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
    }
    
    formatBytes(bytes) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }
    
    startAutoRefresh() {
        // Refresh every 30 seconds
        this.refreshInterval = setInterval(() => {
            this.loadData();
            this.checkDockerConnection();
        }, 30000);
    }
    
    stopAutoRefresh() {
        if (this.refreshInterval) {
            clearInterval(this.refreshInterval);
            this.refreshInterval = null;
        }
    }

    // Container Details View
    showContainerDetails(containerId) {
        console.log('=== showContainerDetails called ===');
        console.log('Container ID:', containerId);
        
        // Find the container in our data
        const container = this.data.containers.find(c => c.id === containerId);
        console.log('Container data found:', !!container);
        console.log('Container details:', container);
        
        if (!container) {
            console.error('Container not found:', containerId);
            return;
        }
        
        // Check if container is running (for some functionality)
        const isRunning = (container.state === 'running' || container.status === 'running');
        console.log('Container is running:', isRunning);
        
        // Hide the containers list view
        document.getElementById('containersView').style.display = 'none';
        
        // Show the details view
        const detailsView = document.getElementById('containerDetailsView');
        detailsView.style.display = 'block';
        
        // Set the title
        document.getElementById('detailsContainerTitle').textContent = `${container.name} Details`;
        
        // Populate overview data
        document.getElementById('detailsContainerName').textContent = container.name;
        document.getElementById('detailsContainerId').textContent = container.id;
        document.getElementById('detailsContainerImage').textContent = container.image;
        document.getElementById('detailsContainerCreated').textContent = this.formatDate(container.created);
        
        const statusElement = document.getElementById('detailsContainerStatus');
        const status = container.status || container.state || 'unknown';
        statusElement.textContent = status;
        statusElement.className = `container-status status-${status.toLowerCase()}`;
        
        const portsElement = document.getElementById('detailsContainerPorts');
        if (container.ports && container.ports.length > 0) {
            portsElement.textContent = container.ports.join(', ');
        } else {
            portsElement.textContent = 'No ports exposed';
        }
        
        // Store current container ID for tab operations
        this.currentContainerId = containerId;
        
        // Set up action buttons
        this.setupContainerDetailsActions(containerId, isRunning);
        
        // Set up tab functionality
        this.setupContainerDetailsTabs(containerId);
        
        // Set up back button
        document.getElementById('backToContainers').onclick = () => this.hideContainerDetails();
        
        // Show overview tab by default
        this.showContainerTab('overview');
        
        console.log('Container details view displayed');
    }
    
    hideContainerDetails() {
        // Hide the details view
        document.getElementById('containerDetailsView').style.display = 'none';
        
        // Show the containers list view
        document.getElementById('containersView').style.display = 'block';
        
        // Clear current container ID
        this.currentContainerId = null;
        
        console.log('Returned to containers list');
    }

    setupContainerDetailsActions(containerId, isRunning) {
        const startBtn = document.getElementById('detailsStartBtn');
        const stopBtn = document.getElementById('detailsStopBtn');
        const restartBtn = document.getElementById('detailsRestartBtn');
        const removeBtn = document.getElementById('detailsRemoveBtn');

        // Show/hide buttons based on container state
        startBtn.style.display = isRunning ? 'none' : 'inline-block';
        stopBtn.style.display = isRunning ? 'inline-block' : 'none';

        startBtn.onclick = () => {
            this.startContainer(containerId);
            setTimeout(() => this.showContainerDetails(containerId), 1000); // Refresh details
        };
        stopBtn.onclick = () => {
            this.stopContainer(containerId);
            setTimeout(() => this.showContainerDetails(containerId), 1000); // Refresh details
        };
        restartBtn.onclick = () => {
            this.restartContainer(containerId);
            setTimeout(() => this.showContainerDetails(containerId), 1000); // Refresh details
        };
        removeBtn.onclick = () => {
            if (confirm('Are you sure you want to remove this container?')) {
                this.removeContainer(containerId);
                this.hideContainerDetails(); // Go back to list
            }
        };
    }

    setupContainerDetailsTabs(containerId) {
        // Tab navigation
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const tabName = e.currentTarget.dataset.tab;
                this.showContainerTab(tabName);
            });
        });
    }
    
    showContainerTab(tabName) {
        // Remove active class from all tabs and content
        document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
        document.querySelectorAll('.tab-content').forEach(content => content.classList.remove('active'));
        
        // Add active class to selected tab and content
        document.querySelector(`[data-tab="${tabName}"]`).classList.add('active');
        document.getElementById(`${tabName}-tab`).classList.add('active');
        
        // Load tab-specific content
        if (this.currentContainerId) {
            this.loadTabContent(tabName, this.currentContainerId);
        }
    }
    
    async loadTabContent(tabName, containerId) {
        try {
            switch (tabName) {
                case 'logs':
                    const logsResponse = await fetch(`${this.apiBase}/api/containers/${containerId}/logs`);
                    const logs = await logsResponse.text();
                    document.getElementById('containerLogs').textContent = logs || 'No logs available';
                    break;
                    
                case 'files':
                    const filesResponse = await fetch(`${this.apiBase}/api/containers/${containerId}/files`);
                    const files = await filesResponse.json();
                    this.renderContainerFiles(files);
                    break;
                    
                case 'env':
                    const envResponse = await fetch(`${this.apiBase}/api/containers/${containerId}/env`);
                    const env = await envResponse.json();
                    this.renderContainerEnv(env);
                    break;
                    
                case 'volumes':
                    this.renderContainerVolumes(containerId);
                    break;
                    
                case 'network':
                    this.renderContainerNetwork(containerId);
                    break;
            }
        } catch (error) {
            console.error(`Error loading ${tabName} content:`, error);
            const contentElement = document.getElementById(`container${tabName.charAt(0).toUpperCase() + tabName.slice(1)}`);
            if (contentElement) {
                contentElement.textContent = `Error loading ${tabName} information`;
            }
        }
    }
    
    renderContainerFiles(files) {
        const filesContainer = document.getElementById('containerFiles');
        if (!files || files.length === 0) {
            filesContainer.innerHTML = '<p>No files found or unable to access container filesystem</p>';
            return;
        }
        
        filesContainer.innerHTML = files.map(file => `
            <div class="file-item">
                <i class="fas fa-${file.type === 'directory' ? 'folder' : 'file'}"></i>
                <span class="file-name">${file.name}</span>
                <span class="file-size">${file.size || '-'}</span>
            </div>
        `).join('');
    }
    
    renderContainerEnv(env) {
        const envContainer = document.getElementById('containerEnv');
        if (!env || Object.keys(env).length === 0) {
            envContainer.innerHTML = '<p>No environment variables found</p>';
            return;
        }
        
        envContainer.innerHTML = Object.entries(env).map(([key, value]) => `
            <div class="env-item">
                <strong>${key}:</strong> <span>${value}</span>
            </div>
        `).join('');
    }
    
    renderContainerVolumes(containerId) {
        const volumesContainer = document.getElementById('containerVolumes');
        const container = this.data.containers.find(c => c.id === containerId);
        
        if (!container || !container.volumes || container.volumes.length === 0) {
            volumesContainer.innerHTML = '<p>No volumes mounted</p>';
            return;
        }
        
        volumesContainer.innerHTML = container.volumes.map(volume => `
            <div class="volume-item">
                <strong>Host:</strong> ${volume.host}<br>
                <strong>Container:</strong> ${volume.container}<br>
                <strong>Mode:</strong> ${volume.mode || 'rw'}
            </div>
        `).join('');
    }
    
    renderContainerNetwork(containerId) {
        const networkContainer = document.getElementById('containerNetwork');
        const container = this.data.containers.find(c => c.id === containerId);
        
        if (!container || !container.networks) {
            networkContainer.innerHTML = '<p>No network information available</p>';
            return;
        }
        
        networkContainer.innerHTML = `
            <div class="network-info">
                <strong>Networks:</strong> ${container.networks.join(', ')}<br>
                <strong>Ports:</strong> ${container.ports && container.ports.length > 0 ? container.ports.join(', ') : 'None'}
            </div>
        `;
    }

    async loadContainerLogs(containerId) {
        try {
            const response = await fetch(`${this.apiBase}/api/containers/${containerId}/logs`);
            const data = await response.json();
            
            if (response.ok) {
                document.getElementById('logsContent').textContent = data.logs || 'No logs available';
            } else {
                document.getElementById('logsContent').textContent = `Error: ${data.error}`;
            }
        } catch (error) {
            document.getElementById('logsContent').textContent = `Error loading logs: ${error.message}`;
        }

        // Set up logs controls
        document.getElementById('refreshLogsBtn').onclick = () => this.loadContainerLogs(containerId);
        document.getElementById('clearLogsBtn').onclick = () => {
            document.getElementById('logsContent').textContent = '';
        };
    }

    async loadContainerFiles(containerId, path = '/') {
        try {
            const response = await fetch(`${this.apiBase}/api/containers/${containerId}/files?path=${encodeURIComponent(path)}`);
            const data = await response.json();
            
            if (response.ok) {
                this.renderFilesList(data.files || [], path, containerId);
                this.updateFilesBreadcrumb(path, containerId);
            } else {
                document.getElementById('filesList').innerHTML = `<div class="error">Error: ${data.error}</div>`;
            }
        } catch (error) {
            document.getElementById('filesList').innerHTML = `<div class="error">Error loading files: ${error.message}</div>`;
        }

        // Set up files controls
        document.getElementById('refreshFilesBtn').onclick = () => this.loadContainerFiles(containerId, path);
    }

    renderFilesList(files, currentPath, containerId) {
        const filesList = document.getElementById('filesList');
        
        if (files.length === 0) {
            filesList.innerHTML = '<div class="empty">No files found</div>';
            return;
        }

        const html = files.map(file => `
            <div class="file-item" onclick="dockerGUI.handleFileClick('${file.name}', ${file.is_directory}, '${currentPath}', '${containerId}')">
                <i class="file-icon fas ${file.is_directory ? 'fa-folder' : 'fa-file'}"></i>
                <span class="file-name">${file.name}</span>
                <span class="file-size">${file.is_directory ? '' : this.formatFileSize(file.size)}</span>
            </div>
        `).join('');

        filesList.innerHTML = html;
    }

    handleFileClick(fileName, isDirectory, currentPath, containerId) {
        if (isDirectory) {
            const newPath = currentPath === '/' ? `/${fileName}` : `${currentPath}/${fileName}`;
            this.loadContainerFiles(containerId, newPath);
        }
    }

    updateFilesBreadcrumb(path, containerId) {
        const breadcrumb = document.getElementById('filesBreadcrumb');
        const parts = path.split('/').filter(p => p);
        
        let html = '<span class="breadcrumb-item" onclick="dockerGUI.loadContainerFiles(\'' + containerId + '\', \'/\')">/</span>';
        
        let currentPath = '';
        parts.forEach(part => {
            currentPath += '/' + part;
            html += ` <span class="breadcrumb-item" onclick="dockerGUI.loadContainerFiles('${containerId}', '${currentPath}')">${part}</span>`;
        });
        
        breadcrumb.innerHTML = html;
    }

    setupContainerTerminal(containerId) {
        const connectBtn = document.getElementById('connectTerminalBtn');
        const disconnectBtn = document.getElementById('disconnectTerminalBtn');
        const terminalInput = document.getElementById('terminalInput');
        const sendBtn = document.getElementById('sendCommandBtn');
        const output = document.getElementById('terminalOutput');

        connectBtn.onclick = () => {
            // Simulate terminal connection
            output.innerHTML = `<div class="terminal-line">Connected to container ${containerId}</div>
                               <div class="terminal-line">Type commands and press Enter or click Send</div>
                               <div class="terminal-prompt">root@${containerId.substring(0, 12)}:/#</div>`;
            
            connectBtn.disabled = true;
            disconnectBtn.disabled = false;
            terminalInput.disabled = false;
            sendBtn.disabled = false;
            terminalInput.focus();
        };

        disconnectBtn.onclick = () => {
            output.innerHTML = '<div class="terminal-welcome">Click "Connect" to start a terminal session in the container.</div>';
            connectBtn.disabled = false;
            disconnectBtn.disabled = true;
            terminalInput.disabled = true;
            sendBtn.disabled = true;
        };

        const sendCommand = async () => {
            const command = terminalInput.value.trim();
            if (!command) return;

            try {
                const response = await fetch(`${this.apiBase}/api/containers/${containerId}/exec`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ command })
                });
                
                const data = await response.json();
                
                output.innerHTML += `<div class="terminal-line">$ ${command}</div>`;
                if (response.ok) {
                    output.innerHTML += `<div class="terminal-output">${data.output}</div>`;
                } else {
                    output.innerHTML += `<div class="terminal-error">Error: ${data.error}</div>`;
                }
                output.innerHTML += `<div class="terminal-prompt">root@${containerId.substring(0, 12)}:/#</div>`;
                
                terminalInput.value = '';
                output.scrollTop = output.scrollHeight;
            } catch (error) {
                output.innerHTML += `<div class="terminal-error">Error: ${error.message}</div>`;
            }
        };

        sendBtn.onclick = sendCommand;
        terminalInput.onkeypress = (e) => {
            if (e.key === 'Enter') {
                sendCommand();
            }
        };
    }

    async loadContainerEnvironment(containerId) {
        try {
            const response = await fetch(`${this.apiBase}/api/containers/${containerId}/env`);
            const data = await response.json();
            
            if (response.ok) {
                this.renderEnvironmentTable(data.environment || []);
            } else {
                document.getElementById('envTableBody').innerHTML = `<tr><td colspan="2" class="error">Error: ${data.error}</td></tr>`;
            }
        } catch (error) {
            document.getElementById('envTableBody').innerHTML = `<tr><td colspan="2" class="error">Error loading environment: ${error.message}</td></tr>`;
        }

        // Set up environment controls
        document.getElementById('refreshEnvBtn').onclick = () => this.loadContainerEnvironment(containerId);
        
        const searchInput = document.getElementById('envSearch');
        searchInput.oninput = () => this.filterEnvironmentTable(searchInput.value);
    }

    renderEnvironmentTable(envVars) {
        const tbody = document.getElementById('envTableBody');
        
        if (envVars.length === 0) {
            tbody.innerHTML = '<tr><td colspan="2" class="empty">No environment variables found</td></tr>';
            return;
        }

        const html = envVars.map(env => `
            <tr class="env-row">
                <td>${env.key}</td>
                <td>${env.value}</td>
            </tr>
        `).join('');

        tbody.innerHTML = html;
    }

    filterEnvironmentTable(query) {
        const rows = document.querySelectorAll('.env-row');
        const lowerQuery = query.toLowerCase();
        
        rows.forEach(row => {
            const key = row.cells[0].textContent.toLowerCase();
            const value = row.cells[1].textContent.toLowerCase();
            const matches = key.includes(lowerQuery) || value.includes(lowerQuery);
            row.style.display = matches ? '' : 'none';
        });
    }

    async loadContainerVolumes(containerId) {
        try {
            const response = await fetch(`${this.apiBase}/api/containers/${containerId}/details`);
            const data = await response.json();
            
            if (response.ok && data.mounts) {
                this.renderVolumesTable(data.mounts);
            } else {
                document.getElementById('volumesTableBody').innerHTML = '<tr><td colspan="4" class="empty">No volume mounts found</td></tr>';
            }
        } catch (error) {
            document.getElementById('volumesTableBody').innerHTML = `<tr><td colspan="4" class="error">Error loading volumes: ${error.message}</td></tr>`;
        }
    }

    renderVolumesTable(mounts) {
        const tbody = document.getElementById('volumesTableBody');
        
        if (mounts.length === 0) {
            tbody.innerHTML = '<tr><td colspan="4" class="empty">No volume mounts found</td></tr>';
            return;
        }

        const html = mounts.map(mount => `
            <tr>
                <td>${mount.source || '-'}</td>
                <td>${mount.destination || '-'}</td>
                <td>${mount.mode || 'rw'}</td>
                <td>${mount.type || 'bind'}</td>
            </tr>
        `).join('');

        tbody.innerHTML = html;
    }

    async loadContainerNetwork(containerId) {
        try {
            const response = await fetch(`${this.apiBase}/api/containers/${containerId}/details`);
            const data = await response.json();
            
            if (response.ok && data.network_settings) {
                this.renderNetworkInfo(data.network_settings);
            } else {
                this.renderNetworkInfo({});
            }
        } catch (error) {
            this.renderNetworkInfo({});
        }
    }

    renderNetworkInfo(networkSettings) {
        document.getElementById('networkMode').textContent = 'bridge';
        document.getElementById('networkIP').textContent = networkSettings.ip_address || '-';
        document.getElementById('networkGateway').textContent = networkSettings.gateway || '-';
        document.getElementById('networkMAC').textContent = networkSettings.mac_address || '-';

        const portsTableBody = document.getElementById('portsTableBody');
        
        if (networkSettings.ports && networkSettings.ports.length > 0) {
            const html = networkSettings.ports.map(port => `
                <tr>
                    <td>${port.container_port}</td>
                    <td>${port.host_port}</td>
                    <td>${port.protocol}</td>
                </tr>
            `).join('');
            portsTableBody.innerHTML = html;
        } else {
            portsTableBody.innerHTML = '<tr><td colspan="3" class="empty">No port bindings found</td></tr>';
        }
    }

    formatFileSize(bytes) {
        if (bytes === 0) return '0 B';
        const sizes = ['B', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(1024));
        return Math.round(bytes / Math.pow(1024, i) * 100) / 100 + ' ' + sizes[i];
    }

    // Container Details Methods
    async showContainerDetails(containerId) {
        console.log('Showing details for container:', containerId);
        this.currentContainerId = containerId;
        
        try {
            // Hide container list, show details view
            document.getElementById('containersList').style.display = 'none';
            document.getElementById('containerDetails').style.display = 'block';
            
            // Find container data
            const container = this.data.containers.find(c => c.id === containerId);
            if (!container) {
                console.error('Container not found:', containerId);
                return;
            }
            
            // Update overview information
            document.getElementById('containerDetailsTitle').textContent = 'Container Details';
            document.getElementById('detailsContainerId').textContent = container.id;
            document.getElementById('detailsContainerName').textContent = container.name;
            
            // Set up tab event listeners
            this.setupTabEventListeners();
            
            // Load default tab (logs)
            this.switchTab('logs');
            
        } catch (error) {
            console.error('Error showing container details:', error);
        }
    }
    
    hideContainerDetails() {
        console.log('Hiding container details');
        document.getElementById('containerDetails').style.display = 'none';
        document.getElementById('containersList').style.display = 'block';
        this.currentContainerId = null;
    }
    
    setupTabEventListeners() {
        // Remove existing listeners to prevent duplicates
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.replaceWith(btn.cloneNode(true));
        });
        
        // Add new event listeners
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.preventDefault();
                const tabName = btn.getAttribute('data-tab');
                this.switchTab(tabName);
            });
        });
        
        // Also set up back button
        const backBtn = document.getElementById('backToContainers');
        if (backBtn) {
            backBtn.addEventListener('click', () => {
                this.hideContainerDetails();
            });
        }
    }
    
    switchTab(tabName) {
        console.log('Switching to tab:', tabName);
        
        // Update tab buttons
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.classList.remove('active');
        });
        document.querySelector(`[data-tab="${tabName}"]`).classList.add('active');
        
        // Update tab content
        document.querySelectorAll('.tab-pane').forEach(pane => {
            pane.classList.remove('active');
        });
        document.getElementById(`${tabName}Tab`).classList.add('active');
        
        // Load tab-specific content
        this.loadTabContent(tabName);
    }
    
    async loadTabContent(tabName) {
        if (!this.currentContainerId) return;
        
        try {
            switch (tabName) {
                case 'logs':
                    await this.loadContainerLogs();
                    break;
                case 'files':
                    await this.loadContainerFiles('/');
                    break;
                case 'env':
                    await this.loadContainerEnvironment();
                    break;
                case 'exec':
                    this.setupContainerExec();
                    break;
                case 'volumes':
                    await this.loadContainerVolumes();
                    break;
                case 'network':
                    await this.loadContainerNetwork();
                    break;
                case 'stats':
                    await this.loadContainerStats();
                    break;
            }
        } catch (error) {
            console.error(`Error loading ${tabName} content:`, error);
        }
    }
    
    async loadContainerLogs() {
        const logsContent = document.getElementById('logsContent');
        logsContent.innerHTML = '<div class="loading">Loading logs...</div>';
        
        try {
            const response = await fetch(`${this.apiBase}/api/containers/${this.currentContainerId}/logs`);
            if (response.ok) {
                const logs = await response.text();
                logsContent.innerHTML = `<pre class="logs-text">${logs || 'No logs available'}</pre>`;
            } else {
                logsContent.innerHTML = '<div class="error">Error loading logs</div>';
            }
        } catch (error) {
            console.error('Error loading logs:', error);
            logsContent.innerHTML = '<div class="error">Error loading logs</div>';
        }
    }
    
    async loadContainerFiles(path = '/') {
        const filesContent = document.getElementById('filesContent');
        const currentPath = document.getElementById('currentPath');
        
        filesContent.innerHTML = '<div class="loading">Loading files...</div>';
        currentPath.textContent = path;
        
        try {
            const response = await fetch(`${this.apiBase}/api/containers/${this.currentContainerId}/files?path=${encodeURIComponent(path)}`);
            if (response.ok) {
                const files = await response.json();
                this.renderFiles(files, path);
            } else {
                filesContent.innerHTML = '<div class="error">Error loading files</div>';
            }
        } catch (error) {
            console.error('Error loading files:', error);
            filesContent.innerHTML = '<div class="error">Error loading files</div>';
        }
    }
    
    renderFiles(files, currentPath) {
        const filesContent = document.getElementById('filesContent');
        
        if (!files || files.length === 0) {
            filesContent.innerHTML = '<div class="empty">No files found</div>';
            return;
        }
        
        let html = '<div class="files-list">';
        
        // Add parent directory link if not root
        if (currentPath !== '/') {
            const parentPath = currentPath.split('/').slice(0, -1).join('/') || '/';
            html += `<div class="file-item directory" onclick="dockerGUI.loadContainerFiles('${parentPath}')">
                <i class="fas fa-level-up-alt"></i>
                <span>..</span>
            </div>`;
        }
        
        files.forEach(file => {
            const icon = file.is_dir ? 'fa-folder' : 'fa-file';
            const fileClass = file.is_dir ? 'directory' : 'file';
            const onclick = file.is_dir ? `dockerGUI.loadContainerFiles('${file.path}')` : '';
            
            html += `<div class="file-item ${fileClass}" ${onclick ? `onclick="${onclick}"` : ''}>
                <i class="fas ${icon}"></i>
                <span>${file.name}</span>
                <span class="file-size">${file.is_dir ? '' : this.formatFileSize(file.size)}</span>
            </div>`;
        });
        
        html += '</div>';
        filesContent.innerHTML = html;
    }
    
    async loadContainerEnvironment() {
        const envContent = document.getElementById('envContent');
        envContent.innerHTML = '<div class="loading">Loading environment variables...</div>';
        
        try {
            const response = await fetch(`${this.apiBase}/api/containers/${this.currentContainerId}/env`);
            if (response.ok) {
                const envVars = await response.json();
                this.renderEnvironmentVariables(envVars);
            } else {
                envContent.innerHTML = '<div class="error">Error loading environment variables</div>';
            }
        } catch (error) {
            console.error('Error loading environment:', error);
            envContent.innerHTML = '<div class="error">Error loading environment variables</div>';
        }
    }
    
    renderEnvironmentVariables(envVars) {
        const envContent = document.getElementById('envContent');
        
        if (!envVars || Object.keys(envVars).length === 0) {
            envContent.innerHTML = '<div class="empty">No environment variables found</div>';
            return;
        }
        
        let html = '<div class="env-list">';
        Object.entries(envVars).forEach(([key, value]) => {
            html += `<div class="env-item">
                <div class="env-key">${key}</div>
                <div class="env-value">${value}</div>
            </div>`;
        });
        html += '</div>';
        
        envContent.innerHTML = html;
    }
    
    setupContainerExec() {
        const execTerminal = document.getElementById('execTerminal');
        execTerminal.innerHTML = `
            <div class="terminal-placeholder">
                <p>Terminal functionality not yet implemented</p>
                <p>Container: ${this.currentContainerId}</p>
            </div>
        `;
    }
    
    async loadContainerVolumes() {
        const volumesContent = document.getElementById('volumesContent');
        volumesContent.innerHTML = '<div class="loading">Loading volume information...</div>';
        
        // For now, show placeholder
        volumesContent.innerHTML = '<div class="placeholder">Volume information not yet implemented</div>';
    }
    
    async loadContainerNetwork() {
        const networkContent = document.getElementById('networkContent');
        networkContent.innerHTML = '<div class="loading">Loading network information...</div>';
        
        // For now, show placeholder
        networkContent.innerHTML = '<div class="placeholder">Network information not yet implemented</div>';
    }
    
    async loadContainerStats() {
        const statsContent = document.getElementById('statsContent');
        statsContent.innerHTML = '<div class="loading">Loading container statistics...</div>';
        
        // For now, show placeholder
        statsContent.innerHTML = '<div class="placeholder">Statistics not yet implemented</div>';
    }

    // Utility Methods
    formatDate(dateString) {
        if (!dateString || dateString === '-') return '-';
        try {
            const date = new Date(dateString);
            return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
        } catch (error) {
            return dateString;
        }
    }

    showToast(message, type = 'info') {
        console.log(`Toast [${type}]: ${message}`);
        // For now, just log to console. Could implement actual toast notifications later.
    }

    updateCounts() {
        // Update sidebar counts if needed
        const runningContainers = this.data.containers.filter(c => c.state === 'running' || c.status === 'running').length;
        console.log(`Running containers: ${runningContainers}/${this.data.containers.length}`);
    }

    showLoading() {
        console.log('Loading data...');
        // Could add loading spinner here
    }

    hideLoading() {
        console.log('Data loaded');
        // Could hide loading spinner here
    }

    startAutoRefresh() {
        // Auto-refresh every 30 seconds
        this.refreshInterval = setInterval(() => {
            this.loadData();
        }, 30000);
    }

    stopAutoRefresh() {
        if (this.refreshInterval) {
            clearInterval(this.refreshInterval);
            this.refreshInterval = null;
        }
    }

    // Container Action Methods (placeholders)
    async startContainer(containerId) {
        console.log('Starting container:', containerId);
        this.showToast('Container start functionality not yet implemented', 'info');
    }

    async stopContainer(containerId) {
        console.log('Stopping container:', containerId);
        this.showToast('Container stop functionality not yet implemented', 'info');
    }

    async restartContainer(containerId) {
        console.log('Restarting container:', containerId);
        this.showToast('Container restart functionality not yet implemented', 'info');
    }

    async removeContainer(containerId) {
        console.log('Removing container:', containerId);
        this.showToast('Container remove functionality not yet implemented', 'info');
    }
}

// Initialize the application
let dockerGUI;
document.addEventListener('DOMContentLoaded', () => {
    console.log('DOM loaded, initializing Docker GUI...');
    dockerGUI = new DockerGUI();
    window.dockerGUI = dockerGUI;  // Make it globally accessible
    
    console.log('Docker GUI initialized');
    console.log('dockerGUI available globally:', !!window.dockerGUI);
});

// Handle page unload
window.addEventListener('beforeunload', () => {
    if (dockerGUI) {
        dockerGUI.stopAutoRefresh();
    }
});
