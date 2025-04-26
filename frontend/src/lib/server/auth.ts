"use server";
import api from "@/modules/API";
import type { APIRes, LoginRes } from "@/types/api";
import axios from "axios";

/*
Copyright © 2025 Kars (github.com/kars1996)

Not to be shared, replicated, or used without prior consent.
Contact me for any enquiries
*/

export async function login(
    username: string,
    password: string,
): Promise<APIRes<LoginRes>> {
    try {
        const res = await api.post<LoginRes>("/login", {
            username,
            password,
        });
        if (!res.data.token) {
            console.log(JSON.stringify(res))
            throw new Error("No token in response");
        }
        return {
            success: true,
            data: res.data,
        };
    } catch (error) {
        if (axios.isAxiosError(error)) {
            return {
                success: false,
                error: error.response?.data?.error || error.message,
            };
        }
        console.log(error);
        return {
            success: false,
            error: (error as Error).message || "Unknown error",
        };
    }
}
