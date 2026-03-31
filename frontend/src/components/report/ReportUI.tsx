import { HelpCircle } from "lucide-react";
import type { ReactNode } from "react";

export const Card = ({
  children,
  className = "",
}: {
  children: ReactNode;
  className?: string;
}) => (
  <div
    className={`bg-card border border-border-custom rounded-3xl shadow-sm overflow-hidden ${className}`}
  >
    {children}
  </div>
);

export const StatusBadge = ({
  ok,
  message,
}: {
  ok: boolean;
  message: string;
}) => (
  <div
    className={`flex items-center gap-1.5 px-2.5 py-1 rounded-full border ${
      ok
        ? "bg-emerald-50/50 border-emerald-100 text-emerald-700 dark:bg-emerald-500/5 dark:border-emerald-500/20 dark:text-emerald-400"
        : "bg-rose-50/50 border-rose-100 text-rose-700 dark:bg-rose-500/5 dark:border-rose-500/20 dark:text-rose-400"
    }`}
  >
    <div
      className={`w-1.5 h-1.5 rounded-full ${ok ? "bg-emerald-500" : "bg-rose-500"}`}
    />
    <span className="text-[10px] uppercase tracking-wider font-bold">
      {message}
    </span>
  </div>
);

interface SectionProps {
  title: string;
  ok: boolean;
  statusMessage: string;
  info: string;
  seoExplanation: string;
  children: ReactNode;
}

export const Section = ({
  title,
  ok,
  statusMessage,
  info,
  seoExplanation,
  children,
}: SectionProps) => (
  <div className="group space-y-4 pt-8 first:pt-0 border-t border-zinc-100 dark:border-zinc-800 first:border-0">
    <div className="flex items-start justify-between">
      <div className="space-y-1">
        <div className="flex items-center gap-3">
          <h3 className="text-[11px] font-black text-zinc-400 uppercase tracking-widest">
            {title}
          </h3>
          <StatusBadge ok={ok} message={statusMessage} />
        </div>
        <p className="text-xs text-zinc-400 dark:text-zinc-500 leading-relaxed max-w-xl italic">
          {seoExplanation}
        </p>
      </div>
      <a
        href={info}
        target="_blank"
        rel="noopener noreferrer"
        className="p-2 text-zinc-300 hover:text-accent transition-colors"
      >
        <HelpCircle size={18} />
      </a>
    </div>
    <div className="relative">
      <div
        className={`absolute -left-4 top-0 bottom-0 w-1 rounded-full ${ok ? "bg-emerald-100 dark:bg-emerald-900/30" : "bg-rose-100 dark:bg-rose-900/30"}`}
      />
      <div className="pl-2">{children}</div>
    </div>
  </div>
);

export const SpeedIndicator = ({ ms }: { ms: number }) => {
  const isFast = ms < 200;
  const isAverage = ms >= 200 && ms < 500;
  return (
    <div className="flex flex-col items-end gap-1.5 min-w-20">
      <div className="h-1 w-full rounded-full bg-zinc-100 dark:bg-zinc-800 overflow-hidden">
        <div
          className={`h-full rounded-full transition-all duration-700 ${isFast ? "bg-emerald-500" : isAverage ? "bg-amber-500" : "bg-rose-500"}`}
          style={{ width: `${Math.min((ms / 1000) * 100, 100)}%` }}
        />
      </div>
      <span
        className={`text-[9px] font-black uppercase tracking-tighter ${isFast ? "text-emerald-600" : isAverage ? "text-amber-600" : "text-rose-600"}`}
      >
        {isFast ? "Excellent" : isAverage ? "Average" : "Slow"}
      </span>
    </div>
  );
};
