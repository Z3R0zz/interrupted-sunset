"use client";
import {
    useState,
    useRef,
    useEffect,
    KeyboardEvent,
    ClipboardEvent,
} from "react";

/*
Copyright © 2025 Kars (github.com/kars1996)

Not to be shared, replicated or used without prior consent.
Contact Kars for any enquiries
*/

interface OTPInputProps {
    length?: number;
    onComplete?: (code: string) => void;
    onSend?: () => Promise<void>;
    cooldownTime?: number; // seconds
    inputClassName?: string;
    containerClassName?: string;
}

export default function OTPInput({
    length = 6,
    onComplete,
    onSend,
    cooldownTime = 60,
    inputClassName,
    containerClassName,
}: OTPInputProps) {
    const [otp, setOtp] = useState<string[]>(Array(length).fill(""));
    const [cooldown, setCooldown] = useState<number>(0);
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const inputRefs = useRef<(HTMLInputElement | null)[]>([]);

    useEffect(() => {
        inputRefs.current = inputRefs.current.slice(0, length);
        setOtp(Array(length).fill(""));
    }, [length]);

    useEffect(() => {
        if (inputRefs.current[0]) {
            inputRefs.current[0]?.focus();
        }
    }, []);

    useEffect(() => {
        let timer: NodeJS.Timeout;
        if (cooldown > 0) {
            timer = setTimeout(() => setCooldown(cooldown - 1), 1000);
        }
        return () => {
            if (timer) clearTimeout(timer);
        };
    }, [cooldown]);

    useEffect(() => {
        const isComplete = otp.every((digit) => digit !== "");
        if (isComplete && onComplete) {
            onComplete(otp.join(""));
        }
    }, [otp, onComplete]);

    const handleChange = (
        e: React.ChangeEvent<HTMLInputElement>,
        index: number,
    ) => {
        const value = e.target.value;
        if (!value) return;

        const newValue = value.slice(-1);

        if (!/^\d+$/.test(newValue)) return;

        const newOtp = [...otp];
        newOtp[index] = newValue;
        setOtp(newOtp);

        if (index < length - 1 && inputRefs.current[index + 1]) {
            inputRefs.current[index + 1]?.focus();
        }
    };

    const handleKeyDown = (
        e: KeyboardEvent<HTMLInputElement>,
        index: number,
    ) => {
        if (e.key === "Backspace") {
            if (otp[index] === "") {
                if (index > 0 && inputRefs.current[index - 1]) {
                    inputRefs.current[index - 1]?.focus();
                }
            } else {
                const newOtp = [...otp];
                newOtp[index] = "";
                setOtp(newOtp);
            }
        }

        if (e.key === "ArrowLeft" && index > 0) {
            inputRefs.current[index - 1]?.focus();
        }

        if (e.key === "ArrowRight" && index < length - 1) {
            inputRefs.current[index + 1]?.focus();
        }
    };

    const handlePaste = (
        e: ClipboardEvent<HTMLInputElement>,
        index: number,
    ) => {
        e.preventDefault();
        const pastedData = e.clipboardData.getData("text/plain").trim();

        if (!pastedData) return;

        if (!/^\d+$/.test(pastedData)) return;

        const chars = pastedData.split("");
        const newOtp = [...otp];

        const maxChars = Math.min(chars.length, length - index);

        for (let i = 0; i < maxChars; i++) {
            if (index + i < length) {
                newOtp[index + i] = chars[i];
            }
        }

        setOtp(newOtp);

        const focusIndex = Math.min(index + maxChars, length - 1);
        inputRefs.current[focusIndex]?.focus();
    };

    const handleSendCode = async () => {
        if (cooldown > 0 || !onSend) return;

        setIsLoading(true);
        try {
            await onSend();
            setCooldown(cooldownTime);
        } catch (error) {
            console.error("Error sending code:", error);
        } finally {
            setIsLoading(false);
        }
    };

    const clearOTP = () => {
        setOtp(Array(length).fill(""));
        if (inputRefs.current[0]) {
            inputRefs.current[0]?.focus();
        }
    };

    return (
        <div className={`flex flex-col items-center ${containerClassName}`}>
            <div className="flex justify-center gap-2">
                {Array.from({ length }, (_, index) => (
                    <input
                        key={index}
                        type="text"
                        inputMode="numeric"
                        maxLength={1}
                        value={otp[index]}
                        ref={(el) => {
                            inputRefs.current[index] = el;
                        }}
                        onChange={(e) => handleChange(e, index)}
                        onKeyDown={(e) => handleKeyDown(e, index)}
                        onPaste={(e) => handlePaste(e, index)}
                        className={`h-16 w-12 rounded border-2 border-neutral-800 bg-neutral-900 text-center text-xl font-bold text-white transition-all focus:border-red-500 focus:outline-none focus:ring-2 focus:ring-red-500 ${inputClassName}`}
                        aria-label={`Digit ${index + 1}`}
                    />
                ))}
            </div>

            <div className="mt-8 flex items-center gap-4">
                <button
                    onClick={clearOTP}
                    className="rounded-md bg-neutral-800 px-4 py-2 text-white transition-colors hover:bg-neutral-700"
                >
                    Clear
                </button>

                <button
                    onClick={handleSendCode}
                    disabled={cooldown > 0 || isLoading}
                    className={`flex items-center gap-2 rounded-md px-6 py-2 font-medium transition-colors ${cooldown > 0 ? "cursor-not-allowed bg-neutral-700 text-neutral-400" : "bg-red-600 text-white hover:bg-red-700"}`}
                >
                    {isLoading ? (
                        <>
                            <span className="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent"></span>
                            Sending...
                        </>
                    ) : cooldown > 0 ? (
                        `Resend in ${cooldown}s`
                    ) : (
                        "Send Code"
                    )}
                </button>
            </div>
        </div>
    );
}
