package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"social-platform-backend/config"
)

var (
	defaultLogger *slog.Logger
	levelVar      slog.LevelVar
	rotateOut     *RotateWriter
	initOnce      sync.Once
	initErr       error
)

// Init configures the application logger (structured slog, optional console + rotating file).
// Safe to call once at process startup; concurrent writers are serialized by RotateWriter.
func Init(cfg *config.Log) error {
	initOnce.Do(func() {
		lvl := parseLevel(cfg.Level)
		levelVar.Set(lvl)

		maxMB := cfg.MaxSizeMB
		if maxMB <= 0 {
			maxMB = 10
		}
		maxBytes := int64(maxMB) * 1024 * 1024

		filePath := strings.TrimSpace(cfg.FilePath)
		if filePath == "" {
			filePath = "logs/app.log"
		}
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			initErr = err
			return
		}

		rw, err := NewRotateWriter(absPath, maxBytes)
		if err != nil {
			initErr = err
			return
		}
		rotateOut = rw

		addSource := strings.EqualFold(strings.TrimSpace(cfg.Level), "debug")
		opts := &slog.HandlerOptions{
			Level:     &levelVar,
			AddSource: addSource,
		}

		fileHandler := slog.NewJSONHandler(rw, opts)

		var handler slog.Handler = fileHandler
		if cfg.Console {
			consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level:     &levelVar,
				AddSource: addSource,
			})
			handler = newTeeHandler(fileHandler, consoleHandler)
		}

		defaultLogger = slog.New(handler)
		slog.SetDefault(defaultLogger)
	})
	return initErr
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// LogFilePath is the absolute path of the active log file (for admin tail API).
func LogFilePath() string {
	if rotateOut != nil {
		return rotateOut.Path()
	}
	return ""
}

// Sync flushes the log file to disk.
func Sync() error {
	if rotateOut != nil {
		return rotateOut.Sync()
	}
	return nil
}

func log() *slog.Logger {
	if defaultLogger != nil {
		return defaultLogger
	}
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

// Debugf logs at debug level after formatting the message.
func Debugf(format string, args ...any) {
	log().Debug(fmt.Sprintf(format, args...))
}

// Infof logs at info level after formatting the message.
func Infof(format string, args ...any) {
	log().Info(fmt.Sprintf(format, args...))
}

// Warnf logs at warn level after formatting the message.
func Warnf(format string, args ...any) {
	log().Warn(fmt.Sprintf(format, args...))
}

// Errorf logs at error level after formatting the message.
func Errorf(format string, args ...any) {
	log().Error(fmt.Sprintf(format, args...))
}

// Fatalf logs at error level, flushes, then exits the process with code 1.
func Fatalf(format string, args ...any) {
	Errorf(format, args...)
	_ = Sync()
	os.Exit(1)
}
