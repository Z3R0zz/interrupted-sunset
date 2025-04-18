"use client";

import { useState, useEffect } from "react";
import {
    ArrowLeft,
    Download,
    FileDown,
    LogOut,
    Users,
    Mail,
    Check,
    AlertCircle,
} from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import type { StatusTypes, UserRes } from "@/types/api";
import {
    getUser,
    createMail,
    verifyMail,
    downloadQueue,
} from "@/lib/server/protected";
import OTPInput from "../ui/otp";

export default function DashboardPage() {
    const router = useRouter();
    const [user, setUser] = useState<UserRes | null>(null);
    const [loading, setLoading] = useState(true);
    const [activeTab, setActiveTab] = useState("data");
    const [email, setEmail] = useState("");
    const [otpCode, setOtpCode] = useState("");
    const [otpSent, setOtpSent] = useState(false);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");

    useEffect(() => {
        fetchUserData();
    }, []);

    const fetchUserData = async () => {
        setLoading(true);
        try {
            const response = await getUser();
            if (response.success && response.data) {
                setUser(response.data);
            } else {
                setError(response.error || "Failed to load user data");
            }
        } catch (error) {
            setError("Failed to fetch user data");
        } finally {
            setLoading(false);
        }
    };

    const handleLogout = () => {
        router.push("/");
    };

    const handleSendOTP = async () => {
        if (!email && !otpSent) {
            setError("Please enter your email address");
            return;
        }
    
        setError("");
        try {
            const response = await createMail(email);
            if (response.success) {
                setOtpSent(true);
            } else {
                setError(response.error || "Failed to send OTP");
            }
        } catch (error) {
            setError("An error occurred while sending OTP");
        }
    };

    const handleVerifyOTP = async () => {
        if (!otpCode) {
            setError("Please enter the OTP code");
            return;
        }
    
        setError("");
        setSuccess("");
        try {
            const response = await verifyMail(otpCode);
            if (response.success) {
                setSuccess("Email verified successfully");
                fetchUserData(); // Refresh user data
            } else {
                setError(response.error || "Failed to verify OTP");
            }
        } catch (error) {
            setError("An error occurred while verifying OTP");
        }
    };

    const handleDownloadQueue = async () => {
        setError("");
        try {
            const response = await downloadQueue();
            if (response.success) {
                setSuccess("Added to download queue successfully");
                fetchUserData(); // Refresh user data
            } else {
                setError(response.error || "Failed to join download queue");
            }
        } catch (error) {
            setError("An error occurred while joining download queue");
        }
    };

    const getStatusDisplay = (status: StatusTypes) => {
        switch (status) {
            case "waiting":
                return (
                    <div className="flex items-center space-x-2 text-yellow-500">
                        <AlertCircle size={16} />
                        <span>Waiting in queue</span>
                    </div>
                );
            case "processing":
                return (
                    <div className="flex items-center space-x-2 text-blue-500">
                        <AlertCircle size={16} />
                        <span>Processing your data</span>
                    </div>
                );
            case "done":
                return (
                    <div className="flex items-center space-x-2 text-green-500">
                        <Check size={16} />
                        <span>Data ready for download</span>
                    </div>
                );
            case "failed":
                return (
                    <div className="flex items-center space-x-2 text-red-500">
                        <AlertCircle size={16} />
                        <span>Processing failed</span>
                    </div>
                );
            default:
                return null;
        }
    };

    return (
        <div className="flex min-h-screen flex-col bg-[#0c0c0c] text-white">
            <header className="border-b border-zinc-900 px-4 py-4">
                <div className="container mx-auto flex max-w-6xl items-center justify-between">
                    <div className="flex items-center">
                        <div className="mr-3 h-8 w-8 rounded-full bg-gradient-to-b from-orange-300 to-red-500"></div>
                        <span className="text-xl font-bold">
                            interrupted.me
                        </span>
                    </div>
                    <button
                        onClick={handleLogout}
                        className="flex items-center rounded-md px-4 py-2 text-zinc-400 hover:bg-zinc-800 hover:text-white"
                    >
                        <LogOut className="mr-2 h-4 w-4" /> Log Out
                    </button>
                </div>
            </header>

            <main className="container mx-auto max-w-6xl flex-1 px-4 py-8">
                <div className="mb-8 flex flex-col items-start justify-between gap-4 md:flex-row md:items-center">
                    <div>
                        <Link
                            href="/"
                            className="mb-2 flex items-center text-zinc-400 hover:text-red-400"
                        >
                            <ArrowLeft className="mr-2 h-4 w-4" /> Back to
                            interrupted.me
                        </Link>
                        <h1 className="text-3xl font-bold">Your Dashboard</h1>
                        <p className="text-zinc-400">
                            Access and download your project data
                        </p>
                    </div>
                </div>

                {error && (
                    <div className="mb-6 flex items-center rounded-md border border-red-800 bg-red-900/20 px-4 py-3 text-red-300">
                        <AlertCircle className="mr-2 h-4 w-4" />
                        {error}
                    </div>
                )}

                {success && (
                    <div className="mb-6 flex items-center rounded-md border border-green-800 bg-green-900/20 px-4 py-3 text-green-300">
                        <Check className="mr-2 h-4 w-4" />
                        {success}
                    </div>
                )}

                {loading ? (
                    <div className="py-12 text-center">Loading...</div>
                ) : (
                    <div className="mb-8">
                        <div className="mb-6 flex border-b border-zinc-800">
                            <button
                                onClick={() => setActiveTab("data")}
                                className={`flex items-center px-4 py-2 ${
                                    activeTab === "data"
                                        ? "border-b-2 border-red-500 text-red-400"
                                        : "text-zinc-400"
                                }`}
                            >
                                <FileDown className="mr-2 h-4 w-4" /> Your Data
                            </button>
                            <button
                                onClick={() => setActiveTab("team")}
                                className={`flex items-center px-4 py-2 ${
                                    activeTab === "team"
                                        ? "border-b-2 border-red-500 text-red-400"
                                        : "text-zinc-400"
                                }`}
                            >
                                <Users className="mr-2 h-4 w-4" /> Team & Thanks
                            </button>
                        </div>

                        {activeTab === "data" && (
                            <div className="space-y-6">
                                <div className="rounded-lg border border-zinc-800 bg-zinc-900/50 p-6">
                                    <h2 className="mb-4 text-xl font-bold">
                                        Download Your Data
                                    </h2>
                                    <p className="mb-6 text-zinc-400">
                                        All your project data is available for
                                        download until December 31, 2025
                                    </p>

                                    {user && (
                                        <div className="mb-4 rounded-lg bg-zinc-800/50 p-4">
                                            <div className="mb-2 flex items-center justify-between">
                                                <div className="font-medium">
                                                    Account Status
                                                </div>
                                                <div className="text-sm text-zinc-400">
                                                    ID: {user.ID}
                                                </div>
                                            </div>
                                            <div className="mb-4 grid grid-cols-1 gap-4 md:grid-cols-2">
                                                <div className="flex items-center space-x-2">
                                                    <div className="text-zinc-400">
                                                        Email:
                                                    </div>
                                                    <div>{user.Email}</div>
                                                </div>
                                                <div className="flex items-center space-x-2">
                                                    <div className="text-zinc-400">
                                                        Email Verified:
                                                    </div>
                                                    <div
                                                        className={
                                                            user.Verified
                                                                ? "text-green-500"
                                                                : "text-red-500"
                                                        }
                                                    >
                                                        {user.Verified
                                                            ? "Yes"
                                                            : "No"}
                                                    </div>
                                                </div>
                                            </div>
                                            {user.Queue && (
                                                <div className="mt-4 rounded border border-zinc-700 bg-zinc-700/20 p-3">
                                                    <div className="mb-2 font-medium">
                                                        Download Status
                                                    </div>
                                                    {getStatusDisplay(
                                                        user.Status,
                                                    )}
                                                </div>
                                            )}
                                        </div>
                                    )}

                                    {user && !user.Verified ? (
                                        <div className="rounded-lg border border-zinc-700 bg-zinc-800/50 p-6">
                                            <h3 className="mb-4 font-medium">
                                                Verify Your Email
                                            </h3>
                                            <p className="mb-4 text-sm text-zinc-400">
                                                You need to verify your email
                                                before downloading your data
                                            </p>

                                            {!otpSent ? (
                                                <div className="space-y-4">
                                                    <div>
                                                        <label
                                                            htmlFor="email"
                                                            className="mb-1 block text-sm font-medium text-zinc-400"
                                                        >
                                                            Email Address
                                                        </label>
                                                        <input
                                                            type="email"
                                                            id="email"
                                                            value={email}
                                                            onChange={(e) =>
                                                                setEmail(
                                                                    e.target
                                                                        .value,
                                                                )
                                                            }
                                                            className="w-full rounded-md border border-zinc-700 bg-zinc-900 p-2 text-white"
                                                            placeholder="Enter your email"
                                                        />
                                                    </div>
                                                    <button
                                                        onClick={handleSendOTP}
                                                        className="flex w-full items-center justify-center rounded-md bg-gradient-to-r from-red-600 to-red-800 py-2 hover:from-red-700 hover:to-red-900"
                                                    >
                                                        <Mail className="mr-2 h-4 w-4" />
                                                        Send OTP
                                                    </button>
                                                </div>
                                            ) : (
                                                <div className="space-y-4">
                                                    <OTPInput
                                                        length={6}
                                                        onComplete={setOtpCode}
                                                        onSend={handleSendOTP}
                                                        cooldownTime={60}
                                                    />
                                                    <button
                                                        onClick={
                                                            handleVerifyOTP
                                                        }
                                                        className="flex w-full items-center justify-center rounded-md bg-gradient-to-r from-red-600 to-red-800 py-2 hover:from-red-700 hover:to-red-900"
                                                    >
                                                        <Check className="mr-2 h-4 w-4" />
                                                        Verify Email
                                                    </button>
                                                    <button
                                                        onClick={() =>
                                                            setOtpSent(false)
                                                        }
                                                        className="w-full rounded-md border border-zinc-700 bg-transparent py-2 text-zinc-400 hover:bg-zinc-800 hover:text-white"
                                                    >
                                                        Change Email
                                                    </button>
                                                </div>
                                            )}
                                        </div>
                                    ) : user && user.Verified && !user.Queue ? (
                                        <div className="rounded-lg border border-zinc-700 bg-zinc-800/50 p-6">
                                            <h3 className="mb-4 font-medium">
                                                Request Your Data
                                            </h3>
                                            <p className="mb-4 text-sm text-zinc-400">
                                                Your email has been verified.
                                                You can now request your data
                                                download.
                                            </p>
                                            <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
                                                <div className="flex flex-col rounded-lg border border-zinc-700 bg-zinc-800/30 p-6">
                                                    <div className="mb-4 flex items-start justify-between">
                                                        <div>
                                                            <h3 className="mb-1 font-medium">
                                                                Complete Data
                                                                Archive
                                                            </h3>
                                                            <p className="text-sm text-zinc-400">
                                                                All your data in
                                                                JSON format
                                                            </p>
                                                        </div>
                                                        <div className="text-sm text-zinc-400">
                                                            ~4.2 MB
                                                        </div>
                                                    </div>
                                                    <div className="mt-auto">
                                                        <button
                                                            onClick={
                                                                handleDownloadQueue
                                                            }
                                                            className="flex w-full items-center justify-center rounded-md bg-gradient-to-r from-red-600 to-red-800 py-2 hover:from-red-700 hover:to-red-900"
                                                        >
                                                            <Download className="mr-2 h-4 w-4" />
                                                            Request Download
                                                        </button>
                                                    </div>
                                                </div>

                                                <div className="flex flex-col rounded-lg border border-zinc-700 bg-zinc-800/30 p-6">
                                                    <div className="mb-4 flex items-start justify-between">
                                                        <div>
                                                            <h3 className="mb-1 font-medium">
                                                                User Activity
                                                                Report
                                                            </h3>
                                                            <p className="text-sm text-zinc-400">
                                                                Your activity
                                                                history in CSV
                                                                format
                                                            </p>
                                                        </div>
                                                        <div className="text-sm text-zinc-400">
                                                            ~1.8 MB
                                                        </div>
                                                    </div>
                                                    <div className="mt-auto">
                                                        <button
                                                            onClick={
                                                                handleDownloadQueue
                                                            }
                                                            className="flex w-full items-center justify-center rounded-md border border-zinc-700 py-2 text-zinc-300 hover:bg-zinc-700"
                                                        >
                                                            <Download className="mr-2 h-4 w-4" />
                                                            Request Download
                                                        </button>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    ) : user && user.Queue ? (
                                        <div className="rounded-lg border border-zinc-700 bg-zinc-800/50 p-6">
                                            <h3 className="mb-4 font-medium">
                                                Download Status
                                            </h3>
                                            <p className="mb-4 text-sm text-zinc-400">
                                                Your download request is being
                                                processed. We'll notify you when
                                                it's ready.
                                            </p>
                                            <div className="rounded-lg bg-zinc-900/50 p-4">
                                                <div className="mb-2 flex items-center justify-between">
                                                    <span>Status</span>
                                                    {getStatusDisplay(
                                                        user.Status,
                                                    )}
                                                </div>
                                                {user.Status === "done" && (
                                                    <div className="mt-4">
                                                        <button className="flex w-full items-center justify-center rounded-md bg-gradient-to-r from-green-600 to-green-800 py-2 hover:from-green-700 hover:to-green-900">
                                                            <Download className="mr-2 h-4 w-4" />
                                                            Download Now
                                                        </button>
                                                    </div>
                                                )}
                                            </div>
                                        </div>
                                    ) : null}

                                    <div className="mt-6 rounded-lg bg-zinc-800/30 p-4 text-sm text-zinc-500">
                                        <p>
                                            Need help with your data? Contact
                                            our support team at{" "}
                                            <span className="text-red-400">
                                                support@interrupted.me
                                            </span>
                                        </p>
                                    </div>
                                </div>
                            </div>
                        )}

                        {activeTab === "team" && (
                            <div className="space-y-6">
                                <div className="rounded-lg border border-zinc-800 bg-zinc-900/50 p-6">
                                    <h2 className="mb-4 text-xl font-bold">
                                        Team & Acknowledgments
                                    </h2>
                                    <p className="mb-6 text-zinc-400">
                                        The amazing people who made this project
                                        possible
                                    </p>

                                    <div className="space-y-6">
                                        <div>
                                            <h3 className="mb-4 text-lg font-medium text-red-400">
                                                Core Team
                                            </h3>
                                            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3">
                                                {[
                                                    {
                                                        name: "Alex Johnson",
                                                        role: "Design Lead",
                                                        years: "5 years",
                                                    },
                                                    {
                                                        name: "Taylor Smith",
                                                        role: "Lead Developer",
                                                        years: "5 years",
                                                    },
                                                    {
                                                        name: "Jordan Parker",
                                                        role: "Project Manager",
                                                        years: "4 years",
                                                    },
                                                    {
                                                        name: "Morgan Lee",
                                                        role: "UX Researcher",
                                                        years: "3 years",
                                                    },
                                                    {
                                                        name: "Riley Brown",
                                                        role: "Support Lead",
                                                        years: "4 years",
                                                    },
                                                    {
                                                        name: "Casey Wilson",
                                                        role: "Marketing",
                                                        years: "3 years",
                                                    },
                                                ].map((person) => (
                                                    <div
                                                        key={person.name}
                                                        className="rounded-lg border border-zinc-700 bg-zinc-800/50 p-4"
                                                    >
                                                        <h4 className="font-medium">
                                                            {person.name}
                                                        </h4>
                                                        <p className="text-sm text-zinc-400">
                                                            {person.role}
                                                        </p>
                                                        <p className="mt-1 text-xs text-zinc-500">
                                                            {person.years}
                                                        </p>
                                                    </div>
                                                ))}
                                            </div>
                                        </div>

                                        <div>
                                            <h3 className="mb-4 text-lg font-medium text-red-400">
                                                Special Thanks
                                            </h3>
                                            <div className="rounded-lg border border-zinc-700 bg-zinc-800/30 p-6">
                                                <p className="mb-4 text-zinc-300">
                                                    We extend our deepest
                                                    gratitude to all users,
                                                    contributors, and partners
                                                    who supported us throughout
                                                    this journey.
                                                </p>
                                                <div className="flex flex-wrap gap-2">
                                                    {[
                                                        "Early Adopters",
                                                        "Beta Testers",
                                                        "Community Moderators",
                                                        "Content Creators",
                                                        "Translators",
                                                        "Open Source Contributors",
                                                    ].map((group) => (
                                                        <span
                                                            key={group}
                                                            className="rounded-full bg-zinc-700/50 px-3 py-1 text-sm text-zinc-300"
                                                        >
                                                            {group}
                                                        </span>
                                                    ))}
                                                </div>
                                            </div>
                                        </div>

                                        <div className="rounded-lg border border-red-900/20 bg-red-900/10 p-6 text-center">
                                            <h3 className="mb-3 text-xl font-medium">
                                                Thank You for Being Part of Our
                                                Journey
                                            </h3>
                                            <p className="mb-4 text-zinc-300">
                                                As we sunset this project, we
                                                want to express our sincere
                                                appreciation for your support
                                                and contributions. The impact of
                                                this project will continue long
                                                after its conclusion.
                                            </p>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        )}
                    </div>
                )}
            </main>

            <footer className="border-t border-zinc-900 px-4 py-6">
                <div className="container mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 md:flex-row">
                    <div className="flex items-center">
                        <div className="mr-3 h-5 w-5 rounded-full bg-gradient-to-b from-orange-300 to-red-500"></div>
                        <p className="text-zinc-400">
                            interrupted.me © 2020-2025
                        </p>
                    </div>
                    <div className="text-xs text-zinc-500">
                        All data will remain available for download until
                        December 31, 2025
                    </div>
                    <div className="flex gap-6 text-zinc-500">
                        <Link
                            href="#"
                            className="transition-colors hover:text-red-500"
                        >
                            Privacy
                        </Link>
                        <Link
                            href="#"
                            className="transition-colors hover:text-red-500"
                        >
                            Terms
                        </Link>
                        <Link
                            href="#"
                            className="transition-colors hover:text-red-500"
                        >
                            Support
                        </Link>
                    </div>
                </div>
            </footer>
        </div>
    );
}
