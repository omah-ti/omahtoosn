import { NextRequest, NextResponse } from "next/server";
import { backendBaseUrl } from "@/lib/backend";

type RouteContext = {
  params: Promise<{
    path: string[];
  }>;
};

export async function GET(request: NextRequest, context: RouteContext) {
  const { path } = await context.params;
  const targetUrl = new URL(`${backendBaseUrl()}/question-assets/${path.join("/")}`);
  targetUrl.search = request.nextUrl.search;

  const response = await fetch(targetUrl, {
    cache: "force-cache",
  });

  return new NextResponse(response.body, {
    status: response.status,
    statusText: response.statusText,
    headers: response.headers,
  });
}
