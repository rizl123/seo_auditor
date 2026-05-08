import type { Metadata } from "next";
import "./globals.css";
import { NextIntlClientProvider } from "next-intl";
import { getMessages } from "next-intl/server";
import { Navbar } from "@/components/Navbar";
import { type Locale, locales } from "@/config/i18n";

export const metadata: Metadata = {
  title: "SEO Analyzer",
  description: "Check your website on-page SEO elements",
};

export function generateStaticParams() {
  return locales.map((l) => ({ locale: l.short }));
}

export type LocaleLayoutProps = Readonly<{
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}>;

export default async function LocaleLayout(props: LocaleLayoutProps) {
  const locale = (await props.params).locale as Locale;
  const messages = await getMessages();

  return (
    <html lang={locale} className="h-full antialiased">
      <body className="min-h-full flex flex-col">
        <NextIntlClientProvider messages={messages}>
          <Navbar />
          {props.children}
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
