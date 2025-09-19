// Dedicated page for the FAQ & Footer section (public pre-auth)
// Route: /faq
// Server Component by default

import FAQFooter from "../sections/FAQFooter";

export const metadata = {
  title: "UploadParty — FAQ",
  description: "Frequently asked questions about UploadParty",
  alternates: {
    canonical: (process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000") + "/faq",
  },
  openGraph: {
    title: "UploadParty — FAQ",
    description: "Frequently asked questions about UploadParty",
    url: (process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000") + "/faq",
    siteName: "UploadParty",
    type: "article",
  },
};

function JsonLd() {
  const base = process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000";
  const json = {
    "@context": "https://schema.org",
    "@type": "WebPage",
    name: "UploadParty — FAQ",
    url: base + "/faq",
  };
  const breadcrumbs = {
    "@context": "https://schema.org",
    "@type": "BreadcrumbList",
    itemListElement: [
      { "@type": "ListItem", position: 1, name: "Home", item: base + "/" },
      { "@type": "ListItem", position: 2, name: "FAQ", item: base + "/faq" }
    ]
  };
  return (
    <>
      <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(json) }} />
      <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(breadcrumbs) }} />
    </>
  );
}

export default function FAQPage() {
  return (
    <main>
      <JsonLd />
      <FAQFooter />
    </main>
  );
}
