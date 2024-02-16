package db

import (
	"adsb-api/internal/global"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

/*
Initialize the PostgreSQL database and return the connection pointer
*/
func InitDatabase() (*sql.DB, error) {

	dbLogin := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		global.Host, global.Port, global.User, global.Password, global.Dbname)
	// Open a SQL connection to the database
	return sql.Open("postgres", dbLogin)

}

/*
Close the connection to the database
*/
func CloseDatabase(db *sql.DB) error {
	return db.Close()
}

/*
Create current_time_aircraft table in database if it does not already exists
*/
func CreateCurrentTimeAircraftTable(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS current_time_aircraft(" +
		"icao VARCHAR(6) NOT NULL PRIMARY KEY ," +
		"callsign VARCHAR(10) NOT NULL," +
		"altitude INT NOT NULL," +
		"lat DECIMAL NOT NULL," +
		"long DECIMAL NOT NULL," +
		"speed INT NOT NULL," +
		"track INT NOT NULL," +
		"vspeed INT NOT NULL," +
		"timestamp TIMESTAMP NOT NULL);")

	return err
}
