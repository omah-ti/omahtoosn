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

async function proxy(request: NextRequest, context: RouteContext) {
  const { path } = await context.params;
  const targetPath = `/api/v1/${path.join("/")}`;
  const body = await cloneRequestBody(request);

  let backendResponse = await forward(request, targetPath, body);

  if (backendResponse.status === 401 && targetPath !== "/api/v1/auth/refresh") {
    const refreshResponse = await forward(request, "/api/v1/auth/refresh", undefined);
    if (refreshResponse.ok) {
      const refreshedCookies = extractSetCookies(refreshResponse.headers);
      backendResponse = await forward(request, targetPath, body, refreshedCookies);
      const response = buildResponse(backendResponse, request);
      for (const cookie of refreshedCookies) {
        response.headers.append("set-cookie", normalizeSetCookie(cookie, request));
      }
      return response;
    }
  }

  return buildResponse(backendResponse, request);
}

async function forward(
  request: NextRequest,
  targetPath: string,
  body: ArrayBuffer | undefined,
  overrideCookies?: string[]
) {
  const targetUrl = new URL(`${backendBaseUrl()}${targetPath}`);
  targetUrl.search = request.nextUrl.search;

  const allowedHeaders = new Set([
    "content-type",
    "accept",
    "authorization",
    "cookie",
    "origin",
    "x-device-id",
  ]);

  const headers = new Headers();
  request.headers.forEach((value, key) => {
    if (allowedHeaders.has(key.toLowerCase())) {
      headers.set(key, value);
    }
  });

  if (overrideCookies && overrideCookies.length > 0) {
    headers.set("cookie", overrideCookies.join("; "));
  }

  return fetch(targetUrl, {
    method: request.method,
    headers,
    body,
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

function extractSetCookies(headers: Headers) {
  return getSetCookies(headers);
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
  return nextCookie;
}
