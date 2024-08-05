package service

import (
	"context"

	"github.com/neb-ep/monorepo/services/auth/internal/entities"
	"github.com/neb-ep/monorepo/shared/telemetry"
)

type Storage interface {
	CreateUser(context.Context, *entities.User) (*entities.User, error)
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) SignUp(ctx context.Context, user *entities.User) (*entities.User, error) {
	_, span := telemetry.Tracer.Start(ctx, "service/signup")
	defer span.End()

	return s.storage.CreateUser(ctx, user)
}
