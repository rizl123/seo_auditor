package main

import (
	"backend/internal/bootstrap"
	"backend/internal/config"
)

func main() {
	cfg := config.Load()
	cfg.Log()

	app := bootstrap.NewApp(cfg)
	app.Run()
}
