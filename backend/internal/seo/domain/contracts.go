package domain

import (
	"context"
	"net/url"
)

type Fetcher interface {
	Scan(ctx context.Context, url *url.URL) (*PageReport, error)
}

type Auditor interface {
	AuditorName() string
	Analyze(ctx context.Context, report *PageReport) (*ScanResult, error)
}

type Runner interface {
	Run(ctx context.Context, url *url.URL) (*AggregatedReport, error)
}
