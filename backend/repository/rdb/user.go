package rdb

import (
	"context"

	"github.com/voyagegroup/treasure-app/model"

	"github.com/jmoiron/sqlx"
	"github.com/voyagegroup/treasure-app/repository"
)

type user struct {
	txHandler
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) repository.User {
	return &user{
		txHandler: txHandler{db: db},
		db:        db,
	}
}

func (r *user) Get(ctx context.Context, id string) (*model.User, error) {
	var u model.User
	if err := r.db.Get(&u, `
SELECT id, firebase_uid, display_name, email, photo_url FROM user WHERE firebase_uid = ? LIMIT 1
	`, id); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *user) Sync(ctx context.Context, fu *model.FirebaseUser) (int64, error) {
	res, err := r.db.Exec(`
INSERT INTO user (firebase_uid, display_name, email, photo_url)
VALUES (?, ?, ?, ?)
ON DUPLICATE KEY
UPDATE display_name = ?, email = ?, photo_url = ?
`, fu.FirebaseUID, fu.DisplayName, fu.Email, fu.PhotoURL, fu.DisplayName, fu.Email, fu.PhotoURL)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
