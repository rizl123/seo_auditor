import { ShieldX } from "lucide-react";
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

export default async function CSRFErrorPage() {
  const t = await getTranslations("AuthErrors.csrf");

  return (
    <ErrorCard
      title={t("title")}
      description={t("description")}
      icon={ShieldX}
      variant="amber"
      actionLabel={t("action")}
      actionHref="/login"
    />
  );
}
