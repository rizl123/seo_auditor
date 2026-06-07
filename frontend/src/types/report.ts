export type DetailItem =
  | { i18n_label: string; value: string; type: "text" }
  | { i18n_label: string; value: number; type: "number" }
  | { i18n_label: string; value: number; type: "duration_ms" }
  | { i18n_label: string; value: string; type: "url" }
  | { i18n_label: string; value: string; type: "image" }
  | { i18n_label: string; value: string; type: "badge" };

export interface Resource {
  title: string;
  url: string;
}

export interface Problem {
  i18n_namespace: string;
  description_vars: Record<string, string | number>;
  resources: Resource[];
}

export interface AuditFail {
  title: string;
  description: string;
}

export interface AuditResult {
  auditor_name: string;
  i18n_namespace: string;
  details?: DetailItem[];
  problems: Problem[];
  is_cached: boolean;
  started_at: string;
  finished_at: string;
  fail: AuditFail | null;
}

export interface PageReport {
  url: string;
  results: AuditResult[];
}
