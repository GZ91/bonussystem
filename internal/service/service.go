package service

import (
	"context"
	"github.com/GZ91/bonussystem/internal/errorsapp"
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
	CreateOrder(ctx context.Context, number, userID string) error
	GetOrders(ctx context.Context, userID string) ([]models.DataOrder, error)
	GetBalance(ctx context.Context, userID string) (current float64, withdrawn float64, err error)
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

func (r *NodeService) DownloadOrder(ctx context.Context, number, userID string) error {

	if !luhnAlgorithm(number) {
		return errorsapp.ErrIncorrectOrderNumber
	}

	err := r.nodeStorage.CreateOrder(ctx, number, userID)
	if err != nil {
		return err
	}

	return nil
}

func luhnAlgorithm(number string) bool {
	sum := 0
	isSecond := false
	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		if isSecond {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		isSecond = !isSecond
	}
	return sum%10 == 0
}

func (r *NodeService) GetOrders(ctx context.Context, userID string) ([]models.DataOrder, error) {
	return r.nodeStorage.GetOrders(ctx, userID)
}

func (r *NodeService) GetBalance(ctx context.Context, userID string) (models.DataBalance, error) {
	current, withdrawn, err := r.nodeStorage.GetBalance(ctx, userID)
	if err != nil {
		return models.DataBalance{}, err
	}
	data := models.DataBalance{Current: current, Withdrawn: withdrawn}
	return data, nil
}
