package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bobgromozeka/yp-diploma1/internal/app/dependencies"
	"github.com/bobgromozeka/yp-diploma1/internal/server/config"
	"github.com/bobgromozeka/yp-diploma1/internal/server/handlers"
)

func Run(shutdownCtx context.Context, d dependencies.D) {
	server := http.Server{Addr: config.Get().RunAddress, Handler: handlers.MakeMux(d)}

	//graceful shutdown
	go func() {
		<-shutdownCtx.Done()

		forceCtx, cancelForceCtx := context.WithTimeout(context.Background(), time.Second*30)
		go func() {
			<-forceCtx.Done()
			d.Logger.Fatal("Shutdown deadline is exceeded. Forcing exit")
		}()

		d.Logger.Info("Shutting down server.....")
		err := server.Shutdown(forceCtx)
		if err != nil {
			d.Logger.Fatal(err)
		}
		cancelForceCtx()
	}()

	fmt.Println("Running server on " + config.Get().RunAddress)
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			d.Logger.Fatalln(err)
		}
	}
}
