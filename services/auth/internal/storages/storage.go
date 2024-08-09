package storages

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neb-ep/monorepo/services/auth/internal/entities"
	"github.com/neb-ep/monorepo/shared/database"
	"github.com/neb-ep/monorepo/shared/telemetry"
)

var (
	ErrUsernameAlreadyUsed = errors.New("username already used")
	ErrUsernameMinLength   = errors.New("username must be between 6 and 255 rune")
	ErrEmailAlreadyUsed    = errors.New("email already used")
	ErrEmailMinLength      = errors.New("email must be between 6 and 255 rune")
)

var (
	ConstraintPrimaryKey             = "users_pkey"
	ConstraintUniqueUsername         = "users_username_key"
	ConstraintUniqueEmail            = "users_email_key"
	ConstraintCheckMinLengthUsername = "users_username_min_length_check"
	ConstraintCheckMinLengthEmail    = "users_email_min_length_check"
)

type Storage struct {
	pool        *pgxpool.Pool
	queries     Queries
	constraints database.ErrorMapper
	logger      slog.Logger
}

func NewStorage(db *pgxpool.Pool) *Storage {
	return &Storage{queries: *New(db), pool: db, constraints: *database.NewErrorDescriber(database.ConstraintMapper{
		ConstraintUniqueUsername:         ErrUsernameAlreadyUsed,
		ConstraintUniqueEmail:            ErrEmailAlreadyUsed,
		ConstraintCheckMinLengthEmail:    ErrEmailMinLength,
		ConstraintCheckMinLengthUsername: ErrUsernameMinLength,
	}),
		logger: *slog.Default().WithGroup("storage"),
	}
}

func (s *Storage) CreateUser(ctx context.Context, in *entities.User) (out *entities.User, err error) {
	_, span := telemetry.Tracer.Start(ctx, "storage/CreateUser")
	defer span.End()

	ctx, cancelFunc := context.WithTimeout(ctx, 120*time.Second)
	defer cancelFunc()

	row, err := s.queries.CreateUser(ctx, CreateUserParams{
		Username:  in.Username,
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Email:     in.Email,
		Passwhash: in.Password,
		IsActive:  true,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})

	if err != nil {
		s.logger.Error("failed to exec query 'CreateUser'", "error", err)
		return nil, s.constraints.Describe(err)
	}

	return &entities.User{
		Username:  row.Username,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Email:     row.Email,
	}, nil
}

func (s *Storage) GetUserCredential(ctx context.Context, email string) (out *entities.UserCreds, err error) {
	_, span := telemetry.Tracer.Start(ctx, "storage/GetUserCredential")
	defer span.End()

	ctx, cancelFunc := context.WithTimeout(ctx, 120*time.Second)
	defer cancelFunc()

	row, err := s.queries.GetUserCredentialByEmail(ctx, email)
	if err != nil {
		return nil, s.constraints.Describe(err)
	}

	return &entities.UserCreds{
		UserId:       int(row.UserID),
		Email:        row.Email,
		PasswordHash: row.Passwhash,
	}, nil
}

func (s *Storage) CreateTokenSession(ctx context.Context, token string, userId int) error {
	_, span := telemetry.Tracer.Start(ctx, "storage/CreateTokenSession")
	defer span.End()

	tx, _ := s.pool.Begin(ctx)
	defer tx.Rollback(ctx)

	q := s.queries.WithTx(tx)

	if err := q.DeactivateTokenSession(ctx, DeactivateTokenSessionParams{
		UsedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UserID: int32(userId),
	}); err != nil {
		return fmt.Errorf("failed to queries.DeactivateTokenSession(userId=%d): %w", userId, err)
	}

	_, err := q.InsertTokenSession(ctx, InsertTokenSessionParams{
		UserID:    int32(userId),
		Token:     token,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		IsActive:  true,
	})
	if err != nil {
		return fmt.Errorf("failed to queries.InsertTokenSession: %w", err)
	}

	return tx.Commit(ctx)
}
