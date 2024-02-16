package global

import "os"

const (
	Host   = "localhost"
	Port   = 5432
	Dbname = "adsb_db"
)

var User = os.Getenv("DB_USER")

var Password = os.Getenv("DB_PASSWORD")
