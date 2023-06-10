package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/GZ91/bonussystem/internal/errorsapp"
	"github.com/GZ91/bonussystem/internal/models"
	"io"
	"net/http"
)

type Service interface {
	CreateNewUser(ctx context.Context, login, password string) (*http.Cookie, error)
	AuthenticationUser(ctx context.Context, login, password string) (*http.Cookie, error)
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
	textBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	var dataRegist models.DataRegisteration
	err = json.Unmarshal(textBody, &dataRegist)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if dataRegist.Login == "" || dataRegist.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("login or password is not filled in"))
		return
	}

	Cook, err := h.NodeService.CreateNewUser(r.Context(), dataRegist.Login, dataRegist.Password)
	if err != nil {
		if errors.Is(err, errorsapp.ErrLoginAlreadyBorrowed) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("login is already taken"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, Cook)
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	textBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	var dataRegist models.DataRegisteration
	err = json.Unmarshal(textBody, &dataRegist)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if dataRegist.Login == "" || dataRegist.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("login or password is not filled in"))
		return
	}
	cook, err := h.NodeService.AuthenticationUser(r.Context(), dataRegist.Login, dataRegist.Password)
	if err != nil {
		if errors.Is(err, errorsapp.ErrNoFoundUser) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, cook)
}

func (h *Handlers) Withdraw(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) WithdrawalsPost(w http.ResponseWriter, r *http.Request) {

}
