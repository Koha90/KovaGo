// Package migrator manages KovaGo database migrations.
package migrator

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/pressly/goose/v3"
)

// migrations contains all SQL migrations shipped with KovaGo.
//
//go:embed migrations
var migrations embed.FS

// Up applies all pending database migrations.
func Up(
	ctx context.Context,
	db *sql.DB,
	log *slog.Logger,
) error {
	migrationFS, err := fs.Sub(
		migrations,
		"migrations",
	)
	if err != nil {
		return fmt.Errorf(
			"open embedded migrations: %w",
			err,
		)
	}

	provider, err := goose.NewProvider(
		goose.DialectPostgres,
		db,
		migrationFS,
		goose.WithSlog(
			log.With("component", "migrator"),
		),
		goose.WithTableName("kovago_schema_migrations"),
	)
	if err != nil {
		return fmt.Errorf(
			"create migration provider: %w",
			err,
		)
	}

	results, err := provider.Up(ctx)
	if err != nil {
		return fmt.Errorf(
			"apply database migrations: %w",
			err,
		)
	}

	if len(results) == 0 {
		log.Debug("database schema is up to date")

		return nil
	}

	log.Info(
		"database migrations applied",
		"count",
		len(results),
	)

	return nil
}
