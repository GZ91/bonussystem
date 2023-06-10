package server

import (
	"context"
	"errors"
	"github.com/GZ91/bonussystem/internal/api/http/handlers"
	"github.com/GZ91/bonussystem/internal/app/config"
	"github.com/GZ91/bonussystem/internal/app/logger"
	"github.com/GZ91/bonussystem/internal/app/signalreception"
	"github.com/GZ91/bonussystem/internal/service"
	"github.com/GZ91/bonussystem/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

type NodeStorager interface {
	service.Storage
	Close() error
}

func Start(ctx context.Context, conf *config.Config) error {

	NodeStorager, err := storage.New(ctx, conf)
	if err != nil {
		return err
	}
	NodeService, err := service.New(ctx, NodeStorager)
	if err != nil {
		return err
	}

	handls, err := handlers.New(ctx, NodeService)

	router := chi.NewRouter()
	router.Use()

	router.Get("/api/user/orders", handls.Orders)
	router.Get("/api/user/balance", handls.Balance)
	router.Get("/api/user/withdrawals", handls.WithdrawalsGet)

	router.Post("/api/user/register", handls.Register)
	router.Post("/api/user/login", handls.Login)
	router.Post("/api/user/balance/withdraw", handls.Withdraw)
	router.Post("/api/user/withdrawals", handls.WithdrawalsPost)

	Server := http.Server{}
	Server.Addr = "localhost"
	Server.Handler = router

	wg := sync.WaitGroup{}
	go signalreception.OnClose([]signalreception.Closer{
		&signalreception.Stopper{CloserInterf: &Server, Name: "server"},
		&signalreception.Stopper{CloserInterf: NodeStorager, Name: "BD"},
	}, &wg)

	if err := Server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Error("server startup error", zap.String("error", err.Error()))
		}
	}

	return nil
}
