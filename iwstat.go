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
	Ifname, Mac string

	Rssi, Snr, ClientInactive, RxPhy, RxMbytes, RxPrr, RxVhtMcsIndex, RxVhtMcsMhz, RxVhtNss, RxPackets, TxPhy, TxMbytes, TxPrr, TxVhtMcsIndex, TxVhtMcsMhz, TxVhtNss, TxPackets, ExpectedThroughput, ChannelUtlization int
}

//Scan reads and parses iwinfo
func Scan(r io.Reader) ([]IWStat, error) {

	s := bufio.NewScanner(r)
	s.Scan()

	var stats []IWStat
	for s.Scan() {

		const nFields = 21
		fields := strings.Fields(string(s.Bytes()))
		if len(fields) != nFields {
			continue
		}
		var times [19]int
		for i, idx := range []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21} {
			v, err := strconv.Atoi(fields[idx])
			if err != nil {
				return nil, err

			}
			times[i] = v
		}

		stats = append(stats, IWStat{
			Ifname:             fields[0],
			Mac:                fields[1],
			Rssi:               times[0],
			Snr:                times[1],
			ClientInactive:     times[2],
			RxPhy:              times[3],
			RxMbytes:           times[4],
			RxPrr:              times[5],
			RxVhtMcsIndex:      times[6],
			RxVhtMcsMhz:        times[7],
			RxVhtNss:           times[8],
			RxPackets:          times[9],
			TxPhy:              times[10],
			TxMbytes:           times[11],
			TxPrr:              times[12],
			TxVhtMcsIndex:      times[13],
			TxVhtMcsMhz:        times[14],
			TxVhtNss:           times[15],
			TxPackets:          times[16],
			ExpectedThroughput: times[17],
			ChannelUtlization:  times[18],
		})
	}

	if err := s.Err(); err != nil {
		return nil, err

	}

	return stats, nil
}
