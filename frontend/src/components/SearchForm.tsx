import { ArrowRight, Globe, Loader2 } from "lucide-react";

interface SearchFormProps {
  onAnalyze: (url: string) => void;
  loading: boolean;
}

export function SearchForm({ onAnalyze, loading }: SearchFormProps) {
  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const url = formData.get("url") as string;
    if (url) onAnalyze(url);
  };

  return (
    <div className="py-20 text-center">
      <h1 className="text-4xl font-bold mb-4 tracking-tight">SEO Analyzer</h1>
      <p className="text-zinc-500 mb-10 text-lg">
        Enter a URL to check your website's health.
      </p>

      <form onSubmit={handleSubmit} className="relative max-w-2xl mx-auto">
        <div className="flex items-center bg-card border border-border-custom rounded-2xl p-2 pl-5 transition-all focus-within:ring-4 focus-within:ring-accent/10 focus-within:border-accent">
          <Globe size={20} className="text-zinc-400 shrink-0" />
          <input
            name="url"
            type="url"
            required
            placeholder="https://yourwebsite.com"
            className="w-full px-4 py-3 bg-transparent outline-none text-base placeholder:text-zinc-400"
          />
          <button
            type="submit"
            disabled={loading}
            className="bg-accent text-white px-6 py-3 rounded-xl font-semibold flex items-center gap-2 hover:opacity-90 disabled:opacity-50 transition-all active:scale-95"
          >
            {loading ? (
              <Loader2 className="animate-spin size={20}" />
            ) : (
              "Analyze"
            )}
            {!loading && <ArrowRight size={18} />}
          </button>
        </div>
      </form>
    </div>
  );
}
