package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger zerolog.Logger
}

func New() *Logger {
	l := zerolog.New(os.Stdout)
	return &Logger{logger: l}
}
