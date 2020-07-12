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
	MAC string

	RSSI int
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
		var times [34]int
		for i, idx := range []int{1, 2} {
			v, err := strconv.Atoi(fields[idx])
			if err != nil {
				return nil, err

			}
			times[i] = v
		}

		stats = append(stats, IWStat{
			MAC:  fields[0],
			RSSI: times[0],
		})
	}

	if err := s.Err(); err != nil {
		return nil, err

	}

	return stats, nil
}
