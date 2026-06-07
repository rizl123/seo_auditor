import * as jose from "jose";
import { cookies } from "next/headers";
import * as authConfig from "@/config/auth";
import * as oidcConfig from "@/config/oidc";
import type { AuthUserInfo } from "@/types/auth";
import { COOKIE_OPTS } from "./context";

export async function getSession(): Promise<AuthUserInfo | null> {
  const cookieStore = await cookies();
  const idToken = cookieStore.get(authConfig.SESSION_COOKIE)?.value;

  if (!idToken) return null;

  try {
    const JWKS = jose.createRemoteJWKSet(new URL("/jwks", oidcConfig.ISSUER));
    const { payload } = await jose.jwtVerify(idToken, JWKS, {
      issuer: oidcConfig.ISSUER.toString(),
      audience: oidcConfig.CLIENT_ID,
    });
    return payload as AuthUserInfo;
  } catch {
    return null;
  }
}

export async function setSession(idToken: string) {
  const cookieStore = await cookies();
  cookieStore.set(authConfig.SESSION_COOKIE, idToken, {
    ...COOKIE_OPTS,
    maxAge: authConfig.SESSION_COOKIE_AGE,
  });
}

export async function deleteSession() {
  const cookieStore = await cookies();
  cookieStore.delete(authConfig.SESSION_COOKIE);
}
