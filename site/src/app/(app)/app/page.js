'use client';

import { useState, useEffect } from 'react';
import apiClient from '../../lib/api';

export default function AppPage() {
  const [projects, setProjects] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    loadProjects();
    
    // Auto-refresh projects every 10 seconds
    const interval = setInterval(() => {
      loadProjects();
    }, 10000);

    return () => clearInterval(interval);
  }, []);

  const loadProjects = async () => {
    try {
      if (loading) setLoading(true);
      const data = await apiClient.getProjects();
      setProjects(data || []);
      setError(null);
    } catch (error) {
      console.error('Failed to load projects:', error);
      setError(error.message);
    } finally {
      if (loading) setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="p-6">
        <div className="flex items-center justify-center min-h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          <span className="ml-4">Loading your projects...</span>
        </div>
      </div>
    );
  }

  return (
    <main className="p-6">
      <h1 className="text-2xl font-semibold mb-4">Your Projects</h1>
      
      {error && (
        <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg">
          <p className="text-red-700">{error}</p>
        </div>
      )}

      {projects.length === 0 ? (
        <div className="text-center py-12">
          <div className="text-4xl mb-4">🎵</div>
          <h3 className="text-lg font-medium mb-2">No projects yet</h3>
          <p className="text-gray-600">Start creating music to see your projects here.</p>
        </div>
      ) : (
        <div className="space-y-4">
          {projects.map((project) => (
            <div key={project.id} className="bg-white p-6 rounded-lg shadow">
              <h3 className="text-lg font-semibold">{project.title}</h3>
              <p className="text-gray-600 text-sm mt-1">
                Created: {new Date(project.created_at).toLocaleDateString()}
              </p>
              <span className={`inline-block mt-2 px-2 py-1 rounded text-xs ${
                project.completed 
                  ? 'bg-green-100 text-green-800' 
                  : 'bg-yellow-100 text-yellow-800'
              }`}>
                {project.completed ? 'Complete' : 'In Progress'}
              </span>
            </div>
          ))}
        </div>
      )}
    </main>
  );
}
