package domain

import (
	"context"
	"net/url"
)

type Runner interface {
	Run(ctx context.Context, url *url.URL) (*PageReport, error)
}
