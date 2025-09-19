"use client";
import React from "react";

export type ErrorReporterProps = {
  title?: string;
  error?: unknown;
  onRetry?: () => void;
  className?: string;
};

/**
 * ErrorReporter
 * - Lightweight UI component to surface an error message.
 * - Safe to use anywhere (client component).
 */
export default function ErrorReporter({
  title = "Something went wrong",
  error,
  onRetry,
  className = "",
}: ErrorReporterProps) {
  const message = React.useMemo(() => {
    if (!error) return null;
    if (typeof error === "string") return error;
    if (error instanceof Error) return error.message;
    try {
      return JSON.stringify(error);
    } catch (_) {
      return String(error);
    }
  }, [error]);

  return (
    <div
      role="alert"
      className={`rounded-lg border border-red-200 bg-red-50 text-red-800 p-4 text-sm ${className}`}
    >
      <div className="font-semibold mb-1">{title}</div>
      {message && <pre className="whitespace-pre-wrap break-words">{message}</pre>}
      {onRetry && (
        <button
          type="button"
          onClick={onRetry}
          className="mt-3 inline-flex items-center gap-2 rounded-md bg-red-600 px-3 py-1.5 text-white hover:bg-red-500 active:bg-red-700 transition-colors"
        >
          Try again
        </button>
      )}
    </div>
  );
}
