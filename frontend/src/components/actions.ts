"use server";

import { getTranslations } from "next-intl/server";
import { API_URL } from "@/config/urls";
import type { ApiErrorItem, ApiErrorResponse } from "@/types/api";
import type { PageReport } from "@/types/report";

export type ScanResponse =
  | { detail: string; errors?: ApiErrorItem[]; success: false }
  | { data: PageReport; success: true };

export async function scanURL(url: string): Promise<ScanResponse> {
  const t = await getTranslations("ScanErrors");

  if (!url) {
    return {
      detail: t("urlRequired"),
      success: false,
    };
  }

  try {
    const apiUrl = `${API_URL}/api/scan?url=${encodeURIComponent(url)}`;

    const res = await fetch(apiUrl, {
      method: "GET",
      cache: "no-store",
    });

    const data = await res.json();

    if (!res.ok) {
      const apiError = data as ApiErrorResponse;
      return {
        detail: apiError.detail || t("defaultError"),
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
      detail: t("connectionFailed"),
      success: false,
    };
  }
}
