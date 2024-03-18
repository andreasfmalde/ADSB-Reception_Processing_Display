package global

import (
	"github.com/joho/godotenv"
	"os"
)

var Env string

const Development = "dev"
const Production = "prod"

func InitEnv() string {
	_ = godotenv.Load("../.env")
	User = os.Getenv("DB_USER")
	Password = os.Getenv("DB_PASSWORD")

	if User == "" {
		InitDevEnv()
		return "could not find DB_USER, setting environment to " + Development + "\n"
	}

	InitProdDev()
	return "found DB_USER, setting environment to " + Production + "\n"
}

func InitDevEnv() {
	User = "test"
	Password = "test"
	Env = Development
	Dbname = "adsb_test_db"
}

func InitProdDev() {
	Env = Production
}
