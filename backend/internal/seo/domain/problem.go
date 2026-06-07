package domain

import (
	"net/url"
)

type Resource struct {
	Title string
	URL   *url.URL
}

func NewRes(title string, url *url.URL) Resource {
	return Resource{Title: title, URL: url}
}

type Problem struct {
	I18nNamespace   string
	DescriptionVars map[string]any
	Resources       []Resource
}

func NewProblem(ns string) Problem {
	return Problem{
		I18nNamespace:   ns,
		DescriptionVars: make(map[string]any),
		Resources:       make([]Resource, 0),
	}
}

func (p *Problem) AddStringVar(key string, val string) {
	p.DescriptionVars[key] = val
}

func (p *Problem) AddIntVar(key string, val int) {
	p.DescriptionVars[key] = val
}

func (p *Problem) AddInt64Var(key string, val int64) {
	p.DescriptionVars[key] = val
}

func (p *Problem) AddResource(title string, urlStr string) {
	parsed, err := url.Parse(urlStr)
	if err == nil {
		p.Resources = append(p.Resources, NewRes(title, parsed))
	}
}
