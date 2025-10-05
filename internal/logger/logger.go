package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var appLogger zerolog.Logger

func Init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	appLogger = log.With().Caller().Logger()
}

func GetLogger() *zerolog.Logger {
	return &appLogger
}
