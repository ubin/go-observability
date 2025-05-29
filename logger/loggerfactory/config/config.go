package config

const (
	LOGRUS = "logrus"
	SLOG   = "slog"
)

// LogEnv type
type LogEnv int

const (
	// LogEnvDev LogEnv
	LogEnvDev LogEnv = iota
	// LogEnvProd LogEnv
	LogEnvProd
)

func (e LogEnv) String() string {
	return [...]string{"dev", "prod"}[e]
}
