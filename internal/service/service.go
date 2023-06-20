package service

import (
	"context"
	"encoding/json"
	"github.com/GZ91/bonussystem/internal/app/logger"
	"github.com/GZ91/bonussystem/internal/errorsapp"
	"github.com/GZ91/bonussystem/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"math"
	"net/http"
	"sync"
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
	Withdraw(ctx context.Context, NewCurrent, NewWithdraw float64, data models.WithdrawData, userID string) error
	Withdrawals(ctx context.Context, userID string) ([]models.WithdrawalsData, error)
	GetOrdersForProcessing(ctx context.Context) ([]models.DataForProcessing, error)
	NewBalance(ctx context.Context, NewCurrent float64, userID string) error
	NewStatusOrder(ctx context.Context, number, status string, accural float64) error
}

type NodeService struct {
	nodeStorage Storage
	conf        Configer
	mutexOrder  sync.RWMutex
	orderLocks  map[string]chan struct{}
	mutexClient sync.RWMutex
	clientLocks map[string]chan struct{}
}

func New(ctx context.Context, storage Storage, conf Configer) (*NodeService, error) {
	Node := &NodeService{nodeStorage: storage, conf: conf}
	Node.orderLocks = make(map[string]chan struct{})
	Node.clientLocks = make(map[string]chan struct{})
	go Node.ProcessingOrders(ctx)
	return Node, nil
}

func (r *NodeService) LockOrder(number string) {
	var val chan struct{}
	for {
		r.mutexOrder.RLock()
		var ok bool
		val, ok = r.orderLocks[number]
		if !ok {
			r.mutexOrder.RUnlock()
			r.mutexOrder.Lock()
			r.orderLocks[number] = make(chan struct{}, 1)
			r.mutexOrder.Unlock()
		} else {
			r.mutexOrder.RUnlock()
			break
		}
	}
	val <- struct{}{}
}

func (r *NodeService) UnclockOrder(number string) {
	r.mutexOrder.RLock()
	defer r.mutexOrder.RUnlock()
	val, ok := r.orderLocks[number]
	if !ok {
		return
	}
	<-val
}

func (r *NodeService) LockClient(userID string) {
	var val chan struct{}
	for {
		r.mutexClient.RLock()
		var ok bool
		val, ok = r.clientLocks[userID]
		if !ok {
			r.mutexClient.RUnlock()
			r.mutexClient.Lock()
			struc := make(chan struct{}, 1)
			r.clientLocks[userID] = struc
			r.mutexClient.Unlock()
		} else {
			r.mutexClient.RUnlock()
			break
		}
	}
	val <- struct{}{}
}

func (r *NodeService) UnclockClient(userID string) {
	r.mutexClient.RLock()
	defer r.mutexClient.RUnlock()
	val, ok := r.clientLocks[userID]
	if !ok {
		return
	}
	<-val
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
	r.LockOrder(number)
	defer r.UnclockOrder(number)
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
	data := models.DataBalance{Current: math.Floor(current*100) / 100, Withdrawn: math.Floor(withdrawn*100) / 100}
	return data, nil
}

func (r *NodeService) Withdraw(ctx context.Context, data models.WithdrawData, userID string) error {
	r.LockClient(userID)
	defer r.UnclockClient(userID)
	current, withdraw, err := r.nodeStorage.GetBalance(ctx, userID)
	if err != nil {
		return err
	}
	if !luhnAlgorithm(data.Order) {
		return errorsapp.ErrIncorrectOrderNumber
	}
	if data.Sum > current {
		return errorsapp.ErrInsufficientFunds
	}
	err = r.nodeStorage.Withdraw(ctx, current-data.Sum, withdraw+data.Sum, data, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *NodeService) Withdrawals(ctx context.Context, userID string) ([]models.WithdrawalsData, error) {
	return r.nodeStorage.Withdrawals(ctx, userID)
}

func (r *NodeService) ProcessingOrders(ctx context.Context) {

	addressServiceProcessing := r.conf.GetAddressAccrual()
	timer := time.NewTimer(time.Second * 10)
	for {
		dataForProcessing, err := r.nodeStorage.GetOrdersForProcessing(ctx)
		if err != nil {
			logger.Log.Error("GetOrdersForProcessing", zap.Error(err))
			break
		}
		for _, val := range dataForProcessing {
			responce, err := http.Get(addressServiceProcessing + "/" + val.Order)
			if err != nil {
				logger.Log.Error("http.Get:"+addressServiceProcessing+"/"+val.Order, zap.Error(err))
				break
			}
			textBody, err := io.ReadAll(responce.Body)
			if err != nil {
				logger.Log.Error("error when reading the request body", zap.Error(err))
				break
			}
			var data models.ResponceAccural
			err = json.Unmarshal(textBody, &data)
			if err != nil {
				logger.Log.Error("http.Get:"+addressServiceProcessing+"/"+val.Order, zap.Error(err))
				logger.Log.Error("error when convert body json in struct", zap.Error(err), zap.String("text", string(textBody)))
				break
			}
			if data.Status == val.Status {
				break
			}
			if data.Status == "PROCESSED" {
				r.LockClient(val.UserID)
				current, _, err := r.nodeStorage.GetBalance(ctx, val.UserID)
				if err != nil {
					logger.Log.Error("error when receiving a balance", zap.Error(err))
					break
				}
				newCurrent := current + data.Accural
				err = r.nodeStorage.NewBalance(ctx, newCurrent, val.UserID)
				if err != nil {
					logger.Log.Error("error when writing a new balance", zap.Error(err))
					break
				}
				r.UnclockClient(val.UserID)
			}
			err = r.nodeStorage.NewStatusOrder(ctx, data.Order, data.Status, data.Accural)
			if err != nil {
				logger.Log.Error("error when writing a new status order", zap.Error(err))
				break
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
		}
	}
}
