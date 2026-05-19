package domain

import (
	"net/url"
	"time"
)

type AuditResult struct {
	AuditorName   string
	I18nNamespace string
	Details       []Detail
	Problems      []Problem
	IsCached      bool
	ScannedAt     time.Time
}

type PageReport struct {
	URL     *url.URL
	Results []AuditResult
}
