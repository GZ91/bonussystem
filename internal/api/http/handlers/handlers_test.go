package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/GZ91/bonussystem/internal/api/http/handlers/mocks"
	"github.com/GZ91/bonussystem/internal/app/logger"
	"github.com/GZ91/bonussystem/internal/errorsapp"
	"github.com/GZ91/bonussystem/internal/models"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
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

func (suite *TestSuite) TestOrdersGet() {
	type fields struct {
		NodeService *mocks.Service
		userID      string
		orders      []models.DataOrder
		err         error
		returnData  string
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		expectedCode int
		expectedBody string
	}{
		{
			name: "Test 1",
			fields: fields{
				NodeService: suite.NodeService,
				userID:      "userID1",
				orders:      nil,
				err:         errorsapp.ErrNoRecords,
			},
			args: args{r: httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("")),
				w: httptest.NewRecorder()},
			expectedCode: http.StatusNoContent,
		},
		{
			name: "Test 2",
			fields: fields{
				NodeService: suite.NodeService,
				userID:      "userID2",
				orders:      nil,
				err:         errors.New("Test"),
			},
			args: args{r: httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("")),
				w: httptest.NewRecorder()},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "Test 3",
			fields: fields{
				NodeService: suite.NodeService,
				userID:      "userID3",
				orders: []models.DataOrder{
					{
						Status:     "NEW",
						UploadedAt: "",
						Accrual:    0,
					},
				},
				err: nil,
			},
			expectedBody: "[{\"number\":\"\",\"status\":\"NEW\",\"uploaded_at\":\"\"}]",
			args: args{r: httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("")),
				w: httptest.NewRecorder()},
			expectedCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			h := &Handlers{
				NodeService: tt.fields.NodeService,
			}
			var userIDCTX models.CtxString = "userID"
			tt.args.r = tt.args.r.WithContext(context.WithValue(tt.args.r.Context(), userIDCTX, tt.fields.userID))
			suite.NodeService.EXPECT().GetOrders(tt.args.r.Context(), tt.fields.userID).Return(tt.fields.orders, tt.fields.err)
			h.OrdersGet(tt.args.w, tt.args.r)
			suite.Assert().Equal(tt.args.w.Code, tt.expectedCode)
			actualBody, _ := io.ReadAll(tt.args.w.Body)
			suite.Assert().Equal(tt.expectedBody, string(actualBody))
			tt.args.r.Body.Close()
		})
	}
}

func (suite *TestSuite) TestHandlers_Balance() {
	type fields struct {
		NodeService *mocks.Service
		userID      string
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		returnstatus int
		err          error
	}{
		{
			name: "test1",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/api/user/orders", strings.NewReader("")),
				w: httptest.NewRecorder(),
			},
			fields: fields{
				userID:      "user1",
				NodeService: suite.NodeService},
			returnstatus: http.StatusOK,
			err:          nil,
		},
		{
			name: "test2",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/api/user/orders", strings.NewReader("")),
				w: httptest.NewRecorder(),
			},
			returnstatus: http.StatusInternalServerError,
			err:          errors.New("test"),
			fields: fields{
				userID:      "user2",
				NodeService: suite.NodeService},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			h := &Handlers{
				NodeService: tt.fields.NodeService,
			}
			var userIDCTX models.CtxString = "userID"
			tt.args.r = tt.args.r.WithContext(context.WithValue(tt.args.r.Context(), userIDCTX, tt.fields.userID))
			dataForEqual := models.DataBalance{Current: 1000, Withdrawn: 0}
			tt.fields.NodeService.EXPECT().GetBalance(tt.args.r.Context(), tt.fields.userID).Return(dataForEqual, tt.err).Maybe()
			h.Balance(tt.args.w, tt.args.r)

			suite.Assert().Equal(tt.args.w.Code, tt.returnstatus)
			if tt.args.w.Code < 400 {
				suite.Equal(tt.returnstatus, tt.args.w.Code)
				bodyText, err := io.ReadAll(tt.args.w.Body)
				if err != nil {
					panic(err)
				}
				var expectedBody models.DataBalance
				err = json.Unmarshal(bodyText, &expectedBody)
				if err != nil {
					panic(err)
				}
				suite.Assert().True(reflect.DeepEqual(expectedBody, dataForEqual))
			}
		})
	}
}

func (suite *TestSuite) TestHandlers_Withdrawals() {
	type fields struct {
		NodeService   *mocks.Service
		userID        string
		dataReturn    []models.WithdrawalsData
		errorReturned error
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		responseCode int
	}{
		{
			name: "Test 1",
			fields: fields{
				userID: "userTest",
				dataReturn: []models.WithdrawalsData{
					{Order: "1234234234", ProcessedAt: "safewefw21321", Sum: 10},
					{Order: "123422141434234", ProcessedAt: "safewefw21123321", Sum: 1033},
				},
				NodeService:   suite.NodeService,
				errorReturned: nil,
			},
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/api/user/orders", strings.NewReader("")),
				w: httptest.NewRecorder(),
			},
			responseCode: http.StatusOK,
		},
		{
			name: "Test 2",
			fields: fields{
				userID: "userTest2",
				dataReturn: []models.WithdrawalsData{
					{Order: "1234234234", ProcessedAt: "safewefw21321", Sum: 10},
					{Order: "123422141434234", ProcessedAt: "safewefw21123321", Sum: 1033},
				},
				NodeService:   suite.NodeService,
				errorReturned: errorsapp.ErrNoRecords,
			},
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/api/user/orders", strings.NewReader("")),
				w: httptest.NewRecorder(),
			},
			responseCode: http.StatusNoContent,
		},
		{
			name: "Test 3",
			fields: fields{
				userID: "userTest3",
				dataReturn: []models.WithdrawalsData{
					{Order: "1234234234", ProcessedAt: "safewefw21321", Sum: 10},
					{Order: "123422141434234", ProcessedAt: "safewefw21123321", Sum: 1033},
				},
				NodeService:   suite.NodeService,
				errorReturned: errors.New("Test Error"),
			},
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/api/user/orders", strings.NewReader("")),
				w: httptest.NewRecorder(),
			},
			responseCode: http.StatusInternalServerError,
		},
		{
			name: "Test 3",
			fields: fields{
				userID:        "userTest4",
				dataReturn:    nil,
				NodeService:   suite.NodeService,
				errorReturned: nil,
			},
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/api/user/orders", strings.NewReader("")),
				w: httptest.NewRecorder(),
			},
			responseCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			h := &Handlers{
				NodeService: tt.fields.NodeService,
			}
			var userIDCTX models.CtxString = "userID"
			tt.args.r = tt.args.r.WithContext(context.WithValue(tt.args.r.Context(), userIDCTX, tt.fields.userID))
			suite.NodeService.EXPECT().Withdrawals(tt.args.r.Context(), tt.fields.userID).Return(tt.fields.dataReturn, tt.fields.errorReturned).Maybe()
			h.Withdrawals(tt.args.w, tt.args.r)
			if tt.fields.dataReturn == nil {
				dataJSON, err := json.Marshal(tt.fields.dataReturn)
				suite.Assert().NoError(err)
				byteReturned, err := io.ReadAll(tt.args.w.Body)
				suite.Assert().NoError(err)
				suite.Equal(dataJSON, byteReturned)
			} else {
				suite.Equal(tt.args.w.Code, tt.responseCode)
			}
		})
	}
}

func (suite *TestSuite) TestHandlers_Register() {
	type fields struct {
		NodeService *mocks.Service
		cook        *http.Cookie
		login       string
		password    string
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		statusReturn int
		errReturn    error
	}{
		{
			name: "Test 1",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/api/user/register",
					strings.NewReader(`{ "login": "t1", "password": "t2"}`)),
			},
			fields: fields{
				NodeService: suite.NodeService,
				cook:        &http.Cookie{},
				login:       "t1",
				password:    "t2"},
			statusReturn: http.StatusOK,
			errReturn:    nil,
		},
		{
			name: "Test 2",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/api/user/register",
					strings.NewReader(`{ "login": "t3", "password":""}`)),
			},
			fields: fields{NodeService: suite.NodeService, cook: &http.Cookie{},
				login:    "t3",
				password: ""},
			statusReturn: http.StatusBadRequest,
			errReturn:    nil,
		},
		{
			name: "Test 3",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/api/user/register",
					strings.NewReader(`{ "login": "t5", "password": "t62"}`)),
			},
			fields: fields{NodeService: suite.NodeService, cook: &http.Cookie{},
				login:    "t5",
				password: "t62"},
			statusReturn: http.StatusConflict,
			errReturn:    errorsapp.ErrLoginAlreadyBorrowed,
		},
		{
			name: "Test 4",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/api/user/register",
					strings.NewReader(`{ "login": "t6", "password": "t7"}`)),
			},
			fields: fields{NodeService: suite.NodeService, cook: &http.Cookie{},
				login:    "t6",
				password: "t7"},
			statusReturn: http.StatusInternalServerError,
			errReturn:    errors.New("Test error"),
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			h := &Handlers{
				NodeService: tt.fields.NodeService,
			}
			suite.NodeService.EXPECT().CreateNewUser(tt.args.r.Context(), tt.fields.login, tt.fields.password).
				Return(tt.fields.cook, tt.errReturn)
			h.Register(tt.args.w, tt.args.r)
			suite.Assert().Equal(tt.statusReturn, tt.args.w.Code)
		})
	}
}

func (suite *TestSuite) TestHandlers_Login() {
	type fields struct {
		NodeService Service
		cook        *http.Cookie
		login       string
		password    string
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		statusReturn int
		errReturn    error
	}{
		{
			name: "Test 1",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/api/user/login", strings.NewReader(`{ "login": "<login>", "password": "<password>"}`)),
				w: httptest.NewRecorder(),
			},
			fields: fields{
				NodeService: suite.NodeService, cook: &http.Cookie{},
				login:    "<login>",
				password: "<password>",
			},
			statusReturn: http.StatusOK,
			errReturn:    nil,
		},
		{
			name: "Test 2",
			args: args{
				r: httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(`{ "login": "<login>2", "password": ""}`)),
				w: httptest.NewRecorder(),
			},
			fields: fields{
				NodeService: suite.NodeService, cook: &http.Cookie{},
				login:    "<login>2",
				password: "<password>2",
			},
			statusReturn: http.StatusBadRequest,
			errReturn:    nil,
		},
		{
			name: "Test 3",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/api/user/login", strings.NewReader(`{ "login": "<login>3", "password": "<password>3"}`)),
				w: httptest.NewRecorder(),
			},
			fields: fields{
				NodeService: suite.NodeService, cook: &http.Cookie{},
				login:    "<login>3",
				password: "<password>3",
			},
			statusReturn: http.StatusUnauthorized,
			errReturn:    errorsapp.ErrNoFoundUser,
		},

		{
			name: "Test 4",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/api/user/login", strings.NewReader(`{ "login": "<login>4", "password": "<password>4"}`)),
				w: httptest.NewRecorder(),
			},
			fields: fields{
				NodeService: suite.NodeService, cook: &http.Cookie{},
				login:    "<login>4",
				password: "<password>4",
			},
			statusReturn: http.StatusInternalServerError,
			errReturn:    errors.New("TestError"),
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			h := &Handlers{
				NodeService: tt.fields.NodeService,
			}
			suite.NodeService.EXPECT().AuthenticationUser(tt.args.r.Context(), tt.fields.login, tt.fields.password).
				Return(tt.fields.cook, tt.errReturn)
			h.Login(tt.args.w, tt.args.r)
			suite.Assert().Equal(tt.statusReturn, tt.args.w.Code)
		})
	}
}

func (suite *TestSuite) TestHandlers_Withdraw() {
	type fields struct {
		NodeService  Service
		statusReturn int
		err          error
		userID       string
		data         models.WithdrawData
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test1",
			fields: fields{
				NodeService:  suite.NodeService,
				statusReturn: http.StatusOK,
				err:          nil,
				userID:       "userTest1",
				data:         models.WithdrawData{Sum: 751, Order: "2377225624"},
			},
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader(`{"sum":751.0, "order":"2377225624"}`)),
			},
		},
		{
			name: "test2",
			fields: fields{
				NodeService:  suite.NodeService,
				statusReturn: http.StatusPaymentRequired,
				err:          errorsapp.ErrInsufficientFunds,
				userID:       "userTest2",
				data:         models.WithdrawData{Sum: 751, Order: "2377225624"},
			},
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader(`{"sum":751.0, "order":"2377225624"}`)),
			},
		},
		{
			name: "test3",
			fields: fields{
				NodeService:  suite.NodeService,
				statusReturn: http.StatusUnprocessableEntity,
				err:          errorsapp.ErrIncorrectOrderNumber,
				userID:       "userTest3",
				data:         models.WithdrawData{Sum: 751, Order: "2377225624"},
			},
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader(`{"sum":751.0, "order":"2377225624"}`)),
			},
		},
		{
			name: "test4",
			fields: fields{
				NodeService:  suite.NodeService,
				statusReturn: http.StatusInternalServerError,
				err:          errors.New("Test error"),
				userID:       "userTest4",
				data:         models.WithdrawData{Sum: 751, Order: "2377225624"},
			},
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader(`{"sum":751.0, "order":"2377225624"}`)),
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			h := &Handlers{
				NodeService: tt.fields.NodeService,
			}
			var userIDCTX models.CtxString = "userID"
			tt.args.r = tt.args.r.WithContext(context.WithValue(tt.args.r.Context(), userIDCTX, tt.fields.userID))
			suite.NodeService.EXPECT().Withdraw(tt.args.r.Context(), tt.fields.data, tt.fields.userID).Return(tt.fields.err)
			h.Withdraw(tt.args.w, tt.args.r)
			suite.Assert().Equal(tt.fields.statusReturn, tt.args.w.Code)
		})
	}
}
