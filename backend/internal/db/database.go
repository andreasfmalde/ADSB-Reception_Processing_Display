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
	DropAircraftCurrentTable() error
	BulkInsertAircraftCurrent(aircraft []models.AircraftCurrentModel) error
	SelectAllColumnsAircraftCurrent() ([]models.AircraftCurrentModel, error)

	CreateAircraftHistoryTable() error
	CreateAircraftHistoryTimestampIndex() error
	InsertHistoryFromCurrent() error
	SelectAllColumnHistoryByIcao(search string) ([]models.AircraftHistoryModel, error)
	SelectAllColumnHistoryByIcaoFilterByTimestamp(search string, hour int) ([]models.AircraftHistoryModel, error)

	DeleteOldHistory(time int) error

	Begin() error
	Commit() error
	Rollback() error

	Close() error
}

type Context struct {
	db *sql.DB
	tx *sql.Tx
}

func (ctx *Context) Exec(query string, args ...interface{}) (sql.Result, error) {
	if ctx.tx != nil {
		return ctx.tx.Exec(query, args...)
	}
	return ctx.db.Exec(query, args...)
}

func (ctx *Context) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if ctx.tx != nil {
		return ctx.tx.Query(query, args...)
	}
	return ctx.db.Query(query, args...)
}

// Begin begins Context transaction
func (ctx *Context) Begin() error {
	tx, err := ctx.db.Begin()
	if err != nil {
		return err
	}
	ctx.tx = tx
	return nil
}

// Commit commits Context transaction
func (ctx *Context) Commit() error {
	err := ctx.tx.Commit()
	if err != nil {
		return err
	}
	ctx.tx = nil
	return nil
}

// Rollback rollbacks Context transaction
func (ctx *Context) Rollback() error {
	err := ctx.tx.Rollback()
	if err != nil {
		return err
	}
	ctx.tx = nil
	return nil
}

// InitDB initializes the PostgresSQL database and returns the connection pointer.
func InitDB() (*Context, error) {
	dbLogin := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		global.DbHost, global.DbPort, global.DbUser, global.DbPassword, global.DbName)

	dbConn, err := sql.Open("postgres", dbLogin)
	if err = dbConn.Ping(); err != nil {
		return nil, err
	}
	return &Context{db: dbConn}, err
}

// Close closes Context db connection
func (ctx *Context) Close() error {
	return ctx.db.Close()
}

// CreateAircraftCurrentTable creates a table for storing current aircraft data if it does not already exist
func (ctx *Context) CreateAircraftCurrentTable() error {
	query := `CREATE TABLE IF NOT EXISTS aircraft_current(
				 icao VARCHAR(6) NOT NULL,
				 callsign VARCHAR(10) NOT NULL,
				 altitude INT NOT NULL,
				 lat DECIMAL NOT NULL,
				 long DECIMAL NOT NULL,
				 speed INT NOT NULL,
				 track INT NOT NULL,
				 vspeed INT NOT NULL,
				 timestamp TIMESTAMP NOT NULL,
				 PRIMARY KEY (icao))`
	_, err := ctx.Exec(query)
	return err
}

// CreateAircraftHistoryTimestampIndex creates an index called timestamp_index on aircraft_history timestamp column
func (ctx *Context) CreateAircraftHistoryTimestampIndex() error {
	query := `CREATE INDEX IF NOT EXISTS timestamp_index ON aircraft_history(timestamp)`

	_, err := ctx.Exec(query)
	return err
}

// CreateAircraftHistoryTable creates a table for storing aircraft history data if it does not already exist
func (ctx *Context) CreateAircraftHistoryTable() error {
	query := `CREATE TABLE IF NOT EXISTS aircraft_history(
				 icao VARCHAR(6) NOT NULL,
				 lat DECIMAL NOT NULL,
				 long DECIMAL NOT NULL,
				 timestamp TIMESTAMP NOT NULL,
				 PRIMARY KEY (icao,timestamp))`

	_, err := ctx.Exec(query)
	return err
}

// DropAircraftCurrentTable drops aircraft_current table
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
		_, err := ctx.Exec(stmt, vals...)
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

// SelectAllColumnsAircraftCurrent retrieves a list of all aircraft from aircraft_current that are older than global.WaitingTime + 2
func (ctx *Context) SelectAllColumnsAircraftCurrent() ([]models.AircraftCurrentModel, error) {
	query := `SELECT * FROM aircraft_current`

	rows, err := ctx.Query(query)
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
	query := `SELECT * FROM aircraft_history WHERE icao = $1`

	rows, err := ctx.Query(query, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aircraft []models.AircraftHistoryModel

	for rows.Next() {
		var ac models.AircraftHistoryModel
		err := rows.Scan(&ac.Icao, &ac.Latitude, &ac.Longitude, &ac.Timestamp)
		if err != nil {
			return nil, err
		}

		aircraft = append(aircraft, ac)
	}

	return aircraft, nil
}

// DeleteOldHistory will delete rows in aircraft_history older than MAX(timestamp) days
func (ctx *Context) DeleteOldHistory(days int) error {
	query := `DELETE FROM aircraft_history
			  WHERE timestamp <
			        (SELECT MAX(timestamp) - ($1 * INTERVAL '1 day') 
  			  		 FROM aircraft_history)`

	_, err := ctx.Exec(query, days)
	return err
}

// SelectAllColumnHistoryByIcaoFilterByTimestamp selects history by aircraft icao
// and filters every row with a newer timestamp than given hour
func (ctx *Context) SelectAllColumnHistoryByIcaoFilterByTimestamp(search string, hour int) ([]models.AircraftHistoryModel, error) {
	query := `SELECT * FROM aircraft_history 
         		 WHERE icao = $1 AND timestamp > (NOW() - ($2 * INTERVAL '1 hour')) 
         		 ORDER BY timestamp DESC;`

	rows, err := ctx.Query(query, search, hour)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aircraft []models.AircraftHistoryModel

	for rows.Next() {
		var ac models.AircraftHistoryModel
		err := rows.Scan(&ac.Icao, &ac.Latitude, &ac.Longitude, &ac.Timestamp)
		if err != nil {
			return nil, err
		}

		aircraft = append(aircraft, ac)
	}

	return aircraft, nil
}
