package logger

import (
	"context"
	"fmt"

	"github.com/ubin/go-telemetry/logger/loggerfactory/config"
	"github.com/ubin/go-telemetry/logger/loggerfactory/defaultlogger"
)

// Log is a package level variable, every program should access logging function through "Log"
var Log Logger

func init() {
	if Log != nil {
		return
	}
	// //set the Log to default logger (slog) with default configuration
	var err error
	if Log, err = defaultlogger.New(config.LogEnvDev, defaultlogger.Config{}); err != nil {
		fmt.Println("Unable to initialize default logger")
		return
	}
	Log.Info("default logger [slog] initialized")
}

type Config interface {
	// log library name
	GetCode() string
	GetFormatter() string
	GetLevel() string
	GetEnableCaller() bool
	GetFileEnabled() bool
	GetFilename() string
	GetMaxSize() int
	GetMaxBackups() int
	GetMaxAge() int
	GetCompress() bool
	GetLocalTime() bool
}

// BasicLogger represents common logging methods without context
type BasicLogger interface {
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
	Debug(msg string, keyvals ...interface{})
	Panic(msg string, keyvals ...interface{})
}

// ContextLogger represents logging methods that accept a context
type ContextLogger interface {
	InfoContext(ctx context.Context, msg string, keyvals ...interface{})
	WarnContext(ctx context.Context, msg string, keyvals ...interface{})
	ErrorContext(ctx context.Context, msg string, keyvals ...interface{})
	DebugContext(ctx context.Context, msg string, keyvals ...interface{})
	PanicContext(ctx context.Context, msg string, keyvals ...interface{})
}

// UnderlyingLoggerProvider represents the method to retrieve the underlying logger library
type UnderlyingLoggerProvider interface {
	UnderlyingLogger() interface{}
}

// Logger combines all three interfaces into a single interface for convenience
type Logger interface {
	BasicLogger
	ContextLogger
	UnderlyingLoggerProvider
}

// SetLogger is the setter for log variable, it should be the only way to assign value to log
func SetLogger(newLogger Logger) {
	Log = newLogger
}
