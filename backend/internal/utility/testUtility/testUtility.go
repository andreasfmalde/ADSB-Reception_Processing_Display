package testUtility

import (
	"adsb-api/internal/db/models"
	"strconv"
	"time"
)

func CreateMockAircraft(n int) []models.AircraftCurrentModel {
	var aircraft []models.AircraftCurrentModel

	for i := 0; i < n; i++ {
		ac := models.AircraftCurrentModel{
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

func CreateMockHistAircraft(n int) []models.AircraftHistoryModel {
	var aircraft []models.AircraftHistoryModel

	for i := 0; i < n; i++ {
		ac := models.AircraftHistoryModel{
			Icao:      strconv.Itoa(i),
			Latitude:  float32(i),
			Longitude: float32(i),
			Timestamp: time.Now().Format(time.DateTime),
		}
		aircraft = append(aircraft, ac)
	}

	return aircraft
}

func CreateMockAircraftWithTimestamp(icao string, timestamp string) models.AircraftCurrentModel {
	return models.AircraftCurrentModel{
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

func CreateMockAircraftWithIcao(n int, icao string) []models.AircraftCurrentModel {
	var aircraft []models.AircraftCurrentModel

	for i := 0; i < n; i++ {
		ac := models.AircraftCurrentModel{
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
