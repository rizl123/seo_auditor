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

	"github.com/PuerkitoBio/goquery"
)

const maxHTMLSize = 2 * 1024 * 1024

type WebFetcher struct {
	client *http.Client
}

func NewWebFetcher(client *http.Client) *WebFetcher {
	return &WebFetcher{client: client}
}

func (s *WebFetcher) Scan(ctx context.Context, url *neturl.URL) (*domain.RawData, error) {
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

	raw := &domain.RawData{
		URL:    url,
		Status: res.StatusCode,
		Network: &domain.NetworkInfo{
			ResponseTime: time.Since(start),
			Server:       res.Header.Get("Server"),
			ContentType:  res.Header.Get("Content-Type"),
		},
	}

	if res.StatusCode == http.StatusOK {
		meta, err := s.parse(res.Body)
		if err == nil {
			raw.Metadata = meta
		} else {
			raw.Metadata = &domain.Metadata{H1: []string{}}
		}
	}

	return raw, nil
}

func (s *WebFetcher) parse(r io.Reader) (*domain.Metadata, error) {
	limited := io.LimitReader(r, maxHTMLSize)

	doc, err := goquery.NewDocumentFromReader(limited)
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	doc.Find("script, style, template, noscript").Remove()
	doc.Find("[hidden], [style*='display:none'], [style*='visibility:hidden']").Remove()

	m := &domain.Metadata{
		H1: make([]string, 0, 4),
	}

	title := strings.TrimSpace(doc.Find("title").First().Text())
	if title != "" {
		m.Title = normalizeText(title)
	}

	if v, ok := doc.Find(`meta[name="description"]`).First().Attr("content"); ok {
		m.Description = v
	}

	if v, ok := doc.Find(`meta[property="og:image"]`).First().Attr("content"); ok {
		m.OgImage = v
	}

	doc.Find(`link[rel="canonical"]`).Each(func(_ int, sel *goquery.Selection) {
		if v, ok := sel.Attr("href"); ok {
			m.Canonical = v
		}
	})

	doc.Find("h1").Each(func(_ int, sel *goquery.Selection) {
		text := normalizeText(sel.Text())

		if len(text) == 0 || len(text) > 300 {
			return
		}

		m.H1 = append(m.H1, text)
	})

	return m, nil
}

func normalizeText(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
