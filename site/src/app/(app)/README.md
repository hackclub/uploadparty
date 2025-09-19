# Authenticated group: app/(app)

Purpose:
- All authenticated (post‑login) routes live here.
- Assumes a valid user session/JWT; enforce auth via middleware or layout guards.

Conventions:
- Use a shared layout at app/(app)/layout.js for nav/shell.
- Keep server components by default; promote to client components only for interactive UI.

Routing examples:
- `/(app)/dashboard/page.js` → `/dashboard`
- `/(app)/settings/page.js` → `/settings`

Error handling:
- not-found.js — Renders a “Page not found” message for unknown routes in the authenticated area.

Notes:
- Coordinate with backend auth endpoints under /auth and protected /api/v1 routes.
- Respect CORS and include Authorization headers for API requests.
