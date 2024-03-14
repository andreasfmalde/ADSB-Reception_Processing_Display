package main

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global"
	"adsb-api/internal/logger"
	"adsb-api/internal/utility/adsbhub"
	"database/sql"
	"time"

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
	timer := time.Now()
	for {
		aircrafts, err := adsbhub.ProcessSBSstream()
		if err != nil {
			logger.Info.Println(err.Error() + "... will try again in 4 seconds...")
			time.Sleep(4 * time.Second)
			continue
		}

		err = db.UpdateCurrentAircraftsTable(dbConn, aircrafts)
		if err != nil {
			logger.Error.Fatalf("Could not load aircrafts in database: %s", err)
		}
		logger.Info.Println("SBS data successfully in local database")
		if diff := time.Since(timer).Seconds(); diff > 120 {
			if e := db.DeleteCurrentTimeAircrafts(dbConn); e == nil {
				timer = time.Now()
				logger.Info.Println("Rows deleted successfully!")
			}
		}
		time.Sleep(4 * time.Second)
	}
}
