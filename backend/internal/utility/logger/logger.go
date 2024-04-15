package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

// InitLogger initializes the logger for the application.
// It sets the logger as the default logger used throughout the application.
// It checks the value of the "ENV" environment variable, and if it is set to "prod"
// or "production", it sets the global log level to WarnLevel.
// Otherwise, it sets the global log level to TraceLevel
// and logs an informational message indicating the logging level that will be used.
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
