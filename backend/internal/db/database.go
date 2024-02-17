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

/*
Update the current_time_aircraft table with the new aircraft records provided from
the parameter 'aircrafts'
*/
func UpdateCurrentAircraftsTable(db *sql.DB, aircrafts []global.Aircraft) error {
	// Delete the current table
	if _, err := db.Exec("DROP TABLE current_time_aircraft"); err != nil {
		return err
	}
	// Create a new current_time_aircraft table
	if err := CreateCurrentTimeAircraftTable(db); err != nil {
		return err
	}
	// Fill the new current_time_aircraft table
	for _, aircraft := range aircrafts {
		_, err := db.Exec("INSERT INTO current_time_aircraft VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)",
			aircraft.Icao, aircraft.Callsign, aircraft.Altitude, aircraft.Latitude, aircraft.Longitude,
			aircraft.Speed, aircraft.Track, aircraft.VerticalRate, aircraft.Timestamp)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
Method to retrieve a list of all current aircrafts in the
current_time_aircraft table
*/
func RetrieveCurrentTimeAircrafts(db *sql.DB) ([]global.Aircraft, error) {
	var aircrafts []global.Aircraft
	// Make the query to the database
	rows, err := db.Query("SELECT * FROM current_time_aircraft")
	if err != nil {
		return []global.Aircraft{}, err
	}
	// Close the rows, preventing further enumeration
	defer rows.Close()

	// Create an aircraft record
	ac := global.Aircraft{}

	// Loop through the results and append each aircraft to the list/slice
	for rows.Next() {

		// Scan each row and save the values in the aircraft record
		if err := rows.Scan(&ac.Icao, &ac.Callsign, &ac.Altitude,
			&ac.Latitude, &ac.Longitude, &ac.Speed, &ac.Track,
			&ac.VerticalRate, &ac.Timestamp); err != nil {
			return []global.Aircraft{}, err
		}
		// Add the aircraft to the list/slice
		aircrafts = append(aircrafts, ac)
	}

	// Return the list of all current aircrafts
	return aircrafts, nil

}
