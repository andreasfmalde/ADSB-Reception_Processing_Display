package main

import (
	"adsb-api/internal/global"
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/logger"
	"adsb-api/internal/service"
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
		err = sbsSvc.InsertNewSbsData()
		if diff := time.Since(timer).Seconds(); diff > 120 {
			if e := sbsSvc.Cleanup(); e == nil {
				timer = time.Now()
			}
		}
		time.Sleep(global.WaitingTime * time.Second)
	}
}
