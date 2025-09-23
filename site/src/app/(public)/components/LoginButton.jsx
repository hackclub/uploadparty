'use client';

import { useState, useEffect } from 'react';

export default function LoginButton() {
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    // Simple check for user session - we'll improve this later
    const checkAuth = async () => {
      try {
        // For now, just check if we have any auth-related cookies or localStorage
        const hasAuth = localStorage.getItem('auth0_user') || document.cookie.includes('auth0');
        
        if (hasAuth) {
          setUser({ name: 'User', email: 'user@example.com' }); // Placeholder
        }
        
        setIsLoading(false);
      } catch (err) {
        setError(err.message);
        setIsLoading(false);
      }
    };

    checkAuth();
  }, []);

  if (isLoading) {
    return (
      <div className="flex flex-col items-center gap-4">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <p>Loading...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center gap-4">
        <p className="text-red-600">Error loading auth state</p>
        <a 
          href="/api/auth/login" 
          className="bg-green-600 hover:bg-green-700 text-white px-8 py-4 rounded-lg font-bold text-lg transition-colors"
        >
          Sign In / Register
        </a>
      </div>
    );
  }

  if (user) {
    return (
      <div className="flex flex-col items-center gap-4">
        <p>Welcome, {user.name}!</p>
        <div className="flex gap-4">
          <a 
            href="/app" 
            className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-medium transition-colors"
          >
            Go to App
          </a>
          <a 
            href="/api/auth/logout" 
            className="bg-gray-600 hover:bg-gray-700 text-white px-6 py-3 rounded-lg font-medium transition-colors"
          >
            Logout
          </a>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col items-center gap-4">
      <h3>Ready to join UploadParty?</h3>
      <p>Sign in to access the platform and start uploading your beats!</p>
      <a 
        href="/api/auth/login" 
        className="bg-green-600 hover:bg-green-700 text-white px-8 py-4 rounded-lg font-bold text-lg transition-colors"
      >
        Sign In / Register
      </a>
    </div>
  );
}