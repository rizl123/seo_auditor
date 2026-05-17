import Image from "next/image";
import { useState } from "react";

export function ImageLoader(props: Parameters<typeof Image>[0]) {
  const [isLoading, setIsLoading] = useState(true);
  return (
    <div className="relative border border-border-custom rounded-2xl overflow-hidden bg-zinc-100 dark:bg-zinc-900 inline-block shadow-sm min-w-75 min-h-37.5">
      {isLoading && (
        <div className="absolute inset-0 z-10 animate-pulse bg-linear-to-r from-zinc-200 via-zinc-300 to-zinc-200 dark:from-zinc-800 dark:via-zinc-700 dark:to-zinc-800" />
      )}

      <Image
        width={600}
        height={315}
        {...props}
        unoptimized
        className={`max-w-full h-auto max-h-64 object-contain transition-all duration-500 ${
          isLoading
            ? "opacity-0 scale-98 blur-sm"
            : "opacity-100 scale-100 blur-0"
        }`}
        onLoad={() => setIsLoading(false)}
      />
    </div>
  );
}
