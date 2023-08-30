package users

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

func TestRegisterNewUserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uStorage := mockstorage.NewMockUsersStorage(ctrl)
	uStorage.
		EXPECT().
		CreateUser(testutils.MatchContext(), gomock.Eq("login"), gomock.Eq("password")).
		Return(storage.ErrUserAlreadyExists)

	body := strings.NewReader(`{"login":"login","password":"password"}`)
	req := httptest.NewRequest("post", "/api/user/register", body)
	req.Header.Add("Content-Type", "application/json")
	httpW := httptest.NewRecorder()

	logger := zap.NewExample().Sugar()

	d := dependencies.D{
		UsersStorage:       uStorage,
		OrdersStorage:      nil,
		WithdrawalsStorage: nil,
		DB:                 nil,
		Logger:             logger,
	}

	h := Register(d)

	h.ServeHTTP(httpW, req)

	responseBody, _ := io.ReadAll(httpW.Body)

	assert.Equal(t, http.StatusConflict, httpW.Code)
	assert.Equal(t, "User already exists\n", string(responseBody))
}

func TestRegisterNewUserBadRequestBody(t *testing.T) {
	body := strings.NewReader(`{bad body}`)
	req := httptest.NewRequest("post", "/api/user/register", body)
	req.Header.Add("Content-Type", "application/json")
	httpW := httptest.NewRecorder()

	logger := zap.NewExample().Sugar()

	d := dependencies.D{
		UsersStorage:       nil,
		OrdersStorage:      nil,
		WithdrawalsStorage: nil,
		DB:                 nil,
		Logger:             logger,
	}

	h := Register(d)

	h.ServeHTTP(httpW, req)
	responseBody, _ := io.ReadAll(httpW.Body)

	assert.Equal(t, http.StatusBadRequest, httpW.Code)
	assert.Equal(t, "Bad request\n", string(responseBody))
}

func TestRegisterUserInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uStorage := mockstorage.NewMockUsersStorage(ctrl)
	uStorage.
		EXPECT().
		CreateUser(testutils.MatchContext(), gomock.Eq("login"), gomock.Eq("password")).
		Return(errors.New("internal server error"))

	body := strings.NewReader(`{"login":"login","password":"password"}`)
	req := httptest.NewRequest("post", "/api/user/register", body)
	req.Header.Add("Content-Type", "application/json")
	httpW := httptest.NewRecorder()

	logger := zap.NewExample().Sugar()

	d := dependencies.D{
		UsersStorage:       uStorage,
		OrdersStorage:      nil,
		WithdrawalsStorage: nil,
		DB:                 nil,
		Logger:             logger,
	}

	h := Register(d)

	h.ServeHTTP(httpW, req)

	responseBody, _ := io.ReadAll(httpW.Body)

	assert.Equal(t, http.StatusInternalServerError, httpW.Code)
	assert.Equal(t, "Internal server error\n", string(responseBody))
}

func TestRegisterUserInternalServerErrorOnLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uStorage := mockstorage.NewMockUsersStorage(ctrl)
	uStorage.
		EXPECT().
		CreateUser(testutils.MatchContext(), gomock.Eq("login"), gomock.Eq("password")).
		Return(nil)
	uStorage.
		EXPECT().
		AuthUser(testutils.MatchContext(), gomock.Eq("login"), gomock.Eq("password")).
		Return(int64(0), errors.New("internal server error"))

	body := strings.NewReader(`{"login":"login","password":"password"}`)
	req := httptest.NewRequest("post", "/api/user/register", body)
	req.Header.Add("Content-Type", "application/json")
	httpW := httptest.NewRecorder()

	logger := zap.NewExample().Sugar()

	d := dependencies.D{
		UsersStorage:       uStorage,
		OrdersStorage:      nil,
		WithdrawalsStorage: nil,
		DB:                 nil,
		Logger:             logger,
	}

	h := Register(d)

	h.ServeHTTP(httpW, req)

	responseBody, _ := io.ReadAll(httpW.Body)

	assert.Equal(t, http.StatusInternalServerError, httpW.Code)
	assert.Equal(t, "Internal server error\n", string(responseBody))
}

func TestUserRegisterSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uStorage := mockstorage.NewMockUsersStorage(ctrl)
	uStorage.
		EXPECT().
		CreateUser(testutils.MatchContext(), gomock.Eq("login"), gomock.Eq("password")).
		Return(nil)
	uStorage.
		EXPECT().
		AuthUser(testutils.MatchContext(), gomock.Eq("login"), gomock.Eq("password")).
		Return(int64(1), nil)
	config.Set(
		config.Config{
			JWTSecret: "secret",
		},
	)

	body := strings.NewReader(`{"login":"login","password":"password"}`)
	req := httptest.NewRequest("post", "/api/user/register", body)
	req.Header.Add("Content-Type", "application/json")
	httpW := httptest.NewRecorder()

	logger := zap.NewExample().Sugar()

	d := dependencies.D{
		UsersStorage:       uStorage,
		OrdersStorage:      nil,
		WithdrawalsStorage: nil,
		DB:                 nil,
		Logger:             logger,
	}

	h := Register(d)

	h.ServeHTTP(httpW, req)

	assert.Equal(t, http.StatusOK, httpW.Code)
	assert.Contains(t, httpW.Header().Get("Authorization"), "Bearer")
}
