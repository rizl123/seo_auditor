import { ShieldCheck } from "lucide-react";
import Link from "next/link";
import { AuthInfo } from "./AuthInfo";
import { LocaleSwitcher } from "./LocaleSwitcher";

export async function Navbar() {
  return (
    <nav className="sticky top-0 z-50 w-full border-b border-border-custom bg-background/80 backdrop-blur-md">
      <div className="max-w-5xl mx-auto px-6 h-16 flex items-center justify-between">
        <Link
          href="/"
          className="flex items-center gap-2.5 group transition-all active:scale-95"
        >
          <div className="w-9 h-9 bg-accent rounded-xl flex items-center justify-center shadow-sm shadow-accent/20">
            <ShieldCheck className="text-white" size={20} />
          </div>
          <span className="font-black text-xl tracking-tighter text-foreground">
            SEO Auditor
          </span>
        </Link>

        <div className="flex items-center gap-3">
          <LocaleSwitcher />
          <div className="h-4 w-px bg-border-custom mx-1" />
          <AuthInfo />
        </div>
      </div>
    </nav>
  );
}
