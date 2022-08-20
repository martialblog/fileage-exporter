package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var build = "development"

type files []string

func (i *files) String() string {
	return fmt.Sprintf("%d", *i)
}

func (i *files) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	listenAddress = flag.String("web.listen-address", ":9123", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics.")
)

func main() {
	var filePatterns files
	flag.Var(&filePatterns, "file", "file to export")
	flag.Parse()

	// Register Collectors
	prometheus.MustRegister(NewFileAgeCollector(filePatterns))

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
			<head><title>File Age Exporter</title></head>
			<body>
			<h1>File Age Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Printf("Version: %s", build)
	log.Printf("Listening on address: %s", *listenAddress)

	log.Printf("File Patterns: %s", filePatterns)

	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
