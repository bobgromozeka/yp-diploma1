package main

import (
	"github.com/bobgromozeka/yp-diploma1/internal/server"
	"github.com/bobgromozeka/yp-diploma1/internal/server/config"
)

func main() {
	c := config.Get()

	parseFlags(&c)
	parseEnv(&c)

	config.Set(c)

	server.Run()
}
