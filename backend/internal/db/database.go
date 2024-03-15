package db

import (
	"adsb-api/internal/global"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Initialize the PostgreSQL database and return the connection pointer
func InitDatabase() (*sql.DB, error) {

	dbLogin := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		global.Host, global.Port, global.User, global.Password, global.Dbname)
	// Open a SQL connection to the database
	return sql.Open("postgres", dbLogin)

}

// Close the connection to the database
func CloseDatabase(db *sql.DB) error {
	return db.Close()
}

// Create current_time_aircraft table in database if it does not already exists
func CreateCurrentTimeAircraftTable(db *sql.DB) error {
	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// Create current_time table
	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS current_time_aircraft(" +
		"icao VARCHAR(6) NOT NULL," +
		"callsign VARCHAR(10) NOT NULL," +
		"altitude INT NOT NULL," +
		"lat DECIMAL NOT NULL," +
		"long DECIMAL NOT NULL," +
		"speed INT NOT NULL," +
		"track INT NOT NULL," +
		"vspeed INT NOT NULL," +
		"timestamp TIMESTAMP NOT NULL," +
		"PRIMARY KEY(icao,timestamp	));")

	if err != nil {
		tx.Rollback()
		return err
	}
	// Create another index on the timestamp column
	_, err = tx.Exec("CREATE INDEX IF NOT EXISTS timestamp_index ON current_time_aircraft(timestamp);")
	if err != nil {
		tx.Rollback()
		return err
	}
	// Commit the transaction
	return tx.Commit()
}

// Update the current_time_aircraft table with the new aircraft records provided from
// the parameter 'aircrafts'
func UpdateCurrentAircraftsTable(db *sql.DB, aircrafts []global.Aircraft) error {

	query := "INSERT INTO current_time_aircraft VALUES "

	for _, aircraft := range aircrafts {
		entry := fmt.Sprintf("('%s','%s','%d','%f','%f','%d','%d','%d','%s'),",
			aircraft.Icao, aircraft.Callsign, aircraft.Altitude, aircraft.Latitude, aircraft.Longitude,
			aircraft.Speed, aircraft.Track, aircraft.VerticalRate, aircraft.Timestamp)
		query = query + entry
	}
	if query[len(query)-1] == ',' {
		query = query[:len(query)-1]
	}
	query = query + ";"

	_, err := db.Exec(query)

	if err != nil {
		return err
	}

	return nil

}

// Method that will delete rows older that 6 seconds
// from the lastest entry.
func DeleteCurrentTimeAircrafts(db *sql.DB) error {
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// Delete all rows older than 6 second from the latest entry
	_, err = tx.Exec("DELETE FROM current_time_aircraft where timestamp" +
		" < (select max(timestamp)-(6 * interval '1 second')" +
		" from current_time_aircraft);")
	// Roll back transaction if failed
	if err != nil {
		tx.Rollback()
		return err
	}
	// Commit transaction
	return tx.Commit()

}

// Method to retrieve a list of all current aircrafts in the
// current_time_aircraft table
func RetrieveCurrentTimeAircrafts(db *sql.DB) (global.GeoJsonFeatureCollection, error) {
	// Make the query to the database
	rows, err := db.Query("select * from current_time_aircraft where timestamp > " +
		"(select max(timestamp)-(6 * interval '1 second') from current_time_aircraft);")
	if err != nil {
		return global.GeoJsonFeatureCollection{}, err
	}
	defer rows.Close()

	featureCollection := global.GeoJsonFeatureCollection{}
	featureCollection.Type = "FeatureCollection"

	for rows.Next() {
		properties := global.AircraftProperties{}
		var lat float32
		var long float32

		err := rows.Scan(&properties.Icao, &properties.Callsign, &properties.Altitude, &lat,
			&long, &properties.Speed, &properties.Track,
			&properties.VerticalRate, &properties.Timestamp)
		if err != nil {
			return global.GeoJsonFeatureCollection{}, err
		}

		feature := global.GeoJsonFeature{}
		feature.Type = "Feature"
		feature.Properties = properties
		feature.Geometry.Coordinates = append(feature.Geometry.Coordinates, lat, long)
		feature.Geometry.Type = "Point"

		featureCollection.Features = append(featureCollection.Features, feature)
	}

	return featureCollection, nil
}
