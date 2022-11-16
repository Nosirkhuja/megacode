package logging

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"
)

func NewLogger() *zerolog.Logger {
	out := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	out.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("[%s]", i))
	}
	logger := zerolog.New(out).Level(zerolog.DebugLevel).
		With().
		Timestamp().
		Caller().
		Logger()
	return &logger
}
