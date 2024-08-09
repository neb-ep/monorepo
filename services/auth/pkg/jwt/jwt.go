package jwt

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/neb-ep/monorepo/shared/telemetry"
)

const tokenRefreshLength = 16

type jwtHelper struct {
	secret string
	TTL    time.Duration
}

func NewJWTHelper(secret string, TTL time.Duration) *jwtHelper {
	return &jwtHelper{secret: secret, TTL: TTL}
}

func (j *jwtHelper) Generate(ctx context.Context, username string) (string, error) {
	_, span := telemetry.Tracer.Start(ctx, "jwt/generate-accessToken")
	defer span.End()

	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "auth",
			"sub": username,
			"exp": time.Now().Add(j.TTL).Unix(),
		})
	s, err := t.SignedString([]byte(j.secret))
	if err != nil {
		return "", fmt.Errorf("Generate.SignedString(): %w", err)
	}
	return s, nil
}

func (j *jwtHelper) GenerateRefreshToken(ctx context.Context) (string, error) {
	_, span := telemetry.Tracer.Start(ctx, "jwt/generate-refreshToken")
	defer span.End()

	b := make([]byte, tokenRefreshLength)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("GenerateRefreshToken.rand.Read(): %w", err)
	}

	return hex.EncodeToString(b), nil
}
