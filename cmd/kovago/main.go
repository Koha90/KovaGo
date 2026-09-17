// Command kovago starts the KovaGo commerce engine.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Koha90/KovaGo/internal/config"
	"github.com/Koha90/KovaGo/internal/logger"
	"github.com/Koha90/KovaGo/internal/migrator"
	"github.com/Koha90/KovaGo/internal/postgres"
)

const defaultConfigPath = "config/local.toml"

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(
			os.Stderr,
			"kovago: %v\n",
			err,
		)

		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	configPath := os.Getenv("KOVAGO_CONFIG")
	if configPath == "" {
		configPath = defaultConfigPath
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	log, closeLogger, err := logger.New(cfg.Log)
	if err != nil {
		return fmt.Errorf(
			"create logger: %w",
			err,
		)
	}

	defer func() {
		if err := closeLogger(); err != nil {
			_, _ = fmt.Fprintf(
				os.Stderr,
				"kovago: close logger: %v\n",
				err,
			)
		}
	}()

	log = log.With(
		"service",
		cfg.App.Name,
		"environment",
		cfg.App.Environment,
	)

	log.Info(
		"starting KovaGo",
	)

	db, err := postgres.Open(
		ctx,
		cfg.Database,
	)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Info("connected to PostgreSQL")

	if err := migrator.Up(
		ctx,
		db,
		log,
	); err != nil {
		return err
	}

	log.Info("KovaGo initialized")

	return nil
}
