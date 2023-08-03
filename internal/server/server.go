package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/bobgromozeka/yp-diploma1/internal/app"
	"github.com/bobgromozeka/yp-diploma1/internal/db"
	"github.com/bobgromozeka/yp-diploma1/internal/server/handlers/user"
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
						"/register", user.Register(app),
					)
					r.Post(
						"/login", user.Login(app),
					)
				},
			)
		},
	)

	return r
}

func Run(c Config) {
	connErr := db.Connect(c.DatabaseURI)
	if connErr != nil {
		log.Fatalln(connErr)
	}

	pgStorage := storage.NewPGStorage(db.Connection())
	application := app.New(pgStorage, db.Connection())
	server := makeServer(application)

	bootstrapError := storage.Bootstrap(db.Connection())
	if bootstrapError != nil {
		log.Fatalln(bootstrapError)
	}

	fmt.Println("Running server on " + c.RunAddress)
	if err := http.ListenAndServe(c.RunAddress, server); err != nil {
		//TODO logging
		log.Fatalln(err)
	}
}
