import * as oidc from "openid-client";
import * as authConfig from "@/config/auth";
import * as oidcConfig from "@/config/oidc";

let cachedConfig: oidc.Configuration | null = null;

export const options: oidc.DiscoveryRequestOptions &
  oidc.AuthorizationCodeGrantOptions = (() => {
  const execute: oidc.DiscoveryRequestOptions["execute"] = [];
  if (!authConfig.COOKIE_SECURE) {
    execute[0] = oidc.allowInsecureRequests;
  }
  return { execute };
})();

export async function getOidcConfig(): Promise<oidc.Configuration> {
  if (cachedConfig) return cachedConfig;
  const config = await oidc.discovery(
    oidcConfig.ISSUER,
    oidcConfig.CLIENT_ID,
    oidcConfig.CLIENT_SECRET,
    undefined,
    options,
  );
  cachedConfig = config;
  return config;
}

export async function getLoginUrl() {
  const config = await getOidcConfig();
  const codeVerifier = oidc.randomPKCECodeVerifier();
  const state = oidc.randomState();

  const authorizationUrl = oidc.buildAuthorizationUrl(config, {
    redirect_uri: oidcConfig.REDIRECT_URI.toString(),
    scope: oidcConfig.SCOPE,
    code_challenge: await oidc.calculatePKCECodeChallenge(codeVerifier),
    code_challenge_method: "S256",
    state: state,
  });

  return { authorizationUrl, codeVerifier, state };
}

export async function exchangeCodeForTokens(
  incomingUrl: URL,
  codeVerifier: string,
  expectedState: string,
) {
  const config = await getOidcConfig();
  return await oidc.authorizationCodeGrant(
    config,
    incomingUrl,
    {
      pkceCodeVerifier: codeVerifier,
      expectedState: expectedState,
    },
    {
      redirect_uri: oidcConfig.REDIRECT_URI.toString(),
    },
    options,
  );
}

export async function getLogoutUrl(idToken: string) {
  const config = await getOidcConfig();
  return oidc.buildEndSessionUrl(config, {
    id_token_hint: idToken,
    post_logout_redirect_uri: oidcConfig.POST_LOGOUT_REDIRECT_URI.toString(),
  });
}
