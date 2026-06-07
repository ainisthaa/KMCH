package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/lumberjack.v2"
)

// Log is the application-wide structured logger.
var Log zerolog.Logger

// Init sets up zerolog to write to both stdout (pretty) and a rotating file.
// Call once from main before starting services.
func Init(logFilePath string) {
	if err := os.MkdirAll(filepath.Dir(logFilePath), 0o755); err != nil {
		panic("logger: cannot create log directory: " + err.Error())
	}

	fileWriter := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxAge:     30, // keep 30 days
		MaxBackups: 30,
		Compress:   true,
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	multi := io.MultiWriter(consoleWriter, fileWriter)
	Log = zerolog.New(multi).With().Timestamp().Logger()
	Log.Info().Str("action", "logger_init").Str("file", logFilePath).Msg("logger initialised")
}
