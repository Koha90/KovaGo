// Package postgres provides PostrgreSQL infrastructure.
package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Koha90/KovaGo/internal/config"
)

// Open creates and verifies a PostrgreSQL connection pool.
func Open(
	ctx context.Context,
	cfg config.Database,
) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("open PostrgreSQL: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnetctionMaxIdleTime)

	pingCtx, cancel := context.WithTimeout(
		ctx,
		cfg.PingTimeout,
	)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()

		return nil, fmt.Errorf(

			"ping PostrgreSQL: %w",
			err,
		)
	}

	return db, nil
}
