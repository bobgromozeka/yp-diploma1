package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/bobgromozeka/yp-diploma1/internal/app"
	"github.com/bobgromozeka/yp-diploma1/internal/server/handlers/user"
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
	application := app.New()
	server := makeServer(application)

	fmt.Println("Running server on " + c.RunAddress)
	if err := http.ListenAndServe(c.RunAddress, server); err != nil {
		//TODO logging
		log.Fatalln(err)
	}
}
