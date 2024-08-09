package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/neb-ep/monorepo/shared/logger"
)

func main() {
	logger.InitLogger(logger.Config{
		Name:        "auth-service-migration",
		Version:     "v1.0.1",
		Environment: "devel",
	})

	down := flag.Bool("down", false, "rollback all migrations")
	flag.Parse()

	var (
		err error
		m   *migrate.Migrate
	)

	m, err = migrate.New(
		"file://migrations",
		"postgres://devel:devel@localhost:5432/devel?sslmode=disable",
	)
	if err != nil {
		slog.Error("failed to create migrator", "error", err)
		os.Exit(1)
	}

	if *down {
		err = m.Down()
		if err != nil {
			slog.Error("failed to rollback migrations", "error", err)
			os.Exit(1)
		}
	} else {
		err = m.Up()
		if err != nil {
			slog.Error("failed to apply migrations", "error", err)
			os.Exit(1)
		}
	}
}
