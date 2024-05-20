package main

import (
	"adsb-api/internal/db"
	"adsb-api/internal/global"
	"adsb-api/internal/global/errorMsg"
	"adsb-api/internal/handler/aircraftCurrentHandler"
	"adsb-api/internal/handler/aircraftHistoryHandler"
	"adsb-api/internal/handler/defaultHandler"
	"adsb-api/internal/service/restService"
	"adsb-api/internal/utility/logger"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

// main method for the RESTFUL API
func main() {
	// Initialize environment variables
	global.InitEnvironment()
	// Initialize logger
	logger.InitLogger()
	// Initialize the database
	database, err := db.InitDB()
	if err != nil {
		log.Fatal().Msgf("error opening database: %q", err)
	}

	defer func() {
		err = database.Close()
		if err != nil {
			log.Fatal().Msgf(errorMsg.ErrorClosingDatabase+": %q", err)
		}
	}()

	log.Info().Msgf("Reception API successfully connected to database with: User: %s | Database: %s | Host: %s | port: %d",
		global.DbUser, global.DbName, global.DbHost, global.DbPort)

	restSvc := restService.InitRestService(database)

	http.HandleFunc(global.DefaultPath, defaultHandler.DefaultHandler)
	http.HandleFunc(global.AircraftCurrentPath, aircraftCurrentHandler.CurrentAircraftHandler(restSvc))
	http.HandleFunc(global.AircraftHistoryPath, aircraftHistoryHandler.HistoryAircraftHandler(restSvc))

	port := os.Getenv("PORT")
	if port == "" {
		port = global.DefaultPort
		log.Info().Msgf("PORT has not been set. Using default port: %s", port)
	}

	log.Info().Msgf("Listening on port: " + port)
	log.Fatal().Msgf(http.ListenAndServe(":"+port, nil).Error())

}
