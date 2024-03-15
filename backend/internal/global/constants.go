package global

import (
	"os"

	"github.com/joho/godotenv"
)

var User string
var Password string

func InitEnvironment() {
	godotenv.Load("../.env")
	User = os.Getenv("DB_USER")
	Password = os.Getenv("DB_PASSWORD")
}

const (
	Host   = "localhost"
	Port   = 5432
	Dbname = "adsb_db"
)

const WaitingTime = 4
