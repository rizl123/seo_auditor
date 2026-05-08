"use client";

import { ChevronDown, Languages } from "lucide-react";
import { useLocale } from "next-intl";
import { type Locale, locales } from "@/config/i18n";
import { usePathname, useRouter } from "@/i18n/navigation";

export function LocaleSwitcher() {
  const currentLocale = useLocale() as Locale;
  const router = useRouter();
  const pathname = usePathname();

  const handleLocaleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const nextLocale = e.target.value as Locale;
    router.replace(pathname, { locale: nextLocale });
  };

  return (
    <div className="relative group">
      <div className="absolute left-3 top-1/2 -translate-y-1/2 pointer-events-none text-accent group-hover:scale-110 transition-transform">
        <Languages size={15} />
      </div>

      <select
        value={currentLocale}
        onChange={handleLocaleChange}
        className="appearance-none bg-accent/5 border border-border-custom hover:border-accent/40 text-foreground text-xs font-bold pl-9 pr-8 py-2 rounded-xl transition-all cursor-pointer focus:outline-none focus:ring-2 focus:ring-accent/20"
      >
        {locales.map(({ short, full }) => (
          <option
            key={short}
            value={short}
            className="bg-background text-foreground"
          >
            {full}
          </option>
        ))}
      </select>

      <div className="absolute right-2 top-1/2 -translate-y-1/2 pointer-events-none text-foreground/40">
        <ChevronDown size={14} />
      </div>
    </div>
  );
}
