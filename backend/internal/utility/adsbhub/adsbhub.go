package adsbhub

import (
	"adsb-api/internal/global"
	"adsb-api/internal/utility/converter"
	"bufio"
	"errors"
	"net"
	"strconv"
	"strings"
	"time"
)

func MakeTCPConnection(url string) (net.Conn, error) {
	return net.Dial("tcp", url)
}

func CloseTCPConnection(connection net.Conn) error {
	return connection.Close()
}

func ProcessSBSstream() ([]global.Aircraft, error) {
	conn, err := MakeTCPConnection("data.adsbhub.org:5002")
	if err != nil {
		return []global.Aircraft{}, err

	}

	defer CloseTCPConnection(conn)
	scanner := bufio.NewScanner(conn)
	var aircrafts []global.Aircraft

	for {
		timer := time.Now()
		scanner.Scan()
		if diff := time.Since(timer).Seconds(); diff > 4 {
			break
		}
		line := strings.Split(scanner.Text(), ",")

		if len(line) > 10 {
			icao := line[4]
			date := line[8]
			time := line[9]
			callsign := line[10]

			scanner.Scan()
			line = strings.Split(scanner.Text(), ",")
			altitudeStr := line[11]
			latStr := line[14]
			longStr := line[15]

			scanner.Scan()
			line = strings.Split(scanner.Text(), ",")
			speedStr := line[12]
			trackStr := line[13]
			vspeedStr := line[16]

			altitude, altERR := strconv.Atoi(altitudeStr)
			speed, spdERR := strconv.ParseFloat(speedStr, 32)
			track, trkERR := strconv.ParseFloat(trackStr, 32)
			vspeed, vspdERR := strconv.Atoi(vspeedStr)
			lat, latERR := strconv.ParseFloat(latStr, 32)
			long, longERR := strconv.ParseFloat(longStr, 32)

			if altERR != nil || spdERR != nil || trkERR != nil ||
				vspdERR != nil || latERR != nil || longERR != nil {
				continue
			}

			timestamp := converter.MakeTimeStamp(date, time)

			aircraft := global.Aircraft{
				Icao:         icao,
				Callsign:     callsign,
				Altitude:     altitude,
				Latitude:     float32(lat),
				Longitude:    float32(long),
				Speed:        int(speed),
				Track:        int(track),
				VerticalRate: vspeed,
				Timestamp:    timestamp,
			}

			aircrafts = append(aircrafts, aircraft)

		} else {
			return []global.Aircraft{}, errors.New("could not connect to stream")
		}

	}

	return aircrafts, nil
}
