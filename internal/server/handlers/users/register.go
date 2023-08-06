package users

import (
	"encoding/json"
	"net/http"

	"github.com/bobgromozeka/yp-diploma1/internal/app"
	"github.com/bobgromozeka/yp-diploma1/internal/constants"
	httphelpers "github.com/bobgromozeka/yp-diploma1/internal/http"
	"github.com/bobgromozeka/yp-diploma1/internal/jwt"
	"github.com/bobgromozeka/yp-diploma1/internal/server/config"
	"github.com/bobgromozeka/yp-diploma1/internal/server/requests"
	"github.com/bobgromozeka/yp-diploma1/internal/storage"
)

func Register(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !httphelpers.CheckContentType(w, r, httphelpers.ContentJSON) {
			return
		}

		reqPayload := requests.Register{}

		jd := json.NewDecoder(r.Body)
		if decodeErr := jd.Decode(&reqPayload); decodeErr != nil || reqPayload.Login == "" || reqPayload.Password == "" {
			app.Logger.Error(decodeErr)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		createUserErr := app.Storage.CreateUser(r.Context(), reqPayload.Login, reqPayload.Password)
		if createUserErr == storage.UserAlreadyExists {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		} else if createUserErr != nil {
			app.Logger.Error(createUserErr)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		ID, authErr := app.Storage.AuthUser(r.Context(), reqPayload.Login, reqPayload.Password)
		if authErr != nil {
			app.Logger.Error(authErr)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		token, jwtErr := jwt.MakeJWT(config.Get().JWTSecret, jwt.MakeJWTPayload(ID))
		if jwtErr != nil {
			app.Logger.Error(jwtErr)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set(constants.AuthorizationHeader, token)
		w.WriteHeader(http.StatusOK)
	}
}
