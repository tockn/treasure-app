package service

import (
	"context"

	"firebase.google.com/go/auth"
)

type Auth interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
	GetUser(ctx context.Context, uid string) (*auth.UserRecord, error)
}
