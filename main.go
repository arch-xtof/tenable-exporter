package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arch-xtof/tenable-exporter/client"
	"github.com/arch-xtof/tenable-exporter/collector"
	"github.com/arch-xtof/tenable-exporter/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getEnv(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

var (
	listenAddress = getEnv("LISTEN_ADDRESS", ":9095")
	metricsPath   = getEnv("METRICS_PATH", "/metrics")
	accessKey     = getEnv("TENABLE_ACCESS_KEY", "")
	secretKey     = getEnv("TENABLE_SECRET_KEY", "")
)

var (
	serverVulnCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tenable_server_vuln_count",
			Help: "Number of vulnerabilities on servers",
		},
		[]string{"server", "severity", "group", "region"},
	)
	workstationVulnCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tenable_workstation_vuln_count",
			Help: "Number of vulnerabilities on workstations",
		},
		[]string{"workstation", "severity", "os", "app"},
	)
)

func main() {
	c, err := client.NewTenableClient(&http.Client{}, client.Auth{AccessKey: accessKey, SecretKey: secretKey})
	if err != nil {
		log.Fatalf("Could not create Tenable Client: %v", err)
	}

	r := prometheus.NewRegistry()
	r.MustRegister(serverVulnCount, workstationVulnCount)

	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
	db := collector.SqlConnect("", "")
	log.Println("Starting Tenable Cardinality Exporter...")

	go func() {
		for {
			vulns, err := c.GetVulns("")
			if err != nil {
				log.Println(err)
			}
			assets, err := c.GetAssets("")
			if err != nil {
				log.Println(err)
			}

			assetMap := models.CreateAssetMap(assets)

			isv := c.GetServerVulnerabilityCount(vulns, assetMap)
			wsv := c.GetWorkstationVulnerabilityCount(vulns, assetMap)

			for i, sv := range isv {
				for s, v := range sv {
					serverVulnCount.WithLabelValues(i, s, v.Group, v.Region).Set(float64(v.Count))
				}
			}
			for w, sv := range wsv {
				for s, v := range sv {
					workstationVulnCount.WithLabelValues(w, s, v.OS, v.App).Set(float64(v.Count))
				}
			}

			collector.SqlDelete(db)
			collector.SqlInsert(db, vulns, assetMap)

			time.Sleep(6 * time.Hour)
		}
	}()

	/*r := prometheus.NewRegistry()

	client, err := client.NewTenableClient(&http.Client{}, client.Auth{AccessKey: accessKey, SecretKey: secretKey})
	if err != nil {
		log.Fatalf("Could not create Tenable Client: %v", err)
	}*/

	http.Handle(metricsPath, handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Tenable Prometheus Exporter</title></head>
			<body>
			<h1>Tenable Prometheus Exporter</h1>
			<p><a href='/metrics'>Metrics</a></p>
			</body>
			</html>`))
	})

	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
