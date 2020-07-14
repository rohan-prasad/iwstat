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
		fmt.Printf("%4s, %4s, rssi:%d, snr: %d, clientInacrive: %d, rxPhy: %d, rxMbytes: %d, rxPrr: %d, rxPackets: %d, txPhy: %d, tx_Mbytes: %d, txPrr: %d, txpackets: %d, expectedThroughput: %d\n",
			s.Ifname, s.Mac, s.Rssi, s.Snr, s.ClientInactive, s.RxPhy, s.RxMbytes, s.RxPrr, s.RxPackets, s.TxPhy, s.TxMbytes, s.TxPrr, s.TxPackets, s.ExpectedThroughput)
	}

}
