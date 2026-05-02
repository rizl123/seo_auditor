import { cookies } from "next/headers";
import * as authConfig from "@/config/auth";
import { exchangeCodeForTokens, getLoginUrl, getLogoutUrl } from "../oidc";
import * as ac from "./context";
import { deleteSession, setSession } from "./session";

export type AuthResult = "success" | "csrf_error" | "unknown_error";

export async function processCallback(url: URL): Promise<AuthResult> {
  const stateFromUrl = url.searchParams.get("state");
  const { verifier, state: savedState } = await ac.getAuthContext();

  if (!verifier || !stateFromUrl || stateFromUrl !== savedState) {
    return "csrf_error";
  }

  try {
    const tokens = await exchangeCodeForTokens(url, verifier, savedState);

    if (!tokens.id_token) {
      return "unknown_error";
    }

    await setSession(tokens.id_token);

    await ac.clearAuthContext();

    return "success";
  } catch (error) {
    console.error("Auth error:", error);
    return "unknown_error";
  }
}

export async function processLogin(): Promise<URL> {
  const { authorizationUrl, codeVerifier, state } = await getLoginUrl();

  await ac.setAuthContext(codeVerifier, state);

  return authorizationUrl;
}

export async function processLogout(): Promise<URL | undefined> {
  const cookieStore = await cookies();
  const idToken = cookieStore.get(authConfig.SESSION_COOKIE)?.value;

  await deleteSession();

  if (idToken) {
    try {
      return await getLogoutUrl(idToken);
    } catch (e) {
      console.error("Logout URL generation failed", e);
    }
  }

  return;
}
