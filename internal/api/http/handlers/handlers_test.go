package handlers

import (
	"context"
	"errors"
	"github.com/GZ91/bonussystem/internal/api/http/handlers/mocks"
	"github.com/GZ91/bonussystem/internal/app/logger"
	"github.com/GZ91/bonussystem/internal/errorsapp"
	"github.com/GZ91/bonussystem/internal/models"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestSuite struct {
	suite.Suite
	NodeService *mocks.Service
}

func (suite *TestSuite) SetupTest() {
	suite.NodeService = new(mocks.Service)
	logger.Initializing("info")
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestHandlers_OrdersPost() {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	type field struct {
		numberOrder string
		userID      string
	}
	tests := []struct {
		name           string
		field          field
		args           args
		expectedStatus int
		error          error
	}{
		{
			name: "Test 1",
			field: field{
				numberOrder: "12345678903",
				userID:      "TestUser1",
			},
			args: args{
				r: httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("12345678903")),
				w: httptest.NewRecorder(),
			},
			expectedStatus: http.StatusAccepted,
			error:          nil,
		},
		{
			name: "Test 2",
			field: field{
				numberOrder: "12545678907",
				userID:      "TestUser2",
			},
			args: args{
				r: httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("12545678907")),
				w: httptest.NewRecorder(),
			},
			expectedStatus: http.StatusUnprocessableEntity,
			error:          errorsapp.ErrIncorrectOrderNumber,
		},
		{
			name: "Test 3",
			field: field{
				numberOrder: "12345678903",
				userID:      "TestUser3",
			},
			args: args{
				r: httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("12345678903")),
				w: httptest.NewRecorder(),
			},
			expectedStatus: http.StatusOK,
			error:          errorsapp.ErrOrderAlreadyThisUser,
		},
		{
			name: "Test 4",
			field: field{
				numberOrder: "12345678903",
				userID:      "TestUser4",
			},
			args: args{
				r: httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("12345678903")),
				w: httptest.NewRecorder(),
			},
			expectedStatus: http.StatusConflict,
			error:          errorsapp.ErrOrderAlreadyAnotherUser,
		},
		{
			name: "Test 5",
			field: field{
				numberOrder: "12345678903",
				userID:      "TestUser5",
			},
			args: args{
				r: httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("12345678903")),
				w: httptest.NewRecorder(),
			},
			expectedStatus: http.StatusInternalServerError,
			error:          errors.New("Test error"),
		},
	}
	ctx := context.Background()

	for _, tt := range tests {
		h := &Handlers{
			NodeService: suite.NodeService,
		}
		suite.Run(tt.name, func() {
			var userIDCTX models.CtxString = "userID"
			tt.args.r = tt.args.r.WithContext(context.WithValue(ctx, userIDCTX, tt.field.userID))
			suite.NodeService.EXPECT().DownloadOrder(tt.args.r.Context(), tt.field.numberOrder, tt.field.userID).Return(tt.error)
			h.OrdersPost(tt.args.w, tt.args.r)
			suite.Assert().Equal(tt.expectedStatus, tt.args.w.Code, "error status actual not status expected")
			tt.args.r.Body.Close()
		})
	}
}
