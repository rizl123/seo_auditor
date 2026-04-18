import { Clock } from "lucide-react";
import Image from "next/image";
import type { ScanResult } from "@/types/report";
import { Card } from "./Card";
import { DetailItem } from "./DetailItem";
import { ProblemItem } from "./ProblemItem";

interface ScannerSectionProps {
  result: ScanResult;
}

export function ScannerSection({ result }: ScannerSectionProps) {
  const hasProblems = result.problems && result.problems.length > 0;
  const imageDetails = result.details.filter((d) => d.type === "image");
  const regularDetails = result.details.filter((d) => d.type !== "image");

  return (
    <div className="space-y-4">
      <div className="flex items-end justify-between px-2">
        <div className="space-y-1">
          <div className="flex items-center gap-2">
            <h3 className="text-lg font-black uppercase tracking-tight text-zinc-800 dark:text-zinc-200">
              {result.name}
            </h3>
            {hasProblems ? (
              <span className="text-[10px] bg-rose-500 text-white px-2 py-0.5 rounded-full font-bold">
                {result.problems.length} ISSUES
              </span>
            ) : (
              <span className="text-[10px] bg-emerald-500 text-white px-2 py-0.5 rounded-full font-bold">
                PASSED
              </span>
            )}
          </div>
          <p className="text-xs text-zinc-500 max-w-xl">{result.description}</p>
        </div>
        <div className="hidden sm:flex items-center gap-1.5 text-[10px] text-zinc-400 font-mono bg-zinc-50 dark:bg-zinc-900 px-2 py-1 rounded-md">
          <Clock size={12} />
          {new Date(result.scanned_at).toLocaleTimeString()}
        </div>
      </div>

      <Card className="divide-y divide-zinc-100 dark:divide-zinc-800 border-t-2 border-t-zinc-200 dark:border-t-zinc-700">
        {regularDetails.length > 0 && (
          <div className="p-6 grid grid-cols-1 sm:grid-cols-2 gap-x-8 gap-y-4 bg-zinc-50/30 dark:bg-zinc-900/10">
            {regularDetails.map((detail) => (
              <DetailItem key={detail.label} item={detail} />
            ))}
          </div>
        )}

        {imageDetails.length > 0 && (
          <div className="p-6 bg-zinc-50/50 dark:bg-zinc-900/20 border-t border-zinc-100 dark:border-zinc-800">
            {imageDetails.map((detail) => (
              <div key={detail.label} className="space-y-3">
                <p className="text-[9px] font-bold text-zinc-400 uppercase tracking-tighter">
                  {detail.label}
                </p>
                <div className="relative border border-border-custom rounded-2xl overflow-hidden bg-white dark:bg-zinc-950 inline-block shadow-sm">
                  <Image
                    src={String(detail.value)}
                    alt={detail.label}
                    width={600}
                    height={315}
                    unoptimized
                    className="max-w-full h-auto max-h-64 object-contain"
                  />
                </div>
              </div>
            ))}
          </div>
        )}

        {hasProblems && (
          <div className="divide-y divide-zinc-50 dark:divide-zinc-800/50">
            {result.problems.map((problem) => (
              <ProblemItem key={problem.name} problem={problem} />
            ))}
          </div>
        )}
      </Card>
    </div>
  );
}
