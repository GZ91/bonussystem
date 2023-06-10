package handlers

import (
	"context"
	"net/http"
)

type Service interface {
}

type Handlers struct {
	NodeService Service
}

func New(ctx context.Context, service Service) (*Handlers, error) {
	return &Handlers{NodeService: service}, nil
}

func (h *Handlers) Orders(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) Balance(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) WithdrawalsGet(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) Withdraw(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) WithdrawalsPost(w http.ResponseWriter, r *http.Request) {

}
