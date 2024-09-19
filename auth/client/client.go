package client

import (
	"context"

	desc "github.com/Artenso/auth-service/pkg/auth_service"
	"google.golang.org/grpc"
)

type Client interface {
	Register(ctx context.Context, username, password string) (int64, error)
	LogIn(ctx context.Context, username, password string) (string, error)
	Verify(ctx context.Context, token string) error
}

type client struct {
	client desc.AuthServiceClient
}

func NewGRPCclient(conn *grpc.ClientConn) Client {
	return &client{
		client: desc.NewAuthServiceClient(conn),
	}
}

func (c *client) Register(ctx context.Context, username, password string) (int64, error) {
	req := &desc.RegisterRequest{
		Username: username,
		Password: password,
	}

	resp, err := c.client.Register(ctx, req)
	if err != nil {
		return 0, err
	}

	return resp.GetId(), nil
}

func (c *client) LogIn(ctx context.Context, username, password string) (string, error) {
	req := &desc.LogInRequest{
		Username: username,
		Password: password,
	}

	resp, err := c.client.LogIn(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.GetToken(), nil
}

func (c *client) Verify(ctx context.Context, token string) error {
	req := &desc.VerifyRequest{
		Token: token,
	}

	_, err := c.client.Verify(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
