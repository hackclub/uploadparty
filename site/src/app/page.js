"use client";
import { useEffect } from "react";
import Link from "next/link";
import Image from "next/image";

export default function Home() {
  // Hide the global header and remove top padding while on the home page
  useEffect(() => {
    const body = typeof document !== "undefined" ? document.body : null;
    if (!body) return;

    // Remove global header padding and hide the header element
    body.classList.remove("pt-16");
    body.classList.add("pt-0");

    const hdr = document.querySelector("header");
    const prevHeaderDisplay = hdr && hdr instanceof HTMLElement ? hdr.style.display : "";
    if (hdr && hdr instanceof HTMLElement) hdr.style.display = "none";

    return () => {
      // Restore padding and header on unmount
      body.classList.remove("pt-0");
      body.classList.add("pt-16");
      if (hdr) (hdr as HTMLElement).style.display = prevHeaderDisplay || "";
    };
  }, []);

  return (
    <div className="min-h-screen bg-[#fad400] text-black">
      {/* Hero */}
      <section className="mx-auto max-w-6xl px-5 pt-16 sm:pt-24 pb-10 sm:pb-16 text-center">
        <div className="mx-auto mb-6 h-12 w-12 rounded-full grid place-items-center border border-black/20">
          <span className="text-2xl">✨</span>
        </div>
        <h1 className="font-semibold tracking-[.25em] text-3xl sm:text-4xl md:text-5xl leading-tight">
          CREATE SOME MAGIC
        </h1>
        <div className="mt-6 flex justify-center">
          <Link
            href="/app"
            className="inline-flex items-center gap-2 rounded-full bg-black text-yellow-300 hover:bg-black/90 transition px-6 py-3 text-sm sm:text-base font-medium"
            aria-label="Enter tool"
          >
            ENTER TOOL
            <span className="translate-y-[1px]">→</span>
          </Link>
        </div>
        <div className="relative mt-6 hidden sm:block">
          <span className="absolute right-16 -rotate-12 text-xs font-semibold">GET<br/>INSPIRED!</span>
        </div>
      </section>

      {/* Gallery */}
      <section className="mx-auto max-w-[1200px] px-5 pb-24">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          <TiltCard rotation="-6">
            <PlaceholderImage color="#fff0f0" label="Cake" />
          </TiltCard>
          <TiltCard rotation="4">
            <Image
              src="/images/whisk.jpg"
              alt="Whisk placeholder"
              width={800}
              height={500}
              className="w-full h-auto rounded-lg"
            />
          </TiltCard>
          <TiltCard rotation="-3">
            <PlaceholderImage color="#ffe8b3" label="Landscape" />
          </TiltCard>
          <TiltCard rotation="2">
            <PlaceholderImage color="#e7f0ff" label="Beach" />
          </TiltCard>
          <TiltCard rotation="-5">
            <PlaceholderImage color="#f2f2f2" label="Pearls" />
          </TiltCard>
          <TiltCard rotation="3">
            <PlaceholderImage color="#ffe0e0" label="Portraits" />
          </TiltCard>
        </div>
      </section>
    </div>
  );
}

function TiltCard({ children, rotation = "0" }) {
  return (
    <div
      className="bg-black/90 text-yellow-300 rounded-xl p-2 shadow-[0_10px_50px_rgba(0,0,0,0.2)]"
      style={{ transform: `rotate(${rotation}deg)` }}
    >
      <div className="bg-black rounded-lg overflow-hidden">
        {children}
      </div>
      <div className="mt-2 text-center text-xs tracking-wide">TRY ANIMATE!</div>
    </div>
  );
}

function PlaceholderImage({ color = "#eee", label = "" }) {
  return (
    <div
      className="aspect-[16/9] w-full"
      style={{ backgroundColor: color }}
      aria-label={label}
    />
  );
}
