"use client";
import { useEffect } from "react";
import Head from "next/head";

/*
Copyright © 2025 Kars (github.com/kars1996)

Not to be shared, replicated or used without prior consent.
Contact Kars for any enquiries
*/

interface PageMeta {
    title: string;
    description?: string;
}

/**
 * @deprecated this component is deprecated. use default NextJS metadata instead.
 */
export default function Meta(meta: PageMeta) {
    useEffect(() => {
        document.title = meta.title;
    }, [meta.title]);

    return (
        <Head>
            <title>{meta.title}</title>
            {meta.description && (
                <meta name="description" content={meta.description} />
            )}
        </Head>
    );
}
