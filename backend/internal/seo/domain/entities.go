package domain

import (
	"net/url"
	"time"
)

type PageReport struct {
	URL       *url.URL
	Status    int
	IsCached  bool
	ScannedAt time.Time
	Metadata  *Metadata
	Network   *NetworkInfo
}

type Metadata struct {
	Title       string
	Description string
	H1          []string
	Canonical   string
	OgImage     string
}

type NetworkInfo struct {
	ResponseTime time.Duration
	Server       string
	ContentType  string
}
