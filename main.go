package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tarm/serial"
)

func getBuildInfo() (string, string, bool) {
	var commitHash, commitTime string
	var dirtyBuild bool
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", "", false
	}
	for _, kv := range info.Settings {
		switch kv.Key {
		case "vcs.revision":
			commitHash = kv.Value
		case "vcs.time":
			commitTime = kv.Value
		case "vcs.modified":
			dirtyBuild = kv.Value == "true"
		}
	}
	return commitHash, commitTime, dirtyBuild
}

func getEnv(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

func initSerialPort(portname string, br int) (*serial.Port, error) {
	serialPort := &serial.Config{Name: portname, Baud: br}
	session, err := serial.OpenPort(serialPort)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return session, nil
}

func writeSerialData(data []byte, session *serial.Port) (err error) {
	_, err = session.Write(data)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func readSerialData(session *serial.Port, channel chan []byte) {
	buf := make([]byte, 32)
	var readCount int
	message := []byte{}
	for {
		data, err := session.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		readCount++
		message = append(message, buf[:data]...)
		if bytes.Equal(buf[data-1:data], []byte{0x0d}) {
			// fmt.Printf("message : %s\n", string(message))
			channel <- message
			message = []byte{}
			// close(channel)
		}
	}
}

func parseUnixSocketAddress(address string) (string, string, error) {
	addressParts := strings.Split(address, ":")
	addressPartsLength := len(addressParts)

	if addressPartsLength > 3 || addressPartsLength < 1 {
		return "", "", fmt.Errorf("address for unix domain socket has wrong format")
	}

	unixSocketPath := addressParts[1]
	requestPath := ""
	if addressPartsLength == 3 {
		requestPath = addressParts[2]
	}
	return unixSocketPath, requestPath, nil
}

func getListener(listenAddress string) (net.Listener, error) {
	var listener net.Listener
	var err error

	if strings.HasPrefix(listenAddress, "unix:") {
		path, _, pathError := parseUnixSocketAddress(listenAddress)
		if pathError != nil {
			return listener, fmt.Errorf("parsing unix domain socket listen address %s failed: %w", listenAddress, pathError)
		}
		listener, err = net.ListenUnix("unix", &net.UnixAddr{Name: path, Net: "unix"})
	} else {
		listener, err = net.Listen("tcp", listenAddress)
	}

	if err != nil {
		return listener, err
	}
	log.Printf("Listening on %s", listenAddress)
	return listener, nil
}

var (
	version string

	// default vars values
	defaultSerialPortName        = getEnv("SERIAL_PORT", "/dev/ttyS0")
	defaultSerialPortBaudRate, _ = strconv.Atoi(getEnv("SERIAL_PORT_BAUD_RATE", "19200"))
	defaultMetricsPath           = getEnv("METRICS_PATH", "/metrics")
	defaultListenAddr            = getEnv("LISTEN_ADDRESS", "0.0.0.0:9678")

	serialPortName     = flag.String("serialPortName", defaultSerialPortName, "Serial port name of the connected inverter")
	serialPortBaudRate = flag.Int("serialPortBaudRate", defaultSerialPortBaudRate, "Serial port speed")
	displayVersion     = flag.Bool("version", false, "Display SmartWatt ECO exporter version")
	metricsPath        = flag.String("metricsPath", defaultMetricsPath, "Url with metrics data")
	listenAddr         = flag.String("webListenAddr", defaultListenAddr, "TCP port where exporter started")
)

func main() {
	flag.Parse()

	commitHash, commitTime, dirtyBuild := getBuildInfo()
	arch := fmt.Sprintf("%v/%v", runtime.GOOS, runtime.GOARCH)

	log.Printf("SmartWatt ECO Prometheus Exporter version=%v commit=%v date=%v, dirty=%v, arch=%v, go=%v\n", version, commitHash, commitTime, dirtyBuild, arch, runtime.Version())

	if *displayVersion {
		os.Exit(0)
	}

	log.Println("Starting...")
	registry := prometheus.NewRegistry()

	// buildInfoMetric := prometheus.NewGauge(
	// 	prometheus.GaugeOpts{
	// 		Name: "nginxexporter_build_info",
	// 		Help: "Exporter build information",
	// 	},
	// )
	// buildInfoMetric.Set(1)

	// registry.MustRegister(buildInfoMetric)

	srv := http.Server{
		ReadHeaderTimeout: 5 * time.Second,
	}

	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, `<!DOCTYPE html>
			<title>SmartWatt ECO Exporter</title>
			<h1>SmartWatt ECO Exporter</h1>
			<p><a href=%q>Metrics</a></p>`,
			*metricsPath)
		if err != nil {
			log.Printf("Error while sending a response for the '/' path: %v", err)
		}
	})

	listener, err := getListener(*listenAddr)
	if err != nil {
		log.Fatalf("Could not create listener: %v", err)
	}

	log.Println("SmartWatt ECO exporter started")
	log.Fatal(srv.Serve(listener))

	// initCommands := []byte{0x51, 0x50, 0x49, 0xBE, 0xAC, 0x0D}
	getVoltages := []byte{0x51, 0x50, 0x49, 0x52, 0x49, 0xF8, 0x54, 0x0D}

	channel := make(chan []byte)
	session, _ := initSerialPort(*serialPortName, *serialPortBaudRate)
	go readSerialData(session, channel)

	go func() {
		for {
			_ = writeSerialData(getVoltages, session)
			for {
				val, ok := <-channel
				if ok {
					time.Sleep(4 * time.Second)
					log.Printf("%s\n", string(val[1:len(val)-1]))
					break
				}
			}
		}
	}()

}
