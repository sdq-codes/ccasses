package config

import (
	"context"
	"gopkg.in/yaml.v3"
	"os"

	"github.com/sdq-codes/ccasses/core/config"
	"github.com/sdq-codes/ccasses/core/errors"
	"github.com/sdq-codes/ccasses/core/logging"
	"go.uber.org/zap"
)

const (
	// ErrValidation is returned when the configuration is invalid.
	ErrValidation = errors.Error("invalid configuration")
	// ErrEnvVars is returned when the environment variables are invalid.
	ErrEnvVars = errors.Error("failed parsing env vars")
	// ErrRead is returned when the configuration file cannot be read.
	ErrRead = errors.Error("failed to read file")
	// ErrUnmarshal is returned when the configuration file cannot be unmarshalled.
	ErrUnmarshal = errors.Error("failed to unmarshal file")
)

var (
	baseConfigPath = "config/config.yaml"
	envConfigPath  = "config/config-%s.yaml"
)

// Config represents the configuration of our application.
type Config struct {
	config.AppConfig `yaml:",inline"`
}

// Load loads the configuration from the config/config.yaml file.
func Load(ctx context.Context) (*Config, error) {
	cfg := &Config{}
	cfg.HTTP.Port = "8081"

	cfg.PSQL.DSN = "postgresql://guest:guest@localhost:5432/luxridr?sslmode=disable"

	if err := env.Parse(cfg); err != nil {
		return nil, ErrEnvVars.Wrap(err)
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, ErrValidation.Wrap(err)
	}

	return cfg, nil
}

func loadYaml(ctx context.Context, filename string, cfg any) error {
	logging.From(ctx).Info("Loading configur", zap.String("path", filename))

	data, err := os.ReadFile(filename)
	if err != nil {
		return ErrRead.Wrap(err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return ErrUnmarshal.Wrap(err)
	}

	return nil
}
