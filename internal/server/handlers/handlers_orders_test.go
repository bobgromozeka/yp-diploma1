package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/bobgromozeka/yp-diploma1/internal/app/dependencies"
	"github.com/bobgromozeka/yp-diploma1/internal/models"
	"github.com/bobgromozeka/yp-diploma1/internal/server/config"
	"github.com/bobgromozeka/yp-diploma1/internal/storage"
	mockstorage "github.com/bobgromozeka/yp-diploma1/internal/storage/mock"
	"github.com/bobgromozeka/yp-diploma1/internal/testutils"
)

func TestCreateWrongOrderNumber(t *testing.T) {
	body := strings.NewReader(WrongOrderNumber) //Wrong order number
	req := httptest.NewRequest("POST", "/api/user/orders", body)
	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("Authorization", "Bearer "+JWT)
	httpW := httptest.NewRecorder()
	config.Set(
		config.Config{
			JWTSecret: JWTSecret,
		},
	)

	d := dependencies.D{
		UsersStorage:       nil,
		OrdersStorage:      nil,
		WithdrawalsStorage: nil,
		DB:                 nil,
		Logger:             zap.NewExample().Sugar(),
	}

	h := MakeMux(d)

	h.ServeHTTP(httpW, req)

	respBody, _ := io.ReadAll(httpW.Body)
	assert.Equal(t, http.StatusUnprocessableEntity, httpW.Code)
	assert.Equal(t, "Wrong order format\n", string(respBody))
}

func TestCreateOrderAlreadyCreated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	oStorage := mockstorage.NewMockOrdersStorage(ctrl)
	oStorage.
		EXPECT().
		CreateOrder(testutils.MatchContext(), gomock.Eq(OrderNumber), gomock.Eq(int64(UserID))).
		Return(storage.ErrOrderAlreadyCreated)

	body := strings.NewReader(OrderNumber) //Wrong order number
	req := httptest.NewRequest("POST", "/api/user/orders", body)
	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("Authorization", "Bearer "+JWT)
	httpW := httptest.NewRecorder()
	config.Set(
		config.Config{
			JWTSecret: JWTSecret,
		},
	)

	d := dependencies.D{
		UsersStorage:       nil,
		OrdersStorage:      oStorage,
		WithdrawalsStorage: nil,
		DB:                 nil,
		Logger:             zap.NewExample().Sugar(),
	}

	m := MakeMux(d)

	m.ServeHTTP(httpW, req)

	respBody, _ := io.ReadAll(httpW.Body)
	assert.Equal(t, http.StatusOK, httpW.Code)
	assert.Equal(t, "Order already created", string(respBody))
}

func TestCreateOrderAlreadyCreatedByAnotherUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	oStorage := mockstorage.NewMockOrdersStorage(ctrl)
	oStorage.
		EXPECT().
		CreateOrder(testutils.MatchContext(), gomock.Eq(OrderNumber), gomock.Eq(int64(UserID))).
		Return(storage.ErrOrderForeign)

	body := strings.NewReader(OrderNumber) //Wrong order number
	req := httptest.NewRequest("POST", "/api/user/orders", body)
	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("Authorization", "Bearer "+JWT)
	httpW := httptest.NewRecorder()
	config.Set(
		config.Config{
			JWTSecret: JWTSecret,
		},
	)

	d := dependencies.D{
		UsersStorage:       nil,
		OrdersStorage:      oStorage,
		WithdrawalsStorage: nil,
		DB:                 nil,
		Logger:             zap.NewExample().Sugar(),
	}

	m := MakeMux(d)

	m.ServeHTTP(httpW, req)

	respBody, _ := io.ReadAll(httpW.Body)
	assert.Equal(t, http.StatusConflict, httpW.Code)
	assert.Equal(t, "Order created by another user\n", string(respBody))
}

func TestCreateOrderInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	oStorage := mockstorage.NewMockOrdersStorage(ctrl)
	oStorage.
		EXPECT().
		CreateOrder(testutils.MatchContext(), gomock.Eq(OrderNumber), gomock.Eq(int64(UserID))).
		Return(errors.New("internal server error"))

	body := strings.NewReader(OrderNumber) //Wrong order number
	req := httptest.NewRequest("POST", "/api/user/orders", body)
	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("Authorization", "Bearer "+JWT)
	httpW := httptest.NewRecorder()
	config.Set(
		config.Config{
			JWTSecret: JWTSecret,
		},
	)

	d := dependencies.D{
		UsersStorage:       nil,
		OrdersStorage:      oStorage,
		WithdrawalsStorage: nil,
		DB:                 nil,
		Logger:             zap.NewExample().Sugar(),
	}

	m := MakeMux(d)

	m.ServeHTTP(httpW, req)

	respBody, _ := io.ReadAll(httpW.Body)
	assert.Equal(t, http.StatusInternalServerError, httpW.Code)
	assert.Equal(t, "Internal Server Error\n", string(respBody))
}

func TestCreateOrderSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	oStorage := mockstorage.NewMockOrdersStorage(ctrl)
	oStorage.
		EXPECT().
		CreateOrder(testutils.MatchContext(), gomock.Eq(OrderNumber), gomock.Eq(int64(UserID))).
		Return(nil)

	body := strings.NewReader(OrderNumber)
	req := httptest.NewRequest("POST", "/api/user/orders", body)
	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("Authorization", "Bearer "+JWT)
	httpW := httptest.NewRecorder()
	config.Set(
		config.Config{
			JWTSecret: JWTSecret,
		},
	)

	d := dependencies.D{
		UsersStorage:       nil,
		OrdersStorage:      oStorage,
		WithdrawalsStorage: nil,
		DB:                 nil,
		Logger:             zap.NewExample().Sugar(),
	}

	m := MakeMux(d)

	m.ServeHTTP(httpW, req)

	respBody, _ := io.ReadAll(httpW.Body)
	assert.Equal(t, http.StatusAccepted, httpW.Code)
	assert.Equal(t, "Order accepted", string(respBody))
}

func TestOrdersGetAllInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	oStorage := mockstorage.NewMockOrdersStorage(ctrl)
	oStorage.
		EXPECT().
		GetUserOrders(testutils.MatchContext(), gomock.Eq(int64(UserID))).
		Return([]models.Order{}, errors.New("internal server error"))

	req := httptest.NewRequest("GET", "/api/user/orders", nil)
	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("Authorization", "Bearer "+JWT)
	httpW := httptest.NewRecorder()
	config.Set(
		config.Config{
			JWTSecret: JWTSecret,
		},
	)

	d := dependencies.D{
		UsersStorage:       nil,
		OrdersStorage:      oStorage,
		WithdrawalsStorage: nil,
		DB:                 nil,
		Logger:             zap.NewExample().Sugar(),
	}

	m := MakeMux(d)

	m.ServeHTTP(httpW, req)

	respBody, _ := io.ReadAll(httpW.Body)
	assert.Equal(t, http.StatusInternalServerError, httpW.Code)
	assert.Equal(t, "Internal Server Error\n", string(respBody))
}

func TestOrdersGetAllZeroOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	oStorage := mockstorage.NewMockOrdersStorage(ctrl)
	oStorage.
		EXPECT().
		GetUserOrders(testutils.MatchContext(), gomock.Eq(int64(UserID))).
		Return([]models.Order{}, nil)

	req := httptest.NewRequest("GET", "/api/user/orders", nil)
	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("Authorization", "Bearer "+JWT)
	httpW := httptest.NewRecorder()
	config.Set(
		config.Config{
			JWTSecret: JWTSecret,
		},
	)

	d := dependencies.D{
		UsersStorage:       nil,
		OrdersStorage:      oStorage,
		WithdrawalsStorage: nil,
		DB:                 nil,
		Logger:             zap.NewExample().Sugar(),
	}

	m := MakeMux(d)

	m.ServeHTTP(httpW, req)

	assert.Equal(t, http.StatusNoContent, httpW.Code)
}

func TestOrdersGetAllSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderUpdateTime := "2023-08-31T19:35:43Z"
	orderUpdateTimeParsed, _ := time.Parse(time.RFC3339, orderUpdateTime)
	orders := []models.Order{
		{
			ID:         1,
			UserID:     UserID,
			Number:     "12345",
			Status:     "pending",
			UploadedAt: orderUpdateTimeParsed,
		},
		{
			ID:         1,
			UserID:     UserID,
			Number:     "12345678",
			Status:     "in_process",
			UploadedAt: orderUpdateTimeParsed,
		},
	}
	oStorage := mockstorage.NewMockOrdersStorage(ctrl)
	oStorage.
		EXPECT().
		GetUserOrders(testutils.MatchContext(), gomock.Eq(int64(UserID))).
		Return(orders, nil)

	req := httptest.NewRequest("GET", "/api/user/orders", nil)
	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("Authorization", "Bearer "+JWT)
	httpW := httptest.NewRecorder()
	config.Set(
		config.Config{
			JWTSecret: JWTSecret,
		},
	)

	d := dependencies.D{
		UsersStorage:       nil,
		OrdersStorage:      oStorage,
		WithdrawalsStorage: nil,
		DB:                 nil,
		Logger:             zap.NewExample().Sugar(),
	}

	m := MakeMux(d)

	m.ServeHTTP(httpW, req)

	respBody, _ := io.ReadAll(httpW.Body)
	assert.Equal(t, http.StatusOK, httpW.Code)
	assert.Equal(
		t, fmt.Sprintf(
			`[{"number":"12345","status":"pending","uploaded_at":"%s"},{"number":"12345678","status":"in_process","uploaded_at":"%s"}]`,
			orderUpdateTime, orderUpdateTime,
		)+"\n",
		string(respBody),
	)
}
