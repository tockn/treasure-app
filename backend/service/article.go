package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/voyagegroup/treasure-app/model"
	"github.com/voyagegroup/treasure-app/repository"
)

type Article interface {
	GetAll(ctx context.Context) ([]*model.Article, error)
	Get(ctx context.Context, id int64) (*model.Article, error)
	Update(ctx context.Context, id int64, newArticle *model.Article) error
	Destroy(ctx context.Context, id int64) error
	Create(ctx context.Context, newArticle *model.Article) (int64, error)
}

type article struct {
	articleRepo repository.Article
}

func NewArticle(ar repository.Article) Article {
	return &article{
		articleRepo: ar,
	}
}

func (a *article) GetAll(ctx context.Context) ([]*model.Article, error) {
	return a.articleRepo.AllArticle(ctx)
}

func (a *article) Get(ctx context.Context, id int64) (*model.Article, error) {
	return a.articleRepo.FindArticle(ctx, id)
}

func (a *article) Update(ctx context.Context, id int64, newArticle *model.Article) error {
	_, err := a.articleRepo.FindArticle(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed find article")
	}

	if err := a.articleRepo.TxHandler(ctx, func(ctx context.Context) error {
		_, err := a.articleRepo.UpdateArticle(ctx, id, newArticle)
		return err
	}); err != nil {
		return errors.Wrap(err, "failed article update transaction")
	}
	return nil
}

func (a *article) Destroy(ctx context.Context, id int64) error {
	_, err := a.articleRepo.FindArticle(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed find article")
	}

	if err := a.articleRepo.TxHandler(ctx, func(ctx context.Context) error {
		_, err := a.articleRepo.DestroyArticle(ctx, id)
		return err
	}); err != nil {
		return errors.Wrap(err, "failed article delete transaction")
	}
	return nil
}

func (a *article) Create(ctx context.Context, newArticle *model.Article) (int64, error) {
	var createdId int64
	if err := a.articleRepo.TxHandler(ctx, func(ctx context.Context) error {
		id, err := a.articleRepo.CreateArticle(ctx, newArticle)
		if err != nil {
			return err
		}
		createdId = id
		return err
	}); err != nil {
		return 0, errors.Wrap(err, "failed article insert transaction")
	}
	return createdId, nil
}
