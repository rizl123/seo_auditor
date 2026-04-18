export type DetailType =
  | "text"
  | "number"
  | "duration_ms"
  | "url"
  | "image"
  | "badge";

export interface DetailItem {
  label: string;
  value: unknown;
  type: DetailType;
}

export interface Resource {
  title: string;
  url: string;
}

export interface Problem {
  name: string;
  description: string;
  solutions: string[];
  resources: Resource[];
}

export interface ScanResult {
  auditor_name: string;
  name: string;
  description: string;
  details: DetailItem[];
  problems: Problem[];
  is_cached: boolean;
  scanned_at: string;
}

export interface PageReport {
  url: string;
  results: ScanResult[];
}
