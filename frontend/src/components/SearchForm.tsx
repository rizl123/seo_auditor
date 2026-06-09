import { ArrowRight, Globe, Loader2 } from "lucide-react";
import { useTranslations } from "next-intl";
import { useState } from "react";

interface SearchFormProps {
  onAnalyze: (url: string) => void;
  loading: boolean;
}

export function SearchForm({ onAnalyze, loading }: SearchFormProps) {
  const t = useTranslations("Search");
  const [url, setUrl] = useState("");

  const handleAnalyze = () => {
    const trimmed = url.trim();
    if (!trimmed) return;
    onAnalyze(trimmed);
  };

  return (
    <div className="py-8 sm:py-14 lg:py-20 px-4">
      <div className="max-w-2xl mx-auto">
        <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-3 sm:p-4 transition-all focus-within:ring-4 focus-within:ring-blue-500/10 focus-within:border-blue-500 dark:focus-within:border-blue-500">
          <div className="flex items-center gap-3">
            <Globe className="text-zinc-400 shrink-0 size-5" />

            <input
              type="url"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter") {
                  e.preventDefault();
                  handleAnalyze();
                }
              }}
              placeholder={t("placeholder")}
              className="w-full bg-transparent outline-none text-sm sm:text-base placeholder:text-zinc-400 py-2 sm:py-3 text-zinc-900 dark:text-zinc-100"
            />
          </div>

          <button
            type="button"
            disabled={loading}
            onClick={handleAnalyze}
            className="mt-3 sm:mt-4 w-full sm:w-auto bg-blue-600 hover:bg-blue-700 text-white dark:bg-blue-500 dark:hover:bg-blue-600 px-4 sm:px-6 py-3 rounded-xl font-medium sm:font-semibold flex items-center justify-center gap-2 disabled:opacity-50 transition-all active:scale-95 min-h-11"
          >
            {loading ? (
              <Loader2 size={20} className="animate-spin" />
            ) : (
              <>
                <span>{t("button")}</span>
                <ArrowRight className="size-4" />
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );
}
