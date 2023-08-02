package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func makeServer() *chi.Mux {
	server := chi.NewMux()

	server.Get(
		"/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)

	return server
}

func Run(c Config) {
	server := makeServer()

	fmt.Println("Running server on " + c.RunAddress)
	if err := http.ListenAndServe(c.RunAddress, server); err != nil {
		//TODO logging
		log.Fatalln(err)
	}
}
