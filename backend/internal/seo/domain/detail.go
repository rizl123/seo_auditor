package domain

type DetailType string

const (
	DetailTypeText     DetailType = "text"
	DetailTypeNumber   DetailType = "number"
	DetailTypeDuration DetailType = "duration_ms"
	DetailTypeURL      DetailType = "url"
	DetailTypeImage    DetailType = "image"
	DetailTypeBadge    DetailType = "badge"
)

type Detail struct {
	I18nLabel string
	Value     any
	Type      DetailType
}

func NewDetail(label string, value any, dtype DetailType) Detail {
	return Detail{
		I18nLabel: label,
		Value:     value,
		Type:      dtype,
	}
}
