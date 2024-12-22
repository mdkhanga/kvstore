package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	Log zerolog.Logger
}

var Globallogger = New()

func New() *Logger {
	l := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &Logger{Log: l}
}

func (log *Logger) Info(msg string) {
	log.Log.Info().Msg(msg)
}

func WithComponent(name string) *Logger {
	child := Globallogger.Log.With().Str("component", name).Logger()
	return &Logger{Log: child}
}
