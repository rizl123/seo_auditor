package main

import (
	"backend/internal/delivery/http"
	"backend/internal/infrastructure"
	"backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	repo := &infrastructure.HttpSeoRepository{}
	seoUsecase := usecase.NewSeoUsecase(repo)
	handler := http.NewSeoHandler(seoUsecase)
	r := gin.Default()
	r.RedirectTrailingSlash = true

	http.SetupRouter(r, handler)

	r.Run(":8080")
}
