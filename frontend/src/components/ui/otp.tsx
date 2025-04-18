"use client";
import {
    useState,
    useRef,
    useEffect,
    KeyboardEvent,
    ClipboardEvent,
} from "react";
import { Mail, AlertCircle, Check } from "lucide-react";

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
    const [error, setError] = useState<string>("");
    const [success, setSuccess] = useState<string>("");
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
        setError("");
        setSuccess("");
        try {
            await onSend();
            setCooldown(cooldownTime);
            setSuccess("OTP sent to your email");
        } catch (error) {
            console.error("Error sending code:", error);
            setError("Failed to send OTP. Please try again.");
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
            {error && (
                <div className="mb-4 flex w-full items-center rounded-md border border-red-800 bg-red-900/20 px-4 py-3 text-red-300">
                    <AlertCircle className="mr-2 h-4 w-4" />
                    {error}
                </div>
            )}

            {success && (
                <div className="mb-4 flex w-full items-center rounded-md border border-green-800 bg-green-900/20 px-4 py-3 text-green-300">
                    <Check className="mr-2 h-4 w-4" />
                    {success}
                </div>
            )}

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
                        className={`h-12 w-10 rounded-md border border-zinc-700 bg-zinc-900 text-center text-lg font-bold text-white transition-all focus:border-red-500 focus:outline-none focus:ring-2 focus:ring-red-500 md:h-14 md:w-12 ${inputClassName}`}
                        aria-label={`Digit ${index + 1}`}
                    />
                ))}
            </div>

            <div className="mt-6 flex items-center gap-4">
                <button
                    onClick={clearOTP}
                    className="rounded-md border border-zinc-700 bg-transparent px-4 py-2 text-zinc-300 transition-colors hover:bg-zinc-800"
                >
                    Clear
                </button>

                <button
                    onClick={handleSendCode}
                    disabled={cooldown > 0 || isLoading}
                    className={`flex items-center gap-2 rounded-md px-6 py-2 font-medium transition-colors ${
                        cooldown > 0
                            ? "cursor-not-allowed bg-zinc-700 text-zinc-400"
                            : "bg-gradient-to-r from-red-600 to-red-800 text-white hover:from-red-700 hover:to-red-900"
                    }`}
                >
                    {isLoading ? (
                        <>
                            <span className="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent"></span>
                            Sending...
                        </>
                    ) : cooldown > 0 ? (
                        `Resend in ${cooldown}s`
                    ) : (
                        <>
                            <Mail className="h-4 w-4" />
                            Send Code
                        </>
                    )}
                </button>
            </div>
        </div>
    );
}
