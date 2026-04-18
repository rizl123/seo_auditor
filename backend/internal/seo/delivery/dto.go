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
	AuditorName string          `json:"auditor_name"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Details     []DetailItemDTO `json:"details,omitempty"`
	Problems    []ProblemDTO    `json:"problems"`
	IsCached    bool            `json:"is_cached"`
	ScannedAt   time.Time       `json:"scanned_at"`
}

type DetailItemDTO struct {
	Label string            `json:"label"`
	Value any               `json:"value"`
	Type  domain.DetailType `json:"type"`
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
	details := make([]DetailItemDTO, len(r.Details))
	for i, d := range r.Details {
		details[i] = DetailItemDTO{
			Label: d.Label,
			Value: d.Value,
			Type:  d.Type,
		}
	}

	problems := make([]ProblemDTO, len(r.Problems))
	for i, p := range r.Problems {
		problems[i] = toProblemDTO(p)
	}

	return ScanResultDTO{
		AuditorName: r.AuditorName,
		Name:        r.Name,
		Description: r.Description,
		Details:     details,
		IsCached:    r.IsCached,
		ScannedAt:   r.ScannedAt,
		Problems:    problems,
	}
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
