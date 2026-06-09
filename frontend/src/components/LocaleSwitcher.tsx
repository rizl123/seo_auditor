"use client";

import { ChevronDown } from "lucide-react";
import { useLocale } from "next-intl";
import { useState } from "react";
import { type Locale, locales } from "@/config/i18n";
import { usePathname, useRouter } from "@/i18n/navigation";

export function LocaleSwitcher() {
  const currentLocale = useLocale() as Locale;
  const router = useRouter();
  const pathname = usePathname();
  const [open, setOpen] = useState(false);

  const current = locales.find((l) => l.short === currentLocale);

  const changeLocale = (locale: Locale) => {
    setOpen(false);
    router.replace(pathname, { locale });
  };

  return (
    <div className="relative">
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        className="flex items-center gap-1 px-2 sm:px-3 py-2 rounded-lg sm:rounded-xl border border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-900 text-xs font-bold text-zinc-700 dark:text-zinc-300 active:scale-95 transition"
      >
        <span>{current?.short.toUpperCase()}</span>
        <ChevronDown size={14} />
      </button>

      {open && (
        <div className="absolute right-0 mt-2 min-w-max rounded-xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 shadow-lg overflow-hidden z-50">
          {locales.map((l) => (
            <button
              key={l.short}
              type="button"
              onClick={() => changeLocale(l.short as Locale)}
              className="w-full text-left px-3 py-2 text-xs font-bold hover:bg-zinc-100 dark:hover:bg-zinc-800 text-zinc-700 dark:text-zinc-300"
            >
              {l.short.toUpperCase()}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
