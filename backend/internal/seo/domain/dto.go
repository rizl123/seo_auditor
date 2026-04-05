package domain

import "time"

type PageReport struct {
	URL       string       `json:"url"`
	Status    int          `json:"status"`
	IsCached  bool         `json:"is_cached"`
	ScannedAt time.Time    `json:"scanned_at"`
	Metadata  *Metadata    `json:"metadata,omitempty"`
	Network   *NetworkInfo `json:"network,omitempty"`
}

type Metadata struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	H1          []string `json:"h1"`
	Canonical   string   `json:"canonical"`
	OgImage     string   `json:"og_image"`
}

type NetworkInfo struct {
	ResponseTimeMs int64  `json:"response_time_ms"`
	Server         string `json:"server"`
	ContentType    string `json:"content_type"`
}
