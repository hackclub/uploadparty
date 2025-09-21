# App Router: Route map

Public (pre‑auth)
- `/` → app/(public)/page.js (landing)
- `/overview` → app/(public)/overview/page.js
- `/prizes` → app/(public)/prizes/page.js
- `/faq` → app/(public)/faq/page.js
- 404 within public → app/(public)/not-found.js (redirects to `/`)

Authenticated (post‑auth)
- Example: `/dashboard` → app/(app)/dashboard/page.js (if/when added)
- Example: `/settings` → app/(app)/settings/page.js (if/when added)
- 404 within authenticated → app/(app)/not-found.js (message + helpful links)

SEO
- `/sitemap.xml` → app/sitemap.js
- `/robots.txt` → app/robots.js

Notes
- Route groups in parentheses are not part of the URL.
- Keep this file updated when you add or remove routes.
