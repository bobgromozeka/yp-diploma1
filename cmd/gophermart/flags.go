package main

import (
	"flag"
	"os"

	"yp_diploma1/internal/server"
)

const (
	RunAddress           = "RUN_ADDRESS"
	DatabaseURI          = "DATABASE_URI"
	AccrualSystemAddress = "ACCRUAL_SYSTEM_ADDRESS"
)

func parseFlags(c *server.Config) {
	flag.StringVar(&c.RunAddress, "a", "localhost:8080", "server address and port")
	flag.StringVar(
		&c.DatabaseURI, "d", "",
		"Postgresql data source name (connection string like postgres://username:password@localhost:5432/database_name)",
	)
	flag.StringVar(&c.AccrualSystemAddress, "r", "", "Accrual system address")
}

func parseEnv(c *server.Config) {
	if runAddr, found := os.LookupEnv(RunAddress); found {
		c.RunAddress = runAddr
	}

	if databaseURI, found := os.LookupEnv(DatabaseURI); found {
		c.DatabaseURI = databaseURI
	}

	if accrualSystemAddress, found := os.LookupEnv(AccrualSystemAddress); found {
		c.AccrualSystemAddress = accrualSystemAddress
	}
}
