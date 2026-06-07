import { notFound } from "next/navigation";
import type { GetRequestConfigParams } from "next-intl/server";
import { getRequestConfig } from "next-intl/server";
import { cache } from "react";
import type { Locale } from "@/config/i18n";
import { API_URL } from "@/config/urls";
import { routing } from "@/i18n/routing";

const getFrontendMessages = async (locale: string) => {
  const module = await import(`./messages/${locale}.json`);
  return module.default as Record<string, unknown>;
};

const getApiMessages = cache(async (locale: Locale) => {
  const response = await fetch(
    new URL(`/api/locales/${locale}.json`, API_URL),
    { next: { revalidate: 1800 } },
  );

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
