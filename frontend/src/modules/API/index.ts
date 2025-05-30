import axios, { AxiosError, AxiosResponse, AxiosRequestConfig } from "axios";
import { cookies } from "next/headers";

interface ApiResponse<T = unknown> {
    status: number;
    data: T;
    error?: string;
}

/*
Copyright © 2025 Kars (github.com/kars1996)

Not to be shared, replicated or used without prior consent.
Contact Kars for any enquiries
*/

export default class api {
    private static baseURL = process.env.BACKEND_URL
    private static userAgent = process.env.API_USER_AGENT || "your-useragent";
    private static instance = axios.create({
        baseURL: this.baseURL,
        withCredentials: true,
        headers: {
            "Content-Type": "application/json",
            "User-Agent": this.userAgent,
        },
    });

    private static async getAuthToken(): Promise<string | undefined> {
        const cookieStore = await cookies();
        return cookieStore.get("sunset_token")?.value;
    }

    private static async getHeaders(): Promise<any> {
        const token = await this.getAuthToken();
        const headers: any = {
            "Content-Type": "application/json",
        };
        if (token) {
            headers["Authorization"] = `Bearer ${token}`;
        }
        return headers;
    }

    private static async handleRequest<T>(
        promise: Promise<AxiosResponse<T>>,
    ): Promise<ApiResponse<T>> {
        try {
            const response = await promise;
            return {
                status: response.status,
                data: response.data,
            };
        } catch (error) {
            const axiosError = error as AxiosError<T>;
            return {
                status: axiosError.response?.status || 500,
                data: axiosError.response?.data as T,
                error: axiosError.message,
            };
        }
    }

    static async get<T = unknown, P = unknown>(
        url: string,
        config?: Omit<AxiosRequestConfig<P>, "url" | "headers">,
    ): Promise<ApiResponse<T>> {
        const headers = await this.getHeaders();
        return this.handleRequest(
            this.instance.get<T>(url, { ...config, headers }),
        );
    }

    static async post<T = unknown, D = unknown>(
        url: string,
        data?: D,
        config?: Omit<AxiosRequestConfig<D>, "url" | "data" | "headers">,
    ): Promise<ApiResponse<T>> {
        const headers = await this.getHeaders();
        return this.handleRequest(
            this.instance.post<T>(url, data, { ...config, headers }),
        );
    }

    static async put<T = unknown, D = unknown>(
        url: string,
        data?: D,
        config?: Omit<AxiosRequestConfig<D>, "url" | "data" | "headers">,
    ): Promise<ApiResponse<T>> {
        const headers = await this.getHeaders();
        return this.handleRequest(
            this.instance.put<T>(url, data, { ...config, headers }),
        );
    }

    static async delete<T = unknown>(
        url: string,
        config?: Omit<AxiosRequestConfig, "url" | "headers">,
    ): Promise<ApiResponse<T>> {
        const headers = await this.getHeaders();
        return this.handleRequest(
            this.instance.delete<T>(url, { ...config, headers }),
        );
    }

    static async patch<T = unknown, D = unknown>(
        url: string,
        data?: D,
        config?: Omit<AxiosRequestConfig<D>, "url" | "data" | "headers">,
    ): Promise<ApiResponse<T>> {
        const headers = await this.getHeaders();
        return this.handleRequest(
            this.instance.patch<T>(url, data, { ...config, headers }),
        );
    }
}
