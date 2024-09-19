package service

import (
	"context"

	"github.com/go-chi/jwtauth"
)

type Service interface {
	GenerateToken(ctx context.Context, name string) (string, error)
	Verify(ctx context.Context, tokenString string) error
}

type service struct {
	tokenAuth *jwtauth.JWTAuth
}

func NewService(tokenAuth *jwtauth.JWTAuth) Service {
	return &service{
		tokenAuth: tokenAuth,
	}
}

func (s *service) GenerateToken(_ context.Context, name string) (string, error) {
	_, tokenString, err := s.tokenAuth.Encode(map[string]interface{}{"username": name})
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *service) Verify(_ context.Context, tokenString string) error {
	_, err := jwtauth.VerifyToken(s.tokenAuth, tokenString)
	if err != nil {
		return err
	}

	return nil
}
