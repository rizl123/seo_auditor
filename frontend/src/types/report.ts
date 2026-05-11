export type DetailItem =
  | { label: string; type: "text"; value: string }
  | { label: string; type: "number"; value: number }
  | { label: string; type: "duration_ms"; value: number }
  | { label: string; type: "url"; value: string }
  | { label: string; type: "image"; value: string }
  | { label: string; type: "badge"; value: string };

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
