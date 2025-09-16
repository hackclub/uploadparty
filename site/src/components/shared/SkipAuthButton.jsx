"use client";
import { useMemo } from "react";
import { useRouter } from "next/navigation";

/**
 * SkipAuthButton
 * - Fixed small button in the top-right corner to quickly skip OAuth during development.
 * - Controlled by the NEXT_PUBLIC_ENABLE_SKIP_AUTH env variable.
 * - To show the button, set NEXT_PUBLIC_ENABLE_SKIP_AUTH=true and rebuild.
 * - To hide it, set the flag to false or remove the component from layout.
 */
export default function SkipAuthButton() {
  const router = useRouter();
  // Env flags are inlined at build time for client components in Next.js
  const enabled = useMemo(() => process.env.NEXT_PUBLIC_ENABLE_SKIP_AUTH === "true", []);

  if (!enabled) return null;

  return (
    <button
      type="button"
      onClick={() => router.push("/dashboard")}
      className="fixed top-3 right-3 z-[60] rounded-full bg-emerald-600 text-white text-xs font-medium px-3 py-2 shadow hover:bg-emerald-500 active:bg-emerald-700 transition-colors"
      aria-label="Skip authentication and go to dashboard"
      title="Skip authentication (dev)"
    >
      Skip Auth
    </button>
  );
}
