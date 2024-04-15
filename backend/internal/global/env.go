package global

import (
	"adsb-api/internal/utility/logger"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func InitEnvironment() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Error().Msgf("error loading .env file: %q", err)
	}

	InitDatabaseEnvVariables()
	InitSbsEnvVariables()
}

func InitDatabaseEnvVariables() {
	DbUser = os.Getenv("DB_USER")
	DbPassword = os.Getenv("DB_PASSWORD")

	dbName, exist := os.LookupEnv("DB_NAME")
	if exist {
		DbName = dbName
	}

	dbHost, exist := os.LookupEnv("DB_HOST")
	if exist {
		DbHost = dbHost
	}

	var err error
	dbPort, exist := os.LookupEnv("DB_PORT")
	if exist {
		DbPort, err = strconv.Atoi(dbPort)
		if err != nil {
			log.Warn().Msgf("error setting environment variable 'DB_PORT': can only be an integer: Error %q", err)
		}
	}
}

func InitSbsEnvVariables() {
	SbsSource = os.Getenv("SBS_SOURCE")

	var err error
	waitingTime, exist := os.LookupEnv("WAITING_TIME")
	if exist {
		WaitingTime, err = strconv.Atoi(waitingTime)
		if err != nil {
			log.Warn().Msgf("error setting environment variable 'WAITING_TIME': can only be an integer: Error %q", err)
		}
	}

	cleanupSchedule, exist := os.LookupEnv("CLEANUP_SCHEDULE")
	if exist {
		CleanupSchedule = cleanupSchedule
	}

	updatingPeriod, exist := os.LookupEnv("UPDATING_PERIOD")
	if exist {
		UpdatingPeriod, err = strconv.Atoi(updatingPeriod)
		if err != nil {
			log.Warn().Msgf("error setting environment variable 'UPDATING_PERIOD': can only be an integer: Error %q", err)
		}
	}

	maxDaysHistory, exist := os.LookupEnv("MAX_DAYS_HISTORY")
	if exist {
		MaxDaysHistory, err = strconv.Atoi(maxDaysHistory)
		if err != nil {
			log.Warn().Msgf("error setting environment variable 'MAX_DAYS_HISTORY': can only be an integer: Error %q", err)
		}
	}
}

func InitTestEnvironment() {
	logger.InitLogger()
	DbUser = "test"
	DbPassword = "test"
	DbName = "adsb_test_db"
	DbHost = "localhost"
	DbPort = 5432

	SbsSource = "localhost:9999"
	WaitingTime = 4
	CleanupSchedule = "0 0 * * *"
	UpdatingPeriod = 10
	MaxDaysHistory = 1
}
