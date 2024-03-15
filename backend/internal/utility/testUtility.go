package utility

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global"
	"adsb-api/internal/logger"
	"database/sql"
	"fmt"
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
