package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(ctx context.Context, config *Config, logger *slog.Logger) (pool *pgxpool.Pool, err error) {
	var pgxConfig *pgxpool.Config

	logger.Debug("pgxpool parse config from connection URL",
		"connectionURL", config.ConnectionURL)
	pgxConfig, err = pgxpool.ParseConfig(config.ConnectionURL)
	if err != nil {
		logger.Debug("failed to parse connection URL",
			"connectionURL", config.ConnectionURL,
			"error", err)
		return nil, fmt.Errorf("failed to parse connection URL: %w", err)
	}

	logger.Debug("configure before aquire connection check")
	pgxConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) (valid bool) {
		valid = c.Ping(ctx) == nil
		logger.Debug("check connection is valid before aquire from the pool",
			"valid", valid)
		return valid
	}

	logger.Debug("creating new pgxpool from config")
	pool, err = pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		logger.Debug("failed to create connection pool",
			"config", pgxConfig,
			"error", err)
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return pool, nil
}
