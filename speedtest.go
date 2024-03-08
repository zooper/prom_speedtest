package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/showwin/speedtest-go/speedtest"
)

var (
	// Define metrics for speed test results.
	downloadSpeed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "internet_download_speed_mbps",
		Help: "Current Internet download speed in Mbps",
	})
	uploadSpeed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "internet_upload_speed_mbps",
		Help: "Current Internet upload speed in Mbps",
	})
	latency = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "internet_latency_ms",
		Help: "Current Internet latency in milliseconds",
	}, []string{"host"})
)

func init() {
	// Register metrics with Prometheus's default registry.
	prometheus.MustRegister(downloadSpeed)
	prometheus.MustRegister(uploadSpeed)
	prometheus.MustRegister(latency)
}

func performSpeedTest() {
	var speedtestClient = speedtest.New()
	serverList, err := speedtestClient.FetchServers()
	if err != nil {
		fmt.Printf("Error fetching server list: %v\n", err)
		return
	}
	targets, err := serverList.FindServer([]int{})
	if err != nil {
		fmt.Printf("Error finding server: %v\n", err)
		return
	}

	for _, s := range targets {
		s.PingTest(nil)
		s.DownloadTest()
		s.UploadTest()

		// Set metrics values based on the speed test results.
		downloadSpeed.Set(s.DLSpeed)
		uploadSpeed.Set(s.ULSpeed)
		latency.WithLabelValues(s.Host).Set(float64(s.Latency.Milliseconds()))
		s.Context.Reset() // Reset counter after each test
		break             // Perform the test with the first server and stop
	}
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	performSpeedTest()                 // Perform speed test on each request to /metrics
	promhttp.Handler().ServeHTTP(w, r) // Serve the default Prometheus metrics handler
}

func main() {
	// Expose the custom metrics via HTTP.
	http.HandleFunc("/metrics", metricsHandler)
	http.ListenAndServe(":8080", nil)
}
