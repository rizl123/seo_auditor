import { cookies } from "next/headers";
import { NextResponse } from "next/server";
import * as oidc from "openid-client";
import { authConfig } from "@/config/auth";
import { getOidcConfig } from "@/lib/oidc";

export async function GET() {
  const config = await getOidcConfig();
  const cookieStore = await cookies();
  const idToken = cookieStore.get("app_session")?.value;

  cookieStore.delete("app_session");

  if (idToken) {
    const logoutUrl = oidc.buildEndSessionUrl(config, {
      id_token_hint: idToken,
      post_logout_redirect_uri: authConfig.postLogoutRedirectUri,
    });
    return NextResponse.redirect(logoutUrl.href);
  }

  return NextResponse.redirect(new URL("/", authConfig.baseUrl));
}
