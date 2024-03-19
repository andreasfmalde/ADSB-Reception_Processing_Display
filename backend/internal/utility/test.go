package utility

import (
	"adsb-api/internal/global"
	"strconv"
	"time"
)

func CreateAircraft(n int) []global.Aircraft {
	var aircraft []global.Aircraft

	for i := 0; i < n; i++ {
		ac := global.Aircraft{
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
