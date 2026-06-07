import { useTranslations } from "next-intl";
import type { AuditResult } from "@/types/report";
import { AuditCard } from "./AuditCard";
import { AuditFailCard } from "./AuditFailCard";

interface AuditSectionProps {
  result: AuditResult;
}

export function AuditSection({ result }: AuditSectionProps) {
  const t = useTranslations("Report");
  const r = useTranslations(`API.${result.i18n_namespace}`);

  const hasProblems = result.problems && result.problems.length > 0;

  return (
    <div className="space-y-4">
      <div className="flex items-end justify-between px-2">
        <div className="space-y-1">
          <div className="flex items-center gap-2">
            <h3 className="text-lg font-black uppercase tracking-tight text-zinc-800 dark:text-zinc-200">
              {r("name")}
            </h3>
            {result.fail ? (
              <span className="text-xs bg-rose-500 text-white px-2 py-0.5 rounded-full font-bold">
                {t("failed")}
              </span>
            ) : hasProblems ? (
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
        <div className="hidden sm:flex flex-col items-end rounded-xl border border-border-custom bg-card px-3 py-2 text-xs">
          <span className="text-muted-foreground">
            {t("startedAt", { time: new Date(result.started_at) })}
          </span>

          <span className="text-muted-foreground">
            {t("finishedAt", { time: new Date(result.finished_at) })}
          </span>
        </div>
      </div>

      {result.fail ? (
        <AuditFailCard fail={result.fail} />
      ) : (
        <AuditCard result={result} />
      )}
    </div>
  );
}
