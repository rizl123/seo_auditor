export const PKCE_COOKIE = process.env.PKCE_COOKIE || "cv";
export const STATE_COOKIE = process.env.STATE_COOKIE || "st";
export const SESSION_COOKIE = process.env.SESSION_COOKIE || "app_session";

export const COOKIE_SECURE =
  (process.env.COOKIE_INSECURE || "false") !== "true";

export const SESSION_COOKIE_AGE =
  Number(process.env.SESSION_COOKIE_AGE) || 60 * 60 * 24;
export const PKCE_COOKIE_AGE = Number(process.env.PKCE_COOKIE_AGE) || 60 * 10;
