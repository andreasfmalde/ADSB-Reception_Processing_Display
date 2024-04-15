package global

import (
	"adsb-api/internal/utility/logger"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
)

// InitEnvironment initializes the environment variables by loading the .env file.
// It then calls the InitDatabaseEnvVariables and InitSbsEnvVariables functions to initialize the database and SBS
// environment variables respectively.
func InitEnvironment() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Error().Msgf("error loading .env file: %q", err)
	}

	InitDatabaseEnvVariables()
	InitSbsEnvVariables()
}

// InitDatabaseEnvVariables initializes the environment variables related to the database.
// It retrieves the values of the DB_USER, DB_PASSWORD, DB_HOST and DB_PORT environment variables and assigns
// them to the respective variables.
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

// InitSbsEnvVariables initializes the environment variables related to the SBS.
// It retrieves the values of the SBS_SOURCE, WAITING_TIME, CLEANUP_SCHEDULE, UPDATING_PERIOD,
// and MAX_DAYS_HISTORY environment variables and assigns them to the respective variables.
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

// InitTestEnvironment initializes the test environment by initializing the logger and setting up the test database
// and SBS environment variables.
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
