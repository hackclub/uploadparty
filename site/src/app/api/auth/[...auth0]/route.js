// Basic Auth0 routes - simplified implementation for testing
export async function GET(request, { params }) {
  const resolvedParams = await params;
  const { auth0 } = resolvedParams;
  const route = auth0[0]; // First segment after /api/auth/
  
  console.log('Auth route:', route);
  console.log('Full URL:', request.url);

  // Derive the current app origin from the incoming request to avoid env mismatches
  const appOrigin = (() => {
    try { return new URL(request.url).origin; } catch { return process.env.AUTH0_BASE_URL || ''; }
  })();
  
  try {
    switch (route) {
      case 'login':
        // Get returnTo parameter to redirect user after login
        const loginRequestUrl = new URL(request.url);
        const returnTo = loginRequestUrl.searchParams.get('returnTo') || '/app';
        const stateWithReturnTo = JSON.stringify({ 
          returnTo, 
          nonce: Math.random().toString(36).substring(7) 
        });
        
        // Redirect to Auth0 login
        const loginUrl = `${process.env.AUTH0_ISSUER_BASE_URL}/authorize?` + 
          `response_type=code&` +
          `client_id=${process.env.AUTH0_CLIENT_ID}&` +
          `redirect_uri=${encodeURIComponent(appOrigin + '/api/auth/callback')}&` +
          `scope=openid profile email&` +
          `state=${encodeURIComponent(stateWithReturnTo)}`;
        
        return Response.redirect(loginUrl);
        
      case 'logout':
        // Clear session cookie and redirect to Auth0 logout
        const logoutUrl = `${process.env.AUTH0_ISSUER_BASE_URL}/v2/logout?` +
          `client_id=${process.env.AUTH0_CLIENT_ID}&` +
          `returnTo=${encodeURIComponent(appOrigin)}`;
        
        const logoutResponse = new Response(null, {
          status: 302,
          headers: {
            'Location': logoutUrl,
            'Set-Cookie': 'auth_session=; Path=/; HttpOnly; SameSite=Lax; Max-Age=0'
          }
        });
        
        return logoutResponse;
        
      case 'callback':
        // Handle the callback from Auth0
        const url = new URL(request.url);
        const code = url.searchParams.get('code');
        const state = url.searchParams.get('state');
        const error = url.searchParams.get('error');
        const errorDescription = url.searchParams.get('error_description');
        
        console.log('Callback params:', { code: !!code, state, error, errorDescription });
        
        if (error) {
          console.error('Auth0 error:', error, errorDescription);
          return new Response(`Auth0 Error: ${error} - ${errorDescription}`, { status: 400 });
        }
        
        if (!code) {
          console.error('No authorization code received');
          return new Response('Authorization code missing. Please try logging in again.', { status: 400 });
        }
        
        // Extract returnTo from state
        let redirectPath = '/app'; // default
        try {
          if (state) {
            const stateData = JSON.parse(decodeURIComponent(state));
            redirectPath = stateData.returnTo || '/app';
          }
        } catch (e) {
          console.warn('Could not parse state, using default redirect path');
        }
        
        console.log('Auth successful, setting session and redirecting to:', redirectPath);
        
        // Create a redirect response with session cookie
        const response = new Response(null, {
          status: 302,
          headers: {
            'Location': appOrigin + redirectPath,
            'Set-Cookie': 'auth_session=logged_in; Path=/; HttpOnly; SameSite=Lax; Max-Age=86400'
          }
        });
        
        return response;
        
      case 'me':
        // Return user profile
        return new Response(JSON.stringify({ user: null }), {
          headers: { 'Content-Type': 'application/json' }
        });
        
      default:
        return new Response('Auth endpoint not found', { status: 404 });
    }
  } catch (error) {
    console.error('Auth error:', error);
    return new Response(`Auth error: ${error.message}`, { status: 500 });
  }
}