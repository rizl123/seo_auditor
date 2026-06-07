import { useTranslations } from "next-intl";
import type { AuditResult } from "@/types/report";
import { Card } from "../Card";
import { ImageLoader } from "../ImageLoader";
import { DetailItemComponent } from "./DetailItemComponent";
import { ProblemComponent } from "./ProblemComponent";

interface AuditCardProps {
  result: AuditResult;
}

export function AuditCard({ result }: AuditCardProps) {
  const t = useTranslations("Report");
  const a = useTranslations("API");

  const imageDetails = (result.details || []).filter((d) => d.type === "image");
  const details = (result.details || []).filter((d) => d.type !== "image");

  const hasProblems = result.problems && result.problems.length > 0;

  return (
    <Card className="divide-y divide-zinc-100 dark:divide-zinc-800 border-t-2 border-t-zinc-200 dark:border-t-zinc-700">
      {details.length > 0 && (
        <div className="p-6 grid grid-cols-1 sm:grid-cols-2 gap-x-8 gap-y-4 bg-zinc-50/30 dark:bg-zinc-900/10">
          {details.map((detail) => (
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
            <ProblemComponent key={problem.i18n_namespace} problem={problem} />
          ))}
        </div>
      )}
    </Card>
  );
}
