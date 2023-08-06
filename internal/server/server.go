package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"

	"github.com/bobgromozeka/yp-diploma1/internal/app"
	"github.com/bobgromozeka/yp-diploma1/internal/db"
	"github.com/bobgromozeka/yp-diploma1/internal/jwt"
	"github.com/bobgromozeka/yp-diploma1/internal/log"
	"github.com/bobgromozeka/yp-diploma1/internal/server/config"
	"github.com/bobgromozeka/yp-diploma1/internal/server/handlers/orders"
	"github.com/bobgromozeka/yp-diploma1/internal/server/handlers/users"
	"github.com/bobgromozeka/yp-diploma1/internal/storage"
)

func makeServer(app app.App) *chi.Mux {
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
						"/register", users.Register(app),
					)
					r.Post(
						"/login", users.Login(app),
					)
					r.Group(
						func(r chi.Router) {
							r.Use(jwtauth.Verifier(jwt.GetTokenAuth(config.Get().JWTSecret)))
							r.Use(jwtauth.Authenticator)

							r.Route(
								"/orders", func(r chi.Router) {
									r.Get(
										"/", func(w http.ResponseWriter, r *http.Request) {

										},
									)
									r.Post("/", orders.Create(app))
								},
							)
						},
					)
				},
			)
		},
	)

	return r
}

func Run() {
	logger, loggerError := log.New()
	if loggerError != nil {
		fmt.Println(loggerError)
		os.Exit(1)
	}
	defer logger.Sync()

	connErr := db.Connect(config.Get().DatabaseURI)
	if connErr != nil {
		logger.Fatalln(connErr)
	}

	pgStorage := storage.NewPGStorage(db.Connection())

	application := app.New(pgStorage, db.Connection(), logger)
	server := makeServer(application)

	bootstrapError := storage.Bootstrap(db.Connection())
	if bootstrapError != nil {
		logger.Fatalln(bootstrapError)
	}

	fmt.Println("Running server on " + config.Get().RunAddress)
	if err := http.ListenAndServe(config.Get().RunAddress, server); err != nil {
		logger.Fatalln(err)
	}
}
