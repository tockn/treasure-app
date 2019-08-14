package rdb

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/voyagegroup/treasure-app/model"
	"github.com/voyagegroup/treasure-app/repository"
)

type article struct {
	txHandler
	db *sqlx.DB
}

func NewArticle(db *sqlx.DB) repository.Article {
	return &article{
		txHandler: txHandler{db: db},
		db:        db,
	}
}

func (r *article) AllArticle(ctx context.Context) ([]*model.Article, error) {
	a := make([]*model.Article, 0)
	if err := r.db.Select(&a, `SELECT id, title, body FROM article`); err != nil {
		return nil, err
	}
	return a, nil
}

func (r *article) FindArticle(ctx context.Context, id int64) (*model.Article, error) {
	a := model.Article{}
	if err := r.db.Get(&a, `
SELECT id, title, body FROM article WHERE id = ?
`, id); err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *article) CreateArticle(ctx context.Context, a *model.Article) (int64, error) {
	var tx Tx
	tx, ok := GetTx(ctx)
	if !ok {
		tx = r.db
	}
	stmt, err := tx.Prepare(`
INSERT INTO article (title, body) VALUES (?, ?)
`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(a.Title, a.Body)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *article) UpdateArticle(ctx context.Context, id int64, a *model.Article) (int64, error) {
	var tx Tx
	tx, ok := GetTx(ctx)
	if !ok {
		tx = r.db
	}
	stmt, err := tx.Prepare(`
UPDATE article SET title = ?, body = ? WHERE id = ?
`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(a.Title, a.Body, id)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *article) DestroyArticle(ctx context.Context, id int64) (int64, error) {
	var tx Tx
	tx, ok := GetTx(ctx)
	if !ok {
		tx = r.db
	}
	stmt, err := tx.Prepare(`
DELETE FROM article WHERE id = ?
`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
