package main

import (
	"flag"
	"os"

	"github.com/bobgromozeka/yp-diploma1/internal/server/config"
)

const (
	RunAddress           = "RUN_ADDRESS"
	DatabaseURI          = "DATABASE_URI"
	AccrualSystemAddress = "ACCRUAL_SYSTEM_ADDRESS"
	JWTSecret            = "JWT_SECRET"
)

func parseFlags(c *config.Config) {
	flag.StringVar(&c.RunAddress, "a", "localhost:8080", "server address and port")
	flag.StringVar(
		&c.DatabaseURI, "d", "",
		"Postgresql data source name (connection string like postgres://practicum:practicum@localhost:5432/practicum)",
	)
	flag.StringVar(&c.AccrualSystemAddress, "r", "", "Accrual system address")
	flag.StringVar(&c.JWTSecret, "j", "secret", "JWT Secret key")

	flag.Parse()
}

func parseEnv(c *config.Config) {
	if runAddr, found := os.LookupEnv(RunAddress); found {
		c.RunAddress = runAddr
	}

	if databaseURI, found := os.LookupEnv(DatabaseURI); found {
		c.DatabaseURI = databaseURI
	}

	if accrualSystemAddress, found := os.LookupEnv(AccrualSystemAddress); found {
		c.AccrualSystemAddress = accrualSystemAddress
	}

	if jwt, found := os.LookupEnv(JWTSecret); found {
		c.JWTSecret = jwt
	}
}
