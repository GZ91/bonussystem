package server

import (
	"context"
	"errors"
	"github.com/GZ91/bonussystem/internal/api/http/handlers"
	"github.com/GZ91/bonussystem/internal/api/http/middleware/authentication"
	"github.com/GZ91/bonussystem/internal/app/logger"
	"github.com/GZ91/bonussystem/internal/app/signalreception"
	"github.com/GZ91/bonussystem/internal/service"
	"github.com/GZ91/bonussystem/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

type Configer interface {
	service.Configer
	storage.Configer
	GetAddressPort() string
}

type NodeStorager interface {
	service.Storage
	Close() error
}

func Start(ctx context.Context, conf Configer) error {

	NodeStorager, err := storage.New(ctx, conf)
	if err != nil {
		return err
	}
	NodeService, err := service.New(ctx, NodeStorager, conf)
	if err != nil {
		return err
	}

	handls, err := handlers.New(ctx, NodeService)

	router := chi.NewRouter()

	NodeAuthentication := authentication.New(conf)
	router.Use(NodeAuthentication.Authentication)

	router.Get("/api/user/orders", handls.OrdersGet)
	router.Get("/api/user/balance", handls.Balance)
	router.Get("/api/user/withdrawals", handls.Withdrawals)

	router.Post("/api/user/orders", handls.OrdersPost)
	router.Post("/api/user/register", handls.Register)
	router.Post("/api/user/login", handls.Login)
	router.Post("/api/user/balance/withdraw", handls.Withdraw)

	Server := http.Server{}
	Server.Addr = conf.GetAddressPort()
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
