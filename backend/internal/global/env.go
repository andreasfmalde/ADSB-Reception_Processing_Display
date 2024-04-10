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
		logger.Error.Printf("error loading .env file: %q", err.Error())
	}

	InitDatabaseEnvVariables()
	InitSbsEnvVariables()
	CheckRequiredEnvVariables()
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
			logger.Warning.Printf("error setting environment variable 'DB_PORT': can only be an integer: Error %q", err.Error())
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
			logger.Warning.Printf("error setting environment variable 'WAITING_TIME': can only be an integer: Error %q", err.Error())
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
			logger.Warning.Printf("error setting environment variable 'UPDATING_PERIOD': can only be an integer: Error %q", err.Error())
		}
	}

	maxDaysHistory, exist := os.LookupEnv("MAX_DAYS_HISTORY")
	if exist {
		MaxDaysHistory, err = strconv.Atoi(maxDaysHistory)
		if err != nil {
			logger.Warning.Printf("error setting environment variable 'MAX_DAYS_HISTORY': can only be an integer: Error %q", err.Error())
		}
	}
}

func CheckRequiredEnvVariables() {
	if DbUser == "" {
		logger.Error.Fatal("required environment variable for database username (DB_USER) was not set")
	}

	if DbPassword == "" {
		logger.Error.Fatal("required environment variable for database password (DB_PASSWORD) was not set")
	}

	if SbsSource == "" {
		logger.Error.Fatal("required environment variable for SBS data source (SBS_SOURCE) was not set")
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

	CheckRequiredEnvVariables()
}
