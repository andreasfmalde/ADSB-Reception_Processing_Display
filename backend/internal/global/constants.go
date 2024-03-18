package global

import (
	"os"

	"github.com/joho/godotenv"
)

var Env string

const Development = "dev"
const Production = "prod"

var User string
var Password string

func InitEnv() string {
	_ = godotenv.Load("../.env")
	User = os.Getenv("DB_USER")
	Password = os.Getenv("DB_PASSWORD")

	if User == "" {
		User = "test"
		Password = "test"
		Env = Development
		Dbname = "adsb_test_db"
		return "could not find DB_USER, setting environment to " + Development + "\n"
	}

	Env = Production
	return "found DB_USER, setting environment to " + Production + "\n"
}

const (
	Host = "localhost"
	Port = 5432
)

var Dbname = "adsb_db"

const DefaultPort = "8080"
const VERSION = "1.0.0"
const DefaultPath = "/"
const CurrentAircraftPath = "/aircraft/current/"
const WaitingTime = 4
