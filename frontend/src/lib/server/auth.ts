"use server";
import api from "@/modules/API";
import type { APIRes, LoginRes } from "@/types/api";
import axios from "axios";

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
            throw new Error("Login failed");
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
        return {
            success: false,
            error: "Failed to get token",
        };
    }
}
