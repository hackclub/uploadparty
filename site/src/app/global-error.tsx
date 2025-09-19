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
        <main>
          <section>
            <h2>Something went wrong</h2>
            <p>{error?.message ?? "Unknown error"}</p>
            {error?.digest && <p>Error digest: {error.digest}</p>}
            <button type="button" onClick={() => reset()}>Try again</button>
          </section>
        </main>
      </body>
    </html>
  );
}
