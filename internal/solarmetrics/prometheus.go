package solarmetrics

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/vanohaker/solar-exporter/pkg/smartwatt"
)

type PrometheusInstance struct {
	buildinfo           prometheus.Gauge
	serialnumber        *prometheus.GaugeVec
	requestsTotal       *prometheus.CounterVec
	inputVoltageAC      *prometheus.GaugeVec
	inputVoltageACfrq   *prometheus.GaugeVec
	outputVoltageAC     *prometheus.GaugeVec
	outputVoltageACfrq  *prometheus.GaugeVec
	outputApparentPower *prometheus.GaugeVec
	outputActivePower   *prometheus.GaugeVec
	loadPercent         *prometheus.GaugeVec
	batareVoltage       *prometheus.GaugeVec
	bataryChargePercent *prometheus.GaugeVec
	bataryChargeCurrent *prometheus.GaugeVec
	invertorTemperature *prometheus.GaugeVec
	solarPanelCurrent   *prometheus.GaugeVec
	solarPanelVoltage   *prometheus.GaugeVec
}

var Prometheus *PrometheusInstance

func StartPrometheus() *PrometheusInstance {

	build_info := promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "solar_invertor_build_info",
			Help: "SmartWatt ECO Exporter build information",
			ConstLabels: prometheus.Labels{
				"arch": fmt.Sprintf("%v/%v", runtime.GOOS, runtime.GOARCH),
				"go":   runtime.Version(),
			},
		},
	)
	constLabels := prometheus.Labels{
		"app": "smartwatteco-exporter",
	}

	// this metric will store all requests sent to the server. Use this to get the rate of requests per minute or second
	requests_total_counter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:        prometheus.BuildFQName("solar_invertor", "api", "requests_total"),
			Help:        "The total number of HTTP requests made",
			ConstLabels: constLabels,
		},
		[]string{"status_code", "method", "path"},
	)

	invertor_serialnumber := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solar_invertor_serial_number",
			Help:        "Inverter serial number",
			ConstLabels: constLabels,
		},
		[]string{"serialnumber"},
	)

	// Входное перемиенное напряжение инвертора
	invertor_input_voltage_ac := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solar_invertor_input_voltage_ac",
			Help:        "Input voltage from city line",
			ConstLabels: constLabels,
		},
		[]string{"port", "serialnumber"},
	)

	// Частота входящего напряжения сети
	invertor_input_voltage_frequency := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solar_invertor_input_voltage_ac_frequency",
			Help:        "Frequency of voltage supplied to the inverter from the city network",
			ConstLabels: constLabels,
		},
		[]string{"port", "serialnumber"},
	)

	invertor_output_voltage_ac := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solar_invertor_output_voltage_ac",
			Help:        "The magnitude of the output voltage towards the load",
			ConstLabels: constLabels,
		},
		[]string{"port", "serialnumber"},
	)

	invertor_output_voltage_ac_frequency := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solar_invertor_output_voltage_ac_frequency",
			Help:        "Inverter output voltage frequency",
			ConstLabels: constLabels,
		},
		[]string{"port", "serialnumber"},
	)

	invertor_output_apparent_power := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solar_invertor_output_apparent_power",
			Help:        "Total integrated total power",
			ConstLabels: constLabels,
		},
		[]string{"port", "serialnumber"},
	)

	invertor_output_active_power := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solar_invertor_output_active_power",
			Help:        "Active power",
			ConstLabels: constLabels,
		},
		[]string{"port", "serialnumber"},
	)

	invertor_load_percent := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solar_invertor_load_percent",
			Help:        "Load percentage per inverter",
			ConstLabels: constLabels,
		},
		[]string{"port", "serialnumber"},
	)

	invertor_batary_voltage := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solar_invertor_batary_voltage",
			Help:        "Inverter battery voltage",
			ConstLabels: constLabels,
		},
		[]string{"port", "serialnumber"},
	)

	invertor_batary_charge_current := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solar_invertor_charge_current",
			Help:        "Battery charge current",
			ConstLabels: constLabels,
		},
		[]string{"port", "serialnumber"},
	)

	invertor_batary_charge_percent := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solar_invertor_batary_charge_percent",
			Help:        "Battery percentage",
			ConstLabels: constLabels,
		},
		[]string{"port", "serialnumber"},
	)

	invertor_temperature := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solar_invertor_temperature",
			Help:        "Inverter temperature",
			ConstLabels: constLabels,
		},
		[]string{"port", "serialnumber"},
	)

	invertor_solar_panel_current := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solar_invertor_solar_panel_current",
			Help:        "Solar panel current",
			ConstLabels: constLabels,
		},
		[]string{"port", "serialnumber"},
	)

	invertor_solar_panel_voltage := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "solar_invertor_solar_panel_voltage",
			Help:        "Solar Panel Voltage",
			ConstLabels: constLabels,
		},
		[]string{"port", "channel", "serialnumber"},
	)

	return &PrometheusInstance{
		buildinfo:           build_info,
		serialnumber:        invertor_serialnumber,
		requestsTotal:       requests_total_counter,
		inputVoltageAC:      invertor_input_voltage_ac,
		inputVoltageACfrq:   invertor_input_voltage_frequency,
		outputVoltageAC:     invertor_output_voltage_ac,
		outputVoltageACfrq:  invertor_output_voltage_ac_frequency,
		outputApparentPower: invertor_output_apparent_power,
		outputActivePower:   invertor_output_active_power,
		loadPercent:         invertor_load_percent,
		batareVoltage:       invertor_batary_voltage,
		bataryChargeCurrent: invertor_batary_charge_current,
		bataryChargePercent: invertor_batary_charge_percent,
		invertorTemperature: invertor_temperature,
		solarPanelCurrent:   invertor_solar_panel_current,
		solarPanelVoltage:   invertor_solar_panel_voltage,
	}
}

// this is the middleware that will run request-based metrics
func (p *PrometheusInstance) Middleware(c *fiber.Ctx) error {
	method := c.Route().Method
	path := c.Route().Path
	port := c.Queries()["port"]
	baudrate, _ := strconv.Atoi(c.Queries()["baudrate"]) // error TODO
	if port == "" || baudrate == 0 {
		return c.SendString("port and baidrate parametrs are required!")
	}
	invertor, _ := smartwatt.Init(smartwatt.SmartWattEco{
		Port:     port,
		BaudRate: baudrate,
	})

	invertor.GetSerialNumber()
	serialnumber := invertor.SerialNumber
	p.serialnumber.WithLabelValues(serialnumber)

	// Получаем данные о напряжениях и частотах
	invertor.GetVoltage()
	p.inputVoltageAC.WithLabelValues(port, serialnumber).Set(invertor.InputACvoltage)
	p.inputVoltageACfrq.WithLabelValues(port, serialnumber).Set(invertor.InputACfrq)
	p.outputVoltageAC.WithLabelValues(port, serialnumber).Set(invertor.OutACvoltage)
	p.outputVoltageACfrq.WithLabelValues(port, serialnumber).Set(invertor.OutACfrq)
	p.outputApparentPower.WithLabelValues(port, serialnumber).Set(invertor.OutApparentPower)
	p.outputActivePower.WithLabelValues(port, serialnumber).Set(invertor.OutActivePower)
	p.loadPercent.WithLabelValues(port, serialnumber).Set(invertor.LoadPercent)
	p.batareVoltage.WithLabelValues(port, serialnumber).Set(invertor.BatVoltage)
	p.bataryChargeCurrent.WithLabelValues(port, serialnumber).Set(invertor.ChargeCurrent)
	p.bataryChargePercent.WithLabelValues(port, serialnumber).Set(invertor.ChargePercent)
	p.invertorTemperature.WithLabelValues(port, serialnumber).Set(invertor.InvertorTemp)
	p.solarPanelCurrent.WithLabelValues(port, serialnumber).Set(invertor.ChargeSolarCurrent)
	p.solarPanelVoltage.WithLabelValues(port, "ch1", serialnumber).Set(invertor.VoltageDCch1)

	var status int

	err := c.Next()
	if err != nil {
		if e, ok := err.(*fiber.Error); ok {
			// Get correct error code from fiber.Error type
			status = e.Code
		}
	} else {
		status = c.Response().StatusCode()
	}

	statusCodeString := strconv.Itoa(status)
	p.requestsTotal.WithLabelValues(statusCodeString, method, path).Inc()

	return err
}
