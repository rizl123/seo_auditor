import { API_URL } from "@/config/env";
import type { ApiErrorResponse } from "@/types/api";
import type { PageReport } from "@/types/report";

type ScanReturnType = { error?: string; data?: PageReport };

export async function scanURL(url: string): Promise<ScanReturnType> {
  if (!url) return { error: "URL is required" };

  try {
    const apiUrl = `${API_URL}/api/scan?url=${encodeURIComponent(url)}`;
    const res = await fetch(apiUrl, { cache: "no-store" });
    const data = await res.json();

    if (!res.ok) {
      const apiError = data as ApiErrorResponse;
      return {
        error:
          apiError.errors?.[0]?.message || apiError.detail || "Scan failed",
      };
    }

    return { data: data as PageReport };
  } catch {
    return { error: "Failed to connect to the server" };
  }
}
