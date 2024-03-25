package db

import (
	"adsb-api/internal/db/models"
	"adsb-api/internal/global"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type Database interface {
	Close() error
	CreateAdsbTables() error
	BulkInsertCurrentTimeAircraftTable(aircraft []models.AircraftCurrentModel) error
	AddHistoryFromCurrent() error
	DeleteOldCurrentAircraft() error
	GetAllCurrentAircraft() ([]models.AircraftCurrentModel, error)
	GetHistoryByIcao(search string) ([]models.AircraftHistoryModel, error)
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

func (db *AdsbDB) CreateAdsbTables() error {
	err := db.createCurrentTimeAircraftTable()
	if err != nil {
		return err
	}

	err = db.createHistoryAircraft()
	if err != nil {
		return err
	}
	return nil
}

// createCurrentTimeAircraftTable creates current_time_aircraft table in database if it does not already exist
func (db *AdsbDB) createCurrentTimeAircraftTable() error {
	// Begin a transaction
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}

	var query = `CREATE TABLE IF NOT EXISTS current_time_aircraft(
				 icao VARCHAR(6) NOT NULL,
				 callsign VARCHAR(10) NOT NULL,
				 altitude INT NOT NULL,
				 lat DECIMAL NOT NULL,
				 long DECIMAL NOT NULL,
				 speed INT NOT NULL,
				 track INT NOT NULL,
				 vspeed INT NOT NULL,
				 timestamp TIMESTAMP NOT NULL,
				 PRIMARY KEY (icao,timestamp))`

	// Create current_time table
	_, err = tx.Exec(query)

	if err != nil {
		tx.Rollback()
		return err
	}

	query = `CREATE INDEX IF NOT EXISTS timestamp_index ON current_time_aircraft(timestamp)`

	// Create another index on the timestamp column
	_, err = tx.Exec(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Commit the transaction
	return tx.Commit()
}

// createHistoryAircraft creates table for storing aircraft history
func (db *AdsbDB) createHistoryAircraft() error {
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}

	var query = `CREATE TABLE IF NOT EXISTS history_aircraft(
				 icao VARCHAR(6) NOT NULL,
				 lat DECIMAL NOT NULL,
				 long DECIMAL NOT NULL,
				 timestamp TIMESTAMP NOT NULL,
				 PRIMARY KEY (icao,timestamp))`

	_, err = tx.Exec(query)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// BulkInsertCurrentTimeAircraftTable updates the current_time_aircraft table with the new aircraft records provided from
// the parameter 'aircraft'
func (db *AdsbDB) BulkInsertCurrentTimeAircraftTable(aircraft []models.AircraftCurrentModel) error {
	// Maximum number of aircraft per query
	// (65535 is the max amount of parameters postgres supports and there are 9 aircraft parameters)
	const maxAircraft = 65535 / 9

	for i := 0; i < len(aircraft); i += maxAircraft {
		end := i + maxAircraft
		if end > len(aircraft) {
			end = len(aircraft)
		}

		var (
			placeholders []string
			vals         []interface{}
		)

		for j, ac := range aircraft[i:end] {
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				j*9+1, j*9+2, j*9+3, j*9+4, j*9+5, j*9+6, j*9+7, j*9+8, j*9+9))

			vals = append(vals, ac.Icao, ac.Callsign, ac.Altitude, ac.Latitude, ac.Longitude,
				ac.Speed, ac.Track, ac.VerticalRate, ac.Timestamp)
		}

		query := `INSERT INTO current_time_aircraft (icao, callsign, altitude, lat, long, speed, track, vspeed, timestamp) VALUES %s`
		stmt := fmt.Sprintf(query, strings.Join(placeholders, ","))
		_, err := db.Conn.Exec(stmt, vals...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *AdsbDB) AddHistoryFromCurrent() error {
	query := `INSERT INTO history_aircraft (icao, lat, long, timestamp) 
			  SELECT icao, lat, long, timestamp
			  FROM current_time_aircraft
			  ON CONFLICT (icao, timestamp) 
			      DO UPDATE SET timestamp = excluded.timestamp`
	_, err := db.Conn.Exec(query)
	return err
}

// DeleteOldCurrentAircraft will delete rows older than 6 seconds from the latest entry.
func (db *AdsbDB) DeleteOldCurrentAircraft() error {
	// Begin transaction
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}

	var query = `DELETE FROM current_time_aircraft 
       			 WHERE timestamp < (select max(timestamp)-($1 * interval '1 second') 
                 FROM current_time_aircraft)`

	// Delete all rows older than 6 second from the latest entry
	_, err = tx.Exec(query, global.WaitingTime+2)
	// Rolls back transaction if failed
	if err != nil {
		tx.Rollback()
		return err
	}
	// Commit transaction
	return tx.Commit()

}

// GetAllCurrentAircraft retrieves a list of all current aircraft in the current_time_aircraft table
func (db *AdsbDB) GetAllCurrentAircraft() ([]models.AircraftCurrentModel, error) {
	// Make the query to the database

	var query = `SELECT * FROM current_time_aircraft 
				 WHERE timestamp > (select max(timestamp)-($1 * interval '1 second') 
				 FROM current_time_aircraft)`

	rows, err := db.Conn.Query(query, global.WaitingTime+2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aircraft []models.AircraftCurrentModel

	for rows.Next() {
		var ac models.AircraftCurrentModel
		err := rows.Scan(&ac.Icao, &ac.Callsign, &ac.Altitude, &ac.Latitude, &ac.Longitude, &ac.Speed, &ac.Track,
			&ac.VerticalRate, &ac.Timestamp)
		if err != nil {
			return nil, err
		}

		aircraft = append(aircraft, ac)
	}

	return aircraft, nil
}

func (db *AdsbDB) GetHistoryByIcao(search string) ([]models.AircraftHistoryModel, error) {
	var query = `SELECT icao, long, lat FROM history_aircraft WHERE icao = $1`
	rows, err := db.Conn.Query(query, search)
	if err != nil {
		return []models.AircraftHistoryModel{}, err
	}
	defer rows.Close()

	var aircraft []models.AircraftHistoryModel

	for rows.Next() {
		var ac models.AircraftHistoryModel
		err := rows.Scan(&ac.Icao, &ac.Latitude, &ac.Longitude)
		if err != nil {
			return nil, err
		}

		aircraft = append(aircraft, ac)
	}

	return aircraft, nil
}
