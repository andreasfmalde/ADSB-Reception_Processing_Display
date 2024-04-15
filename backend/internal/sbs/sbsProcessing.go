package sbs

import (
	"adsb-api/internal/global/models"
	"adsb-api/internal/utility/convert"
	"bufio"
	"net"
	"strings"
	"time"
)

// ProcessSbsStream reads aircraft data from the SBS stream, converts it to AircraftCurrentModel,
// and returns a slice of AircraftCurrentModel. It takes an address string and a waiting time in seconds
// as input parameters.
//
// The address string is used to establish a TCP connection with the SBS stream.
// The waiting time is the maximum duration for which the function will wait for new data from the stream.
//
// The function reads lines from the scanner until either the waiting time is exceeded or an error occurs.
//
// If any of the messages do not contain the expected number of fields, the line is skipped.
// The messages are then passed to the convert.SbsToAircraftCurrent function for conversion.
// If the conversion is successful, the converted aircraft data is appended to the aircraft slice.
// The final aircraft slice is returned as a result.
func ProcessSbsStream(addr string, waitingTime int) ([]models.AircraftCurrentModel, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	scanner := bufio.NewScanner(conn)

	var aircraft []models.AircraftCurrentModel

	timer := time.Now()
	for scanner.Scan() {
		if diff := time.Since(timer).Seconds(); diff > float64(waitingTime) {
			break
		}

		var ac models.AircraftCurrentModel

		msg1 := strings.Split(scanner.Text(), ",")
		if len(msg1) < 11 {
			continue
		}

		if !scanner.Scan() {
			break
		}
		msg3 := strings.Split(scanner.Text(), ",")
		if len(msg3) < 16 {
			continue
		}

		if !scanner.Scan() {
			break
		}
		msg4 := strings.Split(scanner.Text(), ",")
		if len(msg4) < 17 {
			continue
		}

		ac, err := convert.SbsToAircraftCurrent(msg1, msg3, msg4)
		if err != nil {
			continue
		}

		aircraft = append(aircraft, ac)
	}

	return aircraft, nil
}
