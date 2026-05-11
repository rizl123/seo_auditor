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
        className="flex items-center gap-2 px-6 py-2.5 text-sm font-bold bg-foreground text-background hover:opacity-90 rounded-xl transition-all shadow-sm active:scale-95"
      >
        <User size={16} />
        {t("signIn")}
      </a>
    );
  }

  return (
    <div className="flex items-center gap-2 animate-in fade-in slide-in-from-right-4">
      <span className="hidden sm:inline font-bold text-sm text-foreground/80 px-2">
        {user.username || "User"}
      </span>

      <a
        href="/logout"
        className="flex items-center gap-2 px-4 py-2 text-sm font-bold text-rose-500 hover:bg-rose-500/10 rounded-xl transition-all active:scale-95"
      >
        <LogOut size={16} />
        <span className="hidden xs:inline">{t("exit")}</span>
      </a>
    </div>
  );
}
