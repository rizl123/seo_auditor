"use client";

import { AlertCircle } from "lucide-react";
import { useState } from "react";
import { SearchForm } from "@/components/SearchForm";
import { SeoResults } from "@/components/SeoResults";
import type { SeoData } from "@/types/seo";

export default function Home() {
  const [result, setResult] = useState<SeoData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleAnalyze = async (url: string) => {
    setLoading(true);
    setError("");
    setResult(null);

    try {
      const res = await fetch(`/api/analyze?url=${encodeURIComponent(url)}`);
      if (!res.ok) {
        throw new Error("Failed to analyze URL");
      }
      const data = await res.json();
      setResult(data);
    } catch (_err) {
      setError("Error connecting to the analysis server.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 transition-colors">
      <main className="max-w-4xl mx-auto pt-20 px-6">
        <header className="mb-12 text-center">
          <h1 className="text-4xl font-bold tracking-tight text-zinc-900 dark:text-white mb-4">
            SEO Analyzer
          </h1>
          <p className="text-zinc-600 dark:text-zinc-400">
            Enter a URL to audit its on-page SEO elements.
          </p>
        </header>

        <SearchForm onAnalyze={handleAnalyze} loading={loading} />

        {error && (
          <div className="p-4 mb-8 bg-red-50 border border-red-100 text-red-600 rounded-xl flex items-center gap-3">
            <AlertCircle size={20} /> {error}
          </div>
        )}

        {result && <SeoResults data={result} />}
      </main>
    </div>
  );
}
