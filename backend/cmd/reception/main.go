package main

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global"
	"adsb-api/internal/handler/currentAircraftHandler"
	"adsb-api/internal/handler/defaultHandler"
	"adsb-api/internal/logger"
	"database/sql"
	"net/http"

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
		// Create current time aircraft table if it does not already exists
		if err := db.CreateCurrentTimeAircraftTable(dbConn); err != nil {
			logger.Error.Fatalf("Current_time_aircraft table was not created: %q", err)
		}
	}
	// Close the connection to the datase at the end
	defer func(conn *sql.DB) {
		err := db.CloseDatabase(conn)
		if err != nil {
			logger.Error.Fatalf("Could not close database connection: %q", err)
		}
	}(dbConn)

	http.HandleFunc(global.DefaultPath, defaultHandler.DefaultHandler)
	http.HandleFunc(global.CurrentAircraftPath, currentAircraftHandler.CurrentAircraftHandler(dbConn))

	logger.Info.Println("Listening on port: 8080 ...")
	logger.Info.Fatal(http.ListenAndServe(":"+global.DefaultPort, nil))

}
