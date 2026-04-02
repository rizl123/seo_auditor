package delivery

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(api *gin.RouterGroup, handler *ScanHandler) {
	api.GET("/scan", handler.HandleScan)
}
