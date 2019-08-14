package service

import (
	"context"

	"github.com/voyagegroup/treasure-app/model"
	"github.com/voyagegroup/treasure-app/repository"
)

type User interface {
	Get(ctx context.Context, id string) (*model.User, error)
	Sync(ctx context.Context, fu *model.FirebaseUser) (int64, error)
}

type user struct {
	userRepo repository.User
}

func NewUser(ur repository.User) User {
	return &user{
		userRepo: ur,
	}
}

func (s *user) Get(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.Get(ctx, id)
}

func (s *user) Sync(ctx context.Context, fu *model.FirebaseUser) (int64, error) {
	return s.userRepo.Sync(ctx, fu)
}
