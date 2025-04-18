"use client";
import { useState, useEffect, useContext } from "react";
import { ArrowRight, LogIn, Users } from "lucide-react";
import Link from "next/link";
import { SPContext } from "@/modules/Body/Body";

export default function SunsetPage() {
    const [scrollY, setScrollY] = useState(0);
    const { scroll } = useContext(SPContext);

    useEffect(() => {
        if (scroll) {
            const handleScroll = (e: { scroll: number }) => {
                setScrollY(e.scroll);
            };

            scroll.on("scroll", handleScroll);
            return () => {
                scroll.off("scroll", handleScroll);
            };
        } else {
            const handleStandardScroll = () => {
                setScrollY(window.scrollY);
            };

            window.addEventListener("scroll", handleStandardScroll);
            return () =>
                window.removeEventListener("scroll", handleStandardScroll);
        }
    }, [scroll]);

    const sunPosition = Math.min(50 + scrollY / 10, 120);

    return (
        <div className="flex min-h-screen flex-col bg-[#0c0c0c] text-white">
            <div className="relative h-[70vh] overflow-hidden">
                <div className="absolute inset-0 bg-gradient-to-b from-[#0c0c0c] via-[#1a0f23] to-[#3d1635]" />

                <div
                    className="absolute h-24 w-24 rounded-full bg-gradient-to-b from-orange-300 to-red-500 transition-all duration-700"
                    style={{
                        left: "calc(50% - 3rem)",
                        top: `${sunPosition}%`,
                        boxShadow: "0 0 60px rgba(255, 50, 50, 0.5)",
                    }}
                />

                <div className="absolute bottom-0 h-[1px] w-full bg-gradient-to-r from-transparent via-red-800/30 to-transparent" />

                <div className="absolute inset-0 flex items-center justify-center">
                    <div
                        data-aos="fade-in"
                        data-aos-delay="500"
                        className="text-center"
                    >
                        <h1 className="mb-4 text-5xl font-bold tracking-tighter md:text-7xl">
                            <span className="text-red-500">Interrupted.me</span>
                        </h1>
                        <p className="text-xl text-zinc-200">
                            Our journey comes to a beautiful end
                        </p>
                    </div>
                </div>
            </div>

            <div className="container mx-auto max-w-4xl flex-1 px-4 py-12">
                <div
                    data-aos="fade-up"
                    data-aos-duration="800"
                    className="mb-16"
                >
                    <h2 className="mb-6 text-3xl font-light tracking-tighter md:text-5xl">
                        The project{" "}
                        <span className="text-red-500">concludes</span>,<br />
                        but its impact{" "}
                        <span className="text-red-400">lives on</span>
                    </h2>

                    <p className="mb-8 text-lg leading-relaxed font-light text-zinc-400 md:text-xl">
                        After years of dedication and collaboration, we're proud
                        to announce the completion of our journey. We've created
                        something meaningful together, and now it's time to
                        celebrate our achievements.
                    </p>

                    <div className="mb-16 flex flex-wrap items-center gap-4">
                        <Link href="/login" className="inline-block">
                            <button className="group flex items-center rounded-full bg-gradient-to-r from-red-600 to-red-800 px-6 py-2 font-medium transition-all hover:from-red-700 hover:to-red-900">
                                Access Your Data{" "}
                                <LogIn className="ml-2 h-4 w-4 transition-transform group-hover:translate-x-1" />
                            </button>
                        </Link>
                        <Link href="#stats" className="inline-block">
                            <button className="flex items-center rounded-full border border-zinc-800 px-6 py-2 font-medium text-zinc-400 transition-all hover:bg-zinc-800 hover:text-white">
                                View Project Stats
                            </button>
                        </Link>
                    </div>
                </div>

                <section
                    id="stats"
                    data-aos="fade-in"
                    className="mb-20 border-t border-zinc-900 py-12"
                >
                    <h2 className="mb-10 text-center text-3xl font-medium">
                        Our (Small) <span className="text-red-500">Impact</span>
                    </h2>

                    <div className="grid grid-cols-1 gap-8 md:grid-cols-3">
                        <div
                            data-aos="fade-up"
                            data-aos-delay="100"
                            className="rounded-xl border border-zinc-800 bg-zinc-900/50 p-6 text-center"
                        >
                            <div className="mb-2 text-4xl font-bold text-red-500">
                                37K+
                            </div>
                            <div className="text-zinc-400">Total Uploads</div>
                        </div>
                        <div
                            data-aos="fade-up"
                            data-aos-delay="200"
                            className="rounded-xl border border-zinc-800 bg-zinc-900/50 p-6 text-center"
                        >
                            <div className="mb-2 text-4xl font-bold text-red-500">
                                566
                            </div>
                            <div className="text-zinc-400">Active Users</div>
                        </div>
                        <div
                            data-aos="fade-up"
                            data-aos-delay="300"
                            className="rounded-xl border border-zinc-800 bg-zinc-900/50 p-6 text-center"
                        >
                            <div className="mb-2 text-4xl font-bold text-red-500">
                                72.21 GB
                            </div>
                            <div className="text-zinc-400">Storage Used</div>
                        </div>
                    </div>
                </section>

                <section data-aos="fade-in" className="mb-20">
                    <h2 className="mb-10 text-center text-3xl font-medium">
                        Thank <span className="text-red-500">You</span>
                    </h2>

                    <div className="grid grid-cols-1 gap-8 md:grid-cols-2">
                        <div
                            data-aos="fade-right"
                            className="rounded-xl border border-zinc-800 bg-zinc-900/30 p-8"
                        >
                            <h3 className="mb-4 flex items-center text-xl font-medium">
                                <Users className="mr-2 text-red-500" /> To Our
                                Amazing Users
                            </h3>
                            <p className="mb-4 text-zinc-400">
                                Your support, feedback, and enthusiasm made this
                                project possible. We're grateful for every
                                moment you spent with us.
                            </p>
                            
                        </div>

                        <div
                            data-aos="fade-left"
                            className="rounded-xl border border-zinc-800 bg-zinc-900/30 p-8"
                        >
                            <h3 className="mb-4 flex items-center text-xl font-medium">
                                <Users className="mr-2 text-red-500" /> To Our
                                Dedicated Team
                            </h3>
                            <p className="mb-4 text-zinc-400">
                                Your hard work, creativity, and perseverance
                                brought this vision to life. We couldn't have
                                done it without each one of you.
                            </p>
                            <div className="flex flex-wrap gap-2">
                                {[
                                    "dopamine",
                                    "zero",
                                    "Pota",
                                    "edgebug",
                                    "entroxx",
                                ].map((name) => (
                                    <span
                                        key={name}
                                        className="rounded-full bg-zinc-800 px-3 py-1 text-sm text-zinc-300"
                                    >
                                        {name}
                                    </span>
                                ))}
                                <span className="rounded-full bg-zinc-800 px-3 py-1 text-sm text-zinc-300">
                                    and more
                                </span>
                            </div>
                        </div>
                    </div>
                </section>

                <div
                    data-aos="fade-up"
                    className="border-t border-zinc-900 py-12 text-center"
                >
                    <h2 className="mb-6 text-3xl font-medium">
                        Not without saying{" "}
                        <span className="text-red-500">Goodbye</span>
                    </h2>
                    <p className="mx-auto mb-8 max-w-2xl text-zinc-400">
                        Your data will remain available for download until April
                        30, 2025. Log in to access your information and export
                        everything you need.
                    </p>
                    <Link href="/login" className="inline-block">
                        <button className="flex items-center rounded-full bg-gradient-to-r from-red-600 to-red-800 px-8 py-3 text-lg font-medium transition-all hover:from-red-700 hover:to-red-900">
                            Log in to Dashboard{" "}
                            <ArrowRight className="ml-2 h-4 w-4" />
                        </button>
                    </Link>
                </div>
            </div>

            <footer className="border-t border-zinc-900 px-4 py-8">
                <div className="container mx-auto max-w-4xl">
                    <div className="flex flex-col items-center justify-between gap-4 md:flex-row">
                        <div className="flex items-center">
                            <div className="mr-3 h-6 w-6 rounded-full bg-gradient-to-b from-orange-300 to-red-500"></div>
                            <p className="text-zinc-400">
                                interrupted.me © 2023-2025
                            </p>
                        </div>
                        <div className="flex gap-6 text-zinc-500">
                            <Link
                                href="/dashboard"
                                className="transition-colors hover:text-red-500"
                            >
                                Dashboard
                            </Link>
                        </div>
                    </div>
                </div>
            </footer>
        </div>
    );
}
