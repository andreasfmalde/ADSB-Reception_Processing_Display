package cmd

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global"
	"adsb-api/internal/logger"
	"database/sql"
)

// InitializeApiResources initializes the logger, environmental variables and the database.
func InitializeApiResources() (dbConn *sql.DB) {
	logger.InitLogger()
	global.InitEnvironment()
	dbConn, err := db.InitDatabase()
	if err != nil {
		logger.Error.Fatalf("Error opening database: %q", err)
	} else {
		logger.Info.Println("Successfully connected to database!")
		if err := db.CreateCurrentTimeAircraftTable(dbConn); err != nil {
			logger.Error.Fatalf("Current_time_aircraft table was not created: %q", err)
		}
	}
	defer func(conn *sql.DB) {
		err := db.CloseDatabase(conn)
		if err != nil {
			logger.Error.Fatalf("Could not close database connection: %q", err)
		}
	}(dbConn)

	return dbConn
}
