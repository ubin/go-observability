package logrus

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/ubin/go-observability/logger"
	"github.com/ubin/go-observability/logger/loggerfactory/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LoggerWrapper is a logger that uses the logrus package
type LoggerWrapper struct {
	logger *logrus.Logger
}

func (l LoggerWrapper) UnderlyingLogger() interface{} {
	return l.logger

}

func toFields(keyvals ...interface{}) logrus.Fields {
	fields := make(logrus.Fields, len(keyvals)/2)
	for i := 0; i < len(keyvals); i += 2 {
		// Handle odd-length keyvals slices
		if i+1 >= len(keyvals) {
			// If we have an odd number of elements, use the key with a nil value
			if keyStr, ok := keyvals[i].(string); ok {
				fields[keyStr] = nil
			}
			break
		}
		key, val := keyvals[i], keyvals[i+1]
		if keyStr, ok := key.(string); ok {
			fields[keyStr] = val
		}
	}
	return fields
}
func (l LoggerWrapper) Info(msg string, keyvals ...interface{}) {
	l.InfoContext(context.Background(), msg, keyvals...)
}

func (l LoggerWrapper) Warn(msg string, keyvals ...interface{}) {
	l.WarnContext(context.Background(), msg, keyvals...)
}

func (l LoggerWrapper) Error(msg string, keyvals ...interface{}) {
	l.ErrorContext(context.Background(), msg, keyvals...)
}

func (l LoggerWrapper) Debug(msg string, keyvals ...interface{}) {
	l.DebugContext(context.Background(), msg, keyvals...)
}

func (l LoggerWrapper) Panic(msg string, keyvals ...interface{}) {
	l.PanicContext(context.Background(), msg, keyvals...)
}

func (l LoggerWrapper) InfoContext(ctx context.Context, msg string, keyvals ...interface{}) {
	l.logger.WithContext(ctx).WithFields(toFields(keyvals...)).Info(msg)
}

func (l LoggerWrapper) WarnContext(ctx context.Context, msg string, keyvals ...interface{}) {
	l.logger.WithContext(ctx).WithFields(toFields(keyvals...)).Warn(msg)
}

func (l LoggerWrapper) ErrorContext(ctx context.Context, msg string, keyvals ...interface{}) {
	l.logger.WithContext(ctx).WithFields(toFields(keyvals...)).Error(msg)
}

func (l LoggerWrapper) DebugContext(ctx context.Context, msg string, keyvals ...interface{}) {
	l.logger.WithContext(ctx).WithFields(toFields(keyvals...)).Debug(msg)
}
func (l LoggerWrapper) PanicContext(ctx context.Context, msg string, keyvals ...interface{}) {
	l.logger.WithContext(ctx).WithFields(toFields(keyvals...)).Panic(msg)
}

func New(env config.LogEnv, cfg logger.Config) (logger.Logger, error) {
	w := io.Writer(os.Stdout)

	if cfg.GetFileEnabled() {
		logWriter := &lumberjack.Logger{
			Filename:   cfg.GetFilename(),
			MaxSize:    cfg.GetMaxSize(),
			MaxBackups: cfg.GetMaxBackups(),
			MaxAge:     cfg.GetMaxAge(),
			Compress:   cfg.GetCompress(),
			LocalTime:  cfg.GetLocalTime(),
		}

		logrus.RegisterExitHandler(func() {
			_ = logWriter.Close()
		})

		w = io.MultiWriter(os.Stdout, logWriter)
	}

	formatter := getLoggerFormatter(env)

	rus := logrus.New()
	rus.SetOutput(w)
	rus.SetFormatter(formatter)
	rus.SetReportCaller(cfg.GetEnableCaller())

	err := customizeLogFromConfig(rus, cfg)
	if err != nil {
		return nil, err
	}

	// logger := rus.WithField("logger", "app")
	logger := LoggerWrapper{rus}
	logger.Warn("Logrus initialized...")
	log.SetOutput(logger.logger.Writer())

	return logger, nil
}

func getLoggerFormatter(env config.LogEnv) logrus.Formatter {
	if env == config.LogEnvProd {
		return &logrus.JSONFormatter{}
	}
	return &logrus.TextFormatter{
		TimestampFormat:        "02-01-2006 15:04:05",
		FullTimestamp:          true, //prod,
		PadLevelText:           true,
		DisableLevelTruncation: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", formatFilePath(f.File), f.Line)
		},
	}
}

func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]
}

func customizeLogFromConfig(log *logrus.Logger, cfg logger.Config) error {
	l := &log.Level
	err := l.UnmarshalText([]byte(cfg.GetLevel()))
	if err != nil {
		return err
	}
	log.SetLevel(*l)
	return nil
}
