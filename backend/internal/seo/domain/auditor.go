package domain

import (
	"context"
	"net/url"
)

type Auditor interface {
	AuditorName() string
	Analyze(ctx context.Context, url *url.URL) *AuditResult
	I18nNamespace() string
}
