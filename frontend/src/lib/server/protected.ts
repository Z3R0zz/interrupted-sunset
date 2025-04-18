"use server";
import api from "@/modules/API";
import type { APIRes, UserRes } from "@/types/api";
import axios from "axios";

/*
Copyright © 2025 Kars (github.com/kars1996)

Not to be shared, replicated or used without prior consent.
Contact Kars for any enquiries
*/

export async function getUser(): Promise<APIRes<UserRes>> {
    try {
        const res = await api.get<UserRes>("/user");
        if (res.data instanceof Object) {
            return {
                success: true,
                data: res.data,
            };
        }
        console.log(JSON.stringify(res))
        return {
            success: false,
            error: "Failed to get user",
        };
    } catch (error) {
        if (axios.isAxiosError(error)) {
            return {
                success: false,
                error: error.response?.data?.error || error.message,
            };
        }
        console.log(error)
        return {
            success: false,
            error: "Failed to get user",
        };
    }
}

export async function createMail(
    email: string,
): Promise<APIRes<{ message: string }>> {
    try {
        const res = await api.post<{ message: string }>("/mail/new", {
            email,
        });
        if (res.status == 200) {
            return {
                success: true,
                data: res.data,
            };
        }

        console.log(JSON.stringify(res))
        return {
            success: false,
            error: "Failed to create mail",
        };
    } catch (error) {
        if (axios.isAxiosError(error)) {
            return {
                success: false,
                error: error.response?.data?.error || error.message,
            };
        }
        console.log(error)
        return {
            success: false,
            error: "Failed to create mail",
        };
    }
}

export async function verifyMail(
    code: string,
): Promise<APIRes<{ message: string }>> {
    try {
        const res = await api.post<{ message: string }>("/mail/verify", {
            code: Number(code),
        });
        if (res.status == 200) {
            return {
                success: true,
                data: res.data,
            };
        }
        console.log(JSON.stringify(res))
        return {
            success: false,
            error: "Failed to verify mail",
        };
    } catch (error) {
        if (axios.isAxiosError(error)) {
            return {
                success: false,
                error: error.response?.data?.error || error.message,
            };
        }
        console.log(error)
        return {
            success: false,
            error: "Failed to verify mail",
        };
    }
}

export async function downloadQueue(): Promise<APIRes<{ message: string }>> {
    try {
        const res = await api.post<{ message: string }>("/queue");
        if (res.status == 200) {
            return {
                success: true,
                data: res.data,
            };
        }
        console.log(JSON.stringify(res))
        return {
            success: false,
            error: "Failed to download queue",
        };
    } catch (error) {
        if (axios.isAxiosError(error)) {
            return {
                success: false,
                error: error.response?.data?.error || error.message,
            };
        }
        console.log(error)
        return {
            success: false,
            error: "Failed to download queue",
        };
    }
}
