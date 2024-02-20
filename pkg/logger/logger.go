package logger

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

type Interface interface {
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
	Fatal(message string, args ...interface{})
	Fields(fields ...map[string]interface{}) Interface
}

type logger struct {
	logger zerolog.Logger
}

type Kind int

const (
	Console Kind = iota
	JSON
)

func New(level string, kind Kind) Interface {
	var l zerolog.Logger

	switch kind {
	case Console:
		l = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.DateTime,
		}).With().Timestamp().Logger()
	case JSON:
		l = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		panic("error parsing log level: " + err.Error())
	}

	zerolog.SetGlobalLevel(lvl)

	return &logger{
		logger: l,
	}
}

// Debug -
func (l *logger) Debug(message string, args ...interface{}) {
	if len(args) > 0 {
		l.logger.Debug().Msgf(message, args...)
	} else {
		l.logger.Debug().Msg(message)
	}
}

// Info -
func (l *logger) Info(message string, args ...interface{}) {
	if len(args) > 0 {
		l.logger.Info().Msgf(message, args...)
	} else {
		l.logger.Info().Msg(message)
	}
}

// Warn -
func (l *logger) Warn(message string, args ...interface{}) {
	if len(args) > 0 {
		l.logger.Warn().Msgf(message, args...)
	} else {
		l.logger.Warn().Msg(message)
	}
}

// Error -
func (l *logger) Error(message string, args ...interface{}) {
	if len(args) > 0 {
		l.logger.Error().Msgf(message, args...)
	} else {
		l.logger.Error().Msg(message)
	}
}

// Fatal -
func (l *logger) Fatal(message string, args ...interface{}) {
	if len(args) > 0 {
		l.logger.Fatal().Msgf(message, args...)
	} else {
		l.logger.Fatal().Msg(message)
	}
}

// Fields adds fields for the log
func (l *logger) Fields(fields ...map[string]interface{}) Interface {
	if len(fields) == 0 {
		return l
	}

	return &logger{
		logger: l.logger.With().Fields(fields[0]).Logger(),
	}
}
