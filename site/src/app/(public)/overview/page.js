// Dedicated page for the Overview section (public pre-auth)
// Route: /overview
// Server Component by default

import Overview from "../sections/Overview";

export const metadata = {
  title: "UploadParty — Overview",
  description: "Overview of UploadParty",
  alternates: {
    canonical: (process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000") + "/overview",
  },
  openGraph: {
    title: "UploadParty — Overview",
    description: "Overview of UploadParty",
    url: (process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000") + "/overview",
    siteName: "UploadParty",
    type: "article",
  },
};

function JsonLd() {
  const base = process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000";
  const json = {
    "@context": "https://schema.org",
    "@type": "WebPage",
    name: "UploadParty — Overview",
    url: base + "/overview",
  };
  const breadcrumbs = {
    "@context": "https://schema.org",
    "@type": "BreadcrumbList",
    itemListElement: [
      { "@type": "ListItem", position: 1, name: "Home", item: base + "/" },
      { "@type": "ListItem", position: 2, name: "Overview", item: base + "/overview" }
    ]
  };
  return (
    <>
      <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(json) }} />
      <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(breadcrumbs) }} />
    </>
  );
}

export default function OverviewPage() {
  return (
    <main>
      <JsonLd />
      <Overview />
    </main>
  );
}
