import type { DetailItem as IDetailItem } from "@/types/report";

export function DetailItem({ item }: { item: IDetailItem }) {
  const renderValue = () => {
    if (item.value === null || item.value === undefined || item.value === "") {
      return <span className="text-zinc-300 italic">n/a</span>;
    }

    switch (item.type) {
      case "badge":
        return (
          <span className="inline-block px-2 py-0.5 bg-zinc-100 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded text-xs font-bold text-zinc-600 dark:text-zinc-400">
            {String(item.value)}
          </span>
        );
      case "url":
        return (
          <a
            href={String(item.value)}
            target="_blank"
            rel="noopener noreferrer"
            className="text-accent hover:underline truncate block"
          >
            {String(item.value)}
          </a>
        );
      case "duration_ms":
        return <span>{String(item.value)} ms</span>;
      default:
        return (
          <span className="truncate block">
            {typeof item.value === "object"
              ? JSON.stringify(item.value)
              : String(item.value)}
          </span>
        );
    }
  };

  return (
    <div className="overflow-hidden">
      <p className="text-xs font-bold text-zinc-400 uppercase tracking-tighter mb-0.5">
        {item.label}
      </p>
      <div className="text-sm font-semibold text-zinc-700 dark:text-zinc-300">
        {renderValue()}
      </div>
    </div>
  );
}
