import { ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

/*
Copyright © 2025 Kars (github.com/kars1996)

Not to be shared, replicated, or used without prior consent.
Contact me for any enquiries
*/

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs));
}
