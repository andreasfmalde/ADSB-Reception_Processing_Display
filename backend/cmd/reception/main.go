package main

import (
	"adsb-api/internal/global"
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/logger"
	"adsb-api/internal/service"
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
	sbsSvc, err := service.InitSbsService()
	if err != nil {
		logger.Error.Fatalf("error opening database: %q", err)
	}
	logger.Info.Println("successfully connected to database")

	defer func() {
		err := sbsSvc.DB.Close()
		if err != nil {
			logger.Error.Fatalf(errorMsg.ErrorClosingDatabase, err)
		}
	}()

	if err := sbsSvc.CreateAdsbTables(); err != nil {
		logger.Error.Fatalf(errorMsg.ErrorCreatingDatabaseTables, err)
	}

	timer := time.Now()
	for {
		aircraft, err := adsbhub.ProcessSBSstream()
		if err != nil {
			logger.Info.Println(err.Error() + " ... will try again in 4 seconds")
			time.Sleep(global.WaitingTime * time.Second)
			continue
		}
		err = sbsSvc.InsertNewAircraft(aircraft)
		if err != nil {
			logger.Error.Fatalf(errorMsg.ErrorInsertingNewSbsData, err)
		}
		logger.Info.Println("new SBS data inserted")
		err = sbsSvc.UpdateHistory()
		if err != nil {
			logger.Error.Fatalf("could not add history data: %q", err)
		}
		logger.Info.Println("new history data inserted")
		// Delete old rows every 2 minutes (120 seconds)
		if diff := time.Since(timer).Seconds(); diff > 120 {
			if e := sbsSvc.Cleanup(); e == nil {
				timer = time.Now()
				logger.Info.Println("old SBS data deleted")
			}
		}
		time.Sleep(global.WaitingTime * time.Second)
	}
}
