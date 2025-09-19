// Next.js App Router sitemap generator
// Docs: https://nextjs.org/docs/app/api-reference/functions/generate-sitemaps

const BASE_URL = process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000";

export default function sitemap() {
  const now = new Date();
  /** @type {import('next').MetadataRoute.Sitemap} */
  const routes = [
    "",
    "/overview",
    "/prizes",
    "/faq",
  ].map((path) => ({
    url: `${BASE_URL}${path || "/"}`,
    lastModified: now,
    changeFrequency: "weekly",
    priority: path === "" ? 1.0 : 0.8,
  }));

  return routes;
}
