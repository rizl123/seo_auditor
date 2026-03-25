package domain

type SeoData struct {
	URL   string   `json:"url"`
	Title string   `json:"title"`
	H1    []string `json:"h1"`
}

type SeoRepository interface {
	FetchSeoData(urlStr string) (*SeoData, error)
}
