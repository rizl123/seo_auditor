import { useTranslations } from "next-intl";
import type { AuditFail } from "@/types/report";
import { Card } from "../Card";

export function AuditFailCard({ fail }: { fail: AuditFail }) {
  const a = useTranslations("API");

  return (
    <Card className="border-rose-200 dark:border-rose-900/50">
      <div className="p-6 space-y-3">
        <div>
          <h4 className="font-bold text-rose-600 dark:text-rose-400">
            {a(fail.title)}
          </h4>
        </div>

        <p className="text-sm text-zinc-600 dark:text-zinc-400">
          {a(fail.description)}
        </p>
      </div>
    </Card>
  );
}
