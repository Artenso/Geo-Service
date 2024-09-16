package service

import (
	"context"

	"github.com/Artenso/Geo-Service/internal/model"
	storage "github.com/Artenso/Geo-Service/internal/storage/pg"
	"github.com/Artenso/Geo-Service/internal/token"
	"golang.org/x/crypto/bcrypt"
)

type IService interface {
	RegistrateUser(ctx context.Context, user *model.User) error
	AuthenticateUser(ctx context.Context, user *model.User) (string, error)
}

type service struct {
	storage storage.IStorage
}

func NewService(storage storage.IStorage) IService {
	return &service{
		storage: storage,
	}
}

func (s *service) RegistrateUser(ctx context.Context, user *model.User) error {
	password, err := bcrypt.GenerateFromPassword(user.Pass, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	dbUser := &model.User{
		Name: user.Name,
		Pass: password,
	}

	err = s.storage.Create(ctx, dbUser)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) AuthenticateUser(ctx context.Context, user *model.User) (string, error) {
	users, err := s.storage.GetByName(ctx, user)
	if err != nil {
		return "", err
	}

	for _, dbUser := range users {
		if err := bcrypt.CompareHashAndPassword(dbUser.Pass, user.Pass); err != nil {
			continue
		}

		return token.Generate(user.Name)
	}

	return "", model.ErrorUserNotFound
}
