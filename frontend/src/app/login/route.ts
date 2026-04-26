import { cookies } from "next/headers";
import { NextResponse } from "next/server";
import * as oidc from "openid-client";
import { authConfig } from "@/config/auth";
import { getOidcConfig } from "@/lib/oidc";

export async function GET() {
  const config = await getOidcConfig();

  const code_verifier = oidc.randomPKCECodeVerifier();
  const code_challenge = await oidc.calculatePKCECodeChallenge(code_verifier);

  const authorizationUrl = oidc.buildAuthorizationUrl(config, {
    redirect_uri: authConfig.redirectUri,
    scope: authConfig.scope,
    code_challenge,
    code_challenge_method: "S256",
  });

  (await cookies()).set("cv", code_verifier, { httpOnly: true, path: "/" });

  return NextResponse.redirect(authorizationUrl.href);
}
