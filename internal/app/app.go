package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bobgromozeka/yp-diploma1/internal/accrual"
	"github.com/bobgromozeka/yp-diploma1/internal/app/dependencies"
	"github.com/bobgromozeka/yp-diploma1/internal/db"
	"github.com/bobgromozeka/yp-diploma1/internal/log"
	"github.com/bobgromozeka/yp-diploma1/internal/server"
	"github.com/bobgromozeka/yp-diploma1/internal/server/config"
	"github.com/bobgromozeka/yp-diploma1/internal/storage"
)

func Start(c config.Config) {
	ctx := context.Background()
	shutdownCtx, closeShutdownCtx := context.WithCancel(ctx)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	deps := makeDependencies(c)

	bootstrapError := storage.Bootstrap(db.Connection())
	if bootstrapError != nil {
		deps.Logger.Fatalln(bootstrapError)
	}

	go func() {
		<-sig
		deps.Logger.Info("Stopping application.....")
		closeShutdownCtx()
	}()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go server.Run(shutdownCtx, deps, &wg)

	wg.Add(1)
	go accrual.Run(shutdownCtx, deps, &wg)

	wg.Wait()
}

func makeDependencies(c config.Config) dependencies.D {
	logger, loggerError := log.New()
	if loggerError != nil {
		fmt.Println(loggerError)
		os.Exit(1)
	}
	defer logger.Sync()

	connErr := db.Connect(c.DatabaseURI)
	if connErr != nil {
		logger.Fatalln(connErr)
	}

	pgUsersStorage := storage.NewPgUsersStorage(db.Connection())
	pgOrdersStorage := storage.NewPgOrdersStorage(db.Connection())
	pgWithdrawalsStorage := storage.NewPgWithdrawalsStorage(db.Connection())

	return dependencies.D{
		UsersStorage:       pgUsersStorage,
		OrdersStorage:      pgOrdersStorage,
		WithdrawalsStorage: pgWithdrawalsStorage,
		DB:                 db.Connection(),
		Logger:             logger,
	}
}
