/**
 * Misc utilities for hooks and related logic.
 */

/**
 * Basic user agent mobile detector (best-effort; prefer responsive design).
 */
export function isMobileUserAgent(ua: string): boolean {
  const s = ua.toLowerCase();
  return /iphone|ipod|ipad|android|blackberry|iemobile|opera mini/.test(s);
}

export function noop(): void {
  /* no-op */
}
