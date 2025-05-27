package loggerfactory

import (
	"fmt"

	"github.com/ubin/go-telemetry/logger"
	logconfig "github.com/ubin/go-telemetry/logger/loggerfactory/config"
	"github.com/ubin/go-telemetry/logger/loggerfactory/logrus"
	"github.com/ubin/go-telemetry/logger/loggerfactory/slog"
)

// Register initializes the logger based on the configuration.
func Register(cfg logger.Config, env string) error {

	logEnv := logconfig.LogEnvDev
	if env == "production" || env == "prod" {
		logEnv = logconfig.LogEnvProd
	}

	lgr, err := getLogger(cfg, logEnv)
	if err != nil {
		return fmt.Errorf("error initializing logger: %s", err)
	}
	logger.SetLogger(lgr)
	return nil
}

func getLogger(cfg logger.Config, env logconfig.LogEnv) (logger.Logger, error) {
	switch cfg.GetCode() {
	case logconfig.LOGRUS:
		logrusFactory := &LogrusFactory{}
		return logrusFactory.CreateLogger(cfg, env)
	case logconfig.SLOG:
		logrusFactory := &SlogFactory{}
		return logrusFactory.CreateLogger(cfg, env)
	default:
		return nil, fmt.Errorf("unsupported log provider: %s", cfg.GetCode())
	}
}

// Factory is the creator interface.
type Factory interface {
	CreateLogger(cfg logger.Config) (logger.Logger, error)
}

// LogrusFactory is a factory for Logrus logger.
type LogrusFactory struct{}

// CreateLogger creates a new Logrus logger.
func (f *LogrusFactory) CreateLogger(cfg logger.Config, env logconfig.LogEnv) (logger.Logger, error) {
	return logrus.New(env, cfg)
}

// LogrusFactory is a factory for Logrus logger.
type SlogFactory struct{}

// CreateLogger creates a new Logrus logger.
func (f *SlogFactory) CreateLogger(cfg logger.Config, env logconfig.LogEnv) (logger.Logger, error) {
	return slog.New(env, cfg)
}
