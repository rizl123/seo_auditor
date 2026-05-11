import { useTranslations } from "next-intl";
import type { DetailItem } from "@/types/report";

type DetailItemProps = {
  item: DetailItem;
};

export function DetailItemComponent({ item }: DetailItemProps) {
  const t = useTranslations("Report");

  const renderValue = () => {
    if (item.value === null || item.value === undefined || item.value === "") {
      return <span className="text-zinc-300 italic">{t("noData")}</span>;
    }

    switch (item.type) {
      case "badge":
        return (
          <span className="inline-block px-2 py-0.5 bg-zinc-100 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded text-xs font-bold text-zinc-600 dark:text-zinc-400">
            {item.value}
          </span>
        );
      case "url":
        return (
          <a
            href={item.value}
            target="_blank"
            rel="noopener noreferrer"
            className="text-accent hover:underline truncate block"
          >
            {item.value}
          </a>
        );
      case "duration_ms":
        return <span>{t("duration", { value: item.value })}</span>;
      default:
        return <span className="truncate block">{item.value}</span>;
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
