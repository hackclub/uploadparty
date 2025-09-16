"use client";
import { useState, useEffect } from "react";
import Link from "next/link";

export default function Navbar() {
  const [open, setOpen] = useState(false);

  useEffect(() => {
    const body = typeof document !== "undefined" ? document.body : null;
    if (!body) return;
    // Ensure body padding transitions smoothly and responds to sidebar state on md+ screens
    body.classList.add("transition-[padding]", "duration-300", "ease-in-out");
    if (open) {
      body.classList.remove("md:pl-14");
      body.classList.add("md:pl-64");
    } else {
      body.classList.remove("md:pl-64");
      body.classList.add("md:pl-14");
    }
  }, [open]);

  return (
    <aside
      className={`fixed left-0 top-0 h-screen border-r border-black/10 dark:border-white/10 bg-white/70 dark:bg-black/40 backdrop-blur supports-[backdrop-filter]:bg-white/60 z-50 transition-[width] duration-300 ease-in-out ${
        open ? "w-64" : "w-14"
      }`}
    >
      <div className="h-full relative flex flex-col">
        <button
          aria-label={open ? "Collapse sidebar" : "Expand sidebar"}
          onClick={() => setOpen((v) => !v)}
          className={`absolute -right-3 top-4 h-7 w-7 rounded-full border border-black/10 dark:border-white/10 bg-white/90 dark:bg-black/70 shadow hover:shadow-md transition flex items-center justify-center`}
        >
          <span
            className={`inline-block transition-transform duration-200 ${
              open ? "rotate-180" : "rotate-0"
            }`}
          >
            ▸
          </span>
        </button>

        <div className="px-3 py-4 border-b border-black/5 dark:border-white/10">
          <Link href="/" className="flex items-center gap-3">
            <div className="h-6 w-6 rounded bg-black/10 dark:bg-white/10" aria-hidden />
            <span className={`${open ? "block" : "hidden"} text-base font-semibold tracking-tight`}>
              UploadParty
            </span>
          </Link>
        </div>

        <nav className="flex-1 overflow-y-auto p-2 flex flex-col gap-1 text-sm">
          <NavItem href="/dashboard" label="Dashboard" open={open} />
          <NavItem href="/settings" label="Settings" open={open} />
        </nav>
      </div>
    </aside>
  );
}

function NavItem({ href, label, open }) {
  return (
    <Link
      href={href}
      className="flex items-center gap-3 px-2 py-2 rounded hover:bg-black/5 dark:hover:bg-white/10"
      aria-label={label}
    >
      <div className="h-5 w-5 rounded bg-black/10 dark:bg-white/10" aria-hidden />
      <span className={open ? "block" : "hidden"}>{label}</span>
    </Link>
  );
}
