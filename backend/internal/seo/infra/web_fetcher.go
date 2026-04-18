package infra

import (
	"backend/internal/seo/domain"
	"context"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type WebFetcher struct {
	client *http.Client
}

func NewWebFetcher(client *http.Client) *WebFetcher {
	return &WebFetcher{client: client}
}

func (s *WebFetcher) Scan(ctx context.Context, url *neturl.URL) (*domain.PageReport, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	req.Header.Set("User-Agent", "SiteInspector/1.0")

	start := time.Now()
	res, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	}()

	report := &domain.PageReport{
		URL: url, Status: res.StatusCode, ScannedAt: time.Now(),
		Network: &domain.NetworkInfo{
			ResponseTime: time.Since(start),
			Server:       res.Header.Get("Server"),
			ContentType:  res.Header.Get("Content-Type"),
		},
	}

	if res.StatusCode == http.StatusOK {
		report.Metadata = s.parse(io.LimitReader(res.Body, 1024*512))
	}
	return report, nil
}

func (s *WebFetcher) parse(r io.Reader) *domain.Metadata {
	m := &domain.Metadata{H1: []string{}}
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return m
		}
		if tt == html.StartTagToken || tt == html.SelfClosingTagToken {
			s.processToken(z, m)
		}
	}
}

func (s *WebFetcher) processToken(z *html.Tokenizer, m *domain.Metadata) {
	t := z.Token()
	switch t.Data {
	case "title", "h1":
		tagName := t.Data
		if z.Next() == html.TextToken {
			val := strings.Join(strings.Fields(z.Token().Data), " ")
			if tagName == "title" {
				m.Title = val
			} else {
				m.H1 = append(m.H1, val)
			}
		}
	case "meta", "link":
		attrs := make(map[string]string)
		for _, a := range t.Attr {
			attrs[a.Key] = a.Val
		}
		s.fillMetadata(m, attrs)
	}
}

func (s *WebFetcher) fillMetadata(m *domain.Metadata, attr map[string]string) {
	switch {
	case attr["name"] == "description":
		m.Description = attr["content"]
	case attr["property"] == "og:image":
		m.OgImage = attr["content"]
	case attr["rel"] == "canonical":
		m.Canonical = attr["href"]
	}
}
