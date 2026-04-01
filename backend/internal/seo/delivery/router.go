package delivery

import (
	_ "backend/docs"
	"backend/internal/seo/domain"
	"backend/internal/seo/infrastructure"
	"backend/internal/seo/usecase"

	"github.com/gin-gonic/gin"
)

func SetupRouter(api *gin.RouterGroup, cache domain.ReportRepo) {
	scanner := infrastructure.NewWebScanner()
	usecase := usecase.NewScanUsecase(scanner, cache)
	handler := NewScanHandler(usecase)

	api.GET("/scan", handler.HandleScan)
}
