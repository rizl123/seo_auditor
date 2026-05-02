import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";
import { BASE_URL } from "@/config/urls";
import type { AuthResult } from "@/lib/auth/process";
import { processCallback } from "@/lib/auth/process";

const AUTH_ROUTES: Record<AuthResult, string> = {
  success: "/",
  csrf_error: "/auth-failed/csrf",
  unknown_error: "/auth-failed/unknown",
};

export async function GET(request: NextRequest) {
  const lhost3000Url = new URL(request.url);
  const url = new URL(lhost3000Url.pathname + lhost3000Url.search, BASE_URL);

  const result = await processCallback(url);

  return NextResponse.redirect(new URL(AUTH_ROUTES[result], BASE_URL));
}
