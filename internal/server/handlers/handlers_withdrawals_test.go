package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/bobgromozeka/yp-diploma1/internal/app/dependencies"
	"github.com/bobgromozeka/yp-diploma1/internal/models"
	"github.com/bobgromozeka/yp-diploma1/internal/server/config"
	mockstorage "github.com/bobgromozeka/yp-diploma1/internal/storage/mock"
	"github.com/bobgromozeka/yp-diploma1/internal/testutils"
)

func TestWithdrawalsInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var withdrawals []models.Withdrawal
	wStorage := mockstorage.NewMockWithdrawalsStorage(ctrl)
	wStorage.
		EXPECT().
		GetUserWithdrawals(testutils.MatchContext(), gomock.Eq(int64(UserID))).
		Return(withdrawals, errors.New("internal server error"))

	req := httptest.NewRequest("GET", "/api/user/withdrawals", nil)
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

	assert.Equal(t, http.StatusInternalServerError, httpW.Code)
}

func TestWithdrawalsNoContent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var withdrawals []models.Withdrawal
	wStorage := mockstorage.NewMockWithdrawalsStorage(ctrl)
	wStorage.
		EXPECT().
		GetUserWithdrawals(testutils.MatchContext(), gomock.Eq(int64(UserID))).
		Return(withdrawals, nil)

	req := httptest.NewRequest("GET", "/api/user/withdrawals", nil)
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

	assert.Equal(t, http.StatusNoContent, httpW.Code)
}

func TestWithdrawalsGetAllSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	processedAt := "2023-08-31T19:35:43Z"
	processedAtParsed, _ := time.Parse(time.RFC3339, processedAt)
	withdrawals := []models.Withdrawal{
		{
			ID:          1,
			UserID:      UserID,
			OrderNumber: "1234",
			Sum:         1234,
			ProcessedAt: processedAtParsed,
		},
		{
			ID:          2,
			UserID:      UserID,
			OrderNumber: "12345",
			Sum:         12345,
			ProcessedAt: processedAtParsed,
		},
	}
	wStorage := mockstorage.NewMockWithdrawalsStorage(ctrl)
	wStorage.
		EXPECT().
		GetUserWithdrawals(testutils.MatchContext(), gomock.Eq(int64(UserID))).
		Return(withdrawals, nil)

	req := httptest.NewRequest("GET", "/api/user/withdrawals", nil)
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
	assert.Equal(
		t, fmt.Sprintf(
			`[{"order":"1234","sum":1234,"processed_at":"%s"},{"order":"12345","sum":12345,"processed_at":"%s"}]`,
			processedAt, processedAt,
		)+"\n",
		string(respBody),
	)
}
