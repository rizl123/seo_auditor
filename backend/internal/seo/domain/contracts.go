package domain

import "context"

type Scanner interface {
	Scan(ctx context.Context, urlStr string) (*PageReport, error)
}
