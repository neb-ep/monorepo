package hasher

import (
	"context"

	"github.com/neb-ep/monorepo/shared/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinCost     int = 4
	MaxCost     int = 31
	DefaultCost int = 10
)

type hasher struct {
	cost int
}

func NewHasher(cost int) *hasher {
	return &hasher{cost: cost}
}

func (h *hasher) Generate(ctx context.Context, password string) (string, error) {
	_, span := telemetry.Tracer.Start(ctx, "hasher/generate")
	defer span.End()

	span.SetAttributes(attribute.Int("hasher.cost", h.cost))

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	return string(bytes), err
}

func (h *hasher) Compare(ctx context.Context, passwod string, withHash string) error {
	_, span := telemetry.Tracer.Start(ctx, "hasher/compare")
	defer span.End()

	return bcrypt.CompareHashAndPassword([]byte(withHash), []byte(passwod))
}
