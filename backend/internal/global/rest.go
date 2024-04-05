package global

import "strings"

// API constants
const (
	DefaultPort         = "8080"
	VERSION             = "1.0.0"
	DefaultPath         = "/"
	AircraftCurrentPath = "/aircraft/current/"
	AircraftHistoryPath = "/aircraft/history/"
)

var (
	CurrentPathMaxLength = len(strings.Split(AircraftCurrentPath, "/")) - 1
	HistoryPathMaxLength = len(strings.Split(AircraftHistoryPath, "/"))
)
