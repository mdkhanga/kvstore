package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	Log zerolog.Logger
}

func New() *Logger {
	l := zerolog.New(os.Stdout)
	return &Logger{Log: l}
}

func (log *Logger) Info(msg string) {
	log.Log.Info().Msg(msg)
}
