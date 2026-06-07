package domain

import (
	"net/url"
	"time"
)

type AuditFail struct {
	Title       string
	Description string
}

type AuditResult struct {
	AuditorName   string
	I18nNamespace string
	Details       []Detail
	Problems      []Problem
	IsCached      bool
	StartedAt     time.Time
	FinishedAt    time.Time
	Fail          *AuditFail
}

func NewAuditResult(s Auditor) *AuditResult {
	return &AuditResult{
		AuditorName:   s.AuditorName(),
		I18nNamespace: s.I18nNamespace(),
		Details:       []Detail{},
		Problems:      []Problem{},
		StartedAt:     time.Now(),
	}
}

type PageReport struct {
	URL     *url.URL
	Results []*AuditResult
}
