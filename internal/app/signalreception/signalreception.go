package signalreception

import (
	"github.com/GZ91/linkreduct/internal/app/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type CloserInterf interface {
	Close() error
}

type Stopper struct {
	CloserInterf
	Name string
}

func (s Stopper) GetName() string {
	return s.Name
}

type Closer interface {
	Close() error
	GetName() string
}

func OnClose(players []Closer, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for sig := range c {
		switch sig {
		case syscall.SIGINT:
			for _, player := range players {
				if err := player.Close(); err != nil {
					logger.Log.Error(player.GetName()+" close error", zap.String("error", err.Error()))
				} else {
					logger.Log.Info(player.GetName()+" close", zap.String("status", sig.String()))
				}
			}
			return
		default:
			logger.Log.Info("there was a signal", zap.String("status", sig.String()))
		}

	}

}
