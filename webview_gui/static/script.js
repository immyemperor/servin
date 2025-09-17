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
        this.socket = null;
        this.isLogsStreaming = false;
        this.isExecConnected = false;
        this.data = {
            containers: [],
            images: [],
            volumes: []
        };
        
        this.init();
    }
    
    init() {
        this.setupEventListeners();
        this.initializeSocket();
        this.checkDockerConnection();
        this.startAutoRefresh();
        this.loadData();
    }
    
    initializeSocket() {
        // Initialize Socket.IO connection
        this.socket = io();
        
        this.socket.on('connect', () => {
            console.log('Connected to server');
        });
        
        this.socket.on('disconnect', () => {
            console.log('Disconnected from server');
            this.isLogsStreaming = false;
            this.isExecConnected = false;
        });
        
        // Log streaming events
        this.socket.on('log_data', (data) => {
            this.handleLogData(data);
        });
        
        this.socket.on('logs_started', (data) => {
            console.log('Log streaming started for container:', data.container_id);
            this.isLogsStreaming = true;
            this.updateLogsUI();
        });
        
        this.socket.on('logs_stopped', (data) => {
            console.log('Log streaming stopped for container:', data.container_id);
            this.isLogsStreaming = false;
            this.updateLogsUI();
        });
        
        // Exec session events
        this.socket.on('exec_started', (data) => {
            this.handleExecStarted(data);
        });
        
        this.socket.on('exec_stopped', (data) => {
            this.handleExecStopped(data);
        });
        
        this.socket.on('exec_output', (data) => {
            this.handleExecOutput(data);
        });
        
        this.socket.on('error', (data) => {
            console.error('Socket error:', data.message);
            this.showError(data.message);
        });
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
            
            console.log(`Container ${container.name}: isRunning=${isRunning}, status=${status}, state=${container.state}, created=${container.created}`);
            console.log(`Formatted date for ${container.created}:`, this.formatDate(container.created));
            
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
                <td>
                    <small class="text-muted">
                        ${this.formatDate(container.created)}
                    </small>
                </td>
                <td>${container.ports && container.ports.length > 0 ? 
                    container.ports.map(p => `${p.host_port || p.container_port}:${p.container_port}/${p.protocol || 'tcp'}`).join(', ') 
                    : '<span class="text-muted">No ports</span>'}</td>
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
        // Handle relative time strings like "2 days ago", "14 hours ago"
        if (!dateString || dateString === '-' || dateString === 'unknown') {
            return '-';
        }
        
        // If it's already a relative time string (contains "ago"), return as-is
        if (typeof dateString === 'string' && dateString.includes('ago')) {
            return dateString;
        }
        
        // Try to parse as a regular date
        try {
            const date = new Date(dateString);
            if (isNaN(date.getTime())) {
                return dateString; // Return original string if not a valid date
            }
            return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
        } catch (error) {
            return dateString; // Return original string if parsing fails
        }
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

    // Container Details Implementation starts here

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
        console.log('Setting up tab event listeners');
        
        // Remove existing listeners by replacing elements
        document.querySelectorAll('.tab-btn').forEach(btn => {
            const newBtn = btn.cloneNode(true);
            btn.parentNode.replaceChild(newBtn, btn);
        });
        
        // Add new event listeners with proper context binding
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.preventDefault();
                const tabName = btn.getAttribute('data-tab');
                console.log('Tab clicked:', tabName);
                this.switchTab(tabName);
            });
        });
        
        // Also set up back button
        const backBtn = document.getElementById('backToContainers');
        if (backBtn) {
            // Remove existing listener
            const newBackBtn = backBtn.cloneNode(true);
            backBtn.parentNode.replaceChild(newBackBtn, backBtn);
            
            // Add new listener
            document.getElementById('backToContainers').addEventListener('click', () => {
                this.hideContainerDetails();
            });
        }
    }
    
    switchTab(tabName) {
        console.log('Switching to tab:', tabName);
        
        // Ensure we have valid elements
        const allTabBtns = document.querySelectorAll('.tab-btn');
        const allTabPanes = document.querySelectorAll('.tab-pane');
        
        console.log('Found tab buttons:', allTabBtns.length);
        console.log('Found tab panes:', allTabPanes.length);
        
        // Update tab buttons
        allTabBtns.forEach(btn => {
            btn.classList.remove('active');
        });
        
        const targetTab = document.querySelector(`[data-tab="${tabName}"]`);
        if (targetTab) {
            targetTab.classList.add('active');
            console.log('Tab button activated:', tabName);
        } else {
            console.error('Tab button not found:', tabName);
            return; // Exit early if tab button not found
        }
        
        // Update tab content - hide all first
        allTabPanes.forEach(pane => {
            pane.classList.remove('active');
            pane.style.display = 'none'; // Force hide
            console.log('Hidden pane:', pane.id);
        });
        
        // Show target pane
        const targetPane = document.getElementById(`${tabName}Tab`);
        if (targetPane) {
            targetPane.classList.add('active');
            targetPane.style.display = 'block'; // Force show
            console.log('Tab pane activated:', targetPane.id);
            console.log('Tab pane display style:', window.getComputedStyle(targetPane).display);
        } else {
            console.error('Tab pane not found:', `${tabName}Tab`);
            return; // Exit early if tab pane not found
        }
        
        // Load tab-specific content
        this.loadTabContent(tabName);
    }
    
    async loadTabContent(tabName) {
        console.log('Loading tab content for:', tabName);
        console.log('Current container ID:', this.currentContainerId);
        
        if (!this.currentContainerId) {
            console.error('No current container ID set!');
            return;
        }
        
        try {
            switch (tabName) {
                case 'logs':
                    console.log('Loading logs...');
                    await this.loadContainerLogs();
                    break;
                case 'files':
                    console.log('Loading files...');
                    await this.loadContainerFiles('/');
                    break;
                case 'env':
                    console.log('Loading environment...');
                    await this.loadContainerEnvironment();
                    break;
                case 'exec':
                    console.log('Setting up exec...');
                    this.setupContainerExec();
                    // Auto-launch terminal session
                    setTimeout(() => {
                        if (!this.isExecConnected) {
                            console.log('Auto-launching terminal session...');
                            const execShell = document.getElementById('execShell');
                            const shell = execShell ? execShell.value : '/bin/sh';
                            this.startExecSession(shell);
                        }
                    }, 100); // Small delay to ensure UI is ready
                    break;
                case 'volumes':
                    console.log('Loading volumes...');
                    await this.loadContainerVolumes();
                    break;
                case 'network':
                    console.log('Loading network...');
                    await this.loadContainerNetwork();
                    break;
                case 'stats':
                    console.log('Loading stats...');
                    await this.loadContainerStats();
                    break;
                default:
                    console.warn('Unknown tab:', tabName);
            }
        } catch (error) {
            console.error(`Error loading ${tabName} content:`, error);
        }
    }
    
    async loadContainerLogs() {
        const logsContent = document.getElementById('logsContent');
        logsContent.innerHTML = '<div class="loading">Loading logs...</div>';
        
        // Set up logs controls
        this.setupLogsControls();
        
        try {
            // Try HTTP API first for reliability
            const response = await fetch(`${this.apiBase}/api/containers/${this.currentContainerId}/logs`);
            if (response.ok) {
                const data = await response.json();
                const logs = data.logs || 'No logs available';
                logsContent.innerHTML = `<div class="logs-stream"><pre class="logs-text" id="logsText">${logs}</pre></div>`;
                
                // If successful and logs exist, also start real-time streaming
                if (logs && logs !== 'No logs available') {
                    this.startLogStreaming();
                }
            } else {
                // Fallback to real-time streaming
                this.startLogStreaming();
            }
        } catch (error) {
            console.error('Error loading logs via HTTP, trying WebSocket:', error);
            // Fallback to real-time streaming
            this.startLogStreaming();
        }
    }
    
    setupLogsControls() {
        const downloadLogsBtn = document.getElementById('downloadLogsBtn');
        
        if (downloadLogsBtn) {
            downloadLogsBtn.onclick = () => this.downloadLogs();
        }
        
        // Auto-start log streaming - no manual control needed
        if (!this.isLogsStreaming) {
            this.startLogStreaming();
        }
    }
    
    startLogStreaming() {
        if (!this.currentContainerId || this.isLogsStreaming) {
            return;
        }
        
        const logsContent = document.getElementById('logsContent');
        logsContent.innerHTML = '<div class="logs-stream"><pre class="logs-text" id="logsText"></pre></div>';
        
        // Start streaming
        this.socket.emit('start_logs', {
            container_id: this.currentContainerId
        });
    }
    
    stopLogStreaming() {
        if (!this.currentContainerId || !this.isLogsStreaming) {
            return;
        }
        
        this.socket.emit('stop_logs', {
            container_id: this.currentContainerId
        });
    }
    
    handleLogData(data) {
        if (data.container_id !== this.currentContainerId) {
            return; // Ignore logs from other containers
        }
        
        const logsText = document.getElementById('logsText');
        if (!logsText) {
            return;
        }
        
        if (data.type === 'initial') {
            // Replace initial loading with historical logs
            logsText.textContent = data.data || 'No logs available';
        } else if (data.type === 'stream') {
            // Append new log line
            logsText.textContent += '\n' + data.data;
        }
        
        // Auto-scroll to bottom
        const logsContent = document.getElementById('logsContent');
        if (logsContent) {
            logsContent.scrollTop = logsContent.scrollHeight;
        }
    }
    
    updateLogsUI() {
        const autoRefreshLogs = document.getElementById('autoRefreshLogs');
        if (autoRefreshLogs) {
            autoRefreshLogs.checked = this.isLogsStreaming;
        }
    }
    
    downloadLogs() {
        const logsText = document.getElementById('logsText');
        if (!logsText || !logsText.textContent) {
            this.showError('No logs to download');
            return;
        }
        
        const logs = logsText.textContent;
        const blob = new Blob([logs], { type: 'text/plain' });
        const url = URL.createObjectURL(blob);
        
        const a = document.createElement('a');
        a.href = url;
        a.download = `${this.currentContainerId}_logs.txt`;
        a.click();
        
        URL.revokeObjectURL(url);
    }
    
    async loadContainerFiles(path = '/') {
        console.log('loadContainerFiles called with path:', path);
        console.log('Container ID:', this.currentContainerId);
        
        const filesContent = document.getElementById('filesContent');
        const currentPath = document.getElementById('currentPath');
        const refreshFilesBtn = document.getElementById('refreshFilesBtn');
        
        filesContent.innerHTML = '<div class="loading">Loading files...</div>';
        currentPath.textContent = path;
        
        // Set up refresh button
        if (refreshFilesBtn) {
            refreshFilesBtn.onclick = () => this.loadContainerFiles(path);
        }
        
        try {
            const url = `${this.apiBase}/api/containers/${this.currentContainerId}/files?path=${encodeURIComponent(path)}`;
            console.log('Fetching files from URL:', url);
            
            const response = await fetch(url);
            console.log('Files API response status:', response.status);
            
            if (response.ok) {
                const files = await response.json();
                console.log('Files data received:', files);
                this.renderFiles(files, path);
            } else {
                const error = await response.json();
                console.error('Files API error:', error);
                filesContent.innerHTML = `<div class="error">Error loading files: ${error.error || 'Unknown error'}</div>`;
            }
        } catch (error) {
            console.error('Error loading files:', error);
            filesContent.innerHTML = '<div class="error">Error loading files</div>';
        }
    }
    
    renderFiles(files, currentPath) {
        const filesContent = document.getElementById('filesContent');
        
        let html = '<div class="files-list">';
        
        // Add parent directory link if not root (always show for navigation)
        if (currentPath !== '/') {
            const parentPath = currentPath.split('/').slice(0, -1).join('/') || '/';
            html += `<div class="file-item directory parent-dir" onclick="dockerGUI.loadContainerFiles('${parentPath}')">
                <i class="fas fa-level-up-alt"></i>
                <span>..</span>
                <span class="file-size">Parent Directory</span>
            </div>`;
        }
        
        if (!files || files.length === 0) {
            html += '</div><div class="empty">No files found in this directory</div>';
            filesContent.innerHTML = html;
            return;
        }
        
        // Sort files: directories first, then files
        const sortedFiles = [...files].sort((a, b) => {
            if (a.is_directory && !b.is_directory) return -1;
            if (!a.is_directory && b.is_directory) return 1;
            return a.name.localeCompare(b.name);
        });
        
        sortedFiles.forEach(file => {
            const icon = file.is_directory ? 'fa-folder' : 'fa-file';
            const fileClass = file.is_directory ? 'directory' : 'file';
            
            // Construct proper file path
            let filePath;
            if (file.path) {
                filePath = file.path;
            } else {
                // Ensure proper path construction - fix the logic
                if (currentPath === '/') {
                    filePath = `/${file.name}`;
                } else {
                    filePath = `${currentPath}/${file.name}`;
                }
            }
            
            const onclick = file.is_directory ? `dockerGUI.loadContainerFiles('${filePath}')` : '';
            const cursor = file.is_directory ? 'pointer' : 'default';
            
            html += `<div class="file-item ${fileClass}" ${onclick ? `onclick="${onclick}"` : ''} style="cursor: ${cursor}">
                <i class="fas ${icon}"></i>
                <span>${file.name}</span>
                <span class="file-size">${file.is_directory ? 'Directory' : this.formatFileSize(file.size || 0)}</span>
            </div>`;
        });
        
        html += '</div>';
        filesContent.innerHTML = html;
    }
    
    async loadContainerEnvironment() {
        console.log('loadContainerEnvironment called');
        console.log('Container ID:', this.currentContainerId);
        
        const envContent = document.getElementById('envContent');
        envContent.innerHTML = '<div class="loading">Loading environment variables...</div>';
        
        try {
            const url = `${this.apiBase}/api/containers/${this.currentContainerId}/env`;
            console.log('Fetching environment from URL:', url);
            
            const response = await fetch(url);
            console.log('Environment API response status:', response.status);
            
            if (response.ok) {
                const data = await response.json();
                console.log('Environment data received:', data);
                this.renderEnvironmentVariables(data.environment || data);
            } else {
                const errorText = await response.text();
                console.error('Environment API error:', errorText);
                envContent.innerHTML = '<div class="error">Error loading environment variables</div>';
            }
        } catch (error) {
            console.error('Error loading environment:', error);
            envContent.innerHTML = '<div class="error">Error loading environment variables</div>';
        }
    }
    
    renderEnvironmentVariables(envVars) {
        const envContent = document.getElementById('envContent');
        
        if (!envVars || (Array.isArray(envVars) && envVars.length === 0) || (typeof envVars === 'object' && Object.keys(envVars).length === 0)) {
            envContent.innerHTML = '<div class="empty">No environment variables found</div>';
            return;
        }
        
        let html = '<div class="env-list">';
        
        if (Array.isArray(envVars)) {
            // Handle array format: [{key: 'PATH', value: '/usr/bin'}, ...]
            envVars.forEach(envVar => {
                if (envVar.key && envVar.value) {
                    html += `<div class="env-item">
                        <div class="env-key">${envVar.key}</div>
                        <div class="env-value">${envVar.value}</div>
                    </div>`;
                }
            });
        } else {
            // Handle object format: {PATH: '/usr/bin', ...}
            Object.entries(envVars).forEach(([key, value]) => {
                html += `<div class="env-item">
                    <div class="env-key">${key}</div>
                    <div class="env-value">${value}</div>
                </div>`;
            });
        }
        
        html += '</div>';
        envContent.innerHTML = html;
    }
    
    setupContainerExec() {
        const execTerminal = document.getElementById('execTerminal');
        const connectExecBtn = document.getElementById('connectExecBtn');
        const execShell = document.getElementById('execShell');
        
        // Initialize terminal display
        execTerminal.innerHTML = `
            <div class="terminal-placeholder" id="terminalPlaceholder">
                Click "Connect" to open a terminal session in the container
            </div>
            <div class="terminal-output" id="terminalOutput" style="display: none;">
                <div class="terminal-content" id="terminalContent"></div>
                <div class="terminal-input-line">
                    <span class="terminal-prompt" id="terminalPrompt">$</span>
                    <input type="text" class="terminal-input" id="terminalInput" placeholder="Enter command...">
                </div>
            </div>
        `;
        
        // Set up connect button
        if (connectExecBtn) {
            connectExecBtn.onclick = () => this.toggleExecSession();
        }
        
        // Set up terminal input
        const terminalInput = document.getElementById('terminalInput');
        if (terminalInput) {
            terminalInput.addEventListener('keypress', (e) => {
                if (e.key === 'Enter') {
                    this.sendExecCommand(terminalInput.value);
                    terminalInput.value = '';
                }
            });
        }
    }
    
    toggleExecSession() {
        const connectExecBtn = document.getElementById('connectExecBtn');
        const execShell = document.getElementById('execShell');
        
        if (!this.isExecConnected) {
            // Start exec session
            const shell = execShell.value;
            this.startExecSession(shell);
        } else {
            // Stop exec session
            this.stopExecSession();
        }
    }
    
    startExecSession(shell) {
        if (!this.currentContainerId) {
            this.showError('No container selected');
            return;
        }
        
        // Update UI
        const connectExecBtn = document.getElementById('connectExecBtn');
        const terminalPlaceholder = document.getElementById('terminalPlaceholder');
        const terminalOutput = document.getElementById('terminalOutput');
        
        connectExecBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Connecting...';
        connectExecBtn.disabled = true;
        
        // Start exec session via WebSocket
        this.socket.emit('start_exec', {
            container_id: this.currentContainerId,
            shell: shell
        });
    }
    
    stopExecSession() {
        if (!this.currentContainerId) {
            return;
        }
        
        // Stop exec session via WebSocket
        this.socket.emit('stop_exec', {
            container_id: this.currentContainerId
        });
    }
    
    sendExecCommand(command) {
        if (!this.isExecConnected || !command.trim()) {
            return;
        }
        
        // Send command via WebSocket
        this.socket.emit('exec_input', {
            container_id: this.currentContainerId,
            command: command.trim()
        });
    }
    
    handleExecStarted(data) {
        console.log('Exec session started:', data);
        this.isExecConnected = true;
        
        // Update UI
        const connectExecBtn = document.getElementById('connectExecBtn');
        const terminalPlaceholder = document.getElementById('terminalPlaceholder');
        const terminalOutput = document.getElementById('terminalOutput');
        const terminalInput = document.getElementById('terminalInput');
        
        connectExecBtn.innerHTML = '<i class="fas fa-times"></i> Disconnect';
        connectExecBtn.disabled = false;
        
        if (terminalPlaceholder) terminalPlaceholder.style.display = 'none';
        if (terminalOutput) terminalOutput.style.display = 'block';
        if (terminalInput) terminalInput.focus();
    }
    
    handleExecStopped(data) {
        console.log('Exec session stopped:', data);
        this.isExecConnected = false;
        
        // Update UI
        const connectExecBtn = document.getElementById('connectExecBtn');
        const terminalPlaceholder = document.getElementById('terminalPlaceholder');
        const terminalOutput = document.getElementById('terminalOutput');
        
        connectExecBtn.innerHTML = '<i class="fas fa-plug"></i> Connect';
        connectExecBtn.disabled = false;
        
        if (terminalPlaceholder) terminalPlaceholder.style.display = 'block';
        if (terminalOutput) terminalOutput.style.display = 'none';
    }
    
    handleExecOutput(data) {
        if (data.container_id !== this.currentContainerId) {
            return; // Ignore output from other containers
        }
        
        const terminalContent = document.getElementById('terminalContent');
        const terminalPrompt = document.getElementById('terminalPrompt');
        
        if (!terminalContent) return;
        
        // Create output element
        const outputElement = document.createElement('div');
        outputElement.className = `terminal-line terminal-${data.type}`;
        
        if (data.type === 'prompt') {
            if (terminalPrompt) {
                terminalPrompt.textContent = data.data;
            }
        } else {
            outputElement.textContent = data.data;
            terminalContent.appendChild(outputElement);
        }
        
        // Auto-scroll to bottom
        const terminalOutput = document.getElementById('terminalOutput');
        if (terminalOutput) {
            terminalOutput.scrollTop = terminalOutput.scrollHeight;
        }
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

    formatFileSize(bytes) {
        if (bytes === 0 || !bytes) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
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
