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
	Begin() (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error
	CreateAircraftCurrentTable(tx *sql.Tx) error
	CreateAircraftCurrentTimestampIndex(tx *sql.Tx) error
	CreateAircraftHistory(tx *sql.Tx) error
	DropAircraftCurrentTable(tx *sql.Tx) error
	BulkInsertAircraftCurrent(aircraft []models.AircraftCurrentModel, tx *sql.Tx) error
	InsertHistoryFromCurrent(tx *sql.Tx) error
	DeleteOldCurrent(tx *sql.Tx) error
	GetCurrentAircraft(tx *sql.Tx) ([]models.AircraftCurrentModel, error)
	GetHistoryByIcao(search string, tx *sql.Tx) ([]models.AircraftHistoryModel, error)
}

type AdsbDB struct {
	Conn *sql.DB
}

// InitDB initializes the PostgresSQL database and returns the connection pointer.
func InitDB() (*AdsbDB, error) {
	dbLogin := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		global.Host, global.Port, global.DbUser, global.DbPassword, global.Dbname)

	dbConn, err := sql.Open("postgres", dbLogin)
	// TODO: add ping
	return &AdsbDB{Conn: dbConn}, err
}

func (db *AdsbDB) Close() error {
	return db.Conn.Close()
}

func (db *AdsbDB) Begin() (*sql.Tx, error) {
	return db.Conn.Begin()
}

func (db *AdsbDB) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (db *AdsbDB) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

// CreateAircraftCurrentTable creates a table for storing current aircraft data if it does not already exist
func (db *AdsbDB) CreateAircraftCurrentTable(tx *sql.Tx) error {
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
	return db.exec(query, tx)
}

func (db *AdsbDB) CreateAircraftCurrentTimestampIndex(tx *sql.Tx) error {
	var query = `CREATE INDEX IF NOT EXISTS timestamp_index ON aircraft_current(timestamp)`

	return db.exec(query, tx)
}

// CreateAircraftHistory creates a table for storing aircraft history data if it does not already exist
func (db *AdsbDB) CreateAircraftHistory(tx *sql.Tx) error {
	var query = `CREATE TABLE IF NOT EXISTS aircraft_history(
				 icao VARCHAR(6) NOT NULL,
				 lat DECIMAL NOT NULL,
				 long DECIMAL NOT NULL,
				 timestamp TIMESTAMP NOT NULL,
				 PRIMARY KEY (icao,timestamp))`

	return db.exec(query, tx)
}

func (db *AdsbDB) DropAircraftCurrentTable(tx *sql.Tx) error {
	query := `DROP TABLE IF EXISTS aircraft_current CASCADE`

	return db.exec(query, tx)
}

// BulkInsertAircraftCurrent inserts an array of new aircraft data into aircraft_current
func (db *AdsbDB) BulkInsertAircraftCurrent(aircraft []models.AircraftCurrentModel, tx *sql.Tx) error {
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
		if tx != nil {
			_, err := tx.Exec(stmt, vals...)
			if err != nil {
				return err
			}
		} else {
			_, err := db.Conn.Exec(stmt, vals...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// InsertHistoryFromCurrent inserts all data from aircraft_current table to aircraft_history.
func (db *AdsbDB) InsertHistoryFromCurrent(tx *sql.Tx) error {
	query := `INSERT INTO aircraft_history (icao, lat, long, timestamp) 
			  SELECT icao, lat, long, timestamp
			  FROM aircraft_current
			  ON CONFLICT (icao, timestamp) 
			      DO UPDATE SET timestamp = excluded.timestamp`
	return db.exec(query, tx)
}

// DeleteOldCurrent will delete rows in aircraft_current older than global.WaitingTime seconds from the latest entry.
func (db *AdsbDB) DeleteOldCurrent(tx *sql.Tx) error {
	var query = `DELETE FROM aircraft_current 
       			 WHERE timestamp < (select max(timestamp)-($1 * interval '1 second') 
                 FROM aircraft_current)`

	if tx != nil {
		_, err := tx.Exec(query, global.WaitingTime+2)
		return err
	} else {
		_, err := db.Conn.Exec(query, global.WaitingTime+2)
		return err
	}
}

// GetCurrentAircraft retrieves a list of all aircraft from aircraft_current that are older than global.WaitingTime + 2
func (db *AdsbDB) GetCurrentAircraft(tx *sql.Tx) ([]models.AircraftCurrentModel, error) {
	var query = `SELECT * FROM aircraft_current 
				 WHERE timestamp > (select max(timestamp)-($1 * interval '1 second') 
				 FROM aircraft_current)`

	var rows *sql.Rows
	var err error

	if tx != nil {
		rows, err = tx.Query(query, global.WaitingTime+2)
	} else {
		rows, err = db.Conn.Query(query, global.WaitingTime+2)
	}

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
func (db *AdsbDB) GetHistoryByIcao(search string, tx *sql.Tx) ([]models.AircraftHistoryModel, error) {
	var query = `SELECT icao, long, lat FROM aircraft_history WHERE icao = $1`

	var rows *sql.Rows
	var err error

	if tx != nil {
		rows, err = tx.Query(query, search)
	} else {
		rows, err = db.Conn.Query(query, search)
	}

	if err != nil {
		return []models.AircraftHistoryModel{}, err
	}
	defer rows.Close()

	var aircraft []models.AircraftHistoryModel

	for rows.Next() {
		var ac models.AircraftHistoryModel
		err := rows.Scan(&ac.Icao, &ac.Longitude, &ac.Latitude)
		if err != nil {
			return nil, err
		}

		aircraft = append(aircraft, ac)
	}

	return aircraft, nil
}

func (db *AdsbDB) exec(query string, tx *sql.Tx) error {
	if tx != nil {
		_, err := tx.Exec(query)
		return err
	} else {
		_, err := db.Conn.Exec(query)
		return err
	}
}
