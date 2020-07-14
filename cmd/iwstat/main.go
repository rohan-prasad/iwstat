package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rohan-prasad/iwstat"
)

func main() {
	f, err := os.Open("/tmp/iwstat")
	if err != nil {
		log.Fatalf("Failed to open /tmp/iwstat: %v", err)
	}
	defer f.Close()

	stats, err := iwstat.Scan(f)
	if err != nil {
		log.Fatalf("Failed to scan: %v", err)
	}

	for _, s := range stats {
		fmt.Printf("%4s rssi:%6d\n, snr: %6d/n, clientInacrive: %6d/n, rxPhy: %6d/n",
			s.Mac, s.Rssi, s.Snr, s.ClientInactive, s.RxPhy)
	}

}
