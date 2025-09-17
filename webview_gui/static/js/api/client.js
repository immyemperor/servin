/**
 * API Client for Docker Desktop GUI
 * Handles all HTTP requests to the backend API
 */

class APIClient {
    constructor(baseUrl = '') {
        this.baseUrl = baseUrl;
    }

    /**
     * Generic API request method
     */
    async request(endpoint, options = {}) {
        const url = `${this.baseUrl}${endpoint}`;
        const config = {
            headers: {
                'Content-Type': 'application/json',
            },
            ...options,
        };

        try {
            const response = await fetch(url, config);
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            const contentType = response.headers.get('content-type');
            if (contentType && contentType.includes('application/json')) {
                return await response.json();
            }
            
            return await response.text();
        } catch (error) {
            console.error(`API request failed: ${endpoint}`, error);
            throw error;
        }
    }

    /**
     * Container API endpoints
     */
    async getContainers() {
        return await this.request('/api/containers');
    }

    async getContainerDetails(containerId) {
        return await this.request(`/api/containers/${containerId}/details`);
    }

    async startContainer(containerId) {
        return await this.request(`/api/containers/${containerId}/start`, {
            method: 'POST'
        });
    }

    async stopContainer(containerId) {
        return await this.request(`/api/containers/${containerId}/stop`, {
            method: 'POST'
        });
    }

    async restartContainer(containerId) {
        return await this.request(`/api/containers/${containerId}/restart`, {
            method: 'POST'
        });
    }

    async removeContainer(containerId) {
        return await this.request(`/api/containers/${containerId}/remove`, {
            method: 'DELETE'
        });
    }

    async getContainerLogs(containerId) {
        return await this.request(`/api/containers/${containerId}/logs`);
    }

    async getContainerFiles(containerId, path = '/') {
        return await this.request(`/api/containers/${containerId}/files?path=${encodeURIComponent(path)}`);
    }

    async getContainerEnvironment(containerId) {
        return await this.request(`/api/containers/${containerId}/env`);
    }

    /**
     * Image API endpoints
     */
    async getImages() {
        return await this.request('/api/images');
    }

    async pullImage(imageName) {
        return await this.request('/api/images/pull', {
            method: 'POST',
            body: JSON.stringify({ image: imageName })
        });
    }

    async removeImage(imageId) {
        return await this.request(`/api/images/${imageId}/remove`, {
            method: 'DELETE'
        });
    }

    /**
     * Volume API endpoints
     */
    async getVolumes() {
        return await this.request('/api/volumes');
    }

    async createVolume(volumeName) {
        return await this.request('/api/volumes/create', {
            method: 'POST',
            body: JSON.stringify({ name: volumeName })
        });
    }

    async removeVolume(volumeName) {
        return await this.request(`/api/volumes/${volumeName}/remove`, {
            method: 'DELETE'
        });
    }

    /**
     * System API endpoints
     */
    async getSystemInfo() {
        return await this.request('/api/system/info');
    }

    async checkConnection() {
        return await this.request('/api/system/info');
    }
}

// Export for use in other modules
window.APIClient = APIClient;