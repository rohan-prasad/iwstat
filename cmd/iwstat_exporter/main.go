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
		f, err := os.Open("/tmp/metrics")
		if err != nil {
			return nil, fmt.Errorf("Failed to open /tmp/metrics: %v", err)

		}
		defer f.Close()

		return iwstat.Scan(f)

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
	WifiStationSignal  *prometheus.Desc
	SnrStation         *prometheus.Desc
	Clientinactive     *prometheus.Desc
	RxPhy              *prometheus.Desc
	Rxbytes            *prometheus.Desc
	RxPrr              *prometheus.Desc
	RxPackets          *prometheus.Desc
	TxPhy              *prometheus.Desc
	Txbytes            *prometheus.Desc
	TxPrr              *prometheus.Desc
	TxPackets          *prometheus.Desc
	ExpectedThroughput *prometheus.Desc

	// A paremeterized funtion used to gather metrics.
	stats func() ([]iwstat.IWStat, error)
}

// newCollector constructs a collector using a stats function.
func newCollector(stats func() ([]iwstat.IWStat, error)) prometheus.Collector {
	return &collector{
		WifiStationSignal:  prometheus.NewDesc("iw_rssi_of_connected_client", "RSSI address of connected client", []string{"ifname", "mac"}, nil),
		SnrStation:         prometheus.NewDesc("iw_snr_of_connected_client", "SNR of connected client", []string{"ifname", "mac"}, nil),
		Clientinactive:     prometheus.NewDesc("iw_client_inactive_sec", "Client inactive in sec", []string{"ifname", "mac"}, nil),
		RxPhy:              prometheus.NewDesc("iw_rxPhy_of_connected_client", "rx Phy of connected client", []string{"ifname", "mac"}, nil),
		Rxbytes:            prometheus.NewDesc("iw_rxbytes_of_connected_client", "rx bytes of connected client", []string{"ifname", "mac"}, nil),
		RxPrr:              prometheus.NewDesc("iw_rxPrr_inactive_sec", "rx packet retry rate inactive in sec", []string{"ifname", "mac"}, nil),
		RxPackets:          prometheus.NewDesc("iw_rxPackets_connected_client", "Total rx packets address of connected client.", []string{"ifname", "mac"}, nil),
		TxPhy:              prometheus.NewDesc("iw_txPhy_of_connected_client", "tx Phy of connected client", []string{"ifname", "mac"}, nil),
		Txbytes:            prometheus.NewDesc("iw_txbytes_of_connected_client", "tx bytes inactive in sec", []string{"ifname", "mac"}, nil),
		TxPrr:              prometheus.NewDesc("iw_txPrr_of_connected_client", "tx Packet retry rate address of connected client.", []string{"ifname", "mac"}, nil),
		TxPackets:          prometheus.NewDesc("iw_txPackets_of_connected_client", "Total tx packets of connected client", []string{"ifname", "mac"}, nil),
		ExpectedThroughput: prometheus.NewDesc("iw_expectedThroughput_of_connected_client", "Expected Throughput of connected client", []string{"ifname", "mac"}, nil),
		stats:              stats,
	}
}

// Describe implements prometheus.Collector.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	// Gather metadata about each metrics.
	ds := []*prometheus.Desc{
		c.WifiStationSignal,
		c.SnrStation,
		c.Clientinactive,
		c.RxPhy,
		c.Rxbytes,
		c.RxPrr,
		c.RxPackets,
		c.TxPhy,
		c.Txbytes,
		c.TxPrr,
		c.TxPackets,
		c.ExpectedThroughput,
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
		ch <- prometheus.NewInvalidMetric(c.WifiStationSignal, err)
		ch <- prometheus.NewInvalidMetric(c.SnrStation, err)
		ch <- prometheus.NewInvalidMetric(c.Clientinactive, err)
		ch <- prometheus.NewInvalidMetric(c.RxPhy, err)
		ch <- prometheus.NewInvalidMetric(c.Rxbytes, err)
		ch <- prometheus.NewInvalidMetric(c.RxPrr, err)
		ch <- prometheus.NewInvalidMetric(c.RxPackets, err)
		ch <- prometheus.NewInvalidMetric(c.TxPhy, err)
		ch <- prometheus.NewInvalidMetric(c.Txbytes, err)
		ch <- prometheus.NewInvalidMetric(c.TxPrr, err)
		ch <- prometheus.NewInvalidMetric(c.TxPackets, err)
		ch <- prometheus.NewInvalidMetric(c.ExpectedThroughput, err)
		return
	}

	for _, s := range stats {
		tuples := []struct {
			mac string
			v   int
		}{
			{mac: s.Mac, v: s.Rssi},
		}
		for _, t := range tuples {
			// Prometheus.Collector implementation should always use
			// "const metric" constructors.
			ch <- prometheus.MustNewConstMetric(
				c.WifiStationSignal,
				prometheus.CounterValue,
				float64(t.v),
				s.Ifname, t.mac,
			)
		}
	}

	for _, s := range stats {
		tuples := []struct {
			mac string
			v   int
		}{
			{mac: s.Mac, v: s.Snr},
		}
		for _, t := range tuples {
			// Prometheus.Collector implementation should always use
			// "const metric" constructors.
			ch <- prometheus.MustNewConstMetric(
				c.SnrStation,
				prometheus.CounterValue,
				float64(t.v),
				s.Ifname, t.mac,
			)
		}
	}

	for _, s := range stats {
		tuples := []struct {
			mac string
			v   int
		}{
			{mac: s.Mac, v: s.ClientInactive},
		}
		for _, t := range tuples {
			// Prometheus.Collector implementation should always use
			// "const metric" constructors.
			ch <- prometheus.MustNewConstMetric(
				c.Clientinactive,
				prometheus.CounterValue,
				float64(t.v),
				s.Ifname, t.mac,
			)
		}
	}

	for _, s := range stats {
		tuples := []struct {
			mac string
			v   int
		}{
			{mac: s.Mac, v: s.RxPhy},
		}
		for _, t := range tuples {
			// Prometheus.Collector implementation should always use
			// "const metric" constructors.
			ch <- prometheus.MustNewConstMetric(
				c.RxPhy,
				prometheus.CounterValue,
				float64(t.v),
				s.Ifname, t.mac,
			)
		}
	}

	for _, s := range stats {
		tuples := []struct {
			mac string
			v   int
		}{
			{mac: s.Mac, v: s.RxMbytes},
		}
		for _, t := range tuples {
			// Prometheus.Collector implementation should always use
			// "const metric" constructors.
			ch <- prometheus.MustNewConstMetric(
				c.Rxbytes,
				prometheus.CounterValue,
				float64(t.v),
				s.Ifname, t.mac,
			)
		}
	}

	for _, s := range stats {
		tuples := []struct {
			mac string
			v   int
		}{
			{mac: s.Mac, v: s.RxPrr},
		}
		for _, t := range tuples {
			// Prometheus.Collector implementation should always use
			// "const metric" constructors.
			ch <- prometheus.MustNewConstMetric(
				c.RxPrr,
				prometheus.CounterValue,
				float64(t.v),
				s.Ifname, t.mac,
			)
		}
	}

	for _, s := range stats {
		tuples := []struct {
			mac string
			v   int
		}{
			{mac: s.Mac, v: s.RxPackets},
		}
		for _, t := range tuples {
			// Prometheus.Collector implementation should always use
			// "const metric" constructors.
			ch <- prometheus.MustNewConstMetric(
				c.RxPackets,
				prometheus.CounterValue,
				float64(t.v),
				s.Ifname, t.mac,
			)
		}
	}

	for _, s := range stats {
		tuples := []struct {
			mac string
			v   int
		}{
			{mac: s.Mac, v: s.TxPhy},
		}
		for _, t := range tuples {
			// Prometheus.Collector implementation should always use
			// "const metric" constructors.
			ch <- prometheus.MustNewConstMetric(
				c.TxPhy,
				prometheus.CounterValue,
				float64(t.v),
				s.Ifname, t.mac,
			)
		}
	}

	for _, s := range stats {
		tuples := []struct {
			mac string
			v   int
		}{
			{mac: s.Mac, v: s.TxMbytes},
		}
		for _, t := range tuples {
			// Prometheus.Collector implementation should always use
			// "const metric" constructors.
			ch <- prometheus.MustNewConstMetric(
				c.Txbytes,
				prometheus.CounterValue,
				float64(t.v),
				s.Ifname, t.mac,
			)
		}
	}

	for _, s := range stats {
		tuples := []struct {
			mac string
			v   int
		}{
			{mac: s.Mac, v: s.TxPrr},
		}
		for _, t := range tuples {
			// Prometheus.Collector implementation should always use
			// "const metric" constructors.
			ch <- prometheus.MustNewConstMetric(
				c.TxPrr,
				prometheus.CounterValue,
				float64(t.v),
				s.Ifname, t.mac,
			)
		}
	}

	for _, s := range stats {
		tuples := []struct {
			mac string
			v   int
		}{
			{mac: s.Mac, v: s.TxPackets},
		}
		for _, t := range tuples {
			// Prometheus.Collector implementation should always use
			// "const metric" constructors.
			ch <- prometheus.MustNewConstMetric(
				c.TxPackets,
				prometheus.CounterValue,
				float64(t.v),
				s.Ifname, t.mac,
			)
		}
	}

	for _, s := range stats {
		tuples := []struct {
			mac string
			v   int
		}{
			{mac: s.Mac, v: s.ExpectedThroughput},
		}
		for _, t := range tuples {
			// Prometheus.Collector implementation should always use
			// "const metric" constructors.
			ch <- prometheus.MustNewConstMetric(
				c.ExpectedThroughput,
				prometheus.CounterValue,
				float64(t.v),
				s.Ifname, t.mac,
			)
		}
	}

}
