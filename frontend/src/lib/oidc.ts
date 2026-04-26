import * as oidc from "openid-client";
import { authConfig } from "@/config/auth";

let cachedConfig: oidc.Configuration | null = null;

export async function getOidcConfig(): Promise<oidc.Configuration> {
  if (cachedConfig) return cachedConfig;

  const config = await oidc.discovery(
    new URL(authConfig.issuer),
    authConfig.clientId,
    authConfig.clientSecret,
    undefined,
    options,
  );

  cachedConfig = config;
  return config;
}

export const options: oidc.DiscoveryRequestOptions &
  oidc.AuthorizationCodeGrantOptions = {
  execute: [oidc.allowInsecureRequests],
};
