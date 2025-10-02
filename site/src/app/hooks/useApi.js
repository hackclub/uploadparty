'use client';

import { useState, useCallback } from 'react';
import secureApiClient, { APIError, NetworkError } from '../lib/secure-api';
import { logger } from '../lib/logger';

// Custom hook for API calls with loading states and error handling
export function useApi() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const execute = useCallback(async (apiCall, onSuccess, onError) => {
    setLoading(true);
    setError(null);
    
    const startTime = Date.now();
    
    try {
      const result = await apiCall();
      const duration = Date.now() - startTime;
      
      logger.logApiCall(
        'UNKNOWN', // Method not available from hook
        'UNKNOWN', // URL not available from hook
        200,
        duration,
        { success: true }
      );
      
      if (onSuccess) {
        onSuccess(result);
      }
      
      return result;
    } catch (err) {
      const duration = Date.now() - startTime;
      
      if (err instanceof APIError) {
        logger.logApiCall(
          'UNKNOWN',
          'UNKNOWN', 
          err.status,
          duration,
          { error: err.message, code: err.code }
        );
        setError({
          type: 'api',
          message: err.message,
          status: err.status,
          code: err.code,
          details: err.details,
        });
      } else if (err instanceof NetworkError) {
        logger.error('Network error in API call', {
          error: err.message,
          duration,
        });
        setError({
          type: 'network',
          message: 'Unable to connect to server. Please check your connection.',
        });
      } else {
        logger.error('Unexpected error in API call', {
          error: err.message,
          stack: err.stack,
          duration,
        });
        setError({
          type: 'unknown',
          message: 'An unexpected error occurred. Please try again.',
        });
      }
      
      if (onError) {
        onError(err);
      }
      
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  return {
    loading,
    error,
    execute,
    clearError,
  };
}

// Specific hooks for common API operations
export function useHealthCheck() {
  const { loading, error, execute, clearError } = useApi();

  const checkHealth = useCallback(
    () => execute(() => secureApiClient.health()),
    [execute]
  );

  return {
    loading,
    error,
    checkHealth,
    clearError,
  };
}

export function useProjects() {
  const { loading, error, execute, clearError } = useApi();
  const [projects, setProjects] = useState([]);

  const fetchProjects = useCallback(
    (req, res) => execute(
      () => secureApiClient.getProjects(req, res),
      setProjects
    ),
    [execute]
  );

  const markComplete = useCallback(
    (projectId, req, res) => execute(
      () => secureApiClient.markProjectComplete(projectId, req, res),
      () => {
        // Update local state
        setProjects(prev => prev.map(p => 
          p.id === projectId ? { ...p, completed: true } : p
        ));
      }
    ),
    [execute]
  );

  return {
    projects,
    loading,
    error,
    fetchProjects,
    markComplete,
    clearError,
  };
}

export function usePublicProfile() {
  const { loading, error, execute, clearError } = useApi();
  const [profile, setProfile] = useState(null);

  const fetchProfile = useCallback(
    (handle) => execute(
      () => secureApiClient.getPublicProfile(handle),
      setProfile
    ),
    [execute]
  );

  return {
    profile,
    loading,
    error,
    fetchProfile,
    clearError,
  };
}