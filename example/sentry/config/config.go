package config

// LoggerConfig Options
type LoggerConfig struct {
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
func (logger LoggerConfig) GetCode() string {
	return logger.Code
}

// GetLevel returns the level
func (logger LoggerConfig) GetFormatter() string {
	return logger.Formatter
}

// GetLevel returns the level
func (logger LoggerConfig) GetLevel() string {
	return logger.Level
}

// GetEnableCaller returns if we want to enable caller in output
func (logger LoggerConfig) GetEnableCaller() bool {
	return logger.EnableCaller
}

// GetFileEnabled if we want to show file
func (logger LoggerConfig) GetFileEnabled() bool {
	return logger.FileEnabled
}

// GetFilename if we want to show filename
func (logger LoggerConfig) GetFilename() string {
	return logger.Filename
}

// GetMaxSize of log
func (logger LoggerConfig) GetMaxSize() int {
	return logger.MaxSize
}

// GetMaxBackups max backups
func (logger LoggerConfig) GetMaxBackups() int {
	return logger.MaxBackups
}

// GetMaxAge of log
func (logger LoggerConfig) GetMaxAge() int {
	return logger.MaxAge
}

// GetCompress should we compress
func (logger LoggerConfig) GetCompress() bool {
	return logger.Compress
}

// GetLocalTime if we store in local time
func (logger LoggerConfig) GetLocalTime() bool {
	return logger.LocalTime
}
