import { toast } from "sonner";
import { AlertCircle, CheckCircle, Info, XCircle } from "lucide-react";

/*
Copyright © 2025 Kars (github.com/kars1996)

Not to be shared, replicated, or used without prior consent.
Contact me for any enquiries
*/

type ToastType = "success" | "error" | "warning" | "info";

interface ToastOptions {
    duration?: number;
    position?:
        | "top-left"
        | "top-right"
        | "bottom-left"
        | "bottom-right"
        | "top-center"
        | "bottom-center";
}

const ToastIcon = ({ type }: { type: ToastType }) => {
    switch (type) {
        case "success":
            return <CheckCircle className="h-5 w-5 text-green-500" />;
        case "error":
            return <XCircle className="h-5 w-5 text-red-500" />;
        case "warning":
            return <AlertCircle className="h-5 w-5 text-yellow-500" />;
        case "info":
            return <Info className="h-5 w-5 text-blue-500" />;
    }
};

export function showToast(
    type: ToastType,
    message: string,
    options?: ToastOptions,
) {
    return toast.custom(
        (t) => (
            <div
                className={`${
                    (t as any).visible
                        ? "animate-in fade-in slide-in-from-top-3"
                        : "animate-out fade-out slide-out-to-top-3"
                } pointer-events-auto relative max-w-md rounded-lg border shadow-lg`}
                style={{
                    backgroundColor: "rgba(20, 20, 25, 0.9)",
                    borderColor:
                        type === "success"
                            ? "rgba(74, 222, 128, 0.3)"
                            : type === "error"
                              ? "rgba(248, 113, 113, 0.3)"
                              : type === "warning"
                                ? "rgba(250, 204, 21, 0.3)"
                                : "rgba(96, 165, 250, 0.3)",
                }}
            >
                <div className="flex items-center gap-3 p-4">
                    <ToastIcon type={type} />
                    <div className="text-sm font-medium text-white">
                        {message}
                    </div>
                    <button
                        onClick={() => toast.dismiss((t as any).id)}
                        className="ml-auto rounded-md p-1 text-zinc-400 hover:bg-zinc-800 hover:text-white"
                    >
                        <XCircle className="h-4 w-4" />
                    </button>
                </div>
            </div>
        ),
        options,
    );
}
