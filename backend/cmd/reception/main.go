package main

import (
	"adsb-api/internal/global"
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/sbs"
	"adsb-api/internal/service/sbsService"
	"adsb-api/internal/utility/logger"
	"time"
)

// main method and starting point of the reception and processing part of the ADS-B API
func main() {
	// Initialize environment variables
	global.InitProdEnvironment()
	// Initialize the database
	sbsSvc, err := sbsService.InitSbsService()
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
		aircraft, err := sbs.ProcessSbsStream()
		if err != nil {
			logger.Error.Fatalf(errorMsg.ErrorCouldNotConnectToTcpStream)
			return
		} else if len(aircraft) == 0 {
			logger.Warning.Printf("recieved no data from SBS data source, will try again in: %d seconds", global.WaitingTime)
			time.Sleep(global.WaitingTime * time.Second)
			continue
		}

		logger.Info.Printf("retrieved SBS data: %d aircraft", len(aircraft))

		err = sbsSvc.InsertNewSbsData(aircraft)
		if err != nil {
			logger.Error.Fatalf("could not insert new SBS data: %q", err)
		}
		logger.Info.Println("new SBS data inserted")

		if diff := time.Since(timer).Seconds(); diff > global.CleaningPeriod {
			if err = sbsSvc.Cleanup(); err == nil {
				timer = time.Now()
				logger.Info.Println("old SBS data deleted")
			}
		}

		time.Sleep(global.UpdatingPeriod * time.Second)
	}
}
