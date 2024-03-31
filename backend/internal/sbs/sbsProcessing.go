package sbs

import (
	"adsb-api/internal/global"
	"adsb-api/internal/global/models"
	"adsb-api/internal/utility/converter"
	"bufio"
	"net"
	"strconv"
	"strings"
	"time"
)

func ProcessSbsStream() ([]models.AircraftCurrentModel, error) {
	conn, err := net.Dial("tcp", global.SbsSource)
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
		if diff := time.Since(timer).Seconds(); diff > global.WaitingTime {
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

		ac, err := parseSbsToAircraftCurrent(msg1, msg3, msg4)
		if err != nil {
			continue
		}

		aircraft = append(aircraft, ac)
	}

	return aircraft, nil
}

func parseSbsToAircraftCurrent(msg1 []string, msg3 []string, msg4 []string) (models.AircraftCurrentModel, error) {
	icao := msg1[4]
	date := msg1[8]
	hour := msg1[9]
	callsign := msg1[10]
	timestamp := converter.MakeTimeStamp(date, hour)

	altitudeStr := msg3[11]
	latStr := msg3[14]
	longStr := msg3[15]

	altitude, err := strconv.Atoi(altitudeStr)
	if err != nil {
		return models.AircraftCurrentModel{}, err
	}
	lat, err := strconv.ParseFloat(latStr, 32)
	if err != nil {
		return models.AircraftCurrentModel{}, err
	}
	long, err := strconv.ParseFloat(longStr, 32)
	if err != nil {
		return models.AircraftCurrentModel{}, err
	}

	speedStr := msg4[12]
	trackStr := msg4[13]
	vspeedStr := msg4[16]

	speed, err := strconv.ParseFloat(speedStr, 32)
	if err != nil {
		return models.AircraftCurrentModel{}, err
	}
	track, err := strconv.ParseFloat(trackStr, 32)
	if err != nil {
		return models.AircraftCurrentModel{}, err
	}
	vspeed, err := strconv.Atoi(vspeedStr)
	if err != nil {
		return models.AircraftCurrentModel{}, err
	}

	return models.AircraftCurrentModel{
		Icao:         icao,
		Callsign:     callsign,
		Altitude:     altitude,
		Latitude:     float32(lat),
		Longitude:    float32(long),
		Speed:        int(speed),
		Track:        int(track),
		VerticalRate: vspeed,
		Timestamp:    timestamp,
	}, nil
}
