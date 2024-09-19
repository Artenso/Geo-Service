package controller

import (
	"context"

	userClient "github.com/Artenso/auth-service/internal/clients/user/client" // replase with import from user-service in normal case
	"github.com/Artenso/auth-service/internal/service"
	desc "github.com/Artenso/auth-service/pkg/auth_service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Controller struct {
	desc.UnimplementedAuthServiceServer

	service           service.Service
	userServiceClient userClient.Client
}

func NewController(service service.Service, client userClient.Client) *Controller {
	return &Controller{
		service:           service,
		userServiceClient: client,
	}
}

func (c *Controller) Register(ctx context.Context, req *desc.RegisterRequest) (*desc.RegisterResponse, error) {
	id, err := c.userServiceClient.CreateUser(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &desc.RegisterResponse{Id: id}, nil
}

func (c *Controller) LogIn(ctx context.Context, req *desc.LogInRequest) (*desc.LogInResponse, error) {
	user, err := c.userServiceClient.GetUserByNameAndPass(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	token, err := c.service.GenerateToken(ctx, user.Name)
	if err != nil {
		return nil, err
	}

	return &desc.LogInResponse{Token: token}, err
}

func (c *Controller) Verify(ctx context.Context, req *desc.VerifyRequest) (*emptypb.Empty, error) {
	if err := c.service.Verify(ctx, req.GetToken()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
