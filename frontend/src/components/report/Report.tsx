"use client";

import { Globe } from "lucide-react";
import type { PageReport } from "@/types/report";
import { ScannerSection } from "./ScannerSection";

export function Report({ data }: { data: PageReport }) {
  return (
    <div className="max-w-3xl mx-auto space-y-12 py-10 animate-in fade-in slide-in-from-bottom-4 duration-700">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-6 px-2">
        <div className="min-w-0">
          <p className="text-xs font-bold text-zinc-400 uppercase tracking-widest mb-1">
            Target URL
          </p>
          <div className="flex items-center gap-2 text-xl font-bold break-all text-zinc-900 dark:text-white">
            <Globe size={20} className="text-accent shrink-0" />
            {data.url}
          </div>
        </div>
      </div>

      <div className="space-y-10">
        {data.results.map((result) => (
          <ScannerSection key={result.auditor_name} result={result} />
        ))}
      </div>
    </div>
  );
}
