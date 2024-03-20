package main

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global"
	"adsb-api/internal/handler/aircraftCurrentHandler"
	"adsb-api/internal/handler/defaultHandler"
	"adsb-api/internal/logger"
	"log"
	"net/http"
	"os"
)

// main method for the RESTFUL API
func main() {
	// Initialize logger
	logger.InitLogger()
	// Initialize environment variables
	global.InitEnvironment()
	// Initialize the database
	adsbDB, err := db.InitDB()
	if err != nil {
		logger.Error.Fatalf("error opening database: %q", err)
	}
	logger.Info.Println("successfully connected to database")

	defer func() {
		err := adsbDB.Close()
		if err != nil {
			logger.Error.Fatalf("error closing database: %q", err)
		}
	}()

	http.HandleFunc(global.DefaultPath, defaultHandler.DefaultHandler)
	http.HandleFunc(global.CurrentAircraftPath, aircraftCurrentHandler.CurrentAircraftHandler(adsbDB))

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: " + global.DefaultPort)
		port = global.DefaultPort
	}

	logger.Info.Println("Listening on port: " + port)
	logger.Info.Fatal(http.ListenAndServe(":"+port, nil))

}
