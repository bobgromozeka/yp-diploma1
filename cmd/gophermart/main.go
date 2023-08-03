package main

import (
	"github.com/bobgromozeka/yp-diploma1/internal/server"
)

func main() {
	config := server.Config{}

	parseFlags(&config)
	parseEnv(&config)

	server.Run(config)
}
