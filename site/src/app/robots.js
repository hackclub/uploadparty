// Next.js App Router robots route
// Docs: https://nextjs.org/docs/app/api-reference/file-conventions/metadata/robots

export default function robots() {
  const host = process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000";
  return {
    rules: {
      userAgent: "*",
      allow: "/",
    },
    sitemap: `${host}/sitemap.xml`,
    host,
  };
}
