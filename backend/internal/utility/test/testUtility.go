package test

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global"
	"adsb-api/internal/logger"
	"database/sql"
	"fmt"
	"strconv"
)

func InitTestDb() *sql.DB {
	dbLogin := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		global.Host, global.Port, "test", "test", "adsb_test_db")
	conn, err := sql.Open("postgres", dbLogin)
	if err != nil {
		logger.Error.Fatalf("Error opening test database: %q", err)
		return nil
	}
	logger.Info.Println("Successfully connected to test database.")

	err = db.CreateCurrentTimeAircraftTable(conn)
	if err != nil {
		logger.Error.Fatalf("Current_time_aircraft table was not created: %q", err)
		return nil
	}

	return conn
}

func PopulateSeqTestDB(dbConn *sql.DB, nAircraft int) {
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

	err := db.UpdateCurrentAircraftsTable(dbConn, aircraft)
	if err != nil {
		logger.Error.Fatalf(err.Error())
	}
}

func CleanTestDB(dbConn *sql.DB) {
	_, err := dbConn.Exec("DROP TABLE current_time_aircraft;")
	err = db.CreateCurrentTimeAircraftTable(dbConn)
	if err != nil {
		logger.Error.Fatalf(err.Error())
	}
}
