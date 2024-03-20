package main

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global"
	"adsb-api/internal/logger"
	"adsb-api/internal/utility/adsbhub"
	"time"
)

// main method and starting point of the reception and processing part of the ADS-B API
func main() {
	// Initialize logger
	logger.InitLogger()
	// Initialize environment variables
	global.InitEnvironment()
	// Initialize the database
	adsbDB, err := db.InitDB()
	if err != nil {
		logger.Error.Fatalf("error opening database: %q", err)
	}
	logger.Info.Println("successfully connected to database")

	defer func() {
		err := adsbDB.Close()
		if err != nil {
			logger.Error.Fatalf("error closing database: %q", err)
		}
	}()

	if err := adsbDB.CreateAdsbTables(); err != nil {
		logger.Error.Fatalf("error creating tables for database: %q", err)
	}

	timer := time.Now()
	for {
		aircraft, err := adsbhub.ProcessSBSstream()
		if err != nil {
			logger.Info.Println(err.Error() + " ... will try again in 4 seconds")
			time.Sleep(global.WaitingTime * time.Second)
			continue
		}
		err = adsbDB.BulkInsertCurrentTimeAircraftTable(aircraft)
		if err != nil {
			logger.Error.Fatalf("could not insert new SBS data: %q", err)
		}
		logger.Info.Println("new SBS data inserted")
		// Delete old rows every 2 minutes (120 seconds)
		if diff := time.Since(timer).Seconds(); diff > 120 {
			if e := adsbDB.DeleteOldCurrentAircraft(); e == nil {
				timer = time.Now()
				logger.Info.Println("old SBS data deleted")
			}
		}
		time.Sleep(global.WaitingTime * time.Second)
	}
}
