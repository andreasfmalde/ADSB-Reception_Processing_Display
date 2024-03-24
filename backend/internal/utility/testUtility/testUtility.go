package testUtility

import (
	models2 "adsb-api/internal/global/models"
	"strconv"
	"time"
)

func CreateMockAircraft(n int) []models2.AircraftCurrentModel {
	var aircraft []models2.AircraftCurrentModel

	for i := 0; i < n; i++ {
		ac := models2.AircraftCurrentModel{
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

func CreateMockHistAircraft(n int) []models2.AircraftHistoryModel {
	var aircraft []models2.AircraftHistoryModel

	for i := 0; i < n; i++ {
		ac := models2.AircraftHistoryModel{
			Icao:      strconv.Itoa(i),
			Latitude:  float32(i),
			Longitude: float32(i),
			Timestamp: time.Now().Format(time.DateTime),
		}
		aircraft = append(aircraft, ac)
	}

	return aircraft
}

func CreateMockAircraftWithTimestamp(icao string, timestamp string) models2.AircraftCurrentModel {
	return models2.AircraftCurrentModel{
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

func CreateMockAircraftWithIcao(n int, icao string) []models2.AircraftCurrentModel {
	var aircraft []models2.AircraftCurrentModel

	for i := 0; i < n; i++ {
		ac := models2.AircraftCurrentModel{
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
