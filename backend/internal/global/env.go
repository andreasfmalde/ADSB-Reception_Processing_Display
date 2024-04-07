package global

import (
	"adsb-api/internal/utility/logger"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func InitEnvironment() {
	err := godotenv.Load("./.env")
	if err != nil {
		logger.Error.Printf(err.Error())
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
		logger.Warning.Printf("error setting environment variable 'DB_PORT': can only be an integer: Error %q", err.Error())
	}
}

func InitSbsEnvVariables() {
	SbsSource = os.Getenv("SBS_SOURCE")

	var err error
	WaitingTime, err = strconv.Atoi(os.Getenv("WAITING_TIME"))
	if err != nil {
		logger.Warning.Printf("error setting environment variable 'WAITING_TIME': can only be an integer: Error %q", err.Error())
	}

	CleanupSchedule = os.Getenv("CLEANING_SCHEDULE")

	UpdatingPeriod, err = strconv.Atoi(os.Getenv("UPDATING_PERIOD"))
	if err != nil {
		logger.Warning.Printf("error setting environment variable 'UPDATING_PERIOD': can only be an integer: Error %q", err.Error())
	}

	MaxDaysHistory, err = strconv.Atoi(os.Getenv("MAX_DAYS_HISTORY"))
	if err != nil {
		logger.Warning.Printf("error setting environment variable 'MAX_DAYS_HISTORY': can only be an integer: Error %q", err.Error())
	}
}

func CheckEnvVariables() {
	if DbUser == "" {
		logger.Error.Fatal("required environment variable for database username (DB_USER) was not set")
	}

	if DbPassword == "" {
		logger.Error.Fatal("required environment variable for database password (DB_PASSWORD) was not set")
	}

	if DbName == "" {
		DbName = "adsb"
		logger.Warning.Printf("environment variable database name (DB_NAME) was not set, u"+
			"sing default name: %q", DbName)
	}

	if DbHost == "" {
		DbHost = "localhost"
		logger.Warning.Printf("environment variable database host (DB_HOST) was not set, "+
			"using default host: %q", DbHost)
	}

	if DbPort == 0 {
		DbPort = 5432
		logger.Warning.Printf("environment variable database port (DB_PORT) was not set, "+
			"using default port: %q", DbPort)
	}

	if SbsSource == "" {
		logger.Error.Fatal("required environment variable for SBS data source (SBS_SOURCE) was not set")
	}

	if WaitingTime == 0 {
		WaitingTime = 4
		logger.Warning.Printf("environment variable waiting time (WAITING_TIME), "+
			"time period in seconds between each batch of SBS data, was not set, using default time: %q seconds", WaitingTime)
	}

	if CleanupSchedule == "" {
		CleanupSchedule = "0 0 * * *"
		logger.Warning.Printf("environment variable cleanning schedule (CLEANING_SCHEDULE), "+
			"how often cleaning how data should happen to save space, was not set, using default cron schedule: %q", CleanupSchedule)
	}

	if UpdatingPeriod == 0 {
		UpdatingPeriod = 10
		logger.Warning.Printf("environment variable updating periode (UPDATING_PERIOD), "+
			"how often the the programs request for new data, was not set, using defaul time: %q seconds", UpdatingPeriod)
	}

	if MaxDaysHistory == 0 {
		MaxDaysHistory = 1
		logger.Warning.Printf("environment variable max days of history (MAX_DAYS_HISTORY), "+
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
