package infrastructure

import (
	"backend/internal/domain"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"syscall"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type HttpSeoRepository struct {
	client *http.Client
}

func NewHttpSeoRepository() *HttpSeoRepository {
	return &HttpSeoRepository{
		client: createSafeClient(),
	}
}

func createSafeClient() *http.Client {
	return &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
				Control: func(network, address string, c syscall.RawConn) error {
					host, _, err := net.SplitHostPort(address)
					if err != nil {
						return err
					}
					ip := net.ParseIP(host)
					if ip != nil && (ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast()) {
						return errors.New("internal network access denied")
					}
					return nil
				},
			}).DialContext,
		},
	}
}

func (r *HttpSeoRepository) FetchSeoData(urlStr string) (*domain.SeoData, error) {
	if r.client == nil {
		r.client = createSafeClient()
	}

	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		return nil, errors.New("invalid protocol")
	}

	req, err := http.NewRequestWithContext(context.Background(), "GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; SeoBot/1.0)")

	res, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return &domain.SeoData{URL: urlStr, Status: res.StatusCode}, nil
	}

	doc, err := goquery.NewDocumentFromReader(io.LimitReader(res.Body, 2*1024*1024))
	if err != nil {
		return nil, err
	}

	data := &domain.SeoData{
		URL:    urlStr,
		H1:     make([]string, 0),
		Status: res.StatusCode,
	}

	data.Title = strings.TrimSpace(doc.Find("title").First().Text())
	data.Description, _ = doc.Find("meta[name='description']").Attr("content")
	data.Canonical, _ = doc.Find("link[rel='canonical']").Attr("href")
	data.OgImage, _ = doc.Find("meta[property='og:image']").Attr("content")

	doc.Find("h1").Each(func(i int, s *goquery.Selection) {
		val := strings.TrimSpace(s.Text())
		if val != "" {
			data.H1 = append(data.H1, val)
		}
	})

	return data, nil
}
