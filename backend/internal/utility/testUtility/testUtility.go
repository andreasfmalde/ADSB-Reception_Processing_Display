package testUtility

import (
	"adsb-api/internal/global"
	"strconv"
	"time"
)

func CreateMockAircraft(n int) []global.AircraftCurrentModel {
	var aircraft []global.AircraftCurrentModel

	for i := 0; i < n; i++ {
		ac := global.AircraftCurrentModel{
			Icao:         strconv.Itoa(i),
			Callsign:     strconv.Itoa(i),
			Altitude:     i,
			Latitude:     float32(i),
			Longitude:    float32(i),
			Speed:        i,
			Track:        i,
			VerticalRate: i,
			Timestamp:    time.Now().Format(time.DateTime),
		}
		aircraft = append(aircraft, ac)
	}

	return aircraft
}

func CreateMockHistAircraft(n int) []global.AircraftHistoryModel {
	var aircraft []global.AircraftHistoryModel

	for i := 0; i < n; i++ {
		ac := global.AircraftHistoryModel{
			Icao:      strconv.Itoa(i),
			Latitude:  float32(i),
			Longitude: float32(i),
			Timestamp: time.Now().Format(time.DateTime),
		}
		aircraft = append(aircraft, ac)
	}

	return aircraft
}

func CreateMockAircraftWithTimestamp(icao string, timestamp string) global.AircraftCurrentModel {
	return global.AircraftCurrentModel{
		Icao:         icao,
		Callsign:     "TEST",
		Altitude:     10000,
		Latitude:     51.5074,
		Longitude:    0.1278,
		Speed:        450,
		Track:        180,
		VerticalRate: 0,
		Timestamp:    timestamp,
	}
}

func CreateMockAircraftWithIcao(n int, icao string) []global.AircraftCurrentModel {
	var aircraft []global.AircraftCurrentModel

	for i := 0; i < n; i++ {
		ac := global.AircraftCurrentModel{
			Icao:         icao,
			Callsign:     strconv.Itoa(i),
			Altitude:     i,
			Latitude:     float32(i),
			Longitude:    float32(i),
			Speed:        i,
			Track:        i,
			VerticalRate: i,
			Timestamp:    time.Now().Add(time.Duration(i) * time.Second).Format(time.DateTime),
		}
		aircraft = append(aircraft, ac)
	}

	return aircraft
}
