package slog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	// "github.com/Garchen-Archive/garchen-archive/pkg/logger"

	"github.com/ubin/go-observability/logger/loggerfactory/config"
	otelslog "github.com/ubin/go-observability/telemetry/log/slog"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	TEXT_FORMATTER = "TEXT"
	JSON_FORMATTER = "JSON"
)

type Config interface {
	GetFormatter() string
	GetLevel() string
	GetFileEnabled() bool
	//File name is expected to be a path with "/" as separator
	GetFilename() string
	GetMaxSize() int
	GetMaxBackups() int
	GetMaxAge() int
	GetCompress() bool
	GetLocalTime() bool
}

// custom level for panic, as slog doesn't define panic level by default
const LevelPanic = slog.Level(15)

// LoggerWrapper is a logger that uses the logrus package
type LoggerWrapper struct {
	lgr *slog.Logger
}

func (l LoggerWrapper) Info(msg string, keyvals ...interface{}) {
	l.InfoContext(context.Background(), msg, keyvals...)
}
func (l LoggerWrapper) Warn(msg string, keyvals ...interface{}) {
	l.WarnContext(context.Background(), msg, keyvals...)

}
func (l LoggerWrapper) Debug(msg string, keyvals ...interface{}) {
	l.DebugContext(context.Background(), msg, keyvals...)

}
func (l LoggerWrapper) Error(msg string, keyvals ...interface{}) {
	l.ErrorContext(context.Background(), msg, keyvals...)

}
func (l LoggerWrapper) Panic(msg string, keyvals ...interface{}) {
	l.PanicContext(context.Background(), msg, keyvals...)
}

func (l LoggerWrapper) InfoContext(ctx context.Context, msg string, keyvals ...interface{}) {
	otelslog.AddLogToSpan(ctx, slog.LevelInfo, msg, keyvals...)
	l.lgr.Log(ctx, slog.LevelInfo, msg, keyvals...)
}
func (l LoggerWrapper) WarnContext(ctx context.Context, msg string, keyvals ...interface{}) {
	otelslog.AddLogToSpan(ctx, slog.LevelWarn, msg, keyvals...)
	l.lgr.Log(ctx, slog.LevelWarn, msg, keyvals...)
}
func (l LoggerWrapper) DebugContext(ctx context.Context, msg string, keyvals ...interface{}) {
	otelslog.AddLogToSpan(ctx, slog.LevelDebug, msg, keyvals...)
	l.lgr.Log(ctx, slog.LevelDebug, msg, keyvals...)

}
func (l LoggerWrapper) ErrorContext(ctx context.Context, msg string, keyvals ...interface{}) {
	otelslog.AddLogToSpan(ctx, slog.LevelError, msg, keyvals...)
	l.lgr.Log(ctx, slog.LevelError, msg, keyvals...)

}
func (l LoggerWrapper) PanicContext(ctx context.Context, msg string, keyvals ...interface{}) {
	otelslog.AddLogToSpan(ctx, LevelPanic, msg, keyvals...)
	l.lgr.Log(ctx, LevelPanic, msg, keyvals...)
	panic(msg)
}

func (l LoggerWrapper) UnderlyingLogger() interface{} {
	return l.lgr

}

func New(env config.LogEnv, cfg Config) (LoggerWrapper, error) {
	//TODO : customizations
	level, err := ParseLevel(cfg.GetLevel())
	if err != nil {
		level = slog.LevelInfo
	}

	options := &slog.HandlerOptions{
		Level: level,
	}

	w := io.Writer(os.Stdout)
	if cfg.GetFileEnabled() {
		logWriter := &lumberjack.Logger{
			Filename:   filepath.FromSlash(cfg.GetFilename()),
			MaxSize:    cfg.GetMaxSize(),
			MaxBackups: cfg.GetMaxBackups(),
			MaxAge:     cfg.GetMaxAge(),
			Compress:   cfg.GetCompress(),
			LocalTime:  cfg.GetLocalTime(),
		}

		w = io.MultiWriter(os.Stdout, logWriter)
	}

	var handler slog.Handler
	switch strings.ToUpper(cfg.GetFormatter()) {
	case JSON_FORMATTER:
		// handler = slog.NewJSONHandler(w, options)
		handler = otelslog.NewOtelHandler(slog.NewJSONHandler(w, options))
	default:
		// handler = slog.NewTextHandler(w, options)
		handler = otelslog.NewOtelHandler(slog.NewTextHandler(w, options))

	}

	sl := slog.New(handler)

	lgr := LoggerWrapper{sl} //.WithGroup("app")

	lgr.Warn("Slog initialized...")

	return lgr, nil
}

func ParseLevel(s string) (slog.Level, error) {
	var level slog.Level
	var err = level.UnmarshalText([]byte(s))
	return level, err
}
