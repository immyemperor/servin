/**
 * Terminal Component
 * Handles container exec sessions and terminal interactions
 */

class Terminal {
    constructor(apiClient, socketManager) {
        this.apiClient = apiClient;
        this.socketManager = socketManager;
        this.currentContainerId = null;
        this.isConnected = false;
        this.commandHistory = [];
        this.historyIndex = -1;
        
        this.init();
    }

    init() {
        this.setupSocketHandlers();
    }

    setupSocketHandlers() {
        this.socketManager.on('exec_started', (data) => this.handleExecStarted(data));
        this.socketManager.on('exec_stopped', (data) => this.handleExecStopped(data));
        this.socketManager.on('exec_output', (data) => this.handleExecOutput(data));
    }

    setupTerminal(containerId) {
        this.currentContainerId = containerId;
        
        const execTerminal = document.getElementById('execTerminal');
        if (!execTerminal) return;

        // Simple placeholder for direct connection
        if (!execTerminal.querySelector('.terminal-placeholder')) {
            execTerminal.innerHTML = `
                <div class="terminal-placeholder">
                    <i class="fas fa-terminal"></i>
                    <p>Connecting to container shell...</p>
                </div>
            `;
        }

        // Set up event listeners for terminal input
        this.setupTerminalControls();
        
        // Automatically connect when terminal is set up
        this.connect();
    }

    setupTerminalControls() {
        const terminalInput = document.getElementById('terminalInput');

        if (terminalInput) {
            terminalInput.addEventListener('keydown', (e) => this.handleTerminalInput(e));
        }
    }

    connect() {
        if (!this.currentContainerId) return;

        // Use bash as default shell for running containers
        const shell = '/bin/bash';

        // Update placeholder to show connecting status
        const placeholder = document.getElementById('terminalPlaceholder');
        if (placeholder) {
            placeholder.innerHTML = `
                <i class="fas fa-spinner fa-spin"></i>
                <p>Connecting to container shell...</p>
            `;
        }

        // Start exec session via WebSocket
        this.socketManager.emit('start_exec', {
            container_id: this.currentContainerId,
            shell: shell
        });
    }

    disconnect() {
        if (!this.currentContainerId) return;

        // Stop exec session via WebSocket
        this.socketManager.emit('stop_exec', {
            container_id: this.currentContainerId
        });
    }

    handleExecStarted(data) {
        console.log('Terminal: Exec session started:', data);
        this.isConnected = true;

        const terminalPlaceholder = document.getElementById('terminalPlaceholder');
        const terminalOutput = document.getElementById('terminalOutput');
        const terminalInput = document.getElementById('terminalInput');

        if (terminalPlaceholder) {
            terminalPlaceholder.style.display = 'none';
        }

        if (terminalOutput) {
            terminalOutput.style.display = 'block';
        }

        if (terminalInput) {
            terminalInput.disabled = false;
            terminalInput.focus();
        }

        // Simple welcome message
        this.addTerminalLine('system', 'Shell session started. Type commands and press Enter.');
    }

    handleExecStopped(data) {
        console.log('Terminal: Exec session stopped:', data);
        this.isConnected = false;

        const terminalPlaceholder = document.getElementById('terminalPlaceholder');
        const terminalOutput = document.getElementById('terminalOutput');
        const terminalInput = document.getElementById('terminalInput');

        if (terminalPlaceholder) {
            terminalPlaceholder.style.display = 'block';
            terminalPlaceholder.innerHTML = `
                <i class="fas fa-exclamation-triangle"></i>
                <p>Terminal session ended</p>
            `;
        }

        if (terminalOutput) {
            terminalOutput.style.display = 'none';
        }

        if (terminalInput) {
            terminalInput.disabled = true;
        }

        this.addTerminalLine('system', 'Terminal session ended.');
    }

    handleExecOutput(data) {
        if (data.container_id !== this.currentContainerId) return;

        const terminalContent = document.getElementById('terminalContent');
        const terminalPrompt = document.getElementById('terminalPrompt');

        if (!terminalContent) return;

        if (data.type === 'prompt') {
            // Update prompt with proper format: user@hostname:path$
            if (terminalPrompt) {
                terminalPrompt.textContent = data.output || 'root@container:/# ';
            }
        } else if (data.type === 'output') {
            // Add command output to terminal content
            this.addCommandOutput(data.output);
        } else if (data.type === 'command') {
            // Show the executed command in terminal
            this.addCommandLine(data.output);
        } else {
            // Default: add any output to terminal
            this.addTerminalLine(data.type || 'output', data.output);
        }
    }

    handleTerminalInput(e) {
        if (e.key === 'Enter') {
            e.preventDefault();
            this.sendCommand();
        } else if (e.key === 'ArrowUp') {
            e.preventDefault();
            this.navigateHistory(-1);
        } else if (e.key === 'ArrowDown') {
            e.preventDefault();
            this.navigateHistory(1);
        } else if (e.key === 'Tab') {
            e.preventDefault();
            // Tab completion could be implemented here
        }
    }

    sendCommand() {
        const terminalInput = document.getElementById('terminalInput');
        if (!terminalInput || !this.isConnected) return;

        const command = terminalInput.value.trim();
        if (!command) return;

        // Add to command history
        this.commandHistory.push(command);
        this.historyIndex = this.commandHistory.length;

        // Display the command in terminal with proper prompt
        this.addCommandLine(command);

        // Send command via WebSocket
        this.socketManager.emit('exec_input', {
            container_id: this.currentContainerId,
            command: command
        });

        // Clear input
        terminalInput.value = '';
    }

    navigateHistory(direction) {
        const terminalInput = document.getElementById('terminalInput');
        if (!terminalInput || this.commandHistory.length === 0) return;

        this.historyIndex += direction;

        if (this.historyIndex < 0) {
            this.historyIndex = 0;
        } else if (this.historyIndex >= this.commandHistory.length) {
            this.historyIndex = this.commandHistory.length;
            terminalInput.value = '';
            return;
        }

        terminalInput.value = this.commandHistory[this.historyIndex] || '';
    }

    addCommandLine(command) {
        const terminalContent = document.getElementById('terminalContent');
        if (!terminalContent) return;

        // Add command line with prompt
        const commandLine = document.createElement('div');
        commandLine.className = 'terminal-line terminal-command';
        commandLine.innerHTML = `<span class="terminal-prompt-text">root@container:/# </span><span class="terminal-command-text">${this.escapeHtml(command)}</span>`;
        
        terminalContent.appendChild(commandLine);
        this.autoScroll();
    }

    addCommandOutput(output) {
        const terminalContent = document.getElementById('terminalContent');
        if (!terminalContent) return;

        if (output && output.trim()) {
            const outputLine = document.createElement('div');
            outputLine.className = 'terminal-line terminal-output';
            outputLine.textContent = output;
            
            terminalContent.appendChild(outputLine);
            this.autoScroll();
        }
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    autoScroll() {
        const terminalOutput = document.getElementById('terminalOutput');
        if (terminalOutput) {
            terminalOutput.scrollTop = terminalOutput.scrollHeight;
        }
    }

    addTerminalLine(type, content) {
        const terminalContent = document.getElementById('terminalContent');
        if (!terminalContent) return;

        const line = document.createElement('div');
        line.className = `terminal-line terminal-${type}`;
        line.textContent = content;

        terminalContent.appendChild(line);

        // Auto-scroll to bottom
        const terminalOutput = document.getElementById('terminalOutput');
        if (terminalOutput) {
            terminalOutput.scrollTop = terminalOutput.scrollHeight;
        }
    }

    clearTerminal() {
        const terminalContent = document.getElementById('terminalContent');
        if (terminalContent) {
            terminalContent.innerHTML = '';
        }
    }

    cleanup() {
        if (this.isConnected) {
            this.disconnect();
        }
        this.currentContainerId = null;
        this.commandHistory = [];
        this.historyIndex = -1;
    }
}

// Export the component
window.Terminal = Terminal;