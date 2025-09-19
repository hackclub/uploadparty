# Public group: app/(public)

Purpose:
- Houses all public (pre‑auth) routes and components.
- Intended for marketing/landing and other non‑authenticated content.

Key routes:
- `/` — Landing page (implemented in app/(public)/page.js, re‑exported by app/page.js)
- `/overview` — Dedicated Overview section page
- `/prizes` — Dedicated Prizes section page
- `/faq` — Dedicated FAQ & Footer section page

Structure:
- sections/ — Placeholder section components used by landing and dedicated pages
- overview/page.js — Renders the Overview section standalone
- prizes/page.js — Renders the Prizes section standalone
- faq/page.js — Renders the FAQ & Footer section standalone
- not-found.js — Redirects unknown public routes back to `/`

Notes:
- Keep components here as server components unless you need client interactivity.
- Update ../../sitemap.js if you add or remove public routes so SEO stays accurate.
