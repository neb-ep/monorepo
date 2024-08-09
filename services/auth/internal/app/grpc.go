package app

import (
	"context"

	authv1 "github.com/neb-ep/shared/contracts/protos/auth/v1"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/neb-ep/monorepo/services/auth/internal/entities"
)

type Service interface {
	SignUp(ctx context.Context, in *entities.User) (out *entities.User, err error)
	SignIn(ctx context.Context, email, password string) (out *entities.AccessToken, err error)
}

type Api struct {
	service Service
	authv1.UnimplementedAuthServiceServer
}

func NewApi(service Service) *Api {
	return &Api{service: service}
}

// SignIn implements auth.AuthServiceServer.
func (a *Api) SignIn(ctx context.Context, req *authv1.SignInRequest) (*authv1.SignInResponse, error) {
	span := trace.SpanFromContext(ctx)

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	access, err := a.service.SignIn(ctx, req.Email, req.Password)
	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &authv1.SignInResponse{
		AccessToken:  access.Access,
		RefreshToken: access.Refresh,
	}, nil
}

// SignUp implements auth.AuthServiceServer.
func (a *Api) SignUp(ctx context.Context, req *authv1.SignUpRequest) (*authv1.SignUpResponse, error) {
	span := trace.SpanFromContext(ctx)

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Service processing
	u, err := a.service.SignUp(ctx, &entities.User{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
		Email:     req.Email,
	})

	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &authv1.SignUpResponse{
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}, nil
}
