import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import SkipAuthButton from "./components/SkipAuthButton";
import Header from "./components/shared/Header";
import ServerErrorBoundary from "./components/ServerErrorBoundary";
import ErrorBoundary from "./components/ErrorBoundary";
import { UserProvider } from '@auth0/nextjs-auth0/client';
import { headers } from "next/headers";
import { isMobileUserAgent } from "./lib/hooks/utils";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
  display: "swap", // Performance optimization
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
  display: "swap", // Performance optimization
});

export const metadata = {
  title: "UploadParty - Upload beats and audio for community challenges",
  description: "Join UploadParty to upload your beats and audio for exciting community challenges. Connect with producers and showcase your music.",
  keywords: "music, beats, audio, upload, community, challenges, producers, music production",
  authors: [{ name: "UploadParty Team" }],
  creator: "UploadParty",
  publisher: "UploadParty",
  robots: "index, follow",
  openGraph: {
    title: "UploadParty",
    description: "Upload beats and audio for community challenges",
    type: "website",
    locale: "en_US",
    url: process.env.NEXT_PUBLIC_SITE_URL,
    siteName: "UploadParty",
  },
  twitter: {
    card: "summary_large_image",
    title: "UploadParty",
    description: "Upload beats and audio for community challenges",
    creator: "@uploadparty",
  },
  alternates: {
    canonical: process.env.NEXT_PUBLIC_SITE_URL,
  },
  other: {
    "apple-mobile-web-app-capable": "yes",
    "apple-mobile-web-app-status-bar-style": "default",
    "apple-mobile-web-app-title": "UploadParty",
  },
};

export default async function RootLayout({ children }) {
  const headersList = await headers();
  const ua = headersList.get("user-agent") || "";
  const uaMobile = isMobileUserAgent(ua);

  return (
    <html lang="en" data-ua-mobile={uaMobile ? "true" : "false"}>
      <head>
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta name="theme-color" content="#000000" />
      </head>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased transition-[padding] duration-300 ease-in-out pt-16`}
        data-ua-mobile={uaMobile ? "true" : "false"}
      >
        <ErrorBoundary>
          <UserProvider>
            <ServerErrorBoundary>
              {children}
              <SkipAuthButton />
            </ServerErrorBoundary>
          </UserProvider>
        </ErrorBoundary>
      </body>
    </html>
  );
}
