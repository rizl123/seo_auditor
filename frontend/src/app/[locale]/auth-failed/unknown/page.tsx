import { ShieldAlert } from "lucide-react";
import type { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import { ErrorCard } from "@/components/ErrorCard";

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("AuthErrors.unknown");

  return {
    title: t("title"),
    description: t("description"),
  };
}

export default async function UnknownErrorPage() {
  const t = await getTranslations("AuthErrors.unknown");

  return (
    <ErrorCard
      title={t("title")}
      description={t("description")}
      icon={ShieldAlert}
      variant="rose"
      actionLabel={t("action")}
      actionHref="/login"
    />
  );
}
