# (public) route group

This directory contains the public, pre‑authentication experience. Pages here must be accessible without a session and generally render as Server Components.

- URL: The root route `/` is implemented by re-exporting the component from `./(public)/page.js` via `../page.js`.
- Purpose: Marketing/landing, docs, and any content visible before sign-in.
- Notes: Keep bundles small; avoid `"use client"` unless interaction is required.
