import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";
import { BASE_URL } from "@/config/urls";
import { processCallback } from "@/lib/auth/process";

const AUTH_ROUTES = {
  success: "/" as const,
  csrf_error: "/auth-failed/csrf" as const,
  unknown_error: "/auth-failed/unknown" as const,
};

export async function GET(request: NextRequest) {
  const lhost3000Url = new URL(request.url);
  const url = new URL(lhost3000Url.pathname + lhost3000Url.search, BASE_URL);

  const path = await processCallback(url);

  return NextResponse.redirect(new URL(AUTH_ROUTES[path], BASE_URL));
}
