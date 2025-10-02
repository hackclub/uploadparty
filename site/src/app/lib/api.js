// API client for backend communication
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

class ApiClient {
  constructor() {
    this.baseURL = API_BASE_URL;
  }

  // Helper method to get auth token from secure cookie
  getAuthToken() {
    // Token is handled by Auth0 SDK automatically via httpOnly cookies
    // This method is kept for backward compatibility but should not be used
    // for actual authentication - Auth0 handles this automatically
    return null;
  }

  // Generic fetch wrapper with error handling
  async fetch(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const token = this.getAuthToken();
    
    const config = {
      headers: {
        'Content-Type': 'application/json',
        ...(token && { Authorization: `Bearer ${token}` }),
        ...options.headers,
      },
      ...options,
    };

    try {
      const response = await fetch(url, config);
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || `HTTP ${response.status}: ${response.statusText}`);
      }

      // Handle empty responses
      const contentType = response.headers.get('content-type');
      if (contentType && contentType.includes('application/json')) {
        return await response.json();
      }
      
      return response;
    } catch (error) {
      console.error(`API Error (${endpoint}):`, error);
      
      // Check if it's a network error (backend not available)
      if (error.name === 'TypeError' && error.message.includes('fetch')) {
        const serverError = new Error('Server unavailable');
        serverError.isServerError = true;
        throw serverError;
      }
      
      throw error;
    }
  }

  // Health check
  async health() {
    return this.fetch('/health');
  }

  // Auth endpoints
  async register(userData) {
    return this.fetch('/auth/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
  }

  async login(credentials) {
    return this.fetch('/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
  }

  // App endpoints (protected)
  async getProjects() {
    return this.fetch('/api/v1/app/projects');
  }

  async getProjectPlugins(projectId) {
    return this.fetch(`/api/v1/app/projects/${projectId}/plugins`);
  }

  async markProjectComplete(projectId) {
    return this.fetch(`/api/v1/app/projects/${projectId}/complete`, {
      method: 'PATCH',
    });
  }

  // Ingest endpoints (for VST/plugin data)
  async upsertProject(projectData) {
    return this.fetch('/api/v1/ingest/projects', {
      method: 'POST',
      body: JSON.stringify(projectData),
    });
  }

  async upsertPluginForProject(projectId, pluginData) {
    return this.fetch(`/api/v1/ingest/projects/${projectId}/plugins`, {
      method: 'POST',
      body: JSON.stringify(pluginData),
    });
  }

  // Public endpoints
  async getPublicProfile(handle) {
    return this.fetch(`/profiles/${handle}`);
  }
}

// Create singleton instance
const apiClient = new ApiClient();

export default apiClient;

// Named exports for specific API calls
export const {
  health,
  register,
  login,
  getProjects,
  getProjectPlugins,
  markProjectComplete,
  upsertProject,
  upsertPluginForProject,
  getPublicProfile,
} = apiClient;