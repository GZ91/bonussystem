package main

import (
	"context"
	"github.com/GZ91/bonussystem/internal/app"
	"github.com/GZ91/bonussystem/internal/app/config"
	"github.com/GZ91/bonussystem/internal/app/logger"
)

func main() {
	logger.Initializing("error")
	conf := config.New()
	app := app.New(conf)
	app.Run(context.Background())
}
