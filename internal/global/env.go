package global

import (
	"adsb-api/internal/utility/logger"
	"os"

	"github.com/joho/godotenv"
)

func InitProdEnvironment() {
	logger.InitLogger()
	_ = godotenv.Load("../.env")
	DbUser = os.Getenv("DB_USER")
	DbPassword = os.Getenv("DB_PASSWORD")
	SbsSource = os.Getenv("SBS_SOURCE")
}

func InitTestEnvironment() {
	logger.InitLogger()
	DbUser = "test"
	DbPassword = "test"
	Dbname = "adsb_test_db"
	SbsSource = "localhost:9999"
}
