"use client";

export type ApiEnvelope<T> = {
  success: boolean;
  message: string;
  data?: T;
  request_id?: string;
};

type ApiSuccess<T> = { ok: true; data: T };
type ApiFailure = { ok: false; status: number; message: string };
export type ApiResult<T> = ApiSuccess<T> | ApiFailure;

function redirectToLogin() {
  if (typeof window !== "undefined") {
    window.location.href = "/login";
  }
}

export async function apiFetch<T = unknown>(
  path: string,
  init: RequestInit = {},
): Promise<ApiResult<T>> {
  const res = await fetch(path, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...init.headers,
    },
  });

  if (res.status === 401) {
    redirectToLogin();
    return { ok: false, status: 401, message: "Sesi telah berakhir, silakan login kembali" };
  }

  const payload = (await res.json().catch(() => null)) as ApiEnvelope<T> | null;

  if (!res.ok || !payload?.success) {
    return {
      ok: false,
      status: res.status,
      message: payload?.message || "Terjadi kesalahan",
    };
  }

  return { ok: true, data: payload.data as T };
}
