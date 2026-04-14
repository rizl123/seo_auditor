package delivery

import (
	"backend/internal/seo/domain"
	"time"
)

type AggregatedReportDTO struct {
	URL     string          `json:"url"`
	Results []ScanResultDTO `json:"results"`
}

type ScanResultDTO struct {
	AuditorName string         `json:"auditor_name"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Details     map[string]any `json:"details,omitempty"`
	Problems    []ProblemDTO   `json:"problems"`
	IsCached    bool           `json:"is_cached"`
	ScannedAt   time.Time      `json:"scanned_at"`
}

type ProblemDTO struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Solutions   []string      `json:"solutions"`
	Resources   []ResourceDTO `json:"resources"`
}

type ResourceDTO struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func ToAggregatedReportDTO(report *domain.AggregatedReport) *AggregatedReportDTO {
	if report == nil {
		return nil
	}

	dto := &AggregatedReportDTO{
		URL:     report.URL.String(),
		Results: make([]ScanResultDTO, len(report.Results)),
	}

	for i, r := range report.Results {
		dto.Results[i] = toScanResultDTO(r)
	}

	return dto
}

func toScanResultDTO(r domain.ScanResult) ScanResultDTO {
	dto := ScanResultDTO{
		AuditorName: r.AuditorName,
		Name:        r.Name,
		Description: r.Description,
		Details:     r.Details,
		IsCached:    r.IsCached,
		ScannedAt:   r.ScannedAt,
		Problems:    make([]ProblemDTO, len(r.Problems)),
	}

	for i, p := range r.Problems {
		dto.Problems[i] = toProblemDTO(p)
	}

	return dto
}

func toProblemDTO(p domain.Problem) ProblemDTO {
	resources := make([]ResourceDTO, len(p.Resources))
	for i, r := range p.Resources {
		resources[i] = ResourceDTO{Title: r.Title, URL: r.URL}
	}

	solutions := p.Solutions
	if solutions == nil {
		solutions = []string{}
	}

	return ProblemDTO{
		Name:        p.Name,
		Description: p.Description,
		Solutions:   solutions,
		Resources:   resources,
	}
}

type PageReportDTO struct {
	URL       string          `json:"url"`
	Status    int             `json:"status"`
	IsCached  bool            `json:"is_cached"`
	ScannedAt time.Time       `json:"scanned_at"`
	Metadata  *MetadataDTO    `json:"metadata,omitempty"`
	Network   *NetworkInfoDTO `json:"network,omitempty"`
}

type MetadataDTO struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	H1          []string `json:"h1"`
	Canonical   string   `json:"canonical"`
	OgImage     string   `json:"og_image"`
}

type NetworkInfoDTO struct {
	ResponseTimeMs int64  `json:"response_time_ms"`
	Server         string `json:"server"`
	ContentType    string `json:"content_type"`
}
