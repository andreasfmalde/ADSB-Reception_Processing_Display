package global

import (
	"adsb-api/internal/logger"
	"github.com/joho/godotenv"
	"os"
)

func InitEnvironment() {
	_ = godotenv.Load("../.env")
	DbUser = os.Getenv("DB_USER")
	DbPassword = os.Getenv("DB_PASSWORD")
}

func InitTestEnv() {
	logger.InitLogger()
	DbUser = "test"
	DbPassword = "test"
	Dbname = "adsb_test_db"
}
