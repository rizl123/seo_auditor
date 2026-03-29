import { Globe, Loader2 } from "lucide-react";

interface SearchFormProps {
  onAnalyze: (url: string) => void;
  loading: boolean;
}

export function SearchForm({ onAnalyze, loading }: SearchFormProps) {
  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const url = formData.get("url") as string;
    if (url) {
      onAnalyze(url);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="relative mb-12">
      <div className="relative group">
        <Globe
          className="absolute left-4 top-1/2 -translate-y-1/2 text-zinc-400 group-focus-within:text-blue-500 transition-colors"
          size={20}
        />
        <input
          name="url"
          type="url"
          required
          placeholder="https://example.com"
          className="w-full pl-12 pr-32 py-4 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl shadow-sm focus:ring-2 focus:ring-blue-500 outline-none transition-all"
        />
        <button
          type="submit"
          disabled={loading}
          className="absolute right-2 top-2 bottom-2 px-6 bg-zinc-900 dark:bg-white text-white dark:text-black rounded-xl font-medium hover:opacity-90 disabled:opacity-50 transition-all flex items-center gap-2"
        >
          {loading ? <Loader2 className="animate-spin" size={18} /> : "Analyze"}
        </button>
      </div>
    </form>
  );
}
