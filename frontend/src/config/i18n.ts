export type Locale = (typeof locales)[number]["short"];
export const locales = [
  {
    short: "en" as const,
    full: "English" as const,
  },
  {
    short: "ru" as const,
    full: "Russian" as const,
  },
];
