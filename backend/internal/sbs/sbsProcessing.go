package sbs

import (
	"adsb-api/internal/global/models"
	"adsb-api/internal/utility/convert"
	"bufio"
	"net"
	"strings"
	"time"
)

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
	for {
		if diff := time.Since(timer).Seconds(); diff > float64(waitingTime) {
			break
		}

		var ac models.AircraftCurrentModel

		if !scanner.Scan() {
			break
		}
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
