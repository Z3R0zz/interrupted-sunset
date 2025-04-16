/*
Copyright © 2025 Kars (github.com/kars1996)

Not to be shared, replicated or used without prior consent.
Contact Kars for any enquiries
*/

export async function GET(): Promise<Response> {
    return new Response("Example", {
        status: 200,
        headers: {
            "Content-Type": "application/json",
        },
    });
}
