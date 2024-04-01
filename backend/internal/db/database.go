package db

import (
	"adsb-api/internal/global"
	"adsb-api/internal/global/models"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type Database interface {
	Close() error
	CreateAdsbTables() error
	BulkInsertAircraftCurrent(aircraft []models.AircraftCurrentModel) error
	InsertHistoryFromCurrent() error
	DeleteOldCurrent() error
	GetCurrentAircraft() ([]models.AircraftCurrentModel, error)
	GetHistoryByIcao(search string) ([]models.AircraftHistoryModel, error)
}

type AdsbDB struct {
	Conn *sql.DB
}

// InitDB initializes the PostgresSQL database and returns the connection pointer.
func InitDB() (*AdsbDB, error) {
	dbLogin := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		global.Host, global.Port, global.DbUser, global.DbPassword, global.Dbname)

	dbConn, err := sql.Open("postgres", dbLogin)
	return &AdsbDB{Conn: dbConn}, err
}

func (db *AdsbDB) Close() error {
	return db.Conn.Close()
}

// CreateAdsbTables creates all tables for the database schema
func (db *AdsbDB) CreateAdsbTables() error {
	err := db.createAircraftCurrentTable()
	if err != nil {
		return err
	}

	err = db.createAircraftHistory()
	if err != nil {
		return err
	}
	return nil
}

// createAircraftCurrentTable creates a table for storing current aircraft data if it does not already exist
func (db *AdsbDB) createAircraftCurrentTable() error {
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}

	var query = `CREATE TABLE IF NOT EXISTS aircraft_current(
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

	_, err = tx.Exec(query)

	if err != nil {
		tx.Rollback()
		return err
	}

	query = `CREATE INDEX IF NOT EXISTS timestamp_index ON aircraft_current(timestamp)`

	_, err = tx.Exec(query)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// createAircraftHistory creates a table for storing aircraft history data if it does not already exist
func (db *AdsbDB) createAircraftHistory() error {
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}

	var query = `CREATE TABLE IF NOT EXISTS aircraft_history(
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

// BulkInsertAircraftCurrent inserts an array of new aircraft data into aircraft_current
func (db *AdsbDB) BulkInsertAircraftCurrent(aircraft []models.AircraftCurrentModel) error {
	/*
		Maximum number of aircraft per query
		(65535 is the max number of parameters postgres supports and there are 9 aircraft parameters)
	*/
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

		query := `INSERT INTO aircraft_current (icao, callsign, altitude, lat, long, speed, track, vspeed, timestamp) VALUES %s`
		stmt := fmt.Sprintf(query, strings.Join(placeholders, ","))
		_, err := db.Conn.Exec(stmt, vals...)
		if err != nil {
			return err
		}
	}

	return nil
}

// InsertHistoryFromCurrent inserts all data from aircraft_current table to aircraft_history.
func (db *AdsbDB) InsertHistoryFromCurrent() error {
	query := `INSERT INTO aircraft_history (icao, lat, long, timestamp) 
			  SELECT icao, lat, long, timestamp
			  FROM aircraft_current
			  ON CONFLICT (icao, timestamp) 
			      DO UPDATE SET timestamp = excluded.timestamp`
	_, err := db.Conn.Exec(query)
	return err
}

// DeleteOldCurrent will delete rows in aircraft_current older than global.WaitingTime seconds from the latest entry.
func (db *AdsbDB) DeleteOldCurrent() error {
	// Begin transaction
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}

	var query = `DELETE FROM aircraft_current 
       			 WHERE timestamp < (select max(timestamp)-($1 * interval '1 second') 
                 FROM aircraft_current)`

	_, err = tx.Exec(query, global.WaitingTime+2)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// GetCurrentAircraft retrieves a list of all aircraft from aircraft_current that are older than global.WaitingTime + 2
func (db *AdsbDB) GetCurrentAircraft() ([]models.AircraftCurrentModel, error) {
	var query = `SELECT * FROM aircraft_current 
				 WHERE timestamp > (select max(timestamp)-($1 * interval '1 second') 
				 FROM aircraft_current)`

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

// GetHistoryByIcao retrieves a list from aircraft_history of rows matching the icao parameter.
func (db *AdsbDB) GetHistoryByIcao(search string) ([]models.AircraftHistoryModel, error) {
	var query = `SELECT icao, long, lat FROM history_aircraft WHERE icao = $1  order by timestamp asc`
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
