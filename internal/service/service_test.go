package service

import (
	"context"
	"errors"
	"github.com/GZ91/bonussystem/internal/app/logger"
	"github.com/GZ91/bonussystem/internal/errorsapp"
	"github.com/GZ91/bonussystem/internal/service/mocks"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
	"time"
)

type TestSuite struct {
	suite.Suite
	NodeStorage *mocks.Storage
	Config      *mocks.Configer
}

func (suite *TestSuite) SetupTest() {
	suite.NodeStorage = new(mocks.Storage)
	suite.Config = new(mocks.Configer)
	logger.Initializing("info")
	suite.Config.EXPECT().GetAddressAccrual().Return("addressAccrual").Maybe()
	suite.Config.EXPECT().GetSecretKey().Return("SecretKey").Maybe()
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestDownloadOrder() {
	type fields struct {
		nodeStorage *mocks.Storage
		conf        *mocks.Configer
		mutexOrder  *sync.RWMutex
		orderLocks  map[string]chan struct{}
		mutexClient *sync.RWMutex
		clientLocks map[string]chan struct{}
	}
	type args struct {
		ctx    context.Context
		number string
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "Test 1",
			fields: fields{
				nodeStorage: suite.NodeStorage,
				conf:        suite.Config,
				orderLocks:  make(map[string]chan struct{}),
				clientLocks: make(map[string]chan struct{}),
			},
			wantErr: false,
			args: args{
				userID: "user5",
				ctx:    context.Background(),
				number: "12345678903"},
			err: nil,
		},
		{
			name: "Test 2",
			fields: fields{
				nodeStorage: suite.NodeStorage,
				conf:        suite.Config,
				orderLocks:  make(map[string]chan struct{}),
				clientLocks: make(map[string]chan struct{}),
			},
			wantErr: false,
			args: args{
				userID: "user6",
				ctx:    context.Background(),
				number: "12345678902"},
			err: errorsapp.ErrIncorrectOrderNumber,
		},
		{
			name: "Test 3",
			fields: fields{
				nodeStorage: suite.NodeStorage,
				conf:        suite.Config,
				orderLocks:  make(map[string]chan struct{}),
				clientLocks: make(map[string]chan struct{}),
			},
			wantErr: false,
			args: args{
				userID: "user7",
				ctx:    context.Background(),
				number: "12345678903"},
			err: errors.New("test"),
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			r := &NodeService{
				nodeStorage: tt.fields.nodeStorage,
				conf:        tt.fields.conf,
				orderLocks:  tt.fields.orderLocks,
				clientLocks: tt.fields.clientLocks,
			}
			suite.NodeStorage.EXPECT().CreateOrder(tt.args.ctx, tt.args.number, tt.args.userID).Return(tt.err)
			if err := r.DownloadOrder(tt.args.ctx, tt.args.number, tt.args.userID); (err != nil) != tt.wantErr {
				suite.Assert().Equal(tt.err, err)
			}
		})
	}
}

func (suite *TestSuite) Test_luhnAlgorithm() {
	type args struct {
		number string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test 1",
			args: args{number: "12345678903"},
			want: true,
		},
		{
			name: "Test 2",
			args: args{number: "12345678902"},
			want: true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if got := luhnAlgorithm(tt.args.number); got != tt.want {
				suite.Errorf(errors.New("not correct returned"), "not correct returned")
			}
		})
	}
}

func TestNodeService_LockOrder(t *testing.T) {
	var v NodeService
	v.orderLocks = make(map[string]chan struct{})
	order := "sakpfoafskasf"
	go func() {
		v.LockOrder(order)
		time.Sleep(5 * time.Second)
		v.UnclockOrder(order)
	}()
	v.LockOrder(order)
	v.UnclockOrder(order)
}

func TestNodeService_LockClients(t *testing.T) {
	var v NodeService
	v.clientLocks = make(map[string]chan struct{})
	usderID := "sakpfoafskasf"
	go func() {
		v.LockClient(usderID)
		time.Sleep(5 * time.Second)
		v.UnclockClient(usderID)
	}()
	v.LockClient(usderID)
	v.UnclockClient(usderID)
}
