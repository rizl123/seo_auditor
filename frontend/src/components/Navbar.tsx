"use server";

import { LogOut, ShieldCheck, User } from "lucide-react";
import Link from "next/link";
import { getSession } from "@/lib/auth";

export async function Navbar() {
  const user = await getSession();
  const isAuthenticated = !!user;

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
          {isAuthenticated ? (
            <div className="flex items-center gap-2 animate-in fade-in slide-in-from-right-4">
              <span className="hidden sm:inline font-bold text-sm text-foreground/80 px-2">
                {user.username || "User"}
              </span>

              <a
                href="/logout"
                className="flex items-center gap-2 px-4 py-2 text-sm font-bold text-rose-500 hover:bg-rose-500/10 rounded-xl transition-all active:scale-95"
              >
                <LogOut size={16} />
                <span className="hidden xs:inline">Exit</span>
              </a>
            </div>
          ) : (
            <a
              href="/login"
              className="flex items-center gap-2 px-6 py-2.5 text-sm font-bold bg-foreground text-background hover:opacity-90 rounded-xl transition-all shadow-sm active:scale-95"
            >
              <User size={16} />
              Sign In
            </a>
          )}
        </div>
      </div>
    </nav>
  );
}
