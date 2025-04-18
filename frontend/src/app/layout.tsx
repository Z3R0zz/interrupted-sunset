import type { Metadata } from "next";
import { Inter, Poppins } from "next/font/google";
import "./globals.css";
import Body from "@/modules/Body/Body";
import Console from "@/modules/Console/Console";
import AOS from "@/lib/Aos/aos";
import * as Fonts from "../../public/fonts/fontExports";

/*
Copyright © 2025 Kars (github.com/kars1996)

Not to be shared, replicated or used without prior consent.
Contact Kars for any enquiries
*/

const inter = Inter({ subsets: ["latin"] });
// ? Optional Font (You can add more if you want)
const poppins = Poppins({
    weight: ["100", "200", "300", "400", "500", "600", "700", "800", "900"],
    subsets: ["latin"],
    display: "swap",
});

export const metadata: Metadata = {
    title: { template: "%s – interrupted.me", default: "Sunset" },

    openGraph: {
        url: "https://interrupted.me/",
        type: "website",
        title: "interrupted.me - Sunset",
        siteName: "interrupted.me",
        description:
            "After years of dedication and collaboration, we're proud to announce the completion of our journey. We've created something meaningful together, and now it's time to celebrate our achievements.",
    },
    robots: {
        index: true,
        follow: true,
        nocache: true,
        googleBot: {
            index: true,
            follow: true,
            noimageindex: true,
            "max-video-preview": -1,
            "max-image-preview": "large",
            "max-snippet": -1,
        },
    },
};

export const viewport = {
    themeColor: "#ff1000"
}

// ? This check assumes you're hosting on vercel. If you're self-hosting you will need another check
let isProd: boolean = false;
if (process.env.NEXT_PUBLIC_VERCEL_GIT_COMMIT_SHA) {
    isProd = true;
}

export default function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="en" made-by="kars">
            <head>
                <link rel="icon" href="/assets/fav.png" />
            </head>
            <Body className={`${inter.className} __kars`}>
                <main className="relative flex min-h-screen w-full flex-col">
                    <AOS />
                    {children}
                </main>
                <Console isProd={isProd} />
            </Body>
        </html>
    );
}
