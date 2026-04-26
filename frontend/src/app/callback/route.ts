import { cookies } from "next/headers";
import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";
import * as oidc from "openid-client";
import { authConfig } from "@/config/auth";
import { getOidcConfig, options } from "@/lib/oidc";

export async function GET(request: NextRequest) {
  const incomingUrl = new URL(request.url);

  const fixedUrl = new URL(
    incomingUrl.pathname + incomingUrl.search,
    authConfig.baseUrl,
  );

  try {
    const config = await getOidcConfig();

    const cookieStore = await cookies();
    const code_verifier = cookieStore.get("cv")?.value;

    if (!code_verifier) {
      return new Response("Missing code verifier (cv cookie)", { status: 400 });
    }

    const tokens = await oidc.authorizationCodeGrant(
      config,
      fixedUrl,
      { pkceCodeVerifier: code_verifier },
      {
        client_id: authConfig.clientId,
        client_secret: authConfig.clientSecret,
        redirect_uri: authConfig.redirectUri,
      },
      options,
    );

    const response = NextResponse.redirect(new URL("/", authConfig.baseUrl));

    response.cookies.set("app_session", tokens.id_token || "", {
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
      path: "/",
      maxAge: 60 * 60 * 24,
    });

    response.cookies.delete("cv");

    return response;
  } catch {
    return new Response("Authentication failed.", { status: 500 });
  }
}
