package global

import (
	"github.com/joho/godotenv"
	"os"
)

func InitEnvironment() {
	_ = godotenv.Load("../.env")
	User = os.Getenv("DB_USER")
	Password = os.Getenv("DB_PASSWORD")
}

func InitTestEnv() {
	User = "test"
	Password = "test"
	Dbname = "adsb_test_db"
}
