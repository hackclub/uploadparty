// Public (pre‑auth) landing page
// Route: /
// This file lives in the (public) group to make the separation from the authenticated (app) area explicit.
// Keep this page server-side by default (no "use client").

import Overview from "./sections/Overview";
import Prizes from "./sections/Prizes";
import FAQFooter from "./sections/FAQFooter";

export const metadata = {
  title: "UploadParty — Welcome",
  description: "Public landing page for UploadParty",
  alternates: {
    canonical: (process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000") + "/",
  },
  openGraph: {
    title: "UploadParty — Welcome",
    description: "Public landing page for UploadParty",
    url: (process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000") + "/",
    siteName: "UploadParty",
    type: "website",
  },
};

function JsonLd() {
  const base = process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000";
  const json = {
    "@context": "https://schema.org",
    "@type": "WebSite",
    name: "UploadParty",
    url: base + "/",
    potentialAction: {
      "@type": "SearchAction",
      target: `${base}/search?q={search_term_string}`,
      "query-input": "required name=search_term_string"
    }
  };
  const breadcrumbs = {
    "@context": "https://schema.org",
    "@type": "BreadcrumbList",
    itemListElement: [
      { "@type": "ListItem", position: 1, name: "Home", item: base + "/" },
      { "@type": "ListItem", position: 2, name: "Overview", item: base + "/overview" },
      { "@type": "ListItem", position: 3, name: "Prizes", item: base + "/prizes" },
      { "@type": "ListItem", position: 4, name: "FAQ", item: base + "/faq" }
    ]
  };
  return (
    <>
      <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(json) }} />
      <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(breadcrumbs) }} />
    </>
  );
}

export default function PublicLanding() {
  return (
    <main>
      <JsonLd />
      <Overview />
      <Prizes />
      <FAQFooter />
    </main>
  );
}
