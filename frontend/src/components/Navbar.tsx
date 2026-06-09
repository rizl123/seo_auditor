import { ShieldCheck } from "lucide-react";
import Link from "next/link";
import { AuthInfo } from "./AuthInfo";
import { LocaleSwitcher } from "./LocaleSwitcher";

export async function Navbar() {
  return (
    <nav className="sticky top-0 z-50 w-full border-b border-zinc-200 bg-white/80 backdrop-blur-md dark:border-zinc-800 dark:bg-zinc-950/80">
      <div className="mx-auto flex h-16 max-w-5xl flex-wrap items-center justify-between gap-y-2 px-3 sm:px-6">
        <Link
          href="/"
          className="flex items-center gap-2 sm:gap-2.5 transition-all active:scale-95"
        >
          <div className="flex h-8 w-8 sm:h-9 sm:w-9 items-center justify-center rounded-xl bg-blue-600 dark:bg-blue-500 shadow-sm shadow-blue-500/20">
            <ShieldCheck className="text-white" size={18} />
          </div>

          <span className="text-base sm:text-lg font-black tracking-tighter text-zinc-900 dark:text-zinc-50">
            SEO Auditor
          </span>
        </Link>

        <div className="flex items-center gap-1.5 sm:gap-3 ml-auto sm:ml-0">
          <LocaleSwitcher />

          <div className="hidden sm:block mx-1 h-4 w-px bg-zinc-200 dark:bg-zinc-800" />

          <AuthInfo />
        </div>
      </div>
    </nav>
  );
}
