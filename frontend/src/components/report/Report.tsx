import { AlertCircle, Clock, ImageIcon } from "lucide-react";
import Image from "next/image";
import { useOgImage } from "@/hooks/useOgImage";
import type { PageReport } from "@/types/report";
import { Card, Section, SpeedIndicator } from "./ReportUI";

export function Report({ data }: { data: PageReport }) {
  const { metadata, network, scanned_at } = data;
  const og = useOgImage(metadata?.og_image);

  const scanTime = new Date(scanned_at).toLocaleTimeString("en-US", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });

  const h1Count = metadata?.h1?.length || 0;
  const isH1Ok = h1Count === 1;
  const h1Verdict =
    h1Count === 0 ? "Missing" : h1Count > 1 ? `${h1Count} tags` : "Perfect";

  const isTitleOk = (metadata?.title?.length || 0) >= 10;
  const isDescOk = (metadata?.description?.length || 0) > 50;

  return (
    <div className="max-w-3xl mx-auto space-y-8 py-10 animate-in fade-in duration-500">
      <div className="flex items-center justify-between px-2">
        <div className="space-y-1">
          <p className="text-[10px] font-bold text-zinc-400 uppercase tracking-widest">
            Analysis Snapshot
          </p>
          <div className="flex items-center gap-2 text-zinc-900 dark:text-white font-medium">
            <Clock size={14} className="text-zinc-400" />
            <span className="text-sm">{scanTime}</span>
          </div>
        </div>
        <Card className="flex items-center gap-6 px-6 py-3">
          <div className="text-right">
            <p className="text-[10px] font-bold text-zinc-400 uppercase">
              Response
            </p>
            <p className="text-xl font-black text-zinc-900 dark:text-white leading-none">
              {network?.response_time_ms}
              <span className="text-[10px] ml-0.5 text-zinc-400 font-normal">
                ms
              </span>
            </p>
          </div>
          <div className="h-8 w-px bg-zinc-100 dark:bg-zinc-800" />
          <SpeedIndicator ms={network?.response_time_ms || 0} />
        </Card>
      </div>

      <Card className="p-10 space-y-10">
        <Section
          title="Title Tag"
          ok={isTitleOk}
          statusMessage={isTitleOk ? "Optimized" : "Too Short"}
          info="https://mdn.io/title"
          seoExplanation="Primary heading for search results and browser tabs."
        >
          <div className="text-2xl font-bold text-zinc-900 dark:text-zinc-100 leading-tight tracking-tight">
            {metadata?.title || (
              <span className="text-rose-500 font-medium">Missing Title</span>
            )}
          </div>
        </Section>

        <Section
          title="Description"
          ok={isDescOk}
          statusMessage={
            !metadata?.description
              ? "Missing"
              : isDescOk
                ? "Good"
                : "Needs work"
          }
          info="https://mdn.io/meta-description"
          seoExplanation="Summary displayed in search results to influence click-through rates."
        >
          <div
            className={`text-lg leading-relaxed ${isDescOk ? "text-zinc-600 dark:text-zinc-400" : "text-rose-400"}`}
          >
            {metadata?.description || "No meta description found."}
          </div>
        </Section>

        <Section
          title="Heading Structure (H1)"
          ok={isH1Ok}
          statusMessage={h1Verdict}
          info="https://mdn.io/h1"
          seoExplanation="The main semantic topic of your page content."
        >
          <div className="flex flex-wrap gap-3">
            {metadata?.h1?.map((h) => (
              <div
                key={h}
                className="inline-flex items-center gap-3 px-4 py-3 bg-zinc-50 dark:bg-zinc-900/50 border border-zinc-100 dark:border-zinc-800 rounded-xl"
              >
                <span className="text-[10px] font-black text-zinc-300 dark:text-zinc-700 uppercase">
                  H1
                </span>
                <span className="font-bold text-zinc-800 dark:text-zinc-200">
                  {h}
                </span>
              </div>
            ))}
            {!h1Count && (
              <div className="flex items-center gap-2 text-rose-500 text-sm font-semibold">
                <AlertCircle size={16} /> Add an H1 tag to improve semantics.
              </div>
            )}
          </div>
        </Section>
      </Card>

      <Card className="p-10">
        <div className="flex items-center gap-3 mb-8">
          <div className="p-2.5 bg-zinc-900 dark:bg-zinc-100 text-white dark:text-zinc-900 rounded-xl">
            <ImageIcon size={18} />
          </div>
          <h3 className="font-black text-lg uppercase tracking-tight">
            Social Preview
          </h3>
        </div>
        {metadata?.og_image ? (
          <div className="group relative aspect-[1.91/1] rounded-2xl overflow-hidden border border-zinc-100 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-900">
            <Image
              src={metadata.og_image}
              alt="OG"
              fill
              unoptimized
              className="object-contain p-4 group-hover:scale-[1.01] transition-transform"
            />
            {og.status === "sub" && (
              <div className="absolute bottom-4 left-4 right-4 bg-amber-50/90 dark:bg-amber-950/90 backdrop-blur-sm p-3 rounded-lg border border-amber-200 dark:border-amber-900/50 text-[11px] font-bold text-amber-800 dark:text-amber-400 text-center">
                ⚠️ Image size is not optimal (1200x630px recommended)
              </div>
            )}
          </div>
        ) : (
          <div className="py-20 border-2 border-dashed border-zinc-100 dark:border-zinc-800 rounded-2xl text-center text-zinc-300 font-bold uppercase text-[10px] tracking-widest leading-loose">
            No Preview Image Found
            <br />
            <span className="text-[8px] font-normal lowercase tracking-normal">
              Add og:image tags for better sharing
            </span>
          </div>
        )}
      </Card>
    </div>
  );
}
