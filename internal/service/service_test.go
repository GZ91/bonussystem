package service

import (
	"context"
	"errors"
	"github.com/GZ91/bonussystem/internal/app/logger"
	"github.com/GZ91/bonussystem/internal/errorsapp"
	mocksStorager "github.com/GZ91/bonussystem/internal/service/mocks"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
)

type TestSuite struct {
	suite.Suite
	NodeStorage *mocksStorager.Storage
	Config      *mocksStorager.Configer
}

func (suite *TestSuite) SetupTest() {
	suite.NodeStorage = new(mocksStorager.Storage)
	suite.Config = new(mocksStorager.Configer)
	logger.Initializing("info")
	suite.Config.EXPECT().GetAddressAccrual().Return("addressAccrual").Maybe()
	suite.Config.EXPECT().GetSecretKey().Return("SecretKey").Maybe()
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestDownloadOrder() {
	type fields struct {
		nodeStorage *mocksStorager.Storage
		conf        *mocksStorager.Configer
		mutexOrder  sync.RWMutex
		orderLocks  map[string]chan struct{}
		mutexClient sync.RWMutex
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
				mutexOrder:  tt.fields.mutexOrder,
				orderLocks:  tt.fields.orderLocks,
				mutexClient: tt.fields.mutexClient,
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
