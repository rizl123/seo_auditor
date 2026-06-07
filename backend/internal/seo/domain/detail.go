package domain

import "time"

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

func NewTextDetail(label string, value string) Detail {
	return Detail{
		I18nLabel: label,
		Value:     value,
		Type:      DetailTypeText,
	}
}

func NewNumberDetail(label string, value int) Detail {
	return Detail{
		I18nLabel: label,
		Value:     value,
		Type:      DetailTypeNumber,
	}
}

func NewDurationDetail(label string, duration time.Duration) Detail {
	return Detail{
		I18nLabel: label,
		Value:     duration.Milliseconds(),
		Type:      DetailTypeDuration,
	}
}

func NewURLDetail(label string, value string) Detail {
	return Detail{
		I18nLabel: label,
		Value:     value,
		Type:      DetailTypeURL,
	}
}

func NewImageDetail(label string, value string) Detail {
	return Detail{
		I18nLabel: label,
		Value:     value,
		Type:      DetailTypeImage,
	}
}

func NewBadgeDetail(label string, value string) Detail {
	return Detail{
		I18nLabel: label,
		Value:     value,
		Type:      DetailTypeBadge,
	}
}
