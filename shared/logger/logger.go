package logger

import (
	"log/slog"
	"os"
)

const (
	attrServiceName        = "service_name"
	attrServiceVersion     = "service_version"
	attrServiceEnvironment = "service_environment"
)

type Config struct {
	Name        string     `yaml:"name"`
	Version     string     `yaml:"version"`
	Environment string     `yaml:"environment"`
	Level       slog.Level `yaml:"level"`
}

func InitLogger(config Config) {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: true,
	})
	slog.SetDefault(slog.New(handler).With(
		attrServiceName, config.Name,
		attrServiceVersion, config.Version,
		attrServiceEnvironment, config.Environment),
	)
}
