package test

import (
	"adsb-api/internal/global"
	"strconv"
)

/*
func PopulateSeqTestDB(db db.AdsbDB, nAircraft int) {
	var aircraft []global.Aircraft

	for i := 0; i < nAircraft; i++ {
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

	err := db.BulkInsertCurrentTimeAircraftTable(aircraft)
	if err != nil {
		logger.Error.Fatalf(err.Error())
	}
}

*/

func GetAircraft(n int) []global.Aircraft {
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

/*
func CleanTestDB(db *db.AdsbDB) {
	_, err := db.Conn.Exec(`DROP TABLE current_time_aircraft`)
	err = db.CreateCurrentTimeAircraftTable()
	if err != nil {
		logger.Error.Fatalf(err.Error())
	}
}

*/
