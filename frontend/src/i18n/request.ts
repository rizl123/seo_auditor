import { notFound } from "next/navigation";
import type { GetRequestConfigParams, RequestConfig } from "next-intl/server";
import { getRequestConfig } from "next-intl/server";
import type { Locale } from "@/config/i18n";
import { routing } from "@/i18n/routing";

async function createRequestConfig(
  params: GetRequestConfigParams,
): Promise<RequestConfig> {
  const l = params.locale || (await params.requestLocale);
  const locale = (l || routing.defaultLocale) as Locale;
  if (!routing.locales.includes(locale)) notFound();
  return {
    locale,
    messages: (await import(`./messages/${locale}.json`)).default,
  };
}

export default getRequestConfig(createRequestConfig);
