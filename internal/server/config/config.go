package config

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
	JWTSecret            string
}

var configuration Config

func Get() Config {
	return configuration
}

func Set(c Config) {
	configuration = c
}
