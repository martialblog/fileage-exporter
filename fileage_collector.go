package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type fileAgeCollector struct {
	filePatterns       []string
	num_files_matching *prometheus.Desc
	age_in_seconds     *prometheus.Desc
	size_in_bytes      *prometheus.Desc
}

func NewFileAgeCollector(filePatterns []string) *fileAgeCollector {
	return &fileAgeCollector{
		filePatterns:       filePatterns,
		num_files_matching: prometheus.NewDesc("file_age_num_files_matching", "Number of files matching glob patterns", nil, nil),
		age_in_seconds:     prometheus.NewDesc("file_age_age_in_seconds", "File age in seconds", []string{"path"}, nil),
		size_in_bytes:      prometheus.NewDesc("file_age_size_in_bytes", "File size in bytes", []string{"path"}, nil),
	}
}

func (c *fileAgeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.num_files_matching
}

func (c *fileAgeCollector) Collect(ch chan<- prometheus.Metric) {
	err := c.collect(ch)
	if err != nil {
		// TODO Handle Error
		log.Fatal(err)
	}
}

func (c *fileAgeCollector) collect(ch chan<- prometheus.Metric) error {

	var files []string

	// Get list of file from Glob and append to single list
	for _, v := range c.filePatterns {
		f, err := filepath.Glob(v)
		if err != nil {
			// TODO Handle Error
			log.Fatal(err)
		}
		files = append(files, f...)
	}

	// Number of files matching the given patterns
	ch <- prometheus.MustNewConstMetric(c.num_files_matching, prometheus.GaugeValue, float64(len(files)))

	// Get stats for all files
	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, filepath := range files {
		go func(filepath string) {
			finfo, err := os.Stat(filepath)
			if err != nil {
				// TODO Handle Error
				log.Fatal(err)
			}

			fileAgeInSeconds := time.Since(finfo.ModTime()).Seconds()
			fileSizeInBytes := float64(finfo.Size())
			ch <- prometheus.MustNewConstMetric(c.age_in_seconds, prometheus.CounterValue, fileAgeInSeconds, filepath)
			ch <- prometheus.MustNewConstMetric(c.size_in_bytes, prometheus.GaugeValue, fileSizeInBytes, filepath)
			wg.Done()
		}(filepath)
	}

	wg.Wait()
	return nil
}
