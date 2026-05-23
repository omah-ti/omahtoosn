import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

const PROTECTED_PREFIXES = ["/dashboard", "/result", "/tryout"];
const AUTH_ROUTES = ["/login", "/signup", "/forgot-password"];
const DEFAULT_BACKEND_URL = "https://omahtoosn-production-7e28.up.railway.app";
const ACCESS_TOKEN_REFRESH_SKEW_MS = 30_000;

type RefreshResult = {
  ok: boolean;
  setCookies: string[];
};

export async function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const accessToken = request.cookies.get("access_token")?.value;
  const refreshToken = request.cookies.get("refresh_token")?.value;
  const hasRefreshToken = Boolean(refreshToken);
  const accessTokenExpired = isAccessTokenExpired(accessToken);
  const hasUsableAccessToken = Boolean(accessToken) && !accessTokenExpired;
  const shouldRefreshSession = hasRefreshToken && (!accessToken || accessTokenExpired);

  let refreshResult: RefreshResult | null = null;
  let isLoggedIn = hasUsableAccessToken || hasRefreshToken;

  if (shouldRefreshSession) {
    refreshResult = await refreshSession(request);
    isLoggedIn = refreshResult.ok;
  } else if (accessTokenExpired) {
    refreshResult = { ok: false, setCookies: expiredAuthCookies(request) };
    isLoggedIn = false;
  }

  const isProtected = PROTECTED_PREFIXES.some(
    (prefix) => pathname === prefix || pathname.startsWith(prefix + "/")
  );

  const isAuthRoute = AUTH_ROUTES.some(
    (route) => pathname === route || pathname.startsWith(route + "/")
  );

  let response: NextResponse;

  if (isProtected && !isLoggedIn) {
    const loginUrl = new URL("/login", request.url);
    response = NextResponse.redirect(loginUrl);
  } else if (isAuthRoute && isLoggedIn) {
    const dashboardUrl = new URL("/dashboard", request.url);
    response = NextResponse.redirect(dashboardUrl);
  } else {
    response = NextResponse.next();
  }

  if (refreshResult) {
    appendSetCookies(response, refreshResult.setCookies, request);
  }

  return response;
}

function isAccessTokenExpired(accessToken?: string) {
  if (!accessToken) {
    return false;
  }

  const payload = decodeJwtPayload(accessToken);
  if (!payload || typeof payload.exp !== "number") {
    return false;
  }

  return payload.exp * 1000 <= Date.now() + ACCESS_TOKEN_REFRESH_SKEW_MS;
}

function decodeJwtPayload(token: string): Record<string, unknown> | null {
  const parts = token.split(".");
  if (parts.length < 2) {
    return null;
  }

  try {
    const base64 = parts[1].replace(/-/g, "+").replace(/_/g, "/");
    const padded = base64.padEnd(Math.ceil(base64.length / 4) * 4, "=");
    return JSON.parse(atob(padded)) as Record<string, unknown>;
  } catch {
    return null;
  }
}

async function refreshSession(request: NextRequest): Promise<RefreshResult> {
  try {
    const headers = new Headers();
    const cookieHeader = request.headers.get("cookie");
    const deviceId = request.headers.get("x-device-id");

    if (cookieHeader) {
      headers.set("cookie", cookieHeader);
    }

    if (deviceId) {
      headers.set("x-device-id", deviceId);
    }

    const response = await fetch(`${backendBaseUrl()}${"/api/v1/auth/refresh"}`, {
      method: "POST",
      headers,
      cache: "no-store",
      redirect: "manual",
    });

    if (!response.ok) {
      return { ok: false, setCookies: expiredAuthCookies(request) };
    }

    const setCookies = getSetCookies(response.headers);
    if (setCookies.length === 0) {
      return { ok: false, setCookies: expiredAuthCookies(request) };
    }

    return { ok: true, setCookies };
  } catch {
    return { ok: false, setCookies: expiredAuthCookies(request) };
  }
}

function backendBaseUrl() {
  return (process.env.BACKEND_API_URL || process.env.NEXT_PUBLIC_BACKEND_API_URL || DEFAULT_BACKEND_URL).replace(/\/+$/, "");
}

function appendSetCookies(response: NextResponse, setCookies: string[], request: NextRequest) {
  for (const cookie of setCookies) {
    response.headers.append("set-cookie", normalizeSetCookie(cookie, request));
  }
}

function expiredAuthCookies(request: NextRequest) {
  return ["access_token", "refresh_token"].map((name) => {
    let cookie = `${name}=; Path=/; Max-Age=0; HttpOnly`;
    if (request.nextUrl.protocol === "https:") {
      cookie += "; Secure; SameSite=None";
    } else {
      cookie += "; SameSite=Lax";
    }
    return cookie;
  });
}

function getSetCookies(headers: Headers) {
  const withGetSetCookie = headers as Headers & {
    getSetCookie?: () => string[];
  };

  if (typeof withGetSetCookie.getSetCookie === "function") {
    return withGetSetCookie.getSetCookie();
  }

  const value = headers.get("set-cookie");
  if (!value) {
    return [];
  }

  return value.split(/,(?=\s*[^;,]+=)/).map((cookie) => cookie.trim());
}

function normalizeSetCookie(cookie: string, request: NextRequest) {
  let nextCookie = cookie.replace(/;\s*Domain=[^;]*/gi, "");
  if (request.nextUrl.protocol === "http:") {
    nextCookie = nextCookie
      .replace(/;\s*Secure/gi, "")
      .replace(/;\s*SameSite=None/gi, "; SameSite=Lax");
  }

  if (!/;\s*Path=/i.test(nextCookie)) {
    nextCookie += "; Path=/";
  }

  return nextCookie;
}

export const config = {
  matcher: [
    "/dashboard/:path*",
    "/result/:path*",
    "/tryout/:path*",
    "/login",
    "/signup",
    "/forgot-password",
  ],
};
