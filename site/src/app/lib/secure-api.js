// Production-ready API client with Auth0 integration
import { getAccessToken } from '@auth0/nextjs-auth0';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// Enhanced error types for better error handling
export class APIError extends Error {
  constructor(message, status, code, details) {
    super(message);
    this.name = 'APIError';
    this.status = status;
    this.code = code;
    this.details = details;
  }
}

export class NetworkError extends Error {
  constructor(message) {
    super(message);
    this.name = 'NetworkError';
    this.isNetworkError = true;
  }
}

// Retry configuration
const RETRY_CONFIG = {
  maxRetries: 3,
  baseDelay: 1000, // 1 second
  maxDelay: 10000, // 10 seconds
  retryableStatuses: [408, 429, 500, 502, 503, 504],
};

// Exponential backoff delay calculation
function calculateDelay(attempt) {
  const delay = Math.min(
    RETRY_CONFIG.baseDelay * Math.pow(2, attempt),
    RETRY_CONFIG.maxDelay
  );
  // Add jitter to prevent thundering herd
  return delay + Math.random() * 1000;
}

// Sleep helper
function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

class SecureApiClient {
  constructor() {
    this.baseURL = API_BASE_URL;
  }

  // Server-side API calls with Auth0 token
  async serverFetch(endpoint, options = {}, req = null, res = null) {
    try {
      const { accessToken } = await getAccessToken(req, res);
      return this.fetch(endpoint, {
        ...options,
        headers: {
          ...options.headers,
          Authorization: `Bearer ${accessToken}`,
        },
      });
    } catch (error) {
      if (error.code === 'ERR_EXPIRED_ACCESS_TOKEN') {
        throw new APIError('Access token expired', 401, 'TOKEN_EXPIRED');
      }
      throw error;
    }
  }

  // Generic fetch with retry logic and enhanced error handling
  async fetch(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const config = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      credentials: 'include', // Include cookies for Auth0
      ...options,
    };

    let lastError;

    for (let attempt = 0; attempt <= RETRY_CONFIG.maxRetries; attempt++) {
      try {
        const response = await fetch(url, config);

        // Handle non-2xx responses
        if (!response.ok) {
          const errorData = await response.json().catch(() => ({}));
          
          // Check if this is a retryable error
          if (
            attempt < RETRY_CONFIG.maxRetries &&
            RETRY_CONFIG.retryableStatuses.includes(response.status)
          ) {
            const delay = calculateDelay(attempt);
            await sleep(delay);
            continue;
          }

          throw new APIError(
            errorData.error || `HTTP ${response.status}: ${response.statusText}`,
            response.status,
            errorData.code,
            errorData.details
          );
        }

        // Handle successful response
        const contentType = response.headers.get('content-type');
        if (contentType && contentType.includes('application/json')) {
          return await response.json();
        }
        
        return response;

      } catch (error) {
        lastError = error;

        // Don't retry API errors or authentication issues
        if (error instanceof APIError) {
          if (error.status === 401 || error.status === 403) {
            throw error;
          }
          if (attempt < RETRY_CONFIG.maxRetries && RETRY_CONFIG.retryableStatuses.includes(error.status)) {
            const delay = calculateDelay(attempt);
            await sleep(delay);
            continue;
          }
          throw error;
        }

        // Handle network errors
        if (error.name === 'TypeError' && error.message.includes('fetch')) {
          if (attempt < RETRY_CONFIG.maxRetries) {
            const delay = calculateDelay(attempt);
            await sleep(delay);
            continue;
          }
          throw new NetworkError('Unable to connect to server');
        }

        // Re-throw other errors immediately
        throw error;
      }
    }

    throw lastError;
  }

  // Health check with timeout
  async health(timeout = 5000) {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), timeout);

    try {
      const result = await this.fetch('/health', {
        signal: controller.signal,
      });
      clearTimeout(timeoutId);
      return result;
    } catch (error) {
      clearTimeout(timeoutId);
      if (error.name === 'AbortError') {
        throw new NetworkError('Health check timeout');
      }
      throw error;
    }
  }

  // App endpoints (protected) - Server-side versions
  async getProjects(req, res) {
    return this.serverFetch('/api/v1/app/projects', {}, req, res);
  }

  async getProjectPlugins(projectId, req, res) {
    return this.serverFetch(`/api/v1/app/projects/${projectId}/plugins`, {}, req, res);
  }

  async markProjectComplete(projectId, req, res) {
    return this.serverFetch(`/api/v1/app/projects/${projectId}/complete`, {
      method: 'PATCH',
    }, req, res);
  }

  // Public endpoints (no auth required)
  async getPublicProfile(handle) {
    return this.fetch(`/profiles/${handle}`);
  }
}

// Create singleton instance
const secureApiClient = new SecureApiClient();

export default secureApiClient;