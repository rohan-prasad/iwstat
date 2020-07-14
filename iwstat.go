//Package iwstat provides a parser for OpenWRT iwinfo stats
package iwstat

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

//IWStat statistics of associated clients
type IWStat struct {
	ifname             string
	mac                string
	rssi               int
	snr                int
	clientInactive     int
	rxPhy              int
	rxMbytes           int
	rxPrr              int
	rxVhtMcsIndex      int
	rxVhtMcsMhz        int
	rxVhtNss           int
	rxPackets          int
	txPhy              int
	txMbytes           int
	txPrr              int
	txVhtMcsIndex      int
	txVhtMcsMhz        int
	txVhtNss           int
	txPackets          int
	expectedThroughput int
	channelUtlization  int
}

//Scan reads and parses iwinfo
func Scan(r io.Reader) ([]IWStat, error) {

	s := bufio.NewScanner(r)
	s.Scan()

	var stats []IWStat
	for s.Scan() {

		const nFields = 20
		fields := strings.Fields(string(s.Bytes()))
		if len(fields) != nFields {
			continue
		}
		var times [20]int
		for i, idx := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 13, 14, 15, 16, 17, 18, 19, 20} {
			v, err := strconv.Atoi(fields[idx])
			if err != nil {
				return nil, err

			}
			times[i] = v
		}

		stats = append(stats, IWStat{
			ifname:             fields[0],
			mac:                fields[1],
			rssi:               times[0],
			snr:                times[1],
			clientInactive:     times[2],
			rxPhy:              times[3],
			rxMbytes:           times[4],
			rxPrr:              times[5],
			rxVhtMcsIndex:      times[6],
			rxVhtMcsMhz:        times[7],
			rxVhtNss:           times[8],
			rxPackets:          times[9],
			txPhy:              times[10],
			txMbytes:           times[11],
			txPrr:              times[12],
			txVhtMcsIndex:      times[13],
			txVhtMcsMhz:        times[14],
			txVhtNss:           times[15],
			txPackets:          times[16],
			expectedThroughput: times[17],
			channelUtlization:  times[18],
		})
	}

	if err := s.Err(); err != nil {
		return nil, err

	}

	return stats, nil
}
