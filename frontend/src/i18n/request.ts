import { notFound } from "next/navigation";
import type { GetRequestConfigParams } from "next-intl/server";
import { getRequestConfig } from "next-intl/server";
import { cache } from "react";
import type { Locale } from "@/config/i18n";
import { API_URL } from "@/config/urls";
import { routing } from "@/i18n/routing";

async function getFrontendMessages(locale: string) {
  const module = await import(`./messages/${locale}.json`);
  return module.default as Record<string, unknown>;
}

const getApiMessages = cache(async (locale: Locale) => {
  const response = await fetch(`${API_URL}/api/locales/${locale}.json`, {
    next: { revalidate: 1800, tags: [`locale-${locale}`] },
    cache: "no-store",
  });

  if (!response.ok) {
    return {};
  }

  return (await response.json()) as Record<string, unknown>;
});

export default getRequestConfig(async (params: GetRequestConfigParams) => {
  let locale = (params.locale || (await params.requestLocale)) as Locale;

  if (!locale || !routing.locales.includes(locale)) {
    locale = routing.defaultLocale as Locale;
  }

  if (!routing.locales.includes(locale)) {
    notFound();
  }

  const [frontendMessages, apiMessages] = await Promise.all([
    getFrontendMessages(locale),
    getApiMessages(locale),
  ]);

  return {
    locale,
    messages: {
      ...frontendMessages,
      API: apiMessages,
    },
  };
});
