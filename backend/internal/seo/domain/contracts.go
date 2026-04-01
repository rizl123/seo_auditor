package domain

import "context"

type Scanner interface {
	Scan(ctx context.Context, urlStr string) (*PageReport, error)
}

type ReportRepo interface {
	Fetch(ctx context.Context, url string) (*PageReport, error)
	Store(ctx context.Context, url string, report *PageReport) error
}
