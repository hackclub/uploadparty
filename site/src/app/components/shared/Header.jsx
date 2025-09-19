"use client";
import Link from "next/link";

export default function Header() {
  return (
    <header className="fixed top-0 w-full z-[80] bg-[#fad400] text-black flex flex-row items-center justify-between min-h-16">
      <div className="flex justify-between w-full h-full px-5">
        {/* Left group */}
        <div className="pointer-events-auto flex flex-row items-center gap-4 md:gap-6 isolation-isolate">
          <Link href="/site/public" className="font-semibold text-lg leading-none tracking-tight">
            UploadParty
          </Link>
          <span className="hidden sm:inline-block rounded-full border border-black/20 px-2 py-1 text-xs">
            EXPERIMENT
          </span>
        </div>

        {/* Center group */}
        <div className="hidden md:flex items-center">
          <Link
            href="/app"
            className="rounded-full bg-black text-yellow-300 hover:bg-black/90 transition-colors px-5 py-2 font-medium"
            aria-label="Enter tool"
          >
            Enter Tool
          </Link>
        </div>

        {/* Right group */}
        <div className="pointer-events-auto flex flex-row items-center gap-4 md:gap-6 isolation-isolate">
          <Link href="/library" className="hidden sm:inline font-medium">
            My Library
          </Link>
          <button aria-label="Help" className="h-9 w-9 rounded-full border border-black/20 grid place-items-center">
            ?
          </button>
          <Link href="/pro" className="hidden sm:inline font-medium">
            Pro
          </Link>
          <div className="h-9 w-9 rounded-full bg-black/20" aria-hidden />
        </div>
      </div>
    </header>
  );
}
