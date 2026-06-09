import { LogOut, User } from "lucide-react";
import { getTranslations } from "next-intl/server";
import { getSession } from "@/lib/auth/session";

export async function AuthInfo() {
  const t = await getTranslations("Navbar");
  const user = await getSession();

  if (!user) {
    return (
      <a
        href="/login"
        className="flex items-center gap-2 rounded-lg sm:rounded-xl bg-zinc-900 px-2 sm:px-6 py-2 text-white dark:bg-zinc-50 dark:text-zinc-900 active:scale-95"
      >
        <User size={16} />
        <span className="hidden sm:inline font-bold text-sm">
          {t("signIn")}
        </span>
      </a>
    );
  }

  return (
    <div className="flex items-center gap-1 sm:gap-2">
      <span className="hidden sm:inline font-bold text-sm text-zinc-700 dark:text-zinc-300 px-2 truncate max-w-22.5">
        {user.username || "User"}
      </span>

      <a
        href="/logout"
        className="flex items-center gap-1 sm:gap-2 rounded-lg sm:rounded-xl px-2 sm:px-4 py-2 text-rose-500 hover:bg-rose-500/10 active:scale-95"
      >
        <LogOut size={16} />
        <span className="hidden sm:inline font-bold text-sm">{t("exit")}</span>
      </a>
    </div>
  );
}
