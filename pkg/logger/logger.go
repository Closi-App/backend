package logger

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"os"
	"time"
)

type Logger struct {
	zerolog.Logger
}

func NewLogger(cfg *viper.Viper) *Logger {
	lvl, err := zerolog.ParseLevel(cfg.GetString("log.level"))
	if err != nil {
		panic("error parsing log level: " + err.Error())
	}

	zerolog.SetGlobalLevel(lvl)

	var logger zerolog.Logger

	switch cfg.GetString("log.format") {
	case "console":
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.TimeOnly,
		}).With().Timestamp().Logger()
	case "json":
		zerolog.TimeFieldFormat = time.DateTime
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	default:
		panic("log format can be 'console' or 'json'")
	}

	return &Logger{
		Logger: logger,
	}
}
