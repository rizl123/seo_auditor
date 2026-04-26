"use client";

import { AlertCircle } from "lucide-react";
import { useState } from "react";
import { scanURL } from "@/app/actions";
import { Report } from "@/components/report/Report";
import { SearchForm } from "@/components/SearchForm";
import type { PageReport } from "@/types/report";

export function MainClientContainer() {
  const [result, setResult] = useState<PageReport | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleAnalyze = async (url: string) => {
    setLoading(true);
    setError(null);
    setResult(null);

    const response = await scanURL(url);

    if (response.error) {
      setError(response.error);
    } else if (response.data) {
      setResult(response.data);
    }

    setLoading(false);
  };

  return (
    <>
      <SearchForm onAnalyze={handleAnalyze} loading={loading} />

      {error && (
        <div className="max-w-2xl mx-auto p-4 mb-8 bg-rose-50 border border-rose-100 text-rose-600 rounded-2xl flex items-start gap-3 animate-in fade-in zoom-in duration-300">
          <AlertCircle size={20} className="mt-0.5 shrink-0" />
          <div className="flex flex-col gap-1">
            <span className="font-semibold text-sm">Scan Error</span>
            <span className="text-sm opacity-90">{error}</span>
          </div>
        </div>
      )}

      {loading && <Skeletons />}
      {result && <Report data={result} />}
    </>
  );
}

function Skeletons() {
  return (
    <div className="max-w-3xl mx-auto space-y-8 py-10">
      {[1, 2].map((i) => (
        <div key={i} className="animate-pulse space-y-4">
          <div className="h-4 bg-zinc-200 dark:bg-zinc-800 w-1/4 rounded-full" />
          <div className="h-64 bg-zinc-100 dark:bg-zinc-900 rounded-3xl" />
        </div>
      ))}
    </div>
  );
}
