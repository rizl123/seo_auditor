import { API_URL } from "@/config/urls";
import type { ApiErrorItem, ApiErrorResponse } from "@/types/api";
import type { PageReport } from "@/types/report";

export type ScanResponse =
  | { detail: string; errors?: ApiErrorItem[]; success: false }
  | { data: PageReport; success: true };

export async function scanURL(url: string): Promise<ScanResponse> {
  if (!url) {
    return {
      detail: "ScanErrors.urlRequired",
      success: false,
    };
  }

  try {
    const apiUrl = new URL(`/api/scan?url=${encodeURIComponent(url)}`, API_URL);

    const res = await fetch(apiUrl, {
      method: "GET",
      cache: "no-store",
    });

    const data = await res.json();

    if (!res.ok) {
      const apiError = data as ApiErrorResponse;
      return {
        detail: apiError.detail
          ? `API.${apiError.detail}`
          : "ScanErrors.defaultError",
        errors: apiError.errors,
        success: false,
      };
    }

    return {
      data: data as PageReport,
      success: true,
    };
  } catch {
    return {
      detail: "ScanErrors.connectionFailed",
      success: false,
    };
  }
}
