// Global setup for Playwright tests
export default async function globalSetup(config) {
  // Setup test environment
  console.log('Setting up Playwright test environment...');
  
  // Wait for backend to be ready
  const backendUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  
  console.log(`Waiting for backend at ${backendUrl}...`);
  
  // Poll backend health endpoint
  for (let i = 0; i < 30; i++) {
    try {
      const response = await fetch(`${backendUrl}/health`);
      if (response.ok) {
        console.log('Backend is ready!');
        break;
      }
    } catch (error) {
      // Backend not ready yet
    }
    
    if (i === 29) {
      throw new Error('Backend failed to start within timeout');
    }
    
    await new Promise(resolve => setTimeout(resolve, 2000)); // Wait 2 seconds
  }
  
  console.log('Global setup complete!');
}