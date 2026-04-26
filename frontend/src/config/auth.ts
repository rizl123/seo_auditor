import { BASE_URL, OIDC_URL } from "./env";

export const authConfig = {
  issuer: `${OIDC_URL}/oidc`,
  clientId: process.env.OIDC_ID || "",
  clientSecret: process.env.OIDC_SECRET || "",
  baseUrl: BASE_URL,
  redirectUri: `${BASE_URL}/callback`,
  postLogoutRedirectUri: BASE_URL,
  scope: "openid profile email",
};
