package main

import (
	"adsb-api/internal/global"
	"adsb-api/internal/handler/aircraftCurrentHandler"
	"adsb-api/internal/handler/aircraftHistory"
	"adsb-api/internal/handler/defaultHandler"
	"adsb-api/internal/service/restService"
	"adsb-api/internal/utility/logger"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

// main method for the RESTFUL API
func main() {
	// Initialize environment variables
	global.InitEnvironment()
	// Initialize logger
	logger.InitLogger()
	// Initialize the database
	restSvc, err := restService.InitRestService()
	if err != nil {
		log.Fatal().Msgf("error opening database: %q", err)
	}
	log.Info().Msgf("REST API successfully connected to database with database user: %s name: %s host: %s port: %d",
		global.DbUser, global.DbName, global.DbHost, global.DbPort)

	defer func() {
		err := restSvc.DB.Close()
		if err != nil {
			log.Fatal().Msgf("error closing database: %q", err)
		}
	}()

	http.HandleFunc(global.DefaultPath, defaultHandler.DefaultHandler)
	http.HandleFunc(global.AircraftCurrentPath, aircraftCurrentHandler.CurrentAircraftHandler(restSvc))
	http.HandleFunc(global.AircraftHistoryPath, aircraftHistory.HistoryAircraftHandler(restSvc))

	port := os.Getenv("PORT")
	if port == "" {
		port = global.DefaultPort
		log.Info().Msg("PORT has not been set. Using default port: " + port)
	}

	log.Info().Msgf("Listening on port: " + port)
	log.Fatal().Msg(http.ListenAndServe(":"+port, nil).Error())

}
