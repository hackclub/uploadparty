/**
 * @jest-environment jsdom
 */
import { render, screen, waitFor } from '@testing-library/react';
import { act } from 'react';
import secureApiClient, { APIError, NetworkError } from '../lib/secure-api';
import { useHealthCheck, useProjects } from '../hooks/useApi';

// Mock fetch globally
global.fetch = jest.fn();

// Mock Auth0
jest.mock('@auth0/nextjs-auth0', () => ({
  getAccessToken: jest.fn(),
}));

// Test component for useHealthCheck hook
function TestHealthComponent() {
  const { loading, error, checkHealth } = useHealthCheck();
  
  return (
    <div>
      <button onClick={checkHealth} disabled={loading}>
        {loading ? 'Checking...' : 'Check Health'}
      </button>
      {error && <div data-testid="error">{error.message}</div>}
    </div>
  );
}

describe('SecureApiClient', () => {
  beforeEach(() => {
    fetch.mockClear();
    console.error = jest.fn(); // Suppress console.error in tests
  });

  describe('health check', () => {
    it('should return health data on successful response', async () => {
      const mockHealthData = {
        status: 'healthy',
        timestamp: '2023-01-01T00:00:00Z',
        services: { database: 'healthy' }
      };

      fetch.mockResolvedValueOnce({
        ok: true,
        status: 200,
        headers: new Map([['content-type', 'application/json']]),
        json: async () => mockHealthData,
      });

      const result = await secureApiClient.health();
      
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('/health'),
        expect.objectContaining({
          credentials: 'include',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
        })
      );
      expect(result).toEqual(mockHealthData);
    });

    it('should throw APIError on server error', async () => {
      const errorResponse = {
        error: 'Service unavailable',
        code: 'SERVICE_UNAVAILABLE'
      };

      fetch.mockResolvedValueOnce({
        ok: false,
        status: 503,
        headers: new Map([['content-type', 'application/json']]),
        json: async () => errorResponse,
      });

      await expect(secureApiClient.health()).rejects.toThrow(APIError);
      await expect(secureApiClient.health()).rejects.toThrow('Service unavailable');
    });

    it('should throw NetworkError on fetch failure', async () => {
      fetch.mockRejectedValueOnce(new TypeError('Failed to fetch'));

      await expect(secureApiClient.health()).rejects.toThrow(NetworkError);
    });

    it('should timeout after specified duration', async () => {
      // Mock a slow response
      fetch.mockImplementationOnce(() => 
        new Promise(resolve => setTimeout(resolve, 10000))
      );

      await expect(secureApiClient.health(100)).rejects.toThrow(NetworkError);
      await expect(secureApiClient.health(100)).rejects.toThrow('Health check timeout');
    });
  });

  describe('retry logic', () => {
    it('should retry on 500 errors', async () => {
      // First call fails, second succeeds
      fetch
        .mockResolvedValueOnce({
          ok: false,
          status: 500,
          headers: new Map([['content-type', 'application/json']]),
          json: async () => ({ error: 'Internal error' }),
        })
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          headers: new Map([['content-type', 'application/json']]),
          json: async () => ({ status: 'healthy' }),
        });

      const result = await secureApiClient.health();
      
      expect(fetch).toHaveBeenCalledTimes(2);
      expect(result).toEqual({ status: 'healthy' });
    });

    it('should not retry on 401 errors', async () => {
      fetch.mockResolvedValueOnce({
        ok: false,
        status: 401,
        headers: new Map([['content-type', 'application/json']]),
        json: async () => ({ error: 'Unauthorized' }),
      });

      await expect(secureApiClient.health()).rejects.toThrow(APIError);
      expect(fetch).toHaveBeenCalledTimes(1);
    });
  });
});

describe('useHealthCheck hook', () => {
  it('should handle loading states correctly', async () => {
    fetch.mockImplementationOnce(
      () => new Promise(resolve => 
        setTimeout(() => resolve({
          ok: true,
          status: 200,
          headers: new Map([['content-type', 'application/json']]),
          json: async () => ({ status: 'healthy' }),
        }), 100)
      )
    );

    render(<TestHealthComponent />);
    
    const button = screen.getByRole('button');
    expect(button).toHaveTextContent('Check Health');
    
    act(() => {
      button.click();
    });
    
    expect(button).toHaveTextContent('Checking...');
    expect(button).toBeDisabled();
    
    await waitFor(() => {
      expect(button).toHaveTextContent('Check Health');
      expect(button).not.toBeDisabled();
    });
  });

  it('should display error messages', async () => {
    fetch.mockRejectedValueOnce(new TypeError('Failed to fetch'));

    render(<TestHealthComponent />);
    
    const button = screen.getByRole('button');
    
    act(() => {
      button.click();
    });
    
    await waitFor(() => {
      const errorElement = screen.getByTestId('error');
      expect(errorElement).toHaveTextContent('Unable to connect to server');
    });
  });
});

describe('API Error Handling', () => {
  it('should create APIError with correct properties', () => {
    const error = new APIError('Test error', 400, 'BAD_REQUEST', { field: 'invalid' });
    
    expect(error.name).toBe('APIError');
    expect(error.message).toBe('Test error');
    expect(error.status).toBe(400);
    expect(error.code).toBe('BAD_REQUEST');
    expect(error.details).toEqual({ field: 'invalid' });
  });

  it('should create NetworkError with correct properties', () => {
    const error = new NetworkError('Connection failed');
    
    expect(error.name).toBe('NetworkError');
    expect(error.message).toBe('Connection failed');
    expect(error.isNetworkError).toBe(true);
  });
});