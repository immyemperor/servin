/**
 * Logs Component
 * Handles container log streaming and display
 */

class LogsComponent {
    constructor(apiClient, socketManager) {
        this.apiClient = apiClient;
        this.socketManager = socketManager;
        this.currentContainerId = null;
        this.isStreaming = false;
        this.autoScroll = true; // Enable auto-scroll by default
        
        this.init();
    }

    init() {
        this.setupSocketHandlers();
    }

    setupSocketHandlers() {
        this.socketManager.on('log_data', (data) => this.handleLogData(data));
        this.socketManager.on('logs_started', (data) => this.handleLogsStarted(data));
        this.socketManager.on('logs_stopped', (data) => this.handleLogsStopped(data));
    }

    async loadLogs(containerId) {
        this.currentContainerId = containerId;
        const logsContent = document.getElementById('logsContent');
        
        if (!logsContent) return;
        
        // Only setup controls if not already set up for this container
        const existingLogsText = document.getElementById('logsText');
        if (!existingLogsText) {
            logsContent.innerHTML = '<div class="loading">Loading logs...</div>';
            this.setupLogsControls();
        }
        
        try {
            // Start log streaming automatically if not already streaming
            if (!this.isStreaming) {
                this.startStreaming();
            }
        } catch (error) {
            console.error('Failed to load logs:', error);
            logsContent.innerHTML = '<div class="error">Failed to load logs</div>';
        }
    }

    setupLogsControls() {
        const logsContent = document.getElementById('logsContent');
        if (!logsContent) return;

        // Set up the logs stream HTML structure
        logsContent.innerHTML = `
            <div class="logs-stream">
                <pre class="logs-text" id="logsText"></pre>
            </div>
        `;
    }

    startStreaming() {
        if (!this.currentContainerId || this.isStreaming) return;

        const logsContent = document.getElementById('logsContent');
        const existingLogsText = document.getElementById('logsText');
        
        if (!logsContent) return;

        // Only create the structure if it doesn't exist
        if (!existingLogsText) {
            logsContent.innerHTML = `
                <div class="logs-stream">
                    <pre class="logs-text" id="logsText"></pre>
                </div>
            `;
        }

        // Start streaming via WebSocket
        this.socketManager.emit('start_logs', {
            container_id: this.currentContainerId
        });
    }

    stopStreaming() {
        if (!this.currentContainerId || !this.isStreaming) return;

        this.socketManager.emit('stop_logs', {
            container_id: this.currentContainerId
        });
    }

    handleLogData(data) {
        if (data.container_id !== this.currentContainerId) return;

        const logsText = document.getElementById('logsText');
        if (!logsText) return;

        // The server sends log data in the 'data' property, not 'logs'
        const logContent = data.data || '';

        if (data.type === 'initial') {
            // Replace all content with initial logs
            logsText.textContent = logContent;
        } else {
            // Append new log lines
            logsText.textContent += logContent + '\n';
        }

        // Auto-scroll to bottom if enabled
        this.autoScrollToBottom();
    }

    handleLogsStarted(data) {
        console.log('Logs streaming started:', data);
        this.isStreaming = true;
        this.updateStreamingStatus();
    }

    handleLogsStopped(data) {
        console.log('Logs streaming stopped:', data);
        this.isStreaming = false;
        this.updateStreamingStatus();
    }

    updateStreamingStatus() {
        const statusIndicator = document.querySelector('.logs-status .status-indicator');
        const statusText = document.querySelector('.logs-status span:last-child');

        if (statusIndicator) {
            statusIndicator.className = `status-indicator ${this.isStreaming ? 'streaming' : ''}`;
        }

        if (statusText) {
            statusText.textContent = this.isStreaming ? 'Streaming' : 'Stopped';
        }
    }

    autoScrollToBottom() {
        const logsContent = document.getElementById('logsContent');
        if (logsContent && this.autoScroll) {
            logsContent.scrollTop = logsContent.scrollHeight;
        }
    }

    cleanup() {
        if (this.isStreaming) {
            this.stopStreaming();
        }
        this.currentContainerId = null;
    }
}

// Export the component
window.LogsComponent = LogsComponent;