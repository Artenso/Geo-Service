// delete this file and pkg/user_service in normal case
package client

import (
	"context"

	"github.com/Artenso/auth-service/internal/model"
	desc "github.com/Artenso/auth-service/pkg/user_service"

	"google.golang.org/grpc"
)

type Client interface {
	CreateUser(ctx context.Context, username, password string) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, limit, offset int64) ([]*model.User, error)
	GetUserByNameAndPass(ctx context.Context, username, password string) (*model.User, error)
}

type client struct {
	client desc.UserServiceClient
}

func NewGRPCclient(conn *grpc.ClientConn) Client {
	return &client{
		client: desc.NewUserServiceClient(conn),
	}
}

func (c *client) CreateUser(ctx context.Context, username, password string) (int64, error) {
	req := &desc.CreateUserRequest{
		Username: username,
		Password: password,
	}

	resp, err := c.client.CreateUser(ctx, req)
	if err != nil {
		return 0, err
	}

	return resp.GetId(), nil
}

func (c *client) GetUserByNameAndPass(ctx context.Context, username, password string) (*model.User, error) {
	req := &desc.GetUserByNameAndPassRequest{
		Username: username,
		Password: password,
	}

	resp, err := c.client.GetUserByNameAndPass(ctx, req)
	if err != nil {
		return nil, err
	}

	descUser := resp.GetUser()

	return &model.User{
			ID:   descUser.Id,
			Name: descUser.Username,
			Pass: []byte(descUser.Password),
		},
		nil
}

func (c *client) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	req := &desc.GetUserByIDRequest{
		Id: id,
	}

	resp, err := c.client.GetUserByID(ctx, req)
	if err != nil {
		return nil, err
	}

	descUser := resp.GetUser()

	return &model.User{
			ID:   descUser.Id,
			Name: descUser.Username,
			Pass: []byte(descUser.Password),
		},
		nil
}

func (c *client) DeleteUser(ctx context.Context, id int64) error {
	req := &desc.DeleteUserRequest{
		Id: id,
	}

	_, err := c.client.DeleteUser(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) ListUsers(ctx context.Context, limit, offset int64) ([]*model.User, error) {
	req := &desc.ListUsersRequest{
		Limit:  limit,
		Offset: offset,
	}

	resp, err := c.client.ListUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	descUsers := resp.GetUsers()

	users := make([]*model.User, 0, len(descUsers))

	for _, descUser := range descUsers {
		user := &model.User{
			ID:   descUser.Id,
			Name: descUser.Username,
		}

		users = append(users, user)
	}
	return users, nil
}
