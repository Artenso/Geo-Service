package controller

import (
	"context"

	"github.com/Artenso/user-service/internal/model"
	"github.com/Artenso/user-service/internal/service"
	desc "github.com/Artenso/user-service/pkg/user_service"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Controller struct {
	desc.UnimplementedUserServiceServer

	service service.IService
}

func NewController(service service.IService) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	user := &model.User{
		Name: req.GetUsername(),
		Pass: []byte(req.GetPassword()),
	}

	id, err := c.service.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &desc.CreateUserResponse{Id: id}, nil
}

func (c *Controller) GetUserByID(ctx context.Context, req *desc.GetUserByIDRequest) (*desc.GetUserByIDResponse, error) {
	dbUser, err := c.service.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	user := &desc.User{
		Id:       dbUser.ID,
		Username: dbUser.Name,
	}

	return &desc.GetUserByIDResponse{User: user}, nil
}

func (c *Controller) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	if err := c.service.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Controller) ListUsers(ctx context.Context, req *desc.ListUsersRequest) (*desc.ListUsersResponse, error) {
	users, err := c.service.List(ctx, req.GetLimit(), req.GetOffset())
	if err != nil {
		return nil, err
	}

	resp := &desc.ListUsersResponse{
		Users: make([]*desc.User, 0, len(users)),
	}

	for _, user := range users {
		descUser := &desc.User{
			Id:       user.ID,
			Username: user.Name,
		}

		resp.Users = append(resp.Users, descUser)
	}

	return resp, nil
}

func (c *Controller) GetUserByNameAndPass(ctx context.Context, req *desc.GetUserByNameAndPassRequest) (*desc.GetUserByNameAndPassResponse, error) {
	user := &model.User{
		Name: req.GetUsername(),
		Pass: []byte(req.GetPassword()),
	}

	dbUser, err := c.service.GetByNameAndPass(ctx, user)
	if err != nil {
		return nil, err
	}

	descUser := &desc.User{
		Id:       dbUser.ID,
		Username: dbUser.Name,
	}

	return &desc.GetUserByNameAndPassResponse{User: descUser}, nil
}
