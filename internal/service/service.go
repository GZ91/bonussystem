package service

import "context"

type Storage interface {
}

type Node struct {
	NodeStorage Storage
}

func New(ctx context.Context, storage Storage) (*Node, error) {
	return &Node{NodeStorage: storage}, nil
}
