package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/neb-ep/monorepo/services/auth/internal/entities"
	"go.opentelemetry.io/otel/trace"
)

type Storage interface {
	CreateUser(context.Context, *entities.User) (*entities.User, error)
	GetUserCredential(context.Context, string) (*entities.UserCreds, error)
	CreateTokenSession(context.Context, string, int) error
}

type Hasher interface {
	Generate(ctx context.Context, password string) (string, error)
	Compare(ctx context.Context, password string, hash string) error
}

type JWTHelper interface {
	Generate(ctx context.Context, username string) (string, error)
	GenerateRefreshToken(ctx context.Context) (string, error)
}

type Service struct {
	storage   Storage
	hasher    Hasher
	jwtHelper JWTHelper
}

func NewService(storage Storage, hasher Hasher, jwtHelper JWTHelper) *Service {
	return &Service{storage: storage, hasher: hasher, jwtHelper: jwtHelper}
}

func (s *Service) SignUp(ctx context.Context, user *entities.User) (u *entities.User, err error) {
	span := trace.SpanFromContext(ctx)

	slog.Debug("generate password hash")
	user.Password, err = s.hasher.Generate(ctx, user.Password)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to generate hash from password: %w", err)
	}

	u, err = s.storage.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) SignIn(ctx context.Context, email, password string) (u *entities.AccessToken, err error) {
	uc, err := s.storage.GetUserCredential(ctx, email)
	if err != nil {
		return nil, err
	}

	if err = s.hasher.Compare(ctx, password, uc.PasswordHash); err != nil {
		return nil, fmt.Errorf("signin. failed hasher.compare: %s", err)
	}

	// Geneate tokens
	accessToken, err := s.jwtHelper.Generate(ctx, uc.Username)
	if err != nil {
		return nil, fmt.Errorf("signin. failed jwtHelper.Generate: %w", err)
	}
	refreshToken, err := s.jwtHelper.GenerateRefreshToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("signin. failed jwtHelper.GenerateRefreshToken: %w", err)
	}

	if err = s.storage.CreateTokenSession(ctx, refreshToken, uc.UserId); err != nil {
		return nil, fmt.Errorf("signin. failed storage.CreateTokenSession: %w", err)
	}

	return &entities.AccessToken{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}
