package handlers

import (
	"errors"
	"fmt"
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

func TestBalanceWithdrawWrongOrderNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	body := strings.NewReader(fmt.Sprintf(`{"order":"%s","sum":111}`, WrongOrderNumber))
	req := httptest.NewRequest("POST", "/api/user/balance/withdraw", body)
	req.Header.Add("Content-Type", "application/json")
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

	m := MakeMux(d)

	m.ServeHTTP(httpW, req)

	respBody, _ := io.ReadAll(httpW.Body)
	assert.Equal(t, http.StatusUnprocessableEntity, httpW.Code)
	assert.Equal(t, "Wrong order format\n", string(respBody))
}

func TestBalanceWithdrawInsufficientFunds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wStorage := mockstorage.NewMockWithdrawalsStorage(ctrl)
	wStorage.
		EXPECT().
		Withdraw(testutils.MatchContext(), gomock.Eq(int64(UserID)), OrderNumber, float64(111)).
		Return(storage.ErrInsufficientFunds)

	body := strings.NewReader(fmt.Sprintf(`{"order":"%s","sum":111}`, OrderNumber))
	req := httptest.NewRequest("POST", "/api/user/balance/withdraw", body)
	req.Header.Add("Content-Type", "application/json")
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
		WithdrawalsStorage: wStorage,
		DB:                 nil,
		Logger:             zap.NewExample().Sugar(),
	}

	m := MakeMux(d)

	m.ServeHTTP(httpW, req)

	respBody, _ := io.ReadAll(httpW.Body)
	assert.Equal(t, http.StatusPaymentRequired, httpW.Code)
	assert.Equal(t, "Insufficient funds\n", string(respBody))
}

func TestBalanceWithdrawInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wStorage := mockstorage.NewMockWithdrawalsStorage(ctrl)
	wStorage.
		EXPECT().
		Withdraw(testutils.MatchContext(), gomock.Eq(int64(UserID)), OrderNumber, float64(111)).
		Return(errors.New("internal server error"))

	body := strings.NewReader(fmt.Sprintf(`{"order":"%s","sum":111}`, OrderNumber))
	req := httptest.NewRequest("POST", "/api/user/balance/withdraw", body)
	req.Header.Add("Content-Type", "application/json")
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
		WithdrawalsStorage: wStorage,
		DB:                 nil,
		Logger:             zap.NewExample().Sugar(),
	}

	m := MakeMux(d)

	m.ServeHTTP(httpW, req)

	respBody, _ := io.ReadAll(httpW.Body)
	assert.Equal(t, http.StatusInternalServerError, httpW.Code)
	assert.Equal(t, "Internal Server Error\n", string(respBody))
}

func TestBalanceWithdrawSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wStorage := mockstorage.NewMockWithdrawalsStorage(ctrl)
	wStorage.
		EXPECT().
		Withdraw(testutils.MatchContext(), gomock.Eq(int64(UserID)), OrderNumber, float64(111)).
		Return(nil)

	body := strings.NewReader(fmt.Sprintf(`{"order":"%s","sum":111}`, OrderNumber))
	req := httptest.NewRequest("POST", "/api/user/balance/withdraw", body)
	req.Header.Add("Content-Type", "application/json")
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
		WithdrawalsStorage: wStorage,
		DB:                 nil,
		Logger:             zap.NewExample().Sugar(),
	}

	m := MakeMux(d)

	m.ServeHTTP(httpW, req)

	assert.Equal(t, http.StatusOK, httpW.Code)
}

func TestBalanceGetInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wStorage := mockstorage.NewMockWithdrawalsStorage(ctrl)
	wStorage.
		EXPECT().
		GetUserBalance(testutils.MatchContext(), gomock.Eq(int64(UserID))).
		Return(float64(0), float64(0), errors.New("internal server error"))

	req := httptest.NewRequest("GET", "/api/user/balance", nil)
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
		WithdrawalsStorage: wStorage,
		DB:                 nil,
		Logger:             zap.NewExample().Sugar(),
	}

	m := MakeMux(d)

	m.ServeHTTP(httpW, req)

	respBody, _ := io.ReadAll(httpW.Body)
	assert.Equal(t, http.StatusInternalServerError, httpW.Code)
	assert.Equal(t, "Internal Server Error\n", string(respBody))
}

func TestBalanceGetSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	balance := float64(1000)
	withdrawalsSum := float64(5000)
	wStorage := mockstorage.NewMockWithdrawalsStorage(ctrl)
	wStorage.
		EXPECT().
		GetUserBalance(testutils.MatchContext(), gomock.Eq(int64(UserID))).
		Return(balance, withdrawalsSum, nil)

	req := httptest.NewRequest("GET", "/api/user/balance", nil)
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
		WithdrawalsStorage: wStorage,
		DB:                 nil,
		Logger:             zap.NewExample().Sugar(),
	}

	m := MakeMux(d)

	m.ServeHTTP(httpW, req)

	respBody, _ := io.ReadAll(httpW.Body)
	assert.Equal(t, http.StatusOK, httpW.Code)
	assert.Equal(t, fmt.Sprintf(`{"current":%.0f,"withdrawn":%.0f}`, balance, withdrawalsSum)+"\n", string(respBody))
}
