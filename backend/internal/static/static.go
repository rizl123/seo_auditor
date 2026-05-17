package static

import (
	"embed"
	"io/fs"
)

//go:embed locales/*.json
var localesFS embed.FS

func GetLocalesFS() fs.FS {
	subFS, err := fs.Sub(localesFS, "locales")
	if err != nil {
		panic("failed to sub locales directory: " + err.Error())
	}
	return subFS
}
