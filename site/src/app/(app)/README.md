# (app) route group — authenticated area

This directory contains the signed‑in experience. Treat everything here as protected.

- URL: Routes inside `(app)` render without the segment name in the URL (e.g., `/(app)/dashboard` is `/dashboard`).
- Layout: `/(app)/layout.js` defines shared UI for authenticated pages.
- Auth: Add a server‑side guard either via Next.js middleware or in `/(app)/layout.js` to redirect unauthenticated users.
- Data fetching: Prefer Server Components; use Client Components only where interaction is needed.
