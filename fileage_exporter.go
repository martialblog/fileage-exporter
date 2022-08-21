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

type globpatterns []string

func (gp *globpatterns) String() string {
	return fmt.Sprintf("%s", *gp)
}

func (gp *globpatterns) Set(value string) error {
	*gp = append(*gp, value)
	return nil
}

var (
	listenAddress = flag.String("web.listen-address", ":9123", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics.")
)

func main() {
	var filePatterns globpatterns
	flag.Var(&filePatterns, "file", "Files to include. Can be used multiple times with Golang glob patterns.")
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

	if len(filePatterns) == 0 {
		log.Printf("Warning: No File Pattern provided")
	} else {
		log.Printf("File Patterns: %s", filePatterns)
	}

	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
