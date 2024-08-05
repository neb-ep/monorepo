package storage

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neb-ep/monorepo/services/auth/internal/entities"
	"github.com/neb-ep/monorepo/shared/telemetry"
)

var logger = slog.Default().WithGroup("storage")

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

func (s *Storage) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	_, span := telemetry.Tracer.Start(ctx, "storage/create_user")
	defer span.End()
	return user, nil
}
