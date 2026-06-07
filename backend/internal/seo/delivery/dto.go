package delivery

import (
	"backend/internal/seo/domain"
	"time"
)

type PageReportDTO struct {
	URL     string           `json:"url"`
	Results []AuditResultDTO `json:"results"`
}

type AuditResultDTO struct {
	AuditorName   string          `json:"auditor_name"`
	I18nNamespace string          `json:"i18n_namespace"`
	Details       []DetailItemDTO `json:"details,omitempty"`
	Problems      []ProblemDTO    `json:"problems"`
	IsCached      bool            `json:"is_cached"`
	StartedAt     time.Time       `json:"started_at"`
	FinishedAt    time.Time       `json:"finished_at"`
	Fail          *AuditFailDTO   `json:"fail"`
}

type AuditFailDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
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

func ToPageReportDTO(r *domain.PageReport) *PageReportDTO {
	if r == nil {
		return nil
	}

	dto := &PageReportDTO{
		URL:     r.URL.String(),
		Results: make([]AuditResultDTO, len(r.Results)),
	}

	for i, r := range r.Results {
		dto.Results[i] = ToAuditResultDTO(r)
	}

	return dto
}

func ToAuditResultDTO(a *domain.AuditResult) AuditResultDTO {
	dto := AuditResultDTO{
		AuditorName:   a.AuditorName,
		I18nNamespace: a.I18nNamespace,
		Details:       make([]DetailItemDTO, len(a.Details)),
		IsCached:      a.IsCached,
		StartedAt:     a.StartedAt,
		FinishedAt:    a.FinishedAt,
		Problems:      make([]ProblemDTO, len(a.Problems)),
	}

	for i, d := range a.Details {
		dto.Details[i] = ToDetailItemDTO(d)
	}

	for i, p := range a.Problems {
		dto.Problems[i] = ToProblemDTO(p)
	}

	if a.Fail != nil {
		dto.Fail = ToAuditFailDTO(a.Fail)
	}

	return dto
}

func ToAuditFailDTO(f *domain.AuditFail) *AuditFailDTO {
	return &AuditFailDTO{
		Title:       f.Title,
		Description: f.Description,
	}
}

func ToDetailItemDTO(d domain.Detail) DetailItemDTO {
	return DetailItemDTO{
		I18nLabel: d.I18nLabel,
		Value:     d.Value,
		Type:      d.Type,
	}
}

func ToProblemDTO(p domain.Problem) ProblemDTO {
	dto := ProblemDTO{
		I18nNamespace:   p.I18nNamespace,
		DescriptionVars: p.DescriptionVars,
		Resources:       make([]ResourceDTO, len(p.Resources)),
	}

	for i, r := range p.Resources {
		dto.Resources[i] = ToResourceDTO(r)
	}

	return dto
}

func ToResourceDTO(r domain.Resource) ResourceDTO {
	return ResourceDTO{
		Title: r.Title,
		URL:   r.URL.String(),
	}
}
