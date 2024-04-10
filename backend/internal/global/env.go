package global

import (
	"adsb-api/internal/utility/logger"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

func InitEnvironment() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Error().Msgf(err.Error())
	}

	InitDatabaseEnvVariables()
	InitSbsEnvVariables()
	CheckEnvVariables()
}

func InitDatabaseEnvVariables() {
	DbUser = os.Getenv("DB_USER")
	DbPassword = os.Getenv("DB_PASSWORD")
	DbName = os.Getenv("DB_NAME")
	DbHost = os.Getenv("DB_HOST")

	var err error
	DbPort, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Warn().Msgf("error setting environment variable 'DB_PORT': can only be an integer: Error %q", err.Error())
	}
}

func InitSbsEnvVariables() {
	SbsSource = os.Getenv("SBS_SOURCE")

	var err error
	WaitingTime, err = strconv.Atoi(os.Getenv("WAITING_TIME"))
	if err != nil {
		log.Warn().Msgf("error setting environment variable 'WAITING_TIME': can only be an integer: Error %q", err.Error())
	}

	CleaningPeriod, err = strconv.Atoi(os.Getenv("CLEANING_PERIOD"))
	if err != nil {
		log.Warn().Msgf("error setting environment variable 'CLEANING_PERIOD': can only be an integer: Error %q", err.Error())
	}

	UpdatingPeriod, err = strconv.Atoi(os.Getenv("UPDATING_PERIOD"))
	if err != nil {
		log.Warn().Msgf("error setting environment variable 'UPDATING_PERIOD': can only be an integer: Error %q", err.Error())
	}

	MaxDaysHistory, err = strconv.Atoi(os.Getenv("MAX_DAYS_HISTORY"))
	if err != nil {
		log.Warn().Msgf("error setting environment variable 'MAX_DAYS_HISTORY': can only be an integer: Error %q", err.Error())
	}
}

func CheckEnvVariables() {
	if DbUser == "" {
		log.Fatal().Msgf("required environment variable for database username (DB_USER) was not set")
	}

	if DbPassword == "" {
		log.Fatal().Msgf("required environment variable for database password (DB_PASSWORD) was not set")
	}

	if DbName == "" {
		DbName = "adsb"
		log.Warn().Msgf("environment variable database name (DB_NAME) was not set, u"+
			"sing default name: %q", DbName)
	}

	if DbHost == "" {
		DbHost = "localhost"
		log.Warn().Msgf("environment variable database host (DB_HOST) was not set, "+
			"using default host: %q", DbHost)
	}

	if DbPort == 0 {
		DbPort = 5432
		log.Warn().Msgf("environment variable database port (DB_PORT) was not set, "+
			"using default port: %q", DbPort)
	}

	if SbsSource == "" {
		log.Fatal().Msgf("required environment variable for SBS data source (SBS_SOURCE) was not set")
	}

	if WaitingTime == 0 {
		WaitingTime = 4
		log.Fatal().Msgf("environment variable waiting time (WAITING_TIME), "+
			"time period in seconds between each batch of SBS data, was not set, using default time: %q seconds", WaitingTime)
	}

	if CleaningPeriod == 0 {
		CleaningPeriod = 120
		log.Fatal().Msgf("environment variable cleanning periode (CLEANING_PERIOD), "+
			"how often cleaning how data should happen to save space, was not set, using default time: %q sconds", CleaningPeriod)
	}

	if UpdatingPeriod == 0 {
		UpdatingPeriod = 10
		log.Fatal().Msgf("environment variable updating periode (UPDATING_PERIOD), "+
			"how often the the programs request for new data, was not set, using defaul time: %q seconds", UpdatingPeriod)
	}

	if MaxDaysHistory == 0 {
		MaxDaysHistory = 1
		log.Fatal().Msgf("environment variable max days of history (MAX_DAYS_HISTORY), "+
			"how many days of history to keep when cleaning data, was not set, using defaul value: %q", MaxDaysHistory)
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
}
