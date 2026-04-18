export function SpeedIndicator({ ms }: { ms: number }) {
  const isFast = ms < 800;
  const isAverage = ms >= 800 && ms < 1500;
  const percentage = Math.max(5, Math.min(100, 100 - (ms / 3000) * 100));

  return (
    <div className="flex flex-col items-end gap-1.5 min-w-25">
      <div className="h-1.5 w-full rounded-full bg-zinc-100 dark:bg-zinc-800 overflow-hidden">
        <div
          className={`h-full rounded-full transition-all duration-1000 ${
            isFast
              ? "bg-emerald-500"
              : isAverage
                ? "bg-amber-500"
                : "bg-rose-500"
          }`}
          style={{ width: `${percentage}%` }}
        />
      </div>
      <span
        className={`text-[9px] font-black uppercase tracking-tighter ${
          isFast
            ? "text-emerald-600"
            : isAverage
              ? "text-amber-600"
              : "text-rose-600"
        }`}
      >
        {isFast
          ? "Optimal"
          : isAverage
            ? "Needs Improvement"
            : "Critical Latency"}
      </span>
    </div>
  );
}
