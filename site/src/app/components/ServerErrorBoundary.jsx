'use client';

import { useState, useEffect } from 'react';
import ServerError from './ServerError';

export default function ServerErrorBoundary({ children }) {
  const [hasServerError, setHasServerError] = useState(false);

  useEffect(() => {
    const handleUnhandledRejection = (event) => {
      if (event.reason?.isServerError) {
        setHasServerError(true);
        event.preventDefault();
      }
    };

    window.addEventListener('unhandledrejection', handleUnhandledRejection);
    return () => window.removeEventListener('unhandledrejection', handleUnhandledRejection);
  }, []);

  if (hasServerError) {
    return <ServerError />;
  }

  return children;
}