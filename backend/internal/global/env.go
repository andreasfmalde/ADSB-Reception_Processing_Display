package global

import (
	"adsb-api/internal/utility/logger"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func InitProdEnvironment() {
	logger.InitLogger()
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error.Printf(err.Error())
	}

	DbUser = os.Getenv("DB_USER")
	DbPassword = os.Getenv("DB_PASSWORD")
	DbName = os.Getenv("DB_NAME")
	DbHost = os.Getenv("DB_HOST")
	DbPort, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		DbPort = 5432
		logger.Error.Printf("error setting database port from environment: %s Default port: %v", err.Error(), DbPort)
	}

	SbsSource = os.Getenv("SBS_SOURCE")
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
