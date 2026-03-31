"use client";

import { AlertCircle } from "lucide-react";
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
      if (!res.ok) {
        throw new Error("Failed to scan URL");
      }
      const data = await res.json();
      setResult(data);
    } catch (_err) {
      setError("Error connecting to the scanner server.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-background transition-colors">
      <main className="max-w-5xl mx-auto px-6 pb-20">
        <SearchForm onAnalyze={handleAnalyze} loading={loading} />

        {error && (
          <div className="p-4 mb-8 bg-red-50 border border-red-100 text-red-600 rounded-xl flex items-center gap-3">
            <AlertCircle size={20} /> {error}
          </div>
        )}

        {result && <Report data={result} />}
      </main>
    </div>
  );
}
