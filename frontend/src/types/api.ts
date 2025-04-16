export type APIRes<T = any> = {
    success: boolean;
    data?: T,
    error?: string;
}

export type LoginPost = {
    username: string,
    password: string,
}

export type LoginRes = {
    token: string,
}

export type UserRes = {
    Email: string,
    ID: number,
    Queue: boolean, // currently downloading shit
    Status: StatusTypes,
    Verified: boolean, // Email currently verified
}

export type StatusTypes = "waiting" | "processing" | "done" | "failed";

export type RequestEmailPost = {
    email: string,
}

export type RequestEmailRes = {
    message: string, // Please check your email for the OTP code
}

export type VerifyEmailPost = {
    code: string,
}

export type VerifyEmailRes = {
    message: string, // Email verified successfully
}

