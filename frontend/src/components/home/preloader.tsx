"use client";
import { useState, useEffect } from "react";

export const SunsetPreloader = () => {
    const [isLoading, setIsLoading] = useState(true);
    const [isVisible, setIsVisible] = useState(true);
    const duration = 1000;

    useEffect(() => {
        const loadingTimer = setTimeout(() => {
            setIsLoading(false);

            const fadeTimer = setTimeout(() => {
                setIsVisible(false);
            }, duration);

            return () => clearTimeout(fadeTimer);
        }, 2000); // Adjust this time as needed

        return () => clearTimeout(loadingTimer);
    }, [duration]);

    if (!isVisible) return null;

    return (
        <div
            className={`fixed inset-0 z-50 flex flex-col items-center justify-center bg-[#0c0c0c] transition-opacity duration-1000 ease-in-out ${
                isLoading ? "opacity-100" : "opacity-0"
            }`}
            style={{
                transitionDuration: `${duration}ms`,
            }}
        >
            <div className="relative flex flex-col items-center">
                <div className="sun-container mb-12">
                    <div className="sunset-circle rounded-full"></div>
                    <div className="horizon z-[51]"></div>
                </div>

                <div className="-mt-14 flex flex-col items-center justify-center">
                    <h1 className="text-2xl font-bold text-white">
                        interrupted.me sunset
                    </h1>
                    <p className="mt-2 text-zinc-400">
                        Thank you for the journey
                    </p>
                </div>

                <div className="mt-8 flex space-x-2">
                    <div className="loading-dot"></div>
                    <div
                        className="loading-dot"
                        style={{ animationDelay: "0.2s" }}
                    ></div>
                    <div
                        className="loading-dot"
                        style={{ animationDelay: "0.4s" }}
                    ></div>
                </div>
            </div>

            <style jsx>{`
                .sun-container {
                    position: relative;
                    width: 120px;
                    height: 120px;
                    overflow: hidden;
                }

                .sunset-circle {
                    position: absolute;
                    width: 60px;
                    height: 60px;
                    border-radius: 50%;
                    background: linear-gradient(to bottom, #ff7e5f, #feb47b);
                    bottom: 35px; /* Changed from top: 30px to bottom: 30px */
                    left: 30px;
                    animation: sunset 3s infinite ease-in-out;
                    box-shadow: 0 0 15px rgba(255, 126, 95, 0.6);
                }

                .horizon {
                    position: absolute;
                    width: 120px;
                    height: 60px;
                    background: #0c0c0c;
                    bottom: 0;
                    border-top: 1px solid #e33a42;
                }

                @keyframes sunset {
                    0% {
                        transform: translateY(-20px);
                        opacity: 1;
                    }
                    100% {
                        transform: translateY(30px);
                        opacity: 0.6;
                    }
                }

                .loading-dot {
                    width: 10px;
                    height: 10px;
                    border-radius: 50%;
                    background: linear-gradient(to bottom, #ff7e5f, #feb47b);
                    animation: pulse 1.5s infinite ease-in-out;
                }

                @keyframes pulse {
                    0%,
                    100% {
                        transform: scale(0.8);
                        opacity: 0.5;
                    }
                    50% {
                        transform: scale(1.2);
                        opacity: 1;
                    }
                }
            `}</style>
        </div>
    );
};

export default SunsetPreloader;
