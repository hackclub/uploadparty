import axios from 'axios';

const baseURL = process.env.NEXT_PUBLIC_BACKEND_URL || 'http://localhost:8080';

const axiosInstance = axios.create({
  baseURL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
});

// Request interceptor for adding auth token
axiosInstance.interceptors.request.use(
  (config) => {
    // Get token from localStorage or wherever you store it
    const token = typeof window !== 'undefined' ? localStorage.getItem('authToken') : null;

    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor for handling errors
axiosInstance.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (error.response) {
      // Server responded with error status
      console.error('API Error:', error.response.status, error.response.data);

      // Handle unauthorized errors
      if (error.response.status === 401) {
        // Clear token and redirect to login if needed
        if (typeof window !== 'undefined') {
          localStorage.removeItem('authToken');
          // Optionally redirect to login page
          // window.location.href = '/login';
        }
      }
    } else if (error.request) {
      // Request made but no response received
      console.error('Network Error:', error.request);
    } else {
      // Error in request setup
      console.error('Request Error:', error.message);
    }

    return Promise.reject(error);
  }
);

export default axiosInstance;
