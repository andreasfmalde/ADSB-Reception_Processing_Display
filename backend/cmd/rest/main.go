package main

import (
	"adsb-api/cmd"
	"adsb-api/internal/global"
	"adsb-api/internal/handler/currentAircraftHandler"
	"adsb-api/internal/handler/defaultHandler"
	"adsb-api/internal/logger"
	_ "github.com/lib/pq"
	"net/http"
)

// main method for the RESTFUL API
func main() {
	dbConn := cmd.InitializeApiResources()

	http.HandleFunc(global.DefaultPath, defaultHandler.DefaultHandler)
	http.HandleFunc(global.CurrentAircraftPath, currentAircraftHandler.CurrentAircraftHandler(dbConn))

	logger.Info.Println("Listening on port: 8080 ...")
	logger.Info.Fatal(http.ListenAndServe(":"+global.DefaultPort, nil))
}
