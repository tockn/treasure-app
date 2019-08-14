package repository

import (
	"context"

	"github.com/voyagegroup/treasure-app/model"
)

type User interface {
	Tx
	Get(ctx context.Context, id string) (*model.User, error)
	Sync(ctx context.Context, fu *model.FirebaseUser) (int64, error)
}
