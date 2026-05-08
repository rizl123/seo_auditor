import createMiddleware from "next-intl/middleware";
import { routing } from "@/i18n/routing";

export default createMiddleware({
  locales: routing.locales,
  defaultLocale: routing.defaultLocale,
  localeDetection: true,
});

export const config = {
  matcher: ["/((?!api|trpc|_next|_vercel|login|logout|callback|.*\\..*).*)"],
};
