import { notFound } from "next/navigation";
import type {
  GetRequestConfigParams as Params,
  RequestConfig,
} from "next-intl/server";
import { getRequestConfig } from "next-intl/server";
import type { Locale } from "@/config/i18n";
import { API_URL } from "@/config/urls";
import { routing } from "@/i18n/routing";

async function createRequestConfig(params: Params): Promise<RequestConfig> {
  const l = params.locale || (await params.requestLocale);
  const locale = (l || routing.defaultLocale) as Locale;
  if (!routing.locales.includes(locale)) notFound();

  const frontendLocale = (await import(`./messages/${locale}.json`)).default;

  const apiFetched = await fetch(`${API_URL}/api/locales/${locale}.json`);
  const API = await apiFetched.json();

  return {
    locale,
    messages: { ...frontendLocale, API },
  };
}

export default getRequestConfig(createRequestConfig);
