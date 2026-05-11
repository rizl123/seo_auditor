"use client";

import { AlertCircle } from "lucide-react";
import { useState } from "react";
import { Report } from "@/components/report/Report";
import { SearchForm } from "@/components/SearchForm";
import { type ScanResponse, scanURL } from "./actions";

export function MainClientContainer() {
  const [response, setResponse] = useState<ScanResponse | null>(null);
  const [loading, setLoading] = useState(false);

  const handleAnalyze = async (url: string) => {
    setLoading(true);
    setResponse(null);

    try {
      const response = await scanURL(url);
      setResponse(response);
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <SearchForm onAnalyze={handleAnalyze} loading={loading} />

      {loading && <Skeletons />}
      {response && <ResponseComponent response={response} />}
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

function ResponseComponent({ response }: { response: ScanResponse }) {
  if (response.success) {
    return <Report data={response.data} />;
  }

  return (
    <div className="max-w-2xl mx-auto p-4 mb-8 bg-rose-50 border border-rose-100 text-rose-600 rounded-2xl flex items-start gap-3 animate-in fade-in zoom-in duration-300">
      <AlertCircle size={20} className="mt-0.5 shrink-0" />
      <div className="flex flex-col gap-1">
        <p className="font-semibold text-sm">{response.detail}</p>
        <ol>
          {response.errors?.map((error) => (
            <li key={error.message} className="text-sm opacity-90">
              {error.message}
            </li>
          ))}
        </ol>
      </div>
    </div>
  );
}
