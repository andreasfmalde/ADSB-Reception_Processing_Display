package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

var logFile zerolog.Logger

func InitLogger() {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		With().
		Timestamp().
		Caller().
		Logger()

	log.Logger = logger

	env := os.Getenv("ENV")

	if env == "prod" || env == "production" {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		log.Info().Msgf("environment variable ENV was not set to production. Using logging level: %s", zerolog.GlobalLevel().String())
	}
}

func InitLogFile(filePath string) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Error().Msgf("Error opening file for logging: %v", err)
	}

	logFile = zerolog.New(file).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func LogToFile(logString string) {
	logFile.Info().Msg(logString)
}
