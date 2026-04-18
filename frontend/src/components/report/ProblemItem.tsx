import { AlertTriangle, ExternalLink, Lightbulb } from "lucide-react";
import type { Problem } from "@/types/report";

export function ProblemItem({ problem }: { problem: Problem }) {
  return (
    <div className="p-6">
      <div className="flex items-start gap-4">
        <div className="mt-1 p-1.5 bg-rose-50 text-rose-600 dark:bg-rose-500/10 dark:text-rose-400 rounded-lg shrink-0">
          <AlertTriangle size={16} />
        </div>
        <div className="flex-1 space-y-4">
          <div>
            <h4 className="font-bold text-zinc-900 dark:text-zinc-100">
              {problem.name}
            </h4>
            <p className="text-sm text-zinc-500 mt-1 leading-relaxed">
              {problem.description}
            </p>
          </div>
          <div className="grid md:grid-cols-2 gap-6">
            {problem.solutions?.length > 0 && (
              <div className="space-y-2">
                <p className="text-[10px] font-bold text-zinc-400 uppercase flex items-center gap-1.5">
                  <Lightbulb size={12} className="text-amber-500" /> How to fix
                </p>
                <ul className="space-y-1.5">
                  {problem.solutions.map((s) => (
                    <li
                      key={s}
                      className="text-xs text-zinc-600 dark:text-zinc-400 flex items-start gap-2"
                    >
                      <span className="w-1 h-1 rounded-full bg-accent mt-1.5 shrink-0" />
                      {s}
                    </li>
                  ))}
                </ul>
              </div>
            )}
            {problem.resources?.length > 0 && (
              <div className="space-y-2">
                <p className="text-[10px] font-bold text-zinc-400 uppercase tracking-widest">
                  Docs
                </p>
                <div className="flex flex-wrap gap-2">
                  {problem.resources.map((res) => (
                    <a
                      key={res.title}
                      href={res.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center gap-1.5 px-3 py-1.5 bg-zinc-100 dark:bg-zinc-800 hover:bg-accent hover:text-white rounded-lg text-xs font-medium transition-all"
                    >
                      {res.title}
                      <ExternalLink size={10} className="opacity-50" />
                    </a>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
