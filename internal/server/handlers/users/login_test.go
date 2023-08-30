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
	"github.com/bobgromozeka/yp-diploma1/internal/jwt"
	"github.com/bobgromozeka/yp-diploma1/internal/server/config"
	"github.com/bobgromozeka/yp-diploma1/internal/storage"
	mockstorage "github.com/bobgromozeka/yp-diploma1/internal/storage/mock"
	"github.com/bobgromozeka/yp-diploma1/internal/testutils"
)

func TestLoginBadRequestWrongJSON(t *testing.T) {
	type testCase struct {
		Name         string
		Body         string
		Status       int
		ResponseBody string
	}

	testCases := []testCase{
		{
			Name:         "Wrong JSON",
			Body:         "{bad body}",
			Status:       http.StatusBadRequest,
			ResponseBody: "Bad request\n",
		},
		{
			Name:         "Empty login",
			Body:         `{"login":"","password":"password"}`,
			Status:       http.StatusBadRequest,
			ResponseBody: "Bad request\n",
		},
		{
			Name:         "Empty password",
			Body:         `{"login":"login","password":""}`,
			Status:       http.StatusBadRequest,
			ResponseBody: "Bad request\n",
		},
		{
			Name:         "No login",
			Body:         `{"password":"password"}`,
			Status:       http.StatusBadRequest,
			ResponseBody: "Bad request\n",
		},
		{
			Name:         "No password",
			Body:         `{"login":"login"}`,
			Status:       http.StatusBadRequest,
			ResponseBody: "Bad request\n",
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.Name, func(t *testing.T) {
				body := strings.NewReader(tc.Body)
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

				h := Login(d)

				h.ServeHTTP(httpW, req)
				responseBody, _ := io.ReadAll(httpW.Body)

				assert.Equal(t, tc.Status, httpW.Code)
				assert.Equal(t, tc.ResponseBody, string(responseBody))
			},
		)
	}
}

func TestLoginInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uStorage := mockstorage.NewMockUsersStorage(ctrl)
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

	h := Login(d)

	h.ServeHTTP(httpW, req)

	responseBody, _ := io.ReadAll(httpW.Body)

	assert.Equal(t, http.StatusInternalServerError, httpW.Code)
	assert.Equal(t, "Internal server error\n", string(responseBody))
}

func TestLoginUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uStorage := mockstorage.NewMockUsersStorage(ctrl)
	uStorage.
		EXPECT().
		AuthUser(testutils.MatchContext(), gomock.Eq("login"), gomock.Eq("password")).
		Return(int64(0), storage.ErrUserNotFound)

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

	h := Login(d)

	h.ServeHTTP(httpW, req)

	responseBody, _ := io.ReadAll(httpW.Body)

	assert.Equal(t, http.StatusUnauthorized, httpW.Code)
	assert.Equal(t, "Wrong login or password\n", string(responseBody))
}

func TestLoginSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uStorage := mockstorage.NewMockUsersStorage(ctrl)
	uStorage.
		EXPECT().
		AuthUser(testutils.MatchContext(), gomock.Eq("login"), gomock.Eq("password")).
		Return(int64(1), nil)
	JWTSecret := "secret"
	config.Set(
		config.Config{
			JWTSecret: JWTSecret,
		},
	)
	expectedJWT, _ := jwt.MakeJWT(JWTSecret, jwt.MakeJWTPayload(1))

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

	h := Login(d)

	h.ServeHTTP(httpW, req)

	assert.Equal(t, http.StatusOK, httpW.Code)
	assert.Equal(t, "Bearer "+expectedJWT, httpW.Header().Get("Authorization"))
}
