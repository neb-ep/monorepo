package app

import (
	"context"

	authv1 "github.com/neb-ep/shared/contracts/protos/auth/v1"

	"github.com/neb-ep/monorepo/services/auth/internal/entities"
)

type Service interface {
	SignUp(context.Context, *entities.User) (*entities.User, error)
}

type Api struct {
	service Service
	authv1.UnimplementedAuthServiceServer
}

func NewApi(service Service) *Api {
	return &Api{service: service}
}

// SignIn implements auth.AuthServiceServer.
func (a *Api) SignIn(_ context.Context, req *authv1.SignInRequest) (*authv1.SignInResponse, error) {
	panic("unimplemented")
}

// SignUp implements auth.AuthServiceServer.
func (a *Api) SignUp(ctx context.Context, req *authv1.SignUpRequest) (*authv1.SignUpResponse, error) {
	u, _ := a.service.SignUp(ctx, &entities.User{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
		Email:     req.Email,
	})

	return &authv1.SignUpResponse{
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}, nil
}
