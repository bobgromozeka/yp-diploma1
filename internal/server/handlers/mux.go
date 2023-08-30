package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"

	"github.com/bobgromozeka/yp-diploma1/internal/app/dependencies"
	"github.com/bobgromozeka/yp-diploma1/internal/jwt"
	"github.com/bobgromozeka/yp-diploma1/internal/server/config"
	"github.com/bobgromozeka/yp-diploma1/internal/server/handlers/balance"
	"github.com/bobgromozeka/yp-diploma1/internal/server/handlers/orders"
	"github.com/bobgromozeka/yp-diploma1/internal/server/handlers/users"
	"github.com/bobgromozeka/yp-diploma1/internal/server/handlers/withdrawals"
)

func MakeMux(d dependencies.D) *chi.Mux {
	r := chi.NewMux()

	r.Use(
		middleware.StripSlashes,
		middleware.Logger,
		middleware.Recoverer,
	)

	r.Route(
		"/api", func(r chi.Router) {
			r.Use(middleware.Heartbeat("/health"))

			r.Route(
				"/user", func(r chi.Router) {
					r.Post(
						"/register", users.Register(d),
					)

					r.Post(
						"/login", users.Login(d),
					)

					r.Group(
						func(r chi.Router) {
							r.Use(jwtauth.Verifier(jwt.GetTokenAuth(config.Get().JWTSecret)))
							r.Use(jwtauth.Authenticator)

							r.Route(
								"/orders", func(r chi.Router) {
									r.Get("/", orders.GetAll(d))
									r.Post("/", orders.Create(d))
								},
							)

							r.Route(
								"/balance", func(r chi.Router) {
									r.Get(
										"/", balance.Get(d),
									)
									r.Post(
										"/withdraw", balance.Withdraw(d),
									)
								},
							)

							r.Get("/withdrawals", withdrawals.GetAll(d))
						},
					)
				},
			)
		},
	)

	return r
}
