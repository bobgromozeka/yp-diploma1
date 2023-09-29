package main

import (
	"github.com/bobgromozeka/yp-diploma1/internal/app"
	"github.com/bobgromozeka/yp-diploma1/internal/server/config"
)

func main() {
	c := config.Get()

	parseFlags(&c)
	parseEnv(&c)

	config.Set(c)

	app.Start(c)
}
