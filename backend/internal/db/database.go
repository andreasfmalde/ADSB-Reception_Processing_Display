package db

import (
	"adsb-api/internal/global"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type Database interface {
	InitDB() (*AdsbDB, error)
	Close() error
	CreateCurrentTimeAircraftTable() error
	BulkInsertCurrentTimeAircraftTable(aircraft []global.Aircraft) error
	DeleteOldCurrentAircraft() error
	GetAllCurrentAircraft() (global.GeoJsonFeatureCollection, error)
}

type AdsbDB struct {
	Conn *sql.DB
}

// InitDB the PostgresSQL database and return the connection pointer
func InitDB() (*AdsbDB, error) {
	dbLogin := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		global.Host, global.Port, global.User, global.Password, global.Dbname)

	dbConn, err := sql.Open("postgres", dbLogin)
	return &AdsbDB{Conn: dbConn}, err
}

func (db *AdsbDB) Close() error {
	return db.Conn.Close()
}

// CreateCurrentTimeAircraftTable creates current_time_aircraft table in database if it does not already exist
func (db *AdsbDB) CreateCurrentTimeAircraftTable() error {
	// Begin a transaction
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}

	// Create current_time table
	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS current_time_aircraft(
		icao VARCHAR(6) NOT NULL,
		callsign VARCHAR(10) NOT NULL,
		altitude INT NOT NULL,
		lat DECIMAL NOT NULL,
		long DECIMAL NOT NULL,
		speed INT NOT NULL,
		track INT NOT NULL,
		vspeed INT NOT NULL,
		timestamp TIMESTAMP NOT NULL,
		PRIMARY KEY (icao,timestamp))`)

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

// BulkInsertCurrentTimeAircraftTable updates the current_time_aircraft table with the new aircraft records provided from
// the parameter 'aircraft'
func (db *AdsbDB) BulkInsertCurrentTimeAircraftTable(aircraft []global.Aircraft) error {
	var (
		placeholders []string
		vals         []any
	)

	for i, aircraft := range aircraft {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*9+1, i*9+2, i*9+3, i*9+4, i*9+5, i*9+6, i*9+7, i*9+8, i*9+9))

		vals = append(vals, aircraft.Icao, aircraft.Callsign, aircraft.Altitude, aircraft.Latitude, aircraft.Longitude,
			aircraft.Speed, aircraft.Track, aircraft.VerticalRate, aircraft.Timestamp)
	}

	stmt := fmt.Sprintf("INSERT INTO current_time_aircraft (icao, callsign, altitude, lat, long, speed, track, vspeed, timestamp) VALUES %s",
		strings.Join(placeholders, ","))
	_, err := db.Conn.Exec(stmt, vals...)
	return err
}

// DeleteOldCurrentAircraft will delete rows older than 6 seconds from the latest entry.
func (db *AdsbDB) DeleteOldCurrentAircraft() error {
	// Begin transaction
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}
	// Delete all rows older than 6 second from the latest entry
	_, err = tx.Exec(`DELETE FROM current_time_aircraft 
       					WHERE timestamp < (select max(timestamp)-(6 * interval '1 second') 
                        FROM current_time_aircraft)`)
	// Rolls back transaction if failed
	if err != nil {
		tx.Rollback()
		return err
	}
	// Commit transaction
	return tx.Commit()

}

// GetAllCurrentAircraft retrieves a list of all current aircraft in the current_time_aircraft table
func (db *AdsbDB) GetAllCurrentAircraft() (global.GeoJsonFeatureCollection, error) {
	// Make the query to the database

	rows, err := db.Conn.Query(`
							SELECT * FROM current_time_aircraft 
         					WHERE timestamp > (select max(timestamp)-(6 * interval '1 second') 
                            FROM current_time_aircraft)`)
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
