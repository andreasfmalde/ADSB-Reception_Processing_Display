package main

import (
	"adsb-api/internal/global"
	"adsb-api/internal/logger"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	logger.InitLogger()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		global.Host, global.Port, global.User, global.Password, global.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.Error.Fatalf("Error opening database: %q", err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error.Fatalf("Error closing database: %q", err)
		}
	}(db)

	err = db.Ping()
	if err != nil {
		logger.Error.Fatalf("Error pinging database: %q", err)
	}

	logger.Info.Println("Successfully connected to database!")
}
