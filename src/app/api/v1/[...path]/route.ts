import { NextRequest, NextResponse } from "next/server";
import { backendBaseUrl } from "@/lib/backend";

type RouteContext = {
  params: Promise<{
    path: string[];
  }>;
};

const HOP_BY_HOP_HEADERS = new Set([
  "connection",
  "content-length",
  "host",
  "keep-alive",
  "proxy-authenticate",
  "proxy-authorization",
  "te",
  "trailer",
  "transfer-encoding",
  "upgrade",
]);

/**
 * Refresh mutex — mencegah race condition ketika banyak request 401 barengan.
 * Hanya 1 refresh yang jalan; yang lain menunggu hasilnya.
 */
let refreshPromise: Promise<{ setCookies: string[]; mergedCookie: string } | null> | null = null;

async function acquireRefresh(
  request: NextRequest,
): Promise<{ setCookies: string[]; mergedCookie: string } | null> {
  if (refreshPromise) return refreshPromise;

  refreshPromise = (async () => {
    try {
      const incomingCookies = request.headers.get("cookie") || "";
      const refreshResponse = await forward(request, "/api/v1/auth/refresh", {
        method: "POST",
      });

      if (!refreshResponse.ok) {
        return null;
      }

      const setCookies = getSetCookies(refreshResponse.headers);
      const mergedCookie = mergeCookies(incomingCookies, setCookies);

      return { setCookies, mergedCookie };
    } finally {
      refreshPromise = null;
    }
  })();

  return refreshPromise;
}

async function proxy(request: NextRequest, context: RouteContext) {
  const { path } = await context.params;
  const targetPath = `/api/v1/${path.join("/")}`;
  const body = await cloneRequestBody(request);

  let backendResponse = await forward(request, targetPath, { body });

  if (backendResponse.status === 401 && targetPath !== "/api/v1/auth/refresh") {
    const refreshResult = await acquireRefresh(request);

    if (refreshResult) {
      backendResponse = await forward(request, targetPath, {
        body,
        cookieHeader: refreshResult.mergedCookie,
      });
      const response = buildResponse(backendResponse, request);
      for (const cookie of refreshResult.setCookies) {
        response.headers.append("set-cookie", normalizeSetCookie(cookie, request));
      }
      return response;
    }

    const response = buildResponse(backendResponse, request);
    response.headers.append("set-cookie", "access_token=; Path=/; Max-Age=0; HttpOnly");
    response.headers.append("set-cookie", "refresh_token=; Path=/; Max-Age=0; HttpOnly");
    return response;
  }

  return buildResponse(backendResponse, request);
}

async function forward(
  request: NextRequest,
  targetPath: string,
  options: {
    body?: ArrayBuffer;
    cookieHeader?: string;
    method?: string;
  } = {},
) {
  const targetUrl = new URL(`${backendBaseUrl()}${targetPath}`);
  targetUrl.search = request.nextUrl.search;

  const allowedHeaders = new Set([
    "content-type",
    "accept",
    "authorization",
    "origin",
    "x-device-id",
  ]);

  const headers = new Headers();
  request.headers.forEach((value, key) => {
    if (allowedHeaders.has(key.toLowerCase())) {
      headers.set(key, value);
    }
  });

  headers.set("cookie", options.cookieHeader ?? (request.headers.get("cookie") ?? ""));

  return fetch(targetUrl, {
    method: options.method ?? request.method,
    headers,
    body: options.body,
    cache: "no-store",
    redirect: "manual",
  });
}

function buildResponse(backendResponse: Response, request: NextRequest) {
  const responseHeaders = new Headers();
  backendResponse.headers.forEach((value, key) => {
    if (!HOP_BY_HOP_HEADERS.has(key.toLowerCase()) && key.toLowerCase() !== "set-cookie") {
      responseHeaders.set(key, value);
    }
  });

  const response = new NextResponse(backendResponse.body, {
    status: backendResponse.status,
    statusText: backendResponse.statusText,
    headers: responseHeaders,
  });

  for (const cookie of getSetCookies(backendResponse.headers)) {
    response.headers.append("set-cookie", normalizeSetCookie(cookie, request));
  }

  return response;
}

function mergeCookies(originalCookieHeader: string, setCookieStrings: string[]): string {
  const cookies = parseCookieHeader(originalCookieHeader);
  for (const str of setCookieStrings) {
    const firstSemi = str.indexOf(";");
    const nv = firstSemi >= 0 ? str.substring(0, firstSemi).trim() : str.trim();
    const eq = nv.indexOf("=");
    if (eq > 0) {
      cookies.set(nv.substring(0, eq).trim(), nv.substring(eq + 1).trim());
    }
  }
  return Array.from(cookies.entries())
    .map(([name, value]) => `${name}=${value}`)
    .join("; ");
}

function parseCookieHeader(cookieHeader: string): Map<string, string> {
  const map = new Map<string, string>();
  if (!cookieHeader) return map;
  cookieHeader.split(";").forEach((part) => {
    const trimmed = part.trim();
    if (!trimmed) return;
    const eq = trimmed.indexOf("=");
    if (eq > 0) {
      map.set(trimmed.substring(0, eq).trim(), trimmed.substring(eq + 1).trim());
    }
  });
  return map;
}

async function cloneRequestBody(request: NextRequest) {
  if (request.method === "GET" || request.method === "HEAD") {
    return undefined;
  }
  return request.arrayBuffer();
}

export const GET = proxy;
export const POST = proxy;
export const PUT = proxy;
export const OPTIONS = proxy;

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
