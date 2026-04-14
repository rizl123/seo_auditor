package infrastructure

import (
	"backend/internal/seo/domain"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	neturl "net/url"
	"strings"
	"syscall"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type WebFetcher struct {
	client *http.Client
}

func NewWebFetcher(client *http.Client) *WebFetcher {
	return &WebFetcher{client: client}
}

func (s *WebFetcher) Scan(ctx context.Context, url *neturl.URL) (*domain.PageReport, error) {
	if !strings.HasPrefix(url.Scheme, "http") {
		return nil, errors.New("invalid protocol")
	}

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		slog.Error("infrastructure: failed to create http request", "url", neturl.QueryEscape(url.String()), "error", err)
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "SiteInspector/1.0 (Bot)")

	res, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("infrastructure: http request failed: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			slog.Error("infrastructure: failed to close response body", "error", err)
		}
	}()

	report := &domain.PageReport{
		URL:       url,
		Status:    res.StatusCode,
		ScannedAt: time.Now(),
		Network: &domain.NetworkInfo{
			ResponseTime: time.Since(start),
			Server:       res.Header.Get("Server"),
			ContentType:  res.Header.Get("Content-Type"),
		},
	}

	if res.StatusCode != http.StatusOK {
		return report, nil
	}

	doc, err := goquery.NewDocumentFromReader(io.LimitReader(res.Body, 2*1024*1024))
	if err != nil {
		return report, nil
	}

	report.Metadata = &domain.Metadata{
		Title:       strings.TrimSpace(doc.Find("title").First().Text()),
		Description: doc.Find("meta[name='description']").AttrOr("content", ""),
		Canonical:   doc.Find("link[rel='canonical']").AttrOr("href", ""),
		OgImage:     doc.Find("meta[property='og:image']").AttrOr("content", ""),
		H1:          []string{},
	}

	doc.Find("h1").Each(func(_ int, sel *goquery.Selection) {
		if val := strings.TrimSpace(sel.Text()); val != "" {
			report.Metadata.H1 = append(report.Metadata.H1, val)
		}
	})

	return report, nil
}

func CreateSecureClient() *http.Client {
	return &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
				Control: func(_, address string, _ syscall.RawConn) error {
					host, _, _ := net.SplitHostPort(address)
					ip := net.ParseIP(host)
					if ip != nil && (ip.IsLoopback() || ip.IsPrivate()) {
						return errors.New("internal network access denied")
					}
					return nil
				},
			}).DialContext,
		},
	}
}
