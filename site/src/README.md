# Site/src overview

This directory contains the Next.js application source. We follow the App Router conventions with a clear separation between public (pre‑auth) and authenticated areas.

Structure at a glance:
- app/ — App Router entry and routes
  - (public)/ — Public, pre‑auth pages (marketing/landing, overview, prizes, FAQ)
  - (app)/ — Authenticated application area (dashboard, settings, etc.)
  - page.js — Root route that re‑exports (public)/page.js
  - globals.css — Global styles (framework‑agnostic)
  - robots.js, sitemap.js — SEO endpoints
  - global-error.tsx — Global error boundary
- components/ — (optional) shared UI components (add as needed)
- lib/ — (optional) utilities/helpers (add as needed)

Conventions:
- Use server components by default. Add "use client" only where interactivity is required.
- Keep public content inside app/(public). Keep authenticated content inside app/(app).
- Each major public section also has its own dedicated page for SEO: /overview, /prizes, /faq.
- If you add new routes, update app/ROUTES.md to keep the map current.

Adding new code:
- New public route: place under app/(public)/<route>/page.js
- New authenticated route: place under app/(app)/<route>/page.js (guarded by your chosen auth mechanism)
- Shared UI: create a component under site/src/components and import from routes as needed.
