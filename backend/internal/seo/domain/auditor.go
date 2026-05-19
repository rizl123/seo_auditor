package domain

import "context"

type Auditor interface {
	AuditorName() string
	Analyze(ctx context.Context, raw *RawData) (*AuditResult, error)
}
