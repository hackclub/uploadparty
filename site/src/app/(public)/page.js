// Public (pre‑auth) landing page
// Route: /
// This file lives in the (public) group to make the separation from the authenticated (app) area explicit.
// Keep this page server-side by default (no "use client").

import Overview from "./sections/Overview";
import Prizes from "./sections/Prizes";
import FAQFooter from "./sections/FAQFooter";

export const metadata = {
    title: "UploadParty × Native Instruments — Teen Music Production Challenge",
    description: "Join 500 teen creators in UploadParty - a music production challenge backed by Native Instruments. Learn beats, earn Komplete 15 software, and connect with young producers worldwide.",
    keywords: "music production, teen creators, Native Instruments, Komplete 15, beats, music challenge, young producers, VST, music software",
    alternates: {
        canonical: (process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000") + "/",
    },
    openGraph: {
        title: "UploadParty × Native Instruments — Teen Music Production Challenge",
        description: "Join 500 teen creators in UploadParty - a music production challenge backed by Native Instruments. Learn beats, earn Komplete 15 software, and connect with young producers worldwide.",
        url: (process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000") + "/",
        siteName: "UploadParty",
        type: "website",
        images: [
            {
                url: (process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000") + "/og-image.jpg",
                width: 1200,
                height: 630,
                alt: "UploadParty - Teen Music Production Challenge",
            },
        ],
    },
    twitter: {
        card: "summary_large_image",
        title: "UploadParty × Native Instruments — Teen Music Production Challenge",
        description: "Join 500 teen creators learning music production with Native Instruments software",
        images: [(process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000") + "/og-image.jpg"],
    },
};

function JsonLd() {
    const base = process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000";

    const organization = {
        "@context": "https://schema.org",
        "@type": "Organization",
        name: "UploadParty",
        description: "A music production challenge for teen creators backed by Native Instruments",
        url: base,
        logo: base + "/logo.png",
        sameAs: [
            "https://hackclub.com",
            // Add social media URLs when available
        ]
    };

    const event = {
        "@context": "https://schema.org",
        "@type": "Event",
        name: "UploadParty Music Production Challenge",
        description: "A gamified music production challenge for 500 teen creators, sponsored by Native Instruments with $300K in Komplete 15 software licenses",
        startDate: "2025-09-26",
        endDate: "2025-10-26", // Adjust based on your actual end date
        eventStatus: "https://schema.org/EventScheduled",
        eventAttendanceMode: "https://schema.org/OnlineEventAttendanceMode",
        location: {
            "@type": "VirtualLocation",
            url: base
        },
        organizer: {
            "@type": "Organization",
            name: "Hack Club",
            url: "https://hackclub.com"
        },
        sponsor: {
            "@type": "Organization",
            name: "Native Instruments",
            url: "https://native-instruments.com"
        },
        offers: {
            "@type": "Offer",
            price: "0",
            priceCurrency: "USD",
            availability: "https://schema.org/InStock",
            url: base + "/rsvp"
        }
    };

    const breadcrumbs = {
        "@context": "https://schema.org",
        "@type": "BreadcrumbList",
        itemListElement: [
            {
                "@type": "ListItem",
                position: 1,
                name: "Home",
                item: base + "/"
            },
            {
                "@type": "ListItem",
                position: 2,
                name: "How It Works",
                item: base + "#overview"
            },
            {
                "@type": "ListItem",
                position: 3,
                name: "Prizes & Software",
                item: base + "#prizes"
            },
            {
                "@type": "ListItem",
                position: 4,
                name: "FAQ",
                item: base + "#faq"
            }
        ]
    };

    const faq = {
        "@context": "https://schema.org",
        "@type": "FAQPage",
        mainEntity: [
            {
                "@type": "Question",
                name: "Who can participate in UploadParty?",
                acceptedAnswer: {
                    "@type": "Answer",
                    text: "UploadParty is open to teen creators aged 13-18 who are interested in learning music production and working with professional software like Native Instruments Komplete 15."
                }
            },
            {
                "@type": "Question",
                name: "What software do I get?",
                acceptedAnswer: {
                    "@type": "Answer",
                    text: "Successful participants receive Native Instruments Komplete 15 software licenses, valued at $600 each, as part of our $300,000 partnership with Native Instruments."
                }
            },
            {
                "@type": "Question",
                name: "How long does the challenge last?",
                acceptedAnswer: {
                    "@type": "Answer",
                    text: "The UploadParty challenge runs for approximately one month, with rolling challenges and community engagement throughout the program."
                }
            }
        ]
    };

    return (
        <>
            <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(organization) }} />
            <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(event) }} />
            <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(breadcrumbs) }} />
            <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(faq) }} />
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