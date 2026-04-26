import * as jose from "jose";
import { cookies } from "next/headers";
import { NextResponse } from "next/server";
import { authConfig } from "@/config/auth";

export async function GET() {
  const cookieStore = await cookies();
  const idToken = cookieStore.get("app_session")?.value;

  if (!idToken) return NextResponse.json({ isAuthenticated: false });

  try {
    const JWKS = jose.createRemoteJWKSet(new URL(`${authConfig.issuer}/jwks`));

    const { payload } = await jose.jwtVerify(idToken, JWKS, {
      issuer: authConfig.issuer,
      audience: authConfig.clientId,
    });

    return NextResponse.json({
      isAuthenticated: true,
      user: payload,
    });
  } catch {
    return NextResponse.json({ isAuthenticated: false }, { status: 401 });
  }
}
