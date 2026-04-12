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
	"strings"
	"syscall"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type WebScanner struct {
	client *http.Client
}

func NewWebScanner(client *http.Client) *WebScanner {
	return &WebScanner{client: client}
}

func (s *WebScanner) Scan(ctx context.Context, urlStr string) (*domain.PageReport, error) {
	if !strings.HasPrefix(urlStr, "http") {
		return nil, errors.New("invalid protocol")
	}

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		slog.Error("infrastructure: failed to create http request", "url", urlStr, "error", err)
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
		URL:       urlStr,
		Status:    res.StatusCode,
		ScannedAt: time.Now(),
		Network: &domain.NetworkInfo{
			ResponseTimeMs: time.Since(start).Milliseconds(),
			Server:         res.Header.Get("Server"),
			ContentType:    res.Header.Get("Content-Type"),
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
