"use client";

import { AlertCircle, Globe, Search } from "lucide-react";
import { useState } from "react";
import { Report } from "@/components/report/Report";
import { SearchForm } from "@/components/SearchForm";
import type { PageReport } from "@/types/report";

export default function Home() {
  const [result, setResult] = useState<PageReport | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleAnalyze = async (url: string) => {
    setLoading(true);
    setError("");
    setResult(null);

    try {
      const res = await fetch(`/api/scan?url=${encodeURIComponent(url)}`);
      if (!res.ok) throw new Error("Failed to scan URL");
      const data = await res.json();
      setResult(data);
    } catch (_err) {
      setError("Error connecting to the scanner server. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-background transition-colors">
      <main className="max-w-5xl mx-auto px-6 pb-20">
        <SearchForm onAnalyze={handleAnalyze} loading={loading} />

        {error && (
          <div className="max-w-2xl mx-auto p-4 mb-8 bg-rose-50 border border-rose-100 text-rose-600 rounded-2xl flex items-center gap-3 animate-in fade-in zoom-in duration-300">
            <AlertCircle size={20} /> {error}
          </div>
        )}

        {loading && (
          <div className="max-w-3xl mx-auto space-y-8 py-10">
            {[1, 2].map((i) => (
              <div key={i} className="animate-pulse space-y-4">
                <div className="h-4 bg-zinc-200 dark:bg-zinc-800 w-1/4 rounded-full" />
                <div className="h-64 bg-zinc-100 dark:bg-zinc-900 rounded-3xl" />
              </div>
            ))}
          </div>
        )}

        {!loading && !result && !error && (
          <div className="py-20 text-center opacity-30">
            <div className="inline-flex p-8 rounded-full bg-zinc-100 dark:bg-zinc-900 mb-4">
              <Search size={40} className="text-zinc-400" />
            </div>
            <p className="text-sm font-medium tracking-wide uppercase">
              Ready for analysis
            </p>
          </div>
        )}

        {result && <Report data={result} />}
      </main>
    </div>
  );
}
