// Command cpustat provides basic Linux CPU utilization statistics.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rohan-prasad/iwstat"
)

func main() {
	f, err := os.Open("/Users/rp/iwstat")
	if err != nil {
		log.Fatal("Failed to open /Users/rp/iwstat: %v", err)
	}
	defer f.Close()

	stats, err := iwstat.Scan(f)
	if err != nil {
		log.Fatalf("Failed to scan: %v", err)

	}
	for _, s := range stats {
		fmt.Printf("%4s: RSSI: %06d",
			s.MAC, s.RSSI)

	}

}
