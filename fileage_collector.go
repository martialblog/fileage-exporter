package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type filestat struct {
	age_in_seconds float64
	size_in_bytes  float64
}

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

	files, err := getFilePaths(c.filePatterns)
	if err != nil {
		return err
	}

	// Number of files matching the given patterns
	ch <- prometheus.MustNewConstMetric(c.num_files_matching, prometheus.GaugeValue, float64(len(files)))

	// Get stats for all files
	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, filepath := range files {
		go func(filepath string) {
			fs, _ := getFileStats(filepath)
			// TODO Handle Errors
			ch <- prometheus.MustNewConstMetric(c.age_in_seconds, prometheus.CounterValue, fs.age_in_seconds, filepath)
			ch <- prometheus.MustNewConstMetric(c.size_in_bytes, prometheus.GaugeValue, fs.size_in_bytes, filepath)
			wg.Done()
		}(filepath)
	}

	wg.Wait()
	return nil
}

// Get and calculate stats from file
func getFileStats(filepath string) (filestat, error) {
	var fstat filestat

	finfo, err := os.Stat(filepath)
	if err != nil {
		return fstat, err
	}

	fstat.age_in_seconds = time.Since(finfo.ModTime()).Seconds()
	fstat.size_in_bytes = float64(finfo.Size())
	return fstat, nil
}

// Get list of files from Glob and append to single list
func getFilePaths(globpatterns []string) ([]string, error) {
	var files []string

	for _, v := range globpatterns {
		f, err := filepath.Glob(v)
		if err != nil {
			return files, err
		}
		files = append(files, f...)
	}
	return files, nil
}
