"use client";

import {
  AlertTriangle,
  Clock,
  ExternalLink,
  Globe,
  Lightbulb,
} from "lucide-react";
import type { PageReport, Problem, ScanResult } from "@/types/report";
import { Card, SpeedIndicator } from "./ReportUI";

export function Report({ data }: { data: PageReport }) {
  const perfData = data.results.find((r) => r.auditor_name === "performance");

  return (
    <div className="max-w-3xl mx-auto space-y-12 py-10 animate-in fade-in slide-in-from-bottom-4 duration-700">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-6 px-2">
        <div className="min-w-0">
          <p className="text-[10px] font-bold text-zinc-400 uppercase tracking-widest mb-1">
            Target URL
          </p>
          <div className="flex items-center gap-2 text-xl font-bold break-all text-zinc-900 dark:text-white">
            <Globe size={20} className="text-accent shrink-0" />
            {data.url}
          </div>
        </div>

        {perfData && (
          <Card className="flex items-center gap-5 px-6 py-3 shrink-0">
            <div className="text-right">
              <p className="text-[10px] font-bold text-zinc-400 uppercase">
                Response
              </p>
              <p className="text-xl font-black">
                {perfData.details.response_time_ms}
                <span className="text-[10px] ml-0.5 font-normal text-zinc-400">
                  ms
                </span>
              </p>
            </div>
            <div className="h-8 w-px bg-zinc-100 dark:bg-zinc-800" />
            <SpeedIndicator ms={perfData.details.response_time_ms} />
          </Card>
        )}
      </div>

      <div className="space-y-10">
        {data.results.map((result) => (
          <ScannerSection key={result.auditor_name} result={result} />
        ))}
      </div>
    </div>
  );
}

function ScannerSection({ result }: { result: ScanResult }) {
  const hasProblems = result.problems && result.problems.length > 0;

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
        {/* Render Details */}
        <div className="p-6 grid grid-cols-1 sm:grid-cols-2 gap-x-8 gap-y-4 bg-zinc-50/30 dark:bg-zinc-900/10">
          {Object.entries(result.details).map(([key, value]) => (
            <div key={key} className="overflow-hidden">
              <p className="text-[9px] font-bold text-zinc-400 uppercase tracking-tighter mb-0.5">
                {key.replace(/_/g, " ")}
              </p>
              <p className="text-sm font-semibold text-zinc-700 dark:text-zinc-300 truncate">
                {value !== null && value !== undefined ? (
                  String(value)
                ) : (
                  <span className="text-zinc-300 italic">n/a</span>
                )}
              </p>
            </div>
          ))}
        </div>

        {/* Problems */}
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

function ProblemItem({ problem }: { problem: Problem }) {
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
