import { Clock } from "lucide-react";
import { useTranslations } from "next-intl";
import { Card } from "@/components/Card";
import { ImageLoader } from "@/components/ImageLoader";
import type { ScanResult } from "@/types/report";
import { DetailItemComponent } from "./DetailItemComponent";
import { ProblemItem } from "./ProblemItem";

interface ScannerSectionProps {
  result: ScanResult;
}

export function ScannerSection({ result }: ScannerSectionProps) {
  const t = useTranslations("Report");
  const a = useTranslations("API");
  const r = useTranslations(`API.${result.i18n_namespace}`);

  const hasProblems = result.problems && result.problems.length > 0;

  const imageDetails = result.details.filter((d) => d.type === "image");
  const regularDetails = result.details.filter((d) => d.type !== "image");

  return (
    <div className="space-y-4">
      <div className="flex items-end justify-between px-2">
        <div className="space-y-1">
          <div className="flex items-center gap-2">
            <h3 className="text-lg font-black uppercase tracking-tight text-zinc-800 dark:text-zinc-200">
              {r("name")}
            </h3>
            {hasProblems ? (
              <span className="text-xs bg-rose-500 text-white px-2 py-0.5 rounded-full font-bold">
                {t("issues", { count: result.problems.length })}
              </span>
            ) : (
              <span className="text-xs bg-emerald-500 text-white px-2 py-0.5 rounded-full font-bold">
                {t("passed")}
              </span>
            )}
          </div>
          <p className="text-xs text-zinc-500 max-w-xl">{r("description")}</p>
        </div>
        <div className="hidden sm:flex items-center gap-1.5 text-xs text-zinc-400 font-mono bg-zinc-50 dark:bg-zinc-900 px-2 py-1 rounded-md">
          <Clock size={12} />
          {t("scannedAt", { time: new Date(result.scanned_at) })}
        </div>
      </div>

      <Card className="divide-y divide-zinc-100 dark:divide-zinc-800 border-t-2 border-t-zinc-200 dark:border-t-zinc-700">
        {regularDetails.length > 0 && (
          <div className="p-6 grid grid-cols-1 sm:grid-cols-2 gap-x-8 gap-y-4 bg-zinc-50/30 dark:bg-zinc-900/10">
            {regularDetails.map((detail) => (
              <DetailItemComponent key={detail.i18n_label} item={detail} />
            ))}
          </div>
        )}

        {imageDetails.length > 0 && (
          <div className="p-6 bg-zinc-50/50 dark:bg-zinc-900/20 border-t border-zinc-100 dark:border-zinc-800">
            {imageDetails.map((detail) => (
              <div key={detail.i18n_label} className="space-y-3">
                <p className="text-xs font-bold text-zinc-400 uppercase tracking-tighter">
                  {a(detail.i18n_label)}
                </p>
                {detail.value && detail.value.trim() !== "" ? (
                  <ImageLoader src={detail.value} alt={a(detail.i18n_label)} />
                ) : (
                  <div className="text-sm font-semibold">
                    <span className="text-zinc-300 italic">{t("noData")}</span>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}

        {hasProblems && (
          <div className="divide-y divide-zinc-50 dark:divide-zinc-800/50">
            {result.problems.map((problem) => (
              <ProblemItem key={problem.i18n_namespace} problem={problem} />
            ))}
          </div>
        )}
      </Card>
    </div>
  );
}
