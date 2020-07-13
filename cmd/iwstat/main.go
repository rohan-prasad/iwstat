package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rohan-prasad/iwstat"
)

func main() {
	f, err := os.Open("/home/rohan/iwstat")
	if err != nil {
		log.Fatalf("Failed to open /home/rohan/iwstat: %v", err)
	}
	defer f.Close()

	stats, err := iwstat.Scan(f)
	if err != nil {
		log.Fatalf("Failed to scan: %v", err)
	}

	for _, s := range stats {
		fmt.Printf("%4s RSSI:%6d\n",
			s.MAC, s.RSSI)
	}

}
