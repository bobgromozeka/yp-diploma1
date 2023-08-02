package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"yp_diploma1/internal/app"
)

func makeServer(app app.App) *chi.Mux {
	server := chi.NewMux()

	server.Use(
		middleware.Heartbeat("/health"),
		middleware.StripSlashes,
		middleware.Recoverer,
	)

	return server
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
