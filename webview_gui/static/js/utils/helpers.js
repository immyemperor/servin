/**
 * Utility functions and helpers for Docker Desktop GUI
 */

class UIHelpers {
    /**
     * Show toast notification
     */
    static showToast(message, type = 'info') {
        const container = document.getElementById('toastContainer');
        if (!container) return;

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

    /**
     * Get appropriate icon for toast type
     */
    static getToastIcon(type) {
        const icons = {
            success: 'fa-check-circle',
            error: 'fa-exclamation-circle',
            warning: 'fa-exclamation-triangle',
            info: 'fa-info-circle'
        };
        return icons[type] || icons.info;
    }

    /**
     * Show loading overlay
     */
    static showLoading() {
        const overlay = document.getElementById('loadingOverlay');
        if (overlay) {
            overlay.style.display = 'flex';
        }
    }

    /**
     * Hide loading overlay
     */
    static hideLoading() {
        const overlay = document.getElementById('loadingOverlay');
        if (overlay) {
            overlay.style.display = 'none';
        }
    }

    /**
     * Format date string for display
     */
    static formatDate(dateString) {
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
            return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
        } catch (error) {
            return dateString; // Return original if parsing fails
        }
    }

    /**
     * Format bytes to human readable size
     */
    static formatBytes(bytes) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

    /**
     * Format file size for display
     */
    static formatFileSize(bytes) {
        return this.formatBytes(bytes);
    }

    /**
     * Filter table rows based on search term
     */
    static filterTable(type, searchTerm) {
        const tableBody = document.getElementById(`${type}sTableBody`);
        if (!tableBody) return;

        const rows = tableBody.querySelectorAll('tr');
        
        rows.forEach(row => {
            const text = row.textContent.toLowerCase();
            row.style.display = text.includes(searchTerm.toLowerCase()) ? '' : 'none';
        });
    }

    /**
     * Setup search functionality for a table
     */
    static setupSearch(type) {
        const searchInput = document.getElementById(`${type}Search`);
        if (searchInput) {
            searchInput.addEventListener('input', (e) => {
                this.filterTable(type, e.target.value);
            });
        }
    }

    /**
     * Update navigation item counts
     */
    static updateCounts(containers, images, volumes) {
        const containerCount = document.getElementById('containerCount');
        const imageCount = document.getElementById('imageCount');
        const volumeCount = document.getElementById('volumeCount');

        if (containerCount) containerCount.textContent = containers?.length || 0;
        if (imageCount) imageCount.textContent = images?.length || 0;
        if (volumeCount) volumeCount.textContent = volumes?.length || 0;
    }

    /**
     * Switch between main sections
     */
    static switchSection(section) {
        // Update navigation
        document.querySelectorAll('.nav-item').forEach(item => {
            item.classList.remove('active');
        });
        const targetNav = document.querySelector(`[data-section="${section}"]`);
        if (targetNav) {
            targetNav.classList.add('active');
        }
        
        // Update content sections
        document.querySelectorAll('.content-section').forEach(sec => {
            sec.classList.remove('active');
        });
        const targetSection = document.getElementById(`${section}Section`);
        if (targetSection) {
            targetSection.classList.add('active');
        }
    }

    /**
     * Setup modal controls
     */
    static setupModal(modalId, openButtonId, closeButtonId, cancelButtonId) {
        const modal = document.getElementById(modalId);
        const openBtn = document.getElementById(openButtonId);
        const closeBtn = document.getElementById(closeButtonId);
        const cancelBtn = document.getElementById(cancelButtonId);

        if (openBtn && modal) {
            openBtn.addEventListener('click', () => {
                modal.style.display = 'block';
            });
        }

        const closeModal = () => {
            if (modal) modal.style.display = 'none';
        };

        if (closeBtn) closeBtn.addEventListener('click', closeModal);
        if (cancelBtn) cancelBtn.addEventListener('click', closeModal);

        // Close on outside click
        if (modal) {
            modal.addEventListener('click', (e) => {
                if (e.target === modal) {
                    closeModal();
                }
            });
        }
    }
}

/**
 * Socket management utilities
 */
class SocketManager {
    constructor() {
        this.socket = null;
        this.isConnected = false;
        this.eventHandlers = new Map();
    }

    /**
     * Initialize socket connection
     */
    init() {
        this.socket = io();
        
        this.socket.on('connect', () => {
            console.log('Socket connected');
            this.isConnected = true;
            this.trigger('connect');
        });
        
        this.socket.on('disconnect', () => {
            console.log('Socket disconnected');
            this.isConnected = false;
            this.trigger('disconnect');
        });

        this.socket.on('error', (data) => {
            console.error('Socket error:', data);
            this.trigger('error', data);
        });
    }

    /**
     * Add event listener
     */
    on(event, handler) {
        if (!this.eventHandlers.has(event)) {
            this.eventHandlers.set(event, []);
        }
        this.eventHandlers.get(event).push(handler);

        // Also register with socket
        if (this.socket) {
            this.socket.on(event, handler);
        }
    }

    /**
     * Emit event
     */
    emit(event, data) {
        if (this.socket) {
            this.socket.emit(event, data);
        }
    }

    /**
     * Trigger local event handlers
     */
    trigger(event, data) {
        const handlers = this.eventHandlers.get(event);
        if (handlers) {
            handlers.forEach(handler => handler(data));
        }
    }
}

// Export utilities
window.UIHelpers = UIHelpers;
window.SocketManager = SocketManager;