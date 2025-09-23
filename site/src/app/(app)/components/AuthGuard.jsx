'use client';

import { useState, useEffect } from 'react';

export default function AuthGuard({ children }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const checkAuth = () => {
      // Check for session cookie
      const hasAuthSession = document.cookie.includes('auth_session=logged_in');
      
      if (hasAuthSession) {
        setIsAuthenticated(true);
      } else {
        // Redirect to login with returnTo parameter
        window.location.href = '/api/auth/login?returnTo=' + encodeURIComponent(window.location.pathname);
      }
      
      setIsLoading(false);
    };

    checkAuth();
  }, []);

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p>Loading...</p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <p>Redirecting to login...</p>
        </div>
      </div>
    );
  }

  return children;
}