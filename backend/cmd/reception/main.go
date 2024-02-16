package main

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global"
	"adsb-api/internal/logger"
	"database/sql"

	_ "github.com/lib/pq"
)

/*
Main method and starting point of the reception and prosessing part of
the ADS-B API
*/
func main() {
	// Initialize logger
	logger.InitLogger()
	// Initialize environmental variables
	global.InitEnvironment()
	// Initialize the database
	dbConn, err := db.InitDatabase()
	if err != nil {
		logger.Error.Fatalf("Error opening database: %q", err)
	} else {
		logger.Info.Println("Successfully connected to database!")
	}
	// Close the connection to the datase at the end
	defer func(conn *sql.DB) {
		err := db.CloseDatabase(conn)
		if err != nil {
			logger.Error.Fatalf("Could not close database connection: %q", err)
		}
	}(dbConn)

}
