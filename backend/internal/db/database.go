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
