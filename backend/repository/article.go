package repository

import (
	"context"

	"github.com/voyagegroup/treasure-app/model"
)

type Article interface {
	Tx
	AllArticle(ctx context.Context) ([]*model.Article, error)
	FindArticle(ctx context.Context, id int64) (*model.Article, error)
	CreateArticle(ctx context.Context, a *model.Article) (id int64, err error)
	UpdateArticle(ctx context.Context, id int64, a *model.Article) (int64, error)
	DestroyArticle(ctx context.Context, id int64) (int64, error)
}
