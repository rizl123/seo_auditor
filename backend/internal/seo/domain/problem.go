package domain

import (
	"net/url"
)

type Problem struct {
	I18nNamespace   string
	DescriptionVars map[string]any
	Resources       []Resource
}

type Resource struct {
	Title string
	URL   *url.URL
}

func NewRes(title string, urlStr string) Resource {
	parsed, _ := url.Parse(urlStr)
	return Resource{Title: title, URL: parsed}
}
