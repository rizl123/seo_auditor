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
  details: Record<string, any>;
  problems: Problem[];
  is_cached: boolean;
  scanned_at: string;
}

export interface PageReport {
  url: string;
  results: ScanResult[];
}
