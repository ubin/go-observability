package defaultlogger

import (
	"github.com/ubin/go-telemetry/logger/loggerfactory/config"
	"github.com/ubin/go-telemetry/logger/loggerfactory/slog"
)

type Config struct {
	Code         string `koanf:"code"`
	Formatter    string `koanf:"formatter"`
	Level        string `koanf:"level"`
	EnableCaller bool   `koanf:"enable_caller"`
	FileEnabled  bool   `koanf:"file_enabled"`
	Filename     string `koanf:"filename"`
	MaxSize      int    `koanf:"max_size"`
	MaxBackups   int    `koanf:"max_backups"`
	MaxAge       int    `koanf:"max_age"`
	Compress     bool   `koanf:"compress"`
	LocalTime    bool   `koanf:"local_time"`
}

// GetCode returns the code level we filter by
func (cfg Config) GetFormatter() string {
	return cfg.Formatter
}

// GetCode returns the code level we filter by
func (cfg Config) GetCode() string {
	return cfg.Code
}

// GetLevel returns the level
func (cfg Config) GetLevel() string {
	return cfg.Level
}

// GetEnableCaller returns if we want to enable caller in output
func (cfg Config) GetEnableCaller() bool {
	return cfg.EnableCaller
}

// GetFileEnabled if we want to show file
func (cfg Config) GetFileEnabled() bool {
	return cfg.FileEnabled
}

// GetFilename if we want to show filename
func (cfg Config) GetFilename() string {
	return cfg.Filename
}

// GetMaxSize of log
func (cfg Config) GetMaxSize() int {
	return cfg.MaxSize
}

// GetMaxBackups max backups
func (cfg Config) GetMaxBackups() int {
	return cfg.MaxBackups
}

// GetMaxAge of log
func (cfg Config) GetMaxAge() int {
	return cfg.MaxAge
}

// GetCompress should we compress
func (cfg Config) GetCompress() bool {
	return cfg.Compress
}

// GetLocalTime if we store in local time
func (cfg Config) GetLocalTime() bool {
	return cfg.LocalTime
}

func New(env config.LogEnv, cfg Config) (slog.LoggerWrapper, error) {
	//TODO : customizations
	return slog.New(config.LogEnvDev, cfg)
}
