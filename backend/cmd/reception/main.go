package main

import (
	"adsb-api/internal/global"
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/service/sbsService"
	"adsb-api/internal/utility/logger"
	"time"
)

// main method and starting point of the reception and processing part of the ADS-B API
func main() {
	logger.InitLogger()
	// Initialize environment variables
	global.InitEnvironment()
	// Initialize the database
	sbsSvc, err := sbsService.InitSbsService()
	if err != nil {
		logger.Error.Fatalf("error opening database: %q", err)
	}
	logger.Info.Printf("Reception API successfully connected to database with: User: %s | Name: %s | Host: %s | port: %d",
		global.DbUser, global.DbName, global.DbHost, global.DbPort)

	defer func() {
		err := sbsSvc.DB.Close()
		if err != nil {
			logger.Error.Fatalf(errorMsg.ErrorClosingDatabase+": %q", err)
		}
	}()

	if err := sbsSvc.CreateAdsbTables(); err != nil {
		logger.Error.Fatalf(errorMsg.ErrorCreatingDatabaseTables+": %q", err)
	}

	logger.Info.Printf("Starting the process for recieving SBS data. \n"+
		"SBS source : %q | WaitingTime: %d seconds | CleaningPeriod: %d seconds | UpdatingPeriod: %d seconds | MaxDaysHistory: %d",
		global.SbsSource, global.WaitingTime, global.CleaningPeriod, global.UpdatingPeriod, global.MaxDaysHistory)
	timer := time.Now()
	for {
		aircraft, err := sbsSvc.ProcessSbsData()
		if err != nil {
			logger.Error.Printf(errorMsg.ErrorCouldNotConnectToTcpStream)
			time.Sleep(time.Duration(global.WaitingTime) * time.Second)
			continue
		} else if len(aircraft) == 0 {
			logger.Warning.Printf("recieved no data from SBS data source, will try again in: %d seconds", global.WaitingTime)
			time.Sleep(time.Duration(global.WaitingTime) * time.Second)
			continue
		}

		err = sbsSvc.InsertNewSbsData(aircraft)
		if err != nil {
			logger.Error.Fatalf(errorMsg.ErrorInsertingNewSbsData+": %q", err)
		}
		logger.Info.Println("new SBS data inserted")

		if diff := time.Since(timer).Seconds(); diff > float64(global.CleaningPeriod) {
			if err = sbsSvc.Cleanup(); err == nil {
				timer = time.Now()
				logger.Info.Println("old SBS data deleted")
			}
		}

		time.Sleep(time.Duration(global.UpdatingPeriod) * time.Second)
	}
}
