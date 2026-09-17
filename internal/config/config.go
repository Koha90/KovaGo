// Package config loads and validates KovaGo configuration.
package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config contains the complete KovaGo configuration.
type Config struct {
	App      App      `toml:"app"`
	Log      Log      `toml:"log"`
	Database Database `toml:"database"`
}

// App contains general application settings.
type App struct {
	Name        string `toml:"name"        env:"KOVAGO_APP_NAME" env-default:"KovaGo"`
	Environment string `toml:"environment" env:"KOVAGO_ENV"      env-default:"local"`
}

// Log contains logging settings.
type Log struct {
	Level     string  `toml:"level"      env:"KOVAGO_LOG_LEVEL"      env-default:"info"`
	Format    string  `toml:"format"     env:"KOVAGO_LOG_FORMAT"     env-default:"text"`
	AddSource bool    `toml:"add_source" env:"KOVAGO_LOG_ADD_SOURCE" env-default:"false"`
	File      LogFile `toml:"file"`
}

// LogFile contains file logging settings.
type LogFile struct {
	Enabled bool   `toml:"enabled" env:"KOVAGO_LOG_FILE_ENABLED" env-default:"false"`
	Path    string `toml:"path"    env:"KOVAGO_LOG_FILE_PATH"    env-default:"logs/kovago.log"`
	Format  string `toml:"format"  env:"KOVAGO_LOG_FILE_FORMAT"  env-default:"json"`
}

// Database contains PostgreSQL connection settings.
type Database struct {
	URL string `toml:"url" env:"url"`

	MaxOpenConnections int `toml:"max_open_connections" env:"KOVAGO_DATABASE_MAX_OPEN_CONNECTIONS" env-default:"20"`
	MaxIdleConnections int `toml:"max_idle_connections" env:"KOVAGO_DATABASE_MAX_IDLE_CONNECTIONS" env-default:"10"`

	ConnectionMaxLifetime  time.Duration `toml:"connection_max_lifetime"   env:"KOVAGO_DATABASE_CONNECTION_MAX_LIFETIME"   env-default:"30m"`
	ConnetctionMaxIdleTime time.Duration `toml:"connetction_max_idle_time" env:"KOVAGO_DATABASE_CONNETCTION_MAX_IDLE_TIME" env-default:"5m"`

	PingTimeout time.Duration `toml:"ping_timeout" env:"KOVAGO_DATABASE_ping_timeout"`
}

// Load reads configuration from path and applies environment overrides.
func Load(path string) (Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return Config{}, fmt.Errorf("read configuration %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return Config{}, fmt.Errorf("validate configuration: %w", err)
	}

	return cfg, nil
}

func (cfg Config) validate() error {
	if strings.TrimSpace(cfg.App.Name) == "" {
		return errors.New("app name is required")
	}

	if strings.TrimSpace(cfg.Database.URL) == "" {
		return errors.New("database URL is required")
	}

	if cfg.Database.MaxOpenConnections <= 0 {
		return errors.New("database max open connections must be greater than zero")
	}

	if cfg.Database.MaxIdleConnections < 0 {
		return errors.New("database max idle connections must not be negative")
	}

	if cfg.Database.MaxIdleConnections > cfg.Database.MaxOpenConnections {
		return errors.New(
			"database max idle connections must not exceed max open connections",
		)
	}

	if cfg.Database.PingTimeout <= 0 {
		return errors.New("database ping timeout must be greater than zero")
	}

	return nil
}
