package domain

import (
	"context"
	"net/url"
)

type Scanner interface {
	Scan(ctx context.Context, url *url.URL) (*PageReport, error)
}
