package app

import (
	"fmt"
	"runtime"

	"github.com/gofiber/fiber/v2"
	fiberlog "github.com/gofiber/fiber/v2/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vanohaker/solar-exporter/internal/options"
)

func Run() {
	arch := fmt.Sprintf("%v/%v", runtime.GOOS, runtime.GOARCH)
	config, err := options.LoadConfig()
	if err != nil {
		fiberlog.Fatal(err.Error())
	}
	app := fiber.New(fiber.Config{
		Immutable: true,
	})
	fiberlog.Infof("SmartWatt ECO Prometheus Exporter, arch=%v, go=%v\n", arch, runtime.Version())

	registry := prometheus.NewRegistry()

	buildInfoMetric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "smartwatteco_build_info",
			Help: "SmartWatt ECO Exporter build information",
			ConstLabels: prometheus.Labels{
				"arch": arch,
				"go":   runtime.Version(),
			},
		},
	)
	buildInfoMetric.Set(1)
	registry.MustRegister(buildInfoMetric)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("/")
	})

	app.Listen(fmt.Sprintf("%s:%v", config.Core.BindAddress, config.Core.BindPort))

	// var voltage = new(float64)
	// invertor.GetVoltages(serialsession, voltage)

	// go smartwatt.StartColllect(registry)

}
