export interface SeoMetadata {
  title: string;
  description: string;
  h1: string[];
  canonical: string;
  og_image: string;
}

export interface NetworkInfo {
  content_type: string;
  response_time_ms: number;
  server: string;
}

export interface PageReport {
  url: string;
  status: number;
  scanned_at: string;
  metadata?: SeoMetadata;
  network?: NetworkInfo;
}
