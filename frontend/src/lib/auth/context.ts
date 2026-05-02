import { cookies } from "next/headers";
import * as authConfig from "@/config/auth";

export const COOKIE_OPTS = {
  httpOnly: true,
  secure: authConfig.COOKIE_SECURE,
  sameSite: "lax" as const,
  path: "/",
};

export async function setAuthContext(verifier: string, state: string) {
  const cookieStore = await cookies();
  cookieStore.set(authConfig.PKCE_COOKIE, verifier, {
    ...COOKIE_OPTS,
    maxAge: authConfig.PKCE_COOKIE_AGE,
  });
  cookieStore.set(authConfig.STATE_COOKIE, state, {
    ...COOKIE_OPTS,
    maxAge: authConfig.PKCE_COOKIE_AGE,
  });
}

export async function getAuthContext() {
  const cookieStore = await cookies();
  return {
    verifier: cookieStore.get(authConfig.PKCE_COOKIE)?.value,
    state: cookieStore.get(authConfig.STATE_COOKIE)?.value,
  };
}

export async function clearAuthContext() {
  const cookieStore = await cookies();
  cookieStore.delete(authConfig.PKCE_COOKIE);
  cookieStore.delete(authConfig.STATE_COOKIE);
}
