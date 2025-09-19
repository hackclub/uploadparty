// Dedicated page for the Prizes section (public pre-auth)
// Route: /prizes
// Server Component by default

import Prizes from "../sections/Prizes";

export const metadata = {
  title: "UploadParty — Prizes",
  description: "Prizes and rewards at UploadParty",
  alternates: {
    canonical: (process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000") + "/prizes",
  },
  openGraph: {
    title: "UploadParty — Prizes",
    description: "Prizes and rewards at UploadParty",
    url: (process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000") + "/prizes",
    siteName: "UploadParty",
    type: "article",
  },
};

function JsonLd() {
  const base = process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000";
  const json = {
    "@context": "https://schema.org",
    "@type": "WebPage",
    name: "UploadParty — Prizes",
    url: base + "/prizes",
  };
  const breadcrumbs = {
    "@context": "https://schema.org",
    "@type": "BreadcrumbList",
    itemListElement: [
      { "@type": "ListItem", position: 1, name: "Home", item: base + "/" },
      { "@type": "ListItem", position: 2, name: "Prizes", item: base + "/prizes" }
    ]
  };
  return (
    <>
      <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(json) }} />
      <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(breadcrumbs) }} />
    </>
  );
}

export default function PrizesPage() {
  return (
    <main>
      <JsonLd />
      <Prizes />
    </main>
  );
}
