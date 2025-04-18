import type { NextRequest } from "next/server";

/*
Copyright © 2025 Kars (github.com/kars1996)

Not to be shared, replicated or used without prior consent.
Contact Kars for any enquiries
*/

export function middleware(request: NextRequest) {
    const currentUser = request.cookies.get("sunset_token")?.value;

    if (currentUser && request.nextUrl.pathname.startsWith("/login")) {
        return Response.redirect(new URL("/dashboard", request.url));
    }

    if (!currentUser && request.nextUrl.pathname.startsWith("/dashboard")) {
        return Response.redirect(
            new URL(
                `/login?redirect=${encodeURI(request.nextUrl.pathname)}`,
                request.url,
            ),
        );
    }
}

export const config = {
    matcher: ["/((?!api|_next/static|_next/image|.*\\.png$|favicon.ico).*)"],
};
