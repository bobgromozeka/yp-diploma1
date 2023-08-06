package middlewares

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth/v5"

	"github.com/bobgromozeka/yp-diploma1/internal/app"
)

const ContextUserKey = "USER_KEY"

func SetUserFromJWT(app app.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, claims, err := jwtauth.FromContext(r.Context())
				if err != nil {
					app.Logger.Errorw("SetUserFromJWT", "error", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				ctx := context.WithValue(r.Context(), ContextUserKey, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}
