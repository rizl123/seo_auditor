export interface ApiErrorItem {
  message: string;
  location?: string;
  value?: unknown;
}

export interface ApiErrorResponse {
  title: string;
  status: number;
  detail?: string;
  errors?: ApiErrorItem[];
}
