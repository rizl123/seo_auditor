import { BASE_URL, OIDC_URL } from "./urls";

export const ISSUER = `${OIDC_URL}/oidc`;
export const CLIENT_ID = process.env.OIDC_ID || "";
export const CLIENT_SECRET = process.env.OIDC_SECRET || "";
export const REDIRECT_URI = `${BASE_URL}/callback`;
export const POST_LOGOUT_REDIRECT_URI = BASE_URL;
export const SCOPE = "openid profile email";
