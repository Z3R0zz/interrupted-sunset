"use client";
import { useState } from "react";
import { LogIn, ArrowLeft, Mail, Lock } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { login } from "@/lib/server/auth";
import { showToast } from "../ui/toast";

export default function LoginPage() {
    const router = useRouter();
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [isLoading, setIsLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);

        try {
            const response = await login(username, password);
            if (!response.success)
                throw new Error(response.error || "Login failed");
            if (response.success && response.data && response.data.token) {
                const expiryDate = new Date();
                expiryDate.setDate(expiryDate.getDate() + 7);
                document.cookie = `sunset_token=${response.data!.token}; expires=${expiryDate.toUTCString()}; path=/; SameSite=Strict${process.env.NODE_ENV === "production" ? "; Secure" : ""}`;

                showToast("success", "Login successful");

                setTimeout(() => {
                    router.push("/dashboard");
                }, 1000);
            } else {
                throw new Error(response.error || "Login failed");
            }
        } catch (error) {
            showToast(
                "error",
                error instanceof Error
                    ? error.message
                    : "Login failed. Please check your credentials",
            );
            setIsLoading(false);
        }
    };

    return (
        <div className="flex min-h-screen flex-col bg-[#0c0c0c] text-white">
            <div className="absolute inset-0 overflow-hidden">
                <div className="absolute inset-0 bg-gradient-to-b from-[#0c0c0c] via-[#1a0f23] to-[#3d1635]" />
                <div
                    className="absolute h-24 w-24 rounded-full bg-gradient-to-b from-orange-300 to-red-500"
                    style={{
                        left: "calc(50% - 3rem)",
                        top: "50%",
                        boxShadow: "0 0 60px rgba(255, 50, 50, 0.5)",
                    }}
                />
                <div className="absolute bottom-0 h-[1px] w-full bg-gradient-to-r from-transparent via-red-800/30 to-transparent" />
            </div>

            <div className="relative flex min-h-screen items-center justify-center px-4 py-12">
                <div className="w-full max-w-md">
                    <div className="mb-8 text-center">
                        <h1 className="mb-2 text-4xl font-bold tracking-tighter">
                            <span className="text-red-500">Welcome Back</span>
                        </h1>
                        <p className="text-zinc-400">
                            Log in to access your data
                        </p>
                    </div>

                    <div className="rounded-xl border border-zinc-800 bg-zinc-900/70 p-6 backdrop-blur-md">
                        <form onSubmit={handleSubmit}>
                            <div className="mb-4">
                                <label className="mb-2 block text-sm font-medium text-zinc-300">
                                    Username
                                </label>
                                <div className="relative">
                                    <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                                        <Mail className="h-5 w-5 text-zinc-500" />
                                    </div>
                                    <input
                                        type="input"
                                        value={username}
                                        onChange={(e) =>
                                            setUsername(e.target.value)
                                        }
                                        className="block w-full rounded-lg border border-zinc-700 bg-zinc-800/50 p-2.5 pl-10 text-white placeholder-zinc-400 focus:border-red-500 focus:outline-none focus:ring-1 focus:ring-red-500"
                                        placeholder="you"
                                        required
                                    />
                                </div>
                            </div>

                            <div className="mb-6">
                                <label className="mb-2 block text-sm font-medium text-zinc-300">
                                    Password
                                </label>
                                <div className="relative">
                                    <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                                        <Lock className="h-5 w-5 text-zinc-500" />
                                    </div>
                                    <input
                                        type="password"
                                        value={password}
                                        onChange={(e) =>
                                            setPassword(e.target.value)
                                        }
                                        className="block w-full rounded-lg border border-zinc-700 bg-zinc-800/50 p-2.5 pl-10 text-white placeholder-zinc-400 focus:border-red-500 focus:outline-none focus:ring-1 focus:ring-red-500"
                                        placeholder="••••••••"
                                        required
                                    />
                                </div>
                            </div>

                            <button
                                type="submit"
                                disabled={isLoading}
                                className="group flex w-full items-center justify-center rounded-full bg-gradient-to-r from-red-600 to-red-800 px-6 py-2.5 font-medium transition-all hover:from-red-700 hover:to-red-900 disabled:opacity-70"
                            >
                                {isLoading ? (
                                    <div className="h-5 w-5 animate-spin rounded-full border-2 border-white border-t-transparent"></div>
                                ) : (
                                    <>
                                        Log In{" "}
                                        <LogIn className="ml-2 h-4 w-4 transition-transform group-hover:translate-x-1" />
                                    </>
                                )}
                            </button>
                        </form>
                    </div>

                    <div className="mt-8 text-center">
                        <Link
                            href="/"
                            className="group inline-flex items-center text-sm text-zinc-400 transition-colors hover:text-red-500"
                        >
                            <ArrowLeft className="mr-2 h-4 w-4 transition-transform group-hover:-translate-x-1" />
                            Back to Home
                        </Link>
                    </div>

                    <footer className="mt-auto border-t border-zinc-900 py-6">
                        <div className="flex items-center justify-center">
                            <div className="flex items-center">
                                <div className="mr-3 h-4 w-4 rounded-full bg-gradient-to-b from-orange-300 to-red-500"></div>
                                <p className="text-sm text-zinc-500">
                                    interrupted.me © 2023-2025
                                </p>
                            </div>
                        </div>
                    </footer>
                </div>
            </div>
        </div>
    );
}
