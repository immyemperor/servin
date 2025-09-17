/**
 * VM Engine Management Component
 * Handles VM status monitoring and control operations
 */

class VMManager {
    constructor() {
        this.isPolling = false;
        this.pollInterval = null;
        this.currentStatus = null;
        this.isLoading = false; // Prevent concurrent loading operations
        
        this.initializeEventListeners();
        // Don't load VM status immediately, wait for section to be shown
    }

    initializeEventListeners() {
        // Navigation
        const vmNavItem = document.querySelector('[data-section="vm"]');
        if (vmNavItem) {
            vmNavItem.addEventListener('click', () => {
                this.showVMSection();
                this.loadVMStatus(true); // Show loading for initial section load
            });
        }

        // VM Control Buttons
        document.getElementById('startVmBtn')?.addEventListener('click', () => this.startVM());
        document.getElementById('stopVmBtn')?.addEventListener('click', () => this.stopVM());
        document.getElementById('restartVmBtn')?.addEventListener('click', () => this.restartVM());
        document.getElementById('enableVmBtn')?.addEventListener('click', () => this.enableVM());
        document.getElementById('disableVmBtn')?.addEventListener('click', () => this.disableVM());
        document.getElementById('refreshVmBtn')?.addEventListener('click', () => {
            const refreshBtn = document.getElementById('refreshVmBtn');
            if (refreshBtn) {
                // Show mini loading state on button
                const originalHTML = refreshBtn.innerHTML;
                refreshBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Refreshing...';
                refreshBtn.disabled = true;
                
                this.loadVMStatus(false).finally(() => {
                    // Restore button state
                    refreshBtn.innerHTML = originalHTML;
                    refreshBtn.disabled = false;
                });
            }
        }); // Show loading for manual refresh
        document.getElementById('clearVmLogsBtn')?.addEventListener('click', () => this.clearLogs());
    }

    showVMSection() {
        // Hide all sections
        document.querySelectorAll('.content-section').forEach(section => {
            section.classList.remove('active');
        });
        
        // Show VM section
        const vmSection = document.getElementById('vmSection');
        if (vmSection) {
            vmSection.classList.add('active');
        }

        // Update navigation
        document.querySelectorAll('.nav-item').forEach(item => {
            item.classList.remove('active');
        });
        document.querySelector('[data-section="vm"]')?.classList.add('active');

        // Start polling when VM section is active
        this.startPolling();
    }

    async loadVMStatus(showLoading = false) {
        // Prevent concurrent loading operations
        if (this.isLoading) {
            return;
        }
        
        console.log('Loading VM status...');
        this.isLoading = true;
        
        try {
            // Only show loading spinner for manual refreshes, not automatic polling
            if (showLoading) {
                UIHelpers.showLoading();
            }
            const response = await fetch('/api/vm/status');
            const data = await response.json();
            
            console.log('VM status response:', data);
            
            if (data.available) {
                this.updateVMStatus(data);
            } else {
                this.updateVMUnavailable(data.error || 'VM engine not available');
            }
        } catch (error) {
            console.error('Failed to load VM status:', error);
            this.updateVMUnavailable('Failed to connect to VM engine');
        } finally {
            this.isLoading = false;
            if (showLoading) {
                UIHelpers.hideLoading();
            }
        }
    }

    updateVMStatus(status) {
        console.log('Updating VM status:', status);
        this.currentStatus = status;
        
        // Update main status indicator
        const statusIndicator = document.getElementById('vmStatusIndicator');
        const statusText = document.getElementById('vmStatusText');
        
        // Update engine status indicator
        const engineStatus = document.getElementById('vmEngineStatus');
        const engineText = document.getElementById('vmEngineText');
        
        if (statusIndicator && statusText) {
            statusIndicator.className = 'status-indicator-large';
            
            if (status.running) {
                statusIndicator.classList.add('running');
                statusText.textContent = 'VM Engine Running';
            } else if (status.enabled) {
                statusIndicator.classList.add('stopped');
                statusText.textContent = 'VM Engine Stopped';
            } else {
                statusIndicator.classList.add('stopped');
                statusText.textContent = 'VM Mode Disabled';
            }
        }
        
        // Update detailed engine status with color coding
        if (engineStatus && engineText) {
            engineStatus.className = 'engine-status';
            
            if (status.running) {
                engineStatus.classList.add('running');
                engineText.textContent = 'Running';
            } else if (status.enabled) {
                engineStatus.classList.add('stopped');
                engineText.textContent = 'Stopped';
            } else {
                engineStatus.classList.add('disabled');
                engineText.textContent = 'Disabled';
            }
        }

        // Update details
        this.updateVMDetails(status);
        
        // Update button states
        this.updateButtonStates(status);
        
        this.addLogEntry(`VM Status Updated: ${status.running ? 'Running' : (status.enabled ? 'Stopped' : 'Disabled')}`);
    }

    updateVMDetails(status) {
        const details = status.details || {};
        
        document.getElementById('vmProvider').textContent = status.provider || '-';
        document.getElementById('vmPlatform').textContent = status.platform || '-';
        document.getElementById('vmContainers').textContent = status.containers || '0';
        document.getElementById('vmName').textContent = details.name || '-';
        document.getElementById('vmIpAddress').textContent = details.ip || '-';
    }

    updateVMUnavailable(error) {
        const statusIndicator = document.getElementById('vmStatusIndicator');
        const statusText = document.getElementById('vmStatusText');
        const engineStatus = document.getElementById('vmEngineStatus');
        const engineText = document.getElementById('vmEngineText');
        
        if (statusIndicator && statusText) {
            statusIndicator.className = 'status-indicator-large stopped';
            statusText.textContent = 'VM Engine Unavailable';
        }
        
        if (engineStatus && engineText) {
            engineStatus.className = 'engine-status disabled';
            engineText.textContent = 'Unavailable';
        }

        // Clear details
        ['vmProvider', 'vmPlatform', 'vmContainers', 'vmName', 'vmIpAddress'].forEach(id => {
            const element = document.getElementById(id);
            if (element) element.textContent = '-';
        });

        // Disable all buttons
        this.disableAllButtons();
        
        this.addLogEntry(`VM Error: ${error}`, 'error');
    }

    updateButtonStates(status) {
        const startBtn = document.getElementById('startVmBtn');
        const stopBtn = document.getElementById('stopVmBtn');
        const restartBtn = document.getElementById('restartVmBtn');
        const enableBtn = document.getElementById('enableVmBtn');
        const disableBtn = document.getElementById('disableVmBtn');

        // Reset all buttons to enabled first
        [startBtn, stopBtn, restartBtn, enableBtn, disableBtn].forEach(btn => {
            if (btn) btn.disabled = false;
        });

        if (status.available) {
            if (status.enabled) {
                // VM mode is enabled
                if (enableBtn) enableBtn.disabled = true;
                if (disableBtn) disableBtn.disabled = false;
                
                if (status.running) {
                    // VM is running
                    if (startBtn) startBtn.disabled = true;
                    if (stopBtn) stopBtn.disabled = false;
                    if (restartBtn) restartBtn.disabled = false;
                } else {
                    // VM is stopped
                    if (startBtn) startBtn.disabled = false;
                    if (stopBtn) stopBtn.disabled = true;
                    if (restartBtn) restartBtn.disabled = true;
                }
            } else {
                // VM mode is disabled
                if (enableBtn) enableBtn.disabled = false;
                if (disableBtn) disableBtn.disabled = true;
                if (startBtn) startBtn.disabled = true;
                if (stopBtn) stopBtn.disabled = true;
                if (restartBtn) restartBtn.disabled = true;
            }
        } else {
            this.disableAllButtons();
        }
    }

    disableAllButtons() {
        ['startVmBtn', 'stopVmBtn', 'restartVmBtn', 'enableVmBtn', 'disableVmBtn'].forEach(id => {
            const button = document.getElementById(id);
            if (button) button.disabled = true;
        });
    }

    updateEngineTransitionState(state, text) {
        const engineStatus = document.getElementById('vmEngineStatus');
        const engineText = document.getElementById('vmEngineText');
        const statusIndicator = document.getElementById('vmStatusIndicator');
        const statusText = document.getElementById('vmStatusText');
        
        if (engineStatus && engineText) {
            engineStatus.className = `engine-status ${state}`;
            engineText.textContent = text;
        }
        
        if (statusIndicator && statusText) {
            statusIndicator.className = `status-indicator-large ${state}`;
            statusText.textContent = `VM Engine ${text}`;
        }
    }

    async startVM() {
        const startBtn = document.getElementById('startVmBtn');
        const originalHTML = startBtn?.innerHTML;
        
        try {
            // Show button loading state
            if (startBtn) {
                startBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Starting...';
                startBtn.disabled = true;
            }
            
            // Show transitional state
            this.updateEngineTransitionState('starting', 'Starting');
            this.addLogEntry('Starting VM engine...', 'info');
            
            const response = await fetch('/api/vm/start', { method: 'POST' });
            const data = await response.json();
            
            if (data.success) {
                UIHelpers.showToast('VM engine started successfully', 'success');
                this.addLogEntry('VM engine started successfully', 'success');
                setTimeout(() => this.loadVMStatus(false), 1000); // Reduced delay
            } else {
                throw new Error(data.error || 'Failed to start VM');
            }
        } catch (error) {
            console.error('Failed to start VM:', error);
            UIHelpers.showToast(`Failed to start VM: ${error.message}`, 'error');
            this.addLogEntry(`Failed to start VM: ${error.message}`, 'error');
        } finally {
            // Restore button state
            if (startBtn && originalHTML) {
                startBtn.innerHTML = originalHTML;
                // Don't re-enable immediately, let the status update handle it
                setTimeout(() => {
                    if (startBtn) startBtn.disabled = false;
                }, 500);
            }
        }
    }

    async stopVM() {
        const stopBtn = document.getElementById('stopVmBtn');
        const originalHTML = stopBtn?.innerHTML;
        
        try {
            // Show button loading state
            if (stopBtn) {
                stopBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Stopping...';
                stopBtn.disabled = true;
            }
            
            // Show transitional state
            this.updateEngineTransitionState('starting', 'Stopping');
            this.addLogEntry('Stopping VM engine...', 'info');
            
            const response = await fetch('/api/vm/stop', { method: 'POST' });
            const data = await response.json();
            
            if (data.success) {
                UIHelpers.showToast('VM engine stopped successfully', 'success');
                this.addLogEntry('VM engine stopped successfully', 'success');
                setTimeout(() => this.loadVMStatus(false), 500); // Quick refresh for stop
            } else {
                throw new Error(data.error || 'Failed to stop VM');
            }
        } catch (error) {
            console.error('Failed to stop VM:', error);
            UIHelpers.showToast(`Failed to stop VM: ${error.message}`, 'error');
            this.addLogEntry(`Failed to stop VM: ${error.message}`, 'error');
        } finally {
            // Restore button state
            if (stopBtn && originalHTML) {
                stopBtn.innerHTML = originalHTML;
                setTimeout(() => {
                    if (stopBtn) stopBtn.disabled = false;
                }, 500);
            }
        }
    }

    async restartVM() {
        const restartBtn = document.getElementById('restartVmBtn');
        const originalHTML = restartBtn?.innerHTML;
        
        try {
            // Show button loading state
            if (restartBtn) {
                restartBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Restarting...';
                restartBtn.disabled = true;
            }
            
            // Show transitional state
            this.updateEngineTransitionState('starting', 'Restarting');
            this.addLogEntry('Restarting VM engine...', 'info');
            
            const response = await fetch('/api/vm/restart', { method: 'POST' });
            const data = await response.json();
            
            if (data.success) {
                UIHelpers.showToast('VM engine restarted successfully', 'success');
                this.addLogEntry('VM engine restarted successfully', 'success');
                setTimeout(() => this.loadVMStatus(false), 1500); // Moderate delay for restart
            } else {
                throw new Error(data.error || 'Failed to restart VM');
            }
        } catch (error) {
            console.error('Failed to restart VM:', error);
            UIHelpers.showToast(`Failed to restart VM: ${error.message}`, 'error');
            this.addLogEntry(`Failed to restart VM: ${error.message}`, 'error');
        } finally {
            // Restore button state
            if (restartBtn && originalHTML) {
                restartBtn.innerHTML = originalHTML;
                setTimeout(() => {
                    if (restartBtn) restartBtn.disabled = false;
                }, 500);
            }
        }
    }

    async enableVM() {
        try {
            this.addLogEntry('Enabling VM mode...', 'info');
            
            const response = await fetch('/api/vm/enable', { method: 'POST' });
            const data = await response.json();
            
            if (data.success) {
                UIHelpers.showToast('VM mode enabled successfully', 'success');
                this.addLogEntry('VM mode enabled successfully', 'success');
                setTimeout(() => this.loadVMStatus(false), 500); // Quick refresh
            } else {
                throw new Error(data.error || 'Failed to enable VM mode');
            }
        } catch (error) {
            console.error('Failed to enable VM:', error);
            UIHelpers.showToast(`Failed to enable VM: ${error.message}`, 'error');
            this.addLogEntry(`Failed to enable VM: ${error.message}`, 'error');
        }
    }

    async disableVM() {
        try {
            this.addLogEntry('Disabling VM mode...', 'info');
            
            const response = await fetch('/api/vm/disable', { method: 'POST' });
            const data = await response.json();
            
            if (data.success) {
                UIHelpers.showToast('VM mode disabled successfully', 'success');
                this.addLogEntry('VM mode disabled successfully', 'success');
                setTimeout(() => this.loadVMStatus(false), 500); // Quick refresh
            } else {
                throw new Error(data.error || 'Failed to disable VM mode');
            }
        } catch (error) {
            console.error('Failed to disable VM:', error);
            UIHelpers.showToast(`Failed to disable VM: ${error.message}`, 'error');
            this.addLogEntry(`Failed to disable VM: ${error.message}`, 'error');
        }
    }

    addLogEntry(message, type = 'info') {
        const logsContent = document.getElementById('vmLogsContent');
        if (!logsContent) return;

        // Remove placeholder if it exists
        const placeholder = logsContent.querySelector('.log-placeholder');
        if (placeholder) {
            placeholder.remove();
        }

        const timestamp = new Date().toLocaleTimeString();
        const logEntry = document.createElement('div');
        logEntry.className = `log-entry log-${type}`;
        logEntry.innerHTML = `<span class="log-time">[${timestamp}]</span> <span class="log-message">${message}</span>`;
        
        logsContent.appendChild(logEntry);
        logsContent.scrollTop = logsContent.scrollHeight;

        // Limit log entries to prevent memory issues
        const entries = logsContent.querySelectorAll('.log-entry');
        if (entries.length > 100) {
            entries[0].remove();
        }
    }

    clearLogs() {
        const logsContent = document.getElementById('vmLogsContent');
        if (logsContent) {
            logsContent.innerHTML = '<p class="log-placeholder">VM logs will appear here...</p>';
        }
    }

    startPolling() {
        if (this.isPolling) return;
        
        this.isPolling = true;
        this.pollInterval = setInterval(() => {
            // Only poll if VM section is active
            const vmSection = document.getElementById('vmSection');
            if (vmSection && vmSection.classList.contains('active')) {
                this.loadVMStatus(false); // Don't show loading for automatic polling
            } else {
                this.stopPolling();
            }
        }, 5000); // Poll every 5 seconds
    }

    stopPolling() {
        this.isPolling = false;
        if (this.pollInterval) {
            clearInterval(this.pollInterval);
            this.pollInterval = null;
        }
    }

    destroy() {
        this.stopPolling();
    }
}

// Add log entry styles
const style = document.createElement('style');
style.textContent = `
    .log-entry {
        margin-bottom: 4px;
        line-height: 1.4;
    }
    
    .log-time {
        color: var(--text-secondary);
        font-size: 11px;
    }
    
    .log-message {
        color: var(--text-primary);
    }
    
    .log-info .log-message {
        color: var(--text-primary);
    }
    
    .log-success .log-message {
        color: var(--success-color);
    }
    
    .log-error .log-message {
        color: var(--danger-color);
    }
`;
document.head.appendChild(style);

// Initialize VM Manager when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.vmManager = new VMManager();
});

// Export for use in other modules
if (typeof module !== 'undefined' && module.exports) {
    module.exports = VMManager;
}