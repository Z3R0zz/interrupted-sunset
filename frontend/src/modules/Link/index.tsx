"use client";
import { SPContext } from "@/modules/Body/Body";
import Link from "next/link";
import { useContext } from "react";

// Thanks to nitlix for this code

type LinkProp = {
    href: string;
} & Omit<React.ComponentProps<"a">, "href">;

export default function RouterLink({
    children,
    href,
    className = "",
    ...props
}: LinkProp) {
    const lenisContext = useContext(SPContext);

    const { scroll } = lenisContext;

    if (href.startsWith("#")) {
        return (
            <a
                href={href}
                onClick={(e) => {
                    if (scroll) {
                        e.preventDefault();
                        scroll.scrollTo(href);
                    }
                }}
                className={className}
                {...props}
            >
                {children}
            </a>
        );
    }
    return (
        <Link href={href} className={className} {...props}>
            {children}
        </Link>
    );
}
