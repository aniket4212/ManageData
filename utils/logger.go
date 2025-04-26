package utils

import (
	"os"

	"github.com/rs/zerolog"
)

func SetLogLevel() {
	logLevel := zerolog.InfoLevel
	zerolog.SetGlobalLevel(logLevel)
}

func GetLogger() zerolog.Logger {
	return zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()
}
