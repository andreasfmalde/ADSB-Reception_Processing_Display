package global

import (
	"adsb-api/internal/utility/logger"
	"os"

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
	SbsSource = os.Getenv("SBS_SOURCE")
	Host = os.Getenv("DB_HOST")
}

func InitTestEnvironment() {
	logger.InitLogger()
	DbUser = "test"
	DbPassword = "test"
	Dbname = "adsb_test_db"
	SbsSource = "localhost:9999"
}
