# App Router structure

We enforce a clear separation between the public (pre‑auth) experience and the authenticated application area:

- Public (pre‑auth)
  - Route group: `/(public)`
  - Entry point: `page.js` re‑exports `./(public)/page.js` for `/`
  - Purpose: Marketing/landing and other non‑auth content

- Authenticated application
  - Route group: `/(app)`
  - Example routes: `/(app)/dashboard` → `/dashboard`, `/(app)/app` → `/app`
  - Shared layout: `/(app)/layout.js`

See the README.md files inside each group for details.
