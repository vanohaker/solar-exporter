package app

import (
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
	"github.com/vanohaker/solar-exporter/internal/args"
	"github.com/vanohaker/solar-exporter/pkg/smartwatt"
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

func Run() {
	flag.Parse()

	commitHash, commitTime, dirtyBuild := getBuildInfo()
	arch := fmt.Sprintf("%v/%v", runtime.GOOS, runtime.GOARCH)

	log.Printf("SmartWatt ECO Prometheus Exporter version=%v commit=%v date=%v, dirty=%v, arch=%v, go=%v\n", args.Version, commitHash, commitTime, dirtyBuild, arch, runtime.Version())

	if *args.DisplayVersion {
		os.Exit(0)
	}

	log.Println("Starting prometheus metrics builder")

	registry := prometheus.NewRegistry()

	buildInfoMetric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "solarexporter_build_info",
			Help: "Exporter build information",
			ConstLabels: prometheus.Labels{
				"version": args.Version,
				"commit":  commitHash,
				"date":    commitTime,
				"dirty":   strconv.FormatBool(dirtyBuild),
				"arch":    arch,
				"go":      runtime.Version(),
			},
		},
	)
	buildInfoMetric.Set(1)
	registry.MustRegister(buildInfoMetric)

	log.Println("Starting serial port communications")

	// var voltage = new(float64)
	// invertor.GetVoltages(serialsession, voltage)

	go smartwatt.StartColllect(registry)

	srv := http.Server{
		ReadHeaderTimeout: 5 * time.Second,
	}

	http.Handle(*args.MetricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, `<!DOCTYPE html>
			<title>SmartWatt ECO Exporter</title>
			<h1>SmartWatt ECO Exporter</h1>
			<p><a href=%q>Metrics</a></p>`,
			*args.MetricsPath)
		if err != nil {
			log.Printf("Error while sending a response for the '/' path: %v", err)
		}
	})

	listener, err := getListener(*args.ListenAddr)
	if err != nil {
		log.Fatalf("Could not create listener: %v", err)
	}

	log.Println("SmartWatt ECO exporter started")
	log.Fatal(srv.Serve(listener))
}
