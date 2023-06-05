package app

import (
	"context"
	"github.com/GZ91/bonussystem/internal/api/http/server"
	"github.com/GZ91/bonussystem/internal/app/config"
	"github.com/GZ91/bonussystem/internal/app/logger"
	"go.uber.org/zap"
)

var appLink *App

type App struct {
	config *config.Config
}

func New(config *config.Config) *App {
	if appLink == nil {
		appLink = &App{
			config,
		}
		return appLink
	}
	return appLink
}

func (r App) Run(ctx context.Context) {
	if err := server.Start(ctx, r.config); err != nil {
		logger.Log.Error("server startup error", zap.Error(err))
	}
}
