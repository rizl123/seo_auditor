package domain

type SeoData struct {
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	H1          []string `json:"h1"`
	Canonical   string   `json:"canonical"`
	OgImage     string   `json:"og_image"`
	Status      int      `json:"status"`
}

type SeoRepository interface {
	FetchSeoData(urlStr string) (*SeoData, error)
}
