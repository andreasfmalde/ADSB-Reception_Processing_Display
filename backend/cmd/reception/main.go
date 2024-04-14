package main

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global"
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/sbs"
	"adsb-api/internal/service/cronScheduler"
	"adsb-api/internal/service/sbsService"
	"adsb-api/internal/utility/logger"
	"github.com/rs/zerolog/log"
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
		log.Fatal().Msgf("error opening database: %q", err)
	}

	defer func() {
		err = database.Close()
		if err != nil {
			log.Fatal().Msgf(errorMsg.ErrorClosingDatabase+": %q", err)
		}
	}()

	// Initialize cron scheduler
	scheduler := cronScheduler.NewCronScheduler()

	// Initialize SBS service
	sbsSvc := sbsService.InitSbsService(database, scheduler)

	if err := sbsSvc.CreateAdsbTables(); err != nil {
		log.Fatal().Msgf(errorMsg.ErrorCreatingDatabaseTables+": %q", err)
	}

	if err := sbsSvc.ScheduleCleanUpJob(global.CleanupSchedule, global.MaxDaysHistory); err != nil {
		log.Fatal().Msgf("error initiazling cleanupJob job")
	}

	log.Info().Msgf("Reception API successfully connected to database with: User: %s | Name: %s | Host: %s | port: %d",
		global.DbUser, global.DbName, global.DbHost, global.DbPort)

	log.Info().Msgf("Scheduled clean up job with cron schedule: %s", global.CleanupSchedule)

	log.Info().Msgf("Starting the process for receiving SBS data. \n"+
		"SBS source : %q | WaitingTime: %d seconds | CleanupSchedule: %s seconds | UpdatingPeriod: %d seconds | MaxDaysHistory: %d",
		global.SbsSource, global.WaitingTime, global.CleanupSchedule, global.UpdatingPeriod, global.MaxDaysHistory)

	for {
		aircraft, err := sbs.ProcessSbsStream(global.SbsSource, global.WaitingTime)
		if err != nil {
			log.Error().Msgf(errorMsg.ErrorCouldNotConnectToTcpStream)
			time.Sleep(time.Duration(global.WaitingTime) * time.Second)
			continue
		} else if len(aircraft) == 0 {
			log.Warn().Msgf("received no data from SBS data source, will try again in: %d seconds", global.WaitingTime)
			time.Sleep(time.Duration(global.WaitingTime) * time.Second)
			continue
		}

		err = sbsSvc.InsertNewSbsData(aircraft)
		if err != nil {
			log.Error().Msgf(errorMsg.ErrorInsertingNewSbsData+": %q", err)
		}
		log.Info().Msgf("%d new aircraft inserted", len(aircraft))

		time.Sleep(time.Duration(global.UpdatingPeriod) * time.Second)
	}
}
