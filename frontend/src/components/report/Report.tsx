"use client";

import { Globe } from "lucide-react";
import type { PageReport } from "@/types/report";
import { Card } from "./Card";
import { ScannerSection } from "./ScannerSection";
import { SpeedIndicator } from "./SpeedIndicator";

export function Report({ data }: { data: PageReport }) {
  const perfResult = data.results.find((r) => r.auditor_name === "performance");
  const responseTimeDetail = perfResult?.details.find(
    (d) => d.type === "duration_ms",
  );
  const responseTime = responseTimeDetail
    ? Number(responseTimeDetail.value)
    : null;

  return (
    <div className="max-w-3xl mx-auto space-y-12 py-10 animate-in fade-in slide-in-from-bottom-4 duration-700">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-6 px-2">
        <div className="min-w-0">
          <p className="text-[10px] font-bold text-zinc-400 uppercase tracking-widest mb-1">
            Target URL
          </p>
          <div className="flex items-center gap-2 text-xl font-bold break-all text-zinc-900 dark:text-white">
            <Globe size={20} className="text-accent shrink-0" />
            {data.url}
          </div>
        </div>

        {responseTime !== null && (
          <Card className="flex items-center gap-5 px-6 py-3 shrink-0">
            <div className="text-right">
              <p className="text-[10px] font-bold text-zinc-400 uppercase">
                Response
              </p>
              <p className="text-xl font-black">
                {responseTime}
                <span className="text-[10px] ml-0.5 font-normal text-zinc-400">
                  ms
                </span>
              </p>
            </div>
            <div className="h-8 w-px bg-zinc-100 dark:bg-zinc-800" />
            <div className="w-32">
              <SpeedIndicator ms={responseTime} />
            </div>
          </Card>
        )}
      </div>

      <div className="space-y-10">
        {data.results.map((result) => (
          <ScannerSection key={result.auditor_name} result={result} />
        ))}
      </div>
    </div>
  );
}
