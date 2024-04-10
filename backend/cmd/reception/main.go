package main

import (
	"adsb-api/internal/global"
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/service/sbsService"
	"adsb-api/internal/utility/logger"
	"github.com/rs/zerolog/log"
	"time"
)

// main method and starting point of the reception and processing part of the ADS-B API
func main() {
	// Initialize environment variables
	global.InitEnvironment()
	// Initialize logger
	logger.InitLogger()
	// Initialize the database
	sbsSvc, err := sbsService.InitSbsService()
	if err != nil {
		log.Fatal().Msgf("error opening database: %v", err)
	}
	log.Info().Msgf("Reception API successfully connected to database with: User: %s | Name: %s | Host: %s | port: %d",
		global.DbUser, global.DbName, global.DbHost, global.DbPort)

	defer func() {
		err := sbsSvc.DB.Close()
		if err != nil {
			log.Fatal().Msgf(errorMsg.ErrorClosingDatabase+": %q", err)
		}
	}()

	if err := sbsSvc.CreateAdsbTables(); err != nil {
		log.Fatal().Msgf(errorMsg.ErrorCreatingDatabaseTables+": %q", err)
	}

	log.Info().Msgf("Starting the process for recieving SBS data. \n"+
		"SBS source : %q | WaitingTime: %d seconds | CleaningPeriod: %d seconds | UpdatingPeriod: %d seconds | MaxDaysHistory: %d",
		global.SbsSource, global.WaitingTime, global.CleaningPeriod, global.UpdatingPeriod, global.MaxDaysHistory)
	timer := time.Now()
	for {
		aircraft, err := sbsSvc.ProcessSbsData()
		if err != nil {
			log.Error().Msg(errorMsg.ErrorCouldNotConnectToTcpStream)
			time.Sleep(time.Duration(global.WaitingTime) * time.Second)
			continue
		} else if len(aircraft) == 0 {
			log.Warn().Msgf("recieved no data from SBS data source, will try again in: %d seconds", global.WaitingTime)
			time.Sleep(time.Duration(global.WaitingTime) * time.Second)
			continue
		}

		err = sbsSvc.InsertNewSbsData(aircraft)
		if err != nil {
			log.Fatal().Msgf(errorMsg.ErrorInsertingNewSbsData+": %q", err)
		}
		log.Info().Msgf("new SBS data inserted")

		if diff := time.Since(timer).Seconds(); diff > float64(global.CleaningPeriod) {
			if err = sbsSvc.Cleanup(); err == nil {
				timer = time.Now()
				log.Info().Msgf("old SBS data deleted")
			}
		}

		time.Sleep(time.Duration(global.UpdatingPeriod) * time.Second)
	}
}
