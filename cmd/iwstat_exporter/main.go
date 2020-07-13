package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rohan-prasad/iwstat"
)

func main() {
	//called on each collector.collect
	stats := func() ([]iwstat.IWStat, error) {
		f, err := os.Open("/home/rohan/iwstat")
		if err != nil {
			return nil, fmt.Errorf("Failed to open /home/rohan/iwstat: %v", err)

		}
		defer f.Close()

		return iwstat.Scan()

	}

	//Make Prometheus client aware of new collector.
	c := newCollector(stats)
	prometheus.MustRegister(c)

	//Setup HTTP handler for metrics.
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	//Starts listening for HTTP connection.
	const addr = ":9999"
	log.Printf("starting iwstat exporter on %q", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("cannot start iwstat exporter: %s", err)
	}
}

var _ prometheus.Collector = &collector{}

// A colector is a prometheus.Collector for OpenWRT iwinfo stats.
type collector struct {
	// Metrics description.
	ClientConnectionStat *prometheus.Desc

	// A paremeterized funtion used to gather metrics.
	stats func() ([]iwstat.IWStat, error)
}

// newCollector constructs a collector using a stats function.
func newCollector(stats func() ([]iwstat.IWStat, error)) prometheus.Collector {
	return &collector{
		ClientConnectionStat: prometheus.NewDesc(
			//Name of the metrics.
			"mac_address_of_connected_client",
			// The metrics help text.
			"MAC address of connected client.",
			// The metric's variable label dimensions.
			[]string{"client"},
			// The metric's constant label dimensions.
			nil,
		),

		stats: stats,
	}
}

// Describe implements prometheus.Collector.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	// Gather metadata about each metrics.
	ds := []*prometheus.Desc{
		c.ClientConnectionStat,
	}

	for _, d := range ds {
		ch <- d
	}
}

// Collect implements prometheus.Collector.
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	// Take a stat snapshot. Must be concorreny safe,
	stats, err := c.stats()
	if err != nil {
		// If an error occurs, send an invalid metrics to notify
		// Prometheus of the problem.
		ch <- prometheus.NewInvalidMetric(c.ClientConnectionStat, err)
		return
	}

	for _, s := range stats {
		tuples := []struct {
			client string
			v      int
		}{
			{client: "rssi", v: s.RSSI},
		}
		for _, t := range tuples {
			// Prometheus.Collector implementation should always use
			// "const metric" constructors.
			ch <- prometheus.MustNewConstMetric(
				c.ClientConnectionStat,
				prometheus.CounterValue,
				float64(t.v),
				s.MAC, t.client,
			)
		}
	}

}
