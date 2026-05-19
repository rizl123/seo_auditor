package domain

import (
	"context"
	"net/url"
)

type Fetcher interface {
	Scan(ctx context.Context, url *url.URL) (*RawData, error)
}
