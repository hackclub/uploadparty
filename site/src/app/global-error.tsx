"use client";

import React from "react";

export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <html lang="en">
      <body>
        <div className="min-h-screen flex items-center justify-center p-6">
          <div className="max-w-lg w-full rounded-lg border border-red-200 bg-red-50 text-red-800 p-6">
            <h2 className="text-lg font-semibold mb-2">Something went wrong</h2>
            <p className="text-sm opacity-90 mb-4">{error?.message ?? "Unknown error"}</p>
            {error?.digest && (
              <p className="text-xs opacity-70 mb-4">Error digest: {error.digest}</p>
            )}
            <button
              type="button"
              onClick={() => reset()}
              className="inline-flex items-center gap-2 rounded-md bg-red-600 px-3 py-1.5 text-white hover:bg-red-500 active:bg-red-700 transition-colors"
            >
              Try again
            </button>
          </div>
        </div>
      </body>
    </html>
  );
}
