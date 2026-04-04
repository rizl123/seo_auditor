package delivery

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API, handler *ScanHandler) {
	huma.Register(api, huma.Operation{
		OperationID: "get-scan-report",
		Method:      http.MethodGet,
		Path:        "/api/scan",
		Summary:     "Scan page and get report",
		Tags:        []string{"scanner"},
	}, handler.HandleScan)
}
