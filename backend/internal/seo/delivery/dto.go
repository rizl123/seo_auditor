package delivery

import (
	"backend/internal/seo/domain"
	neturl "net/url"
	"time"
)

type PageReportDTO struct {
	URL       string          `json:"url"`
	Status    int             `json:"status"`
	IsCached  bool            `json:"is_cached"`
	ScannedAt time.Time       `json:"scanned_at"`
	Metadata  *MetadataDTO    `json:"metadata,omitempty"`
	Network   *NetworkInfoDTO `json:"network,omitempty"`
}

func ToPageReportDTO(report *domain.PageReport) *PageReportDTO {
	if report == nil {
		return nil
	}
	dto := &PageReportDTO{
		URL:       report.URL.String(),
		Status:    report.Status,
		IsCached:  report.IsCached,
		ScannedAt: report.ScannedAt,
	}

	if report.Metadata != nil {
		dto.Metadata = &MetadataDTO{
			Title:       report.Metadata.Title,
			Description: report.Metadata.Description,
			H1:          report.Metadata.H1,
			Canonical:   report.Metadata.Canonical,
			OgImage:     report.Metadata.OgImage,
		}
	}

	if report.Network != nil {
		dto.Network = &NetworkInfoDTO{
			ResponseTimeMs: report.Network.ResponseTime.Milliseconds(),
			Server:         report.Network.Server,
			ContentType:    report.Network.ContentType,
		}
	}

	return dto
}

func ToPageReport(dto PageReportDTO) domain.PageReport {
	u, _ := neturl.Parse(dto.URL)
	report := domain.PageReport{
		URL:       u,
		Status:    dto.Status,
		IsCached:  dto.IsCached,
		ScannedAt: dto.ScannedAt,
	}

	if dto.Metadata != nil {
		report.Metadata = &domain.Metadata{
			Title:       dto.Metadata.Title,
			Description: dto.Metadata.Description,
			H1:          dto.Metadata.H1,
			Canonical:   dto.Metadata.Canonical,
			OgImage:     dto.Metadata.OgImage,
		}
	}

	if dto.Network != nil {
		report.Network = &domain.NetworkInfo{
			ResponseTime: time.Duration(dto.Network.ResponseTimeMs) * time.Millisecond,
			Server:       dto.Network.Server,
			ContentType:  dto.Network.ContentType,
		}
	}

	return report
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
