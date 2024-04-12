package main

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global"
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/sbs"
	"adsb-api/internal/service/cronScheduler"
	"adsb-api/internal/service/sbsService"
	"adsb-api/internal/utility/logger"
	"time"
)

// main method and starting point of the reception and processing part of the ADS-B API
func main() {
	// Initialize logger
	logger.InitLogger()
	// Initialize environment variables
	global.InitEnvironment()
	// Initialize the database
	database, err := db.InitDB()
	if err != nil {
		logger.Error.Fatalf("error opening database: %q", err)
	}

	defer func() {
		err = database.Close()
		if err != nil {
			logger.Error.Fatalf(errorMsg.ErrorClosingDatabase+": %q", err)
		}
	}()

	// Initialize cron scheduler
	scheduler := cronScheduler.NewCronScheduler()

	// Initialize SBS service
	sbsSvc := sbsService.InitSbsService(database, scheduler)

	if err := sbsSvc.CreateAdsbTables(); err != nil {
		logger.Error.Fatalf(errorMsg.ErrorCreatingDatabaseTables+": %q", err)
	}

	if err := sbsSvc.ScheduleCleanUpJob(global.CleanupSchedule, global.MaxDaysHistory); err != nil {
		logger.Error.Fatalf("error initiazling cleanupJob job")
	}

	logger.Info.Printf("Reception API successfully connected to database with: User: %s | Name: %s | Host: %s | port: %d",
		global.DbUser, global.DbName, global.DbHost, global.DbPort)

	logger.Info.Printf("Scheduled clean up job with cron schedule: %s", global.CleanupSchedule)

	logger.Info.Printf("Starting the process for receiving SBS data. \n"+
		"SBS source : %q | WaitingTime: %d seconds | CleanupSchedule: %s seconds | UpdatingPeriod: %d seconds | MaxDaysHistory: %d",
		global.SbsSource, global.WaitingTime, global.CleanupSchedule, global.UpdatingPeriod, global.MaxDaysHistory)

	for {
		aircraft, err := sbs.ProcessSbsStream(global.SbsSource, global.WaitingTime)
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

		time.Sleep(time.Duration(global.UpdatingPeriod) * time.Second)
	}
}
