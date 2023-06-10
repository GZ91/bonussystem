package main

import (
	"context"
	"github.com/GZ91/bonussystem/internal/app"
	"github.com/GZ91/bonussystem/internal/app/initializing"
)

func main() {
	conf := initializing.Configuration()
	app := app.New(conf)
	app.Run(context.Background())
}
