package utility

import (
	"adsb-api/internal/global"
	"strconv"
)

func CreateAircraft(n int) []global.Aircraft {
	var aircraft []global.Aircraft

	for i := 0; i < n; i++ {
		ac := global.Aircraft{
			Icao:         "AB" + strconv.Itoa(i),
			Callsign:     "ABC" + strconv.Itoa(i),
			Altitude:     i,
			Latitude:     float32(i),
			Longitude:    float32(i),
			Speed:        i,
			Track:        i,
			VerticalRate: i,
			Timestamp:    "2024-01-01 12:00:00",
		}
		aircraft = append(aircraft, ac)
	}

	return aircraft
}
