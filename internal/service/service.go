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
	GetAddrByPart(ctx context.Context, input string) ([]*model.Address, error)
	GetAddrByCoord(ctx context.Context, lat, lng string) ([]*model.Address, error)
}

type service struct {
	storage storage.IStorage
	geoServ GeoProvider
}

func NewService(storage storage.IStorage, geoServ GeoProvider) IService {
	return &service{
		storage: storage,
		geoServ: geoServ,
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

func (s *service) GetAddrByPart(ctx context.Context, input string) ([]*model.Address, error) {
	return s.geoServ.AddressSearch(input)

}

func (s *service) GetAddrByCoord(ctx context.Context, lat, lng string) ([]*model.Address, error) {
	return s.geoServ.GeoCode(lat, lng)
}
