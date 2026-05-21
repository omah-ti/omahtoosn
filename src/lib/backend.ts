import { cookies } from "next/headers";

export type ApiEnvelope<T> = {
  success: boolean;
  message: string;
  data?: T;
  request_id?: string;
};

const DEFAULT_BACKEND_URL = "https://omahtoosn-production-7e28.up.railway.app";

export function backendBaseUrl() {
  return (process.env.BACKEND_API_URL || process.env.NEXT_PUBLIC_BACKEND_API_URL || DEFAULT_BACKEND_URL).replace(/\/+$/, "");
}

export async function backendFetch(path: string, init: RequestInit = {}) {
  const cookieStore = await cookies();
  const cookieHeader = cookieStore
    .getAll()
    .map((cookie) => `${cookie.name}=${cookie.value}`)
    .join("; ");

  const headers = new Headers(init.headers);
  if (cookieHeader && !headers.has("cookie")) {
    headers.set("cookie", cookieHeader);
  }

  return fetch(`${backendBaseUrl()}${path}`, {
    ...init,
    headers,
    cache: init.cache ?? "no-store",
  });
}

export async function backendJson<T>(path: string, init: RequestInit = {}) {
  const response = await backendFetch(path, init);
  const payload = (await response.json().catch(() => null)) as ApiEnvelope<T> | null;

  if (!response.ok || !payload?.success) {
    return null;
  }

  return payload.data ?? null;
}
