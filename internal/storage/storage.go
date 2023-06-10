package storage

import (
	"context"
	"github.com/GZ91/bonussystem/internal/app/config"
)

type Node struct {
	conf *config.Config
}

func New(ctx context.Context, conf *config.Config) (*Node, error) {
	return &Node{conf: conf}, nil
}

func (n *Node) Close() error {
	return nil
}
