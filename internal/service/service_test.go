package service

import (
	"context"
	"github.com/GZ91/bonussystem/internal/app/logger"
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
				suite.Assert().Errorf(err, "DownloadOrder() error = %v, wantErr %v", tt.wantErr)
			}
		})
	}
}
