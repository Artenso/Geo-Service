package storage

import (
	"github.com/Artenso/Geo-Service/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type IStorage interface {
	Create(user *model.User) (int, error)
	IsRegistered(user *model.User) (bool, error)
}

type storage struct {
	data   map[int]*model.User
	prevID int
}

func NewStorage() IStorage {
	return &storage{
		data:   make(map[int]*model.User),
		prevID: 0,
	}
}

func (s *storage) Create(user *model.User) (int, error) {
	id := s.prevID + 1
	s.data[id] = user
	s.prevID = id

	return id, nil
}

func (s *storage) IsRegistered(user *model.User) (bool, error) {
	for _, dbuser := range s.data {
		if dbuser.Name == user.Name {
			if bcrypt.CompareHashAndPassword(dbuser.Pass, user.Pass) == nil {
				return true, nil
			}
		}
	}

	return false, nil
}
