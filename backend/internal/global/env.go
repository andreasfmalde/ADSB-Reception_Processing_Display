package global

import (
	"adsb-api/internal/logger"
	"github.com/joho/godotenv"
	"os"
)

func InitEnvironment() {
	_ = godotenv.Load("../.env")
	User = os.Getenv("DB_USER")
	Password = os.Getenv("DB_PASSWORD")
}

func InitTestEnv() {
	logger.InitLogger()
	User = "test"
	Password = "test"
	Dbname = "adsb_test_db"
}
