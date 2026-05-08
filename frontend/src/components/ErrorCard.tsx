import type { LucideIcon } from "lucide-react";
import { ArrowLeft } from "lucide-react";
import { getTranslations } from "next-intl/server";
import { Card } from "@/components/Card";

interface ErrorCardProps {
  title: string;
  description?: string;
  icon: LucideIcon;
  variant?: "amber" | "rose" | "blue";
  actionLabel?: string;
  actionHref?: string;
}

export async function ErrorCard({
  title,
  description,
  icon: Icon,
  variant = "rose",
  actionLabel,
  actionHref,
}: ErrorCardProps) {
  const t = await getTranslations("AuthErrors");

  const themes = {
    amber: {
      border: "border-t-amber-500",
      iconBg:
        "bg-amber-50 text-amber-600 dark:bg-amber-500/10 dark:text-amber-400",
    },
    rose: {
      border: "border-t-rose-500",
      iconBg: "bg-rose-50 text-rose-600 dark:bg-rose-500/10 dark:text-rose-400",
    },
    blue: {
      border: "border-t-blue-500",
      iconBg: "bg-blue-50 text-blue-600 dark:bg-blue-500/10 dark:text-blue-400",
    },
  };

  const theme = themes[variant];

  return (
    <main className="bg-zinc-50 dark:bg-zinc-950 flex flex-col items-center justify-center p-4">
      <div className="w-full max-w-md space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
        <Card className={`border-t-2 ${theme.border} shadow-lg`}>
          <div className="p-8 flex flex-col items-center text-center">
            <div className={`mb-6 p-4 rounded-full ${theme.iconBg}`}>
              <Icon size={40} strokeWidth={1.5} />
            </div>

            <div className="space-y-2 mb-8">
              <h2 className="font-bold text-lg text-zinc-900 dark:text-zinc-100">
                {title}
              </h2>
              {description && (
                <p className="text-sm text-zinc-500 dark:text-zinc-400">
                  {description}
                </p>
              )}
            </div>

            <div className="w-full space-y-3">
              {actionLabel && actionHref && (
                <a
                  href={actionHref}
                  className="flex items-center justify-center gap-2 w-full py-3 px-4 bg-zinc-900 dark:bg-zinc-100 text-zinc-50 dark:text-zinc-900 rounded-xl font-bold text-sm hover:opacity-90 transition-opacity"
                >
                  {actionLabel}
                </a>
              )}

              <a
                href="/"
                className="flex items-center justify-center gap-2 w-full py-3 px-4 bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 rounded-xl font-bold text-sm hover:bg-zinc-200 dark:hover:bg-zinc-700 transition-colors"
              >
                <ArrowLeft size={16} />
                {t("back")}
              </a>
            </div>
          </div>
        </Card>
      </div>
    </main>
  );
}
