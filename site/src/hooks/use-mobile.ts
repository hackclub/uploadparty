"use client";
import { useEffect, useState } from "react";

/**
 * useMobile
 * - Tracks whether the viewport is considered mobile via a CSS media query.
 * - Server-safe: starts with undefined and resolves on the client.
 */
export function useMobile(breakpointPx: number = 768) {
  const [isMobile, setIsMobile] = useState<boolean | undefined>(undefined);

  useEffect(() => {
    if (typeof window === "undefined") return;

    const mql = window.matchMedia(`(max-width: ${breakpointPx}px)`);
    const onChange = (e: MediaQueryListEvent | MediaQueryList) => {
      setIsMobile("matches" in e ? e.matches : (e as MediaQueryList).matches);
    };

    // Initialize and subscribe
    onChange(mql);
    const listener = (e: MediaQueryListEvent) => onChange(e);
    mql.addEventListener?.("change", listener);

    return () => {
      mql.removeEventListener?.("change", listener);
    };
  }, [breakpointPx]);

  return isMobile;
}
