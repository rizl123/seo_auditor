package http

import (
	"backend/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SeoHandler struct {
	usecase *usecase.SeoUsecase
}

func NewSeoHandler(u *usecase.SeoUsecase) *SeoHandler {
	return &SeoHandler{usecase: u}
}

// Analyze godoc
//
//	@Summary	Analyze SEO data
//	@Tags		seo
//	@Param		url	query		string	true	"URL"
//	@Success	200	{object}	domain.SeoData
//	@Router		/api/analyze [get]
func (h *SeoHandler) Analyze(c *gin.Context) {
	urlStr := c.Query("url")
	if urlStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}

	data, err := h.usecase.Analyze(urlStr)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
