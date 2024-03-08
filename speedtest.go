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

	latency = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "internet_latency_ms",
		Help: "Current Internet latency in milliseconds",
	})
	server_id = "21016" // Starry Internet Speedtest server
)

func init() {
	// Register metrics with Prometheus's default registry.
	prometheus.MustRegister(downloadSpeed)
	prometheus.MustRegister(uploadSpeed)
	prometheus.MustRegister(latency)
}

func performSpeedTest() {
	var speedtestClient = speedtest.New()

	// Fetch the specific server by ID
	server, err := speedtestClient.FetchServerByID(server_id)
	if err != nil {
		fmt.Printf("Error fetching server by ID: %v\n", err)
		return
	}

	// Perform the speed test on the specific server
	server.PingTest(nil)
	server.DownloadTest()
	server.UploadTest()

	// Set metrics values based on the speed test results.
	downloadSpeed.Set(server.DLSpeed)
	uploadSpeed.Set(server.ULSpeed)
	latency.Set(float64(server.Latency.Milliseconds()))

	server.Context.Reset() // Reset counter after the test
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
