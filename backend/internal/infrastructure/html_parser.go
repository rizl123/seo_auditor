package infrastructure

import (
	"backend/internal/domain"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type HttpSeoRepository struct{}

func (r *HttpSeoRepository) FetchSeoData(urlStr string) (*domain.SeoData, error) {
	res, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	data := &domain.SeoData{
		URL:    urlStr,
		H1:     make([]string, 0),
		Status: res.StatusCode,
	}

	data.Title = doc.Find("title").First().Text()
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
