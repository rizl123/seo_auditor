package delivery

import (
	"backend/internal/seo/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScanHandler struct {
	usecase *usecase.ScanUsecase
}

func NewScanHandler(u *usecase.ScanUsecase) *ScanHandler {
	return &ScanHandler{usecase: u}
}

// HandleScan godoc
//
//	@Summary    Scan page and get report
//	@Tags       scanner
//	@Param      url   query     string   true   "URL to scan"
//	@Success    200   {object}  domain.PageReport
//	@Router     /api/scan [get]
func (h *ScanHandler) HandleScan(c *gin.Context) {
	urlStr := c.Query("url")
	if urlStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}

	report, err := h.usecase.Execute(c, urlStr)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}
