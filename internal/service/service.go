package service

import (
	"context"
	"github.com/GZ91/bonussystem/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Configer interface {
	GetAddressAccrual() string
	GetSecretKey() string
}

type Storage interface {
	CreateNewUser(context.Context, string, string, string) error
	AuthenticationUser(ctx context.Context, login, password string) (string, error)
}

type NodeService struct {
	nodeStorage Storage
	conf        Configer
}

func New(ctx context.Context, storage Storage, conf Configer) (*NodeService, error) {
	return &NodeService{nodeStorage: storage, conf: conf}, nil
}

func (r *NodeService) CreateNewUser(ctx context.Context, login, password string) (*http.Cookie, error) {

	userID := uuid.New().String()
	err := r.nodeStorage.CreateNewUser(ctx, userID, login, password)
	if err != nil {
		return nil, err
	}

	cook, err := getAuthorizationCookie(r.conf.GetSecretKey(), userID)
	if err != nil {
		return nil, err
	}
	return cook, nil
}

func getAuthorizationCookie(SecretKey, userID string) (*http.Cookie, error) {
	cook := &http.Cookie{}
	cook.Name = "Authorization"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.Claims{
		UserID:           userID,
		RegisteredClaims: &jwt.RegisteredClaims{ExpiresAt: &jwt.NumericDate{time.Now().Add(time.Hour * 24)}},
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return nil, err
	}
	cook.Value = tokenString
	return cook, nil
}

func (r *NodeService) AuthenticationUser(ctx context.Context, login, password string) (*http.Cookie, error) {
	userID, err := r.nodeStorage.AuthenticationUser(ctx, login, password)
	if err != nil {
		return nil, err
	}
	return getAuthorizationCookie(r.conf.GetSecretKey(), userID)
}
