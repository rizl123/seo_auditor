import * as jose from "jose";
import { cookies } from "next/headers";
import { authConfig } from "@/config/auth";
import type { AuthUserInfo } from "@/types/auth";

export async function getSession(): Promise<AuthUserInfo | null> {
  const cookieStore = await cookies();
  const idToken = cookieStore.get("app_session")?.value;

  if (!idToken) return null;

  try {
    const JWKS = jose.createRemoteJWKSet(new URL(`${authConfig.issuer}/jwks`));
    const { payload } = await jose.jwtVerify(idToken, JWKS, {
      issuer: authConfig.issuer,
      audience: authConfig.clientId,
    });
    return payload as AuthUserInfo;
  } catch {
    return null;
  }
}
