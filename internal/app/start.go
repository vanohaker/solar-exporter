package app

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	fiberlog "github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vanohaker/solar-exporter/internal/options"
	"github.com/vanohaker/solar-exporter/internal/solarmetrics"
)

func Run() {
	app := fiber.New(fiber.Config{
		Immutable: true,
	})
	prometheus_instance := solarmetrics.StartPrometheus()
	app.Use(
		logger.New(
			logger.Config{
				Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${queryParams} | ${error}\n",
			},
		),
		prometheus_instance.Middleware,
	)
	config, err := options.LoadConfig()
	if err != nil {
		fiberlog.Fatal(err.Error())
	}

	app.Get(config.Core.MetricsPath, adaptor.HTTPHandler(promhttp.Handler()))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("/")
	})

	app.Listen(fmt.Sprintf("%s:%v", config.Core.BindAddress, config.Core.BindPort))

	// var voltage = new(float64)
	// invertor.GetVoltages(serialsession, voltage)

	// go smartwatt.StartColllect(registry)

}
