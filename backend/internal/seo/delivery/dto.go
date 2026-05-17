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
	AuditorName   string          `json:"auditor_name"`
	I18nNamespace string          `json:"i18n_namespace"`
	Details       []DetailItemDTO `json:"details,omitempty"`
	Problems      []ProblemDTO    `json:"problems"`
	IsCached      bool            `json:"is_cached"`
	ScannedAt     time.Time       `json:"scanned_at"`
}

type DetailItemDTO struct {
	I18nLabel string            `json:"i18n_label"`
	Value     any               `json:"value"`
	Type      domain.DetailType `json:"type"`
}

type ProblemDTO struct {
	I18nNamespace   string         `json:"i18n_namespace"`
	DescriptionVars map[string]any `json:"description_vars"`
	Resources       []ResourceDTO  `json:"resources"`
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
			I18nLabel: d.I18nLabel,
			Value:     d.Value,
			Type:      d.Type,
		}
	}

	problems := make([]ProblemDTO, len(r.Problems))
	for i, p := range r.Problems {
		problems[i] = toProblemDTO(p)
	}

	return ScanResultDTO{
		AuditorName:   r.AuditorName,
		I18nNamespace: r.I18nNamespace,
		Details:       details,
		IsCached:      r.IsCached,
		ScannedAt:     r.ScannedAt,
		Problems:      problems,
	}
}

func toProblemDTO(p domain.Problem) ProblemDTO {
	resources := make([]ResourceDTO, len(p.Resources))
	for i, r := range p.Resources {
		resources[i] = ResourceDTO{Title: r.Title, URL: r.URL}
	}

	return ProblemDTO{
		I18nNamespace:   p.I18nNamespace,
		DescriptionVars: p.DescriptionVars,
		Resources:       resources,
	}
}
