import { AlertCircle, CheckCircle2 } from "lucide-react";
import type { SeoData } from "@/types/seo";

interface SeoResultsProps {
  data: SeoData;
}

export function SeoResults({ data }: SeoResultsProps) {
  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <StatusCard label="Status Code" value={data.status} />
        <div className="md:col-span-2 p-6 bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-100 dark:border-zinc-800 shadow-sm">
          <p className="text-sm text-zinc-500 mb-1">Title</p>
          <p className="text-lg font-semibold truncate">
            {data.title || "Missing Title"}
          </p>
        </div>
      </div>

      <div className="p-8 bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-100 dark:border-zinc-800 shadow-sm">
        <h3 className="text-lg font-bold mb-4 flex items-center gap-2 text-zinc-900 dark:text-zinc-50">
          <CheckCircle2 className="text-green-500" size={20} /> Meta Description
        </h3>
        <p className="text-zinc-600 dark:text-zinc-400 leading-relaxed mb-8">
          {data.description ||
            "No meta description found. This is bad for CTR."}
        </p>

        <h3 className="text-lg font-bold mb-4 text-zinc-900 dark:text-zinc-50">
          H1 Headers ({(data.h1 ?? []).length})
        </h3>
        <ul className="space-y-2">
          {(data.h1 ?? []).map((h) => (
            <li
              key={h}
              className="p-3 bg-zinc-50 dark:bg-zinc-800/50 rounded-lg border border-zinc-100 dark:border-zinc-800 text-sm text-zinc-700 dark:text-zinc-300"
            >
              {h}
            </li>
          ))}
          {(!data.h1 || data.h1.length === 0) && (
            <li className="flex items-center gap-2 text-red-500 font-medium">
              <AlertCircle size={16} /> Missing H1 header!
            </li>
          )}
        </ul>
      </div>
    </div>
  );
}

function StatusCard({ label, value }: { label: string; value: number }) {
  const isOk = value >= 200 && value < 300;
  return (
    <div className="p-6 bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-100 dark:border-zinc-800 shadow-sm">
      <p className="text-sm text-zinc-500 mb-1">{label}</p>
      <p
        className={`text-2xl font-bold ${isOk ? "text-green-500" : "text-amber-500"}`}
      >
        {value}
      </p>
    </div>
  );
}
