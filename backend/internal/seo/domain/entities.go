package domain

import (
	"net/url"
	"time"
)

type DetailType string

const (
	DetailTypeText     DetailType = "text"
	DetailTypeNumber   DetailType = "number"
	DetailTypeDuration DetailType = "duration_ms"
	DetailTypeURL      DetailType = "url"
	DetailTypeImage    DetailType = "image"
	DetailTypeBadge    DetailType = "badge"
)

type Detail struct {
	Label string
	Value any
	Type  DetailType
}

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

type Problem struct {
	Name        string
	Description string
	Solutions   []string
	Resources   []Resource
}

type Resource struct {
	Title string
	URL   string
}

type ScanResult struct {
	AuditorName string
	Name        string
	Description string
	Details     []Detail
	Problems    []Problem
	IsCached    bool
	ScannedAt   time.Time
}

type AggregatedReport struct {
	URL     *url.URL
	Results []ScanResult
}
