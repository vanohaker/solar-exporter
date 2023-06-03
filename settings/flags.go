package settings

import (
	"flag"
	"os"
	"strconv"
)

func getEnv(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

var (
	Version string

	// default vars values
	defaultSerialPortName        = getEnv("SERIAL_PORT", "/dev/ttyS0")
	defaultSerialPortBaudRate, _ = strconv.Atoi(getEnv("SERIAL_PORT_BAUD_RATE", "9600"))
	defaultMetricsPath           = getEnv("METRICS_PATH", "/metrics")
	defaultListenAddr            = getEnv("LISTEN_ADDRESS", "0.0.0.0:9678")

	SerialPortName     = flag.String("serialPortName", defaultSerialPortName, "Serial port name of the connected inverter")
	SerialPortBaudRate = flag.Int("serialPortBaudRate", defaultSerialPortBaudRate, "Serial port speed")
	DisplayVersion     = flag.Bool("version", false, "Display SmartWatt ECO exporter version")
	DebugMode          = flag.Bool("debug", false, "Debug mode")
	MetricsPath        = flag.String("metricsPath", defaultMetricsPath, "Url with metrics data")
	ListenAddr         = flag.String("webListenAddr", defaultListenAddr, "TCP port where exporter started")
)
