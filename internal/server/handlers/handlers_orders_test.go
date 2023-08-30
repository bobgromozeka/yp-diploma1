package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/bobgromozeka/yp-diploma1/internal/app/dependencies"
	"github.com/bobgromozeka/yp-diploma1/internal/server/config"
	"github.com/bobgromozeka/yp-diploma1/internal/storage"
	mockstorage "github.com/bobgromozeka/yp-diploma1/internal/storage/mock"
	"github.com/bobgromozeka/yp-diploma1/internal/testutils"
)

func TestCreateWrongOrderNumber(t *testing.T) {
	wrongOrderNumber := "123123"
	body := strings.NewReader(wrongOrderNumber) //Wrong order number
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
		CreateOrder(testutils.MatchContext(), gomock.Eq("4561261212345467"), gomock.Eq(int64(1))).
		Return(storage.ErrOrderAlreadyCreated)

	orderNumber := "4561261212345467"
	body := strings.NewReader(orderNumber) //Wrong order number
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
		CreateOrder(testutils.MatchContext(), gomock.Eq("4561261212345467"), gomock.Eq(int64(1))).
		Return(storage.ErrOrderForeign)

	orderNumber := "4561261212345467"
	body := strings.NewReader(orderNumber) //Wrong order number
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
		CreateOrder(testutils.MatchContext(), gomock.Eq("4561261212345467"), gomock.Eq(int64(1))).
		Return(errors.New("internal server error"))

	orderNumber := "4561261212345467"
	body := strings.NewReader(orderNumber) //Wrong order number
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
		CreateOrder(testutils.MatchContext(), gomock.Eq("4561261212345467"), gomock.Eq(int64(1))).
		Return(nil)

	orderNumber := "4561261212345467"
	body := strings.NewReader(orderNumber) //Wrong order number
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
