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
	CreateAircraftCurrentTable() error
	CreateAircraftCurrentTimestampIndex() error
	DropAircraftCurrentTable() error
	BulkInsertAircraftCurrent(aircraft []models.AircraftCurrentModel) error
	SelectAllColumnsAircraftCurrent() ([]models.AircraftCurrentModel, error)

	CreateAircraftHistoryTable() error
	InsertHistoryFromCurrent() error
	SelectAllColumnHistoryByIcao(search string) ([]models.AircraftHistoryModel, error)

	DeleteOldCurrent() error

	BeginTx() error
	Commit() error
	Rollback() error

	Close() error
}

type Context struct {
	conn *sql.DB
	tx   *sql.Tx
}

func (ctx *Context) Exec(query string, args ...interface{}) (sql.Result, error) {
	if ctx.tx != nil {
		return ctx.tx.Exec(query, args...)
	}
	return ctx.conn.Exec(query, args...)
}

func (ctx *Context) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if ctx.tx != nil {
		return ctx.tx.Query(query, args...)
	}
	return ctx.conn.Query(query, args...)
}

func (ctx *Context) BeginTx() error {
	var err error
	ctx.tx, err = ctx.conn.Begin()
	return err
}

func (ctx *Context) Commit() error {
	if ctx.tx != nil {
		err := ctx.tx.Commit()
		ctx.tx = nil
		return err
	}
	return nil
}

func (ctx *Context) Rollback() error {
	if ctx.tx != nil {
		err := ctx.tx.Rollback()
		ctx.tx = nil
		return err
	}
	return nil
}

// InitDB initializes the PostgresSQL database and returns the connection pointer.
func InitDB() (*Context, error) {
	dbLogin := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		global.Host, global.Port, global.DbUser, global.DbPassword, global.Dbname)

	dbConn, err := sql.Open("postgres", dbLogin)
	// TODO: add ping
	return &Context{conn: dbConn}, err
}

func (ctx *Context) Close() error {
	return ctx.conn.Close()
}

// CreateAircraftCurrentTable creates a table for storing current aircraft data if it does not already exist
func (ctx *Context) CreateAircraftCurrentTable() error {
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
	_, err := ctx.Exec(query)
	return err
}

func (ctx *Context) CreateAircraftCurrentTimestampIndex() error {
	var query = `CREATE INDEX IF NOT EXISTS timestamp_index ON aircraft_current(timestamp)`

	_, err := ctx.Exec(query)
	return err
}

// CreateAircraftHistoryTable creates a table for storing aircraft history data if it does not already exist
func (ctx *Context) CreateAircraftHistoryTable() error {
	var query = `CREATE TABLE IF NOT EXISTS aircraft_history(
				 icao VARCHAR(6) NOT NULL,
				 lat DECIMAL NOT NULL,
				 long DECIMAL NOT NULL,
				 timestamp TIMESTAMP NOT NULL,
				 PRIMARY KEY (icao,timestamp))`

	_, err := ctx.Exec(query)
	return err
}

func (ctx *Context) DropAircraftCurrentTable() error {
	query := `DROP TABLE IF EXISTS aircraft_current CASCADE`

	_, err := ctx.Exec(query)
	return err
}

// BulkInsertAircraftCurrent inserts an array of new aircraft data into aircraft_current
func (ctx *Context) BulkInsertAircraftCurrent(aircraft []models.AircraftCurrentModel) error {
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
		_, err := ctx.conn.Exec(stmt, vals...)
		if err != nil {
			return err
		}
	}

	return nil
}

// InsertHistoryFromCurrent inserts all data from aircraft_current table to aircraft_history.
func (ctx *Context) InsertHistoryFromCurrent() error {
	query := `INSERT INTO aircraft_history (icao, lat, long, timestamp) 
			  SELECT icao, lat, long, timestamp
			  FROM aircraft_current
			  ON CONFLICT (icao, timestamp) 
			      DO UPDATE SET timestamp = excluded.timestamp`
	_, err := ctx.Exec(query)
	return err
}

// DeleteOldCurrent will delete rows in aircraft_current older than global.WaitingTime seconds from the latest entry.
func (ctx *Context) DeleteOldCurrent() error {
	var query = `DELETE FROM aircraft_current 
       			 WHERE timestamp < (select max(timestamp)-($1 * interval '1 second') 
                 FROM aircraft_current)`

	_, err := ctx.Exec(query, global.WaitingTime+2)
	return err
}

// SelectAllColumnsAircraftCurrent retrieves a list of all aircraft from aircraft_current that are older than global.WaitingTime + 2
func (ctx *Context) SelectAllColumnsAircraftCurrent() ([]models.AircraftCurrentModel, error) {
	var query = `SELECT * FROM aircraft_current`

	rows, err := ctx.Query(query, global.WaitingTime+2)
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

// SelectAllColumnHistoryByIcao retrieves a list from aircraft_history of rows matching the icao parameter.
func (ctx *Context) SelectAllColumnHistoryByIcao(search string) ([]models.AircraftHistoryModel, error) {
	var query = `SELECT * FROM aircraft_history WHERE icao = $1`

	rows, err := ctx.conn.Query(query, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aircraft []models.AircraftHistoryModel

	for rows.Next() {
		var ac models.AircraftHistoryModel
		err := rows.Scan(&ac.Icao, &ac.Longitude, &ac.Latitude, &ac.Timestamp)
		if err != nil {
			return nil, err
		}

		aircraft = append(aircraft, ac)
	}

	return aircraft, nil
}
