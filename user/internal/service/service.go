package service

import (
	"context"

	"github.com/Artenso/user-service/internal/model"
	"github.com/Artenso/user-service/internal/storage"

	"golang.org/x/crypto/bcrypt"
)

type IService interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int64) ([]*model.User, error)
	GetByNameAndPass(ctx context.Context, user *model.User) (*model.User, error)
}

type service struct {
	storage storage.IStorage
}

func NewService(storage storage.IStorage) IService {
	return &service{
		storage: storage,
	}
}

func (s *service) Create(ctx context.Context, user *model.User) (int64, error) {
	password, err := bcrypt.GenerateFromPassword(user.Pass, bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	dbUser := &model.User{
		Name: user.Name,
		Pass: password,
	}

	id, err := s.storage.Create(ctx, dbUser)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *service) GetByNameAndPass(ctx context.Context, user *model.User) (*model.User, error) {
	users, err := s.storage.GetByName(ctx, user)
	if err != nil {
		return nil, err
	}

	for _, dbUser := range users {
		if err := bcrypt.CompareHashAndPassword(dbUser.Pass, user.Pass); err != nil {
			continue
		}

		return dbUser, nil
	}

	return nil, model.ErrorUserNotFound
}

func (s *service) GetByID(ctx context.Context, id int64) (*model.User, error) {
	return s.storage.GetByID(ctx, id)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.storage.Delete(ctx, id)
}

func (s *service) List(ctx context.Context, limit, offset int64) ([]*model.User, error) {
	return s.storage.List(ctx, limit, offset)
}
