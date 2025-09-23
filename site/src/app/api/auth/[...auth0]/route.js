// Basic Auth0 routes - simplified implementation for testing (hardened)
export async function GET(request, { params }) {
  const resolvedParams = await params;
  const { auth0 } = resolvedParams;
  const route = auth0[0]; // First segment after /api/auth/

  // Derive the current app origin from the incoming request to avoid env mismatches
  const appOrigin = (() => {
    try { return new URL(request.url).origin; } catch { return process.env.AUTH0_BASE_URL || ''; }
  })();

  // Determine allowed origin for redirects
  const allowedBase = process.env.AUTH0_BASE_URL || process.env.NEXT_PUBLIC_SITE_URL || appOrigin;
  const isHttps = allowedBase?.startsWith('https://');

  const cookie = (name, value, opts = {}) => {
    const parts = [
      `${name}=${value}`,
      'Path=/',
      'HttpOnly',
      `SameSite=${opts.sameSite || 'Lax'}`,
      opts.maxAge !== undefined ? `Max-Age=${opts.maxAge}` : undefined,
      // Set Secure for https origins; avoid on localhost/http to keep local dev working
      (opts.secure ?? isHttps) ? 'Secure' : undefined,
      opts.domain ? `Domain=${opts.domain}` : undefined
    ].filter(Boolean);
    return parts.join('; ');
  };

  const safeInternalPath = (p) => {
    if (!p || typeof p !== 'string') return '/app';
    try {
      // Disallow absolute URLs and protocol/scheme
      if (p.startsWith('http://') || p.startsWith('https://')) return '/app';
      // Ensure it starts with a slash and does not contain CRLF
      if (!p.startsWith('/')) p = '/' + p;
      if (/[\r\n]/.test(p)) return '/app';
      return p;
    } catch { return '/app'; }
  };

  try {
    switch (route) {
      case 'login': {
        // Get returnTo parameter to redirect user after login (restricted to internal path)
        const loginRequestUrl = new URL(request.url);
        const requestedReturnTo = loginRequestUrl.searchParams.get('returnTo') || '/app';
        const returnTo = safeInternalPath(requestedReturnTo);
        const nonce = (globalThis.crypto?.randomUUID?.() || Math.random().toString(36).slice(2));
        const stateWithReturnTo = JSON.stringify({ returnTo, nonce, ts: Date.now() });

        // Validate allowed base for redirect_uri
        const base = (() => {
          try { return new URL(allowedBase).origin; } catch { return appOrigin; }
        })();

        // Normalize issuer base (ensure scheme present)
        const rawIssuer = process.env.AUTH0_ISSUER_BASE_URL || '';
        let issuerBase = rawIssuer.trim();
        if (issuerBase && !/^https?:\/\//i.test(issuerBase)) {
          issuerBase = 'https://' + issuerBase;
        }
        let issuerOrigin = '';
        try {
          issuerOrigin = new URL(issuerBase).origin;
        } catch {
          // Fall back: use as-is; Response.redirect will still error if invalid
          issuerOrigin = issuerBase;
        }

        // Redirect to Auth0 login
        const loginUrl = `${issuerOrigin}/authorize?` +
          `response_type=code&` +
          `client_id=${encodeURIComponent(process.env.AUTH0_CLIENT_ID || '')}&` +
          `redirect_uri=${encodeURIComponent(base + '/api/auth/callback')}&` +
          `scope=openid%20profile%20email&` +
          `state=${encodeURIComponent(stateWithReturnTo)}`;

        return Response.redirect(loginUrl);
      }

      case 'logout': {
        // Clear session cookie and redirect to Auth0 logout
        const base = (() => {
          try { return new URL(allowedBase).origin; } catch { return appOrigin; }
        })();
        // Normalize issuer base for logout as well
        const rawIssuer = process.env.AUTH0_ISSUER_BASE_URL || '';
        let issuerBase = rawIssuer.trim();
        if (issuerBase && !/^https?:\/\//i.test(issuerBase)) {
          issuerBase = 'https://' + issuerBase;
        }
        let issuerOrigin = '';
        try {
          issuerOrigin = new URL(issuerBase).origin;
        } catch {
          issuerOrigin = issuerBase;
        }
        const logoutUrl = `${issuerOrigin}/v2/logout?` +
          `client_id=${encodeURIComponent(process.env.AUTH0_CLIENT_ID || '')}&` +
          `returnTo=${encodeURIComponent(base)}`;

        const logoutResponse = new Response(null, {
          status: 302,
          headers: {
            'Location': logoutUrl,
            'Set-Cookie': cookie('auth_session', '', { maxAge: 0 })
          }
        });

        return logoutResponse;
      }

      case 'callback': {
        // Handle the callback from Auth0
        const url = new URL(request.url);
        const code = url.searchParams.get('code');
        const state = url.searchParams.get('state');
        const error = url.searchParams.get('error');
        const errorDescription = url.searchParams.get('error_description');

        if (error) {
          // Avoid reflecting full details in prod
          const msg = process.env.NODE_ENV === 'development' ? `${error} - ${errorDescription}` : 'Authentication failed';
          return new Response(`Auth0 Error: ${msg}`, { status: 400 });
        }

        if (!code) {
          return new Response('Authorization code missing. Please try logging in again.', { status: 400 });
        }

        // Extract and validate returnTo from state (no open redirects)
        let redirectPath = '/app';
        try {
          if (state) {
            const stateData = JSON.parse(decodeURIComponent(state));
            redirectPath = safeInternalPath(stateData.returnTo);
          }
        } catch {
          // ignore and use default
        }

        // Build response with secure session cookie
        const base = (() => {
          try { return new URL(allowedBase).origin; } catch { return appOrigin; }
        })();
        const response = new Response(null, {
          status: 302,
          headers: {
            'Location': base + redirectPath,
            'Set-Cookie': cookie('auth_session', 'logged_in', { maxAge: 60 * 60 * 24 })
          }
        });

        return response;
      }

      case 'me': {
        // Return user profile (placeholder)
        return new Response(JSON.stringify({ user: null }), {
          headers: { 'Content-Type': 'application/json' }
        });
      }

      default:
        return new Response('Auth endpoint not found', { status: 404 });
    }
  } catch (error) {
    // Avoid leaking internals
    const msg = process.env.NODE_ENV === 'development' ? (error?.message || 'Unknown error') : 'Internal error';
    return new Response(`Auth error: ${msg}`, { status: 500 });
  }
}