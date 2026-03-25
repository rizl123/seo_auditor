package infrastructure

import (
	"backend/internal/domain"
	"net/http"

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

	data := &domain.SeoData{URL: urlStr}
	data.Title = doc.Find("title").Text()
	doc.Find("h1").Each(func(i int, s *goquery.Selection) {
		data.H1 = append(data.H1, s.Text())
	})

	return data, nil
}
