import { useEffect, useState } from "react";

export type OgStatus = "ok" | "sub" | "miss";

export function useOgImage(url?: string) {
  const [data, setData] = useState<{ status: OgStatus; w: number; h: number }>({
    status: "miss",
    w: 0,
    h: 0,
  });

  useEffect(() => {
    if (!url) return;
    const img = new window.Image();
    img.src = url;
    img.onload = () => {
      const { naturalWidth: w, naturalHeight: h } = img;
      setData({ status: w / h >= 1.5 ? "ok" : "sub", w, h });
    };
  }, [url]);

  return data;
}
