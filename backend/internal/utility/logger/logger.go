package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

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
