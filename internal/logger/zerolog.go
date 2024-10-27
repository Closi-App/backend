package logger

import (
	"github.com/Closi-App/backend/internal/config"
	"github.com/rs/zerolog"
	"os"
	"time"
)

type zerologLogger struct {
	log zerolog.Logger
}

func NewZerolog(cfg *config.Config) Logger {
	lvl, err := zerolog.ParseLevel(cfg.Log.Level)
	if err != nil {
		panic("error parsing log level: " + err.Error())
	}

	zerolog.SetGlobalLevel(lvl)

	var logger zerolog.Logger

	switch cfg.Log.Encoding {
	case "console":
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.TimeOnly,
		}).With().Timestamp().Logger()
	case "json":
		zerolog.TimeFieldFormat = time.DateTime
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	default:
		panic("log encoding can be 'console' or 'json'")
	}

	return &zerologLogger{
		log: logger,
	}
}

func (l *zerologLogger) Info(msg string) {
	l.log.Info().Msg(msg)
}

func (l *zerologLogger) Warn(msg string) {
	l.log.Warn().Msg(msg)
}

func (l *zerologLogger) Error(msg string) {
	l.log.Error().Msg(msg)
}

func (l *zerologLogger) Debug(msg string) {
	l.log.Debug().Msg(msg)
}

func (l *zerologLogger) WithField(key string, value interface{}) Logger {
	return &zerologLogger{
		log: l.log.With().Interface(key, value).Logger(),
	}
}

func (l *zerologLogger) Logger() interface{} {
	return l.log
}
